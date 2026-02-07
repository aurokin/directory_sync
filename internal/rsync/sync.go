package rsync

import (
	"fmt"
	"strings"
)

type SyncSpec struct {
	Source   string
	Dest     string
	UseSSH   bool
	Mirror   bool
	Excludes []string
	DryRun   bool
}

func BuildArgs(spec SyncSpec) ([]string, error) {
	if spec.Source == "" {
		return nil, fmt.Errorf("source is required")
	}
	if spec.Dest == "" {
		return nil, fmt.Errorf("dest is required")
	}
	if !strings.HasSuffix(spec.Source, "/") {
		return nil, fmt.Errorf("source must end with '/' (contents semantics): %q", spec.Source)
	}
	if !strings.HasSuffix(spec.Dest, "/") {
		return nil, fmt.Errorf("dest must end with '/' (contents semantics): %q", spec.Dest)
	}

	args := []string{
		"-a",
		"--no-owner",
		"--no-group",
		"--mkpath",
		"--protect-args",
		"--partial",
		"--partial-dir=.dsync-partial",
		"--human-readable",
		"--stats",
		"--itemize-changes",
	}

	if spec.Mirror {
		args = append(args, "--delete", "--delete-delay")
	}

	if spec.UseSSH {
		args = append(args, "-e", "ssh")
	}

	for _, ex := range spec.Excludes {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		args = append(args, "--exclude", ex)
	}

	if spec.DryRun {
		args = append(args, "--dry-run")
	}

	args = append(args, spec.Source, spec.Dest)
	return args, nil
}

func BuildPreviewApply(spec SyncSpec) (preview []string, apply []string, err error) {
	spec.DryRun = true
	preview, err = BuildArgs(spec)
	if err != nil {
		return nil, nil, err
	}
	spec.DryRun = false
	apply, err = BuildArgs(spec)
	if err != nil {
		return nil, nil, err
	}
	return preview, apply, nil
}
