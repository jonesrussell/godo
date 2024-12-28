package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// Middleware represents a function that wraps an http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies a sequence of middleware to a handler
func Chain(h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// WithValidation validates the request body against the provided type
func WithValidation[T any](log logger.Logger) Middleware {
	validate := validator.New()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var req T
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Error("failed to decode request body", "error", err)
				writeError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
				return
			}

			if err := validate.Struct(req); err != nil {
				var validationErrors validator.ValidationErrors
				if errors.As(err, &validationErrors) {
					fields := make(map[string]string)
					for _, err := range validationErrors {
						fields[err.Field()] = err.Tag()
					}
					writeValidationError(w, fields)
					return
				}
				writeError(w, http.StatusBadRequest, "validation_failed", "Request validation failed")
				return
			}

			// Store the validated request in the context
			ctx := r.Context()
			ctx = context.WithValue(ctx, requestKey{}, req)
			next(w, r.WithContext(ctx))
		}
	}
}

// WithLogging logs request details
func WithLogging(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Info("handling request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)
			next(w, r)
		}
	}
}

// WithErrorHandling handles errors from the handler
func WithErrorHandling(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered", "error", err)
					writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
				}
			}()
			next(w, r)
		}
	}
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeError writes an error response
func writeError(w http.ResponseWriter, status int, code, message string) {
	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}
	writeJSON(w, status, resp)
}

// writeValidationError writes a validation error response
func writeValidationError(w http.ResponseWriter, fields map[string]string) {
	resp := ValidationErrorResponse{
		Code:    "validation_error",
		Message: "Request validation failed",
		Fields:  fields,
	}
	writeJSON(w, http.StatusBadRequest, resp)
}

// mapError maps storage errors to HTTP status codes and error responses
func mapError(err error) (int, string, string) {
	switch {
	case errors.Is(err, storage.ErrTaskNotFound):
		return http.StatusNotFound, "not_found", "Task not found"
	case errors.Is(err, storage.ErrDuplicateID):
		return http.StatusConflict, "duplicate_id", "Task ID already exists"
	case errors.Is(err, storage.ErrInvalidID):
		return http.StatusBadRequest, "invalid_id", "Invalid task ID"
	default:
		return http.StatusInternalServerError, "internal_error", "Internal server error"
	}
}

// requestKey is a type for request context keys
type requestKey struct{}

// GetRequest retrieves the validated request from the context
func GetRequest[T any](r *http.Request) (T, bool) {
	ctx := r.Context()
	req, ok := ctx.Value(requestKey{}).(T)
	return req, ok
}
