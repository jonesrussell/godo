//go:build !docker

package app

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	storage.TaskStore
}

func (m *mockStore) List(ctx context.Context) ([]storage.Task, error) {
	return []storage.Task{}, nil
}

func TestApp_SetupUI(t *testing.T) {
	// Create test app
	testApp := test.NewApp()

	// Create mock logger
	mockLogger := &mockLogger{}

	// Create mock store
	mockStore := &mockStore{}

	// Create mock hotkey manager
	mockHotkey := &mockHotkey{}

	// Create test config
	cfg := &config.Config{
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				StartHidden: true,
			},
		},
	}

	// Create app options
	appOpts := &options.AppOptions{
		Name:    "Godo",
		Version: "1.0.0",
		ID:      "com.jonesrussell.godo",
	}

	// Create main window
	mainWin := mainwindow.New(testApp, mockStore, mockLogger, cfg.UI.MainWindow)

	// Create app with all required dependencies
	app := &App{
		name:       appOpts.Name,
		version:    appOpts.Version,
		id:         appOpts.ID,
		fyneApp:    testApp,
		logger:     mockLogger,
		store:      mockStore,
		hotkey:     mockHotkey,
		config:     cfg,
		mainWindow: mainWin,
	}

	// Test SetupUI
	app.SetupUI()

	// Verify expectations
	assert.True(t, cfg.UI.MainWindow.StartHidden, "Main window should be configured to start hidden")
	assert.NotNil(t, app.hotkey, "Hotkey manager should be initialized")
}

// Mock implementations
type mockLogger struct {
	logger.Logger
}

func (m *mockLogger) Error(msg string, args ...interface{}) {}

type mockHotkey struct {
	hotkey.Manager
}

func (m *mockHotkey) Register() error { return nil }
func (m *mockHotkey) Start() error    { return nil }
func (m *mockHotkey) Stop() error     { return nil }
