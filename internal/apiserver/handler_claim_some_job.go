package apiserver

import (
	"context"
	"fmt"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) ClaimSomeJob(ctx context.Context, body *goqmodel.ClaimSomeJobRequest) (*goqmodel.Job, error) {
	return nil, fmt.Errorf("unimplemented")
}
