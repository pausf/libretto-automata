# add-queue-command — payload

Targets: payload

The payload half of `add-queue-command`. The decision log, the approach and the checklist
are in `proposal.md`, and the vocabulary is its: *queued*, *dispatched by a branch*, the
*scan*.

## Outcomes

- `payload`'s non-goal "the Go binary knowing about the queue. It is payload, not
  delivery." is gone, replaced by the line that still holds: the binary never *defines*
  what is queued.
- find-work's *The queue — captured, never started* states the tie-breaker on one
  `Queued:` date — change name, ascending — so `/libretto-next`'s recommended pick and
  `libretto queue`'s first line are the same change.
- The same section names `libretto queue` as an implementation of it, and says the
  section is the rule where the two disagree.
- `/libretto-status` and `/libretto-next` are unchanged: they delegate the scan to
  find-work and never call the binary.

## Scope boundaries

**In:** two sentences in `skills/find-work/SKILL.md`; four lines in
`scripts/check-payload`; at landing, the non-goal rewritten in
`.agents/specs/payload/spec.md` and the three criteria below added to it.

**Out:**

- **any change to `/libretto-status`, `/libretto-next` or `/libretto-queue`.**
- **find-work running `libretto queue` in place of its own scan.** The skill works on a
  machine with skills and no binary; the binary is a reader, never the source.
- **restating the branch filter anywhere else in the payload.** find-work owns it.

## Constraints

- A skill may name the binary only where it works without it — the precedent is the
  `libretto init` pointer at the two refusal sites. The pointer here costs a skills-only
  machine nothing: the section is still the whole rule.
- `scripts/` is never referenced from a skill.
- `check-payload` proves wiring only — that the sentence is in the file that owns it, not
  that a session obeyed it.

## Prior decisions

- **The non-goal is reversed, not deleted silently.** 2026-10-22, attacca run. Its
  replacement in *Scope boundaries*: "the Go binary *defining* the queue. find-work's
  queue section is the definition; `libretto queue` and the wiki implement it, and where
  they disagree the binary is the bug." The original was already broken by the wiki's
  flow board, which read `Queued:` lines from the day it landed; this change says so.
- **Ties on one `Queued:` date go by change name, ascending.** Assumed 2026-10-22,
  attacca run. If wrong (capture order), the sentence and the binary change together.
- **`/libretto-status` and `/libretto-next` never call `libretto queue`.** Assumed
  2026-10-22, attacca run. Two paths to "what is queued" is the thing find-work's
  single-owner rule exists to prevent. If wrong, status calls the binary where it is on
  PATH, with the skill's scan as the fallback.

## Task breakdown

- the tie-breaker sentence in find-work's queue section
- the sentence naming `libretto queue` as an implementation, the section winning
- `check_wiring` for each sentence; `check_absent` of `libretto queue` in the status and
  next commands
- at landing: the non-goal rewritten, the three criteria added

## Verification criteria

- **find-work's queue section shall order queued changes that carry the same `Queued:`
  date by change name, ascending.** Proved as wiring only.
  Proof: scripts/check-payload
- **find-work's queue section shall name `libretto queue` as an implementation of its
  definition, and shall state that where the two disagree the section is the rule.**
  Proved as wiring only.
  Proof: scripts/check-payload
- **`/libretto-status` and `/libretto-next` shall not invoke `libretto queue`, and shall
  keep delegating the queue scan to find-work.** Proved as absence of the string.
  Proof: scripts/check-payload
