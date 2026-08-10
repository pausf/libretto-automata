package ui

import "strings"

// ArtWidth is the fixed rune width of the clef-and-wordmark block. It is a
// drawing, so it has an intrinsic size and is never stretched — only centred.
// Everything around it (rails, rules, borders) is fluid.
const ArtWidth = 58

// artRows is the clef, wordmark and subtitle from docs/DESIGN.md, borders and
// rails stripped. Extracted from that file, not retyped. Do not edit by hand —
// edit DESIGN.md and re-extract, or the doc and the binary drift apart.
var artRows = []string{
	"    ▄▀▀▄                                                  ",
	"   ▐▌  ▐▌    █    ▀█▀  █▀▄  █▀▄  █▀▀  ▀█▀  ▀█▀  ▄▀▄       ",
	" ──█▄▄▄▀──   █     █   █▀▄  █▀▄  █▀    █    █   █ █  ──   ",
	" ──█▀▀▀▄──   █▄▄  ▄█▄  █▄▀  █ ▀  █▄▄   █    █   ▀▄▀  ──   ",
	"   ▐▌  ▐▌                                                 ",
	"   ▐▙▄▄▟▘  ▏ A U T O M A T A                              ",
	"    ▐▌     ▏ the libretto is written first ·              ",
	"   ▄▀      ▏ the automaton performs it                    ",
	"           ▏ b y   p a u s f                              ",
}

// SmallMark is the compact clef, for one-line output and spinners.
var SmallMark = []string{
	" ▄▀",
	"▄█▀",
	"▀▄ ",
}

// Row indices into artRows, and the column span the wordmark occupies. Columns
// are relative to the art block, not to the panel, so centring never has to be
// accounted for here.
const (
	rowWordmarkFirst = 1
	rowWordmarkLast  = 3
	colWordmarkFirst = 13
	colWordmarkLast  = 50

	rowAutomata     = 5
	rowTaglineFirst = 6
	rowTaglineLast  = 7

	rowSignature = 8
)

// railFixedCost is the width the rail spends on its own margins and shading
// blocks: "  ░▒▓█ " plus " █▓▒░  ".
const railFixedCost = 14

// clefGlyphs are the block elements the clef is drawn from.
const clefGlyphs = "▄▀█▐▌▙▟▘"

// structureGlyphs are borders, staff lines, rails and thin rules. They share one
// colour on purpose, so a staff line running into a border reads as a single
// continuous stroke.
const structureGlyphs = "╭╮╰╯├┤│─═▏"

// quadrantGlyphs have very high but not universal font coverage. See asciiSafe.
const quadrantGlyphs = "▙▟▘"

// safeReplacements substitute quadrants with half blocks when the terminal font
// may lack them. The clef looks squarer and never renders as tofu.
var safeReplacements = strings.NewReplacer("▙", "█", "▟", "█", "▘", "▀")

// ArtRows returns the art block, applying the quadrant fallback when asked.
func ArtRows(asciiSafe bool) []string {
	rows := make([]string, len(artRows))
	for i, r := range artRows {
		if asciiSafe {
			r = safeReplacements.Replace(r)
		}
		rows[i] = r
	}
	return rows
}

// Rail returns the shading rail at a given content width, uncoloured.
//
// The margins are symmetric: the rail is fluid, so there is no reason for it to
// sit off-centre the way a fixed-width rail with leftover space did.
func Rail(width int) string {
	n := width - railFixedCost
	if n < 1 {
		n = 1
	}
	return "  ░▒▓█ " + strings.Repeat("═", n) + " █▓▒░  "
}

// Logo renders the coloured logo block — rail, art, rail — at a content width.
//
// The art is centred within the width; the rails span it.
func (t Theme) Logo(width int, asciiSafe bool) string {
	if width < ArtWidth {
		width = ArtWidth
	}
	indent := strings.Repeat(" ", (width-ArtWidth)/2)

	lines := []string{t.rail(width), ""}
	for i, row := range ArtRows(asciiSafe) {
		lines = append(lines, indent+t.colourRow(i, row))
	}
	return strings.Join(append(lines, "", t.rail(width)), "\n")
}

// rail colours the shading rail: the ramp on the block glyphs, dim on the ═ run.
func (t Theme) rail(width int) string {
	var b strings.Builder
	for _, ch := range Rail(width) {
		switch {
		case ch == ' ':
			b.WriteByte(' ')
		case ch == '═':
			b.WriteString(Fg(t.Dim).Render("═"))
		default:
			b.WriteString(Fg(t.Ramp[ch]).Render(string(ch)))
		}
	}
	return b.String()
}

// colourRow colours one art row, coalescing runs of the same colour into a
// single escape sequence.
func (t Theme) colourRow(row int, line string) string {
	// the signature is the one italic row, and lipgloss resets styles between
	// segments, so it is rendered whole instead of run by run
	if row == rowSignature {
		i := strings.IndexRune(line, '▏')
		return line[:i] + Fg(t.Dim).Render("▏") +
			Fg(t.Gold).Italic(true).Render(line[i+len("▏"):])
	}

	var b strings.Builder
	runs := []rune(line)
	start, current := 0, ""
	flush := func(end int) {
		if end <= start {
			return
		}
		text := string(runs[start:end])
		if current == "" {
			b.WriteString(text)
		} else {
			b.WriteString(Fg(current).Render(text))
		}
	}

	for col, ch := range runs {
		want := t.colourOf(row, col, ch)
		if want != current {
			flush(col)
			start, current = col, want
		}
	}
	flush(len(runs))
	return b.String()
}

// colourOf implements docs/DESIGN.md's colouring rules, in order. First match
// wins. An empty return means "leave uncoloured".
func (t Theme) colourOf(row, col int, ch rune) string {
	// 1. spaces are never coloured
	if ch == ' ' {
		return ""
	}

	// 2. the wordmark box wins over every glyph rule, which is why the
	//    wordmark's █ reads gradient while the clef's █ reads gold
	if row >= rowWordmarkFirst && row <= rowWordmarkLast &&
		col >= colWordmarkFirst && col <= colWordmarkLast {
		span := float64(colWordmarkLast - colWordmarkFirst)
		return Lerp(t.Parchment, t.Gold, float64(col-colWordmarkFirst)/span)
	}

	// 3. the shading rail ramp
	if c, ok := t.Ramp[ch]; ok {
		return c
	}

	// 4. the clef
	if strings.ContainsRune(clefGlyphs, ch) {
		return t.Gold
	}

	// 5. structure
	if strings.ContainsRune(structureGlyphs, ch) {
		return t.Dim
	}

	// 6. row-specific text
	switch {
	case row == rowAutomata:
		return t.Steel
	case row >= rowTaglineFirst && row <= rowTaglineLast:
		return t.Muted
	}
	return ""
}
