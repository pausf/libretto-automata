package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// The loop runner — one fresh session per task, driven from outside the session.
//
// The reason this is a Go subcommand and not a skill or a command is not taste. A
// loop that relaunches a session cannot BE that session: the whole value is a fresh
// context per iteration, and a skill runs inside the context it would be trying to
// discard. So the engine has to sit outside, which is the binary.
//
// It owns no state of its own. The change's checklist is the state — the `## Checklist`
// section of `.agents/changes/<change>/proposal.md`, or `tasks.md`/`plan.md` in a change
// written before that — and the loop's only job is to notice when the boxes stop moving.
// Everything the iteration does — reading the checklist, building, gating, committing —
// belongs to the flow and stays there.

// planLine matches a markdown task checkbox at any indentation: `- [ ]` or `- [x]`.
// Only the box is read. What follows is the task's own text and this never parses it —
// the flow owns the plan's shape, and a runner that understood it would be a second
// opinion about what a task is.
var planLine = regexp.MustCompile(`^\s*[-*]\s+\[( |x|X)\]`)

type planCount struct{ open, done int }

func countBoxes(r io.Reader) planCount {
	var c planCount
	fenced := false
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for s.Scan() {
		if fenceLine.MatchString(s.Text()) {
			fenced = !fenced
			continue
		}
		if !fenced {
			c = c.add(s.Text())
		}
	}
	return c
}

var (
	checklistHeading = regexp.MustCompile(`^##\s+Checklist\s*$`)
	sectionHeading   = regexp.MustCompile(`^##\s`)
	fenceLine        = regexp.MustCompile("^\\s*```")
)

// countChecklist counts the boxes of a change document: only those under a `## Checklist`
// heading, up to the next `##`, outside fences. The boundaries are spec-drift's, so the
// two never disagree about which lines are the checklist. Everything else in the file —
// a quoted ask, a breakdown written as boxes — is prose that happens to look like a task.
func countChecklist(r io.Reader) planCount {
	var c planCount
	var scope checklistScope
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for s.Scan() {
		if line := s.Text(); scope.holds(line) {
			c = c.add(line)
		}
	}
	return c
}

// checklistScope is that section rule fed one line at a time, so `metrics` can walk a
// diff's pre-image and post-image with the same boundaries `loop` reads a file with.
type checklistScope struct{ in, fenced bool }

// holds advances past line and reports whether it lies inside the checklist. Headings
// and fence markers are boundaries, never content.
func (s *checklistScope) holds(line string) bool {
	switch {
	case fenceLine.MatchString(line):
		s.fenced = !s.fenced
	case s.fenced:
	case sectionHeading.MatchString(line):
		s.in = checklistHeading.MatchString(line)
	default:
		return s.in
	}
	return false
}

func (c planCount) add(line string) planCount {
	m := planLine.FindStringSubmatch(line)
	switch {
	case m == nil:
	case m[1] == " ":
		c.open++
	default:
		c.done++
	}
	return c
}

func readPlan(path string) (planCount, error) {
	f, err := os.Open(path)
	if err != nil {
		return planCount{}, err
	}
	defer f.Close()
	if filepath.Base(path) == "proposal.md" {
		return countChecklist(f), nil
	}
	return countBoxes(f), nil
}

// loopPrompt is what each fresh session is handed. It routes and does not instruct:
// every rule the iteration needs already lives in a skill, and restating one here
// creates a second copy that drifts from the installed one.
func loopPrompt(change, checklist string) string {
	where := ""
	if checklist == "proposal.md" {
		where = " under `## Checklist`"
	}
	return fmt.Sprintf(`/libretto-attacca continue the change %q.

Read .agents/changes/%s/%s, take the FIRST unchecked box%s, and do only that one.
Mark it the moment it is genuinely finished, and commit. Then stop — do not take a
second box, and do not push or open a request. Another session follows this one.

If a gate fails twice on this box, stop and say what you observed.`, change, change, checklist, where)
}

type loopOpts struct {
	change string
	max    int
	dryRun bool
}

func parseLoopArgs(args []string) (loopOpts, error) {
	o := loopOpts{max: 10}
	for i := 0; i < len(args); i++ {
		switch a := args[i]; {
		case a == "--dry-run":
			o.dryRun = true
		case a == "--max":
			if i+1 >= len(args) {
				return o, fmt.Errorf("--max needs a number")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n < 1 {
				return o, fmt.Errorf("--max wants a positive number, got %q", args[i])
			}
			o.max = n
		case strings.HasPrefix(a, "-"):
			return o, fmt.Errorf("unknown flag %q", a)
		case o.change != "":
			return o, fmt.Errorf("one change at a time, got %q and %q", o.change, a)
		default:
			o.change = a
		}
	}
	if o.change == "" {
		return o, fmt.Errorf("which change? `%s loop <change>` — see .agents/changes/", invokedAs())
	}
	return o, nil
}

// loopRunner is the one thing the loop does to the outside world, extracted so the
// tests can drive every branch without a `claude` on PATH and without a real session.
type loopRunner func(prompt string) error

func execClaude(prompt string) error {
	c := exec.Command("claude", "-p", prompt)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func loop(projectDir string, args []string) error {
	o, err := parseLoopArgs(args)
	if err != nil {
		return err
	}
	// The checklist is checked before PATH is, and the order is the point: `loop <typo>`
	// on a machine without `claude` should name the file it could not find, and a queued
	// proposal should be named as having nothing to drive — not as missing a dependency
	// that is irrelevant to either.
	if _, _, err := openChecklist(projectDir, o.change); err != nil {
		return err
	}
	if !o.dryRun {
		if _, err := exec.LookPath("claude"); err != nil {
			return fmt.Errorf("claude is not on PATH — the loop relaunches it once per task\n" +
				"       install Claude Code, or use --dry-run to see what would run")
		}
	}
	return runLoop(os.Stdout, projectDir, o, execClaude)
}

// runLoop is the whole engine. Two things stop it besides finishing:
//
//   - the iteration cap, because a loop with no ceiling that never converges burns
//     a session's budget silently and the user finds out from a bill.
//   - stuck detection. If a round ends with the same number of open boxes it started
//     with, the session did not finish a task — and the next fresh session will read
//     the identical plan and make the identical non-progress. One such round is a
//     hiccup; two in a row is a loop that has stopped being one.
//
// It never pushes, never merges, never tags. Those are attacca's answered questions
// for one branch, and nothing here was asked.
// checklistPath answers "which file holds this change's boxes", and it answers it three
// times because the checklist moved twice. The `## Checklist` section of `proposal.md` is
// where the flow writes it now; `tasks.md` is what the 5→6 seam used to cut; `plan.md` is
// what every change created before that carries. A loop that could not drive one of those
// would have broken work already in flight to tidy a name.
//
// `proposal.md` wins only when its section holds a box. Every change in both layouts has a
// proposal.md, so its existence alone would hide every old change's boxes. When nothing
// holds the checklist it returns `proposal.md`, so the error names the file the change
// should have rather than one the flow no longer writes.
//
// ponytail: the tasks.md and plan.md fallbacks go when no `.agents/changes/*/tasks.md` or
// `.agents/changes/*/plan.md` anywhere holds a checkbox.
func checklistPath(projectDir, change string) string {
	dir := filepath.Join(projectDir, ".agents", "changes", change)
	doc := filepath.Join(dir, "proposal.md")
	if c, err := readPlan(doc); err == nil && c.open+c.done > 0 {
		return doc
	}
	for _, name := range []string{"tasks.md", "plan.md"} {
		p := filepath.Join(dir, name)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return doc
}

// openChecklist is both refusals, shared by loop() — which must run them before the PATH
// check — and runLoop, which the tests drive directly.
func openChecklist(projectDir, change string) (string, planCount, error) {
	checklist := checklistPath(projectDir, change)
	c, err := readPlan(checklist)
	if err != nil {
		return checklist, c, fmt.Errorf("no checklist at %s\n"+
			"       the loop drives a change whose checklist already exists", checklist)
	}
	if c.open == 0 && c.done == 0 {
		return checklist, c, fmt.Errorf("%s has no task boxes — nothing to drive", checklist)
	}
	return checklist, c, nil
}

func runLoop(w io.Writer, projectDir string, o loopOpts, run loopRunner) error {
	checklist, before, err := openChecklist(projectDir, o.change)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n  %s — %d open, %d done, at most %d iteration(s)\n\n",
		o.change, before.open, before.done, o.max)

	stuck := 0
	for i := 1; before.open > 0 && i <= o.max; i++ {
		fmt.Fprintf(w, "  ── %d/%d · %d box(es) open\n", i, o.max, before.open)

		if o.dryRun {
			fmt.Fprintf(w, "\n%s\n\n", loopPrompt(o.change, filepath.Base(checklist)))
			return nil
		}

		// A failed session is not a failed loop. The checklist is the state, so the next
		// round reads whatever this one managed to finish and carries on from there;
		// what an exit code cannot distinguish is a crash from a deliberate stop.
		// Stuck detection is what catches the difference, by looking at the boxes.
		if err := run(loopPrompt(o.change, filepath.Base(checklist))); err != nil {
			fmt.Fprintf(w, "     session exited with an error: %v\n", err)
		}

		after, err := readPlan(checklist)
		if err != nil {
			return fmt.Errorf("checklist disappeared mid-loop: %w", err)
		}
		// Boxes *closed*, not boxes remaining. A session that finishes one task and splits
		// another in two leaves `open` unchanged, and reading that as no progress stops the
		// loop with "two rounds closed nothing" when two things were done — the opposite of
		// what happened. A plan is allowed to grow: phase 6 discovering a task is the flow
		// working, not the loop stalling.
		if after.done <= before.done {
			stuck++
			fmt.Fprintf(w, "     no box closed (%d round(s))\n", stuck)
			if stuck >= 2 {
				fmt.Fprintf(w, "\n  stopped: two rounds closed nothing. %d box(es) still open in\n"+
					"  %s — a third fresh session reads the same checklist and does the same.\n\n",
					after.open, checklist)
				return fmt.Errorf("loop made no progress in 2 rounds")
			}
		} else {
			stuck = 0
		}
		before = after
	}

	if before.open > 0 {
		fmt.Fprintf(w, "\n  stopped at the cap: %d box(es) still open. `--max` raises it.\n\n", before.open)
		return fmt.Errorf("%d task(s) unfinished after %d iteration(s)", before.open, o.max)
	}
	fmt.Fprintf(w, "\n  every box in %s is closed.\n"+
		"  Nothing was pushed — phase 8 is yours.\n\n", o.change)
	return nil
}
