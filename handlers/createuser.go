package handlers

import (
	"context"
	"fmt"

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
		ID     string `json:"id"`
		Email  string `json:"email"`
		Status string `json:"status"`
	}
}

func NewCreateUserHandler(ctx context.Context, input *CreateUserRequest) (*CreateUserOutput, error) {
	usr := User{
		ID:        uuid.New().String(),
		FirstName: input.Body.FirstName,
		LastName:  input.Body.LastName,
		Email:     input.Body.Email,
	}
	fmt.Printf("sending user data to postgres... \n%+v", usr)
	resp := &CreateUserOutput{}
	resp.Body.ID = usr.ID
	resp.Body.Email = usr.Email
	resp.Body.Status = "created"
	return resp, nil
}
