# Project Charter — Godo

Human-edited governance for **Godo**: a local-first desktop companion that will serve as the **system-level interface** to one or more AI harnesses. The historical “quick notes / todos” surface is being **re-centered on durable, queryable state** and an **audit trail**, not on ephemeral UI-only lists.

**Harness work for now:** Treat any remote AI harness as **pluggable behind a small interface**. Use **in-memory fakes, `httptest` servers, and golden/fixture payloads** in tests and local dev until a real integration is explicitly chartered. Production coupling to a specific host (for example Claudriel) is **out of scope** until you add it back in a mission and bump this charter.

Last amended: 2026-05-07

---

## Identity and north star

1. **Local authority:** The desktop client owns capture UX (hotkeys, windows), local policy defaults, and persisted artifacts the user should retain without network access.
2. **Harness boundary:** A harness (local or remote) provides model/session orchestration **only through a documented client boundary** in this codebase—implemented with **swappable adapters**, not ad hoc HTTP scattered across UI.
3. **Truth in SQLite:** User-visible state, conversation transcripts (as the product evolves), and **audit metadata** live in **SQLite** as the **canonical source of truth** for the local agent. Other caches or in-memory views are projections, not competing primaries.
4. **Spec-driven change:** Material pivots (storage schema, harness boundary, trust boundaries) go through Spec Kitty **specify → plan → tasks** and are reflected here when governance changes.

---

## Testing standards

- Target **80%+ test coverage** on domain and application packages; infrastructure may be lower but must be justified in review.
- Use **`go test ./...`** with **`-tags=wireinject`** wherever Wire-generated code or DI graphs participate in tests (see project `CLAUDE.md`).
- Prefer table-driven tests for branching logic; use `t.TempDir()` for filesystem-backed tests.
- **Harness-facing code** is tested against **fixtures and test doubles** by default; contract-style tests may spin a **local `httptest.Server`** instead of calling real external services in CI or default `go test`.

---

## Quality and linting gates

- **`task fmt`** then **`task lint`** (or `task check` for a fast gate) before merge unless a mission documents an exception.
- Primary linter: **golangci-lint** (project toolchain). Fix or suppress with narrow scope and comment.
- **1 approval** required before merging to `main`.

---

## Commits and change hygiene

- **Conventional commits** for changelog and release hygiene unless a one-off hotfix branch is explicitly agreed.
- Cross-cutting changes (storage schema, public API surface, harness adapter interfaces) should be **one logical commit or a small series** with a clear message scope.

---

## Performance targets

- Routine **local** operations (open window, list/query local store, append audit row) should complete in **under 2 seconds** on a typical dev machine, excluding optional network calls to a harness.
- Any future **live harness I/O** must be **off the UI thread** (async, cancellation, bounded timeouts); until wired, this is satisfied by not blocking Fyne on network calls in new code paths.

---

## Branch strategy

1. Branch `main` is the integration branch; default to short-lived feature branches.
2. When Spec Kitty lanes are active, follow lane merge rules from mission docs; do not merge drifted charter bundles.
3. External harness wire-ups: when a real remote is introduced, document request and event assumptions in a mission spec (and extend this charter if governance changes); until then, fixtures only in tree.

---

## Governance activation

```yaml
mission: software-dev
selected_paradigms: []
selected_directives: []
available_tools: [git, spec-kitty]
template_set: software-dev-default
```

---

## Data store: SQLite as source of truth

| Concern | Decision |
| --- | --- |
| Primary persistence | **SQLite** (single canonical DB path per install; migrations versioned in-repo). |
| Todos / notes | Legacy “todo list” semantics migrate toward **first-class entities** (e.g. messages, threads, tasks) stored in normalized tables—not as the only copy in widget state. |
| Audit | Maintain an **append-only audit log** (or append-only event table) for operations that matter for **trust, replay, or future cross-system reconciliation**. Redaction/deletion is a **policy layer**, not silent row loss without trace. |
| API / other backends | Optional HTTP storage remains **secondary** unless explicitly configured; if enabled, **sync rules and conflict policy** must be specified in a mission spec, not ad hoc. |

---

## Deferred: live harness integration

A future mission may attach a concrete remote (for example **Claudriel**). Until then:

- **No required dependency** on any specific harness URL or auth flow in application code paths used for default builds.
- **Fixtures and fakes** are the expected way to exercise parsing, persistence, and UI flows that will eventually talk to a real API.
- When integration lands, add a short **wire + trust** subsection here (or link to an OpenAPI artifact) and add **contract tests** that still run fully offline using recorded responses where appropriate.

---

## Project directives

1. Single primary store: do not introduce a second primary mutable store for the same logical entities without a chartered migration and dual-write/dual-read plan.
2. Audit accountability: mutations that affect user trust or future harness reconciliation must be observable (append-only log or equivalent immutable history).
3. Harness adapter and fixtures: outbound “harness” behavior lives behind an explicit interface; tests use fakes or `httptest`—not live third-party services—unless a mission opts into a gated integration suite.
4. Architecture respect: preserve domain, application, and infrastructure separation unless a mission explicitly redesigns boundaries.
5. Workflow: non-trivial work uses Spec Kitty artifacts so implementation stays reviewable and reversible.

---

## Amendment process

Amendments ship via PR with **at least one reviewer** and an explicit note if **governance.yaml** / **directives.yaml** outputs change behavior for agents.

## Exception policy

Time-bounded exceptions are allowed for experiments; document **rationale**, **expiry**, and **follow-up issue** in the PR or mission dossier.

---

## Reference index

| Reference ID | Kind | Summary | Local Doc |
| --- | --- | --- | --- |
| `USER:PROJECT_PROFILE` | user_profile | Interview-derived defaults when present. | `_LIBRARY/user-project-profile.md` |
| `TEMPLATE_SET:software-dev-default` | template_set | Default software mission doctrine bundle. | `_LIBRARY/template-set-software-dev-default.md` |
