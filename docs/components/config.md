# Configuration System

## Overview

The configuration system uses YAML files for application settings, supporting different environments and runtime configuration updates.

## Configuration Files

### Default Configuration (`configs/default.yaml`)
```yaml
app:
  name: "Godo"
  version: "0.1.0"
  id: "io.github.jonesrussell.godo"

database:
  path: "godo.db"

logger:
  level: "info"
  console: true

hotkeys:
  quick_note: "Ctrl+Alt+G"

ui:
  main_window:
    width: 800
    height: 600
    start_hidden: true
  quick_note:
    width: 400
    height: 200
```

### Test Configuration (`configs/test.yaml`)
```yaml
database:
  path: ":memory:"

logger:
  level: "debug"
  console: true

ui:
  main_window:
    start_hidden: true
```

## Configuration Types

### App Configuration
```go
type AppConfig struct {
    Name    string
    Version string
    ID      string
}
```

### Database Configuration
```go
type DBConfig struct {
    Path string
}
```

### Logger Configuration
```go
type LogConfig struct {
    Level   string
    Console bool
    File    bool
    Path    string
}
```

### UI Configuration
```go
type UIConfig struct {
    MainWindow  WindowConfig
    QuickNote   WindowConfig
}

type WindowConfig struct {
    Width       int
    Height      int
    StartHidden bool
}
```

## Usage

### Loading Configuration
```go
config, err := config.Load("configs/default.yaml")
if err != nil {
    log.Fatal(err)
}
```

### Accessing Configuration
```go
dbPath := config.Database.Path
logLevel := config.Logger.Level
quickNoteHotkey := config.Hotkeys.QuickNote
```

## Environment Support

### Development
- Default configuration
- Debug logging
- Local database

### Testing
- In-memory database
- Debug logging
- Minimal UI

### Production
- Production logging
- Persistent storage
- Full UI features

## Best Practices

1. Configuration Management
   - Use environment-specific files
   - Don't commit sensitive data
   - Use reasonable defaults
   - Validate configuration values

2. Security
   - Protect sensitive values
   - Use environment variables for secrets
   - Validate file permissions

3. Validation
   - Check required fields
   - Validate value ranges
   - Handle missing values

## Testing

### Configuration Tests
```go
func TestConfig(t *testing.T) {
    config, err := config.Load("testdata/config.yaml")
    require.NoError(t, err)
    
    assert.Equal(t, "Godo", config.App.Name)
    assert.Equal(t, ":memory:", config.Database.Path)
}
```

### Validation Tests
```go
func TestConfigValidation(t *testing.T) {
    config, err := config.Load("testdata/invalid.yaml")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "required field")
}
```

## Common Patterns

### Default Values
```go
func (c *Config) setDefaults() {
    if c.UI.MainWindow.Width == 0 {
        c.UI.MainWindow.Width = 800
    }
    if c.Logger.Level == "" {
        c.Logger.Level = "info"
    }
}
```

### Environment Variables
```go
func (c *Config) loadEnv() {
    if path := os.Getenv("GODO_DB_PATH"); path != "" {
        c.Database.Path = path
    }
}
```

### Validation
```go
func (c *Config) validate() error {
    if c.App.Name == "" {
        return errors.New("app name is required")
    }
    return nil
}
```

## Resources

- [YAML Package](https://pkg.go.dev/gopkg.in/yaml.v3)
- [Viper](https://github.com/spf13/viper) (Future consideration)
- [Configuration Best Practices](https://12factor.net/config) 