package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var ansi = regexp.MustCompile("\x1b\\[[0-9;]*m")

func strip(s string) string { return ansi.ReplaceAllString(s, "") }

// forceTrueColor makes rendering deterministic. Without it lipgloss detects the
// test runner's pipe, drops to Ascii, and emits no escapes at all — which would
// let the colouring tests pass vacuously.
func forceTrueColor(t *testing.T) {
	t.Helper()
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })
}

func TestArtRowsGeometry(t *testing.T) {
	for _, safe := range []bool{false, true} {
		rows := ArtRows(safe)
		// The row count is not asserted: rows get added to the art. The width is
		// the invariant — one row off by a column tears the frame.
		if len(rows) == 0 {
			t.Fatalf("asciiSafe=%v: the art is empty", safe)
		}
		for i, r := range rows {
			if w := len([]rune(r)); w != ArtWidth {
				t.Errorf("asciiSafe=%v row %d: width %d, want %d", safe, i, w, ArtWidth)
			}
		}
	}
}

// The 𝄞 rule, made executable: no rune outside the BMP may ever reach a
// terminal. See docs/DESIGN.md.
func TestNoAstralRunesInArt(t *testing.T) {
	sources := map[string][]string{
		"logo":       ArtRows(false),
		"logo safe":  ArtRows(true),
		"small mark": SmallMark,
	}
	for name, rows := range sources {
		for i, r := range rows {
			for _, ch := range r {
				if ch > 0xFFFF {
					t.Errorf("%s row %d contains U+%04X, outside the BMP", name, i, ch)
				}
			}
		}
	}
}

func TestASCIISafeRemovesQuadrants(t *testing.T) {
	t.Run("default art uses quadrants", func(t *testing.T) {
		joined := strings.Join(ArtRows(false), "")
		if !strings.ContainsAny(joined, quadrantGlyphs) {
			t.Skip("art no longer uses quadrants; this test is obsolete")
		}
	})

	t.Run("safe art has none", func(t *testing.T) {
		joined := strings.Join(ArtRows(true), "")
		if strings.ContainsAny(joined, quadrantGlyphs) {
			t.Errorf("safe art still contains a quadrant glyph from %q", quadrantGlyphs)
		}
	})
}

// Colour must never alter geometry. Stripping every escape from a rendered row
// has to reproduce the source art byte for byte.
func TestColouringPreservesArt(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	for _, safe := range []bool{false, true} {
		want := ArtRows(safe)

		// Logo() wraps the art in rail, blank ... blank, rail.
		rows := strings.Split(theme.Logo(ArtWidth, safe), "\n")
		if len(rows) != len(want)+4 {
			t.Fatalf("asciiSafe=%v: got %d rows, want %d", safe, len(rows), len(want)+4)
		}
		got := rows[2 : 2+len(want)]

		for i := range want {
			if plain := strip(got[i]); plain != want[i] {
				t.Errorf("asciiSafe=%v row %d altered by colouring:\n got %q\nwant %q",
					safe, i, plain, want[i])
			}
		}
	}
}

func TestColouringActuallyEmitsColour(t *testing.T) {
	forceTrueColor(t)

	if out := darkTheme().Logo(ArtWidth, false); !strings.Contains(out, "\x1b[") {
		t.Fatal("Logo() emitted no escape sequences; the colour tests would be vacuous")
	}
}

func TestWordmarkGradientRunsAcrossTheWord(t *testing.T) {
	theme := darkTheme()

	first := theme.colourOf(rowWordmarkFirst, colWordmarkFirst, '█')
	last := theme.colourOf(rowWordmarkFirst, colWordmarkLast, '▄')

	if !strings.EqualFold(first, theme.Parchment) {
		t.Errorf("left edge = %s, want %s", first, theme.Parchment)
	}
	if !strings.EqualFold(last, theme.Gold) {
		t.Errorf("right edge = %s, want %s", last, theme.Gold)
	}
	if strings.EqualFold(first, last) {
		t.Error("gradient endpoints are identical; no gradient is being applied")
	}
}

// Rule 2 must beat rule 4: the same glyph reads gradient inside the wordmark box
// and gold outside it.
func TestWordmarkBoxWinsOverClefRule(t *testing.T) {
	theme := darkTheme()

	inside := theme.colourOf(rowWordmarkFirst, colWordmarkFirst+5, '█')
	outside := theme.colourOf(rowWordmarkFirst, colWordmarkFirst-5, '█')

	if strings.EqualFold(inside, outside) {
		t.Errorf("█ got %s both inside and outside the wordmark box", inside)
	}
	if !strings.EqualFold(outside, theme.Gold) {
		t.Errorf("clef █ = %s, want gold %s", outside, theme.Gold)
	}
}

func TestColourOfRules(t *testing.T) {
	theme := darkTheme()

	tests := []struct {
		name string
		row  int
		col  int
		ch   rune
		want string
	}{
		{"space is never coloured", 0, 0, ' ', ""},
		{"light shading", 0, 2, '░', theme.Ramp['░']},
		{"medium shading", 0, 3, '▒', theme.Ramp['▒']},
		{"staff line is structure", 4, 1, '─', theme.Dim},
		{"rail is structure", 0, 8, '═', theme.Dim},
		{"thin rule is structure", 7, 11, '▏', theme.Dim},
		{"clef quadrant is gold", 7, 4, '▙', theme.Gold},
		{"AUTOMATA is steel", rowAutomata, 13, 'A', theme.Steel},
		{"tagline is muted", rowTaglineFirst, 13, 't', theme.Muted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := theme.colourOf(tt.row, tt.col, tt.ch); !strings.EqualFold(got, tt.want) {
				t.Errorf("colourOf(%d,%d,%q) = %q, want %q", tt.row, tt.col, tt.ch, got, tt.want)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	tests := []struct {
		name     string
		from, to string
		t        float64
		want     string
	}{
		{"start", "F5E6C8", "E8B44A", 0, "F5E6C8"},
		{"end", "F5E6C8", "E8B44A", 1, "E8B44A"},
		{"midpoint", "000000", "FFFFFF", 0.5, "808080"},
		{"clamps below zero", "000000", "FFFFFF", -3, "000000"},
		{"clamps above one", "000000", "FFFFFF", 42, "FFFFFF"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lerp(tt.from, tt.to, tt.t); !strings.EqualFold(got, tt.want) {
				t.Errorf("Lerp(%s,%s,%v) = %s, want %s", tt.from, tt.to, tt.t, got, tt.want)
			}
		})
	}
}

func TestThemesCoverEveryLayer(t *testing.T) {
	for name, theme := range map[string]Theme{"dark": darkTheme(), "light": lightTheme()} {
		t.Run(name, func(t *testing.T) {
			fields := map[string]string{
				"Parchment": theme.Parchment, "Gold": theme.Gold,
				"Steel": theme.Steel, "Muted": theme.Muted,
				"Dim": theme.Dim, "Off": theme.Off,
				"Green": theme.Green, "Error": theme.Error,
			}
			for field, hex := range fields {
				if len(hex) != 6 {
					t.Errorf("%s = %q, want a 6-digit hex", field, hex)
				}
			}
			for _, glyph := range []rune{'░', '▒', '▓', '█'} {
				if _, ok := theme.Ramp[glyph]; !ok {
					t.Errorf("Ramp is missing %q", glyph)
				}
			}
		})
	}
}
