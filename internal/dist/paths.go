package dist

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pausf/libretto-automata/internal/repo"
)

// currentLink is the name of the symlink the installed payload is reached through.
const currentLink = "current"

// Base is where installed payloads live, under the given home directory.
//
// It takes home rather than reading it, so no test can reach a real ~/.local/share — the
// same rule CLAUDE_HOME and LIBRETTO_ROOT already carry for the other two roots.
//
// ~/.local/share and not a dotdir in $HOME, which is the opposite of where the clone
// bootstrap put things and correct for the opposite reason: this *is* opaque application
// data. Nobody is expected to cd into it and edit a skill.
func Base(home string) string {
	return filepath.Join(home, ".local", "share", "libretto")
}

// VersionDir is where one version's payload is extracted.
func VersionDir(base, tag string) string { return filepath.Join(base, tag) }

// Current is the stable path ~/.claude links through.
//
// The indirection is what makes an upgrade cheap: `~/.claude/skills/write-spec` points at
// `<base>/current/skills/write-spec`, so activating a new version leaves every existing
// link resolving without touching one of them. Real files at the end of the chain, which is
// what ownership and linking require and what an embedded payload could never give.
func Current(base string) string { return filepath.Join(base, currentLink) }

// Activate points current at tag.
//
// **Symlink to a temporary name, then rename over current.** Never remove-then-symlink:
// os.Symlink fails when the target name already exists, so removing first is the obvious
// implementation — and it leaves a window in which current is absent, every link under
// ~/.claude resolves to nothing, and the user's whole payload has apparently vanished.
// Rename is atomic and has no such window.
//
// The version has to be on disk first. Pointing current at a directory that is not there
// would produce exactly the vanishing above, just permanently.
func Activate(base, tag string) error {
	dir := VersionDir(base, tag)
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return fmt.Errorf("cannot activate %s: %s is not there", tag, dir)
	}

	tmp := filepath.Join(base, ".current-next")
	// A leftover from an interrupted run would make Symlink fail below. Removing a link we
	// named ourselves is safe; this is never the live current.
	os.Remove(tmp)

	// Relative, so the tree can be moved or mounted elsewhere and still resolve.
	if err := os.Symlink(tag, tmp); err != nil {
		return err
	}
	if err := os.Rename(tmp, Current(base)); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// Active is the tag current points at, or "" when nothing is activated.
func Active(base string) string {
	target, err := os.Readlink(Current(base))
	if err != nil {
		return ""
	}
	return filepath.Base(target)
}

// Versions lists the version directories on disk, newest first.
//
// Sorted by the same numeric semver rules the rest of the project uses rather than
// lexically, so v0.10.0 comes before v0.9.0. Anything that is not a plain release tag is not
// a version and is ignored — which is also what keeps `current` and a stray file out of the
// list.
func Versions(base string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if !e.IsDir() || !repo.IsRelease(e.Name()) {
			continue
		}
		out = append(out, e.Name())
	}
	sort.Slice(out, func(i, j int) bool { return repo.IsNewer(out[i], out[j]) })
	return out, nil
}

// Prune removes every version directory except the ones named — and except the active one,
// whatever it was told.
//
// A caller that gets its keep-list wrong should lose a spare version, not the one that is
// running. Deleting the payload out from under a live installation is the kind of failure
// that cannot be apologised for afterwards.
func Prune(base string, keep ...string) error {
	protected := map[string]bool{}
	for _, k := range keep {
		protected[k] = true
	}
	if active := Active(base); active != "" {
		protected[active] = true
	}

	versions, err := Versions(base)
	if err != nil {
		return err
	}
	for _, v := range versions {
		if protected[v] {
			continue
		}
		if err := os.RemoveAll(VersionDir(base, v)); err != nil {
			return err
		}
	}
	return nil
}
