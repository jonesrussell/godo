// Package storage provides storage functionality
package storage

import "context"

// Store defines the interface for note storage implementations
type Store interface {
	// List returns all notes
	List(ctx context.Context) ([]Note, error)

	// Get returns a note by ID
	Get(ctx context.Context, id string) (Note, error)

	// Add creates a new note
	Add(ctx context.Context, note Note) error

	// Update modifies an existing note
	Update(ctx context.Context, note Note) error

	// Delete removes a note by ID
	Delete(ctx context.Context, id string) error

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (Transaction, error)

	// Close cleans up any resources
	Close() error
}

// Transaction represents a database transaction
type Transaction interface {
	// List returns all notes in the transaction
	List(ctx context.Context) ([]Note, error)

	// Get returns a note by ID in the transaction
	Get(ctx context.Context, id string) (Note, error)

	// Add creates a new note in the transaction
	Add(ctx context.Context, note Note) error

	// Update modifies an existing note in the transaction
	Update(ctx context.Context, note Note) error

	// Delete removes a note by ID in the transaction
	Delete(ctx context.Context, id string) error

	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error
}

// Note represents a quick note item
type Note struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Completed bool   `json:"completed"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
