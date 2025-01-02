# Error Handling Guidelines

## Core Principles

1. **Domain-Specific Errors**
   - Create custom error types for each domain
   - Include operation context
   - Support error wrapping
   - Provide clear error messages

2. **Error Wrapping**
   - Always wrap errors with context
   - Use `fmt.Errorf` with `%w` verb
   - Include operation name
   - Add relevant details

3. **Error Logging**
   - Use structured logging
   - Include operation context
   - Log at appropriate levels
   - Avoid sensitive data

## Error Types

### 1. Base Error Type
```go
type Error struct {
    Op   string // Operation that failed
    Kind Kind   // Category of error
    Err  error  // Underlying error
}

func (e *Error) Error() string {
    return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *Error) Unwrap() error {
    return e.Err
}
```

### 2. Domain Errors
```go
// Task domain errors
type TaskError struct {
    Op   string
    Kind TaskErrorKind
    ID   string // Task ID if applicable
    Err  error
}

type TaskErrorKind int

const (
    TaskNotFound TaskErrorKind = iota
    TaskInvalidState
    TaskValidationFailed
)

// Storage domain errors
type StorageError struct {
    Op    string
    Kind  StorageErrorKind
    Table string
    Err   error
}

type StorageErrorKind int

const (
    ConnectionFailed StorageErrorKind = iota
    QueryFailed
    TransactionFailed
)
```

## Error Creation

### 1. Constructor Functions
```go
func NewTaskError(op string, kind TaskErrorKind, id string, err error) *TaskError {
    return &TaskError{
        Op:   op,
        Kind: kind,
        ID:   id,
        Err:  err,
    }
}

// Usage
if err != nil {
    return NewTaskError("CreateTask", TaskValidationFailed, task.ID, err)
}
```

### 2. Error Wrapping
```go
// With operation context
if err != nil {
    return fmt.Errorf("failed to create task: %w", err)
}

// With additional context
if err != nil {
    return fmt.Errorf("failed to update task %s: %w", id, err)
}
```

## Error Handling

### 1. Service Layer
```go
func (s *TaskService) CreateTask(ctx context.Context, content string) (*Task, error) {
    // Validate input
    if content == "" {
        return nil, NewTaskError("CreateTask", TaskValidationFailed, "", 
            errors.New("content cannot be empty"))
    }

    // Create task
    task := NewTask(content)
    if err := s.store.Create(ctx, task); err != nil {
        s.logger.Error("failed to create task",
            "op", "CreateTask",
            "content", content,
            "error", err)
        return nil, fmt.Errorf("failed to create task: %w", err)
    }

    return task, nil
}
```

### 2. Storage Layer
```go
func (s *SQLiteStore) Create(ctx context.Context, task *Task) error {
    query := `INSERT INTO tasks (id, content, done) VALUES (?, ?, ?)`
    
    if err := s.db.ExecContext(ctx, query, task.ID, task.Content, task.Done); err != nil {
        return NewStorageError("Create", QueryFailed, "tasks", err)
    }
    
    return nil
}
```

### 3. API Layer
```go
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
    var req CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("failed to decode request",
            "op", "CreateTask",
            "error", err)
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    task, err := h.service.CreateTask(r.Context(), req.Content)
    if err != nil {
        h.handleError(w, err)
        return
    }

    json.NewEncoder(w).Encode(task)
}
```

## Error Response Handling

### 1. HTTP Responses
```go
func (h *Handler) handleError(w http.ResponseWriter, err error) {
    var taskErr *TaskError
    if errors.As(err, &taskErr) {
        switch taskErr.Kind {
        case TaskNotFound:
            http.Error(w, "task not found", http.StatusNotFound)
        case TaskValidationFailed:
            http.Error(w, "invalid task", http.StatusBadRequest)
        default:
            http.Error(w, "internal error", http.StatusInternalServerError)
        }
        return
    }

    http.Error(w, "internal error", http.StatusInternalServerError)
}
```

### 2. GUI Error Handling
```go
func (w *Window) handleError(err error) {
    var taskErr *TaskError
    if errors.As(err, &taskErr) {
        switch taskErr.Kind {
        case TaskNotFound:
            w.showError("Task not found")
        case TaskValidationFailed:
            w.showError("Invalid task")
        default:
            w.showError("An error occurred")
        }
        return
    }

    w.showError("An unexpected error occurred")
}
```

## Error Logging

### 1. Structured Logging
```go
func (s *Service) logError(op string, err error, fields ...interface{}) {
    s.logger.Error("operation failed",
        append([]interface{}{
            "op", op,
            "error", err,
        }, fields...)...)
}

// Usage
s.logError("CreateTask",
    err,
    "content", content,
    "user", userID)
```

### 2. Log Levels
```go
// Error: Operation failures
logger.Error("failed to create task",
    "op", "CreateTask",
    "error", err)

// Warn: Recoverable issues
logger.Warn("retrying operation",
    "op", "CreateTask",
    "attempt", attempt)

// Info: Normal operations
logger.Info("task created",
    "op", "CreateTask",
    "id", task.ID)

// Debug: Detailed information
logger.Debug("validating task",
    "op", "CreateTask",
    "content", content)
```

## Testing

### 1. Error Testing
```go
func TestCreateTask(t *testing.T) {
    tests := []struct {
        name    string
        content string
        wantErr error
    }{
        {
            name:    "empty content",
            content: "",
            wantErr: &TaskError{
                Op:   "CreateTask",
                Kind: TaskValidationFailed,
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := service.CreateTask(context.Background(), tt.content)
            if tt.wantErr != nil {
                require.Error(t, err)
                var taskErr *TaskError
                require.ErrorAs(t, err, &taskErr)
                assert.Equal(t, tt.wantErr.Op, taskErr.Op)
                assert.Equal(t, tt.wantErr.Kind, taskErr.Kind)
            }
        })
    }
}
```

### 2. Error Mocking
```go
func TestHandleError(t *testing.T) {
    mockLogger := mocks.NewLogger(t)
    mockLogger.EXPECT().Error(
        "failed to create task",
        "op", "CreateTask",
        "error", mock.Anything,
    ).Once()

    service := NewService(mockStore, mockLogger)
    err := service.CreateTask(context.Background(), "")
    require.Error(t, err)
}
```

## Best Practices

1. **Error Creation**
   - Use constructor functions
   - Include operation context
   - Add relevant details
   - Support error wrapping

2. **Error Handling**
   - Check error types with `errors.As`
   - Handle specific error cases
   - Provide clear messages
   - Log appropriately

3. **Error Response**
   - Map errors to status codes
   - Sanitize error messages
   - Include request IDs
   - Handle all error types

4. **Logging**
   - Use structured logging
   - Include context
   - Choose appropriate levels
   - Avoid sensitive data

5. **Testing**
   - Test error conditions
   - Verify error types
   - Check error messages
   - Mock error scenarios 