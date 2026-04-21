package errors

import (
	"encoding/json"
	"net/http"

	"github.com/example/todo-api/internal/domain"
)

// HTTPError represents an HTTP error response
type HTTPError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

// ErrorResponse writes an error response to the http.ResponseWriter
func ErrorResponse(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(HTTPError{
		Status:  status,
		Code:    code,
		Message: message,
	})
}

// FromDomainError converts a domain error to an HTTP error
func FromDomainError(w http.ResponseWriter, err error) {
	switch {
	case domain.IsNotFound(err):
		ErrorResponse(w, http.StatusNotFound, domain.ErrCodeNotFound, err.Error())
	case domain.IsValidationError(err):
		ErrorResponse(w, http.StatusBadRequest, domain.ErrCodeValidationFailed, err.Error())
	default:
		ErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal server error occurred")
	}
}

// Common HTTP errors
var (
	BadRequest = func(message string) HTTPError {
		return HTTPError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: message}
	}
	NotFound = func(message string) HTTPError {
		return HTTPError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: message}
	}
	InternalServerError = func() HTTPError {
		return HTTPError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "An internal server error occurred"}
	}
)
