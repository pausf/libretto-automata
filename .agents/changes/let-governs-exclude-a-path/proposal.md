# let-governs-exclude-a-path

Tracker: none
Targets: payload
Durable decisions: the one in Prior decisions below

## The ask, verbatim

> give-ci-workflow-tests-to-ci: cli's Governs: claims release_workflow_test.go,
> label_workflow_test.go and gates_test.go, which test ci's workflows.
>
> give-readme-test-one-governor: cmd/libretto/readme_test.go is claimed by both the cli and
> readme Governs: globs.

(Two ideas queued 2026-10-22, picked up together 2026-10-23: neither can be done without the
same missing piece, so they are one change. Both captures are retired in it.)

## Reading

`Governs:` is a list of globs and nothing else, so `cmd/libretto/**` cannot give up a file
it covers. Four test files sit under it whose promise belongs elsewhere: three test ci's
workflows and Makefile, one tests the README. spec-drift warns for every capability claiming
a staged path, so each of those files makes a change look like cli drift when the spec that
should move is ci's or readme's — measured in report-drift-trend-in-metrics: 6 of the 10
hand-checked warned landings were overlaps of this kind.

The piece missing: a `!glob` entry in `Governs:` that excludes a path. A path is governed by
a capability when one of its plain globs matches and none of its `!` globs does.

## Decision log

### Session 2026-10-23
- Q: what marks an exclusion? → A: a leading `!` on one `Governs:` entry, as .gitignore does — the reading most authors bring (assumed; if wrong, a separate `Excludes:` line)
- Q: does an exclusion that matches nothing count as a dead claim? → A: yes — it is a glob pointing at nothing, the same staleness (assumed)
- Q: do cli's remaining overlaps go too? → A: only the four the queue named; the rest are a reading for whoever next touches them (assumed)

### Session 2026-10-23, after review-spec
- Q: the Go wiki matches Governs: too (governedSplit, which cli promises mirrors spec-drift) — in or out? → A: in. Two matching authorities is what cli's promise forbids; the footer honours `!` by the same rule, per capability (assumed)
- Q: governedSplit pools every capability's globs into one list — with exclusions that is wrong (cli's `!` would hide a path ci claims) → A: it takes one entry list per capability; a path is governed when some capability governs it (assumed)
- Q: is editing another capability's Governs: line amending it? → A: yes — it changes what that capability owns; spec-cli.md and spec-ci.md carry those edits. readme's line does not change (assumed)
- Q: which readers does the self-test drive? → A: both, not only governed_by: a --trace orphan case and a default-mode case through --block (assumed)
- Q: the two other places that define orphan as "no Governs: glob claims" (record-work's --trace table, spec-drift's header)? → A: both are rewritten to "no Governs: line governs" in this change (assumed)
- Q: which write-spec passage? → A: "And the spec declares what it governs" (assumed)
- Q: how is Outcome 5 observed? → A: with the worktree's spec-drift (./skills/record-work/spec-drift, never the installed one), default mode, gates_test.go staged alone: it names ci and not cli (assumed)

## Outcomes

1. A `!` entry in a capability's `Governs:` removes the paths it matches from that capability.
2. spec-drift's default mode names a staged path's spec by that rule, so a ci workflow test
   staged alone points at ci and no longer at cli.
3. `--trace` counts a path as owned by that rule; an excluded path no other capability
   claims is an orphan.
4. `--trace` reports a `!` entry matching no tracked file as a dead claim.
5. cli gives up `release_workflow_test.go`, `label_workflow_test.go`, `gates_test.go` and
   `readme_test.go`; ci claims the first three. readme already claims the fourth.

## Scope boundaries

**In:** `skills/record-work/spec-drift` (the rule, its self-test cases, its header comment),
`skills/record-work/SKILL.md` (the `--trace` table's orphan row), `write-spec`'s "And the spec
declares what it governs" section; through `spec-cli.md`, the wiki footer's matcher and
cli's `Governs:` line; through `spec-ci.md`, ci's `Governs:` line.

**Out:** negation inside a glob (`[!x]` keeps meaning what `case` means); ordering rules (no
"later wins" — exclusion always beats inclusion); how the wiki *renders* the line (verbatim);
any other overlap in cli.

## Constraints

- The rule lives in one function, `governed_by`, used by every reader in spec-drift.
- Every existing self-test case stays unedited; a `Governs:` with no `!` behaves as today.
- No pipe into `rg -q`.

## Prior decisions

- **`Governs:` can exclude a path with a `!` entry.** Decided 2026-10-23: four test files under
  `cmd/libretto/**` prove other capabilities' promises, and a glob list with no subtraction
  could only keep claiming them, so drift warnings pointed at cli for ci's and readme's tests.
  Exclusion beats inclusion, as in .gitignore; an exclusion matching nothing is a dead claim.
  Lost: a separate `Excludes:` line (a second key to keep beside the first) and listing cli's
  files one by one (a glob nobody maintains).

## Task breakdown

- `governed_by` in spec-drift, used by default mode and `--trace`, with self-test cases;
- the three `Governs:` lines; the `write-spec` sentence.

## Verification criteria

- **Where** a capability's `Governs:` holds a `!` entry matching a path, `governed_by`
  **shall** answer that the capability does not govern it, even when a plain entry matches.
  Proof: skills/record-work/spec-drift --self-test
- **If** a path is excluded by the only capability whose plain entries match it, **then**
  `spec-drift --trace` **shall** list it as an orphan.
  Proof: skills/record-work/spec-drift --self-test
- **When** a path one capability excludes and another claims is staged alone, spec-drift's
  default mode **shall** name only the claiming capability's spec.
  Proof: skills/record-work/spec-drift --self-test
- **While** a `Governs:` line holds no `!` entry, spec-drift **shall** treat its paths as it
  did before this change.
  Proof: skills/record-work/spec-drift --self-test
- **If** a `!` entry in `Governs:` matches no tracked file, **then** `spec-drift --trace`
  **shall** report it as a dead claim.
  Proof: skills/record-work/spec-drift --self-test
- **When** `write-spec` defines `Governs:` in "And the spec declares what it governs", it
  **shall** name the `!` exclusion. Proved as
  wiring only, pinned on `a leading .!. excludes`. Ceiling named: the row pins the phrase.
  Proof: scripts/check-payload

## Approach

1. `governed_by <path> <entries...>`: false unless some plain entry matches the path under
   `case`; false when any `!` entry, with its `!` stripped, matches. Default mode (around line
   212) and the orphan check in `--trace` call it per spec, with that spec's whole entry list
   — the orphan check today matches claim by claim, which cannot see an exclusion, so it is
   regrouped by spec. The dead-claim check strips a leading `!` before matching.
2. Self-test: a spec claiming `src/**` and `!src/gen/**`; `src/gen/x.go` is not governed,
   `src/a.go` is; a `Governs:` with no `!` gives today's answers; a `!` entry matching nothing
   is dead.
3. cli: `Governs: cmd/libretto/** install.sh !cmd/libretto/readme_test.go
   !cmd/libretto/release_workflow_test.go !cmd/libretto/label_workflow_test.go
   !cmd/libretto/gates_test.go`. ci gains the three workflow tests.
4. `write-spec`: its `Governs:` paragraph gains the `!` sentence; one `check-payload` row.

## Checklist

- [x] **`Governs:` excludes with `!` in spec-drift.** `governed_by`, both readers regrouped,
  five self-test cases red first (the helper, the unchanged case, the dead `!`, the orphan,
  the default-mode naming), the header comment, record-work's row, the `write-spec` sentence
  and its row. Closes: the six payload criteria. Waits on: nothing.
- [x] **The wiki footer honours `!` per capability.** `governedSplit` takes one entry list per
  capability. Closes: spec-cli's criteria. Waits on: box 1.
- [x] **The four test files move to their owners.** cli's and ci's `Governs:` lines; observed
  with the worktree's spec-drift: `gates_test.go` staged alone names ci, not cli. Closes:
  Outcome 5, observed. Waits on: boxes 1–2.
  Observed 2026-10-23 with ./skills/record-work/spec-drift: gates_test.go staged alone named cli
  before the move and ci after; readme_test.go named readme alone; --trace reported no dead claim
  from the new lines.
