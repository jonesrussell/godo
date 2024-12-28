package sqlite

import (
	"os"
	"path/filepath"
	"testing"

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

	t.Run("Add and List", func(t *testing.T) {
		task := storage.Task{
			ID:    "1",
			Title: "Test Task",
		}

		err := store.Add(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task, tasks[0])
	})

	t.Run("Update", func(t *testing.T) {
		task := storage.Task{
			ID:    "1",
			Title: "Updated Task",
		}

		err := store.Update(task)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", tasks[0].Title)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete("1")
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)

		err = store.Delete("1")
		assert.ErrorIs(t, err, storage.ErrTaskNotFound)
	})
}
