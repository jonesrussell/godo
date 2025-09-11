// Package storage provides interfaces and implementations for task persistence
package storage

import (
	"context"
	"fmt"

	"github.com/jonesrussell/godo/internal/domain/model"
)

//go:generate mockgen -destination=../../test/mocks/mock_notestore.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage NoteStore
//go:generate mockgen -destination=../../test/mocks/mock_notetx.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage NoteTx
//go:generate mockgen -destination=../../test/mocks/mock_notereader.go -package=mocks github.com/jonesrussell/godo/internal/infrastructure/storage NoteReader

// NoteStore defines the interface for note storage operations
type NoteStore interface {
	Add(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id string) (model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Note, error)
	Close() error
}

// NoteTxStore extends NoteStore with transaction support
type NoteTxStore interface {
	NoteStore
	BeginTx(ctx context.Context) (NoteTx, error)
}

// NoteTx defines the interface for note operations within a transaction
type NoteTx interface {
	Add(ctx context.Context, note *model.Note) error
	GetByID(ctx context.Context, id string) (model.Note, error)
	Update(ctx context.Context, note *model.Note) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Note, error)
	Commit() error
	Rollback() error
}

// Store is deprecated: use NoteStore instead
// Kept for backward compatibility during migration
type Store interface {
	// List returns all stored notes
	List() ([]model.Note, error)

	// Add stores a new note
	Add(note *model.Note) error

	// Update modifies an existing note
	Update(note *model.Note) error

	// Delete removes a note by ID
	Delete(id string) error

	// GetByID retrieves a note by its ID
	GetByID(id string) (*model.Note, error)

	// Close releases any resources held by the store
	Close() error
}

// NoteReader defines the read-only interface for note storage operations
type NoteReader interface {
	GetByID(ctx context.Context, id string) (model.Note, error)
	List(ctx context.Context) ([]model.Note, error)
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
