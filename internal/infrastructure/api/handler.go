// Package api implements the HTTP server and API endpoints
package api

import (
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/domain/model"
)

// NoteHandler defines the interface for note-related HTTP handlers
type NoteHandler interface {
	// List returns all notes
	List(w http.ResponseWriter, r *http.Request)
	// Create creates a new note
	Create(w http.ResponseWriter, r *http.Request)
	// Get returns a specific note
	Get(w http.ResponseWriter, r *http.Request)
	// Update replaces an existing note
	Update(w http.ResponseWriter, r *http.Request)
	// Patch partially updates an existing note
	Patch(w http.ResponseWriter, r *http.Request)
	// Delete removes a note
	Delete(w http.ResponseWriter, r *http.Request)
}

// CreateNoteRequest represents a request to create a new note
type CreateNoteRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// UpdateNoteRequest represents a request to update an existing note
type UpdateNoteRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
	Done    bool   `json:"done"`
}

// PatchNoteRequest represents a request to partially update a note
type PatchNoteRequest struct {
	Content *string `json:"content,omitempty" validate:"omitempty,max=1000"`
	Done    *bool   `json:"done,omitempty"`
}

// NoteResponse represents a note in API responses
type NoteResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewNoteResponse creates a NoteResponse from a model.Note
func NewNoteResponse(note *model.Note) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		Content:   note.Content,
		Done:      note.Done,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

// NoteListResponse represents a list of notes in API responses
type NoteListResponse struct {
	Notes []NoteResponse `json:"notes"`
}

// NewNoteListResponse creates a NoteListResponse from a slice of model.Notes
func NewNoteListResponse(notes []model.Note) NoteListResponse {
	response := NoteListResponse{
		Notes: make([]NoteResponse, len(notes)),
	}
	for i, note := range notes {
		response.Notes[i] = NewNoteResponse(&note)
	}
	return response
}

// ErrorResponse represents an error in API responses
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// ValidationErrorResponse represents a validation error in API responses
type ValidationErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}
