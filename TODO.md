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

### In Progress
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Task priorities
- [ ] Multiple task lists
- [ ] Recurring tasks
- [ ] Task search and filtering

## Code Cleanup

### High Priority
- [ ] Task Creation Consolidation:
  - [ ] Create TaskService interface
  - [ ] Move task creation logic to service layer
  - [ ] Update GUI and API to use TaskService
  - [ ] Add validation at service layer
- [ ] Storage Interface Cleanup:
  - [ ] Remove deprecated Store interface
  - [ ] Consolidate TaskReader into TaskStore
  - [ ] Update all implementations
  - [ ] Add migration guide
- [ ] Error Handling Standardization:
  - [ ] Create common ErrorHandler interface
  - [ ] Implement consistent error types
  - [ ] Add error mapping layer
  - [ ] Update error documentation
- [ ] Documentation Alignment:
  - [ ] Update storage documentation
  - [ ] Standardize task model documentation
  - [ ] Add service layer documentation
  - [ ] Update API documentation

### Technical Debt
- [ ] Code Duplication Removal:
  - [ ] Consolidate window management code
  - [ ] Extract common UI patterns
  - [ ] Standardize test setup code
  - [ ] Create shared test utilities
- [ ] Interface Optimization:
  - [ ] Review and simplify interfaces
  - [ ] Remove unused methods
  - [ ] Add interface documentation
  - [ ] Update interface tests
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
