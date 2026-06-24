package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailStore struct {
	pool *pgxpool.Pool
}

func NewEmailStore(pool *pgxpool.Pool) *EmailStore {
	return &EmailStore{pool: pool}
}

// Need the following...
// orders table - email address, address, items, order_id
// shipments table - tracking number, status, carrier
func (es EmailStore) SendConfirmation(ctx context.Context, job *Job) error {
	return nil
}
