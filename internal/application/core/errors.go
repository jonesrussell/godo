package core

import "errors"

// Application-specific errors
var (
	// ErrWSL2NotSupported indicates that a feature is not supported in WSL2 environment
	ErrWSL2NotSupported = errors.New("feature not supported in WSL2 environment")

	// ErrAPIServerStartTimeout indicates that the API server failed to start within timeout
	ErrAPIServerStartTimeout = errors.New("API server failed to start within timeout")

	// ErrHotkeySetupFailed indicates that hotkey setup failed
	ErrHotkeySetupFailed = errors.New("hotkey setup failed")

	// ErrUISetupFailed indicates that UI setup failed
	ErrUISetupFailed = errors.New("UI setup failed")
)
