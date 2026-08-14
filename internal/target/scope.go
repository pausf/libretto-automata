package target

import (
	"os"
	"path/filepath"
)

// Where the payload goes.
//
// Claude Code reads a project-local `.claude/` as well as the global one, so the
// same items can be installed for one repository without touching the
// configuration every other project shares.
//
// The two are never written in the same run. Two roots in one command is two
// answers to "where is my payload", and a bad day when they disagree.

// Scope names a destination.
type Scope string

const (
	// GlobalScope is ~/.claude, or CLAUDE_HOME when set.
	GlobalScope Scope = "global"

	// ProjectScope is <working directory>/.claude.
	ProjectScope Scope = "project"

	// CodexScope is ~/.agents, or AGENTS_HOME when set.
	CodexScope Scope = "codex"

	// OpencodeScope is ~/.config/opencode, or OPENCODE_HOME when set.
	OpencodeScope Scope = "opencode"
)

// Global is the machine-wide target.
func Global() Target { return NewClaude() }

// Project is the target inside a project directory.
//
// dir is where the user is, not where the repository being installed lives.
// Conflating the two would install libretto's payload into libretto's own
// `.claude/` regardless of where the command was run.
func Project(dir string) Target { return NewProject(dir) }

// Resolve returns the target for a scope. dir is only consulted for
// ProjectScope, and an empty dir falls back to the working directory.
// An unrecognised scope resolves to global, not to nothing.
func Resolve(s Scope, dir string) Target {
	switch s {
	case CodexScope:
		return NewCodex()
	case OpencodeScope:
		return NewOpencode()
	case ProjectScope:
		if dir == "" {
			if wd, err := os.Getwd(); err == nil {
				dir = wd
			}
		}
		return Project(dir)
	}
	return Global()
}

// Project is a project-local `.claude/` directory.
type ProjectTarget struct{ dir string }

// NewProject roots a target at <dir>/.claude.
func NewProject(dir string) ProjectTarget { return ProjectTarget{dir: dir} }

func (p ProjectTarget) Name() string { return "project" }

func (p ProjectTarget) Root() string {
	if p.dir == "" {
		return ""
	}
	return filepath.Join(p.dir, ".claude")
}

// Kinds matches the global target's. A project-local directory holds the same
// three things, and accepting a different set would mean an item installs in one
// scope and silently vanishes in the other.
func (p ProjectTarget) Kinds() []Kind { return []Kind{Skills, Agents, Commands} }

func (p ProjectTarget) Dir(k Kind) string {
	root := p.Root()
	if root == "" {
		return ""
	}
	return dirUnderRoot(root, k)
}

func (p ProjectTarget) Accepts(k Kind) bool { return accepts(p.Kinds(), k) }

// Exists reports whether the project has a `.claude/` yet.
//
// Absent means the project has not opted in, which is a state and not an error.
// Install creates the directory; it is not a precondition.
func (p ProjectTarget) Exists() bool {
	root := p.Root()
	if root == "" {
		return false
	}
	fi, err := os.Stat(root)
	return err == nil && fi.IsDir()
}
