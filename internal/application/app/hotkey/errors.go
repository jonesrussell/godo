package hotkey

import "errors"

// Predefined errors for hotkey operations
var (
	// ErrWSL2NotSupported indicates that hotkeys are not supported in WSL2 environment
	ErrWSL2NotSupported = errors.New("hotkeys not supported in WSL2 environment")

	// ErrX11NotAvailable indicates that X11 server is not available
	ErrX11NotAvailable = errors.New("X11 server not available")

	// ErrHotkeyNotRegistered indicates that the hotkey is not registered
	ErrHotkeyNotRegistered = errors.New("hotkey not registered")

	// ErrQuickNoteNotSet indicates that the quick note service is not set
	ErrQuickNoteNotSet = errors.New("quick note service not set")
)
