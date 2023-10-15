package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ldassonville/beer-puller-api/pkg/api"
	"sort"
	"sync"
)

type MemoryStorage struct {
	registry sync.Map
}

func NewMemoryDao() Storage {

	dao := &MemoryStorage{}
	return dao
}

func (dao *MemoryStorage) Update(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error) {

	dao.registry.Store(dispenser.Ref, dispenser)

	return dispenser, nil
}

func (dao *MemoryStorage) GetByRef(ctx context.Context, ref string) (*api.Dispenser, error) {

	if dispenser, ok := dao.registry.Load(ref); ok {
		return dispenser.(*api.Dispenser), nil
	}
	return nil, errors.New("not_found")
}

func (dao *MemoryStorage) DeleteByRef(ctx context.Context, ref string) error {

	if _, ok := dao.registry.Load(ref); !ok {
		return errors.New("not_found")
	}
	dao.registry.Delete(ref)
	return nil
}

func (dao *MemoryStorage) Create(ctx context.Context, dispenser *api.Dispenser) (*api.Dispenser, error) {

	dispenser.Ref = uuid.NewString()[:8]

	dao.registry.Store(dispenser.Ref, dispenser)

	return dispenser, nil
}

func (dao *MemoryStorage) Search(ctx context.Context) ([]*api.Dispenser, error) {

	var dispensers []*api.Dispenser = nil
	dao.registry.Range(func(key any, value any) bool {
		dispenser := value.(*api.Dispenser)
		dispensers = append(dispensers, dispenser)
		return true
	})

	if len(dispensers) > 0 {

		sort.SliceStable(dispensers[:], func(i, j int) bool {

			iCreatedAt := dispensers[i].Metadata.CreatedAt
			jCreatedAt := dispensers[j].Metadata.CreatedAt

			return iCreatedAt.Before(jCreatedAt)
		})
	}

	return dispensers, nil
}
