package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStore(t *testing.T) (storage.TaskTxStore, func()) {
	t.Helper()

	// Create temporary directory for test database
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create test store
	log := logger.NewTestLogger(t)
	store, err := New(dbPath, log)
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		store.Close()
		os.Remove(dbPath)
	}

	return store, cleanup
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	store, cleanup := setupTestStore(t)
	defer cleanup()

	t.Run("empty store", func(t *testing.T) {
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("add task", func(t *testing.T) {
		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}

		err := store.Add(ctx, task)
		require.NoError(t, err)

		// Verify task was added
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 1)
		assertTaskEqual(t, task, tasks[0])

		// Try to add duplicate
		err = store.Add(ctx, task)
		assert.Error(t, err)
	})

	t.Run("get task", func(t *testing.T) {
		task, err := store.GetByID(ctx, "test-1")
		require.NoError(t, err)
		assert.Equal(t, "test-1", task.ID)
		assert.Equal(t, "Test Task", task.Content)

		// Try to get nonexistent task
		_, err = store.GetByID(ctx, "nonexistent")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})

	t.Run("update task", func(t *testing.T) {
		task := storage.Task{
			ID:        "test-1",
			Content:   "Updated Task",
			Done:      true,
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}

		err := store.Update(ctx, task)
		require.NoError(t, err)

		// Verify task was updated
		updated, err := store.GetByID(ctx, "test-1")
		require.NoError(t, err)
		assert.Equal(t, task.Content, updated.Content)
		assert.Equal(t, task.Done, updated.Done)

		// Try to update nonexistent task
		task.ID = "nonexistent"
		err = store.Update(ctx, task)
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})

	t.Run("delete task", func(t *testing.T) {
		err := store.Delete(ctx, "test-1")
		require.NoError(t, err)

		// Verify task was deleted
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, tasks)

		// Try to delete nonexistent task
		err = store.Delete(ctx, "test-1")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	store, cleanup := setupTestStore(t)
	defer cleanup()

	t.Run("commit", func(t *testing.T) {
		// Add initial task
		task1 := storage.Task{
			ID:        "test-1",
			Content:   "Test Task 1",
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}
		require.NoError(t, store.Add(ctx, task1))

		// Start transaction
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		// Add second task in transaction
		task2 := storage.Task{
			ID:        "test-2",
			Content:   "Test Task 2",
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}
		require.NoError(t, tx.Add(ctx, task2))

		// Update first task in transaction
		task1.Content = "Updated Task 1"
		require.NoError(t, tx.Update(ctx, task1))

		// Verify changes are not visible outside transaction
		tasks, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, "Test Task 1", tasks[0].Content)

		// Commit transaction
		require.NoError(t, tx.Commit())

		// Verify changes are now visible
		tasks, err = store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})

	t.Run("rollback", func(t *testing.T) {
		// Start transaction
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		// Delete all tasks in transaction
		tasks, err := tx.List(ctx)
		require.NoError(t, err)
		for _, task := range tasks {
			require.NoError(t, tx.Delete(ctx, task.ID))
		}

		// Add new task in transaction
		task := storage.Task{
			ID:        "test-3",
			Content:   "Test Task 3",
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
		}
		require.NoError(t, tx.Add(ctx, task))

		// Rollback transaction
		require.NoError(t, tx.Rollback())

		// Verify original state is preserved
		tasks, err = store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, tasks, 2)
	})
}

func assertTaskEqual(t *testing.T, expected, actual storage.Task) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Content, actual.Content)
	assert.Equal(t, expected.Done, actual.Done)
	assert.Equal(t, expected.CreatedAt.Unix(), actual.CreatedAt.Unix())
	assert.Equal(t, expected.UpdatedAt.Unix(), actual.UpdatedAt.Unix())
}
