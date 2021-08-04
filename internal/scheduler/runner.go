package scheduler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/pkg/jsonlog"
	"github.com/adamlouis/goq/pkg/goqmodel"
	cron "github.com/robfig/cron/v3"
)

func NewRunner(schedulerRepo Repository, jobRepo job.Repository, errChan chan<- error) Runner {
	return &runner{
		schedulerRepo:         schedulerRepo,
		jobRepo:               jobRepo,
		errChan:               errChan,
		entryIDsBySchedulerID: map[string]cron.EntryID{},
	}
}

type runner struct {
	schedulerRepo         Repository
	jobRepo               job.Repository
	c                     *cron.Cron
	errChan               chan<- error
	entryIDsBySchedulerID map[string]cron.EntryID
	writeMutex            sync.Mutex
}

func (r *runner) Run(ctx context.Context) error {
	r.writeMutex.Lock()
	all, err := getAllSchedulers(ctx, r.schedulerRepo)
	if err != nil {
		return err
	}

	r.c = cron.New()
	for _, s := range all {
		r.schedule(ctx, s)
	}
	r.writeMutex.Unlock()

	r.c.Run()
	return nil
}

func (r *runner) Update(ctx context.Context, id string) {
	r.writeMutex.Lock()
	s, err := r.schedulerRepo.Get(ctx, id)
	if err != nil {
		r.errChan <- err
	}
	r.schedule(ctx, s)
	r.writeMutex.Unlock()
}

func getAllSchedulers(ctx context.Context, repo Repository) ([]*goqmodel.Scheduler, error) {
	s := []*goqmodel.Scheduler{}
	pageToken := ""
	for {
		result, err := repo.List(ctx, &goqmodel.ListSchedulersRequest{
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}
		s = append(s, result.Schedulers...)
		pageToken = result.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return s, nil
}

func (r *runner) schedule(ctx context.Context, s *goqmodel.Scheduler) {
	_, err := json.Marshal(s.Input)
	if err != nil {
		r.errChan <- err
	}
	if id, ok := r.entryIDsBySchedulerID[s.ID]; ok {
		r.c.Remove(id)
		delete(r.entryIDsBySchedulerID, s.ID)
	}
	entryID, err := r.c.AddFunc(s.Schedule, func() {
		_, err := r.jobRepo.Queue(ctx, &goqmodel.Job{
			Name:  s.JobName,
			Input: goqmodel.JSONObject(s.Input),
		})
		if err != nil {
			r.errChan <- err
		} else {
			jsonlog.Log(
				"name", "ScheduleJob",
				"schedule", s.Schedule,
				"job_name", s.JobName,
				"input", s.Input,
				"timestamp", time.Now(),
			)
		}
	})
	if err != nil {
		r.errChan <- err
	} else {
		jsonlog.Log(
			"name", "AddScheduler",
			"entry_id", entryID,
			"schedule", s.Schedule,
			"job_name", s.JobName,
			"input", s.Input,
			"timestamp", time.Now(),
		)
	}
	r.entryIDsBySchedulerID[s.ID] = entryID
}
