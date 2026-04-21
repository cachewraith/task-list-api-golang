package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/todo-api/internal/domain"
)

// Response is the standard API response structure
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Meta contains metadata for list responses (pagination, etc.)
type Meta struct {
	Total  int `json:"total,omitempty"`
	Page   int `json:"page,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewSuccess creates a success response
func NewSuccess(data interface{}) Response {
	return Response{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

// NewSuccessWithMessage creates a success response with a message
func NewSuccessWithMessage(message string, data interface{}) Response {
	return Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

// NewSuccessWithMeta creates a success response with pagination metadata
func NewSuccessWithMeta(data interface{}, meta *Meta) Response {
	return Response{
		Success:   true,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now().UTC(),
	}
}

// NewError creates an error response
func NewError(code string, message string) Response {
	return Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message},
		Timestamp: time.Now().UTC(),
	}
}

// NewErrorWithData creates an error response with additional data
func NewErrorWithData(code string, message string, data interface{}) Response {
	return Response{
		Success:   false,
		Error:     &ErrorInfo{Code: code, Message: message},
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

// Write sends the response to the client
func (r Response) Write(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}

// --- Helper functions for common responses ---

// OK sends a 200 OK response
func OK(w http.ResponseWriter, data interface{}) {
	NewSuccess(data).Write(w, http.StatusOK)
}

// Created sends a 201 Created response
func Created(w http.ResponseWriter, data interface{}) {
	NewSuccess(data).Write(w, http.StatusCreated)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, code string, message string) {
	NewError(code, message).Write(w, http.StatusBadRequest)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	NewError(domain.ErrCodeNotFound, message).Write(w, http.StatusNotFound)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(w http.ResponseWriter) {
	NewError("INTERNAL_ERROR", "An internal server error occurred").Write(w, http.StatusInternalServerError)
}

// FromDomainError converts a domain error to an appropriate HTTP response
func FromDomainError(w http.ResponseWriter, err error) {
	switch {
	case domain.IsNotFound(err):
		NotFound(w, err.Error())
	case domain.IsValidationError(err):
		BadRequest(w, domain.ErrCodeValidationFailed, err.Error())
	default:
		InternalError(w)
	}
}

// ValidationError sends a validation error response
func ValidationError(w http.ResponseWriter, message string) {
	BadRequest(w, domain.ErrCodeValidationFailed, message)
}

// Pagination creates pagination metadata
func Pagination(total, page, limit int) *Meta {
	return &Meta{
		Total:  total,
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
