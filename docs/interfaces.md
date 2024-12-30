# Core Interfaces and Types

This document provides a comprehensive reference of how Godo's components interact through interfaces and types. It's designed to help quickly understand the system's architecture and data flows.

## System Overview

The application follows a layered architecture where:
1. User interactions come through either the GUI (Fyne-based) or HTTP API
2. These are handled by the Application layer which coordinates all components
3. Tasks are persisted through the Storage layer
4. Cross-cutting concerns like logging and configuration support all layers

## Storage Layer

The storage system uses interface segregation to provide flexible, transactional task storage:

### Task
```go
type Task struct {
    ID        string    // Unique identifier
    Content   string    // Task content
    Done      bool      // Completion status
    CreatedAt time.Time // Creation timestamp
    UpdatedAt time.Time // Last update timestamp
}
```
Core data model representing a todo item. All task operations (API/GUI) ultimately manipulate this type.

### TaskReader
```go
type TaskReader interface {
    GetByID(ctx context.Context, id string) (Task, error)
    List(ctx context.Context) ([]Task, error)
}
```
Read-only interface used by components that only need to query tasks (e.g., task list views, search functionality).
Enables read-only access patterns and supports caching implementations.

### TaskStore
```go
type TaskStore interface {
    Add(ctx context.Context, task Task) error
    GetByID(ctx context.Context, id string) (Task, error)
    Update(ctx context.Context, task Task) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]Task, error)
    Close() error
}
```
Primary interface for task storage. Used by:
- HTTP API handlers for CRUD operations
- GUI components for task management
- Quick-note feature for task creation
Implementations:
- SQLite for production (persistent)
- In-memory for testing
- Mock implementations for unit tests

### TaskTx
```go
type TaskTx interface {
    Add(ctx context.Context, task Task) error
    GetByID(ctx context.Context, id string) (Task, error)
    Update(ctx context.Context, task Task) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context) ([]Task, error)
    Commit() error
    Rollback() error
}
```
Transaction support for atomic operations. Used when:
- Multiple tasks need to be modified together
- Data consistency must be guaranteed
- Operations need rollback capability

## Application Layer

The application layer coordinates all components and manages the application lifecycle:

### Application
```go
type Application interface {
    SetupUI()    // Initializes GUI components and layouts
    Run()        // Starts all services (HTTP, hotkeys, etc.)
    Cleanup()    // Ensures graceful shutdown
}
```
Central coordinator that:
1. Initializes all components (storage, logging, GUI)
2. Sets up global hotkeys
3. Starts the HTTP server
4. Manages application lifecycle

### UI
```go
type UI interface {
    Show()
    Hide()
    SetContent(content fyne.CanvasObject)
    Resize(size fyne.Size)
    CenterOnScreen()
}
```
Common interface for all windows (main window, quick-note). Handles:
- Window lifecycle (show/hide)
- Content management
- Window positioning
- Platform-specific behaviors

## HTTP API Layer

RESTful API providing programmatic access to tasks:

### Server
```go
type Server struct {
    store  storage.TaskStore  // For task persistence
    log    logger.Logger      // Structured logging
    router *mux.Router        // HTTP routing
    srv    *http.Server       // Core HTTP server
}
```
HTTP server that:
1. Handles REST endpoints for tasks
2. Provides health checks
3. Implements middleware (logging, auth)
4. Manages graceful shutdown

### ServerConfig
```go
type ServerConfig struct {
    ReadTimeout       time.Duration
    WriteTimeout      time.Duration
    ReadHeaderTimeout time.Duration
    IdleTimeout       time.Duration
}
```
Configures server behavior for:
- Request timeouts
- Connection management
- Performance tuning

## Logger Layer

Structured logging throughout the application:

### Logger
```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}
```
Used by all components for:
- Operation tracking
- Error reporting
- Debugging
- Audit logging

Implementations:
1. Production (Zap)
   - High-performance
   - Structured logging
   - Log level control
2. Test Logger
   - Enhanced readability
   - Test-specific formatting
3. No-op Logger
   - For benchmarking
   - When logging is disabled

## Error Types

Structured error handling throughout the system:

### ValidationError
```go
type ValidationError struct {
    Field   string  // Which field failed
    Message string  // Why it failed
}
```
Used for:
- Task validation
- API request validation
- Configuration validation

### ConnectionError
```go
type ConnectionError struct {
    Operation string  // What failed
    Message   string  // Why it failed
    Err       error   // Original error
}
```
Handles:
- Database connection issues
- Network problems
- Resource availability

### TransactionError
```go
type TransactionError struct {
    Operation string  // What failed
    Message   string  // Why it failed
    Err       error   // Original error
}
```
For transaction-related failures:
- Commit failures
- Rollback errors
- Deadlocks

## Configuration Types

Configuration management throughout the application:

### AppOptions
```go
type AppOptions struct {
    Logger  LoggerOptions  // Logging configuration
    Core    CoreOptions    // Core app settings
    GUI     GUIOptions     // UI configuration
    Hotkey  HotkeyOptions  // Global hotkey settings
    HTTP    HTTPOptions    // API server config
}
```
Central configuration that:
1. Loads from YAML/environment
2. Provides defaults
3. Validates settings
4. Used during initialization

### CoreOptions
```go
type CoreOptions struct {
    StoragePath string  // Database location
    Debug       bool    // Debug mode flag
}
```
Core settings affecting:
- Data persistence
- Debug capabilities
- Development features

### GUIOptions
```go
type GUIOptions struct {
    Theme string   // UI theme selection
    Scale float64  // Display scaling
}
```
Controls:
- Visual appearance
- Accessibility features
- Platform adaptation

## Data Flows

1. Task Creation Flow:
   ```
   GUI/API → Application → TaskStore.Add → SQLite → Success/Error Response
   ```

2. Quick Note Flow:
   ```
   Hotkey → QuickNote Window → Application → TaskStore.Add → Success/Error
   ```

3. Task Update Flow:
   ```
   GUI/API → Application → TaskStore.Update → SQLite → UI Refresh
   ```

## Implementation Guidelines

1. **Storage Implementations**
   - Must be thread-safe for concurrent access
   - Should support transactions for data consistency
   - Must implement proper cleanup of resources
   - Should validate inputs before persistence
   - Must handle connection pooling efficiently

2. **Logger Implementations**
   - Should support structured logging with context
   - Must be thread-safe for concurrent logging
   - Should handle formatting consistently
   - Must not panic under any circumstances
   - Should support log levels and filtering

3. **GUI Implementations**
   - Should follow platform UI conventions
   - Must handle window lifecycle properly
   - Should support keyboard shortcuts
   - Must clean up resources on close
   - Should handle high-DPI displays

4. **Error Handling**
   - Use domain-specific error types
   - Include context in error messages
   - Support error wrapping for stack traces
   - Provide clear, actionable messages
   - Log errors with appropriate context 