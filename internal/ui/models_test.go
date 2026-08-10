package ui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// selectorModel is a panel wired to a fixed set of agents and a recording apply.
//
// The applied set is captured rather than written anywhere: this package never
// touches the filesystem, and a test that made it do so would be testing the wrong
// thing about the wrong package.
type applied struct {
	names []string
	model string
	calls int
	err   error
}

func selectorModel(t *testing.T, rows []AgentRow) (Model, *applied) {
	t.Helper()

	rec := &applied{}
	agents := make([]AgentRow, len(rows))
	copy(agents, rows)

	choices := []ModelChoice{
		{Name: "", Label: "the session's model"},
		{Name: "haiku", Label: "cheapest"},
		{Name: "sonnet", Label: "everyday"},
		{Name: "opus", Label: "most capable"},
	}

	m := NewModel("v0", []MenuItem{
		{Label: "install", Desc: "link", Enabled: true},
		{Label: "models", Desc: "", Enabled: true},
	}, []TargetRow{{Name: "claude", Active: true}}, false).
		WithAgents(
			choices,
			func() ([]AgentRow, error) {
				out := make([]AgentRow, len(agents))
				copy(out, agents)
				return out, nil
			},
			func(names []string, model string) error {
				rec.calls++
				rec.names, rec.model = names, model
				if rec.err != nil {
					return rec.err
				}
				for i := range agents {
					for _, n := range names {
						if agents[i].Name == n {
							agents[i].Model = model
						}
					}
				}
				return nil
			},
		)
	return m, rec
}

func threeAgents() []AgentRow {
	return []AgentRow{
		{Name: "review-design", Model: ""},
		{Name: "review-tests", Model: ""},
		{Name: "review-security", Model: "opus"},
	}
}

func key(m Model, s string) Model {
	var msg tea.Msg
	switch s {
	case "enter", "esc", "up", "down", "tab":
		msg = tea.KeyMsg{Type: keyTypeFor(s)}
	case " ":
		msg = tea.KeyMsg{Type: tea.KeySpace}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
	next, _ := m.Update(msg)
	return next.(Model)
}

func keyTypeFor(s string) tea.KeyType {
	switch s {
	case "enter":
		return tea.KeyEnter
	case "esc":
		return tea.KeyEsc
	case "up":
		return tea.KeyUp
	case "down":
		return tea.KeyDown
	default:
		return tea.KeyTab
	}
}

// openSelector puts the cursor on the models row and enters.
func openSelector(t *testing.T, m Model) Model {
	t.Helper()
	m = m.SetSelectedForTest(1)
	m = key(m, "enter")
	if !m.InSelector() {
		t.Fatal("enter on the models row did not open the selector")
	}
	return m
}

func TestSelectorOpensFromTheMenuAndEscapeReturns(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())

	m = openSelector(t, m)
	m = key(m, "esc")

	if m.InSelector() {
		t.Error("esc did not return to the menu")
	}
}

// esc leaves the selector rather than quitting the panel. Sharing one key between
// "go back" and "exit the program" is how a user loses the panel trying to back out
// of a screen.
func TestEscapeInTheSelectorDoesNotQuitThePanel(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())

	m = openSelector(t, m)
	m = key(m, "esc")

	if m.Done() {
		t.Error("esc in the selector quit the whole panel")
	}
}

func TestSpaceMarksAndUnmarksTheCurrentRow(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 1 || got[0] != "review-design" {
		t.Fatalf("marked = %v, want [review-design]", got)
	}

	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 0 {
		t.Errorf("marked = %v after unmarking, want none", got)
	}
}

// A key that only ever adds leaves no way back but pressing space once per row.
func TestMarkAllTogglesEveryRow(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "a")
	if got := m.MarkedAgents(); len(got) != 3 {
		t.Fatalf("marked = %v, want all three", got)
	}

	m = key(m, "a")
	if got := m.MarkedAgents(); len(got) != 0 {
		t.Errorf("marked = %v after the second press, want none", got)
	}
}

func TestChosenModelReachesOnlyTheMarkedRows(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ")    // review-design
	m = key(m, "down") // review-tests
	m = key(m, " ")
	m = key(m, "m")

	if !m.ChoosingModel() {
		t.Fatal("m did not open the catalogue")
	}
	m = chooseModel(t, m, "haiku")

	if rec.calls != 1 {
		t.Fatalf("apply called %d times, want once", rec.calls)
	}
	if rec.model != "haiku" {
		t.Errorf("applied model = %q, want haiku", rec.model)
	}
	want := map[string]bool{"review-design": true, "review-tests": true}
	if len(rec.names) != 2 {
		t.Fatalf("applied to %v, want exactly the two marked rows", rec.names)
	}
	for _, n := range rec.names {
		if !want[n] {
			t.Errorf("applied to %q, which was not marked", n)
		}
	}
}

// A selector with a marking mechanism that sometimes ignores it teaches the user
// not to trust the marks, so nothing marked means nothing happens.
func TestApplyingWithNothingMarkedSaysSo(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "m")

	if rec.calls != 0 {
		t.Error("the catalogue opened and applied with nothing marked")
	}
	if m.ChoosingModel() {
		t.Error("the catalogue opened with nothing marked")
	}
	if !strings.Contains(strings.ToLower(m.Notice()), "nothing marked") {
		t.Errorf("notice = %q, want it to say nothing is marked", m.Notice())
	}
}

func TestRowsShowTheNewModelAfterApplying(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "a")
	m = key(m, "m")
	m = chooseModel(t, m, "haiku")

	for _, row := range m.AgentRows() {
		if row.Model != "haiku" {
			t.Errorf("%s still shows %q — the screen needs a reopen to tell the truth", row.Name, row.Model)
		}
	}
}

func TestUndeclaredAgentRendersAsSession(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	view := strip(m.View())
	if !strings.Contains(view, "(session)") {
		t.Errorf("an agent with no declared model should render as the session's:\n%s", view)
	}
}

// Colour alone fails a non-colour terminal. Both signals, the same rule the
// destination strip already follows.
func TestMarkIsLegibleWithoutColour(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)
	m = key(m, " ")

	view := strip(m.View())
	if !strings.Contains(view, "[x]") {
		t.Errorf("a marked row is invisible with colour stripped:\n%s", view)
	}
	if !strings.Contains(view, "[ ]") {
		t.Errorf("an unmarked row is invisible with colour stripped:\n%s", view)
	}
}

func TestSelectorFrameIsFlushAtEveryWidth(t *testing.T) {
	forceTrueColor(t)

	for _, term := range []int{MinPanelWidth, 80, 100, 140, 250} {
		m, _ := selectorModel(t, threeAgents())
		m = openSelector(t, m)
		next, _ := m.Update(tea.WindowSizeMsg{Width: term, Height: 40})
		m = next.(Model)

		want := ContentWidth(term) + 2
		for i, line := range strings.Split(m.View(), "\n") {
			plain := strings.TrimSpace(strip(line))
			if plain == "" || !strings.ContainsAny(string([]rune(plain)[0]), "╭│├╰") {
				continue
			}
			if got := lipgloss.Width(plain); got != want {
				t.Errorf("term %d, row %d: %d columns, want %d: %q", term, i, got, want, plain)
			}
		}
	}
}

func TestFailedApplyIsReportedAndTheScreenSurvives(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	rec.err = errors.New("read-only file")

	m = openSelector(t, m)
	m = key(m, "a")
	m = key(m, "m")
	m = chooseModel(t, m, "haiku")

	if !strings.Contains(m.Notice(), "read-only file") {
		t.Errorf("notice = %q, want the error in it", m.Notice())
	}
	if !m.InSelector() {
		t.Error("a failed apply threw the user out of the selector")
	}
}

// chooseModel moves the catalogue cursor onto name and confirms.
func chooseModel(t *testing.T, m Model, name string) Model {
	t.Helper()

	for i := 0; i <= len(m.ModelChoices()); i++ {
		if m.ChosenModelName() == name {
			return key(m, "enter")
		}
		m = key(m, "down")
	}
	t.Fatalf("model %q never came under the catalogue cursor", name)
	return m
}
