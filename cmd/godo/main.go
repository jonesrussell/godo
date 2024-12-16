package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonesrussell/godo/internal/di"
)

func main() {
	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting Godo application...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize dependency injection container
	log.Println("Initializing dependency injection...")
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	log.Println("Dependency injection initialized successfully")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		cancel()
	}()

	log.Println("Starting application run...")
	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
