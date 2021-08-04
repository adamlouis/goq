package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) GetScheduler(ctx context.Context, pathParams *goqmodel.GetSchedulerPathParams) (*goqmodel.Scheduler, error) {
	repo, _, rollback, err := h.GetSchedulerRepository()
	if err != nil {
		return nil, err
	}
	defer rollback()

	return repo.Get(ctx, pathParams.SchedulerID)
}
