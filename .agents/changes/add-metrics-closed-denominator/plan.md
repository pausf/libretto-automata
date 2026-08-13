# add-metrics-closed-denominator — plan

Executed by `build-and-check` (phase 6). Second: touches the same counters the overlap
fix's test observes.

- [x] track net box total across the plan.md diff walk in `measure`/`churn`
- [x] render the closed cell as `n/total` (net closed / net total); `—` unchanged for
      plan-less changes
      (spec: Verification · Proof: TestClosedShowsItsDenominator)
- [x] reworded box changes neither number; no-plan change never shows `0/0`
      (Proof: TestAChangeWithNoPlanReportsADashNotAZero — its no-"0" assertion already
      forbids `0/0`; the reword case is asserted in TestClosedShowsItsDenominator)
