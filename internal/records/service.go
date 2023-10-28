package records

import (
	"context"
	"github.com/ldassonville/happy-beer-api/internal/records/storage"
	"github.com/ldassonville/happy-beer-api/pkg/api"
	"github.com/ldassonville/happy-beer-api/pkg/core/event"
	"github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	dispatcher event.Dispatcher
	dao        storage.Storage
}

func NewService(dao storage.Storage) *Service {
	return &Service{
		dao:        dao,
		dispatcher: event.Dispatcher{},
	}
}

// Create a record
func (s *Service) Create(ctx context.Context, event *api.Record) (*api.Record, error) {

	event.Date = time.Now()
	record, err := s.dao.Create(ctx, event)
	if err != nil {
		logrus.WithError(err).Error("fail to save record")
	}
	s.dispatcher.Dispatch(ctx, "CREATE", event)

	return record, err
}

func (s *Service) Search(ctx context.Context) ([]*api.Record, error) {

	events, err := s.dao.Search(ctx)

	if err != nil {
		logrus.WithError(err).Error("fail to fetch records")
	}

	return events, err

}

func (s *Service) Subscribe() *event.Subscription {
	return s.dispatcher.Subscribe()
}
