//go:build linux
// +build linux

// Package hotkey provides hotkey functionality
package hotkey

import "golang.design/x/hotkey"

// hotkeyInterface defines the interface for hotkey functionality
type hotkeyInterface interface {
	Register() error
	Unregister() error
	Keydown() <-chan hotkey.Event
	IsRegistered() bool
}
