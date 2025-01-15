package quicknote

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestWindow(t *testing.T) (*Window, *storage.MockStore) {
	store := storage.NewMockStore()
	log := logger.NewTestLogger(t)
	app := test.NewApp()

	// Create main window first
	mainWinCfg := config.WindowConfig{
		Width:       800,
		Height:      600,
		StartHidden: false,
	}
	mainWin := mainwindow.New(app, store, log, mainWinCfg)

	// Create quick note window
	quickNoteCfg := config.WindowConfig{
		Width:       400,
		Height:      300,
		StartHidden: true,
	}
	quickNote := New(app, store, log, quickNoteCfg, mainWin)
	return quickNote, store
}

func TestWindow(t *testing.T) {
	ctx := context.Background()

	t.Run("Creation", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		require.NotNil(t, window)
		assert.NotNil(t, window.input, "Input field should be initialized")
		assert.Equal(t, "", window.input.Text, "Input field should be empty")
		assert.NotNil(t, window.window, "Window should be initialized")
		assert.NotNil(t, window.saveBtn, "Save button should be initialized")
	})

	t.Run("Show", func(t *testing.T) {
		window, _ := setupTestWindow(t)

		// Set some text to verify it gets cleared
		window.input.SetText("test")

		// Show the window
		window.Show()

		// Verify state
		canvas := window.window.Canvas()
		assert.Equal(t, window.input, canvas.Focused(), "Input field should have focus")
		assert.Equal(t, "", window.input.Text, "Input field should be cleared")
	})

	t.Run("Hide", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		window.Show()
		window.Hide()
		// Note: In test environment, we can't directly verify window visibility
		// but we can verify the window exists
		assert.NotNil(t, window.window)
	})

	t.Run("SaveTask", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Simulate entering text
		window.input.SetText("Test Task")

		// Click the save button
		test.Tap(window.saveBtn)

		// Verify task was saved
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		require.Len(t, tasks, 1)
		assert.Equal(t, "Test Task", tasks[0].Content)
		assert.False(t, tasks[0].Done)
	})

	t.Run("SaveEmptyTask", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Click the save button
		test.Tap(window.saveBtn)

		// Verify no task was saved
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("SaveError", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Set up store error
		store.Error = assert.AnError

		// Try to save task
		window.input.SetText("Test Task")

		// Click the save button
		test.Tap(window.saveBtn)

		// Verify window is still shown (not hidden after error)
		assert.NotNil(t, window.window)
		assert.Equal(t, "Test Task", window.input.Text)
	})
}
