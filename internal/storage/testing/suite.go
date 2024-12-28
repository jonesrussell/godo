// Package testing provides test utilities for storage implementations
package testing

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StoreSuite provides a suite of tests that can be run against any Store implementation
type StoreSuite struct {
	NewStore func() storage.Store
}

// Run executes all test cases in the suite
func (s *StoreSuite) Run(t *testing.T) {
	t.Run("EmptyStore", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("AddAndRetrieve", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Add(task)
		require.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Content, tasks[0].Content)
		assert.Equal(t, task.Done, tasks[0].Done)
	})

	t.Run("Update", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		task := storage.Task{
			ID:        "test-1",
			Content:   "Original Content",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Add(task)
		require.NoError(t, err)

		task.Content = "Updated Content"
		task.Done = true
		err = store.Update(task)
		assert.NoError(t, err)

		updated, err := store.GetByID(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Content", updated.Content)
		assert.True(t, updated.Done)
	})

	t.Run("Delete", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Add(task)
		require.NoError(t, err)

		err = store.Delete(task.ID)
		assert.NoError(t, err)

		tasks, err := store.List()
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("GetByID", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		task := storage.Task{
			ID:        "test-1",
			Content:   "Test Task",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Add(task)
		require.NoError(t, err)

		retrieved, err := store.GetByID(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, retrieved.ID)
		assert.Equal(t, task.Content, retrieved.Content)

		_, err = store.GetByID("nonexistent")
		assert.Error(t, err)
	})
}
