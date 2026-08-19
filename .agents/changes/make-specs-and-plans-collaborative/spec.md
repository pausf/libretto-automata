# Delta: make specs and plans collaborative and delegated

Targets: payload

The contract for phases 2 and 5 changes in two directions at once: the user's decisions
enter the documents verbatim through an interview and a fork, and the documents themselves
get written by fresh-context subagents instead of the session that argued them. The
conversation stays in the orchestrator; the drafting leaves it.

## Outcomes

- Every change with a contract carries `decisions.md`: a verbatim, dated Q→A log. It is
  the input contract for the writer subagents and the source the delta's *Prior
  decisions* section is filled from — so landing needs no new mechanism.
- Phase 2 interviews before anything is drafted: questions one at a time, each with a
  reasoned recommendation first and "no more questions" always available. The answers go
  to the log in the user's words, not paraphrased.
- The spec file is written by `spec-writer` in the single-spec case as in the fan-out,
  from `brief.md` and `decisions.md` given by path. Where the log does not reach, the
  writer leaves `[NEEDS CLARIFICATION: question]` and never guesses; the orchestrator
  asks the markers, logs the answers, and patches the marker text.
- Phase 5 presents two or three approaches with tradeoffs and the user chooses. Chosen
  and rejected — with why the rejected lost — go to the log, which is what makes the
  plan's alternatives pillar reconstructible outside the conversation for the first time.
- `plan.md` is drafted by a new `plan-writer` agent modeled on `task-cutter`: fresh
  context, no Write, returns markdown plus what the documents failed to answer; the
  orchestrator writes the file.
- The stop count stays at three. The interview rides phase 2's existing questions; the
  fork is a question inside phase 5, not a stop after it.

## Scope boundaries

**In:** `skills/write-spec/SKILL.md`, `skills/write-plan/SKILL.md`,
`agents/spec-writer.md`, a new `agents/plan-writer.md`, `commands/libretto-flow.md`,
`commands/libretto-attacca.md`, `docs/FLOW.md`, and this delta onto the payload spec.
Also in: the stop-owner drift — `docs/FLOW.md` and the payload spec's stop table still
name phase 5 as the second stop's owner when it moved to the 5→6 seam; both rows get
corrected here, since this change rewrites those rows anyway.

**Out, named:**

- Section-by-section validation of the plan (superpowers-style). It multiplies round
  trips and the specs here are short by design. It comes back if plans grow past what
  one read absorbs.
- Any new approval gate or stop. Three stops is the contract.
- `task-cutter`, `work-reviewer`, `review-spec`, `write-tasks`: untouched. The 5→6 seam
  keeps its shape.
- The interview in the trivial lane. A change with no contract has nothing to interview
  about; step 0's collapse is unchanged.
- `[NEEDS CLARIFICATION]` markers in plans. `plan-writer` returns gaps in its reply,
  like `task-cutter`; markers are a spec mechanism.
- `record-work`: unchanged. The log's entries reach the delta's *Prior decisions* in
  phases 2 and 5, and landing already carries that section.
- Mechanical enforcement of the interview shape. Nothing can count questions or check a
  recommendation's reasoning; the wiring is checkable, the conduct is not, and the
  ceiling is named below.

## Constraints

- The stop is the rule, never the widget: every question is `AskUserQuestion` where the
  native prompt exists and conversation where it does not, as `record-work` states once
  for the whole flow.
- Sub-agents never ask the user. Markers and gaps funnel through the orchestrator; the
  writer's return value is its only channel.
- One author per file holds, with exactly one declared exception: the orchestrator may
  replace a `[NEEDS CLARIFICATION: …]` expression — the brackets and everything inside
  them, nothing outside them — with the logged answer, in a spec the writer authored.
  Anything the answer changes beyond that expression is a relaunch, not a patch.
- `decisions.md` has one writer, the orchestrator — same rule as `tasks.md`.
- `plan-writer` gets `Read, Grep, Glob, Skill`. No Write, no Bash. Same posture as
  `task-cutter`.
- The fan-out's model/effort rule extends to the single-spec launch: pick the tier
  before the writer starts.

## Prior decisions

- **2026-08-19, the user, this change** — four answers, verbatim options chosen:
  - Interview bound: *"Blando ~5, con juicio"* — a soft target of ~5 with judgment both
    ways, "no more questions" always visible. Keeps the spirit of 2026-08-14's lifted
    cap; what changes is one call becoming a sequence.
  - Single-spec input: *"Brief siempre"* — `brief.md` is written for one spec as for
    ten. One prompt contract for `spec-writer`, auditable and re-runnable.
  - Marker resolution: *"Orquestador parchea inline"* — a mechanical replacement, not a
    relaunch. This is the declared exception to one-author-per-file, scoped to marker
    text only.
  - Attacca: *"Sí, marcadas 'assumed'"* — assumed answers are logged in `decisions.md`
    marked as assumed, each naming what changes if it is wrong. One home for decisions
    in both modes.
- **Retired, and declared here: the 2026-08-12 decision that phase 5 asks nothing.**
  "Questions-at-phase-5 were offered and declined" (payload spec, prior decisions) is
  reversed by the user on 2026-08-19: phase 5 now asks exactly one question, the
  approach fork. The original decision's reason — no third stop — survives intact: the
  fork rides inside phase 5 and the flow still stops in three places.
- **2026-08-14's bias-to-ask stands.** No hard cap, zero is legitimate and said in one
  line, never a form-length interrogation of what the code answers.
- The 2026-08-12 companion decision — *"para que el plan se cree entre los 2 y no solo
  tú"* — is the direction this change extends, not a constraint against it.

## Task breakdown

1. The `decisions.md` contract in `write-spec`: when it is created, its format (dated
   sessions, `- Q: … → A: …` verbatim, `(assumed)` marker), its single writer, and how
   its entries reach the delta's *Prior decisions*.
2. The phase 2 interview in `write-spec`: one at a time, soft ~5, recommendation with
   reason first, "no more" option, answers logged before the next question.
3. `spec-writer` generalized: the single-spec launch, `decisions.md` as named input, the
   `[NEEDS CLARIFICATION: …]` rule replacing silent guessing.
4. The marker loop in `write-spec`: ask remaining markers, log, patch marker text only.
5. The fork in `write-plan`: two or three approaches with tradeoffs, recommended first,
   chosen and rejected logged with why.
6. `agents/plan-writer.md`: new agent, cloned from `task-cutter`'s posture — inputs by
   path (delta spec, capability spec, `decisions.md`, proposal), returns the plan
   markdown plus what the inputs failed to answer, writes nothing.
7. `write-plan` delegates drafting to `plan-writer`; the orchestrator writes `plan.md`.
8. Routing and narrative: `commands/libretto-flow.md` phase 2/5 sections and the
   question rules, `docs/FLOW.md` the same, `commands/libretto-attacca.md` the assumed
   entries. The stop-owner drift fixed in the same pass.
9. The `scripts/check-payload` wiring rows behind every `Proof:` above — one row per
   criterion, named in the plan, written before the criteria can claim them. And two
   landing instructions, for the commit that applies this delta: the three-stops
   promise already lives in the payload spec, so the criterion **merges** rather than
   adding a second copy; and the payload spec's stop table still carries the drifted
   row `| 5 · write-plan | yes | …` — the landing commit **rewrites it** to name the
   5→6 seam, matching the row `docs/FLOW.md` already fixed.

## Verification criteria

- **Where** a change reaches phase 2 and a spec is written, `write-spec` **shall**
  direct the session to create `decisions.md` in the change folder and record each
  answer verbatim, dated, before the spec file is written.
  Proof: scripts/check-payload
- **When** phase 2 has questions, the skill **shall** instruct one question per call,
  each carrying a recommended option with its reason and a "no more questions" option,
  bounded by judgment around five.
  Proof: scripts/check-payload
  Ceiling, named: the wiring is checkable — the strings exist in the skill — and the
  conduct is not. Nothing mechanical counts questions or reads a recommendation.
- The spec file **shall** be authored by `spec-writer` in the single-spec case as in
  the fan-out, launched with `brief.md` and `decisions.md` by path.
  Proof: scripts/check-payload
- **Where** the brief and the log do not settle a decision, `spec-writer` **shall**
  write `[NEEDS CLARIFICATION: question]` in place of an answer, and the orchestrator
  **shall** resolve each marker by asking, logging, and patching the marker text only.
  Proof: scripts/check-payload
- **When** phase 5 runs on a change with a contract, `write-plan` **shall** present two
  or three approaches with tradeoffs as a native question, recommended first, and
  **shall** record the chosen and rejected approaches with reasons in `decisions.md`.
  Proof: scripts/check-payload
- `plan.md` **shall** be drafted by `plan-writer` — an agent whose tools are `Read,
  Grep, Glob, Skill` and nothing else — returning markdown the orchestrator writes.
  Proof: scripts/check-payload
- **While** a run is `/libretto-attacca`, the interview, the fork, **and every
  `[NEEDS CLARIFICATION]` marker the writer returns** **shall** become entries in
  `decisions.md` marked as assumed — the recommended answer taken, each naming what
  changes if wrong.
  Proof: scripts/check-payload
- The flow **shall** keep exactly three stops: after the spec, after the tasks are cut,
  and at push.
  Proof: scripts/check-payload
  Ceiling, named: the same one the payload spec already carries for its stop table —
  a script cannot tell a prompt from a paragraph; the check is the string search across
  stop-owning skills.
