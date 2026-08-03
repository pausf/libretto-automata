package ui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func demoPanel() Panel {
	return Panel{
		Version: "v0.1.0",
		Menu: []MenuItem{
			{Label: "install", Desc: "link the score into ~/.claude", Enabled: false},
			{Label: "update", Desc: "git pull · relink · report", Enabled: false},
			{Label: "status", Desc: "what is linked", Enabled: true},
			{Label: "doctor", Desc: "diagnose the orchestra", Enabled: false},
			{Label: "prune", Desc: "drop links whose source is gone", Enabled: false},
		},
		Selected: 2,
		Targets: []TargetRow{
			{Name: "claude", Info: "12 skills · 8 agents · 4 commands", Configured: true},
			{Name: "codex", Info: "not configured", Configured: false},
		},
	}
}

func TestRenderWidthIsStable(t *testing.T) {
	forceTrueColor(t)
	out := darkTheme().Render(demoPanel())

	for i, line := range strings.Split(out, "\n") {
		if w := lipgloss.Width(line); w > MinPanelWidth {
			t.Errorf("row %d is %d columns wide, exceeding the panel's %d: %q",
				i, w, MinPanelWidth, strip(line))
		}
	}
}

// Every bordered row must be exactly the panel width. A short row means a torn box.
func TestBorderedRowsAreFlush(t *testing.T) {
	forceTrueColor(t)
	out := darkTheme().Render(demoPanel())

	for i, line := range strings.Split(out, "\n") {
		plain := strip(line)
		if plain == "" || !strings.ContainsAny(plain[:3], "╭│├╰") {
			continue // the footer sits outside the border
		}
		if w := lipgloss.Width(line); w != MinPanelWidth {
			t.Errorf("bordered row %d is %d columns, want %d: %q", i, w, MinPanelWidth, plain)
		}
	}
}

// Section breaks must be ├───┤, not │───│. lipgloss's Border() cannot produce a
// junction, which is the whole reason frame() exists.
func TestSectionBreaksUseJunctions(t *testing.T) {
	forceTrueColor(t)
	lines := strings.Split(strip(darkTheme().Render(demoPanel())), "\n")

	junctions := 0
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "├"):
			junctions++
			if !strings.HasSuffix(line, "┤") {
				t.Errorf("row %d opens with ├ but does not close with ┤: %q", i, line)
			}
		case strings.HasPrefix(line, "│") && strings.Contains(line, "──────────"):
			t.Errorf("row %d is a rule drawn inside │ │ instead of ├ ┤: %q", i, line)
		}
	}
	if junctions != 2 {
		t.Errorf("found %d section breaks, want 2 (menu and target strip)", junctions)
	}
}

func TestRenderContainsTheDesign(t *testing.T) {
	forceTrueColor(t)
	plain := strip(darkTheme().Render(demoPanel()))

	for _, want := range []string{
		"A U T O M A T A",
		"the libretto is written first",
		"install", "update", "status", "doctor", "prune",
		"claude", "codex",
		"v0.1.0",
	} {
		if !strings.Contains(plain, want) {
			t.Errorf("panel is missing %q", want)
		}
	}
}

func TestSelectedRowGetsTheCursor(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	plain := strip(darkTheme().Render(p))

	for _, line := range strings.Split(plain, "\n") {
		if !strings.Contains(line, "status") {
			continue
		}
		if !strings.Contains(line, "❯") {
			t.Errorf("the selected row has no cursor: %q", line)
		}
		return
	}
	t.Fatal("no row containing the selected label was rendered")
}

func TestNarrowTerminalDegradesWithoutBorders(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.Width = 40

	out := darkTheme().Render(p)
	plain := strip(out)

	if strings.ContainsAny(plain, "╭╮╰╯├┤") {
		t.Error("the narrow layout drew a box it cannot fit")
	}
	if !strings.Contains(plain, "LIBRETTO") {
		t.Error("the narrow layout dropped the title")
	}
	for i, line := range strings.Split(out, "\n") {
		if w := lipgloss.Width(line); w > p.Width {
			t.Errorf("narrow row %d is %d columns, exceeding %d", i, w, p.Width)
		}
	}
}

func TestPanelIsCentredWhenThereIsRoom(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.Width = 100

	lines := strings.Split(strip(darkTheme().Render(p)), "\n")

	var top string
	for _, l := range lines {
		if strings.Contains(l, "╭") {
			top = l
			break
		}
	}
	if top == "" {
		t.Fatal("no top border was rendered")
	}

	want := (p.Width - (ContentWidth(p.Width) + 2)) / 2
	if got := len(top) - len(strings.TrimLeft(top, " ")); got != want {
		t.Errorf("left margin = %d, want %d", got, want)
	}
}

// A piped or preview render has no known size. It must stay flush left rather
// than carry alignment whitespace into a file.
func TestUnknownWidthStaysFlushLeft(t *testing.T) {
	forceTrueColor(t)
	out := strip(darkTheme().Render(demoPanel())) // Width and Height are 0

	if !strings.HasPrefix(out, "╭") {
		t.Errorf("output does not start at column 0: %q", out[:min(20, len(out))])
	}
}

func TestPanelIsCentredVertically(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.Width, p.Height = 100, 60

	lines := strings.Split(darkTheme().Render(p), "\n")
	if len(lines) != p.Height {
		t.Fatalf("rendered %d rows, want %d", len(lines), p.Height)
	}

	blank := 0
	for _, l := range lines {
		if strings.TrimSpace(strip(l)) != "" {
			break
		}
		blank++
	}
	if blank == 0 {
		t.Error("no blank rows above the panel; it was not centred vertically")
	}
}

// A notice must not push the panel off centre — it belongs inside the centred
// block, so the panel does not jump the moment feedback appears.
//
// Asserting that the top row index changes would be wrong: whether it moves by a
// row depends on the parity of the block height, so the test would pass or fail
// by luck. What must hold is that the notice is inside the centred block, which
// shows up as the render still filling the height and the margins staying even.
func TestNoticeStaysInsideTheCentredBlock(t *testing.T) {
	forceTrueColor(t)

	p := demoPanel()
	p.Width, p.Height = 100, 60
	p.Notice = "install is not wired up yet"

	lines := strings.Split(darkTheme().Render(p), "\n")
	if len(lines) != p.Height {
		t.Fatalf("rendered %d rows, want %d", len(lines), p.Height)
	}

	above, below := 0, 0
	for _, l := range lines {
		if strings.TrimSpace(strip(l)) != "" {
			break
		}
		above++
	}
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(strip(lines[i])) != "" {
			break
		}
		below++
	}

	if above == 0 || below == 0 {
		t.Fatalf("margins are %d above and %d below; the block is not centred", above, below)
	}
	if diff := above - below; diff > 1 || diff < -1 {
		t.Errorf("margins are %d above and %d below; the notice pushed the panel off centre", above, below)
	}

	// And the notice must be the last non-blank row, not stranded outside.
	last := ""
	for _, l := range lines {
		if plain := strings.TrimSpace(strip(l)); plain != "" {
			last = plain
		}
	}
	if !strings.Contains(last, "not wired up") {
		t.Errorf("the last rendered row is %q, want the notice", last)
	}
}

func topRow(t *testing.T, p Panel) int {
	t.Helper()
	for i, l := range strings.Split(darkTheme().Render(p), "\n") {
		if strings.Contains(strip(l), "╭") {
			return i
		}
	}
	t.Fatal("no top border was rendered")
	return -1
}

func TestPadNeverTruncates(t *testing.T) {
	long := "a-very-long-label"

	if got := pad(long, 4); !strings.HasPrefix(got, long) {
		t.Errorf("pad() truncated %q to %q", long, got)
	}
	if got := pad("ab", 5); lipgloss.Width(got) != 5 {
		t.Errorf("pad(%q,5) width = %d, want 5", "ab", lipgloss.Width(got))
	}
}

func TestModelNavigationWraps(t *testing.T) {
	m := newDemoModel()
	n := len(m.panel.Menu)

	m.panel.Selected = 0
	up, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if got := up.(Model).Selected(); got != n-1 {
		t.Errorf("↑ from the first row = %d, want %d", got, n-1)
	}

	m.panel.Selected = n - 1
	down, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if got := down.(Model).Selected(); got != 0 {
		t.Errorf("↓ from the last row = %d, want 0", got)
	}
}

func TestModelSelectingADisabledActionRefuses(t *testing.T) {
	m := newDemoModel()
	m.panel.Selected = 0 // install, disabled

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("a disabled action returned a command; it must do nothing")
	}
	if notice := next.(Model).Notice(); !strings.Contains(notice, "not wired up") {
		t.Errorf("notice = %q, want it to say the action is unavailable", notice)
	}
}

// Selecting an action runs it and shows the report, without leaving. Staying is the
// point: the destination, the state and the last report end up on screen together.
//
// An earlier version asserted the notice read "running status…" and passed while
// nothing ran. A test that pins a placeholder in place makes the lie look verified.
func TestSelectingAnEnabledActionRunsAndStays(t *testing.T) {
	m := newDemoModel()
	m.panel.Selected = 2 // status, enabled
	ran := ""
	m = m.WithRunner(func(action string, dest int, _ bool) ([]string, error) {
		ran = action
		return []string{"  linked  skills/alpha", "1 linked"}, nil
	})

	next, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	after := next.(Model)

	if ran != "status" {
		t.Errorf("the runner was asked for %q, want status", ran)
	}
	if cmd != nil {
		t.Error("selecting an action quit the panel; the report is meant to be read in place")
	}
	if got := after.Results(); len(got) != 2 {
		t.Errorf("the report is %v, want the runner's two lines", got)
	}
	if !strings.Contains(after.Notice(), "status") {
		t.Errorf("the notice does not name the action: %q", after.Notice())
	}
}

// The destination is passed in, never captured. A runner closing over the scope the
// panel opened with would send prune at the old destination after a tab — the strip
// reading "project" while links vanish from the global config.
func TestTheRunnerIsToldWhichDestination(t *testing.T) {
	rows := []TargetRow{
		{Name: "global", Configured: true, Active: true},
		{Name: "project", Configured: true},
	}
	got := -1
	m := NewModel("v0", demoPanel().Menu, rows, false).
		WithRefresh(func(i int) ([]MenuItem, []TargetRow, error) {
			out := []TargetRow{{Name: "global"}, {Name: "project"}}
			out[i].Active = true
			return demoPanel().Menu, out, nil
		}).
		WithRunner(func(action string, dest int, _ bool) ([]string, error) {
			got = dest
			return nil, nil
		})
	m.panel.Selected = 2

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = next

	if got != 1 {
		t.Fatalf("the runner was told destination %d after a tab, want 1", got)
	}
}

// A refused action runs nothing and leaves no report behind.
func TestSelectingADisabledActionRunsNothing(t *testing.T) {
	m := newDemoModel()
	for i, item := range m.panel.Menu {
		if !item.Enabled {
			m.panel.Selected = i
			break
		}
	}
	called := false
	m = m.WithRunner(func(string, int, bool) ([]string, error) { called = true; return nil, nil })

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if called {
		t.Error("a disabled action was run")
	}
	if got := next.(Model).Results(); got != nil {
		t.Errorf("a refused action left a report: %v", got)
	}
}

// With no runner wired every action refuses, which is honest.
func TestNoRunnerMeansEveryActionRefuses(t *testing.T) {
	m := newDemoModel()
	m.panel.Selected = 2

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if got := next.(Model).Notice(); !strings.Contains(got, "not wired up") {
		t.Errorf("notice = %q, want it to say the action is unavailable", got)
	}
}

func TestModelWindowSizeReachesThePanel(t *testing.T) {
	m := newDemoModel()

	next, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	if got := next.(Model).panel.Width; got != 40 {
		t.Errorf("panel width = %d, want 40", got)
	}
}

func TestModelQuits(t *testing.T) {
	for _, key := range []string{"q", "esc", "ctrl+c"} {
		t.Run(key, func(t *testing.T) {
			m := newDemoModel()
			_, cmd := m.Update(keyMsg(key))
			if cmd == nil {
				t.Errorf("%q did not return a command", key)
			}
		})
	}
}

func newDemoModel() Model {
	p := demoPanel()
	m := NewModel(p.Version, p.Menu, p.Targets, false)
	m.theme = darkTheme()
	m.panel.Selected = p.Selected
	return m
}

func keyMsg(s string) tea.KeyMsg {
	switch s {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}
	default:
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
	}
}

// The active scope must be visible, and visibly different from the other one.
// A strip that lists both destinations identically is a strip that answers
// "where would install go?" with a shrug.
func TestPanelShowsTheActiveScope(t *testing.T) {
	t.Setenv(EnvTheme, "dark")
	theme := NewTheme()

	rows := []TargetRow{
		{Name: "global", Info: "10 skills", Configured: true, Active: true},
		{Name: "project", Info: "0 skills", Configured: false},
	}
	p := demoPanel()
	p.Targets, p.Width = rows, 90
	out := theme.Render(p)

	if !strings.Contains(out, "global") || !strings.Contains(out, "project") {
		t.Fatal("the strip does not list both scopes")
	}

	// Swapping which one is active must change the render. If it does not, the
	// flag is decorative.
	rows[0].Active, rows[1].Active = false, true
	p.Targets = rows
	other := theme.Render(p)

	if out == other {
		t.Fatal("marking a different scope active rendered identically")
	}
}

// Without this the panel is half a feature: you could see the other destination
// and never choose it.
func TestModelSwitchesScope(t *testing.T) {
	t.Setenv(EnvTheme, "dark")

	rows := []TargetRow{
		{Name: "global", Info: "a", Configured: true, Active: true},
		{Name: "project", Info: "b", Configured: true},
	}
	calls := 0
	m := NewModel("v0", demoPanel().Menu, rows, false).
		WithRefresh(func(i int) ([]MenuItem, []TargetRow, error) {
			calls++
			out := []TargetRow{
				{Name: "global", Info: "a", Configured: true},
				{Name: "project", Info: "b", Configured: true},
			}
			out[i].Active = true
			return demoPanel().Menu, out, nil
		})

	if got := m.ActiveScope(); got != 0 {
		t.Fatalf("starts on %d, want 0", got)
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = next.(Model)
	if got := m.ActiveScope(); got != 1 {
		t.Fatalf("tab moved to %d, want 1", got)
	}
	if calls != 1 {
		t.Fatalf("refresh called %d times, want 1 — figures must match the new destination", calls)
	}

	// And it wraps, so two destinations are reachable with one key.
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if got := next.(Model).ActiveScope(); got != 0 {
		t.Fatalf("tab did not wrap: landed on %d", got)
	}
}

// A refresh that fails must leave the panel as it was. Showing one destination's
// counts under another's name would be a panel that lies about where install goes.
func TestModelKeepsStateWhenRefreshFails(t *testing.T) {
	t.Setenv(EnvTheme, "dark")

	rows := []TargetRow{
		{Name: "global", Info: "a", Configured: true, Active: true},
		{Name: "project", Info: "b", Configured: true},
	}
	m := NewModel("v0", demoPanel().Menu, rows, false).
		WithRefresh(func(int) ([]MenuItem, []TargetRow, error) {
			return nil, nil, errors.New("permission denied")
		})

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = next.(Model)

	if got := m.ActiveScope(); got != 0 {
		t.Fatalf("a failed refresh moved the active scope to %d", got)
	}
	if !strings.Contains(m.Notice(), "permission denied") {
		t.Errorf("the failure is not reported: %q", m.Notice())
	}
}

// With no refresh wired the key is inert rather than moving the mark to a
// destination whose figures nobody recomputed.
func TestModelScopeKeyIsInertWithoutRefresh(t *testing.T) {
	rows := []TargetRow{{Name: "global", Active: true}, {Name: "project"}}
	m := NewModel("v0", demoPanel().Menu, rows, false)

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if got := next.(Model).ActiveScope(); got != 0 {
		t.Fatalf("the mark moved to %d with no refresh wired", got)
	}
}

// tab changes the destination, not the cursor. Anybody watching the `❯` sees
// nothing move, so the key has to be listed and the switch has to announce itself
// — otherwise it gets reported as broken, which is exactly what happened.
func TestFooterListsTheScopeKey(t *testing.T) {
	t.Setenv(EnvTheme, "dark")
	out := NewTheme().Render(demoPanel())

	if !strings.Contains(out, "tab") {
		t.Error("the footer never mentions tab; an unlisted key that appears to do nothing reads as broken")
	}
}

func TestSwitchingScopeSaysSo(t *testing.T) {
	t.Setenv(EnvTheme, "dark")

	rows := []TargetRow{
		{Name: "global", Configured: true, Active: true},
		{Name: "project", Configured: false},
	}
	m := NewModel("v0", demoPanel().Menu, rows, false).
		WithRefresh(func(i int) ([]MenuItem, []TargetRow, error) {
			out := []TargetRow{{Name: "global"}, {Name: "project"}}
			out[i].Active = true
			return demoPanel().Menu, out, nil
		})

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if got := next.(Model).Notice(); !strings.Contains(got, "project") {
		t.Errorf("switching scope said %q; it has to name where it now acts", got)
	}
}

// Selection is carried by colour, not by a glyph — the project's own rule, and the
// menu's. Encoding it in the bullet put two meanings on one channel: a green ●
// (configured) read as "selected" more strongly than a gold ◉ ring, so the inactive
// destination looked active and the correct behaviour got reported as a bug.
func TestActiveDestinationIsGoldEndToEnd(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	rows := []TargetRow{
		{Name: "global", Info: "12 missing", Configured: true},
		{Name: "project", Info: "3 linked", Configured: true, Active: true},
	}
	p := demoPanel()
	p.Targets = rows

	rendered := strings.Split(theme.targets(p), "\n")
	if len(rendered) != 2 {
		t.Fatalf("rendered %d rows, want 2", len(rendered))
	}

	// One row, one colour — the whole active row is gold, bullet and figures too.
	if got := coloursIn(rendered[1]); len(got) != 1 || got[0] != rgbOf(theme.Gold) {
		t.Errorf("the active row is %v, want one colour rgb(%s)", got, rgbOf(theme.Gold))
	}
	if got := coloursIn(rendered[0]); slicesContain(got, rgbOf(theme.Gold)) {
		t.Errorf("an inactive row carries gold, so selection is ambiguous: %v", got)
	}

	// And the bullet is free to mean one thing only.
	if strings.Contains(theme.targets(p), "◉") {
		t.Error("the ring glyph is back; the bullet means configured-ness, nothing else")
	}
}

// Colour cannot be the only signal. Strip it and the rows must still be told apart,
// which is what a non-colour terminal and a colour-blind reader both get.
func TestActiveDestinationIsMarkedWithoutColour(t *testing.T) {
	theme := darkTheme()
	rows := []TargetRow{
		{Name: "global", Info: "12 missing", Configured: true},
		{Name: "project", Info: "3 linked", Configured: true, Active: true},
	}
	p := demoPanel()
	p.Targets = rows

	rendered := strings.Split(stripANSI(theme.targets(p)), "\n")
	if len(rendered) != 2 {
		t.Fatalf("rendered %d rows, want 2", len(rendered))
	}
	if strings.Contains(rendered[0], "❯") {
		t.Errorf("the inactive row carries the cursor: %q", rendered[0])
	}
	if !strings.Contains(rendered[1], "❯") {
		t.Errorf("with colour removed the active row is unmarked: %q", rendered[1])
	}
}

func slicesContain(hay []string, needle string) bool {
	for _, h := range hay {
		if h == needle {
			return true
		}
	}
	return false
}

// stripANSI removes colour so a test can ask what survives without it.
func stripANSI(s string) string {
	return escapes.ReplaceAllString(s, "")
}

// Gold is the only colour that may mean "selected". The inactive destination's
// bullet was green, and green says "on" loudly enough that it kept looking like the
// chosen row even with gold on the active one. Two colours arguing about selection
// is one colour too many.
func TestNoSecondColourCompetesWithSelection(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	rows := []TargetRow{
		{Name: "global", Info: "12 missing", Configured: true},
		{Name: "project", Info: "3 linked", Configured: true, Active: true},
	}
	p := demoPanel()
	p.Targets = rows

	rendered := strings.Split(theme.targets(p), "\n")
	green := rgbOf(theme.Green)

	for i, row := range rendered {
		if slicesContain(coloursIn(row), green) {
			t.Errorf("strip row %d carries green, which competes with gold for meaning: %q", i, row)
		}
	}
	// And the active row is still the only gold one.
	if got := coloursIn(rendered[1]); len(got) != 1 || got[0] != rgbOf(theme.Gold) {
		t.Errorf("the active row is %v, want only gold", got)
	}
}

// ── a destructive action is asked twice, in place ─────────────────────────────

func destructiveModel(t *testing.T, calls *[]bool) Model {
	t.Helper()
	menu := []MenuItem{
		{Label: "install", Desc: "d", Enabled: true},
		{Label: "prune", Desc: "d", Enabled: true, Destructive: true},
	}
	rows := []TargetRow{
		{Name: "global", Configured: true, Active: true},
		{Name: "project", Configured: true},
	}
	m := NewModel("v0", menu, rows, false).
		WithRefresh(func(i int) ([]MenuItem, []TargetRow, error) {
			out := []TargetRow{{Name: "global"}, {Name: "project"}}
			out[i].Active = true
			return menu, out, nil
		}).
		WithRunner(func(_ string, _ int, confirm bool) ([]string, error) {
			*calls = append(*calls, confirm)
			return []string{"  would remove  skills/gone"}, nil
		})
	m.panel.Selected = 1 // prune
	return m
}

func TestDestructiveActionNeedsASecondPress(t *testing.T) {
	var calls []bool
	m := destructiveModel(t, &calls)

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	after := next.(Model)
	if len(calls) != 1 || calls[0] {
		t.Fatalf("the first press ran with confirm=%v, want a dry run", calls)
	}
	if after.Armed() != "prune" {
		t.Fatalf("prune is not armed after showing its plan: %q", after.Armed())
	}
	if !strings.Contains(after.Notice(), "again") {
		t.Errorf("the notice does not ask for a second press: %q", after.Notice())
	}
	if len(after.Results()) == 0 {
		t.Error("the plan it wants confirming is not on screen")
	}

	next, _ = after.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if len(calls) != 2 || !calls[1] {
		t.Fatalf("the second press ran with confirm=%v, want it carried out", calls)
	}
	if got := next.(Model).Armed(); got != "" {
		t.Errorf("still armed after going ahead: %q", got)
	}
}

// Moving the cursor disarms. Otherwise a confirmation could be spent on a row you
// walked to afterwards.
func TestMovingTheCursorDisarms(t *testing.T) {
	var calls []bool
	m := destructiveModel(t, &calls)

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyUp})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(calls) != 2 || calls[1] {
		t.Fatalf("calls=%v — a confirmation survived moving the cursor", calls)
	}
}

// Switching destination disarms, or the confirmation would apply to a plan read for
// somewhere else. This is the one that would delete from the wrong place.
func TestSwitchingDestinationDisarms(t *testing.T) {
	var calls []bool
	m := destructiveModel(t, &calls)

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyTab})
	next, _ = next.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})

	if len(calls) != 2 || calls[1] {
		t.Fatalf("calls=%v — a confirmation followed the user to another destination", calls)
	}
}

// Nothing to remove is not something to confirm.
func TestNothingToRemoveIsNotArmed(t *testing.T) {
	menu := []MenuItem{{Label: "prune", Desc: "d", Enabled: true, Destructive: true}}
	m := NewModel("v0", menu, []TargetRow{{Name: "global", Active: true}}, false).
		WithRunner(func(string, int, bool) ([]string, error) { return nil, nil })

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	after := next.(Model)
	if after.Armed() != "" {
		t.Errorf("armed with an empty plan: %q", after.Armed())
	}
	if !strings.Contains(after.Notice(), "nothing") {
		t.Errorf("notice = %q, want it to say there was nothing to do", after.Notice())
	}
}

// A non-destructive action runs on the first press. Asking twice for everything
// teaches people to press twice for everything.
func TestNonDestructiveActionRunsAtOnce(t *testing.T) {
	var calls []bool
	m := destructiveModel(t, &calls)
	m.panel.Selected = 0 // install

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if len(calls) != 1 {
		t.Fatalf("calls=%v, want one", calls)
	}
	if got := next.(Model).Armed(); got != "" {
		t.Errorf("a non-destructive action armed itself: %q", got)
	}
	if !strings.Contains(next.(Model).Notice(), "done") {
		t.Errorf("notice = %q, want it to say the action finished", next.(Model).Notice())
	}
}

// A report line leads with its verb and its subject and trails off into paths, so an
// over-long one keeps its head. Keeping the tail instead ate the verb: `would remove
// skills/gone → /very/long/path` rendered as `…/001/skills/gone`, which is
// unreadable and looks like a different fact altogether.
func TestReportLinesKeepTheirHead(t *testing.T) {
	long := "  would remove  skills/gone → /private/var/folders/gn/n55rl6v17v3g4wzhvcl/T/x/001/skills/gone"
	got := elideRight(long, 40)

	if len([]rune(got)) > 40 {
		t.Fatalf("elided to %d columns, want at most 40: %q", len([]rune(got)), got)
	}
	if !strings.Contains(got, "would remove") {
		t.Errorf("the verb did not survive: %q", got)
	}
	if !strings.Contains(got, "skills/gone") {
		t.Errorf("the subject did not survive: %q", got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("the cut is not marked: %q", got)
	}
	if short := "  create  skills/alpha"; elideRight(short, 40) != short {
		t.Errorf("a line inside the budget was altered: %q", elideRight(short, 40))
	}
}
