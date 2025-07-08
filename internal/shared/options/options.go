// Package options defines the dependency injection options for the application
package options

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/infrastructure/gui"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	"github.com/jonesrussell/godo/internal/infrastructure/storage"
	"github.com/jonesrussell/godo/internal/shared/common"
	"github.com/jonesrussell/godo/internal/shared/config"
)

// LoggerOptions groups logger configuration options
type LoggerOptions struct {
	Level       common.LogLevel
	Output      common.LogOutputPaths
	ErrorOutput common.ErrorOutputPaths
}

// CoreOptions groups core application dependencies
type CoreOptions struct {
	Logger logger.Logger
	Store  storage.TaskStore
	Config *config.Config
}

// GUIOptions groups GUI dependencies
type GUIOptions struct {
	App        fyne.App
	MainWindow gui.MainWindow
	QuickNote  gui.QuickNote
}

// HotkeyOptions groups hotkey configuration
type HotkeyOptions struct {
	Modifiers common.ModifierKeys
	Key       common.KeyCode
}

// HTTPOptions groups HTTP server configuration
type HTTPOptions struct {
	Config *common.HTTPConfig
}

// AppOptions groups all application dependencies
type AppOptions struct {
	Core    *CoreOptions
	GUI     *GUIOptions
	HTTP    *HTTPOptions
	Hotkey  *HotkeyOptions
	Name    common.AppName
	Version common.AppVersion
	ID      common.AppID
}
