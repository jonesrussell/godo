// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"fmt"

	"github.com/jonesrussell/godo/internal/domain/model"
)

//go:generate mockgen -destination=../../test/mocks/mock_taskstore.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage TaskStore
//go:generate mockgen -destination=../../test/mocks/mock_tasktx.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage TaskTx
//go:generate mockgen -destination=../../test/mocks/mock_taskreader.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage TaskReader

// TaskStore defines the interface for task storage operations
type TaskStore interface {
	Add(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Task, error)
	Close() error
}

// TaskTxStore extends TaskStore with transaction support
type TaskTxStore interface {
	TaskStore
	BeginTx(ctx context.Context) (TaskTx, error)
}

// TaskTx defines the interface for task operations within a transaction
type TaskTx interface {
	Add(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id string) (model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Task, error)
	Commit() error
	Rollback() error
}

// Store is deprecated: use TaskStore instead
// Kept for backward compatibility during migration
type Store interface {
	// List returns all stored tasks
	List() ([]model.Task, error)

	// Add stores a new task
	Add(task *model.Task) error

	// Update modifies an existing task
	Update(task *model.Task) error

	// Delete removes a task by ID
	Delete(id string) error

	// GetByID retrieves a task by its ID
	GetByID(id string) (*model.Task, error)

	// Close releases any resources held by the store
	Close() error
}

// TaskReader defines the read-only interface for task storage operations
type TaskReader interface {
	GetByID(ctx context.Context, id string) (model.Task, error)
	List(ctx context.Context) ([]model.Task, error)
}

// ConnectionError represents a database connection error
type ConnectionError struct {
	Operation string
	Message   string
	Err       error
}

func (e *ConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("connection error: %s: %s: %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("connection error: %s: %s", e.Operation, e.Message)
}

// TransactionError represents a transaction error
type TransactionError struct {
	Operation string
	Message   string
	Err       error
}

func (e *TransactionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("transaction error: %s: %s: %v", e.Operation, e.Message, e.Err)
	}
	return fmt.Sprintf("transaction error: %s: %s", e.Operation, e.Message)
}
