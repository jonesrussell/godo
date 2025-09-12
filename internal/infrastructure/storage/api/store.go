// Package api provides API-based implementation of the unified storage interface
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jonesrussell/godo/internal/domain/model"
	domainstorage "github.com/jonesrussell/godo/internal/domain/storage"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
	storageerrors "github.com/jonesrussell/godo/internal/infrastructure/storage/errors"
)

// Store implements domain.storage.UnifiedNoteStorage using HTTP API calls
type Store struct {
	client     *http.Client
	baseURL    string
	logger     logger.Logger
	retryCount int
	retryDelay time.Duration
}

// New creates a new API store
func New(config domainstorage.APIConfig, log logger.Logger) (*Store, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("API base URL is required")
	}

	timeout := time.Duration(config.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	retryCount := config.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}

	retryDelay := time.Duration(config.RetryDelay) * time.Millisecond
	if retryDelay <= 0 {
		retryDelay = 1000 * time.Millisecond
	}

	client := &http.Client{
		Timeout: timeout,
	}

	return &Store{
		client:     client,
		baseURL:    config.BaseURL,
		logger:     log,
		retryCount: retryCount,
		retryDelay: retryDelay,
	}, nil
}

// CreateNote creates a new note via API
func (s *Store) CreateNote(ctx context.Context, content string) (*model.Note, error) {
	requestBody := map[string]string{
		"content": content,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/notes", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// GetNote retrieves a note by ID via API
func (s *Store) GetNote(ctx context.Context, id string) (*model.Note, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+"/notes/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// GetAllNotes retrieves all notes via API
func (s *Store) GetAllNotes(ctx context.Context) ([]*model.Note, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+"/notes", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIListResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	notes := make([]*model.Note, len(apiResp.Data))
	for i, apiNote := range apiResp.Data {
		notes[i] = s.mapAPINoteToModel(&apiNote)
	}

	return notes, nil
}

// UpdateNote updates a note via API
func (s *Store) UpdateNote(ctx context.Context, id string, content string, done bool) (*model.Note, error) {
	requestBody := map[string]interface{}{
		"content": content,
		"done":    done,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", s.baseURL+"/notes/"+id, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// DeleteNote deletes a note via API
func (s *Store) DeleteNote(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", s.baseURL+"/notes/"+id, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &storageerrors.NotFoundError{ID: id}
	}
	if resp.StatusCode == http.StatusNotFound {
		return &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return s.handleAPIError(resp)
	}

	return nil
}

// ToggleDone toggles the done status of a note via API
func (s *Store) ToggleDone(ctx context.Context, id string) (*model.Note, error) {
	req, err := http.NewRequestWithContext(ctx, "PATCH", s.baseURL+"/notes/"+id+"/toggle", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// MarkDone marks a note as done via API
func (s *Store) MarkDone(ctx context.Context, id string) (*model.Note, error) {
	req, err := http.NewRequestWithContext(ctx, "PATCH", s.baseURL+"/notes/"+id+"/mark-done", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// MarkUndone marks a note as undone via API
func (s *Store) MarkUndone(ctx context.Context, id string) (*model.Note, error) {
	req, err := http.NewRequestWithContext(ctx, "PATCH", s.baseURL+"/notes/"+id+"/mark-undone", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.executeWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, &storageerrors.NotFoundError{ID: id}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, s.handleAPIError(resp)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.mapAPINoteToModel(&apiResp.Data), nil
}

// Close closes the API store (no-op for HTTP client)
func (s *Store) Close() error {
	s.client.CloseIdleConnections()
	return nil
}

// executeWithRetry executes HTTP request with retry logic
func (s *Store) executeWithRetry(req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= s.retryCount; attempt++ {
		if attempt > 0 {
			s.logger.Debug("Retrying API request", "attempt", attempt, "url", req.URL.String())
			time.Sleep(s.retryDelay)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}

		// Don't retry on client errors (4xx) except for 429 (rate limit)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			return resp, nil
		}

		// Don't retry on success or server errors that shouldn't be retried
		if resp.StatusCode < 500 || resp.StatusCode == 501 || resp.StatusCode == 505 {
			return resp, nil
		}

		// Server error (5xx) - retry
		resp.Body.Close()
		lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// handleAPIError handles API error responses
func (s *Store) handleAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var apiErr APIErrorResponse
	if err := json.Unmarshal(body, &apiErr); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if len(apiErr.Errors) > 0 {
		return &storageerrors.ValidationError{
			Message: apiErr.Message,
			Fields:  apiErr.Errors,
		}
	}

	return fmt.Errorf("API error: %s", apiErr.Message)
}

// mapAPINoteToModel converts API note format to domain model
func (s *Store) mapAPINoteToModel(apiNote *APINote) *model.Note {
	return &model.Note{
		ID:        apiNote.ID,
		Content:   apiNote.Content,
		Done:      apiNote.Done,
		CreatedAt: apiNote.CreatedAt,
		UpdatedAt: apiNote.UpdatedAt,
	}
}

// API Response structures

// APIResponse represents a single note API response
type APIResponse struct {
	Data    APINote `json:"data"`
	Message string  `json:"message"`
}

// APIListResponse represents a list of notes API response
type APIListResponse struct {
	Data    []APINote `json:"data"`
	Message string    `json:"message"`
}

// APINote represents a note in API format
type APINote struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIErrorResponse represents an API error response
type APIErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}
