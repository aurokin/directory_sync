# Agent Notes

If you are an agent working on this repo:

Read first
- `dsync_go_documentation/agent/CONTEXT.md`

Key invariants
- rsync-only engine; no standalone remote shell commands
- preview-before-apply always (including with --yes)
- full-root apply requires --all
- non-interactive mutating runs require --yes
- --json emits NDJSON to stdout; human logs to stderr

Local workflow
- `make check`
