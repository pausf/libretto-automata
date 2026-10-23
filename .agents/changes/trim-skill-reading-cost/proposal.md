# trim-skill-reading-cost

Tracker: none
Targets: payload
Durable decisions: the four under Prior decisions below

## The ask, verbatim

> Pendiente aparte: la "dieta" de los skills (se leyeron ~3.000 líneas de skills para 15 líneas de código). ¿La encolo como idea?

> si

Then, picked up the same day:

> Lo que resta (los 3 puntos):
> - El costo de lectura es desproporcionado. 3.000 líneas de proceso para 15 de código es exactamente la ceremonia contra la que ODD está pensado. Hoy un cambio chico paga el precio de uno grande.
> - El flow se audita a sí mismo más de lo que entrega. Mirá el historial: muchas versiones son gates sobre gates y specs sobre el propio flow. Eso tiene valor, pero el payload que usa la gente en otros repos mejora más lento que el aparato que lo vigila.
> - La autonomía choca con los permisos. Asumimos "haz todo", pero el clasificador no lo reconoce para merges ni agentes headless. Ese tramo final todavía depende de vos, y conviene asumirlo en el diseño en vez de toparse con él.
>
> arregla esto, tu hazlo y no pregunte termina todo

## Reading

A trial run of `/libretto-attacca` on a throwaway project read about 3,000 lines of skill
text to build about 15 lines of code. Every invocation loads a skill's whole body, so text
that applies only on a condition still has to be read. That covers a remote, a tracker key,
a fan-out and the `--trace` details, and those paths cost most of the reading. A small
change should pay for the paths it actually takes.

## Decision log

### Session 2026-10-23
- Q: the approach? → A: progressive disclosure. Conditional sections move into sibling files that SKILL.md names, each with the condition that loads it. Considered and rejected: rewriting each skill shorter, which would put 52 guarded rules at risk in one diff. (assumed)
- Q: `writing-plans` — keep or drop? → A: drop it from the payload. Nothing but write-tasks used it, and write-tasks spent 26 lines overriding its format. (assumed)
- Q: points 2 and 3 — a gate or a rule? → A: a rule in AGENTS.md. A gate for "audit less" is the problem it describes. (assumed)
- Q: a line budget? → A: yes, one check-payload row per split skill, at a ceiling below today's size, so the row is red before the split and green after. (assumed)
- Q: (spec-reviewer) which siblings does "names each sibling with its condition" cover? → A: the new `.md` siblings only; `spec-drift` is a tool, not reading. Each pointer has one shape, `**When <condition>, read `<file>` beside this file**`, and the row pins the condition and the file on the same line. (assumed)
- Q: (spec-reviewer) "three things" and the EARS example — which occurrences? → A: record-work's landing list only (the attacca description's three-item list is correct and stays); the fenced after-example in write-spec, which joined two `shall`. (assumed)
- Q: (spec-reviewer) "foreground" — the whole file or the launch? → A: the reader launch only. Running a proof in the foreground stays everywhere. (assumed)
- Q: (spec-reviewer) the checklist's shape — what syntax? → A: each box carries a `Closes:` naming the criteria it closes, the token converge already maps on, and a row pins it in write-tasks. (assumed)
- Q: (spec-reviewer) how do status and next reach `queue.md` without a hardcoded path? → A: they already invoke find-work, and find-work's SKILL.md says when to read `queue.md` beside itself; the commands change nothing. (assumed)
- Q: (spec-reviewer) "labelled request" against record-work's unlabeled-is-honest rule? → A: AGENTS.md says "green, with its bump label when the bump was decided" — record-work's rule stands. (assumed)
- Q: (spec-reviewer) docs/DESIGN.md's writing-plans prose? → A: left; it is a dated design record, not a list of what ships. README's and FLOW.md's counts and lists are updated. (assumed)

## Outcomes

1. When a flow skill is invoked, its SKILL.md carries only what every run needs. The
   material below lives in a sibling file, which SKILL.md names together with the condition
   that sends the reader to it:
   - record-work: the push, the request, the bump and the forge go in `forge.md`; the
     `spec-drift` modes, the `--trace` table and the opt-in hook go in `drift.md`
   - find-work: source 2, the tracker, goes in `tracker.md`; the queue scan goes in `queue.md`
   - write-spec: the fan-out of extra-capability writers goes in `fan-out.md`
2. The payload no longer ships `writing-plans`. `write-tasks` states the checklist's shape
   itself, in a few lines.
3. Four contradictions the trial run measured are gone:
   - present-work's "then commit it"
   - record-work's "three things" over a list of four
   - write-spec's EARS example, which `--anchors` reads as a joined criterion
   - "foreground" for a reviewer agent that the host runs in the background
4. `AGENTS.md` carries three rules:
   - A change to this repository's own apparatus needs a defect seen in a real run.
   - A merge that moved the payload is followed by a trial run. A subagent drives
     `/libretto-attacca` in a throwaway project, and what it finds becomes fixes or queued
     ideas.
   - The agent ends at a green request. The merge and any unattended launch are the
     user's.
5. The capability spec stays true:
   - The criterion that has record-work do "three things" says four.
   - The vendored-reference criterion no longer cites `writing-plans`.
   - The criteria naming write-spec's fan-out section name `fan-out.md`.

## Scope boundaries

In: the four split skills, write-tasks, the four contradictions, writing-plans' removal,
`AGENTS.md`, and the check-payload rows that follow every moved sentence.

Out:
- rewriting the prose of skills not named above
- a gate for the AGENTS.md rules
- changes to the CLI

The vendored skills that remain are untouched.

## Constraints

- Every `check_wiring` and `check_absent` row keeps guarding its sentence wherever that
  sentence moves. No row is deleted to make text movable.
- The bare-dollar-digit scan stays on SKILL.md and commands. Only an invoked body is
  substituted, and a sibling file is read, not invoked. The scan for references to
  unreachable `scripts/` widens to every Libretto skill's `.md` file, so moved text stays
  in its reach.
- The push question (`AskUserQuestion`) and the three `spec-drift` modes the landing runs
  stay in record-work's SKILL.md. `drift.md` takes the `--trace` table and the `--block`
  hook. `forge.md` takes everything that needs a remote: the attacca exception, the return
  to base, the bump and the forge detection.
- A sibling file installs with its skill directory, as `test-driven-development`'s already
  do.
- The `spec-writer` row allows its name inside `write-spec/fan-out.md` and inside the short
  `## Step 2b` left in SKILL.md. The rows on the vocabulary and the rebilling move to
  `fan-out.md`, and "context cannot hold" stays in SKILL.md.

## Prior decisions

- **Conditional guidance lives in a sibling file, named with its condition.** Decided
  2026-10-23 (assumed) from a measured trial run: about 3,000 lines read for a change of
  about 15. A skill body loads whole on every invocation, so a section that applies on a
  condition is paid for on every run. Lost: shortening each skill in place, which would put
  52 guarded rules at risk in one unreviewable diff.
- **Each split SKILL.md has a line budget in check-payload.** Decided 2026-10-23
  (assumed): without one, the bodies grow back one rationale paragraph at a time. The
  budget is the measured size at landing. Raising it is a visible edit, never a drift.
- **`writing-plans` is not shipped.** Decided 2026-10-23 (assumed): only write-tasks used
  it, and write-tasks spent a 26-line table overriding its file path, its header and its
  execution question. A vendored skill whose every instruction is overridden costs reading
  and buys nothing.
- **A reviewer is waited for, not run in the foreground.** Decided 2026-10-23 (assumed):
  the host decides where an agent runs, and the trial run saw both reviewers run in the
  background despite the word. What the flow needs is that nothing proceeds before the
  return is read.

## Task breakdown

- split record-work, find-work and write-spec, and move each guarding row with its text
- rewrite write-tasks' shape paragraph and drop writing-plans from every list
- fix the four contradictions
- add the AGENTS.md rules

Every row below is wiring: it pins a sentence, never that a reader followed it. The
repeated trial run is the behavioural check, and it is not a gate.

- **The SKILL.md of record-work, find-work and write-spec shall carry, for each new `.md`
  sibling, one line of the form `When <condition>, read `<file>` beside this file`.**
  Ceiling: the row pins the line, not that the condition is right.
  Proof: scripts/check-payload
- **If record-work's SKILL.md exceeds 300 lines, find-work's 300 or write-spec's 580, as
  `wc -l` counts the whole file, then check-payload shall fail naming it.** Today's are 490,
  449 and 721, so each row is red before the split.
  Proof: scripts/check-payload
- **The payload shall not ship a `skills/writing-plans/` directory.**
  Proof: scripts/check-payload
- **write-tasks shall state that each box carries a `Closes:` naming the criteria it
  closes.**
  Proof: scripts/check-payload
- **If present-work says "then commit it", then check-payload shall fail.**
  Proof: scripts/check-payload
- **If record-work's landing list is introduced as "three things", then check-payload shall
  fail.**
  Proof: scripts/check-payload
- **If write-spec's EARS after-example joins two `shall`, then check-payload shall fail.**
  Proof: scripts/check-payload
- **If review-spec, review-work or converge launches its reader "foreground", then
  check-payload shall fail.**
  Proof: scripts/check-payload
- **review-spec, review-work and converge shall each say that nothing proceeds until the
  reader's return is read.**
  Proof: scripts/check-payload

## Approach

**Move first, then trim.** Each moved section goes into its sibling file verbatim, and the
check-payload rows that guarded it are re-pointed at the new file. Only then does SKILL.md
gain one line per sibling: the condition and the path, for example *"With a remote, read
`forge.md` beside this file before the push question."*

- **record-work:**
  - keeps: commit per task, the spec ships with the code, the landing, messages and the
    branch backstop
  - keeps, short: the push question
  - moves: everything about the attacca exception, the return to base, the bump and the
    forge goes to `forge.md`; the `--trace` table and the `--block` hook go to `drift.md`
- **find-work:** keeps the phase-1 contract, source 1, the branches and source 3. It moves
  source 2 to `tracker.md` and the queue section to `queue.md`. `libretto-status` and
  `libretto-next` name `queue.md`.
- **write-spec:** keeps a short Step 2b that says when to fan out and points at
  `fan-out.md`.
- **write-tasks:** loses the delegation table and gains the shape it needs: boxes only,
  one per end-to-end unit, each naming the criterion that closes it and what it waits on.
- **writing-plans:** removed from the payload, `THIRD-PARTY.md`, `README.md`,
  `docs/FLOW.md` and `docs/PAYLOAD.md`.

| Considered | Why it lost |
|---|---|
| Rewrite each skill shorter in place | 52 guarded rules in one diff; nothing distinguishes a cut rationale from a cut rule |
| A `--lite` mode the commands pass | Two versions of every rule, and the lite one drifts |
| Keep writing-plans and shrink the table | The table exists only to undo what writing-plans says |

Risks:

| Risk | Caught by |
|---|---|
| A moved sentence loses its guard | Every row is re-pointed, not deleted, and check-payload fails on a missing file or pattern |
| A reader never opens the sibling | The condition is written in SKILL.md where the reader already is, and the trial run is repeated after landing |
| `libretto prune` leaves a dead writing-plans link on users' machines | prune exists for exactly this; noted in the PR |

Validation: `scripts/check-payload` and the six gates; one forced red per new row.
Rollback: one revert.

## Checklist

- [x] record-work split into SKILL.md, forge.md and drift.md, with rows re-pointed and a budget row · closes criteria 1–2
- [x] find-work split into SKILL.md, tracker.md and queue.md, with status and next pointed at queue.md · closes criteria 1–2
- [x] write-spec's fan-out moved to fan-out.md, with the spec-writer row and EARS example fixed · closes criteria 1, 2 and 4
- [ ] write-tasks shape inlined and writing-plans removed from the payload and every list · closes criterion 3
- [ ] present-work, record-work and the reviewer launches de-contradicted · closes criteria 4–5
- [ ] AGENTS.md rules for apparatus, trial runs and where the agent stops · no criterion (repository instructions)
- [ ] payload spec: criteria for "three things", the vendored refs and the fan-out location amended · Closes: outcome 5
