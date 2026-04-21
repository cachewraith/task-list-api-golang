package repository

import (
	"context"

	"github.com/cachewraith/task-list-api-golang/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// GetByID returns a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail returns a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error
}
