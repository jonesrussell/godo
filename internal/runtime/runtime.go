package runtime

// Run is the future single runtime entrypoint that will own process lifecycle
// orchestration.
//
// TODO(WP03+):
// - create root context via signal.NotifyContext
// - coordinate startup and shutdown sequencing
// - normalize panic and exit handling
// - return deterministic exit code to main
func Run(_ Lifecycle) int {
	// TODO: implement runtime lifecycle orchestration.
	return 0
}
