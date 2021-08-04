package apiserver

import (
	"context"
	"time"

	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/job/jobsqlite3"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/internal/scheduler"
	"github.com/adamlouis/goq/internal/scheduler/schedulersqlite3"
	"github.com/jmoiron/sqlx"
)

func NewAPIHandler(ctx context.Context, jobDB, schedulerDB *sqlx.DB) APIHandler {
	updatedSchedulerChan := make(chan string)

	runScheduler(ctx, jobDB, schedulerDB, updatedSchedulerChan)

	return &hdl{
		jobDB:                jobDB,
		schedulerDB:          schedulerDB,
		updatedSchedulerChan: updatedSchedulerChan,
	}
}

type hdl struct {
	jobDB                *sqlx.DB
	schedulerDB          *sqlx.DB
	updatedSchedulerChan chan string
}

type CommitFn func() error
type RollbackFn func() error

func (h *hdl) GetJobRepository() (job.Repository, CommitFn, RollbackFn, error) {
	tx, err := h.jobDB.Beginx()
	if err != nil {
		return nil, nil, nil, err
	}
	return jobsqlite3.NewJobRepository(tx), tx.Commit, tx.Rollback, nil
}

func (h *hdl) GetSchedulerRepository() (scheduler.Repository, CommitFn, RollbackFn, error) {
	tx, err := h.schedulerDB.Beginx()
	if err != nil {
		return nil, nil, nil, err
	}
	return schedulersqlite3.NewRepo(tx), tx.Commit, tx.Rollback, nil
}

func (h *hdl) onUpdateScheduler(id string) {
	h.updatedSchedulerChan <- id
}

func runScheduler(ctx context.Context, jobDB, schedulerDB *sqlx.DB, updatedSchedulerChan chan string) {
	runnerErrChan := make(chan error)

	runner := scheduler.NewRunner(
		schedulersqlite3.NewRepo(schedulerDB),
		jobsqlite3.NewJobRepository(jobDB),
		runnerErrChan,
	)

	go func() {
		for err := range runnerErrChan {
			jsonlog.Log(
				"name", "RunnerError",
				"error", err.Error(),
				"timestamp", time.Now(),
			)
		}
	}()

	go func() {
		for updatedID := range updatedSchedulerChan {
			jsonlog.Log(
				"name", "RunnerUpdate",
				"scheduler_id", updatedID,
				"timestamp", time.Now(),
			)
			runner.Update(ctx, updatedID)
		}
	}()

	go func() {
		for {
			err := runner.Run(ctx)
			runnerErrChan <- err
			time.Sleep(5 * time.Second)
		}
	}()
}
