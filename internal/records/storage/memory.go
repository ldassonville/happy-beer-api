package storage

import (
	"context"
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

func (dao *MemoryStorage) Create(ctx context.Context, event *api.Record) (*api.Record, error) {

	event.Id = uuid.NewString()
	dao.registry.Store(event.Id, event)

	return event, nil
}

func (dao *MemoryStorage) Search(ctx context.Context) ([]*api.Record, error) {

	var events []*api.Record = nil
	dao.registry.Range(func(key any, value any) bool {
		event := value.(*api.Record)
		events = append(events, event)
		return true
	})

	if len(events) > 0 {

		sort.SliceStable(events[:], func(i, j int) bool {
			return events[i].Date.Before(events[j].Date)
		})
	}
	return events, nil
}
