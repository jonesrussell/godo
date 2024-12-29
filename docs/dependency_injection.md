# Dependency Injection Guide

## Overview
Godo uses Google's Wire framework for dependency injection. This document outlines our DI patterns and best practices.

## Core Concepts

### Provider Sets
We organize providers into focused sets:
```go
var CoreSet = wire.NewSet(
    ProvideCoreOptions,
    ProvideLogger,
    ProvideSQLiteStore,
)

var GUISet = wire.NewSet(
    ProvideGUIOptions,
    ProvideFyneApp,
    ProvideMainWindow,
    ProvideQuickNote,
)

var HotkeySet = wire.NewSet(
    ProvideHotkeyOptions,
    ProvideHotkeyManager,
)
```

### Options Pattern
We use options structs to group related dependencies:
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

### Provider Functions
Clean provider functions with proper error handling:
```go
func ProvideLogger(opts *options.LoggerOptions) (*logger.ZapLogger, func(), error) {
    cfg := &logger.Config{
        Level:       string(opts.Level),
        Development: true,
        Encoding:    "console",
    }
    return logger.NewLogger(cfg)
}
```

## Best Practices

### 1. Clean Injector Functions
- Only contain wire.Build call
- No additional logic
```go
func InitializeApp() (app.Application, func(), error) {
    wire.Build(AppSet)
    return nil, nil, nil
}
```

### 2. Resource Cleanup
- Return cleanup functions
- Clean up in reverse order
```go
func ProvideSQLiteStore(log logger.Logger) (*sqlite.Store, func(), error) {
    store, err := sqlite.New("godo.db", log)
    if err != nil {
        return nil, nil, err
    }
    cleanup := func() {
        store.Close()
    }
    return store, cleanup, nil
}
```

### 3. Error Handling
- Return errors from providers
- Validate inputs
- Check dependencies
```go
func ProvideHotkeyManager(opts *options.HotkeyOptions) (*hotkey.Manager, error) {
    if len(opts.Modifiers) == 0 {
        return nil, errors.New("no modifiers specified")
    }
    return hotkey.NewManager(opts.Modifiers, opts.Key)
}
```

### 4. Interface Bindings
- Bind concrete types to interfaces
- Use wire.Bind in provider sets
```go
var StorageSet = wire.NewSet(
    ProvideSQLiteStore,
    wire.Bind(new(storage.TaskStore), new(*sqlite.Store)),
)
```

### 5. Testing Support
- Create mock providers
- Use test-specific provider sets
```go
var TestSet = wire.NewSet(
    ProvideMockStore,
    ProvideMockLogger,
    wire.Bind(new(storage.TaskStore), new(*mock.Store)),
)
```

## Common Patterns

### Two-Step Initialization
```go
type App struct {
    name    string
    logger  logger.Logger
    store   storage.TaskStore
    hotkey  hotkey.Manager
}

func New(params *Params) *App {
    app := &App{
        name:   params.Name,
        logger: params.Logger,
        store:  params.Store,
        hotkey: params.Hotkey,
    }
    return app
}

func (a *App) Initialize() error {
    if err := a.store.Initialize(); err != nil {
        return err
    }
    if err := a.hotkey.Register(); err != nil {
        return err
    }
    return nil
}
```

### Platform-Specific Code
```go
//go:build windows
package hotkey

func ProvideHotkeyManager() (hotkey.Manager, error) {
    // Windows-specific implementation
}
```

### Configuration Management
```go
func ProvideConfig() (*config.Config, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }
    return cfg, nil
}
```

## Troubleshooting

### Common Issues
1. Circular Dependencies
   - Use options pattern
   - Split interfaces
   - Use two-step initialization

2. Missing Providers
   - Check provider sets
   - Verify interface bindings
   - Check build tags

3. Multiple Providers
   - Use distinct types
   - Use wire.Value for constants
   - Use wire.InterfaceValue for interfaces

### Debugging Tips
1. Check wire_gen.go
2. Use wire check command
3. Review provider graph
4. Verify cleanup order

## Testing

### Mock Providers
```go
func ProvideMockStore() storage.TaskStore {
    return &mock.Store{}
}
```

### Test Initialization
```go
func TestApp(t *testing.T) {
    app, cleanup, err := InitializeTestApp()
    require.NoError(t, err)
    defer cleanup()
    
    // Test app functionality
}
```

## Migration Guide

### From Global State
1. Identify global variables
2. Create provider functions
3. Add to provider sets
4. Update consumers

### From Factory Functions
1. Convert to providers
2. Add cleanup functions
3. Update call sites
4. Add error handling 