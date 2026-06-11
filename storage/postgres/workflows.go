package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowStore struct {
	pool *pgxpool.Pool
}

type Workflow struct {
	ID          string
	Status      string `oneOf:"running,completed,failed"`
	CurrentStep string `oneOf:"create_user,create_billing,send_email"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewWorkflowStore(pool *pgxpool.Pool) *WorkflowStore {
	return &WorkflowStore{pool: pool}
}

func (ws *WorkflowStore) CreateWorkflow(ctx context.Context, workflowID string) error {
	if _, err := ws.pool.Exec(ctx, `INSERT INTO workflows 
	(
		id,
		status,
		current_step,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5)`, workflowID, "running", "create_user", time.Now(), time.Now()); err != nil {
		return err
	}
	return nil
}

func (ws *WorkflowStore) Get() bool {
	return true
}

func (ws *WorkflowStore) UpdateCurrentStep() bool {
	return true
}

func (ws *WorkflowStore) UpdateStatus() bool {
	return true
}
