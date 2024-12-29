//go:build windows

package quicknote

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowsQuickNote(t *testing.T) {
	store := storage.NewMockStore()
	log := logger.NewTestLogger(t)
	quickNote := New(store, log)
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