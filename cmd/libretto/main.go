// Command lib is the Libretto Automata CLI.
//
// With a TTY and no subcommand it shows the panel. Otherwise it behaves as a
// plain command, so scripts and CI can use it — see docs/SPEC.md R8.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"

	"github.com/pausf/libretto-automata/internal/agentmodel"
	"github.com/pausf/libretto-automata/internal/link"
	"github.com/pausf/libretto-automata/internal/repo"
	"github.com/pausf/libretto-automata/internal/target"
	"github.com/pausf/libretto-automata/internal/ui"
)

// version is stamped at build time from the nearest git tag:
//
//	go build -ldflags "-X main.version=$(git describe --tags --always --dirty)"
//
// `make build` does that. The fallback matters when someone builds with a plain
// `go build` or `go install`, and it is deliberately not a version number —
// claiming to be v0.1.0 when nothing said so is the drift this replaces. A binary
// that cannot prove its version says so instead of guessing.
var version = "dev"

// EnvASCIISafe swaps the clef's quadrant glyphs for half blocks, for fonts that
// lack them. See docs/DESIGN.md.
const EnvASCIISafe = "LIBRETTO_ASCII"

func main() {
	// Resolved once, into the same variable everything already reads, so a binary built
	// by `go install` reports its module version everywhere the stamped one appears —
	// the panel footer and `version` included. See version.go.
	version = buildVersion(version)

	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", invokedAs(), err)
		os.Exit(1)
	}
}

// scopeFlags pulls --global / --project out of the arguments and returns the rest.
//
// Both at once is an error rather than a precedence rule. Two answers to "where
// should this go" is a mistake worth reporting, not one worth resolving quietly by
// picking the last flag and hoping it was meant.
//
// Neither means global, which is what every invocation meant before this existed.
// Changing what an existing command does is worse than making people type a flag
// for the new thing.
//
// chosen is which flag was seen, empty when none was. The panel needs to tell "no
// flag" from "--global" — they resolve to the same scope and mean different things,
// one being a default and the other an instruction. Returned rather than recomputed
// by a second scanner: two places deciding what a scope flag is would agree only by
// accident, and this file already records that bug happening with the project root.
func scopeFlags(args []string) (target.Scope, string, []string, error) {
	scope, chosen := target.GlobalScope, ""
	rest := make([]string, 0, len(args))

	for _, a := range args {
		switch a {
		case "--global", "-g":
			if chosen == "project" {
				return scope, "", nil, fmt.Errorf("--global and --project are two answers to one question; pick one")
			}
			scope, chosen = target.GlobalScope, "global"
		case "--project", "-p":
			if chosen == "global" {
				return scope, "", nil, fmt.Errorf("--global and --project are two answers to one question; pick one")
			}
			scope, chosen = target.ProjectScope, "project"
		default:
			rest = append(rest, a)
		}
	}
	return scope, chosen, rest, nil
}

// openingScope is the destination the panel opens on.
//
// A flag is about this run and wins. With no flag, the panel opens where it was
// left — which is the whole feature, and it exists as a named function because
// `run`'s panel branch cannot be reached from a test: with no arguments it checks
// isatty, and under `go test` that is false.
//
// Subcommands never call this. That is what keeps a typed command from changing
// meaning because of state left by a session the reader cannot see.
func openingScope(flagged target.Scope, chosen string) target.Scope {
	if chosen != "" {
		return flagged
	}
	return rememberedScope()
}

func run(args []string) error {
	scope, chosen, args, err := scopeFlags(args)
	if err != nil {
		return err
	}

	// `version` and `help` answer without the payload, and they are answered before the
	// clone is even looked for. Cloning a repository into somebody's home because they
	// asked what version they were running would be indefensible.
	if len(args) > 0 {
		switch args[0] {
		case "version", "-v", "--version":
			fmt.Println("libretto-automata", version)
			return nil
		case "help", "-h", "--help":
			usage()
			return nil
		}
	}

	// Everything below links, reads or reports on the payload, so from here a clone has
	// to exist. This is where `go install` gets one.
	root, err := ensureClone()
	if err != nil {
		return err
	}
	// Where "the project" is, decided once.
	//
	// Resolving it separately in each place that needs it means two answers that
	// agree only by accident — and they did disagree: the strip read one root while
	// an action wrote to another. One lookup, threaded down.
	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}
	tg := target.Resolve(scope, projectDir)

	if len(args) == 0 {
		if !isatty.IsTerminal(os.Stdout.Fd()) {
			usage()
			os.Exit(2)
		}
		// Only the panel remembers. `tg` above is deliberately left alone, so the
		// subcommand paths below cannot be reached by this at all — which is the
		// promise, not an implementation detail.
		return panelUI(root, projectDir, openingScope(scope, chosen))
	}

	switch args[0] {
	case "status":
		return status(root, tg)
	case "preview":
		return preview(root, tg)
	case "install":
		return install(root, tg)
	case "doctor":
		return doctor(root, tg)
	case "prune":
		return prune(root, tg, args[1:])
	case "uninstall":
		return uninstall(root, tg, args[1:])
	case "update":
		return update(root, tg)
	case "models":
		return models(root, tg, args[1:])
	default:
		usage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func panelUI(root, projectDir string, scope target.Scope) error {
	menu, targets, err := panelData(root, projectDir, scope)
	if err != nil {
		return err
	}

	model := ui.NewModel(version, menu, targets, asciiSafe()).
		WithRefresh(panelRefresh(root, projectDir)).
		WithReleaseCheck(func() string { return releaseNotice(root, version) })

	// Actions run inside the panel and report there, so the destination, the state
	// and the last report stay on screen together.
	//
	// The destination comes in as an argument, never captured here: a closure over
	// the scope the panel opened with would send `prune` at the old destination after
	// a tab — the strip reading "project" while links disappear from the global
	// config. Destructive and silent, which is the worst pair.
	model = model.WithRunner(func(action string, dest int, confirm bool) ([]string, error) {
		if dest < 0 || dest >= len(scopeOrder) {
			return nil, fmt.Errorf("no destination %d", dest)
		}
		return runCaptured(action, root, target.Resolve(scopeOrder[dest], projectDir), confirm)
	})

	// The selector's two callbacks. The destination comes in as an argument, never
	// captured here — the same rule as the runner above, and the same failure if it
	// were not: the strip would name one destination while the rows came from
	// another, and the write would land where the user could not see it.
	model = model.WithAgents(
		modelChoices(),
		func(dest int) ([]ui.AgentRow, error) {
			tg, err := destination(dest, projectDir)
			if err != nil {
				return nil, err
			}
			return agentRows(root, tg)
		},
		func(dest int, names []string, m string) error {
			tg, err := destination(dest, projectDir)
			if err != nil {
				return err
			}
			return agentmodel.Apply(agentsDir(tg), names, m)
		},
	)

	_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}

// panelRefresh gives the panel a fresh view of the destination at index i, and
// remembers that it moved there.
//
// The panel changes which destination is active by asking for this. scopeOrder is the
// single place that maps a row back to a scope, so the strip and the refresh can never
// disagree about the order.
//
// The write happens *after* panelData succeeds, and only then. A failed refresh leaves
// the panel where it was, and a file that disagrees with the screen is the same class
// of lie as a strip showing one destination's counts under another's name.
//
// It is a named function rather than a closure so a test can hand it a bad index. The
// panel's own path is unreachable from a test — `run` checks isatty first.
func panelRefresh(root, projectDir string) func(int) ([]ui.MenuItem, []ui.TargetRow, error) {
	return func(i int) ([]ui.MenuItem, []ui.TargetRow, error) {
		if i < 0 || i >= len(scopeOrder) {
			return nil, nil, fmt.Errorf("no destination %d", i)
		}

		menu, targets, err := panelData(root, projectDir, scopeOrder[i])
		if err != nil {
			return nil, nil, err
		}

		remember(scopeOrder[i])
		return menu, targets, nil
	}
}

// runCaptured performs an action and returns what it printed, line by line.
//
// The commands write their reports to stdout. Inside a full-screen TUI that output
// would land on top of the panel, so it is redirected for the duration and handed
// back as lines instead. That also means the panel shows the command's own words
// rather than a second rendering of the same facts, which could disagree with it.
func runCaptured(action, root string, tg target.Target, confirm bool) ([]string, error) {
	prev := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = w

	// Drain concurrently: a report longer than the pipe buffer would otherwise block
	// the action half-way through, with the links half-written.
	out := make(chan string, 1)
	go func() {
		var buf strings.Builder
		_, _ = io.Copy(&buf, r)
		out <- buf.String()
	}()

	runErr := dispatch(action, root, tg, confirm)

	_ = w.Close()
	os.Stdout = prev
	text := <-out
	_ = r.Close()

	var lines []string
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, strings.TrimRight(l, " "))
		}
	}
	return lines, runErr
}

// dispatch runs one menu action. The panel's labels are the subcommand names, so
// there is one list of actions and not two to keep in agreement.
func dispatch(action, root string, tg target.Target, confirm bool) error {
	switch action {
	case "install":
		return install(root, tg)
	case "update":
		return update(root, tg)
	case "status":
		return status(root, tg)
	case "doctor":
		return doctor(root, tg)
	case "uninstall":
		if confirm {
			return uninstall(root, tg, []string{"--yes"})
		}
		return uninstall(root, tg, nil)
	case "prune":
		// Dry unless confirmed. Nothing chosen from a menu deletes on one keypress;
		// the panel asks again and only then passes the confirmation through.
		if confirm {
			return prune(root, tg, []string{"--yes"})
		}
		return prune(root, tg, nil)
	default:
		return fmt.Errorf("unknown action %q", action)
	}
}

// scopeOrder is the order destinations appear in the strip, and the order tab
// cycles through them. Global first: it is the default everywhere else.
var scopeOrder = []target.Scope{target.GlobalScope, target.ProjectScope}

// panelData assembles the menu and the target strip for the active scope.
//
// **Both scopes are listed, one is marked active.** The alternative — asking once,
// at startup, which one to use — is worse: an answer given at the top of a session
// is invisible by the time you press a key, and "where did that just install?" is
// the question this whole strip exists to answer before it is asked.
//
// The install row names the active root, so the destination is legible without
// counting bullets.
func panelData(root, projectDir string, scope target.Scope) ([]ui.MenuItem, []ui.TargetRow, error) {
	active := target.Resolve(scope, projectDir)

	rows := make([]ui.TargetRow, 0, len(scopeOrder))
	overall := map[link.State]int{}

	// Each row reports **its own** state, not the repo's contents.
	//
	// It used to show `link.Counts`, which counts items in the repo filtered by the
	// kinds a target accepts. Both scopes accept the same three kinds, so both rows
	// showed identical numbers — always, by construction. The column answered "what
	// does the repo hold?" twice and looked like it was answering "what is installed
	// here?". Two rows that cannot differ are two rows that mislead.
	//
	// Scanning both costs a second read-only pass, which is the correct price for a
	// strip that distinguishes its destinations.
	for _, sc := range scopeOrder {
		tg := target.Resolve(sc, projectDir)

		entries, err := link.Scan(root, tg)
		if err != nil {
			return nil, nil, err
		}
		tally := link.Tally(entries)

		if sc == scope {
			for state, n := range tally {
				overall[state] += n
			}
		}

		rows = append(rows, ui.TargetRow{
			Name:       string(sc),
			Info:       pad(summarise(tally), 24) + shorten(tg.Root()),
			Configured: configured(tg),
			Active:     sc == scope,
		})
	}

	// The status row carries the live tally, exactly as the design mocks it.
	menu := []ui.MenuItem{
		{Label: "install", Desc: "link the score into " + shorten(active.Root()), Enabled: true},
		{Label: "uninstall", Desc: "take it back out of " + shorten(active.Root()), Enabled: true, Destructive: true},
		{Label: "update", Desc: "git pull · relink · report", Enabled: true},
		{Label: "status", Desc: summarise(overall), Enabled: true},
	}

	// Next to status, and reporting like it. "choose agent models" would be the one
	// entry in the menu describing itself instead of saying something, and the
	// tally is the question the screen exists to answer.
	//
	// A repo with no agents/ directory simply has no row. An entry that opens an
	// empty screen is a promise the panel cannot keep.
	if agents, err := agentRows(root, active); err == nil && len(agents) > 0 {
		menu = append(menu, ui.MenuItem{
			Label:   "models",
			Desc:    ui.Tally(agents, modelChoices()),
			Enabled: true,
		})
	}

	menu = append(menu,
		ui.MenuItem{Label: "doctor", Desc: "diagnose the orchestra", Enabled: true},
		ui.MenuItem{Label: "prune", Desc: "drop links whose source is gone", Enabled: true, Destructive: true},
	)

	return menu, rows, nil
}

// agentRows adapts the agent package to the panel's row type.
//
// The adaptation lives here rather than in either package: internal/ui must not
// learn what an agent file is, and internal/agentmodel must not learn what a panel
// row is. One of them knowing about the other is how a package that promises to
// stay off the filesystem stops being able to prove it.
func agentRows(root string, tg target.Target) ([]ui.AgentRow, error) {
	agents, _, err := agentmodel.Agents(agentsDir(tg))
	if err != nil {
		return nil, err
	}
	rows := make([]ui.AgentRow, 0, len(agents))
	for _, a := range agents {
		rows = append(rows, ui.AgentRow{
			Name:  a.Name,
			Model: a.Model,
			// Owned means the file is one of ours, reached from more than one
			// destination — so writing it is not local to this one.
			Shared: link.Owned(root, a.Path),
		})
	}
	return rows, nil
}

// destination maps a strip row back to a target. scopeOrder is the single place that
// mapping lives, so the strip and every callback agree by construction.
func destination(i int, projectDir string) (target.Target, error) {
	if i < 0 || i >= len(scopeOrder) {
		return nil, fmt.Errorf("no destination %d", i)
	}
	return target.Resolve(scopeOrder[i], projectDir), nil
}

func modelChoices() []ui.ModelChoice {
	cat := agentmodel.Catalogue()
	out := make([]ui.ModelChoice, 0, len(cat))
	for _, m := range cat {
		// The version leads the label. `opus` answers which tier; "Opus 5"
		// answers the question the user opened the screen with.
		label := m.Label
		if m.Version != "" {
			label = m.Version + " · " + m.Label
		}
		out = append(out, ui.ModelChoice{Name: m.Name, Label: label})
	}
	return out
}

// pad widens s to n columns so the roots line up in the strip. It never truncates —
// a column that silently cuts its content is a column that lies about it.
func pad(s string, n int) string {
	if len([]rune(s)) >= n {
		return s + " "
	}
	return s + strings.Repeat(" ", n-len([]rune(s)))
}

// pathBudget is how many columns a root may occupy in the panel.
//
// ponytail: a fixed budget, not a computed one. The panel's content is 58–98
// columns and the label and counts take the rest, so 34 fits the narrow case with
// room. If the layout ever needs it exact, the width is known in internal/ui and
// the eliding belongs there instead.
const pathBudget = 34

// shorten makes a root fit the panel without tearing it.
//
// Two problems, and the second one is why this exists at all. A path under the home
// directory reads better as `~/...` — cosmetic. But an arbitrarily long path
// overflows the frame and pushes the right border out of alignment, which the panel
// spec forbids at every width. A temporary directory is enough to do it.
//
// So the tail is kept and the head is replaced with an ellipsis: the last segments
// are what tell you *which* directory this is, and the marker says plainly that
// something was removed. Silently cutting the end would leave a path that looks
// complete and is not.
func shorten(path string) string {
	if path == "" {
		return "not configured"
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" && strings.HasPrefix(path, home) {
		path = "~" + strings.TrimPrefix(path, home)
	}
	if len([]rune(path)) <= pathBudget {
		return path
	}

	// Keep whole segments from the right, so the result is still a readable path
	// rather than a string chopped mid-name.
	parts := strings.Split(path, string(filepath.Separator))
	out := ""
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == "" {
			continue
		}
		next := string(filepath.Separator) + parts[i] + out
		if len([]rune(next))+1 > pathBudget {
			break
		}
		out = next
	}
	if out == "" {
		// One segment longer than the whole budget. Cut it, and say so.
		r := []rune(path)
		return "…" + string(r[len(r)-pathBudget+1:])
	}
	return "…" + out
}

// plural singularises a kind name for a count of one. Every kind is named in the
// plural, so trimming the trailing "s" is enough — no irregulars to handle.
func plural(kind string, n int) string {
	if n == 1 {
		return strings.TrimSuffix(kind, "s")
	}
	return kind
}

// configured reports whether the target is present on disk. Only Claude
// implements Exists today; anything else is assumed present.
func configured(tg target.Target) bool {
	if c, ok := tg.(interface{ Exists() bool }); ok {
		return c.Exists()
	}
	return true
}

// preview prints the panel once and exits, without starting the TUI.
//
// It forces truecolor so the output is identical whether stdout is a terminal or
// a pipe. That is what makes it usable for screenshots, golden files, and eyeing
// the LIBRETTO_ASCII fallback.
func preview(root string, tg target.Target) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}
	menu, targets, err := panelData(root, projectDir, target.GlobalScope)
	if err != nil {
		return err
	}

	lipgloss.SetColorProfile(termenv.TrueColor)
	theme := ui.NewTheme()

	fmt.Println()
	fmt.Println(theme.Render(ui.Panel{
		Version:   version,
		Menu:      menu,
		Selected:  indexOf(menu, "status"),
		Targets:   targets,
		Width:     terminalWidth(),
		ASCIISafe: asciiSafe(),
	}))
	return nil
}

// terminalWidth reports the width to lay out against, or 0 when it is unknown.
//
// COLUMNS is honoured as a fallback so a piped preview can still be rendered at
// a chosen width — otherwise centring would be impossible to check without a
// terminal.
func terminalWidth() int {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil {
			return w
		}
	}
	if w, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && w > 0 {
		return w
	}
	return 0
}

func indexOf(menu []ui.MenuItem, label string) int {
	for i, item := range menu {
		if item.Label == label {
			return i
		}
	}
	return 0
}

// status reports every item's state in every target. Read-only, per SPEC R4.
func status(root string, tg target.Target) error {
	{
		where := "not configured"
		if configured(tg) {
			where = tg.Root()
		}
		fmt.Printf("%s  %s\n", tg.Name(), where)

		entries, err := link.Scan(root, tg)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			fmt.Println("  nothing to link")
			return nil
		}

		for _, e := range entries {
			fmt.Printf("  %-12s %s/%s\n", e.State, e.Kind, e.Name)
		}
		fmt.Printf("  %s\n", summarise(link.Tally(entries)))
	}
	return nil
}

// install links every item into every target.
//
// Idempotent by construction: an already-correct tree produces an empty plan, so
// running it twice is not a special case that needed handling (SPEC R9).
//
// Conflicts are reported and never touched, and they make the exit code non-zero:
// an item that did not get linked is an incomplete install, whatever the reason.
func install(root string, tg target.Target) error {
	var wrote, refused, failed, conflicts int

	{
		entries, err := link.Scan(root, tg)
		if err != nil {
			return err
		}

		plan := link.Plan(entries)
		fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
		if len(plan) == 0 {
			fmt.Println("  already correct")
		}

		for _, r := range link.Apply(root, plan) {
			e := r.Action.Entry
			switch {
			case r.Action.Act == link.Skip:
				fmt.Printf("  skip     %s/%s (%s)\n", e.Kind, e.Name, e.State)
				conflicts++
			case r.Refused:
				fmt.Printf("  refused  %s/%s — %v\n", e.Kind, e.Name, r.Err)
				refused++
			case r.Err != nil:
				fmt.Printf("  FAILED   %s/%s — %v\n", e.Kind, e.Name, r.Err)
				failed++
			default:
				fmt.Printf("  %-8s %s/%s\n", r.Action.Act, e.Kind, e.Name)
				wrote++
			}
		}
	}

	fmt.Printf("\n%d linked · %d skipped · %d refused · %d failed\n",
		wrote, conflicts, refused, failed)

	if failed+refused+conflicts > 0 {
		return fmt.Errorf("%d item(s) were not linked", failed+refused+conflicts)
	}
	return nil
}

// update refreshes the repo, rebuilds the binary when the Go source moved, and
// relinks. SPEC R3.
//
// The order matters. Nothing is pulled over uncommitted work, and nothing is
// relinked from source that has not been compiled — a payload from the new commit
// linked by a binary from the old one is a state nobody asked for and nobody can
// reason about.
func update(root string, tg target.Target) error {
	git := repo.Shell{Root: root}

	dirty, err := git.Dirty()
	if err != nil {
		return fmt.Errorf("cannot read the repo state: %w", err)
	}
	if dirty {
		fmt.Println("the working tree has uncommitted changes")
		fmt.Println("nothing was pulled and no link was touched — commit or stash first")
		return fmt.Errorf("refusing to pull over uncommitted work")
	}

	before, err := git.Head()
	if err != nil {
		return err
	}

	// No remote is not a failure. There is nothing to pull from, and relinking is
	// still worth doing — the items on disk may have moved since the last install.
	remote, err := git.HasRemote()
	if err != nil {
		return err
	}
	if !remote {
		fmt.Println("no remote configured — skipping the pull, relinking anyway")
	} else {
		if err := git.Pull(); err != nil {
			return fmt.Errorf("pull failed: %w", err)
		}
		after, _ := git.Head()
		switch {
		case before == "" || after == "":
			fmt.Println("pulled")
		case before == after:
			fmt.Println("already up to date")
		default:
			fmt.Printf("pulled  %s → %s\n", short(before), short(after))
		}
	}

	changed, err := git.ChangedSince(before)
	if err != nil {
		return err
	}
	if repo.NeedsRebuild(changed) {
		// The binary that is running, not the one in the clone. A $GOBIN install and a
		// `make link` symlink both end up here, and both are what the user's next
		// invocation will execute.
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("the pull landed but the binary could not be located: %w", err)
		}

		// The clone's bin/ is the fallback: the user owns the checkout, so it is the one
		// place a write cannot be refused.
		note, err := rebuildOrReport(root, exe, filepath.Join(root, "bin", "libretto"))
		if err != nil {
			return fmt.Errorf("the pull landed but the rebuild failed: %w", err)
		}
		// The advice depends on whether the running binary was actually replaced. Telling
		// somebody to "run it again to use the new one" when the rename was refused is
		// advice that does nothing and reads as though the upgrade took.
		if note != "" {
			fmt.Println("rebuilt  " + note)
			fmt.Println("         " + exe + " is unchanged and still what your PATH runs")
		} else {
			fmt.Println("rebuilt  " + exe)
			fmt.Println("         this process is still the old binary — run it again to use the new one")
		}
	}

	fmt.Println()
	return install(root, tg)
}

// rebuild compiles the CLI to a temporary file and renames it over dest.
//
// dest is the binary that is *running*, not bin/libretto. Once `go install` can put the
// command in $GOBIN, rebuilding into the clone's bin/ upgrades a file nobody executes:
// `update` would report success and every later invocation would stay on the old version.
// The destination is a parameter rather than read from os.Executable() inside, so a test
// does not have to be the binary under test.
//
// A symlink is resolved and written through. `make link` puts one in ~/.local/bin
// pointing at bin/libretto, and replacing that link with a regular file would sever the
// development setup silently.
//
// Writing straight over the running executable is what produces "text file busy", and a
// half-written binary is worse than a stale one. Rename is atomic, and the process
// already running keeps the inode it started from — which is why the caller has to say so
// rather than pretend the upgrade took effect mid-run.
func rebuild(root, dest string) error {
	// EvalSymlinks fails on a path that does not exist yet, which is fine: there is no
	// link to follow, and the original is the destination.
	if resolved, err := filepath.EvalSymlinks(dest); err == nil {
		dest = resolved
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	// The temporary file goes beside the destination, not in TMPDIR: rename across
	// filesystems fails, and $GOBIN and /tmp are routinely on different ones.
	tmp := dest + ".new"

	// Probe the write before paying for a compile. Without this the permission failure
	// arrives as a `go build` error string — unrecognisable as a permission problem, so
	// rebuildOrReport could not tell "you cannot write there" from "the code is broken",
	// and a locked $GOBIN would read as a compile failure after a three-second wait.
	probe, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		return err
	}
	probe.Close()

	cmd := exec.Command("go", "build", "-o", tmp, "./cmd/libretto")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp)
		if text := strings.TrimSpace(string(out)); text != "" {
			return fmt.Errorf("%s", text)
		}
		return err
	}

	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// rebuildOrReport rebuilds, and turns "I could not write there" into a note instead of an
// error.
//
// The pull already happened and the links are already correct. Failing the whole update
// because a rename was refused would roll back work that succeeded, and leave the user
// with neither the new binary nor the relink. The note has to say where the binary is, or
// it is a dead end dressed as an explanation — and saying that means building it
// somewhere, which is what fallback is for.
//
// fallback is a parameter, not derived. Otherwise the test that proves the note names a
// real binary has to write into the clone's bin/ — the same file another test asserts is
// left alone, so the two would collide under `go test -shuffle`.
func rebuildOrReport(root, dest, fallback string) (string, error) {
	err := rebuild(root, dest)
	if err == nil {
		return "", nil
	}
	if !errors.Is(err, os.ErrPermission) {
		return "", err
	}

	if berr := rebuild(root, fallback); berr != nil {
		return "", err
	}
	return fmt.Sprintf("could not replace %s — no permission\n"+
		"         the new binary is at %s; move it there yourself", dest, fallback), nil
}

func short(rev string) string {
	if len(rev) > 7 {
		return rev[:7]
	}
	return rev
}

// doctor reports what needs attention and what the flow expects to find. It never
// writes.
//
// Two sections with different authority: link problems are this tool's business
// and set the exit code; prerequisites are informational, because the flow works
// without the optional ones and saying otherwise would be a lie.
func doctor(root string, tg target.Target) error {
	problems := 0

	{
		fmt.Printf("%s  %s\n", tg.Name(), tg.Root())

		entries, err := link.Scan(root, tg)
		if err != nil {
			return err
		}

		clean := true
		for _, e := range entries {
			if !e.State.NeedsAttention() {
				continue
			}
			clean = false
			problems++
			fmt.Printf("  %-12s %s/%s  %s\n", e.State, e.Kind, e.Name, remedy(e.State))
		}
		if clean {
			fmt.Println("  no problems")
		}
	}

	// Live, not cached: the user typed a diagnostic and can afford the wait. It never
	// sets the exit code — being a release behind is news, and an unreachable remote is
	// not this tool's fault.
	fmt.Println("\nrelease")
	fmt.Println("  " + releaseLine(version, repo.Shell{Root: root}.LatestTag))

	fmt.Println("\nprerequisites")
	for _, p := range prerequisites() {
		mark := "—"
		if p.Found {
			mark = "ok"
		}
		fmt.Printf("  %-4s %-22s %s\n", mark, p.Name, p.Note)
	}

	if problems > 0 {
		return fmt.Errorf("%d item(s) need attention", problems)
	}
	return nil
}

// remedy names the command that fixes a state, so the report says what to do
// rather than only what is wrong.
//
// The command names itself with invokedAs, for the same reason usage does: this
// binary is linked under more than one name, and a remedy naming a command the
// reader does not have is a remedy they cannot run.
func remedy(s link.State) string {
	n := invokedAs()
	switch s {
	case link.Missing, link.WrongTarget:
		return "→ " + n + " install"
	case link.Stale:
		return "→ " + n + " prune"
	case link.Conflict:
		return "→ yours, not ours; move it or rename it"
	default:
		return ""
	}
}

// prune drops owned links with no item behind them.
//
// Dry by default. It deletes things, and a destructive command that acts before
// being asked twice is a command that eventually deletes the wrong thing. Without
// --yes it prints the plan and changes nothing.
func prune(root string, tg target.Target, args []string) error {
	confirmed := len(args) > 0 && (args[0] == "--yes" || args[0] == "-y")

	var planned, wrote, refused, failed int

	{
		entries, err := link.Scan(root, tg)
		if err != nil {
			return err
		}

		plan := link.PrunePlan(entries)
		planned = len(plan)

		switch {
		case planned == 0:
			// Nothing to remove. The report below says so.

		case !confirmed:
			fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
			for _, a := range plan {
				// The resolved destination is an absolute path that repeats what the
				// item name already says. Shortened, so the line stays readable in a
				// report and in the panel.
				fmt.Printf("  would remove  %s/%s → %s\n",
					a.Entry.Kind, a.Entry.Name, shorten(a.Entry.Actual))
			}

		default:
			fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
			for _, r := range link.Apply(root, plan) {
				e := r.Action.Entry
				switch {
				case r.Refused:
					fmt.Printf("  refused  %s/%s — %v\n", e.Kind, e.Name, r.Err)
					refused++
				case r.Err != nil:
					fmt.Printf("  FAILED   %s/%s — %v\n", e.Kind, e.Name, r.Err)
					failed++
				default:
					fmt.Printf("  removed  %s/%s\n", e.Kind, e.Name)
					wrote++
				}
			}
		}
	}

	if planned == 0 {
		fmt.Println("nothing to prune")
		return nil
	}
	if !confirmed {
		fmt.Printf("\n%d link(s) to remove. Nothing has been changed.\n", planned)
		fmt.Printf("Run `%s prune --yes` to go ahead.\n", invokedAs())
		return nil
	}

	fmt.Printf("\n%d removed · %d refused · %d failed\n", wrote, refused, failed)
	if failed+refused > 0 {
		return fmt.Errorf("%d link(s) were not removed", failed+refused)
	}
	return nil
}

// uninstall undoes this repository's work at one destination.
//
// The pair of install, not of prune. Prune cleans up after the repo changed and
// deliberately spares links that are correct; this removes links that are working,
// because the user changed their mind rather than because the repo did.
//
// Dry by default for the same reason prune is: a destructive command that acts before
// being asked twice eventually acts on the wrong thing, and a pipe is no reason to be
// less careful.
func uninstall(root string, tg target.Target, args []string) error {
	confirmed := len(args) > 0 && (args[0] == "--yes" || args[0] == "-y")

	entries, err := link.Scan(root, tg)
	if err != nil {
		return err
	}
	plan := link.UninstallPlan(entries)
	writes := link.Writes(plan)

	if len(writes) == 0 {
		fmt.Printf("%s  %s\n", tg.Name(), tg.Root())
		fmt.Println("  nothing of ours is installed here")
		return nil
	}

	fmt.Printf("%s  %s\n", tg.Name(), tg.Root())

	var removed, refused, failed, conflicts int
	if !confirmed {
		for _, a := range plan {
			if a.Act == link.Skip {
				fmt.Printf("  keep          %s/%s (%s — not ours)\n",
					a.Entry.Kind, a.Entry.Name, a.Entry.State)
				conflicts++
				continue
			}
			fmt.Printf("  would remove  %s/%s\n", a.Entry.Kind, a.Entry.Name)
		}
		fmt.Printf("\n%d link(s) to remove. Nothing has been changed.\n", len(writes))
		fmt.Printf("Run `%s uninstall --yes` to go ahead.\n", invokedAs())
		return nil
	}

	for _, r := range link.Apply(root, plan) {
		e := r.Action.Entry
		switch {
		case r.Action.Act == link.Skip:
			fmt.Printf("  keep     %s/%s (%s — not ours)\n", e.Kind, e.Name, e.State)
			conflicts++
		case r.Refused:
			fmt.Printf("  refused  %s/%s — %v\n", e.Kind, e.Name, r.Err)
			refused++
		case r.Err != nil:
			fmt.Printf("  FAILED   %s/%s — %v\n", e.Kind, e.Name, r.Err)
			failed++
		default:
			fmt.Printf("  removed  %s/%s\n", e.Kind, e.Name)
			removed++
		}
	}

	fmt.Printf("\n%d removed · %d kept · %d refused · %d failed\n",
		removed, conflicts, refused, failed)

	// Anything that survived means the destination is not in the state the report
	// implies, so a script has to be able to tell.
	if refused+failed > 0 {
		return fmt.Errorf("%d link(s) were not removed", refused+failed)
	}
	return nil
}

// Prereq is one thing the flow looks for on this machine.
type Prereq struct {
	Name  string
	Found bool
	Note  string
}

// prerequisites reports what the payload's own skills expect to find.
//
// Nothing here is required. `find-work` needs the jira CLI, `record-work`
// needs a git host, and ponytail and caveman are companions the flow calls when
// present. Reporting a missing optional as a problem would train people to ignore
// this section.
func prerequisites() []Prereq {
	home, _ := os.UserHomeDir()

	jira := onPath("jira")
	jiraNote := "find-work — brew install ankitpokhrel/jira-cli/jira-cli"
	if jira {
		cfg := os.Getenv("JIRA_CONFIG_FILE")
		if cfg == "" {
			cfg = filepath.Join(home, ".config", ".jira", ".config.yml")
		}
		if _, err := os.Stat(cfg); err == nil {
			jiraNote = "find-work — configured"
		} else {
			jiraNote = "find-work — installed, run `jira init` yourself"
			jira = false
		}
	}

	host := "record-work — install gh or glab"
	gh, gl := onPath("gh"), onPath("glab")
	switch {
	case gh && gl:
		host = "record-work — gh and glab"
	case gh:
		host = "record-work — gh"
	case gl:
		host = "record-work — glab"
	}

	return []Prereq{
		{"jira", jira, jiraNote},
		{"git host", gh || gl, host},
		{"rg", onPath("rg"), "record-work, find-work — brew install ripgrep"},
		{"jq", onPath("jq"), "find-work — brew install jq"},
		{"ponytail", companion(home, "ponytail"), "optional — how much gets built"},
		{"caveman", companion(home, "caveman"), "optional — how much gets said"},
	}
}

func onPath(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// companion reports whether an optional tool is installed anywhere it can
// legitimately live.
//
// All four places are real and in use: ponytail arrives as a plugin, caveman as
// plain commands. Checking only the plugin tree reports an installed companion as
// missing, which is worse than not checking at all — it sends people to install
// something they already have.
func companion(home, name string) bool {
	claude := filepath.Join(home, ".claude")
	for _, pattern := range []string{
		filepath.Join(claude, "plugins", "*", "*"+name+"*"),
		filepath.Join(claude, "plugins", "*", "*", "*"+name+"*"),
		filepath.Join(claude, "skills", "*"+name+"*"),
		filepath.Join(claude, "commands", "*"+name+"*"),
	} {
		if hits, _ := filepath.Glob(pattern); len(hits) > 0 {
			return true
		}
	}
	return false
}

// summarise renders a tally in a fixed state order, omitting zeros, so two runs
// with the same situation produce the same line.
func summarise(counts map[link.State]int) string {
	order := []link.State{link.Linked, link.Missing, link.WrongTarget, link.Conflict, link.Stale}

	var parts []string
	for _, s := range order {
		if n := counts[s]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, s))
		}
	}
	if len(parts) == 0 {
		return "nothing to link"
	}
	return strings.Join(parts, " · ")
}

// repoRoot lives in root.go — it grew a rung and a bug worth its own test file.

func asciiSafe() bool { return os.Getenv(EnvASCIISafe) == "safe" }

// invokedAs reports the name this binary was called by.
//
// The command can be linked under any name — `libretto-automata`, a shorter alias,
// whatever the user prefers — so help that hardcodes one name is help that lies to
// everyone using another.
func invokedAs() string {
	if len(os.Args) == 0 || os.Args[0] == "" {
		return "libretto-automata"
	}
	return filepath.Base(os.Args[0])
}

func usage() {
	n := invokedAs()
	fmt.Fprintf(os.Stderr, `libretto-automata %s

  %[2]s               show the panel (requires a terminal)
  %[2]s install       link every item into each target
  %[2]s update        pull, relink, report
  %[2]s status        what is linked
  %[2]s doctor        diagnose links and repo state
  %[2]s prune         show links whose source is gone, change nothing
  %[2]s prune --yes   remove them
  %[2]s uninstall     show what this repo installed here, change nothing
  %[2]s uninstall --yes  remove it
  %[2]s preview       print the panel once, no TUI
  %[2]s models        which model each agent runs on
  %[2]s models set <model> <agent>…   declare it; --all for every agent

  --global, -g          act on ~/.claude (the default)
  --project, -p         act on <this directory>/.claude

  LIBRETTO_ASCII=safe   swap quadrant glyphs for half blocks
  LIBRETTO_THEME=dark|light  force a palette instead of detecting
  LIBRETTO_ROOT=<path>  the payload clone; default ~/%[3]s
  CLAUDE_HOME=<path>    override Claude Code's root

  installed with:  go install github.com/pausf/libretto-automata/cmd/libretto@latest
  the payload is cloned to ~/%[3]s on the first command that needs it
`, version, n, BootstrapDir)
}
