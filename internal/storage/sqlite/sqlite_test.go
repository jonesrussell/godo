package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.NewZapLogger(&logger.Config{
		Level:    "debug",
		Console:  true,
		File:     false,
		FilePath: "",
	})
	require.NoError(t, err)
	return log
}

func TestSQLiteStore(t *testing.T) {
	log := setupTestLogger(t)

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create store
	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer store.Close()

	// Run common tests
	storage.TestStore(t, store)
}

func TestSQLiteStore_Persistence(t *testing.T) {
	log := setupTestLogger(t)

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persistence_test.db")

	// Create first store instance
	store1, err := New(dbPath, log)
	require.NoError(t, err)

	// Add a todo
	todo := model.NewTodo("Test persistence")
	err = store1.Add(todo)
	assert.NoError(t, err)
	store1.Close()

	// Create second store instance
	store2, err := New(dbPath, log)
	require.NoError(t, err)
	defer store2.Close()

	// Verify todo persists
	got, err := store2.Get(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo, got)
}

func TestSQLiteStore_InvalidPath(t *testing.T) {
	log := setupTestLogger(t)

	// Try to create store with invalid path
	store, err := New("/nonexistent/path/test.db", log)
	assert.Error(t, err)
	assert.Nil(t, store)
}
