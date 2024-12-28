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
	_ "modernc.org/sqlite" // SQLite driver
)

func TestSQLiteStore(t *testing.T) {
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
	ctx := context.Background()

	t.Run("Add and List", func(t *testing.T) {
		task := storage.Task{
			ID:        "1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, task)
		assert.NoError(t, err)

		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)

		// Compare fields individually
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Content, tasks[0].Content)
		assert.Equal(t, task.Done, tasks[0].Done)
		// Compare times with some tolerance for database rounding
		assert.WithinDuration(t, task.CreatedAt, tasks[0].CreatedAt, time.Second)
		assert.WithinDuration(t, task.UpdatedAt, tasks[0].UpdatedAt, time.Second)
	})

	t.Run("Update", func(t *testing.T) {
		task := storage.Task{
			ID:        "1",
			Content:   "Updated Task",
			Done:      true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Update(ctx, task)
		assert.NoError(t, err)

		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", tasks[0].Content)
		assert.True(t, tasks[0].Done)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, "1")
		assert.NoError(t, err)

		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, tasks)

		err = store.Delete(ctx, "1")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})
}
