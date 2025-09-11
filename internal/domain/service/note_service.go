// Package service provides business logic services for the application
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jonesrussell/godo/internal/domain/model"
	"github.com/jonesrussell/godo/internal/domain/repository"
	"github.com/jonesrussell/godo/internal/infrastructure/logger"
)

//go:generate mockgen -destination=../../test/mocks/mock_noteservice.go -package=mocks github.com/jonesrussell/godo/internal/domain/service NoteService

// NoteFilter represents filtering options for note queries
type NoteFilter struct {
	Done          *bool      `json:"done,omitempty"`
	Content       *string    `json:"content,omitempty"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	Limit         *int       `json:"limit,omitempty"`
	Offset        *int       `json:"offset,omitempty"`
}

// NoteUpdateRequest represents a request to update a note
type NoteUpdateRequest struct {
	Content *string `json:"content,omitempty"`
	Done    *bool   `json:"done,omitempty"`
}

// NoteService defines the interface for note business logic operations
type NoteService interface {
	CreateNote(ctx context.Context, content string) (*model.Note, error)
	GetNote(ctx context.Context, id string) (*model.Note, error)
	UpdateNote(ctx context.Context, id string, updates NoteUpdateRequest) (*model.Note, error)
	DeleteNote(ctx context.Context, id string) error
	ListNotes(ctx context.Context, filter *NoteFilter) ([]*model.Note, error)
}

// noteService implements NoteService
type noteService struct {
	repo   repository.NoteRepository
	logger logger.Logger
}

// NewNoteService creates a new NoteService instance
func NewNoteService(repo repository.NoteRepository, log logger.Logger) NoteService {
	return &noteService{
		repo:   repo,
		logger: log,
	}
}

func (s *noteService) validateNoteContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return &model.ValidationError{
			Field:   "content",
			Message: "note content cannot be empty",
		}
	}
	if len(content) > 1000 {
		return &model.ValidationError{
			Field:   "content",
			Message: "note content cannot exceed 1000 characters",
		}
	}
	return nil
}

func (s *noteService) validateNoteID(id string) error {
	if strings.TrimSpace(id) == "" {
		return &model.ValidationError{
			Field:   "id",
			Message: "note ID cannot be empty",
		}
	}
	if _, err := uuid.Parse(id); err != nil {
		return &model.ValidationError{
			Field:   "id",
			Message: "invalid note ID format",
		}
	}
	return nil
}

func (s *noteService) CreateNote(ctx context.Context, content string) (*model.Note, error) {
	s.logger.Info("Creating new note", "content_length", len(content))
	if err := s.validateNoteContent(content); err != nil {
		s.logger.Error("Note content validation failed", "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	note := model.Note{
		ID:        uuid.New().String(),
		Content:   strings.TrimSpace(content),
		Done:      false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Add(ctx, &note); err != nil {
		s.logger.Error("Failed to store note", "note_id", note.ID, "error", err)
		return nil, fmt.Errorf("failed to create note: %w", err)
	}
	s.logger.Info("Note created successfully", "note_id", note.ID)
	return &note, nil
}

func (s *noteService) GetNote(ctx context.Context, id string) (*model.Note, error) {
	s.logger.Info("Retrieving note", "note_id", id)
	if err := s.validateNoteID(id); err != nil {
		s.logger.Error("Note ID validation failed", "note_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	note, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve note", "note_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve note: %w", err)
	}
	s.logger.Info("Note retrieved successfully", "note_id", id)
	return note, nil
}

func (s *noteService) UpdateNote(ctx context.Context, id string, updates NoteUpdateRequest) (*model.Note, error) {
	s.logger.Info("Updating note", "note_id", id, "updates", updates)
	if err := s.validateNoteID(id); err != nil {
		s.logger.Error("Note ID validation failed", "note_id", id, "error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	existingNote, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to retrieve existing note", "note_id", id, "error", err)
		return nil, fmt.Errorf("failed to retrieve note: %w", err)
	}
	if updates.Content != nil {
		if validErr := s.validateNoteContent(*updates.Content); validErr != nil {
			s.logger.Error("Note content validation failed", "note_id", id, "error", validErr)
			return nil, fmt.Errorf("validation failed: %w", validErr)
		}
		existingNote.Content = strings.TrimSpace(*updates.Content)
	}
	if updates.Done != nil {
		existingNote.Done = *updates.Done
	}
	existingNote.UpdatedAt = time.Now()
	if updateErr := s.repo.Update(ctx, existingNote); updateErr != nil {
		s.logger.Error("Failed to update note", "note_id", id, "error", updateErr)
		return nil, fmt.Errorf("failed to update note: %w", updateErr)
	}
	s.logger.Info("Note updated successfully", "note_id", id)
	return existingNote, nil
}

func (s *noteService) DeleteNote(ctx context.Context, id string) error {
	s.logger.Info("Deleting note", "note_id", id)
	if err := s.validateNoteID(id); err != nil {
		s.logger.Error("Note ID validation failed", "note_id", id, "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete note", "note_id", id, "error", err)
		return fmt.Errorf("failed to delete note: %w", err)
	}
	s.logger.Info("Note deleted successfully", "note_id", id)
	return nil
}

func (s *noteService) ListNotes(ctx context.Context, filter *NoteFilter) ([]*model.Note, error) {
	s.logger.Info("Retrieving notes", "filter", filter)
	notes, err := s.repo.List(ctx)
	if err != nil {
		s.logger.Error("Failed to retrieve notes", "error", err)
		return nil, fmt.Errorf("failed to retrieve notes: %w", err)
	}
	if filter != nil {
		notes = s.applyFilters(notes, filter)
	}
	s.logger.Info("Notes retrieved successfully", "count", len(notes))
	return notes, nil
}

// applyFilters applies the given filters to the note list
func (s *noteService) applyFilters(notes []*model.Note, filter *NoteFilter) []*model.Note {
	if filter == nil {
		return notes
	}
	filtered := s.filterByCriteria(notes, filter)
	return s.applyPagination(filtered, filter)
}

// filterByCriteria applies content and date filters
func (s *noteService) filterByCriteria(notes []*model.Note, filter *NoteFilter) []*model.Note {
	var filtered []*model.Note
	for _, note := range notes {
		if s.matchesFilter(note, filter) {
			filtered = append(filtered, note)
		}
	}
	return filtered
}

// applyPagination applies limit and offset to the filtered results
func (s *noteService) applyPagination(notes []*model.Note, filter *NoteFilter) []*model.Note {
	if filter.Limit == nil || *filter.Limit <= 0 {
		return notes
	}
	limit := *filter.Limit
	offset := 0
	if filter.Offset != nil && *filter.Offset > 0 {
		offset = *filter.Offset
	}
	return s.sliceWithBounds(notes, offset, limit)
}

// sliceWithBounds safely slices the notes array with offset and limit
func (s *noteService) sliceWithBounds(notes []*model.Note, offset, limit int) []*model.Note {
	if offset >= len(notes) {
		return []*model.Note{}
	}
	end := offset + limit
	if end > len(notes) {
		end = len(notes)
	}
	return notes[offset:end]
}

// matchesFilter checks if a note matches the given filter criteria
func (s *noteService) matchesFilter(note *model.Note, filter *NoteFilter) bool {
	if filter.Done != nil && note.Done != *filter.Done {
		return false
	}
	if filter.Content != nil && !strings.Contains(strings.ToLower(note.Content), strings.ToLower(*filter.Content)) {
		return false
	}
	if filter.CreatedAfter != nil && note.CreatedAt.Before(*filter.CreatedAfter) {
		return false
	}
	if filter.CreatedBefore != nil && note.CreatedAt.After(*filter.CreatedBefore) {
		return false
	}
	return true
}
