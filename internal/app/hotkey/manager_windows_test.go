//go:build windows && !linux && !darwin && !docker
// +build windows,!linux,!darwin,!docker

package hotkey

import (
	"fmt"
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

// mockTestLogger combines both test logger and mock functionality
type mockTestLogger struct {
	logger.Logger
	mock.Mock
	t *testing.T
}

func newMockTestLogger(t *testing.T) *mockTestLogger {
	l := &mockTestLogger{
		Logger: logger.NewTestLogger(t),
		t:      t,
	}
	// Set up default expectations
	l.On("Debug", mock.Anything, mock.Anything).Return()
	l.On("Info", mock.Anything, mock.Anything).Return()
	l.On("Error", mock.Anything, mock.Anything).Return()
	return l
}

func (m *mockTestLogger) Debug(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Debug(msg, args...)
}

func (m *mockTestLogger) Info(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Info(msg, args...)
}

func (m *mockTestLogger) Error(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Error(msg, args...)
}

func TestWindowsManager_QuickNoteHotkey(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

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
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_InvalidKey(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

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
	log := newMockTestLogger(t)

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
	log := newMockTestLogger(t)

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
	log := newMockTestLogger(t)

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

	// First registration should succeed
	err = manager.Register()
	assert.NoError(t, err, "First registration should succeed")

	// Second registration should fail
	err = manager.Register()
	assert.Error(t, err, "Second registration should fail")
	assert.Contains(t, err.Error(), "already registered", "Error should mention already registered")

	// Clean up
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify expectations
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_StopWithoutStart(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Stop without starting should return error
	err = manager.Stop()
	assert.Error(t, err, "Stop without start should return error")
	assert.Contains(t, err.Error(), "not started", "Error should mention not started")
}

func TestWindowsManager_StartWithoutRegister(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

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
	log := newMockTestLogger(t)

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

func TestWindowsManager_HotkeyTrigger(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service with expectations
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

	// Create mock hotkey that will trigger callback
	mockHk := newMockHotkey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG).(*mockHotkey)
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.hk = mockHk

	// Register should succeed
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Start the manager
	err = manager.Start()
	assert.NoError(t, err, "Should start manager without error")

	// Simulate hotkey trigger
	mockHk.SimulateKeyPress()

	// Stop and cleanup
	err = manager.Stop()
	assert.NoError(t, err, "Should stop manager without error")

	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify the Show method was called
	quickNote.AssertCalled(t, "Show")
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_StateTransitions(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

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

	// Create mock hotkey
	mockHk := newMockHotkey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG).(*mockHotkey)
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.hk = mockHk

	// Test valid state transitions
	// 1. Register
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// 2. Start
	err = manager.Start()
	assert.NoError(t, err, "Should start manager without error")

	// 3. Stop
	err = manager.Stop()
	assert.NoError(t, err, "Should stop manager without error")

	// 4. Unregister
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")

	// Verify all transitions completed
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_SystemErrors(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

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

	// Create mock hotkey that returns errors
	mockHk := newMockHotkey([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyG).(*mockHotkey)
	mockHk.On("Register").Return(fmt.Errorf("system error: hotkey registration failed"))
	manager.hk = mockHk

	// Register should fail with system error
	err = manager.Register()
	assert.Error(t, err, "Should fail to register hotkey")
	assert.Contains(t, err.Error(), "system error", "Error should contain system error message")

	// Verify expectations
	mockHk.AssertExpectations(t)
}
