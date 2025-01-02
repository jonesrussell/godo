package quicknote

import (
	"context"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestQuickNote(t *testing.T) (*Window, *storage.MockStore) {
	store := storage.NewMockStore()
	log := logger.NewMockTestLogger(t)
	app := test.NewApp()
	cfg := config.WindowConfig{
		Width:       200,
		Height:      100,
		StartHidden: false,
	}
	quickNote := New(app, store, log, cfg)
	return quickNote, store
}

func TestQuickNote(t *testing.T) {
	ctx := context.Background()

	t.Run("AddTask", func(t *testing.T) {
		quickNote, store := setupTestQuickNote(t)
		require.NotNil(t, quickNote)

		// Add a task
		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Verify the task was added
		addedTask, err := store.GetByID(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, addedTask.ID)
		assert.Equal(t, task.Content, addedTask.Content)
	})

	t.Run("WindowClose", func(t *testing.T) {
		quickNote, store := setupTestQuickNote(t)
		require.NotNil(t, quickNote)

		// Close the window
		quickNote.Hide()

		// Verify the store is closed
		err := store.Close()
		assert.NoError(t, err)
	})

	t.Run("StoreError", func(t *testing.T) {
		quickNote, store := setupTestQuickNote(t)
		require.NotNil(t, quickNote)

		// Set store error
		store.Error = assert.AnError

		// Try to add a task
		task := storage.Task{
			ID:        "test-2",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(ctx, task)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})

	t.Run("EmptyContent", func(t *testing.T) {
		quickNote, store := setupTestQuickNote(t)
		require.NotNil(t, quickNote)

		// Try to add a task with empty content
		task := storage.Task{
			ID:        "test-3",
			Content:   "",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Verify the task was added
		addedTask, err := store.GetByID(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, addedTask.ID)
		assert.Empty(t, addedTask.Content)
	})
}

type mockStore struct {
	storage.TaskStore
}

func TestQuickNote_Show(t *testing.T) {
	// Create test dependencies
	testApp := test.NewApp()
	store := &mockStore{}
	log := logger.NewTestLogger(t)
	cfg := config.WindowConfig{
		Width:  200,
		Height: 100,
	}

	// Create quick note window
	quickNote := NewQuickNote(testApp, store, log, cfg)

	// Test initial state
	assert.NotNil(t, quickNote.input, "Input field should be initialized")
	assert.Equal(t, "", quickNote.input.Text, "Input field should be empty")

	// Show the window
	quickNote.Show()

	// Test that input has focus
	canvas := quickNote.window.Canvas()
	assert.Equal(t, quickNote.input, canvas.Focused(), "Input field should have focus")
	assert.Equal(t, "", quickNote.input.Text, "Input field should be cleared")
}
