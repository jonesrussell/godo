# Repository audit artifacts

This folder holds a snapshot of lint, static analysis, tests, module justification, and size checks. Regenerate locally from the repository root (`github.com/jonesrussell/godo`).

## Prerequisites

- Go toolchain matching `go.mod` (see `go` directive).
- `golangci-lint` on PATH (project uses v2-style flags).
- `jq` for JSON shaping.
- `staticcheck`: `go install honnef.co/go/tools/cmd/staticcheck@v0.6.1` (then use `$GOPATH/bin/staticcheck` or `/tmp/staticcheck` as in the captured `summary.json` notes). Using `go run ... -- -f=json ./...` is fragile; installing the binary avoids `go run` flag parsing issues.

## Commands

```bash
mkdir -p docs/audit

# Lint (golangci-lint v2)
golangci-lint run ./... --output.json.path=docs/audit/lint.json 2> docs/audit/lint.stderr || true
echo "lint_exit:$?" >> docs/audit/lint.stderr

# Staticcheck (JSON)
staticcheck -f=json ./... > docs/audit/staticcheck.json 2>> docs/audit/staticcheck.stderr || true
echo "staticcheck_exit:$?" >> docs/audit/staticcheck.stderr

# Tests + coverage
go test ./... -coverprofile=docs/audit/cover.out -tags=wireinject -count=1 > docs/audit/test-output.txt 2>&1 || true
echo "test_exit:$?" >> docs/audit/test-output.txt
go tool cover -func=docs/audit/cover.out | grep '^total:' > docs/audit/coverage-total.txt

# If committing: this repo ignores `*.out`; force-add the profile when needed.
# git add -f docs/audit/cover.out

# Dependency package paths (note: use ImportPath, not Path)
go list -deps -json ./... | jq -rs 'map(.ImportPath) | unique | .[]' -r > docs/audit/deps.txt 2> docs/audit/go-list-deps.stderr

# go mod why for every module required in go.mod
go mod edit -json | jq -r '.Require[] | .Path' | sort -u > docs/audit/go-mod-modules.txt
: > docs/audit/mod-why.txt
while IFS= read -r m; do
  echo "==== go mod why -m $m ====" >> docs/audit/mod-why.txt
  go mod why -m "$m" >> docs/audit/mod-why.txt 2>&1 || true
  echo >> docs/audit/mod-why.txt
done < docs/audit/go-mod-modules.txt

# CGO: direct `import "C"` in the repo
grep -R --line-number '^import .*"C"' . > docs/audit/cgo-packages.txt 2>&1 || true

# go vet
go vet ./... > docs/audit/govet.txt 2>&1 || true
echo "govet_exit:$?" >> docs/audit/govet.txt

# Large files: requested pipeline (first 200 lines of each file, then first 200 of that stream)
git ls-files -z | xargs -0 sed -n '1,200p' | sed -n '1,200p' > docs/audit/large-files-user-pipeline.txt

# Large files: meaningful >2000 LOC scan
git ls-files | while read -r f; do
  [ -f "$f" ] || continue
  n=$(wc -l <"$f" 2>/dev/null) || continue
  [[ "$n" =~ ^[0-9]+$ ]] || continue
  (( n > 2000 )) && echo "$n $f"
done | sort -rn > docs/audit/large-files.txt
```

Recompute `summary.json` by re-running the audit steps, then merge metrics and raw file contents as needed (the checked-in `summary.json` was produced with a short Python one-liner in the audit session).

## Reading `summary.json`

- `lint_count` — issues in `lint.json` (`Issues` array).
- `staticcheck_count` — diagnostics lines/objects in `staticcheck.json` (compile failures surface as a single JSON object with `"code":"compile"`).
- `test_coverage_percent` — parsed from `coverage-total.txt` (`go tool cover -func`).
- `cgo_packages` — matches from the `import "C"` grep (transitive CGO via dependencies is not listed here).
- `unused_modules` — modules whose `go mod why -m` output contains `(main module does not need module ...)` (empty in the current snapshot because every listed module is reachable from the main module or tools).
- `large_files` — structured entries from the `wc -l` scan; see `large-files-user-pipeline.txt` for the literal sed pipeline output from the checklist.

## Snapshot caveats (this branch)

`internal/application/container/wire_gen.go` was out of sync with config/domain types at audit time, which produced typecheck noise in `lint.json`, compile output in `staticcheck.json`, and failures in `govet.txt`. This audit does not regenerate or commit Wire output.
