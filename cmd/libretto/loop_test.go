package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeChecklist(t *testing.T, dir, change, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, ".agents", "changes", change)
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(p, name)
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func writeTasks(t *testing.T, dir, change, body string) string {
	t.Helper()
	return writeChecklist(t, dir, change, "tasks.md", body)
}

func TestCountBoxesReadsOnlyTheBox(t *testing.T) {
	in := `# Plan

- [ ] open one
- [x] done one
  - [ ] nested open
	- [X] tab-indented, capital X
* [ ] asterisk bullet
- [] malformed, not a box
- not a box at all
Prose mentioning - [ ] mid-line does not count
`
	got := countBoxes(strings.NewReader(in))
	if got.open != 3 || got.done != 2 {
		t.Fatalf("wanted 3 open / 2 done, got %d open / %d done", got.open, got.done)
	}
}

func TestLoopStopsWhenEveryBoxIsClosed(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n")

	runs := 0
	run := func(string) error {
		runs++
		// Each session closes exactly one box, as the prompt instructs.
		b, _ := os.ReadFile(plan)
		if err := os.WriteFile(plan, []byte(strings.Replace(string(b), "- [ ]", "- [x]", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		return nil
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, run); err != nil {
		t.Fatalf("wanted a clean finish, got %v", err)
	}
	if runs != 2 {
		t.Fatalf("wanted one session per box, got %d for 2 boxes", runs)
	}
}

// The guardrail the loop exists for. A session that closes nothing leaves a plan the
// next fresh session reads identically, so it makes the identical non-progress — and
// without this the cap is the only thing between that and a burnt budget.
func TestLoopStopsAfterTwoRoundsThatCloseNothing(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n")

	runs := 0
	run := func(string) error { runs++; return nil }

	var out strings.Builder
	err := runLoop(&out, dir, loopOpts{change: "c", max: 50}, run)
	if err == nil {
		t.Fatal("wanted an error when nothing moves, got nil")
	}
	if runs != 2 {
		t.Fatalf("wanted exactly 2 rounds before giving up, got %d", runs)
	}
	if !strings.Contains(out.String(), "same checklist") {
		t.Fatalf("the stop must say why a third round is pointless, got:\n%s", out.String())
	}
}

// One bad round is a hiccup, not a stall. Resetting the counter is what lets a session
// that crashed be followed by one that works.
func TestOneBarrenRoundDoesNotStopTheLoop(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n")

	runs := 0
	run := func(string) error {
		runs++
		if runs == 1 {
			return nil // closed nothing
		}
		return os.WriteFile(plan, []byte("- [x] one\n"), 0o644)
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, run); err != nil {
		t.Fatalf("a single barren round must not stop the loop, got %v", err)
	}
	if runs != 2 {
		t.Fatalf("wanted 2 rounds, got %d", runs)
	}
}

// Progress is boxes closed, not boxes remaining. A session that finishes one task and
// splits another in two leaves `open` unchanged — and reading that as no progress stops
// the loop saying "two rounds closed nothing" when two things were done.
func TestAPlanThatGrowsIsProgressNotAStall(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n")

	runs := 0
	run := func(string) error {
		runs++
		switch runs {
		case 1: // closes one, and splits the other in two: open stays at 2
			return os.WriteFile(plan, []byte("- [x] one\n- [ ] two a\n- [ ] two b\n"), 0o644)
		case 2:
			return os.WriteFile(plan, []byte("- [x] one\n- [x] two a\n- [ ] two b\n"), 0o644)
		default:
			return os.WriteFile(plan, []byte("- [x] one\n- [x] two a\n- [x] two b\n"), 0o644)
		}
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, run); err != nil {
		t.Fatalf("a growing plan is the flow working, not the loop stalling: %v", err)
	}
	if runs != 3 {
		t.Fatalf("wanted 3 rounds for 3 eventual boxes, got %d", runs)
	}
}

// `loop <typo>` must name the plan it could not find, whatever else is missing from the
// machine. The PATH check used to run first, so the error was about a dependency
// irrelevant to the mistake actually made.
func TestTheMissingPlanIsReportedBeforeTheMissingDependency(t *testing.T) {
	t.Setenv("PATH", "")
	err := loop(t.TempDir(), []string{"nope"})
	if err == nil {
		t.Fatal("wanted a refusal")
	}
	if !strings.Contains(err.Error(), "no checklist at") {
		t.Fatalf("wanted the missing plan named, got %v", err)
	}
}

// The criterion says --dry-run needs no `claude` on PATH. Calling runLoop directly proves
// nothing about that, because the guard lives above it in loop().
func TestDryRunNeedsNoClaudeOnPath(t *testing.T) {
	t.Setenv("PATH", "")
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n")
	if err := loop(dir, []string{"c", "--dry-run"}); err != nil {
		t.Fatalf("--dry-run must work with an empty PATH, got %v", err)
	}
}

func TestLoopStopsAtTheCapAndSaysSo(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n- [ ] three\n")

	run := func(string) error {
		b, _ := os.ReadFile(plan)
		return os.WriteFile(plan, []byte(strings.Replace(string(b), "- [ ]", "- [x]", 1)), 0o644)
	}
	var out strings.Builder
	err := runLoop(&out, dir, loopOpts{change: "c", max: 2}, run)
	if err == nil {
		t.Fatal("wanted an error when the cap is hit with boxes open")
	}
	if !strings.Contains(out.String(), "--max") {
		t.Fatalf("the cap message must name the flag that raises it, got:\n%s", out.String())
	}
}

// A session that exits non-zero is not a failed loop: the plan is the state, so the
// next round carries on from whatever this one finished. Stuck detection, not the exit
// code, is what notices a session that achieved nothing.
func TestASessionErrorDoesNotAbortTheLoop(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n")

	run := func(string) error {
		if err := os.WriteFile(plan, []byte("- [x] one\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return io.ErrUnexpectedEOF
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5}, run); err != nil {
		t.Fatalf("a box closed is progress whatever the exit code, got %v", err)
	}
}

func TestLoopRefusesAPlanThatIsNotThere(t *testing.T) {
	dir := t.TempDir()
	err := runLoop(io.Discard, dir, loopOpts{change: "nope", max: 1}, func(string) error {
		t.Fatal("no session may start without a plan")
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "no checklist at") {
		t.Fatalf("wanted a refusal naming the missing plan, got %v", err)
	}
}

// A plan with prose and no boxes is not a finished plan — it is a plan phase 5 never
// wrote. Reading it as "0 open, therefore done" would report success having run nothing.
func TestLoopRefusesAPlanWithNoBoxes(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "# Plan\n\nSome prose and no checkboxes.\n")
	err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 1}, func(string) error {
		t.Fatal("no session may start against a plan with no tasks")
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "no task boxes") {
		t.Fatalf("wanted a refusal naming the empty plan, got %v", err)
	}
}

func TestDryRunPrintsThePromptAndStartsNoSession(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n")
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 5, dryRun: true}, func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	}); err != nil {
		t.Fatalf("--dry-run must not fail, got %v", err)
	}
	if !strings.Contains(out.String(), "libretto-attacca") {
		t.Fatalf("--dry-run must show the prompt it would send, got:\n%s", out.String())
	}
}

func TestParseLoopArgs(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want loopOpts
		err  string
	}{
		{name: "bare change", args: []string{"add-thing"}, want: loopOpts{change: "add-thing", max: 10}},
		{name: "max", args: []string{"add-thing", "--max", "3"}, want: loopOpts{change: "add-thing", max: 3}},
		{name: "dry run", args: []string{"--dry-run", "add-thing"}, want: loopOpts{change: "add-thing", max: 10, dryRun: true}},
		{name: "no change", args: nil, err: "which change"},
		{name: "two changes", args: []string{"a", "b"}, err: "one change at a time"},
		{name: "max needs a value", args: []string{"a", "--max"}, err: "needs a number"},
		{name: "max rejects zero", args: []string{"a", "--max", "0"}, err: "positive number"},
		{name: "unknown flag", args: []string{"a", "--force"}, err: "unknown flag"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseLoopArgs(tc.args)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("wanted an error containing %q, got %v", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("wanted %+v, got %+v", tc.want, got)
			}
		})
	}
}

// The checklist used to be `plan.md` and is now `tasks.md`, and a change created before
// that rename still has the old name on disk. The loop reads it — a rename that stranded
// work already in flight would have cost more than the name was worth. It also has to
// name the *current* file when neither exists, or the error teaches the reader to create
// the file the flow no longer writes.
func TestChecklistPathFallsBackToPlan(t *testing.T) {
	dir := t.TempDir()
	writeChecklist(t, dir, "legacy", "plan.md", "- [ ] one\n")
	if got := checklistPath(dir, "legacy"); filepath.Base(got) != "plan.md" {
		t.Errorf("a change with only plan.md must still be drivable, got %s", got)
	}

	writeChecklist(t, dir, "both", "plan.md", "- [ ] stale\n")
	writeChecklist(t, dir, "both", "tasks.md", "- [ ] current\n")
	if got := checklistPath(dir, "both"); filepath.Base(got) != "tasks.md" {
		t.Errorf("tasks.md wins when both are there, got %s", got)
	}

	if got := checklistPath(dir, "neither"); filepath.Base(got) != "proposal.md" {
		t.Errorf("with neither present the error must name the current file, got %s", got)
	}
	writeChecklist(t, dir, "queued", "proposal.md", "## Task breakdown\n- [ ] not the checklist\n")
	if got := checklistPath(dir, "queued"); filepath.Base(got) != "proposal.md" {
		t.Errorf("boxes outside ## Checklist still name the current file, got %s", got)
	}

	// And the prompt has to carry the file that was actually found, not a constant.
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "legacy", max: 1, dryRun: true}, func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	}); err != nil {
		t.Fatalf("a legacy plan must drive the loop, got %v", err)
	}
	if !strings.Contains(out.String(), "plan.md") {
		t.Errorf("the prompt must point the session at the file that exists:\n%s", out.String())
	}
}

// The loop hands work to the flow; it never restates what the flow already decides.
// A prompt that grew instructions would be a second copy of a skill, drifting from the
// installed one — and it must never authorise what attacca answers for one branch only.
func TestTheLoopPromptRoutesAndNeverAuthorises(t *testing.T) {
	p := loopPrompt("add-thing", "tasks.md")
	for _, want := range []string{"libretto-attacca", "add-thing", "tasks.md", "FIRST unchecked box"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt must contain %q:\n%s", want, p)
		}
	}
	for _, forbidden := range []string{"push", "merge", "--force", "release"} {
		if strings.Contains(strings.ToLower(p), forbidden) && forbidden != "push" {
			t.Errorf("prompt must not authorise %q:\n%s", forbidden, p)
		}
	}
	if !strings.Contains(p, "do not push") {
		t.Errorf("prompt must forbid pushing, not merely omit it:\n%s", p)
	}
}

func TestLoopIsNotGatedOnThePayload(t *testing.T) {
	if needsPayload([]string{"loop", "add-thing"}) {
		t.Fatal("loop reads the project's plan, not the payload tree — gating it refuses " +
			"the loop on exactly the machine it is for")
	}
}

func TestOnlyTheChecklistSectionCounts(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		want planCount
	}{
		{"a box under Task breakdown is not the checklist",
			"## Task breakdown\n- [ ] cut\n## Checklist\n", planCount{}},
		{"a fenced box inside the checklist does not count",
			"## Checklist\n```\n- [ ] example\n```\n", planCount{}},
		{"a ### subheading does not close the section",
			"## Checklist\n### Layer 1\n- [ ] one\n", planCount{open: 1}},
		{"the next ## heading closes it",
			"## Checklist\n- [x] in\n## Notes\n- [ ] out\n", planCount{done: 1}},
		{"a # heading does not close it",
			"## Checklist\n- [ ] a\n# Title\n- [ ] b\n", planCount{open: 2}},
		{"every Checklist section counts",
			"## Checklist\n- [ ] a\n## Notes\n- [ ] no\n## Checklist\n- [X] b\n", planCount{open: 1, done: 1}},
		{"a heading inside a fence opens nothing",
			"```\n## Checklist\n```\n- [ ] a\n", planCount{}},
		{"trailing space still names the section",
			"## Checklist  \n- [ ] a\n", planCount{open: 1}},
		{"a longer heading is another section",
			"## Checklists\n- [ ] a\n", planCount{}},
		// C2: no heading means no checklist, whatever boxes the file holds.
		{"no Checklist heading counts nothing",
			"# Proposal\n- [ ] a\n- [x] b\n## Task breakdown\n- [ ] c\n", planCount{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := countChecklist(strings.NewReader(tc.in)); got != tc.want {
				t.Fatalf("wanted %+v, got %+v", tc.want, got)
			}
		})
	}
}

// The change document wins only when it holds the checklist. Its mere existence is never
// the signal: every change in both layouts has a proposal.md.
func TestChecklistPathPrefersTheChangeDocument(t *testing.T) {
	dir := t.TempDir()
	doc := "# Proposal\n\n## Checklist\n- [ ] one\n"
	for _, tc := range []struct {
		change string
		others []string
	}{
		{"alone", nil},
		{"beside-tasks", []string{"tasks.md"}},
		{"beside-plan", []string{"plan.md"}},
		{"beside-both", []string{"tasks.md", "plan.md"}},
	} {
		writeChecklist(t, dir, tc.change, "proposal.md", doc)
		for _, o := range tc.others {
			writeChecklist(t, dir, tc.change, o, "- [ ] old layout\n")
		}
		if got := checklistPath(dir, tc.change); filepath.Base(got) != "proposal.md" {
			t.Errorf("%s: a proposal.md with a checklist must win, got %s", tc.change, got)
		}
	}
}

// A queued idea is a proposal.md with no checklist. Naming it as missing `claude` would
// send the reader to install a dependency for a change that has nothing to drive.
func TestLoopRefusesAQueuedProposal(t *testing.T) {
	t.Setenv("PATH", "")
	dir := t.TempDir()
	p := writeChecklist(t, dir, "c", "proposal.md", "Queued: 2026-10-01\n\n# An idea\n\n- [ ] not a checklist\n")
	err := loop(dir, []string{"c"})
	if err == nil {
		t.Fatal("wanted a refusal")
	}
	if !strings.Contains(err.Error(), "no task boxes") || !strings.Contains(err.Error(), p) {
		t.Fatalf("wanted %s named as having no task boxes, got %v", p, err)
	}
}

func TestTheLoopDrivesTheChangeDocumentsChecklist(t *testing.T) {
	dir := t.TempDir()
	doc := writeChecklist(t, dir, "c", "proposal.md",
		"# Proposal\n\n## Task breakdown\n- [ ] outside a\n- [ ] outside b\n\n## Checklist\n- [ ] one\n")
	want := []string{".agents/changes/c/proposal.md", "## Checklist"}

	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 5, dryRun: true}, func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	}); err != nil {
		t.Fatalf("a proposal.md with a checklist must drive the loop, got %v", err)
	}
	for _, w := range want {
		if !strings.Contains(out.String(), w) {
			t.Errorf("the dry-run prompt must name %q:\n%s", w, out.String())
		}
	}

	// Each session closes a box outside the section. That is not progress, so two rounds
	// of it stop the loop; counting the whole file would have read both as progress.
	var prompts []string
	run := func(p string) error {
		prompts = append(prompts, p)
		b, _ := os.ReadFile(doc)
		return os.WriteFile(doc, []byte(strings.Replace(string(b), "- [ ] outside", "- [x] outside", 1)), 0o644)
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5}, run); err == nil {
		t.Fatal("boxes closed outside ## Checklist must not count as progress")
	}
	if len(prompts) != 2 {
		t.Fatalf("wanted 2 barren rounds, got %d", len(prompts))
	}
	for _, w := range want {
		if !strings.Contains(prompts[0], w) {
			t.Errorf("the session prompt must name %q:\n%s", w, prompts[0])
		}
	}
}

// Nothing cuts the checklist any more, so a refusal naming the cutter sends the reader
// looking for a step that no longer exists.
func TestLoopRefusalsDoNotNameTheCutter(t *testing.T) {
	t.Setenv("PATH", "")
	dir := t.TempDir()
	writeTasks(t, dir, "empty", "# Tasks\n\nno boxes\n")
	never := func(string) error { t.Fatal("no session may start"); return nil }

	var errs []error
	errs = append(errs, loop(dir, []string{"nope"}))
	errs = append(errs, runLoop(io.Discard, dir, loopOpts{change: "nope", max: 1}, never))
	errs = append(errs, runLoop(io.Discard, dir, loopOpts{change: "empty", max: 1}, never))
	for _, err := range errs {
		if err == nil {
			t.Fatal("wanted a refusal")
		}
		if strings.Contains(err.Error(), "cut into tasks") {
			t.Errorf("a refusal must not name the cutter: %v", err)
		}
	}
}
