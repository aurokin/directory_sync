package xdg

import (
	"path/filepath"
	"testing"
)

func TestConfigFilePathUsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(string(filepath.Separator), "tmp", "xdg"))
	got, err := ConfigFilePath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join(string(filepath.Separator), "tmp", "xdg", "dsync", "config.toml")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestConfigSearchPathsIncludesFallback(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(string(filepath.Separator), "tmp", "xdg"))
	paths, err := ConfigSearchPaths()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 search paths, got %d: %v", len(paths), paths)
	}
}
