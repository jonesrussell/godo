package testutil

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockWindow is a mock implementation of a window
type MockWindow struct {
	mock.Mock
}

// Show displays the mock window
func (m *MockWindow) Show() {
	m.Called()
}

// Hide hides the mock window
func (m *MockWindow) Hide() {
	m.Called()
}

// MockHotkeyManager is a mock implementation of a hotkey manager
type MockHotkeyManager struct {
	mock.Mock
}

// Register registers the mock hotkey
func (m *MockHotkeyManager) Register() error {
	args := m.Called()
	return args.Error(0)
}

// Start starts the mock hotkey manager
func (m *MockHotkeyManager) Start() error {
	args := m.Called()
	return args.Error(0)
}

// Stop stops the mock hotkey manager
func (m *MockHotkeyManager) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// MockQuickNote is a mock implementation of a quick note window
type MockQuickNote struct {
	mock.Mock
}

// Show displays the mock quick note window
func (m *MockQuickNote) Show() {
	m.Called()
}

// Hide hides the mock quick note window
func (m *MockQuickNote) Hide() {
	m.Called()
}

// WithMockExpectations sets up mock expectations for testing
func WithMockExpectations(_ *testing.T, _ *TestFixture) {
	// Set up expectations here if needed
}
