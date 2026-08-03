package link

import (
	"fmt"
	"os"
	"path/filepath"
)

// Applying is the only code in this program that writes. Everything it does is
// gated on Owned, re-checked here rather than trusted from the scan.
//
// The re-check is the point. A plan is a snapshot of a moment; between the scan
// and the write, another tool may have replaced a link, a user may have dropped a
// real directory in its place. Acting on a stale classification is how a tool
// deletes something it never examined.
//
// Nothing here removes a real file or a real directory, ever. os.Remove on a
// symlink removes the link and never its destination, and a path that is not a
// symlink fails the ownership re-check before it gets that far.

// Result is what happened to one action.
type Result struct {
	Action Action

	// Err is nil when the action succeeded or was a Skip.
	Err error

	// Refused is set when the action was abandoned because the filesystem no
	// longer matched the plan. This is a safety outcome, not a bug: something
	// changed and the conservative move was to leave it alone.
	Refused bool
}

// Done reports whether the filesystem actually changed.
func (r Result) Done() bool { return r.Err == nil && !r.Refused && r.Action.Writes() }

// Apply executes a plan and reports what happened to every action.
//
// It does not stop at the first failure. One item the tool cannot touch is no
// reason to abandon the nine it can, and a partial result that says exactly which
// is more useful than an early exit that says neither.
func Apply(repoRoot string, actions []Action) []Result {
	results := make([]Result, 0, len(actions))
	for _, a := range actions {
		results = append(results, apply(repoRoot, a))
	}
	return results
}

func apply(repoRoot string, a Action) Result {
	switch a.Act {
	case Skip:
		return Result{Action: a}
	case Create:
		return create(a)
	case Repoint:
		return repoint(repoRoot, a)
	case Remove:
		return remove(repoRoot, a)
	default:
		return Result{Action: a, Err: fmt.Errorf("unknown action %q", a.Act)}
	}
}

// create makes a link where the scan found nothing.
//
// If something is there now, the scan was stale and this is not ours to resolve.
func create(a Action) Result {
	dest := a.Entry.DestPath

	if _, err := os.Lstat(dest); err == nil {
		return Result{
			Action:  a,
			Refused: true,
			Err:     fmt.Errorf("%s appeared since the scan; not touching it", dest),
		}
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return Result{Action: a, Err: err}
	}
	if err := os.Symlink(a.Entry.RepoPath, dest); err != nil {
		return Result{Action: a, Err: err}
	}
	return Result{Action: a}
}

// repoint replaces one of our own links with the right destination.
func repoint(repoRoot string, a Action) Result {
	dest := a.Entry.DestPath

	if !Owned(repoRoot, dest) {
		return Result{
			Action:  a,
			Refused: true,
			Err:     fmt.Errorf("%s is no longer ours; not touching it", dest),
		}
	}

	if err := os.Remove(dest); err != nil {
		return Result{Action: a, Err: err}
	}
	if err := os.Symlink(a.Entry.RepoPath, dest); err != nil {
		return Result{Action: a, Err: err}
	}
	return Result{Action: a}
}

// remove deletes one of our own links whose item is gone.
func remove(repoRoot string, a Action) Result {
	dest := a.Entry.DestPath

	if !Owned(repoRoot, dest) {
		return Result{
			Action:  a,
			Refused: true,
			Err:     fmt.Errorf("%s is not ours; not removing it", dest),
		}
	}

	if err := os.Remove(dest); err != nil {
		return Result{Action: a, Err: err}
	}
	return Result{Action: a}
}

// Summarise counts results into done, refused and failed.
func Summarise(results []Result) (done, refused, failed int) {
	for _, r := range results {
		switch {
		case r.Refused:
			refused++
		case r.Err != nil:
			failed++
		case r.Done():
			done++
		}
	}
	return done, refused, failed
}
