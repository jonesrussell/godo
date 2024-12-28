package storage

import "context"

// StoreAdapter adapts a TaskStore to the legacy Store interface
type StoreAdapter struct {
	store TaskStore
}

// NewStoreAdapter creates a new adapter for TaskStore
func NewStoreAdapter(store TaskStore) Store {
	return &StoreAdapter{store: store}
}

// List implements Store
func (a *StoreAdapter) List() ([]Task, error) {
	return a.store.List(context.Background())
}

// Add implements Store
func (a *StoreAdapter) Add(task Task) error {
	return a.store.Add(context.Background(), task)
}

// Update implements Store
func (a *StoreAdapter) Update(task Task) error {
	return a.store.Update(context.Background(), task)
}

// Delete implements Store
func (a *StoreAdapter) Delete(id string) error {
	return a.store.Delete(context.Background(), id)
}

// GetByID implements Store
func (a *StoreAdapter) GetByID(id string) (*Task, error) {
	return a.store.GetByID(context.Background(), id)
}

// Close implements Store
func (a *StoreAdapter) Close() error {
	return a.store.Close()
}
