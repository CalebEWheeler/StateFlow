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
	conn.Pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id uuid PRIMARY KEY NOT NULL, first_name VARCHAR(50) NOT NULL, last_name VARCHAR(50) NOT NULL, email VARCHAR(255) UNIQUE NOT NULL);")

	router := chi.NewMux()
	config := huma.DefaultConfig("My API", "1.0.0")
	config.RejectUnknownQueryParameters = true
	api := humachi.New(router, config)

	hs := handlers.Handlers{
		CreateUserHandler: handlers.NewCreateUserHandler(conn),
	}

	huma.Register(api, operations.CreateUser, hs.CreateUserHandler.Handle)
	huma.Register(api, operations.CreateBilling, handlers.NewCreateBillingHandler)
	huma.Register(api, operations.SendEmail, handlers.NewSendEmailHandler)

	http.ListenAndServe("127.0.0.1:8080", router)
}
