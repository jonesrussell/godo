# Build System

## Overview

The build system uses Task for automation and supports multiple platforms with different build configurations.

## Prerequisites

### Windows
- Go 1.23+
- MinGW-w64 GCC
- Task runner
- GNU diffutils (optional)

### Linux
- Go 1.23+
- GCC
- Task runner
- Docker (optional)

## Build Commands

### Basic Build
```bash
# Windows build
task build:windows

# Linux build
task build:linux

# Linux build using Docker
task build:linux:docker
```

### Testing
```bash
# Run tests
task test

# Run tests in Docker
task test:linux:docker

# Run linting
task lint
```

## Build Configuration

### Task File (`Taskfile.yaml`)
```yaml
version: '3'

vars:
  BINARY_NAME: godo
  VERSION: 0.1.0

tasks:
  build:windows:
    cmds:
      - go build -tags windows -o bin/{{.BINARY_NAME}}.exe ./cmd/godo

  build:linux:
    cmds:
      - go build -tags linux -o bin/{{.BINARY_NAME}} ./cmd/godo

  test:
    cmds:
      - go test -v ./...
```

### Build Tags

#### Windows Build
```go
//go:build windows
// +build windows
```

#### Linux Build
```go
//go:build linux
// +build linux
```

## Docker Support

### Dockerfile (`build/Dockerfile.linux`)
```dockerfile
# Build stage
FROM golang:1.23-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -tags linux -o bin/godo cmd/godo/main.go

# Runtime stage
FROM debian:bullseye-slim
COPY --from=builder /app/bin/godo /app/godo
ENTRYPOINT ["/app/godo"]
```

### Docker Commands
```bash
# Build Docker image
docker build -f build/Dockerfile.linux -t godo .

# Run in Docker
docker run godo
```

## Directory Structure

```
.
├── bin/                    # Build outputs
├── build/                  # Build configurations
│   ├── Dockerfile.linux
│   └── package/
├── cmd/                    # Entry points
│   └── godo/
├── configs/               # Configuration files
├── internal/              # Internal packages
└── scripts/              # Build scripts
```

## Build Artifacts

### Windows
- `bin/godo.exe`
- Debug symbols
- Resource files

### Linux
- `bin/godo`
- Debug symbols

## Continuous Integration

### GitHub Actions
```yaml
name: Build and Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [windows-latest, ubuntu-latest]
```

## Development Workflow

1. Local Development
   ```bash
   # Build and run
   task build:windows
   ./bin/godo.exe
   ```

2. Testing Changes
   ```bash
   # Run tests
   task test
   
   # Run linting
   task lint
   ```

3. Creating Release
   ```bash
   # Build all platforms
   task build:all
   
   # Run tests
   task test:all
   ```

## Best Practices

1. Build Process
   - Use build tags appropriately
   - Keep build scripts maintainable
   - Document build requirements
   - Version build artifacts

2. Dependencies
   - Vendor dependencies when needed
   - Use go.mod for versioning
   - Document external requirements
   - Handle CGO dependencies

3. Cross-Compilation
   - Test on target platforms
   - Handle platform-specific code
   - Use appropriate build constraints
   - Check CGO requirements

## Common Issues

### Windows
- CGO compilation errors
- MinGW-w64 configuration
- Path issues
- Resource compilation

### Linux
- Library dependencies
- Permission issues
- Path configuration
- Docker build context

## Resources

- [Task Documentation](https://taskfile.dev)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Docker Documentation](https://docs.docker.com)
- [GitHub Actions](https://docs.github.com/en/actions) 