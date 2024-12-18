package ui

// SystrayManager defines the interface for platform-specific systray implementations
type SystrayManager interface {
	Setup() error
}

// SetupSystray is a convenience function for setting up the systray
func SetupSystray() error {
	manager := newSystrayManager()
	return manager.Setup()
}

// Variable to hold the platform-specific systray manager constructor
var newSystrayManager = func() SystrayManager {
	return &defaultSystray{}
}

// defaultSystray provides a fallback implementation
type defaultSystray struct{}

func (s *defaultSystray) Setup() error {
	return nil // No-op implementation
}
