//go:build !windows && !linux && !android && !ios && !wasm && !js

package container

import "errors"

func init() {
	// This will only run if we're on an unsupported platform
	panic(errors.New("unsupported platform"))
}
