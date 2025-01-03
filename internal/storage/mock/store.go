// Package mock provides a mock implementation of the storage interface for testing
package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jonesrussell/godo/internal/storage"
)

// Store implements the storage.Store interface for testing
type Store struct {
	mu    sync.RWMutex
	notes map[string]storage.Note
	err   error
}

// Transaction represents a mock transaction
type Transaction struct {
	store     *Store
	notes     map[string]storage.Note
	committed bool
}

// New creates a new mock store instance
func New() *Store {
	return &Store{
		notes: make(map[string]storage.Note),
	}
}

// SetError sets the error to be returned by store operations
func (s *Store) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	// Create a copy of the notes map for the transaction
	notes := make(map[string]storage.Note)
	for k, v := range s.notes {
		notes[k] = v
	}

	return &Transaction{
		store:     s,
		notes:     notes,
		committed: false,
	}, nil
}

// Add adds a new note
func (s *Store) Add(_ context.Context, note storage.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

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

	if s.err != nil {
		return storage.Note{}, s.err
	}

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

	if s.err != nil {
		return nil, s.err
	}

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

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[note.ID]; !exists {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	note.UpdatedAt = time.Now().Unix()
	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.err != nil {
		return s.err
	}

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

	if s.err != nil {
		return s.err
	}

	s.notes = make(map[string]storage.Note)
	return nil
}

// Transaction methods

// Add adds a new note in the transaction
func (tx *Transaction) Add(_ context.Context, note storage.Note) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.notes[note.ID]; exists {
		return fmt.Errorf("note with ID %s already exists", note.ID)
	}

	tx.notes[note.ID] = note
	return nil
}

// Get retrieves a note by ID in the transaction
func (tx *Transaction) Get(_ context.Context, id string) (storage.Note, error) {
	if tx.committed {
		return storage.Note{}, fmt.Errorf("transaction already committed")
	}

	note, exists := tx.notes[id]
	if !exists {
		return storage.Note{}, fmt.Errorf("note with ID %s not found", id)
	}

	return note, nil
}

// List retrieves all notes in the transaction
func (tx *Transaction) List(_ context.Context) ([]storage.Note, error) {
	if tx.committed {
		return nil, fmt.Errorf("transaction already committed")
	}

	notes := make([]storage.Note, 0, len(tx.notes))
	for _, note := range tx.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Update updates an existing note in the transaction
func (tx *Transaction) Update(_ context.Context, note storage.Note) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.notes[note.ID]; !exists {
		return fmt.Errorf("note with ID %s not found", note.ID)
	}

	note.UpdatedAt = time.Now().Unix()
	tx.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID in the transaction
func (tx *Transaction) Delete(_ context.Context, id string) error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	if _, exists := tx.notes[id]; !exists {
		return fmt.Errorf("note with ID %s not found", id)
	}

	delete(tx.notes, id)
	return nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	tx.store.mu.Lock()
	defer tx.store.mu.Unlock()

	// Copy all changes back to the store
	for k, v := range tx.notes {
		tx.store.notes[k] = v
	}

	tx.committed = true
	return nil
}

// Rollback rolls back the transaction
func (tx *Transaction) Rollback() error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}

	tx.committed = true
	return nil
}
