# Development

Source of truth specification
- `dsync_go_documentation/README.md`
- `dsync_go_documentation/spec/CLI.md`
- `dsync_go_documentation/spec/BEHAVIORS.md`
- `dsync_go_documentation/TASKS.md`

Prerequisites
- Go (see `go.mod`)
- rsync (for integration testing later)

Common tasks
- Install dev tools: `make tools`
- Run format check + lint + tests: `make check`
- Auto-format: `make fmt`
- Build: `make build`
