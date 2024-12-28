// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

const (
	// Server timeouts
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
)

// Server represents the HTTP server
type Server struct {
	router *chi.Mux
	server *http.Server
	store  storage.Store
	logger logger.Logger
}

// NewServer creates a new HTTP server instance
func NewServer(store storage.Store, l logger.Logger) *Server {
	s := &Server{
		router: chi.NewRouter(),
		store:  store,
		logger: l,
	}

	// Set up middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// Set up routes
	s.router.Get("/health", s.handleHealth)
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Get("/tasks", s.handleListTasks)
		r.Post("/tasks", s.handleCreateTask)
		r.Put("/tasks/{id}", s.handleUpdateTask)
		r.Delete("/tasks/{id}", s.handleDeleteTask)
	})

	return s
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	s.logger.Info("Starting HTTP server", "port", port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.List()
	if err != nil {
		s.logger.Error("Failed to list tasks", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, tasks)
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate UUID and set timestamps
	now := time.Now()
	task.ID = uuid.New().String()
	task.CreatedAt = now
	task.UpdatedAt = now

	if err := s.store.Add(task); err != nil {
		s.logger.Error("Failed to create task", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, task)
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.ID = id
	if err := s.store.Update(task); err != nil {
		if err == storage.ErrTaskNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		s.logger.Error("Failed to update task", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, task)
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.store.Delete(id); err != nil {
		if err == storage.ErrTaskNotFound {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		s.logger.Error("Failed to delete task", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
