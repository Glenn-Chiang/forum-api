package services

import "fmt"

type ValidationError struct {
	Field string
	Message string
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.Field, err.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field: field,
		Message: message,
	}
}
