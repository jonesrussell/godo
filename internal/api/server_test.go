//go:build integration
// +build integration

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*Server, *storage.MockStore) {
	store := storage.NewMockStore()
	log := logger.NewTestLogger(t)
	server := NewServer(store, log)
	return server, store
}

func TestHandleCreateTask(t *testing.T) {
	server, store := setupTestServer(t)

	tests := []struct {
		name       string
		task       storage.Task
		setupStore func()
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid task",
			task: storage.Task{
				ID:      "test-1",
				Content: "Test Task",
			},
			setupStore: func() {},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "empty ID",
			task: storage.Task{
				Content: "Test Task",
			},
			setupStore: func() {},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "duplicate ID",
			task: storage.Task{
				ID:      "test-1",
				Content: "Test Task",
			},
			setupStore: func() {
				store.Error = storage.ErrDuplicateID
			},
			wantStatus: http.StatusConflict,
			wantErr:    true,
		},
		{
			name: "store error",
			task: storage.Task{
				ID:      "test-1",
				Content: "Test Task",
			},
			setupStore: func() {
				store.Error = assert.AnError
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Reset()
			tt.setupStore()

			body, err := json.Marshal(tt.task)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleCreateTask(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response storage.Task
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.task.ID, response.ID)
				assert.Equal(t, tt.task.Content, response.Content)
			}
		})
	}
}

func TestHandleGetTask(t *testing.T) {
	server, store := setupTestServer(t)
	ctx := context.Background()

	existingTask := storage.Task{
		ID:        "test-1",
		Content:   "Test Task",
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		taskID     string
		setupStore func()
		wantStatus int
		wantTask   *storage.Task
	}{
		{
			name:   "existing task",
			taskID: existingTask.ID,
			setupStore: func() {
				err := store.Add(ctx, existingTask)
				require.NoError(t, err)
			},
			wantStatus: http.StatusOK,
			wantTask:   &existingTask,
		},
		{
			name:   "nonexistent task",
			taskID: "nonexistent",
			setupStore: func() {
				store.Error = storage.ErrTaskNotFound
			},
			wantStatus: http.StatusNotFound,
			wantTask:   nil,
		},
		{
			name:   "store error",
			taskID: "test-1",
			setupStore: func() {
				store.Error = assert.AnError
			},
			wantStatus: http.StatusInternalServerError,
			wantTask:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Reset()
			tt.setupStore()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+tt.taskID, http.NoBody)
			w := httptest.NewRecorder()

			server.handleGetTask(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantTask != nil {
				var response storage.Task
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.wantTask.ID, response.ID)
				assert.Equal(t, tt.wantTask.Content, response.Content)
				assert.Equal(t, tt.wantTask.Done, response.Done)
			}
		})
	}
}

func TestHandleUpdateTask(t *testing.T) {
	server, store := setupTestServer(t)
	ctx := context.Background()

	existingTask := storage.Task{
		ID:        "test-1",
		Content:   "Test Task",
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		taskID     string
		update     storage.Task
		setupStore func()
		wantStatus int
		wantErr    bool
	}{
		{
			name:   "valid update",
			taskID: existingTask.ID,
			update: storage.Task{
				ID:      existingTask.ID,
				Content: "Updated Task",
				Done:    true,
			},
			setupStore: func() {
				err := store.Add(ctx, existingTask)
				require.NoError(t, err)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:   "nonexistent task",
			taskID: "nonexistent",
			update: storage.Task{
				ID:      "nonexistent",
				Content: "Updated Task",
			},
			setupStore: func() {
				store.Error = storage.ErrTaskNotFound
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:   "store error",
			taskID: existingTask.ID,
			update: storage.Task{
				ID:      existingTask.ID,
				Content: "Updated Task",
			},
			setupStore: func() {
				store.Error = assert.AnError
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Reset()
			tt.setupStore()

			body, err := json.Marshal(tt.update)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+tt.taskID, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.handleUpdateTask(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response storage.Task
				err = json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.update.ID, response.ID)
				assert.Equal(t, tt.update.Content, response.Content)
				assert.Equal(t, tt.update.Done, response.Done)
			}
		})
	}
}

func TestHandleDeleteTask(t *testing.T) {
	server, store := setupTestServer(t)
	ctx := context.Background()

	existingTask := storage.Task{
		ID:        "test-1",
		Content:   "Test Task",
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name       string
		taskID     string
		setupStore func()
		wantStatus int
	}{
		{
			name:   "existing task",
			taskID: existingTask.ID,
			setupStore: func() {
				err := store.Add(ctx, existingTask)
				require.NoError(t, err)
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:   "nonexistent task",
			taskID: "nonexistent",
			setupStore: func() {
				store.Error = storage.ErrTaskNotFound
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "store error",
			taskID: existingTask.ID,
			setupStore: func() {
				store.Error = assert.AnError
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Reset()
			tt.setupStore()

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+tt.taskID, http.NoBody)
			w := httptest.NewRecorder()

			server.handleDeleteTask(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestHandleListTasks(t *testing.T) {
	server, store := setupTestServer(t)
	ctx := context.Background()

	existingTasks := []storage.Task{
		{
			ID:        "test-1",
			Content:   "Test Task 1",
			Done:      false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "test-2",
			Content:   "Test Task 2",
			Done:      true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	tests := []struct {
		name       string
		setupStore func()
		wantStatus int
		wantTasks  []storage.Task
	}{
		{
			name: "list tasks",
			setupStore: func() {
				for _, task := range existingTasks {
					err := store.Add(ctx, task)
					require.NoError(t, err)
				}
			},
			wantStatus: http.StatusOK,
			wantTasks:  existingTasks,
		},
		{
			name: "store error",
			setupStore: func() {
				store.Error = assert.AnError
			},
			wantStatus: http.StatusInternalServerError,
			wantTasks:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Reset()
			tt.setupStore()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", http.NoBody)
			w := httptest.NewRecorder()

			server.handleListTasks(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantTasks != nil {
				var response []storage.Task
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
				assert.Len(t, response, len(tt.wantTasks))
				for i, task := range tt.wantTasks {
					assert.Equal(t, task.ID, response[i].ID)
					assert.Equal(t, task.Content, response[i].Content)
					assert.Equal(t, task.Done, response[i].Done)
				}
			}
		})
	}
}
