//go:build !darwin && !linux && !windows
// +build !darwin,!linux,!windows

package ui

// NewSystrayManager creates a default systray manager
func NewSystrayManager() SystrayManager {
	return &defaultSystray{}
}
