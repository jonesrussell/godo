package sqlite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
	storagetesting "github.com/jonesrussell/godo/internal/storage/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestStore(t *testing.T) *Store {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	log := logger.NewTestLogger(t)
	store, err := New(dbPath, log)
	require.NoError(t, err)
	t.Cleanup(func() {
		store.Close()
		os.Remove(dbPath)
	})
	return store
}

func TestSQLiteStoreWithSuite(t *testing.T) {
	// Run the standard storage test suite
	suite := &storagetesting.StoreSuite{
		NewStore: func() storage.Store {
			return setupTestStore(t)
		},
	}
	suite.Run(t)

	// Additional SQLite-specific tests
	t.Run("InvalidPath", func(t *testing.T) {
		log := logger.NewTestLogger(t)
		_, err := New("", log)
		assert.Error(t, err)
	})

	t.Run("DuplicateID", func(t *testing.T) {
		store := setupTestStore(t)
		defer store.Close()

		task := storage.Task{ID: "test-1", Content: "Test Task"}
		err := store.Add(task)
		require.NoError(t, err)

		err = store.Add(task)
		assert.ErrorIs(t, err, errors.ErrDuplicateID)
	})

	t.Run("UpdateNonexistent", func(t *testing.T) {
		store := setupTestStore(t)
		defer store.Close()

		task := storage.Task{ID: "nonexistent", Content: "Test Task"}
		err := store.Update(task)
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
	})

	t.Run("DeleteNonexistent", func(t *testing.T) {
		store := setupTestStore(t)
		defer store.Close()

		err := store.Delete("nonexistent")
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
	})

	t.Run("GetByIDNonexistent", func(t *testing.T) {
		store := setupTestStore(t)
		defer store.Close()

		_, err := store.GetByID("nonexistent")
		assert.ErrorIs(t, err, errors.ErrTaskNotFound)
	})

	t.Run("EmptyID", func(t *testing.T) {
		store := setupTestStore(t)
		defer store.Close()

		task := storage.Task{ID: "", Content: "Test Task"}
		err := store.Add(task)
		assert.ErrorIs(t, err, errors.ErrEmptyID)

		err = store.Update(task)
		assert.ErrorIs(t, err, errors.ErrEmptyID)

		err = store.Delete("")
		assert.ErrorIs(t, err, errors.ErrEmptyID)

		_, err = store.GetByID("")
		assert.ErrorIs(t, err, errors.ErrEmptyID)
	})

	t.Run("StoreClosed", func(t *testing.T) {
		store := setupTestStore(t)
		err := store.Close()
		require.NoError(t, err)

		task := storage.Task{ID: "test", Content: "Test Task"}
		err = store.Add(task)
		assert.ErrorIs(t, err, errors.ErrStoreClosed)

		err = store.Update(task)
		assert.ErrorIs(t, err, errors.ErrStoreClosed)

		err = store.Delete("test")
		assert.ErrorIs(t, err, errors.ErrStoreClosed)

		_, err = store.List()
		assert.ErrorIs(t, err, errors.ErrStoreClosed)

		_, err = store.GetByID("test")
		assert.ErrorIs(t, err, errors.ErrStoreClosed)

		err = store.Close()
		assert.ErrorIs(t, err, errors.ErrStoreClosed)
	})
}
