//go:build windows
// +build windows

package ui

func init() {
	// Windows uses the default implementation
	newSystrayManager = func() SystrayManager {
		return &defaultSystray{}
	}
}
