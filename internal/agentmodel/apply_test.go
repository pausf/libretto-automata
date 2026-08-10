package agentmodel

import (
	"os"
	"path/filepath"
	"testing"
)

// dirWith builds a throwaway directory of agent files, one per name.
//
// It replaces a helper that built a repo root with an `agents/` inside it. The
// assertions below did not change — only the shape of the fixture, because the
// package no longer knows that a repository is where agents come from.
//
// The payload's own agents are never touched by a test.
func dirWith(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		body := "---\nname: " + name + "\ndescription: A lens.\n---\n\nBody.\n"
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func modelOf(t *testing.T, dir, name string) string {
	t.Helper()
	got, err := ReadModel(filepath.Join(dir, name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestAgentsListsEveryAgentSorted(t *testing.T) {
	dir := dirWith(t, "work-reviewer", "review-design", "spec-writer")

	agents, _, err := Agents(dir)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"review-design", "spec-writer", "work-reviewer"}
	if len(agents) != len(want) {
		t.Fatalf("Agents() returned %d agents, want %d", len(agents), len(want))
	}
	for i := range want {
		if agents[i].Name != want[i] {
			t.Errorf("Agents()[%d].Name = %q, want %q", i, agents[i].Name, want[i])
		}
	}
}

func TestAgentsReportsEachCurrentModel(t *testing.T) {
	dir := dirWith(t, "review-design", "spec-writer")
	if err := SetModel(filepath.Join(dir, "review-design.md"), "haiku"); err != nil {
		t.Fatal(err)
	}

	agents, _, err := Agents(dir)
	if err != nil {
		t.Fatal(err)
	}
	if agents[0].Model != "haiku" {
		t.Errorf("review-design model = %q, want %q", agents[0].Model, "haiku")
	}
	if agents[1].Model != Default {
		t.Errorf("spec-writer model = %q, want the default", agents[1].Model)
	}
}

func TestApplyReachesEveryAgentInTheSet(t *testing.T) {
	dir := dirWith(t, "review-design", "review-tests", "review-security")

	if err := Apply(dir, []string{"review-design", "review-tests"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if got := modelOf(t, dir, "review-design"); got != "haiku" {
		t.Errorf("review-design = %q, want haiku", got)
	}
	if got := modelOf(t, dir, "review-tests"); got != "haiku" {
		t.Errorf("review-tests = %q, want haiku", got)
	}
	if got := modelOf(t, dir, "review-security"); got != Default {
		t.Errorf("review-security = %q, want untouched", got)
	}
}

// The test this design exists for.
//
// A writer that validates as it goes leaves the user with a half-applied set and no
// way to know how far it got. So the whole set is checked before any file is opened,
// and the proof is that the *first* member is untouched when the last one is bad.
func TestApplyWritesNothingWhenAnyAgentIsUnwritable(t *testing.T) {
	dir := dirWith(t, "review-design", "review-tests")

	// A file that is not an agent file: no frontmatter, so SetModel refuses it.
	broken := filepath.Join(dir, "broken.md")
	if err := os.WriteFile(broken, []byte("Just a document.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Apply(dir, []string{"review-design", "broken"}, "haiku")
	if err == nil {
		t.Fatal("Apply() accepted a set containing an unwritable agent")
	}
	if got := modelOf(t, dir, "review-design"); got != Default {
		t.Errorf("review-design = %q — the first member was written before the set failed", got)
	}
}

// The sharper version of the test above, and the one that actually pins the
// ordering.
//
// Here the bad member is a perfectly good agent file that simply cannot be written.
// Nothing catches it while listing — its frontmatter reads fine — so the only thing
// standing between the user and a half-applied set is that Apply asks about every
// file before it writes any. Remove that and the first member gets its model while
// the second does not.
func TestApplyWritesNothingWhenAMemberIsReadOnly(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: the permission bits would not stop a write")
	}
	dir := dirWith(t, "review-design", "review-tests")

	locked := filepath.Join(dir, "review-tests.md")
	if err := os.Chmod(locked, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o644) })

	if err := Apply(dir, []string{"review-design", "review-tests"}, "haiku"); err == nil {
		t.Fatal("Apply() accepted a set containing a read-only agent")
	}
	if got := modelOf(t, dir, "review-design"); got != Default {
		t.Errorf("review-design = %q — written before the set was known to be writable", got)
	}
}

func TestApplyRejectsAnUnknownModel(t *testing.T) {
	dir := dirWith(t, "review-design")

	if err := Apply(dir, []string{"review-design"}, "gpt-4"); err == nil {
		t.Fatal("Apply() accepted a model the catalogue does not list")
	}
	if got := modelOf(t, dir, "review-design"); got != Default {
		t.Errorf("review-design = %q, want untouched", got)
	}
}

func TestApplyRejectsAnUnknownAgent(t *testing.T) {
	dir := dirWith(t, "review-design")

	if err := Apply(dir, []string{"review-design", "no-such-agent"}, "haiku"); err == nil {
		t.Fatal("Apply() accepted an agent the repo does not have")
	}
	if got := modelOf(t, dir, "review-design"); got != Default {
		t.Errorf("review-design = %q — written despite an unknown name in the set", got)
	}
}

// Default is not in Valid()'s list, because removing the key is a different act
// from declaring a model. Apply still has to accept it, or the panel's "back to the
// session's model" choice has no way through.
func TestApplyAcceptsTheDefaultToClearTheKey(t *testing.T) {
	dir := dirWith(t, "review-design")
	if err := Apply(dir, []string{"review-design"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if err := Apply(dir, []string{"review-design"}, Default); err != nil {
		t.Fatal(err)
	}
	if got := modelOf(t, dir, "review-design"); got != Default {
		t.Errorf("review-design = %q, want the key gone", got)
	}
}

func TestApplyRejectsAnEmptySet(t *testing.T) {
	dir := dirWith(t, "review-design")

	if err := Apply(dir, nil, "haiku"); err == nil {
		t.Error("Apply() accepted an empty set — nothing marked must not mean everything")
	}
}

// A target that has never had an agent installed has no agents/ directory. That is a
// state, not a failure — and reporting it as an error would make every caller
// special-case os.IsNotExist to render an empty list.
func TestAgentsOnAMissingDirectoryIsEmptyNotAnError(t *testing.T) {
	agents, _, err := Agents(filepath.Join(t.TempDir(), "no-such-directory"))
	if err != nil {
		t.Fatalf("Agents() on a missing directory returned %v, want no error", err)
	}
	if len(agents) != 0 {
		t.Errorf("Agents() = %d agents on a missing directory, want none", len(agents))
	}
}

// The mechanism behind the whole `shared` marker.
//
// An agent installed by this tool is a symlink into the repository, so writing it
// edits the repository's file and every other target linking to it. That is ordinary
// file behaviour rather than anything this package does, which is exactly why it
// needs a test: the callers promise it to the user, and a promise nothing pins is a
// promise that quietly stops being true.
func TestApplyThroughASymlinkWritesTheDestination(t *testing.T) {
	source := dirWith(t, "review-design")
	installed := t.TempDir()

	link := filepath.Join(installed, "review-design.md")
	if err := os.Symlink(filepath.Join(source, "review-design.md"), link); err != nil {
		t.Fatal(err)
	}

	if err := Apply(installed, []string{"review-design"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if got := modelOf(t, source, "review-design"); got != "haiku" {
		t.Errorf("the symlink's destination = %q, want haiku — the write did not reach it", got)
	}
	if fi, err := os.Lstat(link); err != nil {
		t.Fatal(err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Error("the symlink was replaced by a regular file")
	}
}

// A stale symlink is an ordinary state — `prune` exists for it, and renaming an agent
// creates one in every target that had the old name. One broken entry must not take
// the other eleven down with it.
func TestAgentsSkipsAStaleLinkAndNamesIt(t *testing.T) {
	dir := dirWith(t, "spec-writer", "work-reviewer")
	if err := os.Symlink(filepath.Join(dir, "gone.md"), filepath.Join(dir, "stale.md")); err != nil {
		t.Fatal(err)
	}

	agents, unreadable, err := Agents(dir)
	if err != nil {
		t.Fatalf("a stale link made the whole listing fail: %v", err)
	}
	if len(agents) != 2 {
		t.Errorf("listed %d agents, want the two that are readable", len(agents))
	}
	if len(unreadable) != 1 || unreadable[0] != "stale" {
		t.Errorf("unreadable = %v, want [stale] — skipping in silence hides state", unreadable)
	}
}

// The other half of that distinction. A file that is present but is not an agent is
// still an error: Apply's all-or-nothing guarantee rests on it, and skipping it would
// let `--all` quietly write around a file somebody put there on purpose.
func TestAgentsStillFailsOnAPresentFileWithNoFrontmatter(t *testing.T) {
	dir := dirWith(t, "spec-writer")
	if err := os.WriteFile(filepath.Join(dir, "NOTES.md"), []byte("Just a document.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, _, err := Agents(dir); err == nil {
		t.Error("a present non-agent file was accepted — Apply's refusal depends on this")
	}
}
