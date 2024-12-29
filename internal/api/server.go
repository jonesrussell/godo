// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

const (
	// DefaultReadTimeout is the default timeout for reading the entire request
	DefaultReadTimeout = 30 * time.Second
	// DefaultWriteTimeout is the default timeout for writing the response
	DefaultWriteTimeout = 30 * time.Second
	// DefaultReadHeaderTimeout is the default timeout for reading request headers
	DefaultReadHeaderTimeout = 10 * time.Second
	// DefaultIdleTimeout is the default timeout for idle connections
	DefaultIdleTimeout = 120 * time.Second
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
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		IdleTimeout:       DefaultIdleTimeout,
	}
}

// Server represents the HTTP server
type Server struct {
	store  storage.TaskStore
	log    logger.Logger
	router *mux.Router
	srv    *http.Server
}

// NewServer creates a new Server instance
func NewServer(store storage.TaskStore, log logger.Logger) *Server {
	s := &Server{
		store:  store,
		log:    log,
		router: mux.NewRouter(),
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

	// Tasks
	api.HandleFunc("/tasks", Chain(s.handleListTasks,
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/tasks", Chain(s.handleCreateTask,
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[CreateTaskRequest](s.log),
	)).Methods(http.MethodPost)

	api.HandleFunc("/tasks/{id}", Chain(s.handleGetTask,
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/tasks/{id}", Chain(s.handleUpdateTask,
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[UpdateTaskRequest](s.log),
	)).Methods(http.MethodPut)

	api.HandleFunc("/tasks/{id}", Chain(s.handlePatchTask,
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[PatchTaskRequest](s.log),
	)).Methods(http.MethodPatch)

	api.HandleFunc("/tasks/{id}", Chain(s.handleDeleteTask,
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodDelete)
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.store.List(r.Context())
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewTaskListResponse(tasks))
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	req, ok := GetRequest[CreateTaskRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	task := storage.Task{
		ID:        uuid.New().String(),
		Content:   req.Content,
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.store.Add(r.Context(), task); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, NewTaskResponse(task))
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := s.store.GetByID(r.Context(), id)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewTaskResponse(task))
}

func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, ok := GetRequest[UpdateTaskRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	task := storage.Task{
		ID:        id,
		Content:   req.Content,
		Done:      req.Done,
		UpdatedAt: time.Now(),
	}

	if err := s.store.Update(r.Context(), task); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewTaskResponse(task))
}

func (s *Server) handlePatchTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, ok := GetRequest[PatchTaskRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	task, err := s.store.GetByID(r.Context(), id)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	if req.Content != nil {
		task.Content = *req.Content
	}
	if req.Done != nil {
		task.Done = *req.Done
	}
	task.UpdatedAt = time.Now()

	if err := s.store.Update(r.Context(), task); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewTaskResponse(task))
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.store.Delete(r.Context(), id); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Start starts the server
func (s *Server) Start(port int) error {
	s.srv = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s,
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		IdleTimeout:       DefaultIdleTimeout,
	}

	s.log.Info("starting server", "port", port)
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}
