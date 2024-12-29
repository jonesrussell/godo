# Testing Guidelines

## Overview
This document outlines the testing practices and patterns used in the Godo application.

## Test Types

### Unit Tests
- Test individual components in isolation
- Use mock implementations for dependencies
- Focus on single piece of functionality
- Use table-driven tests for multiple cases
- Example:
```go
func TestQuickNote_Show(t *testing.T) {
    // Create test dependencies
    testApp := test.NewApp()
    store := &mockStore{}
    log := logger.NewTestLogger(t)
    cfg := config.WindowConfig{
        Width:  200,
        Height: 100,
    }

    // Create quick note window
    quickNote := NewQuickNote(testApp, store, log, cfg)

    // Test initial state
    assert.NotNil(t, quickNote.input)
    assert.Equal(t, "", quickNote.input.Text)

    // Test behavior
    quickNote.Show()
    assert.Equal(t, quickNote.input, quickNote.window.Canvas().Focused())
}
```

### Integration Tests
- Test component interactions
- Use real implementations where practical
- Test complete workflows
- Example:
```go
func TestTaskWorkflow(t *testing.T) {
    store := sqlite.NewStore()
    defer store.Close()

    // Add task
    task := storage.Task{ID: "test-1", Content: "Test Task"}
    err := store.Add(context.Background(), task)
    require.NoError(t, err)

    // Verify task
    tasks, err := store.List(context.Background())
    require.NoError(t, err)
    assert.Contains(t, tasks, task)
}
```

### UI Tests
- Test window creation and visibility
- Test input focus behavior
- Test window state management
- Use Fyne's test package
- Example:
```go
func TestWindowVisibility(t *testing.T) {
    app := test.NewApp()
    win := app.NewWindow("Test")
    
    assert.False(t, win.Visible())
    win.Show()
    assert.True(t, win.Visible())
}
```

## Test Helpers

### Mock Store
```go
type mockStore struct {
    storage.TaskStore
}

func (m *mockStore) List(ctx context.Context) ([]storage.Task, error) {
    return []storage.Task{}, nil
}
```

### Test Logger
```go
func NewTestLogger(t *testing.T) Logger {
    cfg := &Config{
        Level:       "debug",
        Development: true,
        Encoding:    "console",
    }
    logger, _, err := NewLogger(cfg)
    if err != nil {
        t.Fatalf("Failed to create test logger: %v", err)
    }
    return logger
}
```

## Best Practices

### Test Organization
- Use descriptive test names
- Group related tests using subtests
- Use setup and teardown helpers
- Clean up resources properly

### Assertions
- Use testify/assert for most checks
- Use testify/require for critical checks
- Include meaningful error messages
- Example:
```go
assert.NotNil(t, obj, "Object should be initialized")
require.NoError(t, err, "Operation should succeed")
```

### Mocking
- Create minimal mock implementations
- Mock only what's necessary
- Use interfaces for dependencies
- Document mock behavior

### Test Coverage
- Aim for high coverage of critical paths
- Test both success and error cases
- Test edge cases and boundaries
- Test concurrent operations where relevant

## Running Tests

### All Tests
```bash
task test
```

### Specific Package
```bash
go test ./internal/gui/quicknote
```

### With Coverage
```bash
go test -cover ./...
```

## Test Tags
- Use build tags to control test execution
- Example: `//go:build !docker`
- Common tags:
  - `windows`
  - `linux`
  - `docker`
  - `integration`

## Debugging Tests
- Use `t.Log` for debug output
- Use `t.Logf` for formatted output
- Run specific test: `go test -run TestName`
- Run with verbose output: `go test -v`

## Test Documentation
- Document test purpose
- Document test prerequisites
- Document test data requirements
- Document expected outcomes

## Continuous Integration
- Tests run on every PR
- Coverage reports generated
- Test results published
- Performance tracked 