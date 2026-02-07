package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestInitWritesConfig(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"init"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d; stderr=%q", code, stderr.String())
	}

	cfgPath := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dsync", "config.toml")
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
	if !bytes.Contains(b, []byte("[global]")) {
		t.Fatalf("expected config to contain [global]")
	}
}

func TestInitRefusesToOverwriteWithoutForce(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := Run([]string{"init"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected first init to succeed: code=%d stderr=%q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"init"}, &stdout, &stderr); code == 0 {
		t.Fatalf("expected second init without --force to fail")
	}
}

func TestInitForceOverwrites(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if code := Run([]string{"init"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected first init to succeed: code=%d stderr=%q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	if code := Run([]string{"init", "--force"}, &stdout, &stderr); code != 0 {
		t.Fatalf("expected init --force to succeed: code=%d stderr=%q", code, stderr.String())
	}
}
