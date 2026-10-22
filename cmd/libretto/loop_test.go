package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

// A box shown as an example inside a fence is not work, in tasks.md as in proposal.md —
// find-work's scan skips it, and the loop must count what phase 1 counts.
func TestCountBoxesSkipsFences(t *testing.T) {
	in := "- [ ] real\n```\n- [ ] example\n- [x] example\n```\n- [x] real\n"
	got := countBoxes(strings.NewReader(in))
	if got.open != 1 || got.done != 1 {
		t.Fatalf("wanted 1 open / 1 done outside the fence, got %d open / %d done", got.open, got.done)
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
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, textOnly(run)); err != nil {
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
	err := runLoop(&out, dir, loopOpts{change: "c", max: 50}, textOnly(run))
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
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, textOnly(run)); err != nil {
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
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10}, textOnly(run)); err != nil {
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
	err := runLoop(&out, dir, loopOpts{change: "c", max: 2}, textOnly(run))
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
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5}, textOnly(run)); err != nil {
		t.Fatalf("a box closed is progress whatever the exit code, got %v", err)
	}
}

func TestLoopRefusesAPlanThatIsNotThere(t *testing.T) {
	dir := t.TempDir()
	err := runLoop(io.Discard, dir, loopOpts{change: "nope", max: 1}, textOnly(func(string) error {
		t.Fatal("no session may start without a plan")
		return nil
	}))
	if err == nil || !strings.Contains(err.Error(), "no checklist at") {
		t.Fatalf("wanted a refusal naming the missing plan, got %v", err)
	}
}

// A plan with prose and no boxes is not a finished plan — it is a plan phase 5 never
// wrote. Reading it as "0 open, therefore done" would report success having run nothing.
func TestLoopRefusesAPlanWithNoBoxes(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "# Plan\n\nSome prose and no checkboxes.\n")
	err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 1}, textOnly(func(string) error {
		t.Fatal("no session may start against a plan with no tasks")
		return nil
	}))
	if err == nil || !strings.Contains(err.Error(), "no task boxes") {
		t.Fatalf("wanted a refusal naming the empty plan, got %v", err)
	}
}

func TestDryRunPrintsThePromptAndStartsNoSession(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n")
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 5, dryRun: true}, textOnly(func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	})); err != nil {
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
	if err := runLoop(&out, dir, loopOpts{change: "legacy", max: 1, dryRun: true}, textOnly(func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	})); err != nil {
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
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 5, dryRun: true}, textOnly(func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	})); err != nil {
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
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5}, textOnly(run)); err == nil {
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
	errs = append(errs, runLoop(io.Discard, dir, loopOpts{change: "nope", max: 1}, textOnly(never)))
	errs = append(errs, runLoop(io.Discard, dir, loopOpts{change: "empty", max: 1}, textOnly(never)))
	for _, err := range errs {
		if err == nil {
			t.Fatal("wanted a refusal")
		}
		if strings.Contains(err.Error(), "cut into tasks") {
			t.Errorf("a refusal must not name the cutter: %v", err)
		}
	}
}

// textOnly adapts a fake that only touches the checklist to the runner, which also
// returns what the session reported. Those fakes report nothing, which is the point.
func textOnly(f func(string) error) loopRunner {
	return func(p string) (session, error) { return session{}, f(p) }
}

// closingRunner closes one box per round and reports whatever out holds for that round,
// parsed the way execClaude parses a real session's stdout.
func closingRunner(t *testing.T, plan string, out ...string) loopRunner {
	t.Helper()
	round := 0
	return func(string) (session, error) {
		b, _ := os.ReadFile(plan)
		if err := os.WriteFile(plan, []byte(strings.Replace(string(b), "- [ ]", "- [x]", 1)), 0o644); err != nil {
			t.Fatal(err)
		}
		s := ""
		if round < len(out) {
			s = out[round]
		}
		round++
		return parseSession([]byte(s)), nil
	}
}

// readRecords decodes the run log into maps, so a test can tell an absent key from a
// zero one. Numbers stay json.Number: the host's value, not a float's rendering of it.
func readRecords(t *testing.T, path string) []map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the run log: %v", err)
	}
	var recs []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(string(b), "\n"), "\n") {
		d := json.NewDecoder(strings.NewReader(line))
		d.UseNumber()
		var r map[string]any
		if err := d.Decode(&r); err != nil {
			t.Fatalf("a run log line is not a JSON object: %q: %v", line, err)
		}
		recs = append(recs, r)
	}
	return recs
}

func logIn(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "libretto", "loop.jsonl")
}

const hostJSON = `{"type":"result","subtype":"success","is_error":false,"duration_ms":4567,` +
	`"num_turns":3,"result":"closed the first box","session_id":"abc-123","total_cost_usd":0.1234}`

func TestEverySessionIsAskedForItsJSONResult(t *testing.T) {
	args := claudeArgs("the prompt", "")
	pair := func(k, v string) bool {
		for i := 0; i+1 < len(args); i++ {
			if args[i] == k && args[i+1] == v {
				return true
			}
		}
		return false
	}
	if !pair("-p", "the prompt") {
		t.Errorf("wanted -p followed by the prompt, got %q", args)
	}
	if !pair("--output-format", "json") {
		t.Errorf("wanted --output-format json, got %q", args)
	}
}

func TestParseSessionResultReadsTheHostsCost(t *testing.T) {
	s := parseSession([]byte(hostJSON + "\n"))
	if s.cost == nil || *s.cost != 0.1234 {
		t.Errorf("cost: wanted 0.1234, got %v", s.cost)
	}
	if s.durationMS == nil || *s.durationMS != 4567 {
		t.Errorf("duration_ms: wanted 4567, got %v", s.durationMS)
	}
	if s.sessionID != "abc-123" || s.subtype != "success" || s.text != "closed the first box" {
		t.Errorf("wanted the host's session_id, subtype and result, got %+v", s)
	}
	if s.isError == nil || *s.isError {
		t.Errorf("is_error: wanted a reported false, got %v", s.isError)
	}

	for _, tc := range []struct{ name, out string }{
		{"not JSON", "plain text from an older host\n"},
		{"JSON with no total_cost_usd", `{"result":"hi","session_id":"x"}`},
		{"empty output", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if s := parseSession([]byte(tc.out)); s.cost != nil {
				t.Fatalf("an unreported cost must be unknown, never %v", *s.cost)
			}
		})
	}
	if s := parseSession([]byte(`{"total_cost_usd":0}`)); s.cost == nil || *s.cost != 0 {
		t.Errorf("a reported zero is a fact, not an unknown: got %v", s.cost)
	}
}

// JSON that carries no result field is not the host's report, so what the session said
// is the raw output — a bare `null` or `{}` must not print nothing.
func TestJSONWithoutAResultPrintsTheRawOutput(t *testing.T) {
	for _, out := range []string{"null", "{}", `{"other":1}`} {
		if got := parseSession([]byte(out)).text; got != out {
			t.Errorf("%s: wanted the raw output printed, got %q", out, got)
		}
	}
}

func TestTheSessionsTextStillReachesTheTerminal(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n")
	var out strings.Builder
	run := closingRunner(t, plan, hostJSON, "raw output, not JSON\n")
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 5, log: logIn(t)}, run); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "closed the first box") {
		t.Errorf("the JSON's result text must be printed:\n%s", out.String())
	}
	if strings.Contains(out.String(), `"total_cost_usd"`) {
		t.Errorf("the JSON itself must not be printed, only its result:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "raw output, not JSON\n") {
		t.Errorf("output that is not JSON must be printed unchanged:\n%s", out.String())
	}
}

func TestEveryRoundAppendsOneRecord(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n- [ ] three\n")
	log := logIn(t)
	if err := os.MkdirAll(filepath.Dir(log), 0o755); err != nil {
		t.Fatal(err)
	}
	// An earlier run's record: the log is appended to, never rewritten.
	if err := os.WriteFile(log, []byte(`{"change":"earlier","round":1}`+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, closingRunner(t, plan)); err != nil {
		t.Fatal(err)
	}
	recs := readRecords(t, log)
	if len(recs) != 4 {
		t.Fatalf("wanted the earlier record plus one per round (4), got %d", len(recs))
	}
	if recs[0]["change"] != "earlier" {
		t.Errorf("the earlier record must survive, got %v", recs[0])
	}
}

func TestTheRecordSaysWhenWhichChangeAndWhichRound(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "add-thing", "- [ ] one\n- [ ] two\n")
	log := logIn(t)
	start := time.Now().Truncate(time.Second)
	if err := runLoop(io.Discard, dir, loopOpts{change: "add-thing", max: 10, log: log}, closingRunner(t, plan)); err != nil {
		t.Fatal(err)
	}
	end := time.Now()
	for i, r := range readRecords(t, log) {
		ts, ok := r["time"].(string)
		if !ok {
			t.Fatalf("record %d has no time: %v", i, r)
		}
		when, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t.Fatalf("record %d: time %q is not RFC 3339: %v", i, ts, err)
		}
		if when.Before(start) || when.After(end) {
			t.Errorf("record %d: time %s is outside the run [%s, %s]", i, when, start, end)
		}
		if r["change"] != "add-thing" {
			t.Errorf("record %d: change %v", i, r["change"])
		}
		if r["round"] != json.Number(fmt.Sprint(i+1)) {
			t.Errorf("record %d: wanted round %d, got %v", i, i+1, r["round"])
		}
	}
}

// The fixture must reopen a box, or an unsigned build passes this test.
func TestClosedIsSignedSoAReopenedBoxIsNegative(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [x] a\n- [x] b\n- [ ] c\n")
	log := logIn(t)
	round := 0
	run := func(string) (session, error) {
		round++
		body := "- [ ] a\n- [x] b\n- [ ] c\n" // reopens a
		if round == 2 {
			body = "- [x] a\n- [x] b\n- [x] c\n"
		}
		return session{}, os.WriteFile(plan, []byte(body), 0o644)
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, run); err != nil {
		t.Fatal(err)
	}
	recs := readRecords(t, log)
	if len(recs) != 2 {
		t.Fatalf("wanted 2 records, got %d", len(recs))
	}
	for i, want := range []struct{ closed, open string }{{"-1", "2"}, {"2", "0"}} {
		if recs[i]["closed"] != json.Number(want.closed) || recs[i]["open"] != json.Number(want.open) {
			t.Errorf("record %d: wanted closed %s open %s, got closed %v open %v",
				i, want.closed, want.open, recs[i]["closed"], recs[i]["open"])
		}
	}
}

func TestTheRecordCarriesTheHostsReportVerbatim(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n")
	log := logIn(t)
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, closingRunner(t, plan, hostJSON)); err != nil {
		t.Fatal(err)
	}
	r := readRecords(t, log)[0]
	for k, want := range map[string]any{
		"cost_usd":    json.Number("0.1234"),
		"duration_ms": json.Number("4567"),
		"session_id":  "abc-123",
		"subtype":     "success",
		"is_error":    false,
	} {
		if r[k] != want {
			t.Errorf("%s: wanted %v, got %v", k, want, r[k])
		}
	}
}

// A plain float64 field writes 0, and a test that only checks the record exists passes.
func TestAnUnpricedRoundOmitsCostRatherThanWritingZero(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n")
	log := logIn(t)
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, closingRunner(t, plan, "not JSON\n")); err != nil {
		t.Fatal(err)
	}
	recs := readRecords(t, log)
	if len(recs) != 1 {
		t.Fatalf("wanted the round recorded, got %d records", len(recs))
	}
	for _, k := range []string{"cost_usd", "is_error"} {
		if v, ok := recs[0][k]; ok {
			t.Errorf("an unreported %s must be absent, got %v", k, v)
		}
	}
}

func TestTheRecordCarriesTheExecErrorOnlyOnFailure(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n")
	log := logIn(t)
	closeOne := closingRunner(t, plan)
	round := 0
	run := func(p string) (session, error) {
		round++
		s, _ := closeOne(p)
		if round == 1 {
			return s, io.ErrUnexpectedEOF
		}
		return s, nil
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, run); err != nil {
		t.Fatal(err)
	}
	recs := readRecords(t, log)
	if recs[0]["error"] != io.ErrUnexpectedEOF.Error() {
		t.Errorf("a non-zero exit must record its message, got %v", recs[0]["error"])
	}
	if v, ok := recs[1]["error"]; ok {
		t.Errorf("a clean exit must record no error, got %v", v)
	}
}

func TestTheLastRecordSaysWhyTheLoopStopped(t *testing.T) {
	for _, tc := range []struct {
		name, plan, want string
		max              int
		run              func(t *testing.T, plan string) loopRunner
	}{
		{"done", "- [ ] one\n- [ ] two\n", "done", 10, func(t *testing.T, plan string) loopRunner {
			return closingRunner(t, plan)
		}},
		{"stuck", "- [ ] one\n", "stuck", 10, func(*testing.T, string) loopRunner {
			return func(string) (session, error) { return session{}, nil }
		}},
		{"cap", "- [ ] one\n- [ ] two\n- [ ] three\n", "cap", 2, func(t *testing.T, plan string) loopRunner {
			return closingRunner(t, plan)
		}},
		{"error", "- [ ] one\n- [ ] two\n", "error", 10, func(t *testing.T, plan string) loopRunner {
			closeOne, round := closingRunner(t, plan), 0
			return func(p string) (session, error) {
				if round++; round == 2 {
					return session{}, os.Remove(plan)
				}
				return closeOne(p)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			plan := writeTasks(t, dir, "c", tc.plan)
			log := logIn(t)
			_ = runLoop(io.Discard, dir, loopOpts{change: "c", max: tc.max, log: log}, tc.run(t, plan))
			recs := readRecords(t, log)
			if len(recs) < 2 {
				t.Fatalf("wanted at least two rounds to compare, got %d", len(recs))
			}
			for i, r := range recs[:len(recs)-1] {
				if v, ok := r["stopped"]; ok {
					t.Errorf("record %d is not the last and says stopped %v", i, v)
				}
			}
			if got := recs[len(recs)-1]["stopped"]; got != tc.want {
				t.Errorf("wanted the last record stopped %q, got %v", tc.want, got)
			}
		})
	}
}

func TestAVanishedChecklistIsStillRecorded(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n")
	log := logIn(t)
	err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 10, log: log}, func(string) (session, error) {
		return parseSession([]byte(hostJSON)), os.Remove(plan)
	})
	if err == nil || !strings.Contains(err.Error(), "checklist disappeared mid-loop") {
		t.Fatalf("wanted the error the loop returns today, got %v", err)
	}
	recs := readRecords(t, log)
	if len(recs) != 1 {
		t.Fatalf("wanted the vanished round recorded, got %d records", len(recs))
	}
	r := recs[0]
	for _, k := range []string{"closed", "open"} {
		if v, ok := r[k]; ok {
			t.Errorf("%s must be absent when the checklist could not be read, got %v", k, v)
		}
	}
	if r["stopped"] != "error" {
		t.Errorf("wanted stopped error, got %v", r["stopped"])
	}
	if e, _ := r["error"].(string); !strings.Contains(e, plan) {
		t.Errorf("wanted the read error naming %s, got %v", plan, r["error"])
	}
}

// Passes trivially if no log was written at all, so the record is asserted first.
func TestTheRunLogLeavesGitStatusClean(t *testing.T) {
	dir := landInit(t)
	writeLand(t, dir, ".agents/changes/c/tasks.md", "- [ ] one\n")
	gitLand(t, dir, "add", "-A")
	gitLand(t, dir, "commit", "-qm", "the change")

	log := loopLogPath(execGit(dir))
	if want := filepath.Join(".git", "libretto", "loop.jsonl"); !strings.HasSuffix(log, want) {
		t.Fatalf("wanted the log at <git-common-dir>/libretto/loop.jsonl, got %q", log)
	}
	plan := filepath.Join(dir, ".agents", "changes", "c", "tasks.md")
	run := func(string) (session, error) {
		// The session commits its box, as the prompt instructs; the tree stays clean.
		if err := os.WriteFile(plan, []byte("- [x] one\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		gitLand(t, dir, "commit", "-qam", "close one")
		return session{}, nil
	}
	if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5, log: log}, run); err != nil {
		t.Fatal(err)
	}
	if recs := readRecords(t, log); len(recs) != 1 {
		t.Fatalf("wanted one record in %s, got %d", log, len(recs))
	}
	if st := gitLand(t, dir, "status", "--porcelain", "--untracked-files=all"); st != "" {
		t.Fatalf("the run log must leave git status clean, got:\n%s", st)
	}
}

func TestAnUnlocatableLogWarnsAndCarriesOn(t *testing.T) {
	notARepo := func(...string) (string, error) { return "", fmt.Errorf("not a git repository") }
	if got := loopLogPath(notARepo); got != "" {
		t.Fatalf("an unlocatable common dir must resolve to no path, got %q", got)
	}

	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n- [ ] three\n")
	rounds := 0
	closeOne := closingRunner(t, plan)
	run := func(p string) (session, error) { rounds++; return closeOne(p) }
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 10}, run); err != nil {
		t.Fatalf("a missing log must not stop the loop, got %v", err)
	}
	if rounds != 3 {
		t.Errorf("wanted every round run, got %d of 3", rounds)
	}
	const warning = "run log not recorded: git common dir not found"
	if n := strings.Count(out.String(), warning); n != 1 {
		t.Errorf("wanted the warning exactly once, got %d:\n%s", n, out.String())
	}
}

func TestAnUnwritableLogDoesNotStopTheLoop(t *testing.T) {
	dir := t.TempDir()
	plan := writeTasks(t, dir, "c", "- [ ] one\n- [ ] two\n- [ ] three\n")
	// A regular file where the log's directory should be: every write fails.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	log := filepath.Join(blocker, "libretto", "loop.jsonl")
	rounds := 0
	closeOne := closingRunner(t, plan)
	run := func(p string) (session, error) { rounds++; return closeOne(p) }
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 10, log: log}, run); err != nil {
		t.Fatalf("an unwritable log must not stop the loop, got %v", err)
	}
	if rounds != 3 {
		t.Errorf("wanted every round run, got %d of 3", rounds)
	}
	if n := strings.Count(out.String(), log); n != 1 {
		t.Errorf("wanted one warning naming %s, got %d:\n%s", log, n, out.String())
	}
}

func TestDryRunWritesNoRecord(t *testing.T) {
	t.Setenv("PATH", "")
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n")
	never := func() string { t.Fatal("the run log path must not be resolved"); return "" }

	if err := loopWith(dir, []string{"nope"}, never); err == nil {
		t.Fatal("wanted the missing checklist refused")
	}
	if err := loopWith(dir, []string{"c"}, never); err == nil || !strings.Contains(err.Error(), "claude") {
		t.Fatalf("wanted the missing claude refused, got %v", err)
	}
	if err := loopWith(dir, []string{"c", "--dry-run"}, never); err != nil {
		t.Fatalf("--dry-run must not fail, got %v", err)
	}

	log := logIn(t)
	var out strings.Builder
	for _, o := range []loopOpts{
		{change: "c", max: 5, dryRun: true},
		{change: "c", max: 5, dryRun: true, log: log},
	} {
		if err := runLoop(&out, dir, o, textOnly(func(string) error {
			t.Fatal("--dry-run must not start a session")
			return nil
		})); err != nil {
			t.Fatal(err)
		}
	}
	if strings.Contains(out.String(), "run log") {
		t.Errorf("--dry-run must print no log warning:\n%s", out.String())
	}
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		t.Errorf("--dry-run must write no record, stat says %v", err)
	}
}

func TestClaudeArgsCarryTheBudgetOnlyWhenGiven(t *testing.T) {
	with := claudeArgs("p", "2.50")
	found := false
	for i := 0; i+1 < len(with); i++ {
		if with[i] == "--max-budget-usd" && with[i+1] == "2.50" {
			found = true
		}
	}
	if !found {
		t.Errorf("wanted --max-budget-usd 2.50 exactly as typed, got %q", with)
	}
	for _, a := range claudeArgs("p", "") {
		if a == "--max-budget-usd" {
			t.Errorf("no budget given, yet the args carry one: %q", claudeArgs("p", ""))
		}
	}
}

func TestParseLoopArgsRefusesABadBudget(t *testing.T) {
	for _, v := range []string{"1e2", "0x1", "+3", "-1", "inf", ".5", "0", "0.00", "1.", "2,50", ""} {
		t.Run(v, func(t *testing.T) {
			_, err := parseLoopArgs([]string{"c", "--max-budget-usd", v})
			if err == nil || !strings.Contains(err.Error(), "--max-budget-usd") {
				t.Fatalf("wanted %q refused naming --max-budget-usd, got %v", v, err)
			}
		})
	}
	t.Run("missing value", func(t *testing.T) {
		_, err := parseLoopArgs([]string{"c", "--max-budget-usd"})
		if err == nil || !strings.Contains(err.Error(), "--max-budget-usd") {
			t.Fatalf("wanted a missing value refused naming --max-budget-usd, got %v", err)
		}
	})
	for _, v := range []string{"2.50", "0.1", "3"} {
		got, err := parseLoopArgs([]string{"c", "--max-budget-usd", v})
		if err != nil || got.budget != v {
			t.Errorf("wanted %q accepted as typed, got %q, %v", v, got.budget, err)
		}
	}
}

func TestTheHeaderStatesTheWorstCaseSpend(t *testing.T) {
	dir := t.TempDir()
	writeTasks(t, dir, "c", "- [ ] one\n")
	var out strings.Builder
	if err := runLoop(&out, dir, loopOpts{change: "c", max: 3, budget: "0.1", dryRun: true}, textOnly(func(string) error {
		t.Fatal("--dry-run must not start a session")
		return nil
	})); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"0.1", "0.30"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the header must state %q:\n%s", want, out.String())
		}
	}

	var none strings.Builder
	if err := runLoop(&none, dir, loopOpts{change: "c", max: 3, dryRun: true}, textOnly(func(string) error { return nil })); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(none.String(), "worst case") {
		t.Errorf("no budget, no worst case:\n%s", none.String())
	}
}

func TestTheRecordCarriesTheBudgetThatApplied(t *testing.T) {
	for _, tc := range []struct {
		name, budget string
		want         any
	}{
		{"set", "2.50", json.Number("2.5")},
		{"unset", "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			plan := writeTasks(t, dir, "c", "- [ ] one\n")
			log := logIn(t)
			if err := runLoop(io.Discard, dir, loopOpts{change: "c", max: 5, budget: tc.budget, log: log}, closingRunner(t, plan)); err != nil {
				t.Fatal(err)
			}
			v, ok := readRecords(t, log)[0]["budget_usd"]
			if tc.want == nil {
				if ok {
					t.Fatalf("no budget set, yet the record carries budget_usd %v", v)
				}
				return
			}
			if v != tc.want {
				t.Fatalf("wanted budget_usd %v, got %v", tc.want, v)
			}
		})
	}
}
