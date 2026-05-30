package handlers

import "context"

type CreateBillingRequest struct{}
type CreateBillingOutput struct{}

func NewCreateBillingHandler(ctx context.Context, input *CreateBillingRequest) (*CreateBillingOutput, error) {
	return &CreateBillingOutput{}, nil
}
