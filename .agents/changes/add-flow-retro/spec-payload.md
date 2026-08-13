# add-flow-retro — payload delta

Targets: payload

The flow learns from its own corrections. Today, when the user says "that is wrong,
the tag does not go like that", the correction fixes one commit and dies in the
scroll. This delta makes it an artifact: captured while the flow runs, spent by a
retro command that turns lessons into fixes.

## Outcomes

- **Corrections get captured where they happen.** `evidence` — the standing skill
  every phase already invokes — gains a rule: when the user corrects the work
  (names a mistake, reverses something the flow did, restates a convention the flow
  broke), append one entry to `.agents/lessons.md` in the project the flow is
  running in, and carry on. Capture never interrupts the work and never asks.
  **A correction is the user saying work already done was wrong.** A changed ask —
  "actually, make it blue instead" — is new work, not a lesson, and it is not
  captured: the ledger measures the flow's errors, not the user's evolving intent.
- **The ledger is append-only markdown with a fixed, countable header.** One entry:

  ```
  ## 2026-08-13 · add-flow-retro · build-and-check
  Said: the tag here is <type>/<scope>, never bare
  Did: wrote a bare tag in the commit subject
  ```

  Header fields: date · change name (`-` when no change is open) · phase skill
  active. `Said:` carries the user's words, not a paraphrase. An entry the retro has
  handled gains one more line — `Resolved: <date> · <what was done>` — and is never
  edited otherwise.
- **A retro command spends the ledger.** `/libretto-retro` invokes a new skill,
  `retro`, which reads every entry without a `Resolved:` line and classifies each:
  - **project knowledge** — the flow lacked a fact of this project. The retro
    records the convention where this project keeps its contract (the owning
    capability spec's prior decisions, or `AGENTS.md`), marks the entry resolved.
  - **flow defect** — the payload skill led the work wrong. The retro names the
    skill and **proposes the exact diff** to it, in the report. It never applies it:
    the payload is the product, and it does not edit itself without eyes. The entry
    is marked resolved only when the user says what they did with the proposal.
  - **one-off** — not worth preventing. Marked resolved with that reading, so it is
    never re-classified.
- **The retro reports what the metrics can then count**: entries found, how each was
  classified, what was written where, and every proposed payload diff.

## Scope boundaries

**In:** the capture rule in `evidence`, the ledger format, the `retro` skill, the
`libretto-retro` command, and the wiring checks that keep all four referenced.

**Out, named:**

- **cross-project mode.** The retro runs where the flow ran. Back the day one
  payload lesson shows up in two projects' ledgers — the user's call, 2026-08-13.
- **the retro applying a payload diff.** Propose only — the user's call, 2026-08-13.
  Back, if ever, as an explicit flag after the proposals have earned trust.
- **capturing the model's self-detected failures.** Gates and `evidence` already
  own those. The ledger records *user* corrections — the signal that is unambiguous
  and otherwise lost.
- **a hook or automation for capture.** Recognising "the user is correcting me" is
  judgment; it lives in a skill instruction, not a trigger.
- **the Go binary writing the ledger.** Delivery never writes payload state; it
  only reads it (see the cli delta).
- **editing or deleting ledger entries.** Append and mark. History is the point.
- **a schema beyond the header line.** Two labelled lines and a resolution marker.
  A richer format returns the day the retro measurably misclassifies for lack of a
  field.

## Constraints

- A skill may only invoke what gets installed: the `retro` skill is self-sufficient
  and names no path from this repository.
- `commands/libretto-retro.md` routes and never implements, like every command here.
- The ledger lives at `.agents/lessons.md` — beside `.agents/changes/` and
  `.agents/specs/`, versioned by the project's own git. Per-change files were
  rejected: phase 8 deletes the change folder at landing, and lessons must outlive
  the change that taught them.
- The header line format is load-bearing: `libretto metrics` parses it (cli delta).
  Changing it is a contract change on both capabilities at once.
- Frontmatter `name:` equals directory or filename, as everywhere.

## Prior decisions

- **Ledger per project at `.agents/lessons.md`** — the user's call, 2026-08-13.
  Central-in-`~/.claude` rejected: outside git, no review possible, mixes projects.
- **Per-project retro only** — the user's call, 2026-08-13. See scope.
- **Propose, never apply, on payload diffs** — the user's call, 2026-08-13. This is
  option A from the design conversation; option B (auto-retro inside phase 8) was
  rejected as too much power without eyes, revisitable once the ledger has proven
  its lessons are good.
- **Capture lives in `evidence`, not in eight phase skills.** One place, already
  invoked at every phase; eight copies of the rule is eight things that drift.
- **Two lesson types plus one-off, decided by where the fix belongs.** A retro that
  mixes them "fixes" the flow with one project's manias and breaks it for the rest.

## Task breakdown

- [ ] `evidence` — the capture rule and the entry format, stated once
- [ ] `skills/retro/` — read, classify, record or propose, mark resolved
- [ ] `commands/libretto-retro.md` — route to the skill, describe nothing
- [ ] `scripts/check-payload` — wiring: the command routes to a skill that exists,
      and the decisive capture words are in `evidence`

## Verification criteria

- the `retro` skill parses, is named correctly, and references only what installs
  Proof: scripts/check-payload
- `libretto-retro` routes to `retro` and the reference resolves
  Proof: scripts/check-payload
- the capture rule's decisive words are in `skills/evidence/` — wiring only, a
  prompt is checked by running it
  Proof: scripts/check-payload

**Claims, not facts, until a run observes them:** a correction mid-flow producing an
entry without interrupting the phase; a retro classifying a real ledger; a proposed
payload diff a user could apply verbatim. The first real flow after this lands is
the test.
