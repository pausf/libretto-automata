package ui

import (
	"strings"
	"testing"
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
