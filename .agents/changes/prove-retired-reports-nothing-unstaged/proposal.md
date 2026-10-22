# prove-retired-reports-nothing-unstaged

Tracker: none
Queued: 2026-10-22

## The ask, verbatim

> No spec-drift self-test case stages nothing for --retired, so its "report nothing" clause has no test of its own.

(Found while measuring check-clause-coverage-in-anchors, 2026-10-22; captured by the orchestrator under the queue-run authorisation.)

## Reading

One --self-test case with an empty index asserting --retired prints its nothing-retired line and exits 0.
