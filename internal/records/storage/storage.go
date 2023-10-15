package storage

import (
	"context"
	"github.com/ldassonville/beer-puller-api/pkg/api"
)

type Storage interface {
	Search(ctx context.Context) ([]*api.Record, error)

	Create(ctx context.Context, dispenser *api.Record) (*api.Record, error)
}
