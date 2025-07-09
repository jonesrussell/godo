package app

import "errors"

// Application-specific errors
var (
	// ErrWSL2NotSupported indicates that a feature is not supported in WSL2 environment
	ErrWSL2NotSupported = errors.New("feature not supported in WSL2 environment")

	// ErrDesktopFeaturesNotAvailable indicates that desktop features are not available
	ErrDesktopFeaturesNotAvailable = errors.New("desktop features not available")

	// ErrAPIServerStartTimeout indicates that the API server failed to start within timeout
	ErrAPIServerStartTimeout = errors.New("API server failed to start within timeout")

	// ErrSystraySetupFailed indicates that systray setup failed
	ErrSystraySetupFailed = errors.New("systray setup failed")

	// ErrHotkeySetupFailed indicates that hotkey setup failed
	ErrHotkeySetupFailed = errors.New("hotkey setup failed")

	// ErrUISetupFailed indicates that UI setup failed
	ErrUISetupFailed = errors.New("UI setup failed")
)
