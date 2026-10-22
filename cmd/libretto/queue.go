package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// queuedIdea is one captured idea: what `libretto queue` prints and the wiki's flow board
// lists.
type queuedIdea struct {
	name, date, reading string
}

// queuedIdeas is find-work's queue scan, minus the branch filter: a proposal is a
// candidate when a line starts with `Queued:` at column 0 and no box sits under its
// `## Checklist`, counted by readPlan so the loop and the queue agree on what a box is.
// The skill owns the definition; this implements it. Oldest date first, ties by name.
func queuedIdeas(projectDir string) []queuedIdea {
	changes := filepath.Join(projectDir, ".agents", "changes")
	entries, err := os.ReadDir(changes)
	if err != nil {
		return nil
	}
	var ideas []queuedIdea
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(changes, e.Name(), "proposal.md")
		b, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		date, ok := queuedDate(string(b))
		if !ok {
			continue
		}
		if c, err := readPlan(path); err == nil && c.done+c.open > 0 {
			continue
		}
		ideas = append(ideas, queuedIdea{name: e.Name(), date: date, reading: proposalReading(string(b))})
	}
	sort.SliceStable(ideas, func(i, j int) bool {
		if ideas[i].date != ideas[j].date {
			return ideas[i].date < ideas[j].date
		}
		return ideas[i].name < ideas[j].name
	})
	return ideas
}

func queuedDate(proposal string) (string, bool) {
	for _, line := range strings.Split(proposal, "\n") {
		if strings.HasPrefix(line, "Queued:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Queued:")), true
		}
	}
	return "", false
}

// proposalReading is the first sentence of `## Reading`, or "" when there is none.
func proposalReading(proposal string) string {
	var para []string
	in := false
	for _, line := range strings.Split(proposal, "\n") {
		if strings.HasPrefix(line, "## ") {
			if in {
				break
			}
			in = strings.TrimSpace(line) == "## Reading"
			continue
		}
		if l := strings.TrimSpace(line); in && l != "" {
			para = append(para, l)
		}
	}
	return firstSentence(para)
}

func unmergedBranches(git gitRunner) ([]string, error) {
	out, err := git("branch", "--format=%(refname:short)", "--no-merged", "HEAD")
	if err != nil {
		return nil, err
	}
	var branches []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			branches = append(branches, l)
		}
	}
	return branches, nil
}

// dispatchedBy names the first branch that took the change, or "". The match is
// attributeBranch's — the same prefix list `metrics` strips, and find-work's.
func dispatchedBy(name string, branches []string) string {
	for _, b := range branches {
		if attributeBranch(b, []string{name}) != "" {
			return b
		}
	}
	return ""
}

type dispatchedIdea struct {
	queuedIdea
	branch string
}

func splitDispatched(ideas []queuedIdea, branches []string) (kept []queuedIdea, dropped []dispatchedIdea) {
	for _, i := range ideas {
		if b := dispatchedBy(i.name, branches); b != "" {
			dropped = append(dropped, dispatchedIdea{i, b})
		} else {
			kept = append(kept, i)
		}
	}
	return kept, dropped
}

// queue is `libretto queue`: read-only, and every empty state exits 0. Changes are read at
// the repository root — `metrics` once trusted the cwd and reported nothing from a
// subdirectory — and git failing to list branches keeps every idea, find-work's fail-safe
// direction.
func queue(w, errw io.Writer, args []string, cwd string, git gitRunner) error {
	if len(args) > 0 {
		return fmt.Errorf("queue takes no arguments, got %q", args[0])
	}
	root := cwd
	if out, err := git("rev-parse", "--show-toplevel"); err == nil && strings.TrimSpace(out) != "" {
		root = strings.TrimSpace(out)
	}
	if info, err := os.Stat(filepath.Join(root, ".agents", "changes")); err != nil || !info.IsDir() {
		fmt.Fprintln(w, "no .agents/changes/ here — nothing is queued")
		return nil
	}
	ideas := queuedIdeas(root)
	branches, err := unmergedBranches(git)
	if err != nil {
		fmt.Fprintf(errw, "cannot list branches (%v) — the queue is shown unfiltered\n", err)
	}
	kept, dropped := splitDispatched(ideas, branches)
	if len(kept) == 0 {
		fmt.Fprintln(w, "nothing is queued")
	}
	for _, i := range kept {
		line := i.date + "  " + i.name
		if i.reading != "" {
			line += "  " + i.reading
		}
		fmt.Fprintln(w, line)
	}
	if len(dropped) > 0 {
		names := make([]string, len(dropped))
		for n, d := range dropped {
			names[n] = fmt.Sprintf("%s (%s)", d.name, d.branch)
		}
		fmt.Fprintln(w, "already taken by a branch: "+strings.Join(names, ", "))
	}
	return nil
}
