package target

import (
	"os"
	"path/filepath"
)

// EnvAgentsHome overrides the Codex target's root directory.
//
// Libretto-only, for the same test safety EnvClaudeHome provides — Codex
// itself does not read this variable.
const EnvAgentsHome = "AGENTS_HOME"

// Codex is OpenAI Codex CLI, rooted at ~/.agents. It accepts only skills:
// Codex discovers Claude-compatible SKILL.md directories under
// ~/.agents/skills, and OpenCode reads the same path, so one link here
// serves both tools.
type Codex struct{ root string }

// NewCodex resolves the root from AGENTS_HOME, falling back to ~/.agents.
// An unresolvable home directory yields an empty root, which Exists reports
// as unconfigured rather than panicking.
func NewCodex() Codex {
	if r := os.Getenv(EnvAgentsHome); r != "" {
		return Codex{root: r}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Codex{}
	}
	return Codex{root: filepath.Join(home, ".agents")}
}

func (c Codex) Name() string { return "codex" }
func (c Codex) Root() string { return c.root }

func (c Codex) Kinds() []Kind { return []Kind{Skills} }

func (c Codex) Dir(k Kind) string {
	if c.root == "" {
		return ""
	}
	return dirUnderRoot(c.root, k)
}

func (c Codex) Accepts(k Kind) bool { return accepts(c.Kinds(), k) }

// Exists reports whether the target root is present on disk — absent means
// the user does not have Codex configured, never an error.
func (c Codex) Exists() bool {
	if c.root == "" {
		return false
	}
	fi, err := os.Stat(c.root)
	return err == nil && fi.IsDir()
}
