package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// The hole the bootstrapper opens: once the binary can live in $GOBIN, rebuilding into
// the clone's bin/ upgrades a file nobody runs. `update` would report success and every
// later invocation would stay on the old version.
func TestRebuildReplacesRunningExecutable(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles the CLI")
	}
	root := repoUnderTest(t)
	binBefore := statOrZero(filepath.Join(root, "bin", "libretto"))
	exe := filepath.Join(t.TempDir(), "bin", "libretto")
	if err := os.MkdirAll(filepath.Dir(exe), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(exe, []byte("stale"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := rebuild(root, exe); err != nil {
		t.Fatalf("rebuild: %v", err)
	}

	body, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) == "stale" {
		t.Error("rebuild did not replace the given executable")
	}
	if info, err := os.Stat(exe); err != nil {
		t.Fatal(err)
	} else if info.Mode().Perm()&0o100 == 0 {
		t.Errorf("the rebuilt binary is not executable: %v", info.Mode())
	}

	// And the clone's bin/ was left alone, which is where it used to go. Compared by
	// modification time rather than existence: a development checkout has usually run
	// `make build`, so asserting the file is absent would fail on a clean repository for
	// the wrong reason.
	if after := statOrZero(filepath.Join(root, "bin", "libretto")); after != binBefore {
		t.Error("rebuild also wrote bin/libretto in the clone")
	}
}

// `make link` puts a symlink in ~/.local/bin pointing at bin/libretto. The rebuild has to
// write the file the link resolves to, not replace the link with a regular file — that
// would silently sever the development setup.
func TestRebuildResolvesSymlinkedExecutable(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles the CLI")
	}
	root := repoUnderTest(t)

	real := filepath.Join(t.TempDir(), "libretto")
	if err := os.WriteFile(real, []byte("stale"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "libretto")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}

	if err := rebuild(root, link); err != nil {
		t.Fatalf("rebuild: %v", err)
	}

	if got, err := os.Lstat(link); err != nil {
		t.Fatal(err)
	} else if got.Mode()&os.ModeSymlink == 0 {
		t.Error("rebuild replaced the symlink with a regular file")
	}
	if body, err := os.ReadFile(real); err != nil {
		t.Fatal(err)
	} else if string(body) == "stale" {
		t.Error("rebuild did not write through to the link's target")
	}
}

// A binary somewhere the user cannot write is a report, not a failure. The pull happened
// and the links are correct; rolling those back because a rename failed loses more than
// it saves. The message has to say where the new binary is, or it is a dead end.
func TestRebuildReportsUnwritableDestinationWithoutFailing(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles the CLI")
	}
	if os.Geteuid() == 0 {
		t.Skip("root can write anywhere")
	}
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions do not work this way")
	}
	root := repoUnderTest(t)

	locked := filepath.Join(t.TempDir(), "locked")
	if err := os.Mkdir(locked, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o700) })

	refused := filepath.Join(locked, "libretto")
	fallback := filepath.Join(t.TempDir(), "bin", "libretto")
	note, err := rebuildOrReport(root, refused, fallback)
	if err != nil {
		t.Fatalf("an unwritable destination failed the update: %v", err)
	}
	if note == "" {
		t.Fatal("no note explaining that the binary was not replaced")
	}

	// Both halves, checked separately. `Contains(note, "libretto")` alone is satisfied by
	// the refused path — every path in this test has `libretto` in it — so it passed
	// without ever proving the note says where the new binary went.
	if !strings.Contains(note, refused) {
		t.Errorf("the note does not name what could not be written: %q", note)
	}
	if !strings.Contains(note, fallback) {
		t.Errorf("the note does not say where the new binary is (%s): %q", fallback, note)
	}
	if _, err := os.Stat(fallback); err != nil {
		t.Errorf("the note names a binary that is not there: %v", err)
	}
}

// statOrZero is the file's modification time, or the zero time when it is not there. Both
// answers are legitimate for bin/libretto: whether it exists depends on whether anyone has
// run `make build` in this checkout.
func statOrZero(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

// repoUnderTest is a clone this test can compile: the real one. The rebuild runs `go
// build ./cmd/libretto`, so it needs actual sources, and copying them would test the copy.
func repoUnderTest(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate the source tree")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(file)))
}
