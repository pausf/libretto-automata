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

// Model is the Bubbletea model driving the panel.
type Model struct {
	theme   Theme
	panel   Panel
	notice  string // one-line feedback under the panel
	done    bool
	refresh Refresh
	run     Runner
}

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
type Runner func(action string, destination int) ([]string, error)

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

// Results exposes the last report for tests.
func (m Model) Results() []string { return m.panel.Results }

func (m Model) Init() tea.Cmd { return nil }

// Update handles navigation and selection. Kept free of I/O so state transitions
// can be tested by calling it directly.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.panel.Width, m.panel.Height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit

		case "up", "k":
			m.notice = ""
			m.panel.Selected = wrap(m.panel.Selected-1, len(m.panel.Menu))
			return m, nil

		case "down", "j":
			m.notice = ""
			m.panel.Selected = wrap(m.panel.Selected+1, len(m.panel.Menu))
			return m, nil

		case "tab", "s":
			return m.nextScope()

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
		m.notice = "cannot read " + m.panel.Targets[next].Name + ": " + err.Error()
		return m, nil
	}

	m.panel.Menu, m.panel.Targets = menu, targets

	// Say it out loud. The switch changes the destination, not the cursor, so the
	// `❯` stays where it was and the only other evidence is a bullet and a path
	// that both take looking for. A key whose effect you have to hunt for is a key
	// people report as broken.
	m.notice = "acting on " + m.panel.Targets[next].Name
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

	if !item.Enabled {
		m.notice, m.panel.Results = item.Label+" is not wired up yet", nil
		return m, nil
	}
	if m.run == nil {
		m.notice, m.panel.Results = item.Label+" is not wired up yet", nil
		return m, nil
	}

	lines, err := m.run(item.Label, m.activeScope())
	m.panel.Results = lines
	if err != nil {
		// The report stays. An action that half-finished has the most to say, and
		// hiding it behind the error is hiding the part that explains it.
		m.notice = item.Label + " · " + err.Error()
		return m, nil
	}
	m.notice = item.Label + " · done"

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
	p.Notice = m.notice
	return m.theme.Render(p)
}

// Selected exposes the cursor for tests.
func (m Model) Selected() int { return m.panel.Selected }

// Notice exposes the feedback line for tests.
func (m Model) Notice() string { return m.notice }

// wrap moves the cursor cyclically, so ↑ from the first row lands on the last.
func wrap(i, n int) int {
	if n == 0 {
		return 0
	}
	return ((i % n) + n) % n
}
