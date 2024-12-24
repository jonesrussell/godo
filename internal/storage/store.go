package storage

// Store defines the interface for storage operations
type Store interface {
	// SaveNote saves a note to storage
	SaveNote(note string) error

	// GetNotes retrieves all notes from storage
	GetNotes() ([]string, error)

	// DeleteNote removes a note from storage
	DeleteNote(note string) error

	// Clear removes all notes from storage
	Clear() error
}
