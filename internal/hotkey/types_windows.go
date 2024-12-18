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

// DefaultConfig provides default hotkey configuration for Windows
var DefaultConfig = BaseHotkeyConfig{
	ID:        1,
	Modifiers: MOD_CONTROL | MOD_ALT,
	Key:       'G',
}
