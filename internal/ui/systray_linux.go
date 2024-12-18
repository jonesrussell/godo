//go:build linux

package ui

func init() {
	setPlatformSpecificTitle = func() {
		// Linux typically doesn't use title
	}
}
