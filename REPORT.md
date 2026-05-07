# PR: fix/duplicate-api-handler (ISSUE-004)

## Summary

Remove the duplicated `http.StatusNotFound` branch in `Store.DeleteNote` and add an `httptest` unit test that asserts `NotFoundError` is returned for a 404 API response.

## Changes

- `internal/infrastructure/storage/api/store.go` — single NotFound handling path in `DeleteNote`.
- `internal/infrastructure/storage/api/store_delete_test.go` — `TestStore_DeleteNote_NotFound`.

## How to verify

```bash
go test ./internal/infrastructure/storage/api/... -tags=wireinject -count=1 -v
go test ./... -tags=wireinject -count=1
```

## Acceptance criteria

- [ ] `DeleteNote` has no duplicate identical `StatusNotFound` checks.
- [ ] New test passes against a local `httptest` server.

## Audit reference

```json
{
  "id": "ISSUE-004",
  "severity": "medium",
  "file": "internal/infrastructure/storage/api/store.go",
  "line": 229,
  "message": "DeleteNote duplicates identical http.StatusNotFound branches..."
}
```
