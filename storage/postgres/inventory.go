package postgres

import (
	"context"
	"errors"

	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInsufficientInventory = errors.New("insufficient inventory")
	ErrOrderItemsNotFound    = errors.New("order items not found")
	items                    []shared.Item
)

type InventoryStore struct {
	pool *pgxpool.Pool
}

func NewInventoryStore(pool *pgxpool.Pool) *InventoryStore {
	return &InventoryStore{pool: pool}
}

func (is InventoryStore) ReserveInventory(ctx context.Context, job *Job) error {
	tx, err := is.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `
		SELECT 
			items 
		FROM orders
		WHERE id = $1`, job.OrderID).Scan(
		&items,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrOrderItemsNotFound
		}

		return err
	}

	// Loop through items, update inventory table for each item...
	for _, item := range items {
		result, err := tx.Exec(ctx, `
			UPDATE inventory
			SET
				quantity = quantity - $1,
				updated_at = NOW()
			WHERE id = $2
			AND sku = $3
			AND quantity >= $1
		`,
			item.Quantity,
			item.ID,
			item.SKU,
		)

		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return ErrInsufficientInventory
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
