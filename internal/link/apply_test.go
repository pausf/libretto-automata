package link

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// installInto scans, plans and applies, returning the results.
func installInto(t *testing.T, repo string, tg target.Target) []Result {
	t.Helper()
	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	return Apply(repo, Plan(entries))
}

func mustBeLinkTo(t *testing.T, path, want string) {
	t.Helper()
	got, ok := LinkTarget(path)
	if !ok {
		t.Fatalf("%s is not a symlink", path)
	}
	if got != normalise(want) {
		t.Fatalf("%s points at %s, want %s", path, got, want)
	}
}

func TestApplyCreatesMissingLinks(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md", "commands/beta.md")

	results := installInto(t, repo, tg)

	done, refused, failed := Summarise(results)
	if done != 2 || refused != 0 || failed != 0 {
		t.Fatalf("done=%d refused=%d failed=%d, want 2/0/0", done, refused, failed)
	}

	mustBeLinkTo(t, filepath.Join(tg.Dir(target.Skills), "alpha"),
		filepath.Join(repo, "skills", "alpha"))
	mustBeLinkTo(t, filepath.Join(tg.Dir(target.Commands), "beta.md"),
		filepath.Join(repo, "commands", "beta.md"))
}

// R9. The important property: install is safe to run again and again.
func TestApplyIsIdempotent(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md")

	if done, _, _ := Summarise(installInto(t, repo, tg)); done != 1 {
		t.Fatal("first install did not create the link")
	}

	second := installInto(t, repo, tg)
	if len(second) != 0 {
		t.Fatalf("second install planned %d action(s), want none: %v", len(second), second)
	}
}

func TestApplyRepointsOurOwnWrongLink(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md", "skills/other/SKILL.md")

	dest := filepath.Join(tg.Dir(target.Skills), "alpha")
	symlink(t, filepath.Join(repo, "skills", "other"), dest) // ours, wrong item

	if done, _, _ := Summarise(installInto(t, repo, tg)); done < 1 {
		t.Fatal("install did not repoint the link")
	}
	mustBeLinkTo(t, dest, filepath.Join(repo, "skills", "alpha"))
}

func TestApplyNeverTouchesAConflict(t *testing.T) {
	t.Run("a real directory survives untouched", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/alpha/SKILL.md")

		dest := filepath.Join(tg.Dir(target.Skills), "alpha")
		if err := os.MkdirAll(dest, 0o755); err != nil {
			t.Fatal(err)
		}
		keep := filepath.Join(dest, "somebody-elses-work.md")
		if err := os.WriteFile(keep, []byte("mine"), 0o644); err != nil {
			t.Fatal(err)
		}

		results := installInto(t, repo, tg)
		if done, _, failed := Summarise(results); done != 0 || failed != 0 {
			t.Fatalf("done=%d failed=%d, want 0/0 — a conflict must be left alone", done, failed)
		}
		if _, err := os.Stat(keep); err != nil {
			t.Fatalf("install destroyed a real directory: %v", err)
		}
	})

	t.Run("a foreign symlink survives untouched", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/alpha/SKILL.md")
		outside := t.TempDir()

		dest := filepath.Join(tg.Dir(target.Skills), "alpha")
		symlink(t, outside, dest)

		installInto(t, repo, tg)

		mustBeLinkTo(t, dest, outside)
	})
}

func TestApplyRefusesWhenTheWorldChangedAfterTheScan(t *testing.T) {
	t.Run("create refuses if something appeared", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/alpha/SKILL.md")

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		plan := Plan(entries) // says: create

		// Somebody drops a real directory in the way between scan and apply.
		dest := filepath.Join(tg.Dir(target.Skills), "alpha")
		if err := os.MkdirAll(dest, 0o755); err != nil {
			t.Fatal(err)
		}

		results := Apply(repo, plan)
		if _, refused, _ := Summarise(results); refused != 1 {
			t.Fatalf("apply did not refuse a destination that appeared: %+v", results)
		}
		if fi, err := os.Lstat(dest); err != nil || !fi.IsDir() {
			t.Fatal("apply replaced the directory it should have refused")
		}
	})

	t.Run("repoint refuses if the link became foreign", func(t *testing.T) {
		repo, tg := sandbox(t, "skills/alpha/SKILL.md", "skills/other/SKILL.md")
		outside := t.TempDir()

		dest := filepath.Join(tg.Dir(target.Skills), "alpha")
		symlink(t, filepath.Join(repo, "skills", "other"), dest)

		entries, err := Scan(repo, tg)
		if err != nil {
			t.Fatal(err)
		}
		plan := Plan(entries) // says: repoint, it is ours

		// Another tool takes it over before we get there.
		if err := os.Remove(dest); err != nil {
			t.Fatal(err)
		}
		symlink(t, outside, dest)

		results := Apply(repo, plan)
		if _, refused, _ := Summarise(results); refused != 1 {
			t.Fatalf("apply repointed a link that stopped being ours: %+v", results)
		}
		mustBeLinkTo(t, dest, outside)
	})
}

func TestPruneRemovesOnlyOurOwnStaleLinks(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md")
	outside := t.TempDir()

	dir := tg.Dir(target.Skills)
	stale := filepath.Join(dir, "gone")
	foreign := filepath.Join(dir, "theirs")

	// Ours, pointing at an item that no longer exists.
	symlink(t, filepath.Join(repo, "skills", "gone"), stale)
	// Somebody else's link.
	symlink(t, outside, foreign)

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}

	plan := PrunePlan(entries)
	if len(plan) != 1 {
		t.Fatalf("prune planned %d action(s), want 1: %v", len(plan), plan)
	}

	if done, _, failed := Summarise(Apply(repo, plan)); done != 1 || failed != 0 {
		t.Fatalf("done=%d failed=%d, want 1/0", done, failed)
	}

	if _, err := os.Lstat(stale); !os.IsNotExist(err) {
		t.Fatal("prune left our stale link behind")
	}
	if _, err := os.Lstat(foreign); err != nil {
		t.Fatal("prune deleted somebody else's link")
	}
}

// A stale link's destination must never be followed. Removing the link is the
// whole job; removing what it aimed at would be catastrophic.
func TestPruneRemovesTheLinkAndNotItsDestination(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md")

	// An owned link aimed at a real directory inside the repo that is not an item.
	real := filepath.Join(repo, "docs")
	if err := os.MkdirAll(real, 0o755); err != nil {
		t.Fatal(err)
	}
	precious := filepath.Join(real, "SPEC.md")
	if err := os.WriteFile(precious, []byte("the spec"), 0o644); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(tg.Dir(target.Skills), "docs")
	symlink(t, real, dest)

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	Apply(repo, PrunePlan(entries))

	if _, err := os.Lstat(dest); !os.IsNotExist(err) {
		t.Fatal("prune left the stale link behind")
	}
	if _, err := os.Stat(precious); err != nil {
		t.Fatalf("prune followed the link and destroyed its destination: %v", err)
	}
}

func TestApplyContinuesAfterAFailure(t *testing.T) {
	repo, tg := sandbox(t, "skills/alpha/SKILL.md", "skills/beta/SKILL.md")

	entries, err := Scan(repo, tg)
	if err != nil {
		t.Fatal(err)
	}
	plan := Plan(entries)

	// Make the first action impossible without disturbing the second.
	blocked := plan[0].Entry.DestPath
	if err := os.MkdirAll(blocked, 0o755); err != nil {
		t.Fatal(err)
	}

	results := Apply(repo, plan)
	if len(results) != len(plan) {
		t.Fatalf("apply reported %d of %d actions", len(results), len(plan))
	}
	if done, refused, _ := Summarise(results); done != 1 || refused != 1 {
		t.Fatalf("done=%d refused=%d, want 1/1 — one failure must not abandon the rest", done, refused)
	}
}
