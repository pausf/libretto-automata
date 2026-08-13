# add-metrics-legend — delta

Targets: cli

## Outcomes

`libretto metrics` explains its own columns: a legend names what commits, span, closed,
reopen, state and the `—` cell each mean, so the table needs no translator.

## Scope boundaries

In: a static legend in the report, printed with the existing footer. Out: any flag to
toggle it, colour, localisation.

## Constraints

Same visual style as the existing `flowCeiling` footer (two-space indent, `·` separators
where useful). English, like all CLI output.

## Prior decisions

- Printed always, next to `flowCeiling` — the report already always prints what it
  cannot measure; what it does measure deserves no less. No flag: a legend behind a flag
  is a legend nobody reads.

## Task breakdown

1. A `flowLegend` const and its print, beside `flowCeiling`.
2. Test that the report carries the legend.

## Verification criteria

- the report names every column it prints, and what a `—` means.
  Proof: cmd/libretto/metrics_test.go TestTheReportExplainsItsColumns
