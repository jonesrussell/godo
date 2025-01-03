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
- [ ] Split Store interface (storage/store.go)
- [ ] Split Transaction interface (storage/store.go)
- [ ] Improve error handling
- [ ] Split NoteHandler interface (api/handler.go)
- [ ] Improve logging
- [ ] Add metrics
- [ ] Add tracing

## Storage
- [x] SQLite implementation
- [x] Transaction support
- [x] Migration support
- [x] Note storage schema
- [ ] Backup support
- [ ] Export/Import
- [ ] Cloud sync

## API
- [x] Basic REST endpoints
- [x] Health check
- [ ] OpenAPI/Swagger docs
- [ ] Authentication
- [ ] Rate limiting
- [ ] Metrics endpoint
- [ ] WebSocket support

## UI
- [x] Basic note list view
- [x] Quick note window
- [x] Main window
- [ ] Note editing interface
- [ ] Settings window
- [ ] System tray
- [ ] Note filtering UI
- [ ] Dark mode support
- [ ] Custom themes
- [ ] Keyboard shortcuts

## Testing
- [x] Unit tests
- [x] Integration tests
- [ ] End-to-end tests
- [ ] Performance tests
- [ ] Load tests
- [ ] UI tests

## CI/CD
- [x] GitHub Actions
- [x] Cross-platform builds
- [ ] Automated releases
- [ ] Docker images
- [ ] Installation packages
