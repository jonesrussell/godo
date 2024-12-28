package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/jonesrussell/godo/internal/container"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("=== Godo Starting (main) ===")

	// Add panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			fmt.Printf("Stack trace:\n%s\n", debug.Stack())
			os.Exit(1)
		}
	}()

	if err := run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== Entering run() ===")

	// Try creating logger first
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return err
	}
	defer func() {
		_ = logger.Sync() // Ignore sync errors
	}()

	fmt.Println("=== Logger created successfully ===")
	fmt.Println("=== Initializing app... ===")

	// Create app
	app, cleanup, err := container.InitializeApp()
	if err != nil {
		fmt.Printf("Failed to initialize app: %v\n", err)
		logger.Error("Failed to initialize app", zap.Error(err))
		cleanup()
		return err
	}
	defer cleanup()

	fmt.Println("=== Setting up UI... ===")
	// Setup UI
	app.SetupUI()

	fmt.Println("=== Running app... ===")
	// Run app
	if err := app.Run(); err != nil {
		fmt.Printf("Failed to run app: %v\n", err)
		logger.Error("Failed to run app", zap.Error(err))
		return err
	}

	fmt.Println("=== App completed successfully ===")
	return nil
}
