package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/repo"
)

// The destination is printed before the clone runs, not after. A tool that has already
// written a directory into someone's home and then mentions it is a tool they cannot
// decline.
func TestBootstrapAnnouncesDestinationBeforeCloning(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "libretto-automata")

	var log strings.Builder
	cloned := false
	clone := func(context.Context, string, string) error {
		if !strings.Contains(log.String(), dest) {
			t.Errorf("the clone started before the destination was announced; log so far: %q",
				log.String())
		}
		cloned = true
		return nil
	}

	if err := bootstrap(context.Background(), &log, dest, "https://example.invalid/x.git", clone); err != nil {
		t.Fatal(err)
	}
	if !cloned {
		t.Fatal("bootstrap did not clone")
	}
	if !strings.Contains(log.String(), dest) {
		t.Errorf("the destination was never announced: %q", log.String())
	}
}

// A destination that exists and is not ours is refused and nothing is touched — the same
// promise the linker makes, applied to the tool's own directory.
func TestBootstrapRefusesForeignDestination(t *testing.T) {
	dest := t.TempDir()
	stranger := filepath.Join(dest, "not-ours")
	if err := os.WriteFile(stranger, []byte("keep me"), 0o644); err != nil {
		t.Fatal(err)
	}

	var log strings.Builder
	clone := func(context.Context, string, string) error {
		t.Error("bootstrap cloned over a directory it did not create")
		return nil
	}

	err := bootstrap(context.Background(), &log, dest, "https://example.invalid/x.git", clone)
	if err == nil {
		t.Fatal("bootstrap accepted a foreign destination")
	}
	if !strings.Contains(err.Error(), dest) {
		t.Errorf("the refusal does not name the destination: %v", err)
	}
	if body, rerr := os.ReadFile(stranger); rerr != nil || string(body) != "keep me" {
		t.Error("the refused bootstrap touched the existing file")
	}
}

// A destination that is already our clone is not an error and not a re-clone. This is the
// ordinary second run, and it has to be silent.
func TestBootstrapIsIdempotent(t *testing.T) {
	dest := t.TempDir()
	if err := os.Mkdir(filepath.Join(dest, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	var log strings.Builder
	clone := func(context.Context, string, string) error {
		t.Error("bootstrap re-cloned over an existing clone")
		return nil
	}

	if err := bootstrap(context.Background(), &log, dest, "https://example.invalid/x.git", clone); err != nil {
		t.Fatalf("bootstrap over an existing clone: %v", err)
	}
	if log.String() != "" {
		t.Errorf("bootstrap narrated a no-op: %q", log.String())
	}
}

// Bootstrap is a step on the way to the command the user typed, not a command of its own.
// `libretto status` on a machine with no clone clones and then reports status — one
// invocation, no "now run it again".
func TestBootstrapContinuesIntoRequestedCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("real-git integration")
	}

	// A clonable source: a repository with one commit and the directories the payload
	// scan expects.
	upstream := t.TempDir()
	for _, args := range [][]string{
		{"init", "-q", "-b", "main"},
		{"config", "user.email", "test@example.invalid"},
		{"config", "user.name", "Test"},
	} {
		out, err := exec.Command("git", append([]string{"-C", upstream}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	writeFile(t, filepath.Join(upstream, "go.mod"), moduleLine+"\n")
	for _, args := range [][]string{{"add", "-A"}, {"commit", "-q", "-m", "payload"}} {
		out, err := exec.Command("git", append([]string{"-C", upstream}, args...)...).CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	dest := filepath.Join(t.TempDir(), "libretto-automata")
	t.Setenv(EnvRoot, dest)
	t.Setenv("CLAUDE_HOME", t.TempDir())

	var log strings.Builder
	if err := bootstrap(context.Background(), &log, dest, upstream, repo.Clone); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	// ensureClone is what commands call. After the clone above it must find the repo and
	// return it, rather than announcing or cloning a second time.
	root, err := ensureClone()
	if err != nil {
		t.Fatalf("ensureClone after bootstrap: %v", err)
	}
	if root != dest {
		t.Errorf("ensureClone = %q, want the clone it just made at %q", root, dest)
	}

	// And the requested command runs against it. `status` is the read-only one, so this
	// proves the hand-off without writing a link anywhere.
	if err := run([]string{"status"}); err != nil {
		t.Errorf("the requested command did not run after bootstrap: %v", err)
	}
}

// `version` and `help` answer without the payload, so they must not go looking for a
// clone — never mind creating one. Cloning a repository into somebody's home because they
// asked what version they were running would be indefensible.
func TestVersionAndHelpDoNotBootstrap(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(EnvRoot, filepath.Join(home, "nothing-here"))

	for _, arg := range []string{"version", "--version", "help", "--help"} {
		if err := run([]string{arg}); err != nil {
			t.Errorf("run(%q): %v", arg, err)
		}
	}

	entries, err := os.ReadDir(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("version/help wrote into the home directory: %v", entries)
	}
}

// A failed clone leaves nothing behind that a later run would refuse. Half a clone in
// ~/.libretto-automata would be a foreign destination forever, and the user would have to
// work out that deleting it is the fix.
func TestBootstrapCleansUpAfterAFailedClone(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "libretto-automata")

	var log strings.Builder
	clone := func(_ context.Context, _, dst string) error {
		if err := os.MkdirAll(filepath.Join(dst, "partial"), 0o755); err != nil {
			t.Fatal(err)
		}
		return os.ErrDeadlineExceeded
	}

	if err := bootstrap(context.Background(), &log, dest, "https://example.invalid/x.git", clone); err == nil {
		t.Fatal("a failed clone was reported as success")
	}
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Errorf("a failed clone left %s behind", dest)
	}
}
