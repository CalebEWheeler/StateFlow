package configs

import (
	"io"
	"log"
	"os"
	"syscall"
	"time"
)

type Config struct {
	ExitSignals         []os.Signal
	LogWriter           io.Writer
	Name                string
	GracefulExitTimeout time.Duration
	Port                int
	PostgresURL         string
}

// process env configs here later...
// create servers config when metrics come into play...
func New() Config {
	config := Config{
		ExitSignals: []os.Signal{
			os.Interrupt,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		},
		LogWriter:           os.Stderr,
		Name:                "stateflow",
		GracefulExitTimeout: 10,
		Port:                8080,
		PostgresURL:         "postgres://postgres:example@localhost:5432/stateflow",
	}

	configureGlobalLogger(config)

	return config
}

func configureGlobalLogger(config Config) {
	log.SetOutput(config.LogWriter)
}
