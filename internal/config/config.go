package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aurokin/directory_sync/internal/xdg"
	"github.com/pelletier/go-toml/v2"
)

type EndpointType string

const (
	EndpointTypeLocal EndpointType = "local"
	EndpointTypeSSH   EndpointType = "ssh"
)

type Global struct {
	Excludes []string
}

type Endpoint struct {
	Name     string
	Type     EndpointType
	Host     string
	Path     string
	RootPath string // always trailing '/'
}

// RsyncRoot returns the root suitable for rsync SRC/DEST.
//
// Contents semantics: the returned string always ends in '/'.
func (e Endpoint) RsyncRoot() string {
	if e.Type == EndpointTypeSSH {
		return e.Host + ":" + e.RootPath
	}
	return e.RootPath
}

type Link struct {
	Name        string
	LocalName   string
	RemoteName  string
	Mirror      bool
	PartialOnly bool
	Paths       []string
	Excludes    []string

	Local  Endpoint
	Remote Endpoint
}

type Config struct {
	FilePath  string
	Global    Global
	Endpoints map[string]Endpoint
	Links     map[string]Link
}

type rawConfig struct {
	Global    rawGlobal              `toml:"global"`
	Endpoints map[string]rawEndpoint `toml:"endpoints"`
	Links     map[string]rawLink     `toml:"links"`
}

type rawGlobal struct {
	Excludes []string `toml:"excludes"`
}

type rawEndpoint struct {
	Type string `toml:"type"`
	Path string `toml:"path"`
	Host string `toml:"host"`
}

type rawLink struct {
	Local       string   `toml:"local"`
	Remote      string   `toml:"remote"`
	Mirror      *bool    `toml:"mirror"`
	PartialOnly *bool    `toml:"partial_only"`
	Paths       []string `toml:"paths"`
	Excludes    []string `toml:"excludes"`
}

type NotFoundError struct {
	Searched []string
}

func (e NotFoundError) Error() string {
	if len(e.Searched) == 0 {
		return "config not found (run 'dsync init' to create one)"
	}
	return fmt.Sprintf(
		"config not found (searched: %s). Run 'dsync init' to create one.",
		strings.Join(e.Searched, ", "),
	)
}

type ValidationError struct {
	Issues []string
}

func (e ValidationError) Error() string {
	if len(e.Issues) == 0 {
		return "config validation failed"
	}
	var b strings.Builder
	b.WriteString("config validation failed:\n")
	for _, issue := range e.Issues {
		b.WriteString("- ")
		b.WriteString(issue)
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func Load() (*Config, error) {
	path, err := FindConfigFile()
	if err != nil {
		return nil, err
	}
	return LoadFromPath(path)
}

func FindConfigFile() (string, error) {
	paths, err := xdg.ConfigSearchPaths()
	if err != nil {
		return "", err
	}
	for _, p := range paths {
		info, statErr := os.Stat(p)
		if statErr == nil {
			if info.Mode().IsRegular() {
				return p, nil
			}
			continue
		}
		if errors.Is(statErr, os.ErrNotExist) {
			continue
		}
		return "", fmt.Errorf("stat config %s: %w", p, statErr)
	}
	return "", NotFoundError{Searched: paths}
}

func LoadFromPath(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	cfg, err := Parse(b)
	if err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	cfg.FilePath = path
	return cfg, nil
}

func Parse(b []byte) (*Config, error) {
	var raw rawConfig
	if err := toml.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return normalize(raw)
}

func normalize(raw rawConfig) (*Config, error) {
	issues := make([]string, 0)

	cfg := &Config{
		Global: Global{Excludes: nilToEmpty(raw.Global.Excludes)},
	}

	// Endpoints
	endpoints := make(map[string]Endpoint, len(raw.Endpoints))
	for name, re := range raw.Endpoints {
		e, epIssues := normalizeEndpoint(name, re)
		issues = append(issues, epIssues...)
		if len(epIssues) == 0 {
			endpoints[name] = e
		}
	}
	if len(endpoints) == 0 {
		issues = append(issues, "no endpoints defined (missing [endpoints.<name>] sections)")
	}
	cfg.Endpoints = endpoints

	// Links
	links := make(map[string]Link, len(raw.Links))
	for name, rl := range raw.Links {
		l, linkIssues := normalizeLink(name, rl, endpoints)
		issues = append(issues, linkIssues...)
		if len(linkIssues) == 0 {
			links[name] = l
		}
	}
	cfg.Links = links

	if len(issues) > 0 {
		return nil, ValidationError{Issues: issues}
	}
	return cfg, nil
}

func normalizeEndpoint(name string, re rawEndpoint) (Endpoint, []string) {
	issues := make([]string, 0)

	name = strings.TrimSpace(name)
	if name == "" {
		return Endpoint{}, []string{"endpoint name cannot be empty"}
	}

	t := strings.TrimSpace(strings.ToLower(re.Type))
	path := strings.TrimSpace(re.Path)
	host := strings.TrimSpace(re.Host)

	var et EndpointType
	switch EndpointType(t) {
	case EndpointTypeLocal:
		et = EndpointTypeLocal
		if host != "" {
			issues = append(issues, fmt.Sprintf("endpoints.%s.host is set but type is local", name))
		}
	case EndpointTypeSSH:
		et = EndpointTypeSSH
		if host == "" {
			issues = append(issues, fmt.Sprintf("endpoints.%s.host is required for ssh endpoints", name))
		}
	default:
		issues = append(issues, fmt.Sprintf("endpoints.%s.type must be 'local' or 'ssh'", name))
	}

	if path == "" {
		issues = append(issues, fmt.Sprintf("endpoints.%s.path is required", name))
		return Endpoint{}, issues
	}
	if !strings.HasPrefix(path, string(filepath.Separator)) {
		issues = append(issues, fmt.Sprintf("endpoints.%s.path must be an absolute path", name))
		return Endpoint{}, issues
	}

	clean := filepath.Clean(path)
	if clean == string(filepath.Separator) {
		issues = append(issues, fmt.Sprintf("endpoints.%s.path must not be '/'", name))
		return Endpoint{}, issues
	}
	root := clean
	if !strings.HasSuffix(root, string(filepath.Separator)) {
		root += string(filepath.Separator)
	}

	return Endpoint{
		Name:     name,
		Type:     et,
		Host:     host,
		Path:     clean,
		RootPath: root,
	}, issues
}

func normalizeLink(name string, rl rawLink, endpoints map[string]Endpoint) (Link, []string) {
	issues := make([]string, 0)

	name = strings.TrimSpace(name)
	if name == "" {
		return Link{}, []string{"link name cannot be empty"}
	}

	localName := strings.TrimSpace(rl.Local)
	remoteName := strings.TrimSpace(rl.Remote)
	if localName == "" {
		issues = append(issues, fmt.Sprintf("links.%s.local is required", name))
	}
	if remoteName == "" {
		issues = append(issues, fmt.Sprintf("links.%s.remote is required", name))
	}

	localEP, okLocal := endpoints[localName]
	if localName != "" && !okLocal {
		issues = append(issues, fmt.Sprintf("links.%s.local references unknown endpoint %q", name, localName))
	}
	remoteEP, okRemote := endpoints[remoteName]
	if remoteName != "" && !okRemote {
		issues = append(issues, fmt.Sprintf("links.%s.remote references unknown endpoint %q", name, remoteName))
	}

	if okLocal {
		if localEP.Type != EndpointTypeLocal {
			issues = append(issues, fmt.Sprintf("links.%s.local endpoint %q must be type local", name, localName))
		}
	}
	if okRemote {
		if remoteEP.Type != EndpointTypeSSH {
			issues = append(issues, fmt.Sprintf("links.%s.remote endpoint %q must be type ssh", name, remoteName))
		}
	}

	// Exactly one remote side for a link (MVP).
	if okLocal && okRemote {
		if localEP.Type == remoteEP.Type {
			issues = append(issues, fmt.Sprintf("links.%s must connect one local endpoint and one ssh endpoint", name))
		}
		if localName == remoteName {
			issues = append(issues, fmt.Sprintf("links.%s.local and links.%s.remote must be different endpoints", name, name))
		}
	}

	// Defaults
	mirror := true
	if rl.Mirror != nil {
		mirror = *rl.Mirror
	}
	partialOnly := false
	if rl.PartialOnly != nil {
		partialOnly = *rl.PartialOnly
	}

	paths, pathIssues := normalizeScopes(name, rl.Paths)
	issues = append(issues, pathIssues...)

	return Link{
		Name:        name,
		LocalName:   localName,
		RemoteName:  remoteName,
		Mirror:      mirror,
		PartialOnly: partialOnly,
		Paths:       paths,
		Excludes:    nilToEmpty(rl.Excludes),
		Local:       localEP,
		Remote:      remoteEP,
	}, issues
}

func normalizeScopes(linkName string, in []string) ([]string, []string) {
	issues := make([]string, 0)
	if in == nil {
		return []string{}, nil
	}

	out := make([]string, 0, len(in))
	for i, p := range in {
		raw := strings.TrimSpace(p)
		if raw == "" {
			issues = append(issues, fmt.Sprintf("links.%s.paths[%d] must not be empty", linkName, i))
			continue
		}
		if strings.HasPrefix(raw, string(filepath.Separator)) {
			issues = append(issues, fmt.Sprintf("links.%s.paths[%d] must be a relative path", linkName, i))
			continue
		}

		clean := filepath.Clean(raw)
		if clean == "." {
			issues = append(issues, fmt.Sprintf("links.%s.paths[%d] resolves to empty scope ('.'); omit it to use full-root with --all", linkName, i))
			continue
		}
		if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			issues = append(issues, fmt.Sprintf("links.%s.paths[%d] must not traverse outside the root", linkName, i))
			continue
		}
		out = append(out, clean)
	}
	return out, issues
}

func nilToEmpty[T any](in []T) []T {
	if in == nil {
		return []T{}
	}
	return in
}
