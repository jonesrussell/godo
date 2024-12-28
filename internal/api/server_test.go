package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name           string
	method         string
	path           string
	body           interface{}
	setupFn        func(*testutil.MockStore)
	expectedStatus int
	validateFn     func(*testing.T, *httptest.ResponseRecorder)
}

func setupTestServer(t *testing.T) (*Server, *testutil.MockStore) {
	t.Helper()
	log, err := logger.New(&common.LogConfig{Level: "debug", Output: []string{"stdout"}, ErrorOutput: []string{"stderr"}})
	require.NoError(t, err)
	store := testutil.NewMockStore()

	config := &common.HTTPConfig{
		Port:              8080,
		ReadTimeout:       30,
		WriteTimeout:      30,
		ReadHeaderTimeout: 10,
		IdleTimeout:       120,
	}

	server := NewServer(store, log, config)
	require.NotNil(t, server)
	return server, store
}

func executeRequest(t *testing.T, server *Server, tc testCase) *httptest.ResponseRecorder {
	t.Helper()
	var body []byte
	var err error
	if tc.body != nil {
		body, err = json.Marshal(tc.body)
		require.NoError(t, err)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
	server.router.ServeHTTP(w, req)
	return w
}

func TestHealthCheck(t *testing.T) {
	server, _ := setupTestServer(t)
	tc := testCase{
		name:           "Health check returns OK",
		method:         "GET",
		path:           "/health",
		expectedStatus: http.StatusOK,
		validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
			var response map[string]string
			require.NoError(t, json.NewDecoder(w.Body).Decode(&response))
			assert.Equal(t, "ok", response["status"])
		},
	}
	w := executeRequest(t, server, tc)
	assert.Equal(t, tc.expectedStatus, w.Code)
	tc.validateFn(t, w)
}

func TestTaskOperations(t *testing.T) {
	now := time.Now()
	testTask := storage.Task{
		ID:        "1",
		Content:   "Test Task",
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []testCase{
		{
			name:   "List tasks",
			method: "GET",
			path:   "/api/v1/tasks",
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var tasks []storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&tasks))
				assert.Len(t, tasks, 1)
				assert.Equal(t, testTask.ID, tasks[0].ID)
			},
		},
		{
			name:           "Create task",
			method:         "POST",
			path:           "/api/v1/tasks",
			body:           testTask,
			expectedStatus: http.StatusCreated,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, testTask.Content, task.Content)
			},
		},
		{
			name:   "Update task",
			method: "PUT",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: storage.Task{
				ID:        testTask.ID,
				Content:   "Updated Task",
				Done:      true,
				UpdatedAt: now,
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, "Updated Task", task.Content)
			},
		},
		{
			name:   "Delete task",
			method: "DELETE",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			expectedStatus: http.StatusNoContent,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Body.String())
			},
		},
		{
			name:   "Patch task - update title only",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				Title: testutil.StringPtr("Updated Title Only"),
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, "Updated Title Only", task.Title)
				assert.Equal(t, testTask.Description, task.Description)
			},
		},
		{
			name:   "Patch task - update description only",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				Description: testutil.StringPtr("Updated Description Only"),
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, testTask.Title, task.Title)
				assert.Equal(t, "Updated Description Only", task.Description)
			},
		},
		{
			name:   "Patch task - mark as completed",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				CompletedAt: &now,
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, testTask.Title, task.Title)
				assert.Equal(t, testTask.Description, task.Description)
				assert.True(t, task.CompletedAt.Equal(now), "CompletedAt time should match")
			},
		},
		{
			name:   "Patch task - not found",
			method: "PATCH",
			path:   "/api/v1/tasks/nonexistent",
			body: TaskPatch{
				Title: testutil.StringPtr("Updated Title"),
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Patch task - invalid body",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body:           "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server, store := setupTestServer(t)
			if tc.setupFn != nil {
				tc.setupFn(store)
			}
			w := executeRequest(t, server, tc)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.validateFn != nil {
				tc.validateFn(t, w)
			}
		})
	}
}
