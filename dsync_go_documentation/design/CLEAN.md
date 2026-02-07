# dsync clean

This document is normative.

`dsync clean` removes rsync partial-transfer directories named `.dsync-partial`.

Motivation
- `--partial --partial-dir=.dsync-partial` enables resumable transfers.
- Interrupted runs may leave behind `.dsync-partial/` directories and files.
- Clean is the "I am done, remove leftover partial state" command.

## Scope

Clean follows the same scope resolution rules as sync commands:
- `relative_path` applies to both ends
- link mode may infer scope from CWD
- link `paths=[...]` can be used with `--use-link-paths`

## Local clean (default)

Default behavior
- Clean local `.dsync-partial/` directories under the resolved local destination root.
- This deletes the directory and its contents.

Preview behavior
- By default, dsync SHOULD show the count of directories/files that would be removed and sample paths.
- `--verbose` MAY print every matched path.

Apply behavior
- Remove matched directories via filesystem operations.
- This is local-only and does not violate the "no remote shell" constraint.

## Remote clean (`--remote`)

Goal
- Delete `.dsync-partial/` directories on the remote side without executing remote shell commands.

Approach (rsync-only deletion)

dsync runs an rsync command where:
- Remote endpoint is the destination.
- Transfer is effectively disabled (no creates/updates).
- Deletion is enabled.
- Receiver-side filter rules protect everything except `.dsync-partial/***`.

Key rsync options (conceptual)
- `--existing --ignore-existing` (skip creating/updating files)
- `--delete --delete-delay`
- Sender-side hide: `--filter 'H .dsync-partial/'` (do not include partial dirs in sender file list)
- Receiver-side risk+protect ordering:
  - `--filter 'R .dsync-partial/***'` (unprotect these so they can be deleted)
  - `--filter 'P *'` (protect everything else from deletion)

IMPORTANT: filter rule order matters
- First match wins. The risk rule MUST come before the protect-all rule.

Limitations
- This strategy will only delete `.dsync-partial` directories that are reachable in the directory traversal.
- If a remote-only directory is entirely protected and does not exist on the sender side, rsync may not descend into it, so `.dsync-partial` inside it may remain.
- This is acceptable for MVP; the main goal is cleaning leftover partial dirs in active scopes.

Preview/apply requirements
- Remote clean MUST follow the same preview-before-apply behavior as sync:
  - preview: `--dry-run --itemize-changes`
  - apply: same command without `--dry-run`

## Excludes interaction

- `dsync clean` MUST ignore user excludes for targeting.
- For remote clean, dsync MUST NOT include `--partial-dir` (it auto-adds protective rules that interfere with deletion).
