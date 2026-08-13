# add-metrics-prefix-match — plan

Executed by `build-and-check` (phase 6). Independent of the other three — may run any
time after the overlap fix lands its test.

- [x] `filterName`: exact match wins, else unique prefix, else refusal naming the
      candidates; unknown name keeps its refusal
      (spec: Verification · Proof: TestAPrefixSelectsAChangeUnlessAmbiguous)
