# Runtime Lifecycle Baseline

This document tracks lifecycle responsibilities while the runtime layer is being
introduced incrementally.

## Existing Signal Pattern (Current `main.go`)
- `main.go` currently creates a dedicated `sigChan` using `signal.Notify`.
- A signal goroutine waits for `SIGINT`/`SIGTERM`.
- On signal, that goroutine triggers app termination flow and a force-kill timer.

## WP03: Root Context + Signal Handling Integration
WP03 introduces `internal/runtime` signal/context primitives that replace the
ad-hoc `sigChan` ownership pattern with context-first wiring:

- `NewRootContext()` creates a root context backed by `signal.NotifyContext`
  and listens for `SIGINT` + `SIGTERM`.
- `WithSignals(parent, ...)` provides a testable helper for signal-bound
  context creation.
- Callers now receive `(context.Context, context.CancelFunc)` and can compose
  lifecycle cancellation through context instead of directly managing signal
  channels.

This WP does **not** move cleanup orchestration or `main.go` wiring yet; it
only establishes the runtime-owned signal/context boundary used by later WPs.

