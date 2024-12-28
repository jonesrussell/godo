//go:build windows

package quicknote

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindowsQuickNote(t *testing.T) {
	store := mock.New()
	log := logger.NewTestLogger(t)
	quickNote := New(store, log)
	require.NotNil(t, quickNote)

	t.Run("AddTask", func(t *testing.T) {
		// Add a task
		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(task)
		require.NoError(t, err)

		// Verify the task is added
		assert.True(t, store.AddCalled)
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
		err := store.Add(task)
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
		err := store.Add(task)
		require.NoError(t, err)

		// Verify the task is added
		assert.True(t, store.AddCalled)
	})
}
