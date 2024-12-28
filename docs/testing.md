# Testing Guide

## Overview

This guide covers testing practices for the Godo application, including unit tests, integration tests, and testing utilities.

## Test Structure

### Directory Organization
```
internal/
  ├── app/
  │   └── app_test.go
  ├── storage/
  │   ├── sqlite/
  │   │   ├── store_test.go
  │   │   ├── migrations_test.go
  │   │   └── testing.go
  │   └── memory/
  │       └── memory_test.go
  └── testutil/
      └── store.go
```

## Coverage Goals

### Critical Packages
- Storage implementations: >65% coverage
- API handlers: >60% coverage
- Core business logic: >80% coverage

### Current Coverage (as of latest)
- `internal/storage/sqlite`: 66.3%
- `internal/model`: 100%
- `internal/config`: 89%
- `internal/common`: 65%
- `internal/api`: 61.9%

## Test Types

### Unit Tests
- Individual package functionality
- Mocked dependencies
- Fast execution
- Comprehensive error case coverage

### Integration Tests
- Cross-package functionality
- Real dependencies
- Database operations
- GUI interactions

### End-to-End Tests
- Complete workflows
- System integration
- User scenarios

## Testing Tools

### Standard Library
- `testing` package
- `httptest` package
- `iotest` package

### Third-Party
- `testify` for assertions
- `go-sqlmock` for database tests
- `gomock` for interface mocking

## Running Tests

### All Tests
```bash
task test
```

### With Coverage Report
```bash
task test:cover
```

### Specific Package
```bash
go test ./internal/storage/...
```

### With Race Detection
```bash
task test:race
```

## Writing Tests

### Comprehensive Test Example
```go
func TestFeatureComprehensive(t *testing.T) {
    // Setup with cleanup
    store, cleanup := setupTestStore(t)
    defer cleanup()

    // Test normal operation
    t.Run("Success Case", func(t *testing.T) {
        result, err := store.Operation()
        assert.NoError(t, err)
        assert.NotNil(t, result)
    })

    // Test validation
    t.Run("Validation", func(t *testing.T) {
        // Test empty input
        _, err := store.Operation("")
        assert.ErrorIs(t, err, ErrEmptyInput)

        // Test invalid input
        _, err = store.Operation("invalid")
        assert.Error(t, err)
    })

    // Test error conditions
    t.Run("Error Cases", func(t *testing.T) {
        // Test not found
        _, err := store.Get("nonexistent")
        assert.ErrorIs(t, err, ErrNotFound)

        // Test closed connection
        store.Close()
        _, err = store.Operation()
        assert.ErrorIs(t, err, ErrStoreClosed)
    })
}
```

### State Management Tests
```go
func TestStateManagement(t *testing.T) {
    store := NewStore()
    
    // Test initial state
    assert.False(t, store.IsClosed())
    
    // Test after close
    store.Close()
    assert.True(t, store.IsClosed())
    
    // Test operations after close
    _, err := store.Operation()
    assert.ErrorIs(t, err, ErrStoreClosed)
}
```

## Test Utilities

### Mock Store
```go
type MockStore struct {
    tasks []storage.Task
}

func NewMockStore() *MockStore {
    return &MockStore{
        tasks: make([]storage.Task, 0),
    }
}
```

### Test Helpers
```go
func setupTestDB(t *testing.T) (*sql.DB, func()) {
    t.Helper()
    
    db, err := sql.Open("sqlite", ":memory:")
    require.NoError(t, err)
    
    return db, func() {
        db.Close()
    }
}
```

## Best Practices

1. Test Organization
   - One test file per source file
   - Clear test names
   - Proper setup and cleanup
   - Use subtests for organization

2. Test Coverage
   - Maintain minimum coverage targets per package
   - Test all error conditions
   - Test state transitions
   - Test concurrent operations
   - Test resource cleanup

3. Error Testing
   - Test all custom error types
   - Verify error wrapping
   - Test error conditions in order
   - Use ErrorIs for error comparison

4. Resource Management
   - Use defer for cleanup
   - Test cleanup operations
   - Verify resource state
   - Test connection handling

## Common Patterns

### Setup and Teardown
```go
func TestMain(m *testing.M) {
    // Setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Cleanup
    teardown()
    
    os.Exit(code)
}
```

### Context Testing
```go
func TestWithContext(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    
    result, err := operationWithContext(ctx)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Concurrent Testing
```go
func TestConcurrent(t *testing.T) {
    t.Parallel()
    
    store := testutil.NewMockStore()
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Test concurrent operations
        }()
    }
    
    wg.Wait()
}
```

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify](https://pkg.go.dev/github.com/stretchr/testify)
- [Go Testing Blog](https://blog.golang.org/cover)
- [Advanced Testing Tips](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests) 