package main

import (
	"bytes"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// The one piece of glue nothing else covered: a real tea.Program is run, keys are
// fed in, and the model that comes back out has to carry the choice and the
// destination. Without this, "the panel installs" rests on the assumption that
// Run() hands back the model as of the last Update.
//
// tea.WithInput and tea.WithOutput are what make it possible without a terminal.
func TestPanelRunReturnsTheChoiceAndScope(t *testing.T) {
	f := newFixture(t)
	f.skill(t, "alpha")

	menu, targets, err := panelData(f.Repo, target.GlobalScope)
	if err != nil {
		t.Fatal(err)
	}
	model := ui.NewModel("v0", menu, targets, false).
		WithRefresh(func(i int) ([]ui.MenuItem, []ui.TargetRow, error) {
			return panelData(f.Repo, scopeOrder[i])
		})

	// tab → project, then enter on install (row 0, selected at start).
	in := bytes.NewReader([]byte{0x09, 0x0d})
	p := tea.NewProgram(model, tea.WithInput(in), tea.WithOutput(io.Discard))

	done := make(chan tea.Model, 1)
	errc := make(chan error, 1)
	go func() {
		final, err := p.Run()
		if err != nil {
			errc <- err
			return
		}
		done <- final
	}()

	select {
	case err := <-errc:
		t.Fatalf("the program failed: %v", err)
	case final := <-done:
		m, ok := final.(ui.Model)
		if !ok {
			t.Fatalf("Run returned %T, not ui.Model — the dispatch would silently do nothing", final)
		}
		if got := m.Chosen(); got != "install" {
			t.Fatalf("Chosen() = %q, want install", got)
		}
		if got := m.ActiveScope(); got != 1 {
			t.Fatalf("ActiveScope() = %d, want 1 (project)", got)
		}
	case <-time.After(10 * time.Second):
		p.Kill()
		t.Fatal("the panel never exited after enter")
	}
}

// And the choice, run, installs where the strip pointed.
func TestPanelChoiceInstallsIntoTheChosenScope(t *testing.T) {
	f := newFixture(t)
	item := f.skill(t, "alpha")

	if _, _, err := capture(t, func() error {
		return dispatch("install", f.Repo, target.Resolve(scopeOrder[1], f.Project))
	}); err != nil {
		t.Fatal(err)
	}
	if !isSymlinkTo(t, f.projectDest("skills", "alpha"), item) {
		t.Fatal("the panel's choice installed nothing into the project")
	}
}
