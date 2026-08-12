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
//
// Every agent stays reachable whether or not it is on screen — the window follows the
// cursor rather than bounding it. A cursor that could not leave the window would make
// the rows below the fold unmarkable, which is the bug rather than the fix.
func selectorRows(p Panel) int { return len(p.Agents) + 1 }

// windowFloor is the fewest agent rows worth showing. Below this the screen is chrome
// with a peephole in it, and the honest failure is a panel that overflows a tiny
// terminal rather than one that pretends to be usable in it.
const windowFloor = 3

// selectorChrome is how many lines the selector spends on everything that is not an
// agent row: the frame, the wordmark, the destination strip, the footer, the notice, the
// `all` row and its rule.
//
// ponytail: one number, and the arithmetic behind it is measured rather than estimated —
// 25 lines observed for the frame, wordmark, strip, notice, footer, `all` row and its
// rule, plus 2 for the out-of-view indicators, plus one rule per model boundary in the
// window, of which the catalogue has at most three beyond the one already counted.
//
// Deliberately generous, because the two ways of being wrong do not cost the same: a
// line too many is a blank line under the list, a line too few is the torn screen this
// exists to fix. The gate is TestTheSelectorFitsTheTerminalHeight, not this comment —
// the first guess here was 20 and the test said 26 before the arithmetic was done.
//
// The upgrade path, the day a destination is added or the chrome moves: render the panel
// with no agents, measure it, and give the rows what is left. That costs a second render
// per keystroke, which is why it is not here yet.
const selectorChrome = 30

// agentWindow is how many agent rows fit, and it is every one of them when the height is
// unknown.
//
// Height 0 means nobody told us — `preview`, a pipe, most tests — and bounding a list
// against a height we do not have would hide rows for no reason.
func agentWindow(p Panel) int {
	if p.Height <= 0 {
		return len(p.Agents)
	}
	n := p.Height - selectorChrome
	// An open catalogue draws its own rows below the list, and they are the thing the
	// user is looking at while it is open.
	if p.ChoosingModel {
		n -= len(p.ModelChoices) + 2
	}
	if p.ChoosingEffort {
		n -= len(p.EffortChoices) + 2
	}
	if n < windowFloor {
		n = windowFloor
	}
	if n > len(p.Agents) {
		n = len(p.Agents)
	}
	return n
}

// scrollToCursor moves the window the least it can to bring the cursor inside it.
//
// Called after every cursor move rather than inside the renderer, and it is what makes
// the rows below the fold reachable — the whole point. `all` is cursor 0 and is never
// windowed: it is one row, it is the master checkbox for the list under it, and a
// checkbox that scrolls away from the list it governs is a checkbox you cannot find.
func (m Model) scrollToCursor() Model {
	size := agentWindow(m.panel)
	if size >= len(m.panel.Agents) {
		m.panel.AgentTop = 0
		return m
	}

	// The `all` row sits at cursor 0, so agent i is at cursor i+1.
	at := m.panel.AgentCursor - 1
	if at < 0 {
		m.panel.AgentTop = 0
		return m
	}

	switch {
	case at < m.panel.AgentTop:
		m.panel.AgentTop = at
	case at >= m.panel.AgentTop+size:
		m.panel.AgentTop = at - size + 1
	}
	if max := len(m.panel.Agents) - size; m.panel.AgentTop > max {
		m.panel.AgentTop = max
	}
	if m.panel.AgentTop < 0 {
		m.panel.AgentTop = 0
	}
	return m
}

// AgentRow is one agent in the selector: what it is called, what it runs on, and
// whether the user has marked it.
type AgentRow struct {
	Name   string
	Model  string // empty means the session's model
	Effort string // empty means the session's effort
	Marked bool

	// Efforts is the levels this row's model can run, weakest first. Empty means the
	// model has none — Haiku, today.
	//
	// It arrives as data for the same reason the model catalogue does: this package
	// deciding which models support effort would put that rule in a second place, and
	// the copy on screen would be the one nobody updated. What it buys is the
	// difference between offering a choice and refusing one after it was made — the
	// first version had only the apply-time error to go on, so `e` on two Haiku rows
	// opened a menu of five levels none of which could be written.
	Efforts []string

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

// EffortChoice is one entry of the catalogue the `e` key opens.
//
// An alias rather than a second struct: an effort choice is a name and a label, which
// is exactly what a model choice is. Declaring the same two fields twice would be two
// things to keep in step for no property gained, and the alias lets both call sites
// read as what they are.
type EffortChoice = ModelChoice

// ApplyModel declares one model on a set of agents in one destination.
type ApplyModel func(destination int, names []string, model string) error

// ApplyEffort declares one effort level on a set of agents in one destination.
//
// Separate from ApplyModel because the two keys are independent — an agent staying on
// Opus while its effort drops is the case the feature exists for, and one callback
// taking both would make every effort change restate a model.
type ApplyEffort func(destination int, names []string, effort string) error

// WithAgents wires the selector. Without it the menu row refuses, the same way an
// unwired action does.
func (m Model) WithAgents(models []ModelChoice, efforts []EffortChoice, list ListAgents, apply ApplyModel, applyEffort ApplyEffort) Model {
	m.modelChoices = models
	m.effortChoices = efforts
	m.listAgents = list
	m.applyModel = apply
	m.applyEffort = applyEffort
	return m
}

// InSelector reports whether the second screen is up.
func (m Model) InSelector() bool { return m.panel.InSelector }

// ChoosingModel reports whether the model catalogue is open over the selector.
func (m Model) ChoosingModel() bool { return m.panel.ChoosingModel }

// ChoosingEffort reports whether the effort catalogue is open over the selector.
func (m Model) ChoosingEffort() bool { return m.panel.ChoosingEffort }

// EffortChoices exposes the full effort catalogue for tests.
func (m Model) EffortChoices() []EffortChoice { return m.effortChoices }

// OpenEffortChoices exposes the narrowed list actually on screen. Distinct from
// EffortChoices because the difference between the two is the whole of the fix: a test
// asserting against the full catalogue would pass over a screen offering levels it
// cannot write.
func (m Model) OpenEffortChoices() []EffortChoice { return m.panel.EffortChoices }

// ChosenEffortName is the effort entry under the cursor.
//
// It reads the panel's list rather than the full catalogue: the open list is narrowed to
// what the marked set can actually run, so the cursor indexes that one. Reading the full
// catalogue here would pick a different entry than the one under the cursor as soon as
// the two differed, which is a wrong write with a correct-looking screen.
func (m Model) ChosenEffortName() string {
	if len(m.panel.EffortChoices) == 0 {
		return ""
	}
	return m.panel.EffortChoices[m.panel.EffortCursor].Name
}

// offerableEfforts narrows the catalogue to the levels every marked row can run.
//
// The intersection, not the union. A level offered because *some* marked row can run it
// is a level whose apply is refused for the whole set — and the refusal is all-or-nothing
// by contract, so the union would offer choices guaranteed to fail.
//
// Empty means no level applies to the whole set, which is a state the caller reports
// rather than a menu it opens.
func (m Model) offerableEfforts() []EffortChoice {
	marked := 0
	allowed := map[string]int{}
	for _, a := range m.panel.Agents {
		if !a.Marked {
			continue
		}
		marked++
		for _, e := range a.Efforts {
			allowed[e]++
		}
	}
	if marked == 0 {
		return nil
	}

	out := make([]EffortChoice, 0, len(m.effortChoices))
	for _, c := range m.effortChoices {
		// The session default is the absence of the key, not a level, so no row has to
		// support it — removing an effort is legal on any agent, including one whose
		// model has none.
		if c.Name == "" || allowed[c.Name] == marked {
			out = append(out, c)
		}
	}

	// Only the default survived: there is no level the whole set can run, and offering
	// a one-entry menu whose only choice is "remove" would read as the feature being
	// broken rather than inapplicable.
	if len(out) <= 1 {
		return nil
	}
	return out
}

// modelsWithoutEffort names the marked rows' models that have no levels, for the notice
// that replaces the menu. The models, not the agent names: the reason is the model, and
// eleven agent names is a sentence nobody reads to the end.
func (m Model) modelsWithoutEffort() []string {
	seen := map[string]bool{}
	var out []string
	for _, a := range m.panel.Agents {
		if !a.Marked || len(a.Efforts) > 0 {
			continue
		}
		name := describeModel(a.Model)
		if !seen[name] {
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

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
		m = m.refuse("models is not wired up yet")
		return m
	}

	rows, err := m.listAgents(m.activeScope())
	if err != nil {
		m = m.refuse("cannot read the agents: " + err.Error())
		return m
	}

	m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
	m.panel.ModelChoices = m.modelChoices
	// Not pre-populated: the open list is narrowed to what the marked set can run, and
	// nothing is marked yet.
	m.panel.EffortChoices = nil
	m.panel.InSelector, m.panel.ChoosingModel, m.panel.ChoosingEffort = true, false, false
	m.panel.AgentCursor, m.panel.ModelCursor, m.panel.EffortCursor = 0, 0, 0
	m.panel.AgentTop = 0
	m = m.say("space mark · a all · m model · e effort · esc back")
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

	// The same four keys over the other catalogue. Two blocks rather than one
	// parametrised by which is open: the cursor, the choices and the apply are three
	// different fields, and a single block would spend more lines choosing between
	// them than it saves.
	if m.panel.ChoosingEffort {
		switch k {
		case "esc":
			m.panel.ChoosingEffort = false
			return m, true
		case "up", "k":
			m.panel.EffortCursor = wrap(m.panel.EffortCursor-1, len(m.panel.EffortChoices))
			return m, true
		case "down", "j":
			m.panel.EffortCursor = wrap(m.panel.EffortCursor+1, len(m.panel.EffortChoices))
			return m, true
		case "enter":
			return m.applyChosenEffort(), true
		}
		return m, true
	}

	switch k {
	case "esc", "q":
		m.panel.InSelector = false
		m.panel.Agents = nil
		m = m.say("")
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
			before = before.refuse("cannot read the agents: " + err.Error())
			return before, true
		}
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
		m.panel.AgentCursor, m.panel.ChoosingModel = 0, false
		m.panel.AgentTop = 0
		return m, true

	case "up", "k":
		m.panel.AgentCursor = wrap(m.panel.AgentCursor-1, selectorRows(m.panel))
		return m.scrollToCursor(), true

	case "down", "j":
		m.panel.AgentCursor = wrap(m.panel.AgentCursor+1, selectorRows(m.panel))
		return m.scrollToCursor(), true

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

	// enter stays the model. It is what this screen has always meant by it, and what
	// the footer has always said; rebinding a reflex to pick up a second key would be
	// a silent change to the one gesture nobody reads the legend for.
	case "m", "enter":
		if len(m.MarkedAgents()) == 0 {
			m = m.refuse("nothing marked — space marks a row, a marks all")
			return m, true
		}
		m.panel.ChoosingModel, m.panel.ModelCursor = true, 0
		return m, true

	case "e":
		if len(m.MarkedAgents()) == 0 {
			m = m.refuse("nothing marked — space marks a row, a marks all")
			return m, true
		}

		// Measured before the menu, not after the choice. A catalogue of five levels
		// over two Haiku rows was the reported bug, and it is worse than a plain
		// refusal: the user picks, waits, and is told no by a screen that offered yes.
		offer := m.offerableEfforts()
		if len(offer) == 0 {
			if without := m.modelsWithoutEffort(); len(without) > 0 {
				// "haiku and sonnet has" is what the first version said, seen in a
				// rendered panel rather than reasoned about. A refusal is read at the
				// moment something went wrong, which is the worst moment to be reading
				// broken English.
				verb := " has "
				if len(without) > 1 {
					verb = " have "
				}
				m = m.refuse(strings.Join(without, " and ") + verb + "no effort levels — unmark those rows, or change the model with m")
			} else {
				m = m.refuse("no effort level applies to every marked row")
			}
			return m, true
		}

		m.panel.EffortChoices = offer
		m.panel.ChoosingEffort, m.panel.EffortCursor = true, 0
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
		m = m.say(strconv.Itoa(already) + " agent(s) already on " + describeModel(model) + " — nothing to change")
		return m
	}

	if err := m.applyModel(m.activeScope(), names, model); err != nil {
		m = m.refuse("could not apply: " + err.Error())
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

	m = m.say(strconv.Itoa(len(names)) + " agent(s) → " + describeModel(model))
	if m.refresh != nil {
		if menu, targets, err := m.refresh(m.activeScope()); err == nil {
			m.panel.Menu, m.panel.Targets = menu, targets
		}
	}
	return m
}

// applyChosenEffort sends the marked set and the chosen level to the caller.
//
// The model's twin, and deliberately its twin rather than a shared function taking a
// key: the no-op check reads a different field, the refusal reads differently, and the
// notice says a different word. What they share is the shape, which is worth keeping
// visible.
func (m Model) applyChosenEffort() Model {
	names := m.MarkedAgents()
	effort := m.ChosenEffortName()

	m.panel.ChoosingEffort = false

	already := 0
	for _, a := range m.panel.Agents {
		if a.Marked && a.Effort == effort {
			already++
		}
	}
	if already == len(names) {
		m = m.say(strconv.Itoa(already) + " agent(s) already at " + describeModel(effort) + " — nothing to change")
		return m
	}

	// A refusal is the interesting path here in a way it is not for the model: a
	// marked row on a model with no effort levels sends the whole apply back, and the
	// screen has to keep every mark so the user can unmark that one row rather than
	// start again.
	if err := m.applyEffort(m.activeScope(), names, effort); err != nil {
		m = m.refuse("could not apply: " + err.Error())
		return m
	}

	if rows, err := m.listAgents(m.activeScope()); err == nil {
		marked := make(map[string]bool, len(names))
		for _, n := range names {
			marked[n] = true
		}
		for i := range rows {
			rows[i].Marked = marked[rows[i].Name]
		}
		// Still grouped by model. The effort did not move any row between groups, and
		// re-sorting on a second axis is the twenty-group screen the spec refused.
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
	}

	m = m.say(strconv.Itoa(len(names)) + " agent(s) → effort " + describeModel(effort))
	if m.refresh != nil {
		if menu, targets, err := m.refresh(m.activeScope()); err == nil {
			m.panel.Menu, m.panel.Targets = menu, targets
		}
	}
	return m
}

// describeModel renders a model or an effort for a human. The empty string is a state
// for both keys, and an empty column would read as a bug.
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

	cw := ContentWidth(p.Width)

	// The name column is measured, not borrowed. It used to reuse the main menu's
	// constant — sized for `install` and `prune` — while agent names run past twenty
	// characters, so the model and the `shared` warning started wherever each name
	// happened to end. pad never truncates, by design, so a too-narrow column does
	// not clip: it shifts everything after it.
	//
	// Measured *against the frame*, though, and that is what the effort column
	// changed. Two value columns plus `shared` plus the longest name the payload
	// ships came to 64 against a 58-column interior, and pad's refusal to truncate
	// turned that into six columns of torn border on every row of the screen. So the
	// budget is computed and the name yields to it — visibly, with an ellipsis. A
	// name that says it was cut is a smaller lie than a frame that came apart.
	width := nameColumn(p.Agents, cw)

	rows := make([]string, 0, len(p.Agents)+len(p.ModelChoices)+4)

	// The `all` row is a control, not an agent. Its box says what is true of the set
	// below it — `[x]` only when every row is marked, because a master checkbox that
	// stays ticked after one row is cleared is a checkbox that lies.
	allBox := "[ ]"
	if countMarked(p.Agents) == len(p.Agents) {
		allBox = "[x]"
	}
	allColour, allCursor := t.Steel, " "
	if p.AgentCursor == allRow && !p.choosing() {
		allColour, allCursor = t.Gold, "❯"
	}
	rows = append(rows,
		"  "+Fg(allColour).Render(allCursor+" "+allBox+" all"),
		t.groupRule(cw),
	)

	// The window, and what is outside it. A bounded list that does not admit it was
	// bounded is the lie the action report already refuses to tell — and here it is
	// worse, because the rows out of view are rows you have to reach to mark.
	size := agentWindow(p)
	top := p.AgentTop
	if top > len(p.Agents)-size {
		top = len(p.Agents) - size
	}
	if top < 0 {
		top = 0
	}
	if above := top; above > 0 {
		rows = append(rows, "  "+Fg(t.Dim).Render("↑ "+strconv.Itoa(above)+" more"))
	}

	for i := top; i < top+size; i++ {
		a := p.Agents[i]
		// Against the previous *visible* row: a rule drawn because the row above it is
		// off screen is a rule separating nothing from nothing.
		if i > top && a.Model != p.Agents[i-1].Model {
			rows = append(rows, t.groupRule(cw))
		}

		colour, cursor := t.Steel, " "
		if i+1 == p.AgentCursor && !p.choosing() {
			colour, cursor = t.Gold, "❯"
		}

		box := "[ ]"
		if a.Marked {
			box = "[x]"
		}
		shared := ""
		if a.Shared {
			shared = sharedNote
		}
		// width-1, because pad adds a trailing space rather than truncating when its
		// input already fills the column — so a name elided to exactly `width` would
		// come out one column wider and put the tear back.
		line := cursor + " " + box + " " + pad(elide(a.Name, width-1), width) +
			pad(describeModel(a.Model), modelCol) + pad(describeModel(a.Effort), effortCol) + shared
		rows = append(rows, "  "+Fg(colour).Render(line))
	}

	if below := len(p.Agents) - top - size; below > 0 {
		rows = append(rows, "  "+Fg(t.Dim).Render("↓ "+strconv.Itoa(below)+" more"))
	}

	if p.ChoosingModel {
		rows = append(rows, t.catalogue("model for "+strconv.Itoa(countMarked(p.Agents))+" marked:", p.ModelChoices, p.ModelCursor)...)
	}
	if p.ChoosingEffort {
		rows = append(rows, t.catalogue("effort for "+strconv.Itoa(countMarked(p.Agents))+" marked:", p.EffortChoices, p.EffortCursor)...)
	}
	return strings.Join(rows, "\n")
}

// The selector row's fixed parts, in columns: the cursor, its space, the `[x]` box,
// its space, and the two value columns. `(session)` is the longest thing either value
// column holds, at nine, so ten gives each a single space of gap.
const (
	rowFixed   = 1 + 1 + 3 + 1
	modelCol   = 12
	effortCol  = 10
	sharedNote = "   shared"

	// nameFloor is the narrowest a name column may become. Below this the ellipsis is
	// most of what is left, and a list of `revi…` rows is a list you cannot use — at
	// that point the honest failure is a name that overflows one row, not a screen of
	// unreadable stubs.
	nameFloor = 12
)

// nameColumn is how many columns the agent names get: what they measure, or what the
// frame can spare, whichever is smaller.
//
// `shared` is counted only when a row actually carries it. Reserving it always would
// narrow the names on every screen to protect a warning most screens do not show.
func nameColumn(rows []AgentRow, cw int) int {
	measured := 0
	shared := 0
	for _, a := range rows {
		if n := len([]rune(a.Name)); n > measured {
			measured = n
		}
		if a.Shared {
			shared = len([]rune(sharedNote))
		}
	}

	// cw-2 because the caller indents every row by two columns.
	budget := cw - 2 - rowFixed - modelCol - effortCol - shared
	if budget < nameFloor {
		budget = nameFloor
	}
	if want := measured + 2; want < budget {
		return want
	}
	return budget
}

// elide shortens s to n columns, ending in `…` when it had to cut.
//
// The counterpart to pad, which never truncates because a column that silently cuts
// its content lies about it. This one is not silent, which is the whole difference:
// the ellipsis is the row saying there is more name than there is room.
func elide(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n < 1 {
		return ""
	}
	return string(r[:n-1]) + "…"
}

// choosing reports whether either catalogue is open over the rows.
//
// The row cursor hides while one is, because two gold cursors on one screen is two
// answers to "where am I". Adding the second catalogue without this would have left
// the row cursor lit under the effort list — visible only by opening it, which is why
// it is one predicate rather than a condition repeated at each of the three sites.
func (p Panel) choosing() bool { return p.ChoosingModel || p.ChoosingEffort }

// catalogue renders one open choice list. Shared by both because they are the same
// list of the same two fields, and two copies would drift in the indentation nobody
// notices until the screens sit next to each other.
func (t Theme) catalogue(title string, choices []ModelChoice, cursorAt int) []string {
	out := []string{"", "  " + Fg(t.Muted).Render(title)}
	for i, c := range choices {
		colour, cursor := t.Steel, " "
		if i == cursorAt {
			colour, cursor = t.Gold, "❯"
		}
		name := c.Name
		if name == "" {
			name = "default"
		}
		line := cursor + "     " + pad(name, menuDescCol-menuLabelCol) + c.Label
		out = append(out, "  "+Fg(colour).Render(line))
	}
	return out
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

// VisibleAgentsForTest is the slice of rows the window is currently showing.
//
// Exported for tests because the alternative is asserting on rendered text, and a test
// that counts names in a string cannot tell "the window shrank" from "the name changed".
func (m Model) VisibleAgentsForTest() []AgentRow {
	size := agentWindow(m.panel)
	top := m.panel.AgentTop
	if top+size > len(m.panel.Agents) {
		top = len(m.panel.Agents) - size
	}
	if top < 0 {
		top = 0
	}
	return m.panel.Agents[top : top+size]
}
