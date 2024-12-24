package app

// HotkeyManager defines the interface for global hotkey functionality
type HotkeyManager interface {
	Register() error
	Unregister() error
}
