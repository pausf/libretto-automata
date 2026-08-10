package ui

import (
	"sort"
	"strconv"
	"strings"
)

// The model selector: the panel's second screen.
//
// Everything here obeys the rule the rest of the package already lives by — no
// filesystem, no knowledge of what an agent file is. Rows come in through a
// callback and the chosen model goes out through another, exactly as Refresh and
// Runner already do. That is what keeps the panel testable by calling Update.

// modelsAction is the menu label that opens the selector rather than running an
// action. Named once, so the menu and the dispatch cannot drift apart.
const modelsAction = "models"

// allRow is the cursor position of the `all` row, which sits above every group.
//
// The consequence is that screen position and agent index stop being the same number:
// agent i is at cursor i+1. Every read of p.Agents from the cursor goes through
// AgentCursor-1, and getting that wrong marks the neighbouring agent silently — which
// is why it has a test of its own.
const allRow = 0

// selectorRows is how many rows the cursor can reach: the agents, plus `all`.
func selectorRows(p Panel) int { return len(p.Agents) + 1 }

// AgentRow is one agent in the selector: what it is called, what it runs on, and
// whether the user has marked it.
type AgentRow struct {
	Name   string
	Model  string // empty means the session's model
	Marked bool

	// Shared marks a row whose file is one this repository owns, reached from more
	// than one destination. Writing it changes every project on the machine;
	// writing an unmarked row changes this destination only.
	//
	// The caller decides what sets it — this package knows nothing about symlinks,
	// and the flag is deliberately named for the consequence rather than for the
	// mechanism, because the consequence is what the reader has to act on.
	Shared bool
}

// ModelChoice is one entry of the catalogue the `m` key opens.
//
// The catalogue is supplied rather than known. This package deciding which models
// exist would put the list in two places, and the copy nobody edited would be the
// one on screen.
type ModelChoice struct {
	Name  string // empty is the session's model — declaring nothing
	Label string
}

// ListAgents returns the agents of one destination and their current models.
//
// The destination index is passed in, never captured when the panel opened — the
// same rule Runner follows, for the same reason: a closure over the starting scope
// would list one destination's agents while the strip named another, and the write
// that followed would land somewhere the user could not see.
type ListAgents func(destination int) ([]AgentRow, error)

// ApplyModel declares one model on a set of agents in one destination.
type ApplyModel func(destination int, names []string, model string) error

// WithAgents wires the selector. Without it the menu row refuses, the same way an
// unwired action does.
func (m Model) WithAgents(choices []ModelChoice, list ListAgents, apply ApplyModel) Model {
	m.modelChoices = choices
	m.listAgents = list
	m.applyModel = apply
	return m
}

// InSelector reports whether the second screen is up.
func (m Model) InSelector() bool { return m.panel.InSelector }

// ChoosingModel reports whether the catalogue is open over the selector.
func (m Model) ChoosingModel() bool { return m.panel.ChoosingModel }

// AgentRows exposes the rows for tests.
func (m Model) AgentRows() []AgentRow { return m.panel.Agents }

// ModelChoices exposes the catalogue for tests.
func (m Model) ModelChoices() []ModelChoice { return m.modelChoices }

// Done reports whether the panel is quitting, so a test can tell "left the
// selector" from "left the program".
func (m Model) Done() bool { return m.done }

// ChosenModelName is the catalogue entry under the cursor.
func (m Model) ChosenModelName() string {
	if len(m.modelChoices) == 0 {
		return ""
	}
	return m.modelChoices[m.panel.ModelCursor].Name
}

// MarkedAgents lists the marked rows, in screen order.
func (m Model) MarkedAgents() []string {
	var out []string
	for _, a := range m.panel.Agents {
		if a.Marked {
			out = append(out, a.Name)
		}
	}
	return out
}

// openSelector loads the rows and shows the second screen.
func (m Model) openSelector() Model {
	if m.listAgents == nil || m.applyModel == nil {
		m.notice = "models is not wired up yet"
		return m
	}

	rows, err := m.listAgents(m.activeScope())
	if err != nil {
		m.notice = "cannot read the agents: " + err.Error()
		return m
	}

	m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
	m.panel.ModelChoices = m.modelChoices
	m.panel.InSelector, m.panel.ChoosingModel = true, false
	m.panel.AgentCursor, m.panel.ModelCursor = 0, 0
	m.notice = "space mark · a all · m model · esc back"
	return m
}

// updateSelector handles the second screen's keys.
//
// esc closes the catalogue if it is open, and otherwise returns to the menu. It
// never quits: sharing one key between "go back" and "exit the program" is how
// somebody loses the panel trying to back out of a screen.
func (m Model) updateSelector(k string) (Model, bool) {
	if m.panel.ChoosingModel {
		switch k {
		case "esc":
			m.panel.ChoosingModel = false
			return m, true
		case "up", "k":
			m.panel.ModelCursor = wrap(m.panel.ModelCursor-1, len(m.modelChoices))
			return m, true
		case "down", "j":
			m.panel.ModelCursor = wrap(m.panel.ModelCursor+1, len(m.modelChoices))
			return m, true
		case "enter":
			return m.applyChosenModel(), true
		}
		return m, true
	}

	switch k {
	case "esc", "q":
		m.panel.InSelector = false
		m.panel.Agents = nil
		m.notice = ""
		return m, true

	// tab changes the destination, so the rows have to change with it. A selector
	// still showing the previous destination's agents under the new name is the
	// same lie the target strip exists to prevent.
	case "tab", "s":
		// Move nothing until the new destination's rows are in hand.
		//
		// Switching first and loading second leaves the strip naming one
		// destination while the rows below it belong to another when the load
		// fails — which is the exact lie the strip exists to prevent, produced by
		// the code meant to honour it. So the whole switch is abandoned on error
		// and the screen stays where it was.
		before := m
		next, _ := m.nextScope()
		m = next.(Model)

		rows, err := m.listAgents(m.activeScope())
		if err != nil {
			before.notice = "cannot read the agents: " + err.Error()
			return before, true
		}
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
		m.panel.AgentCursor, m.panel.ChoosingModel = 0, false
		return m, true

	case "up", "k":
		m.panel.AgentCursor = wrap(m.panel.AgentCursor-1, selectorRows(m.panel))
		return m, true

	case "down", "j":
		m.panel.AgentCursor = wrap(m.panel.AgentCursor+1, selectorRows(m.panel))
		return m, true

	case " ":
		if len(m.panel.Agents) == 0 {
			return m, true
		}
		if m.panel.AgentCursor == allRow {
			return m.markAll(), true
		}
		rows := append([]AgentRow(nil), m.panel.Agents...)
		rows[m.panel.AgentCursor-1].Marked = !rows[m.panel.AgentCursor-1].Marked
		m.panel.Agents = rows
		return m, true

	// The key and the row are one behaviour, so they are one function. Two copies
	// would drift, and the drift would be a row and a key disagreeing about what
	// "all" means.
	case "a":
		return m.markAll(), true

	case "m", "enter":
		if len(m.MarkedAgents()) == 0 {
			m.notice = "nothing marked — space marks a row, a marks all"
			return m, true
		}
		m.panel.ChoosingModel, m.panel.ModelCursor = true, 0
		return m, true
	}
	return m, true
}

// markAll marks every row, or clears them when they are already all marked.
//
// One gesture both ways: a control that only ever adds leaves no way back but
// pressing space once per row.
func (m Model) markAll() Model {
	marking := len(m.MarkedAgents()) < len(m.panel.Agents)
	rows := append([]AgentRow(nil), m.panel.Agents...)
	for i := range rows {
		rows[i].Marked = marking
	}
	m.panel.Agents = rows
	return m
}

// applyChosenModel sends the marked set and the chosen model to the caller.
//
// A failure keeps the screen and says what went wrong. Throwing the user back to
// the menu would hide the marks they would have to make again.
func (m Model) applyChosenModel() Model {
	names := m.MarkedAgents()
	model := m.ChosenModelName()

	m.panel.ChoosingModel = false

	// A no-op has to say so. SetModel does not rewrite a file that already declares
	// the model — deliberately, so the tool never dirties a git tree for nothing —
	// but from the outside "nothing happened because nothing needed to" and "nothing
	// happened because it is broken" look identical. Twice in one session that was
	// read as a bug.
	already := 0
	for _, a := range m.panel.Agents {
		if a.Marked && a.Model == model {
			already++
		}
	}
	if already == len(names) {
		m.notice = strconv.Itoa(already) + " agent(s) already on " + describeModel(model) + " — nothing to change"
		return m
	}

	if err := m.applyModel(m.activeScope(), names, model); err != nil {
		m.notice = "could not apply: " + err.Error()
		return m
	}

	// Ask for the rows again rather than editing them here. A screen that patches
	// its own state after a write is a second answer to what the files say, and the
	// two disagree the first time a write half-succeeds.
	if rows, err := m.listAgents(m.activeScope()); err == nil {
		marked := make(map[string]bool, len(names))
		for _, n := range names {
			marked[n] = true
		}
		for i := range rows {
			rows[i].Marked = marked[rows[i].Name]
		}
		// The models just changed, so the groups have to move with them. A screen
		// that keeps the old grouping under the new models is the same lie as one
		// that keeps the old models.
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
	}

	m.notice = strconv.Itoa(len(names)) + " agent(s) → " + describeModel(model)
	if m.refresh != nil {
		if menu, targets, err := m.refresh(m.activeScope()); err == nil {
			m.panel.Menu, m.panel.Targets = menu, targets
		}
	}
	return m
}

// describeModel renders a model for a human. The empty string is a state, and an
// empty column would read as a bug.
func describeModel(model string) string {
	if model == "" {
		return "(session)"
	}
	return model
}

// selector renders the second screen: one row per agent, and the catalogue over it
// when the user is choosing.
//
// The mark is `[x]` / `[ ]` and gold, never gold alone. Colour is the signal that
// disappears first — on a mono terminal, in a pipe, in a screenshot somebody
// pasted into a ticket — and a marking UI whose marks vanish is worse than one
// with no marks, because it still applies to them.
func (t Theme) selector(p Panel) string {
	if len(p.Agents) == 0 {
		return "  " + Fg(t.Muted).Render("no agents in this destination")
	}

	// The name column is measured, not borrowed. It used to reuse the main menu's
	// constant — sized for `install` and `prune` — while agent names run past twenty
	// characters, so the model and the `shared` warning started wherever each name
	// happened to end. pad never truncates, by design, so a too-narrow column does
	// not clip: it shifts everything after it.
	width := 0
	for _, a := range p.Agents {
		if n := len([]rune(a.Name)); n > width {
			width = n
		}
	}

	cw := ContentWidth(p.Width)

	rows := make([]string, 0, len(p.Agents)+len(p.ModelChoices)+4)

	// The `all` row is a control, not an agent. Its box says what is true of the set
	// below it — `[x]` only when every row is marked, because a master checkbox that
	// stays ticked after one row is cleared is a checkbox that lies.
	allBox := "[ ]"
	if countMarked(p.Agents) == len(p.Agents) {
		allBox = "[x]"
	}
	allColour, allCursor := t.Steel, " "
	if p.AgentCursor == allRow && !p.ChoosingModel {
		allColour, allCursor = t.Gold, "❯"
	}
	rows = append(rows,
		"  "+Fg(allColour).Render(allCursor+" "+allBox+" all"),
		t.groupRule(cw),
	)

	for i, a := range p.Agents {
		if i > 0 && a.Model != p.Agents[i-1].Model {
			rows = append(rows, t.groupRule(cw))
		}

		colour, cursor := t.Steel, " "
		if i+1 == p.AgentCursor && !p.ChoosingModel {
			colour, cursor = t.Gold, "❯"
		}

		box := "[ ]"
		if a.Marked {
			box = "[x]"
		}
		shared := ""
		if a.Shared {
			shared = "   shared"
		}
		line := cursor + " " + box + " " + pad(a.Name, width+2) +
			pad(describeModel(a.Model), 12) + shared
		rows = append(rows, "  "+Fg(colour).Render(line))
	}

	if p.ChoosingModel {
		rows = append(rows, "", "  "+Fg(t.Muted).Render("model for "+strconv.Itoa(countMarked(p.Agents))+" marked:"))
		for i, c := range p.ModelChoices {
			colour, cursor := t.Steel, " "
			if i == p.ModelCursor {
				colour, cursor = t.Gold, "❯"
			}
			name := c.Name
			if name == "" {
				name = "default"
			}
			line := cursor + "     " + pad(name, menuDescCol-menuLabelCol) + c.Label
			rows = append(rows, "  "+Fg(colour).Render(line))
		}
	}
	return strings.Join(rows, "\n")
}

// modelRank orders models the way the catalogue does: cheapest first, then the
// session default, and a model this build does not know about after both — absent
// from the map, so every caller has to decide what to do with a miss.
//
// Catalogue order is cheapest first, but it opens with the session default, which is
// not a price — it is an unknown. It goes last, so a list reads as an answer to "how
// much of this is still expensive?" rather than starting with the one entry that
// cannot answer it.
//
// One function rather than two, because the menu tally and the selector below it are
// one screen. Two orderings of the same models would be read as a bug, and it would
// be one.
func modelRank(order []ModelChoice) map[string]int {
	rank := make(map[string]int, len(order))
	for i, c := range order {
		rank[c.Name] = i
	}
	rank[""] = len(order)
	return rank
}

// sortRowsByModel groups the rows by model so the screen reads at a glance. Names
// sort inside a group: two agents on one model in an order nobody chose is a list
// that reorders itself between sessions.
func sortRowsByModel(rows []AgentRow, order []ModelChoice) []AgentRow {
	out := append([]AgentRow(nil), rows...)
	rank := modelRank(order)
	sort.SliceStable(out, func(i, j int) bool {
		ri, oki := rank[out[i].Model]
		rj, okj := rank[out[j].Model]
		if oki != okj {
			return oki
		}
		if ri != rj {
			return ri < rj
		}
		// Two models the catalogue does not know rank equally, and leaving it there
		// would interleave them by name — two groups shuffled into one, which is the
		// one thing this function exists to prevent.
		if out[i].Model != out[j].Model {
			return out[i].Model < out[j].Model
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// groupRule divides one model's rows from the next.
//
// An indented dim rule, deliberately not the frame's ├───┤ junction: that glyph means
// "a new section of the panel", and three of them inside one list would read as three
// panels. A blank line was the other candidate and it is worse — inside a bordered
// frame a blank line reads as the end of the list.
func (t Theme) groupRule(width int) string {
	n := width - 4
	if n < 1 {
		n = 1
	}
	return "  " + Fg(t.Dim).Render(strings.Repeat("─", n))
}

func countMarked(rows []AgentRow) int {
	n := 0
	for _, r := range rows {
		if r.Marked {
			n++
		}
	}
	return n
}

// Tally summarises a set of agents by model, cheapest first, for the menu row.
//
// It lives here rather than in the caller because the caller would have to know the
// rendering, and the menu row is a rendering. `order` is the catalogue order the
// panel was given, which is already cheapest first.
func Tally(rows []AgentRow, order []ModelChoice) string {
	counts := make(map[string]int, len(rows))
	for _, r := range rows {
		counts[r.Model]++
	}

	rank := modelRank(order)

	names := make([]string, 0, len(counts))
	for name := range counts {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		ri, oki := rank[names[i]]
		rj, okj := rank[names[j]]
		if oki != okj {
			return oki
		}
		if ri != rj {
			return ri < rj
		}
		return names[i] < names[j]
	})

	parts := make([]string, 0, len(names))
	for _, name := range names {
		label := name
		if label == "" {
			label = "session"
		}
		parts = append(parts, strconv.Itoa(counts[name])+" on "+label)
	}
	return strings.Join(parts, " · ")
}
