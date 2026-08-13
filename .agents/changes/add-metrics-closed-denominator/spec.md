# add-metrics-closed-denominator — delta

Targets: cli

## Outcomes

The `closed` cell reads `n/total` — boxes currently closed over boxes the plan currently
has — so `5/5` (about to land) and `5/18` (bogged down) stop looking identical.

## Scope boundaries

In: the `closed` cell and the counters behind it. Out: the `reopen` column, which keeps
its cumulative meaning — churn is the signal it exists for; a plan-less change keeps its
`—`.

## Constraints

Derived from the same `plan.md` diff walk `measure` already does — no new git queries.
The existing criteria hold: rewording nets zero, reopenings stay counted apart.

## Prior decisions

- **Assumed (attacca, 2026-08-13): the numerator becomes the *net* currently-closed
  count (closed minus reopened), and the denominator is the net box count over the
  plan's history.** The cell answers "how far along is this plan now"; cumulative
  transitions belong to `reopen`. If wrong, the numerator returns to cumulative closes
  and only the denominator was the ask.

## Task breakdown

1. Track net box count (any state) across the diff walk; expose plan total.
2. Render `n/total` in the closed cell; `—` unchanged for plan-less changes.
3. Tests for both.

## Verification criteria

- the closed cell carries the denominator, and a reworded box changes neither number.
  Proof: cmd/libretto/metrics_test.go TestClosedShowsItsDenominator
- a change with no plan still reports a dash, never `0/0`.
  Proof: cmd/libretto/metrics_test.go TestAChangeWithNoPlanReportsADashNotAZero
