# Runtime Lifecycle Baseline

This document tracks lifecycle responsibilities while the runtime layer is being
introduced incrementally.

## Existing Signal Pattern (Current `main.go`)
- `main.go` currently creates a dedicated `sigChan` using `signal.Notify`.
- A signal goroutine waits for `SIGINT`/`SIGTERM`.
- On signal, that goroutine triggers app termination flow and a force-kill timer.

## WP03: Root Context + Signal Handling Integration
WP03 adds `internal/runtime` helpers that use `signal.NotifyContext` so callers
can own cancellation through `context.Context` instead of a bespoke signal
channel (`NewRootContext`, `WithSignals`). Integration into `main.go` is a
separate follow-up.

## WP04: Panic Recovery + Exit Normalization
Today `main.go` installs **two** separate `defer func() { recover(); ... }()`
blocks around `myapp.Run()`: both print to stdout and call `cleanup()` then
`os.Exit(1)`. That duplicates behavior and keeps panic handling outside a single
runtime boundary.

WP04 centralizes the **mechanics** of panic handling in `internal/runtime`:
- `WithPanicRecovery(logger, fn)` and `RecoverFromPanic(logger, recovered)` turn
  panics into typed errors (`RecoveredPanicError`) and log via `PanicLogger`
  (implemented by the app logger or a no-op when nil). There is no `fmt.Printf`
  and no `os.Exit` inside `internal/runtime`.
- `NormalizeExit(err)` maps errors to integer exit codes (`ExitOK`, `ExitError`,
  `ExitPanic`, `ExitForced`) so `main` can call `os.Exit` once with a consistent
  policy.

Replacing the duplicated `main.go` defer blocks with these helpers is deferred
until a later wiring change; this WP only adds the runtime API and tests.
