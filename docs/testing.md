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
  │   │   └── sqlite_test.go
  │   └── memory/
  │       └── memory_test.go
  └── testutil/
      └── store.go
```

## Test Types

### Unit Tests
- Individual package functionality
- Mocked dependencies
- Fast execution

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

### Specific Package
```bash
go test ./internal/storage/...
```

### With Race Detection
```bash
go test -race ./...
```

### With Coverage
```bash
go test -cover ./...
```

## Writing Tests

### Basic Test Structure
```go
func TestFeature(t *testing.T) {
    // Setup
    store := testutil.NewMockStore()
    
    // Test
    result, err := store.Add(task)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Table-Driven Tests
```go
func TestOperation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid case", "input", "expected", false},
        {"error case", "bad", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Operation(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### Integration Test Example
```go
func TestDatabaseIntegration(t *testing.T) {
    // Setup temporary database
    db, cleanup := setupTestDB(t)
    defer cleanup()
    
    // Run tests
    store := sqlite.New(db)
    task := storage.Task{Title: "Test"}
    
    err := store.Add(task)
    assert.NoError(t, err)
    
    tasks, err := store.List()
    assert.NoError(t, err)
    assert.Len(t, tasks, 1)
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
   - Aim for high coverage
   - Test edge cases
   - Test error conditions
   - Test concurrent operations

3. Test Performance
   - Fast unit tests
   - Parallel test execution
   - Efficient setup/teardown
   - Use test caching

4. Test Maintenance
   - Keep tests simple
   - Don't test implementation details
   - Use test helpers
   - Document complex tests

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