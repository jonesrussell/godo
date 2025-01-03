package note

import "context"

// Reader provides read operations for notes
type Reader interface {
	// Get retrieves a note by ID
	// Returns NotFound error if the note doesn't exist
	Get(ctx context.Context, id string) (*Note, error)

	// List returns all notes
	// Returns empty slice if no notes exist
	List(ctx context.Context) ([]*Note, error)
}

// Writer provides write operations for notes
type Writer interface {
	// Add creates a new note
	// Returns ValidationFailed error if the note is invalid
	Add(ctx context.Context, note *Note) error

	// Update modifies an existing note
	// Returns NotFound error if the note doesn't exist
	Update(ctx context.Context, note *Note) error

	// Delete removes a note by ID
	// Returns NotFound error if the note doesn't exist
	Delete(ctx context.Context, id string) error
}

// Store combines read and write operations with lifecycle management
type Store interface {
	Reader
	Writer

	// Close releases any resources held by the store
	Close() error

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (Transaction, error)
}

// Transaction represents a storage transaction
type Transaction interface {
	Reader
	Writer

	// Commit commits the transaction
	Commit() error

	// Rollback aborts the transaction
	Rollback() error
}
