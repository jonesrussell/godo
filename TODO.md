# Godo Development Tasks

## Current Focus
- [ ] Fix build system
  - [x] Add proper Docker support
  - [x] Add build tags for Docker/non-Docker environments
  - [ ] Implement platform-specific hotkey managers
    - [x] Windows implementation
    - [ ] Linux implementation
    - [ ] macOS implementation
    - [x] Docker mock implementation
  - [ ] Add proper cross-compilation support
  - [ ] Add release packaging
  - [ ] Add CI/CD pipeline

## Core Implementation Progress
- [x] Step 1: Add Basic Logging
  - [x] Add zap logger initialization
  - [x] Add basic logging to track application lifecycle
  - [x] Log startup, shutdown, and main operations

- [x] Step 2: System Tray Integration
  - [x] Add system tray icon
  - [x] Add application icon
  - [x] Hide main window by default
  - [x] Move quick note trigger to system tray menu
  - [x] Add logging for system tray events

- [x] Step 3: Quick Note Implementation
  - [x] Move quick note logic to separate package
  - [x] Keep the same functionality but callable from system tray
  - [x] Add logging for quick note operations
  - [x] Add ESC key to close window
  - [x] Auto-focus input field on show

- [x] Step 4: Basic Todo Storage
  - [x] Add simple in-memory todo storage
  - [x] Create basic todo model with UUID and timestamps
  - [x] Add logging for todo operations

- [x] Step 5: Persistence
  - [x] Add SQLite storage implementation
  - [x] Implement basic CRUD operations
  - [x] Add migration support
  - [x] Add logging for database operations

## Next Steps
- [ ] Step 6: Polish
  - [x] Add keyboard shortcuts
  - [x] Add quick note menu item
  - [ ] Add error dialogs for operation failures
  - [ ] Add task completion animations
  - [ ] Improve UI layout
  - [ ] Add basic error handling
  - [ ] Enhance logging with contextual information
  - [ ] Add log rotation
  - [ ] Add auto-start capability
  - [ ] Add update mechanism

- [ ] Step 7: Todo List UI
  - [x] Add a list view to display all todos
  - [x] Add ability to mark todos as done
  - [x] Add ability to delete todos
  - [ ] Add task sorting (by date, completion)
  - [ ] Add task filtering
  - [ ] Show todo creation time and last update time
  - [ ] Add task editing capability

## Future Enhancements
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Data export/import
- [ ] Task priority levels
- [ ] Recurring tasks
- [ ] Multiple todo lists
- [ ] Cloud sync support
