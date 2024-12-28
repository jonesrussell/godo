// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"io"
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

// TaskReader defines methods for reading tasks
type TaskReader interface {
	// List returns all tasks
	List(ctx context.Context) ([]Task, error)
	// GetByID returns a task by its ID
	GetByID(ctx context.Context, id string) (*Task, error)
}

// TaskWriter defines methods for writing tasks
type TaskWriter interface {
	// Add creates a new task
	Add(ctx context.Context, task Task) error
	// Update replaces an existing task
	Update(ctx context.Context, task Task) error
	// Delete removes a task
	Delete(ctx context.Context, id string) error
}

// TaskStore combines TaskReader and TaskWriter with io.Closer
type TaskStore interface {
	TaskReader
	TaskWriter
	io.Closer
}

// TaskTx represents a transaction for task operations
type TaskTx interface {
	TaskReader
	TaskWriter
	// Commit commits the transaction
	Commit() error
	// Rollback rolls back the transaction
	Rollback() error
}

// TaskTxStore extends TaskStore with transaction support
type TaskTxStore interface {
	TaskStore
	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (TaskTx, error)
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
