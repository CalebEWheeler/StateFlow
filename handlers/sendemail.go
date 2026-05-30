package handlers

import "context"

type SendEmailRequest struct{}
type SendEmailOutput struct{}

func NewSendEmailHandler(ctx context.Context, input *SendEmailRequest) (*SendEmailOutput, error) {
	return &SendEmailOutput{}, nil
}
