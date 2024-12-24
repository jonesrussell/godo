package app

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
)

// MockQuickNoteService implements QuickNoteService for testing
type MockQuickNoteService struct {
	shown     bool
	showCount int
}

func (m *MockQuickNoteService) Show() {
	m.shown = true
	m.showCount++
}

func (m *MockQuickNoteService) Hide() {
	m.shown = false
}

func setupTestApp(t *testing.T) (*App, *MockQuickNoteService) {
	t.Helper()

	log, err := logger.New(&common.LogConfig{
		Level:   "debug",
		Console: true,
	})
	assert.NoError(t, err)

	cfg := &config.Config{
		App: config.AppConfig{
			Name:    "Test App",
			Version: "0.1.0",
		},
		Hotkeys: config.HotkeyConfig{
			QuickNote: "Ctrl+Alt+G",
		},
	}

	store := mock.NewStore()
	app := NewApp(cfg, store, log)
	qn := &MockQuickNoteService{}
	app.SetQuickNoteService(qn)

	return app, qn
}

func TestHotkeyIntegration(t *testing.T) {
	t.Run("lifecycle and hotkey setup", func(t *testing.T) {
		app, qn := setupTestApp(t)
		defer app.Cleanup()

		// Setup UI which includes lifecycle events
		app.SetupUI()

		// Verify initial state
		assert.False(t, qn.shown)
		assert.Equal(t, 0, qn.showCount)

		// Simulate hotkey by calling Show directly
		app.ShowQuickNote()
		assert.True(t, qn.shown)
		assert.Equal(t, 1, qn.showCount)

		// Hide and verify
		qn.Hide()
		assert.False(t, qn.shown)

		// Trigger multiple shows to verify counter
		app.ShowQuickNote()
		app.ShowQuickNote()
		assert.Equal(t, 3, qn.showCount)
	})

	t.Run("hotkey display in UI", func(t *testing.T) {
		app, _ := setupTestApp(t)
		defer app.Cleanup()

		app.SetupUI()
		win := test.NewWindow(app.GetMainWindow().Content())
		defer win.Close()

		// Find the label with hotkey text
		var found bool
		expectedText := "Press Ctrl+Alt+G for quick notes"

		// Walk through all containers to find the label
		walkContainers(win.Content(), func(w fyne.CanvasObject) {
			if label, ok := w.(*widget.Label); ok {
				if label.Text == expectedText {
					found = true
				}
			}
		})

		assert.True(t, found, "Hotkey text should be visible in the UI")
	})
}

// walkContainers recursively walks through containers to find widgets
func walkContainers(o fyne.CanvasObject, fn func(fyne.CanvasObject)) {
	fn(o)

	switch cont := o.(type) {
	case *fyne.Container:
		for _, item := range cont.Objects {
			walkContainers(item, fn)
		}
	case *container.Split:
		walkContainers(cont.Leading, fn)
		walkContainers(cont.Trailing, fn)
	case *container.Scroll:
		walkContainers(cont.Content, fn)
	}
}
