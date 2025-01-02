# Interface Design Guidelines

## Core Principles

1. **Interface Segregation**
   - Keep interfaces small and focused (max 5 methods)
   - Split based on client needs, not implementation details
   - Prefer multiple small interfaces over one large interface

2. **Interface Location**
   - Define interfaces in the package that uses them
   - Avoid central interface packages
   - Group related interfaces together

3. **Naming Conventions**
   - Use 'er' suffix for capability interfaces (Reader, Writer)
   - Use 'Service' suffix for business logic interfaces
   - Names should describe behavior, not implementation

4. **Documentation**
   - Document interface contracts clearly
   - Include usage examples
   - Specify thread-safety guarantees
   - Document error conditions

## Interface Patterns

### Reader/Writer Pattern
```go
type TaskReader interface {
    GetByID(ctx context.Context, id string) (*Task, error)
    List(ctx context.Context) ([]*Task, error)
}

type TaskWriter interface {
    Create(ctx context.Context, task *Task) error
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, id string) error
}

// Combine when both capabilities are needed
type TaskStore interface {
    TaskReader
    TaskWriter
}
```

### Service Pattern
```go
type TaskService interface {
    CreateTask(ctx context.Context, content string) (*Task, error)
    CompleteTask(ctx context.Context, id string) error
    ReopenTask(ctx context.Context, id string) error
}
```

### Transaction Pattern
```go
type TaskTx interface {
    TaskReader
    TaskWriter
    Commit() error
    Rollback() error
}
```

## Error Handling

1. **Error Types**
```go
type TaskError struct {
    Op   string // Operation that failed
    ID   string // Task ID if applicable
    Err  error  // Original error
}

func (e *TaskError) Error() string {
    return fmt.Sprintf("task operation %s failed for id %s: %v", e.Op, e.ID, e.Err)
}

func (e *TaskError) Unwrap() error {
    return e.Err
}
```

2. **Error Wrapping**
```go
if err != nil {
    return fmt.Errorf("failed to create task %s: %w", task.ID, err)
}
```

3. **Error Logging**
```go
if err != nil {
    log.Error("failed to update task",
        "id", task.ID,
        "error", err)
    return fmt.Errorf("failed to update task %s: %w", task.ID, err)
}
```

## Testing

1. **Interface Testing**
```go
func TestTaskStore(t *testing.T) {
    stores := map[string]TaskStore{
        "memory": NewMemoryStore(),
        "sqlite": NewSQLiteStore(),
    }

    for name, store := range stores {
        t.Run(name, func(t *testing.T) {
            // Test interface compliance
            var _ TaskReader = store
            var _ TaskWriter = store

            // Test operations
            testTaskCreation(t, store)
            testTaskUpdate(t, store)
            testTaskDeletion(t, store)
        })
    }
}
```

2. **Mock Generation**
```go
//go:generate mockgen -destination=mock_taskstore.go -package=mocks . TaskStore
```

## Implementation Guidelines

1. **Package Organization**
```
internal/
  ├── task/
  │   ├── service.go     # TaskService interface and implementation
  │   ├── repository.go  # TaskStore interface
  │   └── model.go      # Task type and validation
  ├── storage/
  │   ├── memory/       # In-memory implementation
  │   └── sqlite/       # SQLite implementation
  └── api/
      └── handler.go    # Uses TaskService
```

2. **Dependency Injection**
```go
type Handler struct {
    tasks TaskService
    log   Logger
}

func NewHandler(tasks TaskService, log Logger) *Handler {
    return &Handler{tasks: tasks, log: log}
}
```

## Migration Guidelines

When splitting large interfaces:

1. Create new focused interfaces
2. Update implementations
3. Create adapters if needed
4. Update clients gradually
5. Remove old interface

Example:
```go
// Step 1: Create new interfaces
type TaskReader interface {
    GetByID(ctx context.Context, id string) (*Task, error)
    List(ctx context.Context) ([]*Task, error)
}

// Step 2: Update implementation
type SQLiteStore struct { ... }
func (s *SQLiteStore) GetByID(...) { ... }
func (s *SQLiteStore) List(...) { ... }

// Step 3: Create adapter if needed
type legacyAdapter struct {
    reader TaskReader
    writer TaskWriter
}

// Step 4: Update clients to use new interfaces
func NewHandler(tasks TaskReader) *Handler {
    return &Handler{tasks: tasks}
}
``` 