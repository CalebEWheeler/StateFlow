package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderStore struct {
	pool *pgxpool.Pool
}

func NewOrderStore(pool *pgxpool.Pool) *OrderStore {
	return &OrderStore{pool: pool}
}

func (o *OrderStore) CreateOrder(ctx context.Context, job *Job) error {
	var req shared.OrderRequestBody

	if err := json.Unmarshal(job.Payload, &req); err != nil {
		return err
	}

	_, err := o.pool.Exec(ctx, `
		INSERT INTO orders (
			id,
			address,
			currency,
			customer_id,
			email,
			items,
			created_at,
			updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`,
		job.OrderID,
		req.Address,
		req.Currency,
		req.CustomerID,
		req.Email,
		req.Items,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}
