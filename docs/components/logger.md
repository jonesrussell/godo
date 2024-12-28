# Logger Component

## Overview

The logger component provides a clean abstraction for application-wide logging with multiple implementations to suit different needs. It follows the interface segregation principle and allows for easy extension with new implementations.

## Core Interface

The `Logger` interface defines the minimum logging functionality needed by the application:

```go
type Logger interface {
    Debug(msg string, keysAndValues ...interface{})
    Info(msg string, keysAndValues ...interface{})
    Warn(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
    Fatal(msg string, keysAndValues ...interface{})
    WithError(err error) Logger
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
}
```

## Implementations

### Production Logger (Zap-based)

The production logger uses Uber's Zap logging library for high-performance structured logging:
- JSON or console output formats
- Configurable log levels
- Structured logging with key-value pairs
- Error context preservation
- Field attachments for context

### Test Logger

Specialized logger for testing that integrates with Go's testing package:
- Writes directly to test output
- Prefixes log levels for easy scanning
- Maintains test helper functionality
- Supports test context and cleanup

### No-op Logger

A no-op implementation for benchmarks and scenarios where logging should be disabled:
- Zero allocation overhead
- Implements full interface
- Safe for concurrent use
- Useful for performance testing

## Configuration

Logger configuration is handled via the `LogConfig` struct:

```go
type LogConfig struct {
    Level       string   // debug, info, warn, error
    Console     bool     // Use console formatting
    File        bool     // Enable file output
    FilePath    string   // Log file location
    Output      []string // Output destinations
    ErrorOutput []string // Error output destinations
}
```

## Usage Examples

### Basic Usage
```go
log, err := logger.New(config)
if err != nil {
    return err
}

log.Info("Starting application", "version", "1.0.0")
log.Debug("Configuration loaded", "config", config)
```

### With Context
```go
log = log.WithField("requestID", "123")
log.Info("Processing request")

// Multiple fields
log = log.WithFields(map[string]interface{}{
    "user":   "john",
    "action": "login",
})
```

### Error Handling
```go
if err != nil {
    log.WithError(err).Error("Operation failed")
}
```

## Testing Support

The test logger provides special support for testing scenarios:

```go
func TestFeature(t *testing.T) {
    log := logger.NewTestLogger(t)
    
    // Logs will appear in test output
    log.Info("Test started")
    
    // Test your code...
}
```

## Best Practices

1. Use structured logging with key-value pairs
2. Add context using WithField/WithFields
3. Use appropriate log levels
4. Include relevant error context
5. Keep log messages clear and actionable
6. Use consistent key names across the application

## Future Enhancements

- [ ] Add log rotation support
- [ ] Implement async logging option
- [ ] Add support for custom formatters
- [ ] Add metrics integration
- [ ] Support for log aggregation systems 