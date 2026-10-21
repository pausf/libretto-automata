# Brief — adopt-organic-driven-development

## Conventions in use

- The payload is `skills/`, `agents/`, `commands/`; the flow's reasoning is `docs/FLOW.md`
  (not a spec, not installed). Skills never reference `scripts/` or `docs/`.
- The CLI lives in `cmd/libretto/`, one file per command (`loop.go`, `metrics.go`,
  `wiki.go`, `land.go`). Tests: standard library, table-driven, real temporary git repos
  for anything touching git.
- `spec-drift` (`skills/record-work/spec-drift`, bash) is payload: `--anchors` resolves
  `Proof:` citations, checks EARS in deltas, and `--retired` fails a landing that deletes
  `.agents/changes/*/plan.md` unless a capability spec's *Prior decisions* moved or the plan
  says `Durable decisions: none` (spec-drift:331-392, change name derived at :381).
- Deltas carry `Targets: <capability>`; `targets_of` reads `Targets:` lines from any file
  it is given, outside fences (spec-drift:122). Criteria are EARS with `Proof: <file> <Test>`.
- Layout readers today: `loop` checklist = `tasks.md`, fallback `plan.md`
  (loop.go:164-175, with a `ponytail:` naming when the fallback goes); `metrics` box churn
  from `tasks.md`/`plan.md` diffs; `wiki` in-flight = boxes in `tasks.md`/`plan.md`, queue =
  `Queued:` in `proposal.md` (wiki.go:834-857); `find-work` and `/libretto-status` scan boxes in
  `tasks.md`/`plan.md` and `Queued:` in `proposal.md`.
- Precedent for reading two layouts: the `plan.md` → `tasks.md` rename of 2026-08-17. Both
  names are read, the newer wins when both exist, and a `ponytail:` comment states when
  the older one can go.

## Constraints everyone inherits

- Today the payload spec promises **two stops with a contract, none without, plus phase
  8's question** (`.agents/specs/payload/spec.md` ~line 119), and the chain proposal →
  brief → spec delta (fresh `spec-writer`) → plan.md (fresh `plan-writer`) → tasks.md
  (fresh `task-cutter`). This change **reverses** those promises. Write the reversal
  explicitly. Do not soften it into an addition.
- Kept as-is, and not relitigated: one flow (no size flag, no second flow; the gear is
  knowable from the change, `payload/spec.md` ~160-164); the trivial lane and its test
  ("could anyone reasonably disagree about what done means"); the bug rule (criterion +
  red test before the fix); capability specs with `Governs:`/`Proof:`/EARS; `--anchors`;
  landing in one commit (code + delta applied + durable decisions retired + folder
  deleted); `review-work`'s fresh `work-reviewer`; `/libretto-attacca` answering every
  stop and question as `(assumed)`; phase 8's push question; `release:` label question.
- Nothing in flight is migrated. Existing changes in the old layout (12 queued proposals,
  `add-init-command`, `add-land-command`) must keep working in every reader.
- No new dependency. Go 1.26.5. No hardcoded version.
- Vocabulary rule (decision 6): the payload never says "ODD" or "Organic Driven
  Development". `docs/FLOW.md` may credit gentle-ai's Organic Driven Development as
  inspiration in one line, and nowhere else.

## Decisions already settled

Read `decisions.md` in this folder: six answers, verbatim. In short:
1. Substantial work gets **one change document** (intent, EARS criteria, approach,
   checklist). The living capability specs and the landing stay.
2. **One stop** before code: the whole document, reviewed once. Then only the push question.
3. **The session writes the document inline.** One fresh reviewer (`review-spec`) reads
   it before the stop. `spec-writer`, `plan-writer` and `task-cutter` stop being launched.
4. **Read both layouts.** New changes use the one document. Old changes keep being read.
5. **`proposal.md` grows.** The capture (queue or find-work) *is* the document, and
   criteria, approach and checklist are added to the same file.
6. **No acronym.** The method is named in Libretto's own vocabulary.

Derived from 1 + 3 + 5, open to challenge. If a writer finds it wrong, mark it:
- `brief.md` and `decisions.md` existed to feed writers that never had the conversation.
  When the session writes inline, that reason is gone. The expected shape is that
  decisions become a section of `proposal.md` and new changes write no brief. Treat this
  as a derived consequence, not as a user answer.
- The delta's `Targets:` line and its EARS criteria live inside `proposal.md`, so
  `--anchors`/EARS and `--retired` have to read `proposal.md` for new-layout changes.
  `--retired` must not fire when a *queued* proposal (an abandoned idea) is deleted.
- Known inconsistency to resolve, not preserve: find-work requires `proposal.md` on disk
  for every source-3 request, but the trivial lane is said to create no `changes/`
  directory (find-work:16-32 vs :132-134).

### Corrections after the first draft (2026-10-21)

The first two writers found these. They override anything above that contradicts them:

- **`land` exists only on `feat/add-land-command` (PR #71, unmerged).** It is not on this
  branch. Decision: #71 merges first and this change rebases onto it. The cli delta
  specifies that `land` reads `Targets:` from `proposal.md` (at HEAD) and from
  `spec-<cap>.md`, alongside old-layout deltas. The `land` behaviour to amend is the one
  in PR #71's `cmd/libretto/land.go` and its `cli` spec section; read them with
  `git show feat/add-land-command:<path>`.
- **No old-layout change folder exists in this tree.** All 12 changes on disk are
  `proposal.md`-only captures. `add-init-command` is in flight on an unpushed branch
  (`feat/add-init-command`, in the old layout there). Old-layout support is still required
  for those branches and for git history (`metrics`).
- **The checklist is the boxes under a level-2 `## Checklist` heading, and nothing
  else.** Only `## ` headings are seen by spec-drift's parsers, and boxes elsewhere would
  be read as criteria by `--ears`. The *Task breakdown* pillar is not boxes. Every reader
  (loop, metrics, wiki, find-work, /libretto-status, spec-drift default mode) applies the
  same definition.
- **Which file wins:** `proposal.md` counts as the checklist only when it holds at least
  one box under `## Checklist`. Otherwise fall back to `tasks.md`, then `plan.md`. Mere
  existence never wins, because every change in both layouts has a `proposal.md`.
- **`metrics`' old-layout queries are permanent** (git history never stops containing
  `plan.md`). The fallback `ponytail:` rule applies to `loop`/`wiki`, not to `metrics`.
- **The fresh reviewer also inherits `task-cutter`'s question.** It reads contract and
  checklist and reports, per box, what the document fails to answer.
- **Multi-capability changes:** `proposal.md` carries the primary capability's contract
  (`Targets:` + six pillars). Each extra capability gets a `spec-<cap>.md` beside it with
  its own `Targets:`. `write-spec`'s fan-out survives only as "one writer per extra
  `spec-<cap>.md`" where the context requires it.
- **The acronym rule scopes to the installed payload** (`skills/`, `agents/`,
  `commands/`). Verbatim decision quotes in specs may contain it.
- **`readme` is a third capability affected**: `readme/spec.md` ~251-253 and
  `cmd/libretto/readme_test.go TestWhatYouGetCarriesTheDiagrams` promise "the three stops
  — spec, plan and push".
- **`cmd/libretto/recommend.go:39` keys a model recommendation to `spec-writer`.** That is
  cli scope: the key follows the agent set (retired writers out, the new reviewer in).

## Vocabulary

- **change document**: `.agents/changes/<change>/proposal.md` in the new layout. It runs
  from capture to landing and holds the ask, the reading, the contract (six pillars, EARS
  criteria with `Proof:`, `Targets:`), the approach, `Durable decisions:`, the decisions
  log, and the checklist.
- **old layout**: a change folder with `spec.md` and/or `plan.md` and/or `tasks.md` (and
  maybe `brief.md`/`decisions.md`) beside `proposal.md`.
- **trivial lane**: step 0 said no contract is needed. No change document beyond the
  capture, and no stop except the push question.
- **contract lane**: everything else. One change document, one stop.
- **the stop**: the single pre-code go-ahead on the whole change document.
- **fresh reviewer**: `review-spec` launched with none of the conversation, reading the
  change document before the stop.
- **capability spec**: `.agents/specs/<capability>/spec.md`, the living contract.
- **landing**: the one commit of code + delta applied + durable decisions retired +
  change folder deleted.

## Six pillars

Every delta fills all six, with headings: Outcomes · Scope boundaries · Constraints ·
Prior decisions · Task breakdown · Verification criteria (EARS, each with a
`Proof: <file> <TestName>` naming a test that does not exist yet where needed).
