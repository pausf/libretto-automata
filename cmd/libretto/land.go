package main

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// `libretto land` — a read-only verifier of the landing commit.
//
// A landing is one commit doing four things — final code, delta applied onto each
// Targets: capability spec, durable decisions retired, change folder deleted — and it
// failed twice by half-finishing silently. land reads the staged index, before the
// mistake is history, and names which of the two mechanically checkable parts the
// commit about to exist is missing: the folder fully deleted (part 4) and each delta's
// capability spec moved in the same diff (part 2). It verifies; it never performs.

// landChangeRoots mirrors spec-drift's CHANGE_ROOTS verbatim
// (skills/record-work/spec-drift) — the script is the authority. Three strings in Go
// against three in bash is cheaper than the binary parsing the script or the script
// reading Go; unify them, or add a test that greps the script, the day they diverge.
var landChangeRoots = []string{".agents/changes", "changes", "openspec/changes"}

// landFenceRe mirrors spec-drift's `defenced`: a line that is a fence toggles the
// in-fence state and is itself dropped, so a Targets: shown as an example never reads
// as a declaration — the two tools must not disagree about what a delta targets.
var (
	landFenceRe   = regexp.MustCompile("^\\s*```")
	landTargetsRe = regexp.MustCompile(`^\s*Targets:\s*(.+)`)
)

// parseCachedDiff reads `git diff --cached --name-status -z` into the two sets
// everything downstream asks about: paths gone from the index after the commit, and
// paths added or modified by it. A rename counts its source as removed — the contract
// is that nothing under the folder survives, and a rename-out leaves it empty — and
// its destination as touched; a copy touches only its destination.
func parseCachedDiff(out string) (removed, touched map[string]bool, err error) {
	removed, touched = map[string]bool{}, map[string]bool{}
	tok := strings.Split(out, "\x00")
	// The output ends in a NUL, so a missing token arrives as an empty string, not a
	// short slice — both are the same truncation and both must refuse.
	at := func(i int) (string, error) {
		if i >= len(tok) || tok[i] == "" {
			return "", fmt.Errorf("truncated staged-diff entry near %q", tok[i-1])
		}
		return tok[i], nil
	}
	for i := 0; i < len(tok); {
		status := tok[i]
		if status == "" {
			i++
			continue
		}
		switch status[0] {
		case 'D':
			p, err := at(i + 1)
			if err != nil {
				return nil, nil, err
			}
			removed[p] = true
			i += 2
		case 'A', 'M', 'T':
			p, err := at(i + 1)
			if err != nil {
				return nil, nil, err
			}
			touched[p] = true
			i += 2
		case 'R', 'C':
			src, err := at(i + 1)
			if err != nil {
				return nil, nil, err
			}
			dst, err := at(i + 2)
			if err != nil {
				return nil, nil, err
			}
			if status[0] == 'R' {
				removed[src] = true
			}
			touched[dst] = true
			i += 3
		default:
			// 'U' among them: an in-progress merge is out of scope, per the plan.
			return nil, nil, fmt.Errorf("unrecognised staged-diff status %q", status)
		}
	}
	return removed, touched, nil
}

func land(w, ew io.Writer, args []string, git gitRunner) error {
	var name string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			return fmt.Errorf("unknown flag %q — `%s land [<change>]`", a, invokedAs())
		}
		if name != "" {
			return fmt.Errorf("one change at a time, got %q and %q", name, a)
		}
		name = a
	}

	// Everything is asked of the repository root, never of the cwd — metrics records
	// the run-from-a-subdirectory bug this avoids.
	rootOut, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	root := strings.TrimSpace(rootOut)

	diffOut, err := git("diff", "--cached", "--name-status", "-z")
	if err != nil {
		return fmt.Errorf("git could not read the staged index: %w", err)
	}
	removed, touched, err := parseCachedDiff(diffOut)
	if err != nil {
		return err
	}

	folders, err := landingFolders(git, root, name, removed)
	if err != nil {
		return err
	}
	if len(folders) == 0 {
		// A verifier that exits zero having verified nothing is the silent
		// half-landing wearing a green light.
		if name != "" {
			return fmt.Errorf("nothing is landing: %q has no staged deletions — stage the landing, then run `%s land %s`",
				name, invokedAs(), name)
		}
		return fmt.Errorf("nothing is landing: no staged deletion under %s — `%s land [<change>]` verifies a staged landing",
			strings.Join(landChangeRoots, ", "), invokedAs())
	}

	specsRel := ""
	if dir, ok := findSpecsDir(root); ok {
		if rel, err := filepath.Rel(root, dir); err == nil {
			specsRel = filepath.ToSlash(rel)
		}
	}

	missing := 0
	for _, folder := range folders {
		fmt.Fprintf(w, "landing %s (%s)\n", path.Base(folder), folder)

		tracked, err := lsTreeUnder(git, root, folder)
		if err != nil {
			return err
		}

		// Part 4 — nothing under the folder survives the commit: every tracked file
		// a staged deletion, nothing untracked left on disk.
		var survivors []string
		for _, f := range tracked {
			if !removed[f] {
				survivors = append(survivors, f)
			}
		}
		leftovers, err := untrackedUnder(git, root, folder)
		if err != nil {
			return err
		}
		if len(survivors)+len(leftovers) == 0 {
			fmt.Fprintf(w, "  part 4: the folder is fully deleted\n")
		} else {
			missing++
			for _, f := range survivors {
				fmt.Fprintf(w, "  part 4: %s survives the commit — not a staged deletion\n", f)
			}
			for _, f := range leftovers {
				fmt.Fprintf(w, "  part 4: %s is untracked and still on disk — a commit cannot delete a folder it leaves a file in\n", f)
			}
		}

		// Part 2 — each capability a delta targets, read from HEAD because the
		// landing has already deleted the delta from disk and index.
		caps, err := headTargets(git, tracked)
		if err != nil {
			return err
		}
		if len(caps) == 0 {
			fmt.Fprintf(w, "  part 2: no delta — an abandoned proposal, not a landing; passes vacuously\n")
		} else {
			broken := false
			for _, c := range caps {
				spec := specsRel + "/" + c + "/spec.md"
				if specsRel != "" && touched[spec] {
					fmt.Fprintf(w, "  part 2: %s spec moved (%s)\n", c, spec)
					continue
				}
				broken = true
				fmt.Fprintf(w, "  part 2: %s spec did not move — %s is not in the staged diff\n", c, spec)
			}
			if broken {
				missing++
			}
		}
	}

	// Said on every run, pass or fail: a green land must never read as the whole
	// landing contract passing.
	fmt.Fprintf(w, "part 3 (decisions retired) is owned by spec-drift --anchors — not checked here\n")

	warnStaleWiki(ew, root, specsRel, touched)

	if missing > 0 {
		return fmt.Errorf("%d part(s) of the landing contract are missing — fix them, then run `%s land` again",
			missing, invokedAs())
	}
	return nil
}

// landingFolders resolves which change folders this commit lands, repo-root-relative.
// With a name: located in HEAD, not on disk — a fully staged landing may already have
// removed the folder from the working tree. With none: inferred from the staged
// deletions, one folder or several, each on its own contract.
func landingFolders(git gitRunner, root, name string, removed map[string]bool) ([]string, error) {
	if name == "" {
		seen := map[string]bool{}
		for p := range removed {
			for _, cr := range landChangeRoots {
				rest, ok := strings.CutPrefix(p, cr+"/")
				if !ok {
					continue
				}
				if n, _, ok := strings.Cut(rest, "/"); ok && n != "" {
					seen[path.Join(cr, n)] = true
				}
			}
		}
		folders := make([]string, 0, len(seen))
		for f := range seen {
			folders = append(folders, f)
		}
		sort.Strings(folders)
		return folders, nil
	}

	var lastErr error
	failures := 0
	for _, cr := range landChangeRoots {
		folder := path.Join(cr, name)
		files, err := lsTreeUnder(git, root, folder)
		if err != nil {
			failures++
			lastErr = err
			continue
		}
		if len(files) == 0 {
			continue
		}
		for p := range removed {
			if strings.HasPrefix(p, folder+"/") {
				return []string{folder}, nil
			}
		}
		return nil, nil // found in HEAD, nothing staged — the caller refuses
	}
	// Every root failing means HEAD itself would not answer — an unborn branch, not
	// an absent folder — and that is a git error, never an empty green report.
	if failures == len(landChangeRoots) {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no change folder named %q in HEAD under %s",
		name, strings.Join(landChangeRoots, ", "))
}

// lsTreeUnder lists every tracked file HEAD holds under folder. The pathspec is
// absolute so it means the same thing from any working directory; the output is
// repo-root-relative regardless.
func lsTreeUnder(git gitRunner, root, folder string) ([]string, error) {
	out, err := git("ls-tree", "-r", "--name-only", "-z", "HEAD", "--", filepath.Join(root, folder)+"/")
	if err != nil {
		return nil, fmt.Errorf("git could not read HEAD's tree under %s: %w", folder, err)
	}
	return nulSplit(out), nil
}

// untrackedUnder lists untracked files still on disk under folder. --exclude-standard
// keeps ignored files out: an ignored scratch file survives into nobody's commit.
func untrackedUnder(git gitRunner, root, folder string) ([]string, error) {
	out, err := git("ls-files", "--others", "--exclude-standard", "--full-name", "-z", "--", filepath.Join(root, folder)+"/")
	if err != nil {
		return nil, fmt.Errorf("git could not list untracked files under %s: %w", folder, err)
	}
	return nulSplit(out), nil
}

// headTargets reads every capability the folder's delta files target, from HEAD —
// never from disk or the index, where the landing has already deleted them — with
// fenced blocks stripped as defenced strips them. Sorted, for a deterministic report.
func headTargets(git gitRunner, tracked []string) ([]string, error) {
	seen := map[string]bool{}
	for _, f := range tracked {
		if !strings.HasSuffix(f, ".md") {
			continue
		}
		out, err := git("show", "HEAD:"+f)
		if err != nil {
			return nil, fmt.Errorf("git could not read HEAD:%s: %w", f, err)
		}
		fence := false
		for _, l := range strings.Split(out, "\n") {
			if landFenceRe.MatchString(l) {
				fence = !fence
				continue
			}
			if fence {
				continue
			}
			if m := landTargetsRe.FindStringSubmatch(l); m != nil {
				for _, c := range strings.Fields(m[1]) {
					seen[c] = true
				}
			}
		}
	}
	caps := make([]string, 0, len(seen))
	for c := range seen {
		caps = append(caps, c)
	}
	sort.Strings(caps)
	return caps, nil
}

// warnStaleWiki is record-work's own clause read mechanically: the refreshed view
// rides the same commit as the delta that changed it. Stale means a marked view
// exists, a capability spec is in the staged diff, and the view is not — one line to
// stderr per view, and the exit code untouched on every path. A view without its
// marker is somebody's and is ignored silently, wiki's own precedent.
func warnStaleWiki(ew io.Writer, root, specsRel string, touched map[string]bool) {
	if specsRel == "" {
		return
	}
	specStaged := false
	for p := range touched {
		rest, ok := strings.CutPrefix(p, specsRel+"/")
		if !ok {
			continue
		}
		if parts := strings.Split(rest, "/"); len(parts) == 2 && parts[1] == "spec.md" {
			specStaged = true
			break
		}
	}
	if !specStaged {
		return
	}
	for _, v := range []struct{ name, marker string }{
		{"README.md", wikiMarker},
		{"wiki.html", wikiHTMLMarker},
	} {
		rel := specsRel + "/" + v.name
		if !ownsFile(filepath.Join(root, filepath.FromSlash(rel)), v.marker) {
			continue
		}
		if touched[rel] {
			continue
		}
		fmt.Fprintf(ew, "warning: %s is stale — a capability spec is staged and the view is not; run `%s wiki` and stage it\n",
			rel, invokedAs())
	}
}

func nulSplit(out string) []string {
	var ps []string
	for _, p := range strings.Split(out, "\x00") {
		if p != "" {
			ps = append(ps, p)
		}
	}
	return ps
}
