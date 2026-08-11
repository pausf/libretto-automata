package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSubcommandNoticeGoesToStderr(t *testing.T) {
	stubNotice(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	var err bytes.Buffer
	noticeAfter(&err, "status", "/root", "v0.1.0")

	if !strings.Contains(err.String(), "v0.2.0 available") {
		t.Fatalf("stderr = %q, want the notice", err.String())
	}
	if !strings.HasSuffix(err.String(), "\n") {
		t.Fatalf("stderr = %q, want a trailing newline", err.String())
	}
}

func TestSubcommandNoticeIsSilentWithNothingToSay(t *testing.T) {
	stubNotice(t, func(root, running string) string { return "" })

	var err bytes.Buffer
	noticeAfter(&err, "status", "/root", "v0.1.0")

	if err.String() != "" {
		t.Fatalf("stderr = %q, want nothing", err.String())
	}
}

func TestDoctorAndUpdateDoNotRepeatTheNotice(t *testing.T) {
	stubNotice(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	for _, cmd := range []string{"doctor", "update", "version", "help", "unknown"} {
		var err bytes.Buffer
		noticeAfter(&err, cmd, "/root", "v0.1.0")
		if err.String() != "" {
			t.Errorf("%s printed %q, want nothing", cmd, err.String())
		}
	}
}

func TestSubcommandNoticeHasNoEscapeCodes(t *testing.T) {
	stubNotice(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })

	for _, cmd := range []string{"status", "preview", "install", "prune", "uninstall", "models"} {
		var err bytes.Buffer
		noticeAfter(&err, cmd, "/root", "v0.1.0")
		if err.Len() == 0 {
			t.Errorf("%s printed nothing, want the notice", cmd)
		}
		if strings.Contains(err.String(), "\x1b") {
			t.Errorf("%s printed an escape code: %q", cmd, err.String())
		}
	}
}

func TestSubcommandNoticeDoesNotChangeTheExitCode(t *testing.T) {
	t.Setenv("CLAUDE_HOME", t.TempDir())

	var asked []string
	stubNotice(t, func(root, running string) string {
		asked = append(asked, running)
		return "v0.1.0 → v0.2.0 available"
	})

	if err := run([]string{"status"}); err != nil {
		t.Fatalf("status returned %v, want nil", err)
	}
	// Without this the test would pass on a build where the notice is never wired in
	// at all — green for the wrong reason, which proves nothing about the exit code.
	if len(asked) != 1 {
		t.Fatalf("the notice ran %d times on `status`, want 1", len(asked))
	}

	err := run([]string{"nonesuch"})
	if err == nil || !strings.Contains(err.Error(), "nonesuch") {
		t.Fatalf("unknown command returned %v, want its own error", err)
	}
	if len(asked) != 1 {
		t.Fatalf("the notice ran on an unknown command")
	}
}

// stubNotice swaps the notice source for the duration of one test. A package var is
// the smallest seam that keeps these tests off the network.
func stubNotice(t *testing.T, f func(root, running string) string) {
	t.Helper()
	prev := noticeSource
	noticeSource = f
	t.Cleanup(func() { noticeSource = prev })
}
