package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Integration(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	store, err := New(dbPath)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	t.Run("Add and Get Note", func(t *testing.T) {
		now := time.Now()
		n := &note.Note{
			ID:        "test-id",
			Content:   "test content",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, n)
		require.NoError(t, err)

		retrieved, err := store.Get(ctx, n.ID)
		require.NoError(t, err)
		assert.Equal(t, n.ID, retrieved.ID)
		assert.Equal(t, n.Content, retrieved.Content)
		assert.Equal(t, n.Completed, retrieved.Completed)
		assert.Equal(t, n.CreatedAt.Unix(), retrieved.CreatedAt.Unix())
		assert.Equal(t, n.UpdatedAt.Unix(), retrieved.UpdatedAt.Unix())
	})

	t.Run("List Notes", func(t *testing.T) {
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, notes)
	})

	t.Run("Update Note", func(t *testing.T) {
		now := time.Now()
		n := &note.Note{
			ID:        "update-test",
			Content:   "original content",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, n)
		require.NoError(t, err)

		n.Content = "updated content"
		n.Completed = true
		n.UpdatedAt = time.Now()

		err = store.Update(ctx, n)
		require.NoError(t, err)

		retrieved, err := store.Get(ctx, n.ID)
		require.NoError(t, err)
		assert.Equal(t, "updated content", retrieved.Content)
		assert.True(t, retrieved.Completed)
		assert.Equal(t, n.UpdatedAt.Unix(), retrieved.UpdatedAt.Unix())
	})

	t.Run("Delete Note", func(t *testing.T) {
		now := time.Now()
		n := &note.Note{
			ID:        "delete-test",
			Content:   "to be deleted",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, n)
		require.NoError(t, err)

		err = store.Delete(ctx, n.ID)
		require.NoError(t, err)

		_, err = store.Get(ctx, n.ID)
		assert.Error(t, err)
		var nErr *note.Error
		assert.ErrorAs(t, err, &nErr)
		assert.Equal(t, note.NotFound, nErr.Kind)
	})

	t.Run("Transaction", func(t *testing.T) {
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		now := time.Now()
		n := &note.Note{
			ID:        "tx-test",
			Content:   "transaction content",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = tx.Add(ctx, n)
		require.NoError(t, err)

		retrieved, err := tx.Get(ctx, n.ID)
		require.NoError(t, err)
		assert.Equal(t, n.ID, retrieved.ID)

		notes, err := tx.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, notes)

		n.Content = "updated in transaction"
		n.UpdatedAt = time.Now()
		err = tx.Update(ctx, n)
		require.NoError(t, err)

		err = tx.Delete(ctx, n.ID)
		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	})

	t.Run("Transaction Rollback", func(t *testing.T) {
		tx, err := store.BeginTx(ctx)
		require.NoError(t, err)

		now := time.Now()
		n := &note.Note{
			ID:        "rollback-test",
			Content:   "to be rolled back",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err = tx.Add(ctx, n)
		require.NoError(t, err)

		err = tx.Rollback()
		require.NoError(t, err)

		_, err = store.Get(ctx, n.ID)
		assert.Error(t, err)
		var nErr *note.Error
		assert.ErrorAs(t, err, &nErr)
		assert.Equal(t, note.NotFound, nErr.Kind)
	})

	t.Run("Empty Content Validation", func(t *testing.T) {
		now := time.Now()
		n := &note.Note{
			ID:        "invalid-test",
			Content:   "",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, n)
		assert.Error(t, err)
		var nErr *note.Error
		assert.ErrorAs(t, err, &nErr)
		assert.Equal(t, note.ValidationFailed, nErr.Kind)
	})
}
