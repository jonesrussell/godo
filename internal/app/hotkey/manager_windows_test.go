//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

package hotkey

import (
	"testing"
	"time"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockQuickNoteService struct {
	mock.Mock
}

func (m *mockQuickNoteService) Show() {
	m.Called()
}

func (m *mockQuickNoteService) Hide() {
	m.Called()
}

type mockLogger struct {
	logger.Logger
	mock.Mock
}

func (m *mockLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *mockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func TestWindowsManager_QuickNoteHotkey(t *testing.T) {
	// Create mock logger
	log := &mockLogger{}
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Info", mock.Anything, mock.Anything).Return()

	// Create mock quick note service
	quickNote := &mockQuickNoteService{}
	quickNote.On("Show").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Set quick note service
	manager.SetQuickNote(quickNote, binding)

	// Register hotkey
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Start hotkey listener
	err = manager.Start()
	assert.NoError(t, err, "Should start hotkey listener without error")

	// Give some time for the listener to start
	time.Sleep(100 * time.Millisecond)

	// Clean up
	err = manager.Stop()
	assert.NoError(t, err, "Should stop hotkey listener without error")

	// Verify expectations
	quickNote.AssertExpectations(t)
	log.AssertExpectations(t)
}

func TestWindowsManager_InvalidKey(t *testing.T) {
	// Create mock logger
	log := &mockLogger{}
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Error", mock.Anything, mock.Anything).Return()

	// Create mock quick note service
	quickNote := &mockQuickNoteService{}

	// Create hotkey binding with invalid key
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "InvalidKey",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Set quick note service
	manager.SetQuickNote(quickNote, binding)

	// Register hotkey should fail
	err = manager.Register()
	assert.Error(t, err, "Should fail to register hotkey with invalid key")
	assert.Contains(t, err.Error(), "unsupported key", "Error should mention unsupported key")

	// Verify expectations
	log.AssertExpectations(t)
}

func TestWindowsManager_NilBinding(t *testing.T) {
	// Create mock logger
	log := &mockLogger{}
	log.On("Debug", mock.Anything, mock.Anything).Return()

	// Create mock quick note service
	quickNote := &mockQuickNoteService{}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Set quick note service without binding
	manager.SetQuickNote(quickNote, nil)

	// Register should fail due to nil binding
	err = manager.Register()
	assert.Error(t, err, "Should fail to register with nil binding")
	assert.Contains(t, err.Error(), "binding not set", "Error should mention binding not set")

	// Verify expectations
	log.AssertExpectations(t)
}
