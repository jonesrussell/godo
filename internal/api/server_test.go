package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
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
			name:   "Patch task - update content only",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				Content: testutil.StringPtr("Updated Content"),
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, "Updated Content", task.Content)
				assert.Equal(t, testTask.Done, task.Done)
			},
		},
		{
			name:   "Patch task - update done status only",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				Done: testutil.BoolPtr(true),
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, testTask.Content, task.Content)
				assert.True(t, task.Done)
			},
		},
		{
			name:   "Patch task - update both fields",
			method: "PATCH",
			path:   fmt.Sprintf("/api/v1/tasks/%s", testTask.ID),
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			body: TaskPatch{
				Content: testutil.StringPtr("Updated Content"),
				Done:    testutil.BoolPtr(true),
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				var task storage.Task
				require.NoError(t, json.NewDecoder(w.Body).Decode(&task))
				assert.Equal(t, "Updated Content", task.Content)
				assert.True(t, task.Done)
			},
		},
		{
			name:   "Patch task - not found",
			method: "PATCH",
			path:   "/api/v1/tasks/nonexistent",
			body: TaskPatch{
				Content: testutil.StringPtr("Updated Content"),
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

func TestHandlePatchTask(t *testing.T) {
	store := testutil.NewMockStore()
	logger := logger.NewTestLogger(t)
	config := &common.HTTPConfig{Port: 8080}
	server := NewServer(store, logger, config)

	// Add a test task
	testTask := storage.Task{
		ID:        "test-id",
		Content:   "Original content",
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	store.Add(testTask)

	t.Run("successfully patches a task", func(t *testing.T) {
		newContent := "Updated content"
		newDone := true
		patch := TaskPatch{
			Content: &newContent,
			Done:    &newDone,
		}
		body, _ := json.Marshal(patch)
		req := httptest.NewRequest(http.MethodPatch, "/tasks/test-id", bytes.NewReader(body))
		w := httptest.NewRecorder()

		// Set up chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-id")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		server.handlePatchTask(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var task storage.Task
		json.NewDecoder(w.Body).Decode(&task)
		assert.Equal(t, testTask.ID, task.ID)
		assert.Equal(t, newContent, task.Content)
		assert.Equal(t, newDone, task.Done)
	})

	t.Run("returns 404 for non-existent task", func(t *testing.T) {
		patch := TaskPatch{
			Content: testutil.StringPtr("Updated content"),
		}
		body, _ := json.Marshal(patch)
		req := httptest.NewRequest(http.MethodPatch, "/tasks/non-existent", bytes.NewReader(body))
		w := httptest.NewRecorder()

		// Set up chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "non-existent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		server.handlePatchTask(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/tasks/test-id", bytes.NewReader([]byte("invalid json")))
		w := httptest.NewRecorder()

		// Set up chi router context
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "test-id")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		server.handlePatchTask(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestServerErrorHandling(t *testing.T) {
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
			name:   "List tasks - store closed",
			method: "GET",
			path:   "/api/v1/tasks",
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Close())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Create task - invalid JSON",
			method:         "POST",
			path:           "/api/v1/tasks",
			body:           "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Create task - duplicate ID",
			method: "POST",
			path:   "/api/v1/tasks",
			body:   testTask,
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Add(testTask))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Update task - not found",
			method:         "PUT",
			path:           "/api/v1/tasks/nonexistent",
			body:           testTask,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Update task - invalid JSON",
			method:         "PUT",
			path:           "/api/v1/tasks/1",
			body:           "{invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Update task - store closed",
			method: "PUT",
			path:   "/api/v1/tasks/1",
			body:   testTask,
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Close())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Delete task - not found",
			method:         "DELETE",
			path:           "/api/v1/tasks/nonexistent",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Delete task - store closed",
			method: "DELETE",
			path:   "/api/v1/tasks/1",
			setupFn: func(store *testutil.MockStore) {
				require.NoError(t, store.Close())
			},
			expectedStatus: http.StatusInternalServerError,
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
		})
	}
}

func TestServerStartAndShutdown(t *testing.T) {
	server, _ := setupTestServer(t)
	require.NotNil(t, server)

	// Test server start
	go func() {
		err := server.Start(0) // Use port 0 to let the OS choose a free port
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server start error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServerMiddleware(t *testing.T) {
	server, _ := setupTestServer(t)
	require.NotNil(t, server)

	tests := []testCase{
		{
			name:           "Request ID middleware",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.NotEmpty(t, w.Header().Get("X-Request-Id"))
			},
		},
		{
			name:           "Content type middleware",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := executeRequest(t, server, tc)
			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.validateFn != nil {
				tc.validateFn(t, w)
			}
		})
	}
}

func TestServerTimeouts(t *testing.T) {
	log, err := logger.New(&common.LogConfig{Level: "debug"})
	require.NoError(t, err)
	store := testutil.NewMockStore()

	// Test with very short timeouts
	config := &common.HTTPConfig{
		Port:              8080,
		ReadTimeout:       1,
		WriteTimeout:      1,
		ReadHeaderTimeout: 1,
		IdleTimeout:       1,
	}

	server := NewServer(store, log, config)
	require.NotNil(t, server)

	// Start server
	go func() {
		startErr := server.Start(0)
		if startErr != nil && startErr != http.ErrServerClosed {
			t.Errorf("Server start error: %v", startErr)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test timeouts
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdownErr := server.Shutdown(ctx)
	assert.NoError(t, shutdownErr)
}

func TestServerInvalidRequests(t *testing.T) {
	server, _ := setupTestServer(t)
	require.NotNil(t, server)

	tests := []testCase{
		{
			name:           "Invalid method",
			method:         "INVALID",
			path:           "/health",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Invalid path",
			method:         "GET",
			path:           "/invalid/path",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid content type",
			method:         "POST",
			path:           "/api/v1/tasks",
			body:           []byte("plain text"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := executeRequest(t, server, tc)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestHandleDeleteTask(t *testing.T) {
	store := storage.NewMemoryStore()

	// Add a test task
	task := storage.Task{
		ID:      "test-id",
		Content: "Test Task",
	}
	err := store.Add(task)
	require.NoError(t, err)

	// Delete the task
	deleteErr := store.Delete("test-id")
	assert.NoError(t, deleteErr)
}
