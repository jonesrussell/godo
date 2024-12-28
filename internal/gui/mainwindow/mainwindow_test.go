package mainwindow

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLogger struct {
	logger.Logger
}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}

func setupTestWindow(t *testing.T) (*Window, *mock.Store) {
	store := mock.New()
	log := &mockLogger{}
	mainWindow := New(store, log)
	return mainWindow, store
}

func TestMainWindow(t *testing.T) {
	mainWindow, store := setupTestWindow(t)
	require.NotNil(t, mainWindow)

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

	t.Run("UpdateTask", func(t *testing.T) {
		// Add a task
		task := storage.Task{
			ID:        "test-2",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(task)
		require.NoError(t, err)

		// Update the task
		task.Done = true
		err = store.Update(task)
		require.NoError(t, err)

		// Verify the task is updated
		assert.True(t, store.UpdateCalled)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		// Add a task
		task := storage.Task{
			ID:        "test-3",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(task)
		require.NoError(t, err)

		// Delete the task
		err = store.Delete(task.ID)
		require.NoError(t, err)

		// Verify the task is deleted
		assert.True(t, store.DeleteCalled)
	})

	t.Run("ListTasks", func(t *testing.T) {
		// List tasks
		tasks, err := store.List()
		require.NoError(t, err)

		// Verify tasks are listed
		assert.True(t, store.ListCalled)
		assert.NotNil(t, tasks)
	})

	t.Run("WindowClose", func(t *testing.T) {
		// Close the window
		mainWindow.Hide()

		// Verify the store is closed
		err := store.Close()
		assert.NoError(t, err)
	})
}
