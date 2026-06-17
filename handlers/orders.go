package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CalebEWheeler/StateFlow/connections"
	"github.com/CalebEWheeler/StateFlow/shared"
	"github.com/CalebEWheeler/StateFlow/storage/postgres"
	"github.com/google/uuid"
)

type OrderHandler struct {
	db    *connections.DB
	store *postgres.Store
}

type OrderContext struct {
}

type Order struct {
	ID         uuid.UUID
	CustomerID string
	Status     string `oneOf:"pending,processing,completed,failed"`

	Subtotal     float64
	Tax          float64
	ShippingCost float64
	Total        float64
	Currency     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderRequest struct {
	Body shared.OrderRequestBody `json:"body"`
}

type OrderResponse struct{}

func NewOrderHandler(conn *connections.DB, store *postgres.Store) *OrderHandler {
	return &OrderHandler{
		db:    conn,
		store: store,
	}
}

func (h *OrderHandler) Handle(ctx context.Context, input *OrderRequest) (*OrderResponse, error) {
	workflowID := uuid.New()
	if err := h.store.Workflow.CreateWorkflow(ctx, workflowID); err != nil {
		return &OrderResponse{}, err
	}

	// If I can migrate payload data from using json.Marshal to []bytes instead, that would improve performance...
	payload, err := json.Marshal(input.Body)
	if err != nil {
		return &OrderResponse{}, err
	}
	h.store.Job.CreateJob(ctx, workflowID, payload)

	return &OrderResponse{}, nil
}
