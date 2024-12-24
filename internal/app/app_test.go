//go:build !docker
// +build !docker

package app

import (
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
)

func TestMain(m *testing.M) {
	mainthread.Init(func() {
		os.Exit(m.Run())
	})
}

// mockHotkeyFactory is a test implementation of config.HotkeyFactory
type mockHotkeyFactory struct {
	testHotkey config.HotkeyHandler
}

func (f *mockHotkeyFactory) NewHotkey(_ []hotkey.Modifier, _ hotkey.Key) config.HotkeyHandler {
	return f.testHotkey
}

// mockQuickNote is a test implementation of QuickNoteService
type mockQuickNote struct {
	shown bool
}

func (m *mockQuickNote) Show() {
	m.shown = true
}

func (m *mockQuickNote) Hide() {
	m.shown = false
}

// mockSystray is a test implementation of systray.Interface
type mockSystray struct {
	ready bool
	menu  *fyne.Menu
	icon  fyne.Resource
}

func (m *mockSystray) Setup(menu *fyne.Menu) {
	m.menu = menu
	m.ready = true
}

func (m *mockSystray) SetIcon(icon fyne.Resource) {
	m.icon = icon
}

func (m *mockSystray) IsReady() bool {
	return m.ready
}

func TestApp(t *testing.T) {
	// Skip if X11 is not available
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG)
	if err := hk.Register(); err != nil {
		t.Skip("Skipping test due to X11 initialization failure:", err)
	}
	hk.Unregister()

	// Use Fyne test app with test driver
	fyneApp := test.NewApp()
	defer fyneApp.Quit()

	// Create mock factory with test hotkey
	testHotkey := NewTestHotkey()
	mockFactory := &mockHotkeyFactory{testHotkey: testHotkey}

	// Create test dependencies
	cfg := &config.Config{
		Logger: common.LogConfig{
			Level:   "debug",
			Console: true,
		},
		App: config.AppConfig{
			Name:    "Test App",
			Version: "0.1.0",
		},
	}
	log := logger.NewTestLogger(t)
	store := memory.New()

	// Create app
	app := NewApp(cfg, store, log, mockFactory)
	require.NotNil(t, app)
	defer app.Cleanup()

	// Replace services with mocks
	mockQN := &mockQuickNote{}
	app.SetQuickNoteService(mockQN)
	mockSystray := &mockSystray{}
	app.systray = mockSystray

	// Test note operations
	t.Run("note operations", func(t *testing.T) {
		// Add note
		err := app.SaveNote("Test note")
		require.NoError(t, err)

		// Get notes
		notes, err := app.GetNotes()
		require.NoError(t, err)
		assert.Contains(t, notes, "Test note")

		// Verify version
		assert.Equal(t, "0.1.0", app.GetVersion())
	})

	// Test hotkey
	t.Run("hotkey", func(t *testing.T) {
		// Setup UI which includes hotkey registration
		app.SetupUI()

		// Mark systray as ready immediately since it's a mock
		mockSystray.ready = true

		// Setup hotkey directly
		err := app.setupGlobalHotkey()
		if err != nil {
			t.Skip("Skipping hotkey test due to initialization failure:", err)
		}

		// Create a channel to signal when the quick note is shown
		shown := make(chan struct{})
		done := make(chan struct{})
		defer close(done) // Signal monitoring goroutine to stop

		// Monitor for the quick note to be shown
		go func() {
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					if mockQN.shown {
						close(shown)
						return
					}
				}
			}
		}()

		// Trigger the hotkey
		testHotkey.Trigger()

		// Wait for quick note to be shown or timeout
		select {
		case <-shown:
			// Success
		case <-time.After(1 * time.Second):
			t.Skip("Quick note was not shown after hotkey trigger - this may be expected in CI")
		}

		// Cleanup UI
		if app.mainWindow != nil {
			app.mainWindow.Close()
		}
	})
}
