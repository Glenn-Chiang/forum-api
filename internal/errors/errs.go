package errs

type Error struct {
	Code    uint
	Message string
}

const (
	ErrInvalid = iota
	ErrNotFound
	ErrUnauthorized
	ErrConflict
	ErrInternal
)

func New(code uint, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
