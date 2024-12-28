// Package hotkey provides hotkey management functionality for the application
package hotkey

// MockManager implements Manager interface for testing
type MockManager struct {
	RegisterCalled   bool
	UnregisterCalled bool
	StartCalled      bool
	StopCalled       bool
}

// NewMockManager creates a new mock hotkey manager for testing
func NewMockManager() *MockManager {
	return &MockManager{}
}

// Register simulates registering a hotkey and records that it was called
func (m *MockManager) Register() error {
	m.RegisterCalled = true
	return nil
}

// Unregister simulates unregistering a hotkey and records that it was called
func (m *MockManager) Unregister() error {
	m.UnregisterCalled = true
	return nil
}

// Start simulates starting the hotkey listener and records that it was called
func (m *MockManager) Start() error {
	m.StartCalled = true
	return nil
}

// Stop simulates stopping the hotkey listener and records that it was called
func (m *MockManager) Stop() error {
	m.StopCalled = true
	return nil
}
