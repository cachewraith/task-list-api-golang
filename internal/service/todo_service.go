package service

import (
	"context"
	"fmt"

	"github.com/cachewraith/task-list-api-golang/internal/domain"
	"github.com/cachewraith/task-list-api-golang/internal/repository"
	"github.com/cachewraith/task-list-api-golang/pkg/logger"
	"github.com/google/uuid"
)

// TodoService defines the interface for todo business logic
type TodoService interface {
	GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error)
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Todo, error)
	Create(ctx context.Context, input domain.CreateTodoInput, userID uuid.UUID) (*domain.Todo, error)
	CreateMany(ctx context.Context, inputs []domain.CreateTodoInput, userID uuid.UUID) ([]domain.Todo, error)
	Update(ctx context.Context, id uuid.UUID, input domain.UpdateTodoInput, userID uuid.UUID) (*domain.Todo, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// todoService implements TodoService
type todoService struct {
	repo repository.TodoRepository
}

// ServiceOption is a functional option for configuring the service
type ServiceOption func(*todoService)

// NewTodoService creates a new todo service with the given repository
func NewTodoService(repo repository.TodoRepository, opts ...ServiceOption) TodoService {
	svc := &todoService{
		repo: repo,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// ensure todoService implements TodoService
var _ TodoService = (*todoService)(nil)

// GetAll returns all todos for a user
func (s *todoService) GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	logger.Info("TODO_SERVICE", "GetAll", "Fetching all todos for user: "+userID.String())

	todos, err := s.repo.GetAll(ctx, userID)
	if err != nil {
		logger.Error("TODO_SERVICE", "GetAll", "Failed to fetch todos: "+err.Error())
		return nil, err
	}

	logger.Info("TODO_SERVICE", "GetAll", fmt.Sprintf("Retrieved %d todos for user: %s", len(todos), userID.String()))
	return todos, nil
}

// GetByID returns a todo by its ID (with ownership verification)
func (s *todoService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Todo, error) {
	logger.Info("TODO_SERVICE", "GetByID", fmt.Sprintf("Fetching todo %s for user %s", id.String(), userID.String()))

	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		logger.Warn("TODO_SERVICE", "GetByID", fmt.Sprintf("Todo not found: %s, error: %v", id.String(), err))
		return nil, err
	}

	logger.Info("TODO_SERVICE", "GetByID", "Successfully retrieved todo: "+id.String())
	return todo, nil
}

// Create creates a new todo after validating the input
func (s *todoService) Create(ctx context.Context, input domain.CreateTodoInput, userID uuid.UUID) (*domain.Todo, error) {
	logger.Info("TODO_SERVICE", "Create", fmt.Sprintf("Creating todo for user %s with title: %s", userID.String(), input.Title))

	if err := input.Validate(); err != nil {
		logger.Warn("TODO_SERVICE", "Create", "Validation failed: "+err.Error())
		return nil, err
	}

	todo := domain.NewTodo(input, userID)
	if err := s.repo.Create(ctx, todo); err != nil {
		logger.Error("TODO_SERVICE", "Create", "Failed to create todo: "+err.Error())
		return nil, err
	}

	logger.Info("TODO_SERVICE", "Create", "Successfully created todo: "+todo.ID.String())
	return todo, nil
}

// CreateMany creates multiple todos in a batch
func (s *todoService) CreateMany(ctx context.Context, inputs []domain.CreateTodoInput, userID uuid.UUID) ([]domain.Todo, error) {
	logger.Info("TODO_SERVICE", "CreateMany", fmt.Sprintf("Creating %d todos for user: %s", len(inputs), userID.String()))

	if len(inputs) == 0 {
		logger.Warn("TODO_SERVICE", "CreateMany", "Empty input array")
		return nil, domain.ErrValidationFailed("at least one todo is required")
	}

	todos := make([]domain.Todo, 0, len(inputs))
	for i, input := range inputs {
		if err := input.Validate(); err != nil {
			logger.Warn("TODO_SERVICE", "CreateMany", fmt.Sprintf("Validation failed at index %d: %s", i, err.Error()))
			return nil, domain.ErrValidationFailed(fmt.Sprintf("todo at index %d: %s", i, err.Error()))
		}
		todos = append(todos, *domain.NewTodo(input, userID))
	}

	for i := range todos {
		if err := s.repo.Create(ctx, &todos[i]); err != nil {
			logger.Error("TODO_SERVICE", "CreateMany", fmt.Sprintf("Failed to create todo at index %d: %s", i, err.Error()))
			return nil, err
		}
		logger.Debug("TODO_SERVICE", "CreateMany", fmt.Sprintf("Created todo %d/%d: %s", i+1, len(todos), todos[i].ID.String()))
	}

	logger.Info("TODO_SERVICE", "CreateMany", fmt.Sprintf("Successfully created %d todos", len(todos)))
	return todos, nil
}

// Update updates an existing todo after validating the input
func (s *todoService) Update(ctx context.Context, id uuid.UUID, input domain.UpdateTodoInput, userID uuid.UUID) (*domain.Todo, error) {
	logger.Info("TODO_SERVICE", "Update", fmt.Sprintf("Updating todo %s for user %s", id.String(), userID.String()))

	if err := input.Validate(); err != nil {
		logger.Warn("TODO_SERVICE", "Update", "Validation failed: "+err.Error())
		return nil, err
	}

	todo, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		logger.Warn("TODO_SERVICE", "Update", fmt.Sprintf("Todo not found for update: %s", id.String()))
		return nil, err
	}

	todo.Apply(input)

	if err := s.repo.Update(ctx, todo); err != nil {
		logger.Error("TODO_SERVICE", "Update", fmt.Sprintf("Failed to update todo %s: %s", id.String(), err.Error()))
		return nil, err
	}

	logger.Info("TODO_SERVICE", "Update", fmt.Sprintf("Successfully updated todo: %s", id.String()))
	return todo, nil
}

// Delete deletes a todo by its ID (with ownership verification)
func (s *todoService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	logger.Info("TODO_SERVICE", "Delete", fmt.Sprintf("Deleting todo %s for user %s", id.String(), userID.String()))

	if err := s.repo.Delete(ctx, id, userID); err != nil {
		logger.Error("TODO_SERVICE", "Delete", fmt.Sprintf("Failed to delete todo %s: %s", id.String(), err.Error()))
		return err
	}

	logger.Info("TODO_SERVICE", "Delete", fmt.Sprintf("Successfully deleted todo: %s", id.String()))
	return nil
}
