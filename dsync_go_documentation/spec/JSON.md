# NDJSON Event Output (`--json`)

This document is normative.

## High-level rules

- When `--json` is enabled, dsync emits newline-delimited JSON objects (NDJSON) on stdout.
- Human-readable logs go to stderr.
- Each JSON line is one event with a stable `type`.
- Fields MUST be additive-only in minor versions. Removing/renaming fields is a breaking change.

## Common fields

All events MUST include:
- `type` (string)
- `ts` (string; RFC3339 timestamp)
- `cmd` (string; top-level command, e.g. `pull`, `push`, `clean`, `ls`, `doctor`)

Many events SHOULD include:
- `mode` (string; `endpoint` or `link`)
- `name` (string; endpoint or link name)
- `scope` (string; resolved scope, or empty)
- `scope_source` (string; `cli`, `cwd`, `link_paths`, `empty`)
- `source` (string; resolved rsync SRC, with trailing `/` when applicable)
- `dest` (string; resolved rsync DEST, with trailing `/` when applicable)
- `argv` (array of string; rsync arguments excluding the program name)

## Event types

### `resolve`

Emitted after config and arguments are resolved.

Additional fields
- `mirror` (bool)
- `deletes_enabled` (bool)
- `excludes_count` (number)
- `batch` (object)
  - `enabled` (bool)
  - `paths` (array of string)

### `preview_done`

Emitted after the dry-run completes.

Additional fields
- `exit_code` (number)
- `would_delete` (number; best-effort)
- `would_transfer_files` (number; best-effort)
- `would_transfer_bytes` (number; best-effort)

### `prompt`

Emitted immediately before prompting in interactive mode.

Additional fields
- `message` (string)

### `apply_start`

Emitted immediately before running the apply command.

Additional fields
- `yes` (bool)

### `apply_done`

Emitted after the apply completes.

Additional fields
- `exit_code` (number)

### `note`

Used for mismatch notices and other structured warnings.

Additional fields
- `level` (string; `info`, `warn`, `error`)
- `message` (string)
- `hints` (array of string; optional)

## Example output

```json
{"type":"resolve","ts":"2026-02-07T12:00:00Z","cmd":"pull","mode":"link","name":"photos","scope":"2026/portraits","scope_source":"cwd","source":"photo-box:/srv/photos/2026/portraits/","dest":"/Users/you/photos/2026/portraits/","mirror":true,"deletes_enabled":true,"excludes_count":4,"argv":["-a","--no-owner","--no-group","--mkpath","--protect-args","--partial","--partial-dir=.dsync-partial","--human-readable","--stats","--itemize-changes","--delete","--delete-delay","-e","ssh","--dry-run","photo-box:/srv/photos/2026/portraits/","/Users/you/photos/2026/portraits/"],"batch":{"enabled":false,"paths":[]}}
{"type":"preview_done","ts":"2026-02-07T12:00:01Z","cmd":"pull","exit_code":0,"would_delete":3,"would_transfer_files":12,"would_transfer_bytes":10485760}
{"type":"prompt","ts":"2026-02-07T12:00:01Z","cmd":"pull","message":"Type y to apply"}
{"type":"apply_start","ts":"2026-02-07T12:00:02Z","cmd":"pull","yes":false}
{"type":"apply_done","ts":"2026-02-07T12:00:10Z","cmd":"pull","exit_code":0}
```

Notes
- The argv array must contain real arguments. Do not include placeholder control characters (the example uses `...` only as shorthand).
- dsync should avoid emitting sensitive data (e.g., if later adding credentials for other backends).

argv convention
- dsync runs `rsync` from PATH.
- `argv` contains the arguments passed to `rsync` (i.e., excludes the program name).
