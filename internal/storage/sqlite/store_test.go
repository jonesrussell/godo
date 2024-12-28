package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
		closeErr := store.Close()
		if closeErr != nil {
			t.Errorf("Failed to close store: %v", closeErr)
		}
		removeErr := os.Remove(dbPath)
		if removeErr != nil {
			t.Errorf("Failed to remove test database: %v", removeErr)
		}
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

		// Compare fields with stricter time tolerance
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Content, tasks[0].Content)
		assert.Equal(t, task.Done, tasks[0].Done)
		assert.WithinDuration(t, task.CreatedAt, tasks[0].CreatedAt, 10*time.Millisecond)
		assert.WithinDuration(t, task.UpdatedAt, tasks[0].UpdatedAt, 10*time.Millisecond)

		// Test adding task with empty ID
		emptyIDTask := storage.Task{
			Content:   "Empty ID Task",
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Add(emptyIDTask)
		assert.ErrorIs(t, err, ErrEmptyID, "Should return ErrEmptyID for empty ID")
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

	t.Run("GetByID", func(t *testing.T) {
		// Create a new store for this test to avoid interference
		testStore, initErr := New(filepath.Join(tempDir, "getbyid.db"), log)
		require.NoError(t, initErr)
		defer func() {
			closeErr := testStore.Close()
			if closeErr != nil {
				t.Errorf("Failed to close store: %v", closeErr)
			}
			removeErr := os.Remove(filepath.Join(tempDir, "getbyid.db"))
			if removeErr != nil {
				t.Errorf("Failed to remove test database: %v", removeErr)
			}
		}()

		// Add a task
		task := storage.Task{
			ID:        "get-by-id-test",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		addErr := testStore.Add(task)
		assert.NoError(t, addErr)

		// Test successful retrieval
		retrieved, getErr := testStore.GetByID(task.ID)
		assert.NoError(t, getErr)
		assert.NotNil(t, retrieved)
		assert.Equal(t, task.ID, retrieved.ID)
		assert.Equal(t, task.Content, retrieved.Content)
		assert.Equal(t, task.Done, retrieved.Done)
		assert.WithinDuration(t, task.CreatedAt, retrieved.CreatedAt, time.Second)
		assert.WithinDuration(t, task.UpdatedAt, retrieved.UpdatedAt, time.Second)

		// Test non-existent task
		retrieved, notFoundErr := testStore.GetByID("nonexistent")
		assert.ErrorIs(t, notFoundErr, storage.ErrTaskNotFound)
		assert.Nil(t, retrieved)

		// Test empty ID
		retrieved, emptyIDErr := testStore.GetByID("")
		assert.ErrorIs(t, emptyIDErr, ErrEmptyID)
		assert.Nil(t, retrieved)

		// Test with closed store
		testStore.Close()
		retrieved, closedErr := testStore.GetByID(task.ID)
		assert.ErrorIs(t, closedErr, ErrStoreClosed)
		assert.Nil(t, retrieved)
	})
}

func TestConcurrentOperations(t *testing.T) {
	log := logger.NewTestLogger(t)
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "concurrent_test.db")

	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer func() {
		err := store.Close()
		if err != nil {
			t.Errorf("Failed to close store: %v", err)
		}
		err = os.Remove(dbPath)
		if err != nil {
			t.Errorf("Failed to remove test database: %v", err)
		}
	}()

	const numGoroutines = 10
	const opsPerGoroutine = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*opsPerGoroutine)
	start := time.Now()

	// Launch multiple goroutines to perform concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			for j := 0; j < opsPerGoroutine; j++ {
				// Create unique task ID for this operation
				taskID := fmt.Sprintf("task-%d-%d", routineID, j)

				// Add task
				task := storage.Task{
					ID:        taskID,
					Content:   fmt.Sprintf("Concurrent Task %s", taskID),
					Done:      false,
					CreatedAt: start,
					UpdatedAt: start,
				}

				if err := store.Add(task); err != nil {
					errors <- fmt.Errorf("failed to add task %s: %w", taskID, err)
					continue
				}

				// Read task
				_, err := store.GetByID(taskID)
				if err != nil {
					errors <- fmt.Errorf("failed to get task %s: %w", taskID, err)
					continue
				}

				// Update task
				task.Done = true
				task.UpdatedAt = time.Now()
				if err := store.Update(task); err != nil {
					errors <- fmt.Errorf("failed to update task %s: %w", taskID, err)
					continue
				}

				// Delete task
				if err := store.Delete(taskID); err != nil {
					errors <- fmt.Errorf("failed to delete task %s: %w", taskID, err)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)

	// Check for any errors
	var errCount int
	for err := range errors {
		errCount++
		t.Errorf("Concurrent operation error: %v", err)
	}

	assert.Zero(t, errCount, "Expected no errors in concurrent operations")

	// Verify final state
	tasks, err := store.List()
	assert.NoError(t, err)
	assert.Empty(t, tasks, "All tasks should have been deleted")
}
