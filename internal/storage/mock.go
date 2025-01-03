// Package storage provides interfaces and implementations for note persistence
package storage

import "github.com/jonesrussell/godo/internal/storage/mock"

// NewMockStore creates a new mock store for testing
func NewMockStore() Store {
	return mock.New()
}

// MockStore returns the concrete mock store type for advanced testing scenarios
func MockStore() *mock.Store {
	return mock.New()
}
