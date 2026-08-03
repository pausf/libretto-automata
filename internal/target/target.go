// Package target describes the agent installations that consume repo items.
//
// Nothing outside this package may assume where a target lives on disk or that
// every target accepts every kind. See docs/SPEC.md R10.
package target

import "path/filepath"

// Kind is a category of item the repo can hold.
type Kind string

const (
	Skills   Kind = "skills"
	Agents   Kind = "agents"
	Commands Kind = "commands"
)

// ItemsAreDirs reports whether items of this kind are directories rather than
// files. Skills are directories containing a SKILL.md; agents and commands are
// single .md files.
func (k Kind) ItemsAreDirs() bool { return k == Skills }

// Ext is the required file extension for items of this kind, or "" when items
// are directories.
func (k Kind) Ext() string {
	if k.ItemsAreDirs() {
		return ""
	}
	return ".md"
}

// Target is one agent installation.
type Target interface {
	// Name is the short identifier shown in the UI, e.g. "claude".
	Name() string

	// Root is the absolute path of the installation directory.
	Root() string

	// Kinds are the item categories this target accepts, in display order.
	Kinds() []Kind

	// Dir is the absolute directory where items of this kind are installed.
	Dir(Kind) string

	// Accepts reports whether this target takes items of the given kind.
	Accepts(Kind) bool
}

// dirUnderRoot is the common Dir implementation: one subdirectory per kind,
// named after the kind. Shared by every target that follows that convention.
func dirUnderRoot(root string, k Kind) string {
	return filepath.Join(root, string(k))
}

// accepts is the common Accepts implementation.
func accepts(kinds []Kind, k Kind) bool {
	for _, have := range kinds {
		if have == k {
			return true
		}
	}
	return false
}
