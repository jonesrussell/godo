package sqlite

import (
	"path/filepath"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	// Setup test logger and database
	log, err := logger.New(&common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	})
	require.NoError(t, err)

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer store.Close()

	// Test todo operations
	t.Run("todo operations", func(t *testing.T) {
		// Add todo
		todo := model.NewTodo("Test todo")
		err := store.Add(todo)
		require.NoError(t, err)

		// Get todo
		retrieved, err := store.Get(todo.ID)
		require.NoError(t, err)
		assert.Equal(t, todo.Content, retrieved.Content)

		// Update todo
		todo.Content = "Updated content"
		err = store.Update(todo)
		require.NoError(t, err)

		// List todos
		todos := store.List()
		assert.Len(t, todos, 1)
		assert.Equal(t, "Updated content", todos[0].Content)

		// Delete todo
		err = store.Delete(todo.ID)
		require.NoError(t, err)
		assert.Empty(t, store.List())
	})

	// Test note operations
	t.Run("note operations", func(t *testing.T) {
		// Save note
		note := "Test note"
		err := store.SaveNote(note)
		require.NoError(t, err)

		// Get notes
		notes, err := store.GetNotes()
		require.NoError(t, err)
		assert.Contains(t, notes, note)
	})
}

func TestConcurrency(t *testing.T) {
	// Setup test logger and database
	log, err := logger.New(&common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	})
	require.NoError(t, err)

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	store, err := New(dbPath, log)
	require.NoError(t, err)
	defer store.Close()

	done := make(chan bool)
	const numGoroutines = 10

	// Test concurrent todo operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			todo := model.NewTodo("Concurrent todo")
			err := store.Add(todo)
			assert.NoError(t, err)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify results
	todos := store.List()
	assert.Len(t, todos, numGoroutines)
}
