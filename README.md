# Godo

A minimalist todo application with quick-note hotkey support.

## Overview

Godo combines two powerful features:
1. A global hotkey that triggers a lightweight graphical popup for instantly capturing thoughts and tasks
2. A terminal-based (TUI) interface for detailed todo management

The quick-note feature uses a minimal graphical window that appears when you press the hotkey - type your note, hit enter, and it disappears. The main todo management interface uses a terminal UI for efficient keyboard-driven task organization.

## Features

- Instant note capture with global hotkey
  - Press hotkey → Graphical popup appears
  - Type note → Press enter → Window disappears
  - Zero-friction task capture with minimal visual interruption
- Terminal-based todo management interface
  - Organize and manage tasks
  - Mark tasks as complete
  - Delete tasks when done
- Runs as a system service
- SQLite database for reliable data storage
- Cross-platform compatibility
  - Windows: Native Win32 API for quick-note window
  - macOS: Cocoa/NSWindow for quick-note window
  - Linux: GTK for quick-note window
  - Terminal UI works consistently across all platforms

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
- Being cool 😎

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

- This project uses CGO dependencies (specifically `github.com/robotn/gohook`)
- Do not use `go mod vendor` as it may break CGO dependencies
- Always use `go mod tidy` to manage dependencies

## Installation

1. Clone the repository
```bash
git clone https://github.com/jonesrussell/godo.git
cd godo
```

2. Build the application
```bash
go build
```

3. Run the application
```bash
./godo
```

## Usage

[Add specific hotkey combinations and commands here]

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
