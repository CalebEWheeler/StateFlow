package handlers

import (
	"context"
	"fmt"

	"github.com/CalebEWheeler/StateFlow/connections"
	"github.com/google/uuid"
)

type CreateUserHandler struct {
	db *connections.DB
}

type CreateUserRequest struct {
	Body struct {
		FirstName string `json:"firstName" required:"true" example:"John" minLength:"2" maxLength:"50" doc:"First name"`
		LastName  string `json:"lastName" required:"true" example:"Smith" minLength:"2" maxLength:"50" doc:"Last name"`
		Email     string `json:"email" required:"true" example:"johnsmith@gmail.com" format:"email" doc:"Email address"`
	}
}

type OutputBody struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type CreateUserOutput struct {
	Body OutputBody
}

// NewCreateUserHandler initializes a new CreateUserHandler with injected dependencies (e.g., database connection)
func NewCreateUserHandler(conn *connections.DB) *CreateUserHandler {
	return &CreateUserHandler{
		db: conn,
	}
}

func (h *CreateUserHandler) Handle(ctx context.Context, input *CreateUserRequest) (*CreateUserOutput, error) {
	uuid := uuid.New().String()
	fmt.Printf("sending user data to postgres...\n")
	if _, err := h.db.Pool.Exec(ctx, "INSERT INTO users (id, first_name, last_name, email) VALUES($1, $2, $3, $4)", uuid, input.Body.FirstName, input.Body.LastName, input.Body.Email); err != nil {
		// Returning the error and letting Huma handle the response. In a production application, you might want to wrap this error in a custom error type or add more context.
		return nil, fmt.Errorf("insert user: %w", err)
	} else {
		fmt.Printf("user data stored in postgres...\n")
	}
	resp := &CreateUserOutput{}
	resp.Body.ID = uuid
	resp.Body.Email = input.Body.Email
	resp.Body.Status = "created"
	return resp, nil
}
