package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CalebEWheeler/StateFlow/servers"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/CalebEWheeler/StateFlow/workers"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Initialize store
	// Can move URL to env var
	store, err := postgres.NewStore(ctx, "postgres://postgres:example@localhost:5432/stateflow")
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
			30*time.Second,
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
}
