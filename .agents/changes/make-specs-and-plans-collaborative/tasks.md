# Tasks: make specs and plans collaborative and delegated

Execution: `build-and-check`, one box per session. Every box lands its prose change
and its `check_wiring` row(s) in the same commit, and every new row is forced red
once (pattern edited to nonsense, FAIL observed, reverted) before it is believed —
the plan mandates this per row.

- [x] **1. The `decisions.md` contract in `write-spec`.**
  A new section: the first write creates it — normally phase 2's first answer, and a
  phase 5 opening on a change without one creates it the same way; format
  `### Session YYYY-MM-DD` then `- Q: … → A: …` verbatim, `(assumed)` suffix under
  attacca; one writer, the orchestrator; entries copied into the delta's *Prior
  decisions* by the spec's author.
  Closes: delta criterion "decisions.md created and verbatim" — `check_wiring` row
  `skills/write-spec/SKILL.md · decisions.md`, forced red once; all six gates green.
  Waits on: nothing.

- [x] **2. The phase 2 interview in `write-spec`.**
  Step 4 rewritten from one batched call to one question per `AskUserQuestion` call,
  recommendation with its reason first, "no more questions" always an option, soft
  bound around five with judgment both ways, each answer logged before the next
  question. Zero questions stays legitimate in one line. Step 3 becomes "think the
  pillars, interview, then hand to the writer".
  Closes: delta criterion "one question per call, recommendation, no-more" —
  `check_wiring` row `skills/write-spec/SKILL.md · one question per call`, forced red
  once; all six gates green.
  Waits on: 1 (answers land in the log the contract defines).

- [x] **3. `spec-writer` generalized, and the marker loop in `write-spec`.**
  One box, both sides of the seam: `spec-writer`'s launch contract gains the
  single-spec case (brief always, same five headings, shorter), `decisions.md` as a
  named input beside the brief, the tier picked before launch as in the fan-out, and
  the rule "where brief and log do not settle a decision, write
  `[NEEDS CLARIFICATION: question]`, never a guess". `write-spec`'s side: read the
  returned spec, ask each marker, log the answer, replace the bracket expression
  only — anything beyond it is a relaunch. Under attacca, markers take the
  recommended answer, logged `(assumed)`. Stated on both sides, per the plan's risk
  table.
  Closes: delta criteria "spec-writer authors the single-spec case" and "markers
  instead of guessing; patch scoped" — `check_wiring` rows
  `skills/write-spec/SKILL.md · single-spec`,
  `agents/spec-writer.md · NEEDS CLARIFICATION`,
  `skills/write-spec/SKILL.md · NEEDS CLARIFICATION`, each forced red once; all six
  gates green.
  Waits on: 1 (the log is a named input), 2 (same file, the handoff it rewrites).

- [ ] **4. The fork in `write-plan`.**
  Before any document: two or three approaches with tradeoffs, recommended first, one
  `AskUserQuestion`; chosen and rejected with why → `decisions.md`. Phase 5 still does
  not stop.
  Closes: delta criterion "phase 5 fork, chosen+rejected logged" — `check_wiring` rows
  `skills/write-plan/SKILL.md · approaches` and `· decisions.md`, forced red once;
  all six gates green.
  Waits on: 1 (the fork writes to the log).

- [ ] **5. `agents/plan-writer.md`, and `write-plan` delegates to it.**
  One box: the new agent (tools `Read, Grep, Glob, Skill` and nothing else, cloned
  from `task-cutter`'s you-know-nothing posture; inputs by path — delta spec,
  capability spec, `proposal.md`, `decisions.md`; invokes `write-plan` for the
  pillars, whose orchestrator-only conduct is marked as not the writer's to follow;
  returns plan markdown plus what the inputs failed to answer; writes nothing) and
  the skill's delegation (orchestrator launches one fresh `plan-writer`, reads the
  gaps before writing, writes `plan.md`). An agent nothing launches and a skill
  naming an agent that does not exist each fail reachability alone — they merge
  together. `check-payload --index` regenerated so the agent appears in
  `docs/SPEC.md`.
  Closes: delta criterion "plan-writer drafts, orchestrator writes" — `check_wiring`
  row `skills/write-plan/SKILL.md · plan-writer` (agent frontmatter covered by the
  existing agents section), forced red once; all six gates green.
  Waits on: 4 (same file; the fork precedes the delegation in the skill's own order).

- [ ] **6. Routing, narrative, and the stop-owner drift.**
  `commands/libretto-flow.md` §4: "Phase 2 is the exception" becomes phases 2 and 5,
  each described. `docs/FLOW.md`: the same story, plus its stop-table row corrected —
  the 5→6 seam owns the second stop, not phase 5. `commands/libretto-attacca.md`:
  interview, fork and markers become `(assumed)` entries in `decisions.md`, each
  naming what changes if wrong; the fork rides attacca's pre-answered table,
  recommended option taken. (The payload spec's own stop-table row is corrected at
  landing, with the delta — per the plan's blast radius, never in this box.)
  Closes: delta criteria "attacca logs assumed entries" — `check_wiring` row
  `commands/libretto-attacca.md · decisions.md`, forced red once — and "three stops",
  already proven by the existing stop-table wiring; all six gates green.
  Waits on: 2, 3, 4, 5 (it narrates the shape they build).
