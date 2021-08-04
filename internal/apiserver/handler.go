package apiserver

import (
	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/internal/job/jobsqlite3"
	"github.com/jmoiron/sqlx"
)

func NewAPIHandler(db *sqlx.DB) APIHandler {
	return &hdl{
		db: db,
	}
}

type hdl struct {
	db *sqlx.DB
}

type CommitFn func() error
type RollbackFn func() error

func (h *hdl) GetJobRepository() (job.Repository, CommitFn, RollbackFn, error) {
	tx, err := h.db.Beginx()
	if err != nil {
		return nil, nil, nil, err
	}
	return jobsqlite3.NewJobRepository(tx), tx.Commit, tx.Rollback, nil
}
