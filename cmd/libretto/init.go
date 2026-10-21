package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// `libretto init` bootstraps the minimal .agents/ layout — and nothing else.
//
// The layout used to appear as a side effect of the first change, and git drops an
// empty changes/ directory on every fresh checkout, which is how a queued idea got
// lost. The keep-files exist so git carries the directories; a project outside git
// has nothing to carry them, which is why a non-repo is a refusal and not a warning.

// The foreign consolidation roots, mirroring SPEC_ROOTS/CHANGE_ROOTS in
// skills/record-work/spec-drift, minus the .agents/ pair — the layout
// being bootstrapped. spec-drift is the single authority for these lists; a
// divergence is a bug in this copy.
var (
	foreignSpecRoots   = []string{"specs", "openspec", "docs/specs", "spec"}
	foreignChangeRoots = []string{"changes", "openspec/changes"}
)

// lessonsHeader is byte-for-byte how this repository's own .agents/lessons.md opens.
const lessonsHeader = "# Lessons\n\n" +
	"Append-only. Entries are spent by `libretto-retro` and gain a `Resolved:` line there;\n" +
	"they are never edited otherwise.\n"

// initCmd creates what is missing of the .agents/ layout at the repository root.
// Named to dodge the `init` keyword; wired as `libretto init`.
func initCmd(w io.Writer, git gitRunner) error {
	// The root, never the cwd — metrics.go records the cwd-relative bug this avoids.
	root, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	root = strings.TrimSpace(root)

	// The probe runs before any write, so the refusal wins over a partial .agents/
	// by ordering alone. IsDir(), not a bare stat: a *file* named specs is not a
	// consolidation root.
	var found []string
	for _, p := range append(append([]string{}, foreignSpecRoots...), foreignChangeRoots...) {
		if info, err := os.Stat(filepath.Join(root, filepath.FromSlash(p))); err == nil && info.IsDir() {
			found = append(found, p+"/")
		}
	}
	if len(found) > 0 {
		return fmt.Errorf("this project already consolidates elsewhere: %s — init bootstraps only the .agents/ layout",
			strings.Join(found, ", "))
	}

	fmt.Fprintln(w, root)

	// Present in any form means untouched — init cannot tell its own .gitkeep from
	// anyone else's and does not need to. A symlink is present too, dangling or not:
	// Lstat sees it where Stat would follow it, and O_EXCL refuses to write through
	// one created in between, so nothing ever lands outside the root.
	artifacts := []struct {
		rel     string
		content string
	}{
		{filepath.Join(".agents", "specs", ".gitkeep"), ""},
		{filepath.Join(".agents", "changes", ".gitkeep"), ""},
		{filepath.Join(".agents", "lessons.md"), lessonsHeader},
	}
	for _, a := range artifacts {
		path := filepath.Join(root, a.rel)
		if _, err := os.Lstat(path); err == nil {
			fmt.Fprintf(w, "  exists   %s\n", filepath.ToSlash(a.rel))
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return err
		}
		_, werr := f.WriteString(a.content)
		if cerr := f.Close(); werr == nil {
			werr = cerr
		}
		if werr != nil {
			return werr
		}
		fmt.Fprintf(w, "  created  %s\n", filepath.ToSlash(a.rel))
	}
	return nil
}
