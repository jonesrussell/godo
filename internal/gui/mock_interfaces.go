package gui

import "fyne.io/fyne/v2"

// MockMainWindow implements MainWindow for testing
type MockMainWindow struct {
	ShowCalled   bool
	HideCalled   bool
	ContentSet   fyne.CanvasObject
	ResizeCalled bool
	CenterCalled bool
	Window       fyne.Window
}

func (m *MockMainWindow) Show()                                { m.ShowCalled = true }
func (m *MockMainWindow) Hide()                                { m.HideCalled = true }
func (m *MockMainWindow) SetContent(content fyne.CanvasObject) { m.ContentSet = content }
func (m *MockMainWindow) Resize(size fyne.Size)                { m.ResizeCalled = true }
func (m *MockMainWindow) CenterOnScreen()                      { m.CenterCalled = true }
func (m *MockMainWindow) GetWindow() fyne.Window               { return m.Window }

// MockQuickNote implements QuickNote for testing
type MockQuickNote struct {
	ShowCalled bool
	HideCalled bool
}

func (m *MockQuickNote) Show() { m.ShowCalled = true }
func (m *MockQuickNote) Hide() { m.HideCalled = true }
