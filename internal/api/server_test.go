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

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskHandler(t *testing.T) {
	store := mock.New()
	handler := NewTaskHandler(store)

	t.Run("create task", func(t *testing.T) {
		tests := []struct {
			name       string
			task       storage.Task
			wantStatus int
			wantErr    bool
		}{
			{
				name: "valid task",
				task: storage.Task{
					Title:       "Test Task",
					Description: "Test Description",
					CreatedAt:   time.Now().Unix(),
					UpdatedAt:   time.Now().Unix(),
				},
				wantStatus: http.StatusCreated,
			},
			{
				name: "empty title",
				task: storage.Task{
					Title:       "",
					Description: "Test Description",
					CreatedAt:   time.Now().Unix(),
					UpdatedAt:   time.Now().Unix(),
				},
				wantStatus: http.StatusBadRequest,
				wantErr:    true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				body, err := json.Marshal(tt.task)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
				w := httptest.NewRecorder()

				handler.CreateTask(w, req)

				assert.Equal(t, tt.wantStatus, w.Code)

				if !tt.wantErr {
					var response storage.Task
					err := json.NewDecoder(w.Body).Decode(&response)
					require.NoError(t, err)
					assert.Equal(t, tt.task.Title, response.Title)
					assert.Equal(t, tt.task.Description, response.Description)
				}
			})
		}
	})

	t.Run("get task", func(t *testing.T) {
		tests := []struct {
			name       string
			taskID     string
			wantTask   *storage.Task
			wantStatus int
		}{
			{
				name:   "existing task",
				taskID: "test-1",
				wantTask: &storage.Task{
					ID:          "test-1",
					Title:       "Test Task",
					Description: "Test Description",
					Completed:   false,
					CreatedAt:   time.Now().Unix(),
					UpdatedAt:   time.Now().Unix(),
				},
				wantStatus: http.StatusOK,
			},
			{
				name:       "nonexistent task",
				taskID:     "nonexistent",
				wantStatus: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.wantTask != nil {
					err := store.Add(context.Background(), *tt.wantTask)
					require.NoError(t, err)
				}

				req := httptest.NewRequest(http.MethodGet, "/tasks/"+tt.taskID, nil)
				w := httptest.NewRecorder()

				handler.GetTask(w, req)

				assert.Equal(t, tt.wantStatus, w.Code)

				if tt.wantTask != nil {
					var response storage.Task
					err := json.NewDecoder(w.Body).Decode(&response)
					require.NoError(t, err)
					assert.Equal(t, tt.wantTask.Title, response.Title)
					assert.Equal(t, tt.wantTask.Description, response.Description)
					assert.Equal(t, tt.wantTask.Completed, response.Completed)
				}
			})
		}
	})

	t.Run("update task", func(t *testing.T) {
		tests := []struct {
			name       string
			taskID     string
			update     storage.Task
			wantStatus int
		}{
			{
				name:   "valid update",
				taskID: "test-1",
				update: storage.Task{
					ID:          "test-1",
					Title:       "Updated Task",
					Description: "Updated Description",
					Completed:   true,
					CreatedAt:   time.Now().Unix(),
					UpdatedAt:   time.Now().Unix(),
				},
				wantStatus: http.StatusOK,
			},
			{
				name:   "nonexistent task",
				taskID: "nonexistent",
				update: storage.Task{
					ID:          "nonexistent",
					Title:       "Updated Task",
					Description: "Updated Description",
					Completed:   true,
					CreatedAt:   time.Now().Unix(),
					UpdatedAt:   time.Now().Unix(),
				},
				wantStatus: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.taskID == "test-1" {
					existingTask := storage.Task{
						ID:          "test-1",
						Title:       "Original Task",
						Description: "Original Description",
						Completed:   false,
						CreatedAt:   time.Now().Unix(),
						UpdatedAt:   time.Now().Unix(),
					}
					err := store.Add(context.Background(), existingTask)
					require.NoError(t, err)
				}

				body, err := json.Marshal(tt.update)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewReader(body))
				w := httptest.NewRecorder()

				handler.UpdateTask(w, req)

				assert.Equal(t, tt.wantStatus, w.Code)

				if tt.wantStatus == http.StatusOK {
					var response storage.Task
					err := json.NewDecoder(w.Body).Decode(&response)
					require.NoError(t, err)
					assert.Equal(t, tt.update.Title, response.Title)
					assert.Equal(t, tt.update.Description, response.Description)
					assert.Equal(t, tt.update.Completed, response.Completed)
				}
			})
		}
	})

	t.Run("delete task", func(t *testing.T) {
		tests := []struct {
			name       string
			taskID     string
			wantStatus int
		}{
			{
				name:       "existing task",
				taskID:     "test-1",
				wantStatus: http.StatusNoContent,
			},
			{
				name:       "nonexistent task",
				taskID:     "nonexistent",
				wantStatus: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.taskID == "test-1" {
					task := storage.Task{
						ID:          "test-1",
						Title:       "Test Task",
						Description: "Test Description",
						Completed:   false,
						CreatedAt:   time.Now().Unix(),
						UpdatedAt:   time.Now().Unix(),
					}
					err := store.Add(context.Background(), task)
					require.NoError(t, err)
				}

				req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)
				w := httptest.NewRecorder()

				handler.DeleteTask(w, req)

				assert.Equal(t, tt.wantStatus, w.Code)
			})
		}
	})

	t.Run("list tasks", func(t *testing.T) {
		// Clear any existing tasks
		store.Close()

		// Add some test tasks
		tasks := []storage.Task{
			{
				ID:          "test-1",
				Title:       "Task 1",
				Description: "Description 1",
				Completed:   false,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			},
			{
				ID:          "test-2",
				Title:       "Task 2",
				Description: "Description 2",
				Completed:   true,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   time.Now().Unix(),
			},
		}

		for _, task := range tasks {
			err := store.Add(context.Background(), task)
			require.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		handler.ListTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []storage.Task
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Len(t, response, len(tasks))

		for i, task := range tasks {
			assert.Equal(t, task.Title, response[i].Title)
			assert.Equal(t, task.Description, response[i].Description)
			assert.Equal(t, task.Completed, response[i].Completed)
		}
	})
}
