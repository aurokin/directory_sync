# CLI Contract

This document defines the `dsync` CLI interface for MVP.

## Common concepts

- "endpoint mode": `dsync <cmd> <endpoint>` syncs between CWD and the endpoint.
- "link mode": `dsync <cmd> --link <link>` syncs between the link's endpoints.
- `relative_path` is always a scope that applies to BOTH source and destination roots.

All mutating commands follow two-phase behavior:
- Preview (dry-run) and summary output always run first.
- Apply runs after a prompt, unless `--yes` is provided.

## Commands

### `dsync init`

Creates `$XDG_CONFIG_HOME/dsync/config.toml` with a starter template.

Flags
- `--force`: overwrite existing config file

### `dsync doctor`

Validates config + checks prerequisites.

Expected checks (MVP)
- Local rsync exists and is sufficiently new.
- For ssh endpoints, attempt a harmless rsync operation to validate connectivity and that remote rsync is runnable.

Flags
- `--endpoint <name>`: only validate one endpoint
- `--link <name>`: only validate one link
- `--json`

### `dsync ls <name> [relative_path]`

Lists a directory using `rsync --list-only`.

Flags
- `--link`: treat `<name>` as link name
- `--json`
- `--verbose` (optional: stream full list-only output)

Link mode behavior
- `dsync ls --link photos` lists both local and remote roots (labeled).

### `dsync pull <name> [relative_path]`

Syncs from remote -> local.

Endpoint mode
- source: endpoint root (possibly ssh)
- destination: CWD

Link mode
- source: link.remote
- destination: link.local

Flags
- `--link`
- `--yes`: non-interactive apply; still runs preview+summary first
- `--dry-run`: preview only, no prompt, no apply
- `--all`: allow full-root operations when scope is empty
- `--dangerous`: override high-risk destination blocklist checks (see `dsync_go_documentation/design/SAFETY.md`)
- `--use-link-paths`: when link has `paths=[...]`, run the configured batch (ignored if a scope is provided)
- `--json`
- `--verbose`

### `dsync push <name> [relative_path]`

Syncs from local -> remote.

Endpoint mode
- source: CWD
- destination: endpoint root (possibly ssh)

Link mode
- source: link.local
- destination: link.remote

Flags: same as `pull`.

### `dsync clean <name> [relative_path]`

Deletes `.dsync-partial/` directories.

Default behavior
- cleans locally only

Flags
- `--link`
- `--remote`: also delete `.dsync-partial/` directories on the remote endpoint using rsync-only deletion
- `--yes`
- `--dry-run`
- `--use-link-paths`
- `--json`
- `--verbose`

Behavior is fully specified in `dsync_go_documentation/design/CLEAN.md`.

## Examples

Pull one subtree (link mode, inferred scope)
```sh
cd ~/photos/2026/portraits
dsync pull --link photos
```

Push that subtree back
```sh
dsync push --link photos
```

Preview full-root sync (requires `--all`)
```sh
dsync pull --link photos --all --dry-run
```

Run configured batch paths
```sh
dsync pull --link photos --use-link-paths
```
