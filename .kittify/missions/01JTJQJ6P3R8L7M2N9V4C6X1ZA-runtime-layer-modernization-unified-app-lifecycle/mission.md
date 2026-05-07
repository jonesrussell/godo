# Runtime Layer Modernization & Unified App Lifecycle

## Mission Metadata
- Mission ID: `01JTJQJ6P3R8L7M2N9V4C6X1ZA`
- Mission Slug: `runtime-layer-modernization-unified-app-lifecycle`
- Mission Type: `software-dev`
- Status: `draft`
- Scope: `lifecycle-architecture-only`

## Mission Goal
Refactor startup and shutdown orchestration by extracting lifecycle ownership from `main.go` into a dedicated runtime layer that is context-driven, testable, and extensible for future lifecycle participants.

## In Scope
- Replace manual signal handling with `signal.NotifyContext` and a single root `context.Context`.
- Remove duplicate panic-recovery paths from `main.go`.
- Remove direct `os.Exit` usage outside the final exit point in `main.go`.
- Remove runtime type assertions (for example, `myapp.(*core.App)`) by introducing a lifecycle interface contract.
- Move force-kill timeout behavior into the runtime layer.
- Define a unified shutdown sequence for:
  - UI shutdown
  - hotkey shutdown
  - storage cleanup
  - audit flush (future extension point)
  - harness shutdown (future extension point)
- Expose one runtime entrypoint (`runtime.Run(app)` or equivalent).
- Reduce `main.go` to: initialize DI -> call runtime -> return exit code.

## Out of Scope
- Business logic changes.
- Feature work in UI components.
- Storage behavior changes beyond lifecycle boundary wiring.
- New product features unrelated to lifecycle architecture.

## Constraints
- Do not modify runtime code in this mission PR.
- Do not include generated files in commits.
- Keep this mission strictly focused on lifecycle architecture.
- Implementation lands in follow-up task PRs only.

## Deliverables
- Mission bundle under `.kittify/missions/` containing:
  - `mission.md`
  - `plan.md`
  - `tasks.md`
- Draft PR branch `mission/runtime-layer` containing mission bundle only.

## Mission Complete

### Architectural Improvements
- Introduced a dedicated `internal/runtime` layer for lifecycle responsibilities:
  - root context + signal integration
  - panic-to-error recovery
  - coordinated shutdown with force-kill timeout policy
  - exit-code normalization
  - single runtime entrypoint (`runtime.Run(...)`)
- Added lifecycle contracts so runtime integration no longer depends on concrete type assertions in `main.go`.
- Added runtime smoke/testing harness coverage for lifecycle orchestration behavior.

### Lifecycle Before/After
- **Before**: `main.go` owned signal channels, duplicate panic recovery blocks, manual force-kill goroutine, ad hoc shutdown flow, and multiple early-exit paths.
- **After**: `main.go` is a thin bootstrap wrapper (DI init -> runtime root context -> `runtime.Run(...)` -> single final `os.Exit(code)`), while lifecycle orchestration is centralized in runtime.

### Future Harness Integration Notes
- Runtime now exposes stable orchestration seams for deeper lifecycle harness testing.
- Next harness expansion can add scenario matrices for ordered participant shutdown (UI/hotkey/storage/audit/harness), timeout fault injection, and integration-level cancellation behavior.
