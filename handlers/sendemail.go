package handlers

import (
	"context"

	"github.com/CalebEWheeler/StateFlow/connections"
)

type SendEmailHandler struct {
	db *connections.DB
}

type SendEmailRequest struct{}
type SendEmailOutput struct{}

func NewSendEmailHandler(conn *connections.DB) *SendEmailHandler {
	return &SendEmailHandler{
		db: conn,
	}
}

func (h *SendEmailHandler) Handle(ctx context.Context, input *SendEmailRequest) (*SendEmailOutput, error) {
	return &SendEmailOutput{}, nil
}
