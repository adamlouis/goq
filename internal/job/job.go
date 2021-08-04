package job

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

type Repository interface {
	// crud
	Get(ctx context.Context, id string) (*goqmodel.Job, error)
	List(ctx context.Context, args *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, error)
	Delete(ctx context.Context, id string) error
	// job queue semantics
	Queue(ctx context.Context, j *goqmodel.Job) (*goqmodel.Job, error)
	Claim(ctx context.Context, opts *ClaimOptions) (*goqmodel.Job, error)
	Release(ctx context.Context, id string) (*goqmodel.Job, error)
	Success(ctx context.Context, id string, output goqmodel.JSONObject) (*goqmodel.Job, error)
	Error(ctx context.Context, id string, output goqmodel.JSONObject) (*goqmodel.Job, error)
}

type Reporter interface {
	GetCountByStatus(ctx context.Context) (map[JobStatus]int64, error)
	GetCountByName(ctx context.Context) (map[string]int64, error)
}

type JobStatus string

const (
	JobStatusQueued  JobStatus = "QUEUED"
	JobStatusClaimed JobStatus = "CLAIMED"
	JobStatusSuccess JobStatus = "SUCCESS"
	JobStatusError   JobStatus = "ERROR"
)

type ClaimOptions struct {
	JobID string
	Names []string
}
