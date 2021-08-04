package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) SetJobSuccess(ctx context.Context, pathParams *goqmodel.SetJobSuccessPathParams) (*goqmodel.Job, error) {
	repo, commit, rollback, err := h.GetJobRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	out, err := repo.Success(ctx, pathParams.JobID)
	if err != nil {
		return nil, err
	}

	if err = commit(); err != nil {
		return nil, err
	}

	return out, nil
}
