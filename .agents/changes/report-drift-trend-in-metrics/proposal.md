# report-drift-trend-in-metrics

Tracker: none
Targets: cli
Durable decisions: every entry in Prior decisions below

## The ask, verbatim

> Drift trend en `metrics`: cuántas veces el spec de cada capability se movió DESPUÉS
> del código. Convierte `spec-drift` (warn-only) en una señal de salud medible. Ninguna
> tool de la industria tiene esto todavía.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Derive per-capability drift events from git history — spec edits landing after the code
they govern — and report the trend in `libretto metrics`.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed or measured)

- Q: what is a drift event, precisely? → A: **the ask's own reading — a commit touching a
  capability's `spec.md` after earlier commits touched its governed code without it — was
  prototyped on this repository's whole history, and it is mostly noise.** At commit
  granularity (498 non-merge commits) it yields 47 events, and **37 of the 47 catch-up
  commits are landing commits**: the flow working as designed, since a change's code
  commits precede the commit that applies its delta. Read at landing granularity instead
  (70 first-parent commits of `main`) it yields 13 events; all 13 were read, and **one** —
  `readme`, PR #72 after PR #63 — moved the spec for a reason related to the code before
  it. Numbers under *Approach*. Measured, not assumed.
- Q: does the "moved after" count ship? → A: **no** (assumed). A spec moving later is a
  time order, not a cause: in 12 of 13 landing-level events the later spec edit was about
  something else — a badge, a diagram, a new command. A health signal wrong twelve times
  in thirteen is the warning that gets muted. If wrong — the user wants the event count
  anyway — it is one more counter over the same walk, shipped with its measured precision
  in the legend.
- Q: what ships instead? → A: **`spec-drift`'s own warning, replayed over every landing**:
  per capability, how many landings touched its governed code, how many of those
  `spec-drift` would have warned on (the spec did not move and no delta targeting it
  moved), and how many warned landings came after the spec last moved (assumed). That is
  the ask's second sentence taken literally — the warn-only check turned into a number —
  and it is exact about what it counts. Measured: 27 warned of 120 capability-landings.
  Of ten warned landings read against their diffs, **none** changed a promise without its
  spec: six touched a path whose semantic owner is another capability whose spec moved in
  that same landing, three were edits inside an existing contract, one is arguable. So the
  number measures *where the warning fires*, which is where `Governs:` is wider than the
  contract, and the block says so in as many words. If wrong — nothing should ship — the
  walk and its numbers go into *Prior decisions* alone and no code is written.
- Q: what is a landing? → A: **one first-parent commit of `HEAD`**, a merge read as its
  diff against its first parent (assumed). That is the unit at which code reaches the
  branch, so a change's code commits and its later landing commit inside one request are
  one landing — which is what removes the 37 false positives above. If wrong — every
  commit is a landing — the predicate is unchanged and the 37-of-47 noise comes back.
- Q: a path two capabilities claim — counted for each, or covered when either spec moved? →
  A: **counted for each**, because `spec-drift` warns for each claiming spec independently
  and this replays `spec-drift` (assumed). Measured: the shared alternative halves `cli`'s
  warned landings (10 → 5) and moves nothing else. All ten of `cli`'s are test files under
  `cmd/libretto/` that belong semantically to `readme` (seven) or `ci` (three) — the
  overlap the block exists to make visible. If wrong, a landing is covered when any
  claimant's spec moved: one predicate, and the commit-time warning would have to change
  with it or the two tools disagree.
- Q: does an in-flight delta count as the spec moving? → A: **yes, as in `spec-drift`'s
  default mode**: a landing touching a file in a change folder whose document declares
  `Targets:` naming the capability, outside fences, is not warned (assumed). Measured: on
  this history it changes no number, because every request here landed its delta with its
  code; it matters for a change spanning several requests. **Not mirrored:**
  `spec-drift`'s exclusion of a checklist-only edit — any file in the folder counts. If
  wrong, the folder's non-checklist content is diffed per landing, which is `churnIn`'s
  whole-file read once per landing.
- Q: do landings before the capability's spec existed count? → A: **no** (assumed). The
  walk starts at the first landing that touches the spec. Measured: without it `readme`
  reports six more warned landings, all README edits made before any `readme` spec
  existed — a warning on a contract nobody had written. If wrong, those landings count and
  every new capability opens with a backlog it could not have prevented.
- Q: which `Governs:` — the one at `HEAD`, or the one each landing saw? → A: **`HEAD`'s**,
  read from disk (assumed). The historical reading is one `git show` per spec per landing.
  **Ceiling named:** a glob added later reaches back over history, and a path renamed away
  from a glob drops out of it. If wrong, the walk reads each spec at each landing that
  touched it — bounded by spec edits, not by landings.
- Q: how does Go read `Governs:` and match its globs? → A: **as `spec-drift` does, and a
  test proves it against `spec-drift` itself** (assumed). `governs_of` reads every
  `Governs:` line outside a fence, leading whitespace allowed, word-split; `matches_any`
  matches with a bash `case` under `set -f`, where `*` and `?` cross `/`. Go's
  `path.Match` is rejected: there `*` stops at `/`, so `cmd/*.go` claims
  `cmd/libretto/a/b.go` in `spec-drift` and does not in `path.Match` — measured with bash.
  The Go reading shares `land.go`'s fence rule, which already mirrors `defenced`. If wrong
  — Go execs `spec-drift` instead — `metrics` needs bash and `rg` at runtime, where it
  needs only git today.
- Q: where are the capability specs? → A: **`findSpecsDir`, `wiki`'s discovery** —
  `spec-drift`'s order, first hit wins — and `<dir>/*/spec.md` within it (assumed). One Go
  definition already exists. **Divergence named:** `spec-drift`'s `md_under` reads every
  root that exists, not the first; a project keeping specs in two roots gets two answers,
  as it already does from `wiki`.
- Q: does `metrics <change>` filter the block? → A: **no, the block is corpus-wide**, like
  the token block's totals (assumed). A capability is not a change. If wrong, the block
  keeps the rows the named change's `Targets:` name.
- Q: the JSON shape? → A: **a top-level `spec_drift` object**: `capabilities`, one object
  per row of the text block, each `name`, `landings`, `warned`, `outstanding`; a capability
  with no `Governs:` carries `name` alone; `error` in place of `capabilities` when the
  history cannot be read; `spec_drift` absent when no capability spec exists (assumed). The
  spec-moved count is `landings − warned` and is not carried, per "nothing derivable from
  the rest". Schema stays 1: an added key never bumps it. If wrong, it is a renamed key
  before anybody consumes it.
- Q: how is "the text has the block and the JSON does not" made to fail? → A: **one test
  over the whole report**, `TestMetricsJSONCarriesEveryTextBlock`: a report in which every
  block prints, rendered both ways, and a table mapping each block's text header to its
  JSON key — the corrections, token, loop and spec-drift blocks (assumed). It fails when a
  header prints and its key is absent. **Ceiling named:** a block added later with no row
  in that table passes it. This replaces the recorded "whether it was mapped is review's to
  see, not a gate's" for every block the table holds. If wrong, the parity is asserted for
  the spec-drift block alone.
- Q: what does the agreement test do where `rg` is missing? → A: **fails, and is skipped
  only under `-short`**, the real-git tests' precedent (assumed). CI installs `rg` before
  `go test`, and a skip on a missing tool is a test green for the wrong reason. If wrong,
  it skips with the reason printed.
- Q: does `AGENTS.md`'s or the README's description of `metrics` change? → A: **no**
  (assumed). `AGENTS.md`'s line names what `metrics` costs per change; the README is
  `readme`'s. If wrong, one line each, in the landing commit.
- Q: the fork — how is it built? → A: **A, a Go replay inside `metrics`**: one
  `git log --first-parent --name-only` over the history, the `Governs:` reading mirrored in
  Go and pinned against `spec-drift --trace`, one field on `metricsReport`, a text block and
  a JSON key rendered from it (assumed, the recommended option). **B** — a `--history` mode
  in `spec-drift`, exec'd by `metrics` — lost: one glob definition instead of a mirrored
  one, but `metrics` would need bash and `rg` at runtime and would be the first Go command
  to run a payload script, and the shell walk is a `case` per path per glob per landing.
  **C** — ship nothing and record the measurement — lost because the replay is exact about
  what it counts and its first run already names an overlap worth fixing.

## Outcomes

- `libretto metrics` prints a spec-drift block: one row per capability spec, sorted by
  name, holding the landings that touched its governed code, how many of those
  `spec-drift` would have warned on, and how many warned landings came after the spec last
  moved.
- The block says, beside it, that a warned landing is one the warning would have fired on
  — not measured drift — and the report's ceiling names what no walk over git can tell:
  whether a warned landing changed a promise.
- `metrics --json` carries the same rows under `spec_drift`, under the schema-1 rules: an
  unknown value absent, a known zero present.
- A capability with no `Governs:`, a project with no capability specs, and a history git
  cannot read each print a state rather than vanishing or crashing.
- The Go reading of `Governs:` agrees with `spec-drift`'s, and a test fails the day they
  diverge.
- The "moved after" count the ask named is not shipped, and the measurement that rejected
  it is recorded in `cli`'s *Prior decisions*.

## Scope boundaries

**In:**

- `cmd/libretto/specdrift.go`: `governsOf`, `governedBy`, the landing walk, the
  per-capability tally
- `cmd/libretto/land.go`: its fence-and-key reading lifted into one helper that
  `headTargets` and `governsOf` share
- `cmd/libretto/metrics.go`: a `specDrift` field on `metricsReport`, gathered in
  `gatherMetrics`, its text block after the loop block, one bullet in `flowCeiling`
- `cmd/libretto/json.go`: `specDriftJSON` and its mapping in `metricsDocOf`
- `cmd/libretto/specdrift_test.go` and `cmd/libretto/metrics_test.go`: the cases named
  under *Verification criteria*
- `.agents/specs/cli/spec.md`, at landing: the criteria, the `metrics --json` example and
  key table gaining `spec_drift`, and the *Prior decisions* entries below. The existing
  entry "One gathered report per command, rendered twice" has its ceiling sentence
  amended in place — the parity is now a test for every block in the test's table.

**Out:**

- **the "moved after" event count**, in any form. Rejected on the measurement; see *Prior
  decisions*.
- **a time series.** No per-landing list, no window, no bucket per month. The block is
  totals; a consumer wanting a trend needs per-landing rows, which is an additive key.
- **judging whether a warned landing changed a promise.** That is a reading; `review-work`
  and the human answer to the commit-time warning stay the mechanism.
- **historical `Governs:`.** `HEAD`'s globs are applied to every landing.
- **changing `spec-drift`.** No new mode, no flag; it is the reference this mirrors.
- **fixing the overlaps the first run found** — `cmd/libretto/readme_test.go` claimed by
  `cli` and `readme`, and `release_workflow_test.go`, `label_workflow_test.go` and
  `gates_test.go` claimed by `cli` while testing `ci`'s workflows. Each is a `Governs:`
  edit and a `/libretto-queue` capture, not this change.
- **resolving a base branch.** The walk is `HEAD`'s first-parent history; see
  *Constraints*.
- **the panel, `status`, `AGENTS.md`, the README, the wiki.**
- **a seventh gate, or anything in `--anchors`.**

## Constraints

- Stdlib and exec'd git only, through `gitRunner`, as `metrics` already is. No dependency,
  and no `rg` or bash at runtime.
- The walk costs **one** `git log --first-parent --name-only` over `HEAD`'s history, plus
  **at most one** `git grep` per landing and change folder, only for a landing that touches
  governed code and a change folder without moving the spec — the only landings whose
  `Targets:` matter. Matching is in memory: landings × changed paths × globs. Measured on
  this repository: 70 landings, 856 lines of log output, the log call 0.17 s.
- `--full-history` is not passed: there is no pathspec, so history simplification prunes
  nothing.
- `Governs:` is read as `governs_of` reads it and matched as `matches_any` matches: bash
  `case` semantics under `set -f` — `*` any run including `/`, `?` any one character
  including `/`, `[…]` a bracket expression with a leading `!` or `^` negating it, a
  backslash quoting the next character, everything else literal. Never `path.Match` or
  `filepath.Match`.
- Paths under the specs directory and under `.agents/changes/` are never governed code,
  as `spec-drift`'s default mode skips them.
- The JSON follows `cli`'s schema-1 rules: snake_case, an unknown value absent, a known
  zero present, a known-empty list `[]`, no prose. `jsonSchema` does not move.
- The block's legend line is prose and stays out of the JSON.
- **On a branch the walk counts the branch's own commits as landings**, because they are
  `HEAD`'s first parents. The legend line says the walk is `HEAD`'s; run on the default
  branch, it is the landed picture.

## Prior decisions

- **"The spec moved after the code" is not reported.** Assumed 2026-10-22, attacca run, on
  a measurement. Prototyped over this repository's history: at commit granularity 47
  events, 37 of whose catch-up commits were landing commits; at landing granularity 13,
  all read, one related to the code before it. **What would reopen it:** a way to tie a
  spec edit to the code it answers — a commit trailer, a delta naming its landing — which
  the flow does not write today.
- **What `metrics` reports is `spec-drift`'s warning, replayed per landing.** Assumed
  2026-10-22, attacca run. A landing is a first-parent commit of `HEAD`; it is warned for a
  capability when it touches a path that capability's `Governs:` claims, outside the specs
  and changes directories, and moves neither the capability's `spec.md` nor a file in a
  change folder whose `Targets:` names it. Counted from the first landing that touches the
  spec. Each claimant of a co-claimed path counts it. **Ceiling named:** measured on this
  repository, 0 of 10 warned landings read changed a promise; the number measures where
  the warning fires, which is where `Governs:` is wider than the contract, and the block
  says so.
- **Go mirrors `spec-drift`'s `Governs:` reading and globbing, and a test runs
  `spec-drift --trace` to prove the two agree.** Assumed 2026-10-22, attacca run. `*`
  crosses `/` in a bash `case` and does not in `path.Match`, so the stdlib matcher would be
  a second definition. **What would reopen it:** a Go binary that already needs bash and
  `rg` at runtime, at which point execing `spec-drift` is one definition instead of two.
- **`HEAD`'s `Governs:` is applied to all history.** Assumed 2026-10-22, attacca run, for
  cost: one read per spec instead of one per spec per landing. **Ceiling named:** a glob
  added later reaches back, and a renamed-away path drops out.
- **The spec-drift block is corpus-wide whatever `<change>` names.** Assumed 2026-10-22,
  attacca run: a capability is not a change, as the token block's totals are not.
- **Every block the metrics text prints from the report has a JSON key, and a test holds
  the mapping.** Assumed 2026-10-22, attacca run, replacing "whether it was mapped is
  review's to see, not a gate's" for the corrections, token, loop and spec-drift blocks.
  **Ceiling named:** a new block passes until its row is added to the test's table.

## Task breakdown

1. The `Governs:` reading in Go — `governsOf` sharing `land.go`'s fence rule, `governedBy`
   with bash `case` semantics — and the test running `spec-drift --trace` on a fixture to
   prove agreement.
2. The landing walk: first-parent log parsed into landings, the per-capability tally of
   landings and warned ones, co-claimed paths counted per claimant.
3. The text block and the `spec_drift` JSON key from one field on `metricsReport`, and the
   parity test over every block.
4. The in-flight delta rule and its `git grep`, within the stated call bound.
5. The born rule and the outstanding count.
6. The absent states: a capability with no `Governs:`, no capability specs, an unreadable
   history, and the change filter leaving the block whole.

## Verification criteria

`Governs:`, read as `spec-drift` reads it

- **When** `governsOf` reads a spec holding two `Governs:` lines, one indented, and a third
  inside a fence, it **shall** return the words of the first two, in order, and nothing from
  the fence.
  Proof: cmd/libretto/specdrift_test.go TestGovernsOfReadsLikeSpecDrift
- **When** `governedBy` is asked whether a path matches a glob, it **shall** answer as a
  bash `case` under `set -f` answers, over a table holding `*` across `/`, `?`, a bracket
  range, a `!`-negated bracket, a prefix sibling, an exact path and an empty glob.
  Proof: cmd/libretto/specdrift_test.go TestGovernsGlobReadsLikeTheShell
- **When** `spec-drift --trace` runs in a fixture repository whose spec carries a glob
  where `*` crosses `/`, a `?`, a negated bracket, an indented `Governs:` line and a fenced
  one, the set of paths it reports as orphans **shall** equal the set of tracked paths the
  Go reading leaves unclaimed.
  Proof: cmd/libretto/specdrift_test.go TestGovernsAgreesWithSpecDriftTrace

The walk

- **When** a branch of two commits is merged into the mainline, the walk **shall** count
  the merge as one landing.
  Proof: cmd/libretto/specdrift_test.go TestAMergedBranchIsOneLanding
- **Where** a landing touches a path a capability's `Governs:` claims, the walk **shall**
  count it in that capability's landings.
  Proof: cmd/libretto/specdrift_test.go TestSpecDriftCountsLandingsThatTouchGovernedCode
- **If** a landing touches only paths under the specs directory or `.agents/changes/`,
  **then** the walk **shall** count it in no capability's landings, whatever a glob claims.
  Proof: cmd/libretto/specdrift_test.go TestSpecsAndChangesAreNotGovernedCode
- **Where** a landing touches governed code and the capability's `spec.md`, the walk
  **shall** not count it as warned.
  Proof: cmd/libretto/specdrift_test.go TestALandingThatMovesTheSpecIsNotWarned
- **If** a landing touches governed code and moves neither the capability's `spec.md` nor a
  delta targeting it, **then** the walk **shall** count it as warned.
  Proof: cmd/libretto/specdrift_test.go TestALandingWithoutItsSpecIsWarned
- **Where** a landing touches governed code and a file in a change folder whose document
  declares `Targets:` naming the capability outside a fence, the walk **shall** not count
  it as warned.
  Proof: cmd/libretto/specdrift_test.go TestAnInFlightDeltaCountsAsTheSpecMoving
- **Where** two capabilities' globs claim the same path, the walk **shall** count a landing
  touching it for each claimant independently.
  Proof: cmd/libretto/specdrift_test.go TestACoClaimedPathCountsForEveryClaimant
- **While** no landing has yet touched a capability's `spec.md`, the walk **shall** count
  no landing for that capability.
  Proof: cmd/libretto/specdrift_test.go TestLandingsBeforeTheSpecExistedAreNotCounted
- **When** the walk ends, a capability's outstanding count **shall** equal its warned
  landings after the last landing that moved its `spec.md`.
  Proof: cmd/libretto/specdrift_test.go TestOutstandingCountsOnlySinceTheSpecLastMoved
- **When** the walk runs over a history in which no landing needs `Targets:`, it **shall**
  call git exactly once.
  Proof: cmd/libretto/specdrift_test.go TestSpecDriftCostsOneHistoryCall

The text block

- **Where** a capability spec exists, `metrics` **shall** print a spec-drift block after
  the loop block, one row per capability sorted by name, holding its landings, warned and
  outstanding counts.
  Proof: cmd/libretto/metrics_test.go TestMetricsPrintsTheSpecDriftBlock
- The spec-drift block **shall** print a line stating that a warned landing is one
  `spec-drift` would have warned on, not measured drift.
  Proof: cmd/libretto/metrics_test.go TestTheSpecDriftBlockSaysWarnedIsNotDrift
- **If** a capability spec declares no `Governs:`, **then** its row **shall** print a dash
  in every count cell.
  Proof: cmd/libretto/metrics_test.go TestACapabilityWithoutGovernsPrintsDashes
- **If** no capability spec exists, **then** `metrics` **shall** print the one line `no
  capability specs — spec-drift not replayed` in place of the block.
  Proof: cmd/libretto/metrics_test.go TestNoCapabilitySpecsPrintsOneLine
- **If** git cannot read the landing history, **then** `metrics` **shall** print
  `spec-drift history unreadable: <error>` in place of the block and exit 0.
  Proof: cmd/libretto/metrics_test.go TestAnUnreadableLandingHistoryIsReportedNotFatal
- **Where** a change is named, the spec-drift block **shall** print the same rows it
  prints with no change named.
  Proof: cmd/libretto/metrics_test.go TestTheSpecDriftBlockIgnoresTheChangeFilter

The JSON

- **When** `metrics --json` runs and a capability spec exists, the object **shall** carry
  `spec_drift.capabilities`, one object per row of the text block, in its order.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONSpecDriftListsTheBlockRows
- **Where** a capability declares `Governs:`, its `landings`, `warned` and `outstanding`
  **shall** equal the text block's cells on the same history, zero included.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONSpecDriftAgreesWithTheText
- **If** a capability spec declares no `Governs:`, **then** its object **shall** carry
  `name` and no other key.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONSpecDriftOmitsCountsWithoutGoverns
- **If** no capability spec exists, **then** the object **shall** omit `spec_drift`.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsSpecDriftWithoutSpecs
- **If** git cannot read the landing history, **then** `spec_drift` **shall** carry `error`
  in place of `capabilities`.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONSpecDriftUnreadableNamesTheError
- **When** one report in which every block prints is rendered as text and as JSON, each
  block header in the test's table that the text prints **shall** have its JSON key present
  in the document — the corrections, token, loop and spec-drift blocks. Removing the
  `spec_drift` mapping from `metricsDocOf` fails it.
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCarriesEveryTextBlock

## Approach

### Summary

A Go replay of `spec-drift`'s default-mode warning over every landing in `HEAD`'s
first-parent history, tallied per capability and rendered as one more `metrics` block and
one more JSON key, both from a single field on `metricsReport`. The ask's "spec moved after
the code" count was prototyped first, measured, and not built: most of what it counts is
the flow landing a delta after its code, and what is left is a spec edited later for an
unrelated reason. What ships is the measurement that proved exact — where the warning
fires — and the `Governs:` reading it rests on is mirrored from `spec-drift` and pinned
against it by a test that runs the script.

### Technical context

Go 1.26.5, stdlib only, exec'd git through `gitRunner` (`metrics.go`). Tests are standard
library; real-git tests skip under `-short` (`init_test.go`, `update_test.go`); a test
already runs bash (`label_workflow_test.go`). CI installs `rg` before `go test`
(`.github/workflows/gates.yml`). Gates: the six in `AGENTS.md`.

Already present and reused: `findSpecsDir` and `wikiDirOrder` (`wiki.go`), the fence rule
`landFenceRe` and the `Targets:` reader `headTargets` (`land.go`), `writeJSON`, `intPtr`
and the omitted-when-unknown convention (`json.go`), `trunc` (`metrics.go`).

Blast radius: two new files (`specdrift.go`, `specdrift_test.go`), four edited
(`metrics.go`, `json.go`, `land.go`, `metrics_test.go`), and `cli`'s spec at landing.

**The prototype**, a throwaway Python script over this repository at `dd6683d` (565
commits, 498 non-merge, 70 first-parent; 14 capability specs, all with `Governs:`):

| Reading | Units | Events ("spec moved after") | Catch-ups that are landing commits |
|---|---|---|---|
| per commit, a spec move only covers | 498 | 80 | — |
| per commit, a delta in a targeting folder also covers | 498 | 50 | — |
| per commit, + born rule, a co-claim covered by either spec | 498 | 47 | 37 |
| per landing, + born rule, each claimant counted | 70 | 13 | 0 |

The replay that ships — per landing, born rule, each claimant counted:

| Capability | Landings | Warned | Outstanding |
|---|---|---|---|
| agent-models | 4 | 0 | 0 |
| ci | 5 | 0 | 0 |
| cli | 33 | 10 | 0 |
| contributing | 2 | 0 | 0 |
| distribution | 1 | 0 | 0 |
| link-state | 3 | 1 | 0 |
| linking | 2 | 0 | 0 |
| ownership | 2 | 0 | 0 |
| panel | 7 | 0 | 0 |
| payload | 32 | 2 | 0 |
| readme | 18 | 12 | 2 |
| repo-sync | 4 | 0 | 0 |
| review-project | 3 | 2 | 2 |
| targets | 4 | 0 | 0 |
| **total** | **120** | **27** | **4** |

Ten warned landings read against their diffs:

| Landing | Capability | Governed path | Reading |
|---|---|---|---|
| #25 tag-main-on-merge | cli | `release_workflow_test.go` | tests `ci`'s workflow; `ci`'s spec moved in it |
| #31 rewrite-readme | cli | `readme_test.go` | co-claimed with `readme`, whose spec moved in it |
| #40 require-release-label | cli | `label_workflow_test.go` | tests `ci`'s workflow; `ci`'s spec moved in it |
| #69 dedupe-gates-runs | cli | `gates_test.go` | tests `ci`'s workflow; `ci`'s spec moved in it |
| #60 add-contributing-guide | readme | `readme_test.go` | the new test is `contributing`'s, whose spec moved in it |
| #35 agent-effort-level | payload | `agents/review-lens-design.md` | a `model:` value; `agent-models`' spec moved in it |
| #53 multi-tool targets | link-state | `state_test.go` | inside the contract: two more home overrides sandboxed |
| #74 loop budget | payload | `skills/record-work/spec-drift` | inside the contract: a SIGPIPE bug fix |
| #35 agent-effort-level | readme | `README.md` | inside the contract: the new `effort` verb's row |
| #63 add-design-phase | readme | `README.md` | arguable: the README's account of the stops changed, and `readme`'s spec caught up in #72 |

All 13 landing-level "moved after" events were read: the later spec edit was related to
the code before it once (#72 after #63) and twelve times was about something else.

Cost: the whole Python walk ran in 0.64 s wall, git included; `git log --first-parent
--name-only` alone in 0.17 s.

### The approach

1. **One fence-aware key reader.** Lift the loop in `headTargets` into
   `declared(text, key string) []string` — every line matching `^\s*<key>:` outside a
   fence, as `defenced` strips them, the remainder word-split. `headTargets` calls it with
   `Targets`; `governsOf(text)` is `declared(text, "Governs")`. One fence rule for both.
2. **`governedBy(path, glob string) bool`** translates the glob to an anchored regexp:
   `*` → `.*`, `?` → `.`, `[` … `]` copied as a class with a leading `!` turned to `^`, a
   backslash quoting the next rune, everything else `regexp.QuoteMeta`. An empty glob
   matches nothing. Compiled once per glob per run.
3. **Capabilities.** `findSpecsDir(root)`; every `<dir>/*/spec.md`, sorted, each with its
   repo-relative spec path and `governsOf` its content. A spec with no globs is kept and
   tallied as unknown.
4. **The walk.** One call:
   `git log --first-parent --reverse --name-only --format=C%x09%H` — oldest first, a
   merge's paths against its first parent (git 2.31 and later make `--first-parent` imply
   that diff). Each `C` record is a landing with its path list.
5. **Per landing, per capability**, in order: if the landing touches the spec path, the
   capability is born if it was not, and its `lastMoved` is this landing. Skip the
   capability while unborn. Governed code is the landing's paths outside the specs
   directory and `.agents/changes/` that any of its globs match; none → next capability.
   Count a landing. If the spec moved → not warned. Else, if the landing touched folders
   under `.agents/changes/`, ask once per folder per landing, cached:
   `git grep -h -e Targets: <sha> -- .agents/changes/<name>/`, read through `declared` —
   a folder whose document names the capability covers it. Else → warned, its index kept.
6. **Outstanding** is the count of warned indices after `lastMoved`.
7. **The report.** `specDriftReport{specsFound bool; err error; caps []capDrift}` on
   `metricsReport`, gathered in `gatherMetrics` after the loop, independent of `only`;
   `capDrift{name string; governed bool; landings, warned, outstanding int}`.
8. **Text.** `specDriftBlock(w, r)` after `loopBlock`: the no-specs line; or the
   unreadable line; or a header — `spec-drift`, `landings`, `warned`, `outstanding` — one
   row per capability, dashes for an ungoverned one, then the line `warned = a landing
   spec-drift would have warned on (its code moved, its spec and any delta targeting it
   did not) — where the warning fires, not measured drift; landings are HEAD's
   first-parent commits`. `flowCeiling` gains one bullet: whether a warned landing changed
   a promise.
9. **JSON.** In `json.go`:

        type specDriftCapJSON struct {
            Name        string `json:"name"`
            Landings    *int   `json:"landings,omitempty"`
            Warned      *int   `json:"warned,omitempty"`
            Outstanding *int   `json:"outstanding,omitempty"`
        }
        type specDriftJSON struct {
            Capabilities *[]specDriftCapJSON `json:"capabilities,omitempty"`
            Error        string              `json:"error,omitempty"`
        }

   and `SpecDrift *specDriftJSON` tagged `spec_drift,omitempty` on `metricsDoc`, set only
   when `specsFound`. The document then reads:

        "spec_drift": {
          "capabilities": [
            { "name": "cli", "landings": 33, "warned": 10, "outstanding": 0 },
            { "name": "notes" }
          ]
        }

   or `"spec_drift": { "error": "…" }`.
10. **At landing**, `cli`'s `metrics --json` example gains the object above and its key
    table gains:

    | Key | Meaning |
    |---|---|
    | `spec_drift` | *omitted when* the project has no capability spec; corpus-wide whatever `<change>` names |
    | `spec_drift.capabilities` | one object per row of the text block, by name; never empty, since a found spec is a row |
    | `landings`, `warned`, `outstanding` | the block's cells; *omitted when* the capability declares no `Governs:` |
    | `spec_drift.error` | present, in place of `capabilities`, only when git could not read the landing history |

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| **The ask as worded: count each spec edit that follows uncovered code** | Measured: 37 of 47 commit-level events are landings; 12 of 13 landing-level events are an unrelated edit. A number wrong that often gets muted |
| **B — a `--history` mode in `spec-drift`, exec'd by `metrics`** | One glob definition instead of a mirrored one, but `metrics` gains bash and `rg` as runtime needs, and the shell walk is a `case` per path per glob per landing |
| **C — ship nothing, record the measurement** | The replay is exact about what it counts, and its first run names an overlap worth fixing |
| **A co-claimed path covered when any claimant's spec moved** | Halves `cli`'s count, but disagrees with `spec-drift`'s commit-time warning, which this replays |
| **`path.Match` for the globs** | `*` stops at `/`; `cmd/*.go` would govern different files in the two tools |
| **`Governs:` as each landing saw it** | One read per spec per landing for a correction this history never needed |
| **Every commit as a landing** | Brings back the 37 landing-commit false positives |

### Risks

| Risk | What catches it |
|---|---|
| The Go glob drifts from `spec-drift`'s | `TestGovernsAgreesWithSpecDriftTrace` runs the script on a fixture holding every glob feature in use |
| The text grows the block and the JSON does not | `TestMetricsJSONCarriesEveryTextBlock`, forced red once by dropping the mapping |
| The number is read as "drift" | the block's own line and the ceiling bullet, the line with its own test |
| Run on a branch, its commits inflate the counts | named in the legend line; accepted |
| A long history | one log call, matching in memory; measured at 0.17 s for 70 landings; accepted |
| `git grep` calls multiply | bounded by landings touching code and a folder without the spec; cached per landing and folder |
| A git older than 2.31, where `--first-parent` does not imply the first-parent diff | `TestAMergedBranchIsOneLanding` runs real git; accepted below 2.31 |

### Validation and rollback

Gates: the six. New tests as named. **Forced red before believed:**
`TestGovernsAgreesWithSpecDriftTrace` with `governedBy` swapped for `path.Match` once — the
fixture path where `*` crosses `/` must diverge; `TestMetricsJSONCarriesEveryTextBlock`
with the `spec_drift` mapping removed once; `TestLandingsBeforeTheSpecExistedAreNotCounted`
with the born check removed once. The agreement test passes for the wrong reason if its
fixture holds no path where `*` crosses `/`, which is why the criterion names one.

Rollback: one revert; nothing persists, nothing migrates. The JSON key disappears, which
is a removed key — `!` on that revert's commit.

### Complexity deliberately kept

- **A glob translator instead of `path.Match`.** A few lines, because the stdlib's `*`
  differs from the shell's at `/`, and two tools disagreeing about ownership is the bug
  the anchors exist to prevent.
- **The in-flight delta rule**, though it moved no number here: it is `spec-drift`'s rule,
  and a replay without it would warn on every request of a change spanning several.
- **A Go test that runs a bash script.** The only proof two definitions agree is to ask
  both.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [ ] **The block, end to end.** `declared` lifted from `headTargets` (its existing tests
      untouched and green), `governsOf`, `governedBy`, the capability read, the
      first-parent walk with spec moves and co-claims, `specDriftReport` on
      `metricsReport`, the text block with its legend line and ceiling bullet, and
      `spec_drift` in the JSON with `landings` and `warned`. The agreement test forced red
      once with `path.Match`; the parity test forced red once without the mapping.
      Closes: TestGovernsOfReadsLikeSpecDrift, TestGovernsGlobReadsLikeTheShell,
      TestGovernsAgreesWithSpecDriftTrace, TestAMergedBranchIsOneLanding,
      TestSpecDriftCountsLandingsThatTouchGovernedCode,
      TestSpecsAndChangesAreNotGovernedCode, TestALandingThatMovesTheSpecIsNotWarned,
      TestALandingWithoutItsSpecIsWarned, TestACoClaimedPathCountsForEveryClaimant,
      TestSpecDriftCostsOneHistoryCall, TestMetricsPrintsTheSpecDriftBlock (landings and
      warned), TestTheSpecDriftBlockSaysWarnedIsNotDrift,
      TestMetricsJSONSpecDriftListsTheBlockRows, TestMetricsJSONSpecDriftAgreesWithTheText
      (landings and warned), TestMetricsJSONCarriesEveryTextBlock. Waits on: nothing.
- [ ] **A delta in flight covers its capability.** The per-landing, per-folder `git grep`
      read through `declared`, cached, inside the stated bound. Closes:
      TestAnInFlightDeltaCountsAsTheSpecMoving. Waits on: box 1.
- [ ] **Born, and outstanding.** The walk starts at a capability's first spec move; the
      `outstanding` column and JSON key added together; forced red once without the born
      check. Closes: TestLandingsBeforeTheSpecExistedAreNotCounted,
      TestOutstandingCountsOnlySinceTheSpecLastMoved, and the outstanding halves of
      TestMetricsPrintsTheSpecDriftBlock and TestMetricsJSONSpecDriftAgreesWithTheText.
      Waits on: box 1. Independent of box 2.
- [ ] **Absent is absent.** A spec with no `Governs:` prints dashes and maps to `name`
      alone; no capability specs prints its line and omits the key; an unreadable history
      prints its line, exits 0, and maps to `error`; a named change leaves the block whole.
      Closes: TestACapabilityWithoutGovernsPrintsDashes,
      TestNoCapabilitySpecsPrintsOneLine,
      TestAnUnreadableLandingHistoryIsReportedNotFatal,
      TestTheSpecDriftBlockIgnoresTheChangeFilter,
      TestMetricsJSONSpecDriftOmitsCountsWithoutGoverns,
      TestMetricsJSONOmitsSpecDriftWithoutSpecs,
      TestMetricsJSONSpecDriftUnreadableNamesTheError. Waits on: box 1. Independent of
      boxes 2 and 3.
- [ ] **At landing: `cli` carries it.** The criteria above under a spec-drift heading in
      `metrics`' section; the `metrics --json` example and key table gaining `spec_drift`
      as in *Approach*, step 10; the *Prior decisions* above; the ceiling sentence of "One
      gathered report per command, rendered twice" amended in place, no parallel entry.
      Same commit as the last code. Closes: every criterion's citation resolving under
      `spec-drift --anchors`. Waits on: boxes 1–4.
