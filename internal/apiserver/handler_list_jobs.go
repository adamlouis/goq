package apiserver

import (
	"context"
	"fmt"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) ListJobs(ctx context.Context, queryParams *goqmodel.ListJobsQueryParams) (*goqmodel.ListJobsResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
