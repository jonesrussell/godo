# Godo Architecture

## Overview

Godo is a cross-platform Todo application with quick-note capabilities and REST API support, built in Go using the Fyne UI toolkit. The application is designed with a clean architecture that separates concerns and supports multiple platforms.

## Core Components

### Application Layer (`internal/app`)
- Main application lifecycle management
- Component coordination
- Event handling
- Global hotkey management
- Platform-specific implementations

### Storage Layer (`internal/storage`)
- Task persistence
- SQLite implementation
- In-memory implementation for testing
- Repository pattern implementation

### GUI Layer (`internal/gui`)
- Main window management
- Quick note window
- System tray integration (Windows)
- Event handling
- Cross-platform UI components

### HTTP API Layer (`internal/api`)
- RESTful endpoints
- Chi router with middleware
- JSON response handling
- Health check endpoint
- Error handling
- Future WebSocket support

### Configuration (`internal/config`)
- YAML-based configuration
- Environment-specific settings
- Runtime configuration management
- API configuration

### Logging (`internal/logger`)
- Clean logger abstraction
- Multiple implementations:
  - Production logger (Zap-based)
  - Test logger for better test output
  - No-op logger for benchmarks
- Structured logging support
- Log level management
- Operation tracking
- Easy extensibility for new implementations

### Common (`internal/common`)
- Shared types and utilities
- Cross-cutting concerns
- Common interfaces
- Platform-specific utilities

## Dependency Management

The application uses Wire for dependency injection, configured in:
- `internal/container/wire.go`
- `internal/container/wire_gen.go`

## Data Flow

1. User Interaction
   - GUI events
   - Global hotkeys
   - System tray actions
   - HTTP API requests

2. Application Logic
   - Event handling
   - Task management
   - State updates
   - API request processing

3. Storage
   - Task persistence
   - Data retrieval
   - Transaction management

## Build System

- Task-based build automation
- Platform-specific build tags
- Docker support for Linux builds
- Windows-native compilation
- GitHub Actions CI/CD pipeline
- Automated releases
- Cross-platform binary distribution

## Testing Strategy

- Unit tests for core components
- Integration tests for storage and API
- GUI testing utilities
- Mock implementations for testing
- API endpoint testing
- Cross-platform testing

## API Architecture

The HTTP API includes:

1. HTTP Server Layer
   - Chi router for endpoints
   - JSON response handling
   - Middleware pipeline
   - Health checks
   - Error handling

2. Future Enhancements
   - WebSocket support for real-time updates
   - Connection management
   - Event broadcasting
   - Rate limiting
   - Authentication

3. API Documentation
   - OpenAPI/Swagger specs
   - Usage examples
   - Integration guides 