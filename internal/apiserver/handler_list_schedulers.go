package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) ListSchedulers(ctx context.Context, queryParams *goqmodel.ListSchedulersRequest) (*goqmodel.ListSchedulersResponse, error) {
	repo, _, rollback, err := h.GetSchedulerRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	return repo.List(ctx, queryParams)
}
