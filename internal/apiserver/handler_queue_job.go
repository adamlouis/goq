package apiserver

import (
	"context"
	"fmt"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) QueueJob(ctx context.Context, body *goqmodel.Job) (*goqmodel.Job, error) {
	return nil, fmt.Errorf("unimplemented")
}
