package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// The glue nothing else covers: a real tea.Program, keys fed in, and the action
// actually carried out — with the report coming back inside the panel rather than
// after it.
func TestPanelRunsInstallAndReportsInPlace(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	menu, targets, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	model := ui.NewModel("v0", menu, targets, false).
		WithRefresh(func(i int) ([]ui.MenuItem, []ui.TargetRow, error) {
			return panelData(f.Repo, f.Project, scopeOrder[i])
		}).
		WithRunner(func(action string, dest int, _ bool) ([]string, error) {
			return runCaptured(action, f.Repo, target.Resolve(scopeOrder[dest], f.Project), false)
		})

	// tab → project, enter → install (row 0), q → leave.
	in := keys(0x09, 0x0d, 'q')
	p := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(io.Discard))

	done := make(chan tea.Model, 1)
	go func() { m, _ := p.Run(); done <- m }()

	select {
	case final := <-done:
		m := final.(ui.Model)
		if got := m.Results(); len(got) == 0 {
			t.Fatal("the panel shows no report after installing")
		}
		if !strings.Contains(strings.Join(m.Results(), "\n"), "alpha") {
			t.Errorf("the report does not name what was installed: %v", m.Results())
		}
	case <-time.After(10 * time.Second):
		p.Kill()
		t.Fatal("the panel never exited")
	}

	// And the work really happened, in the project.
	if !isSymlinkTo(t, f.projectDest("skills", "alpha"), item) {
		t.Fatal("the panel reported an install that did not happen")
	}
	if _, err := os.Lstat(f.dest("skills", "alpha")); !os.IsNotExist(err) {
		t.Fatal("the panel installed into the global config after tabbing to the project")
	}
}

// The report is the command's own output, captured rather than re-rendered, and
// stdout is put back afterwards — otherwise every later print goes to a dead pipe.
func TestRunCapturedReturnsLinesAndRestoresStdout(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	before := os.Stdout

	lines, err := runCaptured("install", f.Repo, f.global(), false)
	if err != nil {
		t.Fatal(err)
	}
	if os.Stdout != before {
		t.Fatal("stdout was left redirected")
	}
	if len(lines) == 0 {
		t.Fatal("no lines captured")
	}
	joined := strings.Join(lines, "\n")
	if !strings.Contains(joined, "alpha") || !strings.Contains(joined, "linked") {
		t.Errorf("the capture missed the report:\n%s", joined)
	}
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			t.Error("blank lines reached the panel, which would pad the frame with nothing")
		}
	}
}

// A failing action still reports. What half-happened is the part that explains it.
func TestRunCapturedKeepsOutputOnFailure(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.putReal(t, "skills", "alpha", "somebody else's file")

	lines, err := runCaptured("install", f.Repo, f.global(), false)
	if err == nil {
		t.Fatal("a conflict returned no error")
	}
	if len(lines) == 0 {
		t.Fatal("the failure hid the report that explains it")
	}
}

// The destructive path, driven for real: prune shows its plan, the second press
// carries it out, and only the destination the strip pointed at loses anything.
func TestPanelPruneConfirmsInPlace(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	gone := f.Repo + "/skills/gone"
	f.link(t, gone, f.dest("skills", "gone"))
	f.link(t, gone, f.projectDest("skills", "gone"))

	menu, targets, err := panelData(f.Repo, f.Project, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}
	model := ui.NewModel("v0", menu, targets, false).
		WithRefresh(func(i int) ([]ui.MenuItem, []ui.TargetRow, error) {
			return panelData(f.Repo, f.Project, scopeOrder[i])
		}).
		WithRunner(func(action string, dest int, confirm bool) ([]string, error) {
			return runCaptured(action, f.Repo, target.Resolve(scopeOrder[dest], f.Project), confirm)
		}).
		SetSelectedForTest(rowOf(t, menu, "prune"))

	// enter → plan and question, y → go ahead, q → leave.
	in := keys(0x0d, 'y', 'q')
	p := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(io.Discard))

	done := make(chan tea.Model, 1)
	go func() { m, _ := p.Run(); done <- m }()

	select {
	case final := <-done:
		if got := final.(ui.Model).Results(); len(got) == 0 {
			t.Fatal("the panel shows no report after pruning")
		}
	case <-time.After(10 * time.Second):
		p.Kill()
		t.Fatal("the panel never exited")
	}

	if _, err := os.Lstat(f.projectDest("skills", "gone")); !os.IsNotExist(err) {
		t.Error("the second press did not carry the prune out")
	}
	if _, err := os.Lstat(f.dest("skills", "gone")); err != nil {
		t.Fatal("pruning the project removed a link from the global config")
	}
}

// Showing the plan removes nothing. The whole point of asking.
func TestPanelPruneOnOnePressRemovesNothing(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")
	f.link(t, f.Repo+"/skills/gone", f.projectDest("skills", "gone"))

	menu, targets, err := panelData(f.Repo, f.Project, target.ProjectScope)
	if err != nil {
		t.Fatal(err)
	}
	model := ui.NewModel("v0", menu, targets, false).
		WithRunner(func(action string, dest int, confirm bool) ([]string, error) {
			return runCaptured(action, f.Repo, target.Resolve(scopeOrder[dest], f.Project), confirm)
		}).
		SetSelectedForTest(rowOf(t, menu, "prune"))

	in := keys(0x0d, 'q')
	p := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(io.Discard))
	done := make(chan tea.Model, 1)
	go func() { m, _ := p.Run(); done <- m }()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		p.Kill()
		t.Fatal("the panel never exited")
	}

	if _, err := os.Lstat(f.projectDest("skills", "gone")); err != nil {
		t.Fatal("a single press deleted a link")
	}
}

// keys feeds bytes one at a time, with a pause between them.
//
// bubbletea batches whatever is available on a single read into one message, so
// handing it `y` and `q` together arrives as Runes{'y','q'} — a "key" called "yq"
// that matches nothing, and a program that never quits. Real typing has gaps in it;
// the test has to as well.
func keys(b ...byte) io.Reader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		for _, c := range b {
			time.Sleep(20 * time.Millisecond)
			if _, err := w.Write([]byte{c}); err != nil {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}()
	return r
}

func TestStripShowsAllFourDestinations(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	_, rows, err := panelData(f.Repo, f.Project, target.CodexScope)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"global", "project", "codex", "opencode"}
	if len(rows) != len(want) {
		t.Fatalf("the strip has %d rows, want %d", len(rows), len(want))
	}
	active := 0
	for i, row := range rows {
		if row.Name != want[i] {
			t.Errorf("row %d is %q, want %q", i, row.Name, want[i])
		}
		if row.Active {
			active++
			if row.Name != "codex" {
				t.Errorf("the active row is %q, want codex", row.Name)
			}
		}
	}
	if active != 1 {
		t.Fatalf("%d rows are active, want exactly 1", active)
	}
}

func TestUnconfiguredDestinationRow(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	// A codex root that does not exist: the row must say so, never error.
	t.Setenv(target.EnvAgentsHome, filepath.Join(t.TempDir(), "nope"))

	_, rows, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if row.Name == "codex" {
			if row.Configured {
				t.Error("a missing codex root reports as configured")
			}
			return
		}
	}
	t.Fatal("the strip has no codex row")
}

func TestModelsRowAbsentForSkillsOnlyDestination(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	agent := "---\nname: demo\ndescription: fixture\nmodel: haiku\n---\n"
	if err := os.MkdirAll(filepath.Join(f.Claude, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.Claude, "agents", "demo.md"), []byte(agent), 0o644); err != nil {
		t.Fatal(err)
	}

	menu, _, err := panelData(f.Repo, f.Project, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	if !hasRow(menu, "models") {
		t.Fatal("the global destination has agents installed and no models row — the test cannot see the row it wants absent")
	}

	menu, _, err = panelData(f.Repo, f.Project, target.CodexScope)
	if err != nil {
		t.Fatal(err)
	}
	if hasRow(menu, "models") {
		t.Error("a skills-only destination offers a models row over agents it cannot hold")
	}
}

func hasRow(menu []ui.MenuItem, label string) bool {
	for _, m := range menu {
		if m.Label == label {
			return true
		}
	}
	return false
}
