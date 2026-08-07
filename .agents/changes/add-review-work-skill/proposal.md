# add-review-work-skill

Tracker: none

## What was asked

> vamos hacer la skill de revisión independiente

Said in conversation on 2026-08-07, accepting this framing from the same session's
audit of the payload:

> Añadir al payload la skill de revisión independiente (review-work): un subagente
> fresco que lee el spec y el diff sin contexto de la sesión, verifica que cada
> `Proof:` corrió de verdad, y devuelve hallazgos antes de present-work.

## Why

The flow's own philosophy is that nothing is true until observed — and yet
`present-work` is a self-report by the same agent that wrote the code. No phase looks
at the work with fresh eyes. `docs/STATE.md` already names this as the last pending
item: *"an independent verifier, never run by whoever wrote the code — worth having
after the flow has been run once, not before."* This change is that run: the flow
builds its own reviewer.

## Shape, as asked

- A new payload skill, `review-work`, sitting between phase 6 (build) and phase 7
  (present): phase 7 presents *including* the reviewer's verdict.
- The reviewer is a **fresh subagent with no session context** — it reads the spec and
  the diff, not the conversation that produced them.
- It verifies that each `Proof:` citation names a test that actually ran, and returns
  findings.

Everything above the dash is the ask; how the phase is numbered, what the reviewer
reads, and what a finding blocks are for the spec to settle.
