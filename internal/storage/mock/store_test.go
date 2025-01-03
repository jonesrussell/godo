package mock

import (
	"context"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("Add and Get", func(t *testing.T) {
		note := storage.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err := store.Add(ctx, note)
		require.NoError(t, err)

		got, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note, got)
	})

	t.Run("List", func(t *testing.T) {
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 1)
	})

	t.Run("Update", func(t *testing.T) {
		note := storage.Note{
			ID:        "1",
			Content:   "Updated Note",
			Completed: true,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err := store.Update(ctx, note)
		require.NoError(t, err)

		got, err := store.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.Content, got.Content)
		assert.Equal(t, note.Completed, got.Completed)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, "1")
		require.NoError(t, err)

		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("Error handling", func(t *testing.T) {
		store.SetError(assert.AnError)

		note := storage.Note{ID: "2"}
		err := store.Add(ctx, note)
		assert.Error(t, err)

		_, err = store.Get(ctx, "2")
		assert.Error(t, err)

		_, err = store.List(ctx)
		assert.Error(t, err)

		err = store.Update(ctx, note)
		assert.Error(t, err)

		err = store.Delete(ctx, "2")
		assert.Error(t, err)

		err = store.Close()
		assert.Error(t, err)
	})
}

func TestTransaction(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("Successful transaction", func(t *testing.T) {
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		note := storage.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err = tx.Add(ctx, note)
		require.NoError(t, err)

		got, err := tx.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note, got)

		notes, err := tx.List(ctx)
		require.NoError(t, err)
		assert.Len(t, notes, 1)

		note.Content = "Updated Note"
		err = tx.Update(ctx, note)
		require.NoError(t, err)

		got, err = tx.Get(ctx, note.ID)
		require.NoError(t, err)
		assert.Equal(t, note.Content, got.Content)

		err = tx.Delete(ctx, note.ID)
		require.NoError(t, err)

		notes, err = tx.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)

		err = tx.Commit()
		require.NoError(t, err)
	})

	t.Run("Rollback transaction", func(t *testing.T) {
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		note := storage.Note{
			ID:        "2",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		err = tx.Add(ctx, note)
		require.NoError(t, err)

		err = tx.Rollback()
		require.NoError(t, err)

		// Verify note was not added to store
		_, err = store.Get(ctx, note.ID)
		assert.Error(t, err)
	})

	t.Run("Transaction error handling", func(t *testing.T) {
		store.SetError(assert.AnError)

		_, err := store.BeginTx(ctx)
		assert.Error(t, err)
	})
}
