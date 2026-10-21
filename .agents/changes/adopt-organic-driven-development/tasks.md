# Tasks — adopt-organic-driven-development

Phase 6 runs these with `build-and-check`. Contract: spec-payload.md · spec-cli.md ·
spec-readme.md. Approach: plan.md.

Before box 1: build in a separate git worktree of this branch; the main checkout stays on `main`
with add-init-command's files untouched. All six gates run in the worktree before every commit,
after reading `git diff --cached --name-only`. Every test ships in the commit of the code it proves.
A skill whose instructions change bumps its own `metadata.version`.

## Layer 1 — readers learn both layouts

- [x] **1. `spec-drift` reads the change document** — payload unit 1.
  `--retired` also covers a deleted `proposal.md` whose `HEAD` copy carries `Targets:` outside
  fences; a change deleting both it and `plan.md` is reported once. Default mode ignores a staged
  `proposal.md` whose only difference is inside `## Checklist`. `criteria_of` skips `## Checklist`
  boxes. Self-test cases on real temporary repos through the existing `land()` helper; header
  layout tree updated.
  Closes: the six spec-payload `--self-test` criteria; each new case observed red first; removing
  `Targets:` from the contract fixture turns exit 1 into 0; the three existing `land` cases unedited.
  Waits on: nothing.
  Evidence: `spec-drift --self-test` 0, nine new cases (three red before the fix, the rest force-red), land cases unedited.

- [x] **2. Section counter and `loop` read the change document** — cli units 1-2.
  Section-scoped counter beside `countBoxes`, `planLine` the only box definition. `checklistPath`:
  `proposal.md` when its `## Checklist` holds a box, then `tasks.md`, then `plan.md`. Queued refusal
  before the PATH check; refusals stop saying `cut into tasks`; the prompt names `proposal.md` and
  `## Checklist`. `ponytail:` at loop.go:164; loop.go:22-25 and root.go:130 amended.
  Closes: C1-C2, L1-L5 — `TestOnlyTheChecklistSectionCounts`,
  `TestChecklistPathPrefersTheChangeDocument`, `TestChecklistPathFallsBackToPlan` (amended per L2),
  `TestLoopRefusesAQueuedProposal`, `TestTheLoopDrivesTheChangeDocumentsChecklist`,
  `TestLoopRefusalsDoNotNameTheCutter`; each red first; spec-cli's old-layout proofs unedited.
  Waits on: nothing.
  Evidence: five new tests red first, L2 amended as specified, old-layout proofs unedited; anchors 17 → 11.

- [x] **3. `metrics` reads churn from the change document** — cli unit 3.
  A third `git log` for `proposal.md`, first, whole-file `--unified=`; found only with an in-scope
  box; section tracked in pre- and post-image for `proposal.md` only; `fakeGit` refuses a
  `proposal.md` query without `--unified=`; `flowLegend`'s `—` line; plain comment replaces the
  `ponytail:` at metrics.go:172. `go run ./cmd/libretto metrics` before and after, compared.
  Closes: M1-M4 — `TestMetricsReadsChurnFromTheChangeDocument`,
  `TestAProposalWithoutAChecklistDoesNotShadowTheOldLayout`,
  `TestChangeDocumentChurnIsScopedToTheChecklist`, `TestTheReportExplainsItsColumns` (amended per
  M4); M2 observed failing under the naive order; no landed change's numbers move.
  Waits on: nothing.
  Evidence: M1-M3 red first, M2 red under the naive order; real-history metrics identical but the legend line; anchors 11 → 8.

- [x] **4. `wiki` in-flight strip uses the same rule** — cli unit 4.
  `readInFlight` through `checklistPath`, the section counter and `planLine`; `Queued:` untouched;
  comment at wiki.go:833-835.
  Closes: W1-W3 — `TestWikiInFlightReadsTheChangeDocument`,
  `TestWikiInFlightCountsBoxesAsTheLoopDoes` (W3 red first); the two existing wiki proofs unedited.
  Waits on: 2.
  Evidence: W1-W3 red first; both existing wiki proofs unedited; anchors 8 → 5.

- [x] **5. `find-work` and `/libretto-status` read both layouts; the trivial lane deletes its
  capture** — payload unit 2.
  Both scans and the three-row table apply the Outcomes 7 rule with the CLI's box definition;
  find-work 16-32 vs 132-134 reconciled; `write-spec` step 0 deletes the capture in the trivial
  lane's own commit and drops "a brief and a plan"; `/libretto-status` names `proposal.md`.
  Closes: the three `check-payload` wiring criteria, each row red first.
  Waits on: 2, 3, 4.
  Evidence: nine check-payload rows, eight red on today's prose, the delegation row forced red by hand; scan hand-tested on fences, nested and out-of-section boxes.

- [ ] **6. `land` proves the new layout** — cli unit 5. Tests only.
  Closes: `TestLandReadsTargetsFromTheChangeDocument`,
  `TestLandReadsTargetsFromAnExtraCapabilitySpec`, `TestLandOldLayoutTargetsAreUnchanged`, each
  forced red by excluding `proposal.md`/`spec-*.md` from `headTargets` and reverted; #71's own
  cases unedited.
  Waits on: PR #71 merged by the user, and this branch rebased onto it. Blocks only box 11.

## Layer 2 — the writers move, as one commit

- [ ] **7. One change document, one stop, one fresh reviewer** — payload units 3, 4, 5, 8 and cli
  unit 6, one commit: between any two of them the payload launches writers against inputs nobody
  writes any more.
  - `write-spec` grows `proposal.md` (`Targets:`, `## Decision log`, six level-2 pillars; creates it
    with `Tracker:` from a tracker); step 4b and *Resolving the markers* retired for the single
    case; step 2b becomes one `spec-writer` per extra `spec-<capability>.md`, only when context
    requires it, given `proposal.md`; *Where specs live* shows both layouts; its go-ahead removed.
  - `write-plan` writes `## Approach` and `Durable decisions:` into `proposal.md`; nothing launches
    `plan-writer`. `write-tasks` writes `## Checklist` last, inline; nothing launches
    `task-cutter`; its go-ahead removed. Both agent files deleted.
  - `review-spec` launches the new `agents/spec-reviewer.md` (`Read, Grep, Glob, Bash, Skill`,
    never edits), relays unedited, carries the per-box gap question, asks the one stop natively.
  - `agents/spec-writer.md` reads `proposal.md` by path; recommendation reason reworded to the
    fan-out role.
  - `record-work`, `review-work`, `work-reviewer.md` (incl. :54), `build-and-check` read either layout.
  - Stop tables: `libretto-flow`, `libretto-attacca`, `libretto-next`.
  - `check-payload`: retire 283, 288, 296, 324; reword 274, 275, 281, 287, 289, 335; keep 284, 285;
    293-294 become one `review-spec` row and two `check_absent` rows; a section-aware awk row for
    `spec-writer` outside the fan-out section; the `spec-writer`-reads-`brief.md`/`decisions.md` row.
  - `recommend.go`: `spec-reviewer` at sonnet/high, `spec-writer` kept. `docs/PAYLOAD.md` regenerated.
  Closes: every `check-payload` wiring criterion of payload units 3, 4, 5 and 8, new rows red first
  on today's tree; R1-R2 (`TestEveryRecommendationNamesAShippedAgent`), R1 forced red with a
  `plan-writer` key and reverted, R2 red before the agent file exists; the four existing
  `recommend_test.go` cases unedited.
  Waits on: 1, 5.

- [ ] **8. No method acronym in what installs** — payload unit 6.
  One `check-payload` row over `skills/`, `agents/`, `commands/`: the acronym case-sensitive whole
  word, the spelled-out name any case.
  Closes: that criterion, watched red on a hand-broken copy and reverted.
  Waits on: 7.

## Layer 3 — prose

- [ ] **9. Docs describe the result** — payload unit 7.
  `docs/FLOW.md` (phases 2-5, the stops, delegation, attacca, one gentle-ai credit line),
  `docs/STATE.md` 18 and 68, `AGENTS.md` 58, 112, 131, 159-163. If the rewrite drops
  `more expensive than this one`, the `TestMovedReasoningLandedInDocs` anchor is re-chosen here,
  under readme's two anchor rules — never deleted.
  Closes: six gates green, `TestMovedReasoningLandedInDocs` included.
  Waits on: 7.

- [ ] **10. README: one stop before code** — readme units 1-3.
  Amend `TestWhatYouGetCarriesTheDiagrams`, `TestReadmeWalksAFirstRun`; add
  `TestReadmeNamesNoMethodAcronym`; each red on today's README (the acronym one on a hand-broken
  copy). Rewrite the four places at net ≤ 0 lines.
  Closes: the spec-readme criteria — two decision nodes, the last naming `push`, ASCII fences; the
  walk names `/libretto-flow`, `change document`, `push`; no acronym; ≤ 380 lines, ceiling unedited;
  no moved-reasoning subject removed.
  Waits on: 9.

## Layer 4 — landing

- [ ] **11. Land: three deltas applied, this folder deleted** — one commit, `!` in the subject.
  Apply spec-payload, spec-cli, spec-readme onto their capability specs (incl. cli/spec.md 885-896,
  961-969 and readme Outcomes 3 and 7); retire the durable Prior decisions; regenerate the wiki
  with `go run ./cmd/libretto wiki`; delete `.agents/changes/adopt-organic-driven-development/`.
  Closes: six gates green, `--retired`'s `plan.md` arm satisfied; `libretto land` reports payload,
  cli and readme moved; every old-layout proof unedited.
  Waits on: 1-10.
