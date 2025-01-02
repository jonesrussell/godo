//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

package hotkey

import (
	"testing"

	"github.com/jonesrussell/godo/internal/common"
	"github.com/jonesrussell/godo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.design/x/hotkey"
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
	m.Called(msg, mock.Anything)
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, mock.Anything)
}

func (m *mockLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, mock.Anything)
}

func setupMockLogger() *mockLogger {
	log := &mockLogger{}
	// Accept any log message with any arguments
	log.On("Info", mock.Anything, mock.Anything).Return()
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Error", mock.Anything, mock.Anything).Return()
	return log
}

func TestWindowsManager_QuickNoteHotkey(t *testing.T) {
	// Create mock logger
	log := setupMockLogger()
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Error", mock.Anything, mock.Anything).Return()

	// Create mock quick note service
	quickNote := &mockQuickNoteService{}

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

	// Create mock hotkey
	mockHk := newMockHotkey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG).(*mockHotkey)
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.hk = mockHk

	// Register should succeed
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Unregister should succeed
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify expectations
	log.AssertExpectations(t)
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_InvalidKey(t *testing.T) {
	// Create mock logger
	log := setupMockLogger()

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
	log := setupMockLogger()

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

func TestWindowsManager_UnregisterHotkey(t *testing.T) {
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

	// Unregister hotkey
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify expectations
	log.AssertExpectations(t)
}

func TestWindowsManager_MultipleRegistrations(t *testing.T) {
	// Create mock logger
	log := setupMockLogger()

	// Create mock quick note service
	quickNote := &mockQuickNoteService{}

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

	// Create mock hotkey
	mockHk := newMockHotkey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG).(*mockHotkey)
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)

	// First registration should succeed
	manager.hk = mockHk
	err = manager.Register()
	assert.Error(t, err, "First registration should fail because hotkey is already set")
	assert.Contains(t, err.Error(), "already registered", "Error should mention already registered")

	// Clean up
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify expectations
	log.AssertExpectations(t)
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_StopWithoutStart(t *testing.T) {
	// Create mock logger
	log := setupMockLogger()

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Stop without starting should not panic but return error
	err = manager.Stop()
	assert.Error(t, err, "Stop without start should return error")
	assert.Contains(t, err.Error(), "not started", "Error should mention not started")

	// Verify expectations
	log.AssertExpectations(t)
}

func TestWindowsManager_StartWithoutRegister(t *testing.T) {
	// Create mock logger
	log := &mockLogger{}
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Error", mock.Anything, mock.Anything).Return()

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Start without registering should return error
	err = manager.Start()
	assert.Error(t, err, "Start without register should return error")
	assert.Contains(t, err.Error(), "not registered", "Error should mention not registered")

	// Verify expectations
	log.AssertExpectations(t)
}

func TestWindowsManager_UnregisterWithoutRegister(t *testing.T) {
	// Create mock logger
	log := &mockLogger{}
	log.On("Debug", mock.Anything, mock.Anything).Return()
	log.On("Error", mock.Anything, mock.Anything).Return()

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Unregister without registering should return error
	err = manager.Unregister()
	assert.Error(t, err, "Unregister without register should return error")
	assert.Contains(t, err.Error(), "not registered", "Error should mention not registered")

	// Verify expectations
	log.AssertExpectations(t)
}