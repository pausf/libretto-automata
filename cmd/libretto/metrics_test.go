package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
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
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil), "")
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
	if err := metrics(&out, nil, g, ""); err != nil {
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
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil), ""); err != nil {
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
	if err := metrics(&out, []string{"a"}, g, ""); err != nil {
		t.Fatal(err)
	}
	rowFor(t, out.String(), "a") // fails the test if absent
	if strings.Contains(out.String(), "\n  b ") {
		t.Fatalf("wanted only change a:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "1 change(s)") {
		t.Fatalf("the footer must count the filtered set, not the whole history:\n%s", out.String())
	}

	if err := metrics(&strings.Builder{}, []string{"nope"}, g, ""); err == nil {
		t.Fatal("wanted a refusal for a change git never saw")
	}
	if err := metrics(&strings.Builder{}, []string{"--all"}, g, ""); err == nil {
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
	if err := metrics(&out, nil, g, ""); err != nil {
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
	if err := metrics(&out, nil, fakeGit(added, logs, nil), ""); err != nil {
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
	if err := metrics(&out, []string{"add"}, g, ""); err != nil {
		t.Fatalf("a unique prefix must select: %v", err)
	}
	rowFor(t, out.String(), "add-thing")

	err := metrics(&strings.Builder{}, []string{"dra"}, g, "")
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
	if err := metrics(&out, []string{"drain"}, g, ""); err != nil {
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

// The report explains its own columns. It already prints what it cannot measure;
// what it does measure deserves no less, or the table needs a translator.
func TestTheReportExplainsItsColumns(t *testing.T) {
	now := time.Now().Unix()
	var out strings.Builder
	if err := metrics(&out, nil, fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil), ""); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"calendar clock",               // span
		"called done before they were", // reopen
		"boxes the plan holds",         // closed's denominator
		"no plan.md in its history",    // the — cell
		"still on disk",                // state
		"unreadable",                   // the state the error row prints
		"touched the change",           // commits
		"lessons.md",                   // corr
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("the legend must say %q:\n%s", want, out.String())
		}
	}
}

// gitAtRoot points --show-toplevel at a real directory so the ledger tests can put a
// real lessons.md under it; everything else stays the fake's answer.
func gitAtRoot(root string, g gitRunner) gitRunner {
	return func(args ...string) (string, error) {
		if strings.Contains(strings.Join(args, " "), "--show-toplevel") {
			return root + "\n", nil
		}
		return g(args...)
	}
}

func ledgerAt(t *testing.T, content string) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(path.Join(root, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(root, ".agents/lessons.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// corrField returns the corrections cell of one change's row: the field before the
// state, so the assertion cannot be satisfied by the commit count or the span.
func corrField(t *testing.T, out, change string) string {
	t.Helper()
	f := strings.Fields(rowFor(t, out, change))
	if len(f) < 2 {
		t.Fatalf("row for %q has no fields:\n%s", change, out)
	}
	return f[len(f)-2]
}

func TestCorrectionsAreCountedPerChange(t *testing.T) {
	root := ledgerAt(t, `## 2026-08-13 · a · build-and-check
Said: the tag here is scoped, never bare
Did: wrote a bare tag

## 2026-08-13 · a · record-work
Said: no AI attribution, ever
Did: added a Co-Authored-By trailer

## 2026-08-14 · b · write-spec
Said: that criterion cannot fail
Did: wrote an unfalsifiable criterion
`)
	now := time.Now().Unix()
	g := gitAtRoot(root, fakeGit(".agents/changes/a/proposal.md\n.agents/changes/b/proposal.md\n",
		map[string]string{
			"a": fmt.Sprintf("%d\n", now),
			"b": fmt.Sprintf("%d\n", now),
		}, nil))

	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if got := corrField(t, out.String(), "a"); got != "2" {
		t.Fatalf("wanted 2 corrections for a, got %q in:\n%s", got, out.String())
	}
	if got := corrField(t, out.String(), "b"); got != "1" {
		t.Fatalf("wanted 1 correction for b, got %q in:\n%s", got, out.String())
	}
}

func TestCorrectionsCountsByPhase(t *testing.T) {
	// The breakdown answers *where* corrections surface, spelled exactly as the
	// ledger spells them — a normaliser would be a second spelling that drifts.
	root := ledgerAt(t, `## 2026-08-13 · a · phase 2
Said: one
## 2026-08-13 · a · phase 2
Said: two
## 2026-08-14 · b · 6→7
Said: a reviewer finding
## 2026-08-14 · - · phase 8
Said: an orphan still surfaced somewhere
`)
	counts, byPhase, orphans, seen := corrections(root)
	if !seen {
		t.Fatal("the ledger exists and must be seen")
	}
	if byPhase["phase 2"] != 2 || byPhase["6→7"] != 1 || byPhase["phase 8"] != 1 {
		t.Fatalf("wanted phase 2:2, 6→7:1, phase 8:1, got %v", byPhase)
	}
	if orphans != 1 {
		t.Fatalf("wanted 1 orphan, got %d", orphans)
	}
	if counts["b"] != 0 {
		t.Fatalf("a 6→7 entry is a reviewer finding, not a user correction: got %d for b", counts["b"])
	}
}

func TestOrphanFindingIsNotAnOrphanCorrection(t *testing.T) {
	// A 6→7 entry with no change open is still a finding: it counts by phase and
	// never in the "correction(s) outside any change" line, whose word is correction.
	root := ledgerAt(t, `## 2026-08-14 · - · 6→7
Said: a finding with no change open
`)
	_, byPhase, orphans, seen := corrections(root)
	if !seen {
		t.Fatal("the ledger exists and must be seen")
	}
	if orphans != 0 {
		t.Fatalf("a 6→7 orphan is a finding, not a correction: got %d orphans", orphans)
	}
	if byPhase["6→7"] != 1 {
		t.Fatalf("the finding must still count by phase, got %v", byPhase)
	}
}

func TestReviewerFindingsStayOutOfCorrections(t *testing.T) {
	// The per-change corr column keeps meaning user corrections only; the seam's
	// findings arrive under 6→7 and count in the breakdown, never in the column.
	root := ledgerAt(t, `## 2026-08-13 · a · build-and-check
Said: a user correction
## 2026-08-13 · a · 6→7
Said: a reviewer finding
`)
	now := time.Now().Unix()
	g := gitAtRoot(root, fakeGit(".agents/changes/a/proposal.md\n",
		map[string]string{"a": fmt.Sprintf("%d\n", now)}, nil))

	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if got := corrField(t, out.String(), "a"); got != "1" {
		t.Fatalf("wanted 1 user correction for a, got %q in:\n%s", got, out.String())
	}
	if !strings.Contains(out.String(), "6→7") {
		t.Fatalf("the finding must still appear in the phase breakdown:\n%s", out.String())
	}
}

func TestMetricsReportsCorrectionsByPhase(t *testing.T) {
	root := ledgerAt(t, `## 2026-08-13 · a · phase 2
Said: early
## 2026-08-14 · a · phase 8
Said: late
`)
	now := time.Now().Unix()
	g := gitAtRoot(root, fakeGit(".agents/changes/a/proposal.md\n",
		map[string]string{"a": fmt.Sprintf("%d\n", now)}, nil))

	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"corrections by phase", "phase 2", "phase 8", "where the flow got corrected"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("the breakdown must say %q:\n%s", want, out.String())
		}
	}
}

func TestMetricsPhaseBreakdownAbsentLedger(t *testing.T) {
	// Absent, empty and populated are three facts. Absent means capture is not in
	// use; present-but-empty means it is in use and nothing was captured yet.
	now := time.Now().Unix()
	g := gitAtRoot(t.TempDir(), fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil))
	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "corrections not captured") {
		t.Fatalf("an absent ledger must read as capture-not-in-use:\n%s", out.String())
	}
	if strings.Contains(out.String(), "corrections by phase") {
		t.Fatalf("an absent ledger must render no phase rows:\n%s", out.String())
	}

	root := ledgerAt(t, "")
	g = gitAtRoot(root, fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil))
	out.Reset()
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "no corrections captured yet") {
		t.Fatalf("a present-but-empty ledger must say so in one line:\n%s", out.String())
	}
	if strings.Contains(out.String(), "corrections by phase") {
		t.Fatalf("an empty ledger must render no phase rows:\n%s", out.String())
	}
}

func TestNoLedgerReportsADashNotAZero(t *testing.T) {
	// No ledger and a ledger with nothing for this change are different facts: the
	// first means capture is not in use, the second means the flow ran uncorrected.
	now := time.Now().Unix()
	g := gitAtRoot(t.TempDir(), fakeGit(".agents/changes/c/proposal.md\n",
		map[string]string{"c": fmt.Sprintf("%d\n", now)}, nil))

	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if got := corrField(t, out.String(), "c"); got != "—" {
		t.Fatalf("no ledger must read as a dash, got %q in:\n%s", got, out.String())
	}
}

func TestMalformedAndChangelessEntriesDoNotCrashTheCount(t *testing.T) {
	// One separator is not the contract's header; it is skipped, never a crash. An
	// entry whose change field is `-` belongs to no row, and the report names it in
	// one line so it is not silently lost.
	root := ledgerAt(t, `## 2026-08-13 · build-and-check
Said: malformed, one separator only
Did: this line must be skipped

## 2026-08-13 · a · build-and-check
Said: a real one
Did: something wrong

## 2026-08-14 · - · find-work
Said: corrected outside any change
Did: something wrong with no change open
`)
	now := time.Now().Unix()
	g := gitAtRoot(root, fakeGit(".agents/changes/a/proposal.md\n",
		map[string]string{"a": fmt.Sprintf("%d\n", now)}, nil))

	var out strings.Builder
	if err := metrics(&out, nil, g, ""); err != nil {
		t.Fatal(err)
	}
	if got := corrField(t, out.String(), "a"); got != "1" {
		t.Fatalf("wanted the malformed header skipped and 1 counted, got %q in:\n%s", got, out.String())
	}
	if !strings.Contains(out.String(), "1 correction(s) outside any change") {
		t.Fatalf("the report must name the changeless entries:\n%s", out.String())
	}
}

func TestMetricsIsNotGatedOnThePayload(t *testing.T) {
	if needsPayload([]string{"metrics"}) {
		t.Fatal("metrics reads the project's git history, not the payload tree")
	}
}

// tokenLine returns the token block row starting with label, or "" when absent.
//
// Scoped to everything after the token block's header on purpose. Searching the whole
// output matched the change table's own row for a change of the same name, and the
// dash test passed against the churn column while the token block did not exist —
// a false green found by reading the failure list rather than the pass list.
func tokenLine(out, label string) string {
	_, block, ok := strings.Cut(out, "\n  tokens ")
	if !ok {
		return ""
	}
	for _, l := range strings.Split(block, "\n") {
		if t := strings.TrimSpace(l); strings.HasPrefix(t, label+" ") || t == label {
			return t
		}
	}
	return ""
}

func twoChangeGit(now int64) gitRunner {
	return fakeGit(".agents/changes/alpha/proposal.md\n.agents/changes/beta/proposal.md\n",
		map[string]string{
			"alpha": fmt.Sprintf("%d\n", now),
			"beta":  fmt.Sprintf("%d\n", now),
		}, nil)
}

func TestTheTokenFooterIsCorpusWideUnderAFilter(t *testing.T) {
	projects := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("feat/alpha", "write-spec", 1, 10, 100, 1000),
			assistantLine("feat/beta", "write-spec", 2, 20, 200, 2000),
			assistantLine("main", "write-spec", 4, 40, 400, 4000),
		},
	})
	now := time.Now().Unix()

	var all, one strings.Builder
	if err := metrics(&all, nil, twoChangeGit(now), projects); err != nil {
		t.Fatal(err)
	}
	if err := metrics(&one, []string{"alpha"}, twoChangeGit(now), projects); err != nil {
		t.Fatal(err)
	}

	// The invariant attributed + unattributed = corpus has to be readable off the
	// output. It cannot be if the totals move when a filter is applied.
	for _, label := range []string{"attributed", "unattributed"} {
		a, o := tokenLine(all.String(), label), tokenLine(one.String(), label)
		if a == "" {
			t.Fatalf("no %q row in the unfiltered report:\n%s", label, all.String())
		}
		if a != o {
			t.Errorf("%q moved under a filter:\n  all: %s\n  one: %s", label, a, o)
		}
	}
}

func TestAChangeWithNoTokensReportsADashNotAZero(t *testing.T) {
	// A transcript root exists and was read; nothing in it names beta. That is "could
	// not attribute", not "ran and cost nothing", and the two must not print alike.
	projects := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {assistantLine("feat/alpha", "write-spec", 1, 10, 100, 1000)},
	})
	var out strings.Builder
	if err := metrics(&out, []string{"beta"}, twoChangeGit(time.Now().Unix()), projects); err != nil {
		t.Fatal(err)
	}

	row := tokenLine(out.String(), "beta")
	if row == "" {
		t.Fatalf("no token row for beta:\n%s", out.String())
	}
	if strings.Contains(row, "0") {
		t.Errorf("beta has no attributable session and must dash, not zero: %q", row)
	}
	if !strings.Contains(row, "—") {
		t.Errorf("want a dash in beta's token row, got %q", row)
	}
}

func TestPerPhaseCostCarriesAnUnattributedRow(t *testing.T) {
	projects := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {
			assistantLine("feat/alpha", "write-spec", 1, 10, 100, 1000),
			assistantLine("feat/alpha", "build-and-check", 2, 20, 200, 2000),
			assistantLine("feat/alpha", "", 4, 40, 400, 4000), // named no skill
		},
	})
	var out strings.Builder
	if err := metrics(&out, []string{"alpha"}, twoChangeGit(time.Now().Unix()), projects); err != nil {
		t.Fatal(err)
	}
	s := out.String()

	for _, want := range []string{"by phase", "write-spec", "build-and-check"} {
		if !strings.Contains(s, want) {
			t.Errorf("the phase block must name %q:\n%s", want, s)
		}
	}
	// 53% of entries carry no attributionSkill. Sharing them out across the phases that
	// did would invent a number; the row is what keeps the arithmetic honest.
	phases := 0
	for _, l := range strings.Split(s, "\n") {
		f := strings.Fields(l)
		if len(f) == 5 && (f[0] == "write-spec" || f[0] == "build-and-check") {
			phases++
			if f[1] == "3" || f[1] == "7" {
				t.Errorf("the unattributed entry was distributed into a phase row: %q", l)
			}
		}
	}
	if phases != 2 {
		t.Errorf("got %d phase rows, want 2", phases)
	}
	// The corpus block has an unattributed row of its own, so this one is read from the
	// phase block. Two rows legitimately carry the label; each means it in its own block.
	_, phaseBlock, _ := strings.Cut(s, "\n  by phase ")
	if u := tokenLine("\n  tokens "+phaseBlock, "unattributed"); !strings.Contains(u, "4") {
		t.Errorf("the skill-less entry must reach the phase block's unattributed row, got %q", u)
	}
}

func TestNoTranscriptRootStillReportsTheGitMetrics(t *testing.T) {
	now := time.Now().Unix()
	for _, projects := range []string{"", filepath.Join(t.TempDir(), "gone")} {
		var out strings.Builder
		if err := metrics(&out, nil, twoChangeGit(now), projects); err != nil {
			t.Fatalf("projects=%q: %v", projects, err)
		}
		s := out.String()
		// The git-derived report is untouched by a measurement that was unavailable.
		if !strings.Contains(s, "2 change(s)") {
			t.Errorf("projects=%q: the git metrics stopped printing:\n%s", projects, s)
		}
		if !strings.Contains(s, "no session transcripts") {
			t.Errorf("projects=%q: want the token block replaced by a line saying so:\n%s", projects, s)
		}
		if tokenLine(s, "attributed") != "" {
			t.Errorf("projects=%q: printed token totals with no transcripts to read", projects)
		}
	}
}

func TestTheCeilingSeparatesCostFromDuration(t *testing.T) {
	var out strings.Builder
	if err := metrics(&out, nil, twoChangeGit(time.Now().Unix()), ""); err != nil {
		t.Fatal(err)
	}
	s := out.String()

	// Still true, and still said.
	if !strings.Contains(s, "per-phase duration") {
		t.Errorf("duration is still unmeasurable and must still be named:\n%s", s)
	}
	// Nothing here was ever false — the ceiling named duration and review-work findings,
	// and both are still unmeasurable. What it lacked is any statement about cost, which
	// this change starts measuring. Both assertions below were red before it.
	if !strings.Contains(s, "attributionSkill") {
		t.Errorf("the ceiling must say what made per-phase cost measurable:\n%s", s)
	}
	if !strings.Contains(s, "unattributed row is the measurement's own") {
		t.Errorf("the ceiling must name the token block's own limit:\n%s", s)
	}
}

// zeroLine is an attributed entry that genuinely cost nothing — the <synthetic> model
// shape, which the constraints name and which really occurs.
func zeroLine(branch, skill string) string {
	return fmt.Sprintf(`{"type":"assistant","gitBranch":%q,"attributionSkill":%q,`+
		`"message":{"model":"<synthetic>","usage":{"input_tokens":0,"output_tokens":0,`+
		`"cache_creation_input_tokens":0,"cache_read_input_tokens":0,`+
		`"service_tier":null,"iterations":null}}}`, branch, skill)
}

func TestAnAttributedChangeWithZeroTokensIsNotADash(t *testing.T) {
	// The third state of outcome 3, and it is reachable rather than theoretical: a
	// <synthetic> entry attributes to the change and costs nothing. Sessions reached it,
	// so it must not render like a change no session reached.
	projects := transcriptRoot(t, "/repo", map[string][]string{
		"a.jsonl": {zeroLine("feat/alpha", "build-and-check")},
	})
	var out strings.Builder
	if err := metrics(&out, []string{"alpha"}, twoChangeGit(time.Now().Unix()), projects); err != nil {
		t.Fatal(err)
	}
	s := out.String()

	row := tokenLine(s, "alpha")
	if row == "" {
		t.Fatalf("no token row for alpha:\n%s", s)
	}
	if strings.Contains(row, "—") {
		t.Errorf("alpha was attributed and must print zeros, not a dash: %q", row)
	}
	// And the phase surface survives: the entry named a skill, so suppressing the block
	// on a zero total loses a phase that was recorded.
	if !strings.Contains(s, "by phase") || !strings.Contains(s, "build-and-check") {
		t.Errorf("the phase block vanished for a change whose entry named a phase:\n%s", s)
	}
}

func TestASmallMissRateDoesNotRoundAwayToZero(t *testing.T) {
	// The ceiling calls this number the measurement's own error bar. An error bar that
	// truncates to 0 while entries are genuinely unattributed reports the opposite of
	// what it is for.
	lines := []string{assistantLine("main", "write-spec", 1, 1, 1, 1)}
	for i := 0; i < 200; i++ {
		lines = append(lines, assistantLine("feat/alpha", "write-spec", 1, 1, 1, 1))
	}
	projects := transcriptRoot(t, "/repo", map[string][]string{"a.jsonl": lines})

	var out strings.Builder
	if err := metrics(&out, nil, twoChangeGit(time.Now().Unix()), projects); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if strings.Contains(s, "0% of session entries") {
		t.Errorf("1 unattributed entry in 201 reported as 0%%:\n%s", s)
	}
	if !strings.Contains(s, "<1% of session entries") {
		t.Errorf("want a sub-1%% miss rate said as <1%%:\n%s", s)
	}
}
