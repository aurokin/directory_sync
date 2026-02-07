package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/aurokin/directory_sync/internal/version"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		printHelp(stdout)
		return 0
	case "-v", "--version", "version":
		fmt.Fprintf(stdout, "dsync %s\n", version.Version)
		return 0
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "init":
		return runInit(cmdArgs, stdout, stderr)
	case "doctor", "ls", "pull", "push", "clean":
		fmt.Fprintf(stderr, "%s: not implemented yet\n", cmd)
		fmt.Fprintf(stderr, "See dsync_go_documentation/TASKS.md\n")
		return 2
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		if suggestion := suggest(cmd); suggestion != "" {
			fmt.Fprintf(stderr, "Did you mean: %s?\n", suggestion)
		}
		printHelp(stderr)
		return 2
	}
}

func suggest(cmd string) string {
	known := []string{"init", "doctor", "ls", "pull", "push", "clean"}
	cmd = strings.ToLower(cmd)
	for _, k := range known {
		if strings.HasPrefix(k, cmd) {
			return k
		}
	}
	return ""
}
