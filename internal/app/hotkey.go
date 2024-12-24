//go:build !docker
// +build !docker

package app

// HotkeyManager handles global hotkey registration and management
type HotkeyManager interface {
	Setup() error
}
