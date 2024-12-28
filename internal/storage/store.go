// Package storage provides interfaces and implementations for task persistence
package storage

import "time"

// Task represents a todo task
type Task struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store defines the interface for data storage
type Store interface {
	// List returns all stored tasks
	List() ([]Task, error)

	// Add stores a new task
	Add(task Task) error

	// Update modifies an existing task
	Update(task Task) error

	// Delete removes a task by ID
	Delete(id string) error

	// GetByID retrieves a task by its ID
	GetByID(id string) (*Task, error)

	// Close releases any resources held by the store
	Close() error
}
