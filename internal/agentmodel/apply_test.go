package agentmodel

import (
	"os"
	"path/filepath"
	"testing"
)

// repoWith builds a throwaway repo root holding agents/, one file per name.
// The payload's own agents are never touched by a test.
func repoWith(t *testing.T, names ...string) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		body := "---\nname: " + name + "\ndescription: A lens.\n---\n\nBody.\n"
		if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func modelOf(t *testing.T, root, name string) string {
	t.Helper()
	got, err := ReadModel(filepath.Join(root, "agents", name+".md"))
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestAgentsListsEveryAgentSorted(t *testing.T) {
	root := repoWith(t, "work-reviewer", "review-design", "spec-writer")

	agents, err := Agents(root)
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
	root := repoWith(t, "review-design", "spec-writer")
	if err := SetModel(filepath.Join(root, "agents", "review-design.md"), "haiku"); err != nil {
		t.Fatal(err)
	}

	agents, err := Agents(root)
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
	root := repoWith(t, "review-design", "review-tests", "review-security")

	if err := Apply(root, []string{"review-design", "review-tests"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if got := modelOf(t, root, "review-design"); got != "haiku" {
		t.Errorf("review-design = %q, want haiku", got)
	}
	if got := modelOf(t, root, "review-tests"); got != "haiku" {
		t.Errorf("review-tests = %q, want haiku", got)
	}
	if got := modelOf(t, root, "review-security"); got != Default {
		t.Errorf("review-security = %q, want untouched", got)
	}
}

// The test this design exists for.
//
// A writer that validates as it goes leaves the user with a half-applied set and no
// way to know how far it got. So the whole set is checked before any file is opened,
// and the proof is that the *first* member is untouched when the last one is bad.
func TestApplyWritesNothingWhenAnyAgentIsUnwritable(t *testing.T) {
	root := repoWith(t, "review-design", "review-tests")

	// A file that is not an agent file: no frontmatter, so SetModel refuses it.
	broken := filepath.Join(root, "agents", "broken.md")
	if err := os.WriteFile(broken, []byte("Just a document.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := Apply(root, []string{"review-design", "broken"}, "haiku")
	if err == nil {
		t.Fatal("Apply() accepted a set containing an unwritable agent")
	}
	if got := modelOf(t, root, "review-design"); got != Default {
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
	root := repoWith(t, "review-design", "review-tests")

	locked := filepath.Join(root, "agents", "review-tests.md")
	if err := os.Chmod(locked, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(locked, 0o644) })

	if err := Apply(root, []string{"review-design", "review-tests"}, "haiku"); err == nil {
		t.Fatal("Apply() accepted a set containing a read-only agent")
	}
	if got := modelOf(t, root, "review-design"); got != Default {
		t.Errorf("review-design = %q — written before the set was known to be writable", got)
	}
}

func TestApplyRejectsAnUnknownModel(t *testing.T) {
	root := repoWith(t, "review-design")

	if err := Apply(root, []string{"review-design"}, "gpt-4"); err == nil {
		t.Fatal("Apply() accepted a model the catalogue does not list")
	}
	if got := modelOf(t, root, "review-design"); got != Default {
		t.Errorf("review-design = %q, want untouched", got)
	}
}

func TestApplyRejectsAnUnknownAgent(t *testing.T) {
	root := repoWith(t, "review-design")

	if err := Apply(root, []string{"review-design", "no-such-agent"}, "haiku"); err == nil {
		t.Fatal("Apply() accepted an agent the repo does not have")
	}
	if got := modelOf(t, root, "review-design"); got != Default {
		t.Errorf("review-design = %q — written despite an unknown name in the set", got)
	}
}

// Default is not in Valid()'s list, because removing the key is a different act
// from declaring a model. Apply still has to accept it, or the panel's "back to the
// session's model" choice has no way through.
func TestApplyAcceptsTheDefaultToClearTheKey(t *testing.T) {
	root := repoWith(t, "review-design")
	if err := Apply(root, []string{"review-design"}, "haiku"); err != nil {
		t.Fatal(err)
	}

	if err := Apply(root, []string{"review-design"}, Default); err != nil {
		t.Fatal(err)
	}
	if got := modelOf(t, root, "review-design"); got != Default {
		t.Errorf("review-design = %q, want the key gone", got)
	}
}

func TestApplyRejectsAnEmptySet(t *testing.T) {
	root := repoWith(t, "review-design")

	if err := Apply(root, nil, "haiku"); err == nil {
		t.Error("Apply() accepted an empty set — nothing marked must not mean everything")
	}
}
