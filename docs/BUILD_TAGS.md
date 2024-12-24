# Build Tags in Godo

## Platform-Specific Build Tags

### Operating Systems
```go
//go:build windows
// +build windows
```
Used for Windows-specific implementations.

```go
//go:build linux
// +build linux
```
Used for Linux-specific implementations.

## Usage

### Building for Different Platforms
```bash
# Windows build
go build -tags "windows,gl"

# Linux build
go build -tags "linux,gl"
```

### File Naming Conventions
- `*_windows.go` - Windows-specific implementations
- `*_linux.go` - Linux-specific implementations
- `*_common.go` - Shared interfaces and types

## Notes
- Always use both new (`//go:build`) and legacy (`// +build`) build tag formats for compatibility
- Platform-specific files should only contain platform-specific code
- Common interfaces should be in `*_common.go` files 