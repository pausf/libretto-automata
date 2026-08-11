package main

import (
	"context"
	"fmt"
	"io"

	"github.com/pausf/libretto-automata/internal/repo"
)

// carriesNotice is which subcommands end by saying a newer release exists.
//
// `doctor` is absent because it already says something on every path, live and
// uncached — a second line would print the same fact twice, once stale. `update` is
// absent because announcing the thing it is doing is noise. `version` and `help`
// answer before the payload is even located, so neither has a root to ask about.
func carriesNotice(cmd string) bool {
	switch cmd {
	case "status", "preview", "install", "prune", "uninstall", "models":
		return true
	}
	return false
}

// subcommandNotice is the line, or "" when there is nothing to say.
//
// Cached, through the same path the panel uses: a subcommand is not a diagnostic and
// does not get to pay for a network call on every run. Silence on a failed check is
// the panel's rule too — "could not check for a newer release" printed by a command
// about something else is the noise that gets a notice ignored.
func subcommandNotice(root, running string) string {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()

	latest, err := askLatest(ctx, root)
	if err != nil || !repo.IsNewer(latest, running) {
		return ""
	}
	return fmt.Sprintf("%s → %s available · run `%s update`", running, latest, invokedAs())
}

// noticeSource is a var so the tests can answer without reaching a remote.
var noticeSource = subcommandNotice

// noticeAfter writes the notice for cmd, if cmd carries one and there is one to write.
//
// Never coloured. The panel's row is gold because it sits inside a rendered frame; a
// line appended to a command's output is not, and piped `status` carries no escape
// codes by an existing promise.
func noticeAfter(w io.Writer, cmd, root, running string) {
	if !carriesNotice(cmd) {
		return
	}
	if line := noticeSource(root, running); line != "" {
		fmt.Fprintln(w, line)
	}
}
