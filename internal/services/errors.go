package services

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}
