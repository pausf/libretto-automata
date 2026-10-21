Targets: cli

# CLI — reading the change document beside the old layout

In the new layout each change has one change document,
`.agents/changes/<change>/proposal.md`. It starts as the capture and grows to hold the
contract, the approach, `Durable decisions:` and, last, a `## Checklist` section. Today four
places in `cmd/libretto/` assume the old layout, and none of them reads that section:

- `loop`: `checklistPath`, loop.go:165-176
- `metrics`: `measure`, metrics.go:173-180
- `wiki`: `readInFlight`, wiki.go:836-882
- `land`: PR #71, which this change rebases onto

There is also the recommendation table in `recommend.go`, which must name only agents the
payload ships: it gains the new reviewer's key and keeps `spec-writer`, which survives for the
fan-out. This delta teaches all of them the new layout without breaking the old one.

## Outcomes

1. `libretto loop <change>` drives a new-layout change by the boxes under `## Checklist` in its
   `proposal.md`. The session prompt names `proposal.md` and that section.
2. `libretto loop <change>` refuses a queued proposal (a capture with no `## Checklist` boxes,
   no `tasks.md` and no `plan.md`) as having no checklist. The refusal comes before the PATH
   check for `claude`.
3. `libretto loop` on an old-layout change behaves exactly as it does today. That covers a
   change with `tasks.md` or `plan.md` whose `proposal.md` has no `## Checklist` boxes.
4. `libretto metrics` reads box churn for a new-layout change from the `## Checklist` lines in
   `proposal.md`'s history. For an old-layout change it reads `tasks.md`, then `plan.md`,
   counting every box as today. Each change is read from one file. Counts from two files are
   never added together.
5. `libretto wiki` lists a new-layout change with open checklist boxes in the in-flight strip,
   counted from `proposal.md`. A queued proposal stays in the queue and out of the strip.
6. `libretto land` takes a change's capabilities from three sources: `Targets:` in
   `proposal.md` at `HEAD`, `Targets:` in every `spec-<cap>.md` beside it at `HEAD`, and the
   old-layout delta files it reads today.
7. `libretto models` recommends a model for the new `spec-reviewer` and keeps `spec-writer`'s
   recommendation (sonnet/high), because `spec-writer` survives for the fan-out. The table has
   no `plan-writer` or `task-cutter` key. Every agent the recommendation table names is an
   agent the payload ships.
8. The CLI uses one rule for "which file holds this change's checklist", one rule for "which
   lines of a change document are the checklist", and one rule for "what a box is".

## Scope boundaries

**In:**

- `checklistPath`, the refusals in `loop`/`runLoop`, the prompt, and the header comment in
  `cmd/libretto/loop.go`; the `loop` comment at `cmd/libretto/root.go:130`
- a section-scoped box counter for the change document, beside `countBoxes`
- the churn queries in `measure`, section scoping in `churn` for `proposal.md` diffs, the `—`
  line of `flowLegend`, and the `ponytail:` at metrics.go:172 in `cmd/libretto/metrics.go`
- `readInFlight` and its comment in `cmd/libretto/wiki.go`
- `land`'s target discovery in `cmd/libretto/land.go` (as merged from PR #71)
- the `recommendations` table in `cmd/libretto/recommend.go`
- tests in `loop_test.go`, `metrics_test.go`, `wiki_test.go`, `land_test.go`,
  `recommend_test.go`
- the `ponytail:` comments that say when each fallback can be removed

**Out (non-goals):**

- **skills, agents, commands, `spec-drift`, `docs/FLOW.md`.** The sibling `payload` delta owns
  them. That includes `find-work`'s and `/libretto-status`'s scans, `--anchors`/`--retired`
  reading `proposal.md`, the shape of the change document, and creating or deleting any file
  under `agents/`.
- **`cmd/libretto/readme_test.go`** and the README. The sibling `readme` delta owns them.
- **migrating any change.** No reader rewrites a file. The queued proposals on disk, and
  `add-init-command` on its own branch in the old layout, stay as they are.
- **removing the `tasks.md`/`plan.md` fallbacks in `loop` and `wiki`.** They stay until the
  condition in their `ponytail:` comment holds. `metrics`' old-layout queries are permanent.
- **parsing task text.** The rule from the loop's existing criterion stands: only the box is
  read. The only addition is the section boundary that decides which boxes are in scope.
- **what `Queued:` means.** `wiki` keeps reading it exactly as today (decision 5).
- **any `land` behaviour other than where it finds `Targets:`.**
- **`cmd/libretto/init.go` / `init_test.go`.** These sit uncommitted in this tree, but they
  belong to `add-init-command`.

## Constraints

- Go 1.26.5, standard library only, no new dependency. Tests use the standard library and are
  table-driven where the case is a table. `metrics` tests go through the existing `fakeGit`
  seam (metrics_test.go:17-65). That seam's `-p` branch refuses a query that lacks
  `--full-history` or `--diff-filter=AM`, so every new query has to pass both.
- **What counts as the checklist in a change document.** It is the lines of `proposal.md` that
  meet all three conditions:
  - outside fenced code blocks (a line matching `^\s*```` opens or closes a fence)
  - after a line matching `^##\s+Checklist\s*$`
  - before the next line matching `^##\s`

  `###` does not close the section and neither does `#`. These boundaries are exactly
  spec-drift's: `defenced` at spec-drift:119 drops fences, and its section parsers open and
  close only on `^##[[:space:]]`. If `## Checklist` appears more than once, every such section
  counts. Boxes anywhere else in `proposal.md` (a quoted ask, a *Task breakdown* written as
  boxes) are not the checklist.
- **Old-layout files count every box.** `tasks.md` and `plan.md` are read whole by
  `countBoxes`, as today.
- **One definition of a box.** It is `planLine` (loop.go:31). `metrics` already delegates to
  it through `boxIn` (metrics.go:254). `wiki`'s `readInFlight` does not. It matches
  `strings.HasPrefix(t, "- [ ]")` and `"- [x]"` (wiki.go:864-868), so `* [ ]` and `[X]` count
  in `loop` but not in `wiki`, even though its comment claims "the loop's own box definition".
  This delta closes that gap.
- **Which file wins.** `proposal.md` is the checklist only when its `## Checklist` section
  holds at least one box. Otherwise the order is `tasks.md`, then `plan.md`. The existence of
  `proposal.md` is never the signal, because every change in both layouts has one. If it were,
  every old change's `proposal.md` of prose would hide its `tasks.md`/`plan.md` and blank its
  loop, its strip entry and its churn. This follows the 2026-08-17 `plan.md` → `tasks.md`
  precedent: the newer layout wins when both are present.
- **`metrics` reads history, not disk, and its first non-empty diff is not a signal.** Every
  change's `proposal.md` has a non-empty `--diff-filter=AM` history, because the capture adds
  it. A query therefore counts as found only when its diff carries at least one box line in
  scope: under `## Checklist` for `proposal.md`, anywhere for `tasks.md`/`plan.md`. Today's
  test is `strings.TrimSpace(diff) != ""` (metrics.go:175). Put `proposal.md` first under that
  test and it would set `planSeen` and blank the churn of every old-layout change.
- **A diff hunk does not carry its section heading.** With git's default 3 lines of context,
  a box line far below `## Checklist` arrives without the heading, so section scope cannot be
  decided from a default patch. The `proposal.md` query asks for enough context that each
  commit's diff carries the whole file. `churn` then tracks the section separately in each
  image:
  - the post-image is context lines and `+` lines; the pre-image is context lines and `-`
    lines
  - a `+` box counts when the post-image section is `## Checklist`
  - a `-` box counts when the pre-image section is `## Checklist`
  - the section state resets at each commit SHA line

  Named limit: a commit that renames a heading so that untouched boxes move into or out of
  `## Checklist` changes no `+`/`-` box line, so it moves neither the churn nor the
  denominator.
- `metrics` keeps one `git log` call per file name. One pathspec covering several names is
  never used (existing Prior decision). This delta adds a third call ahead of the two.
- Error-message prefixes that existing tests assert stay as they are: `no checklist at` and
  `no task boxes`.
- **`land` reads at `HEAD`** for the new layout. The landing commit deletes the change folder
  in the index, so the working tree and the index no longer hold `proposal.md` or
  `spec-<cap>.md` when `land` needs their `Targets:`. `Targets:` lines are read outside fences,
  as spec-drift's `targets_of` does (spec-drift:122).
- The CLI does not own the agent set. A recommendation key is valid only when the payload
  ships `agents/<key>.md`. Adding the new reviewer's key therefore lands in the same commit as
  the sibling's agent-file changes.

## Prior decisions

From the decision log, verbatim, the user, 2026-10-21:

- Q: for substantial work, what replaces the proposal → brief → spec → plan → tasks chain? → A: Un doc de cambio + specs vivos (Recommended) — one document per change (intent, EARS criteria, approach, checklist); capability specs with Governs:/Proof: and the landing stay
- Q: who writes the change document, and where does the fresh eye stay? → A: Sesión escribe, revisor fresco lee (Recommended) — the session writes the document inline; one fresh reviewer (review-spec) reads it before the stop; spec-writer, plan-writer and task-cutter launches go from three to one review
- Q: what happens to changes and tooling that read proposal/plan/tasks today? → A: Leer ambos formatos (Recommended) — new changes are born as one document; loop/metrics/wiki/find-work keep reading the old layout for what already exists; nothing in flight is touched
- Q: is the change document the same file growing from the capture, or a new one? → A: proposal.md crece (Recommended) — the queue/find-work capture is the document; criteria, approach and checklist are added to the same file; Queued: and proposal.md scans do not change
- Q: PR #71 (libretto land) reads Targets: from old deltas; how are the two changes ordered? → A: #71 primero, este se adapta (Recommended) — #71 merges as is; this change rebases on it and its cli delta makes land read Targets: from proposal.md
- Q: how does a change touching two or more capabilities carry its contract? → A: proposal.md + spec-<cap>.md extra (Recommended) — proposal.md carries the primary capability contract; each extra capability gets a spec-<cap>.md beside it with its own Targets:

**Derived, not a user answer.** Recorded by the orchestrator, 2026-10-21, and marked as such in
the log:

- Q: which boxes in proposal.md count as the checklist? → A: only boxes under a level-2 `## Checklist` heading — settled by the code, not asked: spec-drift parsers react to `## ` only, and boxes elsewhere would be read as criteria by --ears (derived, orchestrator)

The brief records two more things as corrections, not as user answers, and this delta does
not relitigate them:

- `metrics`' old-layout queries are permanent.
- The fallback `ponytail:` rule applies to `loop`/`wiki` and not to `metrics`.

These are already in `cli`'s *Prior decisions* and this delta does not reopen them:

- `metrics` reads churn with separate `git log` calls, not one pathspec covering several
  names, and the legacy path can never go.
- The newer checklist name wins when both are present (`TestChecklistPathFallsBackToPlan`).
- The commands read a checkbox identically (`TestBoxInReadsOnlyRealCheckboxes`). This delta
  extends that from two commands to three.

This delta records the following, to be retired into `cli`'s *Prior decisions* at landing:

- **`proposal.md` is the checklist when, and only when, its `## Checklist` section holds a
  box.** Existence cannot be the signal, because every change in both layouts has a
  `proposal.md`.
- **The `proposal.md` churn query carries whole-file context.** Section scope is undecidable
  from a default 3-line hunk.

## Task breakdown

1. **Section counter.** Add one function that counts `planLine` boxes in a reader, restricted
   to the change document's checklist as defined under *Constraints*. `countBoxes` stays the
   whole-file counter. Tests: C1-C2.
2. **`loop`: the change document wins.**
   - **Which file.** `checklistPath` returns `proposal.md` when its `## Checklist` holds a
     box. Otherwise it returns `tasks.md`, then `plan.md`, as today. When none of the three
     holds the checklist and neither `tasks.md` nor `plan.md` exists, it returns
     `proposal.md`, so the refusal names the file the flow now writes.
   - **Counting.** `runLoop` counts with the section counter when the selected file is
     `proposal.md`, and with `countBoxes` otherwise. This applies both before and after each
     round.
   - **Refusal order.** `loop` refuses a change whose selected file holds no in-scope box
     *before* the PATH check. A queued proposal is then named as having no checklist rather
     than as missing `claude`.
   - **Refusal wording.** The refusals stop naming the cutting step ("cut into tasks"),
     because no agent cuts the checklist any more. The 5→6 seam stays, as the unnumbered place
     the checklist is written. They say instead that the loop drives a change whose checklist
     exists.
   - **Prompt.** When the file is `proposal.md`, the prompt tells the session to take the
     first unchecked box under `## Checklist`. It keeps the substring `FIRST unchecked box`.
   - **Comments.** Replace the `ponytail:` at loop.go:164 with one saying the `tasks.md` and
     `plan.md` fallbacks go when no `.agents/changes/*/tasks.md` or `.agents/changes/*/plan.md`
     anywhere holds a checkbox. Amend loop.go:22-25 and root.go:130, which name `tasks.md` as
     the state.

   Tests: L1-L5.
3. **`metrics`: three queries, first with an in-scope box wins.**
   - **Query order.** In `measure`, query `proposal.md` first, passing a `--unified=<n>` wide
     enough to print the whole file, then `tasks.md`, then `plan.md`.
   - **Found.** A query is found only when its diff yields at least one in-scope box line.
     Otherwise it does not set `planSeen` and the search continues.
   - **Churn.** `churn` applies per-image section tracking to `proposal.md` diffs only.
   - **`fakeGit`.** Teach it to answer `<change>/proposal.md` keys, and to refuse a
     `proposal.md` query that lacks `--unified=`, the way it already refuses a missing
     `--diff-filter=AM`.
   - **Legend.** Update the `—` line of `flowLegend` to say the change has no checklist in its
     history, naming the three files.
   - **Comment.** Correct the `ponytail:` at metrics.go:172. It says the second call goes;
     replace it with a plain comment saying the old-layout queries are permanent.
   - **Real-history check.** Before landing, run `libretto metrics` against this repository's
     real history before and after the change. Confirm that no landed change's
     `closed`/`reopen` cells moved, and report the comparison.

   Tests: M1-M4.
4. **`wiki`: in flight by the same rule.** `readInFlight` picks each change's file with
   `checklistPath`. It counts `proposal.md` with the section counter and the other files with
   `countBoxes`. The `Queued:` scan is untouched. Update the comment at wiki.go:833-835.
   Tests: W1-W3.
5. **`land`: `Targets:` from the new layout.** After rebasing onto PR #71, `land` takes the
   change's capabilities from three sources: the `Targets:` lines of `proposal.md` at `HEAD`,
   the `Targets:` lines of every `spec-*.md` in the change folder at `HEAD`, and the old-layout
   delta files it reads today. A `proposal.md` with no `Targets:` line (a queued or
   trivial-lane capture) contributes no capability. Tests: N1-N3.
   Resolved 2026-10-21 by reading `git show feat/add-land-command:cmd/libretto/land.go`: `headTargets` (land.go:293-328) already reads `Targets:` lines, fences stripped, from **every** `.md` file the folder holds at `HEAD`, never from the index or working tree. `proposal.md` and `spec-<cap>.md` are therefore already covered, and this task adds **no code to `land`**: only the three tests N1-N3 that prove the new layout through `land`'s report, the `part 2: <cap> spec moved` / `did not move` lines and the exit code.
6. **`recommend`: the key follows the agent set.**
   - Keep the `spec-writer` entry in `recommendations`, sonnet/high as today: the agent
     survives for the fan-out. `plan-writer` and `task-cutter` have no entry today, and none
     is added.
   - Add an entry for `spec-reviewer`, the fresh reviewer the `payload` sibling adds as
     `agents/spec-reviewer.md`, read from its draft spec.
   - Answered 2026-10-21 by the user: "sonnet/high, como work-reviewer" — `spec-reviewer` is recommended sonnet at high effort, with work-reviewer's reason: a cheap miss here is a false green.
   - Add a test that every key names a shipped agent file. This task depends on the sibling's
     `agents/` changes landing in the same commit.

   Tests: R1-R2.

## Verification criteria

### The change document's checklist

- **C1. When counting a change document, the section counter shall count only `planLine`
  boxes that lie after a `## Checklist` line, before the next line matching `^##\s`, and
  outside fenced code blocks. A box under `## Task breakdown`, a box in a fence inside
  `## Checklist`, and a box under a `### ` subheading inside `## Checklist` shall count as
  zero, zero and one respectively.**
  Proof: cmd/libretto/loop_test.go TestOnlyTheChecklistSectionCounts
- **C2. Where a change document has no `## Checklist` heading, the section counter shall
  return zero open and zero done, whatever boxes the file holds elsewhere.**
  Proof: cmd/libretto/loop_test.go TestOnlyTheChecklistSectionCounts

### `loop`

- **L1. Where a change's `proposal.md` holds at least one box under `## Checklist`,
  `checklistPath` shall return that `proposal.md`, whether or not `tasks.md` or `plan.md`
  also exist in the folder.**
  Proof: cmd/libretto/loop_test.go TestChecklistPathPrefersTheChangeDocument
- **L2. Where a change's `proposal.md` holds no box under `## Checklist`, `checklistPath`
  shall return `tasks.md` when it exists and otherwise `plan.md` when it exists. Where
  neither exists, it shall return `proposal.md`, including when `proposal.md` holds boxes
  outside `## Checklist`.** The first two cases are today's behaviour. The third replaces
  today's "names `tasks.md`" assertion: the criterion it amends says "names the current one",
  and the current one has moved. This is **one of the two existing assertions this delta
  changes.**
  Proof: cmd/libretto/loop_test.go TestChecklistPathFallsBackToPlan
- **L3. If a change holds a `proposal.md` with no box under `## Checklist` and has neither
  `tasks.md` nor `plan.md`, then `loop` shall return an error containing `no task boxes` and
  the `proposal.md` path while PATH is empty, and shall start no session.** The empty PATH
  proves that the refusal comes before the dependency check.
  Proof: cmd/libretto/loop_test.go TestLoopRefusesAQueuedProposal
- **L4. While `proposal.md` is the selected checklist, the loop shall count progress from
  `## Checklist` boxes only. The prompt printed by `--dry-run` and the prompt handed to each
  session shall name `.agents/changes/<change>/proposal.md` and `## Checklist`. A box closed
  outside the section shall not count as progress.**
  Proof: cmd/libretto/loop_test.go TestTheLoopDrivesTheChangeDocumentsChecklist
- **L5. The `no checklist at` and `no task boxes` refusals shall not contain `cut into
  tasks`.**
  Proof: cmd/libretto/loop_test.go TestLoopRefusalsDoNotNameTheCutter

### `metrics`

- **M1. Where a change's `proposal.md` history (`--full-history --diff-filter=AM -p`) carries
  at least one box line under `## Checklist`, `measure` shall report `checked`, `uncheck` and
  `boxes` from that history alone, even when the `tasks.md` history for the same change also
  carries boxes.** The fixture gives the two files different numbers, so a sum or the wrong
  file fails.
  Proof: cmd/libretto/metrics_test.go TestMetricsReadsChurnFromTheChangeDocument
- **M2. Where a change's `proposal.md` history is non-empty but carries no box line under
  `## Checklist`, `measure` shall report churn from `tasks.md`, or from `plan.md` when
  `tasks.md` has none. Where neither carries a box, it shall leave `planSeen` false so the
  row prints a dash.** The fixture includes a `proposal.md` diff whose only boxes sit outside
  `## Checklist`.
  Proof: cmd/libretto/metrics_test.go TestAProposalWithoutAChecklistDoesNotShadowTheOldLayout
- **M3. When `measure` queries `proposal.md` history, it shall pass `--unified=`, and `churn`
  shall count a `+` box only when its post-image section is `## Checklist` and a `-` box only
  when its pre-image section is `## Checklist`.** The fixture's diff closes a box 20 lines
  below the heading in one commit and adds a box under `## Task breakdown` in another; only
  the first counts. `fakeGit` refuses a `proposal.md` query without `--unified=`.
  Proof: cmd/libretto/metrics_test.go TestChangeDocumentChurnIsScopedToTheChecklist
- **M4. The legend's `—` line shall state that the change has no checklist in its history and
  name `proposal.md`, `tasks.md` and `plan.md`.** This amends the existing assertion on
  `"no plan.md in its history"`, which has been stale since `tasks.md` (2026-08-17). This is
  **the second of the two existing assertions this delta changes.**
  Proof: cmd/libretto/metrics_test.go TestTheReportExplainsItsColumns

### `wiki`

- **W1. Where a change's `proposal.md` holds an open box under `## Checklist`, the in-flight
  strip shall name that change with its closed/total count read from that section, including
  when the folder also holds a `tasks.md` with different counts and `proposal.md` holds boxes
  outside the section.**
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightReadsTheChangeDocument
- **W2. Where a change's `proposal.md` carries a `Queued:` line and no box under
  `## Checklist`, the wiki shall count it in the queue and shall not name it in the in-flight
  strip.**
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightReadsTheChangeDocument
- **W3. The in-flight strip shall count boxes with the loop's definition: a `* [ ]` line and a
  `- [X]` line in a `tasks.md` shall count as one open box and one closed box.**
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightCountsBoxesAsTheLoopDoes

### `land`

- **N1. When a change's `proposal.md` at `HEAD` carries `Targets: <cap>` outside fences,
  `land` shall treat `<cap>` as one of the change's capabilities, including when the index
  has deleted the change folder.**
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromTheChangeDocument
- **N2. When a change folder at `HEAD` holds a `spec-<cap>.md` carrying `Targets: <cap>`,
  `land` shall treat `<cap>` as one of the change's capabilities alongside those named by
  `proposal.md`.**
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromAnExtraCapabilitySpec
- **N3. If a change's `proposal.md` carries no `Targets:` line outside fences, then `land`
  shall derive no capability from it, and `land` on an old-layout change shall find the same
  capabilities it finds before this delta.**
  Proof: cmd/libretto/land_test.go TestLandOldLayoutTargetsAreUnchanged

### `recommend`

- **R1. The recommendation table shall have no `plan-writer` or `task-cutter` key.**
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationNamesAShippedAgent
- **R2. For every key in the recommendation table, the repository shall hold
  `agents/<key>.md`.** The file is read from `../../agents`, as `gates_test.go` and
  `readme_test.go` already read repository files.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationNamesAShippedAgent

The existing proofs that cover old-layout behaviour stay green with no edits: `TestCountBoxesReadsOnlyTheBox`, `TestLoopRefusesAPlanThatIsNotThere`, `TestLoopRefusesAPlanWithNoBoxes`, `TestTheMissingPlanIsReportedBeforeTheMissingDependency`, `TestTheLoopPromptRoutesAndNeverAuthorises`, `TestMetricsFallsBackToLegacyPlan`, `TestAChangeWithNoPlanReportsADashNotAZero`, `TestBoxInReadsOnlyRealCheckboxes`, `TestWikiHomeCarriesTheInFlightStrip`, `TestWikiFlowBoardListsTheQueue`, the four existing `recommend_test.go` cases, and PR #71's `land_test.go` cases. They are listed in prose, not bullets, because a bullet in this section is read as a criterion.

The only existing assertions this delta changes are the two named in L2 and M4, each because
the promise it checks moved. The capability spec's text moves with them in three places: the
`loop` "owns no state" and "reads a change created before the checklist was renamed" bullets
(cli/spec.md:885-896), and the `metrics` "churn is read from `tasks.md`, then from `plan.md`"
bullet (cli/spec.md:961-969).
