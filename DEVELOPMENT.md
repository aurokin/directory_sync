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

Scripts
- Local health checks: `scripts/health.sh` (runs `make check`)
- Run dsync from source: `scripts/dev.sh <dsync args>`

Current CLI status
- `dsync init` works
- `dsync doctor` loads/validates config (probes are planned)
- `dsync pull` / `dsync push` print rsync plans; execution is Phase 5
