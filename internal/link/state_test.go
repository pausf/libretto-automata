package link

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// sandbox builds a repo and an empty target root, and points CLAUDE_HOME at the
// target so nothing in this suite can reach the real ~/.claude.
func sandbox(t *testing.T, repoPaths ...string) (repo string, tg target.Target) {
	t.Helper()

	repo = repoWith(t, repoPaths...)
	home := t.TempDir()
	t.Setenv(target.EnvClaudeHome, home)

	return repo, target.NewClaude()
}

func stateOf(t *testing.T, entries []Entry, kind target.Kind, name string) Entry {
	t.Helper()
	for _, e := range entries {
		if e.Kind == kind && e.Name == name {
			return e
		}
	}
	t.Fatalf("no entry for %s/%s in %v", kind, name, entries)
	return Entry{}
}

func TestScanStates(t *testing.T) {
	t.Run("missing when the target has nothing", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/po-assistant/")

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		if got := stateOf(t, entries, target.Skills, "po-assistant").State; got != Missing {
			t.Errorf("State = %q, want %q", got, Missing)
		}
	})

	t.Run("linked when an owned link points at the item", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/po-assistant/")
		symlink(t, filepath.Join(repo, "skills", "po-assistant"),
			filepath.Join(tg.Dir(target.Skills), "po-assistant"))

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		e := stateOf(t, entries, target.Skills, "po-assistant")
		if e.State != Linked {
			t.Errorf("State = %q, want %q", e.State, Linked)
		}
		if e.Actual != normalise(e.RepoPath) {
			t.Errorf("Actual = %q, want %q", e.Actual, normalise(e.RepoPath))
		}
	})

	t.Run("wrong target when an owned link points elsewhere in the repo", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/po-assistant/", "skills/ci-cd-manager/")
		symlink(t, filepath.Join(repo, "skills", "ci-cd-manager"),
			filepath.Join(tg.Dir(target.Skills), "po-assistant"))

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		if got := stateOf(t, entries, target.Skills, "po-assistant").State; got != WrongTarget {
			t.Errorf("State = %q, want %q", got, WrongTarget)
		}
	})

	t.Run("conflict when a real directory is in the way", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/po-assistant/")
		occupied := filepath.Join(tg.Dir(target.Skills), "po-assistant")
		if err := os.MkdirAll(occupied, 0o755); err != nil {
			t.Fatal(err)
		}

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		if got := stateOf(t, entries, target.Skills, "po-assistant").State; got != Conflict {
			t.Errorf("State = %q, want %q", got, Conflict)
		}
	})

	t.Run("conflict when a real file is in the way", func(t *testing.T) {
		repo, tg := sandbox(t, "agents/jd-judge-a.md")
		dir := tg.Dir(target.Agents)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "jd-judge-a.md"), []byte("theirs"), 0o644); err != nil {
			t.Fatal(err)
		}

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		if got := stateOf(t, entries, target.Agents, "jd-judge-a.md").State; got != Conflict {
			t.Errorf("State = %q, want %q", got, Conflict)
		}
	})

	// A symlink someone else made is a deliberate choice, not a repair job.
	t.Run("conflict when a foreign symlink is in the way", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/po-assistant/")
		elsewhere := repoWith(t, "skills/po-assistant/")
		symlink(t, filepath.Join(elsewhere, "skills", "po-assistant"),
			filepath.Join(tg.Dir(target.Skills), "po-assistant"))

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		if got := stateOf(t, entries, target.Skills, "po-assistant").State; got != Conflict {
			t.Errorf("State = %q, want %q", got, Conflict)
		}
	})

	t.Run("stale when an owned link has no item behind it", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/")
		symlink(t, filepath.Join(repo, "skills", "deleted-skill"),
			filepath.Join(tg.Dir(target.Skills), "deleted-skill"))

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		e := stateOf(t, entries, target.Skills, "deleted-skill")
		if e.State != Stale {
			t.Errorf("State = %q, want %q", e.State, Stale)
		}
		if e.RepoPath != "" {
			t.Errorf("RepoPath = %q, want empty: no item exists", e.RepoPath)
		}
	})
}

// Foreign entries must not appear in the scan at all. gentle-ai's skills live in
// the same directory and are none of this tool's business.
func TestScanIgnoresForeignEntries(t *testing.T) {
	repo, tg := sandbox(t, "skills/po-assistant/")

	dir := tg.Dir(target.Skills)
	if err := os.MkdirAll(filepath.Join(dir, "sdd-apply"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.Name == "sdd-apply" || e.Name == "notes.txt" {
			t.Errorf("the scan reported a foreign entry %q as %q", e.Name, e.State)
		}
	}
}

func TestScanIsReadOnly(t *testing.T) {
	repo, tg := sandbox(t, "skills/po-assistant/", "agents/one.md")

	if _, err := Scan(repo, tg); err != nil {
		t.Fatal(err)
	}
	for _, kind := range tg.Kinds() {
		if _, err := os.Stat(tg.Dir(kind)); err == nil {
			t.Errorf("Scan() created %s", tg.Dir(kind))
		}
	}
}

func TestScanCoversEveryAcceptedKind(t *testing.T) {
	repo, tg := sandbox(t, "skills/a/", "agents/b.md", "commands/c.md")

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}

	seen := map[target.Kind]bool{}
	for _, e := range entries {
		seen[e.Kind] = true
	}
	for _, kind := range tg.Kinds() {
		if !seen[kind] {
			t.Errorf("the scan skipped kind %q", kind)
		}
	}
}

// Ordering must be stable so two runs produce identical output and a golden file
// or a diff stays meaningful.
func TestScanOrderIsStable(t *testing.T) {
	repo, tg := sandbox(t, "skills/zebra/", "skills/alpha/", "skills/middle/")

	first, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}

	if len(first) != len(second) {
		t.Fatalf("two scans returned %d and %d entries", len(first), len(second))
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("entry %d differs between runs: %+v vs %+v", i, first[i], second[i])
		}
	}

	var skills []string
	for _, e := range first {
		if e.Kind == target.Skills {
			skills = append(skills, e.Name)
		}
	}
	if !equal(skills, []string{"alpha", "middle", "zebra"}) {
		t.Errorf("skills came back as %v, want them sorted", skills)
	}
}

func TestTallyAndByState(t *testing.T) {
	repo, tg := sandbox(t, "skills/linked-one/", "skills/missing-one/")
	symlink(t, filepath.Join(repo, "skills", "linked-one"),
		filepath.Join(tg.Dir(target.Skills), "linked-one"))

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}

	counts := Tally(entries)
	if counts[Linked] != 1 {
		t.Errorf("Tally()[linked] = %d, want 1", counts[Linked])
	}
	if counts[Missing] != 1 {
		t.Errorf("Tally()[missing] = %d, want 1", counts[Missing])
	}

	if got := ByState(entries, Missing); len(got) != 1 || got[0].Name != "missing-one" {
		t.Errorf("ByState(missing) = %v, want just missing-one", got)
	}
}

func TestNeedsAttention(t *testing.T) {
	tests := map[State]bool{
		Linked:      false,
		Missing:     true,
		WrongTarget: true,
		Conflict:    true,
		Stale:       true,
	}
	for state, want := range tests {
		t.Run(string(state), func(t *testing.T) {
			if got := state.NeedsAttention(); got != want {
				t.Errorf("NeedsAttention() = %v, want %v", got, want)
			}
		})
	}
}

// The scan must work when the repo is reached through a symlinked path, which is
// the normal case on macOS temp directories.
func TestScanThroughAnAliasedRepoPath(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real-repo")
	if err := os.MkdirAll(filepath.Join(real, "skills", "po-assistant"), 0o755); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(base, "aliased-repo")
	symlink(t, real, alias)

	home := t.TempDir()
	t.Setenv(target.EnvClaudeHome, home)
	tg := target.NewClaude()
	symlink(t, filepath.Join(real, "skills", "po-assistant"),
		filepath.Join(tg.Dir(target.Skills), "po-assistant"))

	entries, err := Scan(alias, tg)
	if err != nil {
		t.Fatal(err)
	}
	if got := stateOf(t, entries, target.Skills, "po-assistant").State; got != Linked {
		t.Errorf("State = %q, want %q: the aliased repo path broke the comparison", got, Linked)
	}
}
