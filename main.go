// Package main is the entry point for the Godo application
package main

import (
	"os"

	"github.com/jonesrussell/godo/internal/application/container"
	"github.com/jonesrussell/godo/internal/runtime"
)

func main() {
	code := runtime.ExitOK

	app, cleanup, err := container.InitializeApp()
	if err != nil {
		code = runtime.NormalizeExit(err)
	} else {
		rootCtx, cancel := runtime.NewRootContext()
		defer cancel()
		code = runtime.Run(rootCtx, app, cleanup, nil)
	}

	os.Exit(code)
}
