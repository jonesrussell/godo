//go:build windows

package ui

func init() {
	setPlatformSpecificTitle = func() {
		// Windows doesn't support SetTitle
	}
}
