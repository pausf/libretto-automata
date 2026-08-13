package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// fakeGit answers the three shapes metrics asks for, keyed by the argument that
// distinguishes them. Driving real parsing against fixed output beats building
// repositories: the parsing is the logic, and a fixture makes a reopened box
// expressible in two lines.
func fakeGit(added string, logs, diffs map[string]string) gitRunner {
	return func(args ...string) (string, error) {
		joined := strings.Join(args, " ")
		if strings.Contains(joined, "--show-toplevel") {
			return "/repo\n", nil
		}
		// Git's default history simplification prunes a path whose final state is
		// "deleted" — which is every landed change, i.e. every change worth measuring.
		// Refusing here means forgetting the flag fails a test rather than silently
		// reporting a repository with almost no history.
		if !strings.Contains(joined, "--full-history") {
			return "", fmt.Errorf("every query must pass --full-history: %q", joined)
		}
		switch {
		case strings.Contains(joined, "--diff-filter=A "), strings.HasSuffix(joined, "--diff-filter=A"):
			return added, nil
		case strings.Contains(joined, "-p"):
			// The deletion filter is not cosmetic: without it a landed change's folder
			// removal shows every closed box as a `-[x]`. The fake refuses to answer a
			// query that forgot it, so dropping the flag fails a test instead of quietly
			// inverting every number.
			if !strings.Contains(joined, "--diff-filter=AM") {
				return "", fmt.Errorf("the churn query must exclude deletions: %q", joined)
			}
			for k, v := range diffs {
				if strings.Contains(joined, k+"/plan.md") {
					return v, nil
				}
			}
			return "", nil
		default:
			for k, v := range logs {
				if strings.HasSuffix(joined, "changes/"+k) {
					return v, nil
				}
			}
			return "", fmt.Errorf("no log for %q", joined)
		}
	}
}

func TestChangeNamesSeesLandedChangesNotJustOpenOnes(t *testing.T) {
	// The point of reading git rather than the filesystem: a landed change's folder is
	// deleted, so a directory listing reports only what is in flight — which is the half
	// `status` already shows, and the half with nothing to learn from.
	added := `.agents/changes/add-thing/proposal.md
.agents/changes/add-thing/plan.md
.agents/changes/fix-other/proposal.md
docs/FLOW.md
`
	got, err := changeNames(fakeGit(added, nil, nil), "/repo")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"add-thing", "fix-other"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("wanted %v, got %v", want, got)
	}
}

func TestChangeNamesIgnoresPathsOutsideTheChangeTree(t *testing.T) {
	got, err := changeNames(fakeGit("docs/FLOW.md\nREADME.md\n.agents/specs/cli/spec.md\n", nil, nil), "/repo")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("wanted nothing, got %v", got)
	}
}

func TestMeasureReadsSpanAndCommitsOldestLast(t *testing.T) {
	// git log is newest-first. Reading it as oldest-first inverts the span into a
	// negative duration, which prints as a plausible-looking number.
	now := time.Now().Unix()
	logs := map[string]string{"c": fmt.Sprintf("%d\n%d\n%d\n", now, now-3600, now-7200)}
	m, err := measureAt(fakeGit("", logs, nil), "c")
	if err != nil {
		t.Fatal(err)
	}
	if m.commits != 3 {
		t.Fatalf("wanted 3 commits, got %d", m.commits)
	}
	if m.span() < 0 {
		t.Fatalf("span must run oldest to newest, got %v", m.span())
	}
	if int(m.span().Hours()) != 2 {
		t.Fatalf("wanted a 2h span, got %v", m.span())
	}
}

// The metric the whole report is for. A box closed and reopened means a task was called
// done before it was; a net count hides exactly that and reports a tidy plan.
//
// The fixture carries commit SHAs because the arithmetic is per-commit and the boundary
// is what makes it correct. An earlier fixture ran the whole history together as one
// diff, which is not what `git log -p` emits, and it made the test pass while the code
// under it counted a reworded task as a reopening.
func TestReopenedBoxesAreCountedSeparatelyFromClosedOnes(t *testing.T) {
	now := time.Now().Unix()
	diff := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
+- [x] one
+- [x] two
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
-- [x] two
+- [ ] two
cccccccccccccccccccccccccccccccccccccccc
-- [ ] two
+- [x] two
`
	m, err := measureAt(fakeGit("",
		map[string]string{"c": fmt.Sprintf("%d\n", now)},
		map[string]string{"c": diff}), "c")
	if err != nil {
		t.Fatal(err)
	}
	if m.uncheck != 1 {
		t.Fatalf("wanted 1 reopened box, got %d", m.uncheck)
	}
	// Commit a closed two boxes, commit c closed the reopened one again: three closes.
	if m.checked != 3 {
		t.Fatalf("wanted 3 boxes closed, got %d", m.checked)
	}
}

// The two false accusations review found. A reopening means "a task was called done
// before it was" — rewording a finished task's text is not that, and neither is deleting
// one. Both emit a removed `[x]` and both netted a reopening under line-counting.
func TestRewordingAndDeletingAFinishedTaskAreNotReopenings(t *testing.T) {
	reword := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
+- [x] one
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
-- [x] one
+- [x] one, reworded
`
	if closed, reopened, _ := churn(reword); closed != 1 || reopened != 0 {
		t.Errorf("rewording a done task: wanted 1 closed / 0 reopened, got %d / %d", closed, reopened)
	}

	// A deleted done task nets -1 in its own commit. That IS a loss of finished work, so
	// reporting it is honest — but it must not also inflate the close count.
	del := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
+- [x] one
+- [x] two
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
-- [x] two
`
	if closed, reopened, boxes := churn(del); closed != 2 || reopened != 1 || boxes != 1 {
		t.Errorf("deleting a done task: wanted 2 closed / 1 reopened / 1 box left, got %d / %d / %d", closed, reopened, boxes)
	}
}

// `+++ b/plan.md` and `--- a/plan.md` are diff headers, not content. Reading them as
// lines starting with + or - is harmless only until a filename contains a checkbox.
func TestChurnIgnoresDiffHeaders(t *testing.T) {
	d := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
--- a/plan.md
+++ b/plan.md
@@ -1 +1 @@
+- [x] one
`
	if closed, reopened, _ := churn(d); closed != 1 || reopened != 0 {
		t.Errorf("wanted 1 closed / 0 reopened, got %d / %d", closed, reopened)
	}
}

func TestBoxInReadsOnlyRealCheckboxes(t *testing.T) {
	for _, tc := range []struct{ line, want string }{
		{"+- [x] done", "x"},
		{"+- [X] done", "x"},
		{"+  - [ ] nested open", " "},
		{"+- [] malformed", ""},
		{"+nothing here", ""},

		// Anchored, so the two commands cannot disagree about what a box is. `loop`'s
		// parser has always required the box at the start of the line; this counted one
		// anywhere, so the same plan.md gave two different totals.
		{"+prose about - [x] a box", ""},
		{"+  - [x] indented is still a box", "x"},

		// The truncation boundary. The old bounds guard was `i+4 > len(l)` where it needed
		// `>=`, so every one of these panicked with index out of range. The old table
		// stopped at "+- [", exactly one character short of the crash — a boundary case
		// that reads as covered and is not.
		{"+- [x", ""},
		{"+- [ ", ""},
		{"+- []", ""},
		{"+-", ""},
		{"+", ""},
		{"", ""},
	} {
		if got := boxIn(tc.line); got != tc.want {
			t.Errorf("boxIn(%q) = %q, wanted %q", tc.line, got, tc.want)
		}
	}
}

func TestAChangeWithNoPlanReportsADashNotAZero(t *testing.T) {
	// Zero closed boxes and no plan at all are different facts. Printing 0 for both
	// says a plan existed and nothing got done.
	//
	// The assertion is on the change's own row, not on the page. An earlier version
	// searched the whole output for an em dash — which flowCeiling prints on every run,
	// so mutating planCell to return "0" left the test green. A test whose subject also
	// appears in boilerplate is asserting on the boilerplate.
	now := time.Now().Unix()
	var out strings.Builder
	err := metrics(&out, nil, fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil))
	if err != nil {
		t.Fatal(err)
	}
	row := rowFor(t, out.String(), "c")
	if strings.Contains(row, "0") {
		t.Fatalf("a change with no plan must not report 0 closed boxes, got row:\n%s", row)
	}
	if !strings.Contains(row, "—") {
		t.Fatalf("wanted a dash in the row, got:\n%s", row)
	}
}

// measureAt pins the repository root the tests measure against. The root is a real
// parameter now: git pathspecs and os.Stat are both cwd-relative, so reading them from
// the process cwd made `metrics` report "no changes yet" when run from a subdirectory.
func measureAt(git gitRunner, name string) (changeMetrics, error) {
	return measure(git, "/repo", name)
}

// rowFor returns the single output line reporting on one change, so an assertion cannot
// accidentally be satisfied by the header, the footer or the ceiling note.
func rowFor(t *testing.T, out, change string) string {
	t.Helper()
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), change+" ") {
			return l
		}
	}
	t.Fatalf("no row for %q in:\n%s", change, out)
	return ""
}

// The silent-skip guard. `measure` failing must still print a row: the footer counts
// names, so a skip makes the total disagree with the rows above it.
func TestAnUnreadableChangePrintsARowRatherThanVanishing(t *testing.T) {
	now := time.Now().Unix()
	// `b` is listed as added but has no log, so measure returns an error.
	g := fakeGit(".agents/changes/a/proposal.md\n.agents/changes/b/proposal.md\n",
		map[string]string{"a": fmt.Sprintf("%d\n", now)}, nil)

	var out strings.Builder
	if err := metrics(&out, nil, g); err != nil {
		t.Fatal(err)
	}
	row := rowFor(t, out.String(), "b")
	if !strings.Contains(row, "unreadable") {
		t.Fatalf("wanted b marked unreadable, got:\n%s", row)
	}
	if !strings.Contains(out.String(), "2 change(s)") {
		t.Fatalf("the footer must count both changes it listed:\n%s", out.String())
	}
}

// The report says what it cannot see. A metrics command that silently omits two of the
// three things asked of it reads as having measured them and found nothing.
func TestTheReportNamesWhatItCannotMeasure(t *testing.T) {
	now := time.Now().Unix()
	var out strings.Builder
	if err := metrics(&out, nil, fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil)); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"per-phase duration", "review-work findings", "Not measured"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the report must name %q as unmeasured:\n%s", want, out.String())
		}
	}
}

func TestMetricsFiltersToOneChangeAndRefusesAnUnknownOne(t *testing.T) {
	now := time.Now().Unix()
	g := fakeGit(".agents/changes/a/proposal.md\n.agents/changes/b/proposal.md\n",
		map[string]string{
			"a": fmt.Sprintf("%d\n", now),
			"b": fmt.Sprintf("%d\n", now),
		}, nil)

	var out strings.Builder
	if err := metrics(&out, []string{"a"}, g); err != nil {
		t.Fatal(err)
	}
	rowFor(t, out.String(), "a") // fails the test if absent
	if strings.Contains(out.String(), "\n  b ") {
		t.Fatalf("wanted only change a:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "1 change(s)") {
		t.Fatalf("the footer must count the filtered set, not the whole history:\n%s", out.String())
	}

	if err := metrics(&strings.Builder{}, []string{"nope"}, g); err == nil {
		t.Fatal("wanted a refusal for a change git never saw")
	}
	if err := metrics(&strings.Builder{}, []string{"--all"}, g); err == nil {
		t.Fatal("wanted a refusal for an unknown flag")
	}
}

// 5/5 and 5/18 are opposite facts about a change in flight, and a bare 5 hides which
// one is happening. The numerator is what is closed now, not cumulative closes —
// churn already has its own column.
func TestClosedShowsItsDenominator(t *testing.T) {
	now := time.Now().Unix()
	diff := `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
+- [x] one
+- [ ] two
+- [ ] three
bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
-- [x] one
+- [x] one, reworded
`
	g := fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)},
		map[string]string{"c": diff})
	var out strings.Builder
	if err := metrics(&out, nil, g); err != nil {
		t.Fatal(err)
	}
	row := rowFor(t, out.String(), "c")
	// The reword must move neither number: still one closed of three.
	if !strings.Contains(row, "1/3") {
		t.Fatalf("wanted the closed cell as 1/3, got row:\n%s", row)
	}
}

// The footer claims wall clock, so each calendar hour is counted once however many
// changes were open during it. Summing per-change spans reported two weeks for two
// changes open the same week — a number that was nobody's clock.
func TestTotalSpanMergesOverlappingChanges(t *testing.T) {
	now := time.Now().Unix()
	logs := map[string]string{
		// a: [now-10h, now], b: [now-5h, now] — inside a. c: [now-30h, now-28h] — disjoint.
		"a": fmt.Sprintf("%d\n%d\n", now, now-10*3600),
		"b": fmt.Sprintf("%d\n%d\n", now, now-5*3600),
		"c": fmt.Sprintf("%d\n%d\n", now-28*3600, now-30*3600),
	}
	added := ".agents/changes/a/proposal.md\n.agents/changes/b/proposal.md\n.agents/changes/c/proposal.md\n"
	var out strings.Builder
	if err := metrics(&out, nil, fakeGit(added, logs, nil)); err != nil {
		t.Fatal(err)
	}
	// Union: 10h + 2h = 12h. The naive sum says 17h.
	if !strings.Contains(out.String(), "12h of wall clock") {
		t.Fatalf("wanted the merged 12h in the footer, got:\n%s", out.String())
	}
}

// A name is typed from memory of how it starts. A unique prefix selects; an ambiguous
// one is refused naming the candidates; and an exact name always wins, or a change
// whose full name prefixes a sibling's becomes unreachable.
func TestAPrefixSelectsAChangeUnlessAmbiguous(t *testing.T) {
	now := time.Now().Unix()
	g := fakeGit(
		".agents/changes/drain-six/proposal.md\n.agents/changes/drain/proposal.md\n.agents/changes/add-thing/proposal.md\n",
		map[string]string{
			"drain-six": fmt.Sprintf("%d\n", now),
			"drain":     fmt.Sprintf("%d\n", now),
			"add-thing": fmt.Sprintf("%d\n", now),
		}, nil)

	var out strings.Builder
	if err := metrics(&out, []string{"add"}, g); err != nil {
		t.Fatalf("a unique prefix must select: %v", err)
	}
	rowFor(t, out.String(), "add-thing")

	err := metrics(&strings.Builder{}, []string{"dra"}, g)
	if err == nil {
		t.Fatal("an ambiguous prefix must be refused")
	}
	for _, name := range []string{"drain", "drain-six"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("the refusal must name candidate %q, got: %v", name, err)
		}
	}

	// "drain" is exact and also a prefix of drain-six — exact wins.
	out.Reset()
	if err := metrics(&out, []string{"drain"}, g); err != nil {
		t.Fatalf("an exact name must win over a longer sibling: %v", err)
	}
	rowFor(t, out.String(), "drain")
	if strings.Contains(out.String(), "drain-six") {
		t.Fatalf("wanted only the exact match:\n%s", out.String())
	}
}

func TestHumanSpanDropsPrecisionItDoesNotHave(t *testing.T) {
	// A change spanning four days did not get four days of attention. Minutes on that
	// span is a number implying an accuracy the measurement cannot support.
	for _, tc := range []struct {
		d    time.Duration
		want string
	}{
		{30 * time.Second, "<1m"},
		{45 * time.Minute, "45m"},
		{5 * time.Hour, "5h"},
		{96 * time.Hour, "4d"},
	} {
		if got := humanSpan(tc.d); got != tc.want {
			t.Errorf("humanSpan(%v) = %q, wanted %q", tc.d, got, tc.want)
		}
	}
}

func TestMetricsIsNotGatedOnThePayload(t *testing.T) {
	if needsPayload([]string{"metrics"}) {
		t.Fatal("metrics reads the project's git history, not the payload tree")
	}
}
