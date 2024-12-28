# Build System

This directory contains the build system for Godo, enabling cross-compilation for both Windows and Linux targets using Docker.

## Directory Structure

- `docker/` - Contains Docker-related files for building
  - `Dockerfile` - Main Dockerfile for cross-compilation environment
- `build.ps1` - PowerShell script to build both Windows and Linux versions

## Requirements

- Docker
- PowerShell
- Git (for development)

## Building

To build both Windows and Linux versions:

```powershell
.\build\build.ps1
```

This will:
1. Create a Docker image with all necessary build dependencies
2. Build the Linux version
3. Build the Windows version
4. Place the output in the `dist/` directory:
   - Windows: `dist/godo.exe`
   - Linux: `dist/godo`

## Build Environment

The build environment is based on the official Golang Docker image and includes:
- GCC and required development libraries
- X11 development files for Linux GUI support
- MinGW-w64 for Windows cross-compilation

## Notes

- The build process uses Docker to ensure consistent builds across different development machines
- All builds are statically linked with CGO enabled
- Build artifacts are placed in the `dist/` directory (git-ignored)
- Each build is tagged with the appropriate platform tag (`windows` or `linux`) 