package app

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aurokin/directory_sync/internal/config"
	"github.com/aurokin/directory_sync/internal/xdg"
)

func runInit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)

	force := fs.Bool("force", false, "overwrite existing config")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	cfgPath, err := xdg.ConfigFilePath()
	if err != nil {
		fmt.Fprintf(stderr, "init: unable to resolve config path: %v\n", err)
		return 1
	}

	if _, statErr := os.Stat(cfgPath); statErr == nil && !*force {
		fmt.Fprintf(stderr, "init: config already exists: %s\n", cfgPath)
		fmt.Fprintf(stderr, "init: re-run with --force to overwrite\n")
		return 1
	} else if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		fmt.Fprintf(stderr, "init: unable to stat config: %v\n", statErr)
		return 1
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0o755); err != nil {
		fmt.Fprintf(stderr, "init: unable to create config directory: %v\n", err)
		return 1
	}

	if err := os.WriteFile(cfgPath, []byte(config.DefaultConfigTemplate), 0o644); err != nil {
		fmt.Fprintf(stderr, "init: unable to write config: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Wrote %s\n", cfgPath)
	return 0
}
