// Package gui defines interfaces for the graphical user interface components
package gui

import "fyne.io/fyne/v2"

// QuickNote defines the interface for quick note functionality
type QuickNote interface {
	Show()
	Hide()
}

// MainWindow defines the interface for main window functionality
type MainWindow interface {
	Show()
	Setup() error
	GetWindow() fyne.Window
}
