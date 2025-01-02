//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

// Package hotkey provides hotkey functionality
package hotkey

import "github.com/jonesrussell/godo/internal/common"

// HotkeyHandler defines the interface for Windows hotkey functionality
type HotkeyHandler interface {
	Register(binding *common.HotkeyBinding) error
	Unregister(binding *common.HotkeyBinding) error
	Start() error
	Stop() error
}
