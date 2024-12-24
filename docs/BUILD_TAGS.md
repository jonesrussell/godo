# Build Tags in Godo

This document explains the build tags used in the Godo project, following Fyne's conventions.

## Platform Tags

### Desktop (Windows/Unix)
```go
//go:build !ci && !android && !ios && !wasm && !test_web_driver
```
Used for Windows-specific code.

```go
//go:build !windows && !android && !ios && !wasm && !js
```
Used for Unix-specific code (Linux, BSD, macOS).

### Docker
```go
//go:build docker
```
Used for Docker-specific implementations.

```go
//go:build !docker
```
Used for non-Docker implementations.

## Test Tags
```go
//go:build !docker && wireinject
```
Used for wire injection in tests.

```go
//go:build !docker
```
Used for regular tests.

## Build Environments

### CI Environment
- `!ci` - Used to exclude code from CI builds
- `ci` - Used for CI-specific implementations

### Development
- `!test_web_driver` - Excludes web driver test code
- `!mobile` - Excludes mobile platform code

## Usage Examples

1. Windows Desktop App:
   ```go
   //go:build !ci && !android && !ios && !wasm && !test_web_driver
   ```

2. Unix Desktop App:
   ```go
   //go:build !windows && !android && !ios && !wasm && !js
   ```

3. Docker Environment:
   ```go
   //go:build docker
   ```

## Build Commands

### Windows Build
```bash
task build:windows
```
Uses: `!ci,!android,!ios,!wasm,!test_web_driver`

### Linux Build
```bash
task build:linux
```
Uses: `!windows,!android,!ios,!wasm,!js`

### Docker Build
```bash
task build:docker
```
Uses: `docker` 