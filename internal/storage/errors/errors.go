// Package errors provides error definitions for the storage package
package errors

import "errors"

// Common errors
var (
	// ErrNoteNotFound is returned when a note is not found
	ErrNoteNotFound = errors.New("note not found")
	// ErrNoteAlreadyExists is returned when attempting to add a note that already exists
	ErrNoteAlreadyExists = errors.New("note already exists")
	// ErrInvalidNoteID is returned when a note ID is invalid
	ErrInvalidNoteID = errors.New("invalid note ID")
	// ErrInvalidNoteTitle is returned when a note title is invalid
	ErrInvalidNoteTitle = errors.New("invalid note title")
	// ErrTransactionClosed is returned when attempting to use a closed transaction
	ErrTransactionClosed = errors.New("transaction is closed")
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
