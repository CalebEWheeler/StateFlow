package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CalebEWheeler/StateFlow/configs"
	"github.com/CalebEWheeler/StateFlow/servers"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/CalebEWheeler/StateFlow/workers"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Config configs.Config
	// Context Context
}

func New() App {
	return App{
		Config: configs.New(),
	}
}

// more logs for steps...
func (a App) Run() error {
	log.Infof("starting %s", a.Config.Name)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	// Initialize store
	log.Info("starting new store...")
	store, err := postgres.NewStore(ctx, a.Config.PostgresURL)
	if err != nil {
		return err
	}

	// Initialize worker engine
	log.Info("initializing worker engine...")
	worker := workers.NewWorker(store)
	go worker.Start(ctx)

	// Create Server, Router and register endpoints with handlers
	log.Info("initializing server...")
	server := servers.New(store)

	server.Start(ctx)

	// Wait until context is cancelled
	<-ctx.Done()

	log.Info("application shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(
		ctx,
		a.Config.GracefulExitTimeout,
	)
	defer shutdownCancel()

	server.Stop(shutdownCtx)

	return nil
}
