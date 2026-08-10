package repo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestNeedsRebuild(t *testing.T) {
	cases := []struct {
		name  string
		paths []string
		want  bool
	}{
		{"nothing changed", nil, false},
		{"only payload", []string{"skills/write-spec/SKILL.md", "docs/FLOW.md"}, false},
		{"a go file", []string{"internal/link/plan.go"}, true},
		{"go file among markdown", []string{"README.md", "cmd/libretto/main.go"}, true},
		{"go.mod", []string{"go.mod"}, true},
		{"go.sum", []string{"go.sum"}, true},
		// A path that merely mentions go must not trigger a compile.
		{"markdown about go", []string{"docs/going-further.md"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NeedsRebuild(c.paths); got != c.want {
				t.Errorf("NeedsRebuild(%v) = %v, want %v", c.paths, got, c.want)
			}
		})
	}
}

// The git-backed tests below are the "real-git integration" make test-short has
// promised to skip since before this change — a promise that was empty, because no
// test in the repository called testing.Short(). They honour it now.

// gitRepo builds a real repository in a temp dir and returns its root.
//
// Real git rather than a fake: Shell exists precisely so the invocation lives in one
// place, and replacing it in tests would prove the fake works. The cost is a slower
// test that needs git on the machine, which is what -short is for.
func gitRepo(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("real-git integration")
	}

	root := t.TempDir()
	run(t, root, "init", "-q", "-b", "main")
	// A runner has no global git identity, and `git commit` fails without one — as a
	// failure that reads like the code broke rather than the environment.
	run(t, root, "config", "user.email", "test@example.invalid")
	run(t, root, "config", "user.name", "Test")
	return root
}

func run(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commit(t *testing.T, root, message string) {
	t.Helper()
	run(t, root, "add", "-A")
	run(t, root, "commit", "-q", "-m", message)
}

func TestDirtyReportsTheWorkingTree(t *testing.T) {
	root := gitRepo(t)
	write(t, root, "a.txt", "one\n")
	commit(t, root, "first")

	s := Shell{Root: root}

	if dirty, err := s.Dirty(); err != nil || dirty {
		t.Errorf("Dirty() = %v, %v on a clean tree; want false, nil", dirty, err)
	}

	write(t, root, "a.txt", "two\n")
	if dirty, err := s.Dirty(); err != nil || !dirty {
		t.Errorf("Dirty() = %v, %v with an uncommitted change; want true, nil", dirty, err)
	}
}

// No remote is a state, not a failure — `update` has to tell the two apart to decide
// whether pulling is even possible.
func TestHasRemoteDistinguishesNoRemoteFromAnError(t *testing.T) {
	root := gitRepo(t)
	s := Shell{Root: root}

	has, err := s.HasRemote()
	if err != nil {
		t.Fatalf("HasRemote() errored on a repository with no remote: %v", err)
	}
	if has {
		t.Error("HasRemote() = true before any remote was added")
	}

	run(t, root, "remote", "add", "origin", "https://example.invalid/x.git")
	if has, err := s.HasRemote(); err != nil || !has {
		t.Errorf("HasRemote() = %v, %v after adding one; want true, nil", has, err)
	}
}

func TestHeadIsTheCurrentCommit(t *testing.T) {
	root := gitRepo(t)
	write(t, root, "a.txt", "one\n")
	commit(t, root, "first")

	want := run(t, root, "rev-parse", "HEAD")
	got, err := Shell{Root: root}.Head()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Head() = %q, want %q", got, want)
	}
}

// Head is the one call that answers "" with no error when git fails, deliberately:
// a repository with no commits yet has no HEAD, and that is a state rather than a
// failure. Pinning it because a future tidy-up would otherwise "fix" it into an error
// and break the empty-repository path silently.
func TestHeadOnARepositoryWithNoCommitsIsEmptyNotAnError(t *testing.T) {
	got, err := Shell{Root: gitRepo(t)}.Head()
	if err != nil {
		t.Errorf("Head() errored on a repository with no commits: %v", err)
	}
	if got != "" {
		t.Errorf("Head() = %q, want empty", got)
	}
}

func TestChangedSinceNamesOnlyWhatChanged(t *testing.T) {
	root := gitRepo(t)
	write(t, root, "kept.txt", "one\n")
	write(t, root, "changed.txt", "one\n")
	commit(t, root, "first")
	before := run(t, root, "rev-parse", "HEAD")

	write(t, root, "changed.txt", "two\n")
	write(t, root, "added.txt", "new\n")
	commit(t, root, "second")

	got, err := Shell{Root: root}.ChangedSince(before)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{"changed.txt": true, "added.txt": true}
	if len(got) != len(want) {
		t.Fatalf("ChangedSince() = %v, want exactly %v", got, want)
	}
	for _, p := range got {
		if !want[p] {
			t.Errorf("ChangedSince() named %q, which did not change", p)
		}
	}
}

// An empty rev means "no previous state to compare against", which is not an error
// and must not become a diff against everything.
func TestChangedSinceWithNoRevisionIsEmpty(t *testing.T) {
	root := gitRepo(t)
	write(t, root, "a.txt", "one\n")
	commit(t, root, "first")

	got, err := Shell{Root: root}.ChangedSince("")
	if err != nil || got != nil {
		t.Errorf("ChangedSince(\"\") = %v, %v; want nil, nil", got, err)
	}
}

// Pull against a bare repository on disk rather than a network remote.
//
// A test that reaches the internet is a test that fails on a train, and this one has
// no reason to: `git` treats a local path as a remote like any other, so the real
// fast-forward path runs with nothing to flake on.
func TestPullFetchesFromALocalRemote(t *testing.T) {
	origin := gitRepo(t)
	write(t, origin, "a.txt", "one\n")
	commit(t, origin, "first")

	bare := t.TempDir()
	run(t, bare, "init", "-q", "--bare", "-b", "main")
	run(t, origin, "remote", "add", "origin", bare)
	run(t, origin, "push", "-q", "origin", "main")

	clone := t.TempDir()
	if out, err := exec.Command("git", "clone", "-q", bare, clone).CombinedOutput(); err != nil {
		t.Fatalf("clone: %v\n%s", err, out)
	}
	run(t, clone, "config", "user.email", "test@example.invalid")
	run(t, clone, "config", "user.name", "Test")

	// A commit that exists upstream and not yet in the clone.
	write(t, origin, "b.txt", "two\n")
	commit(t, origin, "second")
	run(t, origin, "push", "-q", "origin", "main")

	before, err := Shell{Root: clone}.Head()
	if err != nil {
		t.Fatal(err)
	}
	if err := (Shell{Root: clone}).Pull(); err != nil {
		t.Fatalf("Pull() = %v", err)
	}
	after, err := Shell{Root: clone}.Head()
	if err != nil {
		t.Fatal(err)
	}

	if after == before {
		t.Error("Pull() left HEAD where it was — the upstream commit did not arrive")
	}
	if _, err := os.Stat(filepath.Join(clone, "b.txt")); err != nil {
		t.Errorf("the pulled commit's file is not in the working tree: %v", err)
	}
}

// Pull refuses anything that is not a fast-forward, deliberately: a merge commit
// created by a background command is a merge commit nobody chose. Diverged histories
// are the user's to resolve, in their own shell, seeing what they are doing.
func TestPullRefusesADivergedHistory(t *testing.T) {
	origin := gitRepo(t)
	write(t, origin, "a.txt", "one\n")
	commit(t, origin, "first")

	bare := t.TempDir()
	run(t, bare, "init", "-q", "--bare", "-b", "main")
	run(t, origin, "remote", "add", "origin", bare)
	run(t, origin, "push", "-q", "origin", "main")

	clone := t.TempDir()
	if out, err := exec.Command("git", "clone", "-q", bare, clone).CombinedOutput(); err != nil {
		t.Fatalf("clone: %v\n%s", err, out)
	}
	run(t, clone, "config", "user.email", "test@example.invalid")
	run(t, clone, "config", "user.name", "Test")

	// Both sides commit. Neither is an ancestor of the other.
	write(t, origin, "b.txt", "upstream\n")
	commit(t, origin, "upstream work")
	run(t, origin, "push", "-q", "origin", "main")

	write(t, clone, "c.txt", "local\n")
	commit(t, clone, "local work")

	if err := (Shell{Root: clone}).Pull(); err == nil {
		t.Error("Pull() accepted a diverged history — --ff-only exists to refuse exactly this")
	}
}

// The error path nobody reads until the day it fires. Outside a repository git exits
// non-zero, and the calls that report state have to say so rather than answering
// "clean" and "no remote" — an update that trusted those answers would act on a
// directory it does not understand.
func TestOutsideARepositoryTheReadsError(t *testing.T) {
	if testing.Short() {
		t.Skip("real-git integration")
	}
	s := Shell{Root: t.TempDir()}

	if _, err := s.Dirty(); err == nil {
		t.Error("Dirty() succeeded outside a repository")
	}
	if _, err := s.HasRemote(); err == nil {
		t.Error("HasRemote() succeeded outside a repository")
	}
	if _, err := s.ChangedSince("HEAD~1"); err == nil {
		t.Error("ChangedSince() succeeded outside a repository")
	}
}
