# tier-verification-by-risk

Tracker: none
Targets: payload
Durable decisions: the one in Prior decisions below

## The ask, verbatim

> La verificación no escala con el riesgo: spec-reviewer y work-reviewer corren siempre. Medir primero en qué tipo de cambio encontró algo cada revisor, después decidir niveles (passive / medium / high).

(Proposed 2026-10-23 after draining the queue, comparing the flow with gentle-ai's organic workflow; accepted with "has todo".)

## Reading

Measure, from the 2026-10-22/23 run's ledger and reviews, which reviewer findings came from which kind of change, then decide whether the fresh readers run by named risk instead of always. The measurement may decide against tiers: they found real defects in every change that week.

## Decision log

### Session 2026-10-23
- Q: measure first — which reviewer found what, on which kind of change? → A: measured from git, below; the measurement decides against tiers (decided by the orchestrator under the user's standing authorisation; if wrong, a tier rule needs a predictor this data does not show)

## Outcomes

1. The measurement is recorded where the next proposal of risk-tiered verification will meet it.

## Scope boundaries

**In:** one entry in payload's *Prior decisions*. **Out:** any change to when `spec-reviewer` or
`work-reviewer` run; a size or risk classifier.

## Constraints

- No code, no skill text.

## Prior decisions

- **The fresh readers run on every contract-lane change; verification is not tiered by risk.**
  Measured 2026-10-23 over the ten changes merged 2026-10-22/23 (`measurement.tsv` in this
  change's history): per PR, the cold read's decision-log questions, the 6→7 entries in
  `.agents/lessons.md`, and the PR's lines outside `.agents/`. `spec-reviewer` found 3–17
  consequential gaps in **9 of 9** contract-lane changes, and size did not predict it — the
  81-line change had 5, the 2,195-line one 6. `work-reviewer` found something in **4 of 10**,
  and size did not predict that either: nothing in the two largest (2,195 and 691 lines), a
  real regex gap in the 81-line gate change, and `init` writing outside `.agents/` in a
  472-line one. A tier rule by size would have skipped exactly the reviews that found bugs.
  What would reopen it: reviewer cost measured per change, and a predictor that names in
  advance which reviews come back empty — neither exists here, and ten changes is a small
  sample. Prompted by comparing the flow with gentle-ai's risk-tiered verification.

## Task breakdown

- land the measurement.

## Verification criteria

None: this change ships no behaviour, only a recorded decision.

## Approach

Counted from git: each PR's merge diff (`merge^1..merge^2`, `.agents/` excluded), the last
pre-landing copy of its `proposal.md` (questions in any review session of its decision log), and
`.agents/lessons.md` 6→7 headers naming the change. `add-init-command` predates `spec-reviewer`
and is excluded from the cold-read column.

## Checklist

- [x] **Land the measurement into payload's Prior decisions.**
