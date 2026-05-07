# PR: fix/tls-defaults (ISSUE-005)

## Summary

Make API TLS defaults explicit in `NewDefaultConfig`, reject userinfo in `storage.api.base_url` during `ValidateConfig`, and add unit tests for the default and validation rule.

## Changes

- `NewDefaultConfig` — populate `Storage` with explicit `API.InsecureSkipVerify: false`.
- `validateStorageAPI` — `url.Parse` on base URL; error if `User` is set.
- `internal/config/config_storage_api_test.go` — asserts default TLS verify flag and credential rejection.

## How to verify

```bash
go test ./internal/config/... -tags=wireinject -count=1 -v
go test ./... -tags=wireinject -count=1
```

## Acceptance criteria

- [ ] Default `InsecureSkipVerify` is `false` for new default config.
- [ ] `ValidateConfig` fails when API base URL contains embedded credentials.
- [ ] Full test suite passes (stacked on prior fix branches).

## Audit reference

```json
{
  "id": "ISSUE-005",
  "severity": "medium",
  "file": "internal/config/config.go",
  "line": 358,
  "message": "NewDefaultConfig omits an explicit Storage.API block..."
}
```
