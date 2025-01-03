// Package options defines the dependency injection options for the application
package options

import (
	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// LoggerOptions defines options for logger configuration
type LoggerOptions struct {
	Level       common.LogLevel
	Output      common.LogOutputPaths
	ErrorOutput common.ErrorOutputPaths
}

// HotkeyOptions defines options for hotkey configuration
type HotkeyOptions struct {
	Modifiers common.ModifierKeys
	Key       common.KeyCode
}

// CoreOptions groups core application dependencies
type CoreOptions struct {
	Logger logger.Logger
	Store  types.Store
	Config *config.Config
}

// GUIOptions groups GUI dependencies
type GUIOptions struct {
	App        fyne.App
	MainWindow gui.MainWindowManager
	QuickNote  gui.QuickNoteManager
}

// AppOptions groups all application dependencies
type AppOptions struct {
	Core    *CoreOptions
	GUI     *GUIOptions
	Name    common.AppName
	Version common.AppVersion
	ID      common.AppID
}
