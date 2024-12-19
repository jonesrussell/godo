# Simple Todo Integration Plan

## Current State
- [x] Basic Fyne app with quick note functionality
- [x] Main window hidden by default
- [x] System tray integration with proper icons
- [x] Lifecycle logging implemented
- [x] Quick note in separate package
- [x] SQLite persistence working

## Step 1: Add Basic Logging ✅
- [x] Add zap logger initialization
- [x] Add basic logging to track application lifecycle
- [x] Log startup, shutdown, and main operations
- [x] Reference: `internal/logger/logger.go`

## Step 2: Add System Tray ✅
- [x] Add system tray icon (using favicon.ico)
- [x] Add application icon (using Icon.png)
- [x] Hide main window by default
- [x] Move quick note trigger to system tray menu
- [x] Add logging for system tray events
- [x] Reference: `cmd/godo/main.go`

## Step 3: Refactor Quick Note ✅
- [x] Move quick note logic to separate package
- [x] Keep the same functionality but callable from system tray
- [x] Add logging for quick note operations
- [x] Reference: `internal/gui/quicknote/quicknote.go`

## Step 4: Basic Todo Storage ✅
- [x] Add simple in-memory todo storage
- [x] Create basic todo model with UUID and timestamps
- [x] Add logging for todo operations
- [x] Reference: `internal/model/todo.go`, `internal/storage/memory/memory.go`

## Step 5: Persistence ✅
- [x] Add SQLite storage implementation
- [x] Implement basic CRUD operations
- [x] Add migration support
- [x] Add logging for database operations
- [x] Reference: `internal/storage/sqlite/sqlite.go`

## Step 6: Polish 👈 (Next)
- [ ] Add keyboard shortcuts
- [ ] Improve UI layout
- [ ] Add basic error handling
- [ ] Enhance logging with contextual information
- [ ] Add log rotation
- [ ] Add auto-start capability
- [ ] Add update mechanism

## Step 7: Todo List UI
- [ ] Add a list view to display all todos
- [ ] Add ability to mark todos as done
- [ ] Add ability to delete todos
- [ ] Show todo creation time and last update time
- [ ] Reference: `internal/gui/todolist/todolist.go`

Would you like to proceed with Step 6 and start polishing the application?

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