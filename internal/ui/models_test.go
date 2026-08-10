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
			func(int) ([]AgentRow, error) {
				out := make([]AgentRow, len(agents))
				copy(out, agents)
				return out, nil
			},
			func(_ int, names []string, model string) error {
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
		{Name: "review-design", Model: "", Shared: true},
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

// Every other menu row reports rather than describing itself: `status` says
// "21 linked · 2 missing", not "show the status". The tally is the question the
// panel was opened to answer — how much of this is still on the expensive model.
func TestMenuRowReportsTheModelTally(t *testing.T) {
	choices := []ModelChoice{
		{Name: "", Label: "session"},
		{Name: "haiku", Label: "cheapest"},
		{Name: "sonnet", Label: "everyday"},
		{Name: "opus", Label: "most capable"},
	}
	rows := []AgentRow{
		{Name: "a", Model: "haiku"},
		{Name: "b", Model: "haiku"},
		{Name: "c", Model: ""},
		{Name: "d", Model: "opus"},
	}

	got := Tally(rows, choices)
	want := "2 on haiku · 1 on opus · 1 on session"
	if got != want {
		t.Errorf("Tally() = %q, want %q", got, want)
	}
}

// The session default goes last even though the catalogue lists it first. It is
// not a price, and a row that opens with the one entry that cannot answer "how
// much is still expensive?" is not answering the question.
func TestTallyPutsTheSessionDefaultLast(t *testing.T) {
	choices := []ModelChoice{{Name: ""}, {Name: "haiku"}, {Name: "opus"}}
	rows := []AgentRow{{Name: "a", Model: ""}, {Name: "b", Model: "opus"}}

	if got, want := Tally(rows, choices), "1 on opus · 1 on session"; got != want {
		t.Errorf("Tally() = %q, want %q", got, want)
	}
}

// A screen that requires a reopen to tell the truth is a screen that lies for as
// long as it is open — and the menu row behind it lies too.
func TestMenuRowTallyRefreshesAfterApplying(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())

	// The caller owns the menu, so the refresh is what rebuilds the row. Here it
	// recomputes the tally from whatever the agents callback now reports.
	m = m.WithRefresh(func(int) ([]MenuItem, []TargetRow, error) {
		rows, err := m.listAgents(0)
		if err != nil {
			return nil, nil, err
		}
		return []MenuItem{
			{Label: "install", Desc: "link", Enabled: true},
			{Label: "models", Desc: Tally(rows, m.ModelChoices()), Enabled: true},
		}, []TargetRow{{Name: "claude", Active: true}}, nil
	})

	m = openSelector(t, m)
	m = key(m, "a")
	m = key(m, "m")
	m = chooseModel(t, m, "haiku")

	row := m.MenuItemForTest(1)
	if row.Desc != "3 on haiku" {
		t.Errorf("menu row = %q, want %q", row.Desc, "3 on haiku")
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

// `shared` is a warning, not decoration: applying to a row that is a symlink into the
// repository changes every project on the machine, and applying to a real file in the
// target changes that target only. The two are indistinguishable without the marker.
func TestSharedAgentsAreMarked(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	for _, line := range strings.Split(strip(m.View()), "\n") {
		switch {
		case strings.Contains(line, "review-design"):
			if !strings.Contains(line, "shared") {
				t.Errorf("a shared agent was not marked: %q", line)
			}
		case strings.Contains(line, "review-tests"):
			if strings.Contains(line, "shared") {
				t.Errorf("a target-local agent was marked shared: %q", line)
			}
		}
	}
}

// The strip already shipped the other version of this: selection encoded in a colour,
// correct behaviour reported as a bug. Colour is the signal that vanishes first.
func TestSharedMarkerIsLegibleWithoutColour(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	if !strings.Contains(strip(m.View()), "shared") {
		t.Error("the shared marker is invisible with colour stripped")
	}
}

// A screen still showing one destination's agents under another's name is exactly the
// lie the destination strip exists to prevent.
func TestTabReloadsTheSelectorForTheNewDestination(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())

	// Open first, with the original rows loaded. Swapping the callback before
	// opening made the screen show the new rows without tab doing anything — the
	// test passed with the reload deleted, which a mutation caught.
	m = openSelector(t, m)
	if len(m.AgentRows()) != 3 {
		t.Fatalf("opened with %d rows, want the original three", len(m.AgentRows()))
	}

	other := []AgentRow{{Name: "sdd-apply", Model: "sonnet"}}
	m = m.WithAgents(m.ModelChoices(),
		func(int) ([]AgentRow, error) { return other, nil },
		func(int, []string, string) error { return nil })

	m = key(m, "tab")

	got := m.AgentRows()
	if len(got) != 1 || got[0].Name != "sdd-apply" {
		t.Errorf("rows after tab = %v, want the other destination's agents", got)
	}
	if !m.InSelector() {
		t.Error("tab left the selector")
	}
}

// Showing the previous destination's rows under the new name would be worse than
// showing an error, so a failed reload keeps what it had and says so.
func TestAFailedReloadKeepsTheRowsAndSaysSo(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = m.WithAgents(m.ModelChoices(),
		func(int) ([]AgentRow, error) { return nil, errors.New("unreadable") },
		func(int, []string, string) error { return nil })
	m = key(m, "tab")

	if len(m.AgentRows()) != 3 {
		t.Errorf("rows = %d after a failed reload, want the three it had", len(m.AgentRows()))
	}
	if !strings.Contains(m.Notice(), "unreadable") {
		t.Errorf("notice = %q, want the error in it", m.Notice())
	}
}

func TestAnEmptyAgentSetSaysSo(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, nil)
	m = openSelector(t, m)

	if !strings.Contains(strip(m.View()), "no agents") {
		t.Errorf("an empty destination should say so:\n%s", strip(m.View()))
	}
}
