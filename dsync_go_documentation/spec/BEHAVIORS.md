# Behavioral Specification

This document is normative. Use it as the acceptance criteria for implementation.

Keywords
- MUST, SHOULD, MAY are used in the RFC sense.

## Rsync-only contract

- dsync MUST use `rsync` as the engine for all sync operations.
- For remote endpoints, dsync MUST use rsync-over-ssh (via `rsync -e ssh ...`).
- dsync MUST NOT execute standalone remote shell commands (e.g., `ssh host <cmd>`).

## Paths and scope

### Contents semantics

- Endpoint `path` is treated as a directory root whose *contents* are synced.
- dsync MUST construct rsync source/destination roots with a trailing slash.

Example
- Endpoint path `/srv/photos` becomes rsync root `/srv/photos/`.

### Scope applies to both ends

- When a `relative_path` scope is used, dsync MUST append it to both source and destination roots.
- The effective rsync roots MUST end with `/` after scoping.

Example
- root: `/srv/photos/`
- scope: `2026/portraits`
- effective: `/srv/photos/2026/portraits/`

## Scope resolution rules

Given an operation in link mode, dsync resolves a scope in this order:

1) Explicit CLI `relative_path` argument (if provided)
2) CWD inference (if CWD is inside the link's local endpoint root)
3) Link `paths=[...]` batch (only when no scope in steps 1-2 and `--use-link-paths` is not required)
4) Empty scope (full-root)

Rules
- If a scope is resolved via steps 1 or 2, it MUST override link `paths` for that run.
- When a scope overrides `paths`, dsync MUST print a mismatch notice and MUST print alternate commands:
  - `--use-link-paths` (to run configured `paths`)
  - `--all` (to run full-root)

## Full-root guardrail (`--all`)

- If the effective scope is empty (full-root), dsync MUST refuse to apply unless `--all` is provided.
- This applies to both interactive and `--yes` execution.
- `--dry-run` MAY be allowed without `--all` (implementation choice), but the default is:
  - allow preview without `--all`
  - refuse apply without `--all`

Rationale
- Full-root mirroring is often dangerous when run from the wrong directory or against the wrong endpoint.

## Non-interactive execution guardrail

- If stdin is not a TTY (non-interactive), dsync MUST refuse to run mutating operations unless `--yes` is provided.
- If a full-root apply would occur, `--yes` MUST be paired with `--all`.
- In non-interactive mode, dsync MUST NOT block waiting for input.

## Preview-before-apply

- Mutating operations (`push`, `pull`, `clean`) MUST run a preview phase first.
- Preview phase uses rsync `--dry-run` and prints a summary.
- Apply phase runs the same rsync invocation without `--dry-run`.

`--yes` behavior
- When `--yes` is provided, dsync MUST still run the preview phase and print its summary.
- Then dsync MUST immediately run the apply phase without prompting.

`--dry-run` behavior
- When `--dry-run` is provided, dsync MUST run only the preview phase and MUST NOT prompt.

## Mirror semantics (default)

- Links and endpoints default to "mirror" behavior.
- Mirror means destination is made to match source within the selected scope.
- dsync MUST use rsync deletion flags consistent with safer deletes:
  - MUST use `--delete` and SHOULD use `--delete-delay`.

Important
- Mirror MUST NOT be described as "copy everything".
- Mirror still transfers only differences; unchanged files are skipped per rsync rules.

## Output requirements (human)

For `push`/`pull`, dsync MUST print:
- Direction label: `PUSH` or `PULL`
- Source -> destination roots (fully resolved, with trailing `/`)
- Scope source (CLI, inferred from CWD, link paths, or empty)
- Mirror status (`mirror=true`) and whether deletions are enabled
- Exclude summary (count of global + link excludes)

Preview summary MUST include:
- "Would transfer" count and bytes (best-effort; derived from rsync output)
- "Would delete" count (best-effort)
- A sample of itemized changes (first N)
- A sample of deletions (first N)

In interactive mode, dsync MUST prompt "type y" to proceed.

## Output requirements (machine / agents)

`--json` flag
- When `--json` is set, dsync MUST emit NDJSON events on stdout.
- Human-oriented text MUST go to stderr.
- Event schema is defined in `dsync_go_documentation/spec/JSON.md`.

`--json` + `--verbose`
- When both are set, raw rsync output MUST go to stderr (so stdout remains valid NDJSON).

## `--verbose`

- Default output is summary-level.
- With `--verbose`, dsync MUST stream the full underlying rsync output.

## Additional safety checks

dsync MUST enforce additional safety rules defined in:
- `dsync_go_documentation/design/SAFETY.md`
