// Package memory provides an in-memory implementation of the storage interface
package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements the storage.Store interface using in-memory storage
type Store struct {
	mu    sync.RWMutex
	notes map[string]storage.Note
}

// New creates a new memory store instance
func New() *Store {
	return &Store{
		notes: make(map[string]storage.Note),
	}
}

// Add adds a new note
func (s *Store) Add(_ context.Context, note storage.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; exists {
		return fmt.Errorf("note with ID %s already exists", note.ID)
	}

	s.notes[note.ID] = note
	return nil
}

// Get retrieves a note by ID
func (s *Store) Get(_ context.Context, id string) (storage.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	note, exists := s.notes[id]
	if !exists {
		return storage.Note{}, fmt.Errorf("note with ID %s not found", id)
	}

	return note, nil
}

// List retrieves all notes
func (s *Store) List(_ context.Context) ([]storage.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notes := make([]storage.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Update updates an existing note
func (s *Store) Update(_ context.Context, note storage.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[note.ID]; !exists {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.notes[id]; !exists {
		return fmt.Errorf("note with ID %s not found", id)
	}

	delete(s.notes, id)
	return nil
}

// Close cleans up any resources
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notes = make(map[string]storage.Note)
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.Transaction, error) {
	return nil, fmt.Errorf("transactions not supported in memory store")
}
