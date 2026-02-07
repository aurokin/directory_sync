# TASKS (Implementation Plan)

This file is the actionable plan for implementing `dsync`.

Read first
- `dsync_go_documentation/agent/CONTEXT.md`
- `dsync_go_documentation/STATUS.md` (what is already implemented)

Normative behavior
- `dsync_go_documentation/spec/CLI.md`
- `dsync_go_documentation/spec/BEHAVIORS.md`
- `dsync_go_documentation/design/RSYNC.md`
- `dsync_go_documentation/design/SAFETY.md`

Guiding rule
- If behavior changes, update the spec docs in the same change set.

## Implemented (do not redo)

- Go scaffold + CI + dev workflow
  - `go.mod`, `Makefile`, `.golangci.yml`, `.github/workflows/ci.yml`
- Config discovery + parse + validation
  - `internal/config/config.go`, `internal/xdg/xdg.go`
- Scope engine (including `--use-link-paths` override and conflicts)
  - `internal/scope/scope.go`
- Canonical rsync argv builder (preview/apply)
  - `internal/rsync/sync.go`
- pull/push planning output (prints SRC/DEST and argv)
  - `internal/app/sync.go`

## Next: Phase 5 (Run rsync + preview/apply UX)

Goal
- Make `pull`/`push` execute rsync safely and predictably.

### 5.1 Run rsync (exec, streaming, and capture)

- Add a small runner that executes `rsync` with argv from `internal/rsync/sync.go`.
- Support 2 output modes:
  - default: capture output for parsing, print a concise summary
  - `--verbose`: stream raw rsync output to stderr

Acceptance criteria
- Exit code reflects rsync exit code.
- No shell invocation.

### 5.2 Preview summary parsing

- From preview output (`--dry-run`):
  - `would_delete`: count `*deleting ` lines
  - best-effort transfer counts/bytes: parse `--stats`
- Print a clear summary:
  - direction banner
  - SRC and DEST
  - mirror/deletes enabled
  - exclude count (global + link)
  - would delete / would transfer
  - sample deletions + sample itemize lines (first N)

Acceptance criteria
- Output makes it hard to miss deletions.

### 5.3 Prompting and non-interactive guardrails

- Interactive mode:
  - preview always runs first
  - prompt "type y" to apply
- Non-interactive mode:
  - refuse apply unless `--yes` is set
  - never block waiting for input
- `--yes`:
  - still runs preview and prints summary
  - then applies without prompting
- `--dry-run`:
  - preview only, no prompt, no apply
- Full-root apply:
  - require `--all`

Acceptance criteria
- `--yes` works in scripts without hanging.
- Full-root apply is refused without `--all`.

## Phase 6 (Finish core commands)

### 6.1 `doctor` probes

- Keep config validation.
- Add local rsync version check.
- Add a harmless remote probe for ssh endpoints (rsync-over-ssh only).

### 6.2 `ls`

- Implement `rsync --list-only`.
- Link mode lists both local and remote (labeled).

### 6.3 Safety blocklist + `--dangerous`

- Enforce high-risk destination blocklist from `dsync_go_documentation/design/SAFETY.md`.
- Implement `--dangerous` override (should require `--yes`).

## Phase 7 (`--json` NDJSON)

- Emit NDJSON events per `dsync_go_documentation/spec/JSON.md`.
- Ensure stdout is NDJSON-only; all human logs and raw rsync output go to stderr.

## Phase 8 (`clean`)

- Local clean: delete `.dsync-partial/` dirs and contents under the resolved local roots.
- Remote clean: implement rsync-only deletion per `dsync_go_documentation/design/CLEAN.md`.
- Preview-before-apply and all guardrails apply.

## Phase 9 (Integration tests)

- Add integration tests that run real rsync locally.
- Add ssh-based integration tests using an sshd container that has rsync.
- Test dangerous behaviors explicitly:
  - mirror deletes within scope
  - excludes protect files
  - `--all` required for full-root apply
  - `partial_only` forbids full-root

## Phase 10 (Packaging)

- Add GoReleaser to build macOS + Linux binaries.
- Document install: copy `dsync` into a PATH directory.
