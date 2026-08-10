# refine-model-selector Implementation Plan

> **For agentic workers:** this plan is **live state**, not a handoff. Phase 6
> (`build-and-check`) executes it one task at a time. The orchestrator is the only
> writer — sub-agents report, the orchestrator marks the box.

**Goal:** the model selector groups its rows by model with a rule between groups,
carries an `all` row that marks everything, and shows a footer legend that is true
for the screen it is on.

**Architecture:** all of it inside `internal/ui`. Rows still arrive through the
`ListAgents` callback and the catalogue through `WithAgents`; the sort is applied to
the rows the package was handed, using the catalogue it was handed. No new import, no
filesystem, no knowledge of what an agent file is. The `all` row is a rendering and a
cursor position — it never becomes a name in `MarkedAgents` or in `ApplyModel`.

**Tech Stack:** Go 1.26.5, bubbletea 1.3.10, lipgloss 1.1.0. Standard library testing.

## Global Constraints

- Six gates pass before any commit: `gofmt -l .` (silent), `go vet ./...`,
  `go test ./... -count=1`, `scripts/check-payload`,
  `skills/record-work/spec-drift --self-test`, `skills/record-work/spec-drift --anchors`.
- `internal/ui` imports neither `internal/target` nor `internal/agentmodel`.
- Legible without colour: `[x]`/`[ ]` glyphs and a rule, never emphasis alone.
- No astral-plane runes, no East Asian Ambiguous Width glyphs. `─` (U+2500) only.
- Contrast floor 4.5:1 text, 3:1 borders — enforced by `contrast_test.go`.
- Fluid frame 58–98 content columns; padding never truncates.
- Delta spec: `.agents/changes/refine-model-selector/spec.md`.

---

### Task 1: Rank models once, and sort the rows by it

**Spec:** outcome "The rows are grouped by model, with a rule between groups".
**Closes:** criteria 1 and 2.
**Waits on:** nothing. Start here.

**Files:**
- Modify: `internal/ui/models.go` — extract `modelRank`, add `sortRowsByModel`, call it in `openSelector` and in the `tab` branch of `updateSelector`
- Test: `internal/ui/models_test.go`

**Interfaces:**
- Produces: `func modelRank(order []ModelChoice) map[string]int` — catalogue index for
  a known model, `len(order)` for the session default `""`, and a model absent from
  the map sorts after everything. `func sortRowsByModel(rows []AgentRow, order []ModelChoice) []AgentRow`
  — returns a new slice, does not mutate the argument.
- Consumes: nothing.

- [x] **Step 1: Write the failing tests**

```go
// order joins the row names so a wrong order reads as a sentence rather than as two
// slices to diff by eye.
func order(rows []AgentRow) string {
	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.Name
	}
	return strings.Join(names, " ")
}

func TestRowsAreGroupedByModel(t *testing.T) {
	choices := []ModelChoice{{Name: ""}, {Name: "haiku"}, {Name: "sonnet"}}
	rows := []AgentRow{
		{Name: "work-reviewer"},
		{Name: "review-lens-tests", Model: "haiku"},
		{Name: "spec-writer", Model: "sonnet"},
		{Name: "review-lens-design", Model: "haiku"},
	}

	// Catalogue order, cheapest first, the session default last. Names sort inside
	// a group.
	want := "review-lens-design review-lens-tests spec-writer work-reviewer"
	if got := order(sortRowsByModel(rows, choices)); got != want {
		t.Fatalf("grouped order = %q, want %q", got, want)
	}
}

func TestAnUnknownModelGetsItsOwnGroup(t *testing.T) {
	choices := []ModelChoice{{Name: ""}, {Name: "haiku"}}
	rows := []AgentRow{
		{Name: "b", Model: "some-future-model"},
		{Name: "a", Model: "haiku"},
		{Name: "c"},
	}

	// haiku, then the session default, then the model this build does not know —
	// the same position Tally gives it, from the same ranking.
	if got := order(sortRowsByModel(rows, choices)); got != "a c b" {
		t.Fatalf("order = %q, want %q", got, "a c b")
	}
}
```

`sortRowsByModel` is a pure function on rows the package was handed, so these two need
no panel and no helper.

- [x] **Step 2: Run them and watch them fail**

Run: `go test ./internal/ui/ -run 'TestRowsAreGroupedByModel|TestAnUnknownModelGetsItsOwnGroup' -v`
Expected: FAIL — `undefined: sortRowsByModel`.

- [x] **Step 3: Extract the ranking and write the sort**

```go
// modelRank orders models the way the catalogue does: cheapest first, the session
// default after them, and a model this build does not know about last.
//
// One function rather than two, because the menu tally and the selector below it are
// one screen. Two orderings of the same models would be read as a bug, and it would
// be one.
func modelRank(order []ModelChoice) map[string]int {
	rank := make(map[string]int, len(order))
	for i, c := range order {
		rank[c.Name] = i
	}
	// The session default is not a price, so it cannot lead a list answering "how
	// much of this is still expensive?".
	rank[""] = len(order)
	return rank
}

// sortRowsByModel groups the rows by model so the screen reads at a glance. Names
// sort inside a group: two agents on one model in an order nobody chose is a list
// that reorders itself between sessions.
func sortRowsByModel(rows []AgentRow, order []ModelChoice) []AgentRow {
	out := append([]AgentRow(nil), rows...)
	rank := modelRank(order)
	sort.SliceStable(out, func(i, j int) bool {
		ri, oki := rank[out[i].Model]
		rj, okj := rank[out[j].Model]
		if oki != okj {
			return oki
		}
		if ri != rj {
			return ri < rj
		}
		return out[i].Name < out[j].Name
	})
	return out
}
```

Then rewrite `Tally`'s inline rank as `rank := modelRank(order)`, deleting the six
lines it duplicates.

- [x] **Step 4: Apply it where rows enter the selector**

In `openSelector`, after the successful `m.listAgents(...)`:

```go
	m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
```

And in the `tab`/`s` branch of `updateSelector`, at the same place:

```go
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
```

And in `applyChosenModel`'s reload, after the marks are restored — the models just
changed, so the groups have to move with them:

```go
		m.panel.Agents = sortRowsByModel(rows, m.modelChoices)
```

- [x] **Step 5: Run the package**

Run: `go test ./internal/ui/ -count=1`
Expected: PASS. `TestTallyPutsTheSessionDefaultLast` and
`TestRowsShowTheNewModelAfterApplying` both cover ground this step touched — if
either fails, the ranking was changed, not extracted.

- [x] **Step 6: Commit**

```bash
git add internal/ui/models.go internal/ui/models_test.go
git commit -m "feat(panel): group the selector's rows by model"
```

---

### Task 2: Draw the rule between groups

**Spec:** outcome "The rows are grouped by model, with a rule between groups".
**Closes:** criterion 3.
**Waits on:** Task 1 — the rule is drawn where the model changes, which is only
meaningful once the rows are sorted.

**Files:**
- Modify: `internal/ui/models.go` — `func (t Theme) selector(p Panel) string`
- Test: `internal/ui/models_test.go`

**Interfaces:**
- Consumes: `sortRowsByModel` from Task 1 (through `p.Agents`, already sorted).
- Produces: `func (t Theme) groupRule(width int) string` — the indented dim rule,
  used by Task 3 as well.

- [x] **Step 1: Write the failing test**

```go
// rules counts the group rules in a rendered selector. The frame's own ├───┤ rows are
// not counted: they start with a junction glyph, the group rule is indented.
func rules(rendered string) int {
	n := 0
	for _, line := range strings.Split(strip(rendered), "\n") {
		if strings.HasPrefix(line, "  ─") {
			n++
		}
	}
	return n
}

func TestGroupRuleSitsOnlyBetweenGroups(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	// threeAgents is one opus row and two session rows: two groups, one rule.
	mixed := theme.selector(Panel{Width: 90, Agents: sortRowsByModel(threeAgents(), catalogue())})
	if got := rules(mixed); got != 1 {
		t.Fatalf("two groups drew %d rules, want 1", got)
	}
	lines := strings.Split(strings.TrimRight(strip(mixed), "\n"), "\n")
	if strings.HasPrefix(lines[len(lines)-1], "  ─") {
		t.Fatal("a rule was drawn after the last group")
	}

	uniform := theme.selector(Panel{Width: 90, Agents: []AgentRow{
		{Name: "a", Model: "haiku"}, {Name: "b", Model: "haiku"},
	}})
	// A division of nothing is not a division.
	if got := rules(uniform); got != 0 {
		t.Fatalf("one group drew %d rules, want none", got)
	}
}
```

`catalogue()` is the four `ModelChoice` values `selectorModel` already builds inline —
lift them into a helper there and call it from both, so the fixture cannot drift.

**Task 3 adds a second rule** — the one under the `all` row — and these two counts
become 2 and 1. That is the feature changing under the test, not the test being
softened: the assertion still says "one rule per boundary and none anywhere else", and
Task 3 adds a boundary. Update the numbers there, in the same commit as the row.

- [x] **Step 2: Run it and watch it fail**

Run: `go test ./internal/ui/ -run TestGroupRuleSitsOnlyBetweenGroups -v`
Expected: FAIL — the count is 0, no rules are drawn.

- [x] **Step 3: Draw it**

Add beside `selector`:

```go
// groupRule divides one model's rows from the next.
//
// An indented dim rule, deliberately not the frame's ├───┤ junction: that glyph means
// "a new section of the panel", and three of them inside one list would read as three
// panels. A blank line was the other candidate and it is worse — inside a bordered
// frame a blank line reads as the end of the list.
func (t Theme) groupRule(width int) string {
	n := width - 4
	if n < 1 {
		n = 1
	}
	return "  " + Fg(t.Dim).Render(strings.Repeat("─", n))
}
```

In `selector`, track the previous row's model and emit the rule when it changes:

```go
	cw := ContentWidth(p.Width)
	for i, a := range p.Agents {
		if i > 0 && a.Model != p.Agents[i-1].Model {
			rows = append(rows, t.groupRule(cw))
		}
		...
	}
```

`i > 0` is what keeps the rule out from above the first group; the loop ending is what
keeps it out from below the last.

- [x] **Step 4: Run the package**

Run: `go test ./internal/ui/ -count=1`
Expected: PASS, including `TestSelectorFrameIsFlushAtEveryWidth` — the rule is a
content row like any other and `padTo` measures it.

- [x] **Step 5: Commit**

```bash
git add internal/ui/models.go internal/ui/models_test.go
git commit -m "feat(panel): rule the selector between model groups"
```

---

### Task 3: The `all` row

**Spec:** outcome "`all` is a row, and space on it marks everything".
**Closes:** criteria 4, 5, 6, 7, 9, 10.
**Waits on:** Task 2 — it reuses `groupRule` for the line beneath it.

**Files:**
- Modify: `internal/ui/models.go` — `updateSelector` (`up`, `down`, `" "`), `selector`
- Test: `internal/ui/models_test.go`

**Interfaces:**
- Consumes: `t.groupRule(width)` from Task 2.
- Produces: `const allRow = 0` — the cursor position of the `all` row. Agent `i` is at
  cursor `i+1`; `p.AgentCursor-1` indexes `p.Agents`. This offset is the whole risk of
  the task.

- [x] **Step 1: Write the failing tests**

Sorted, `threeAgents()` is `review-security` (opus) first, then `review-design` and
`review-tests` (session). Every cursor count below assumes that order — it is Task 1's
output, and writing the counts against the unsorted fixture is the easiest way to get
these wrong.

```go
func TestSpaceOnTheAllRowMarksEveryAgent(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	// The selector opens on `all`, so the first space marks everything.
	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 3 {
		t.Fatalf("space on the all row marked %v, want all three", got)
	}

	m = key(m, " ")
	if got := m.MarkedAgents(); len(got) != 0 {
		t.Fatalf("space on the all row again left %v marked, want none", got)
	}
}

func TestTheAllRowBoxFollowsEveryAgent(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ")     // every agent marked
	m = key(m, "down")  // onto the first agent
	m = key(m, " ")     // unmark it

	first := strings.Split(strip(darkTheme().selector(m.panel)), "\n")[0]
	if strings.Contains(first, "[x]") {
		t.Fatalf("the all row still reads %q with an agent unmarked", first)
	}
}

func TestTheAllRowIsNeverAppliedAsAnAgent(t *testing.T) {
	m, rec := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, " ")      // cursor on all: mark everything
	m = key(m, "m")      // open the catalogue from the all row
	m = key(m, "down")   // off the session default, onto haiku
	m = key(m, "enter")

	for _, n := range rec.names {
		if n == "all" {
			t.Fatal("the all row reached ApplyModel as an agent name")
		}
	}
	if len(rec.names) != 3 {
		t.Fatalf("applied to %v, want the three agents", rec.names)
	}
}

func TestTheCursorMarksTheRowItPointsAt(t *testing.T) {
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)

	m = key(m, "down")  // all -> review-security
	m = key(m, "down")  // -> review-design
	m = key(m, " ")

	got := m.MarkedAgents()
	if len(got) != 1 || got[0] != "review-design" {
		t.Fatalf("marked %v, want [review-design] — the all row shifted the cursor", got)
	}
}

func TestTheAllRowIsLegibleWithoutColour(t *testing.T) {
	forceTrueColor(t)
	m, _ := selectorModel(t, threeAgents())
	m = openSelector(t, m)
	m = key(m, " ")

	if plain := strip(m.View()); !strings.Contains(plain, "[x] all") {
		t.Fatalf("the all row is not legible without colour:\n%s", plain)
	}
}
```

`MarkedAgents` returns names in screen order, which after Task 1 is grouped order —
that is what `TestTheCursorMarksTheRowItPointsAt` relies on.

- [x] **Step 2: Run them and watch them fail**

Run: `go test ./internal/ui/ -run 'AllRow|TestTheCursorMarksTheRowItPointsAt' -v`
Expected: FAIL — `undefined: allRow`.

- [x] **Step 3: Move the cursor over one extra row**

```go
// allRow is the cursor position of the `all` row, which sits above every group.
//
// The consequence is that screen position and agent index stop being the same number:
// agent i is at cursor i+1. Every read of p.Agents from the cursor goes through
// AgentCursor-1, and getting that wrong marks the neighbouring agent silently — which
// is why it has a test of its own.
const allRow = 0

// selectorRows is how many rows the cursor can reach: the agents, plus `all`.
func selectorRows(p Panel) int { return len(p.Agents) + 1 }
```

In `updateSelector`, `up`/`down` wrap over `selectorRows(m.panel)` instead of
`len(m.panel.Agents)`, and `" "` splits:

```go
	case " ":
		if len(m.panel.Agents) == 0 {
			return m, true
		}
		if m.panel.AgentCursor == allRow {
			return m.markAll(), true
		}
		rows := append([]AgentRow(nil), m.panel.Agents...)
		rows[m.panel.AgentCursor-1].Marked = !rows[m.panel.AgentCursor-1].Marked
		m.panel.Agents = rows
		return m, true
```

Pull the body of the existing `a` case into `markAll` so the key and the row cannot
drift apart:

```go
// markAll marks every row, or clears them when they are already all marked. One
// gesture both ways: a control that only ever adds leaves no way back but pressing
// space once per row.
func (m Model) markAll() Model {
	marking := len(m.MarkedAgents()) < len(m.panel.Agents)
	rows := append([]AgentRow(nil), m.panel.Agents...)
	for i := range rows {
		rows[i].Marked = marking
	}
	m.panel.Agents = rows
	return m
}
```

`case "a"` becomes `return m.markAll(), true`.

- [x] **Step 4: Render it**

At the top of `selector`, before the agent loop and after the empty guard:

```go
	box := "[ ]"
	if countMarked(p.Agents) == len(p.Agents) {
		box = "[x]"
	}
	colour, cursor := t.Steel, " "
	if p.AgentCursor == allRow && !p.ChoosingModel {
		colour, cursor = t.Gold, "❯"
	}
	rows = append(rows,
		"  "+Fg(colour).Render(cursor+" "+box+" all"),
		t.groupRule(cw),
	)
```

and the agent loop's cursor test becomes `i+1 == p.AgentCursor`.

- [x] **Step 5: Run the package**

Run: `go test ./internal/ui/ -count=1`
Expected: PASS. `TestSpaceMarksAndUnmarksTheCurrentRow` and
`TestChosenModelReachesOnlyTheMarkedRows` are the two that fail if the offset is
wrong — they open the selector, press `space`, and expect the first agent. They will
need `down` before `space` now, which is a real change to what the test asserts and
must be made deliberately, not to get green.

- [x] **Step 6: Commit**

```bash
git add internal/ui/models.go internal/ui/models_test.go
git commit -m "feat(panel): give the selector an all row"
```

---

### Task 4: A footer that belongs to its screen

**Spec:** outcome "The footer legend tells the truth about the screen it is on".
**Closes:** criterion 8.
**Waits on:** nothing — independent of Tasks 1–3, may be done first or in parallel.

**Files:**
- Modify: `internal/ui/panel.go:308-322` — `func (t Theme) footer(p Panel, width int) string`
- Test: `internal/ui/panel_test.go`

**Interfaces:**
- Consumes: `p.InSelector`, `p.ChoosingModel`, `p.Confirm` — all already on `Panel`.
- Produces: nothing other tasks use.

- [x] **Step 1: Write the failing test**

```go
func TestTheFooterFollowsTheScreen(t *testing.T) {
	forceTrueColor(t)
	theme := darkTheme()

	menu := strip(theme.footer(Panel{Version: "v0"}, 90))
	if !strings.Contains(menu, "⏎ select") {
		t.Fatalf("the menu footer changed: %q", menu)
	}

	sel := strip(theme.footer(Panel{Version: "v0", InSelector: true}, 90))
	for _, want := range []string{"space mark", "a all", "m model", "tab scope"} {
		if !strings.Contains(sel, want) {
			t.Fatalf("the selector footer does not mention %q: %q", want, sel)
		}
	}
	if strings.Contains(sel, "⏎ select") {
		t.Fatalf("the selector footer still promises ⏎ select: %q", sel)
	}

	cat := strip(theme.footer(Panel{Version: "v0", InSelector: true, ChoosingModel: true}, 90))
	if !strings.Contains(cat, "⏎ apply") || strings.Contains(cat, "space mark") {
		t.Fatalf("the catalogue footer lists keys that do nothing: %q", cat)
	}

	confirm := strip(theme.footer(Panel{Version: "v0", InSelector: true, Confirm: "sure?"}, 90))
	if !strings.Contains(confirm, "y yes") {
		t.Fatalf("a confirm no longer wins the footer: %q", confirm)
	}
}
```

- [x] **Step 2: Run it and watch it fail**

Run: `go test ./internal/ui/ -run TestTheFooterFollowsTheScreen -v`
Expected: FAIL — the selector footer still says `⏎ select`.

- [x] **Step 3: Switch on the screen**

Replace the two-branch `hints` with:

```go
	// The hints belong to the screen, not to the program. A legend that lists ⏎ select
	// while ⏎ opens a catalogue is worse than no legend: it is read once, believed, and
	// never re-read. The correct keys used to live in the opening notice, where the
	// first apply overwrote them and they were gone for the session.
	hints := "↑↓ · ⏎ select · tab scope · q quit"
	switch {
	case p.Confirm != "":
		// While a question is open the only keys that matter are its answers.
		// Listing the others invites pressing one by reflex.
		hints = "y yes · n no"
	case p.ChoosingModel:
		hints = "↑↓ · ⏎ apply · esc back"
	case p.InSelector:
		hints = "↑↓ · space mark · a all · m model · tab scope · esc back"
	}
```

`p.Confirm` stays first: it is the existing precedence and a confirm can be open over
either screen.

- [x] **Step 4: Check the width at the floor**

The selector's hint string is the longest in the program. Render at
`MinContentWidth`:

Run: `LIBRETTO_ASCII=safe COLUMNS=60 go run ./cmd/libretto preview`
Expected: the footer sits under the box without widening the centred block. If it
does widen it, drop `↑↓` from the **selector** line only — movement is already taught
on the menu screen, and `space`/`a`/`m` are what this change exists to teach. Do not
shorten by dropping `tab scope`: a key that changes the destination while the cursor
does not move is the one that looks broken when unlisted.

- [x] **Step 5: Run the package**

Run: `go test ./internal/ui/ -count=1`
Expected: PASS, `TestPanelIsCentredWhenThereIsRoom` included.

- [x] **Step 6: Commit**

```bash
git add internal/ui/panel.go internal/ui/panel_test.go
git commit -m "fix(panel): tell the truth in the selector's footer"
```

---

### Task 5: Land the delta and close the change

**Spec:** the whole delta.
**Closes:** nothing new — it is what makes criteria 1–10 part of the standing contract.
**Waits on:** Tasks 1–4, and all six gates.

**Files:**
- Modify: `.agents/specs/panel/spec.md` — the "The model selector" outcome, its scope
  boundaries, its prior decisions, its verification criteria
- Delete: `.agents/changes/refine-model-selector/`

- [ ] **Step 1: Fold the delta into the capability spec**

Under `### The model selector — the panel's second screen`: update the mock to show
the `all` row and the group rules, and add the bullets for grouping, the `all` row and
the per-screen footer. Add the ordering decision to **Prior decisions**. Add criteria
1–10 to **Verification criteria** with their `Proof:` lines.

- [ ] **Step 2: Delete the change folder**

```bash
git rm -r .agents/changes/refine-model-selector
```

- [ ] **Step 3: Run all six gates, foreground, output read**

```bash
gofmt -l .
go vet ./...
go test ./... -count=1
scripts/check-payload
skills/record-work/spec-drift --self-test
skills/record-work/spec-drift --anchors
```

Every `Proof:` added in Step 1 must resolve. `--anchors` checks the **test name**, not
just the file — an invented name passes a file-level check and has done so here twice.

- [ ] **Step 4: Commit**

```bash
git add .agents/
git commit -m "docs(spec): land refine-model-selector onto panel"
```
