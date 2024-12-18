//go:build darwin

package ui

import "github.com/getlantern/systray"

func init() {
	setPlatformSpecificTitle = func() {
		systray.SetTitle("Godo")
	}
}
