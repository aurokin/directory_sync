# Overview

`dsync` is a CLI that provides a config-driven, safety-focused wrapper around `rsync`.

The goal is a predictable "named endpoints + push/pull" experience for syncing directories between:
- your current working directory (CWD) and a named endpoint, or
- two named endpoints via a named link.

`dsync` is intentionally opinionated.

## Goals

- Simple mental model: endpoints, links, push/pull.
- Rsync-native semantics: avoid inventing new path rules.
- Mirror-by-default: destination matches source within the chosen scope.
- Strong guardrails:
  - always preview first
  - require explicit opt-in for full-root operations (`--all`)
  - refuse to run in non-interactive contexts unless `--yes` is provided
- Clear logs:
  - show direction (source -> destination)
  - show what would transfer and what would delete
  - provide machine-readable NDJSON events (`--json`) for agents and debugging

## Non-goals (MVP)

- No backup/rollback features (e.g., rsync `--backup`, `--backup-dir`). Backups belong at a different layer.
- No non-rsync transport backends.
- No Windows support guarantees.
- No Google Drive backend in MVP.

## Required environment assumptions

- `rsync` available on both ends, kept up to date (assume rsync 3.2.x).
- For remote endpoints, SSH connectivity is configured via `~/.ssh/config`.
  - dsync uses ssh host aliases (e.g., `photo-box`) and relies on SSH config for port/user/identity.

## Glossary

- Endpoint: A named root directory. Type is `local` or `ssh`.
- Link: A named mapping between one local endpoint and one remote endpoint (1:1 in MVP).
- Scope: A relative path that is appended to both endpoint roots to target a subtree.
- Mirror: Destination is made to match source inside the scope; includes deletion of destination extras.
- Preview: A dry-run rsync invocation that prints "what would happen".
- Apply: The real rsync invocation that performs the changes.
