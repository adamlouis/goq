package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) DeleteScheduler(ctx context.Context, pathParams *goqmodel.DeleteSchedulerPathParams) error {
	repo, commit, rollback, err := h.GetSchedulerRepository()
	if err != nil {
		return err
	}
	defer rollback()

	err = repo.Delete(ctx, pathParams.SchedulerID)
	if err != nil {
		return err
	}

	return commit()
}
