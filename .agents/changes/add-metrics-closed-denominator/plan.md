# add-metrics-closed-denominator — plan

Executed by `build-and-check` (phase 6). Second: touches the same counters the overlap
fix's test observes.

- [ ] track net box total across the plan.md diff walk in `measure`/`churn`
- [ ] render the closed cell as `n/total` (net closed / net total); `—` unchanged for
      plan-less changes
      (spec: Verification · Proof: TestClosedShowsItsDenominator)
- [ ] reworded box changes neither number; no-plan change never shows `0/0`
      (Proof: TestAChangeWithNoPlanReportsADashNotAZero, extended)
