package postgres

import (
	"context"

	"github.com/google/uuid"
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

func (s *Store) Complete(ctx context.Context, id uuid.UUID) error {
	_, err := s.Job.pool.Exec(ctx, `
	UPDATE jobs
	SET 
		status = 'completed',
		updated_at = NOW()
	WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	return nil
}

// Update later on to schedule a retry after X attempts...
func (s *Store) Fail(ctx context.Context, id uuid.UUID, pe error) error {
	// on fail...update job status, and last_error columns
	_, err := s.Job.pool.Exec(ctx, `
	UPDATE jobs
	SET
		status = 'failed',
		last_error = $2,
		updated_at = NOW()
	WHERE id = $1
	`, id, pe.Error())
	if err != nil {
		return err
	}

	return nil
}
