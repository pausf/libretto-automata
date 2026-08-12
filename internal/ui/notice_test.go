package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

const sampleNotice = "v0.2.0 → v0.3.0 available · choose update"

// The banner goes above the logo, which is the first thing read. It used to sit between
// the menu and the destination strip, where it was correct, unmissable in a test and
// invisible in practice: in the same register as a target row, it read as one more line
// rather than as news.
//
// Not the footer, still — the footer drops p.Version when the legend and the version
// cannot both fit, so a notice there would vanish at 96 columns.
func TestUpdateNoticeRendersAboveTheLogo(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice

	lines := strings.Split(strip(darkTheme().Render(p)), "\n")

	notice, logo := -1, -1
	for i, l := range lines {
		switch {
		case strings.Contains(l, sampleNotice) && notice == -1:
			notice = i
		// The shading rail opens the logo block, and nothing above it is the logo.
		case strings.Contains(l, "░▒▓█") && logo == -1:
			logo = i
		}
	}

	if notice == -1 {
		t.Fatalf("the notice is not in the panel:\n%s", strings.Join(lines, "\n"))
	}
	if logo == -1 {
		t.Fatalf("the logo is not in the panel:\n%s", strings.Join(lines, "\n"))
	}
	if notice > logo {
		t.Errorf("notice at row %d, logo at row %d — want the notice first", notice, logo)
	}
}

// A bare gold line was already gold and already inside the frame, and it still read as a
// row. The box is what separates news from content: its own border, and a caption in the
// top edge saying what the box is for before the version numbers have to be parsed.
func TestUpdateNoticeBannerIsBoxedAndCaptioned(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice

	lines := strings.Split(strip(darkTheme().Render(p)), "\n")

	text := -1
	for i, l := range lines {
		if strings.Contains(l, sampleNotice) {
			text = i
			break
		}
	}
	if text == -1 || text == 0 || text+1 >= len(lines) {
		t.Fatalf("no room for a box around the notice:\n%s", strings.Join(lines, "\n"))
	}

	if !strings.Contains(lines[text-1], bannerCaption) {
		t.Errorf("the row above the notice carries no caption:\n%s", lines[text-1])
	}
	if !strings.Contains(lines[text-1], "╭") {
		t.Errorf("the notice has no box above it:\n%s", lines[text-1])
	}
	if !strings.Contains(lines[text+1], "╰") {
		t.Errorf("the notice has no box below it:\n%s", lines[text+1])
	}
	// The text sits inside its own border, not just between two rules.
	if !strings.Contains(lines[text], "│ "+sampleNotice) {
		t.Errorf("the notice is not inside the box:\n%s", lines[text])
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
//
// So it keeps the notice and it keeps the *order* — first, as in the wide layout — and
// drops only the box, which is the one thing this width has no room for.
func TestNarrowLayoutKeepsUpdateNotice(t *testing.T) {
	forceTrueColor(t)
	p := demoPanel()
	p.UpdateNotice = sampleNotice
	p.Width = MinPanelWidth - 10

	out := strip(darkTheme().Render(p))
	if !strings.Contains(out, sampleNotice) {
		t.Errorf("the narrow layout dropped the notice:\n%s", out)
	}

	lines := strings.Split(out, "\n")
	notice, mark := -1, -1
	for i, l := range lines {
		switch {
		case strings.Contains(l, sampleNotice) && notice == -1:
			notice = i
		case strings.Contains(l, "LIBRETTO") && mark == -1:
			mark = i
		}
	}
	if mark == -1 {
		t.Fatalf("no wordmark in the narrow layout:\n%s", out)
	}
	if notice > mark {
		t.Errorf("notice at row %d, wordmark at row %d — want the notice first", notice, mark)
	}
	if strings.Contains(lines[notice], "╭") || strings.Contains(lines[notice], "│") {
		t.Errorf("the narrow layout drew a box it has no room for:\n%s", lines[notice])
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
