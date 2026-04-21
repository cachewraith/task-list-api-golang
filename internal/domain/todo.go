package domain

import (
	"time"

	"github.com/google/uuid"
)

// Todo represents the domain entity for a todo item
type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTodoInput represents the input for creating a new todo
type CreateTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// UpdateTodoInput represents the input for updating an existing todo
type UpdateTodoInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

// Validate performs validation on CreateTodoInput
func (c CreateTodoInput) Validate() error {
	if c.Title == "" {
		return ErrValidationFailed("title is required")
	}
	if len(c.Title) > 255 {
		return ErrValidationFailed("title must be less than 255 characters")
	}
	return nil
}

// Validate performs validation on UpdateTodoInput
func (u UpdateTodoInput) Validate() error {
	if u.Title != nil && *u.Title == "" {
		return ErrValidationFailed("title cannot be empty if provided")
	}
	if u.Title != nil && len(*u.Title) > 255 {
		return ErrValidationFailed("title must be less than 255 characters")
	}
	return nil
}

// Apply updates the todo fields based on UpdateTodoInput
func (t *Todo) Apply(input UpdateTodoInput) {
	if input.Title != nil {
		t.Title = *input.Title
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if input.Completed != nil {
		t.Completed = *input.Completed
	}
	t.UpdatedAt = time.Now().UTC()
}

// NewTodo creates a new Todo from CreateTodoInput with user ID
func NewTodo(input CreateTodoInput, userID uuid.UUID) *Todo {
	now := time.Now().UTC()
	return &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
