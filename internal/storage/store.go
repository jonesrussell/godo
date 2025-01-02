// Package storage defines interfaces and types for task storage
package storage

import "context"

// Store defines the interface for task storage implementations
type Store interface {
	// List returns all tasks
	List(ctx context.Context) ([]Task, error)

	// Get returns a task by ID
	Get(ctx context.Context, id string) (Task, error)

	// Add creates a new task
	Add(ctx context.Context, task Task) error

	// Update modifies an existing task
	Update(ctx context.Context, task Task) error

	// Delete removes a task by ID
	Delete(ctx context.Context, id string) error

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (Transaction, error)

	// Close cleans up any resources
	Close() error
}

// Transaction represents a database transaction
type Transaction interface {
	// List returns all tasks in the transaction
	List(ctx context.Context) ([]Task, error)

	// Get returns a task by ID in the transaction
	Get(ctx context.Context, id string) (Task, error)

	// Add creates a new task in the transaction
	Add(ctx context.Context, task Task) error

	// Update modifies an existing task in the transaction
	Update(ctx context.Context, task Task) error

	// Delete removes a task by ID in the transaction
	Delete(ctx context.Context, id string) error

	// Commit commits the transaction
	Commit() error

	// Rollback aborts the transaction
	Rollback() error
}

// Task represents a todo item
type Task struct {
	ID        string `json:"id"`
	Title     string `json:"title"`     // Changed from Content
	Completed bool   `json:"completed"` // Changed from Done
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
