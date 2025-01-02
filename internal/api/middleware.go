package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/storage"
)

// ErrorResponse represents a standard error response format for the API
type ErrorResponse struct {
	Code    string `json:"code"`    // Machine-readable error code
	Message string `json:"message"` // Human-readable error message
}

// ValidationErrorResponse represents a validation error response format for the API
type ValidationErrorResponse struct {
	Code    string            `json:"code"`    // Machine-readable error code
	Message string            `json:"message"` // Human-readable error message
	Fields  map[string]string `json:"fields"`  // Map of field names to validation error messages
}

// Middleware represents a function that wraps an http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies a sequence of middleware to a handler
func Chain(h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	// Apply middleware in reverse order so they execute in the order they were passed
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
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

// WithErrorHandling adds panic recovery and error handling to a handler
func WithErrorHandling(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered",
						"error", err,
						"method", r.Method,
						"path", r.URL.Path,
						"remote_addr", r.RemoteAddr,
						"user_agent", r.UserAgent(),
					)

					var status int
					var code string
					var msg string

					switch e := err.(type) {
					case error:
						if errors.Is(e, storage.ErrTaskNotFound) {
							status = http.StatusNotFound
							code = "task_not_found"
							msg = "Task not found"
							log.Info("task not found error", "error", e)
						} else {
							status = http.StatusInternalServerError
							code = "internal_error"
							msg = "Internal server error"
							log.Error("internal server error",
								"error", e,
							)
						}
					default:
						status = http.StatusInternalServerError
						code = "internal_error"
						msg = "Internal server error"
						log.Error("unknown panic value",
							"type", fmt.Sprintf("%T", err),
							"value", fmt.Sprintf("%+v", err),
						)
					}

					writeError(w, status, code, msg)
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

// mapError maps an error to an HTTP status code and error message
func mapError(err error) (code int, msg, details string) {
	switch {
	case errors.Is(err, storage.ErrTaskNotFound):
		return http.StatusNotFound, "Task not found", err.Error()
	case errors.Is(err, storage.ErrDuplicateID):
		return http.StatusConflict, "Task ID already exists", err.Error()
	default:
		return http.StatusInternalServerError, "Internal server error", err.Error()
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
