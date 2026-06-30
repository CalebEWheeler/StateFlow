package handlers

import (
	"context"
)

type SendEmailHandler struct{}

type SendEmailRequest struct{}
type SendEmailOutput struct{}

func NewSendEmailHandler() *SendEmailHandler {
	return &SendEmailHandler{}
}

func (h *SendEmailHandler) Handle(ctx context.Context, input *SendEmailRequest) (*SendEmailOutput, error) {
	return &SendEmailOutput{}, nil
}
