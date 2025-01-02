# Architecture Overview

## Core Principles

1. **Clean Architecture**
   - Business logic in `internal/`
   - Dependencies point inward
   - Interfaces defined by consumers
   - Platform-specific code isolated

2. **Interface Segregation**
   - Small, focused interfaces
   - Split by client needs
   - Maximum 5 methods per interface
   - Clear responsibility boundaries

3. **Error Handling**
   - Domain-specific error types
   - Error wrapping with context
   - Structured logging
   - Clear error hierarchies

4. **Dependency Management**
   - Wire for dependency injection
   - Explicit dependencies
   - No global state
   - Testable components

## Layer Organization

```
cmd/
  ├── godo/              # Main application
  └── godo-linter/       # Custom linter tool

internal/
  ├── app/              # Application core
  │   ├── service.go    # Business logic
  │   └── app.go        # App lifecycle
  │
  ├── task/            # Task domain
  │   ├── model.go     # Task entity
  │   ├── service.go   # Task operations
  │   └── errors.go    # Domain errors
  │
  ├── storage/         # Data persistence
  │   ├── memory/      # In-memory store
  │   └── sqlite/      # SQLite store
  │
  ├── api/            # HTTP API
  │   ├── server.go   # API server
  │   ├── handler.go  # Request handlers
  │   └── middleware/ # API middleware
  │
  ├── gui/           # User interface
  │   ├── window.go  # Window management
  │   ├── task/      # Task-related UI
  │   └── theme/     # UI theming
  │
  └── platform/      # Platform-specific
      ├── win/       # Windows features
      └── linux/     # Linux features

pkg/                 # Reusable packages
  ├── config/        # Configuration
  ├── log/          # Logging
  └── errors/       # Error utilities
```

## Component Interactions

### 1. Task Management Flow
```
GUI/API → TaskService → TaskStore → Database
   ↑          ↓            ↓
   └──────── Events ←──────┘
```

### 2. Quick Note Flow
```
Hotkey → QuickNote UI → TaskService → TaskStore
   ↑          ↓             ↓
   └──────── Events ←───────┘
```

### 3. Data Flow
```
User Input → Validation → Business Logic → Storage
    ↑           ↓             ↓             ↓
    └───────── Error Handling & Logging ────┘
```

## Interface Design

### 1. Service Layer
```go
type TaskService interface {
    CreateTask(ctx context.Context, content string) (*Task, error)
    CompleteTask(ctx context.Context, id string) error
    ReopenTask(ctx context.Context, id string) error
}
```

### 2. Storage Layer
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

type TaskStore interface {
    TaskReader
    TaskWriter
    Close() error
}
```

### 3. UI Layer
```go
type Window interface {
    Show()
    Hide()
    SetContent(content fyne.CanvasObject)
    Close() error
}

type TaskView interface {
    Refresh()
    SetTasks([]*Task)
    OnTaskSelected(func(*Task))
}
```

## Error Handling

### 1. Domain Errors
```go
type TaskError struct {
    Op   string
    Kind TaskErrorKind
    Err  error
}

type TaskErrorKind int

const (
    TaskNotFound TaskErrorKind = iota
    TaskInvalidState
    TaskValidationFailed
)
```

### 2. Error Wrapping
```go
if err != nil {
    return &TaskError{
        Op:   "CreateTask",
        Kind: TaskValidationFailed,
        Err:  err,
    }
}
```

### 3. Error Logging
```go
logger.Error("failed to create task",
    "op", "CreateTask",
    "content", content,
    "error", err)
```

## Testing Strategy

### 1. Unit Tests
```go
func TestTaskService(t *testing.T) {
    suite.Run(t, new(TaskServiceSuite))
}

type TaskServiceSuite struct {
    suite.Suite
    store  *mocks.TaskStore
    logger *mocks.Logger
    svc    *TaskService
}
```

### 2. Integration Tests
```go
func TestTaskIntegration(t *testing.T) {
    store := sqlite.NewStore(":memory:")
    defer store.Close()
    
    svc := task.NewService(store)
    // Test full flow
}
```

### 3. API Tests
```go
func TestTaskAPI(t *testing.T) {
    srv := api.NewTestServer()
    defer srv.Close()
    
    // Test HTTP endpoints
}
```

## Configuration Management

### 1. Application Config
```yaml
app:
  name: "Godo"
  version: "1.0.0"
  
storage:
  type: "sqlite"
  path: "tasks.db"
  
api:
  port: 8080
  timeout: 30s
  
gui:
  theme: "dark"
  scale: 1.0
```

### 2. Loading Config
```go
type Config struct {
    App     AppConfig
    Storage StorageConfig
    API     APIConfig
    GUI     GUIConfig
}

func LoadConfig() (*Config, error) {
    // Load and validate configuration
}
```

## Dependency Injection

### 1. Wire Setup
```go
//+build wireinject

func InitializeApp(cfg *Config) (*App, error) {
    wire.Build(
        NewLogger,
        NewTaskStore,
        NewTaskService,
        NewAPIServer,
        NewGUI,
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

### 2. Component Setup
```go
type App struct {
    cfg     *Config
    logger  Logger
    store   TaskStore
    service TaskService
    api     *APIServer
    gui     *GUI
}
```

## Platform Integration

### 1. Windows Support
```go
//go:build windows

func (app *App) setupHotkeys() error {
    // Windows-specific hotkey registration
}
```

### 2. Linux Support
```go
//go:build linux

func (app *App) setupHotkeys() error {
    // Linux-specific hotkey registration
}
```

## Security Considerations

1. **Data Protection**
   - Secure storage of tasks
   - Proper file permissions
   - Safe error messages

2. **API Security**
   - Input validation
   - Rate limiting
   - CORS configuration

3. **Platform Security**
   - Safe hotkey registration
   - Protected IPC
   - Resource cleanup

## Performance Optimization

1. **Database**
   - Connection pooling
   - Prepared statements
   - Index optimization

2. **API**
   - Response caching
   - Compression
   - Connection reuse

3. **GUI**
   - Lazy loading
   - Resource pooling
   - Event debouncing 