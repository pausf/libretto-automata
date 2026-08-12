package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writePlan(t *testing.T, dir, change string, body string) string {
	t.Helper()
	p := filepath.Join(dir, ".agents", "changes", change)
	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatal(err)
	}
	f := filepath.Join(p, "plan.md")
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
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
	plan := writePlan(t, dir, "c", "- [ ] one\n- [ ] two\n")

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
	writePlan(t, dir, "c", "- [ ] one\n- [ ] two\n")

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
	if !strings.Contains(out.String(), "same plan") {
		t.Fatalf("the stop must say why a third round is pointless, got:\n%s", out.String())
	}
}

// One bad round is a hiccup, not a stall. Resetting the counter is what lets a session
// that crashed be followed by one that works.
func TestOneBarrenRoundDoesNotStopTheLoop(t *testing.T) {
	dir := t.TempDir()
	plan := writePlan(t, dir, "c", "- [ ] one\n")

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

func TestLoopStopsAtTheCapAndSaysSo(t *testing.T) {
	dir := t.TempDir()
	plan := writePlan(t, dir, "c", "- [ ] one\n- [ ] two\n- [ ] three\n")

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
	plan := writePlan(t, dir, "c", "- [ ] one\n")

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
	if err == nil || !strings.Contains(err.Error(), "no plan at") {
		t.Fatalf("wanted a refusal naming the missing plan, got %v", err)
	}
}

// A plan with prose and no boxes is not a finished plan — it is a plan phase 5 never
// wrote. Reading it as "0 open, therefore done" would report success having run nothing.
func TestLoopRefusesAPlanWithNoBoxes(t *testing.T) {
	dir := t.TempDir()
	writePlan(t, dir, "c", "# Plan\n\nSome prose and no checkboxes.\n")
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
	writePlan(t, dir, "c", "- [ ] one\n")
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

// The loop hands work to the flow; it never restates what the flow already decides.
// A prompt that grew instructions would be a second copy of a skill, drifting from the
// installed one — and it must never authorise what attacca answers for one branch only.
func TestTheLoopPromptRoutesAndNeverAuthorises(t *testing.T) {
	p := loopPrompt("add-thing")
	for _, want := range []string{"libretto-attacca", "add-thing", "plan.md", "FIRST unchecked box"} {
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
