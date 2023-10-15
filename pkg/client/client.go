package client

import (
	"context"
	"github.com/ldassonville/beer-puller-api/pkg/api"
)

type Client interface {

	// Dispensers
	GetDispenser(ctx context.Context, ref string) (*api.Dispenser, error)
	SearchDispensers(ctx context.Context) ([]*api.Dispenser, error)
	CreateDispenser(ctx context.Context, component *api.DispenserEditable) (*api.Dispenser, error)
	UpdateDispenser(ctx context.Context, component *api.Dispenser) (*api.Dispenser, error)
	DeleteDispenser(ctx context.Context, name string) error

	// Records
	SearchRecords(ctx context.Context) ([]*api.Record, error)
	CreateRecord(ctx context.Context, record *api.Record) (*api.Record, error)
}
