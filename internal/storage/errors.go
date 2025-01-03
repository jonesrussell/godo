// Package storage provides interfaces and implementations for note persistence
package storage

import (
	"context"
	"errors"

	"github.com/jonesrussell/godo/internal/storage/types"
)

// Common errors returned by storage operations
var (
	// ErrNoteNotFound is returned when a note cannot be found
	ErrNoteNotFound = errors.New("note not found")

	// ErrStoreClosed is returned when attempting to use a closed store
	ErrStoreClosed = errors.New("store is closed")

	// ErrEmptyID is returned when an empty note ID is provided
	ErrEmptyID = errors.New("note ID cannot be empty")

	// ErrDuplicateID is returned when attempting to add a note with an existing ID
	ErrDuplicateID = errors.New("note ID already exists")

	// ErrInvalidPath is returned when an invalid database path is provided
	ErrInvalidPath = errors.New("invalid database path")

	// ErrInvalidID is returned when a note ID is invalid
	ErrInvalidID = errors.New("invalid note ID")

	// ErrTransactionNotSupported is returned when transactions are not supported
	ErrTransactionNotSupported = errors.New("transactions not supported")
)

// NotFoundError is returned when a note cannot be found
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return "note not found: " + e.ID
}

// Is implements errors.Is interface to match against ErrNoteNotFound
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNoteNotFound
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

// NoteReader is an interface for read-only note operations
type NoteReader interface {
	Get(ctx context.Context, id string) (types.Note, error)
	List(ctx context.Context) ([]types.Note, error)
}
