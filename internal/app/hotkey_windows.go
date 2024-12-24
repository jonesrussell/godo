//go:build windows
// +build windows

package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"golang.design/x/hotkey"
)

func getHotkeyModifiers() []hotkey.Modifier {
	return config.GetDefaultQuickNoteModifiers()
}
