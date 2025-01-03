package testutil

import (
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// CreateTestNote creates a note for testing
func CreateTestNote(id, title string, completed bool) storage.Note {
	now := time.Now().Unix()
	return storage.Note{
		ID:        id,
		Title:     title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateTestNotes creates multiple notes for testing
func CreateTestNotes(count int) []storage.Note {
	notes := make([]storage.Note, count)
	for i := range notes {
		notes[i] = CreateTestNote(
			fmt.Sprintf("test-%d", i+1),
			fmt.Sprintf("Test Note %d", i+1),
			false,
		)
	}
	return notes
}
