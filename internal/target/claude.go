package target

import (
	"os"
	"path/filepath"
)

// EnvClaudeHome overrides Claude Code's root directory.
//
// This is not a user feature. It is what lets the whole test suite run against
// t.TempDir() so nothing ever reads or writes the real ~/.claude.
const EnvClaudeHome = "CLAUDE_HOME"

// Claude is Claude Code, rooted at ~/.claude.
type Claude struct{ root string }

// NewClaude resolves the root from CLAUDE_HOME, falling back to ~/.claude.
// An unresolvable home directory yields an empty root, which Exists reports as
// unconfigured rather than panicking.
func NewClaude() Claude {
	if r := os.Getenv(EnvClaudeHome); r != "" {
		return Claude{root: r}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Claude{}
	}
	return Claude{root: filepath.Join(home, ".claude")}
}

func (c Claude) Name() string { return "claude" }
func (c Claude) Root() string { return c.root }

func (c Claude) Kinds() []Kind { return []Kind{Skills, Agents, Commands} }

func (c Claude) Dir(k Kind) string {
	if c.root == "" {
		return ""
	}
	return dirUnderRoot(c.root, k)
}

func (c Claude) Accepts(k Kind) bool { return accepts(c.Kinds(), k) }

// Exists reports whether the target root is present on disk. A target that is
// not present is shown as unconfigured rather than treated as an error — the
// user may simply not have that agent installed.
func (c Claude) Exists() bool {
	if c.root == "" {
		return false
	}
	fi, err := os.Stat(c.root)
	return err == nil && fi.IsDir()
}

// All returns every known target, in display order.
//
// One entry today. The slice is what makes adding Codex a new file instead of a
// refactor — see docs/DESIGN.md.
func All() []Target { return []Target{NewClaude()} }
