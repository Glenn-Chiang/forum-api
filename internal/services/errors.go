package services

import "fmt"

type NotFoundError struct {
	Name string 
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", err.Name)
}

func NewNotFoundError(name string) *NotFoundError {
	return &NotFoundError{Name: name}
}

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

type AlreadyInUseError struct {
	Field string
}

func (err *AlreadyInUseError) Error() string {
	return fmt.Sprintf("%s already in use", err.Field)
}

func NewAlreadyInUseError(field string) *AlreadyInUseError {
	return &AlreadyInUseError{Field: field}
}

type UnauthorizedError struct {
	Message string
}

func (err *UnauthorizedError) Error() string {
	return err.Message
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}
