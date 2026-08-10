// Package repo answers questions about the repository this tool lives in.
//
// It shells out to git rather than linking a git library. That is deliberate:
// the pull has to work with the user's ssh agent, credential helper, signing
// keys, proxy and .gitconfig, and the only implementation guaranteed to honour
// all of those is the git that made the repository in the first place. A library
// covers a subset, and the day it fails on a credential real git resolves, the
// bug is unfixable from here.
//
// Everything is behind Git so the update flow can be exercised against a fake,
// with no network and no temp repository. See docs/PLAN.md phase 4.1.
package repo

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Git is the repository state the update flow needs. Nothing more.
type Git interface {
	// Dirty reports whether the working tree has uncommitted changes,
	// untracked files included. A repository with no commits yet is dirty.
	Dirty() (bool, error)

	// HasRemote reports whether any remote is configured.
	HasRemote() (bool, error)

	// Head is the current commit, or "" when the repository has no commits.
	Head() (string, error)

	// Pull fast-forwards from the tracking branch.
	Pull() error

	// ChangedSince lists paths that changed between rev and HEAD, relative to
	// the repository root.
	ChangedSince(rev string) ([]string, error)
}

// Shell is the real implementation, running git in a subprocess.
type Shell struct{ Root string }

func (s Shell) git(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", s.Root}, args...)...)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text == "" {
			return "", err
		}
		return "", fmt.Errorf("%s", text)
	}
	return text, nil
}

// Dirty treats a repository with no commits as dirty.
//
// That is not a corner case to tidy away — it is the honest answer. Every file
// is untracked, so pulling over it could not be reconciled with anything, and
// R3's promise is that uncommitted work is never at risk.
func (s Shell) Dirty() (bool, error) {
	out, err := s.git("status", "--porcelain")
	if err != nil {
		return true, err
	}
	return out != "", nil
}

func (s Shell) HasRemote() (bool, error) {
	out, err := s.git("remote")
	if err != nil {
		return false, err
	}
	return out != "", nil
}

// Head returns "" rather than an error when there are no commits yet. Having no
// history is a state, not a failure.
func (s Shell) Head() (string, error) {
	out, err := s.git("rev-parse", "HEAD")
	if err != nil {
		return "", nil
	}
	return out, nil
}

// Pull refuses anything that is not a fast-forward.
//
// A merge commit created by a background command is a merge commit nobody chose.
// If the histories have diverged, that is the user's to resolve, in their own
// shell, seeing what they are doing.
func (s Shell) Pull() error {
	_, err := s.git("pull", "--ff-only")
	return err
}

func (s Shell) ChangedSince(rev string) ([]string, error) {
	if rev == "" {
		return nil, nil
	}
	out, err := s.git("diff", "--name-only", rev, "HEAD")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, nil
	}
	return strings.Split(out, "\n"), nil
}

// NeedsRebuild reports whether any changed path invalidates the compiled binary.
//
// Pure, so the decision is testable without a repository. A pull that carries
// only markdown must not pay for a compile — that is most pulls, in a project
// whose payload is prose.
func NeedsRebuild(paths []string) bool {
	for _, p := range paths {
		switch {
		case filepath.Ext(p) == ".go":
			return true
		case p == "go.mod", p == "go.sum":
			return true
		}
	}
	return false
}
