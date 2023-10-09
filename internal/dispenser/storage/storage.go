package storage

import (
	"context"
	"github.com/ldassonville/beer-puller-api/pkg/model"
)

type Storage interface {
	Search(ctx context.Context) ([]*model.Dispenser, error)

	Create(ctx context.Context, product *model.Dispenser) (*model.Dispenser, error)

	GetByRef(ctx context.Context, ref string) (*model.Dispenser, error)

	DeleteByRef(ctx context.Context, ref string) error
}
