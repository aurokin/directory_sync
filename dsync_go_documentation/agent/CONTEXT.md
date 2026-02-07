# Agent Context (dsync MVP)

If you are an agent implementing dsync, read these files first:

1) `dsync_go_documentation/spec/OVERVIEW.md`
2) `dsync_go_documentation/spec/CONFIG.md`
3) `dsync_go_documentation/spec/CLI.md`
4) `dsync_go_documentation/spec/BEHAVIORS.md`
5) `dsync_go_documentation/spec/JSON.md`
6) `dsync_go_documentation/design/RSYNC.md`
7) `dsync_go_documentation/design/SAFETY.md`
8) `dsync_go_documentation/design/CLEAN.md`

Implementation status
- `dsync_go_documentation/STATUS.md`

Implementation invariants (do not violate)
- Rsync-only engine; no standalone `ssh` remote commands.
- Always preview before apply, including when `--yes` is used.
- Full-root apply requires `--all`.
- Non-interactive runs require `--yes` (and `--all` when scope is empty).
- Link-mode scope inference from CWD is a core workflow.
- `--json` emits NDJSON on stdout; human text goes to stderr.

Where to put work breakdown
- See `dsync_go_documentation/TASKS.md`.

Legacy reference
- The old Rust code is in `rust/`.

Code map (current)
- CLI entrypoint: `cmd/dsync/main.go`
- Command routing: `internal/app/app.go`
- pull/push planning + flag parsing: `internal/app/sync.go`
- config discovery/parse/validate: `internal/config/config.go`
- scope resolution: `internal/scope/scope.go`
- rsync argv builder: `internal/rsync/sync.go`
- XDG config paths: `internal/xdg/xdg.go`
