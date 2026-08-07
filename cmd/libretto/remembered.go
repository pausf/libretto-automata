package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pausf/libretto-automata/internal/target"
)

// Where the panel's last destination is remembered.
//
// One word, one file, at the global target's root — so CLAUDE_HOME relocates it and
// the test suite never touches a real ~/.claude. Deliberately *not* under `skills/`,
// `agents/` or `commands/`, and a regular file rather than a symlink this tool
// created, so `prune` and `uninstall` never have any reason to look at it.
//
// One machine-wide value, not one per directory. Decided with the case against it
// stated — the panel can open on `project` inside a repository that never asked for
// it — and bounded by the fact that no subcommand reads this. It decides which side
// the panel *opens* on and nothing else.
const preferenceFile = "libretto-scope"

func preferencePath() string {
	root := target.Global().Root()
	if root == "" {
		return ""
	}
	return filepath.Join(root, preferenceFile)
}

// rememberedScope is the destination the panel was left on.
//
// Absent, empty, unreadable, or holding anything but the two words all give global.
// That is target.Resolve's rule and it is followed rather than reinvented: a value
// nobody recognises must not produce a destination nobody chose.
func rememberedScope() target.Scope {
	path := preferencePath()
	if path == "" {
		return target.GlobalScope
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return target.GlobalScope
	}
	if target.Scope(strings.TrimSpace(string(b))) == target.ProjectScope {
		return target.ProjectScope
	}
	return target.GlobalScope
}

// remember writes the destination, and says nothing when it cannot.
//
// The error is dropped on purpose. Remembering is a convenience, and a convenience
// that can fail the switch it decorates is worse than no convenience — the panel has
// already moved by the time this runs, so there is nothing left to abort.
//
// ponytail: one word, no format. The moment a second thing wants remembering this
// becomes a real file with a real shape, not a delimiter invented in place.
func remember(s target.Scope) {
	path := preferencePath()
	if path == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, []byte(string(s)+"\n"), 0o644)
}
