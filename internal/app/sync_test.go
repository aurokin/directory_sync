package app

import (
	"bytes"
	"testing"
)

func TestParseSyncArgsAllowsInterspersedFlags(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	flags, pos, code, printedHelp := parseSyncArgs("pull", []string{"photos", "--dry-run"}, &stdout, &stderr)
	if printedHelp {
		t.Fatalf("unexpected help")
	}
	if code != 0 {
		t.Fatalf("unexpected code %d stderr=%q", code, stderr.String())
	}
	if !flags.dryRun {
		t.Fatalf("expected dryRun true")
	}
	if len(pos) != 1 || pos[0] != "photos" {
		t.Fatalf("unexpected positionals: %v", pos)
	}
}

func TestParseSyncArgsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	_, _, code, printedHelp := parseSyncArgs("pull", []string{"photos", "--nope"}, &stdout, &stderr)
	if printedHelp {
		t.Fatalf("unexpected help")
	}
	if code == 0 {
		t.Fatalf("expected non-zero code")
	}
}
