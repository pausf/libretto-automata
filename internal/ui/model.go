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
}

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
	m.notice = ""
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

// selectCurrent runs — or refuses to run — the highlighted action.
func (m Model) selectCurrent() (tea.Model, tea.Cmd) {
	if len(m.panel.Menu) == 0 {
		return m, nil
	}
	item := m.panel.Menu[m.panel.Selected]

	if !item.Enabled {
		m.notice = item.Label + " is not wired up yet"
		return m, nil
	}
	m.notice = "running " + item.Label + "…"
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
