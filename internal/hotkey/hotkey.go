package hotkey

// HotkeyManager handles global hotkey registration and events
type HotkeyManager struct {
	eventChan chan struct{}
}

// HotkeyHandler defines the interface for platform-specific hotkey implementations
type HotkeyHandler interface {
	Register(keys []string) error
	Unregister() error
}

func NewHotkeyManager() (*HotkeyManager, error) {
	return &HotkeyManager{
		eventChan: make(chan struct{}, 1),
	}, nil
}

func (h *HotkeyManager) GetEventChannel() <-chan struct{} {
	return h.eventChan
}

func (h *HotkeyManager) Cleanup() error {
	return nil
}
