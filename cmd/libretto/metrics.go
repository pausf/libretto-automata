package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"
)

// Flow metrics, derived from git — the loop's run log is the one thing recorded, by the binary.
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
    —        the change has no checklist in its history: proposal.md, tasks.md, plan.md`

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
	//
	// Three paths, tried in order, because the checklist moved twice and this reads history
	// rather than disk: `proposal.md`'s `## Checklist` is where it lives now, `tasks.md` is
	// where it lived after the first rename, and `plan.md` is what every change that landed
	// before that carries — and their folders are gone, so they will never move. One
	// pathspec covering all three was the shorter version and it is wrong: the other two
	// are prose documents allowed to show a checkbox, and counting their diffs would
	// inflate the churn of every change.
	//
	// A query is found only when it yields a box in scope, never merely a diff: every
	// change's `proposal.md` has history, because the capture adds it, and a non-empty
	// diff as the signal would blank the churn of every change that kept its boxes
	// elsewhere. The change document is read with whole-file context, because a hunk does
	// not carry its section heading.
	//
	// The `tasks.md` and `plan.md` queries are permanent: history is never rewritten, so
	// the changes that carry them will always be there to measure.
	for _, q := range []struct {
		file      string
		checklist bool
	}{{"proposal.md", true}, {"tasks.md", false}, {"plan.md", false}} {
		args := []string{"log", "--full-history", "--diff-filter=AM", "-p", "--format=%H"}
		if q.checklist {
			args = append(args, "--unified=1000000")
		}
		diff, err := git(append(args, "--", path.Join(dir, q.file))...)
		if err != nil {
			continue
		}
		if closed, reopened, boxes, seen := churnIn(diff, q.checklist); seen {
			m.planSeen = true
			m.checked, m.uncheck, m.boxes = closed, reopened, boxes
			break
		}
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
	closed, reopened, boxes, _ = churnIn(diff, false)
	return closed, reopened, boxes
}

// churnIn is churn, optionally restricted to a change document's `## Checklist`, and it
// also reports whether any box line was in scope — the only thing that makes a query
// found.
//
// Restricted, the section is tracked twice per commit, because a diff is two files at
// once: the post-image is the context and `+` lines, the pre-image the context and `-`
// lines. A `+` box counts when the post-image is in the checklist, a `-` box when the
// pre-image was. Named limit: a commit that renames a heading so untouched boxes move in
// or out changes no box line, so it moves neither the churn nor the denominator.
func churnIn(diff string, checklist bool) (closed, reopened, boxes int, seen bool) {
	flush := func(net int) {
		if net > 0 {
			closed += net
		} else if net < 0 {
			reopened += -net
		}
	}
	net := 0
	var pre, post checklistScope
	for _, l := range strings.Split(diff, "\n") {
		if isCommitSHA(l) {
			flush(net)
			net = 0
			pre, post = checklistScope{}, checklistScope{}
			continue
		}
		// A diff header line begins with +++ or ---; it is not content.
		if strings.HasPrefix(l, "+++") || strings.HasPrefix(l, "---") {
			continue
		}
		inPre, inPost := true, true
		if checklist && l != "" {
			switch l[0] {
			case ' ':
				inPre, inPost = pre.holds(l[1:]), post.holds(l[1:])
			case '+':
				inPost = post.holds(l[1:])
			case '-':
				inPre = pre.holds(l[1:])
			}
		}
		switch {
		case strings.HasPrefix(l, "+") && inPost && boxIn(l) != "":
			seen = true
			boxes++
			if boxIn(l) == "x" {
				net++
			}
		case strings.HasPrefix(l, "-") && inPre && boxIn(l) != "":
			seen = true
			boxes--
			if boxIn(l) == "x" {
				net--
			}
		}
	}
	flush(net)
	return closed, reopened, boxes, seen
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
// Returns per-change counts of user corrections, per-phase counts of everything,
// the entries whose change field is `-` (a correction with no change open — in no
// change's count, named in the report), and whether the file exists at all. Absent
// and empty are different facts: absent means capture is not in use, and printing 0
// for that would read as a flow that was never corrected.
//
// Phase values pass through exactly as the ledger spells them — a normaliser would
// be a second spelling that drifts. `6→7` is the one value with meaning here: the
// review seam ledgers its findings under it, and a finding is not a user
// correction, so those entries count by phase and never in a change's column.
func corrections(root string) (counts, byPhase map[string]int, orphans int, seen bool) {
	b, err := os.ReadFile(path.Join(root, ".agents/lessons.md"))
	if err != nil {
		return nil, nil, 0, false
	}
	counts, byPhase = map[string]int{}, map[string]int{}
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
		byPhase[phase]++
		// The exclusion outranks the orphan count: a 6→7 entry is a reviewer
		// finding wherever it sits, and the orphans line says "correction(s)".
		if phase == "6→7" {
			continue
		}
		if change == "-" {
			orphans++
			continue
		}
		counts[change]++
	}
	return counts, byPhase, orphans, true
}

// The repository root comes from git rather than a parameter: the caller's cwd is where
// git was pointed, and `--show-toplevel` is the only thing that knows how far up the root
// is. A projectDir parameter here would be a second answer that agrees only by accident.
func metrics(w io.Writer, args []string, git gitRunner, projects string) error {
	return metricsWith(w, args, git, projects, "")
}

func metricsWith(w io.Writer, args []string, git gitRunner, projects, loopLog string) error {
	asJSON, args := jsonFlag(args)
	r, err := gatherMetrics(args, git, projects, loopLog)
	if err != nil {
		return err
	}
	if asJSON {
		return writeJSON(w, metricsDocOf(r))
	}

	if len(r.rows) == 0 {
		fmt.Fprintf(w, "\n  no changes in this repository's history yet\n\n")
		return nil
	}

	fmt.Fprintf(w, "\n  %-34s %7s %8s %7s %7s %7s  %s\n", "change", "commits", "span", "closed", "reopen", "corr", "state")
	var totalCommits, totalReopen int
	var spans []changeMetrics
	for _, row := range r.rows {
		corr := corrCell(r.ledger.seen, r.ledger.counts[row.name])
		if row.err != nil {
			// Never skip silently. The footer counts names, so a skipped row makes the
			// total disagree with what is above it — and a report that says twelve while
			// showing ten is worse than one that admits it could not read two.
			fmt.Fprintf(w, "  %-34s %7s %8s %7s %7s %7s  unreadable\n", trunc(row.name, 34), "?", "?", "?", "?", corr)
			continue
		}
		m := row.m
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
			planCell(m.planSeen, m.checked-m.uncheck, m.boxes), churn, corr, state)
		totalCommits += m.commits
		totalReopen += m.uncheck
		spans = append(spans, m)
	}
	fmt.Fprintf(w, "\n  %d change(s), %d commit(s), %s of wall clock, %d box(es) reopened\n",
		len(r.rows), totalCommits, humanSpan(mergedSpan(spans)), totalReopen)
	if r.ledger.orphans > 0 {
		fmt.Fprintf(w, "  %d correction(s) outside any change\n", r.ledger.orphans)
	}
	phaseBreakdown(w, r.ledger.byPhase, r.ledger.seen)
	tokenBlock(w, r.tokens)
	loopBlock(w, r.loop)
	fmt.Fprintf(w, "\n%s\n\n%s\n\n", flowLegend, flowCeiling)
	return nil
}

// metricsReport is everything one metrics run read, gathered once. The table, the blocks
// under it and the JSON document all render from it, so a block cannot exist in one
// rendering and be computed differently, or not at all, in the other.
//
// The ledger, the tokens and the loop are read even when git has seen no change: none of
// them depends on a change folder, and the JSON reports whichever exist. The text still
// prints only its no-changes line.
type metricsReport struct {
	rows   []metricsRow
	ledger ledgerReport
	tokens tokenReport
	loop   loopReport
}

// metricsRow is one change: measured, or err when git would not answer for it.
type metricsRow struct {
	name string
	m    changeMetrics
	err  error
}

type ledgerReport struct {
	counts, byPhase map[string]int
	orphans         int
	seen            bool
}

func gatherMetrics(args []string, git gitRunner, projects, loopLog string) (metricsReport, error) {
	var r metricsReport
	var only string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			return r, fmt.Errorf("unknown flag %q — `%s metrics [<change>]`", a, invokedAs())
		}
		if only != "" {
			return r, fmt.Errorf("one change at a time, got %q and %q", only, a)
		}
		only = a
	}

	// Everything below is asked of the repository root, never of the cwd. git pathspecs
	// and os.Stat are both cwd-relative, so run from a subdirectory this reported "no
	// changes in this repository's history yet" — a plausible answer that was false.
	root, err := git("rev-parse", "--show-toplevel")
	if err != nil {
		return r, fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	root = strings.TrimSpace(root)

	names, err := changeNames(git, root)
	if err != nil {
		return r, fmt.Errorf("not a git repository, or git is unavailable: %w", err)
	}
	// Kept before the filter narrows it. The token block attributes against every change
	// git has seen, so its totals stay corpus-wide however the table is scoped.
	allNames := names
	if only != "" {
		matches := filterName(names, only)
		if len(matches) == 0 {
			return r, fmt.Errorf("git has never seen a change called %q", only)
		}
		if len(matches) > 1 {
			return r, fmt.Errorf("%q is ambiguous: %s", only, strings.Join(matches, ", "))
		}
		names = matches
		// A prefix resolves to one name here, so the token block reports what the table
		// reports rather than the letters that were typed.
		only = matches[0]
	}

	for _, n := range names {
		m, err := measure(git, root, n)
		r.rows = append(r.rows, metricsRow{name: n, m: m, err: err})
	}
	r.ledger.counts, r.ledger.byPhase, r.ledger.orphans, r.ledger.seen = corrections(root)
	r.tokens = gatherTokens(projects, root, allNames, only)
	r.loop = gatherLoop(loopLog, only)
	return r, nil
}

// phaseBreakdown prints where corrections surfaced, one row per phase value the
// ledger holds, spelled as the ledger spells it. Corrections at phase 2 are
// questions working; findings at 6→7 are questions that were missing; corrections
// after phase 8 are the expensive kind. The counting stays here and the reading
// stays human. Sorted by count then name, so the ordering is deterministic and the
// biggest number is the first thing read.
func phaseBreakdown(w io.Writer, byPhase map[string]int, seen bool) {
	if !seen {
		fmt.Fprintf(w, "  corrections not captured — no .agents/lessons.md in this repository\n")
		return
	}
	if len(byPhase) == 0 {
		fmt.Fprintf(w, "  no corrections captured yet — the ledger exists and is empty\n")
		return
	}
	phases := make([]string, 0, len(byPhase))
	for p := range byPhase {
		phases = append(phases, p)
	}
	sort.Slice(phases, func(i, j int) bool {
		if byPhase[phases[i]] != byPhase[phases[j]] {
			return byPhase[phases[i]] > byPhase[phases[j]]
		}
		return phases[i] < phases[j]
	})
	fmt.Fprintf(w, "\n  corrections by phase — where the flow got corrected, as the ledger spells it:\n")
	for _, p := range phases {
		fmt.Fprintf(w, "    %-20s %3d\n", p, byPhase[p])
	}
}

// loopTally is one change's rounds in the loop's run log.
type loopTally struct {
	rounds, net, unpriced int
	cost                  float64
	priced                bool
}

// loopReport is what the run log held. exists is false only when there is no file; err
// is set when there is one and it cannot be opened. names follows the filter; rows and
// skipped never do.
type loopReport struct {
	exists  bool
	err     error
	only    string
	rows    map[string]*loopTally
	names   []string
	skipped int
}

// gatherLoop reads the loop's run log. It opens the log for reading only: the loop is
// its one writer.
//
// The skipped count is taken before the filter, always. A malformed line names no change,
// so filtering it would make it vanish under every filter — the silent skip the count
// exists to prevent.
func gatherLoop(path, only string) loopReport {
	r := loopReport{only: only}
	f, err := os.Open(path)
	if path == "" || errors.Is(err, os.ErrNotExist) {
		return r
	}
	r.exists = true
	if err != nil {
		r.err = err
		return r
	}
	defer f.Close()

	r.rows = map[string]*loopTally{}
	// bufio.Reader, not Scanner: the transcript reader's rule, no line-length cap.
	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if t := strings.TrimSpace(line); t != "" {
			var rec loopRecord
			if json.Unmarshal([]byte(t), &rec) != nil || rec.Change == "" {
				r.skipped++
			} else {
				row := r.rows[rec.Change]
				if row == nil {
					row = &loopTally{}
					r.rows[rec.Change] = row
				}
				row.rounds++
				if rec.Closed != nil {
					row.net += *rec.Closed
				}
				if rec.CostUSD != nil {
					row.cost += *rec.CostUSD
					row.priced = true
				} else {
					row.unpriced++
				}
			}
		}
		if err != nil {
			break
		}
	}

	r.names = make([]string, 0, len(r.rows))
	for n := range r.rows {
		if only == "" || n == only {
			r.names = append(r.names, n)
		}
	}
	sort.Strings(r.names)
	return r
}

// loopBlock prints what the loop's run log says each change's rounds did and cost.
func loopBlock(w io.Writer, r loopReport) {
	if !r.exists {
		fmt.Fprintf(w, "\n  no loop runs are recorded in this repository\n")
		return
	}
	if r.err != nil {
		fmt.Fprintf(w, "\n  run log unreadable: %v\n", r.err)
		return
	}

	switch {
	case len(r.rows) == 0 && r.only == "":
		fmt.Fprintf(w, "\n  no loop runs are recorded in this repository\n")
	case len(r.names) == 0 && r.only != "":
		fmt.Fprintf(w, "\n  no loop rounds recorded for %s\n", r.only)
	case len(r.names) > 0:
		fmt.Fprintf(w, "\n  %-34s %7s %10s %12s %9s\n", "loop", "rounds", "net closed", "host est. $", "unpriced")
		for _, n := range r.names {
			row := r.rows[n]
			cost := "—"
			if row.priced {
				cost = fmt.Sprintf("%.2f", row.cost)
			}
			fmt.Fprintf(w, "  %-34s %7d %10d %12s %9d\n", trunc(n, 34), row.rounds, row.net, cost, row.unpriced)
		}
		fmt.Fprintf(w, "  cost is the host's reported estimate in dollars, summed as reported — never priced here\n")
	}
	if r.skipped > 0 {
		fmt.Fprintf(w, "  %d run log line(s) skipped — not a loop record\n", r.skipped)
	}
}

// tokenReport is what the session transcripts say this repository's work cost.
// configured is false when no transcript root was given, found when none exists for this
// repository; either way nothing below them was read.
type tokenReport struct {
	configured, found        bool
	attributed, unattributed usageTotals
	entries, offEntries      int
	only                     string
	attributedHere           bool
	total, noPhase           usageTotals
	byPhase                  map[string]usageTotals
	noPhaseSeen              bool
}

// gatherTokens takes allNames rather than the filtered set, always: the corpus totals must
// not move when a change filter is applied, or "attributed + unattributed = corpus" stops
// being readable off the output and the unattributed row stops meaning anything.
func gatherTokens(projects, root string, allNames []string, only string) tokenReport {
	r := tokenReport{configured: projects != "", only: only}
	if !r.configured {
		return r
	}
	es, found := readUsage(projects, root)
	r.found = found
	if !found {
		return r
	}

	byChange, unattributed := attribute(es, allNames)
	for _, t := range byChange {
		r.attributed = r.attributed.plus(t)
	}
	r.unattributed = unattributed
	r.entries = len(es)
	for _, e := range es {
		if attributeBranch(e.branch, allNames) == "" {
			r.offEntries++
		}
	}

	if only == "" {
		return r
	}
	// Presence in the map, never a zero total. A change can be attributed and cost
	// nothing — a <synthetic> entry does exactly that — and rendering it like a change
	// no session reached would collapse the two states outcome 3 exists to separate.
	r.total, r.attributedHere = byChange[only]
	if !r.attributedHere {
		return r
	}
	r.byPhase = map[string]usageTotals{}
	for _, e := range es {
		if attributeBranch(e.branch, allNames) != only {
			continue
		}
		t := usageTotals{e.in, e.out, e.cacheW, e.cacheR}
		if e.skill == "" {
			// Never distributed across the phases that did name a skill: 4,216 of 8,944
			// entries carry no attributionSkill, and sharing them out would invent the number.
			r.noPhase, r.noPhaseSeen = r.noPhase.plus(t), true
			continue
		}
		r.byPhase[e.skill] = r.byPhase[e.skill].plus(t)
	}
	return r
}

// tokenBlock prints the token report.
func tokenBlock(w io.Writer, r tokenReport) {
	if !r.configured {
		fmt.Fprintf(w, "\n  no session transcripts configured — token cost not measured\n")
		return
	}
	if !r.found {
		fmt.Fprintf(w, "\n  no session transcripts for this repository — token cost not measured\n")
		return
	}

	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "tokens", "input", "output", "cache-w", "cache-r")
	tokenRow(w, "attributed", r.attributed, true)
	tokenRow(w, "unattributed", r.unattributed, true)
	if n, off := r.entries, r.offEntries; n > 0 {
		// Integer division truncates, so one miss in two hundred reads as 0% — the
		// opposite of what a number the ceiling calls an error bar is for. Anything
		// under one percent says so rather than rounding itself away.
		pct := fmt.Sprintf("%d%%", off*100/n)
		if off > 0 && off*100/n == 0 {
			pct = "<1%"
		}
		fmt.Fprintf(w, "\n  %s of session entries could not be attributed to a change\n", pct)
	}

	if r.only == "" {
		return
	}
	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "by change", "input", "output", "cache-w", "cache-r")
	tokenRow(w, trunc(r.only, 16), r.total, r.attributedHere)

	if !r.attributedHere {
		return
	}
	fmt.Fprintf(w, "\n  %-16s %14s %14s %14s %16s\n", "by phase", "input", "output", "cache-w", "cache-r")
	phases := make([]string, 0, len(r.byPhase))
	for p := range r.byPhase {
		phases = append(phases, p)
	}
	sort.Strings(phases)
	for _, p := range phases {
		tokenRow(w, trunc(p, 16), r.byPhase[p], true)
	}
	tokenRow(w, "unattributed", r.noPhase, r.noPhaseSeen)
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
