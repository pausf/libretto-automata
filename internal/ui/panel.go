package ui

import (
	"fmt"
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

	// Destructive marks an action that removes things, and so must be asked twice:
	// the first press reports what it would do, the second carries it out.
	Destructive bool
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
	Version  string
	Menu     []MenuItem
	Selected int
	Targets  []TargetRow

	// Results is the last action's report, one line per row, shown inside the frame
	// below the strip. Empty until something has run.
	//
	// It lives in the panel rather than being printed after quitting so the loop
	// stays closed: install, read what happened, tab to the other destination,
	// install again — without relaunching between each step.
	Results []string

	// Confirm is the question a destructive action is waiting on. Rendered inside the
	// frame, in the attention colour, because a question you have to hunt for is a
	// question that gets answered by accident.
	Confirm string

	// The model selector, the panel's second screen. When InSelector is set the
	// menu is replaced rather than overlaid: two lists of rows on one frame, each
	// with its own cursor, is two things to misread at a glance.
	InSelector    bool
	Agents        []AgentRow
	AgentCursor   int
	ChoosingModel bool
	ModelChoices  []ModelChoice
	ModelCursor   int

	// The effort catalogue, opened by `e` over the same rows. A mode rather than a
	// third screen: a screen that navigates to another screen to change one line is a
	// screen with a hallway in it.
	ChoosingEffort bool
	EffortChoices  []EffortChoice
	EffortCursor   int

	// AgentTop is the first agent index the window shows. The window exists because
	// the list is unbounded and the terminal is not: 29 agents drew 29 rows, and
	// PlaceVertical centring a block taller than the screen pushed the wordmark, the
	// strip and the first rows off the top — so the rows you could not reach took the
	// ones you could with them.
	//
	// It lives in the panel rather than being derived while rendering, because the
	// renderer is pure and the cursor is what moves the window. A renderer that
	// scrolled would be a renderer with state.
	AgentTop int

	// UpdateNotice says a newer release exists, and what to do about it. Empty until the
	// check answers, and empty forever when it cannot.
	//
	// Deliberately not Notice: that field is action feedback, and the first `install`
	// overwrites it. The same overwrite once ate the selector's key legend — footer()
	// still carries the comment about it — and news that one keypress deletes is news
	// nobody finishes reading.
	//
	// Nor the footer, where the version already gets dropped when the legend and the
	// version cannot both fit. A notice that disappears at 96 columns was never read.
	UpdateNotice string

	Notice string // one-line feedback under the panel

	// Refused marks the notice as something the panel declined to do, rather than
	// something it did. A refusal is drawn in the error colour inside its own box; every
	// other notice stays the muted line it was.
	//
	// The distinction is the whole of why this field exists. Most notices are outcomes —
	// `install · done`, `acting on project`, `3 agent(s) → haiku`, the selector's key
	// legend — and painting those red would make a successful apply read as a failure,
	// which is worse than the grey it replaced. It is set through refuse() and say() so
	// the message and its kind cannot be set apart from each other and go stale.
	Refused bool

	Width     int // terminal width; 0 means lay out at the minimum
	Height    int // terminal height; 0 means do not centre vertically
	ASCIISafe bool
}

// MaxResultRows bounds the report shown in the panel.
//
// ponytail: a constant, not a function of the terminal height. Twelve items is the
// current payload and fits; a hundred would push the frame off a short screen, which
// is the same tearing the path budget exists to prevent. Anything past the cap is
// counted rather than dropped silently — a truncated list that does not say it was
// truncated is a list that lies.
const MaxResultRows = 14

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
	if p.UpdateNotice != "" {
		rows = append(rows, t.banner(p.UpdateNotice, cw)...)
	}
	rows = append(rows, strings.Split(t.Logo(cw, p.ASCIISafe), "\n")...)
	rows = append(rows, separator, "")
	if p.InSelector {
		rows = append(rows, strings.Split(t.selector(p), "\n")...)
	} else {
		rows = append(rows, strings.Split(t.menu(p), "\n")...)
	}
	rows = append(rows, "", separator)
	rows = append(rows, strings.Split(t.targets(p), "\n")...)

	if len(p.Results) > 0 {
		rows = append(rows, separator)
		rows = append(rows, t.results(p, cw)...)
	}
	if p.Confirm != "" {
		rows = append(rows, "", "  "+Fg(t.Gold).Render(p.Confirm))
	}

	parts := []string{t.frame(rows, cw), t.footer(p, cw+2)}
	if p.Notice != "" {
		// Part of the block, not appended after it, so centring accounts for it
		// and the panel does not jump when a notice appears.
		if p.Refused {
			parts = append(parts, t.refusal(p.Notice, cw))
		} else {
			parts = append(parts, Fg(t.Muted).Render("      "+p.Notice))
		}
	}

	return t.centre(lipgloss.JoinVertical(lipgloss.Left, parts...), p)
}

// refusal draws a notice the panel declined to act on: the error colour, in its own box,
// exactly as wide as the frame above it.
//
// Its own box rather than the muted line, because a refusal has to survive being glanced
// past. The messages it carries are the ones that answer "why did nothing happen?" — a
// keypress that refused, a marked set that cannot take the value — and as one grey line
// under a bordered panel they read as decoration.
//
// **Exactly cw+2, the frame's own outer width.** Not narrower and not wider: a box one
// column off under a box that is flush at every width reads as the frame having broken
// rather than as a second box. It borrows the frame's border glyphs for the same reason,
// which also means the existing flush tests measure this box for free — they filter rows
// on those characters, so a width mistake here fails a test that already exists.
//
// The text wraps rather than eliding. A refusal names the way out — *unmark those rows,
// or move them off it with m* — and the way out is the half an ellipsis would eat.
func (t Theme) refusal(msg string, cw int) string {
	edge := Fg(t.Error)

	// padTo is what the frame pads its rows with, so the two boxes agree by construction
	// rather than by two pieces of arithmetic that happen to match. Two columns of that
	// width are the indent that keeps the text off the border.
	inner := cw - 4
	if inner < 1 {
		inner = 1
	}

	body := lipgloss.NewStyle().Width(inner).Foreground(lipgloss.Color("#" + t.Error)).Render(msg)

	rule := strings.Repeat("─", cw)
	rows := []string{edge.Render("╭" + rule + "╮")}
	for _, line := range strings.Split(body, "\n") {
		side := edge.Render("│")
		rows = append(rows, side+padTo("  "+line, cw)+side)
	}
	rows = append(rows, edge.Render("╰"+rule+"╯"))
	return strings.Join(rows, "\n")
}

// bannerCaption names the box before its contents have to be parsed. English, like every
// other string on the panel.
const bannerCaption = "NEW VERSION"

// bannerMargin is the gap between the banner's box and the frame it sits inside, on each
// side. Two columns, the same indent every content row already uses.
const bannerMargin = 2

// banner draws the update notice as its own box at the top of the frame.
//
// A box rather than a filled tag, and that is the whole of why this is not one gold line
// any more. The line was correct — inside the frame, in the attention colour, elided so it
// could not tear anything — and it still read as a row, because below the menu it sat in
// the same register as a target. Weight had to come from somewhere, and the two candidates
// were a background colour or a border.
//
// The border wins on a constraint rather than on taste: Fg is the only helper the theme
// has, and every value in both palettes was measured at 4.5:1 as a *foreground*. A filled
// tag needs an ink/background pair measured again, in two palettes, for the loudest
// element on the screen — and the palette that shipped at 1.4:1 is why that is not done in
// passing.
//
// ponytail: no new palette values. The day somebody wants the filled sticker, it is a
// Tag(ink, bg) helper and two colours, each measured against its terminal background
// before it ships — not a Background() call bolted onto Fg.
func (t Theme) banner(text string, cw int) []string {
	gold := Fg(t.Gold)
	box := cw - 2*bannerMargin
	pad := strings.Repeat(" ", bannerMargin)

	// The caption lives in the top edge, costing its own length plus the corner, the
	// lead-in rule and a space either side.
	//
	// ponytail: no fallback for a box too narrow to carry it. `box` is `ContentWidth`
	// minus the margins, and `ContentWidth` floors at `MinContentWidth` — which is
	// `ArtWidth`, 58 — so the narrowest box this can be handed is 54 against a caption
	// needing 16. Anything narrower went to `renderNarrow`, which draws no box at all.
	// The guard that was here was unreachable by construction, and unreachable code that
	// looks defensive is worse than none: it reads as a case somebody handled.
	// Drop `MinContentWidth` below 20 and `TestBannerBoxHoldsAtEveryWidth` is what says so.
	top := "╭─ " + bannerCaption + " " +
		strings.Repeat("─", box-len(bannerCaption)-5) + "╮"

	return []string{
		pad + gold.Render(top),
		pad + gold.Render("│ "+padTo(elideRight(text, box-4), box-4)+" │"),
		pad + gold.Render("╰"+strings.Repeat("─", box-2)+"╯"),
	}
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

// results renders the last action's report inside the frame.
//
// Every line recedes to the muted colour: this is a record of what happened, not
// something to act on, and painting it like the menu would make the panel look as
// though it had two sets of controls.
//
// Lines longer than the content area keep their **head** and mark the cut at the end.
// A report line leads with its verb and its subject — `create skills/alpha`, `would
// remove skills/gone` — and trails off into paths. Keeping the tail instead ate the
// verb and left only a directory, so `would remove skills/gone → /very/long/path`
// rendered as `…/001/skills/gone`: unreadable, and it looked like a different fact.
//
// Roots in the strip use the opposite rule, and correctly: there the tail is what
// identifies the directory. Different content, different end to keep.
func (t Theme) results(p Panel, width int) []string {
	shown := p.Results
	extra := 0
	if len(shown) > MaxResultRows {
		extra = len(shown) - MaxResultRows
		shown = shown[:MaxResultRows]
	}

	out := make([]string, 0, len(shown)+1)
	for _, line := range shown {
		out = append(out, "  "+Fg(t.Muted).Render(elideRight(line, width-4)))
	}
	if extra > 0 {
		out = append(out, "  "+Fg(t.Off).Render(fmt.Sprintf("… and %d more", extra)))
	}
	return out
}

// elideRight trims a line to n columns by dropping the tail, marking the cut.
func elideRight(s string, n int) string {
	r := []rune(s)
	if n < 2 || len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
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
	// The hints belong to the screen, not to the program. A legend that lists `⏎
	// select` while ⏎ opens a catalogue is worse than no legend: it is read once,
	// believed, and never re-read. The selector's real keys used to live in the
	// opening notice, where the first apply overwrote them and they were gone for the
	// rest of the session.
	hints, tight := "↑↓ · ⏎ select · tab scope · q quit", "⏎ select · tab · q"
	switch {
	case p.Confirm != "":
		// While a question is open the only keys that matter are its answers.
		// Listing the others invites pressing one by reflex.
		hints, tight = "y yes · n no", "y yes · n no"
	case p.choosing():
		hints, tight = "↑↓ · ⏎ apply · esc back", "⏎ apply · esc"
	case p.InSelector:
		// `e effort` earns its place in both variants: a key that appears to do
		// nothing reads as broken, which is the whole reason this legend exists. What
		// gives way instead is the arrows in the tight form — they are the one hint a
		// user does not need told.
		hints, tight = "↑↓ · space mark · a all · m model · e effort · tab scope · esc back",
			"space · a all · m model · e effort · esc"
	}

	left := strings.Repeat(" ", footerIndent) + p.Version
	gap := width - lipgloss.Width(left) - lipgloss.Width(hints) - footerIndent

	// The selector's legend is half again as long as the menu's, and a footer wider
	// than the frame drags the centred block off the terminal — the same tearing the
	// fluid frame exists to prevent, arriving from underneath it.
	//
	// ponytail: one fallback string per screen, not a truncation ladder. What survives
	// the drop is the keys that exist only on this screen; ↑↓ and `tab` are taught on
	// the menu, which is the only way in.
	if gap < 1 {
		hints = tight
		gap = width - lipgloss.Width(left) - lipgloss.Width(hints) - footerIndent
	}

	// A tag name has no length limit, so no pair of legends can be short enough on its
	// own — `v0.10.0-17-g96c04e3-dirty` is two tags away and already one column over.
	// The legend is what the panel is operated with and the version is reference that
	// `libretto version` still prints, so the version is what goes.
	if gap < 1 {
		left = strings.Repeat(" ", footerIndent)
		gap = width - lipgloss.Width(left) - lipgloss.Width(hints) - footerIndent
	}
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

	var parts []string
	if p.UpdateNotice != "" {
		// The narrow layout drops the border, not the content. Degrading what the panel
		// is telling you as well as how it is drawn is two losses for one terminal width.
		//
		// It keeps the order too — first, as in the wide layout. What it drops is the box,
		// which is the one part of the banner this width genuinely has no room for.
		parts = append(parts, Fg(t.Gold).Render(p.UpdateNotice), "")
	}
	parts = append(parts, head, "", strings.Join(rows, "\n"))
	return t.centre(lipgloss.JoinVertical(lipgloss.Left, parts...), p)
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
