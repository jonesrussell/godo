//go:build linux
// +build linux

package ui

func init() {
	// Linux uses the default implementation
	newSystrayManager = func() SystrayManager {
		return &defaultSystray{}
	}
}
