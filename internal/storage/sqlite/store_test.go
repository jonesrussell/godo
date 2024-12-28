package sqlite

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // SQLite driver
)

func TestSQLiteStoreComprehensive(t *testing.T) {
	log := logger.NewTestLogger(t)
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer func() {
		store.Close()
		os.Remove(dbPath)
	}()

	now := time.Now()

	t.Run("New Store", func(t *testing.T) {
		// Test empty path
		invalidStore, err := New("", log)
		assert.ErrorIs(t, err, ErrInvalidPath)
		assert.Nil(t, invalidStore)

		// Test whitespace path
		invalidStore, err = New("   ", log)
		assert.ErrorIs(t, err, ErrInvalidPath)
		assert.Nil(t, invalidStore)

		// Test with valid path
		validStore, err := New(filepath.Join(tempDir, "valid.db"), log)
		assert.NoError(t, err)
		assert.NotNil(t, validStore)
		defer validStore.Close()
	})

	t.Run("Add and List", func(t *testing.T) {
		task := storage.Task{
			ID:        "1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)

		// Compare fields individually
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Content, tasks[0].Content)
		assert.Equal(t, task.Done, tasks[0].Done)
		assert.WithinDuration(t, task.CreatedAt, tasks[0].CreatedAt, time.Second)
		assert.WithinDuration(t, task.UpdatedAt, tasks[0].UpdatedAt, time.Second)

		// Test adding duplicate ID
		err = store.Add(task)
		assert.Error(t, err, "Should not allow duplicate IDs")

		// Test adding task with empty ID
		emptyIDTask := storage.Task{
			Content:   "Empty ID Task",
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Add(emptyIDTask)
		assert.Error(t, err, "Should not allow empty ID")
	})

	t.Run("Update", func(t *testing.T) {
		task := storage.Task{
			ID:        "1",
			Content:   "Updated Task",
			Done:      true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Test successful update
		err := store.Update(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", tasks[0].Content)
		assert.True(t, tasks[0].Done)

		// Test updating non-existent task
		nonExistentTask := storage.Task{
			ID:        "999",
			Content:   "Non-existent Task",
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Update(nonExistentTask)
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)

		// Test updating with empty ID
		emptyIDTask := storage.Task{
			Content:   "Empty ID Task",
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Update(emptyIDTask)
		assert.Error(t, err, "Should not allow empty ID")
	})

	t.Run("Delete", func(t *testing.T) {
		// Test successful delete
		err := store.Delete("1")
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)

		// Test deleting non-existent task
		err = store.Delete("1")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)

		// Test deleting with empty ID
		err = store.Delete("")
		assert.Error(t, err, "Should not allow empty ID")
	})

	t.Run("List with Multiple Tasks", func(t *testing.T) {
		// Add multiple tasks
		tasks := []storage.Task{
			{
				ID:        "1",
				Content:   "Task 1",
				Done:      false,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "2",
				Content:   "Task 2",
				Done:      true,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "3",
				Content:   "Task 3",
				Done:      false,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		for _, task := range tasks {
			err := store.Add(task)
			assert.NoError(t, err)
		}

		// Verify list returns all tasks
		listTasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, listTasks, len(tasks))

		// Verify task order (should be in insertion order)
		for i, task := range tasks {
			assert.Equal(t, task.ID, listTasks[i].ID)
			assert.Equal(t, task.Content, listTasks[i].Content)
			assert.Equal(t, task.Done, listTasks[i].Done)
		}
	})

	t.Run("Close and Reopen", func(t *testing.T) {
		// Add a task
		task := storage.Task{
			ID:        "close-test",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(task)
		assert.NoError(t, err)

		// Close the store
		err = store.Close()
		assert.NoError(t, err)

		// Try to use closed store
		_, err = store.List()
		assert.Error(t, err, "Should not allow operations on closed store")

		// Reopen the store
		newStore, err := New(dbPath, log)
		assert.NoError(t, err)
		defer newStore.Close()

		// Verify data persisted
		tasks, err := newStore.List()
		assert.NoError(t, err)
		assert.NotEmpty(t, tasks)
		found := false
		for _, t := range tasks {
			if t.ID == "close-test" {
				found = true
				break
			}
		}
		assert.True(t, found, "Task should persist after close and reopen")
	})
}
