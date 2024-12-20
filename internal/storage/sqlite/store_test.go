package sqlite

import (
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteStorage(t *testing.T) {
	// Setup test database
	log, _ := logger.NewZapLogger(&logger.Config{Level: "debug"})
	store, err := New(":memory:", log) // Use in-memory SQLite for testing
	require.NoError(t, err)
	defer store.Close()

	// Test saving a note
	note := "Test quick note"
	err = store.SaveNote(note)
	require.NoError(t, err)

	// Test retrieving notes
	notes, err := store.GetNotes()
	require.NoError(t, err)
	assert.Contains(t, notes, note)
}
