package sqlite

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore(t *testing.T) {
	// Setup
	log, err := logger.NewZapLogger(&logger.Config{
		Level:   "debug",
		Console: true,
	})
	require.NoError(t, err)

	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer store.Close()

	// Test CRUD operations
	t.Run("CRUD operations", func(t *testing.T) {
		// Create
		todo := &model.Todo{
			ID:        uuid.New().String(),
			Content:   "Test todo",
			Done:      false,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		require.NoError(t, store.Add(todo))

		// Read
		retrieved, err := store.Get(todo.ID)
		require.NoError(t, err)
		assert.Equal(t, todo.Content, retrieved.Content)

		// Update
		todo.Done = true
		require.NoError(t, store.Update(todo))
		updated, err := store.Get(todo.ID)
		require.NoError(t, err)
		assert.True(t, updated.Done)

		// Delete
		require.NoError(t, store.Delete(todo.ID))
		_, err = store.Get(todo.ID)
		assert.ErrorIs(t, err, storage.ErrTodoNotFound)
	})
}
