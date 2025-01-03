# Storage Component

The storage component provides a unified interface for storing and retrieving notes in the application.

## Interface

The storage interface is defined in `internal/storage/store.go`:

```go
type Store interface {
    // Add adds a new note to the store
    Add(ctx context.Context, note Note) error

    // Get retrieves a note by its ID
    Get(ctx context.Context, id string) (Note, error)

    // Update updates an existing note
    Update(ctx context.Context, note Note) error

    // Delete removes a note by its ID
    Delete(ctx context.Context, id string) error

    // List returns all notes in the store
    List(ctx context.Context) ([]Note, error)

    // Close closes the store and releases any resources
    Close() error
}
```

## Data Model

The core data model is the `Note` struct:

```go
type Note struct {
    ID        string `json:"id"`
    Content   string `json:"content"`
    Completed bool   `json:"completed"`
    CreatedAt int64  `json:"created_at"`
    UpdatedAt int64  `json:"updated_at"`
}
```

## Implementations

### SQLite Store

The SQLite implementation uses the `modernc.org/sqlite` package to store notes in a SQLite database. The implementation is in `internal/storage/sqlite/store.go`.

Key features:
- Uses prepared statements for all operations
- Handles concurrent access with proper locking
- Supports automatic schema migrations
- Provides efficient indexing for common queries

### Memory Store

The in-memory implementation stores notes in memory using a map. This is primarily used for testing and development. The implementation is in `internal/storage/memory/store.go`.

Key features:
- Fast access for all operations
- No persistence between application restarts
- Thread-safe operations using a mutex
- Useful for testing and prototyping

### Mock Store

A mock implementation is provided for testing in `internal/storage/mock/store.go`. This allows tests to verify storage operations without requiring a real database.

Key features:
- Configurable behavior for testing edge cases
- Records method calls for verification
- Can be used to simulate errors and edge cases

## Error Handling

The storage component defines several error types in `internal/storage/errors/errors.go`:

- `ErrNotFound`: Returned when a note is not found
- `ErrInvalidID`: Returned when an invalid ID is provided
- `ErrInvalidNote`: Returned when a note is invalid
- `ErrDuplicate`: Returned when trying to add a note with an existing ID

## Usage Example

```go
// Create a new SQLite store
store, err := sqlite.NewStore("notes.db")
if err != nil {
    log.Fatal(err)
}
defer store.Close()

// Add a new note
note := storage.Note{
    ID:        uuid.New().String(),
    Content:   "Buy groceries",
    CreatedAt: time.Now().Unix(),
    UpdatedAt: time.Now().Unix(),
}
if err := store.Add(context.Background(), note); err != nil {
    log.Fatal(err)
}

// List all notes
notes, err := store.List(context.Background())
if err != nil {
    log.Fatal(err)
}
for _, note := range notes {
    fmt.Printf("Note: %s\n", note.Content)
}
``` 