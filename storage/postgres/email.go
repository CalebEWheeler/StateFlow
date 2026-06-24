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

func (es EmailStore) SendConfirmation(ctx context.Context, job *Job) error {
	return nil
}
