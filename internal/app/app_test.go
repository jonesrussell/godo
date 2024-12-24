//go:build !docker
// +build !docker

package app

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	// Use Fyne test app with test driver
	fyneApp := test.NewApp()
	defer fyneApp.Quit()

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

	// Create app without hotkey factory
	app := NewApp(cfg, store, log, nil)
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

	// Test UI setup without hotkey
	t.Run("ui setup", func(t *testing.T) {
		app.SetupUI()
		assert.NotNil(t, app.mainWindow)
		assert.True(t, mockSystray.ready)
	})
}
