// Package api implements the HTTP server and API endpoints
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	store storage.Store
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(store storage.Store) *TaskHandler {
	return &TaskHandler{
		store: store,
	}
}

// CreateTask handles task creation requests
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	if err := h.store.Add(r.Context(), task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTask handles task retrieval requests
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing task id", http.StatusBadRequest)
		return
	}

	task, err := h.store.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// UpdateTask handles task update requests
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.UpdatedAt = time.Now().Unix()

	if err := h.store.Update(r.Context(), task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(task)
}

// DeleteTask handles task deletion requests
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing task id", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTasks handles task listing requests
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}
