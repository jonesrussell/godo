// Package errors provides error types for storage operations
package errors

import "errors"

var (
	// ErrNoteNotFound is returned when a note is not found
	ErrNoteNotFound = errors.New("note not found")

	// ErrNoteExists is returned when a note already exists
	ErrNoteExists = errors.New("note already exists")

	// ErrTransactionClosed is returned when attempting to use a closed transaction
	ErrTransactionClosed = errors.New("transaction is closed")

	// ErrStoreClosed is returned when attempting to use a closed store
	ErrStoreClosed = errors.New("store is closed")
)
