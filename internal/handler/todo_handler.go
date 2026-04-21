package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/example/todo-api/internal/domain"
	"github.com/example/todo-api/internal/middleware"
	"github.com/example/todo-api/internal/service"
	"github.com/example/todo-api/pkg/logger"
	"github.com/example/todo-api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service service.TodoService
}

// HandlerOption is a functional option for configuring the handler
type HandlerOption func(*TodoHandler)

// NewTodoHandler creates a new todo handler with the given service
func NewTodoHandler(svc service.TodoService, opts ...HandlerOption) *TodoHandler {
	h := &TodoHandler{
		service: svc,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Routes returns the router with all todo routes registered
func (h *TodoHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetAll)
	r.Post("/", h.Create)
	r.Post("/bulk", h.CreateMany)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetByID)
		r.Put("/", h.Update)
		r.Delete("/", h.Delete)
	})

	return r
}

// GetAll handles GET /todos
func (h *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "GetAll", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "GetAll", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	todos, err := h.service.GetAll(ctx, userID)
	if err != nil {
		logger.Error("TODO_HANDLER", "GetAll", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "GetAll", "Successfully retrieved todos")
	meta := &response.Meta{Total: len(todos)}
	response.NewSuccessWithMeta(todos, meta).Write(w, http.StatusOK)
}

// GetByID handles GET /todos/{id}
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "GetByID", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "GetByID", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		logger.Warn("TODO_HANDLER", "GetByID", "Invalid ID format: "+chi.URLParam(r, "id"))
		response.BadRequest(w, "INVALID_ID", "Invalid todo ID format")
		return
	}
	logger.Debug("TODO_HANDLER", "GetByID", "Fetching todo: "+id.String())

	todo, err := h.service.GetByID(ctx, id, userID)
	if err != nil {
		logger.Error("TODO_HANDLER", "GetByID", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "GetByID", "Successfully retrieved todo: "+id.String())
	response.OK(w, todo)
}

// Create handles POST /todos
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "Create", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "Create", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	var input domain.CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("TODO_HANDLER", "Create", "Invalid JSON payload: "+err.Error())
		response.BadRequest(w, "INVALID_JSON", "Invalid JSON payload")
		return
	}
	defer r.Body.Close()
	logger.Debug("TODO_HANDLER", "Create", "Creating todo with title: "+input.Title)

	todo, err := h.service.Create(ctx, input, userID)
	if err != nil {
		logger.Error("TODO_HANDLER", "Create", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "Create", "Successfully created todo: "+todo.ID.String())
	response.Created(w, todo)
}

// CreateMany handles POST /todos/bulk for creating multiple todos
func (h *TodoHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "CreateMany", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "CreateMany", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	var inputs []domain.CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		logger.Warn("TODO_HANDLER", "CreateMany", "Invalid JSON payload: "+err.Error())
		response.BadRequest(w, "INVALID_JSON", "Invalid JSON payload - expected array of todos")
		return
	}
	defer r.Body.Close()
	logger.Debug("TODO_HANDLER", "CreateMany", fmt.Sprintf("Creating %d todos", len(inputs)))

	todos, err := h.service.CreateMany(ctx, inputs, userID)
	if err != nil {
		logger.Error("TODO_HANDLER", "CreateMany", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "CreateMany", fmt.Sprintf("Successfully created %d todos", len(todos)))
	meta := &response.Meta{Total: len(todos)}
	response.NewSuccessWithMeta(todos, meta).Write(w, http.StatusCreated)
}

// Update handles PUT /todos/{id}
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "Update", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "Update", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		logger.Warn("TODO_HANDLER", "Update", "Invalid ID format: "+chi.URLParam(r, "id"))
		response.BadRequest(w, "INVALID_ID", "Invalid todo ID format")
		return
	}
	logger.Debug("TODO_HANDLER", "Update", "Updating todo: "+id.String())

	var input domain.UpdateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("TODO_HANDLER", "Update", "Invalid JSON payload: "+err.Error())
		response.BadRequest(w, "INVALID_JSON", "Invalid JSON payload")
		return
	}
	defer r.Body.Close()

	todo, err := h.service.Update(ctx, id, input, userID)
	if err != nil {
		logger.Error("TODO_HANDLER", "Update", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "Update", "Successfully updated todo: "+id.String())
	response.OK(w, todo)
}

// Delete handles DELETE /todos/{id}
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Info("TODO_HANDLER", "Delete", "Request received")

	userID, ok := getUserIDFromContext(ctx)
	if !ok {
		logger.Warn("TODO_HANDLER", "Delete", "Unauthorized access attempt")
		response.NewError("UNAUTHORIZED", "User not authenticated").Write(w, http.StatusUnauthorized)
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		logger.Warn("TODO_HANDLER", "Delete", "Invalid ID format: "+chi.URLParam(r, "id"))
		response.BadRequest(w, "INVALID_ID", "Invalid todo ID format")
		return
	}
	logger.Debug("TODO_HANDLER", "Delete", "Deleting todo: "+id.String())

	if err := h.service.Delete(ctx, id, userID); err != nil {
		logger.Error("TODO_HANDLER", "Delete", "Service error: "+err.Error())
		response.FromDomainError(w, err)
		return
	}

	logger.Info("TODO_HANDLER", "Delete", "Successfully deleted todo: "+id.String())
	response.NewSuccessWithMessage("Todo deleted successfully", nil).Write(w, http.StatusOK)
}

// respondJSON sends a raw JSON response (kept for backward compatibility)
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// parseUUID parses a string into a UUID
func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userIDStr, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		return uuid.UUID{}, false
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, false
	}
	return userID, true
}
