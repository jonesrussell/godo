// Package note defines the core domain types and interfaces for note management
package note

import "time"

// Note represents a task or quick note in the system
type Note struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks if the note is valid
func (n *Note) Validate() error {
	if n.Content == "" {
		return &Error{
			Op:   "Note.Validate",
			Kind: ValidationFailed,
			Msg:  "content cannot be empty",
		}
	}
	return nil
}

// MarkComplete marks the note as completed
func (n *Note) MarkComplete() {
	n.Completed = true
	n.UpdatedAt = time.Now()
}

// MarkIncomplete marks the note as incomplete
func (n *Note) MarkIncomplete() {
	n.Completed = false
	n.UpdatedAt = time.Now()
}

// UpdateContent updates the note content
func (n *Note) UpdateContent(content string) error {
	if content == "" {
		return &Error{
			Op:   "Note.UpdateContent",
			Kind: ValidationFailed,
			Msg:  "content cannot be empty",
		}
	}
	n.Content = content
	n.UpdatedAt = time.Now()
	return nil
}
