package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusComplete string = "completed"
	StatusFailed   string = "failed"
	StatusPending  string = "pending"
	StatusRunning  string = "running"
)

type WorkflowStore struct {
	pool *pgxpool.Pool
}

func NewWorkflowStore(pool *pgxpool.Pool) *WorkflowStore {
	return &WorkflowStore{pool: pool}
}

type Workflow struct {
	ID          uuid.UUID
	Status      string `oneOf:"running,completed,failed"`
	CurrentStep string `oneOf:"create_user,create_billing,send_email"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ws *WorkflowStore) CreateWorkflow(ctx context.Context, workflowID uuid.UUID) error {
	if _, err := ws.pool.Exec(ctx, `INSERT INTO workflows 
	(
		id,
		status,
		current_step,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5)`, workflowID, StatusRunning, "create_user", time.Now(), time.Now()); err != nil {
		return err
	}
	return nil
}

func (ws *WorkflowStore) CompleteWorkflow(ctx context.Context, workflowID uuid.UUID) error {
	_, err := ws.pool.Exec(ctx, `
		UPDATE workflows
		SET
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`, StatusComplete, workflowID)

	if err != nil {
		return fmt.Errorf(
			"failed to mark workflow %s complete: %w",
			workflowID,
			err,
		)
	}
	return nil
}

func (ws *WorkflowStore) FailWorkflow(ctx context.Context, workflowID uuid.UUID) error {
	_, err := ws.pool.Exec(ctx, `
		UPDATE workflows
		SET 
			status = $1,
			updated_at = NOW()
		WHERE id = $2
	`, StatusFailed, workflowID)
	if err != nil {
		return fmt.Errorf("failed to mark workflow %s failed: %w", workflowID, err)
	}
	return nil
}
