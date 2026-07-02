package servers

import (
	"context"
	"net"
	"net/http"

	"github.com/CalebEWheeler/StateFlow/handlers"
	"github.com/CalebEWheeler/StateFlow/operations"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
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

func (s *Server) Start(ctx context.Context) {
	log.Infof("starting server (%s)", s.server.Addr)
	s.server.BaseContext = func(_ net.Listener) context.Context { return ctx }

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()
	log.Info("started all servers")
}

func (s *Server) Stop(ctx context.Context) {
	log.Info("stopping all servers")
	if err := s.server.Shutdown(ctx); err != nil {
		log.Warn("Error shutting down server: ", err)
	}
	log.Info("Done stopping all servers")
}
