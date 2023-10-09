package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ldassonville/beer-puller-api/pkg/model"
	"sync"
)

type MemoryStorage struct {
	registry sync.Map
}

func (dao *MemoryStorage) GetByRef(ctx context.Context, ref string) (*model.Dispenser, error) {

	if dispenser, ok := dao.registry.Load(ref); ok {
		return dispenser.(*model.Dispenser), nil
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

func (dao *MemoryStorage) Create(ctx context.Context, dispenser *model.Dispenser) (*model.Dispenser, error) {

	dispenser.Ref = uuid.NewString()

	dao.registry.Store(dispenser.Ref, dispenser)

	return dispenser, nil
}

func (dao *MemoryStorage) initDispenser() {

	var dispensers = []*model.Dispenser{
		{
			Ref:  uuid.NewString(),
			Size: "large",
			Beer: "Brewdog - Hazy Janes",
		},
		{
			Ref:  uuid.NewString(),
			Size: "large",
			Beer: "Brewdog - Hazy Janes",
		},
	}

	for _, dispenser := range dispensers {
		dao.registry.Store(dispenser.Ref, dispenser)
	}
}

func NewMemoryDao() Storage {

	dao := &MemoryStorage{}
	dao.initDispenser()

	return dao
}

func (dao *MemoryStorage) Search(ctx context.Context) ([]*model.Dispenser, error) {

	var products []*model.Dispenser = nil
	dao.registry.Range(func(key any, value any) bool {

		product := value.(*model.Dispenser)
		products = append(products, product)
		return true
	})

	return products, nil
}
