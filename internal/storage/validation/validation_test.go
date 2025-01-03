package validation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/storage/types"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	ctx := context.Background()
	mockStore := testutil.NewMockStore()
	store := New(mockStore)

	t.Run("Add", func(t *testing.T) {
		// Valid note
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

		// Empty ID
		invalidNote := note
		invalidNote.ID = ""
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)

		// Empty content
		invalidNote = note
		invalidNote.Content = ""
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)

		// Invalid timestamps
		invalidNote = note
		invalidNote.CreatedAt = 0
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = 0
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = invalidNote.CreatedAt - 1
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = time.Now().Unix() + 3600
		err = store.Add(ctx, invalidNote)
		assert.Error(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		// Valid ID
		_, err := store.Get(ctx, uuid.New().String())
		assert.NoError(t, err)

		// Empty ID
		_, err = store.Get(ctx, "")
		assert.Error(t, err)
	})

	t.Run("Update", func(t *testing.T) {
		// Valid note
		now := time.Now().Unix()
		note := types.Note{
			ID:        uuid.New().String(),
			Content:   "Test Note",
			Completed: false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := store.Update(ctx, note)
		require.NoError(t, err)

		// Empty ID
		invalidNote := note
		invalidNote.ID = ""
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)

		// Empty content
		invalidNote = note
		invalidNote.Content = ""
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)

		// Invalid timestamps
		invalidNote = note
		invalidNote.CreatedAt = 0
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = 0
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = invalidNote.CreatedAt - 1
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)

		invalidNote = note
		invalidNote.UpdatedAt = time.Now().Unix() + 3600
		err = store.Update(ctx, invalidNote)
		assert.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		// Valid ID
		err := store.Delete(ctx, uuid.New().String())
		assert.NoError(t, err)

		// Empty ID
		err = store.Delete(ctx, "")
		assert.Error(t, err)
	})

	t.Run("List", func(t *testing.T) {
		_, err := store.List(ctx)
		assert.NoError(t, err)
	})
}
