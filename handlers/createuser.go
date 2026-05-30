package handlers

import (
	"context"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Body struct {
		FirstName string `json:"firstName" example:"John" doc:"First Name"`
		LastName  string `json:"lastName" example:"Smith" doc:"Last Name"`
		Email     string `json:"email" example:"johnsmith@gmail.com" doc:"Email"`
		Password  string `json:"password" example:"" doc:"Password"`
	}
}

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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
		Password:  input.Body.Password,
	}
	resp.Body.Message = msg
	return resp, nil
}
