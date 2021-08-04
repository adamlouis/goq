package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) PostScheduler(ctx context.Context, body *goqmodel.Scheduler) (*goqmodel.Scheduler, error) {
	repo, commit, rollback, err := h.GetSchedulerRepository()
	if err != nil {
		return nil, err
	}
	defer rollback()

	out, err := repo.Put(ctx, body)
	if err != nil {
		return nil, err
	}

	if err := commit(); err != nil {
		return nil, err
	}

	h.onUpdateScheduler(body.ID)

	return out, nil
}
