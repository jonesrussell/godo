# TODO List

## Core Features

### Completed ✓
- [x] Quick note capture via global hotkey
- [x] Main window hidden on startup
- [x] Systray icon and menu
- [x] Basic task management
- [x] SQLite storage
- [x] Human-readable logging in development
- [x] Regression tests for core functionality
- [x] Hotkey lifecycle management
- [x] Proper cleanup on exit
- [x] Custom linter implementation

### In Progress
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Task priorities
- [ ] Multiple task lists
- [ ] Recurring tasks
- [ ] Task search and filtering

## Code Cleanup

### High Priority
- [ ] Interface Segregation (Linter Findings):
  - [ ] Split TaskStore interface (storage/store.go)
  - [ ] Split TaskTx interface (storage/store.go)
  - [ ] Split Store interface (storage/store.go)
  - [ ] Split TaskHandler interface (api/handler.go)
  - [ ] Split MainWindow interface (gui/interfaces.go)
  - [ ] Add interface documentation
  - [ ] Update interface tests
- [ ] Error Handling Standardization (Linter Findings):
  - [ ] Add error wrapping in storage layer (sqlite/store.go)
  - [ ] Add error wrapping in config layer (config/config.go)
  - [ ] Add error wrapping in hotkey layer (app/hotkey/)
  - [ ] Add error wrapping in container layer (container/wire_gen.go)
  - [ ] Create error wrapping guidelines
  - [ ] Update error documentation
- [ ] Test Assertions Standardization (Linter Findings):
  - [ ] Update API tests to use testify
  - [ ] Update app tests to use testify
  - [ ] Update config tests to use testify
  - [ ] Update GUI tests to use testify
  - [ ] Update storage tests to use testify
  - [ ] Create test assertion guidelines

### Technical Debt
- [ ] Code Duplication Removal:
  - [ ] Consolidate window management code
  - [ ] Extract common UI patterns
  - [ ] Standardize test setup code
  - [ ] Create shared test utilities
- [ ] Test Coverage:
  - [ ] Add service layer tests
  - [ ] Update storage tests
  - [ ] Add integration tests
  - [ ] Improve error case coverage

## Testing

### Completed ✓
- [x] Unit tests for core packages
- [x] Integration tests for storage
- [x] Mock implementations for testing
- [x] Test coverage for UI components
- [x] Regression tests for systray and hotkeys
- [x] Hotkey lifecycle tests
- [x] Cleanup and resource management tests
- [x] Error path testing
- [x] Edge case coverage
- [x] Platform-specific test coverage (Windows)

### In Progress
- [ ] Performance benchmarks
- [ ] Load testing
- [ ] End-to-end tests
- [ ] Test coverage improvements:
  - [ ] Window positioning tests
  - [ ] Keyboard shortcut tests
  - [ ] Visual regression tests
  - [ ] Component interaction tests
- [ ] Cross-platform test suite:
  - [ ] Linux-specific features
  - [ ] macOS preparation
  - [ ] Platform-specific UI behavior
- [ ] Integration tests:
  - [ ] Systray integration
  - [ ] Hotkey registration
  - [ ] Window management
- [ ] Stress testing:
  - [ ] Memory leak detection
  - [ ] Race condition tests
  - [ ] Hotkey reliability
  - [ ] Resource cleanup under load

## UI Improvements

### Completed ✓
- [x] Quick note window focus handling
- [x] Systray icon and menu
- [x] Main window visibility control
- [x] Basic task list view

### In Progress
- [ ] Task editing interface
- [ ] Drag and drop support
- [ ] Dark mode support
- [ ] Custom themes
- [ ] Keyboard shortcuts
- [ ] Task filtering UI
- [ ] Settings dialog

## Storage

### Completed ✓
- [x] SQLite implementation
- [x] Basic CRUD operations
- [x] Task storage schema
- [x] Storage interface

### In Progress
- [ ] Data migration system
- [ ] Backup/restore functionality
- [ ] Data export/import
- [ ] Cloud sync support
- [ ] Encryption support
- [ ] Storage interface cleanup
- [ ] Error handling improvements
- [ ] Transaction management enhancements

## Configuration

### Completed ✓
- [x] YAML configuration
- [x] Logger configuration
- [x] Hotkey configuration
- [x] UI configuration

### In Progress
- [ ] User preferences
- [ ] Theme configuration
- [ ] Sync settings
- [ ] Backup settings
- [ ] Plugin configuration

## Documentation

### Completed ✓
- [x] Dependency injection guide
- [x] Code style guidelines
- [x] Basic README
- [x] Installation guide

### In Progress
- [ ] API documentation
- [ ] User guide
- [ ] Developer guide
- [ ] Architecture overview
- [ ] Contributing guidelines
- [ ] Interface documentation
- [ ] Error handling guide
- [ ] Testing guide
- [ ] Service layer documentation

## Infrastructure

### Completed ✓
- [x] GitHub Actions CI
- [x] Basic build system
- [x] Development environment setup

### In Progress
- [ ] Release automation
- [ ] Cross-platform builds
- [ ] Code signing
- [ ] Auto-update system
- [ ] Telemetry system

## Future Considerations

### Features
- [ ] Mobile app support
- [ ] Web interface
- [ ] Calendar integration
- [ ] Email notifications
- [ ] Plugin system
- [ ] API endpoints
- [ ] Collaboration features
- [ ] Data analytics

### Technical
- [ ] GraphQL API
- [ ] Real-time sync
- [ ] Offline support
- [ ] Performance optimizations
- [ ] Security audit
- [ ] Accessibility improvements
- [ ] Localization support
