# Godo

A minimalist todo application with global hotkey support.

## Overview

Godo is a lightweight todo application that runs as a system service, allowing quick access to your tasks through convenient hotkeys. Built with Go and SQLite, it provides a seamless task management experience.

## Features

- Runs as a system service
- Global hotkey support for quick access
- SQLite database for reliable data storage
- Minimalist interface
- Cross-platform compatibility

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
