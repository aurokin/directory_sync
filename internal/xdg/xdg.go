package xdg

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigFilePath returns the preferred config path for writing.
//
// Discovery (reading) should use ConfigSearchPaths.
func ConfigFilePath() (string, error) {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, "dsync", "config.toml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "dsync", "config.toml"), nil
}

// ConfigSearchPaths returns candidate config paths in order.
//
// Order (MVP):
// - $XDG_CONFIG_HOME/dsync/config.toml (if set)
// - ~/.config/dsync/config.toml
func ConfigSearchPaths() ([]string, error) {
	var paths []string
	seen := map[string]bool{}

	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		p := filepath.Join(v, "dsync", "config.toml")
		if !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home dir: %w", err)
	}
	fallback := filepath.Join(home, ".config", "dsync", "config.toml")
	if !seen[fallback] {
		paths = append(paths, fallback)
	}

	return paths, nil
}
