package service

import (
	"context"

	"github.com/enterprise/enterprise-3tier/backend/internal/domain"
	"github.com/google/uuid"
)

// UserStore abstracts persistence for users (implemented by *repository.UserRepository).
type UserStore interface {
	Create(ctx context.Context, u domain.User) error
	List(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, name, email *string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
