package handlers

import (
	"context"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Body struct {
		FirstName string `json:"firstName" required:"true" example:"John" minLength:"2" maxLength:"50" doc:"First name"`
		LastName  string `json:"lastName" required:"true" example:"Smith" minLength:"2" maxLength:"50" doc:"Last name"`
		Email     string `json:"email" required:"true" example:"johnsmith@gmail.com" format:"email" doc:"Email address"`
	}
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserOutput struct {
	Body struct {
		Message User
		Success bool
	}
}

func NewCreateUserHandler(ctx context.Context, input *CreateUserRequest) (*CreateUserOutput, error) {
	resp := &CreateUserOutput{}
	resp.Body.Success = true
	msg := User{
		ID:        uuid.New().String(),
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Email:     input.Body.Email,
	}
	resp.Body.Message = msg
	return resp, nil
}
