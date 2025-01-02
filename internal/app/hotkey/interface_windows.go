// Package hotkey provides hotkey functionality
//go:build windows && !linux && !darwin
// +build windows,!linux,!darwin

package hotkey

import "golang.design/x/hotkey"

// hotkeyInterface defines the interface for hotkey functionality
type hotkeyInterface interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
}
