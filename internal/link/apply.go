package link

import (
	"errors"
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
// Nothing here removes a real directory, ever, and the only regular file it removes is
// one this tool generated and can prove it generated — by a marker naming a source
// inside the repo, for a kind the destination installs by transform. os.Remove on a
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

// ownedFor asks the right ownership question for this entry: the marker arm applies
// only to a kind whose target installs it by transform.
func ownedFor(repoRoot string, e Entry) bool {
	if e.GeneratedKind {
		return OwnedEither(repoRoot, e.DestPath)
	}
	return Owned(repoRoot, e.DestPath)
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
	if a.Entry.Generated != nil {
		if err := createGenerated(dest, a.Entry.Generated); err != nil {
			if errors.Is(err, os.ErrExist) {
				// The same refusal a symlink gets. os.Rename would have replaced it
				// silently, which made `create` destroy a file for a generated kind
				// where it reports a failure for a linked one — the opposite of what
				// this function promises.
				return Result{
					Action:  a,
					Refused: true,
					Err:     fmt.Errorf("%s appeared since the scan; not touching it", dest),
				}
			}
			return Result{Action: a, Err: err}
		}
		return Result{Action: a}
	}
	if err := os.Symlink(a.Entry.RepoPath, dest); err != nil {
		return Result{Action: a, Err: err}
	}
	return Result{Action: a}
}

// createGenerated writes content to dest, and fails with os.ErrExist rather than
// replacing anything already there.
//
// os.Rename replaces silently, so the Lstat check in `create` leaves a window in which
// a file that appears is destroyed. os.Link refuses when the destination exists, which
// is the same guarantee os.Symlink gives the linked path — so the two kinds behave
// alike instead of one being safe and the other destructive.
func createGenerated(dest string, content []byte) error {
	tmp, err := writeTemp(dest, content)
	if err != nil {
		return err
	}
	defer os.Remove(tmp)
	return os.Link(tmp, dest)
}

// writeGenerated replaces dest atomically. Used where replacing is the intent.
//
// Temp file then rename, so the destination goes from old bytes to new with nothing
// in between. Not a nicety: OpenCode throws on a malformed agent rather than skipping
// it, so a torn write does not degrade one agent — it breaks the host's config load.
//
// The temp file goes in the destination directory, never the system temp directory.
// os.Rename fails across filesystems and a target root may well be on another one.
// It is removed on any failure, so a scan never finds a half-written neighbour and
// reports it as a conflict.
func writeGenerated(dest string, content []byte) error {
	tmp, err := writeTemp(dest, content)
	if err != nil {
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// writeTemp writes content to a temporary file beside dest and returns its path.
//
// Beside dest, never in the system temp directory: os.Rename and os.Link both fail
// across filesystems, and a target root is often on a different one. Removed on any
// failure, so a scan never finds a half-written neighbour and calls it a conflict.
func writeTemp(dest string, content []byte) (string, error) {
	f, err := os.CreateTemp(filepath.Dir(dest), ".libretto-*")
	if err != nil {
		return "", err
	}
	tmp := f.Name()

	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(tmp)
		return "", err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return "", err
	}
	// CreateTemp makes the file 0o600. Items are readable, like every link's target.
	if err := os.Chmod(tmp, 0o644); err != nil {
		os.Remove(tmp)
		return "", err
	}
	return tmp, nil
}

// repoint replaces one of our own links with the right destination.
func repoint(repoRoot string, a Action) Result {
	dest := a.Entry.DestPath

	if !ownedFor(repoRoot, a.Entry) {
		return Result{
			Action:  a,
			Refused: true,
			Err:     fmt.Errorf("%s is no longer ours; not touching it", dest),
		}
	}

	// A generated file is replaced in place, atomically. No unlink first: removing
	// it would leave the destination absent for as long as the write takes, and the
	// rename below is what makes the swap invisible. Nothing in the old file is the
	// user's — every byte of it was derived — so there is nothing to preserve.
	if a.Entry.Generated != nil {
		if err := writeGenerated(dest, a.Entry.Generated); err != nil {
			return Result{Action: a, Err: err}
		}
		return Result{Action: a}
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

	if !ownedFor(repoRoot, a.Entry) {
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
