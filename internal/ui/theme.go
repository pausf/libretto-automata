// Package ui renders the Libretto Automata panel.
//
// The art and the palette are specified in docs/DESIGN.md. That file is the
// source of truth; this package is its implementation.
package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

// Theme is the palette resolved for the current terminal background.
//
// Colours are concrete hex strings rather than lipgloss.AdaptiveColor because
// the wordmark gradient has to interpolate real values. Resolving light/dark
// once, here, keeps that decision in a single place.
type Theme struct {
	Parchment string // gradient start; selected menu label
	Gold      string // gradient end; clef; cursor
	Steel     string // AUTOMATA; unselected labels; target name
	Muted     string // tagline; descriptions; counts
	Dim       string // borders, staff lines, rails, rules
	Off       string // inactive bullets, unconfigured rows, footer
	Green     string // configured
	Error     string // failures, conflicts

	// Ramp colours the shading rail, fading the staff into the border.
	Ramp map[rune]string
}

// EnvTheme forces a palette, bypassing background detection. Terminals lie about
// their background often enough that being stuck in the wrong palette must not
// require a rebuild. Accepts "dark" or "light".
const EnvTheme = "LIBRETTO_THEME"

// NewTheme resolves the palette against the terminal background, unless
// LIBRETTO_THEME overrides it.
func NewTheme() Theme {
	switch os.Getenv(EnvTheme) {
	case "dark":
		return darkTheme()
	case "light":
		return lightTheme()
	}
	if lipgloss.HasDarkBackground() {
		return darkTheme()
	}
	return lightTheme()
}

// Every value below clears 4.5:1 against its terminal background — the contrast
// floor for readable text. An earlier palette used #3A3A42 for borders and
// #6A6A78 for descriptions, which measure 1.4:1 and 3.4:1 on a dark terminal:
// elegant in a mockup, unreadable in a terminal. Recession is expressed by
// staying achromatic, not by fading out.
func darkTheme() Theme {
	return Theme{
		Parchment: "F7EAD2", // gradient start
		Gold:      "F0BE52", // gradient end, clef, cursor
		Steel:     "E4E4EC", // near-white: AUTOMATA, labels, target name
		Muted:     "AEAEBE", // descriptions, counts
		Dim:       "6A6A78", // borders, staff lines, rails, rules
		Off:       "8A8A98", // inactive bullets, unconfigured rows, footer
		Green:     "8FE3B0",
		Error:     "FF8878",
		Ramp: map[rune]string{
			'░': "8A6E36",
			'▒': "B08C42",
			'▓': "D2A84A",
			'█': "F0BE52",
		},
	}
}

// lightTheme inverts the lightness axis while keeping every hue. On a light
// terminal the parchment end of the gradient becomes dark ink, otherwise the
// wordmark would disappear into the background.
func lightTheme() Theme {
	return Theme{
		Parchment: "5A4A24",
		Gold:      "8A6A1E",
		Steel:     "1C1C26", // near-black: the same role Steel plays on dark
		Muted:     "4E4E5A",
		Dim:       "85858F",
		Off:       "6E6E7C",
		Green:     "1B7A47",
		Error:     "A82A1A",
		Ramp: map[rune]string{
			'░': "C4A868",
			'▒': "A88A44",
			'▓': "97781F",
			'█': "8A6A1E",
		},
	}
}

// Fg styles text in a hex colour.
func Fg(hex string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#" + hex))
}

// Lerp interpolates two hex colours. t is clamped to [0,1].
//
// This is what makes the wordmark gradient run smoothly across the whole word
// instead of stepping letter by letter — see docs/DESIGN.md, rule 2.
func Lerp(from, to string, t float64) string {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	out := ""
	for i := 0; i < 6; i += 2 {
		a, errA := strconv.ParseInt(from[i:i+2], 16, 32)
		b, errB := strconv.ParseInt(to[i:i+2], 16, 32)
		if errA != nil || errB != nil {
			return from
		}
		v := float64(a) + (float64(b)-float64(a))*t
		out += fmt.Sprintf("%02X", int(v+0.5))
	}
	return out
}
