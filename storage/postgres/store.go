package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *Store) ClaimNextPendingJob(ctx context.Context) (*Job, error) {
	// Start transaction block...
	tx, err := s.Job.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Discards all changes made since the Begin statement and restores DB to previous state.
	defer tx.Rollback(ctx)

	var job Job

	// not sure why chatgpt has .Scan for var job...I'd think at this stage you'd simply query the jobs table for the latest job with a status of pending...
	err = tx.QueryRow(ctx, `
		SELECT
			id,
			workflow_id,
			step,
			status,
			retry_count,
			last_error,
			payload,
			created_at,
			updated_at
	`).Scan(
		&job.ID,
		&job.WorkflowID,
		&job.Step,
		&job.Status,
		&job.RetryCount,
		&job.LastError,
		&job.Payload,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoJobs
		}

		return nil, err
	}

	_, err = tx.Exec(ctx, `
	UPDATE jobs
	SET 
		status = 'running'
		updated_at = NOW()
	WHERE id = $1
	`, job.ID)

	if err != nil {
		return nil, err
	}

	job.Status = "running"

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &job, nil
}
