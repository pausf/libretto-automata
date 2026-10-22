# prove-init-refuses-without-git

Tracker: none
Queued: 2026-10-22

## The ask, verbatim

> TestInitRefusesOutsideAGitRepository never exercises the "git is unavailable" half of its trigger.

(Found while measuring check-clause-coverage-in-anchors, 2026-10-22; captured by the orchestrator under the queue-run authorisation.)

## Reading

A test that runs init with git absent from PATH and asserts the refusal names git, so the second half of the criterion can fail.
