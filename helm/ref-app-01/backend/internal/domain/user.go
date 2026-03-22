package domain

import (
	"time"

	"github.com/google/uuid"
)

// User is the core entity for the users table.
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserInput is used when creating a new user.
type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserInput is used for partial updates.
type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
