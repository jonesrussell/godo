// Package testing provides test utilities for storage implementations
package testing

import (
	"context"
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

		ctx := context.Background()
		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})

	t.Run("AddAndRetrieve", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, task)
		require.NoError(t, err)

		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, task.ID, tasks[0].ID)
		assert.Equal(t, task.Title, tasks[0].Title)
		assert.Equal(t, task.Completed, tasks[0].Completed)
		assert.Equal(t, task.CreatedAt, tasks[0].CreatedAt)
		assert.Equal(t, task.UpdatedAt, tasks[0].UpdatedAt)
	})

	t.Run("Update", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Original Title",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, task)
		require.NoError(t, err)

		task.Title = "Updated Title"
		task.Completed = true
		task.UpdatedAt = time.Now().Unix()
		err = store.Update(ctx, task)
		assert.NoError(t, err)

		updated, err := store.Get(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", updated.Title)
		assert.True(t, updated.Completed)
		assert.Equal(t, task.UpdatedAt, updated.UpdatedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		task := storage.Task{
			ID:        "test-1",
			Title:     "Test Task",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, task)
		require.NoError(t, err)

		err = store.Delete(ctx, task.ID)
		assert.NoError(t, err)

		tasks, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, tasks)
	})
}
