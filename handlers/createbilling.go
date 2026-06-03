package handlers

import (
	"context"

	"github.com/CalebEWheeler/StateFlow/connections"
)

type CreateBillingHandler struct {
	db *connections.DB
}

type CreateBillingRequest struct {
	Body struct {
		ID    string `json:"id" example:"365c2b56-5bd1-49e9-8f32-c66fc16881e7" format:"uuid" doc:"Id"`
		Email string `json:"email" example:"johnsmith@gmail.com" format:"email" doc:"Email address"`
		Plan  string `json:"plan" enum:"Basic,Intermediate,Advanced" example:"Basic" doc:"Plan type"`
	}
}
type CreateBillingOutput struct{}

func NewCreateBillingHandler(conn *connections.DB) *CreateBillingHandler {
	return &CreateBillingHandler{
		db: conn,
	}
}

func (h *CreateBillingHandler) Handle(ctx context.Context, input *CreateBillingRequest) (*CreateBillingOutput, error) {
	return &CreateBillingOutput{}, nil
}
