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

### Testing
- **TestSet**: Mock implementations for testing
  - Mock store
  - Mock windows
  - Mock hotkey manager
  - Mock Fyne app
  - Test configuration

## Dependency Flow
1. Core services are initialized first (logging, storage, config)
2. UI components are created using core services
3. Platform-specific features are initialized
4. Options are created and validated
5. Main application is assembled

## Options Pattern
We use a layered options pattern to prevent circular dependencies:

### Configuration Options
Pure configuration without dependencies:
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

type HTTPOptions struct {
    Config *common.HTTPConfig
}
```

### Component Options
Group related dependencies:
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
```go
// All application options
type AppOptions struct {
    Core    *CoreOptions
    GUI     *GUIOptions
    HTTP    *HTTPOptions
    Hotkey  *HotkeyOptions
    Name    common.AppName
    Version common.AppVersion
    ID      common.AppID
}

// Final application parameters
type Params struct {
    Options *AppOptions
    Hotkey  hotkey.Manager
}
```

## Testing Strategy

### Mock Providers
Each component has a corresponding mock provider:
```go
func ProvideMockStore() storage.TaskStore
func ProvideMockMainWindow() *gui.MockMainWindow
func ProvideMockQuickNote() *gui.MockQuickNote
func ProvideMockHotkey() *apphotkey.MockManager
func ProvideMockFyneApp() fyne.App
```

### Test App Assembly
The test app is assembled using the TestSet:
```go
func ProvideTestAppParams(
    logger logger.Logger,
    store storage.TaskStore,
    mainWindow gui.MainWindow,
    quickNote gui.QuickNote,
    hotkey apphotkey.Manager,
    httpConfig *common.HTTPConfig,
    name common.AppName,
    version common.AppVersion,
    id common.AppID,
) *app.TestApp
```

## Best Practices
1. Keep provider sets small and focused
2. Use interfaces over concrete types
3. Separate configuration from instances
4. Use layered options pattern to prevent cycles
5. Document dependencies explicitly
6. Follow Wire naming conventions:
   - Provider functions: `Provide` prefix
   - Provider sets: `Set` suffix
7. Test with mock implementations
8. Use cleanup functions for resource management

## Resource Cleanup
Providers that allocate resources must provide cleanup functions:
```go
func ProvideLogger(opts *LoggerOptions) (*logger.ZapLogger, func(), error)
func ProvideSQLiteStore(log logger.Logger) (*sqlite.Store, func(), error)
```

## Breaking Circular Dependencies
1. Separate configuration from instances
2. Use layered options pattern
3. Group related dependencies
4. Use interfaces to decouple components
5. Consider dependency direction (core → ui → platform) 