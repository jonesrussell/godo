// Package errors provides error definitions for the storage package
package errors

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNoteNotFound = errors.New("note not found")
	ErrStoreClosed  = errors.New("store is closed")
	ErrEmptyID      = errors.New("empty note ID")
	ErrDuplicateID  = errors.New("duplicate note ID")
)

// NotFoundError represents a note not found error
type NotFoundError struct {
	ID string
}

// Error returns the error message
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("note not found: %s", e.ID)
}

// Is reports whether the target error matches this error
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNoteNotFound
}
