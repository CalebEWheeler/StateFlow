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
	connections.New(context.Background(), "postgres://postgres:example@localhost:5432/stateflow")

	router := chi.NewMux()
	config := huma.DefaultConfig("My API", "1.0.0")
	config.RejectUnknownQueryParameters = true
	api := humachi.New(router, config)

	huma.Register(api, operations.CreateUser, handlers.NewCreateUserHandler)
	huma.Register(api, operations.CreateBilling, handlers.NewCreateBillingHandler)
	huma.Register(api, operations.SendEmail, handlers.NewSendEmailHandler)

	http.ListenAndServe("127.0.0.1:8080", router)
}
