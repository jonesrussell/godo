package storage

import "errors"

var (
	// ErrTodoNotFound is returned when a todo item is not found
	ErrTodoNotFound = errors.New("todo not found")
)
