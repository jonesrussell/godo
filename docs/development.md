# Development Guide

## Prerequisites

### Required Software
- Go 1.23 or higher
- MinGW-w64 GCC (Windows)
- Task (task runner)
- Git

### Optional Tools
- Docker (for Linux builds)
- VSCode with Go extension
- GNU diffutils (Windows)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/jonesrussell/godo.git
cd godo
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
task test
```

4. Build the application:
```bash
# Windows
task build:windows

# Linux
task build:linux
```

## Development Workflow

### Code Organization
- `cmd/` - Application entry points
- `internal/` - Internal packages
- `build/` - Build configurations
- `configs/` - Configuration files
- `docs/` - Documentation
- `scripts/` - Development scripts

### Making Changes
1. Create a new branch
2. Make changes
3. Run tests
4. Run linter
5. Submit PR

### Testing
- Run all tests: `task test`
- Run specific tests: `go test ./internal/...`
- Run with race detector: `go test -race ./...`

### Linting
- Run linter: `task lint`
- Auto-fix: `task lint:fix`

### Building
- Use Task commands for building
- Check `Taskfile.yaml` for available commands
- Use appropriate build tags

## Debugging

### VSCode
1. Open workspace
2. Use launch configurations in `.vscode/`
3. Set breakpoints
4. Start debugging

### Common Issues
- CGO requirements on Windows
- Build tag configuration
- Dependency issues

## Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## Code Style

### General Guidelines
- Follow standard Go conventions
- Use meaningful names
- Keep functions focused
- Add tests for new features

### Project-Specific
- Use dependency injection
- Follow repository pattern
- Use structured logging
- Handle errors appropriately

## Resources

### Documentation
- [Architecture Overview](architecture.md)
- [API Documentation](api/README.md)
- [Build System](build.md)

### External Links
- [Go Documentation](https://golang.org/doc/)
- [Fyne Toolkit](https://developer.fyne.io/)
- [Chi Router](https://go-chi.io/) 