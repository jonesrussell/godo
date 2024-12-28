# Storage Layer

## Overview

The storage layer provides persistent data storage for tasks using SQLite, with an in-memory implementation for testing.

## Components

### Storage Interface (`internal/storage/store.go`)
```go
type Store interface {
    Add(task Task) error
    List() ([]Task, error)
    Update(task Task) error
    Delete(id string) error
}
```

### Implementations

#### SQLite (`internal/storage/sqlite/`)
- Primary storage implementation
- Uses modernc.org/sqlite driver
- Supports migrations
- Transaction support
- Thread-safe operations

#### In-Memory (`internal/storage/memory/`)
- Used for testing
- No persistence
- Fast operations
- Thread-safe implementation

## Task Model

```go
type Task struct {
    ID        string
    Title     string
    Completed bool
}
```

## Usage Examples

### Adding a Task
```go
task := storage.Task{
    ID:    uuid.New().String(),
    Title: "New Task",
}
err := store.Add(task)
```

### Listing Tasks
```go
tasks, err := store.List()
```

### Updating a Task
```go
task.Completed = true
err := store.Update(task)
```

### Deleting a Task
```go
err := store.Delete(taskID)
```

## Error Handling

Common errors:
- `ErrTaskNotFound`: Task doesn't exist
- Database errors: Connection, constraint violations
- Transaction errors: Rollback scenarios

## Testing

### Unit Tests
- CRUD operation tests
- Error case tests
- Transaction tests

### Integration Tests
- Database connection tests
- Migration tests
- Concurrent operation tests

## Configuration

Storage configuration in `configs/default.yaml`:
```yaml
database:
  path: "godo.db"  # SQLite database path
```

## Best Practices

1. Always use transactions for multi-operation changes
2. Handle database errors appropriately
3. Use prepared statements for repeated operations
4. Close database connections properly
5. Use appropriate indices for performance
6. Follow repository pattern conventions 