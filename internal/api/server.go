// Package api implements the HTTP server and API endpoints
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/service"
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
	store   storage.TaskStore
	service service.TaskService
	log     logger.Logger
	router  *mux.Router
	srv     *http.Server
}

// NewServer creates a new Server instance
func NewServer(store storage.TaskStore, taskService service.TaskService, log logger.Logger) *Server {
	s := &Server{
		store:   store,
		service: taskService,
		log:     log,
		router:  mux.NewRouter(),
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

	// Protected Task endpoints (JWT auth required)
	api.HandleFunc("/tasks", Chain(s.handleListTasks,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/tasks", Chain(s.handleCreateTask,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[CreateTaskRequest](s.log),
	)).Methods(http.MethodPost)

	api.HandleFunc("/tasks/{id}", Chain(s.handleGetTask,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodGet)

	api.HandleFunc("/tasks/{id}", Chain(s.handleUpdateTask,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[UpdateTaskRequest](s.log),
	)).Methods(http.MethodPut)

	api.HandleFunc("/tasks/{id}", Chain(s.handlePatchTask,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
		WithValidation[PatchTaskRequest](s.log),
	)).Methods(http.MethodPatch)

	api.HandleFunc("/tasks/{id}", Chain(s.handleDeleteTask,
		WithJWTAuth(s.log),
		WithLogging(s.log),
		WithErrorHandling(s.log),
	)).Methods(http.MethodDelete)
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.service.ListTasks(r.Context(), nil)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	// Convert service tasks back to storage tasks for response
	storageTasks := make([]storage.Task, len(tasks))
	for i, task := range tasks {
		storageTasks[i] = *task
	}

	writeJSON(w, http.StatusOK, NewTaskListResponse(storageTasks))
}

func (s *Server) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	req, ok := GetRequest[CreateTaskRequest](r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request")
		return
	}

	task, err := s.service.CreateTask(r.Context(), req.Content)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusCreated, NewTaskResponse(task))
}

func (s *Server) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	task, err := s.service.GetTask(r.Context(), id)
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

	updates := service.TaskUpdateRequest{
		Content: &req.Content,
		Done:    &req.Done,
	}

	task, err := s.service.UpdateTask(r.Context(), id, updates)
	if err != nil {
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

	updates := service.TaskUpdateRequest{
		Content: req.Content,
		Done:    req.Done,
	}

	task, err := s.service.UpdateTask(r.Context(), id, updates)
	if err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, NewTaskResponse(task))
}

func (s *Server) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.service.DeleteTask(r.Context(), id); err != nil {
		status, code, msg := mapError(err)
		writeError(w, status, code, msg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, health)
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
