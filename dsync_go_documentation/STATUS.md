# Implementation Status (Informative)

This document describes what is implemented in the Go codebase right now.

Normative behavior is defined in:
- `dsync_go_documentation/spec/CLI.md`
- `dsync_go_documentation/spec/BEHAVIORS.md`
- `dsync_go_documentation/design/RSYNC.md`

## Implemented

- Go scaffold + CI + local dev workflow
  - `go.mod`, `Makefile`, `.golangci.yml`, `.github/workflows/ci.yml`
  - `make check` runs format check + lint + tests

- `dsync init`
  - Writes a starter config to `$XDG_CONFIG_HOME/dsync/config.toml`
  - Template source: `internal/config/template.go`

- Config discovery + parsing + validation
  - Discovery order: `$XDG_CONFIG_HOME/...` then `~/.config/...`
  - Code: `internal/config/config.go`, `internal/xdg/xdg.go`
  - Validates key invariants (e.g., endpoint paths absolute and not `/`, link endpoint types)

- Scope engine (link mode)
  - CLI scope -> CWD inference -> full-root
  - `--use-link-paths` overrides CWD inference
  - `--use-link-paths` conflicts with explicit scope and `--all`
  - `partial_only=true` forbids full-root
  - Code: `internal/scope/scope.go`

- Canonical rsync argv builder
  - Builds preview/apply argv; differ only by `--dry-run`
  - Global + per-link excludes supported
  - Remote transport uses `-e ssh` (ssh-config Host alias)
  - Code: `internal/rsync/sync.go`

- `dsync pull` / `dsync push` planning output
  - Prints resolved scope and rsync SRC/DEST
  - Prints the preview/apply argv
  - Does NOT execute rsync yet
  - Code: `internal/app/sync.go`

- CLI parsing for `pull`/`push`
  - Flags are allowed before/after positionals (interspersed)
  - Use `--` to end flag parsing
  - Rationale: the Go `flag` package stops parsing at first positional, which breaks UX and safety
  - Code: `internal/app/sync.go`

## Not implemented yet

- Running rsync for real
  - preview execution, summary parsing, prompting, and apply phase
- `dsync ls`
- `dsync clean` (local and `--remote`)
- `dsync doctor` rsync probes (currently only validates config load + prints counts)
- `--json` NDJSON output
- Integration tests that exercise rsync over ssh

## Quick smoke checks

```sh
make check

go run ./cmd/dsync init
go run ./cmd/dsync doctor

# Planning only (prints rsync argv)
go run ./cmd/dsync pull --link photos --dry-run
go run ./cmd/dsync push --link photos 2026/portraits --dry-run
```
