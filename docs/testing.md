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

### Platform-Specific Tests
- Test platform-specific features separately
- Use build tags to control test execution
- Test platform-specific error cases
- Example:
```go
//go:build windows && !linux && !darwin && !docker

func TestWindowsHotkey(t *testing.T) {
    // Test Windows-specific hotkey behavior
    manager := NewWindowsManager(logger)
    err := manager.Register()
    assert.NoError(t, err)
    
    // Test cleanup
    err = manager.Unregister()
    assert.NoError(t, err)
}
```

### Lifecycle Tests
- Test component initialization
- Test proper cleanup and shutdown
- Test resource management
- Example:
```go
func TestComponentLifecycle(t *testing.T) {
    // Test initialization
    component := New(deps...)
    assert.NotNil(t, component)

    // Test operation
    err := component.Start()
    assert.NoError(t, err)

    // Test cleanup
    err = component.Stop()
    assert.NoError(t, err)
    // Verify resources are released
}
```

### Error Path Tests
- Test all error conditions
- Verify error messages
- Test recovery from errors
- Example:
```go
func TestErrorHandling(t *testing.T) {
    cases := []struct {
        name          string
        input         Input
        expectedErr   error
        errorContains string
    }{
        {
            name:          "invalid state",
            input:         invalidInput,
            expectedErr:   ErrInvalidState,
            errorContains: "invalid state",
        },
        // ... more cases ...
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            err := component.Operation(tc.input)
            assert.Error(t, err)
            assert.ErrorIs(t, err, tc.expectedErr)
            assert.Contains(t, err.Error(), tc.errorContains)
        })
    }
}

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
- Test cleanup and resource management
- Test component lifecycles
- Test platform-specific features
- Test error messages and recovery
- Test invalid state transitions

### Resource Management
- Always clean up resources in tests
- Use defer for cleanup operations
- Verify cleanup was successful
- Test resource leak scenarios
- Example:
```go
func TestResourceManagement(t *testing.T) {
    // Create resources
    resource := NewResource()
    defer func() {
        err := resource.Close()
        assert.NoError(t, err, "Cleanup should succeed")
    }()

    // Use resource
    err := resource.Operation()
    assert.NoError(t, err)

    // Verify no leaks
    assert.Zero(t, resource.ActiveConnections())
}
```

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