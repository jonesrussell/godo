// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage/types"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

// Server represents the API server
type Server struct {
	store  types.Store
	logger logger.Logger
	srv    *http.Server
}

// NewServer creates a new API server
func NewServer(store types.Store, logger logger.Logger) *Server {
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
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	// Start server
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("Failed to start server", "error", err)
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop stops the API server
func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to stop server", "error", err)
		return fmt.Errorf("failed to stop server: %w", err)
	}
	return nil
}
