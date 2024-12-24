//go:build docker
// +build docker

package app

// mockHotkeyManager provides a no-op implementation for Docker environments
type mockHotkeyManager struct{}

func NewHotkeyManager() HotkeyManager {
	return &mockHotkeyManager{}
}

// No-op implementations...
