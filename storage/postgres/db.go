package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createTables(pool *pgxpool.Pool) error {
	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS workflows
	(
		id UUID PRIMARY KEY, 
		status VARCHAR(20) NOT NULL, 
		current_step VARCHAR(50), 
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`); err != nil {
		return fmt.Errorf("failed to create workflows table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS jobs
	(
		id UUID PRIMARY KEY,
		workflow_id UUID NOT NULL,
		order_id UUID,
		shipment_id UUID,
		step VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		retry_count INT NOT NULL DEFAULT 0,
		last_error TEXT,
		payload JSONB,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`); err != nil {
		return fmt.Errorf("failed to create jobs table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders
	(
		id UUID PRIMARY KEY, 
		address JSONB NOT NULL,
		currency VARCHAR(10) NOT NULL,
		customer_id VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		items JSONB NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`); err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS inventory
	(
		id VARCHAR(255) PRIMARY KEY,
		sku VARCHAR(50) UNIQUE NOT NULL,
		quantity INT NOT NULL,
		msrp NUMERIC(10,2) NOT NULL,
		price NUMERIC(10,2) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`); err != nil {
		return fmt.Errorf("failed to create inventory table: %w", err)
	}

	if _, err := pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS shipments (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL,
		tracking_number VARCHAR(100) NOT NULL,
		carrier VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`); err != nil {
		return fmt.Errorf("failed to create shipments table: %w", err)
	}

	return nil
}

func seedTables(pool *pgxpool.Pool) error {
	if _, err := pool.Exec(context.Background(), `INSERT INTO inventory (
		id,
		sku,
		quantity,
		msrp,
		price
	)
	VALUES
		(1234567890, 'ABC123', 100, 29.99, 24.99),
		(9876543210, 'DEF456', 50, 49.99, 39.99),
		(1234509876, 'GHI789', 25, 99.99, 89.99),
		(9876012345, 'JKL012', 10, 199.99, 179.99),
		(1256903478, 'MNO345', 500, 9.99, 7.99)
	ON CONFLICT (sku) DO NOTHING;
	`); err != nil {
		return fmt.Errorf("failed to seed inventory table: %w", err)
	}

	return nil
}
