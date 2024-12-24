//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package container

import "errors"

func init() {
	// This will only run if we're on an unsupported platform
	panic(errors.New("unsupported platform"))
}
