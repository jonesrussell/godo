# Godo

A cross-platform note-taking application built with Go and Fyne. The global hotkey feature enables instant note capture without interrupting your workflow. The management interface provides efficient note organization with a clean, modern design. The REST API allows for integration with other tools and services.

## Overview
Godo is a cross-platform note-taking application built with Go and Fyne. The global hotkey feature enables instant note capture without interrupting your workflow. The management interface provides efficient note organization with a clean, modern design. The REST API allows for integration with other tools and services.

## Features
- Zero-friction note capture with minimal visual interruption
- Global hotkey support
- SQLite storage
- Organize and manage notes
- Mark notes as complete
- Delete notes when done
- REST API
- Full CRUD operations for notes
- Cross-platform support

## Architecture
- Clean Architecture
- Dependency Injection
- Interface segregation (NoteReader, NoteWriter, NoteStore)
- Testable components
- Modular design

### Notes
- List all notes: `GET /api/v1/notes`
- Create note: `POST /api/v1/notes`
- Update note: `PUT /api/v1/notes/{id}`
- Delete note: `DELETE /api/v1/notes/{id}`

### Examples
```bash
# List notes
http :8080/api/v1/notes

# Create note
http POST :8080/api/v1/notes content="Buy groceries"

# Update note
http PUT :8080/api/v1/notes/{id} content="Updated content"

# Delete note
http DELETE :8080/api/v1/notes/{id}
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
