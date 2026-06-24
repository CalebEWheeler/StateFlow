package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Shipment struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	TrackingNumber string
	Carrier        string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

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
		"label_created",
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create shipment: %w", err)
	}

	return nil
}

// read from shipments table - tracking number, status, carrier
// scan data to and return Shipment struct...
func GetByOrderID(ctx context.Context, id uuid.UUID) error {
	return nil
}
