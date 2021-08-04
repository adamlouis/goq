package scheduler

import (
	"context"

	"github.com/adamlouis/goq/pkg/goqmodel"
)

type Repository interface {
	Put(ctx context.Context, scheduler *goqmodel.Scheduler) (*goqmodel.Scheduler, error)
	Get(ctx context.Context, id string) (*goqmodel.Scheduler, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, body *goqmodel.ListSchedulersRequest) (*goqmodel.ListSchedulersResponse, error)
}

type Runner interface {
	Run(ctx context.Context) error
	Update(ctx context.Context, id string)
}
