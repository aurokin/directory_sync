# TASKS (Agent Work Breakdown)

This file is for agents and humans coordinating implementation work.

Source of truth
- Behavioral requirements: `dsync_go_documentation/spec/BEHAVIORS.md`
- CLI contract: `dsync_go_documentation/spec/CLI.md`
- Config: `dsync_go_documentation/spec/CONFIG.md`
- rsync flags: `dsync_go_documentation/design/RSYNC.md`
- Safety rules: `dsync_go_documentation/design/SAFETY.md`
- Clean semantics: `dsync_go_documentation/design/CLEAN.md`
- NDJSON events: `dsync_go_documentation/spec/JSON.md`

## Phase 1: Project scaffold

- [ ] Create new Go module for `dsync` (single binary)
- [ ] Add basic command runner that execs `rsync` with argv (no shell)
- [ ] Add `--help` output that documents guardrails and rsync-only contract

Acceptance criteria
- `dsync --help` exists and documents key flags and safety behavior.

## Phase 2: Config + resolution

- [ ] Implement config discovery: `$XDG_CONFIG_HOME/dsync/config.toml` then `~/.config/dsync/config.toml`
- [ ] Parse TOML into in-memory structs; validate schema
- [ ] Resolve endpoint roots and ensure trailing-slash contents semantics
- [ ] Validate constraints: exactly one remote side for a link; no endpoint root of `/`

References
- `dsync_go_documentation/spec/CONFIG.md`
- `dsync_go_documentation/spec/BEHAVIORS.md`

Acceptance criteria
- Invalid configs fail with actionable errors.

## Phase 3: Scope engine

- [ ] Implement scope resolution rules (CLI -> CWD inference -> link paths -> empty)
- [ ] Implement mismatch notice: scope overrides `paths` and print alternate commands
- [ ] Implement full-root guardrail: require `--all` for apply

References
- `dsync_go_documentation/spec/BEHAVIORS.md`

Acceptance criteria
- Link-mode CWD inference works for the photos workflow.
- Apply is refused without `--all` when scope is empty.

## Phase 4: rsync command builder

- [ ] Build canonical argv for preview/apply (mirror default)
- [ ] Implement global + per-link excludes
- [ ] Remote endpoints add `-e ssh` and host alias path form
- [ ] Default compress OFF

References
- `dsync_go_documentation/design/RSYNC.md`

Acceptance criteria
- Preview argv and apply argv differ only by `--dry-run`.

## Phase 5: Preview summary + prompts

- [ ] Run preview (`--dry-run --itemize-changes --stats`) and capture output
- [ ] Parse best-effort counts:
  - `would_delete` from `*deleting` lines
  - transfer counts/bytes from `--stats` where possible
- [ ] Print direction banner, resolved roots, and summary; print samples (first N)
- [ ] Implement interactive prompt and refusal in non-interactive mode unless `--yes`
- [ ] `--yes` still runs preview and prints summary, then applies without prompt
- [ ] `--dry-run` exits after preview (no prompt, no apply)

References
- `dsync_go_documentation/spec/BEHAVIORS.md`
- `dsync_go_documentation/design/SAFETY.md`

Acceptance criteria
- Output makes it clear what would be deleted.
- Non-interactive runs never block.

## Phase 6: Commands

- [ ] `dsync init` writes a starter config template
- [ ] `dsync ls` uses `rsync --list-only` (link mode lists both sides)
- [ ] `dsync pull` and `dsync push` work in endpoint and link modes
- [ ] `dsync doctor` validates prerequisites via harmless rsync probes

References
- `dsync_go_documentation/spec/CLI.md`

Acceptance criteria
- Can pull/push a scoped subtree with clear preview/apply behavior.

## Phase 7: NDJSON (`--json`)

- [ ] Implement `--json` mode:
  - NDJSON events to stdout
  - human logs to stderr
- [ ] Emit events: `resolve`, `preview_done`, `prompt`, `apply_start`, `apply_done`, `note`

References
- `dsync_go_documentation/spec/JSON.md`

Acceptance criteria
- `--json` output is valid NDJSON and not interleaved with human logs.

## Phase 8: clean

- [ ] Implement local clean (delete `.dsync-partial/` dirs and contents)
- [ ] Implement `--remote` clean using rsync-only deletion (filters + delete)
- [ ] clean follows preview-before-apply, `--yes`, `--dry-run`, and `--json`

References
- `dsync_go_documentation/design/CLEAN.md`

Acceptance criteria
- Clean previews show deletions and never deletes non-`.dsync-partial` paths.

## Phase 9: Tests + CI

- [ ] Unit tests for scope resolution and command building
- [ ] Integration tests (Linux CI): local-local rsync
- [ ] Integration tests (Linux CI): local-remote using sshd container + rsync

Acceptance criteria
- CI runs on PR and validates core invariants.

## Phase 10: Packaging

- [ ] Add GoReleaser config to build macOS + Linux binaries
- [ ] Document installation (copy binary into PATH)
