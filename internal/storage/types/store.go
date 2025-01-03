// Package types defines interfaces and types for note storage
package types

import (
	"context"
)

// Note represents a stored note
type Note struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// Store defines the interface for note storage
type Store interface {
	// Add adds a new note to the store
	Add(ctx context.Context, note Note) error

	// Get retrieves a note by ID
	Get(ctx context.Context, id string) (Note, error)

	// List returns all notes
	List(ctx context.Context) ([]Note, error)

	// Update updates an existing note
	Update(ctx context.Context, note Note) error

	// Delete removes a note by ID
	Delete(ctx context.Context, id string) error

	// Close closes the store and releases resources
	Close() error

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (Transaction, error)
}

// Transaction represents a storage transaction
type Transaction interface {
	// Add adds a new note in the transaction
	Add(ctx context.Context, note Note) error

	// Get retrieves a note by ID in the transaction
	Get(ctx context.Context, id string) (Note, error)

	// List returns all notes in the transaction
	List(ctx context.Context) ([]Note, error)

	// Update updates an existing note in the transaction
	Update(ctx context.Context, note Note) error

	// Delete removes a note by ID in the transaction
	Delete(ctx context.Context, id string) error

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}
