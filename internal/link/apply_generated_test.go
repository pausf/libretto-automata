package link

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// installOnce scans, plans and applies, and fails on anything refused or errored.
func installOnce(t *testing.T, repo string, tg target.Target) []Result {
	t.Helper()
	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	results := Apply(repo, Plan(entries))
	for _, r := range results {
		if r.Err != nil || r.Refused {
			t.Fatalf("apply refused or failed: %+v", r)
		}
	}
	return results
}

func TestCreateWritesGeneratedContent(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	installOnce(t, repo, tg)

	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("nothing was written: %v", err)
	}
	if !bytes.Equal(got, want(t, tg, source)) {
		t.Errorf("written bytes are not the transform's output:\n%q", got)
	}

	// A real file, never a symlink — a symlinked agent breaks OpenCode's config load.
	fi, err := os.Lstat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !fi.Mode().IsRegular() {
		t.Errorf("mode = %v, want a regular file", fi.Mode())
	}
	if perm := fi.Mode().Perm(); perm != 0o644 {
		t.Errorf("perm = %o, want 644 — CreateTemp makes 600 and an unreadable agent is no agent", perm)
	}
	if !strings.Contains(string(got), target.MarkerKey+`: "`+source+`"`) {
		t.Error("the written file carries no marker, so nothing can prove it is ours")
	}
	if !OwnedEither(repo, dest) {
		t.Error("the file we just wrote is not recognised as ours")
	}
}

func TestRepointRewritesGeneratedContent(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	installOnce(t, repo, tg)
	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")

	// The source's prose changes, so what belongs at the destination changes.
	if err := os.WriteFile(source, []byte(strings.Replace(agentSource, "The prompt body.", "A different prompt body.", 1)), 0o644); err != nil {
		t.Fatal(err)
	}

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	e := stateOf(t, entries, target.Agents, "work-reviewer.md")
	if e.State != WrongTarget {
		t.Fatalf("state = %s, want wrong target after the source changed", e.State)
	}

	installOnce(t, repo, tg)
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want(t, tg, source)) {
		t.Error("repoint did not rewrite the file to the transform's current output")
	}
	if !strings.Contains(string(got), "A different prompt body.") {
		t.Error("the new prose did not reach the destination")
	}
}

// The plan is computed from a scan that is already stale by the time it runs.
// Ownership is re-checked at apply time, and that re-check is the last thing between
// a race and somebody's file.
func TestRepointRefusesAForeignFileAtApplyTime(t *testing.T) {
	repo, tg, source := opencodeSandbox(t)
	installOnce(t, repo, tg)
	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")

	// Make the destination drift so the plan says Repoint...
	if err := os.WriteFile(source, []byte(strings.Replace(agentSource, "fixture", "changed", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	plan := Plan(entries)

	// ...then somebody replaces it with their own file between the scan and the apply.
	mine := "---\nname: mine\n---\n\nhands off\n"
	if err := os.WriteFile(dest, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	var sawRefusal bool
	for _, r := range Apply(repo, plan) {
		if r.Action.Entry.Name == "work-reviewer.md" && r.Refused {
			sawRefusal = true
		}
	}
	if !sawRefusal {
		t.Error("apply did not refuse a file that stopped being ours after the scan")
	}
	if got, _ := os.ReadFile(dest); string(got) != mine {
		t.Error("somebody else's file was overwritten")
	}
}

func TestRemoveDeletesAnOwnedGeneratedFile(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	dir := tg.Dir(target.Agents)
	orphan := generated(t, dir, "retired.md", filepath.Join(repo, "agents", "retired.md"))
	foreign := generatedRaw(t, dir, "theirs.md", "---\nname: theirs\n---\n\nnot ours\n")

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range Apply(repo, PrunePlan(entries)) {
		if r.Err != nil {
			t.Fatalf("prune failed: %+v", r)
		}
	}

	if _, err := os.Lstat(orphan); !os.IsNotExist(err) {
		t.Error("the orphaned generated file survived prune")
	}
	if _, err := os.Lstat(foreign); err != nil {
		t.Error("prune removed a file that was not ours")
	}
}

// The promise that makes install safe to re-run. A generated tree that rewrote itself
// every run would also mean `status` never reads clean.
func TestGeneratedApplyIsIdempotent(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	installOnce(t, repo, tg)

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	if got := stateOf(t, entries, target.Agents, "work-reviewer.md").State; got != Linked {
		t.Fatalf("state after install = %s, want linked", got)
	}
	if plan := Plan(entries); len(plan) != 0 {
		t.Errorf("a second install plans %d actions, want none: %+v", len(plan), plan)
	}
}

func TestGeneratedWriteLeavesNoTempFile(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)
	installOnce(t, repo, tg)

	dir := tg.Dir(target.Agents)
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, de := range des {
		if strings.HasPrefix(de.Name(), ".libretto-") {
			t.Errorf("a temporary file survived the write: %s", de.Name())
		}
	}
	if len(des) != 1 {
		t.Errorf("the agents directory holds %d entries, want exactly the one agent", len(des))
	}
}

// The temp file must be created in the destination directory: os.Rename fails across
// filesystems, and a target root is often on a different one from the system temp dir.
func TestGeneratedWriteUsesTheDestinationDirectory(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "a.md")
	if err := writeGenerated(dest, []byte("content\n")); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil || string(got) != "content\n" {
		t.Fatalf("content = %q, err = %v", got, err)
	}

	// A destination whose directory does not exist fails rather than falling back
	// somewhere else.
	if err := writeGenerated(filepath.Join(dir, "nope", "b.md"), []byte("x")); err == nil {
		t.Error("writing into a missing directory succeeded; the temp file went somewhere unintended")
	}
}

// A generated Create must refuse a file that appeared since the scan, exactly as a
// symlinked one does. os.Rename replaces silently, so the first version of this write
// destroyed the file for a generated kind while reporting a failure for a linked one.
func TestCreateRefusesAFileThatAppearedSinceTheScan(t *testing.T) {
	repo, tg, _ := opencodeSandbox(t)

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	plan := Plan(entries)

	dest := filepath.Join(tg.Dir(target.Agents), "work-reviewer.md")
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		t.Fatal(err)
	}
	mine := "---\nname: mine\n---\n\nappeared between the scan and the apply\n"
	if err := os.WriteFile(dest, []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}

	var refused bool
	for _, r := range Apply(repo, plan) {
		if r.Action.Entry.Name == "work-reviewer.md" && r.Refused {
			refused = true
		}
	}
	if !refused {
		t.Error("create did not refuse a file that appeared after the scan")
	}
	if got, _ := os.ReadFile(dest); string(got) != mine {
		t.Error("a file that appeared after the scan was overwritten")
	}
}

// The refusal above is caught by create's own Lstat, which runs before the write — so it
// passes whether the write replaces or refuses, and it does not test the fix. This does:
// the write primitive itself must refuse an existing destination, which is what closes
// the window between that Lstat and the write.
//
// os.Rename replaces silently; os.Link returns ErrExist. That asymmetry is the finding:
// a symlinked kind reported a failure in this window while a generated kind destroyed
// the file.
func TestCreateGeneratedRefusesAnExistingDestination(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "a.md")
	mine := []byte("somebody else wrote this\n")
	if err := os.WriteFile(dest, mine, 0o644); err != nil {
		t.Fatal(err)
	}

	err := createGenerated(dest, []byte("ours\n"))
	if err == nil {
		t.Fatal("createGenerated replaced an existing file")
	}
	if !errors.Is(err, os.ErrExist) {
		t.Errorf("err = %v, want os.ErrExist so create can report the same refusal a symlink does", err)
	}
	if got, _ := os.ReadFile(dest); string(got) != string(mine) {
		t.Error("the existing file was modified")
	}

	// And it leaves no temp file behind on that refusal.
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, de := range des {
		if strings.HasPrefix(de.Name(), ".libretto-") {
			t.Errorf("a temp file survived a refused create: %s", de.Name())
		}
	}

	// writeGenerated, by contrast, is allowed to replace — that is repoint's intent.
	if err := writeGenerated(dest, []byte("ours\n")); err != nil {
		t.Fatalf("writeGenerated refused to replace: %v", err)
	}
	if got, _ := os.ReadFile(dest); string(got) != "ours\n" {
		t.Error("writeGenerated did not replace the destination")
	}
}

// Prune must not offer a hand-written file in a destination that generates nothing.
func TestPruneSparesAMarkedFileInANonGeneratingDestination(t *testing.T) {
	repo, tgClaude := sandbox(t, "agents/work-reviewer.md")
	source := filepath.Join(repo, "agents", "work-reviewer.md")
	notes := generated(t, filepath.Join(tgClaude.Dir(target.Agents)), "notes.md", source)

	entries, err := Scan(repo, tgClaude)
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range PrunePlan(entries) {
		if a.Entry.Name == "notes.md" {
			t.Fatal("prune planned to remove a hand-written file in the claude destination")
		}
	}
	for _, r := range Apply(repo, PrunePlan(entries)) {
		if r.Err != nil {
			t.Fatalf("prune failed: %+v", r)
		}
	}
	if _, err := os.Lstat(notes); err != nil {
		t.Error("prune deleted a hand-written file that merely carried a marker line")
	}
}
