package servers

import (
	"context"
	"net/http"

	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
}

func New(store *postgres.Store) *Server {
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

	// migrate address and any other fields to an app config...
	return &Server{
		server: &http.Server{
			Addr:    "127.0.0.1:8080",
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
