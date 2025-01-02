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

// mockQuickNoteManager implements QuickNoteManager for testing
type mockQuickNoteManager struct {
	mock.Mock
}

func (m *mockQuickNoteManager) Show() {
	m.Called()
}

func (m *mockQuickNoteManager) Hide() {
	m.Called()
}

func (m *mockQuickNoteManager) CenterOnScreen() {
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
	l.On("Debug", mock.Anything, mock.Anything).Maybe()
	l.On("Info", mock.Anything, mock.Anything).Maybe()
	l.On("Error", mock.Anything, mock.Anything).Maybe()
	l.On("Warn", mock.Anything, mock.Anything).Maybe()
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

func (m *mockTestLogger) Warn(msg string, args ...interface{}) {
	m.Called(msg, args)
	m.Logger.Warn(msg, args...)
}

func TestWindowsManager_QuickNoteHotkey(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	mockQuickNote := &mockQuickNoteManager{}
	mockQuickNote.On("Show").Return()
	mockQuickNote.On("CenterOnScreen").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Create mock hotkey
	mockHk := &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift},
		key:         hotkey.KeyG,
	}
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.SetHotkey(mockHk)

	// Set quick note service
	manager.SetQuickNote(mockQuickNote, binding)

	// Register hotkey
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Start the manager
	err = manager.Start()
	assert.NoError(t, err, "Should start manager without error")

	// Verify registration state
	assert.True(t, mockHk.IsRegistered(), "Hotkey should be registered")

	// Simulate hotkey trigger
	mockHk.SimulateKeyPress()

	// Stop and cleanup
	err = manager.Stop()
	assert.NoError(t, err, "Should stop manager without error")

	// Verify final state
	assert.False(t, mockHk.IsRegistered(), "Hotkey should be unregistered")
	mockQuickNote.AssertExpectations(t)
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_InvalidKey(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	mockQuickNote := &mockQuickNoteManager{}
	mockQuickNote.On("CenterOnScreen").Return()

	// Create hotkey binding with invalid key
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "InvalidKey",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Set quick note service
	manager.SetQuickNote(mockQuickNote, binding)

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
	mockQuickNote := &mockQuickNoteManager{}
	mockQuickNote.On("CenterOnScreen").Return()

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Set quick note service without binding
	manager.SetQuickNote(mockQuickNote, nil)

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
	mockQuickNote := &mockQuickNoteManager{}
	mockQuickNote.On("CenterOnScreen").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Create mock hotkey
	mockHk := &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift},
		key:         hotkey.KeyG,
	}
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.SetHotkey(mockHk)

	// Set quick note service
	manager.SetQuickNote(mockQuickNote, binding)

	// Register hotkey
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	// Verify registration state
	assert.True(t, mockHk.IsRegistered(), "Hotkey should be registered")

	// Unregister hotkey
	err = manager.Unregister()
	assert.NoError(t, err, "Should unregister hotkey without error")
	assert.False(t, mockHk.IsRegistered(), "Hotkey should be unregistered")

	// Verify expectations
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_MultipleRegistrations(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	mockQuickNote := &mockQuickNoteManager{}
	mockQuickNote.On("CenterOnScreen").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Create mock hotkey
	mockHk := &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift},
		key:         hotkey.KeyG,
	}
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.SetHotkey(mockHk)

	// Set quick note service
	manager.SetQuickNote(mockQuickNote, binding)

	// First registration should succeed
	err = manager.Register()
	assert.NoError(t, err, "First registration should succeed")

	// Second registration should fail
	err = manager.Register()
	assert.Error(t, err, "Second registration should fail")
	assert.Contains(t, err.Error(), "already registered", "Error should mention already registered")

	// Clean up
	err = manager.Stop()
	assert.NoError(t, err, "Should stop manager without error")

	// Verify expectations
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_StopWithoutStart(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Stop without starting should not error (safe to call)
	err = manager.Stop()
	assert.NoError(t, err, "Stop without start should not error")
}

func TestWindowsManager_StartWithoutRegister(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	quickNote := &mockQuickNoteManager{}

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

	// Start without registering should fail
	err = manager.Start()
	assert.Error(t, err, "Start without register should fail")
	assert.Contains(t, err.Error(), "not registered", "Error should mention not registered")
}

func TestWindowsManager_UnregisterWithoutRegister(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Unregister without registering should not error (no-op)
	err = manager.Unregister()
	assert.NoError(t, err, "Unregister without register should not error")
}

func TestWindowsManager_HotkeyTrigger(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	quickNote := &mockQuickNoteManager{}
	quickNote.On("Show").Return()

	// Create hotkey binding
	binding := &common.HotkeyBinding{
		Modifiers: []string{"Ctrl", "Shift"},
		Key:       "G",
	}

	// Create manager
	manager, err := NewWindowsManager(log)
	assert.NoError(t, err, "Should create manager without error")

	// Create mock hotkey
	mockHk := &mockHotkey{
		keydownChan: make(chan hotkey.Event),
		modifiers:   []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift},
		key:         hotkey.KeyG,
	}
	mockHk.On("Register").Return(nil)
	mockHk.On("Unregister").Return(nil)
	manager.SetHotkey(mockHk)

	// Set quick note service
	manager.SetQuickNote(quickNote, binding)

	// Register and start
	err = manager.Register()
	assert.NoError(t, err, "Should register hotkey without error")

	err = manager.Start()
	assert.NoError(t, err, "Should start manager without error")

	// Verify registration state
	assert.True(t, mockHk.IsRegistered(), "Hotkey should be registered")

	// Simulate hotkey trigger
	mockHk.SimulateKeyPress()

	// Stop and cleanup
	err = manager.Stop()
	assert.NoError(t, err, "Should stop manager without error")

	// Verify expectations
	quickNote.AssertExpectations(t)
	mockHk.AssertExpectations(t)
}

func TestWindowsManager_StateTransitions(t *testing.T) {
	// Create mock logger
	log := newMockTestLogger(t)

	// Create mock quick note service
	quickNote := &mockQuickNoteManager{}
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
	quickNote := &mockQuickNoteManager{}

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
