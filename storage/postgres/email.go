package postgres

import (
	"context"
	"time"

	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Email struct {
	Address        shared.Address
	Carrier        string
	CreatedAt      time.Time
	Currency       string
	CustomerID     string
	Email          string
	Items          []shared.Item
	OrderID        uuid.UUID
	ShipmentID     uuid.UUID
	Status         string
	TrackingNumber string
	UpdatedAt      time.Time
}

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
