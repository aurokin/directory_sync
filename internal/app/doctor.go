package app

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/aurokin/directory_sync/internal/config"
)

func runDoctor(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(stderr)

	endpointName := fs.String("endpoint", "", "validate a single endpoint")
	linkName := fs.String("link", "", "validate a single link")
	_ = fs.Bool("json", false, "emit NDJSON on stdout (not implemented yet)")

	if err := fs.Parse(args); err != nil {
		return 2
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
		fmt.Fprintf(stderr, "doctor: %v\n", err)
		return 1
	}

	if *endpointName != "" {
		if _, ok := cfg.Endpoints[*endpointName]; !ok {
			fmt.Fprintf(stderr, "doctor: unknown endpoint %q\n", *endpointName)
			return 1
		}
	}
	if *linkName != "" {
		if _, ok := cfg.Links[*linkName]; !ok {
			fmt.Fprintf(stderr, "doctor: unknown link %q\n", *linkName)
			return 1
		}
	}

	fmt.Fprintf(stdout, "Config: %s\n", cfg.FilePath)
	fmt.Fprintf(stdout, "Endpoints: %d\n", len(cfg.Endpoints))
	fmt.Fprintf(stdout, "Links: %d\n", len(cfg.Links))
	fmt.Fprintf(stdout, "OK\n")
	return 0
}
