package config

import (
	"strings"
	"testing"
)

func TestParseValidConfigNormalizesDefaultsAndRoots(t *testing.T) {
	toml := `
[global]
excludes = [".DS_Store"]

[endpoints.laptop]
type = "local"
path = "/tmp/photos"

[endpoints.server]
type = "ssh"
host = "photo-box"
path = "/srv/photos/"

[links.photos]
local = "laptop"
remote = "server"
# mirror omitted on purpose (defaults true)
`

	cfg, err := Parse([]byte(toml))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	laptop := cfg.Endpoints["laptop"]
	if laptop.RootPath != "/tmp/photos/" {
		t.Fatalf("expected trailing-slash root, got %q", laptop.RootPath)
	}

	server := cfg.Endpoints["server"]
	if server.RootPath != "/srv/photos/" {
		t.Fatalf("expected cleaned root, got %q", server.RootPath)
	}
	if got := server.RsyncRoot(); got != "photo-box:/srv/photos/" {
		t.Fatalf("unexpected rsync root: %q", got)
	}

	link := cfg.Links["photos"]
	if !link.Mirror {
		t.Fatalf("expected mirror default true")
	}
	if link.PartialOnly {
		t.Fatalf("expected partial_only default false")
	}
}

func TestParseRejectsEndpointRootSlash(t *testing.T) {
	toml := `
[endpoints.bad]
type = "local"
path = "/"
`
	_, err := Parse([]byte(toml))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "must not be '/'") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsLinkWithWrongEndpointTypes(t *testing.T) {
	toml := `
[endpoints.a]
type = "ssh"
host = "x"
path = "/srv/a"

[endpoints.b]
type = "local"
path = "/tmp/b"

[links.l]
local = "a"
remote = "b"
`
	_, err := Parse([]byte(toml))
	if err == nil {
		t.Fatalf("expected error")
	}

	msg := err.Error()
	if !strings.Contains(msg, "must be type local") || !strings.Contains(msg, "must be type ssh") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsLinkPathsThatTraverseUp(t *testing.T) {
	toml := `
[endpoints.laptop]
type = "local"
path = "/tmp/photos"

[endpoints.server]
type = "ssh"
host = "photo-box"
path = "/srv/photos"

[links.photos]
local = "laptop"
remote = "server"
paths = ["../oops"]
`
	_, err := Parse([]byte(toml))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "must not traverse") {
		t.Fatalf("unexpected error: %v", err)
	}
}
