package main

import (
	"context"
	"net/http"

	"github.com/CalebEWheeler/StateFlow/connections"
	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/CalebEWheeler/StateFlow/workers"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Establish connection to Postgres
	// Can move URL to env var
	conn, err := connections.New(context.Background(), "postgres://postgres:example@localhost:5432/stateflow")
	if err != nil {
		panic(err)
	}

	// TODO: move to storage/postgres package...maybe call as an initialize function that returns the store?
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS workflows
	(
		id UUID PRIMARY KEY, 
		status VARCHAR(20) NOT NULL, 
		current_step VARCHAR(50), 
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS jobs
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
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders
	(
		id UUID PRIMARY KEY, 
		address JSONB NOT NULL,
		currency VARCHAR(10) NOT NULL,
		customer_id VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		items JSONB NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users
	(
		id UUID PRIMARY KEY, 
		first_name VARCHAR(50) NOT NULL, 
		last_name VARCHAR(50) NOT NULL, 
		email VARCHAR(255) UNIQUE NOT NULL
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS inventory
	(
		id VARCHAR(255) PRIMARY KEY,
		sku VARCHAR(50) UNIQUE NOT NULL,
		quantity INT NOT NULL,
		msrp NUMERIC(10,2) NOT NULL,
		price NUMERIC(10,2) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	conn.Pool.Exec(context.Background(), `INSERT INTO inventory (
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
	`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS shipments (
		id UUID PRIMARY KEY,
		order_id UUID NOT NULL,
		tracking_number VARCHAR(100) NOT NULL,
		carrier VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)

	// Initialize store
	store := postgres.NewStore(conn.Pool)

	// Initialize workers...
	worker := workers.NewWorker(store)
	go worker.Start(context.Background())

	// Create Router and register endpoints with handlers...
	router := chi.NewMux()
	config := huma.DefaultConfig("My API", "1.0.0")
	config.RejectUnknownQueryParameters = true
	api := humachi.New(router, config)

	// Initialize handlers with database connection
	hs := handlers.Handlers{
		OrderHandler: handlers.NewOrderHandler(conn, store),
	}

	huma.Register(api, operations.Order, hs.OrderHandler.Handle)
	huma.Register(api, operations.CreateUser, hs.CreateUserHandler.Handle)
	huma.Register(api, operations.CreateBilling, hs.CreateBillingHandler.Handle)
	huma.Register(api, operations.SendEmail, hs.SendEmailHandler.Handle)

	http.ListenAndServe("127.0.0.1:8080", router)
}
