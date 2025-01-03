package quicknote

import (
	"context"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/jonesrussell/godo/internal/config"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindow(t *testing.T) {
	ctx := context.Background()
	store := storage.NewMockStore()
	logger := logger.NewMockTestLogger(t)
	cfg := config.WindowConfig{
		Width:       400,
		Height:      300,
		StartHidden: true,
	}

	app := test.NewApp()
	quickNote := New(app, store, logger, cfg)
	require.NotNil(t, quickNote)

	t.Run("SaveNote", func(t *testing.T) {
		// Set input text
		quickNote.input.SetText("Test Note")

		// Save the note
		quickNote.saveBtn.OnTapped()

		// Verify note was saved
		notes, err := store.List(ctx)
		require.NoError(t, err)
		require.Len(t, notes, 1)
		assert.Equal(t, "Test Note", notes[0].Content)
		assert.False(t, notes[0].Completed)
		assert.NotZero(t, notes[0].CreatedAt)
		assert.NotZero(t, notes[0].UpdatedAt)
		assert.Equal(t, notes[0].CreatedAt, notes[0].UpdatedAt)
	})

	t.Run("SaveEmptyNote", func(t *testing.T) {
		// Try to save empty note
		quickNote.input.SetText("")
		quickNote.saveBtn.OnTapped()

		// Verify no note was saved
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)
	})

	t.Run("SaveWithError", func(t *testing.T) {
		// Set up store to return error
		store.SetError(assert.AnError)

		// Try to save note
		quickNote.input.SetText("Test Note")
		quickNote.saveBtn.OnTapped()

		// Verify input text remains
		assert.Equal(t, "Test Note", quickNote.input.Text)

		// Verify no note was saved
		notes, err := store.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, notes)

		// Reset store error
		store.SetError(nil)
	})
}
