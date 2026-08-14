package target

import (
	"os"
	"path/filepath"
)

// Where the payload goes: a tool crossed with a scope.
//
// The two axes are orthogonal and deliberately never fused into one list. A flat
// list of destinations grows as tools × scopes; two axes grow as tools + scopes,
// and the panel's strip stays one row per tool whatever happens to scopes.
//
// A command still acts on exactly one (tool, scope) pair per run. Two roots in one
// command is two answers to "where is my payload", and a bad day when they disagree.

// Tool names an agent installation family.
type Tool string

const (
	// ClaudeTool is Claude Code.
	ClaudeTool Tool = "claude"

	// CodexTool is OpenAI Codex CLI.
	CodexTool Tool = "codex"

	// OpencodeTool is OpenCode.
	OpencodeTool Tool = "opencode"
)

// Tools is the display order: Claude first, because it is configured on every
// machine this tool has users on today.
var Tools = []Tool{ClaudeTool, CodexTool, OpencodeTool}

// Scope names which side of a tool a command acts on.
type Scope string

const (
	// GlobalScope is the tool's machine-wide root.
	GlobalScope Scope = "global"

	// ProjectScope is the tool's directory inside the working directory.
	ProjectScope Scope = "project"
)

// Global is the machine-wide Claude target.
func Global() Target { return NewClaude() }

// Project is the Claude target inside a project directory.
//
// dir is where the user is, not where the repository being installed lives.
// Conflating the two would install libretto's payload into libretto's own
// `.claude/` regardless of where the command was run.
func Project(dir string) Target { return NewProject(dir) }

// Resolve returns the target for a tool and scope. dir is only consulted for
// ProjectScope, and an empty dir falls back to the working directory. An
// unrecognised tool resolves to Claude and an unrecognised scope to global —
// a typo must not silently produce a rootless target that writes nowhere and
// reports success.
func Resolve(tool Tool, s Scope, dir string) Target {
	if s == ProjectScope {
		if dir == "" {
			if wd, err := os.Getwd(); err == nil {
				dir = wd
			}
		}
		switch tool {
		case CodexTool:
			return NewCodexProject(dir)
		case OpencodeTool:
			return NewOpencodeProject(dir)
		}
		return Project(dir)
	}
	switch tool {
	case CodexTool:
		return NewCodex()
	case OpencodeTool:
		return NewOpencode()
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
