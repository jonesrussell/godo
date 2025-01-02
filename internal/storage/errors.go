// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"errors"
)

// Common errors returned by storage operations
var (
	// ErrTaskNotFound is returned when a task cannot be found
	ErrTaskNotFound = errors.New("task not found")

	// ErrStoreClosed is returned when attempting to use a closed store
	ErrStoreClosed = errors.New("store is closed")

	// ErrEmptyID is returned when an empty task ID is provided
	ErrEmptyID = errors.New("task ID cannot be empty")

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

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// ConnectionError represents a database connection error
type ConnectionError struct {
	Operation string
	Err       error
}

func (e *ConnectionError) Error() string {
	return "connection error in " + e.Operation + ": " + e.Err.Error()
}

// TransactionError represents a transaction error
type TransactionError struct {
	Operation string
	Err       error
}

func (e *TransactionError) Error() string {
	return "transaction error in " + e.Operation + ": " + e.Err.Error()
}

// TaskReader is an interface for read-only task operations
type TaskReader interface {
	Get(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
}
