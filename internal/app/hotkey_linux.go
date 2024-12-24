//go:build linux
// +build linux

package app

import (
	"github.com/jonesrussell/godo/internal/config"
	"golang.design/x/hotkey"
)

func getHotkeyModifiers() []hotkey.Modifier {
	return config.GetDefaultQuickNoteModifiers()
}
