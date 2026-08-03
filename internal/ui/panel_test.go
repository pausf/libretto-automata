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

func TestModelSelectingAnEnabledActionRuns(t *testing.T) {
	m := newDemoModel()
	m.panel.Selected = 2 // status, enabled

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if notice := next.(Model).Notice(); !strings.Contains(notice, "status") {
		t.Errorf("notice = %q, want it to mention the action", notice)
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
