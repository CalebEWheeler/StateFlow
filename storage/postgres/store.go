package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Inventory InventoryStore
	Job       JobStore
	Order     OrderStore
	Shipments ShipmentStore
	Workflow  WorkflowStore
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Inventory: *NewInventoryStore(pool),
		Job:       *NewJobStore(pool),
		Order:     *NewOrderStore(pool),
		Shipments: *NewShipmentStore(pool),
		Workflow:  *NewWorkflowStore(pool),
	}
}
