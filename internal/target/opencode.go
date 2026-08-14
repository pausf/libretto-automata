package target

import (
	"os"
	"path/filepath"
)

// EnvOpencodeHome overrides the OpenCode target's root directory.
//
// Libretto-only, for the same test safety EnvClaudeHome provides — OpenCode
// itself does not read this variable.
const EnvOpencodeHome = "OPENCODE_HOME"

// Opencode is OpenCode (opencode.ai), rooted at ~/.config/opencode. It accepts
// skills and commands: OpenCode discovers Claude-compatible SKILL.md directories
// under ~/.config/opencode/skills, and globs "{command,commands}/**/*.md" with
// symlinks followed, so the plural directory every kind already uses is one of the
// two names it looks for. Agents are absent — those need a frontmatter transform
// rather than a link.
type Opencode struct{ root string }

// NewOpencode resolves the root from OPENCODE_HOME, falling back to
// ~/.config/opencode. An unresolvable home directory yields an empty root,
// which Exists reports as unconfigured rather than panicking.
func NewOpencode() Opencode {
	if r := os.Getenv(EnvOpencodeHome); r != "" {
		return Opencode{root: r}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return Opencode{}
	}
	return Opencode{root: filepath.Join(home, ".config", "opencode")}
}

// NewOpencodeProject roots the opencode target at <dir>/.opencode — OpenCode
// reads a project's .opencode/skills (and .claude/skills, which the claude
// project target already serves). An empty dir is inert.
func NewOpencodeProject(dir string) Opencode {
	if dir == "" {
		return Opencode{}
	}
	return Opencode{root: filepath.Join(dir, ".opencode")}
}

func (o Opencode) Name() string { return "opencode" }
func (o Opencode) Root() string { return o.root }

func (o Opencode) Kinds() []Kind { return []Kind{Skills, Commands} }

func (o Opencode) Dir(k Kind) string {
	if o.root == "" {
		return ""
	}
	return dirUnderRoot(o.root, k)
}

func (o Opencode) Accepts(k Kind) bool { return accepts(o.Kinds(), k) }

// Exists reports whether the target root is present on disk — absent means
// the user does not have OpenCode configured, never an error.
func (o Opencode) Exists() bool {
	if o.root == "" {
		return false
	}
	fi, err := os.Stat(o.root)
	return err == nil && fi.IsDir()
}
