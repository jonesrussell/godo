package hotkey

import "context"

// HotkeyManager defines the interface for platform-specific hotkey implementations
type HotkeyManager interface {
	Start(ctx context.Context) error
	Stop() error
	GetEventChannel() <-chan struct{}
}

// BaseHotkeyConfig holds the common configuration for hotkeys
type BaseHotkeyConfig struct {
	ID        int
	Key       uint
	Modifiers uint
}

// Variable to hold the platform-specific hotkey manager constructor
var newPlatformHotkeyManager = func() (HotkeyManager, error) {
	return &defaultHotkeyManager{
		eventChan: make(chan struct{}),
	}, nil
}

// NewHotkeyManager creates a new platform-specific hotkey manager
func NewHotkeyManager() (HotkeyManager, error) {
	return newPlatformHotkeyManager()
}

// defaultHotkeyManager provides a default implementation
type defaultHotkeyManager struct {
	eventChan chan struct{}
}

func (m *defaultHotkeyManager) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (m *defaultHotkeyManager) Stop() error {
	return nil
}

func (m *defaultHotkeyManager) GetEventChannel() <-chan struct{} {
	return m.eventChan
}
