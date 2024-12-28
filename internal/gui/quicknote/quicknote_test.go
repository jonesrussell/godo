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

type testLogger struct {
	logger.Logger
}

func (m *testLogger) Debug(_ string, _ ...interface{}) {}
func (m *testLogger) Info(_ string, _ ...interface{})  {}
func (m *testLogger) Warn(_ string, _ ...interface{})  {}
func (m *testLogger) Error(_ string, _ ...interface{}) {}

func setupTestQuickNote(t *testing.T) (*Window, *mock.Store) {
	store := mock.New()
	log := &testLogger{}
	quickNote := New(store, log)
	return quickNote, store
}

func TestQuickNote(t *testing.T) {
	quickNote, store := setupTestQuickNote(t)
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
