package service

import (
	"context"
	"errors"
	"time"

	"github.com/enterprise/enterprise-3tier/backend/internal/domain"
	"github.com/enterprise/enterprise-3tier/backend/internal/repository"
	"github.com/google/uuid"
)

// UserService implements business rules for users.
type UserService struct {
	repo UserStore
}

func NewUserService(repo UserStore) *UserService {
	return &UserService{repo: repo}
}

// Create validates input and persists a new user.
func (s *UserService) Create(ctx context.Context, in domain.CreateUserInput) (domain.User, error) {
	u := domain.User{
		ID:        uuid.New(),
		Name:      in.Name,
		Email:     in.Email,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		if repository.IsUniqueViolation(err) {
			return domain.User{}, ErrEmailTaken
		}
		return domain.User{}, err
	}
	return u, nil
}

// List returns all users.
func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}

// Get returns one user by id.
func (s *UserService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if repository.IsNotFound(err) {
			return domain.User{}, ErrNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

// Update applies partial updates.
func (s *UserService) Update(ctx context.Context, id uuid.UUID, in domain.UpdateUserInput) (domain.User, error) {
	u, err := s.repo.Update(ctx, id, in.Name, in.Email)
	if err != nil {
		if repository.IsNotFound(err) {
			return domain.User{}, ErrNotFound
		}
		if repository.IsUniqueViolation(err) {
			return domain.User{}, ErrEmailTaken
		}
		return domain.User{}, err
	}
	return u, nil
}

// Delete removes a user.
func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if repository.IsNotFound(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// Sentinel errors for handlers.
var (
	ErrNotFound   = errors.New("user not found")
	ErrEmailTaken = errors.New("email already registered")
)
