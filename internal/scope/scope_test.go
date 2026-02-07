package scope

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/aurokin/directory_sync/internal/config"
)

func TestResolveForLinkConflicts(t *testing.T) {
	link := config.Link{
		Name:  "photos",
		Paths: []string{"2026/portraits"},
		Local: config.Endpoint{Path: "/tmp/photos"},
	}

	_, err := ResolveForLink(link, Options{
		Command:         "pull",
		LinkName:        "photos",
		UseLinkPaths:    true,
		HasRelativePath: true,
		RelativePath:    "2026/portraits",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveForLinkUseLinkPathsOverridesCWD(t *testing.T) {
	link := config.Link{
		Name:  "photos",
		Paths: []string{"2026/portraits", "2026/events"},
		Local: config.Endpoint{Path: "/tmp/photos"},
	}

	res, err := ResolveForLink(link, Options{
		Command:      "pull",
		LinkName:     "photos",
		UseLinkPaths: true,
		CWD:          filepath.Join(string(filepath.Separator), "tmp", "photos", "2026", "portraits"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Batch {
		t.Fatalf("expected batch true")
	}
	if res.Source != SourceLinkPath {
		t.Fatalf("expected source %q, got %q", SourceLinkPath, res.Source)
	}
	if len(res.Scopes) != 2 {
		t.Fatalf("expected 2 scopes, got %v", res.Scopes)
	}
}

func TestResolveForLinkInfersScopeFromCWD(t *testing.T) {
	link := config.Link{
		Name:  "photos",
		Paths: []string{"2026/portraits"},
		Local: config.Endpoint{Path: "/tmp/photos"},
	}

	res, err := ResolveForLink(link, Options{
		Command:  "pull",
		LinkName: "photos",
		CWD:      "/tmp/photos/2026/portraits",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != SourceCWD {
		t.Fatalf("expected source %q, got %q", SourceCWD, res.Source)
	}
	if res.ScopeString != filepath.Join("2026", "portraits") {
		t.Fatalf("unexpected inferred scope: %q", res.ScopeString)
	}
	if !res.NoticesHaveMessageContaining("ignored") {
		t.Fatalf("expected mismatch notice due to configured link paths")
	}
}

func TestResolveForLinkUsesCLIScope(t *testing.T) {
	link := config.Link{
		Name:  "photos",
		Paths: []string{"2026/portraits"},
		Local: config.Endpoint{Path: "/tmp/photos"},
	}

	res, err := ResolveForLink(link, Options{
		Command:         "pull",
		LinkName:        "photos",
		HasRelativePath: true,
		RelativePath:    "2026/portraits",
		CWD:             "/tmp/photos/2026/portraits",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Source != SourceCLI {
		t.Fatalf("expected source %q, got %q", SourceCLI, res.Source)
	}
	if res.ScopeString != filepath.Join("2026", "portraits") {
		t.Fatalf("unexpected scope: %q", res.ScopeString)
	}
}

func TestResolveForEndpointNormalizesScope(t *testing.T) {
	res, err := ResolveForEndpoint(Options{HasRelativePath: true, RelativePath: "./2026/portraits/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ScopeString != filepath.Join("2026", "portraits") {
		t.Fatalf("unexpected scope: %q", res.ScopeString)
	}
}

func TestResolveForLinkPartialOnlyRejectsFullRoot(t *testing.T) {
	link := config.Link{
		Name:        "photos",
		PartialOnly: true,
		Paths:       []string{"2026/portraits"},
		Local:       config.Endpoint{Path: "/tmp/photos"},
	}

	if _, err := ResolveForLink(link, Options{Command: "pull", LinkName: "photos"}); err == nil {
		t.Fatalf("expected error")
	}

	if _, err := ResolveForLink(link, Options{Command: "pull", LinkName: "photos", All: true}); err == nil {
		t.Fatalf("expected error")
	}

	if _, err := ResolveForLink(link, Options{Command: "pull", LinkName: "photos", HasRelativePath: true, RelativePath: "2026/portraits"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := ResolveForLink(link, Options{Command: "pull", LinkName: "photos", UseLinkPaths: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func (r Result) NoticesHaveMessageContaining(s string) bool {
	for _, n := range r.Notices {
		if strings.Contains(n.Message, s) {
			return true
		}
	}
	return false
}
