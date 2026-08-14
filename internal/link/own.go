package link

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/pausf/libretto-automata/internal/target"
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
//
// **This does not consider the generated-file marker, and that is the point.** A
// destination that never generates anything cannot own a generated file, so asking
// about a marker there would add a way for a hand-written file to be deleted and buy
// nothing. Callers that install a kind by transform ask OwnedEither instead, and only
// for that kind.
func Owned(repoRoot, path string) bool {
	dest, ok := linkDestination(path)
	if !ok {
		return false
	}
	return within(normalise(repoRoot), normalise(dest))
}

// OwnedEither is Owned widened by the marker arm, for a kind whose target installs it
// by transform.
//
// The narrowing matters more than the widening. A review found that a single Owned
// covering both arms made `prune --claude --yes` offer to delete a hand-written file in
// a Claude destination that happened to carry a marker line — a destination this tool
// never writes a generated file into. Gating on the kind removes that surface entirely
// and costs nothing: no destination that generates loses anything.
func OwnedEither(repoRoot, path string) bool {
	if dest, ok := linkDestination(path); ok {
		return within(normalise(repoRoot), normalise(dest))
	}
	return OwnedGenerated(repoRoot, path)
}

// OwnedGenerated reports whether path is a regular file this repo generated.
//
// Narrower than the symlink arm, deliberately. A symlink is ours when its
// destination lands anywhere inside the repo — a whole subtree. A generated file is
// ours only when it carries an exact marker naming an exact path that is itself
// inside the repo. Anything else is foreign: no marker, a marker outside the
// frontmatter block, a relative or empty marker, a directory, an unreadable file.
//
// A marker whose source no longer exists is still ours. Orphaned is `Stale`, which
// prune removes; treating it as foreign would strand it forever.
func OwnedGenerated(repoRoot, path string) bool {
	fi, err := os.Lstat(path)
	if err != nil || !fi.Mode().IsRegular() {
		return false
	}
	source, ok := GeneratedSource(path)
	if !ok {
		return false
	}

	root, src := normalise(repoRoot), normalise(source)
	if !within(root, src) {
		return false
	}

	// Strictly inside, never the root itself. A marker naming the repository is not
	// naming an item, and a review found both that and a marker naming a directory
	// inside it being accepted and then deleted by `prune --yes`.
	if src == root {
		return false
	}

	// A directory is not an item of a generated kind — those are single files. When
	// the path is gone we cannot tell, and that is deliberate: a marker whose source
	// was deleted must stay ours or `prune` could never remove the orphan.
	if fi, err := os.Stat(src); err == nil && fi.IsDir() {
		return false
	}
	return true
}

// GeneratedSource returns the source path a generated file claims to come from, and
// whether it made such a claim at all. The mirror of LinkTarget, so link-state can
// tell "ours, wrong content" from "ours, source gone" without parsing again.
//
// Only the frontmatter block is read — the region between an opening `---` on line 1
// and the next `---`. A marker key in an agent's prose is prose, not a claim of
// ownership, and reading the whole file would let one be typed into a body and grant
// ownership of somebody's file.
//
// The value must be absolute. A symlink may be relative because there is an
// unambiguous base to resolve it against — the directory holding the link. A marker
// has none: the file sits in the target, so resolving a relative source against it
// would name something inside the target rather than inside the repo.
func GeneratedSource(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	// Frontmatter is short by construction. A bounded buffer keeps a pathological
	// single-line file from being read into memory to answer a predicate.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	if !sc.Scan() || sc.Text() != "---" {
		return "", false
	}
	for sc.Scan() {
		line := sc.Text()
		if line == "---" {
			return "", false
		}
		value, found := strings.CutPrefix(line, target.MarkerKey+":")
		if !found {
			continue
		}
		source := unquoteYAML(strings.TrimSpace(value))
		if source == "" || !filepath.IsAbs(source) {
			return "", false
		}
		return source, true
	}
	// The frontmatter never closed, so nothing in it can be trusted.
	return "", false
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

// unquoteYAML undoes quoteYAML for a double-quoted scalar, and returns anything else
// unchanged.
//
// Both forms are accepted on read. The transform always quotes, but a bare value is
// still a legitimate marker — and refusing one would mean a file this tool wrote
// before the quoting landed stopped being recognised as ours, which turns an owned
// file into a foreign one that nothing can clean up.
func unquoteYAML(s string) string {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return s
	}
	inner := s[1 : len(s)-1]
	var b strings.Builder
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' && i+1 < len(inner) {
			i++
		}
		b.WriteByte(inner[i])
	}
	return b.String()
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
