// Package runtime defines lifecycle orchestration boundaries for application
// startup and shutdown.
//
// Goals:
// - move process/runtime coordination out of main.go
// - centralize signal/context cancellation ownership
// - coordinate ordered shutdown across lifecycle participants
// - normalize panic and exit behavior at one runtime boundary
//
// Non-goals for this scaffold:
// - no lifecycle logic movement yet
// - no behavior changes
// - no runtime implementation in WP02
package runtime
