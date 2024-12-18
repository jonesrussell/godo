# Godo Development Tasks

## Immediate Tasks
- [x] Run wire command to generate DI implementation  ```bash
  go install github.com/google/wire/cmd/wire@latest
  cd internal/app && wire  ```
- [x] Add necessary dependencies to go.mod
  - github.com/mattn/go-sqlite3
  - github.com/google/wire
  - go.uber.org/zap
  - fyne.io/fyne/v2

## User Interface
- [x] Research and choose GUI library (Fyne)
- [x] Design minimal UI layout
- [ ] Implement basic UI with Fyne
  - [ ] Main window layout
  - [ ] Task list view
  - [ ] Task input form
  - [ ] Task actions (complete/delete)
- [x] Implement system tray integration
  - [x] Add system tray icon
  - [x] Add basic menu (Open Manager, Quit)
  - [x] Implement clean shutdown
- [x] Add global hotkey support
  - [x] Research cross-platform hotkey libraries
  - [x] Implement hotkey registration
  - [ ] Add user-configurable shortcuts

## Future Enhancements
- [ ] Terminal User Interface (TUI) alternative
  - [ ] Research TUI libraries (Bubble Tea, etc.)
  - [ ] Design TUI layout
  - [ ] Implement basic TUI functionality
  - [ ] Ensure feature parity with GUI
- [ ] Task categories/tags
- [ ] Due dates and reminders
- [ ] Data export/import
- [ ] Task priority levels
- [ ] Recurring tasks
- [ ] Multiple todo lists
- [ ] Cloud sync support

## System Tray Integration (Completed)
- [x] Add system tray icon support
- [x] Create basic menu structure
- [x] Implement proper shutdown handling
- [x] Add Open Manager option
- [x] Handle icon loading
- [x] Implement clean exit
