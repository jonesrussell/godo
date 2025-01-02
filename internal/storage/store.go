// Package storage defines interfaces and types for task storage
package storage

import "context"

// Store combines TaskStore with transaction support
type Store interface {
	TaskStore
	BeginTx(ctx context.Context) (TaskTx, error)
}

// TaskStore defines the interface for task storage operations
type TaskStore interface {
	Add(ctx context.Context, task Task) error
	Get(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id string) error
	Close() error
}

// TaskTx defines the interface for task transactions
type TaskTx interface {
	Add(ctx context.Context, task Task) error
	Get(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id string) error
	Commit() error
	Rollback() error
}

// Task represents a todo task
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
}
