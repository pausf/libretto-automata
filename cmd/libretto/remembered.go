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

// rememberedScope is the (tool, scope) pair the panel was left on.
//
// Absent, empty, unreadable, or holding anything unrecognised gives claude/global.
// That is target.Resolve's rule and it is followed rather than reinvented: a value
// nobody recognises must not produce a destination nobody chose. The legacy
// one-word forms survive: "project" reads as claude/project.
func rememberedScope() (target.Tool, target.Scope) {
	tool, scope := target.ClaudeTool, target.GlobalScope

	path := preferencePath()
	if path == "" {
		return tool, scope
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return tool, scope
	}

	for _, word := range strings.Fields(string(b)) {
		for _, t := range target.Tools {
			if word == string(t) {
				tool = t
			}
		}
		if word == string(target.ProjectScope) {
			scope = target.ProjectScope
		}
	}
	return tool, scope
}

// remember writes the destination, and says nothing when it cannot.
//
// The error is dropped on purpose. Remembering is a convenience, and a convenience
// that can fail the switch it decorates is worse than no convenience — the panel has
// already moved by the time this runs, so there is nothing left to abort.
//
// ponytail: two space-separated words, no format beyond that. The second thing
// that wanted remembering arrived (the scope axis); a third gets a real file.
func remember(t target.Tool, s target.Scope) {
	path := preferencePath()
	if path == "" {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(path, []byte(string(t)+" "+string(s)+"\n"), 0o644)
}
