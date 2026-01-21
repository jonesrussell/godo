// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/service"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
}

// NewServerConfig creates a new server configuration with default values
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// Server represents the HTTP server
type Server struct {
	service   service.NoteService
	log       logger.Logger
	router    *mux.Router
	srv       *http.Server
	jwtSecret string
}

// NewServer creates a new Server instance
func NewServer(noteService service.NoteService, log logger.Logger, jwtSecret string) *Server {
	s := &Server{
		service:   noteService,
		log:       log,
		router:    mux.NewRouter(),
		jwtSecret: jwtSecret,
	}
	s.routes()
	return s
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// routes sets up the server routes
func (s *Server) routes() {
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint (no auth required)
	api.HandleFunc("/health", Chain(s.handleHealth,
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	// Protected Note endpoints (JWT auth required)
	api.HandleFunc("/notes", Chain(s.handleListNotes,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/notes", Chain(s.handleCreateNote,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[CreateNoteRequest](s.log),
	)).Methods(http.MethodPost)

	api.HandleFunc("/notes/{id}", Chain(s.handleGetNote,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/notes/{id}", Chain(s.handleUpdateNote,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[UpdateNoteRequest](s.log),
	)).Methods(http.MethodPut)

	api.HandleFunc("/notes/{id}", Chain(s.handlePatchNote,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[PatchNoteRequest](s.log),
	)).Methods(http.MethodPatch)

	api.HandleFunc("/notes/{id}", Chain(s.handleDeleteNote,
		WithJWTAuth(s.log, s.jwtSecret),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodDelete)
}

func (s *Server) handleListNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := s.service.ListNotes(r.Context(), nil)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	// Convert service notes to model notes for response
	modelNotes := make([]model.Note, len(notes))
	for i, note := range notes {
		modelNotes[i] = *note
	}

	writeJSON(w, http.StatusOK, NewNoteListResponse(modelNotes))
}

func (s *Server) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	req, ok := GetRequest[CreateNoteRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	note, err := s.service.CreateNote(r.Context(), req.Content)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, NewNoteResponse(note))
}

func (s *Server) handleGetNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	note, err := s.service.GetNote(r.Context(), id)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewNoteResponse(note))
}

func (s *Server) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, ok := GetRequest[UpdateNoteRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	updates := service.NoteUpdateRequest{
		Content: &req.Content,
		Done:    &req.Done,
	}

	note, err := s.service.UpdateNote(r.Context(), id, updates)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewNoteResponse(note))
}

func (s *Server) handlePatchNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, ok := GetRequest[PatchNoteRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	updates := service.NoteUpdateRequest{
		Content: req.Content,
		Done:    req.Done,
	}

	note, err := s.service.UpdateNote(r.Context(), id, updates)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewNoteResponse(note))
}

func (s *Server) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.service.DeleteNote(r.Context(), id); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, response)
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	s.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	s.log.Info("Starting HTTP server", "port", port)

	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		// Provide a clearer error message for port conflicts
		if strings.Contains(err.Error(), "address already in use") {
			return fmt.Errorf("HTTP server failed to start: port %d is already in use. "+
				"Please configure a different port via environment variable GODO_HTTP_PORT, "+
				"config file, or CLI flag", port)
		}
		return fmt.Errorf("HTTP server failed to start: %w", err)
	}

	return err
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		s.log.Info("Shutting down HTTP server")
		return s.srv.Shutdown(ctx)
	}
	return nil
}
