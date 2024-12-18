//go:build linux
// +build linux

package ui

func init() {
	newSystrayManager = func() SystrayManager {
		return &linuxSystray{}
	}
}

type linuxSystray struct {
	defaultSystray
}
