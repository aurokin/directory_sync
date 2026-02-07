package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfigFilePrefersXDG(t *testing.T) {
	home := t.TempDir()
	xdgHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	primary := filepath.Join(xdgHome, "dsync", "config.toml")
	fallback := filepath.Join(home, ".config", "dsync", "config.toml")

	if err := os.MkdirAll(filepath.Dir(primary), 0o755); err != nil {
		t.Fatalf("mkdir primary: %v", err)
	}
	if err := os.WriteFile(primary, []byte("[global]\n"), 0o644); err != nil {
		t.Fatalf("write primary: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(fallback), 0o755); err != nil {
		t.Fatalf("mkdir fallback: %v", err)
	}
	if err := os.WriteFile(fallback, []byte("[global]\n"), 0o644); err != nil {
		t.Fatalf("write fallback: %v", err)
	}

	got, err := FindConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != primary {
		t.Fatalf("expected %q, got %q", primary, got)
	}
}

func TestFindConfigFileFallsBackToHome(t *testing.T) {
	home := t.TempDir()
	xdgHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	fallback := filepath.Join(home, ".config", "dsync", "config.toml")
	if err := os.MkdirAll(filepath.Dir(fallback), 0o755); err != nil {
		t.Fatalf("mkdir fallback: %v", err)
	}
	if err := os.WriteFile(fallback, []byte("[global]\n"), 0o644); err != nil {
		t.Fatalf("write fallback: %v", err)
	}

	got, err := FindConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != fallback {
		t.Fatalf("expected %q, got %q", fallback, got)
	}
}

func TestFindConfigFileNotFound(t *testing.T) {
	home := t.TempDir()
	xdgHome := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", xdgHome)

	_, err := FindConfigFile()
	if err == nil {
		t.Fatalf("expected error")
	}
	var nf NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected NotFoundError, got %T: %v", err, err)
	}
}
