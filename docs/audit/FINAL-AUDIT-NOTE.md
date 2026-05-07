# Go 1.26 audit — backlog closeout

**Date:** 2026-05-07

## Pull requests merged (complete set #78–#86)

| PR | Title |
|----|--------|
| #78 | Go 1.26 upgrade and audit |
| #79 | fix: CI hardening, tests, Wire drift gate, lint fixes |
| #80 | docs(readme): clarify positioning, charter link, Go 1.26 audit |
| #81 | chore: audit lint, staticcheck, tests, and modules |
| #82 | fix: Wire API storage config and track wire_gen (ISSUE-001) |
| #83 | test: SQLite-backed repository CRUD (ISSUE-003) |
| #84 | fix: API TLS defaults and base URL validation (ISSUE-005) |
| #85 | fix: dedupe API DeleteNote 404 branch (ISSUE-004) |
| #86 | refactor: extract hotkey listenChannels helper (ISSUE-006) |

## Skipped

None. All listed PRs reached `main` via squash merge (stack: ISSUE branches → `chore/audit-deps-tests`, then #81 → `main`).

## Blockers

None remaining.

## Trivial CI fixes applied (unblocking quality)

- **#85 / #86:** `store_delete_test.go` used removed `InsecureSkipVerify` on `domainstorage.APIConfig`; updated to `TLSInsecureSkipVerify`. Regenerated `wire_gen.go` and checksum via `./scripts/regen-wire.sh` where needed.
- **#86:** Removed legacy `// +build` duplicate line in `manager_listen_channels_test.go` (satisfies `go fix -diff` gate).
- **#86:** After #85 merged, base branch update required; merged `chore/audit-deps-tests` into `refactor/hotkey-start` before final squash merge.

## Status

**all-merged**
