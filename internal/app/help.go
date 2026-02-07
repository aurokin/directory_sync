package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer) {
	fmt.Fprint(w, helpText)
}

const helpText = `dsync - rsync-first directory sync (WIP)

Usage:
  dsync <command> [args]

Commands:
  init      Create starter config at $XDG_CONFIG_HOME/dsync/config.toml
  doctor    Validate config and prerequisites (planned)
  ls        List directory via rsync --list-only (planned)
  pull      Sync remote -> local (planned)
  push      Sync local -> remote (planned)
  clean     Remove .dsync-partial directories (planned)

Global flags:
  -h, --help       Show help
  -v, --version    Show version

Project docs:
  dsync_go_documentation/README.md
  dsync_go_documentation/spec/CLI.md
  dsync_go_documentation/spec/BEHAVIORS.md
  dsync_go_documentation/TASKS.md
`
