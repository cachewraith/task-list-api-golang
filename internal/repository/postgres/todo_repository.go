package postgres

import (
	"context"
	"errors"

	"github.com/cachewraith/task-list-api-golang/internal/domain"
	"github.com/cachewraith/task-list-api-golang/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// todoRepository implements repository.TodoRepository using PostgreSQL
type todoRepository struct {
	pool *pgxpool.Pool
}

// NewTodoRepository creates a new PostgreSQL todo repository
func NewTodoRepository(pool *pgxpool.Pool) repository.TodoRepository {
	return &todoRepository{pool: pool}
}

// ensure todoRepository implements repository.TodoRepository
var _ repository.TodoRepository = (*todoRepository)(nil)

// GetAll returns all todos for a specific user
func (r *todoRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var todo domain.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// GetByID returns a todo by its ID and user ID (for ownership verification)
func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1 AND user_id = $2
	`

	var todo domain.Todo
	err := r.pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound("todo")
		}
		return nil, err
	}

	return &todo, nil
}

// Create inserts a new todo into the database
func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `
		INSERT INTO todos (id, user_id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		todo.ID,
		todo.UserID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CreatedAt,
		todo.UpdatedAt,
	)

	return err
}

// Update updates an existing todo in the database (with user_id verification)
func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
		UPDATE todos
		SET title = $2, description = $3, completed = $4, updated_at = $5
		WHERE id = $1 AND user_id = $6
	`

	result, err := r.pool.Exec(ctx, query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.UpdatedAt,
		todo.UserID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound("todo")
	}

	return nil
}

// Delete removes a todo from the database (with user_id verification)
func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`

	result, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound("todo")
	}

	return nil
}
