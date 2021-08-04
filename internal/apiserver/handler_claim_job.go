package apiserver

import (
	"context"

	"github.com/adamlouis/goq/internal/job"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) ClaimJob(ctx context.Context, pathParams *goqmodel.ClaimJobPathParams) (*goqmodel.Job, error) {
	repo, commit, rollback, err := h.GetJobRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	out, err := repo.Claim(ctx, &job.ClaimOptions{
		JobID: pathParams.JobID,
	})
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, nil
	}

	if err = commit(); err != nil {
		return nil, err
	}

	return out, nil

}
