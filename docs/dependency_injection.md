# Dependency Injection Architecture

## Overview
Godo uses Google's Wire framework for dependency injection. The architecture is designed to be modular, testable, and free of circular dependencies.

## Provider Sets

### Core Services
- **BaseSet**: Application metadata (name, version, ID)
- **LoggingSet**: Logging infrastructure
- **StorageSet**: Data persistence layer
- **ConfigSet**: Application configuration

### UI Layer
- **UISet**: User interface components
  - Fyne application instance
  - Main window
  - Quick note window
  - System tray

### Platform Features
- **HotkeySet**: Global hotkey management
  - Configuration (modifiers, key bindings)
  - Platform-specific manager implementation
- **HTTPSet**: HTTP server configuration

## Dependency Flow
1. Core services are initialized first
2. UI components are created using core services
3. Platform-specific features are initialized
4. All options are combined into AppOptions
5. HotkeyManager is created from HotkeyOptions
6. AppParams combines AppOptions and HotkeyManager
7. Main application is assembled

## Options Pattern
We use a layered options pattern to prevent circular dependencies:

### Configuration Options
Options that only contain configuration data, no dependencies:
```go
type LoggerOptions struct {
    Level       common.LogLevel
    Output      common.LogOutputPaths
    ErrorOutput common.ErrorOutputPaths
}

type HotkeyOptions struct {
    Modifiers common.ModifierKeys
    Key       common.KeyCode
}
```

### Component Options
Options that group related dependencies:
```go
type CoreOptions struct {
    Logger logger.Logger
    Store  storage.TaskStore
    Config *config.Config
}

type GUIOptions struct {
    App        fyne.App
    MainWindow gui.MainWindow
    QuickNote  gui.QuickNote
}
```

### Application Assembly
Two-step assembly to avoid circular dependencies:
```go
// Step 1: Combine all options
type AppOptions struct {
    Core    *CoreOptions
    GUI     *GUIOptions
    HTTP    *HTTPOptions
    Hotkey  *HotkeyOptions
    Name    common.AppName
    Version common.AppVersion
    ID      common.AppID
}

// Step 2: Combine options with instances
type AppParams struct {
    Options *AppOptions
    Hotkey  hotkey.Manager
}
```

## Best Practices
1. Keep provider sets small and focused
2. Use interfaces over concrete types
3. Separate configuration from instances
4. Use two-step assembly for circular dependencies
5. Document dependencies explicitly
6. Follow Wire naming conventions:
   - Provider functions: `Provide` prefix
   - Provider sets: `Set` suffix

## Testing
- **TestSet**: Provides mock implementations
- Each component has corresponding mock providers
- Test-specific options structs when needed

## Cleanup
All providers that allocate resources provide cleanup functions:
```go
func ProvideLogger(opts *LoggerOptions) (*logger.ZapLogger, func(), error)
func ProvideSQLiteStore(log logger.Logger) (*sqlite.Store, func(), error)
```

## Breaking Circular Dependencies
When faced with circular dependencies:
1. Separate configuration from instances
2. Use two-step assembly with intermediate types
3. Consider if the dependency is truly needed
4. Use interfaces to decouple components

Example:
```go
// Instead of:
type HotkeyOptions struct {
    Manager   hotkey.Manager  // Creates circular dependency
    Modifiers ModifierKeys
}

// Do:
type HotkeyOptions struct {
    Modifiers ModifierKeys    // Configuration only
}

// Then combine in AppParams:
type AppParams struct {
    Options *AppOptions
    Hotkey  hotkey.Manager    // Instance provided separately
}
``` 