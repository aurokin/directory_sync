package xdg

import (
	"fmt"
	"os"
	"path/filepath"
)

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
