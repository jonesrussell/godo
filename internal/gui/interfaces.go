// Package gui defines interfaces for the graphical user interface components
package gui

import "fyne.io/fyne/v2"

// WindowManager defines the core window management capabilities
type WindowManager interface {
	Show()
	Hide()
	CenterOnScreen()
}

// ContentManager defines content management capabilities
type ContentManager interface {
	SetContent(content fyne.CanvasObject)
}

// SizeManager defines window size management capabilities
type SizeManager interface {
	Resize(size fyne.Size)
}

// WindowAccessor defines window access capabilities
type WindowAccessor interface {
	GetWindow() fyne.Window
}

// QuickNoteManager defines quick note functionality
type QuickNoteManager interface {
	WindowManager
}

// MainWindowManager defines the complete main window functionality by composing smaller interfaces
type MainWindowManager interface {
	WindowManager
	ContentManager
	SizeManager
	WindowAccessor
}

// MainWindow is an alias for MainWindowManager for backward compatibility
type MainWindow = MainWindowManager

// QuickNote is an alias for QuickNoteManager for backward compatibility
type QuickNote = QuickNoteManager

// QuickNoteService defines the service interface for quick notes
type QuickNoteService interface {
	Show()
	Hide()
}
