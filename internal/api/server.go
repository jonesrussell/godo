// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jonesrussell/godo/internal/storage"
	"go.uber.org/zap"
)

// Server represents the API server
type Server struct {
	store  storage.Store
	logger *zap.Logger
	srv    *http.Server
}

// NewServer creates a new API server
func NewServer(store storage.Store, logger *zap.Logger) *Server {
	return &Server{
		store:  store,
		logger: logger,
	}
}

// Start starts the API server
func (s *Server) Start(addr string) error {
	router := mux.NewRouter()

	// Create handlers
	handler := NewNoteHandler(s.store)

	// Register routes
	router.HandleFunc("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListNotes(w, r)
		case http.MethodPost:
			handler.CreateNote(w, r)
		case http.MethodPut:
			handler.UpdateNote(w, r)
		case http.MethodDelete:
			handler.DeleteNote(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create server
	s.srv = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop stops the API server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}
	return nil
}

// NoteHandler handles note-related requests
type NoteHandler struct {
	store storage.Store
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(store storage.Store) *NoteHandler {
	return &NoteHandler{store: store}
}

// ListNotes handles GET /api/v1/notes
func (h *NoteHandler) ListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// CreateNote handles POST /api/v1/notes
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note storage.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Add(r.Context(), note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// UpdateNote handles PUT /api/v1/notes/:id
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	var note storage.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Update(r.Context(), note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// DeleteNote handles DELETE /api/v1/notes/:id
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.store.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
