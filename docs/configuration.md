# Configuration Guide

## Overview

Godo uses YAML configuration files with environment-specific overrides and runtime configuration options.

## Configuration Files

### Location
- `configs/default.yaml`: Default configuration
- `configs/test.yaml`: Test environment configuration
- Environment-specific files (e.g., `configs/production.yaml`)

### File Structure
```yaml
# Application configuration
app:
  name: "Godo"
  version: "0.1.0"
  id: "io.github.jonesrussell.godo"

# Database settings
database:
  path: "godo.db"

# Logging configuration
logger:
  level: "info"    # debug, info, warn, error
  console: true    # Enable console logging
  file: true       # Enable file logging
  path: "logs/"    # Log file directory

# Hotkey configuration
hotkeys:
  quick_note: "Ctrl+Alt+G"

# UI settings
ui:
  main_window:
    width: 800
    height: 600
    start_hidden: true
  quick_note:
    width: 400
    height: 200
```

## Environment Variables

### Available Overrides
```bash
GODO_CONFIG_PATH    # Configuration file path
GODO_DB_PATH       # Database file location
GODO_LOG_LEVEL     # Logging level
GODO_LOG_PATH      # Log file directory
```

### Usage Example
```bash
export GODO_LOG_LEVEL=debug
export GODO_DB_PATH=/custom/path/godo.db
./godo
```

## Configuration Loading

### Priority Order
1. Environment variables
2. Command-line flags
3. Environment-specific config file
4. Default config file

### Loading Process
```go
config, err := config.Load("configs/default.yaml")
if err != nil {
    log.Fatal(err)
}
```

## Configuration Types

### Application
```go
type AppConfig struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`
    ID      string `yaml:"id"`
}
```

### Database
```go
type DBConfig struct {
    Path string `yaml:"path"`
}
```

### Logging
```go
type LogConfig struct {
    Level   string `yaml:"level"`
    Console bool   `yaml:"console"`
    File    bool   `yaml:"file"`
    Path    string `yaml:"path"`
}
```

### UI
```go
type UIConfig struct {
    MainWindow  WindowConfig `yaml:"main_window"`
    QuickNote   WindowConfig `yaml:"quick_note"`
}

type WindowConfig struct {
    Width       int  `yaml:"width"`
    Height      int  `yaml:"height"`
    StartHidden bool `yaml:"start_hidden"`
}
```

## Environment-Specific Configurations

### Development
```yaml
logger:
  level: "debug"
  console: true
  file: true

database:
  path: "godo.db"
```

### Testing
```yaml
database:
  path: ":memory:"

logger:
  level: "debug"
  console: true
  file: false
```

### Production
```yaml
logger:
  level: "info"
  console: false
  file: true
  path: "/var/log/godo/"

database:
  path: "/var/lib/godo/godo.db"
```

## Configuration Validation

### Required Fields
```go
func (c *Config) Validate() error {
    if c.App.Name == "" {
        return errors.New("app name is required")
    }
    return nil
}
```

### Value Validation
```go
func (c *Config) ValidateLogLevel() error {
    validLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    
    if !validLevels[c.Logger.Level] {
        return fmt.Errorf("invalid log level: %s", c.Logger.Level)
    }
    return nil
}
```

## Best Practices

1. Configuration Management
   - Use environment variables for sensitive data
   - Keep configurations versioned
   - Document all options
   - Provide sensible defaults

2. Security
   - Don't commit sensitive data
   - Use appropriate file permissions
   - Validate all inputs
   - Handle missing values gracefully

3. Maintainability
   - Keep configurations organized
   - Use consistent naming
   - Document changes
   - Version control configurations

## Common Issues

### File Permissions
```bash
# Set appropriate permissions
chmod 600 configs/production.yaml
```

### Missing Configuration
```go
// Set defaults if config is missing
if config.UI.MainWindow.Width == 0 {
    config.UI.MainWindow.Width = 800
}
```

### Environment Variables
```go
// Override with environment variables
if envPath := os.Getenv("GODO_DB_PATH"); envPath != "" {
    config.Database.Path = envPath
}
```

## Resources

- [YAML Specification](https://yaml.org/spec/)
- [Go YAML Package](https://pkg.go.dev/gopkg.in/yaml.v3)
- [12-Factor App Config](https://12factor.net/config)
- [Go Configuration Patterns](https://dave.cheney.net/2014/10/22/configuration-in-go) 