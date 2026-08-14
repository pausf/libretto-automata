package link

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pausf/libretto-automata/internal/target"
)

// generated writes a file carrying a marker naming source, and returns its path.
func generated(t *testing.T, dir, name, source string) string {
	t.Helper()
	return generatedRaw(t, dir, name,
		"---\nname: "+name+"\nmode: subagent\n"+target.MarkerKey+": "+source+"\n---\n\nbody\n")
}

func generatedRaw(t *testing.T, dir, name, content string) string {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// Every refusal here asks OwnedEither, never Owned. Owned is symlink-only, so asking it
// about a regular file returns false whatever the file contains — the refusals would all
// pass for the wrong reason, which is the worst kind of green.
//
// The refusals come first, and on purpose. This predicate is conservative by
// construction: the cost of a false negative is a repair that does not happen, and
// the cost of a false positive is destroyed work. So what it must refuse is pinned
// before what it must accept.
func TestGeneratedOwnershipRefusesWhatItCannotProve(t *testing.T) {
	repo := repoWith(t, "agents/work-reviewer.md")
	dir := t.TempDir()
	source := filepath.Join(repo, "agents", "work-reviewer.md")

	t.Run("no marker at all", func(t *testing.T) {
		p := generatedRaw(t, dir, "a.md", "---\nname: a\nmode: subagent\n---\n\nbody\n")
		if OwnedEither(repo, p) {
			t.Error("a file with no marker was called ours")
		}
	})

	t.Run("marker names a path outside the repo", func(t *testing.T) {
		p := generated(t, dir, "b.md", filepath.Join(t.TempDir(), "elsewhere.md"))
		if OwnedEither(repo, p) {
			t.Error("a marker naming a path outside the repo was accepted")
		}
	})

	t.Run("marker names a prefix sibling of the repo", func(t *testing.T) {
		// /repo-backup is not inside /repo. A string prefix test would say it is,
		// which is how a tool deletes the wrong directory.
		p := generated(t, dir, "c.md", repo+"-backup/agents/x.md")
		if OwnedEither(repo, p) {
			t.Error("a sibling directory sharing the repo's prefix was accepted")
		}
	})

	t.Run("marker is in the prose body, not the frontmatter", func(t *testing.T) {
		p := generatedRaw(t, dir, "d.md",
			"---\nname: d\n---\n\nProse that mentions "+target.MarkerKey+": "+source+"\n")
		if OwnedEither(repo, p) {
			t.Error("a marker in the body granted ownership — anybody could type that into a file")
		}
	})

	t.Run("marker after the frontmatter closes", func(t *testing.T) {
		p := generatedRaw(t, dir, "e.md",
			"---\nname: e\n---\n"+target.MarkerKey+": "+source+"\n")
		if OwnedEither(repo, p) {
			t.Error("a marker past the closing --- was read as frontmatter")
		}
	})

	t.Run("frontmatter never closes", func(t *testing.T) {
		p := generatedRaw(t, dir, "f.md", "---\nname: f\nno close here\n")
		if OwnedEither(repo, p) {
			t.Error("an unclosed frontmatter block was trusted")
		}
	})

	t.Run("no frontmatter at all", func(t *testing.T) {
		p := generatedRaw(t, dir, "g.md", target.MarkerKey+": "+source+"\n")
		if OwnedEither(repo, p) {
			t.Error("a marker on line 1 with no --- was trusted")
		}
	})

	t.Run("relative marker", func(t *testing.T) {
		p := generated(t, dir, "h.md", "agents/work-reviewer.md")
		if OwnedEither(repo, p) {
			t.Error("a relative marker was accepted — there is no unambiguous base to resolve it against")
		}
	})

	t.Run("empty marker", func(t *testing.T) {
		p := generated(t, dir, "i.md", "")
		if OwnedEither(repo, p) {
			t.Error("an empty marker was accepted")
		}
	})

	t.Run("a directory", func(t *testing.T) {
		d := filepath.Join(dir, "adir")
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if OwnedEither(repo, d) {
			t.Error("a directory was called a generated file")
		}
	})

	t.Run("a path that does not exist", func(t *testing.T) {
		if OwnedEither(repo, filepath.Join(dir, "nope.md")) {
			t.Error("a missing path was called ours")
		}
	})
}

func TestOwnedGeneratedFile(t *testing.T) {
	repo := repoWith(t, "agents/work-reviewer.md")
	source := filepath.Join(repo, "agents", "work-reviewer.md")
	p := generated(t, t.TempDir(), "work-reviewer.md", source)

	if !OwnedEither(repo, p) {
		t.Error("a file carrying a marker naming a source inside the repo is ours")
	}
	if !OwnedGenerated(repo, p) {
		t.Error("OwnedGenerated disagreed with OwnedEither")
	}
	if Owned(repo, p) {
		t.Error("plain Owned accepted a regular file; the marker arm is scoped to generated kinds")
	}
}

// Orphaned is Stale, never foreign. Treating a marker with a missing source as
// somebody else's would strand the file forever — the same reason a broken symlink
// is still ours.
func TestGeneratedOwnershipSurvivesAMissingSource(t *testing.T) {
	repo := repoWith(t, "agents/")
	p := generated(t, t.TempDir(), "gone.md", filepath.Join(repo, "agents", "deleted.md"))

	if !OwnedEither(repo, p) {
		t.Error("a marker naming a deleted source was called foreign — prune could never remove it")
	}
}

func TestGeneratedSourceReporting(t *testing.T) {
	repo := repoWith(t, "agents/a.md")
	source := filepath.Join(repo, "agents", "a.md")

	got, ok := GeneratedSource(generated(t, t.TempDir(), "a.md", source))
	if !ok {
		t.Fatal("GeneratedSource reported no marker on a file that has one")
	}
	if got != source {
		t.Errorf("GeneratedSource = %q, want %q", got, source)
	}

	if _, ok := GeneratedSource(generatedRaw(t, t.TempDir(), "b.md", "---\nname: b\n---\n")); ok {
		t.Error("GeneratedSource claimed a marker on a file with none")
	}
}

func TestUnreadableFileIsForeign(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root reads anything, so the refusal cannot be observed")
	}
	repo := repoWith(t, "agents/a.md")
	p := generated(t, t.TempDir(), "a.md", filepath.Join(repo, "agents", "a.md"))
	if err := os.Chmod(p, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(p, 0o644) })

	if OwnedEither(repo, p) {
		t.Error("an unreadable file was called ours — what cannot be proven is foreign")
	}
}

// The arm that existed before this change must answer identically. Reaching the
// generated arm requires the path not to be a symlink, so a symlink's answer cannot
// move — pinned rather than assumed, because every write is gated on it.
func TestOwnedSymlinkArmIsUnchanged(t *testing.T) {
	repo := repoWith(t, "skills/alpha/")
	item := filepath.Join(repo, "skills", "alpha")
	dir := t.TempDir()

	inside := filepath.Join(dir, "inside")
	symlink(t, item, inside)
	if !Owned(repo, inside) {
		t.Error("a link into the repo stopped being ours")
	}

	outside := filepath.Join(dir, "outside")
	symlink(t, t.TempDir(), outside)
	if Owned(repo, outside) {
		t.Error("a link leaving the repo became ours")
	}

	// A symlink pointing at a generated file that *would* be owned is still judged
	// by its destination, not by the content it resolves to.
	marked := generated(t, dir, "marked.md", item)
	viaLink := filepath.Join(dir, "via-link.md")
	symlink(t, marked, viaLink)
	if Owned(repo, viaLink) {
		t.Error("a symlink leaving the repo was judged by the marker in the file it points at")
	}
}

// The marker arm must not reach a destination this tool never generates into.
//
// A review found `prune --claude --yes` offering to delete a hand-written file in a
// Claude destination because it happened to carry a marker line. Owned is symlink-only;
// OwnedEither is the widened question, and only a generated kind asks it.
func TestOwnedIgnoresTheMarkerForNonGeneratedKinds(t *testing.T) {
	repo := repoWith(t, "agents/work-reviewer.md")
	source := filepath.Join(repo, "agents", "work-reviewer.md")
	p := generated(t, t.TempDir(), "notes.md", source)

	if Owned(repo, p) {
		t.Error("Owned accepted a marked regular file — that is how prune deletes somebody's notes")
	}
	if !OwnedEither(repo, p) {
		t.Error("OwnedEither refused a marked file; the widened arm is what a generated kind needs")
	}
}

// A marker must name an item, and the repository is not one. Both of these were
// accepted and then deleted by `prune --yes` before the fix.
func TestGeneratedOwnershipRefusesTheRootAndDirectories(t *testing.T) {
	repo := repoWith(t, "agents/work-reviewer.md", "skills/alpha/")
	dir := t.TempDir()

	if OwnedEither(repo, generated(t, dir, "a.md", repo)) {
		t.Error("a marker naming the repository root was accepted")
	}
	if OwnedEither(repo, generated(t, dir, "b.md", filepath.Join(repo, "agents"))) {
		t.Error("a marker naming a directory inside the repo was accepted")
	}
	if OwnedEither(repo, generated(t, dir, "c.md", filepath.Join(repo, "skills", "alpha"))) {
		t.Error("a marker naming a skill directory was accepted")
	}
	// The orphan case still has to work: a missing source cannot be stat'd, so it
	// stays ours and prune can remove it.
	if !OwnedEither(repo, generated(t, dir, "d.md", filepath.Join(repo, "agents", "gone.md"))) {
		t.Error("a marker naming a deleted file stopped being ours; prune could never remove it")
	}
}

// The transform quotes the marker, so the reader has to accept both forms. Refusing a
// bare value would turn a file written before quoting landed into a foreign one that
// nothing can clean up.
func TestGeneratedSourceAcceptsQuotedAndBareMarkers(t *testing.T) {
	repo := repoWith(t, "agents/a.md")
	source := filepath.Join(repo, "agents", "a.md")
	dir := t.TempDir()

	quoted := generatedRaw(t, dir, "q.md",
		"---\nname: q\n"+target.MarkerKey+`: "`+source+"\"\n---\n\nbody\n")
	got, ok := GeneratedSource(quoted)
	if !ok || got != source {
		t.Errorf("quoted marker: got %q ok=%v, want %q", got, ok, source)
	}

	bare := generated(t, dir, "b.md", source)
	got, ok = GeneratedSource(bare)
	if !ok || got != source {
		t.Errorf("bare marker: got %q ok=%v, want %q", got, ok, source)
	}
}
