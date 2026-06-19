package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CalebEWheeler/StateFlow/storage/postgres"
)

type Worker struct {
	store *postgres.Store
}

func NewWorker(store *postgres.Store) *Worker {
	return &Worker{store: store}
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

		err = w.ProcessJob(ctx, job)
		if err != nil {
			if failErr := w.store.Job.Fail(
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

		err = w.store.Job.Complete(ctx, job.ID)
		if err != nil {
			if failErr := w.store.Job.Fail(
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
		// Update inventory table
		// 1. GET items data by orderID from 'orders' table
		// 2. UPDATE inventory quantity for each item from 'orders' table
		// 3. step := "create_shipment"
		// 4. create new job...
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
		//

		// step := "send_confirmation"
		// if err = w.store.Job.CreateJob(ctx, job.WorkflowID, step, orderID); err != nil {
		// 	return err
		// }
		return nil
	case "send_confirmation":
		// Get email from 'orders' table
		// Build and send email
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

func (w *Worker) ReserveInventory(ctx context.Context, job *postgres.Job) error {
	if err := w.store.Inventory.ReserveInventory(ctx, job); err != nil {
		return err
	}
	return nil
}
