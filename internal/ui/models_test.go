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

	// The effort half, recorded separately. One set of fields for both keys would
	// make "the model apply was called" and "the effort apply was called"
	// indistinguishable, which is the only thing several of these tests assert.
	effort      string
	effortCalls int
	effortErr   error

	// listedFor records the destination index each listing was asked for. Without
	// it the tab tests pass with the index hardcoded, which is how they shipped.
	listedFor []int
}

func selectorModel(t *testing.T, rows []AgentRow) (Model, *applied) {
	t.Helper()

	rec := &applied{}
	agents := make([]AgentRow, len(rows))
	copy(agents, rows)

	m := NewModel("v0", []MenuItem{
		{Label: "install", Desc: "link", Enabled: true},
		{Label: "models", Desc: "", Enabled: true},
	}, []TargetRow{
		{Name: "global", Active: true},
		{Name: "project"},
	}, false).
		WithRefresh(func(i int) ([]MenuItem, []TargetRow, error) {
			rows := []TargetRow{{Name: "global"}, {Name: "project"}}
			rows[i].Active = true
			return []MenuItem{
				{Label: "install", Desc: "link", Enabled: true},
				{Label: "models", Desc: "", Enabled: true},
			}, rows, nil
		}).
		WithAgents(
			catalogue(),
			effortCatalogue(),
			func(dest int) ([]AgentRow, error) {
				rec.listedFor = append(rec.listedFor, dest)
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
			func(_ int, names []string, effort string) error {
				rec.effortCalls++
				rec.names, rec.effort = names, effort
				if rec.effortErr != nil {
					return rec.effortErr
				}
				for i := range agents {
					for _, n := range names {
						if agents[i].Name == n {
							agents[i].Effort = effort
						}
					}
				}
				return nil
			},
		)
	return m, rec
}

// catalogue is the model list the panel is handed. One helper rather than a literal
// per test: the ordering rules under test are the catalogue's, so a fixture that
// drifted would be testing a catalogue nothing ships.
func catalogue() []ModelChoice {
	return []ModelChoice{
		{Name: "", Label: "the session's model"},
		{Name: "haiku", Label: "cheapest"},
		{Name: "sonnet", Label: "everyday"},
		{Name: "opus", Label: "most capable"},
	}
}

func effortCatalogue() []EffortChoice {
	return []EffortChoice{
		{Name: "", Label: "whatever the session runs at"},
		{Name: "low", Label: "short, scoped work"},
		{Name: "medium", Label: "cost-sensitive"},
		{Name: "high", Label: "the balance point"},
		{Name: "xhigh", Label: "deeper reasoning"},
		{Name: "max", Label: "the deepest"},
	}
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

	// The selector opens on the `all` row, so one down reaches the first agent — the
	// lone opus one, because the rows are grouped by model.
	m = key(m, "down")
	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 1 || got[0] != "review-security" {
		t.Fatalf("marked = %v, want [review-security]", got)
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

	// Past the `all` row, then the grouped order: the opus row, then the two on the
	// session default.
	m = key(m, "down") // review-security
	m = key(m, " ")
	m = key(m, "down") // review-design
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
	want := map[string]bool{"review-security": true, "review-design": true}
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
	// One agent, not the `all` row — the screen has to show both states at once for
	// this to prove anything.
	m = key(m, "down")
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
// order joins the row names so a wrong order reads as a sentence rather than as two
// slices to diff by eye.
func order(rows []AgentRow) string {
	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.Name
	}
	return strings.Join(names, " ")
}

func TestRowsAreGroupedByModel(t *testing.T) {
	choices := []ModelChoice{{Name: ""}, {Name: "haiku"}, {Name: "sonnet"}}
	rows := []AgentRow{
		{Name: "work-reviewer"},
		{Name: "review-lens-tests", Model: "haiku"},
		{Name: "spec-writer", Model: "sonnet"},
		{Name: "review-lens-design", Model: "haiku"},
	}

	// Catalogue order, cheapest first, the session default last. Names sort inside
	// a group.
	want := "review-lens-design review-lens-tests spec-writer work-reviewer"
	if got := order(sortRowsByModel(rows, choices)); got != want {
		t.Fatalf("grouped order = %q, want %q", got, want)
	}
}

func TestAnUnknownModelGetsItsOwnGroup(t *testing.T) {
	choices := []ModelChoice{{Name: ""}, {Name: "haiku"}}
	rows := []AgentRow{
		{Name: "b", Model: "some-future-model"},
		{Name: "a", Model: "haiku"},
		{Name: "c"},
	}

	// haiku, then the session default, then the model this build does not know —
	// the same position Tally gives it, from the same ranking.
	if got := order(sortRowsByModel(rows, choices)); got != "a c b" {
		t.Fatalf("order = %q, want %q", got, "a c b")
	}
}

// rules counts the group rules in a rendered selector. The frame's own ├───┤ rows are
// not counted: they start with a junction glyph, the group rule is indented.
func rules(rendered string) int {
	n := 0
	for _, line := range strings.Split(strip(rendered), "\n") {
		if strings.HasPrefix(line, "  ─") {
			n++
		}
	}
	return n
}

func TestGroupRuleSitsOnlyBetweenGroups(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	// threeAgents is one opus row and two session rows: two groups. Two rules — one
	// under the `all` row, one between the groups. One rule per boundary, and the
	// `all` row is a boundary.
	mixed := theme.selector(Panel{Width: 90, Agents: sortRowsByModel(threeAgents(), catalogue())})
	if got := rules(mixed); got != 2 {
		t.Fatalf("two groups drew %d rules, want 2", got)
	}
	lines := strings.Split(strings.TrimRight(strip(mixed), "\n"), "\n")
	if strings.HasPrefix(lines[len(lines)-1], "  ─") {
		t.Fatal("a rule was drawn after the last group")
	}

	uniform := theme.selector(Panel{Width: 90, Agents: []AgentRow{
		{Name: "a", Model: "haiku"}, {Name: "b", Model: "haiku"},
	}})
	// Only the `all` row's rule survives: a division of nothing is not a division.
	if got := rules(uniform); got != 1 {
		t.Fatalf("one group drew %d rules, want just the all row's", got)
	}
}

// Sorted, threeAgents is review-security (opus) first, then review-design and
// review-tests (session). Every cursor count below assumes that order.

func TestSpaceOnTheAllRowMarksEveryAgent(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	// The selector opens on `all`, so the first space marks everything.
	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 3 {
		t.Fatalf("space on the all row marked %v, want all three", got)
	}

	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 0 {
		t.Fatalf("space on the all row again left %v marked, want none", got)
	}
}

func TestTheAllRowBoxFollowsEveryAgent(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ")    // every agent marked
	m = key(m, "down") // onto the first agent
	m = key(m, " ")    // unmark it

	first := strings.Split(strip(darkTheme().selector(m.panel)), "\n")[0]
	if strings.Contains(first, "[x]") {
		t.Fatalf("the all row still reads %q with an agent unmarked", first)
	}
}

func TestTheAllRowIsNeverAppliedAsAnAgent(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ") // cursor on all: mark everything
	m = key(m, "m") // open the catalogue from the all row
	m = chooseModel(t, m, "haiku")

	for _, n := range rec.names {
		if n == "all" {
			t.Fatal("the all row reached ApplyModel as an agent name")
		}
	}
	if len(rec.names) != 3 {
		t.Fatalf("applied to %v, want the three agents", rec.names)
	}
}

// The all row means screen position and agent index stop being the same number.
// Getting that offset wrong marks the neighbouring agent and says nothing.
func TestTheCursorMarksTheRowItPointsAt(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "down") // all -> review-security
	m = key(m, "down") // -> review-design
	m = key(m, " ")

	got := m.MarkedAgents()
	if len(got) != 1 || got[0] != "review-design" {
		t.Fatalf("marked %v, want [review-design] — the all row shifted the cursor", got)
	}
}

func TestTheAllRowIsLegibleWithoutColour(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)
	m = key(m, " ")

	if plain := strip(m.View()); !strings.Contains(plain, "[x] all") {
		t.Fatalf("the all row is not legible without colour:\n%s", plain)
	}
}

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
	m, rec := selectorModel(t, threeAgents())

	m = openSelector(t, m)
	if got := m.ActiveScope(); got != 0 {
		t.Fatalf("opened on destination %d, want 0", got)
	}

	other := []AgentRow{{Name: "sdd-apply", Model: "sonnet"}}
	m = m.WithAgents(m.ModelChoices(), m.EffortChoices(),
		func(dest int) ([]AgentRow, error) {
			rec.listedFor = append(rec.listedFor, dest)
			return other, nil
		},
		func(int, []string, string) error { return nil },
		func(int, []string, string) error { return nil })

	m = key(m, "tab")

	// The rows, the strip and the index the caller was asked for all have to move
	// together. Asserting only the rows passed with the index hardcoded.
	if got := m.ActiveScope(); got != 1 {
		t.Errorf("destination = %d after tab, want 1", got)
	}
	if len(rec.listedFor) == 0 || rec.listedFor[len(rec.listedFor)-1] != 1 {
		t.Errorf("agents were listed for %v, want the last request to name destination 1", rec.listedFor)
	}
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

	m = m.WithAgents(m.ModelChoices(), m.EffortChoices(),
		func(int) ([]AgentRow, error) { return nil, errors.New("unreadable") },
		func(int, []string, string) error { return nil },
		func(int, []string, string) error { return nil })
	m = key(m, "tab")

	if len(m.AgentRows()) != 3 {
		t.Errorf("rows = %d after a failed reload, want the three it had", len(m.AgentRows()))
	}
	// The whole switch is abandoned, not just the rows. Moving the strip while the
	// rows stay behind is the divergence the marker exists to prevent — produced by
	// the code meant to honour it.
	if got := m.ActiveScope(); got != 0 {
		t.Errorf("destination = %d after a failed reload, want it to have stayed at 0", got)
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

// The name column used to reuse the main menu's constant, sized for `install` and
// `prune`. Agent names run past twenty characters, and pad never truncates — so a
// long name did not clip, it shoved the model and the `shared` warning rightwards and
// every row started somewhere different.
func TestTheModelColumnLinesUpWhateverTheNamesAre(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, []AgentRow{
		{Name: "review-lens-reliability", Model: "sonnet", Shared: true},
		{Name: "spec-writer", Model: "haiku"},
		{Name: "jd-judge-a", Model: ""},
	})
	m = openSelector(t, m)

	at := map[string]int{}
	for _, line := range strings.Split(strip(m.View()), "\n") {
		for _, model := range []string{"sonnet", "haiku", "(session)"} {
			if i := strings.Index(line, model); i >= 0 && strings.Contains(line, "[") {
				// Rune count, not byte offset: the cursor is `❯`, three bytes wide
				// and one column wide, so a byte index reports the selected row two
				// columns left of where it renders.
				at[model] = len([]rune(line[:i]))
			}
		}
	}
	if len(at) != 3 {
		t.Fatalf("found %d model columns, want 3: %v", len(at), at)
	}
	first := -1
	for model, col := range at {
		if first == -1 {
			first = col
			continue
		}
		if col != first {
			t.Errorf("%q starts at column %d, another model starts at %d — the column moves with the name", model, col, first)
		}
	}
}

// "Nothing happened because nothing needed to" and "nothing happened because it is
// broken" look identical from outside. Twice in one session the first was read as the
// second.
func TestApplyingTheModelTheyAlreadyHaveSaysNothingChanged(t *testing.T) {
	m, rec := selectorModel(t, []AgentRow{
		{Name: "spec-writer", Model: "sonnet"},
		{Name: "work-reviewer", Model: "sonnet"},
	})
	m = openSelector(t, m)

	m = key(m, "a")
	m = key(m, "m")
	m = chooseModel(t, m, "sonnet")

	if rec.calls != 0 {
		t.Errorf("apply was called %d times for a no-op", rec.calls)
	}
	if !strings.Contains(m.Notice(), "nothing to change") {
		t.Errorf("notice = %q, want it to say nothing changed", m.Notice())
	}
}

func chooseEffort(t *testing.T, m Model, name string) Model {
	t.Helper()

	for i := 0; i <= len(m.EffortChoices()); i++ {
		if m.ChosenEffortName() == name {
			return key(m, "enter")
		}
		m = key(m, "down")
	}
	t.Fatalf("effort %q never came under the catalogue cursor", name)
	return m
}

// markAllAndOpenEffort is the shortest route to the effort catalogue over every row.
func markAllAndOpenEffort(t *testing.T, m Model) Model {
	t.Helper()
	m = key(m, "a")
	m = key(m, "e")
	if !m.ChoosingEffort() {
		t.Fatal("e did not open the effort catalogue")
	}
	return m
}

func TestRowsShowTheirEffort(t *testing.T) {
	forceTrueColor(t)

	rows := []AgentRow{
		{Name: "deep", Model: "opus", Effort: "xhigh"},
		{Name: "inherits", Model: "opus"},
	}
	out := strip(darkTheme().selector(Panel{
		Width:        90,
		Agents:       rows,
		ModelChoices: catalogue(),
	}))

	if !strings.Contains(out, "xhigh") {
		t.Errorf("a declared effort never reaches the screen:\n%s", out)
	}
	// An empty column would read as a bug. The word is the same one the model column
	// already uses, because it is the same state.
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, "inherits") && strings.Count(line, "(session)") != 1 {
			t.Errorf("a row declaring no effort should say so once:\n%s", line)
		}
	}
}

func TestEOpensTheEffortCatalogueAndEscapeReturns(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = markAllAndOpenEffort(t, m)
	if m.ChoosingModel() {
		t.Error("e opened the model catalogue as well")
	}

	m = key(m, "esc")
	if m.ChoosingEffort() {
		t.Error("esc did not close the effort catalogue")
	}
	if !m.InSelector() {
		t.Error("esc left the selector instead of closing the catalogue")
	}
}

// enter is the one gesture nobody reads the legend for. Rebinding it to pick up a
// second key would be a silent change to a reflex.
func TestEnterStillOpensTheModelCatalogue(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)
	m = key(m, "a")

	m = key(m, "enter")
	if !m.ChoosingModel() {
		t.Error("enter no longer opens the model catalogue")
	}
	if m.ChoosingEffort() {
		t.Error("enter opened the effort catalogue")
	}
}

func TestChoosingEffortWithNothingMarkedSaysSo(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "e")

	if m.ChoosingEffort() {
		t.Error("e opened the catalogue with nothing marked")
	}
	if rec.effortCalls != 0 {
		t.Errorf("apply called %d times with nothing marked, want none", rec.effortCalls)
	}
	if !strings.Contains(m.Notice(), "nothing marked") {
		t.Errorf("notice = %q, want it to say nothing is marked", m.Notice())
	}
}

func TestChosenEffortReachesOnlyTheMarkedRows(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "down") // review-security, the lone opus row
	m = key(m, " ")
	m = key(m, "e")
	if !m.ChoosingEffort() {
		t.Fatal("e did not open the catalogue")
	}
	m = chooseEffort(t, m, "xhigh")

	if rec.effortCalls != 1 {
		t.Fatalf("apply called %d times, want once", rec.effortCalls)
	}
	if rec.effort != "xhigh" {
		t.Errorf("applied effort = %q, want xhigh", rec.effort)
	}
	if len(rec.names) != 1 || rec.names[0] != "review-security" {
		t.Errorf("applied to %v, want exactly [review-security]", rec.names)
	}
	// The model apply must not have been reached. Two keys, two callbacks.
	if rec.calls != 0 {
		t.Errorf("the model apply was called %d times by an effort change", rec.calls)
	}
}

// A screen that needs a reopen to tell the truth is a screen the user stops believing.
func TestRowsShowTheNewEffortAfterApplying(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = markAllAndOpenEffort(t, m)
	m = chooseEffort(t, m, "low")

	for _, r := range m.AgentRows() {
		if r.Effort != "low" {
			t.Errorf("%s still reads effort %q after applying low", r.Name, r.Effort)
		}
	}
	if m.ChoosingEffort() {
		t.Error("the catalogue stayed open after applying")
	}
}

// A marked row on a model with no effort levels sends the whole apply back, and the
// screen keeps every mark so the user can unmark that one row rather than start again.
func TestARefusedEffortApplyChangesNoRow(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	rec.effortErr = errors.New("review-security runs on haiku, which has no effort levels")
	m = openSelector(t, m)

	m = markAllAndOpenEffort(t, m)
	m = chooseEffort(t, m, "max")

	for _, r := range m.AgentRows() {
		if r.Effort != "" {
			t.Errorf("%s reads effort %q after a refused apply, want untouched", r.Name, r.Effort)
		}
	}
	if len(m.MarkedAgents()) != 3 {
		t.Errorf("marked = %v after a refusal, want the three still marked", m.MarkedAgents())
	}
	if !strings.Contains(m.Notice(), "no effort levels") {
		t.Errorf("notice = %q, want the refusal it was handed", m.Notice())
	}
}

// Rows group by model, cheapest first, and that stays the only grouping. Grouping by
// the pair turns four groups into twenty on a screen whose argument is that a reader
// sees the shape at a glance.
func TestRowsStillGroupByModelAlone(t *testing.T) {
	rows := sortRowsByModel([]AgentRow{
		{Name: "b", Model: "opus", Effort: "low"},
		{Name: "a", Model: "opus", Effort: "max"},
		{Name: "c", Model: "haiku"},
	}, catalogue())

	want := []string{"c", "a", "b"}
	for i, name := range want {
		if rows[i].Name != name {
			t.Fatalf("row %d = %q, want %q — grouped by model, then by name, and by nothing else",
				i, rows[i].Name, name)
		}
	}
}
