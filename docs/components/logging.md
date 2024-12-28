# Logging System

## Overview

The logging system uses Uber's Zap logger for structured, high-performance logging with different output formats and log levels.

## Components

### Logger Interface (`internal/logger/logger.go`)
```go
type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warn(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Fatal(msg string, fields ...interface{})
}
```

### Implementation

#### Zap Logger (`internal/logger/zap.go`)
- Production-ready logging
- Structured JSON output
- Multiple output targets
- Log level control
- Performance optimized

## Configuration

Logger configuration in `configs/default.yaml`:
```yaml
logger:
  level: "info"    # debug, info, warn, error
  console: true    # stdout logging
  file: true       # file logging
  path: "logs/"    # log file directory
```

## Usage Examples

### Basic Logging
```go
logger.Info("Application started")
logger.Debug("Processing task", "taskID", task.ID)
logger.Error("Failed to save task", "error", err)
```

### Structured Fields
```go
logger.Info("Task updated",
    "taskID", task.ID,
    "title", task.Title,
    "completed", task.Completed,
)
```

### Error Logging
```go
if err != nil {
    logger.Error("Operation failed",
        "operation", "SaveTask",
        "error", err,
        "stack", debug.Stack(),
    )
}
```

## Log Levels

1. **Debug**
   - Detailed debugging information
   - Development use
   - Disabled in production

2. **Info**
   - General operational events
   - State changes
   - User actions

3. **Warn**
   - Warning conditions
   - Recoverable errors
   - Potential issues

4. **Error**
   - Error conditions
   - Application can continue
   - Requires attention

5. **Fatal**
   - Severe errors
   - Application cannot continue
   - Immediate shutdown

## Best Practices

1. Logging Guidelines
   - Use appropriate log levels
   - Include relevant context
   - Avoid sensitive information
   - Use structured fields

2. Performance
   - Use sugar logger for simple cases
   - Avoid expensive operations in log calls
   - Use sampling in high-volume scenarios

3. Error Handling
   - Log errors with stack traces
   - Include operation context
   - Use error wrapping

## Testing

### Log Capture
```go
// Example test with log capture
func TestWithLogs(t *testing.T) {
    logs := captureTestLogs()
    defer logs.Restore()
    
    // Test code
    
    assert.Contains(t, logs.All(), "Expected log message")
}
```

### Log Level Tests
- Verify level filtering
- Check structured fields
- Validate output format

## Common Patterns

### Application Events
```go
logger.Info("Application event",
    "event", "StartUp",
    "version", version,
    "config", configPath,
)
```

### User Actions
```go
logger.Info("User action",
    "action", "CreateTask",
    "userID", userID,
    "taskID", taskID,
)
```

### Error Context
```go
logger.Error("Operation failed",
    "operation", op,
    "error", err,
    "context", ctx,
    "stack", debug.Stack(),
)
```

## Resources

- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Logging Best Practices](https://dave.cheney.net/2015/11/05/lets-talk-about-logging)
- [Go Error Handling](https://blog.golang.org/error-handling-and-go) 