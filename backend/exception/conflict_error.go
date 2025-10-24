// exception/conflict_error.go
package exception

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(msg string) error {
	return &ConflictError{Message: msg}
}
