package postgres

import (
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

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Email:     *NewEmailStore(pool),
		Inventory: *NewInventoryStore(pool),
		Job:       *NewJobStore(pool),
		Order:     *NewOrderStore(pool),
		Shipment:  *NewShipmentStore(pool),
		Workflow:  *NewWorkflowStore(pool),
	}
}
