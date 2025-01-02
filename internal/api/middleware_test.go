package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=0,lte=150"`
}

func TestWithValidation(t *testing.T) {
	log := logger.NewTestLogger(t)

	tests := []struct {
		name       string
		req        any
		wantStatus int
	}{
		{
			name: "valid request",
			req: testRequest{
				Name: "Test",
				Age:  30,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "missing required field",
			req: testRequest{
				Age: 30,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid age",
			req: testRequest{
				Name: "Test",
				Age:  -1,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			req:        "{invalid json}",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if s, ok := tt.req.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler := WithValidation[testRequest](log)(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			handler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestWithLogging(t *testing.T) {
	log := logger.NewTestLogger(t)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	called := false
	handler := WithLogging(log)(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler(w, req)

	assert.True(t, called, "handler should be called")
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestWithErrorHandling(t *testing.T) {
	log := logger.NewTestLogger(t)

	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
	}{
		{
			name: "no error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "task not found error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(storage.ErrTaskNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("unexpected error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			handler := WithErrorHandling(log)(tt.handler)
			handler(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestChain(t *testing.T) {
	order := make([]string, 0)

	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1_before")
			next(w, r)
			order = append(order, "m1_after")
		}
	}

	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2_before")
			next(w, r)
			order = append(order, "m2_after")
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	Chain(handler, middleware1, middleware2)(w, req)

	expected := []string{
		"m1_before",
		"m2_before",
		"handler",
		"m2_after",
		"m1_after",
	}
	assert.Equal(t, expected, order, "middleware should be executed in correct order")
	assert.Equal(t, http.StatusOK, w.Code)
}
