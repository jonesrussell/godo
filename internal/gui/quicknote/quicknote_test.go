package quicknote

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestWindow(t *testing.T) (*Window, *storage.MockStore) {
	fixture := testutil.NewTestFixture(t)
	cfg := config.WindowConfig{
		Width:       testutil.DefaultQuickNoteWidth,
		Height:      testutil.DefaultQuickNoteHeight,
		StartHidden: false,
	}
	window := New(test.NewApp(), fixture.Store, fixture.Logger, cfg)
	return window, fixture.Store
}

func TestWindow(t *testing.T) {
	ctx := context.Background()

	t.Run("Creation", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		require.NotNil(t, window)
		assert.NotNil(t, window.input, "Input field should be initialized")
		assert.Equal(t, "", window.input.Text, "Input field should be empty")
		assert.NotNil(t, window.fyneWindow, "Window should be initialized")
	})

	t.Run("Show", func(t *testing.T) {
		window, _ := setupTestWindow(t)

		// Set some text to verify it gets cleared
		window.input.SetText("test")

		// Show the window
		window.Show()

		// Verify state
		assert.Equal(t, "", window.input.Text, "Input field should be cleared")
	})

	t.Run("Hide", func(t *testing.T) {
		window, _ := setupTestWindow(t)
		window.Show()
		window.Hide()
		// Note: In test environment, we can't directly verify window visibility
		// but we can verify the window exists
		assert.NotNil(t, window.fyneWindow)
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
		assert.Equal(t, "Test Task", tasks[0].Title)
		assert.False(t, tasks[0].Completed)
		assert.NotZero(t, tasks[0].CreatedAt)
		assert.NotZero(t, tasks[0].UpdatedAt)
		assert.Equal(t, tasks[0].CreatedAt, tasks[0].UpdatedAt)
	})

	t.Run("SaveEmptyTask", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Click the save button without entering text
		test.Tap(window.saveBtn)

		// Verify no task was saved
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("SaveError", func(t *testing.T) {
		window, store := setupTestWindow(t)

		// Set up store error
		store.SetError(assert.AnError)

		// Try to save task
		window.input.SetText("Test Task")

		// Click the save button
		test.Tap(window.saveBtn)

		// Verify input was not cleared (indicating save failed)
		assert.Equal(t, "Test Task", window.input.Text)

		// Verify no task was saved
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})
}
