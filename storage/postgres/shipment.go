package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShipmentStore struct {
	pool *pgxpool.Pool
}

func NewShipmentStore(pool *pgxpool.Pool) *ShipmentStore {
	return &ShipmentStore{pool: pool}
}

func (ss ShipmentStore) CreateShipment(ctx context.Context, job Job) (uuid.UUID, error) {
	return uuid.Nil, nil
}
