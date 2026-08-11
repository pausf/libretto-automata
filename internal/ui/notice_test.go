package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

const sampleNotice = "v0.2.0 → v0.3.0 available · choose update"

// The row goes between the menu and the destination strip. Not in the footer: the footer
// already drops p.Version when the legend and the version cannot both fit, so a notice
// there would vanish at 96 columns. Not on a target row either — the strip is about where
// links go.
func TestPanelRendersUpdateNoticeBetweenMenuAndStrip(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice

	lines := strings.Split(strip(darkTheme().Render(p)), "\n")

	notice, menu, targets := -1, -1, -1
	for i, l := range lines {
		switch {
		case strings.Contains(l, sampleNotice):
			notice = i
		case strings.Contains(l, "status") && menu == -1:
			menu = i
		// `codex`, not `claude`: the install row's description contains `~/.claude` and
		// would match eight rows above the strip.
		case strings.Contains(l, "codex") && targets == -1:
			targets = i
		}
	}

	if notice == -1 {
		t.Fatalf("the notice is not in the panel:\n%s", strings.Join(lines, "\n"))
	}
	if !(menu < notice && notice < targets) {
		t.Errorf("notice at row %d, want between the menu (%d) and the strip (%d)",
			notice, menu, targets)
	}
}

func TestPanelOmitsUpdateNoticeWhenEmpty(t *testing.T) {
	forceTrueColor(t)
	before := darkTheme().Render(demoPanel())

	p := demoPanel()
	p.UpdateNotice = ""
	if got := darkTheme().Render(p); got != before {
		t.Error("an empty notice changed the panel")
	}
}

// The narrow layout drops the border, not the content. A degraded layout that also
// degrades what it is telling you is two losses for one terminal width.
func TestNarrowLayoutKeepsUpdateNotice(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice
	p.Width = MinPanelWidth - 10

	out := strip(darkTheme().Render(p))
	if !strings.Contains(out, sampleNotice) {
		t.Errorf("the narrow layout dropped the notice:\n%s", out)
	}
}

// The check runs as a command, so the panel paints complete without it and re-renders when
// it lands. A panel that waits on the network before its first frame hangs on bad DNS, and
// the user's only recourse is ⌃C on a tool that looks broken.
func TestInitReturnsReleaseCheckCommand(t *testing.T) {
	m := NewModel("v0.2.0", nil, nil, false).
		WithReleaseCheck(func() string { return sampleNotice })

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init returned no command with a check configured")
	}

	msg, ok := cmd().(releaseMsg)
	if !ok {
		t.Fatalf("the command produced %T, want releaseMsg", cmd())
	}
	if string(msg) != sampleNotice {
		t.Errorf("releaseMsg = %q, want %q", string(msg), sampleNotice)
	}
}

// No check configured is the ordinary case for `preview` and for every test that does not
// care. It must produce no command and no notice — not an empty row, and not a nil-call
// panic.
func TestNoReleaseCheckMeansNoCommandAndNoNotice(t *testing.T) {
	m := NewModel("v0.2.0", nil, nil, false)

	if cmd := m.Init(); cmd != nil {
		t.Error("Init returned a command with no check configured")
	}
	if m.panel.UpdateNotice != "" {
		t.Errorf("UpdateNotice = %q with no check", m.panel.UpdateNotice)
	}
}

func TestUpdateNoticeSetFromMessage(t *testing.T) {
	m := NewModel("v0.2.0", nil, nil, false)

	next, _ := m.Update(releaseMsg(sampleNotice))
	if got := next.(Model).panel.UpdateNotice; got != sampleNotice {
		t.Errorf("UpdateNotice = %q, want %q", got, sampleNotice)
	}

	// An empty answer — up to date, or the check failed — sets nothing.
	blank, _ := NewModel("v0.2.0", nil, nil, false).Update(releaseMsg(""))
	if got := blank.(Model).panel.UpdateNotice; got != "" {
		t.Errorf("an empty releaseMsg set %q", got)
	}
}

// The model-level half of the two-fields decision: running an action writes feedback and
// leaves the news alone. Panel.Notice is overwritten by every action, and that overwrite
// once ate the selector's key legend.
func TestActionFeedbackDoesNotOverwriteUpdateNotice(t *testing.T) {
	m := NewModel("v0.2.0", []MenuItem{{Label: "status", Enabled: true}}, nil, false).
		WithRunner(func(string, int, bool) ([]string, error) { return []string{"12 linked"}, nil })

	withNews, _ := m.Update(releaseMsg(sampleNotice))
	afterAction, _ := withNews.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})

	final := afterAction.(Model)
	if final.panel.UpdateNotice != sampleNotice {
		t.Errorf("the action cleared the update notice: %q", final.panel.UpdateNotice)
	}
	if final.notice == "" {
		t.Error("the action produced no feedback of its own")
	}
}

// Moving the cursor clears action feedback — that is existing behaviour — and must not
// take the news with it.
func TestNavigationDoesNotClearUpdateNotice(t *testing.T) {
	m := NewModel("v0.2.0", []MenuItem{{Label: "a"}, {Label: "b"}}, nil, false)

	withNews, _ := m.Update(releaseMsg(sampleNotice))
	moved, _ := withNews.(Model).Update(tea.KeyMsg{Type: tea.KeyDown})

	if got := moved.(Model).panel.UpdateNotice; got != sampleNotice {
		t.Errorf("moving the cursor cleared the notice: %q", got)
	}
}

// The notice is news, not feedback, so it does not live in Panel.Notice — which the first
// install overwrites. That overwrite once ate the selector's key legend, and footer()
// still carries the comment about it.
func TestUpdateNoticeAndActionFeedbackCoexist(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice
	p.Notice = "install · done"

	out := strip(darkTheme().Render(p))
	if !strings.Contains(out, sampleNotice) {
		t.Error("action feedback displaced the update notice")
	}
	if !strings.Contains(out, "install · done") {
		t.Error("the update notice displaced the action feedback")
	}
}
