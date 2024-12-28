// Package storage provides interfaces and implementations for task persistence
package storage

import "errors"

var (
	// ErrTodoNotFound is returned when a todo item is not found
	ErrTodoNotFound = errors.New("todo not found")
)
