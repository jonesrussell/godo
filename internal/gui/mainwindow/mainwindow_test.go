package mainwindow

import (
	"context"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/storage/types"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindow(t *testing.T) {
	ctx := context.Background()
	store := testutil.NewMockStore()

	window := test.NewWindow(nil)
	defer window.Close()

	t.Run("AddNote", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Get the note to verify it was added
		addedNote, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.ID, addedNote.ID)
		assert.Equal(t, note.Content, addedNote.Content)
		assert.Equal(t, note.Completed, addedNote.Completed)
		assert.Equal(t, note.CreatedAt, addedNote.CreatedAt)
		assert.Equal(t, note.UpdatedAt, addedNote.UpdatedAt)
	})

	t.Run("UpdateNote", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        "test-2",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		// Update the note
		note.Completed = true
		note.UpdatedAt = time.Now().Unix()
		err = store.Update(ctx, note)
		require.NoError(t, err)

		// Get the note to verify it was updated
		updatedNote, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.True(t, updatedNote.Completed)
		assert.Equal(t, note.UpdatedAt, updatedNote.UpdatedAt)
	})

	t.Run("DeleteNote", func(t *testing.T) {
		// Add a note
		now := time.Now().Unix()
		note := types.Note{
			ID:        "test-3",
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

		// Verify the note is deleted
		_, err = store.Get(ctx, note.ID)
		assert.Error(t, err)
	})

	t.Run("ListNotes", func(t *testing.T) {
		// Add some notes
		now := time.Now().Unix()
		note1 := types.Note{
			ID:        "test-4",
			Content:   "Test Note 1",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err := store.Add(ctx, note1)
		require.NoError(t, err)

		note2 := types.Note{
			ID:        "test-5",
			Content:   "Test Note 2",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		err = store.Add(ctx, note2)
		require.NoError(t, err)

		// List notes
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 2)

		// Verify note order and fields
		assert.Equal(t, note1.ID, notes[0].ID)
		assert.Equal(t, note1.Content, notes[0].Content)
		assert.Equal(t, note1.Completed, notes[0].Completed)
		assert.Equal(t, note2.ID, notes[1].ID)
		assert.Equal(t, note2.Content, notes[1].Content)
		assert.Equal(t, note2.Completed, notes[1].Completed)
	})
}
