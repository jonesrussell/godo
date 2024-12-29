package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStore(t *testing.T) (store storage.TaskTxStore, cleanup func()) {
	t.Helper()

	// Create a temporary database file
	dbFile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	dbPath := dbFile.Name()
	require.NoError(t, dbFile.Close())

	// Create test logger
	log := logger.NewTestLogger(t)

	// Create the store
	store, err = New(dbPath, log)
	require.NoError(t, err)

	cleanup = func() {
		require.NoError(t, store.Close())
		require.NoError(t, os.Remove(dbPath))
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
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Content, tasks[0].Content)
	})

	t.Run("get task", func(t *testing.T) {
		task, err := store.GetByID(ctx, "test-1")
		require.NoError(t, err)
		assert.Equal(t, "test-1", task.ID)
		assert.Equal(t, "Test Task", task.Content)

		// Try to get nonexistent task
		_, err = store.GetByID(ctx, "nonexistent")
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
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
		assertTaskEqual(t, task, updated)

		// Try to update nonexistent task
		task.ID = "nonexistent"
		err = store.Update(ctx, task)
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
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
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
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
		assertTaskEqual(t, task1, tasks[0])

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
