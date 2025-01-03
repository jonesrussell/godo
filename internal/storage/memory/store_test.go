package memory

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	store := New()

	t.Run("empty store", func(t *testing.T) {
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("add note", func(t *testing.T) {
		now := time.Now().UTC().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Verify note was added
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, note.ID, notes[0].ID)
		assert.Equal(t, note.Content, notes[0].Content)
		assert.Equal(t, note.Completed, notes[0].Completed)
		assert.Equal(t, note.CreatedAt, notes[0].CreatedAt)
		assert.Equal(t, note.UpdatedAt, notes[0].UpdatedAt)
	})

	t.Run("get note", func(t *testing.T) {
		note, err := store.Get(ctx, "test-1")
		require.NoError(t, err)
		assert.Equal(t, "test-1", note.ID)
		assert.Equal(t, "Test Note", note.Content)

		// Try to get nonexistent note
		_, err = store.Get(ctx, "nonexistent")
		assert.ErrorIs(t, err, errors.ErrNoteNotFound)
	})

	t.Run("update note", func(t *testing.T) {
		now := time.Now().UTC().Unix()
		note := storage.Note{
			ID:        "test-1",
			Content:   "Updated Note",
			Completed: true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Update(ctx, note)
		require.NoError(t, err)

		// Verify note was updated
		updated, err := store.Get(ctx, "test-1")
		require.NoError(t, err)
		assert.Equal(t, note.Content, updated.Content)
		assert.Equal(t, note.Completed, updated.Completed)
		assert.Equal(t, note.UpdatedAt, updated.UpdatedAt)

		// Try to update nonexistent note
		note.ID = "nonexistent"
		err = store.Update(ctx, note)
		assert.ErrorIs(t, err, errors.ErrNoteNotFound)
	})

	t.Run("delete note", func(t *testing.T) {
		err := store.Delete(ctx, "test-1")
		require.NoError(t, err)

		// Verify note was deleted
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)

		// Try to delete nonexistent note
		err = store.Delete(ctx, "test-1")
		assert.ErrorIs(t, err, errors.ErrNoteNotFound)
	})
}

func TestTransaction(t *testing.T) {
	ctx := context.Background()
	store := New()

	t.Run("commit", func(t *testing.T) {
		// Add initial note
		now := time.Now().UTC().Unix()
		note1 := storage.Note{
			ID:        "test-1",
			Content:   "Test Note 1",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, store.Add(ctx, note1))

		// Start transaction
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		// Add second note in transaction
		note2 := storage.Note{
			ID:        "test-2",
			Content:   "Test Note 2",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, tx.Add(ctx, note2))

		// Update first note in transaction
		note1.Content = "Updated Note 1"
		require.NoError(t, tx.Update(ctx, note1))

		// Verify changes are not visible outside transaction
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, "Test Note 1", notes[0].Content)

		// Commit transaction
		require.NoError(t, tx.Commit())

		// Verify changes are now visible
		notes, err = store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 2)
		// Find and verify the updated note
		var found bool
		for _, note := range notes {
			if note.ID == "test-1" {
				assert.Equal(t, "Updated Note 1", note.Content)
				found = true
				break
			}
		}
		assert.True(t, found, "Updated note should be found")
	})

	t.Run("rollback", func(t *testing.T) {
		// Start transaction
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		// Delete all notes in transaction
		notes, err := tx.List(ctx)
		require.NoError(t, err)
		for _, note := range notes {
			require.NoError(t, tx.Delete(ctx, note.ID))
		}

		// Add new note in transaction
		now := time.Now().UTC().Unix()
		note := storage.Note{
			ID:        "test-3",
			Content:   "Test Note 3",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, tx.Add(ctx, note))

		// Rollback transaction
		require.NoError(t, tx.Rollback())

		// Verify original state is preserved
		notes, err = store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 2)
	})
}
