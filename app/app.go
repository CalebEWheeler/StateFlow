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

func (a App) Run() error {
	log.Infof("starting %s", a.Config.Name)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Initialize store
	// Can move URL to env var
	store, err := postgres.NewStore(ctx, a.Config.PostgresURL)
	if err != nil {
		panic(err)
	}

	// Initialize workers...
	worker := workers.NewWorker(store)
	go worker.Start(ctx)

	// Create Router and register endpoints with handlers...
	server := servers.New(store)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			ctx,
			a.Config.GracefulExitTimeout,
		)
		defer cancel()

		if err := server.Stop(shutdownCtx); err != nil {
			log.Printf("failed to shutdown server: %v", err)
		}
	}()

	log.Println("starting server...")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	return nil
}
