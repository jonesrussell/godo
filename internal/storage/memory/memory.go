package memory

import (
	"github.com/jonesrussell/godo/internal/storage"
)

// Store is an in-memory implementation of storage.Store
type Store struct {
	notes []string
}

// New creates a new in-memory store
func New() *Store {
	return &Store{
		notes: make([]string, 0),
	}
}

// SaveNote adds a note to the store
func (s *Store) SaveNote(note string) error {
	s.notes = append(s.notes, note)
	return nil
}

// GetNotes returns all notes in the store
func (s *Store) GetNotes() ([]string, error) {
	// Return a copy to prevent modification of internal state
	result := make([]string, len(s.notes))
	copy(result, s.notes)
	return result, nil
}

// DeleteNote removes a note from the store
func (s *Store) DeleteNote(note string) error {
	for i, n := range s.notes {
		if n == note {
			// Remove the note by slicing
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return nil
		}
	}
	return storage.ErrTodoNotFound
}

// Clear removes all notes from the store
func (s *Store) Clear() error {
	s.notes = make([]string, 0)
	return nil
}
