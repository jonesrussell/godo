// Package storage provides interfaces and implementations for note persistence
package storage

import "github.com/jonesrussell/godo/internal/storage/testing"

// NewMockStore creates a new mock store for testing
func NewMockStore() Store {
	return testing.NewMockStore()
}

// MockStore returns the concrete mock store type for advanced testing scenarios
func MockStore() *testing.MockStore {
	return testing.NewMockStore()
}
