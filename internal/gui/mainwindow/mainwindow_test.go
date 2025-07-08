package mainwindow

import (
	"context"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

func createTestTask(content string) storage.Task {
	return storage.Task{
		ID:        "test-id",
		Content:   content,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func setupTestWindow(t *testing.T) (*Window, *storage.MockStore) {
	t.Helper()
	store := storage.NewMockStore()
	log := logger.NewTestLogger(t)
	app := test.NewApp()
	cfg := config.WindowConfig{
		Width:       800,
		Height:      600,
		StartHidden: false,
	}
	mainWindow := New(app, store, log, cfg)
	mainWindow.Show()
	return mainWindow, store
}

func TestTaskManager(t *testing.T) {
	window, store := setupTestWindow(t)
	tm := window.TaskManager
	ctx := context.Background()

	t.Run("LoadTasks", func(t *testing.T) {
		// Add some tasks to store
		task := createTestTask("Test Task")
		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Load tasks
		err = tm.loadTasks(ctx)
		require.NoError(t, err)
		assert.Len(t, tm.tasks, 1)
		assert.Equal(t, task.Content, tm.tasks[0].Content)
	})

	t.Run("AddTask", func(t *testing.T) {
		// Add task
		err := tm.addTask(ctx, "New Task")
		require.NoError(t, err)

		// Verify task was added
		assert.Len(t, tm.tasks, 2)
		assert.Equal(t, "New Task", tm.tasks[1].Content)
		assert.False(t, tm.tasks[1].Done)

		// Verify task is in store
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("UpdateTask", func(t *testing.T) {
		// Update first task
		err := tm.updateTask(ctx, 0, true)
		require.NoError(t, err)

		// Verify task was updated
		assert.True(t, tm.tasks[0].Done)
		task, err := store.GetByID(ctx, tm.tasks[0].ID)
		require.NoError(t, err)
		assert.True(t, task.Done)
	})

	t.Run("EmptyTask", func(t *testing.T) {
		initialCount := len(tm.tasks)
		err := tm.addTask(ctx, "")
		require.NoError(t, err)
		assert.Len(t, tm.tasks, initialCount)
	})
}

func TestWindow(t *testing.T) {
	window, _ := setupTestWindow(t)
	defer window.Hide()

	t.Run("InitialState", func(t *testing.T) {
		assert.NotNil(t, window.TaskManager)
		assert.NotNil(t, window.window)
		assert.NotNil(t, window.list)
		assert.Empty(t, window.tasks)

		// Test window properties
		assert.Equal(t, "Godo - Task Manager", window.window.Title())
		// Skip icon test as it's not reliable in test environment
	})

	t.Run("UIComponents", func(t *testing.T) {
		// Test input creation
		input := window.createInput()
		assert.NotNil(t, input)
		assert.Equal(t, "Enter a new task...", input.PlaceHolder)

		// Test input submission
		test.Type(input, "Test Task")
		input.OnSubmitted(input.Text)
		assert.Empty(t, input.Text)
		assert.Len(t, window.tasks, 1)

		// Test button creation
		button := window.createAddButton(input)
		assert.NotNil(t, button)
		assert.Equal(t, "Add", button.Text)

		// Test list creation
		list := window.createTaskList()
		assert.NotNil(t, list)
		assert.Equal(t, 1, list.Length()) // Should have the task we added

		// Test task completion visual feedback
		window.updateTask(context.Background(), 0, true)
		list.Refresh()
		item := list.CreateItem()
		list.UpdateItem(0, item)
		label := item.(*fyne.Container).Objects[1].(*widget.Label)
		assert.True(t, label.TextStyle.Monospace)
	})

	t.Run("KeyboardShortcuts", func(t *testing.T) {
		input := window.createInput()

		// Test Enter key submission
		test.Type(input, "Task via Enter")
		input.OnSubmitted(input.Text)
		assert.Empty(t, input.Text)
		assert.Contains(t, window.tasks[len(window.tasks)-1].Content, "Task via Enter")
	})

	t.Run("WindowMethods", func(t *testing.T) {
		// Test window interface methods
		assert.NotPanics(t, func() {
			window.Show()
			window.Hide()
			window.CenterOnScreen()
			window.Resize(fyne.NewSize(100, 100))
			assert.NotNil(t, window.GetWindow())
		})
	})
}
