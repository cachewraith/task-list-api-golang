package repository

import (
	"context"

	"github.com/cachewraith/task-list-api-golang/internal/domain"
	"github.com/google/uuid"
)

// TodoRepository defines the interface for todo data access
type TodoRepository interface {
	// GetAll returns all todos for a user
	GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error)

	// GetByID returns a todo by its ID and user ID (for ownership verification)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Todo, error)

	// Create creates a new todo
	Create(ctx context.Context, todo *domain.Todo) error

	// Update updates an existing todo
	Update(ctx context.Context, todo *domain.Todo) error

	// Delete deletes a todo by its ID and user ID (for ownership verification)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
