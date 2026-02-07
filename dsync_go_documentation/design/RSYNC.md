# Canonical rsync Command Construction

This document is normative. It defines how dsync builds rsync commands.

## Principles

- dsync invokes `rsync` via exec (argv array), never via a shell.
- All remote operations are performed by rsync itself (rsync-over-ssh).
- Remote SSH configuration is delegated to `~/.ssh/config`.

## Remote path forms

Local
- `/abs/path/.../`

SSH
- `<ssh-host-alias>:/abs/path/.../`

Rsync MUST be invoked with `-e ssh` when either endpoint is type `ssh`.

## Trailing slash (contents semantics)

- dsync MUST treat endpoint roots and scoped paths as directories whose contents are synced.
- Effective SRC and DEST MUST end with `/`.

Examples
- Root: `/srv/photos` => `/srv/photos/`
- Root+scope: `/srv/photos/` + `2026/portraits` => `/srv/photos/2026/portraits/`

## Flag sets

### Base flags (always)

These are included in both preview and apply.

- `-a`
- `--no-owner`
- `--no-group`
- `--mkpath`
- `--protect-args`
- `--partial`
- `--partial-dir=.dsync-partial`
- `--human-readable`
- `--stats`
- `--itemize-changes`

Remote-only addition
- `-e ssh` (only when either side is `ssh`)

### Mirror flags (default)

When mirror is enabled (default), dsync adds:

- `--delete`
- `--delete-delay`

Note
- dsync intentionally does NOT use `--delete-excluded` (rsync defaults for exclude behavior).

### Preview flags

Preview phase adds:

- `--dry-run`

Note
- dsync includes `--itemize-changes` in both preview and apply so that logging and parsing are consistent.

### Apply phase

Apply phase uses the same flags as preview minus `--dry-run`.

## Excludes

- Global excludes: `--exclude <pattern>` for each pattern in `[global].excludes`.
- Per-link excludes: appended after global excludes.

Rationale
- Rsync evaluates rules in order; we want per-link excludes to override global rules where relevant.

Implementation note
- `--partial-dir=.dsync-partial` auto-adds a protective exclude for `.dsync-partial/`.
- `dsync clean` must NOT rely on these excludes and should use specialized filter rules instead.

## Compression

- Default compression is OFF.
- Rationale: common LAN usage and media workflows (photos/videos) generally do not compress well.
- Future: add per-link `compress=true` and/or CLI `--compress`.

## Parsing preview output

dsync produces the human summary by parsing rsync output, best-effort.

Primary data sources
- `--itemize-changes` lines
  - deletions: lines starting with `*deleting `
- `--stats` block
  - use it to derive "would transfer" counts/bytes when possible

Best-effort rule
- If parsing fails, dsync MUST still show the raw preview output when `--verbose` is set.
- In default mode, dsync SHOULD print a clear note that parsing was incomplete.
