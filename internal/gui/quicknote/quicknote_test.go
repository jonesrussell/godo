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

type testLogger struct {
	logger.Logger
}

func (m *testLogger) Debug(_ string, _ ...interface{}) {}
func (m *testLogger) Info(_ string, _ ...interface{})  {}
func (m *testLogger) Warn(_ string, _ ...interface{})  {}
func (m *testLogger) Error(_ string, _ ...interface{}) {}

func setupTestQuickNote() (*Window, *storage.MockStore) {
	store := storage.NewMockStore()
	log := &testLogger{}
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
	quickNote, store := setupTestQuickNote()
	require.NotNil(t, quickNote)
	ctx := context.Background()

	t.Run("AddTask", func(t *testing.T) {
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
		// Close the window
		quickNote.Hide()

		// Verify the store is closed
		err := store.Close()
		assert.NoError(t, err)
	})

	t.Run("StoreError", func(t *testing.T) {
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
		// Reset store error
		store.Error = nil

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
