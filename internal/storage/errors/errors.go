// Package errors provides error definitions for the storage package
package errors

import "errors"

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
)