# Delta — a finished change that never landed is reported, not omitted

Targets: payload cli

## Outcomes

Phase 1 sees the state that produced this change, instead of walking past it.

- **A change whose checklist has boxes and none of them open is reported as finished
  and not landed.** Not as work in flight — there is nothing to continue — and not as
  silence, which is what it got.
- The two changes already in that state are landed properly: their durable decisions
  retired into the *Prior decisions* of the capability that governs each — five to
  `payload`, one to `cli` — and their folders deleted.
- `/libretto-status` says the same thing, because it delegates to the same scan.

## Scope boundaries

**In:** `skills/find-work/SKILL.md`, `commands/libretto-status.md`,
`skills/write-plan/SKILL.md`, the three `check_wiring` rows that hold them up, `.agents/specs/payload/spec.md`,
`.agents/specs/cli/spec.md`, and the deletion of `.agents/changes/add-design-phase/` and
`.agents/changes/retire-plan-decisions/`.

**Six decisions are retired, not four.** Both this spec and the plan said four before the
5→6 cutter counted them: three sit in each change's `spec.md` *Prior decisions*. Five are
`payload`'s — `plan.md` reused rather than a new `design.md`, the cutter being a subagent
rather than a phase, the section-not-the-file comparison, the declaration living in the
plan rather than in a flag, and the check riding inside `--anchors` rather than becoming a
seventh gate. **The sixth is `cli`'s**: the `metrics` fallback is a second `git log` and
not a two-path pathspec, and it is about `cmd/libretto/**`, which `payload`'s `Governs:`
does not claim. Filing it under `payload` would put a decision where the anchor that
exists to find it cannot. Settled by the user, 2026-08-17.

**Out — and each of these was considered:**

- **Checking whether the delta was actually applied.** The reporter would have to read a
  capability spec and decide whether a delta is present in it, which is a reading, and a
  wrong one either accuses a correct landing or clears a broken one. Zero open boxes is
  the signal; what it means is the user's to interpret.
- **Failing on it.** This is a report, not a gate. `--retired` already fails the commit
  that gets the landing wrong; this catches the commit that never happened, and there is
  no commit to fail.
- **Migrating either capability spec to EARS.** This delta adds prior decisions, not
  criteria. The 545-criterion rewrite was declined on its own merits and nothing here
  changes that.
- **A CLI subcommand for it.** `libretto status` reports symlink state and is a different
  question; the flow's own status is a slash command, and it already exists.

## Constraints

- `find-work` may not reference `scripts/` or `docs/` — it installs into projects where
  neither exists.
- The scan must survive a change with **no checklist at all** (a proposal captured by
  `/libretto-queue`) and one with a `plan.md`-era checklist, per the existing two-name
  read.
- `payload`'s *Prior decisions* must gain the retired decisions **in the same commit**
  that deletes the folders, or `spec-drift --retired` refuses it. That is the gate
  working, not an obstacle to route around.

## Prior decisions

- **Zero-open-boxes is the whole signal.** Settled by the user, 2026-08-17, choosing
  *limpiar + detectar* over cleaning only: the alternative was for the same omission to
  recur with a human eye as the only detector.
- **A report, never a gate.** `--retired` is blind to this state by construction — it
  fires on a deletion, and the failure here is that no deletion happened. Two mechanisms,
  two moments, and neither can cover the other's.
- **The retirement reads the change's `spec.md` *Prior decisions*, never the plan's
  `Durable decisions:` line.** The two disagree in the tree right now — one plan says
  "the two" over a section holding three, and the other carries no line at all, having
  been written before the line existed. They are not two enumerations of one thing: the
  section is the list, and the line is a **claim about whether the list is empty**, which
  is exactly what `--retired` reads it for. `write-plan` gains that sentence, because a
  contract that permitted this reading permitted the contradiction. Found by the 5→6
  cutter, which is what it is for.
- **The proof is `check-payload`, and it proves the mandate is present.** Whether a
  session actually reports what the skill tells it to is behaviour, and behaviour is
  checked by running it. Named here rather than discovered by somebody trusting the
  citation further than it goes.

## Task breakdown

The scan and its wording in the two payload files; the sentence in `write-plan` that
settles which document the retirement reads; the three wiring rows that hold all of them;
the retirement of six decisions — five into `payload`, one into `cli`; the deletion of
both folders.

## Verification criteria

- **When** a change folder holds `tasks.md` or `plan.md` with at least one box in it and
  none of those boxes open, phase 1 **shall** report that change as finished and not
  landed, naming it. **Both filenames**, on the same two-name reasoning the box scan
  already carries — a change from before the rename is exactly as unlanded. **Where** the
  folder holds neither file, it **shall not** be reported that way: a captured idea is
  not an unlanded change.
  Proof: scripts/check-payload
- The status command **shall** carry the same report, because it delegates to the same
  scan rather than describing its own.
  Proof: scripts/check-payload
- The plan skill **shall** state that a plan's `Durable decisions:` line is a claim about
  whether the list is empty and never the list itself, which is the change's `spec.md`
  *Prior decisions*. Added because both readings looked right and the tree already held
  the contradiction; it is a mandate like the other two and gets a row like the other two.
  Proof: scripts/check-payload

**Ceiling named, and it covers all three criteria above.** `check-payload` proves the mandate
is present in the file. It cannot prove a session obeyed it — that is behaviour, and
behaviour is checked by running it. A phase that reads the rule and reports nothing
anyway surfaces in the 6→7 review or not at all, and citing this gate as though it
covered the whole clause is the failure this paragraph exists to prevent.

**The retirement gate is not restated here.** `payload`'s own criteria already promise
that a commit deleting a `plan.md` without moving *Prior decisions* is refused, with the
same `Proof:`. Two copies of one promise drift, and the copy nobody edits is the one that
reads as authoritative. What belongs to *this* change is that it is the first real
landing to exercise that promise — which is a fact about the work, recorded in the plan,
and not a second contract.
