module github.com/pausf/libretto-automata

go 1.26.5

require (
	github.com/charmbracelet/bubbletea v1.3.10
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/charmbracelet/x/term v0.2.2
	github.com/mattn/go-isatty v0.0.24
	github.com/muesli/termenv v0.16.0
)

require (
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/x/ansi v0.10.1 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.3.8 // indirect
)

// This project is pre-1.0 and has not declared its contract stable. These versions were
// published from a bump table read mechanically: the merge reversed a promise in
// .agents/specs/ci/spec.md, which describes how this repository releases itself, while the
// tool's contract — install, prune, the flags, the payload's skills — did not move by a
// line. In 0.x that is a minor, not a major.
//
// The code they carried is in v0.5.2. Nothing is missing from the 0.5.x line.
retract (
	v1.0.0 // not a major; the tool's contract never changed. Use v0.5.2.
	v1.0.1 // not a major; the tool's contract never changed. Use v0.5.2.
	v1.0.2 // exists only to carry these retractions. Never a release.
)
