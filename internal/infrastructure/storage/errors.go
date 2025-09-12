// Package storage provides interfaces and implementations for note persistence
package storage

import (
	"errors"
	"fmt"
)

// Common errors returned by storage operations
var (
	// ErrNoteNotFound is returned when a note cannot be found
	ErrNoteNotFound = errors.New("note not found")

	// ErrDuplicateID is returned when attempting to add a note with an existing ID
	ErrDuplicateID = errors.New("note ID already exists")

	// ErrInvalidPath is returned when an invalid database path is provided
	ErrInvalidPath = errors.New("invalid database path")

	// ErrInvalidID is returned when a note ID is invalid
	ErrInvalidID = errors.New("invalid note ID")
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

// ValidationError represents a validation error from API
type ValidationError struct {
	Message string
	Fields  map[string]string
}

func (e *ValidationError) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf("validation error: %s - %v", e.Message, e.Fields)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}
