package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) DeleteJob(ctx context.Context, pathParams *goqmodel.DeleteJobPathParams) error {
	repo, commit, rollback, err := h.GetJobRepository()
	if err != nil {
		return err
	}
	defer rollback() //nolint

	err = repo.Delete(ctx, pathParams.JobID)
	if err != nil {
		return err
	}

	return commit()
}
