package storage

import (
	"testing"

	"github.com/jonesrussell/godo/internal/model"
	"github.com/stretchr/testify/assert"
)

// TestStore runs a suite of tests against any Store implementation
func TestStore(t *testing.T, s Store) {
	t.Helper()

	t.Run("Add", func(t *testing.T) {
		todo := model.NewTodo("Test todo")
		err := s.Add(todo)
		assert.NoError(t, err)

		// Verify retrieval
		got, err := s.Get(todo.ID)
		assert.NoError(t, err)
		assert.Equal(t, todo, got)
	})

	t.Run("Get non-existent", func(t *testing.T) {
		got, err := s.Get("non-existent-id")
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("List", func(t *testing.T) {
		// Clear existing todos
		todos := s.List()
		for _, todo := range todos {
			s.Delete(todo.ID)
		}

		// Add new todos
		todo1 := model.NewTodo("First todo")
		todo2 := model.NewTodo("Second todo")
		s.Add(todo1)
		s.Add(todo2)

		// Get list
		todos = s.List()
		assert.Len(t, todos, 2)
		assert.Contains(t, todos, todo1)
		assert.Contains(t, todos, todo2)
	})

	t.Run("Update", func(t *testing.T) {
		todo := model.NewTodo("Original content")
		s.Add(todo)

		todo.UpdateContent("Updated content")
		err := s.Update(todo)
		assert.NoError(t, err)

		got, err := s.Get(todo.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated content", got.Content)
	})

	t.Run("Delete", func(t *testing.T) {
		todo := model.NewTodo("To be deleted")
		s.Add(todo)

		err := s.Delete(todo.ID)
		assert.NoError(t, err)

		got, err := s.Get(todo.ID)
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
