# GUI Components

## Overview

The GUI layer is built using the Fyne toolkit and consists of three main components:
- Main window for task management
- Quick note window for rapid task creation
- System tray integration

## Components

### Main Window (`internal/gui/mainwindow/`)
- Task list display
- Task completion toggling
- Task deletion
- Window state management

### Quick Note Window (`internal/gui/quicknote/`)
- Global hotkey activation
- Rapid task entry
- Auto-focus input
- ESC key handling

### System Tray (`internal/gui/systray/`)
- Application icon
- Menu items
- Quick note trigger
- Window visibility control

## Window Management

### Main Window
```go
type Window interface {
    Setup() error
    Show()
    Hide()
    Close()
}
```

### Quick Note Window
```go
type QuickNote interface {
    Setup() error
    Show()
    Hide()
    SetText(string)
}
```

## Event Handling

### Keyboard Events
- Global hotkeys (Windows-specific)
- Window-specific shortcuts
- Input field events

### Mouse Events
- Task completion clicks
- Delete button clicks
- System tray menu clicks

## Configuration

GUI configuration in `configs/default.yaml`:
```yaml
ui:
  main_window:
    width: 800
    height: 600
    start_hidden: true
  quick_note:
    width: 400
    height: 200
```

## Testing

### Unit Tests
- Window creation
- Event handling
- State management

### Integration Tests
- User interaction flows
- Window lifecycle
- System tray integration

## Platform Specifics

### Windows
- Global hotkey support
- System tray integration
- Native window management

### Future Platforms
- Linux support planned
- macOS support planned
- Platform-specific window management

## Best Practices

1. Use dependency injection for window creation
2. Handle window lifecycle properly
3. Implement proper cleanup
4. Use consistent styling
5. Follow Fyne guidelines
6. Handle platform differences appropriately

## Common Issues

1. Window Focus
   - Proper focus management
   - Z-order handling
   - Modal dialog handling

2. Event Handling
   - Event propagation
   - Key event conflicts
   - Mouse event handling

3. System Tray
   - Icon loading
   - Menu updates
   - Click handling

## Resources

- [Fyne Documentation](https://developer.fyne.io/)
- [Fyne Examples](https://github.com/fyne-io/examples)
- [Windows API Integration](https://pkg.go.dev/golang.org/x/sys/windows) 