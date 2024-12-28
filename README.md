# Godo

A cross-platform todo application with quick-note hotkey support and REST API.

## Overview

Godo combines three powerful features:
1. A global hotkey that triggers a lightweight graphical popup for instantly capturing thoughts and tasks
2. A full-featured graphical interface for detailed todo management
3. A REST API for programmatic task management and integration

The quick-note feature uses a minimal graphical window that appears when you press the hotkey - type your note, hit enter, and it disappears. The main todo management interface provides efficient task organization with a clean, modern design. The REST API allows for integration with other tools and services.

## Features

- Cross-platform support
  - Windows: Full support with system tray integration
  - Linux: Full support (except system tray)
  - macOS: Coming soon
- Instant note capture with global hotkey
  - Press hotkey → Graphical popup appears
  - Type note → Press enter → Window disappears
  - Zero-friction task capture with minimal visual interruption
- Graphical todo management interface
  - Organize and manage tasks
  - Mark tasks as complete
  - Delete tasks when done
- REST API
  - Full CRUD operations for tasks
  - JSON responses
  - Health check endpoint
  - Proper error handling
- Robust storage system
  - SQLite-based persistence with comprehensive validation
  - Prevents data inconsistencies (empty IDs, duplicates)
  - Connection state management
  - Automated migrations
  - High test coverage (66%+)
- Robust logging system
  - Structured logging with multiple implementations
  - Test-friendly logging for better test output
  - Comprehensive operation tracking
  - Clean abstraction for easy customization
- Automated builds and releases
  - GitHub Actions CI/CD pipeline
  - Cross-platform binary releases
  - Docker support for development and testing
- Built with Fyne toolkit for native look and feel

## API Endpoints

All endpoints return JSON responses. The base URL is `http://localhost:8080`.

### Health Check
```
GET /health
Response: {"status": "ok"}
```

### Tasks
- List all tasks: `GET /api/v1/tasks`
- Create task: `POST /api/v1/tasks`
- Update task: `PUT /api/v1/tasks/{id}`
- Delete task: `DELETE /api/v1/tasks/{id}`

Example using HTTPie:
```bash
# List tasks
http :8080/api/v1/tasks

# Create task
http POST :8080/api/v1/tasks title="Buy groceries" description="Milk, bread, eggs"

# Update task
http PUT :8080/api/v1/tasks/{id} title="Updated title" description="New description"

# Delete task
http DELETE :8080/api/v1/tasks/{id}
```

## Prerequisites

- Go 1.23 or higher
- SQLite3
- MinGW-w64 GCC (for Windows users)
  - Recommended version: [MinGW-w64 GCC 14.2.0 or later](https://github.com/niXman/mingw-builds-binaries/releases)
  - Choose the appropriate version based on your system architecture (x86_64 or i686)
  - Installation steps:
    1. Download the appropriate version (e.g., x86_64-14.2.0-release-posix-seh-ucrt-rt_v12-rev0.7z)
    2. Extract to C:\mingw64 (or your preferred location)
    3. Add C:\mingw64\bin to your system's PATH environment variable
    4. Verify installation by running `gcc --version` in Command Prompt
- Task (task runner)
  ```powershell
  # Install using Chocolatey
  choco install go-task
  ```

### Additional Development Prerequisites

For Windows developers:
- GNU diffutils (required for code linting)
  ```powershell
  # Install using Chocolatey
  choco install diffutils
  
  # Add to PATH (in PowerShell)
  $env:PATH += ";C:\ProgramData\chocolatey\lib\diffutils\tools\bin"
  refreshenv
  ```

## Development Notes

- This project uses CGO dependencies (specifically `golang.design/x/hotkey`)
- Do not use `go mod vendor` as it may break CGO dependencies
- Always use `go mod tidy` to manage dependencies

## Building

1. Clone the repository
```bash
git clone https://github.com/jonesrussell/godo.git
cd godo
```

2. Build the application
```powershell
# For Windows
task build:windows

# For Linux
task build:linux

# For Linux using Docker
task build:linux:docker
```

3. Run tests
```powershell
# Run tests for current platform
task test

# Run tests in Docker
task test:linux:docker
```

4. Run linting
```powershell
# Run linting for current platform
task lint

# Run linting in Docker
task lint:linux:docker
```

## Usage

Default hotkey: Ctrl+Shift+N

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to all contributors
- Inspired by the need for a simple, efficient todo system
