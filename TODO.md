# Godo Development Tasks

## Immediate Tasks
- [x] Run wire command to generate DI implementation  ```bash
  go install github.com/google/wire/cmd/wire@latest
  cd internal/di && wire  ```
- [x] Add necessary dependencies to go.mod
  - github.com/mattn/go-sqlite3
  - github.com/google/wire

## User Interface
- [ ] Research and choose TUI library (e.g., Bubble Tea, termui)
- [ ] Design minimal UI layout
- [ ] Implement system tray integration
- [ ] Add global hotkey support
  - [ ] Research cross-platform hotkey libraries
  - [ ] Implement hotkey registration
  - [ ] Add user-configurable shortcuts

## Core Functionality
- [x] Set up basic database structure
  - [x] Define schema
  - [x] Implement repository pattern
  - [x] Create service layer
- [ ] Set up logging system
  - [ ] Choose logging library
  - [ ] Implement structured logging
  - [ ] Add log rotation
- [ ] Create configuration management
  - [ ] Define config file structure
  - [ ] Implement config file loading
  - [ ] Add config validation

## System Integration
- [ ] Implement system service functionality
  - [ ] Windows service support
  - [ ] Linux systemd support
  - [ ] macOS launchd support
- [ ] Add auto-start capability
- [ ] Implement update mechanism

## Testing
- [ ] Write unit tests for existing components
  - [ ] Repository tests
  - [ ] Service tests
  - [ ] Database tests
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