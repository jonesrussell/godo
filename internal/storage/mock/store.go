// Package mock provides mock implementations for testing
package mock

import (
	"context"
	"sync"

	"github.com/jonesrussell/godo/internal/storage"
	"github.com/jonesrussell/godo/internal/storage/errors"
)

// Store implements the storage.Store interface for testing
type Store struct {
	sync.RWMutex
	notes map[string]storage.Note
	err   error
}

// New creates a new mock store
func New() *Store {
	return &Store{
		notes: make(map[string]storage.Note),
	}
}

// SetError sets the error to be returned by store operations
func (s *Store) SetError(err error) {
	s.Lock()
	defer s.Unlock()
	s.err = err
}

// List retrieves all notes
func (s *Store) List(_ context.Context) ([]storage.Note, error) {
	s.RLock()
	defer s.RUnlock()

	if s.err != nil {
		return nil, s.err
	}

	notes := make([]storage.Note, 0, len(s.notes))
	for _, note := range s.notes {
		notes = append(notes, note)
	}

	return notes, nil
}

// Get retrieves a note by ID
func (s *Store) Get(_ context.Context, id string) (storage.Note, error) {
	s.RLock()
	defer s.RUnlock()

	if s.err != nil {
		return storage.Note{}, s.err
	}

	note, exists := s.notes[id]
	if !exists {
		return storage.Note{}, errors.ErrNoteNotFound
	}

	return note, nil
}

// Add adds a new note
func (s *Store) Add(_ context.Context, note storage.Note) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[note.ID]; exists {
		return errors.ErrNoteExists
	}

	s.notes[note.ID] = note
	return nil
}

// Update updates an existing note
func (s *Store) Update(_ context.Context, note storage.Note) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[note.ID]; !exists {
		return errors.ErrNoteNotFound
	}

	s.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID
func (s *Store) Delete(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	if _, exists := s.notes[id]; !exists {
		return errors.ErrNoteNotFound
	}

	delete(s.notes, id)
	return nil
}

// BeginTx starts a new transaction
func (s *Store) BeginTx(_ context.Context) (storage.Transaction, error) {
	s.Lock()
	defer s.Unlock()

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

// Close cleans up any resources
func (s *Store) Close() error {
	s.Lock()
	defer s.Unlock()

	if s.err != nil {
		return s.err
	}

	s.notes = make(map[string]storage.Note)
	return nil
}

// Transaction represents a mock transaction
type Transaction struct {
	store     *Store
	notes     map[string]storage.Note
	committed bool
}

// Add adds a new note in the transaction
func (tx *Transaction) Add(_ context.Context, note storage.Note) error {
	if tx.committed {
		return errors.ErrTransactionClosed
	}

	if _, exists := tx.notes[note.ID]; exists {
		return errors.ErrNoteExists
	}

	tx.notes[note.ID] = note
	return nil
}

// Get retrieves a note by ID in the transaction
func (tx *Transaction) Get(_ context.Context, id string) (storage.Note, error) {
	if tx.committed {
		return storage.Note{}, errors.ErrTransactionClosed
	}

	note, exists := tx.notes[id]
	if !exists {
		return storage.Note{}, errors.ErrNoteNotFound
	}

	return note, nil
}

// List retrieves all notes in the transaction
func (tx *Transaction) List(_ context.Context) ([]storage.Note, error) {
	if tx.committed {
		return nil, errors.ErrTransactionClosed
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
		return errors.ErrTransactionClosed
	}

	if _, exists := tx.notes[note.ID]; !exists {
		return errors.ErrNoteNotFound
	}

	tx.notes[note.ID] = note
	return nil
}

// Delete removes a note by ID in the transaction
func (tx *Transaction) Delete(_ context.Context, id string) error {
	if tx.committed {
		return errors.ErrTransactionClosed
	}

	if _, exists := tx.notes[id]; !exists {
		return errors.ErrNoteNotFound
	}

	delete(tx.notes, id)
	return nil
}

// Commit commits the transaction
func (tx *Transaction) Commit() error {
	if tx.committed {
		return errors.ErrTransactionClosed
	}

	tx.store.Lock()
	defer tx.store.Unlock()

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
		return errors.ErrTransactionClosed
	}

	tx.committed = true
	return nil
}
