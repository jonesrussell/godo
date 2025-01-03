// Package api implements the HTTP server and API endpoints
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/domain/note"
	"github.com/jonesrussell/godo/internal/storage/types"
)

// NoteHandler handles HTTP requests for notes
type NoteHandler struct {
	store types.Store
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(store types.Store) *NoteHandler {
	return &NoteHandler{
		store: store,
	}
}

// CreateNote handles note creation requests
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note note.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.CreatedAt = time.Now().Unix()
	note.UpdatedAt = time.Now().Unix()

	if err := h.store.Add(r.Context(), note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetNote handles note retrieval requests
func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing note id", http.StatusBadRequest)
		return
	}

	note, err := h.store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

// UpdateNote handles note update requests
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var note note.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.UpdatedAt = time.Now().Unix()

	if err := h.store.Update(r.Context(), note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

// DeleteNote handles note deletion requests
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing note id", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListNotes handles note listing requests
func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(notes)
}
