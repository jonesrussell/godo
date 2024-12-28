// Package gui provides graphical user interface components for the application
package gui

import "fyne.io/fyne/v2"

// MockMainWindow implements MainWindow for testing purposes
type MockMainWindow struct {
	ShowCalled   bool
	HideCalled   bool
	ContentSet   fyne.CanvasObject
	ResizeCalled bool
	CenterCalled bool
	Window       fyne.Window
}

// Show simulates showing the window and records that it was called
func (m *MockMainWindow) Show() { m.ShowCalled = true }

// Hide simulates hiding the window and records that it was called
func (m *MockMainWindow) Hide() { m.HideCalled = true }

// SetContent simulates setting window content and stores the provided content
func (m *MockMainWindow) SetContent(content fyne.CanvasObject) { m.ContentSet = content }

// Resize simulates resizing the window and records that it was called
func (m *MockMainWindow) Resize(size fyne.Size) { m.ResizeCalled = true }

// CenterOnScreen simulates centering the window and records that it was called
func (m *MockMainWindow) CenterOnScreen() { m.CenterCalled = true }

// GetWindow returns the mock window instance
func (m *MockMainWindow) GetWindow() fyne.Window { return m.Window }

// MockQuickNote implements QuickNote for testing purposes
type MockQuickNote struct {
	ShowCalled bool
	HideCalled bool
}

// Show simulates showing the quick note window and records that it was called
func (m *MockQuickNote) Show() { m.ShowCalled = true }

// Hide simulates hiding the quick note window and records that it was called
func (m *MockQuickNote) Hide() { m.HideCalled = true }
