package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/stretchr/testify/require"
)

func TestSQLiteStore(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	store, err := New(dbPath)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	t.Run("Add and Get Note", func(t *testing.T) {
		n := &note.Note{
			ID:        "test-id",
			Content:   "test content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := store.Add(ctx, n)
		require.NoError(t, err)

		retrieved, err := store.Get(ctx, n.ID)
		require.NoError(t, err)
		require.Equal(t, n.ID, retrieved.ID)
		require.Equal(t, n.Content, retrieved.Content)
		require.WithinDuration(t, n.CreatedAt, retrieved.CreatedAt, time.Second)
	})

	// ... rest of tests ...
}
