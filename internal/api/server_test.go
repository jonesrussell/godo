package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*Server, *testutil.MockStore) {
	t.Helper()

	// Create logger
	logConfig := &common.LogConfig{
		Level:       "debug",
		Output:      []string{"stdout"},
		ErrorOutput: []string{"stderr"},
	}
	log, err := logger.New(logConfig)
	require.NoError(t, err)

	// Create mock store
	store := testutil.NewMockStore()

	// Create server
	server := NewServer(store, log)
	require.NotNil(t, server)

	return server, store
}

func TestHealthCheck(t *testing.T) {
	server, _ := setupTestServer(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.handleHealth(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestListTasks(t *testing.T) {
	server, store := setupTestServer(t)

	// Add a test task
	task := storage.Task{
		ID:    "1",
		Title: "Test Task",
	}
	err := store.Add(task)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	server.handleListTasks(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []storage.Task
	err = json.NewDecoder(w.Body).Decode(&tasks)
	require.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, task, tasks[0])
}

func TestCreateTask(t *testing.T) {
	server, _ := setupTestServer(t)

	task := storage.Task{
		Title: "New Task",
	}
	body, err := json.Marshal(task)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.handleCreateTask(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdTask storage.Task
	err = json.NewDecoder(w.Body).Decode(&createdTask)
	require.NoError(t, err)
	assert.Equal(t, task.Title, createdTask.Title)
}

func TestUpdateTask(t *testing.T) {
	server, store := setupTestServer(t)

	// Add a test task
	task := storage.Task{
		ID:    "1",
		Title: "Test Task",
	}
	err := store.Add(task)
	require.NoError(t, err)

	// Update the task
	task.Title = "Updated Task"
	body, err := json.Marshal(task)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/tasks/1", bytes.NewReader(body))
	req = req.WithContext(req.Context())
	w := httptest.NewRecorder()

	server.handleUpdateTask(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedTask storage.Task
	err = json.NewDecoder(w.Body).Decode(&updatedTask)
	require.NoError(t, err)
	assert.Equal(t, "Updated Task", updatedTask.Title)
}

func TestDeleteTask(t *testing.T) {
	server, store := setupTestServer(t)

	// Add a test task
	task := storage.Task{
		ID:    "1",
		Title: "Test Task",
	}
	err := store.Add(task)
	require.NoError(t, err)

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	req = req.WithContext(req.Context())
	w := httptest.NewRecorder()

	server.handleDeleteTask(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify task is deleted
	tasks, err := store.List()
	require.NoError(t, err)
	assert.Empty(t, tasks)
}
