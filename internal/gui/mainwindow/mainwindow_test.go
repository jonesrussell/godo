package mainwindow

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

func setupTestWindow(t *testing.T) (*Window, *storage.MockStore) {
	store := storage.NewMockStore()
	log := logger.NewMockTestLogger(t)
	app := test.NewApp()
	cfg := config.WindowConfig{
		Width:       800,
		Height:      600,
		StartHidden: false,
	}
	mainWindow := New(app, store, log, cfg)
	return mainWindow, store
}

func TestMainWindow(t *testing.T) {
	mainWindow, store := setupTestWindow(t)
	require.NotNil(t, mainWindow)
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

		// Get the task to verify it was added
		addedTask, err := store.GetByID(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, addedTask.ID)
		assert.Equal(t, task.Content, addedTask.Content)
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
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Update the task
		task.Done = true
		err = store.Update(ctx, task)
		require.NoError(t, err)

		// Get the task to verify it was updated
		updatedTask, err := store.GetByID(ctx, task.ID)
		require.NoError(t, err)
		assert.True(t, updatedTask.Done)
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
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Delete the task
		err = store.Delete(ctx, task.ID)
		require.NoError(t, err)

		// Verify the task is deleted
		_, err = store.GetByID(ctx, task.ID)
		assert.Error(t, err)
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Add some tasks
		task1 := storage.Task{
			ID:        "test-4",
			Content:   "Test Task 1",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		task2 := storage.Task{
			ID:        "test-5",
			Content:   "Test Task 2",
			Done:      true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := store.Add(ctx, task1)
		require.NoError(t, err)
		err = store.Add(ctx, task2)
		require.NoError(t, err)

		// List tasks
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tasks), 2)
	})

	t.Run("WindowClose", func(t *testing.T) {
		// Close the window
		mainWindow.Hide()

		// Verify the store is closed
		err := store.Close()
		assert.NoError(t, err)
	})
}
