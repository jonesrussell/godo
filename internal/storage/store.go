// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"time"
)

// Task represents a todo item
type Task struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskStore defines the interface for task storage operations
type TaskStore interface {
	Add(ctx context.Context, task Task) error
	GetByID(ctx context.Context, id string) (Task, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Task, error)
	Close() error
}

// TaskTxStore extends TaskStore with transaction support
type TaskTxStore interface {
	TaskStore
	BeginTx(ctx context.Context) (TaskTx, error)
}

// TaskTx defines the interface for task operations within a transaction
type TaskTx interface {
	Add(ctx context.Context, task Task) error
	GetByID(ctx context.Context, id string) (Task, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]Task, error)
	Commit() error
	Rollback() error
}

// Store is deprecated: use TaskStore instead
// Kept for backward compatibility during migration
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
