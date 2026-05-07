# PR: refactor/hotkey-start (ISSUE-006)

## Summary

Extract `HotkeyManager.listenChannels` from the `Start` goroutine so channel construction is isolated and unit-testable without changing listener behavior.

## Changes

- `internal/infrastructure/hotkey/manager.go` — `listenChannels()` helper; `Start` calls it.
- `internal/infrastructure/hotkey/manager_listen_channels_test.go` — whitebox test for nil `hotkey` entries (build-tagged `linux || windows` like the manager).

## How to verify

```bash
go test ./internal/infrastructure/hotkey/... -tags=wireinject -count=1 -v
go test ./... -tags=wireinject -count=1
```

## Acceptance criteria

- [ ] `Start` still builds the same channel slice as before for registered hotkeys.
- [ ] `listenChannels` returns a slice aligned with `len(m.hotkeys)`.

## Audit reference

```json
{
  "id": "ISSUE-006",
  "severity": "low",
  "file": "internal/infrastructure/hotkey/manager.go",
  "line": 200,
  "message": "HotkeyManager.Start nests channel setup..."
}
```
