package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Workflow WorkflowStore
	Job      JobStore
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{
		Workflow: *NewWorkflowStore(pool),
		Job:      *NewJobStore(pool),
	}
}
