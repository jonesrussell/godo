# Godo Development Tasks

## Immediate Tasks
- [x] Run wire command to generate DI implementation  ```bash
  go install github.com/google/wire/cmd/wire@latest
  cd internal/app && wire  ```
- [x] Add necessary dependencies to go.mod
  - github.com/mattn/go-sqlite3
  - github.com/google/wire
  - go.uber.org/zap
  - fyne.io/fyne/v2

## Core Implementation Progress
- [x] Step 1: Add Basic Logging
  - [x] Add zap logger initialization
  - [x] Add basic logging to track application lifecycle
  - [x] Log startup, shutdown, and main operations
  - [x] Reference: `internal/logger/logger.go`

- [x] Step 2: System Tray Integration
  - [x] Add system tray icon (using favicon.ico)
  - [x] Add application icon (using Icon.png)
  - [x] Hide main window by default
  - [x] Move quick note trigger to system tray menu
  - [x] Add logging for system tray events
  - [x] Reference: `cmd/godo/main.go`

- [x] Step 3: Quick Note Implementation
  - [x] Move quick note logic to separate package
  - [x] Keep the same functionality but callable from system tray
  - [x] Add logging for quick note operations
  - [x] Reference: `internal/gui/quicknote/quicknote.go`

- [x] Step 4: Basic Todo Storage
  - [x] Add simple in-memory todo storage
  - [x] Create basic todo model with UUID and timestamps
  - [x] Add logging for todo operations
  - [x] Reference: `internal/model/todo.go`, `internal/storage/memory/memory.go`

- [x] Step 5: Persistence
  - [x] Add SQLite storage implementation
  - [x] Implement basic CRUD operations
  - [x] Add migration support
  - [x] Add logging for database operations
  - [x] Reference: `internal/storage/sqlite/sqlite.go`

## Current Focus: Step 6 - Polish
- [ ] Add keyboard shortcuts
- [ ] Improve UI layout
- [ ] Add basic error handling
- [ ] Enhance logging with contextual information
- [ ] Add log rotation
- [ ] Add auto-start capability
- [ ] Add update mechanism

## User Interface
- [ ] Step 7: Todo List UI
  - [ ] Add a list view to display all todos
  - [ ] Add ability to mark todos as done
  - [ ] Add ability to delete todos
  - [ ] Show todo creation time and last update time
  - [ ] Reference: `internal/gui/todolist/todolist.go`
- [x] Research and choose GUI library (Fyne)
- [x] Design minimal UI layout
- [ ] Implement basic UI with Fyne
  - [ ] Main window layout
  - [ ] Task list view
  - [ ] Task input form
  - [ ] Task actions (complete/delete)
- [x] Add global hotkey support
  - [x] Research cross-platform hotkey libraries
  - [x] Implement hotkey registration
  - [ ] Add user-configurable shortcuts

## Future Enhancements
- [ ] macOS Support
  - [ ] Set up OSXCross for cross-compilation
  - [ ] Add Darwin-specific build tags
  - [ ] Test on macOS
  - [ ] Add macOS-specific features (if any)
- [ ] Terminal User Interface (TUI) alternative
  - [ ] Research TUI libraries (Bubble Tea, etc.)
  - [ ] Design TUI layout
  - [ ] Implement basic TUI functionality
  - [ ] Ensure feature parity with GUI
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Data export/import
- [ ] Task priority levels
- [ ] Recurring tasks
- [ ] Multiple todo lists
- [ ] Cloud sync support
