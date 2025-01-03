package testutil

import (
	"fmt"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// CreateTestNote creates a test note with the given parameters
func CreateTestNote(id, content string, completed bool) storage.Note {
	now := time.Now().Unix()
	return storage.Note{
		ID:        id,
		Content:   content,
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
