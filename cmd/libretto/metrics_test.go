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
	got, err := changeNames(fakeGit(added, nil, nil))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"add-thing", "fix-other"}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("wanted %v, got %v", want, got)
	}
}

func TestChangeNamesIgnoresPathsOutsideTheChangeTree(t *testing.T) {
	got, err := changeNames(fakeGit("docs/FLOW.md\nREADME.md\n.agents/specs/cli/spec.md\n", nil, nil))
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
	m, err := measure(fakeGit("", logs, nil), "c")
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
func TestReopenedBoxesAreCountedSeparatelyFromClosedOnes(t *testing.T) {
	now := time.Now().Unix()
	diff := `+- [x] one
+- [x] two
-- [x] two
+- [ ] two
+- [x] two
`
	m, err := measure(fakeGit("",
		map[string]string{"c": fmt.Sprintf("%d\n", now)},
		map[string]string{"c": diff}), "c")
	if err != nil {
		t.Fatal(err)
	}
	if m.uncheck != 1 {
		t.Fatalf("wanted 1 reopened box, got %d", m.uncheck)
	}
	// three `+[x]` lines, one of them a re-close of the reopened box: two boxes closed.
	if m.checked != 2 {
		t.Fatalf("wanted 2 boxes closed, got %d", m.checked)
	}
}

func TestBoxInReadsOnlyRealCheckboxes(t *testing.T) {
	for _, tc := range []struct{ line, want string }{
		{"+- [x] done", "x"},
		{"+- [X] done", "x"},
		{"+  - [ ] nested open", " "},
		{"+- [] malformed", ""},
		{"+prose about - [x] a box", "x"}, // a diff line is a diff line; this is honest noise
		{"+nothing here", ""},
		{"+- [", ""},
	} {
		if got := boxIn(tc.line); got != tc.want {
			t.Errorf("boxIn(%q) = %q, wanted %q", tc.line, got, tc.want)
		}
	}
}

func TestAChangeWithNoPlanReportsADashNotAZero(t *testing.T) {
	// Zero closed boxes and no plan at all are different facts. Printing 0 for both
	// says a plan existed and nothing got done.
	now := time.Now().Unix()
	var out strings.Builder
	err := metrics(&out, ".", nil, fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "—") {
		t.Fatalf("a change with no plan must not report 0 closed boxes:\n%s", out.String())
	}
}

// The report says what it cannot see. A metrics command that silently omits two of the
// three things asked of it reads as having measured them and found nothing.
func TestTheReportNamesWhatItCannotMeasure(t *testing.T) {
	now := time.Now().Unix()
	var out strings.Builder
	if err := metrics(&out, ".", nil, fakeGit(".agents/changes/c/proposal.md\n",
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
	if err := metrics(&out, ".", []string{"a"}, g); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "a ") || strings.Contains(out.String(), "\n  b ") {
		t.Fatalf("wanted only change a:\n%s", out.String())
	}

	if err := metrics(&strings.Builder{}, ".", []string{"nope"}, g); err == nil {
		t.Fatal("wanted a refusal for a change git never saw")
	}
	if err := metrics(&strings.Builder{}, ".", []string{"--all"}, g); err == nil {
		t.Fatal("wanted a refusal for an unknown flag")
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
