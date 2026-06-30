package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusLabelCreated string = "label_created"
)

type ShipmentStore struct {
	pool *pgxpool.Pool
}

func NewShipmentStore(pool *pgxpool.Pool) *ShipmentStore {
	return &ShipmentStore{pool: pool}
}

func (ss ShipmentStore) CreateShipment(ctx context.Context, job *Job) error {
	_, err := ss.pool.Exec(ctx, `INSERT INTO shipments (
		id,
		order_id,
		tracking_number,
		carrier,
		status,
		created_at,
		updated_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		uuid.New(),
		job.OrderID,
		fmt.Sprintf("SF-%s",
			strings.ToUpper(uuid.New().String()[:8])),
		"UPS",
		StatusLabelCreated,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}

	return nil
}
