package app

// UIManager defines the user interface management capabilities
type UIManager interface {
	SetupUI() error
}

// ApplicationService defines the core application functionality
type ApplicationService interface {
	UIManager
	Run()
	Cleanup()
}
