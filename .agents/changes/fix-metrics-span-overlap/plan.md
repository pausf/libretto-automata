# fix-metrics-span-overlap — plan

Executed by `build-and-check` (phase 6). Runs first of the four metrics changes: it is
the only wrong number, and its test must be red before the fix.

- [x] red test: three changes with overlapping [first, last] ranges — footer total is
      the union, disjoint ranges still sum
      (spec: Verification · Proof: TestTotalSpanMergesOverlappingChanges)
- [x] merge overlapping ranges in `metrics()` before summing; test goes green
