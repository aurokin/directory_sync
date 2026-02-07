package scope

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aurokin/directory_sync/internal/config"
)

type Source string

const (
	SourceCLI      Source = "cli"
	SourceCWD      Source = "cwd"
	SourceLinkPath Source = "link_paths"
	SourceEmpty    Source = "empty"
)

type Notice struct {
	Level   string
	Message string
	Hints   []string
}

type Result struct {
	Scopes      []string
	Source      Source
	Batch       bool
	IsFullRoot  bool
	Notices     []Notice
	ScopeString string // resolved scope for single-scope operations ("" for full-root)
}

type Options struct {
	Command         string
	LinkName        string
	RelativePath    string
	HasRelativePath bool
	UseLinkPaths    bool
	All             bool
	CWD             string
}

func ResolveForEndpoint(opts Options) (Result, error) {
	res := Result{}

	if opts.UseLinkPaths {
		return Result{}, fmt.Errorf("--use-link-paths is only valid with --link")
	}

	scope, notices, err := resolveCLIScope(opts)
	if err != nil {
		return Result{}, err
	}
	res.Notices = append(res.Notices, notices...)

	res.Batch = false
	res.Source = SourceEmpty
	if opts.HasRelativePath {
		res.Source = SourceCLI
	}
	res.ScopeString = scope
	res.Scopes = []string{scope}
	res.IsFullRoot = scope == ""
	return res, nil
}

func ResolveForLink(link config.Link, opts Options) (Result, error) {
	res := Result{}

	if opts.UseLinkPaths {
		if opts.HasRelativePath {
			return Result{}, fmt.Errorf("--use-link-paths conflicts with a relative_path argument")
		}
		if opts.All {
			return Result{}, fmt.Errorf("--use-link-paths conflicts with --all")
		}
		if len(link.Paths) == 0 {
			return Result{}, fmt.Errorf("link %q has no configured paths", link.Name)
		}

		res.Batch = true
		res.Source = SourceLinkPath
		res.Scopes = append([]string{}, link.Paths...)
		res.IsFullRoot = false
		res.Notices = append(res.Notices, Notice{
			Level:   "info",
			Message: "Using configured link paths; ignoring CWD scope inference",
		})
		return res, nil
	}

	scope, notices, err := resolveCLIScope(opts)
	if err != nil {
		return Result{}, err
	}
	res.Notices = append(res.Notices, notices...)

	if opts.HasRelativePath {
		res.Batch = false
		res.Source = SourceCLI
		res.ScopeString = scope
		res.Scopes = []string{scope}
		res.IsFullRoot = scope == ""
		if res.IsFullRoot && link.PartialOnly {
			return Result{}, fmt.Errorf("link %q is partial_only; provide a non-empty scope or use --use-link-paths", link.Name)
		}
		res.Notices = append(res.Notices, mismatchNoticesIfAny(link, opts)...)
		return res, nil
	}

	// When --all is set without an explicit relative_path, treat it as an explicit
	// full-root intent. This does not grant permission to apply; it only makes the
	// user's intent explicit in the plan.
	if opts.All {
		if link.PartialOnly {
			return Result{}, fmt.Errorf("link %q is partial_only; full-root operations are forbidden", link.Name)
		}
		res.Batch = false
		res.Source = SourceEmpty
		res.ScopeString = ""
		res.Scopes = []string{""}
		res.IsFullRoot = true
		res.Notices = append(res.Notices, Notice{
			Level:   "info",
			Message: "Full-root operation explicitly requested via --all",
		})
		return res, nil
	}

	if opts.CWD != "" {
		if inferred, ok, err := inferScopeFromCWD(link.Local.Path, opts.CWD); err != nil {
			return Result{}, err
		} else if ok {
			res.Batch = false
			res.Source = SourceCWD
			res.ScopeString = inferred
			res.Scopes = []string{inferred}
			res.IsFullRoot = inferred == ""
			if res.IsFullRoot && link.PartialOnly {
				return Result{}, fmt.Errorf("link %q is partial_only; provide a non-empty scope or use --use-link-paths", link.Name)
			}
			res.Notices = append(res.Notices, mismatchNoticesIfAny(link, opts)...)
			return res, nil
		}
	}

	// No CLI scope, no CWD inference, no link-path batch requested.
	res.Batch = false
	res.Source = SourceEmpty
	res.ScopeString = ""
	res.Scopes = []string{""}
	res.IsFullRoot = true
	if link.PartialOnly {
		return Result{}, fmt.Errorf("link %q is partial_only; provide a non-empty scope or use --use-link-paths", link.Name)
	}
	return res, nil
}

func resolveCLIScope(opts Options) (string, []Notice, error) {
	if !opts.HasRelativePath {
		return "", nil, nil
	}
	clean, err := normalizeScope(opts.RelativePath)
	if err != nil {
		return "", nil, err
	}

	var notices []Notice
	if opts.All && clean != "" {
		notices = append(notices, Notice{
			Level:   "info",
			Message: "Ignoring --all because a non-empty scope was provided",
		})
	}
	return clean, notices, nil
}

func normalizeScope(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	if strings.HasPrefix(raw, string(filepath.Separator)) {
		return "", fmt.Errorf("scope must be a relative path")
	}

	clean := filepath.Clean(raw)
	if clean == "." {
		return "", nil
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("scope must not traverse outside the root")
	}
	return clean, nil
}

func inferScopeFromCWD(localRoot, cwd string) (string, bool, error) {
	if localRoot == "" || cwd == "" {
		return "", false, nil
	}

	root := filepath.Clean(localRoot)
	cur := filepath.Clean(cwd)
	rel, err := filepath.Rel(root, cur)
	if err != nil {
		return "", false, fmt.Errorf("infer scope from cwd: %w", err)
	}
	if rel == "." {
		return "", true, nil
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false, nil
	}
	return rel, true, nil
}

func mismatchNoticesIfAny(link config.Link, opts Options) []Notice {
	if len(link.Paths) == 0 {
		return nil
	}
	if opts.Command == "" {
		return nil
	}
	if opts.LinkName == "" {
		return nil
	}

	// Only emit mismatch notices when a user scope (CLI or inferred) is used.
	if opts.HasRelativePath || opts.CWD != "" {
		cmdUse := fmt.Sprintf("dsync %s --link %s --use-link-paths", opts.Command, opts.LinkName)
		cmdAll := fmt.Sprintf("dsync %s --link %s --all", opts.Command, opts.LinkName)
		return []Notice{{
			Level:   "info",
			Message: fmt.Sprintf("Link %q has configured paths (%d) that were ignored for this run", link.Name, len(link.Paths)),
			Hints: []string{
				"To sync the configured link paths instead:",
				cmdUse,
				"To sync the full link root (mirror/delete across everything):",
				cmdAll,
			},
		}}
	}
	return nil
}
