package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) GetJob(ctx context.Context, pathParams *goqmodel.GetJobPathParams) (*goqmodel.Job, error) {
	repo, _, rollback, err := h.GetJobRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	return repo.Get(ctx, pathParams.JobID)
}
