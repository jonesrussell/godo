// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"errors"
	"time"
)

// ErrTaskNotFound is returned when a task cannot be found
var ErrTaskNotFound = errors.New("task not found")

// Task represents a todo task in the storage layer
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// Store defines the interface for task storage
type Store interface {
	Add(task Task) error
	List() ([]Task, error)
	Update(task Task) error
	Delete(id string) error
}
