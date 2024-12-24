package testing

import (
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/memory"
	"github.com/jonesrussell/godo/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestStoreInterface(t *testing.T) {
	t.Run("memory store implements Store interface", func(_ *testing.T) {
		var _ storage.Store = (*memory.Store)(nil)
	})

	t.Run("sqlite store implements Store interface", func(_ *testing.T) {
		var _ storage.Store = (*sqlite.Store)(nil)
	})
}

func TestStoreBehavior(t *testing.T) {
	log := logger.NewTestLogger(t)

	stores := map[string]func() storage.Store{
		"memory": func() storage.Store {
			return memory.New()
		},
		"sqlite": func() storage.Store {
			store, err := sqlite.New("file::memory:?cache=shared", log)
			assert.NoError(t, err)
			return store
		},
	}

	for name, createStore := range stores {
		t.Run(name+" store behavior", func(t *testing.T) {
			store := createStore()

			// Test Todo operations
			t.Run("todo operations", func(t *testing.T) {
				// Add
				todo := model.NewTodo("Test todo")
				err := store.Add(todo)
				assert.NoError(t, err)

				// Get
				retrieved, err := store.Get(todo.ID)
				assert.NoError(t, err)
				assert.Equal(t, todo.Content, retrieved.Content)

				// List
				todos := store.List()
				assert.Len(t, todos, 1)
				assert.Equal(t, todo.ID, todos[0].ID)

				// Update
				todo.Content = "Updated todo"
				err = store.Update(todo)
				assert.NoError(t, err)

				// Delete
				err = store.Delete(todo.ID)
				assert.NoError(t, err)

				// Verify deletion
				todos = store.List()
				assert.Empty(t, todos)
			})

			// Test Note operations
			t.Run("note operations", func(t *testing.T) {
				// Save note
				err := store.SaveNote("Test note")
				assert.NoError(t, err)

				// Get notes
				notes, err := store.GetNotes()
				assert.NoError(t, err)
				assert.Len(t, notes, 1)
				assert.Equal(t, "Test note", notes[0])
			})
		})
	}
}
