package link

import (
	"os"
	"path/filepath"
	"testing"
)

func symlink(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(newname), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(oldname, newname); err != nil {
		t.Fatal(err)
	}
}

// This is the test that protects everything else. Ownership decides whether
// `install` may overwrite a path and whether `prune` may delete it, so both
// directions get pinned explicitly.
func TestOwned(t *testing.T) {
	t.Run("an absolute link into the repo is ours", func(t *testing.T) {
		repo := repoWith(t, "skills/po-assistant/")
		dest := filepath.Join(t.TempDir(), "po-assistant")
		symlink(t, filepath.Join(repo, "skills", "po-assistant"), dest)

		if !Owned(repo, dest) {
			t.Error("Owned() = false for a link into the repo")
		}
	})

	t.Run("a link outside the repo is foreign", func(t *testing.T) {
		repo := repoWith(t, "skills/po-assistant/")
		elsewhere := repoWith(t, "skills/other/")
		dest := filepath.Join(t.TempDir(), "other")
		symlink(t, filepath.Join(elsewhere, "skills", "other"), dest)

		if Owned(repo, dest) {
			t.Error("Owned() = true for a link into another repo")
		}
	})

	t.Run("a relative link is resolved against the link's own directory", func(t *testing.T) {
		// The classic mistake is resolving against the process working
		// directory, which would misjudge this.
		base := t.TempDir()
		repo := filepath.Join(base, "repo")
		if err := os.MkdirAll(filepath.Join(repo, "skills", "po-assistant"), 0o755); err != nil {
			t.Fatal(err)
		}
		home := filepath.Join(base, "home", "skills")
		symlink(t, filepath.Join("..", "..", "repo", "skills", "po-assistant"),
			filepath.Join(home, "po-assistant"))

		if !Owned(repo, filepath.Join(home, "po-assistant")) {
			t.Error("Owned() = false for a relative link into the repo")
		}
	})

	t.Run("a relative link pointing away is foreign", func(t *testing.T) {
		base := t.TempDir()
		repo := filepath.Join(base, "repo")
		outside := filepath.Join(base, "outside", "thing")
		for _, d := range []string{filepath.Join(repo, "skills"), outside} {
			if err := os.MkdirAll(d, 0o755); err != nil {
				t.Fatal(err)
			}
		}
		home := filepath.Join(base, "home", "skills")
		symlink(t, filepath.Join("..", "..", "outside", "thing"),
			filepath.Join(home, "thing"))

		if Owned(repo, filepath.Join(home, "thing")) {
			t.Error("Owned() = true for a relative link leaving the repo")
		}
	})

	t.Run("a broken link is still ours", func(t *testing.T) {
		// This is what makes prune possible: the source is gone, but the link
		// is unmistakably one we made.
		repo := repoWith(t, "skills/")
		dest := filepath.Join(t.TempDir(), "deleted-skill")
		symlink(t, filepath.Join(repo, "skills", "deleted-skill"), dest)

		if _, err := os.Stat(dest); err == nil {
			t.Fatal("the fixture link is not broken; the test proves nothing")
		}
		if !Owned(repo, dest) {
			t.Error("Owned() = false for a broken link into the repo")
		}
	})

	t.Run("a real directory is not a link and so not ours", func(t *testing.T) {
		repo := repoWith(t, "skills/po-assistant/")
		dest := filepath.Join(t.TempDir(), "po-assistant")
		if err := os.MkdirAll(dest, 0o755); err != nil {
			t.Fatal(err)
		}

		if Owned(repo, dest) {
			t.Error("Owned() = true for a real directory")
		}
	})

	t.Run("a real file is not ours", func(t *testing.T) {
		repo := repoWith(t, "agents/")
		dest := filepath.Join(t.TempDir(), "agent.md")
		if err := os.WriteFile(dest, []byte("mine"), 0o644); err != nil {
			t.Fatal(err)
		}

		if Owned(repo, dest) {
			t.Error("Owned() = true for a real file")
		}
	})

	t.Run("a missing path is not ours", func(t *testing.T) {
		repo := repoWith(t, "skills/")

		if Owned(repo, filepath.Join(t.TempDir(), "nothing-here")) {
			t.Error("Owned() = true for a path that does not exist")
		}
	})
}

// A sibling directory sharing a name prefix must never count as inside the repo.
// A strings.HasPrefix check would get this wrong, and getting it wrong means
// prune deleting from the wrong tree.
func TestOwnedRejectsPrefixSiblings(t *testing.T) {
	base := t.TempDir()
	repo := filepath.Join(base, "repo")
	sibling := filepath.Join(base, "repo-backup")
	for _, d := range []string{filepath.Join(repo, "skills"), filepath.Join(sibling, "skills", "thing")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	dest := filepath.Join(t.TempDir(), "thing")
	symlink(t, filepath.Join(sibling, "skills", "thing"), dest)

	if Owned(repo, dest) {
		t.Errorf("Owned(%q, → %q) = true; a name prefix is not containment", repo, sibling)
	}
}

// The repo reached through a symlinked path must still recognise its own links.
// This is not hypothetical: on macOS /tmp is a symlink to /private/tmp and
// /var/folders (where t.TempDir lives) sits under a symlinked /var.
func TestOwnedWhenTheRepoPathIsItselfSymlinked(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real-repo")
	if err := os.MkdirAll(filepath.Join(real, "skills", "po-assistant"), 0o755); err != nil {
		t.Fatal(err)
	}

	alias := filepath.Join(base, "aliased-repo")
	symlink(t, real, alias)

	dest := filepath.Join(t.TempDir(), "po-assistant")
	symlink(t, filepath.Join(real, "skills", "po-assistant"), dest)

	t.Run("repo named by its real path", func(t *testing.T) {
		if !Owned(real, dest) {
			t.Error("Owned() = false")
		}
	})

	t.Run("repo named through the alias", func(t *testing.T) {
		if !Owned(alias, dest) {
			t.Error("Owned() = false; the two spellings of one directory did not compare equal")
		}
	})
}

// And the mirror case: the link points at the aliased spelling while the repo is
// named by its real path.
func TestOwnedWhenTheLinkUsesAnAliasedPath(t *testing.T) {
	base := t.TempDir()
	real := filepath.Join(base, "real-repo")
	if err := os.MkdirAll(filepath.Join(real, "skills", "po-assistant"), 0o755); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(base, "aliased-repo")
	symlink(t, real, alias)

	dest := filepath.Join(t.TempDir(), "po-assistant")
	symlink(t, filepath.Join(alias, "skills", "po-assistant"), dest)

	if !Owned(real, dest) {
		t.Error("Owned() = false for a link written through an aliased repo path")
	}
}

func TestLinkTarget(t *testing.T) {
	repo := repoWith(t, "skills/po-assistant/")
	want := normalise(filepath.Join(repo, "skills", "po-assistant"))

	dest := filepath.Join(t.TempDir(), "po-assistant")
	symlink(t, filepath.Join(repo, "skills", "po-assistant"), dest)

	got, ok := LinkTarget(dest)
	if !ok {
		t.Fatal("LinkTarget() reported the path is not a symlink")
	}
	if got != want {
		t.Errorf("LinkTarget() = %q, want %q", got, want)
	}

	t.Run("a real directory is not a link", func(t *testing.T) {
		if _, ok := LinkTarget(repo); ok {
			t.Error("LinkTarget() reported a directory as a symlink")
		}
	})
}

func TestWithin(t *testing.T) {
	tests := []struct {
		name string
		root string
		path string
		want bool
	}{
		{"the root itself", "/repo", "/repo", true},
		{"a child", "/repo", "/repo/skills", true},
		{"a deep descendant", "/repo", "/repo/skills/a/b/c", true},
		{"the parent", "/repo/skills", "/repo", false},
		{"a sibling", "/repo", "/other", false},
		{"a prefix sibling", "/repo", "/repo-backup", false},
		{"a prefix sibling's child", "/repo", "/repo-backup/skills", false},
		{"an escaping traversal", "/repo", "/repo/../other", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := within(tt.root, filepath.Clean(tt.path)); got != tt.want {
				t.Errorf("within(%q, %q) = %v, want %v", tt.root, tt.path, got, tt.want)
			}
		})
	}
}

func TestNormaliseHandlesMissingComponents(t *testing.T) {
	base := t.TempDir()
	missing := filepath.Join(base, "not", "there", "yet")

	got := normalise(missing)
	want := filepath.Join(normalise(base), "not", "there", "yet")

	if got != want {
		t.Errorf("normalise(%q) = %q, want %q", missing, got, want)
	}
}
