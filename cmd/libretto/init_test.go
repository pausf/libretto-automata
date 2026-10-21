package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The command is a straight line — one git call, one probe, three stats — so most
// of these run against a fake gitRunner pointed at a t.TempDir(). The two cases
// that need real git (a subdirectory cwd, a directory that is not a repo) build a
// throwaway repository and skip under -short.

// initGit returns root for rev-parse and logs every call, so a test can assert
// both what init asked git and that it asked nothing else.
func initGit(root string, log *[]string) gitRunner {
	return func(args ...string) (string, error) {
		if log != nil {
			*log = append(*log, strings.Join(args, " "))
		}
		return root + "\n", nil
	}
}

// initArtifacts is what a bootstrap must leave behind, relative to the root.
var initArtifacts = []string{
	filepath.Join(".agents", "specs", ".gitkeep"),
	filepath.Join(".agents", "changes", ".gitkeep"),
	filepath.Join(".agents", "lessons.md"),
}

func TestInitBootstrapsAFreshProject(t *testing.T) {
	root := t.TempDir()

	out, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err != nil {
		t.Fatalf("init on a fresh project: %v", err)
	}

	for _, rel := range initArtifacts {
		if !exists(filepath.Join(root, rel)) {
			t.Errorf("%s was not created", rel)
		}
	}
	if got := strings.Count(out, "created"); got != 3 {
		t.Errorf("want 3 created lines, got %d:\n%s", got, out)
	}
}

// A dangling symlink is present, not missing: os.Stat follows it and reports
// not-exist, and a write through it lands wherever it points — outside .agents/.
func TestInitNeverWritesThroughASymlink(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "elsewhere.md")
	if err := os.MkdirAll(filepath.Join(root, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".agents", "lessons.md")); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err != nil {
		t.Fatalf("init over a dangling symlink: %v", err)
	}
	if _, err := os.Lstat(outside); err == nil {
		t.Fatalf("init wrote through the symlink to %s", outside)
	}
	if !strings.Contains(out, "exists   .agents/lessons.md") {
		t.Errorf("the symlink must read as present:\n%s", out)
	}
}

func TestInitIsIdempotent(t *testing.T) {
	root := t.TempDir()

	if _, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) }); err != nil {
		t.Fatalf("first run: %v", err)
	}
	lessons := filepath.Join(root, ".agents", "lessons.md")
	before, err := os.ReadFile(lessons)
	if err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err != nil {
		t.Fatalf("second run over its own output: %v", err)
	}
	if got := strings.Count(out, "exists"); got != 3 {
		t.Errorf("want 3 exists lines on the second run, got %d:\n%s", got, out)
	}
	if strings.Contains(out, "created") {
		t.Errorf("second run reported a creation:\n%s", out)
	}
	after, err := os.ReadFile(lessons)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("second run modified lessons.md")
	}
}

func TestInitFillsAPartialLayoutWithoutTouchingWhatExists(t *testing.T) {
	root := t.TempDir()

	// A lessons.md with entries — content the header alone would not reproduce, so
	// an overwrite cannot pass as a no-op.
	seeded := "# Lessons\n\ncustom preamble\n\n## 2026-08-20 · some-change · phase 6\n\n> an entry that must survive\n"
	lessons := filepath.Join(root, ".agents", "lessons.md")
	if err := os.MkdirAll(filepath.Dir(lessons), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lessons, []byte(seeded), 0o644); err != nil {
		t.Fatal(err)
	}

	out, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err != nil {
		t.Fatalf("init on a partial layout: %v", err)
	}

	body, err := os.ReadFile(lessons)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != seeded {
		t.Errorf("lessons.md is not byte-identical after the run:\n%q", body)
	}
	for _, rel := range initArtifacts[:2] {
		if !exists(filepath.Join(root, rel)) {
			t.Errorf("missing artifact %s was not filled in", rel)
		}
	}
	if got := strings.Count(out, "created"); got != 2 {
		t.Errorf("want 2 created lines, got %d:\n%s", got, out)
	}
	if got := strings.Count(out, "exists"); got != 1 {
		t.Errorf("want 1 exists line, got %d:\n%s", got, out)
	}
}

func TestInitRefusesAForeignLayoutAndNamesIt(t *testing.T) {
	for _, foreign := range []string{"specs", "openspec", "docs/specs", "spec", "changes", "openspec/changes"} {
		t.Run(foreign, func(t *testing.T) {
			root := t.TempDir()
			if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(foreign)), 0o755); err != nil {
				t.Fatal(err)
			}

			_, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
			if err == nil {
				t.Fatalf("init did not refuse with %s/ present", foreign)
			}
			if !strings.Contains(err.Error(), foreign) {
				t.Errorf("the refusal does not name %s: %v", foreign, err)
			}
			if exists(filepath.Join(root, ".agents")) {
				t.Error("init created .agents/ while refusing")
			}
		})
	}

	// A *file* named like a root is not a consolidation root and must not refuse.
	t.Run("a file named specs is not a layout", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "specs"), []byte("just a file"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) }); err != nil {
			t.Errorf("init refused over a plain file named specs: %v", err)
		}
	})
}

func TestInitRefusalWinsOverAPartialLayout(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	// The partial .agents/ that would otherwise invite a fill.
	if err := os.MkdirAll(filepath.Join(root, ".agents", "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(root, ".agents", "specs", ".gitkeep")
	if err := os.WriteFile(keep, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	_, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err == nil {
		t.Fatal("init filled a partial layout beside a foreign root instead of refusing")
	}
	if exists(filepath.Join(root, ".agents", "changes")) {
		t.Error("init created .agents/changes/ while refusing")
	}
	if exists(filepath.Join(root, ".agents", "lessons.md")) {
		t.Error("init created lessons.md while refusing")
	}
}

func TestInitRefusesOutsideAGitRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("runs real git against a non-repo directory")
	}

	dir := t.TempDir()
	err := initCmd(os.Stdout, execGit(dir))
	if err == nil {
		t.Fatal("init did not refuse outside a git repository")
	}
	if !strings.Contains(err.Error(), "git") {
		t.Errorf("the refusal does not name git: %v", err)
	}
	if exists(filepath.Join(dir, ".agents")) {
		t.Error("init created .agents/ outside a repository")
	}
}

// tempRepo builds a real git repository and returns its root, symlinks resolved —
// on macOS t.TempDir() lives under /var, which git reports as /private/var.
func tempRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if out, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func TestInitActsOnTheRepositoryRootNotTheCwd(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a real git repository")
	}

	root := tempRepo(t)
	sub := filepath.Join(root, "deep", "inside")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, _, err := capture(t, func() error { return initCmd(os.Stdout, execGit(sub)) }); err != nil {
		t.Fatalf("init from a subdirectory: %v", err)
	}

	for _, rel := range initArtifacts {
		if !exists(filepath.Join(root, rel)) {
			t.Errorf("%s missing at the repository root", rel)
		}
	}
	if exists(filepath.Join(sub, ".agents")) {
		t.Error("init created .agents/ at the cwd instead of the root")
	}
}

func TestInitNeverCommitsAndRunsOnlyRevParse(t *testing.T) {
	root := t.TempDir()
	var log []string

	if _, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, &log)) }); err != nil {
		t.Fatal(err)
	}
	if len(log) != 1 || log[0] != "rev-parse --show-toplevel" {
		t.Errorf("git calls = %q, want exactly [rev-parse --show-toplevel]", log)
	}
}

func TestInitLessonsHeaderIsTheStandardOne(t *testing.T) {
	root := t.TempDir()
	if _, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) }); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(filepath.Join(root, ".agents", "lessons.md"))
	if err != nil {
		t.Fatal(err)
	}
	want := "# Lessons\n\nAppend-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there;\nthey are never edited otherwise.\n"
	if string(body) != want {
		t.Errorf("lessons.md = %q, want the standard append-only header %q", body, want)
	}
}

func TestInitWritesNoReadmeIntoSpecs(t *testing.T) {
	root := t.TempDir()
	if _, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) }); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(filepath.Join(root, ".agents", "specs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name() != ".gitkeep" {
			t.Errorf("init wrote %s under .agents/specs/ — only the keep-file is allowed", e.Name())
		}
	}
}

// The needsPayload exemption, proven the way models proves its own — through run,
// with a payload root that holds nothing. The cwd is a temp directory with no git,
// so the error that comes back must be init's git refusal and never the payload gate's.
func TestInitWorksWithNoPayload(t *testing.T) {
	t.Setenv(EnvRoot, filepath.Join(t.TempDir(), "never-installed"))
	t.Chdir(t.TempDir())

	_, _, err := capture(t, func() error { return run([]string{"init"}) })
	if err != nil && strings.Contains(err.Error(), "no payload") {
		t.Fatalf("init was stopped by the payload gate: %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "git") {
		t.Errorf("expected init's own git refusal from a non-repo cwd, got: %v", err)
	}
	if exists(filepath.Join(".", ".agents")) {
		t.Error("init created .agents/ while refusing")
	}
}

func TestInitNamesTheRootAndCarriesNoEscapeCodes(t *testing.T) {
	root := t.TempDir()

	out, _, err := capture(t, func() error { return initCmd(os.Stdout, initGit(root, nil)) })
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "\x1b") {
		t.Errorf("piped init emitted ANSI escape codes:\n%q", out)
	}
	rootAt := strings.Index(out, root)
	createdAt := strings.Index(out, "created")
	if rootAt == -1 {
		t.Fatalf("the report does not name the root %s:\n%s", root, out)
	}
	if createdAt != -1 && rootAt > createdAt {
		t.Errorf("the root is named after the per-artifact lines:\n%s", out)
	}
}
