// Package gui provides mock implementations for testing
package gui

import "fyne.io/fyne/v2"

// MockWindowManager implements WindowManager for testing
type MockWindowManager struct {
	ShowCalled   bool
	HideCalled   bool
	CenterCalled bool
}

func (m *MockWindowManager) Show() {
	m.ShowCalled = true
}

func (m *MockWindowManager) Hide() {
	m.HideCalled = true
}

func (m *MockWindowManager) CenterOnScreen() {
	m.CenterCalled = true
}

// MockContentManager implements ContentManager for testing
type MockContentManager struct {
	ContentSet fyne.CanvasObject
}

func (m *MockContentManager) SetContent(content fyne.CanvasObject) {
	m.ContentSet = content
}

// MockSizeManager implements SizeManager for testing
type MockSizeManager struct {
	ResizeCalled bool
	LastSize     fyne.Size
}

func (m *MockSizeManager) Resize(size fyne.Size) {
	m.ResizeCalled = true
	m.LastSize = size
}

// MockWindowAccessor implements WindowAccessor for testing
type MockWindowAccessor struct {
	Window fyne.Window
}

func (m *MockWindowAccessor) GetWindow() fyne.Window {
	return m.Window
}

// MockMainWindowManager implements MainWindowManager for testing
type MockMainWindowManager struct {
	*MockWindowManager
	*MockContentManager
	*MockSizeManager
	*MockWindowAccessor
}

// NewMockMainWindowManager creates a new mock main window manager
func NewMockMainWindowManager(window fyne.Window) *MockMainWindowManager {
	return &MockMainWindowManager{
		MockWindowManager:  &MockWindowManager{},
		MockContentManager: &MockContentManager{},
		MockSizeManager:    &MockSizeManager{},
		MockWindowAccessor: &MockWindowAccessor{Window: window},
	}
}

// MockQuickNoteManager implements QuickNoteManager for testing
type MockQuickNoteManager struct {
	*MockWindowManager
}

// NewMockQuickNoteManager creates a new mock quick note manager
func NewMockQuickNoteManager() *MockQuickNoteManager {
	return &MockQuickNoteManager{
		MockWindowManager: &MockWindowManager{},
	}
}
