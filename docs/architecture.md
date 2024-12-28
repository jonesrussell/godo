# Godo Architecture

## Overview

Godo is a Windows-focused Todo application with quick-note capabilities, built in Go using the Fyne UI toolkit. The application is designed with a clean architecture that separates concerns and allows for future cross-platform support.

## Core Components

### Application Layer (`internal/app`)
- Main application lifecycle management
- Component coordination
- Event handling
- Global hotkey management

### Storage Layer (`internal/storage`)
- Task persistence
- SQLite implementation
- In-memory implementation for testing
- Repository pattern implementation

### GUI Layer (`internal/gui`)
- Main window management
- Quick note window
- System tray integration
- Event handling

### Configuration (`internal/config`)
- YAML-based configuration
- Environment-specific settings
- Runtime configuration management

### Logging (`internal/logger`)
- Structured logging with Zap
- Log level management
- Operation tracking

### Common (`internal/common`)
- Shared types and utilities
- Cross-cutting concerns
- Common interfaces

## Dependency Management

The application uses Wire for dependency injection, configured in:
- `internal/container/wire.go`
- `internal/container/wire_gen.go`

## Data Flow

1. User Interaction
   - GUI events
   - Global hotkeys
   - System tray actions

2. Application Logic
   - Event handling
   - Task management
   - State updates

3. Storage
   - Task persistence
   - Data retrieval
   - Transaction management

## Build System

- Task-based build automation
- Platform-specific build tags
- Docker support for Linux builds
- Windows-native compilation

## Testing Strategy

- Unit tests for core components
- Integration tests for storage
- GUI testing utilities
- Mock implementations for testing

## Future Architecture (HTTP API)

The upcoming HTTP API will add:

1. HTTP Server Layer
   - Chi router for endpoints
   - JSON response handling
   - Middleware pipeline

2. WebSocket Support
   - Real-time updates
   - Connection management
   - Event broadcasting

3. API Documentation
   - OpenAPI/Swagger specs
   - Usage examples
   - Integration guides 