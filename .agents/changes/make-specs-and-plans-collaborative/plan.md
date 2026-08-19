# Plan: make specs and plans collaborative and delegated

Durable decisions: the five in the delta spec's Prior decisions

## Summary

Prose surgery on the payload, no Go. Phase 2's batched question call becomes a
one-at-a-time interview whose answers land verbatim in a new `decisions.md`; the spec
file's authorship moves from the orchestrator to `spec-writer` in every case, with
`[NEEDS CLARIFICATION]` markers replacing silent guessing. Phase 5 gains its first
question — the approach fork — and loses its authorship to a new `plan-writer` agent
cloned from `task-cutter`'s posture. Each new promise gets a `check_wiring` row in
`scripts/check-payload`, the mechanism that already proves this kind of claim.

## Technical context

Markdown payload only: two skills, two agents (one new), two commands, `docs/FLOW.md`,
plus `scripts/check-payload` (bash) and the generated `docs/SPEC.md` index. Gates: the
six in `AGENTS.md`; the ones that bite here are `scripts/check-payload` (frontmatter,
references, wiring rows, regenerated index — the new agent must appear in it via
`--index`) and `spec-drift --anchors`/`--ears` on the delta. No dependency moves, no
binary changes, `make build` untouched.

Blast radius, exhaustively: `skills/write-spec/SKILL.md`, `skills/write-plan/SKILL.md`,
`agents/spec-writer.md`, `agents/plan-writer.md` (new), `commands/libretto-flow.md`,
`commands/libretto-attacca.md`, `docs/FLOW.md`, `scripts/check-payload`, `docs/SPEC.md`
(regenerated), and at landing `.agents/specs/payload/spec.md`. Ten files, one new.

## The approach

Five moves, in dependency order:

1. **`decisions.md` contract into `write-spec`.** A new section defining: **the first
   write creates it** — normally phase 2's first answer, but a phase 5 opening on a
   change without one (resumed from before this change, or a phase 2 that asked
   nothing) creates it the same way; format `### Session YYYY-MM-DD` then
   `- Q: … → A: …` verbatim, `(assumed)` suffix under attacca; one writer, the
   orchestrator; entries copied into the delta's *Prior decisions* by the spec's
   author so landing needs no new mechanism and the log dies with the folder.
2. **Step 4 of `write-spec` rewritten** from one-call-with-everything to the interview:
   one question per `AskUserQuestion` call, recommendation with its reason first, "no
   more questions" always an option, soft bound around five with judgment both ways,
   each answer logged before the next question is asked. Zero questions stays legitimate
   and said in one line. Step 3 changes from "draft then hold the file" to "think the
   pillars, interview, then hand to the writer".
3. **`spec-writer` generalized.** Its launch contract gains the single-spec case (brief
   always written, even for one spec — same five headings, shorter), `decisions.md` as a
   named input beside the brief, and the marker rule: where brief and log do not settle
   a decision, write `[NEEDS CLARIFICATION: question]`, never a plausible guess. The
   orchestrator's side in `write-spec`: read the returned spec, ask each marker, log,
   replace the bracket expression only — anything beyond it is a relaunch.
4. **`write-plan` forks, then delegates.** Before any document: two or three approaches
   with tradeoffs, recommended first, one `AskUserQuestion`; chosen and rejected with
   why → `decisions.md`. Then one fresh `plan-writer` — new agent, front-matter tools
   `Read, Grep, Glob, Skill`, cloned from `task-cutter`'s you-know-nothing posture —
   gets the delta spec, the capability spec, `proposal.md` and `decisions.md` by path,
   returns the plan markdown plus what those inputs failed to answer; the orchestrator
   writes `plan.md`. Phase 5 still does not stop. **The document's shape reaches the
   agent the way it reaches `task-cutter`**: `plan-writer` invokes the `write-plan`
   skill for the pillars, and `write-plan` marks its orchestrator-only conduct — the
   fork and the delegation — as not the writer's to follow, the same split `write-tasks`
   already draws for its cutter.
5. **Routing, narrative, and wiring.** `commands/libretto-flow.md` §4 ("Phase 2 is the
   exception" becomes phases 2 and 5, each described), `docs/FLOW.md` same story plus
   the stop-owner drift fix (its stop table and the payload spec's both still name
   phase 5 as the second stop's owner; the 5→6 seam owns it), `libretto-attacca.md`
   marks both tranches as assumed entries in the log. New `check_wiring` rows — the
   map below — and `check-payload --index` regenerated for the new agent.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Section-by-section plan validation (superpowers-style) | Multiplies round trips; the specs here are short by design. Comes back if plans outgrow one read. |
| The interview run inside the spec-writer subagent | A subagent has no channel to the user; questions funnel through the orchestrator by standing rule. The conversation stays in the main context, the drafting leaves it. |
| Relaunch `spec-writer` to resolve markers | A full launch to replace a bracket expression is a round trip for a mechanical edit. User chose the inline patch, scoped to the expression. |
| No brief in the single-spec case (pass proposal + log directly) | Two prompt contracts for one agent, and the vocabulary section loses its home. User chose brief-always. |
| A hard cap of five questions (Spec Kit's rule verbatim) | Reintroduces the cap the user lifted on 2026-08-14. Soft bound with judgment keeps the bias-to-ask. |
| Status quo: phase 5 writes silently, inline | The rejected-alternatives pillar lives only in the conversation, which is why plans read as AI-authored — the finding that opened this change. |
| Markers in plans too | `plan-writer` returns gaps in its reply like `task-cutter`; a plan is not a contract, and two marker mechanisms is one too many. |

## Risks

| Risk | What catches it |
|---|---|
| The interview drifts back to one batched call — old habit in the skill's own history | `check_wiring` row on the literal "one question per call" string; conduct itself is uncheckable and the delta names that ceiling |
| Attacca stalls on the new phase 5 question | The fork is in attacca's pre-answered table: recommended option taken, logged `(assumed)`; wiring row on `libretto-attacca.md` |
| The marker exception erodes one-author-per-file | The exception is scoped to the bracket expression in both `write-spec` and `spec-writer`, stated on both sides of the seam |
| `decisions.md` and the delta's *Prior decisions* drift apart — two homes | The spec's copy is written from the log at authoring time and is the durable one; the log dies with the folder at landing. `review-work` reads both while the change is alive |
| Generalizing `spec-writer` breaks the fan-out contract | One prompt contract for both cases is the design; the fan-out path keeps its exact five-heading brief, and `check-payload` reachability covers the agent either way |
| `check_wiring` patterns are literal strings — a later rewording silently breaks the gate red, not green | Accepted: that is `check_wiring`'s existing design, and a red gate is the loud failure mode |

## Validation and rollback

The six gates, with `scripts/check-payload` carrying the new promises. The
criterion→row map, one row per delta criterion:

| Delta criterion | `check_wiring` row (file · pattern) |
|---|---|
| decisions.md created and verbatim | `skills/write-spec/SKILL.md` · `decisions.md` |
| One question per call, recommendation, "no more" | `skills/write-spec/SKILL.md` · `one question per call` |
| spec-writer authors the single-spec case | `skills/write-spec/SKILL.md` · `single-spec` (near `spec-writer`) |
| Markers instead of guessing; patch scoped | `agents/spec-writer.md` · `NEEDS CLARIFICATION` and `skills/write-spec/SKILL.md` · `NEEDS CLARIFICATION` |
| Phase 5 fork, chosen+rejected logged | `skills/write-plan/SKILL.md` · `approaches` / `decisions.md` |
| plan-writer drafts, orchestrator writes | `skills/write-plan/SKILL.md` · `plan-writer` (agent file's existence and frontmatter already checked by the agents section) |
| Attacca logs assumed entries | `commands/libretto-attacca.md` · `decisions.md` |
| Three stops | already proven by the existing stop-table wiring; at landing this criterion merges with the payload spec's copy, never duplicates it |

**Forced red before believed**: every new `check_wiring` row gets its pattern broken on
purpose once (edit the pattern to a nonsense string, watch `FAIL`, revert) — two rows in
this repository's history passed green for the wrong reason, and a wiring row that has
never been red proves only that grep runs.

Rollback: one revert. Nothing migrates — `decisions.md` files exist only in future
changes' folders, created by the skills at run time.

## Complexity deliberately kept

- **A brief for a single spec.** More ceremony than passing two paths, kept because one
  prompt contract for `spec-writer` is what keeps the agent auditable and re-runnable in
  both cases, and the vocabulary heading needs a home even at N=1.
- **`decisions.md` as its own file** rather than a section of `proposal.md`: the writers
  take inputs by path, and a log with one writer cannot share a file with a document the
  orchestrator also edits for other reasons.
- **The one-at-a-time interview** costs a round trip per question where one call cost
  one. That cost is the feature: each answer can redirect the next question, which is
  the difference between an interview and a form.
