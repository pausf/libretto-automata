package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestContentWidthClamping(t *testing.T) {
	tests := []struct {
		name string
		term int
		want int
	}{
		{"unknown width falls back to the minimum", 0, MinContentWidth},
		{"a negative width is treated as unknown", -1, MinContentWidth},
		{"narrower than the art clamps up", 40, MinContentWidth},
		{"exactly at the minimum", MinPanelWidth, MinContentWidth},
		{"fluid in between", 80, 78},
		{"exactly at the maximum", MaxPanelWidth, MaxContentWidth},
		{"wider than the maximum clamps down", 250, MaxContentWidth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContentWidth(tt.term); got != tt.want {
				t.Errorf("ContentWidth(%d) = %d, want %d", tt.term, got, tt.want)
			}
		})
	}
}

// The rail stretches with the panel and keeps symmetric margins. The fixed-width
// rail had six leftover columns on the right; a fluid one has no excuse for that.
func TestRailIsFluidAndSymmetric(t *testing.T) {
	for _, width := range []int{MinContentWidth, 78, MaxContentWidth} {
		rail := Rail(width)

		if got := len([]rune(rail)); got != width {
			t.Errorf("Rail(%d) is %d runes wide", width, got)
		}

		left := rail[:strings.Index(rail, "░")]
		right := rail[strings.LastIndex(rail, "░")+len("░"):]
		if left != right {
			t.Errorf("Rail(%d) margins differ: left %q, right %q", width, left, right)
		}
	}
}

func TestRailSurvivesAbsurdlyNarrowWidths(t *testing.T) {
	// Not reachable through Render, which clamps first, but a panic here would be
	// a landmine for any future caller.
	for _, width := range []int{0, 1, railFixedCost - 1} {
		if got := Rail(width); got == "" {
			t.Errorf("Rail(%d) returned empty", width)
		}
	}
}

// The art is a drawing: it is centred at every width, never stretched.
func TestArtIsCentredNotStretched(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	for _, width := range []int{MinContentWidth, 78, MaxContentWidth} {
		rows := strings.Split(theme.Logo(width, false), "\n")
		art := rows[2 : 2+len(artRows)]

		wantIndent := (width - ArtWidth) / 2
		for i, row := range art {
			plain := strip(row)
			gotIndent := len(plain) - len(strings.TrimLeft(plain, " "))
			wantForRow := wantIndent + (len(artRows[i]) - len(strings.TrimLeft(artRows[i], " ")))

			if gotIndent != wantForRow {
				t.Errorf("width %d, art row %d: indent %d, want %d",
					width, i, gotIndent, wantForRow)
			}
			if strings.TrimRight(plain, " ") != strings.TrimRight(strings.Repeat(" ", wantIndent)+artRows[i], " ") {
				t.Errorf("width %d, art row %d was altered: %q", width, i, plain)
			}
		}
	}
}

// Every bordered row must match the panel width the terminal implies, at every
// width. One short row is a torn box.
func TestFrameIsFlushAtEveryWidth(t *testing.T) {
	forceTrueColor(t)

	for _, term := range []int{MinPanelWidth, 80, 100, 140, 250} {
		p := demoPanel()
		p.Width = term
		want := ContentWidth(term) + 2

		for i, line := range strings.Split(darkTheme().Render(p), "\n") {
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

// The update notice is a frame row like any other, so it is measured and padded like any
// other. A row that is one column short is a torn box, and a notice long enough to be
// worth reading is exactly the row most likely to be it.
func TestFrameHoldsWithUpdateNotice(t *testing.T) {
	forceTrueColor(t)

	for _, term := range []int{MinPanelWidth, 80, 140} {
		p := demoPanel()
		p.Width = term
		p.UpdateNotice = "v0.2.0 → v0.3.0 available · choose update"
		want := ContentWidth(term) + 2

		for i, line := range strings.Split(darkTheme().Render(p), "\n") {
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

// Past the ceiling the panel stops growing and starts centring instead.
func TestPanelStopsGrowingAtTheCeiling(t *testing.T) {
	forceTrueColor(t)

	p := demoPanel()
	p.Width = 200

	var top string
	for _, l := range strings.Split(strip(darkTheme().Render(p)), "\n") {
		if strings.Contains(l, "╭") {
			top = strings.TrimSpace(l)
			break
		}
	}
	if got := lipgloss.Width(top); got != MaxPanelWidth {
		t.Errorf("at 200 columns the panel is %d wide, want the %d ceiling", got, MaxPanelWidth)
	}
}
