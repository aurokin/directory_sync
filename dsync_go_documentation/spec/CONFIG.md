# Configuration

Config file path
- Primary: `$XDG_CONFIG_HOME/dsync/config.toml`
- Fallback: `~/.config/dsync/config.toml`

`dsync init` will generate a starter config at that location.

## TOML schema (MVP)

Top-level sections
- `[global]`
- `[endpoints.<name>]`
- `[links.<name>]`

### `[global]`

Fields
- `excludes` (array of string patterns)

Notes
- Global excludes apply to push/pull operations.
- `dsync clean` ignores excludes (clean targets `.dsync-partial` explicitly).

Recommended defaults for `global.excludes`
- `.DS_Store`
- `.git/`
- `node_modules/`
- `.dsync-partial/`

### `[endpoints.<name>]`

Fields
- `type` (string; `"local"` or `"ssh"`)
- `path` (string; absolute path on that machine)
- `host` (string; required when `type="ssh"`; SSH config host alias)

Semantics
- `path` is a directory root. dsync uses rsync "contents semantics": the effective rsync argument ends with `/`.
  - Example: `path = "/srv/photos"` means the rsync root is `/srv/photos/`.

### `[links.<name>]`

Fields
- `local` (string; endpoint name)
- `remote` (string; endpoint name)
- `mirror` (bool; default: `true`)
- `partial_only` (bool; default: `false`)
- `paths` (array of strings; default: `[]`)
- `excludes` (array of string patterns; default: `[]`)

Constraints
- 1:1 only: one local endpoint and one remote endpoint.
- Exactly one side must be `type="ssh"` in MVP.

Semantics
- `paths` provides a default batch of scopes when no scope is provided by the user.
- If a scope is provided (CLI `relative_path` or inferred from CWD), it overrides `paths` for that run.
- When `partial_only=true`, full-root link operations require extra explicit intent (see `dsync_go_documentation/spec/BEHAVIORS.md`).

## Example config

```toml
[global]
excludes = [
  ".DS_Store",
  ".git/",
  "node_modules/",
  ".dsync-partial/",
]

[endpoints.laptop_photos]
type = "local"
path = "/Users/you/photos"

[endpoints.server_photos]
type = "ssh"
host = "photo-box"
path = "/srv/photos"

[links.photos]
local = "laptop_photos"
remote = "server_photos"
mirror = true
partial_only = false
paths = []
excludes = ["*.tmp"]
```
