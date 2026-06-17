package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Job      JobStore
	Order    OrderStore
	Workflow WorkflowStore
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Job:      *NewJobStore(pool),
		Order:    *NewOrderStore(pool),
		Workflow: *NewWorkflowStore(pool),
	}
}
