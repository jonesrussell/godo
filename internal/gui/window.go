package gui

import "fyne.io/fyne/v2"

// WindowImplementation provides the concrete implementation of window management
type WindowImplementation struct {
	window fyne.Window
}

// NewWindow creates a new window implementation
func NewWindow(window fyne.Window) *WindowImplementation {
	return &WindowImplementation{
		window: window,
	}
}

// Show displays the window
func (w *WindowImplementation) Show() {
	w.window.Show()
}

// Hide hides the window
func (w *WindowImplementation) Hide() {
	w.window.Hide()
}

// Close closes the window
func (w *WindowImplementation) Close() {
	w.window.Close()
}

// GetWindow returns the underlying fyne.Window
func (w *WindowImplementation) GetWindow() fyne.Window {
	return w.window
}

// SetOnClosed sets the callback to be called when the window is closed
func (w *WindowImplementation) SetOnClosed(callback func()) {
	w.window.SetOnClosed(callback)
}
