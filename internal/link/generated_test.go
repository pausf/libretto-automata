package link

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

const agentSource = `---
name: work-reviewer
description: fixture
tools: Read, Grep, Glob
model: sonnet
---

The prompt body.
`

// opencodeSandbox is a repo with one agent and an opencode target rooted in a
// temporary directory. The agent is written with real frontmatter because the whole
// classification turns on transforming it.
func opencodeSandbox(t *testing.T) (repo string, tg target.Target, source string) {
	t.Helper()

	repo = t.TempDir()
	source = filepath.Join(repo, "agents", "work-reviewer.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(source, []byte(agentSource), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv(target.EnvOpencodeHome, t.TempDir())
	return repo, target.NewOpencode(), source
}

// want is the bytes the transform produces for the fixture agent.
func want(t *testing.T, tg target.Target, source string) []byte {
	t.Helper()
	tr, ok := tg.(target.Transformer)
	if !ok {
		t.Fatal("the opencode target does not transform")
	}
	content, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	out, err := tr.Transform(target.Agents, source, content)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestGeneratedMatchingContentIsLinked(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, want(t, tg, source), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	e := stateOf(t, entries, target.Agents, "work-reviewer.md")
	if e.State != Linked {
		t.Errorf("state = %s, want linked", e.State)
	}
	if e.Actual != source {
		t.Errorf("Actual = %q, want the marker's source %q", e.Actual, source)
	}
}

// Drift is WrongTarget, not a sixth state. Ours, at the right path, with the wrong
// content, fixable by rewriting — which is what WrongTarget already means.
func TestGeneratedDriftIsWrongTarget(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	drifted := append(want(t, tg, source), []byte("\nsomething the transform would not produce\n")...)
	if err := os.WriteFile(dest, drifted, 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	if got := stateOf(t, entries, target.Agents, "work-reviewer.md").State; got != WrongTarget {
		t.Errorf("state = %s, want wrong target", got)
	}
}

func TestMarkerlessFileIsAConflict(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	// Somebody's own agent, at the path ours would take.
	if err := os.WriteFile(dest, []byte("---\nname: mine\n---\n\nhands off\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	if got := stateOf(t, entries, target.Agents, "work-reviewer.md").State; got != Conflict {
		t.Fatalf("state = %s, want conflict", got)
	}

	// And the plan refuses it rather than resolving it.
	for _, a := range Plan(entries) {
		if a.Entry.Name == "work-reviewer.md" && a.Act != Skip {
			t.Errorf("plan says %s for a conflict; must be skip", a.Act)
		}
	}
}

func TestGeneratedOrphanIsStale(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	dir := tg.Dir(target.Agents)
	// A generated file whose marker names a source that is not an item any more.
	generated(t, dir, "retired.md", filepath.Join(repo, "agents", "retired.md"))

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	e := stateOf(t, entries, target.Agents, "retired.md")
	if e.State != Stale {
		t.Fatalf("state = %s, want stale", e.State)
	}
	if e.Actual == "" {
		t.Error("a stale generated entry reports no source — the field carries where it claims to come from")
	}

	// Prune removes it; install does not.
	if len(PrunePlan(entries)) != 1 {
		t.Errorf("prune plan has %d actions, want 1", len(PrunePlan(entries)))
	}
	for _, a := range Plan(entries) {
		if a.Entry.Name == "retired.md" {
			t.Error("install planned something for a stale entry; removing is prune's job")
		}
	}
}

// A source we cannot transform is a Conflict: we cannot say what belongs at the
// destination, so nothing is touched. Never Linked, never a crash.
func TestUntransformableSourceIsAConflict(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	if err := os.WriteFile(source, []byte("no frontmatter here at all\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatalf("scan failed instead of reporting a conflict: %v", err)
	}
	if got := stateOf(t, entries, target.Agents, "work-reviewer.md").State; got != Conflict {
		t.Errorf("state = %s, want conflict", got)
	}
}

// A target that does not transform is on exactly the path it was on before.
func TestNonTransformingTargetIsUnaffected(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/")
	if _, ok := tg.(target.Transformer); ok {
		t.Fatal("the claude target transforms; this test cannot see what it means to")
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	e := stateOf(t, entries, target.Skills, "alpha")
	if e.State != Missing {
		t.Errorf("state = %s, want missing", e.State)
	}
	if e.Generated != nil {
		t.Error("a linked kind carries generated bytes")
	}
}

// A transforming target still symlinks the kinds it does not transform. This is the
// bug the first Transformer interface had: Transform errored for skills, and every
// skill in the target was classified as a conflict.
func TestTransformingTargetStillLinksItsOtherKinds(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	for _, p := range []string{"skills/alpha/SKILL.md", "commands/beta.md"} {
		full := filepath.Join(repo, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		kind target.Kind
		name string
	}{{target.Skills, "alpha"}, {target.Commands, "beta.md"}} {
		e := stateOf(t, entries, c.kind, c.name)
		if e.State != Missing {
			t.Errorf("%s/%s state = %s, want missing", c.kind, c.name, e.State)
		}
		if e.Generated != nil {
			t.Errorf("%s/%s carries generated bytes; that kind is symlinked", c.kind, c.name)
		}
	}
}
