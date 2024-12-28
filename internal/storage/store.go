// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"io"
	"time"
)

// Task represents a todo task
type Task struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskReader defines read-only operations for tasks
type TaskReader interface {
	// List returns all stored tasks
	List(ctx context.Context) ([]Task, error)

	// GetByID retrieves a task by its ID
	GetByID(ctx context.Context, id string) (*Task, error)
}

// TaskWriter defines write operations for tasks
type TaskWriter interface {
	// Add stores a new task
	Add(ctx context.Context, task Task) error

	// Update modifies an existing task
	Update(ctx context.Context, task Task) error

	// Delete removes a task by ID
	Delete(ctx context.Context, id string) error
}

// TaskStore combines read and write operations with resource management
type TaskStore interface {
	TaskReader
	TaskWriter
	io.Closer
}

// TaskTx defines a transaction interface for atomic operations
type TaskTx interface {
	TaskStore

	// Begin starts a new transaction
	Begin(ctx context.Context) (TaskTx, error)

	// Commit commits the transaction
	Commit() error

	// Rollback aborts the transaction
	Rollback() error
}

// ValidationError represents a task validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ConnectionError represents a database connection error
type ConnectionError struct {
	Operation string
	Err       error
}

func (e *ConnectionError) Error() string {
	return "storage connection error during " + e.Operation + ": " + e.Err.Error()
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// TransactionError represents a transaction-related error
type TransactionError struct {
	Operation string
	Err       error
}

func (e *TransactionError) Error() string {
	return "transaction error during " + e.Operation + ": " + e.Err.Error()
}

func (e *TransactionError) Unwrap() error {
	return e.Err
}

// NotFoundError represents a missing task error
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return "task not found: " + e.ID
}

// DuplicateIDError represents a duplicate task ID error
type DuplicateIDError struct {
	ID string
}

func (e *DuplicateIDError) Error() string {
	return "duplicate task ID: " + e.ID
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
