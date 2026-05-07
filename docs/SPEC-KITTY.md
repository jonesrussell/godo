# Spec Kitty governance

Spec Kitty provides **spec-driven governance** for this repository: a human-edited charter (`.kittify/charter/charter.md`) is the source of truth, and extracted YAML under `.kittify/charter/` keeps machine-readable policy aligned with that charter. CI verifies the install and detects charter drift so merges do not silently desync governance from the document maintainers edit.

## Local commands

From the repository root (with [Task](https://taskfile.dev/) installed):

- **`task spec:verify`** — runs `spec-kitty verify-setup` (see `.kittify/config.yaml` `verify_setup_cmds`).
- **`task spec:sync`** — runs `spec-kitty charter sync` after you change `charter.md`; commit the updated extracted files when the charter meaningfully changes.

Install the CLI if missing (recommended):

```bash
pip install spec-kitty-cli
```

## Policy

**Do not wire external harnesses into `main` without a mission and charter amendment.** Remote integrations, live third-party services, and new trust boundaries belong in Spec Kitty mission artifacts and must be reflected in the charter before production code depends on them.

## CI

The workflow `.github/workflows/spec-kitty.yml` runs on pushes and pull requests to `main`. It installs `spec-kitty-cli` via pip, runs `scripts/check-charter.sh`, `spec-kitty verify-setup`, and a read-only drift check using `spec-kitty charter status --json` (`status` must be `synced`, and `current_hash` must equal `stored_hash`).
