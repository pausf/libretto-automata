# prove-sparkline-window-not-the-clock

Tracker: none
Queued: 2026-10-22

## The ask, verbatim

> TestWikiPageCarriesTheActivitySparkline's fixture dates fall inside both possible windows, so "never the clock's" cannot fail until about 2027-03.

(Found while measuring check-clause-coverage-in-anchors, 2026-10-22; captured by the orchestrator under the queue-run authorisation.)

## Reading

Move the fixture dates so the latest-commit window and the clock's window disagree today, and the criterion fails if the clock is used.
