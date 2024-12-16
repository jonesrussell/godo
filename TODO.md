# Godo Development Tasks

## Immediate Tasks
- [x] Run wire command to generate DI implementation  ```bash
  go install github.com/google/wire/cmd/wire@latest
  cd internal/di && wire  ```
- [x] Add necessary dependencies to go.mod
  - github.com/mattn/go-sqlite3
  - github.com/google/wire
  - github.com/charmbracelet/bubbletea
  - go.uber.org/zap
  - github.com/getlantern/systray

## User Interface
- [x] Research and choose TUI library (e.g., Bubble Tea, termui)
- [x] Design minimal UI layout
- [x] Implement basic UI with Bubble Tea
- [x] Implement system tray integration
  - [x] Add system tray icon
  - [x] Add basic menu (Open Manager, Quit)
  - [x] Implement clean shutdown
- [x] Add global hotkey support
  - [x] Research cross-platform hotkey libraries
  - [x] Implement hotkey registration
  - [ ] Add user-configurable shortcuts
- [ ] Separate quick-note and management UIs
  - [x] Quick-note: Minimal, instant input
  - [x] Management: Full featured todo interface

## Core Functionality
- [x] Set up basic database structure
  - [x] Define schema
  - [x] Implement repository pattern
  - [x] Create service layer
- [x] Set up logging system
  - [x] Choose logging library (zap)
  - [x] Implement structured logging
  - [ ] Add log rotation
- [ ] Create configuration management
  - [ ] Define config file structure
  - [ ] Implement config file loading
  - [ ] Add config validation
- [ ] Add version information
  - [ ] Create version package
  - [ ] Add build-time version injection
  - [ ] Display version in UI and logs

## System Integration
- [x] Implement graceful shutdown
  - [x] Handle context cancellation
  - [x] Clean systray removal
  - [x] Proper process termination
- [ ] Implement system service functionality
  - [ ] Windows service support
  - [ ] Linux systemd support
  - [ ] macOS launchd support
- [ ] Add auto-start capability
- [ ] Implement update mechanism

## Testing
- [x] Write initial unit tests
  - [x] Service tests
  - [x] Repository tests
  - [ ] Database tests
  - [x] Hotkey system tests
    - [x] Constants validation
    - [x] Type definitions
    - [x] Win32 API wrappers
- [ ] Set up integration tests
- [ ] Add CI pipeline

## Documentation
- [x] Initial README setup
- [ ] Complete API documentation
- [ ] Add usage examples
- [ ] Create user guide
- [ ] Document hotkey combinations
- [ ] Add installation instructions for all platforms

## Future Enhancements
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Data export/import
- [ ] Task priority levels
- [ ] Recurring tasks
- [ ] Multiple todo lists
- [ ] Cloud sync support

## Quick-Note Feature (High Priority)
- [x] Create QuickNoteUI component
  - [x] Single-line input field
  - [x] Minimal window decoration
  - [ ] Transparent/floating window
  - [ ] Position near cursor
- [x] Implement quick-note workflow
  - [x] Hotkey triggers QuickNoteUI
  - [x] Focus input immediately
  - [x] Enter saves note and closes
  - [x] Esc cancels and closes
  - [x] Return to previous window focus
- [ ] Add quick-note settings
  - [ ] Window position (cursor/center/custom)
  - [ ] Window opacity
  - [ ] Custom hotkey binding
    - [ ] Make Ctrl+Alt+G configurable via config file
    - [ ] Add hotkey configuration UI
    - [ ] Validate hotkey combinations
    - [ ] Save/load hotkey preferences
  - [ ] Auto-categorization rules
- [x] Error Handling Improvements
  - [x] Handle "The operation completed successfully" error message
  - [ ] Add retry mechanism for failed note saves
  - [ ] Improve error messages in UI
  - [ ] Add notification for successful saves

## System Tray Integration (Completed)
- [x] Add system tray icon support
- [x] Create basic menu structure
- [x] Implement proper shutdown handling
- [x] Add Open Manager option
- [x] Handle icon loading
- [x] Implement clean exit
