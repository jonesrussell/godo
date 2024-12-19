# Simple Todo Integration Plan

## Current State (main.go)
- [x] Basic Fyne app with quick note functionality
- [x] Main window hidden by default
- [x] System tray integration with proper icons
- [x] Lifecycle logging implemented

## Step 1: Add Basic Logging
- [x] Add zap logger initialization
- [x] Add basic logging to track application lifecycle
- [x] Log startup, shutdown, and main operations
- [x] Reference: `internal/logger/logger.go`

## Step 2: Add System Tray
- [x] Add system tray icon (using favicon.ico)
- [x] Add application icon (using Icon.png)
- [x] Hide main window by default
- [x] Move quick note trigger to system tray menu
- [x] Add logging for system tray events
- [x] Reference: `cmd/godo/main.go`

## Step 3: Refactor Quick Note
- [ ] Move quick note logic to separate package
- [ ] Keep the same functionality but make it callable from system tray
- [ ] Add logging for quick note operations (open, save, cancel)
- [ ] Reference: `internal/gui/quicknote.go`

## Step 4: Basic Todo Storage
- [ ] Add simple in-memory todo storage initially
- [ ] Create basic todo model
- [ ] Add logging for todo operations (create, read, update, delete)
- [ ] Reference: `internal/model/todo.go` (lines 9-16)

## Step 5: Todo UI Components
- [ ] Add minimal todo list view
- [ ] Add task input field
- [ ] Add complete/delete actions
- [ ] Keep it hidden by default, accessible from tray
- [ ] Add logging for UI interactions

## Step 6: Persistence
- [ ] Add SQLite storage
- [ ] Implement basic CRUD operations
- [ ] Add logging for database operations
- [ ] Keep it simple initially

## Step 7: Polish
- [ ] Add keyboard shortcuts
- [ ] Improve UI layout
- [ ] Add basic error handling
- [ ] Enhance logging with contextual information
- [ ] Add log rotation

## Implementation Plan

### Step 1 Implementation
Let's modify main.go to add logging and system tray support:

```go:main.go
package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/driver/desktop"
    "go.uber.org/zap"
)

var logger *zap.Logger

func init() {
    var err error
    logger, err = zap.NewDevelopment()
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
}

func main() {
    logger.Info("Starting application")
    
    myApp := app.New()
    mainWindow := myApp.NewWindow("Godo")
    
    // Hide main window initially
    mainWindow.Hide()
    logger.Debug("Main window hidden")
    
    // Set up system tray
    if desk, ok := myApp.(desktop.App); ok {
        menu := fyne.NewMenu("Godo",
            fyne.NewMenuItem("Quick Note", showQuickNote),
            fyne.NewMenuItemSeparator(),
            fyne.NewMenuItem("Quit", func() {
                logger.Info("Application shutdown requested")
                myApp.Quit()
            }),
        )
        desk.SetSystemTrayMenu(menu)
        logger.Info("System tray initialized")
    }
    
    logger.Info("Starting main event loop")
    myApp.Run()
}

func showQuickNote() {
    logger.Debug("Opening quick note window")
    // Keep existing quick note implementation
    // Will be refactored in Step 3
}