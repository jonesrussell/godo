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
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Initialize dependency injection container
    app, err := di.InitializeApp()
    if err != nil {
        log.Fatalf("Failed to initialize application: %v", err)
    }

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        cancel()
    }()

    if err := app.Run(ctx); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}
