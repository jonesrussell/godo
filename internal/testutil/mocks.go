package testutil

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/jonesrussell/godo/internal/app/hotkey"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/mock"
)

// MockWindow is a mock implementation of fyne.Window
type MockWindow struct {
	mock.Mock
	fyne.Window
}

func (m *MockWindow) Show() {
	m.Called()
}

func (m *MockWindow) Hide() {
	m.Called()
}

// MockHotkeyManager is a mock implementation of hotkey.Manager
type MockHotkeyManager struct {
	mock.Mock
	hotkey.Manager
}

func (m *MockHotkeyManager) Register() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockHotkeyManager) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockHotkeyManager) Stop() error {
	args := m.Called()
	return args.Error(0)
}

// MockQuickNote is a mock implementation of the QuickNote interface
type MockQuickNote struct {
	mock.Mock
}

func (m *MockQuickNote) Show() {
	m.Called()
}

func (m *MockQuickNote) Hide() {
	m.Called()
}

// WithMockExpectations sets up common mock expectations
func WithMockExpectations(t *testing.T, f *TestFixture) {
	// Setup store expectations - using concrete MockStore methods
	f.Store.Reset()       // Clear any existing state
	f.Store.SetError(nil) // Ensure no errors are set

	// Setup logger expectations
	logger := f.Logger.(*logger.MockTestLogger)
	logger.On("Debug", mock.Anything, mock.Anything).Return()
	logger.On("Info", mock.Anything, mock.Anything).Return()
	logger.On("Error", mock.Anything, mock.Anything).Return()
	logger.On("Warn", mock.Anything, mock.Anything).Return()
}
