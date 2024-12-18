//go:build windows

package hotkey

import "syscall"

type Handle syscall.Handle

// MSG represents a Windows message structure
type MSG struct {
	Hwnd    syscall.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// HotkeyConfig defines the configuration for a hotkey
type HotkeyConfig struct {
	WindowHandle syscall.Handle
	ID           int
	Modifiers    uint
	Key          rune
}
