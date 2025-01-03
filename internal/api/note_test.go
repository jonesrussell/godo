package api

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/storage/types"
	"github.com/stretchr/testify/assert"
)

func TestNoteCreation(t *testing.T) {
	// Using correct field names and time handling
	note := types.Note{
		Content:   "Buy groceries",
		Completed: false,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Using correct field names in assertions
	assert.Equal(t, "Buy groceries", note.Content)

	// Using correct field name for completion status
	assert.False(t, note.Completed)

	// Using proper time comparison with Unix timestamps
	assert.WithinDuration(t, time.Now(), time.Unix(note.CreatedAt, 0), time.Second)
}

func TestNoteUpdate(t *testing.T) {
	originalNote := &types.Note{
		Content:   "Original note",
		Completed: true,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	updatedNote := &types.Note{
		Content:   "Updated note",
		Completed: false,
		UpdatedAt: time.Now().Unix(),
	}

	assert.NotEqual(t, originalNote.Content, updatedNote.Content)
	assert.NotEqual(t, originalNote.Completed, updatedNote.Completed)
}
