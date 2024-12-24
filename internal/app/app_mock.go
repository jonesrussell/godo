package app

import (
	"fyne.io/fyne/v2"
)

type MockUI struct {
	content fyne.CanvasObject
	size    fyne.Size
	show    bool
}

func (m *MockUI) Show() {
	m.show = true
}

func (m *MockUI) Hide() {
	m.show = false
}

func (m *MockUI) SetContent(content fyne.CanvasObject) {
	m.content = content
}

func (m *MockUI) Resize(size fyne.Size) {
	m.size = size
}

func (m *MockUI) CenterOnScreen() {
	// No-op for mock
}

// Getter methods for testing
func (m *MockUI) IsShown() bool {
	return m.show
}

func (m *MockUI) Content() fyne.CanvasObject {
	return m.content
}

func (m *MockUI) Size() fyne.Size {
	return m.size
}

// MockApplication implements Application interface for testing
type MockApplication struct {
	setupUICalled bool
	runCalled     bool
	cleanupCalled bool
}

func (m *MockApplication) SetupUI() { m.setupUICalled = true }
func (m *MockApplication) Run()     { m.runCalled = true }
func (m *MockApplication) Cleanup() { m.cleanupCalled = true }

// Test helper methods
func (m *MockApplication) WasSetupUICalled() bool { return m.setupUICalled }
func (m *MockApplication) WasRunCalled() bool     { return m.runCalled }
func (m *MockApplication) WasCleanupCalled() bool { return m.cleanupCalled }
