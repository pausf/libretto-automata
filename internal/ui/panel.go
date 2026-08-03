package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Content width bounds. The panel is fluid between them.
//
// The floor is the art's own width — narrower than that and there is nothing to
// draw. The ceiling exists because past roughly a hundred columns the layout gets
// worse, not better: a label at column 4 and its description forty columns away
// with a desert between them forces the eye to travel. Newspapers use narrow
// columns for the same reason.
const (
	MinContentWidth = ArtWidth
	MaxContentWidth = 98

	MinPanelWidth = MinContentWidth + 2 // plus one border column each side
	MaxPanelWidth = MaxContentWidth + 2
)

// Column offsets inside the panel's content area, taken from the art in
// docs/DESIGN.md. Menu labels and target names start at the same column in every
// row, and stay left-aligned at every width: a centred list is hard to read
// because the eye loses the edge it returns to.
const (
	menuLabelCol  = 6
	menuDescCol   = 18
	targetNameCol = 4
	targetInfoCol = 13
)

// separator marks a place where the frame draws ├───┤ rather than a content row.
const separator = "\x00separator"

// MenuItem is one action row.
type MenuItem struct {
	Label   string
	Desc    string
	Enabled bool // a disabled row recedes and refuses to run
}

// TargetRow is one line of the target strip.
type TargetRow struct {
	Name       string
	Info       string
	Configured bool

	// Active marks the scope the menu's actions will write to.
	//
	// The strip lists every destination; exactly one of them is where `install`
	// would go. A menu whose actions write somewhere the user cannot see is a menu
	// that surprises them, and a surprise about which config just changed is an
	// expensive one.
	Active bool
}

// Panel is everything the view needs to render.
type Panel struct {
	Version   string
	Menu      []MenuItem
	Selected  int
	Targets   []TargetRow
	Notice    string // one-line feedback under the panel
	Width     int    // terminal width; 0 means lay out at the minimum
	Height    int    // terminal height; 0 means do not centre vertically
	ASCIISafe bool
}

// ContentWidth is the width the panel's interior will use for a given terminal
// width, clamped between the bounds above.
func ContentWidth(termWidth int) int {
	if termWidth <= 0 {
		return MinContentWidth
	}
	switch w := termWidth - 2; {
	case w < MinContentWidth:
		return MinContentWidth
	case w > MaxContentWidth:
		return MaxContentWidth
	default:
		return w
	}
}

// Render composes the panel. Below the width the art needs, it degrades to the
// small mark and a plain title rather than tearing the borders.
func (t Theme) Render(p Panel) string {
	if p.Width > 0 && p.Width < MinPanelWidth {
		return t.renderNarrow(p)
	}
	cw := ContentWidth(p.Width)

	var rows []string
	rows = append(rows, strings.Split(t.Logo(cw, p.ASCIISafe), "\n")...)
	rows = append(rows, separator, "")
	rows = append(rows, strings.Split(t.menu(p), "\n")...)
	rows = append(rows, "", separator)
	rows = append(rows, strings.Split(t.targets(p), "\n")...)

	parts := []string{t.frame(rows, cw), t.footer(p, cw+2)}
	if p.Notice != "" {
		// Part of the block, not appended after it, so centring accounts for it
		// and the panel does not jump when a notice appears.
		parts = append(parts, Fg(t.Muted).Render("      "+p.Notice))
	}

	return t.centre(lipgloss.JoinVertical(lipgloss.Left, parts...), p)
}

// frame draws the box around the content rows at a given content width.
//
// lipgloss's Border() wraps a block in a uniform border and has no notion of an
// internal division, so it renders a section break as │───│. The design calls for
// ├───┤. Drawing the frame here is what buys those junctions; every content row
// is still measured and padded by lipgloss.
func (t Theme) frame(rows []string, width int) string {
	edge := Fg(t.Dim)
	rule := strings.Repeat("─", width)

	out := make([]string, 0, len(rows)+2)
	out = append(out, edge.Render("╭"+rule+"╮"))

	for _, row := range rows {
		if row == separator {
			out = append(out, edge.Render("├"+rule+"┤"))
			continue
		}
		side := edge.Render("│")
		out = append(out, side+padTo(row, width)+side)
	}

	out = append(out, edge.Render("╰"+rule+"╯"))
	return strings.Join(out, "\n")
}

// menu renders the action rows.
//
// One row, one colour. Every row reads in the ordinary text colour and the
// selected row turns gold end to end — cursor, bullet, label and description
// together. Splitting a row across two or three colours made the cursor compete
// with the label for attention; a single sweep of gold is unmistakable and needs
// no legend.
func (t Theme) menu(p Panel) string {
	rows := make([]string, len(p.Menu))
	for i, item := range p.Menu {
		colour, cursor := t.Steel, " "
		if i == p.Selected {
			colour, cursor = t.Gold, "❯"
		}

		line := cursor + " ▸ " + pad(item.Label, menuDescCol-menuLabelCol) + item.Desc
		rows[i] = "  " + Fg(colour).Render(line)
	}
	return strings.Join(rows, "\n")
}

func (t Theme) targets(p Panel) string {
	// Two facts about a destination, on two separate channels, because putting them
	// both on the bullet is what made the panel unreadable:
	//
	//   the bullet  ● configured · ○ not yet
	//   the colour  gold end to end = this is what the keys act on
	//
	// Selection is the menu's idiom, used verbatim: a `❯` cursor **and** a sweep of
	// gold. Both, not either.
	//
	// Colour alone is not enough — strip it and the rows become identical, which is
	// what a non-colour terminal and a colour-blind reader both get. Encoding it in
	// the *bullet* instead was worse still: a green ● (configured) reads as "on"
	// more strongly than a gold ◉ ring, which reads as an unticked radio button, so
	// the inactive destination looked selected and correct behaviour got reported as
	// a bug. The bullet now means one thing only.
	rows := make([]string, len(p.Targets))
	for i, tg := range p.Targets {
		bullet, cursor := "○", " "
		if tg.Configured {
			bullet = "●"
		}
		if tg.Active {
			cursor = "❯"
		}
		line := cursor + " " + bullet + " " + pad(tg.Name, targetInfoCol-targetNameCol) + tg.Info

		switch {
		case tg.Active:
			rows[i] = "  " + Fg(t.Gold).Render(line)
		case !tg.Configured:
			// Recedes end to end — the whole row, not just its bullet.
			rows[i] = "  " + Fg(t.Off).Render(line)
		default:
			// Achromatic. The bullet used to be green, and green says "on" loudly
			// enough that the inactive destination still looked like the chosen one
			// even with gold on the active row. Two colours arguing about selection
			// is one colour too many — the shape of the glyph carries configured-ness
			// and nothing carries selection but gold.
			rows[i] = "  " + Fg(t.Steel).Render(cursor+" "+bullet+" "+
				pad(tg.Name, targetInfoCol-targetNameCol)) +
				Fg(t.Muted).Render(tg.Info)
		}
	}
	return strings.Join(rows, "\n")
}

// footerIndent is the footer's margin on both sides. Six columns is what the art
// in docs/DESIGN.md uses, and mirroring it on the right keeps the line balanced
// under the box instead of running to its corner.
const footerIndent = 6

// footer sits outside the border: version left, key hints right.
//
// `tab` is listed because it has to be. It changes the destination rather than the
// cursor, so anybody watching the `❯` sees nothing move and concludes the key does
// nothing — an unlisted key that appears broken is worse than no key.
func (t Theme) footer(p Panel, width int) string {
	const hints = "↑↓ · ⏎ select · tab scope · q quit"

	left := strings.Repeat(" ", footerIndent) + p.Version
	gap := width - lipgloss.Width(left) - lipgloss.Width(hints) - footerIndent
	if gap < 1 {
		gap = 1
	}
	return Fg(t.Off).Render(left + strings.Repeat(" ", gap) + hints)
}

// renderNarrow is the degraded layout: the small mark beside a plain title, no
// box. A torn border is worse than no border.
func (t Theme) renderNarrow(p Panel) string {
	mark := Fg(t.Gold).Render(strings.Join(SmallMark, "\n"))
	title := lipgloss.JoinVertical(lipgloss.Left,
		Fg(t.Gold).Render("LIBRETTO"),
		Fg(t.Steel).Render("AUTOMATA"),
		Fg(t.Off).Render(p.Version),
	)
	head := lipgloss.JoinHorizontal(lipgloss.Top, mark, "  ", title)

	rows := make([]string, len(p.Menu))
	for i, item := range p.Menu {
		colour, cursor := t.Steel, "  "
		if i == p.Selected {
			colour, cursor = t.Gold, "❯ "
		}
		rows[i] = Fg(colour).Render(cursor + item.Label)
	}

	return t.centre(lipgloss.JoinVertical(lipgloss.Left, head, "", strings.Join(rows, "\n")), p)
}

// centre places the block in the terminal.
//
// Each axis is centred only when there is room to spare on it. With no known
// size — a pipe, or `preview` — the block stays flush left and unpadded, so
// piped output is never full of alignment whitespace.
func (t Theme) centre(block string, p Panel) string {
	if w := lipgloss.Width(block); p.Width > w {
		block = lipgloss.PlaceHorizontal(p.Width, lipgloss.Center, block)
	}
	if h := lipgloss.Height(block); p.Height > h {
		block = lipgloss.PlaceVertical(p.Height, lipgloss.Center, block)
	}
	return block
}

// pad right-pads a label to a column width, never truncating: a label longer
// than its column pushes its description right rather than losing characters.
// It always leaves at least one space, so label and description never collide.
func pad(s string, width int) string {
	if n := width - lipgloss.Width(s); n > 0 {
		return s + strings.Repeat(" ", n)
	}
	return s + " "
}

// padTo right-pads to exactly width, adding nothing when already there. Used for
// frame rows, where one extra space would push the border a column out.
func padTo(s string, width int) string {
	if n := width - lipgloss.Width(s); n > 0 {
		return s + strings.Repeat(" ", n)
	}
	return s
}
