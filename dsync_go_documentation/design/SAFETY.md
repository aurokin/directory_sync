# Safety Model

This document is normative.

dsync is mirror-by-default. Mirror implies deletion. The safety model exists to prevent catastrophic mistakes.

## Layer 1: Always preview

- All mutating commands MUST run a preview rsync invocation first.
- Preview MUST include `--dry-run --itemize-changes --stats`.
- Preview output MUST make direction and deletions obvious.

## Layer 2: Explicit full-root opt-in (`--all`)

- If the effective scope is empty, dsync MUST require `--all` to apply.
- This applies even in interactive mode.
- `--yes` does not override this; `--yes` must be paired with `--all`.

Rationale
- "Wrong directory" mistakes are common in endpoint mode.
- Full-root mirroring is too destructive to allow without explicit intent.

## Layer 3: Non-interactive guardrail

- If stdin is not a TTY, dsync MUST refuse mutating commands unless `--yes` is set.
- In non-interactive mode, dsync MUST never prompt or block.

## Layer 4: High-risk destination blocklist

dsync MUST refuse to apply if the destination root resolves to a high-risk path, unless overridden.

High-risk local destinations (MVP)
- `/`
- The current user's home directory (exact path match)

High-risk configured endpoint roots
- Any configured endpoint `path` equal to `/` MUST be rejected.

Override
- dsync SHOULD provide `--dangerous` to override this blocklist.
- `--dangerous` SHOULD require `--yes` (to avoid a casual interactive override).

Notes
- Remote home is not reliably discoverable under the no-remote-shell constraint.
- Therefore, the remote blocklist is intentionally limited to the explicit configured root.

## Mismatch notices

If link `paths=[...]` exist and a user-provided scope overrides them, dsync MUST:
- print an explicit notice
- print alternate command hints:
  - `--use-link-paths`
  - `--all`

## Clean command safety

- `dsync clean` MUST only target `.dsync-partial/` directories.
- It MUST NOT delete any other file/directory names.
- It MUST still do preview-before-apply behavior.
