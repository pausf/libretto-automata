package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestMain silences the release check for the whole package.
//
// Without it every test that reaches `run` inherits the real one, which asks the module
// proxy and writes `.update-check` into the user's module cache — a live network call and
// a write to a real directory from `go test`, which is the thing CLAUDE_HOME exists to
// prevent for links. A test that wants the notice opts in with stubNotice.
func TestMain(m *testing.M) {
	noticeSource = func(root, running string) string { return "" }
	os.Exit(m.Run())
}

// The stream matters as much as the line: `status` output is parseable, and an existing
// promise says piping it carries no escape codes. So this runs a real command through the
// dispatch and compares stdout against the same command with nothing to say.
func TestSubcommandNoticeGoesToStderr(t *testing.T) {
	t.Setenv("CLAUDE_HOME", t.TempDir())

	quiet, _ := captureStd(t, func() { _ = run([]string{"status"}) })

	stubNotice(t, func(root, running string) string { return "v0.1.0 → v0.2.0 available" })
	out, errOut := captureStd(t, func() { _ = run([]string{"status"}) })

	if out != quiet {
		t.Errorf("stdout changed when the notice fired:\n with = %q\nwithout = %q", out, quiet)
	}
	if !strings.Contains(errOut, "v0.2.0 available") {
		t.Errorf("stderr = %q, want the notice", errOut)
	}
	if !strings.HasSuffix(errOut, "\n") {
		t.Errorf("stderr = %q, want a trailing newline", errOut)
	}
}

// The silence conditions live in subcommandNotice, so they are tested there rather than
// through a stub that can only prove noticeAfter prints nothing when handed "".
func TestSubcommandNoticeIsSilentWithNothingToSay(t *testing.T) {
	cases := map[string]struct{ cached, running string }{
		"up to date":          {"v0.3.0", "v0.3.0"},
		"ahead of the remote": {"v0.2.0", "v0.3.0"},
		"unidentifiable":      {"v0.3.0", "dev"},
		"nothing cached":      {"", "v0.2.0"},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := subcommandNotice(rootWithCachedTag(t, c.cached), c.running); got != "" {
				t.Errorf("notice = %q, want silence", got)
			}
		})
	}
}

func TestSubcommandNoticeDoesNotChangeTheExitCode(t *testing.T) {
	t.Setenv("CLAUDE_HOME", t.TempDir())

	var asked int
	stubNotice(t, func(root, running string) string {
		asked++
		return "v0.1.0 → v0.2.0 available"
	})

	// A command that succeeds still succeeds, and the notice really fired — without this
	// the test would pass on a build where nothing is wired in at all.
	_, errOut := captureStd(t, func() {
		if err := run([]string{"status"}); err != nil {
			t.Errorf("status returned %v, want nil", err)
		}
	})
	if asked != 1 {
		t.Fatalf("the notice ran %d times on `status`, want 1", asked)
	}
	if !strings.Contains(errOut, "available") {
		t.Fatalf("stderr = %q, want the notice", errOut)
	}

	// A failing command that carries the notice keeps its own error, and still says it.
	var err error
	_, errOut = captureStd(t, func() { err = run([]string{"models", "nonesuch"}) })
	if err == nil || !strings.Contains(err.Error(), "nonesuch") {
		t.Fatalf("`models nonesuch` returned %v, want its own error", err)
	}
	if asked != 2 {
		t.Fatalf("the notice ran %d times, want it on the failing command too", asked)
	}
	if !strings.Contains(errOut, "available") {
		t.Fatalf("stderr = %q, want the notice on the failing path too", errOut)
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

// Against the real formatter, not a string the test supplied — colour introduced in
// subcommandNotice is exactly what this has to catch.
func TestSubcommandNoticeHasNoEscapeCodes(t *testing.T) {
	// The real formatter, put back for this test — TestMain silences it package-wide.
	stubNotice(t, subcommandNotice)

	root := rootWithCachedTag(t, "v0.3.0")
	line := subcommandNotice(root, "v0.2.0")
	if line == "" {
		t.Fatal("the notice is empty, want a line to check")
	}
	if strings.Contains(line, "\x1b") {
		t.Fatalf("the notice carries an escape code: %q", line)
	}

	for _, cmd := range []string{"status", "preview", "install", "prune", "uninstall", "models"} {
		var err bytes.Buffer
		noticeAfter(&err, cmd, root, "v0.2.0")
		if err.Len() == 0 {
			t.Errorf("%s printed nothing, want the notice", cmd)
		}
		if strings.Contains(err.String(), "\x1b") {
			t.Errorf("%s printed an escape code: %q", cmd, err.String())
		}
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

// captureStd runs fn with both standard streams on pipes and returns what each got.
//
// Drained concurrently: a report longer than the pipe buffer would block the writer
// forever, which is the same reason runCaptured drains in a goroutine.
func captureStd(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	prevOut, prevErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW

	outC, errC := make(chan string, 1), make(chan string, 1)
	go func() { b, _ := io.ReadAll(outR); outC <- string(b) }()
	go func() { b, _ := io.ReadAll(errR); errC <- string(b) }()

	func() {
		defer func() {
			os.Stdout, os.Stderr = prevOut, prevErr
			outW.Close()
			errW.Close()
		}()
		fn()
	}()

	return <-outC, <-errC
}
