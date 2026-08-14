package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// The fixtures every test in this package builds on.
//
// Two temporary directories and two environment variables. `CLAUDE_HOME` is what
// keeps the whole suite away from the author's real ~/.claude, and it is the
// reason these tests can run `install` and `prune` — commands that write and
// delete — without a scratch VM.

// fixture is a repo and a target, both temporary, both empty.
type fixture struct {
	Repo     string
	Claude   string
	Project  string
	Codex    string
	Opencode string
}

func newFixture(t *testing.T) fixture {
	t.Helper()

	f := fixture{
		Repo: t.TempDir(), Claude: t.TempDir(), Project: t.TempDir(),
		Codex: t.TempDir(), Opencode: t.TempDir(),
	}
	t.Setenv(claudeHomeEnv, f.Claude)
	t.Setenv(target.EnvAgentsHome, f.Codex)
	t.Setenv(target.EnvOpencodeHome, f.Opencode)
	return f
}

// global is the scope every existing test means: ~/.claude, redirected by
// CLAUDE_HOME to a temporary directory.
func (f fixture) global() target.Target { return target.Global() }

// project is the other scope, rooted at <Project>/.claude. Kept separate from
// Claude so a test can assert that writing one leaves the other untouched — the
// whole point of having scopes.
func (f fixture) project() target.Target { return target.Project(f.Project) }

// projectDest is a path inside the project scope's .claude.
func (f fixture) projectDest(parts ...string) string {
	return filepath.Join(append([]string{f.Project, ".claude"}, parts...)...)
}

// link plants a symlink, creating its parent. Used to set up the states a scan
// classifies — a stale link, or one aimed at the wrong item.
func (f fixture) link(t *testing.T, from, to string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(from, to); err != nil {
		t.Fatal(err)
	}
}

// claudeHomeEnv is target.EnvClaudeHome, named here so the tests read as tests
// rather than as a tour of the packages they exercise.
const claudeHomeEnv = "CLAUDE_HOME"

// skill writes a skill into the repo and returns its path.
func (f fixture) skill(t *testing.T, name string) string {
	t.Helper()

	dir := filepath.Join(f.Repo, "skills", name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: " + name + "\ndescription: fixture\n---\n"
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// command writes a command into the repo and returns its path. A file, not a
// directory — commands are single .md files, which is what target.Kind already says.
func (f fixture) command(t *testing.T, name string) string {
	t.Helper()

	dir := filepath.Join(f.Repo, "commands")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name+".md")
	body := "---\ndescription: fixture\n---\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// dest is the path a repo item of this kind and name would be linked to.
func (f fixture) dest(kind, name string) string {
	return filepath.Join(f.Claude, kind, name)
}

// putReal drops a real file where a link would go, so a test can build the one
// situation the tool must never touch.
func (f fixture) putReal(t *testing.T, kind, name, content string) string {
	t.Helper()

	path := f.dest(kind, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// capture runs fn with stdout and stderr redirected, and returns what was written
// to each.
//
// A pipe has a buffer and a test that fills it deadlocks, so the reading happens
// in goroutines rather than after fn returns. None of these commands write
// anywhere near that much, but a test suite that hangs on a future verbose
// command is a bad afternoon.
func capture(t *testing.T, fn func() error) (stdout, stderr string, err error) {
	t.Helper()

	outR, outW, e := os.Pipe()
	if e != nil {
		t.Fatal(e)
	}
	errR, errW, e := os.Pipe()
	if e != nil {
		t.Fatal(e)
	}

	origOut, origErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW

	outCh, errCh := make(chan string, 1), make(chan string, 1)
	go func() { outCh <- readAll(outR) }()
	go func() { errCh <- readAll(errR) }()

	err = fn()

	os.Stdout, os.Stderr = origOut, origErr
	outW.Close()
	errW.Close()

	return <-outCh, <-errCh, err
}

func readAll(f *os.File) string {
	var b strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		b.Write(buf[:n])
		if err != nil {
			break
		}
	}
	f.Close()
	return b.String()
}

// isSymlinkTo reports whether path is a symlink resolving to want.
func isSymlinkTo(t *testing.T, path, want string) bool {
	t.Helper()

	fi, err := os.Lstat(path)
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		return false
	}
	got, err := os.Readlink(path)
	if err != nil {
		return false
	}
	gotResolved, _ := filepath.EvalSymlinks(got)
	wantResolved, _ := filepath.EvalSymlinks(want)
	return gotResolved == wantResolved
}

func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// rowOf finds a menu row by label.
//
// Tests used to hardcode the index, and adding `uninstall` to the menu broke one of
// them — a test keyed to a position fails when the position changes, which says
// nothing about the behaviour it was meant to protect.
func rowOf(t *testing.T, menu []ui.MenuItem, label string) int {
	t.Helper()
	for i, m := range menu {
		if m.Label == label {
			return i
		}
	}
	t.Fatalf("the menu has no %q row", label)
	return 0
}
