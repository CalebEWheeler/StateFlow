package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// const maxRetries = 3

var (
	ErrNoJobs = errors.New("no pending jobs")
	job       Job
)

type JobStore struct {
	pool *pgxpool.Pool
}

type Job struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	OrderID    uuid.UUID
	ShipmentID uuid.UUID
	// Can add these steps as a bonus: retry_failed_job, cancel_order, reconcile_order
	Step       string `oneOf:"create_order,reserve_inventory,charge_payment,create_shipment,send_email"`
	Status     string `oneOf:"pending,completed,failed"`
	RetryCount int
	LastError  *string
	Payload    []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewJobStore(pool *pgxpool.Pool) *JobStore {
	return &JobStore{pool: pool}
}

func (j *JobStore) Complete(ctx context.Context, id uuid.UUID) error {
	_, err := j.pool.Exec(ctx, `
	UPDATE jobs
	SET 
		status = $2,
		updated_at = NOW()
	WHERE id = $1
	`, id, StatusComplete)

	if err != nil {
		return err
	}

	return nil
}

func (js *JobStore) CreateJob(ctx context.Context, job Job) error {
	_, err := js.pool.Exec(context.Background(), `INSERT INTO jobs (
		id, 
		workflow_id,
		order_id,
		step,
		status,
		payload,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		uuid.New(),
		job.WorkflowID,
		job.OrderID,
		job.Step,
		"pending",
		job.Payload,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

func (j *JobStore) ClaimNextPendingJob(ctx context.Context) (*Job, error) {
	// Start transaction block...
	tx, err := j.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	// Discards all changes made since the Begin statement and restores DB to previous state.
	defer tx.Rollback(ctx)

	// TODO: After completion, look to see which columns are no longer needed here.
	err = tx.QueryRow(ctx, `
		SELECT
			id,
			workflow_id,
			order_id,
			shipment_id,
			step,
			status,
			retry_count,
			last_error,
			payload,
			created_at,
			updated_at
		FROM jobs
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`, StatusPending).Scan(
		&job.ID,
		&job.WorkflowID,
		&job.OrderID,
		&job.ShipmentID,
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
		status = $1,
		updated_at = NOW()
	WHERE id = $2
	`, StatusRunning, job.ID)

	if err != nil {
		return nil, err
	}

	job.Status = StatusRunning

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *JobStore) Fail(ctx context.Context, job *Job, pe error) error {
	_, err := j.pool.Exec(ctx, `
	UPDATE jobs
	SET
		status = $1,
		retry_count = $2,
		last_error = $3,
		updated_at = NOW()
	WHERE id = $4
	`,
		job.Status,
		job.RetryCount,
		pe.Error(),
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark job %s failed: %w", job.ID, err)
	}

	return nil
}
