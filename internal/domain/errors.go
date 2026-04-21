package domain

import "fmt"

// DomainError represents an error in the domain layer
type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

// Error codes
const (
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeValidationFailed  = "VALIDATION_FAILED"
	ErrCodeDuplicate         = "DUPLICATE"
	ErrCodeInvalidInput      = "INVALID_INPUT"
)

// ErrNotFound returns a not found error
func ErrNotFound(resource string) error {
	return DomainError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
	}
}

// ErrValidationFailed returns a validation error
func ErrValidationFailed(message string) error {
	return DomainError{
		Code:    ErrCodeValidationFailed,
		Message: message,
	}
}

// ErrDuplicate returns a duplicate resource error
func ErrDuplicate(resource string) error {
	return DomainError{
		Code:    ErrCodeDuplicate,
		Message: fmt.Sprintf("%s already exists", resource),
	}
}

// ErrInvalidInput returns an invalid input error
func ErrInvalidInput(message string) error {
	return DomainError{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	if domainErr, ok := err.(DomainError); ok {
		return domainErr.Code == ErrCodeNotFound
	}
	return false
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if domainErr, ok := err.(DomainError); ok {
		return domainErr.Code == ErrCodeValidationFailed
	}
	return false
}
