package main

import (
	"context"
	"net/http"

	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/CalebEWheeler/StateFlow/workers"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize store
	// Can move URL to env var
	store, err := postgres.NewStore(context.Background(), "postgres://postgres:example@localhost:5432/stateflow")
	if err != nil {
		panic(err)
	}

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
		OrderHandler: handlers.NewOrderHandler(store),
	}

	huma.Register(api, operations.Order, hs.OrderHandler.Handle)
	huma.Register(api, operations.SendEmail, hs.SendEmailHandler.Handle)

	http.ListenAndServe("127.0.0.1:8080", router)
}
