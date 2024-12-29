# Dependency Injection in Godo

## Overview

Godo uses Google's Wire framework for dependency injection. This document outlines our DI architecture, best practices, and common patterns.

## Provider Sets

We organize providers into focused sets based on functionality:

### CoreSet
- Essential services that don't depend on UI or platform features
- Includes logging, storage, and configuration
- Base application metadata (name, version, ID)

### UISet
- UI components that depend on core services
- Fyne application and windows
- Interface bindings for GUI components

### HotkeySet
- Platform-specific hotkey functionality
- Modifier and key bindings
- Hotkey manager implementation

### HTTPSet
- HTTP server configuration
- Timeout settings
- Port configuration

### AppSet
- Main application wiring
- Combines all other sets
- Final application assembly

## Options Pattern

We use an options-based approach to prevent circular dependencies and improve configuration:

```go
type AppOptions struct {
    Core    *CoreOptions
    GUI     *GUIOptions
    HTTP    *HTTPOptions
    Hotkey  *HotkeyOptions
    Name    common.AppName
    Version common.AppVersion
    ID      common.AppID
}
```

### Benefits
- Breaks circular dependencies
- Groups related configuration
- Makes dependencies explicit
- Simplifies testing
- Improves maintainability

## Best Practices

### Provider Functions
1. Use clear naming with `Provide` prefix
2. Return cleanup functions when needed
3. Validate inputs and handle errors
4. Document dependencies clearly

Example:
```go
// ProvideSQLiteStore provides a SQLite store instance
func ProvideSQLiteStore(log logger.Logger) (*sqlite.Store, func(), error) {
    store, err := sqlite.New(dbPath, log)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create store: %w", err)
    }

    cleanup := func() {
        store.Close()
    }

    return store, cleanup, nil
}
```

### Interface Bindings
1. Use `wire.Bind` in provider sets
2. Bind concrete types to interfaces
3. Keep interface bindings close to implementations

Example:
```go
var StorageSet = wire.NewSet(
    ProvideSQLiteStore,
    wire.Bind(new(storage.TaskStore), new(*sqlite.Store)),
)
```

### Testing
1. Create separate provider sets for testing
2. Use mock implementations
3. Validate dependency initialization
4. Test cleanup functions

### Error Handling
1. Return meaningful errors from providers
2. Clean up resources on error
3. Use error wrapping
4. Validate dependencies at runtime

### Platform-Specific Code
1. Use build tags to separate implementations
2. Create platform-specific provider sets
3. Use consistent build tags across project
4. Document platform requirements

## Common Patterns

### Two-Step Initialization
For complex dependencies that require configuration:

```go
// Step 1: Configure options
func ProvideLoggerOptions(level, output, errorOutput string) *LoggerOptions {
    return &LoggerOptions{
        Level:       level,
        Output:      output,
        ErrorOutput: errorOutput,
    }
}

// Step 2: Create instance
func ProvideLogger(opts *LoggerOptions) (Logger, func(), error) {
    // Create logger using options
}
```

### Resource Cleanup
Always provide cleanup functions for resources:

```go
func ProvideResource() (*Resource, func(), error) {
    r, err := NewResource()
    if err != nil {
        return nil, nil, err
    }

    cleanup := func() {
        r.Close()
    }

    return r, cleanup, nil
}
```

### Validation
Validate dependencies and configuration:

```go
func ProvideService(dep Dependency, config Config) (*Service, error) {
    if dep == nil {
        return nil, errors.New("dependency cannot be nil")
    }
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    // Create service
}
```

## Troubleshooting

### Common Issues

1. Circular Dependencies
   - Use options pattern to break cycles
   - Split large dependencies into smaller ones
   - Use interfaces to decouple components

2. Provider Conflicts
   - Use distinct types for common values
   - Keep provider sets focused
   - Document provider dependencies

3. Testing Issues
   - Create separate test provider sets
   - Use mock implementations
   - Test cleanup functions
   - Validate initialization

### Best Practices for Maintainability

1. Keep provider sets small and focused
2. Document provider dependencies
3. Use consistent naming conventions
4. Implement proper cleanup
5. Validate at runtime
6. Use build tags appropriately
7. Keep injector functions clean (wire.Build only)
8. Follow options pattern for complex dependencies 