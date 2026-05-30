package main

import (
	"net/http"

	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Register(api, operations.CreateUser, handlers.NewCreateUserHandler)
	huma.Register(api, operations.CreateBilling, handlers.NewCreateBillingHandler)
	huma.Register(api, operations.SendEmail, handlers.NewSendEmailHandler)

	http.ListenAndServe("127.0.0.1:8080", router)
}
