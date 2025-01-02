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
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Get the task to verify it was added
		addedTask, err := store.Get(ctx, task.ID)
		require.NoError(t, err)
		assert.Equal(t, task.ID, addedTask.ID)
		assert.Equal(t, task.Title, addedTask.Title)
		assert.Equal(t, task.Completed, addedTask.Completed)
		assert.Equal(t, task.CreatedAt, addedTask.CreatedAt)
		assert.Equal(t, task.UpdatedAt, addedTask.UpdatedAt)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		// Add a task
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-2",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Update the task
		task.Completed = true
		task.UpdatedAt = time.Now().Unix()
		err = store.Update(ctx, task)
		require.NoError(t, err)

		// Get the task to verify it was updated
		updatedTask, err := store.Get(ctx, task.ID)
		require.NoError(t, err)
		assert.True(t, updatedTask.Completed)
		assert.Equal(t, task.UpdatedAt, updatedTask.UpdatedAt)
	})

	t.Run("DeleteTask", func(t *testing.T) {
		// Add a task
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-3",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Delete the task
		err = store.Delete(ctx, task.ID)
		require.NoError(t, err)

		// Verify the task is deleted
		_, err = store.Get(ctx, task.ID)
		assert.Error(t, err)
	})

	t.Run("ListTasks", func(t *testing.T) {
		// Add some tasks
		now := time.Now().Unix()
		task1 := storage.Task{
			ID:        "test-4",
			Title:     "Test Task 1",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, task1)
		require.NoError(t, err)

		task2 := storage.Task{
			ID:        "test-5",
			Title:     "Test Task 2",
			Completed: true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Add(ctx, task2)
		require.NoError(t, err)

		// List tasks
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)

		// Verify task order and fields
		assert.Equal(t, task1.ID, tasks[0].ID)
		assert.Equal(t, task1.Title, tasks[0].Title)
		assert.Equal(t, task1.Completed, tasks[0].Completed)
		assert.Equal(t, task2.ID, tasks[1].ID)
		assert.Equal(t, task2.Title, tasks[1].Title)
		assert.Equal(t, task2.Completed, tasks[1].Completed)
	})

	t.Run("WindowClose", func(t *testing.T) {
		// Close the window
		mainWindow.Hide()

		// Verify the store is closed
		err := store.Close()
		assert.NoError(t, err)
	})
}
