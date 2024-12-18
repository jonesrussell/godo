package hotkey

import "context"

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	eventChan chan struct{}
}

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register(keys []string) error
	Unregister() error
}

// GetEventChannel returns the channel that emits hotkey events
func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

// Cleanup performs any necessary cleanup
func (h *HotkeyManager) Cleanup() error {
	return nil
}

// Start begins listening for hotkey events
func (h *HotkeyManager) Start(ctx context.Context) error {
	// This will be implemented by platform-specific files
	return nil
}
