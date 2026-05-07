# PR: fix/tests-repo-integration (ISSUE-003)

## Summary

Add `internal/domain/testhelpers` SQLite harness (modernc.org/sqlite, no CGO) and `repository_sqlite_test.go` with deterministic CRUD coverage for `NoteRepository` backed by a temp database.

## Changes

- `internal/domain/testhelpers/sqltest.go` — `NewTempSQLiteUnified` opens SQLite under `t.TempDir()` and returns `storage.UnifiedNoteStorage` + cleanup.
- `internal/domain/repository/repository_sqlite_test.go` — integration-style CRUD test using the helper.

## How to verify

```bash
go test ./internal/domain/... -tags=wireinject -count=1 -v
go test ./... -tags=wireinject -count=1
```

## Acceptance criteria

- [ ] `TestNoteRepository_SQLite_CRUD` passes on Linux without CGO.
- [ ] Tests use only `t.TempDir()` paths (no shared global DB).
- [ ] `go test ./... -tags=wireinject` passes when stacked on the Wire fix branch.

## Audit reference

```json
{
  "id": "ISSUE-003",
  "severity": "medium",
  "file": "internal/domain/repository",
  "line": 0,
  "message": "Repository layer has no SQLite-backed integration tests..."
}
```
