package apiserver

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) SearchJobs(ctx context.Context, body *goqmodel.SearchJobsRequest) (*goqmodel.SearchJobsResponse, error) {
	r, _, rollback, err := h.GetJobRepository()
	defer rollback() //nolint

	if err != nil {
		return nil, err
	}
	return r.Search(ctx, body)
}
