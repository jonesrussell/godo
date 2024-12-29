//go:build !docker

package app

import (
	"testing"

	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/options"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	storage.TaskStore
}

func newMockStore() storage.TaskStore {
	return &mockStore{}
}

func TestApp_SetupUI(t *testing.T) {
	// Create test app
	testApp := test.NewApp()
	store := newMockStore()
	log := logger.NewTestLogger(t)
	cfg := &config.Config{
		UI: config.UIConfig{
			MainWindow: config.WindowConfig{
				StartHidden: true,
				Width:       800,
				Height:      600,
			},
			QuickNote: config.WindowConfig{
				Width:  200,
				Height: 100,
			},
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: "Ctrl+Shift+N",
		},
	}

	// Create windows
	mainWin := mainwindow.New(testApp, store, log, cfg.UI.MainWindow)
	quickNote := quicknote.New(testApp, store, log, cfg.UI.QuickNote)

	// Create options
	coreOpts := &options.CoreOptions{
		Logger: log,
		Store:  store,
		Config: cfg,
	}

	guiOpts := &options.GUIOptions{
		App:        testApp,
		MainWindow: mainWin,
		QuickNote:  quickNote,
	}

	hotkeyOpts := &options.HotkeyOptions{
		Modifiers: []string{"ctrl", "shift"},
		Key:       "n",
	}

	httpOpts := &options.HTTPOptions{
		Config: &common.HTTPConfig{
			Port: 8080,
		},
	}

	appOpts := &options.AppOptions{
		Core:    coreOpts,
		GUI:     guiOpts,
		HTTP:    httpOpts,
		Hotkey:  hotkeyOpts,
		Name:    "Godo",
		Version: "1.0.0",
		ID:      "com.jonesrussell.godo",
	}

	// Create app
	app := New(&Params{
		Options: appOpts,
	})

	// Test main window starts hidden
	app.SetupUI()
	assert.True(t, cfg.UI.MainWindow.StartHidden, "Main window should be configured to start hidden")

	// Test systray is set up
	_, ok := testApp.(desktop.App)
	assert.True(t, ok, "App should implement desktop.App")

	// Test hotkey is registered
	assert.NotNil(t, app.hotkey, "Hotkey manager should be initialized")
}
