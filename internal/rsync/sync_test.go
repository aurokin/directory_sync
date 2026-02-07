package rsync

import "testing"

func TestBuildPreviewApplyDiffersOnlyByDryRun(t *testing.T) {
	spec := SyncSpec{
		Source:   "/tmp/src/",
		Dest:     "/tmp/dst/",
		UseSSH:   false,
		Mirror:   true,
		Excludes: []string{".DS_Store", ".git/"},
	}

	preview, apply, err := BuildPreviewApply(spec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(preview) != len(apply)+1 {
		t.Fatalf("expected preview to have exactly one extra arg (dry-run)")
	}

	applyFromPreview := removeFirst(preview, "--dry-run")
	if len(applyFromPreview) != len(apply) {
		t.Fatalf("unexpected removal result")
	}
	for i := range apply {
		if apply[i] != applyFromPreview[i] {
			t.Fatalf("args differ at %d: want %q got %q", i, apply[i], applyFromPreview[i])
		}
	}
}

func TestBuildArgsIncludesSSHTransportWhenEnabled(t *testing.T) {
	args, err := BuildArgs(SyncSpec{
		Source: "/tmp/src/",
		Dest:   "photo-box:/srv/photos/",
		UseSSH: true,
		Mirror: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsSubseq(args, []string{"-e", "ssh"}) {
		t.Fatalf("expected -e ssh in args: %v", args)
	}
}

func TestBuildArgsDoesNotEnableCompressionByDefault(t *testing.T) {
	args, err := BuildArgs(SyncSpec{
		Source: "/tmp/src/",
		Dest:   "photo-box:/srv/photos/",
		UseSSH: true,
		Mirror: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, a := range args {
		if a == "-z" || a == "--compress" {
			t.Fatalf("unexpected compression flag present")
		}
	}
}

func removeFirst(args []string, needle string) []string {
	out := make([]string, 0, len(args))
	removed := false
	for _, a := range args {
		if !removed && a == needle {
			removed = true
			continue
		}
		out = append(out, a)
	}
	return out
}

func containsSubseq(haystack []string, subseq []string) bool {
	if len(subseq) == 0 {
		return true
	}
	for i := 0; i+len(subseq) <= len(haystack); i++ {
		ok := true
		for j := range subseq {
			if haystack[i+j] != subseq[j] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
