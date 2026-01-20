# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Godo is a cross-platform note application combining three powerful features:
1. Global hotkey-triggered quick note capture (Ctrl+Shift+1)
2. Full-featured graphical note management interface (Ctrl+Shift+2)
3. REST API for programmatic note management

Built with Go 1.25+ using Fyne toolkit for native GUI, supporting Windows and Linux (macOS coming soon).

## Essential Commands

### Development Workflow
```bash
# Install development tools (uses Go 1.25 tool directive in go.mod)
task install-tools

# Format code (runs gofmt, goimports, go vet)
task fmt

# Run linting (depends on fmt)
task lint

# Quick format and lint check (uses --fast flag)
task check

# Run tests
task test                    # Basic tests with wireinject tag
task test:race              # With race detection
task test:cover             # With coverage report (output to coverage/)

# WSL2-specific testing with X11 forwarding
task test:wsl2

# Run the application
task run                    # Normal mode
task run-debug              # With LOG_LEVEL=debug

# Build for current platform
task build                   # Uses main Taskfile.yml (cross-platform)

# Complete development cycle
task dev                    # Format, lint, test
task all                    # deps, fmt, lint, test, build
```

### Platform-Specific Build Notes
- Uses Task's OS-specific taskfile feature: `Taskfile_{{OS}}.yml`
- `Taskfile_windows.yml`: Windows-specific variables (BUILD_TIME uses PowerShell)
- `Taskfile_linux.yml`: Linux-specific variables (BUILD_TIME uses date command)
- Both main Taskfile.yml and Taskfile.build.yml include OS-specific taskfiles
- BUILD_TIME variable automatically uses correct command for the platform

### Additional Build Commands
The main Taskfile.yml includes Taskfile.build.yml, so you can access build tasks with the `build:` prefix:

```bash
# Native builds (current platform)
task build                          # Build for current platform (alias for build:native)
task build:native:windows          # Windows native build
task build:native:linux            # Linux native build

# Cross-compilation
task build:cross-windows           # Cross-compile for Windows from Linux/WSL2
task build:cross-linux             # Cross-compile for Linux from Windows

# Docker builds
task build:docker:build-all        # Build both Windows and Linux using Docker
task build:docker:build-linux      # Build Linux binary using Docker
task build:docker:build-windows    # Build Windows binary using Docker

# Wire generation for specific platforms
task build:wire                    # Generate Wire code for current platform
task build:wire:windows           # Generate Wire code for Windows
task build:wire:linux             # Generate Wire code for Linux

# Cleanup
task build:clean                   # Clean build artifacts (dist/, wire_gen.go)

# WSL2-specific workflows
task build:wsl2:all               # Complete WSL2 workflow (build + copy to Windows)
task build:wsl2:build-windows     # Build Windows executable from WSL2
```

### Testing Single Files or Packages
```bash
# Run tests for specific package
go test ./internal/infrastructure/storage/sqlite/... -v

# Run specific test
go test ./internal/domain/service/... -run TestNoteService_CreateNote -v

# Run with race detection
go test -race ./internal/infrastructure/hotkey/...
```

### Dependency Management
```bash
# IMPORTANT: Do NOT use `go mod vendor` - breaks CGO dependencies
# Always use:
go mod download
go mod tidy

# Development tools are managed via Go 1.25 tool directive
# Tools are defined in go.mod:
#   tool (
#     github.com/golangci/golangci-lint/cmd/golangci-lint
#     golang.org/x/tools/cmd/goimports
#     github.com/google/wire/cmd/wire
#   )
# Install with: task install-tools (runs `go get` and `go mod tidy`)
```

### Mock Generation
```bash
task mocks              # Generate all mocks
task mocks:storage      # Storage interface mocks only
task mocks:service      # Service interface mocks only
task mocks:gui          # GUI interface mocks only
task mocks:logger       # Logger interface mocks only
```

### Wire Dependency Injection
```bash
# Regenerate Wire code for current platform
task wire

# Or directly:
wire ./internal/application/container
```

### OS-Specific Taskfiles
The project uses Task's OS-specific taskfile feature for platform-dependent variables:

```yaml
# Taskfile.yml includes OS-specific taskfiles
includes:
  os:
    taskfile: ./Taskfile_{{OS}}.yml  # Automatically loads Taskfile_windows.yml or Taskfile_linux.yml
```

**Files:**
- `Taskfile_windows.yml`: Windows-specific variables and tasks (BUILD_TIME uses PowerShell)
- `Taskfile_linux.yml`: Linux-specific variables and tasks (BUILD_TIME uses date command)

**Usage:**
- Variables from OS-specific taskfiles are referenced with `.os.` prefix
- Example: `BUILD_TIME: {ref: .os.BUILD_TIME}` in main taskfile
- Task automatically selects the correct file based on runtime OS

## Architecture

### Clean Architecture Layers

```
internal/
├── domain/              # Pure business logic, no external dependencies
│   ├── model/          # Core entities (Note)
│   ├── repository/     # Repository interfaces
│   ├── service/        # Business logic services (NoteService)
│   └── storage/        # Storage interfaces (UnifiedNoteStorage)
│
├── application/         # Application orchestration
│   ├── core/           # Main app logic and interfaces
│   └── container/      # Wire dependency injection
│       ├── wire.go     # Wire provider sets
│       ├── build_*.go  # Platform-specific builds
│       └── build_constraints.go  # Unsupported platform checks
│
├── infrastructure/      # External concerns and implementations
│   ├── storage/        # Storage implementations
│   │   ├── sqlite/     # SQLite backend (modernc.org/sqlite - pure Go)
│   │   ├── api/        # HTTP API backend
│   │   ├── memory/     # In-memory (testing)
│   │   └── factory/    # Factory pattern for storage selection
│   ├── gui/            # Fyne GUI components
│   │   ├── mainwindow/ # Full note management interface
│   │   ├── quicknote/  # Minimal quick capture popup
│   │   └── theme/      # Custom styling
│   ├── hotkey/         # Global hotkey management (csturiale/hotkey)
│   ├── api/            # HTTP server (Gorilla Mux)
│   ├── logger/         # Structured logging (Zap)
│   └── platform/       # Platform detection (WSL2, headless, GUI support)
│
└── config/             # Viper-based configuration with validation
```

### Key Architectural Patterns

**Dependency Injection (Wire)**: All dependencies managed through Google Wire. Provider sets organized by concern:
- `ConfigSet`: Configuration loading
- `LoggingSet`: Zap logger setup
- `StorageSet`: Unified storage with factory pattern
- `ServiceSet`: Business logic (NoteService, NoteRepository)
- `UISet`: Fyne GUI components
- `CoreSet`: Combines Config + Logging + Storage + Service
- `AppSet`: Main application instance

**Storage Abstraction**: Three-layer architecture supporting multiple backends:
1. Domain interface: `UnifiedNoteStorage` in [internal/domain/storage/interfaces.go](internal/domain/storage/interfaces.go)
2. Factory pattern: `NewUnifiedStorage()` creates backend based on config (`StorageTypeSQLite` or `StorageTypeAPI`)
3. Backend implementations: SQLite (pure Go), API (HTTP client with retry logic), Memory (testing)

**Platform-Specific Builds**: Uses build tags and runtime detection:
- Build tags: `//go:build linux || windows` for hotkey manager
- Build tags: `//go:build docker` for no-op GUI in containers
- Runtime detection: `IsWSL2()`, `IsHeadless()`, `SupportsGUI()` in [internal/infrastructure/platform/wsl.go](internal/infrastructure/platform/wsl.go)

**Hotkey Management**: Global hotkey system with graceful degradation:
- Manager interface with Register/Unregister/Start/Stop lifecycle
- Platform-specific modifier handling (Alt differs Windows vs Linux)
- Channel-based event loop in goroutine
- Validates not running in WSL2 or headless environments

**GUI Components**: Thread-safe Fyne widgets:
- All UI operations use `fyne.Do()` for thread safety
- MainWindow: Full CRUD with list view, dialogs, status bar
- QuickNote: Minimal popup with Ctrl+Enter submit, Escape hide
- Docker build: No-op implementation prevents Fyne init in containers

**Testing Strategy**: mockgen-based mocking:
- All mocks generated to [internal/test/mocks/](internal/test/mocks/) via `//go:generate mockgen`
- Test logger implementations: testing.go (test-friendly), noop.go (silent)
- Run `task mocks` to regenerate all mocks

## Important Development Notes

### CGO Dependencies
- **CRITICAL**: This project uses CGO dependencies (`golang.design/x/hotkey`)
- Never use `go mod vendor` - it breaks CGO
- Always use `go mod download` and `go mod tidy`

### Wire Code Generation
- Wire files use `//go:build wireinject` tag
- Tests must include `-tags=wireinject` flag: `go test ./... -tags=wireinject`
- After modifying Wire providers, run `task wire` to regenerate

### Configuration System
- Viper-based with YAML support ([config.yaml](config.yaml))
- Environment variable override: `GODO_*` prefix
- Storage type: `sqlite` or `api` (configured in config.yaml)
- Hotkey bindings: Configurable modifiers + key

### Platform Support
- **Windows**: Full support with system tray integration
- **Linux**: Full support (no system tray)
- **macOS**: Coming soon
- **WSL2**: GUI works with X11 forwarding, but hotkeys don't work (use `task test:wsl2`)

### Build Prerequisites
- Go 1.25+
- SQLite3
- MinGW-w64 GCC (Windows builds - required for CGO)
- GNU diffutils (Windows linting)
- Task runner (`choco install go-task`)
- Development tools are managed via Go 1.25 `tool` directive in go.mod

### Thread Safety in GUI Code
Always wrap GUI operations in `fyne.Do()`:
```go
fyne.Do(func() {
    window.Show()
    canvas.Focus(widget)
})
```

### Naming Conventions
- **IMPORTANT**: Never use "Enhanced" prefix in any naming
- Use descriptive names that reflect functionality
- Follow Go conventions: PascalCase for exported, camelCase for private

## API Endpoints

Base URL: `http://localhost:8080`

- `GET /health` - Health check
- `GET /api/v1/notes` - List all notes
- `POST /api/v1/notes` - Create note
- `PUT /api/v1/notes/{id}` - Update note
- `DELETE /api/v1/notes/{id}` - Delete note

Default HTTP port: 8008 (configured in [config.yaml](config.yaml))

## Coding Standards

### Error Handling
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Use sentinel errors for expected conditions: `var ErrNotFound = errors.New("not found")`
- Log detailed errors internally, return user-friendly messages

### Logging
- Use structured logging with Zap logger from [internal/infrastructure/logger/](internal/infrastructure/logger/)
- Include context in log messages with key-value pairs
- Levels: Debug, Info, Warn, Error

### Testing
- Test files end with `_test.go`
- Use table-driven tests for multiple scenarios
- Naming: `TestType_Method_Scenario`
- Aim for 80%+ coverage on business logic
- Use `t.TempDir()` for temporary files

### Formatting
- Run `task fmt` before committing (runs gofmt, goimports, go vet)
- Use `gofumpt` for code formatting
- Enable `golangci-lint` with `--fast` flag for quick checks

## Application Startup Flow

1. `main.go` calls `container.InitializeApp()`
2. Wire resolves dependency graph
3. Config loaded from YAML/environment
4. Logger initialized
5. Storage factory creates backend (SQLite/API)
6. Repository wraps storage
7. Service layer created
8. Fyne App initialized
9. GUI components created (MainWindow, QuickNote)
10. Hotkey Manager configured
11. API Server started in background
12. UI setup based on config
13. Hotkeys registered with OS
14. `App.Run()` starts Fyne event loop (blocks)

## Request Flow Example

User creates note via QuickNote:
```
QuickNote.addNote()
  → model.NewNote(content)
  → store.Add(ctx, note)              # NoteStoreAdapter
    → unifiedStorage.CreateNote()      # SQLite or API
      → sqlite: Direct DB insert
      → api: HTTP POST /notes
  → QuickNote.showStatus()             # UI feedback
```

## Configuration

Storage backend selection in [config.yaml](config.yaml):
```yaml
storage:
  type: "sqlite"  # Options: "sqlite" or "api"
  sqlite:
    file_path: "$HOME/.config/godo/godo.db"
  api:
    base_url: "https://lame.ddev.site/api"
    timeout_seconds: 30
    retry_count: 3
    retry_delay_ms: 1000
    insecure_skip_verify: true
```

Hotkey configuration:
```yaml
hotkeys:
  quick_note:
    modifiers: ["Ctrl", "Shift"]
    key: "1"
  main_window:
    modifiers: ["Ctrl", "Shift"]
    key: "2"
```

## Cross-Cutting Concerns

**Structured Logging**: Zap-based logging throughout with Debug/Info/Warn/Error levels, key-value pairs for context, coordinated shutdown via `Sync()`

**Error Handling**: Custom error types (ValidationError, NotFoundError), error mapping at boundaries (storage → domain), HTTP status code mapping in API

**Thread Safety**: `fyne.Do()` for UI operations, `sync.RWMutex` in memory storage, context propagation for cancellation

**Graceful Shutdown**: API server with context-based shutdown, cleanup functions from Wire, signal handling in main.go
