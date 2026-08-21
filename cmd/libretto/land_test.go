package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// These tests build real temporary git repositories: the staged index is a git fact,
// and a faked one proves nothing about the command that reads the real thing. Only
// the -z diff parser is driven on fixed strings — parsing is the one piece whose
// logic is independent of git.

func gitLand(t *testing.T, dir string, args ...string) string {
	t.Helper()
	c := exec.Command("git", args...)
	c.Dir = dir
	// Pinned identity and neutered config, so the suite does not depend on — or
	// touch — anything of the person running it.
	c.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=land-test", "GIT_AUTHOR_EMAIL=land@test",
		"GIT_COMMITTER_NAME=land-test", "GIT_COMMITTER_EMAIL=land@test",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func landInit(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("builds a real git repository")
	}
	dir := t.TempDir()
	gitLand(t, dir, "init", "-q")
	return dir
}

func writeLand(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func appendLand(t *testing.T, dir, rel, content string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, append(b, []byte(content)...), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitLand(t *testing.T, dir string) {
	t.Helper()
	gitLand(t, dir, "add", "-A")
	gitLand(t, dir, "commit", "-q", "-m", "seed")
}

// landFixture: one capability spec and one change whose delta targets it, committed.
func landFixture(t *testing.T) string {
	dir := landInit(t)
	writeLand(t, dir, ".agents/specs/cli/spec.md", "# cli\n\n- the CLI shall verify\n")
	writeLand(t, dir, ".agents/changes/add-x/spec.md", "Targets: cli\n\n# delta\n\n- it shall land\n")
	writeLand(t, dir, ".agents/changes/add-x/proposal.md", "Tracker: none\n")
	commitLand(t, dir)
	return dir
}

// stageLanding stages the complete landing: folder deleted, capability spec modified.
func stageLanding(t *testing.T, dir string) {
	t.Helper()
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- and it shall stay verified\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")
}

func runLand(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	var out, errw strings.Builder
	err := land(&out, &errw, args, execGit(dir))
	return out.String(), err
}

func TestLandPassesACompleteLanding(t *testing.T) {
	dir := landFixture(t)
	stageLanding(t, dir)
	out, err := runLand(t, dir, "add-x")
	if err != nil {
		t.Fatalf("a complete landing must pass, got: %v\n%s", err, out)
	}
	if !strings.Contains(out, "add-x") {
		t.Errorf("the report must name the change:\n%s", out)
	}
}

func TestLandFailsAPartialFolderDeletion(t *testing.T) {
	// A tracked file whose deletion was not staged survives the commit, by name.
	dir := landFixture(t)
	writeLand(t, dir, ".agents/changes/add-x/notes.md", "left behind\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-q", ".agents/changes/add-x/spec.md", ".agents/changes/add-x/proposal.md")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")
	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("a partial deletion must fail:\n%s", out)
	}
	if !strings.Contains(out, "notes.md") || !strings.Contains(out, "part 4") {
		t.Errorf("the survivor must be named under part 4:\n%s", out)
	}

	// A rename out of the folder is not a survivor: its source leaves the index,
	// and the folder is empty after the commit.
	dir = landFixture(t)
	gitLand(t, dir, "mv", ".agents/changes/add-x/proposal.md", "kept-proposal.md")
	gitLand(t, dir, "rm", "-q", ".agents/changes/add-x/spec.md")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")
	if out, err := runLand(t, dir); err != nil {
		t.Errorf("a rename out of the folder must not read as a survivor: %v\n%s", err, out)
	}
}

func TestLandFailsAnUntrackedLeftover(t *testing.T) {
	dir := landFixture(t)
	writeLand(t, dir, ".gitignore", "*.scratch\n")
	commitLand(t, dir)
	stageLanding(t, dir)
	// The staged deletion emptied the folder on disk; these two put files back.
	// The spaced name is deliberate — line-splitting would tear it.
	writeLand(t, dir, ".agents/changes/add-x/a file with spaces.md", "leftover\n")
	writeLand(t, dir, ".agents/changes/add-x/tmp.scratch", "ignored\n")
	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("an untracked leftover must fail:\n%s", out)
	}
	if !strings.Contains(out, "a file with spaces.md") {
		t.Errorf("the leftover must be named:\n%s", out)
	}
	if strings.Contains(out, "tmp.scratch") {
		t.Errorf("an ignored scratch file survives into nobody's commit:\n%s", out)
	}
}

func TestLandFailsWhenTheCapabilitySpecDidNotMove(t *testing.T) {
	dir := landFixture(t)
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x")
	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("a landing whose delta was not applied must fail:\n%s", out)
	}
	if !strings.Contains(out, "cli") || !strings.Contains(out, "part 2") {
		t.Errorf("the unmoved capability must be named under part 2:\n%s", out)
	}
}

func TestLandNamesEveryMissingPart(t *testing.T) {
	// Both parts broken, one run, both named — stopping at the first would make
	// the repair iterative for no reason.
	dir := landFixture(t)
	writeLand(t, dir, ".agents/changes/add-x/notes.md", "left behind\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-q", ".agents/changes/add-x/spec.md", ".agents/changes/add-x/proposal.md")
	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("two missing parts must fail:\n%s", out)
	}
	if !strings.Contains(out, "notes.md") {
		t.Errorf("part 4's survivor must be named:\n%s", out)
	}
	if !strings.Contains(out, "cli") {
		t.Errorf("part 2's capability must be named:\n%s", out)
	}
}

func TestLandReadsTargetsFromHeadAndSkipsFences(t *testing.T) {
	// The delta is deleted from the working tree and the index by the landing
	// itself, so HEAD is the one place its Targets: lines still are. And a
	// Targets: inside a fence is an example, not a declaration — spec-drift's
	// own fixture shows the convention in a fence to document it.
	dir := landInit(t)
	writeLand(t, dir, ".agents/specs/cli/spec.md", "# cli\n")
	writeLand(t, dir, ".agents/changes/add-x/spec.md",
		"Targets: cli\n\nAn example:\n\n```\nTargets: bogus\n```\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")
	out, err := runLand(t, dir)
	if err != nil {
		t.Fatalf("the real target moved and the fenced one is not a target: %v\n%s", err, out)
	}
	if strings.Contains(out, "bogus") {
		t.Errorf("a fenced Targets: must not be read as a declaration:\n%s", out)
	}
}

func TestLandChecksEveryDeltasTarget(t *testing.T) {
	dir := landInit(t)
	writeLand(t, dir, ".agents/specs/cli/spec.md", "# cli\n")
	writeLand(t, dir, ".agents/specs/payload/spec.md", "# payload\n")
	writeLand(t, dir, ".agents/changes/add-x/spec.md", "Targets: cli\n")
	writeLand(t, dir, ".agents/changes/add-x/spec-payload.md", "Targets: payload\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")

	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("payload's spec did not move and must fail the landing:\n%s", out)
	}
	if !strings.Contains(out, "payload") {
		t.Errorf("the failing capability must be named:\n%s", out)
	}

	appendLand(t, dir, ".agents/specs/payload/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/payload/spec.md")
	if out, err := runLand(t, dir); err != nil {
		t.Errorf("both targets moved, the landing must pass: %v\n%s", err, out)
	}
}

func TestLandInfersTheChangeFromStagedDeletions(t *testing.T) {
	dir := landFixture(t)
	stageLanding(t, dir)
	out, err := runLand(t, dir) // no name given
	if err != nil {
		t.Fatalf("the change must be inferred from the staged deletions: %v\n%s", err, out)
	}
	if !strings.Contains(out, "add-x") {
		t.Errorf("the report must name the inferred change:\n%s", out)
	}
}

func TestLandVerifiesEveryStagedChangeFolder(t *testing.T) {
	// Two folders landing in one commit is unusual and legal; each has its own
	// contract, and refusing the commit helps nobody.
	dir := landInit(t)
	writeLand(t, dir, ".agents/specs/cli/spec.md", "# cli\n")
	writeLand(t, dir, ".agents/specs/payload/spec.md", "# payload\n")
	writeLand(t, dir, ".agents/changes/add-x/spec.md", "Targets: cli\n")
	writeLand(t, dir, ".agents/changes/add-y/spec.md", "Targets: payload\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x", ".agents/changes/add-y")
	appendLand(t, dir, ".agents/specs/cli/spec.md", "- moved\n")
	gitLand(t, dir, "add", ".agents/specs/cli/spec.md")

	out, err := runLand(t, dir)
	if err == nil {
		t.Fatalf("add-y's spec did not move and must fail:\n%s", out)
	}
	for _, want := range []string{"add-x", "add-y", "payload"} {
		if !strings.Contains(out, want) {
			t.Errorf("the report must carry %q:\n%s", want, out)
		}
	}
}

func TestLandWithNothingStagedRefuses(t *testing.T) {
	dir := landFixture(t)
	_, err := runLand(t, dir)
	if err == nil {
		t.Fatal("a verifier that verified nothing must not show green")
	}
	if !strings.Contains(err.Error(), "nothing is landing") {
		t.Errorf("the refusal must say nothing is landing, got: %v", err)
	}
}

func TestLandRefusesANamedChangeWithNothingStaged(t *testing.T) {
	dir := landFixture(t)
	_, err := runLand(t, dir, "add-x")
	if err == nil {
		t.Fatal("a named change with no staged deletions must be refused")
	}
	if !strings.Contains(err.Error(), "nothing is landing") {
		t.Errorf("the refusal must say nothing is landing, got: %v", err)
	}
}

func TestLandAllowsADeltalessFolderDeletion(t *testing.T) {
	// Deleting a queued proposal is abandoning an idea, which the flow says
	// costs nothing — not a broken landing.
	dir := landInit(t)
	writeLand(t, dir, ".agents/changes/an-idea/proposal.md", "Queued: 2026-08-20\n")
	commitLand(t, dir)
	gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/an-idea")
	out, err := runLand(t, dir)
	if err != nil {
		t.Fatalf("an abandoned proposal must pass part 2 vacuously: %v\n%s", err, out)
	}
	if !strings.Contains(out, "no delta") {
		t.Errorf("the vacuous pass must say so:\n%s", out)
	}
}

func TestLandDiscoversEveryChangeRoot(t *testing.T) {
	for _, root := range []string{".agents/changes", "changes", "openspec/changes"} {
		t.Run(root, func(t *testing.T) {
			dir := landInit(t)
			writeLand(t, dir, root+"/an-idea/proposal.md", "Queued: 2026-08-20\n")
			commitLand(t, dir)
			gitLand(t, dir, "rm", "-r", "-q", root+"/an-idea")
			out, err := runLand(t, dir)
			if err != nil {
				t.Fatalf("a deletion under %s must be discovered: %v\n%s", root, err, out)
			}
			if !strings.Contains(out, "an-idea") {
				t.Errorf("the report must name the change found under %s:\n%s", root, out)
			}
		})
	}
}

func TestLandLeavesPartThreeToSpecDrift(t *testing.T) {
	// Part 3 is owned by spec-drift --retired inside --anchors; a second
	// implementation would be two sources of truth. The report attributes it on
	// every run — a green land must never read as the whole contract passing.
	dir := landFixture(t)
	stageLanding(t, dir)
	out, err := runLand(t, dir)
	if err != nil {
		t.Fatalf("a landing that would fail --retired still passes land: %v\n%s", err, out)
	}
	if !strings.Contains(out, "spec-drift --anchors") || !strings.Contains(out, "part 3") {
		t.Errorf("a passing run must attribute part 3 to spec-drift --anchors:\n%s", out)
	}

	gitLand(t, dir, "restore", "--staged", "--worktree", ".agents/specs/cli/spec.md")
	out, err = runLand(t, dir)
	if err == nil {
		t.Fatalf("the unstaged spec must fail part 2:\n%s", out)
	}
	if !strings.Contains(out, "spec-drift --anchors") {
		t.Errorf("a failing run carries the same attribution:\n%s", out)
	}
}

// snapshotLand reads every byte the repository holds — worktree, index, refs, all of
// .git — so any write land performs shows up as a differing key.
func snapshotLand(t *testing.T, dir string) map[string]string {
	t.Helper()
	m := map[string]string{}
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			dst, err := os.Readlink(p)
			if err != nil {
				return err
			}
			m[p] = "link:" + dst
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		m[p] = string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestLandChangesNothing(t *testing.T) {
	// land runs immediately before the most destructive commit in the flow, so
	// read-only is tested rather than assumed — on the passing path and failing ones.
	for _, tc := range []struct {
		name  string
		stage func(t *testing.T, dir string)
	}{
		{"passing", stageLanding},
		{"spec did not move", func(t *testing.T, dir string) {
			gitLand(t, dir, "rm", "-r", "-q", ".agents/changes/add-x")
		}},
		{"nothing staged", func(t *testing.T, dir string) {}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := landFixture(t)
			tc.stage(t, dir)
			before := snapshotLand(t, dir)
			_, _ = runLand(t, dir)
			after := snapshotLand(t, dir)
			if !reflect.DeepEqual(before, after) {
				for k := range before {
					if before[k] != after[k] {
						t.Errorf("land changed %s", k)
					}
				}
				for k := range after {
					if _, ok := before[k]; !ok {
						t.Errorf("land created %s", k)
					}
				}
			}
		})
	}
}

func TestLandWorksWithNoPayload(t *testing.T) {
	if needsPayload([]string{"land"}) {
		t.Fatal("land reads the project being landed, not the payload tree")
	}
}

func TestLandOutsideARepositoryFails(t *testing.T) {
	if testing.Short() {
		t.Skip("runs real git")
	}
	// No repository at all: an error naming git, never an empty report.
	_, err := runLand(t, t.TempDir())
	if err == nil {
		t.Fatal("outside a repository land must refuse")
	}
	if !strings.Contains(err.Error(), "git") {
		t.Errorf("the error must name git, got: %v", err)
	}

	// Unborn HEAD: a repository with no commits cannot answer for HEAD, and the
	// failure surfaces as the git-naming error path, never a panic.
	dir := landInit(t)
	_, err = runLand(t, dir, "add-x")
	if err == nil {
		t.Fatal("an unborn HEAD must fail, not pass silently")
	}
	if !strings.Contains(err.Error(), "git") {
		t.Errorf("the error must name git, got: %v", err)
	}
}

// The -z parser is the one piece testable without a repository: fixed runner output,
// including a rename (source removed, destination touched) and spaced paths.
func TestLandParsesTheCachedDiff(t *testing.T) {
	out := strings.Join([]string{
		"A", "new.md",
		"M", ".agents/specs/cli/spec.md",
		"D", ".agents/changes/x/plan.md",
		"R100", ".agents/changes/x/a file.md", "docs/a file.md",
		"T", "mode-change.md",
		"C75", "src.md", "copy of src.md",
	}, "\x00") + "\x00"

	removed, touched, err := parseCachedDiff(out)
	if err != nil {
		t.Fatal(err)
	}
	wantRemoved := []string{".agents/changes/x/a file.md", ".agents/changes/x/plan.md"}
	wantTouched := []string{".agents/specs/cli/spec.md", "copy of src.md", "docs/a file.md", "mode-change.md", "new.md"}
	for _, p := range wantRemoved {
		if !removed[p] {
			t.Errorf("%q must be in removed, got %v", p, removed)
		}
	}
	for _, p := range wantTouched {
		if !touched[p] {
			t.Errorf("%q must be in touched, got %v", p, touched)
		}
	}
	if len(removed) != len(wantRemoved) || len(touched) != len(wantTouched) {
		t.Errorf("sets carry extras: removed %v touched %v", removed, touched)
	}

	if _, _, err := parseCachedDiff("R100\x00only-a-source\x00"); err == nil {
		t.Error("a truncated rename entry must be an error, not a silent miss")
	}
}
