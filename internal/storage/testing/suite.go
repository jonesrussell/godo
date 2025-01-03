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
		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("AddAndRetrieve", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, note.ID, notes[0].ID)
		assert.Equal(t, note.Content, notes[0].Content)
		assert.Equal(t, note.Completed, notes[0].Completed)
		assert.Equal(t, note.CreatedAt, notes[0].CreatedAt)
		assert.Equal(t, note.UpdatedAt, notes[0].UpdatedAt)
	})

	t.Run("Update", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Original Content",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		note.Content = "Updated Content"
		note.Completed = true
		note.UpdatedAt = time.Now().Unix()
		err = store.Update(ctx, note)
		assert.NoError(t, err)

		updated, err := store.Get(ctx, note.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Content", updated.Content)
		assert.True(t, updated.Completed)
		assert.Equal(t, note.UpdatedAt, updated.UpdatedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		store := s.NewStore()
		defer store.Close()

		ctx := context.Background()
		now := time.Now().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		err = store.Delete(ctx, note.ID)
		assert.NoError(t, err)

		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, notes)
	})
}
