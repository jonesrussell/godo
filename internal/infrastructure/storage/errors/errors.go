// Package errors provides error definitions for the storage package
package errors

import "errors"

// Common errors
var (
	// ErrTaskNotFound is returned when a task cannot be found
	ErrTaskNotFound = errors.New("task not found")
	// ErrDuplicateID is returned when trying to add a task with an ID that already exists
	ErrDuplicateID = errors.New("task ID already exists")
	// ErrStoreClosed is returned when attempting to use a closed store
	ErrStoreClosed = errors.New("store is closed")
	// ErrEmptyID is returned when an empty task ID is provided
	ErrEmptyID = errors.New("task ID cannot be empty")
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
