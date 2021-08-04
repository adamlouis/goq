package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) SetJobError(ctx context.Context, pathParams *goqmodel.SetJobErrorPathParams, j *goqmodel.Job) (*goqmodel.Job, error) {
	repo, commit, rollback, err := h.GetJobRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	out, err := repo.Error(ctx, pathParams.JobID, j.Output)
	if err != nil {
		return nil, err
	}

	if err = commit(); err != nil {
		return nil, err
	}

	return out, nil
}
