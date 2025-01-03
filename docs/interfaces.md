# Interface Design Guidelines

## Core Principles

1. **Domain-Driven Design**
   - Domain types belong in `internal/domain/{entity}`
   - Domain types must implement `Validate()`
   - Domain types own their error types
   - Use domain-specific error wrapping

2. **Interface Segregation**
   - Keep interfaces small and focused (max 5 methods)
   - Split based on capabilities (e.g., Reader/Writer)
   - Interfaces belong in the domain package
   - Name interfaces after capabilities with 'er' suffix

3. **Error Handling**
   - Use domain-specific error types
   - Include operation context in errors
   - Wrap all errors with domain errors
   - Provide error kind constants

4. **Documentation**
   - Document interface contracts
   - Specify error conditions
   - Include thread-safety guarantees
   - Provide usage examples

## Current Implementation

### Domain Types
```go
// Note represents a task or quick note in the system
type Note struct {
    ID        string    `json:"id"`
    Content   string    `json:"content"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Error represents a note-related error
type Error struct {
    Op   string    // Operation that failed
    Kind ErrorKind // Type of error
    Msg  string    // Error details
    Err  error     // Underlying error
}

// ErrorKind represents the type of note error
type ErrorKind int

const (
    Unknown ErrorKind = iota
    ValidationFailed
    NotFound
    StorageError
)
```

### Core Interfaces
```go
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
```

## Implementation Rules

1. **Domain Package Structure**
```
internal/domain/note/
  ├── note.go       # Note type and methods
  ├── errors.go     # Error types and constants
  └── store.go      # Store interfaces
```

2. **Error Handling Pattern**
```go
// Always wrap errors with domain error types
if err != nil {
    return &Error{
        Op:   "Store.Add",
        Kind: StorageError,
        Msg:  "failed to insert note",
        Err:  err,
    }
}

// Check error types using errors.As
var nErr *note.Error
if errors.As(err, &nErr) {
    switch nErr.Kind {
    case note.NotFound:
        // Handle not found
    case note.ValidationFailed:
        // Handle validation error
    }
}
```

3. **Implementation Example**
```go
// Store implements the note.Store interface
type Store struct {
    db *sql.DB
}

// Add adds a new note
func (s *Store) Add(ctx context.Context, n *note.Note) error {
    // Validate domain object first
    if err := n.Validate(); err != nil {
        return err
    }

    // Perform storage operation
    if err := s.insert(ctx, n); err != nil {
        return &note.Error{
            Op:   "Store.Add",
            Kind: note.StorageError,
            Msg:  "failed to insert note",
            Err:  err,
        }
    }

    return nil
}
```

## Testing Requirements

1. **Interface Compliance**
```go
var _ note.Store = (*Store)(nil)
var _ note.Transaction = (*Transaction)(nil)
```

2. **Error Testing**
```go
t.Run("Not Found Error", func(t *testing.T) {
    _, err := store.Get(ctx, "nonexistent")
    assert.Error(t, err)
    var nErr *note.Error
    assert.ErrorAs(t, err, &nErr)
    assert.Equal(t, note.NotFound, nErr.Kind)
})
```

3. **Transaction Testing**
```go
t.Run("Transaction Rollback", func(t *testing.T) {
    tx, err := store.BeginTx(ctx)
    require.NoError(t, err)
    
    // ... perform operations ...
    
    err = tx.Rollback()
    require.NoError(t, err)
    
    // Verify rollback
    _, err = store.Get(ctx, id)
    assert.ErrorIs(t, err, note.NotFound)
})
```

## Linter Rules

The custom linter enforces:
1. Domain types must be in domain package
2. Domain types must implement Validate()
3. Interfaces must have ≤ 5 methods
4. Interface names must end with 'er' or 'Service'
5. Errors must be wrapped with domain error types
6. Store implementations must support transactions 