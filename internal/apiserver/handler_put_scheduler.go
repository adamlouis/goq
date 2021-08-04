package apiserver

import (
	"context"
	"fmt"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) PutScheduler(ctx context.Context, pathParams *goqmodel.PutSchedulerPathParams, body *goqmodel.Scheduler) (*goqmodel.Scheduler, error) {
	if pathParams.SchedulerID != body.ID {
		return nil, fmt.Errorf("id in path does not match id in request body")
	}
	return h.PostScheduler(ctx, body)
}
