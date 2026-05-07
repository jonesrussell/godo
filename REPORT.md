# PR: fix/wire-api-config (ISSUE-001)

## Summary

Align `ProvideUnifiedStorage` generated Wire output with `config.APIConfig` / `domain/storage.APIConfig` by using `InsecureSkipVerify` only, add `scripts/regen-wire.sh`, and track `wire_gen.go` via a `.gitignore` exception so fresh clones compile.

## Changes

- `internal/application/container/wire_gen.go` — regenerated with `go run github.com/google/wire/cmd/wire@v0.7.0` from `internal/application/container` (matches `wire.go`).
- `scripts/regen-wire.sh` — runs Wire v0.7.0 for the container package.
- `.gitignore` — exception `!internal/application/container/wire_gen.go` so DI output is versioned (previously matched `*_gen.go` and was absent from git).

## How to verify

```bash
./scripts/regen-wire.sh
go test ./... -tags=wireinject -count=1
```

## Acceptance criteria

- [ ] `go test ./... -tags=wireinject` passes on Linux without a pre-existing untracked `wire_gen.go`.
- [ ] `internal/application/container/wire_gen.go` contains `InsecureSkipVerify: cfg.Storage.API.InsecureSkipVerify` in the API config literal (no `TLSInsecureSkipVerify`).
- [ ] `./scripts/regen-wire.sh` completes without error and leaves a clean `git diff` for `wire_gen.go` when providers are unchanged.

## Audit reference

```json
{
  "id": "ISSUE-001",
  "severity": "critical",
  "file": "internal/application/container/wire_gen.go",
  "line": 183,
  "message": "Generated Wire code references TLSInsecureSkipVerify..."
}
```
