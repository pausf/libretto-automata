package main

import (
	"bytes"
	"io"
	"os"
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
		WithRunner(func(action string, dest int) ([]string, error) {
			return runCaptured(action, f.Repo, target.Resolve(scopeOrder[dest], f.Project))
		})

	// tab → project, enter → install (row 0), q → leave.
	in := bytes.NewReader([]byte{0x09, 0x0d, 'q'})
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

	lines, err := runCaptured("install", f.Repo, f.global())
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

	lines, err := runCaptured("install", f.Repo, f.global())
	if err == nil {
		t.Fatal("a conflict returned no error")
	}
	if len(lines) == 0 {
		t.Fatal("the failure hid the report that explains it")
	}
}
