package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CalebEWheeler/StateFlow/storage/postgres"
)

type Worker struct {
	store      *postgres.Store
	jobTimeout int
	maxRetries int
}

func NewWorker(store *postgres.Store) *Worker {
	return &Worker{
		store:      store,
		jobTimeout: 30,
		maxRetries: 3,
	}
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := w.store.Job.ClaimNextPendingJob(ctx)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		err = w.process(ctx, job)
		if err != nil {
			if err := w.HandleExecutionError(ctx, job, err); err != nil {
				log.Print(err)
			}

			continue
		}

		if err := w.HandleExecutionSuccess(ctx, job); err != nil {
			log.Print(err)
		}
	}
}

func (w *Worker) HandleExecutionError(ctx context.Context, job *postgres.Job, pe error) error {
	job.RetryCount++

	if job.RetryCount >= w.maxRetries {
		job.Status = postgres.StatusFailed
	} else {
		job.Status = postgres.StatusPending
	}

	if failErr := w.store.Job.Fail(
		ctx,
		job,
		pe,
	); failErr != nil {
		return fmt.Errorf(
			"failed to persist failed state for job=%s: %w",
			job.ID,
			failErr,
		)
	}

	if job.Status == postgres.StatusFailed {
		if err := w.store.Workflow.FailWorkflow(
			ctx,
			job.WorkflowID,
		); err != nil {
			return fmt.Errorf(
				"failed to handle job failure: %w",
				err,
			)
		}

		log.Printf(
			"workflowID= %s status=%s job=%s step=%s retries=%d",
			job.WorkflowID,
			postgres.StatusFailed,
			job.ID,
			job.Step,
			job.RetryCount,
		)
	}

	return nil
}

func (w *Worker) HandleExecutionSuccess(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Job.Complete(ctx, job.ID); err != nil {
		return fmt.Errorf(
			"failed to persist completion for job=%s workflow=%s: %w",
			job.ID,
			job.WorkflowID,
			err,
		)
	}

	log.Printf(
		"job=%s workflow=%s step=%s status=%s",
		job.ID,
		job.WorkflowID,
		job.Step,
		postgres.StatusComplete,
	)

	if job.Step != "send_confirmation" {
		return nil
	}

	if err := w.store.Workflow.CompleteWorkflow(ctx, job.WorkflowID); err != nil {
		return err
	}

	log.Printf(
		"workflow=%s status=%s",
		job.WorkflowID,
		postgres.StatusComplete,
	)

	return nil
}

func (w *Worker) process(ctx context.Context, job *postgres.Job) error {
	jobCtx, cancel := context.WithTimeout(
		ctx,
		time.Duration(w.jobTimeout)*time.Second,
	)
	defer cancel()

	return w.ProcessJob(jobCtx, job)
}

func (w *Worker) ProcessJob(ctx context.Context, job *postgres.Job) error {
	switch job.Step {
	case "create_order":
		err := w.CreateOrder(ctx, job)
		if err != nil {
			return err
		}

		if err = w.store.Job.CreateJob(ctx, postgres.Job{
			OrderID:    job.OrderID,
			Step:       "reserve_inventory",
			WorkflowID: job.WorkflowID,
		}); err != nil {
			return err
		}
		return nil
	case "reserve_inventory":
		err := w.ReserveInventory(ctx, job)
		if err != nil {
			return err
		}

		if err = w.store.Job.CreateJob(ctx, postgres.Job{
			OrderID:    job.OrderID,
			Step:       "create_shipment",
			WorkflowID: job.WorkflowID,
		}); err != nil {
			return err
		}
		return nil
	case "create_shipment":
		err := w.CreateShipment(ctx, job)
		if err != nil {
			return err
		}

		if err = w.store.Job.CreateJob(ctx, postgres.Job{
			OrderID:    job.OrderID,
			ShipmentID: job.ShipmentID,
			Step:       "send_confirmation",
			WorkflowID: job.WorkflowID,
		}); err != nil {
			return err
		}
		return nil
	case "send_confirmation":
		err := w.SendConfirmation(ctx, job)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown job type: %s", job.Step)
	}
}

func (w *Worker) CreateOrder(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Order.CreateOrder(ctx, job); err != nil {
		return err
	}
	return nil
}

func (w *Worker) CreateShipment(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Shipment.CreateShipment(ctx, job); err != nil {
		return err
	}
	return nil
}

func (w *Worker) ReserveInventory(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Inventory.ReserveInventory(ctx, job); err != nil {
		return err
	}
	return nil
}

func (w *Worker) SendConfirmation(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Email.SendConfirmation(ctx, job); err != nil {
		return err
	}
	return nil
}
