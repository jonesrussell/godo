package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite" // SQLite driver
)

func TestSQLiteStore(t *testing.T) {
	log := logger.NewTestLogger(t)
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := New(dbPath)
	require.NoError(t, err)
	defer func() {
		store.Close()
		os.Remove(dbPath)
	}()

	now := time.Now().Unix()
	ctx := context.Background()

	t.Run("Add and List", func(t *testing.T) {
		note := storage.Note{
			ID:        "1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Add(ctx, note)
		assert.NoError(t, err)

		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)

		// Compare fields individually
		assert.Equal(t, note.ID, notes[0].ID)
		assert.Equal(t, note.Content, notes[0].Content)
		assert.Equal(t, note.Completed, notes[0].Completed)
		assert.Equal(t, note.CreatedAt, notes[0].CreatedAt)
		assert.Equal(t, note.UpdatedAt, notes[0].UpdatedAt)
	})

	t.Run("Update", func(t *testing.T) {
		note := storage.Note{
			ID:        "1",
			Content:   "Updated Note",
			Completed: true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Update(ctx, note)
		assert.NoError(t, err)

		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Note", notes[0].Content)
		assert.True(t, notes[0].Completed)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Delete(ctx, "1")
		assert.NoError(t, err)

		notes, err := store.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, notes)

		err = store.Delete(ctx, "1")
		assert.ErrorIs(t, err, errors.ErrNoteNotFound)
	})
}
