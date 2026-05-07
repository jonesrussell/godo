# Plan: Runtime Layer Modernization & Unified App Lifecycle

## Objective
Create a runtime orchestration layer that owns application lifecycle responsibilities currently concentrated in `main.go`, while preserving behavior and minimizing risk.

## Architecture Direction
- Introduce a runtime package with a single public entrypoint: `Run(app LifecycleApp) int` (exact naming finalized during implementation).
- Define a narrow lifecycle interface to eliminate concrete type assertions.
- Centralize root context creation, signal cancellation, panic recovery, timeout-based forced termination, and cleanup ordering in runtime.
- Keep dependency injection and process exit ownership in `main.go`.

## Planned Runtime Responsibilities
1. Build root context via `signal.NotifyContext`.
2. Execute app lifecycle startup.
3. Listen for cancellation or fatal lifecycle errors.
4. Run orderly shutdown pipeline (idempotent and best-effort):
   - UI
   - hotkeys
   - storage
   - audit flush hook (future)
   - harness hook (future)
5. Enforce force-kill timeout policy from runtime.
6. Return normalized exit code to `main.go`.

## Integration Boundaries
- `main.go` remains responsible for:
  - DI/container initialization
  - invoking runtime
  - final process exit
- Runtime remains responsible for:
  - lifecycle control flow
  - cancellation propagation
  - cleanup orchestration
  - panic recovery normalization

## Risks and Mitigations
- Risk: shutdown ordering regressions.
  - Mitigation: explicit lifecycle contract and sequence tests in follow-up PRs.
- Risk: hidden reliance on concrete app type.
  - Mitigation: define interface-first boundary and migrate call sites.
- Risk: timeout behavior drift.
  - Mitigation: preserve current timeout defaults while relocating ownership.

## Non-Goals
- No runtime implementation changes in this mission bundle PR.
- No cross-cutting refactors outside startup/shutdown lifecycle path.
