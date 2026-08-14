package link

import (
	"bytes"
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

	// Actual is where an existing symlink currently points, or the source a
	// generated file claims to come from. Empty when DestPath is neither.
	Actual string

	// GeneratedKind records that this entry's target installs this kind by transform.
	//
	// It is what scopes the marker arm of ownership. Without it `apply` would have to
	// ask a widened Owned about every entry, and a review found that offering to
	// delete a hand-written file in a destination this tool never generates into.
	GeneratedKind bool

	// Generated is the exact bytes that belong at DestPath, for a kind the target
	// installs by transform. Nil for a symlinked kind, which is what `apply` reads
	// to decide whether to write a file or make a link.
	//
	// It is carried on the entry rather than recomputed later because `classify`
	// already produced it to decide the state, and because `Plan` is a pure
	// function of entries with no target to ask — keeping it that way is what makes
	// the dangerous half testable in memory.
	Generated []byte
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

	// A kind the target installs by transform is judged by its content, not by a link
	// destination. Everything else falls through to the symlink path unchanged.
	if want, ok := transformFor(t, kind, item); ok {
		e.GeneratedKind = true
		return classifyGenerated(repoRoot, e, want)
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

// transformFor returns the bytes that should be at the destination, when this target
// installs this kind by transform.
//
// A source that cannot be transformed returns ok with nil content, which classify
// turns into a Conflict: we cannot say what belongs there, so nothing is touched. A
// target that does not transform returns false and the caller takes the symlink path.
func transformFor(t target.Target, kind target.Kind, item Item) ([]byte, bool) {
	tr, ok := t.(target.Transformer)
	if !ok || !tr.Transforms(kind) {
		// Either the target links everything, or it links this kind. Asking first is
		// what keeps a transforming target's linked kinds on the symlink path.
		return nil, false
	}
	content, err := os.ReadFile(item.Path)
	if err != nil {
		return nil, true
	}
	want, err := tr.Transform(kind, item.Path, content)
	if err != nil {
		return nil, true
	}
	return want, true
}

// classifyGenerated decides the state of one generated item.
//
// The five states carry it with nothing added. want is nil when the source could not
// be transformed, and the honest answer there is Conflict — the state that means
// "reported, never touched".
func classifyGenerated(repoRoot string, e Entry, want []byte) Entry {
	if want == nil {
		e.State = Conflict
		return e
	}
	e.Generated = want

	fi, err := os.Lstat(e.DestPath)
	if err != nil {
		e.State = Missing
		return e
	}
	if !fi.Mode().IsRegular() || !OwnedEither(repoRoot, e.DestPath) {
		// A directory, a symlink, or a file with no marker. Somebody put it there.
		e.State = Conflict
		if source, ok := GeneratedSource(e.DestPath); ok {
			e.Actual = source
		}
		return e
	}

	// Actual is where this thing claims to come from — the same meaning the field
	// carries for a symlink.
	if source, ok := GeneratedSource(e.DestPath); ok {
		e.Actual = source
	}

	have, err := os.ReadFile(e.DestPath)
	if err != nil {
		// Owned but unreadable. Ours to fix, and rewriting is the fix.
		e.State = WrongTarget
		return e
	}
	if bytes.Equal(have, want) {
		e.State = Linked
		return e
	}
	e.State = WrongTarget
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

	// The marker arm of ownership applies only to a kind this target generates. In
	// every other directory a regular file is somebody's own work, whatever it
	// contains.
	generatedKind := false
	if tr, ok := t.(target.Transformer); ok && tr.Transforms(kind) {
		generatedKind = true
	}

	var stale []Entry
	for _, de := range dirEntries {
		name := de.Name()
		if known[name] {
			continue
		}
		path := filepath.Join(dir, name)
		owned := Owned(repoRoot, path)
		if generatedKind {
			owned = OwnedEither(repoRoot, path)
		}
		if !owned {
			// Somebody else's entry. Not our business.
			continue
		}
		actual, isLink := LinkTarget(path)
		if !isLink {
			// A generated orphan: owned by its marker, with no item behind it. Same
			// concept as a broken link, same remedy, so the same state.
			actual, _ = GeneratedSource(path)
		}
		stale = append(stale, Entry{
			Target:        t.Name(),
			Kind:          kind,
			Name:          name,
			State:         Stale,
			DestPath:      path,
			Actual:        actual,
			GeneratedKind: generatedKind,
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
