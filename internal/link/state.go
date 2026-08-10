package link

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pausf/libretto-automata/internal/target"
)

// State is what one item looks like in one target. See docs/SPEC.md R4.
type State string

const (
	// Linked: an owned symlink pointing at the right item. Nothing to do.
	Linked State = "linked"

	// Missing: the repo has the item, the target has nothing. Install creates it.
	Missing State = "missing"

	// WrongTarget: an owned symlink pointing somewhere else inside the repo.
	// Ours to fix, so install repoints it.
	WrongTarget State = "wrong target"

	// Conflict: something is in the way that is not ours — a real file, a real
	// directory, or a symlink leaving the repo. Never touched, always reported.
	Conflict State = "conflict"

	// Stale: an owned symlink with no matching item in the repo. Prune's job.
	//
	// This subsumes "broken link": an owned link whose source was deleted has no
	// matching item, so it lands here. It also covers an owned link aimed at
	// something inside the repo that was never an item. One concept, because the
	// remedy is identical.
	Stale State = "stale"
)

// NeedsAttention reports whether a state calls for action or a warning. Used to
// decide a command's exit code.
func (s State) NeedsAttention() bool { return s != Linked }

// Entry is one item's situation in one target.
type Entry struct {
	Target string
	Kind   target.Kind
	Name   string
	State  State

	// RepoPath is the item's path in the repo. Empty for Stale, where no item
	// exists.
	RepoPath string

	// DestPath is the path in the target.
	DestPath string

	// Actual is where an existing symlink currently points. Empty when DestPath
	// is not a symlink.
	Actual string
}

// Scan classifies every item of every kind the target accepts, and finds owned
// links in the target with no item behind them.
//
// It never writes. Results are ordered by kind, then by name, so output is stable
// between runs.
func Scan(repoRoot string, t target.Target) ([]Entry, error) {
	var entries []Entry

	for _, kind := range t.Kinds() {
		items, err := Items(repoRoot, kind)
		if err != nil {
			return nil, err
		}

		known := make(map[string]bool, len(items))
		for _, item := range items {
			known[item.Name] = true
			entries = append(entries, classify(repoRoot, t, kind, item))
		}

		stale, err := staleIn(repoRoot, t, kind, known)
		if err != nil {
			return nil, err
		}
		entries = append(entries, stale...)
	}

	return entries, nil
}

// classify decides the state of one repo item in one target.
func classify(repoRoot string, t target.Target, kind target.Kind, item Item) Entry {
	dest := filepath.Join(t.Dir(kind), item.Name)
	e := Entry{
		Target:   t.Name(),
		Kind:     kind,
		Name:     item.Name,
		RepoPath: item.Path,
		DestPath: dest,
	}

	if _, err := os.Lstat(dest); err != nil {
		e.State = Missing
		return e
	}

	actual, isLink := LinkTarget(dest)
	switch {
	case !isLink:
		// A real file or directory. Somebody put it there deliberately.
		e.State = Conflict
	case !Owned(repoRoot, dest):
		// A symlink leaving the repo is also a deliberate choice, not ours to
		// overwrite.
		e.State, e.Actual = Conflict, actual
	case actual == normalise(item.Path):
		e.State, e.Actual = Linked, actual
	default:
		e.State, e.Actual = WrongTarget, actual
	}
	return e
}

// staleIn finds owned links in a target directory that no current item explains.
func staleIn(repoRoot string, t target.Target, kind target.Kind, known map[string]bool) ([]Entry, error) {
	dir := t.Dir(kind)
	if dir == "" {
		return nil, nil
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// The target has no directory for this kind yet. Nothing can be
			// stale in a directory that does not exist.
			return nil, nil
		}
		return nil, err
	}

	var stale []Entry
	for _, de := range dirEntries {
		name := de.Name()
		if known[name] {
			continue
		}
		path := filepath.Join(dir, name)
		if !Owned(repoRoot, path) {
			// Somebody else's entry. Not our business.
			continue
		}
		actual, _ := LinkTarget(path)
		stale = append(stale, Entry{
			Target:   t.Name(),
			Kind:     kind,
			Name:     name,
			State:    Stale,
			DestPath: path,
			Actual:   actual,
		})
	}

	sort.Slice(stale, func(i, j int) bool { return stale[i].Name < stale[j].Name })
	return stale, nil
}

// Tally counts entries by state.
func Tally(entries []Entry) map[State]int {
	counts := make(map[State]int, 5)
	for _, e := range entries {
		counts[e.State]++
	}
	return counts
}

// ByState returns the entries in a given state, preserving scan order.
func ByState(entries []Entry, s State) []Entry {
	var out []Entry
	for _, e := range entries {
		if e.State == s {
			out = append(out, e)
		}
	}
	return out
}
