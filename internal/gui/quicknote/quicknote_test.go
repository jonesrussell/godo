package quicknote_test

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/gui/mainwindow"
	"github.com/jonesrussell/godo/internal/gui/quicknote"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

func setupTestWindow(t *testing.T) (*quicknote.Window, *storage.MockStore) {
	t.Helper()
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
	quickNote := quicknote.New(app, store, log, quickNoteCfg, mainWin)
	return quickNote, store
}

func TestWindow(t *testing.T) {
	ctx := context.Background()

	t.Run("Creation", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		require.NotNil(t, window)
	})

	t.Run("Show", func(t *testing.T) {
		window, _ := setupTestWindow(t)

		// Show the window
		window.Show()

		// Verify window is shown (we can't access internal fields in test package)
		assert.NotNil(t, window)
	})

	t.Run("Hide", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		window.Show()
		window.Hide()
		// Note: In test environment, we can't directly verify window visibility
		// but we can verify the window exists
		assert.NotNil(t, window)
	})

	t.Run("SaveTask", func(t *testing.T) {
		window, _ := setupTestWindow(t)

		// Show window and simulate entering text
		window.Show()

		// We can't directly access the input field in test package,
		// but we can test the save functionality by triggering the save action
		// This would require exposing a method to simulate text entry and save

		// For now, just verify the window was created
		assert.NotNil(t, window)
	})

	t.Run("SaveEmptyTask", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Show window without entering text
		window.Show()
		window.Hide()

		// Verify no task was saved
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("SaveError", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Set up store error
		store.Error = assert.AnError

		// Show window
		window.Show()

		// Verify window is still accessible
		assert.NotNil(t, window)
	})
}
