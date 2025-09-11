// Package model defines the core domain models for the application
package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Note represents a todo item - the core domain model
type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewNote creates a new Note item
func NewNote(content string) *Note {
	now := time.Now()
	return &Note{
		ID:        uuid.New().String(),
		Content:   content,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToggleDone toggles the done status of the note
func (n *Note) ToggleDone() {
	n.Done = !n.Done
	n.UpdatedAt = time.Now()
}

// UpdateContent updates the content of the note
func (n *Note) UpdateContent(content string) {
	n.Content = content
	n.UpdatedAt = time.Now()
}

// MarkDone marks the note as done
func (n *Note) MarkDone() {
	n.Done = true
	n.UpdatedAt = time.Now()
}

// MarkUndone marks the note as not done
func (n *Note) MarkUndone() {
	n.Done = false
	n.UpdatedAt = time.Now()
}

// IsValid validates the note content
func (n *Note) IsValid() error {
	if n.Content == "" {
		return &ValidationError{
			Field:   "content",
			Message: "note content cannot be empty",
		}
	}
	if len(n.Content) > 1000 {
		return &ValidationError{
			Field:   "content",
			Message: "note content cannot exceed 1000 characters",
		}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return "validation error: " + e.Field + ": " + e.Message
}

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrDuplicateID  = errors.New("note ID already exists")
)
