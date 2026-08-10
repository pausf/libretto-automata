package link

import (
	"os"
	"path/filepath"
	"strings"
)

// Ownership is the single most consequential predicate in this program.
//
// It decides whether a symlink in a target directory belongs to this repo. Say
// "yes" about somebody else's link and `install` overwrites their work and
// `prune` deletes it. Say "no" about our own and the tool can never repair
// anything. Everything in apply.go is gated on it.
//
// It is deliberately conservative: anything it cannot prove is ours is treated as
// foreign.

// Owned reports whether path is a symlink resolving inside repoRoot.
//
// The destination need not exist: a link whose source was deleted is still ours,
// and recognising that is what makes `prune` possible.
func Owned(repoRoot, path string) bool {
	dest, ok := linkDestination(path)
	if !ok {
		return false
	}
	return within(normalise(repoRoot), normalise(dest))
}

// LinkTarget returns the normalised path a symlink resolves to, and whether path
// was a symlink at all.
func LinkTarget(path string) (string, bool) {
	dest, ok := linkDestination(path)
	if !ok {
		return "", false
	}
	return normalise(dest), true
}

// linkDestination reads a symlink and makes its destination absolute.
//
// A relative link is resolved against the directory holding the link, not the
// process working directory — getting that backwards is the classic way to
// misjudge ownership.
func linkDestination(path string) (string, bool) {
	fi, err := os.Lstat(path)
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		return "", false
	}

	dest, err := os.Readlink(path)
	if err != nil {
		return "", false
	}
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(filepath.Dir(path), dest)
	}
	return filepath.Clean(dest), true
}

// normalise resolves every symlink it can in a path and returns it absolute.
//
// This matters more than it looks. On macOS `t.TempDir()` hands back
// /var/folders/... while /var is itself a symlink to /private/var, and /tmp
// likewise points at /private/tmp. Two spellings of the same directory would
// compare unequal and the tool would call its own links foreign.
//
// The path need not exist. Missing trailing components are resolved as far as
// possible and the remainder is appended verbatim, so a broken link is still
// judged against the right root.
func normalise(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}

	remainder := ""
	current := abs
	for {
		if resolved, err := filepath.EvalSymlinks(current); err == nil {
			return filepath.Join(resolved, remainder)
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached the root without resolving anything.
			return abs
		}
		remainder = filepath.Join(filepath.Base(current), remainder)
		current = parent
	}
}

// within reports whether path is root itself or lies beneath it.
//
// It compares path segments rather than string prefixes. A prefix test would
// declare /repo-backup to be inside /repo, which is how a tool ends up deleting
// the wrong directory.
func within(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
