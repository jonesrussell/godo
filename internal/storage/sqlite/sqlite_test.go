package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore(t *testing.T) {
	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	if _, err := logger.Initialize(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Create store
	store, err := New(dbPath)
	require.NoError(t, err)
	defer store.Close()

	// Run common tests
	storage.TestStore(t, store)
}

func TestSQLiteStore_Persistence(t *testing.T) {
	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	if _, err := logger.Initialize(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "persistence_test.db")

	// Create first store instance
	store1, err := New(dbPath)
	require.NoError(t, err)

	// Add a todo
	todo := model.NewTodo("Test persistence")
	err = store1.Add(todo)
	assert.NoError(t, err)
	store1.Close()

	// Create second store instance
	store2, err := New(dbPath)
	require.NoError(t, err)
	defer store2.Close()

	// Verify todo persists
	got, err := store2.Get(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo, got)
}

func TestSQLiteStore_InvalidPath(t *testing.T) {
	// Initialize logger with test config
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}

	if _, err := logger.Initialize(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	// Try to create store with invalid path
	store, err := New("/nonexistent/path/test.db")
	assert.Error(t, err)
	assert.Nil(t, store)
}
