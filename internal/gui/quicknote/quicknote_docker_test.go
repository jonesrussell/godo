//go:build docker && !windows
// +build docker,!windows

package quicknote

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerQuickNote(t *testing.T) {
	store := mock.New()
	log := logger.NewTestLogger(t)
	quickNote := New(store, log)
	require.NotNil(t, quickNote)

	t.Run("AddNote", func(t *testing.T) {
		// Add a note
		note := storage.Note{
			ID:        "test-1",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		err := store.Add(note)
		require.NoError(t, err)

		// Verify the note is added
		assert.True(t, store.AddCalled)
	})

	t.Run("StoreError", func(t *testing.T) {
		// Set store error
		store.Error = assert.AnError

		// Try to add a note
		note := storage.Note{
			ID:        "test-2",
			Content:   "Test Note",
			Completed: false,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		err := store.Add(note)
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}
