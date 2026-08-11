package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/dist"
)

// The bug this closes: repoRoot accepted any directory with a go.mod beside the
// compile-time source path. Under `go install` that path is the read-only module cache,
// which has a go.mod and would win — and then ~/.claude gets symlinked into a versioned
// cache directory that the next install orphans.
//
// Every operation that matters needs git: the pull, the rebuild decision, the release
// check. So the probe is .git.
func TestPayloadRootRequiresGitDirectory(t *testing.T) {
	// A module cache entry: go.mod, sources, no .git.
	cache := t.TempDir()
	writeFile(t, filepath.Join(cache, "go.mod"), "module github.com/pausf/libretto-automata\n")

	if isRepo(cache) {
		t.Error("a directory with go.mod and no .git was accepted as the repo")
	}

	if err := os.Mkdir(filepath.Join(cache, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !isRepo(cache) {
		t.Error("a directory with .git was rejected")
	}
}

// A .git file rather than a directory is what a worktree has. Refusing it would mean the
// flow's own worktree advice breaks the tool inside the worktree.
func TestPayloadRootAcceptsAWorktreeGitFile(t *testing.T) {
	wt := t.TempDir()
	writeFile(t, filepath.Join(wt, ".git"), "gitdir: /somewhere/.git/worktrees/x\n")

	if !isRepo(wt) {
		t.Error("a worktree's .git file was rejected")
	}
}

// The override wins over every other rung, and it is taken as given with no .git check.
// It is the escape hatch, and a hatch that validates is a hatch that can refuse the one
// case you needed it for.
func TestPayloadRootPrefersEnvOverride(t *testing.T) {
	repo := gitDir(t)
	want := filepath.Join(t.TempDir(), "does-not-exist-yet")

	got, err := resolveRoot(want, repo, repo, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("resolveRoot = %q, want the override %q", got, want)
	}
}

// The real payloadRoot() reads the environment, and the override is the one rung it can be
// steered down from outside. Worth one test, so the wiring between the two is not taken
// on trust.
func TestPayloadRootReadsTheOverrideFromTheEnvironment(t *testing.T) {
	want := t.TempDir()
	t.Setenv(EnvRoot, want)

	got, err := payloadRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("repoRoot = %q, want %q", got, want)
	}
}

// With no override and no checkout anywhere, the answer is the activated release — not the
// working directory, which is what the old fallback returned and would have linked
// ~/.claude against whatever the user happened to be cd'd into.
func TestPayloadRootFallsBackToTheActivatedRelease(t *testing.T) {
	home := t.TempDir()

	got, err := resolveRoot("", t.TempDir(), t.TempDir(), home)
	if err != nil {
		t.Fatal(err)
	}
	if want := dist.Current(dist.Base(home)); got != want {
		t.Errorf("resolveRoot = %q, want the activated release %q", got, want)
	}
	// And emphatically not the directory the clone bootstrap used, which this replaced.
	if strings.Contains(got, ".libretto-automata") {
		t.Errorf("resolveRoot still points at the removed clone directory: %q", got)
	}
}

// A checkout you are standing in still wins over the activated release. That is what keeps
// editing a skill and seeing it live working, and it is the reason the payload is not
// embedded in the binary.
func TestPayloadRootStillPrefersACheckoutYouAreStandingIn(t *testing.T) {
	home := t.TempDir()
	wd := gitDir(t)
	writeFile(t, filepath.Join(wd, "go.mod"), moduleLine+"\n")

	got, err := resolveRoot("", t.TempDir(), wd, home)
	if err != nil {
		t.Fatal(err)
	}
	if got != wd {
		t.Errorf("resolveRoot = %q, want the checkout %q", got, wd)
	}
}

// Rescued from the deleted bootstrap_test.go. The behaviour it guarded is gone — nothing
// clones any more — but the promise underneath it is not: `version` and `help` answer
// without the payload, so neither may write anything anywhere.
func TestVersionAndHelpTouchNothing(t *testing.T) {
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

// The compile-time path is rung 2 and beats the working directory — that is development,
// where `make build` puts the binary under the clone it was built from.
func TestPayloadRootPrefersTheCompileTimePathOverTheWorkingDirectory(t *testing.T) {
	compileTime, wd := gitDir(t), gitDir(t)
	writeFile(t, filepath.Join(wd, "go.mod"), moduleLine+"\n")

	got, err := resolveRoot("", compileTime, wd, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if got != compileTime {
		t.Errorf("resolveRoot = %q, want the compile-time path %q", got, compileTime)
	}
}

// The working directory is only accepted when it is genuinely a clone of this module.
// Any git repository would otherwise do, and `libretto install` run inside an unrelated
// project would go looking for a payload that project does not have.
func TestPayloadRootAcceptsTheWorkingDirectoryOnlyForThisModule(t *testing.T) {
	home := t.TempDir()

	mine := gitDir(t)
	writeFile(t, filepath.Join(mine, "go.mod"), moduleLine+"\n")
	got, err := resolveRoot("", t.TempDir(), mine, home)
	if err != nil {
		t.Fatal(err)
	}
	if got != mine {
		t.Errorf("resolveRoot = %q, want this module's clone %q", got, mine)
	}

	other := gitDir(t)
	writeFile(t, filepath.Join(other, "go.mod"), "module example.com/something-else\n")
	got, err = resolveRoot("", t.TempDir(), other, home)
	if err != nil {
		t.Fatal(err)
	}
	if got == other {
		t.Error("an unrelated git repository was accepted as this tool's clone")
	}
}

func gitDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
