package mock

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("Add", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Verify note was added
		addedNote, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.ID, addedNote.ID)
		assert.Equal(t, note.Content, addedNote.Content)
		assert.Equal(t, note.Completed, addedNote.Completed)
		assert.Equal(t, note.CreatedAt, addedNote.CreatedAt)
		assert.Equal(t, note.UpdatedAt, addedNote.UpdatedAt)

		// Try to add same note again
		err = store.Add(ctx, note)
		assert.Error(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Update the note
		note.Content = "Updated Note"
		note.Completed = true
		note.UpdatedAt = time.Now().Unix()

		err = store.Update(ctx, note)
		require.NoError(t, err)

		// Verify note was updated
		updatedNote, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.Content, updatedNote.Content)
		assert.Equal(t, note.Completed, updatedNote.Completed)
		assert.Equal(t, note.UpdatedAt, updatedNote.UpdatedAt)

		// Try to update non-existent note
		nonExistentNote := types.Note{
			ID:      uuid.New().String(),
			Content: "Non-existent Note",
		}
		err = store.Update(ctx, nonExistentNote)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Delete the note
		err = store.Delete(ctx, note.ID)
		require.NoError(t, err)

		// Verify note was deleted
		_, err = store.Get(ctx, note.ID)
		assert.Error(t, err)

		// Try to delete non-existent note
		err = store.Delete(ctx, uuid.New().String())
		assert.Error(t, err)
	})

	t.Run("List", func(t *testing.T) {
		// Add some notes
		now := time.Now().Unix()
		note1 := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note 1",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, note1)
		require.NoError(t, err)

		note2 := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note 2",
			Completed: true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Add(ctx, note2)
		require.NoError(t, err)

		// List notes
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 2)

		// Verify notes are in list
		var found1, found2 bool
		for _, note := range notes {
			if note.ID == note1.ID {
				found1 = true
				assert.Equal(t, note1.Content, note.Content)
				assert.Equal(t, note1.Completed, note.Completed)
			}
			if note.ID == note2.ID {
				found2 = true
				assert.Equal(t, note2.Content, note.Content)
				assert.Equal(t, note2.Completed, note.Completed)
			}
		}
		assert.True(t, found1)
		assert.True(t, found2)
	})

	t.Run("Error", func(t *testing.T) {
		testErr := assert.AnError
		store.SetError(testErr)

		// Test all operations return the error
		_, err := store.List(ctx)
		assert.Equal(t, testErr, err)

		_, err = store.Get(ctx, "test")
		assert.Equal(t, testErr, err)

		err = store.Add(ctx, types.Note{})
		assert.Equal(t, testErr, err)

		err = store.Update(ctx, types.Note{})
		assert.Equal(t, testErr, err)

		err = store.Delete(ctx, "test")
		assert.Equal(t, testErr, err)

		err = store.Close()
		assert.Equal(t, testErr, err)

		_, err = store.BeginTx(ctx)
		assert.Equal(t, testErr, err)

		// Reset error and verify operations work again
		store.Reset()
		_, err = store.List(ctx)
		assert.NoError(t, err)
	})
}
