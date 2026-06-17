package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoJobs = errors.New("no pending jobs")

type JobStore struct {
	pool *pgxpool.Pool
}

type Job struct {
	ID         uuid.UUID
	WorkflowID uuid.UUID
	// Can add these steps as a bonus: retry_failed_job, cancel_order, reconcile_order
	Step       string `oneOf:"create_order,reserve_inventory,charge_payment,create_shipment,send_email"`
	Status     string `oneOf:"pending,completed,failed"`
	RetryCount int
	LastError  string
	Payload    []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewJobStore(pool *pgxpool.Pool) *JobStore {
	return &JobStore{pool: pool}
}

func (js *JobStore) CreateJob(ctx context.Context, workflowID uuid.UUID, payload interface{}) error {
	js.pool.Exec(context.Background(), `INSERT INTO jobs (
		id, 
		workflow_id,
		step,
		status,
		payload,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`, uuid.New(), workflowID, "create_order", "pending", payload, time.Now(), time.Now())

	return nil
}

func (js *JobStore) GetPending() bool {
	return true
}

func (j *JobStore) ClaimNextPendingJob(ctx context.Context) (*Job, error) {
	// Start transaction block...
	tx, err := j.pool.Begin(ctx)
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
