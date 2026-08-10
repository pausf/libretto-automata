package ui

import (
	"math"
	"strconv"
	"testing"
)

// WCAG contrast floors. 4.5:1 is the threshold for normal-size text; 3:1 is
// allowed for large text and non-text graphics such as borders and rules.
const (
	minTextContrast   = 4.5
	minBorderContrast = 3.0
)

// relativeLuminance implements the WCAG 2.1 definition.
func relativeLuminance(hex string) float64 {
	channel := func(i int) float64 {
		v, err := strconv.ParseInt(hex[i:i+2], 16, 32)
		if err != nil {
			return 0
		}
		c := float64(v) / 255
		if c <= 0.03928 {
			return c / 12.92
		}
		return math.Pow((c+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(0) + 0.7152*channel(2) + 0.0722*channel(4)
}

func contrast(a, b string) float64 {
	la, lb := relativeLuminance(a), relativeLuminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// A palette that only looks good in a mockup is a bug. The first dark palette
// put borders at 1.4:1 and descriptions at 3.4:1, which read as "nothing is
// visible" on a real terminal. This test is why that cannot recur.
func TestPalettesAreReadable(t *testing.T) {
	cases := []struct {
		theme      string
		background string // the terminal background the palette targets
		palette    Theme
	}{
		{"dark", "1E1E1E", darkTheme()},
		{"dark on pure black", "000000", darkTheme()},
		{"light", "FFFFFF", lightTheme()},
		{"light on off-white", "F5F5F5", lightTheme()},
	}

	for _, c := range cases {
		t.Run(c.theme, func(t *testing.T) {
			text := map[string]string{
				"Parchment": c.palette.Parchment,
				"Gold":      c.palette.Gold,
				"Steel":     c.palette.Steel,
				"Muted":     c.palette.Muted,
				"Green":     c.palette.Green,
				"Error":     c.palette.Error,
			}
			for name, hex := range text {
				if got := contrast(hex, c.background); got < minTextContrast {
					t.Errorf("%s (#%s) on #%s = %.2f:1, want >= %.1f:1",
						name, hex, c.background, got, minTextContrast)
				}
			}

			// Borders, rules and inactive rows are non-text graphics: they may
			// recede, but they must still be perceivable.
			graphics := map[string]string{"Dim": c.palette.Dim, "Off": c.palette.Off}
			for name, hex := range graphics {
				if got := contrast(hex, c.background); got < minBorderContrast {
					t.Errorf("%s (#%s) on #%s = %.2f:1, want >= %.1f:1",
						name, hex, c.background, got, minBorderContrast)
				}
			}

			// The shading ramp fades toward the border, so only its brightest
			// step carries meaning and needs graphic-level contrast.
			if got := contrast(c.palette.Ramp['█'], c.background); got < minBorderContrast {
				t.Errorf("Ramp['█'] (#%s) on #%s = %.2f:1, want >= %.1f:1",
					c.palette.Ramp['█'], c.background, got, minBorderContrast)
			}
		})
	}
}

// Every step of the ramp must be distinguishable from its neighbour, or the fade
// reads as a flat block.
func TestRampStepsAreDistinct(t *testing.T) {
	for name, theme := range map[string]Theme{"dark": darkTheme(), "light": lightTheme()} {
		t.Run(name, func(t *testing.T) {
			steps := []rune{'░', '▒', '▓', '█'}
			for i := 1; i < len(steps); i++ {
				prev, cur := theme.Ramp[steps[i-1]], theme.Ramp[steps[i]]
				if got := contrast(prev, cur); got < 1.2 {
					t.Errorf("%q and %q are only %.2f:1 apart; the ramp reads flat",
						steps[i-1], steps[i], got)
				}
			}
		})
	}
}

func TestThemeEnvOverride(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"dark", darkTheme().Steel},
		{"light", lightTheme().Steel},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv(EnvTheme, tt.value)
			if got := NewTheme().Steel; got != tt.want {
				t.Errorf("NewTheme().Steel = %s, want %s", got, tt.want)
			}
		})
	}

	t.Run("an unknown value falls back to detection", func(t *testing.T) {
		t.Setenv(EnvTheme, "banana")
		got := NewTheme().Steel
		if got != darkTheme().Steel && got != lightTheme().Steel {
			t.Errorf("Steel = %s, which is neither palette", got)
		}
	})
}
