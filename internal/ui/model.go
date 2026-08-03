package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the Bubbletea model driving the panel.
type Model struct {
	theme  Theme
	panel  Panel
	notice string // one-line feedback under the panel
	done   bool
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

		case "enter":
			return m.selectCurrent()
		}
	}
	return m, nil
}

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
