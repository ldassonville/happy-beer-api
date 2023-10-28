package client

import (
	"context"
	"github.com/ldassonville/happy-beer-api/pkg/api"
)

type Client interface {

	// GetDispenser provide a dispenser ref
	GetDispenser(ctx context.Context, ref string) (*api.Dispenser, error)
	SearchDispensers(ctx context.Context) ([]*api.Dispenser, error)
	CreateDispenser(ctx context.Context, component *api.DispenserEditable) (*api.Dispenser, error)
	UpdateDispenser(ctx context.Context, component *api.Dispenser) (*api.Dispenser, error)
	DeleteDispenser(ctx context.Context, name string) error

	// SearchRecords give the business records
	SearchRecords(ctx context.Context) ([]*api.Record, error)
	CreateRecord(ctx context.Context, record *api.Record) (*api.Record, error)
}
