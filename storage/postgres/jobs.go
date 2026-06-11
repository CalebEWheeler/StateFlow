package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobStore struct {
	pool *pgxpool.Pool
}

type Job struct {
	ID         string
	WorkflowID string
	CreatedAt  time.Time
	Status     string `oneOf:"pending,completed,failed"`
	// Can add these steps as a bonus: retry_failed_job, cancel_order, reconcile_order
	Step    string `oneOf:"create_order,reserve_inventory,charge_payment,create_shipment,send_email"`
	Payload []byte
}

func NewJobStore(pool *pgxpool.Pool) *JobStore {
	return &JobStore{pool: pool}
}

func (js *JobStore) CreateJob(ctx context.Context, workflowID string, payload interface{}) error {
	js.pool.Exec(context.Background(), `INSERT INTO jobs (
		id, 
		workflow_id,
		step,
		status,
		payload,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`, uuid.New().String(), workflowID, "create_order", "pending", payload, time.Now(), time.Now())

	return nil
}

func (js *JobStore) GetPending() bool {
	return true
}

func (js *JobStore) Complete() bool {
	return true
}

func (js *JobStore) Fail() bool {
	return true
}
