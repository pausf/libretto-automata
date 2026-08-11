package dist

import (
	"os"
	"path/filepath"
	"testing"
)

// Base takes the home directory rather than reading it, so no test can reach a real
// ~/.local/share.
func TestBaseIsUnderLocalShare(t *testing.T) {
	home := t.TempDir()

	if got, want := Base(home), filepath.Join(home, ".local", "share", "libretto"); got != want {
		t.Errorf("Base = %q, want %q", got, want)
	}
}

func TestCurrentPointsAtTheActivatedVersion(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.4.0")

	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}

	resolved, err := filepath.EvalSymlinks(Current(base))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(VersionDir(base, "v0.4.0"))
	if err != nil {
		t.Fatal(err)
	}
	if resolved != want {
		t.Errorf("current resolves to %q, want %q", resolved, want)
	}

	// And it resolves through to a real file, which is what the whole symlink chain exists
	// for: ~/.claude links through current, and a link needs a real file at the end.
	if _, err := os.Stat(filepath.Join(Current(base), "skills", "marker.md")); err != nil {
		t.Errorf("current does not resolve through to the payload: %v", err)
	}
}

// The swap has to be a rename, never remove-then-symlink. os.Symlink fails when the target
// name already exists, so the obvious implementation is to remove `current` first — and
// that leaves a window in which ~/.claude points at nothing and every skill has vanished.
func TestActivateIsAtomicAndNeverLeavesCurrentMissing(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.3.0")
	mkVersion(t, base, "v0.4.0")

	if err := Activate(base, "v0.3.0"); err != nil {
		t.Fatal(err)
	}

	// Activating over an existing current must succeed rather than failing with "file
	// exists" — which is what a bare os.Symlink does.
	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatalf("activating over an existing current: %v", err)
	}

	// current is a symlink throughout, never a directory and never absent.
	info, err := os.Lstat(Current(base))
	if err != nil {
		t.Fatalf("current is missing after the swap: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("current is not a symlink")
	}

	target, err := os.Readlink(Current(base))
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(target) != "v0.4.0" {
		t.Errorf("current points at %q, want v0.4.0", target)
	}

	// No temporary link was left behind for a later run to trip over.
	entries, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if name := e.Name(); name != "v0.3.0" && name != "v0.4.0" && name != "current" {
			t.Errorf("the swap left %q behind", name)
		}
	}
}

func TestActivateRefusesAVersionThatIsNotThere(t *testing.T) {
	base := t.TempDir()

	if err := Activate(base, "v9.9.9"); err == nil {
		t.Error("Activate accepted a version with no directory")
	}
	if _, err := os.Lstat(Current(base)); !os.IsNotExist(err) {
		t.Error("a refused Activate created current anyway")
	}
}

// Two versions on disk is a few hundred kilobytes of markdown. Keeping every version ever
// installed is a directory nobody prunes; keeping none makes rollback a download.
func TestOnlyThePreviousVersionIsKept(t *testing.T) {
	base := t.TempDir()
	for _, v := range []string{"v0.1.0", "v0.2.0", "v0.3.0", "v0.4.0"} {
		mkVersion(t, base, v)
	}
	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}

	if err := Prune(base, "v0.4.0", "v0.3.0"); err != nil {
		t.Fatal(err)
	}

	for _, kept := range []string{"v0.4.0", "v0.3.0"} {
		if _, err := os.Stat(VersionDir(base, kept)); err != nil {
			t.Errorf("%s was removed: %v", kept, err)
		}
	}
	for _, gone := range []string{"v0.1.0", "v0.2.0"} {
		if _, err := os.Stat(VersionDir(base, gone)); !os.IsNotExist(err) {
			t.Errorf("%s survived the prune", gone)
		}
	}
	// And current is untouched — pruning is not activating.
	if target, _ := os.Readlink(Current(base)); filepath.Base(target) != "v0.4.0" {
		t.Errorf("prune moved current to %q", target)
	}
}

// Prune must never remove what current points at, whatever it was told to keep. A caller
// that gets its keep-list wrong should lose a spare version, not the running one.
func TestPruneNeverRemovesTheActiveVersion(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.3.0")
	mkVersion(t, base, "v0.4.0")
	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}

	if err := Prune(base, "v0.3.0"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(VersionDir(base, "v0.4.0")); err != nil {
		t.Errorf("prune removed the active version: %v", err)
	}
}

// Rolling back is the swap and nothing else. The previous version is already on disk, which
// is the entire reason it is kept.
func TestRollbackIsASwapNotADownload(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.3.0")
	mkVersion(t, base, "v0.4.0")

	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}
	if err := Activate(base, "v0.3.0"); err != nil {
		t.Fatalf("rolling back: %v", err)
	}

	if target, _ := os.Readlink(Current(base)); filepath.Base(target) != "v0.3.0" {
		t.Errorf("current points at %q after the rollback, want v0.3.0", target)
	}
	if _, err := os.Stat(filepath.Join(Current(base), "skills", "marker.md")); err != nil {
		t.Errorf("the rolled-back version does not resolve: %v", err)
	}
}

func TestVersionsListsWhatIsOnDisk(t *testing.T) {
	base := t.TempDir()
	mkVersion(t, base, "v0.4.0")
	mkVersion(t, base, "v0.3.0")
	if err := Activate(base, "v0.4.0"); err != nil {
		t.Fatal(err)
	}
	// A stray file and the current link must not be mistaken for versions.
	if err := os.WriteFile(filepath.Join(base, "notes.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Versions(base)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "v0.4.0" || got[1] != "v0.3.0" {
		t.Errorf("Versions = %v, want [v0.4.0 v0.3.0] newest first", got)
	}
}

// mkVersion creates a version directory with something in it, so the symlink chain has a
// real file to resolve through.
func mkVersion(t *testing.T, base, tag string) {
	t.Helper()
	dir := filepath.Join(VersionDir(base, tag), "skills")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "marker.md"), []byte(tag+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
