//go:build windows
// +build windows

package hotkey

import (
	"syscall"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procRegisterHotKey   = user32.NewProc("RegisterHotKey")
	procUnregisterHotKey = user32.NewProc("UnregisterHotKey")
	procPeekMessage      = user32.NewProc("PeekMessageW")
)

func registerHotkey(config HotkeyConfig) (bool, error) {
	ret, _, err := procRegisterHotKey.Call(
		uintptr(config.WindowHandle),
		uintptr(config.ID),
		uintptr(config.Modifiers),
		uintptr(config.Key),
	)
	return ret != 0, err
}

func unregisterHotkey(windowHandle syscall.Handle, id int) (bool, error) {
	ret, _, err := procUnregisterHotKey.Call(
		uintptr(windowHandle),
		uintptr(id),
	)
	return ret != 0, err
}

func peekMessage(msg *MSG) (bool, error) {
	ret, _, err := procPeekMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		0,
		0,
		0,
		PM_REMOVE,
	)
	return ret != 0, err
}
