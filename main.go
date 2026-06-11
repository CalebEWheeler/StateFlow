package main

import (
	"context"
	"net/http"

	"github.com/CalebEWheeler/StateFlow/connections"
	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Can move URL to env var
	conn, err := connections.New(context.Background(), "postgres://postgres:example@localhost:5432/stateflow")
	if err != nil {
		panic(err)
	}
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
		step VARCHAR(50) NOT NULL,
		status VARCHAR(20) NOT NULL,
		retry_count INT NOT NULL DEFAULT 0,
		last_error TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS orders
	(
		id uuid PRIMARY KEY NOT NULL, 
		customer_id VARCHAR(50) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		address TEXT NOT NULL,
		items TEXT NOT NULL,
		currency VARCHAR(10) NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	conn.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users
	(
		id uuid PRIMARY KEY NOT NULL, 
		first_name VARCHAR(50) NOT NULL, 
		last_name VARCHAR(50) NOT NULL, 
		email VARCHAR(255) UNIQUE NOT NULL
	);`)

	router := chi.NewMux()
	config := huma.DefaultConfig("My API", "1.0.0")
	config.RejectUnknownQueryParameters = true
	api := humachi.New(router, config)

	// Initialize handlers with database connection
	hs := handlers.Handlers{
		CreateUserHandler:    handlers.NewCreateUserHandler(conn),
		CreateBillingHandler: handlers.NewCreateBillingHandler(conn),
		OrderHandler:         handlers.NewOrderHandler(conn),
		SendEmailHandler:     handlers.NewSendEmailHandler(conn),
	}

	huma.Register(api, operations.Order, hs.OrderHandler.Handle)
	huma.Register(api, operations.CreateUser, hs.CreateUserHandler.Handle)
	huma.Register(api, operations.CreateBilling, hs.CreateBillingHandler.Handle)
	huma.Register(api, operations.SendEmail, hs.SendEmailHandler.Handle)

	http.ListenAndServe("127.0.0.1:8080", router)
}
