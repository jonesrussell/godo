package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"

	"github.com/jonesrussell/godo/internal/logger"
	"github.com/jonesrussell/godo/internal/model"
)

const internalServerErrorMsg = "Internal server error"

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

// WithErrorHandling adds error handling to a handler
func WithErrorHandling(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered", "error", err)

					var status int
					var code string
					var msg string

					switch e := err.(type) {
					case error:
						if errors.Is(e, model.ErrTaskNotFound) {
							status = http.StatusNotFound
							code = "task_not_found"
							msg = "Task not found"
						} else {
							status = http.StatusInternalServerError
							code = "internal_error"
							msg = internalServerErrorMsg
						}
					default:
						status = http.StatusInternalServerError
						code = "internal_error"
						msg = internalServerErrorMsg
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
	case errors.Is(err, model.ErrTaskNotFound):
		return http.StatusNotFound, "Task not found", err.Error()
	case errors.Is(err, model.ErrDuplicateID):
		return http.StatusConflict, "Task ID already exists", err.Error()
	default:
		return http.StatusInternalServerError, internalServerErrorMsg, err.Error()
	}
}

// requestKey is a type for request context keys
type requestKey struct{}

// userIDKey is a type for user ID context keys
type userIDKey struct{}

// GetRequest retrieves the validated request from the context
func GetRequest[T any](r *http.Request) (T, bool) {
	ctx := r.Context()
	req, ok := ctx.Value(requestKey{}).(T)
	return req, ok
}

// GetUserID retrieves the user ID from the JWT token context
func GetUserID(r *http.Request) (string, bool) {
	ctx := r.Context()
	userID, ok := ctx.Value(userIDKey{}).(string)
	return userID, ok
}

// WithJWTAuth validates JWT tokens and extracts user information
func WithJWTAuth(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Get JWT secret from environment
			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				log.Error("JWT_SECRET environment variable not set")
				writeError(w, http.StatusInternalServerError, "server_error", "Authentication configuration error")
				return
			}

			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing_token", "Authorization header required")
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeError(w, http.StatusUnauthorized, "invalid_token_format", "Authorization header must be Bearer token")
				return
			}

			tokenString := parts[1]

			// Parse and validate JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				log.Error("JWT token validation failed", "error", err)
				writeError(w, http.StatusUnauthorized, "invalid_token", "Invalid or expired token")
				return
			}

			// Extract claims
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Extract user ID from claims
				userIDInterface, exists := claims["user_id"]
				if !exists {
					writeError(w, http.StatusUnauthorized, "invalid_token_claims", "Token missing user_id claim")
					return
				}

				userID, userOk := userIDInterface.(string)
				if !userOk {
					writeError(w, http.StatusUnauthorized, "invalid_token_claims", "Invalid user_id claim type")
					return
				}

				// Add user ID to context
				ctx := context.WithValue(r.Context(), userIDKey{}, userID)
				next(w, r.WithContext(ctx))
			} else {
				writeError(w, http.StatusUnauthorized, "invalid_token", "Invalid token claims")
				return
			}
		}
	}
}

// WithOptionalJWTAuth validates JWT tokens if present, but allows requests without tokens
func WithOptionalJWTAuth(log logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No auth header, continue without user context
				next(w, r)
				return
			}

			// If auth header is present, validate it using the same logic as WithJWTAuth
			jwtAuthMiddleware := WithJWTAuth(log)
			jwtAuthMiddleware(next)(w, r)
		}
	}
}
