package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CalebEWheeler/StateFlow/storage/postgres"
)

// loop:
// claim job
// execute job
// update job status
// repeat

type Worker struct {
	store *postgres.Store
}

// Polls for jobs.
// Claims jobs.
// Hands jobs off for processing.
// Repeats.
func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := w.store.ClaimNextPendingJob(ctx)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		err = w.ProcessJob(ctx, job)
		if err != nil {
			if failErr := w.store.Fail(
				ctx,
				job.ID,
				err,
			); failErr != nil {
				log.Printf(
					"failed to mark job failed: %v",
					failErr,
				)
			}
			continue
		}

		err = w.store.Complete(ctx, job.ID)
		if err != nil {
			if failErr := w.store.Fail(
				ctx,
				job.ID,
				fmt.Errorf("failed to complete job: %w", err),
			); failErr != nil {
				log.Printf(
					"failed to mark job failed: %v",
					failErr,
				)
			}
		}
	}
}

func (w *Worker) ProcessJob(ctx context.Context, job *postgres.Job) error {
	switch job.Step {
	case "create_order":
		return w.CreateOrder(ctx, job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Step)
	}
}

func (w *Worker) CreateOrder(ctx context.Context, job *postgres.Job) error {
	return nil
}

// Implement FIFO, filter by job = "pending", sort by "created_at" ascending
// Also prevent race condition using FOR UPDATE SKIP LOCKED
// SELECT *
// FROM jobs
// WHERE status = 'pending'
// ORDER BY created_at ASC
// FOR UPDATE SKIP LOCKED
// LIMIT 1;
// Consider adding "available_at" field to table for retries
// available_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

// SELECT *
// FROM jobs
// WHERE status = 'pending'
// AND available_at <= NOW()
// ORDER BY available_at ASC
// LIMIT 1;
