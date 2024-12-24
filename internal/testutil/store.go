package testutil

import (
	"testing"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RunStoreTests runs a standard suite of tests against any Store implementation
func RunStoreTests(t *testing.T, s storage.Store) {
	t.Helper()

	// Clear any existing data
	require.NoError(t, s.Clear())

	t.Run("SaveNote", func(t *testing.T) {
		err := s.SaveNote("test note")
		assert.NoError(t, err)
	})

	t.Run("GetNotes", func(t *testing.T) {
		notes, err := s.GetNotes()
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, "test note", notes[0])
	})

	t.Run("DeleteNote", func(t *testing.T) {
		err := s.DeleteNote("test note")
		assert.NoError(t, err)

		notes, err := s.GetNotes()
		assert.NoError(t, err)
		assert.Empty(t, notes)

		err = s.DeleteNote("non-existent note")
		assert.ErrorIs(t, err, storage.ErrTodoNotFound)
	})

	t.Run("Clear", func(t *testing.T) {
		// Add some notes
		require.NoError(t, s.SaveNote("note 1"))
		require.NoError(t, s.SaveNote("note 2"))

		// Clear them
		err := s.Clear()
		assert.NoError(t, err)

		// Verify they're gone
		notes, err := s.GetNotes()
		assert.NoError(t, err)
		assert.Empty(t, notes)
	})
}
