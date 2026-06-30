package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Email     EmailStore
	Inventory InventoryStore
	Job       JobStore
	Order     OrderStore
	Shipment  ShipmentStore
	Workflow  WorkflowStore
}

func NewStore(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	// verify connection actually works
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err := createTables(pool); err != nil {
		return nil, err
	}

	if err := seedTables(pool); err != nil {
		return nil, err
	}

	return &Store{
		Email:     *NewEmailStore(pool),
		Inventory: *NewInventoryStore(pool),
		Job:       *NewJobStore(pool),
		Order:     *NewOrderStore(pool),
		Shipment:  *NewShipmentStore(pool),
		Workflow:  *NewWorkflowStore(pool),
	}, nil
}
