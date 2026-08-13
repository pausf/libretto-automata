package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

// Flow metrics, derived and never recorded.
//
// The obvious build is instrumentation: each phase writes a timestamp somewhere, and
// the report reads them back. That buys precision and costs a new artifact in every
// change folder, written by eight skills that would each have to remember — and a
// metric only collected when somebody remembers is a metric with holes exactly where
// the interesting runs are.
//
// Git already holds most of the answer. A change folder's first commit and its last
// are the wall clock; the commits between are the iterations; plan.md's history is the
// churn. Nothing new is written, nothing can be forgotten, and it works retroactively
// on every change this repository ever landed.
//
// What that cannot see is named in the report rather than guessed at. See flowCeiling.

const flowLegend = `  Columns:
    commits  commits that touched the change's folder
    span     first commit to last — calendar clock, not attention
    closed   boxes now closed, over the boxes the plan holds
    reopen   boxes that went back open — tasks called done before they were
    corr     user corrections recorded in .agents/lessons.md — a dash when no ledger exists
    state    landed = the folder was deleted on landing · in flight = still on disk ·
             unreadable = git would not answer for it, so its row is all ?
    —        the change has no plan.md in its history`

const flowCeiling = `  Not measured, and not derivable from git:
    · per-phase duration — phases 1 to 7 happen inside one session and leave one commit
    · review-work findings — reported in a session and repaired before anything lands,
      so the repair is in the diff and the count never was
  Both need a phase to write them down, and neither is worth an artifact eight skills
  have to remember until somebody wants the number badly enough to say what for.

  Per-phase *cost* was never on that list and is now measured: the transcripts record
  attributionSkill, so the phase is written down already and this is a free rider on it.
  Duration stays off — the same entries carry timestamps, but a phase's wall clock
  includes every wait for a human and would report attention it never had.

  What the token block cannot do is attribute every entry. A change is recognised by its
  branch name, so work done on main, in a detached HEAD, or on a branch named unlike its
  change lands in the unattributed row. That unattributed row is the measurement's own
  error bar — read it before trusting the rest.`

type changeMetrics struct {
	name     string
	first    time.Time
	last     time.Time
	commits  int
	checked  int // boxes that went from open to closed
	uncheck  int // boxes that went back — the churn worth seeing
	boxes    int // boxes the plan holds now, the closed cell's denominator
	landed   bool
	planSeen bool
}

func (m changeMetrics) span() time.Duration { return m.last.Sub(m.first) }

// gitRunner is exec'd git, injected so the tests drive real parsing against fixed
// output instead of building repositories.
type gitRunner func(args ...string) (string, error)

func execGit(dir string) gitRunner {
	return func(args ...string) (string, error) {
		c := exec.Command("git", args...)
		c.Dir = dir
		out, err := c.Output()
		return string(out), err
	}
}

// changeNames lists every change folder git has ever seen, landed or in flight. `git
// log --diff-filter=A` on the whole directory is what makes a landed change visible:
// the folder is deleted when it lands, so a filesystem listing sees only what is open,
// which is the half already visible in `libretto status`.
//
// --full-history on every query here, and it is not optional. Git's default history
// simplification prunes commits that do not change a path's *final* state — and a landed
// change's final state is "does not exist", so the whole history of every change this
// report exists to measure is simplified away. Measured, not assumed: without it, two of
// twelve changes returned no history at all and no plan.md diff was found for any of
// them, so every churn column read as "no plan" on changes that plainly had one.
func changeNames(git gitRunner, root string) ([]string, error) {
	out, err := git("log", "--full-history", "--diff-filter=A", "--name-only", "--format=",
		"--", path.Join(root, ".agents/changes")+"/")
	if err != nil {
		return nil, err
	}
	// git reports paths relative to the repository root whatever the pathspec was, so the
	// prefix to strip is the repo-relative one and never the absolute one just passed in.
	seen := map[string]bool{}
	for _, l := range strings.Split(out, "\n") {
		l = strings.TrimSpace(l)
		if !strings.HasPrefix(l, ".agents/changes/") {
			continue
		}
		rest := strings.TrimPrefix(l, ".agents/changes/")
		n, _, ok := strings.Cut(rest, "/")
		if !ok || n == "" {
			continue
		}
		seen[n] = true
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names, nil
}

// measure reads one change's whole history. Commits arrive newest-first from git, so
// last is the first line seen and first is the last.
func measure(git gitRunner, root, name string) (changeMetrics, error) {
	m := changeMetrics{name: name}
	dir := path.Join(root, ".agents/changes", name)

	out, err := git("log", "--full-history", "--format=%ct", "--", dir)
	if err != nil {
		return m, err
	}
	var stamps []time.Time
	for _, l := range strings.Fields(out) {
		var sec int64
		if _, err := fmt.Sscanf(l, "%d", &sec); err != nil {
			continue
		}
		stamps = append(stamps, time.Unix(sec, 0))
	}
	if len(stamps) == 0 {
		return m, fmt.Errorf("no history for %s", name)
	}
	m.commits = len(stamps)
	m.last, m.first = stamps[0], stamps[len(stamps)-1]

	// Landed means the folder was deleted — the delta applied onto the capability spec
	// and the change removed, per the flow. A folder still on disk is work in flight.
	if _, err := os.Stat(dir); err != nil {
		m.landed = true
	}

	// Checkbox churn, straight out of plan.md's diffs. A `+- [x]` is a box closed in
	// that commit; a `-- [x]` is one that stopped being closed, which is a task called
	// done before it was. Counting reopenings separately is the whole point — a net
	// count hides exactly the signal worth having.
	//
	// --diff-filter=AM excludes the deletion, and without it the numbers invert. A change
	// lands by deleting its folder, and that commit's diff removes every line in plan.md
	// — so every box ever closed also appears as a `-[x]`, every change reads as
	// completely reopened and nothing reads as closed. Measured here on real history:
	// 52 reopenings and 0 closes on a change that had neither.
	plan := path.Join(dir, "plan.md")
	diff, err := git("log", "--full-history", "--diff-filter=AM", "-p", "--format=%H", "--", plan)
	if err == nil && strings.TrimSpace(diff) != "" {
		m.planSeen = true
		m.checked, m.uncheck, m.boxes = churn(diff)
	}
	return m, nil
}

// churn reads `git log -p --format=%H` over one plan.md and returns boxes closed and
// boxes reopened.
//
// **Per commit, not per line, and that is the whole correctness argument.** Counting
// added and removed `[x]` lines across the whole history makes a reworded task look like
// a reopening: rewording a closed box emits `-[x] one` and `+[x] one, reworded` in the
// same commit, which line-counting reports as one close and one reopen. So does deleting
// a finished task. Both were reported by review as false accusations, because a
// reopening means "a task was called done before it was" and a typo fix is not that.
//
// Within one commit the net change in closed boxes is what actually happened: +1 means a
// box closed, -1 means one stopped being closed, 0 means the text moved and the state did
// not. A reword nets zero and vanishes, which is correct.
// It also returns how many boxes the plan holds at the end of its history — added box
// lines minus removed ones, whatever their state. A reword nets zero here too, and a
// deleted task shrinks the denominator along with whatever it did to the numerator.
func churn(diff string) (closed, reopened, boxes int) {
	flush := func(net int) {
		if net > 0 {
			closed += net
		} else if net < 0 {
			reopened += -net
		}
	}
	net := 0
	for _, l := range strings.Split(diff, "\n") {
		if isCommitSHA(l) {
			flush(net)
			net = 0
			continue
		}
		// A diff header line begins with +++ or ---; it is not content.
		if strings.HasPrefix(l, "+++") || strings.HasPrefix(l, "---") {
			continue
		}
		switch {
		case strings.HasPrefix(l, "+") && boxIn(l) != "":
			boxes++
			if boxIn(l) == "x" {
				net++
			}
		case strings.HasPrefix(l, "-") && boxIn(l) != "":
			boxes--
			if boxIn(l) == "x" {
				net--
			}
		}
	}
	flush(net)
	return closed, reopened, boxes
}

func isCommitSHA(l string) bool {
	if len(l) != 40 {
		return false
	}
	for _, c := range l {
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}

// boxIn returns "x", " " or "" for a diff line carrying a task checkbox.
//
// It strips the diff marker and then asks `loop`'s parser, so the two commands cannot
// disagree about what a box is. They did: this counted a box anywhere in a line while
// planLine anchors to the start, so the same plan.md yielded two different totals
// depending on which command read it.
func boxIn(l string) string {
	if l == "" {
		return ""
	}
	m := planLine.FindStringSubmatch(l[1:]) // drop the +/- diff marker
	if m == nil {
		return ""
	}
	if m[1] == " " {
		return " "
	}
	return "x"
}

// corrections reads the lessons ledger — the one artifact the payload writes and this
// command only counts. `evidence` appends an entry per user correction; the header is
// the contract: `## <date> · <change> · <phase>`, exactly two ` · ` separators, three
// non-empty fields, date not validated. Anything else is skipped, never a crash — the
// ledger is written by prompts, and a parser that dies on one bad line loses the report.
//
// Returns per-change counts, the entries whose change field is `-` (a correction with
// no change open — counted nowhere, named in the report), and whether the file exists
// at all. Absent and empty are different facts: absent means capture is not in use,
// and printing 0 for that would read as a flow that was never corrected.
func corrections(root string) (counts map[string]int, orphans int, seen bool) {
	b, err := os.ReadFile(path.Join(root, ".agents/lessons.md"))
	if err != nil {
		return nil, 0, false
	}
	counts = map[string]int{}
	for _, l := range strings.Split(string(b), "\n") {
		rest, ok := strings.CutPrefix(l, "## ")
		if !ok {
			continue
		}
		f := strings.Split(rest, " · ")
		if len(f) != 3 {
			continue
		}
		date, change, phase := strings.TrimSpace(f[0]), strings.TrimSpace(f[1]), strings.TrimSpace(f[2])
		if date == "" || change == "" || phase == "" {
			continue
		}
		if change == "-" {
			orphans++
			continue
		}
		counts[change]++
	}
	return counts, orphans, true
}

// The repository root comes from git rather than a parameter: the caller's cwd is where
// git was pointed, and `--show-toplevel` is the only thing that knows how far up the root
// is. A projectDir parameter here would be a second answer that agrees only by accident.
func metrics(w io.Writer, args []string, git gitRunner, projects string) error {
	var only string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			return fmt.Errorf("unknown flag %q — `%s metrics [<change>]`", a, invokedAs())
		}
		if only != "" {
			return fmt.Errorf("one change at a time, got %q and %q", only, a)
		}
		only = a
	}

	// Everything below is asked of the repository root, never of the cwd. git pathspecs
	// and os.Stat are both cwd-relative, so run from a subdirectory this reported "no
	// changes in this repository's history yet" — a plausible answer that was false.
	root, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	root = strings.TrimSpace(root)

	names, err := changeNames(git, root)
	if err != nil {
		return fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	// Kept before the filter narrows it. The token block attributes against every change
	// git has seen, so its totals stay corpus-wide however the table is scoped.
	allNames := names
	if only != "" {
		matches := filterName(names, only)
		if len(matches) == 0 {
			return fmt.Errorf("git has never seen a change called %q", only)
		}
		if len(matches) > 1 {
			return fmt.Errorf("%q is ambiguous: %s", only, strings.Join(matches, ", "))
		}
		names = matches
		// A prefix resolves to one name here, so the token block reports what the table
		// reports rather than the letters that were typed.
		only = matches[0]
	}
	if len(names) == 0 {
		fmt.Fprintf(w, "\n  no changes in this repository's history yet\n\n")
		return nil
	}

	corrCounts, orphans, ledgerSeen := corrections(root)

	fmt.Fprintf(w, "\n  %-34s %7s %8s %7s %7s %7s  %s\n", "change", "commits", "span", "closed", "reopen", "corr", "state")
	var totalCommits, totalReopen int
	var spans []changeMetrics
	for _, n := range names {
		m, err := measure(git, root, n)
		if err != nil {
			// Never skip silently. The footer counts names, so a skipped row makes the
			// total disagree with what is above it — and a report that says twelve while
			// showing ten is worse than one that admits it could not read two.
			fmt.Fprintf(w, "  %-34s %7s %8s %7s %7s %7s  unreadable\n", trunc(n, 34), "?", "?", "?", "?", corrCell(ledgerSeen, corrCounts[n]))
			continue
		}
		state := "in flight"
		if m.landed {
			state = "landed"
		}
		churn := fmt.Sprint(m.uncheck)
		if !m.planSeen {
			churn = "—"
		}
		fmt.Fprintf(w, "  %-34s %7d %8s %7s %7s %7s  %s\n",
			trunc(m.name, 34), m.commits, humanSpan(m.span()),
			planCell(m.planSeen, m.checked-m.uncheck, m.boxes), churn, corrCell(ledgerSeen, corrCounts[n]), state)
		totalCommits += m.commits
		totalReopen += m.uncheck
		spans = append(spans, m)
	}
	fmt.Fprintf(w, "\n  %d change(s), %d commit(s), %s of wall clock, %d box(es) reopened\n",
		len(names), totalCommits, humanSpan(mergedSpan(spans)), totalReopen)
	if orphans > 0 {
		fmt.Fprintf(w, "  %d correction(s) outside any change\n", orphans)
	}
	tokenBlock(w, projects, root, allNames, only)
	fmt.Fprintf(w, "\n%s\n\n%s\n\n", flowLegend, flowCeiling)
	return nil
}

// tokenBlock prints what the session transcripts say this repository's work cost.
//
// allNames rather than the filtered set, always: the corpus totals must not move when a
// change filter is applied, or "attributed + unattributed = corpus" stops being readable
// off the output and the unattributed row stops meaning anything.
func tokenBlock(w io.Writer, projects, root string, allNames []string, only string) {
	if projects == "" {
		fmt.Fprintf(w, "\n  no session transcripts configured — token cost not measured\n")
		return
	}
	es, found := readUsage(projects, root)
	if !found {
		fmt.Fprintf(w, "\n  no session transcripts for this repository — token cost not measured\n")
		return
	}

	byChange, unattributed := attribute(es, allNames)
	var attributed usageTotals
	for _, t := range byChange {
		attributed = attributed.plus(t)
	}

	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "tokens", "input", "output", "cache-w", "cache-r")
	tokenRow(w, "attributed", attributed, true)
	tokenRow(w, "unattributed", unattributed, true)
	if n := len(es); n > 0 {
		var off int
		for _, e := range es {
			if attributeBranch(e.branch, allNames) == "" {
				off++
			}
		}
		// Integer division truncates, so one miss in two hundred reads as 0% — the
		// opposite of what a number the ceiling calls an error bar is for. Anything
		// under one percent says so rather than rounding itself away.
		pct := fmt.Sprintf("%d%%", off*100/n)
		if off > 0 && off*100/n == 0 {
			pct = "<1%"
		}
		fmt.Fprintf(w, "\n  %s of session entries could not be attributed to a change\n", pct)
	}

	if only == "" {
		return
	}

	// Presence in the map, never a zero total. A change can be attributed and cost
	// nothing — a <synthetic> entry does exactly that — and rendering it like a change
	// no session reached would collapse the two states outcome 3 exists to separate.
	t, attributedHere := byChange[only]
	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "by change", "input", "output", "cache-w", "cache-r")
	tokenRow(w, trunc(only, 16), t, attributedHere)

	if !attributedHere {
		return
	}
	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "by phase", "input", "output", "cache-w", "cache-r")
	byPhase := map[string]usageTotals{}
	var noPhase usageTotals
	var noPhaseSeen bool
	for _, e := range es {
		if attributeBranch(e.branch, allNames) != only {
			continue
		}
		t := usageTotals{e.in, e.out, e.cacheW, e.cacheR}
		if e.skill == "" {
			noPhase, noPhaseSeen = noPhase.plus(t), true
			continue
		}
		byPhase[e.skill] = byPhase[e.skill].plus(t)
	}
	phases := make([]string, 0, len(byPhase))
	for p := range byPhase {
		phases = append(phases, p)
	}
	sort.Strings(phases)
	for _, p := range phases {
		tokenRow(w, trunc(p, 16), byPhase[p], true)
	}
	// Never distributed across the phases that did name a skill: 4,216 of 8,944 entries
	// carry no attributionSkill, and sharing them out would invent the number.
	tokenRow(w, "unattributed", noPhase, noPhaseSeen)
}

func tokenRow(w io.Writer, label string, t usageTotals, real bool) {
	if !real {
		fmt.Fprintf(w, "  %-16s %14s %14s %14s %16s\n", label, "—", "—", "—", "—")
		return
	}
	fmt.Fprintf(w, "  %-16s %14d %14d %14d %16d\n", label, t.in, t.out, t.cacheW, t.cacheR)
}

// mergedSpan sums [first, last] ranges counting each calendar hour once. The footer
// says "wall clock", and two changes open the same week are one week of clock — the
// plain sum reported both, which was nobody's calendar.
func mergedSpan(ms []changeMetrics) time.Duration {
	sort.Slice(ms, func(i, j int) bool { return ms[i].first.Before(ms[j].first) })
	var total time.Duration
	var end time.Time
	for _, m := range ms {
		start := m.first
		if start.Before(end) {
			start = end
		}
		if m.last.After(end) {
			total += m.last.Sub(start)
			end = m.last
		}
	}
	return total
}

// planCell reads n/total: boxes closed now over boxes the plan has. The cumulative
// close count belongs to the reopen column's story, not this cell's.
func planCell(seen bool, n, total int) string {
	if !seen {
		return "—"
	}
	return fmt.Sprintf("%d/%d", n, total)
}

// corrCell prints the corrections column: a dash when no ledger exists — absent is
// not zero — and the change's count when one does, zero included.
func corrCell(seen bool, n int) string {
	if !seen {
		return "—"
	}
	return fmt.Sprint(n)
}

// filterName resolves a typed name: exact first, then prefix. Exact must win even when
// it also prefixes a sibling, or a change whose full name starts another's becomes
// unreachable. A multi-element return is an ambiguity for the caller to refuse.
func filterName(names []string, only string) []string {
	var prefixed []string
	for _, n := range names {
		if n == only {
			return []string{n}
		}
		if strings.HasPrefix(n, only) {
			prefixed = append(prefixed, n)
		}
	}
	return prefixed
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// humanSpan reports what the number is worth. A change that took four days did not get
// four days of attention, so minutes on a multi-day span is precision the number does
// not have.
func humanSpan(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "<1m"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 48*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
