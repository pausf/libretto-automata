package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Refresh rebuilds the menu and strip for the destination at index i.
//
// The model owns which destination is active; it does not know what a destination
// *is*. Counts, roots and tallies come from the caller, which keeps this package
// free of the filesystem and free of any dependency on the target package.
type Refresh func(i int) ([]MenuItem, []TargetRow, error)

// ToggleScope flips the scope every row shows — global to project and back — and
// returns fresh figures for the tool at index i plus the new scope's name. The
// model owns neither axis's meaning: tools arrive as rows, the scope as a label.
type ToggleScope func(i int) ([]MenuItem, []TargetRow, string, error)

// Model is the Bubbletea model driving the panel.
type Model struct {
	theme   Theme
	panel   Panel
	notice  string // one-line feedback under the panel
	refused bool   // the notice is something the panel declined, not something it did
	done    bool
	refresh Refresh
	toggle  ToggleScope
	run     Runner

	// pending is a destructive action that has shown its plan and is waiting on an
	// answer. Scoped to the destination it was planned for, so a confirmation can
	// never be spent somewhere else.
	//
	// An explicit question beats a second press: "⏎ again to go ahead" is a rule you
	// have to know, and the key that means "go ahead" is the same one that meant
	// "show me" a moment ago. `y` and `n` are an answer to something asked.
	pending      string
	pendingScope int

	// The model selector's wiring. Callbacks rather than a dependency, so this
	// package still cannot read a file — see models.go.
	modelChoices  []ModelChoice
	effortChoices []EffortChoice
	listAgents    ListAgents
	applyModel    ApplyModel
	applyEffort   ApplyEffort

	// releaseCheck asks whether a newer release exists. Nil is the ordinary case.
	releaseCheck ReleaseCheck
}

// ReleaseCheck returns the notice to show, or "" for nothing to say.
//
// It returns the finished row rather than a version, so the comparison stays in
// internal/repo — one implementation of "is this newer", not two that can disagree. Which
// also keeps this package free of any knowledge of what a version is.
type ReleaseCheck func() string

// releaseMsg carries the check's answer back into Update.
type releaseMsg string

// Runner performs a menu action and returns its report, one line per row.
//
// The model knows action *labels* and nothing about what they do. Keeping the work
// on the caller's side is what stops this package from touching the filesystem, and
// it is why the report can be the command's own output rather than a second
// rendering of the same facts.
// The destination index is passed in, never captured at startup. A runner that
// closed over the scope the panel opened with would send `prune` at the old
// destination after a tab — the panel saying "project" while links vanish from the
// global config. Destructive, silent, and entirely avoidable.
// confirm is true only on the second press of a destructive action. Without it the
// second press would repeat the dry run and the panel would promise a confirmation
// it never delivered.
type Runner func(action string, destination int, confirm bool) ([]string, error)

// NewModel builds the model. Menu and targets are supplied by the caller so this
// package never reads the filesystem.
func NewModel(version string, menu []MenuItem, targets []TargetRow, asciiSafe bool) Model {
	return Model{
		theme: NewTheme(),
		panel: Panel{
			Version:   version,
			Menu:      menu,
			Targets:   targets,
			ASCIISafe: asciiSafe,
		},
	}
}

// WithRefresh lets the panel change which destination is active.
//
// Without it the panel can only ever act on whatever it was handed, which for a
// tool that writes into one of two places is half a feature: you could see the
// other destination and never choose it.
func (m Model) WithRefresh(r Refresh) Model {
	m.refresh = r
	return m
}

// WithRunner lets the panel carry actions out. Without it every action refuses,
// which is honest: a menu that cannot run anything says so rather than pretending.
func (m Model) WithRunner(r Runner) Model {
	m.run = r
	return m
}

// WithScopeToggle lets `s` flip the scope. Without it the key is inert, which is
// what preview and every fixture that predates scopes want.
func (m Model) WithScopeToggle(t ToggleScope) Model {
	m.toggle = t
	return m
}

// WithScopeLabel sets the strip's scope line for the first render; the toggle
// keeps it current afterwards.
func (m Model) WithScopeLabel(s string) Model {
	m.panel.Scope = s
	return m
}

// Results exposes the last report for tests.
func (m Model) Results() []string { return m.panel.Results }

// Pending reports the destructive action waiting on an answer, if any.
func (m Model) Pending() string { return m.pending }

// Confirm exposes the open question for tests.
func (m Model) Confirm() string { return m.panel.Confirm }

// forget drops an unanswered question without acting on it.
func (m Model) forget() Model {
	m.pending, m.pendingScope, m.panel.Confirm = "", 0, ""
	return m
}

// cancelPending answers no.
func (m Model) cancelPending() (tea.Model, tea.Cmd) {
	label := m.pending
	m = m.forget()
	m = m.say(label + " · cancelled, nothing changed")
	return m, nil
}

// carryOut answers yes.
//
// The destination is the one the plan was made for, not whatever is active now. They
// cannot differ — switching destination cancels the question — but reading it from
// the answer rather than from the current state is what makes that true by
// construction instead of by discipline.
func (m Model) carryOut() (tea.Model, tea.Cmd) {
	label, dest := m.pending, m.pendingScope
	m = m.forget()

	lines, err := m.run(label, dest, true)
	m.panel.Results = lines
	if err != nil {
		m = m.refuse(label + " · " + err.Error())
		return m, nil
	}
	m = m.say(label + " · done")

	if m.refresh != nil {
		if menu, targets, rerr := m.refresh(m.activeScope()); rerr == nil {
			m.panel.Menu, m.panel.Targets = menu, targets
		}
	}
	return m, nil
}

// WithReleaseCheck lets the panel say a newer release exists. Without it the panel is
// silent on the subject, which is what `preview` and every test that does not care want.
func (m Model) WithReleaseCheck(c ReleaseCheck) Model {
	m.releaseCheck = c
	return m
}

// Init runs the release check, if there is one, as a command.
//
// A command and not a call: the panel paints complete without the answer and re-renders
// when it arrives. Asking inline would put a subprocess and a network round trip in front
// of the first frame, and on bad DNS the user's only recourse is ⌃C on a tool that looks
// broken.
func (m Model) Init() tea.Cmd {
	if m.releaseCheck == nil {
		return nil
	}
	return func() tea.Msg { return releaseMsg(m.releaseCheck()) }
}

// Update handles navigation and selection. Kept free of I/O so state transitions
// can be tested by calling it directly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.panel.Width, m.panel.Height = msg.Width, msg.Height
		return m, nil

	case releaseMsg:
		// Empty means up to date, or the check could not answer. Both are silence, and
		// neither is worth a row.
		m.panel.UpdateNotice = string(msg)
		return m, nil

	case tea.KeyMsg:
		// A question is open: answer it, or cancel it and let the key do its usual
		// job. Anything other than "yes" cancels — a destructive action never
		// proceeds on a key that meant something else.
		if m.pending != "" {
			switch msg.String() {
			case "y", "Y":
				return m.carryOut()
			case "n", "N", "esc":
				return m.cancelPending()
			default:
				m = m.forget()
			}
		}

		// The selector owns every key while it is up, including esc and q. Falling
		// through to the menu's handlers would quit the panel on the key that means
		// "go back one screen".
		if m.panel.InSelector {
			if msg.String() == "ctrl+c" {
				m.done = true
				return m, tea.Quit
			}
			next, _ := m.updateSelector(msg.String())
			return next, nil
		}

		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit

		case "up", "k":
			m = m.say("")
			m.panel.Selected = wrap(m.panel.Selected-1, len(m.panel.Menu))
			return m, nil

		case "down", "j":
			m = m.say("")
			m.panel.Selected = wrap(m.panel.Selected+1, len(m.panel.Menu))
			return m, nil

		case "tab":
			return m.nextScope()

		case "s":
			return m.toggleScope()

		case "enter":
			return m.selectCurrent()
		}
	}
	return m, nil
}

// nextScope moves the active mark to the following destination and asks the caller
// for figures that match it.
//
// A refresh that fails leaves the previous state alone and says so. Showing counts
// from one destination under the name of another would be worse than showing an
// error — it would be a panel that lies about where install goes.
func (m Model) nextScope() (tea.Model, tea.Cmd) {
	if m.refresh == nil || len(m.panel.Targets) < 2 {
		return m, nil
	}

	next := wrap(m.activeScope()+1, len(m.panel.Targets))
	menu, targets, err := m.refresh(next)
	if err != nil {
		m = m.refuse("cannot read " + m.panel.Targets[next].Name + ": " + err.Error())
		return m, nil
	}

	m.panel.Menu, m.panel.Targets = menu, targets

	// Say it out loud. The switch changes the destination, not the cursor, so the
	// `❯` stays where it was and the only other evidence is a bullet and a path
	// that both take looking for. A key whose effect you have to hunt for is a key
	// people report as broken.
	m = m.say("acting on " + m.panel.Targets[next].Name)
	return m, nil
}

// toggleScope flips every row to the other scope and asks for figures that match.
//
// The same refusal rule as nextScope: a failed toggle leaves the previous state
// alone, because figures from one scope under another's label are a lie.
func (m Model) toggleScope() (tea.Model, tea.Cmd) {
	if m.toggle == nil {
		return m, nil
	}

	menu, targets, scope, err := m.toggle(m.activeScope())
	if err != nil {
		m = m.refuse("cannot switch scope: " + err.Error())
		return m, nil
	}

	m.panel.Menu, m.panel.Targets, m.panel.Scope = menu, targets, scope
	m = m.say("scope · " + scope)
	return m, nil
}

// activeScope is the index of the marked destination, or 0 when none is.
func (m Model) activeScope() int {
	for i, t := range m.panel.Targets {
		if t.Active {
			return i
		}
	}
	return 0
}

// ActiveScope exposes the marked destination for tests.
func (m Model) ActiveScope() int { return m.activeScope() }

// selectCurrent runs the highlighted action and shows what it did, without leaving.
//
// Staying is the point. Quitting to print the report closes the loop on the terminal
// and reopens it in the user's head: install, relaunch, tab, install again. Keeping
// the panel up means the destination, the state and the last report are all on screen
// at once, which is the only arrangement where "did that go where I meant?" has an
// answer you can see rather than remember.
//
// An earlier version set the notice to "running install…" and ran nothing at all. A
// panel that says it is working and is not is worse than a panel with no actions.
func (m Model) selectCurrent() (tea.Model, tea.Cmd) {
	if len(m.panel.Menu) == 0 {
		return m, nil
	}
	item := m.panel.Menu[m.panel.Selected]

	// The one action that opens a screen instead of running something.
	if item.Label == modelsAction {
		return m.openSelector(), nil
	}

	if !item.Enabled {
		m.notice, m.panel.Results = item.Label+" is not wired up yet", nil
		return m, nil
	}
	if m.run == nil {
		m.notice, m.panel.Results = item.Label+" is not wired up yet", nil
		return m, nil
	}

	// A destructive action reports what it would do and then asks.
	//
	// Telling the user to leave the panel and type `prune --yes` was a confirmation
	// step that threw away the plan it was confirming. Asking with a second press was
	// worse in a different way: the key that means "go ahead" was the same one that
	// meant "show me" a moment earlier, and the rule lived in a notice you had to
	// read. A question with two answers is a question.
	lines, err := m.run(item.Label, m.activeScope(), false)
	m.panel.Results = lines

	if err != nil {
		// The report stays. An action that half-finished has the most to say, and
		// hiding it behind the error is hiding the part that explains it.
		m = m.refuse(item.Label + " · " + err.Error())
		return m, nil
	}

	if item.Destructive {
		if len(lines) == 0 {
			m = m.say(item.Label + " · nothing to do")
			return m, nil
		}
		m.pending, m.pendingScope = item.Label, m.activeScope()
		m.panel.Confirm = "Go ahead and " + item.Label + " " +
			m.panel.Targets[m.activeScope()].Name + "?   y / n"
		m = m.refuse("nothing has changed yet")
		return m, nil
	}
	m = m.say(item.Label + " · done")

	// The figures on screen describe the state before the action. Ask for them again
	// or the panel contradicts its own report.
	if m.refresh != nil {
		if menu, targets, rerr := m.refresh(m.activeScope()); rerr == nil {
			m.panel.Menu, m.panel.Targets = menu, targets
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.done {
		return ""
	}
	p := m.panel
	p.Notice, p.Refused = m.notice, m.refused
	return m.theme.Render(p)
}

// Selected exposes the cursor for tests.
func (m Model) Selected() int { return m.panel.Selected }

// Notice exposes the feedback line for tests.
func (m Model) Notice() string { return m.notice }

// Refused reports whether the current notice is a refusal, for tests.
func (m Model) Refused() bool { return m.refused }

// say records a notice about something that happened.
//
// refuse records one about something that did not, and the panel draws that in the error
// colour inside its own box.
//
// Two setters rather than a field anybody may assign, because the message and its kind
// going out of step is a silent failure with a bad shape: a refusal left set would paint
// the next success red, and a success left set would draw the next refusal as one grey
// line under a bordered panel — which is the state this pair exists to end. There are
// twenty-odd notice sites, and the two that matter are always the ones nobody re-read.
func (m Model) say(msg string) Model {
	m.notice, m.refused = msg, false
	return m
}

func (m Model) refuse(msg string) Model {
	m.notice, m.refused = msg, true
	return m
}

// wrap moves the cursor cyclically, so ↑ from the first row lands on the last.
func wrap(i, n int) int {
	if n == 0 {
		return 0
	}
	return ((i % n) + n) % n
}

// MenuItemForTest exposes one menu row, so a test can assert what it says rather
// than scraping it back out of the rendered frame.
func (m Model) MenuItemForTest(i int) MenuItem { return m.panel.Menu[i] }

// SetSelectedForTest moves the cursor from outside the package. Tests in cmd need it
// to reach a row without simulating keypresses to get there.
func (m Model) SetSelectedForTest(i int) Model {
	m.panel.Selected = i
	return m
}
