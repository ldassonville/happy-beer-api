package dispenser

import (
	"context"
	"github.com/ldassonville/beer-puller-api/internal/dispenser/storage"
	"github.com/ldassonville/beer-puller-api/pkg/model"
)

type Service struct {
	dao storage.Storage
}

func (s *Service) Search(ctx context.Context) ([]*model.Dispenser, error) {

	return s.dao.Search(ctx)

}

func (s *Service) GetByRef(ctx context.Context, ref string) (*model.Dispenser, error) {

	return s.dao.GetByRef(ctx, ref)
}

func (s *Service) Create(ctx context.Context, dispenser *model.Dispenser) (*model.Dispenser, error) {

	return s.dao.Create(ctx, dispenser)
}

func (s *Service) DeleteByRef(ctx context.Context, ref string) error {

	return s.dao.DeleteByRef(ctx, ref)
}
