// Package storage provides interfaces and implementations for task persistence
package storage

import "errors"

// Common errors returned by storage operations
var (
	// ErrTaskNotFound is returned when a task cannot be found
	ErrTaskNotFound = errors.New("task not found")

	// ErrDuplicateID is returned when attempting to add a task with an existing ID
	ErrDuplicateID = errors.New("task ID already exists")

	// ErrInvalidPath is returned when an invalid database path is provided
	ErrInvalidPath = errors.New("invalid database path")

	// ErrInvalidID is returned when a task ID is invalid
	ErrInvalidID = errors.New("invalid task ID")
)

// NotFoundError is returned when a task cannot be found
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return "task not found: " + e.ID
}

// Is implements errors.Is interface to match against ErrTaskNotFound
func (e *NotFoundError) Is(target error) bool {
	return target == ErrTaskNotFound
}
