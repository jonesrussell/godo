package gui

// QuickNote defines the interface for quick note functionality
type QuickNote interface {
	Show()
	Hide()
}

// MainWindow defines the interface for main window functionality
type MainWindow interface {
	Show()
	Setup()
}
