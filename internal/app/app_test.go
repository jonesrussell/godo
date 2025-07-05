//go:build !docker

package app_test

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/godo/internal/app"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
)

type mockStore struct {
	storage.TaskStore
	tasks []storage.Task
}

func (m *mockStore) List(_ context.Context) ([]storage.Task, error) {
	return m.tasks, nil
}

func (m *mockStore) Close() error {
	return nil
}

// Create a new mock logger for testing
func newMockLogger() *mockLogger {
	return &mockLogger{}
}

func TestApp_SetupUI(t *testing.T) {
	// Create test app
	testApp := test.NewApp()

	// Create mock logger
	mockLogger := newMockLogger()

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
		Core: &options.CoreOptions{
			Logger: mockLogger,
			Store:  mockStore,
			Config: cfg,
		},
		GUI: &options.GUIOptions{
			App:        testApp,
			MainWindow: mainwindow.New(testApp, mockStore, mockLogger, cfg.UI.MainWindow),
		},
	}

	// Create app with all required dependencies
	appInstance := app.New(&app.Params{
		Options: appOpts,
		Hotkey:  mockHotkey,
	})

	// Test SetupUI
	err := appInstance.SetupUI()
	require.NoError(t, err)

	// Verify expectations
	assert.True(t, cfg.UI.MainWindow.StartHidden, "Main window should be configured to start hidden")
	assert.NotNil(t, appInstance.Logger(), "Logger should be accessible")
}

func TestApp_Run(t *testing.T) {
	// Create test app
	testApp := test.NewApp()

	// Create mock logger
	mockLogger := newMockLogger()

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
		Core: &options.CoreOptions{
			Logger: mockLogger,
			Store:  mockStore,
			Config: cfg,
		},
		GUI: &options.GUIOptions{
			App:        testApp,
			MainWindow: mainwindow.New(testApp, mockStore, mockLogger, cfg.UI.MainWindow),
		},
	}

	// Create app with all required dependencies
	appInstance := app.New(&app.Params{
		Options: appOpts,
		Hotkey:  mockHotkey,
	})

	// Run app in a goroutine since it blocks
	go appInstance.Run()

	// Give it time to initialize
	testApp.Quit()

	// Verify the app initialized correctly
	assert.NotNil(t, appInstance.Logger(), "Logger should be accessible")
}

func TestApp_Cleanup(t *testing.T) {
	// Create test app
	testApp := test.NewApp()

	// Create mock logger
	mockLogger := newMockLogger()

	// Create mock store with cleanup tracking
	mockStore := &mockStoreWithCleanup{}

	// Create mock hotkey manager with cleanup tracking
	mockHotkey := &mockHotkeyWithCleanup{}

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
		Core: &options.CoreOptions{
			Logger: mockLogger,
			Store:  mockStore,
			Config: cfg,
		},
		GUI: &options.GUIOptions{
			App:        testApp,
			MainWindow: mainwindow.New(testApp, mockStore, mockLogger, cfg.UI.MainWindow),
		},
	}

	// Create app with all required dependencies
	appInstance := app.New(&app.Params{
		Options: appOpts,
		Hotkey:  mockHotkey,
	})

	// Test cleanup
	appInstance.Cleanup()

	// Verify cleanup was called
	assert.True(t, mockStore.cleanupCalled, "Store cleanup should be called")
	assert.True(t, mockHotkey.stopCalled, "Hotkey stop should be called")
}

// Additional mock implementations for cleanup testing
type mockStoreWithCleanup struct {
	mockStore
	cleanupCalled bool
}

func (m *mockStoreWithCleanup) Close() error {
	m.cleanupCalled = true
	return nil
}

type mockHotkeyWithCleanup struct {
	mockHotkey
	stopCalled bool
}

func (m *mockHotkeyWithCleanup) Stop() error {
	m.stopCalled = true
	return nil
}

type mockLogger struct {
	logger.Logger
}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}

type mockHotkey struct {
	hotkey.Manager
}

func (m *mockHotkey) Register() error { return nil }
func (m *mockHotkey) Start() error    { return nil }
func (m *mockHotkey) Stop() error     { return nil }
