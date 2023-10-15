package storage

import (
	"context"
	"github.com/ldassonville/beer-puller-api/pkg/api"
)

type Storage interface {
	Search(ctx context.Context) ([]*api.Dispenser, error)

	Create(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error)

	GetByRef(ctx context.Context, ref string) (*api.Dispenser, error)

	DeleteByRef(ctx context.Context, ref string) error

	Update(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error)
}
