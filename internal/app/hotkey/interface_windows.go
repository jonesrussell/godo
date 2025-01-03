//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

// Package hotkey provides hotkey functionality
package hotkey

import "github.com/jonesrussell/godo/internal/common"

// Handler defines the interface for hotkey handling
type Handler interface {
	Register(*common.HotkeyBinding) error
	Unregister(*common.HotkeyBinding) error
	Start() error
	Stop() error
	IsRegistered() bool
}
