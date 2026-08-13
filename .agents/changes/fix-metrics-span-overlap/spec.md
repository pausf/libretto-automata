# fix-metrics-span-overlap — delta

Targets: cli

This is a bug: the footer's total claims to be wall clock and is not. The criterion
below is written before the fix and its proof must be red first.

## Outcomes

`libretto metrics`' footer total counts each calendar hour once, however many changes
were open during it. Two changes open the same week contribute one week, not two.

## Scope boundaries

In: the footer total only. Out: per-row spans (already correct), any new column, any
change to how a single change's span is measured.

## Constraints

Derived from git only — no instrumentation. The merge is over [first, last] ranges the
rows already carry; no new git queries.

## Prior decisions

- **Assumed (attacca, 2026-08-13): fix the number by merging overlapping ranges, rather
  than relabelling the footer as "accumulated span".** Merging keeps the footer meaning
  what it says; relabelling keeps a number nobody can act on. If wrong, the union code
  goes and one word in the footer changes instead.

## Task breakdown

1. A red test: three changes with overlapping ranges, footer total = the union.
2. Merge the [first, last] ranges before summing in `metrics()`.

## Verification criteria

- overlapping changes are not double-counted; adjacent-but-disjoint ranges still sum.
  Proof: cmd/libretto/metrics_test.go TestTotalSpanMergesOverlappingChanges
