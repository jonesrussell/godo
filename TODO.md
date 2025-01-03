# TODO

## Core Features
- [x] Basic note management
- [x] SQLite storage
- [x] Global hotkey support
- [x] Quick note window
- [x] Main window
- [ ] Note categories/tags
- [ ] Due dates
- [ ] Note priorities
- [ ] Multiple note lists
- [ ] Recurring notes
- [ ] Note search and filtering

## Architecture
- [x] Clean architecture
- [x] Dependency injection
- [x] Interface segregation
- [x] Split Store interface (storage/store.go)
- [x] Split Transaction interface (storage/store.go)
- [x] Domain-driven design
- [x] Custom error types
- [ ] Improve logging with structured logging
- [ ] Add metrics collection
- [ ] Add distributed tracing
- [ ] Add service health checks

## Storage
- [x] SQLite implementation
- [x] Transaction support
- [x] Migration support
- [x] Note storage schema
- [x] Error handling with domain types
- [ ] Backup support
- [ ] Export/Import functionality
- [ ] Cloud sync support
- [ ] Data versioning
- [ ] Conflict resolution

## API
- [x] Basic REST endpoints
- [x] Health check
- [ ] OpenAPI/Swagger docs
- [ ] Authentication with JWT
- [ ] Rate limiting
- [ ] Metrics endpoint
- [ ] WebSocket support for real-time updates
- [ ] API versioning
- [ ] Request validation
- [ ] Response caching

## UI
- [x] Basic note list view
- [x] Quick note window
- [x] Main window
- [ ] Note editing interface
- [ ] Settings window
- [ ] System tray integration
- [ ] Note filtering UI
- [ ] Dark mode support
- [ ] Custom themes
- [ ] Keyboard shortcuts
- [ ] Accessibility support

## Testing
- [x] Unit tests
- [x] Integration tests
- [x] Store implementation tests
- [ ] End-to-end tests
- [ ] Performance tests
- [ ] Load tests
- [ ] UI tests
- [ ] API contract tests
- [ ] Property-based tests
- [ ] Fuzzing tests

## CI/CD
- [x] GitHub Actions
- [x] Cross-platform builds
- [ ] Automated releases
- [ ] Docker images
- [ ] Installation packages
- [ ] Automated dependency updates
- [ ] Security scanning
- [ ] Code coverage reports
- [ ] Performance regression tests

## Documentation
- [x] Architecture documentation
- [x] Interface guidelines
- [x] Error handling patterns
- [ ] API documentation
- [ ] User guide
- [ ] Developer guide
- [ ] Deployment guide
- [ ] Contributing guidelines
- [ ] Security policy
- [ ] Release notes template

## Code Quality
- [x] Custom linter rules
- [x] Interface segregation checks
- [x] Error handling checks
- [ ] Cyclomatic complexity limits
- [ ] Code duplication detection
- [ ] Dead code elimination
- [ ] Import organization
- [ ] Comment coverage
- [ ] API documentation coverage
- [ ] Test coverage thresholds
