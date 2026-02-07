# dsync (Go) - Project Specification

This folder contains the authoritative, human- and agent-readable specification for a new project named `dsync`.

Legacy reference
- The original Rust `directory-sync` implementation is preserved under `rust/` for reference.

Design contract
- These docs are the source of truth for intended behavior.
- If behavior changes, update these docs in the same change set.
- If an implementation diverges, treat it as a bug unless/until the docs are updated.

Quick navigation
- Product overview + goals: `dsync_go_documentation/spec/OVERVIEW.md`
- Configuration format + semantics: `dsync_go_documentation/spec/CONFIG.md`
- CLI contract + examples: `dsync_go_documentation/spec/CLI.md`
- Behavioral rules (scope, guardrails, prompts, output): `dsync_go_documentation/spec/BEHAVIORS.md`
- Machine-readable output (NDJSON events): `dsync_go_documentation/spec/JSON.md`
- Canonical rsync command construction: `dsync_go_documentation/design/RSYNC.md`
- `dsync clean` (local + remote via rsync-only deletes): `dsync_go_documentation/design/CLEAN.md`
- Safety model and guardrails: `dsync_go_documentation/design/SAFETY.md`
- Agent execution plan / work breakdown: `dsync_go_documentation/TASKS.md`
- Minimal context for agents: `dsync_go_documentation/agent/CONTEXT.md`
- Current implementation status: `dsync_go_documentation/STATUS.md`

Core decisions (locked for MVP)
- Name: `dsync`
- Language: Go; ship a single binary
- Platforms: macOS + Linux required; Windows optional later
- Engine: `rsync` only; remote endpoints use rsync-over-ssh (via ssh-config Host aliases)
- No direct remote shell: dsync never executes standalone `ssh` commands (rsync may use ssh as its transport)
- Semantics: rsync "contents semantics" (trailing slash); `relative_path` scopes both source and destination
- Default: mirror (destination matches source within scope), includes deletions
- Safety: always run a preview (rsync dry-run) and print a clear summary before applying
- Full-root operations require `--all` (even with `--yes`)
- Non-interactive mode requires `--yes` (and `--all` when needed)
- Output: concise by default; `--verbose` streams full rsync output
- Agents: `--json` outputs NDJSON events on stdout; human text goes to stderr

Implementation note
- The Go rewrite lives at the repo root; the original Rust project is under `rust/`.
