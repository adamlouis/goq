package apiserver

import (
	"context"

	"github.com/adamlouis/goq/internal/domain/job"
	"github.com/adamlouis/goq/pkg/goqmodel"
)

func (h *hdl) ClaimSomeJob(ctx context.Context, body *goqmodel.ClaimSomeJobRequest) (*goqmodel.Job, error) {
	repo, commit, rollback, err := h.GetJobRepository()
	if err != nil {
		return nil, err
	}
	defer rollback() //nolint

	out, err := repo.Claim(ctx, &job.ClaimOptions{
		Names: body.Names,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}
	if err = commit(); err != nil {
		return nil, err
	}
	return out, nil
}
