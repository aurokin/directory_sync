package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aurokin/directory_sync/internal/config"
	"github.com/aurokin/directory_sync/internal/rsync"
	"github.com/aurokin/directory_sync/internal/scope"
)

type syncArgs struct {
	link         bool
	useLinkPaths bool
	all          bool
	dryRun       bool
	yes          bool
	verbose      bool
	json         bool
	dangerous    bool
}

func parseSyncArgs(cmd string, args []string, stdout, stderr io.Writer) (syncArgs, []string, int, bool) {
	var s syncArgs
	pos := make([]string, 0, 2)

	stopFlags := false
	for _, a := range args {
		if stopFlags {
			pos = append(pos, a)
			continue
		}
		if a == "--" {
			stopFlags = true
			continue
		}
		if a == "-h" || a == "--help" {
			printSyncHelp(cmd, stdout)
			return syncArgs{}, nil, 0, true
		}
		if strings.HasPrefix(a, "--") {
			switch a {
			case "--link":
				s.link = true
			case "--use-link-paths":
				s.useLinkPaths = true
			case "--all":
				s.all = true
			case "--dry-run":
				s.dryRun = true
			case "--yes":
				s.yes = true
			case "--verbose":
				s.verbose = true
			case "--json":
				s.json = true
			case "--dangerous":
				s.dangerous = true
			default:
				fmt.Fprintf(stderr, "%s: unknown flag %s\n", cmd, a)
				printSyncUsage(cmd, stderr)
				return syncArgs{}, nil, 2, false
			}
			continue
		}
		if strings.HasPrefix(a, "-") {
			fmt.Fprintf(stderr, "%s: unknown flag %s\n", cmd, a)
			printSyncUsage(cmd, stderr)
			return syncArgs{}, nil, 2, false
		}
		pos = append(pos, a)
	}

	if len(pos) < 1 || len(pos) > 2 {
		printSyncUsage(cmd, stderr)
		return syncArgs{}, nil, 2, false
	}

	return s, pos, 0, false
}

func runSync(cmd string, args []string, stdout, stderr io.Writer) int {
	flags, pos, code, printedHelp := parseSyncArgs(cmd, args, stdout, stderr)
	if printedHelp {
		return code
	}
	if code != 0 {
		return code
	}

	name := pos[0]
	var rel string
	hasRel := false
	if len(pos) == 2 {
		rel = pos[1]
		hasRel = true
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "%s: unable to resolve cwd: %v\n", cmd, err)
		return 1
	}

	cfg, err := config.Load()
	if err != nil {
		var nf config.NotFoundError
		if errors.As(err, &nf) {
			fmt.Fprintf(stderr, "%v\n", err)
			return 1
		}
		var ve config.ValidationError
		if errors.As(err, &ve) {
			fmt.Fprintf(stderr, "%v\n", err)
			return 1
		}
		fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
		return 1
	}

	if flags.link {
		lnk, ok := cfg.Links[name]
		if !ok {
			fmt.Fprintf(stderr, "%s: unknown link %q\n", cmd, name)
			return 1
		}
		res, err := scope.ResolveForLink(lnk, scope.Options{
			Command:         cmd,
			LinkName:        name,
			RelativePath:    rel,
			HasRelativePath: hasRel,
			UseLinkPaths:    flags.useLinkPaths,
			All:             flags.all,
			CWD:             cwd,
		})
		if err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
			return 2
		}

		printScopePlan(cmd, "link", name, res, stdout, stderr)

		if err := printRsyncPlansForLink(cmd, cfg, lnk, res, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
			return 1
		}
		if !flags.dryRun && res.IsFullRoot && !flags.all {
			fmt.Fprintf(stderr, "%s: full-root operation requires --all\n", cmd)
			fmt.Fprintf(stderr, "hint: re-run with --all, or provide a scope (relative_path or CWD inference)\n")
			return 2
		}

		if flags.dryRun {
			fmt.Fprintf(stdout, "%s: dry-run only (rsync execution not implemented yet)\n", cmd)
			return 0
		}
		fmt.Fprintf(stderr, "%s: rsync execution not implemented yet (use --dry-run for planning)\n", cmd)
		return 2
	}

	ep, ok := cfg.Endpoints[name]
	if !ok {
		fmt.Fprintf(stderr, "%s: unknown endpoint %q\n", cmd, name)
		return 1
	}
	res, err := scope.ResolveForEndpoint(scope.Options{
		Command:         cmd,
		RelativePath:    rel,
		HasRelativePath: hasRel,
		UseLinkPaths:    flags.useLinkPaths,
		All:             flags.all,
		CWD:             cwd,
	})
	if err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
		return 2
	}

	printScopePlan(cmd, "endpoint", name, res, stdout, stderr)
	if err := printRsyncPlansForEndpoint(cmd, cfg, ep, res, cwd, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "%s: %v\n", cmd, err)
		return 1
	}
	if !flags.dryRun && res.IsFullRoot && !flags.all {
		fmt.Fprintf(stderr, "%s: full-root operation requires --all\n", cmd)
		fmt.Fprintf(stderr, "hint: re-run with --all, or provide a scope (relative_path)\n")
		return 2
	}

	_ = ep
	if flags.dryRun {
		fmt.Fprintf(stdout, "%s: dry-run only (rsync execution not implemented yet)\n", cmd)
		return 0
	}
	fmt.Fprintf(stderr, "%s: rsync execution not implemented yet (use --dry-run for planning)\n", cmd)
	return 2
}

func printSyncUsage(cmd string, w io.Writer) {
	fmt.Fprintf(w, "usage: dsync %s [flags] <name> [relative_path]\n", cmd)
	fmt.Fprintf(w, "       dsync %s [flags] --link <link> [relative_path]\n", cmd)
	fmt.Fprintf(w, "hint: use '--' to separate flags from a scope that starts with '-'\n")
}

func printSyncHelp(cmd string, w io.Writer) {
	fmt.Fprintf(w, "dsync %s\n\n", cmd)
	printSyncUsage(cmd, w)
	fmt.Fprintf(w, "\nFlags:\n")
	fmt.Fprintf(w, "  --link            Treat <name> as a link name\n")
	fmt.Fprintf(w, "  --use-link-paths  Use link's configured paths batch (conflicts with scope and --all)\n")
	fmt.Fprintf(w, "  --all             Allow full-root operations\n")
	fmt.Fprintf(w, "  --dry-run         Preview only (no prompt, no apply)\n")
	fmt.Fprintf(w, "  --yes             Apply without prompting (still runs preview first)\n")
	fmt.Fprintf(w, "  --verbose         Stream full rsync output (future)\n")
	fmt.Fprintf(w, "  --json            Emit NDJSON events on stdout (future)\n")
	fmt.Fprintf(w, "  --dangerous       Override high-risk destination blocklist (future)\n")
}

func printRsyncPlansForLink(cmd string, cfg *config.Config, link config.Link, res scope.Result, stdout, stderr io.Writer) error {
	srcEP, dstEP, err := resolveLinkDirection(cmd, link)
	if err != nil {
		return err
	}

	excludes := append([]string{}, cfg.Global.Excludes...)
	excludes = append(excludes, link.Excludes...)

	useSSH := srcEP.Type == config.EndpointTypeSSH || dstEP.Type == config.EndpointTypeSSH

	for _, sc := range res.Scopes {
		src := appendScope(srcEP.RsyncRoot(), sc)
		dst := appendScope(dstEP.RsyncRoot(), sc)

		preview, apply, err := rsync.BuildPreviewApply(rsync.SyncSpec{
			Source:   src,
			Dest:     dst,
			UseSSH:   useSSH,
			Mirror:   link.Mirror,
			Excludes: excludes,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(stdout, "Rsync plan (%s):\n", scLabel(sc))
		fmt.Fprintf(stdout, "  SRC : %s\n", src)
		fmt.Fprintf(stdout, "  DEST: %s\n", dst)
		fmt.Fprintf(stdout, "  Preview argv:\n")
		fmt.Fprintf(stdout, "    rsync %s\n", strings.Join(preview, " "))
		fmt.Fprintf(stdout, "  Apply argv:\n")
		fmt.Fprintf(stdout, "    rsync %s\n", strings.Join(apply, " "))
	}
	return nil
}

func printRsyncPlansForEndpoint(cmd string, cfg *config.Config, ep config.Endpoint, res scope.Result, cwd string, stdout, stderr io.Writer) error {
	cwdRoot := ensureTrailingSlash(filepath.Clean(cwd))

	excludes := append([]string{}, cfg.Global.Excludes...)

	for _, sc := range res.Scopes {
		var srcRoot string
		var dstRoot string
		switch cmd {
		case "pull":
			srcRoot = ep.RsyncRoot()
			dstRoot = cwdRoot
		case "push":
			srcRoot = cwdRoot
			dstRoot = ep.RsyncRoot()
		default:
			return fmt.Errorf("unknown sync command %q", cmd)
		}

		src := appendScope(srcRoot, sc)
		dst := appendScope(dstRoot, sc)

		useSSH := ep.Type == config.EndpointTypeSSH

		preview, apply, err := rsync.BuildPreviewApply(rsync.SyncSpec{
			Source:   src,
			Dest:     dst,
			UseSSH:   useSSH,
			Mirror:   true,
			Excludes: excludes,
		})
		if err != nil {
			return err
		}

		fmt.Fprintf(stdout, "Rsync plan (%s):\n", scLabel(sc))
		fmt.Fprintf(stdout, "  SRC : %s\n", src)
		fmt.Fprintf(stdout, "  DEST: %s\n", dst)
		fmt.Fprintf(stdout, "  Preview argv:\n")
		fmt.Fprintf(stdout, "    rsync %s\n", strings.Join(preview, " "))
		fmt.Fprintf(stdout, "  Apply argv:\n")
		fmt.Fprintf(stdout, "    rsync %s\n", strings.Join(apply, " "))
	}

	return nil
}

func resolveLinkDirection(cmd string, link config.Link) (src config.Endpoint, dst config.Endpoint, err error) {
	switch cmd {
	case "pull":
		return link.Remote, link.Local, nil
	case "push":
		return link.Local, link.Remote, nil
	default:
		return config.Endpoint{}, config.Endpoint{}, fmt.Errorf("unknown sync command %q", cmd)
	}
}

func appendScope(root, scope string) string {
	root = ensureTrailingSlash(root)
	if scope == "" {
		return root
	}
	// scopes are normalized elsewhere; treat them as slash-separated relative paths.
	if strings.HasSuffix(root, "/") {
		return root + scope + "/"
	}
	return root + "/" + scope + "/"
}

func ensureTrailingSlash(s string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

func scLabel(scope string) string {
	if scope == "" {
		return "full-root"
	}
	return scope
}

func printScopePlan(cmd, mode, name string, res scope.Result, stdout, stderr io.Writer) {
	fmt.Fprintf(stdout, "%s (%s %s)\n", strings.ToUpper(cmd), mode, name)
	fmt.Fprintf(stdout, "Scope source: %s\n", res.Source)
	if res.Batch {
		fmt.Fprintf(stdout, "Scopes: %d\n", len(res.Scopes))
	} else {
		if res.IsFullRoot {
			fmt.Fprintf(stdout, "Scope: <full-root>\n")
		} else {
			fmt.Fprintf(stdout, "Scope: %s\n", res.ScopeString)
		}
	}

	for _, n := range res.Notices {
		w := stderr
		if n.Level == "info" {
			w = stdout
		}
		fmt.Fprintf(w, "%s: %s\n", n.Level, n.Message)
		for _, h := range n.Hints {
			fmt.Fprintf(w, "%s\n", h)
		}
	}
}
