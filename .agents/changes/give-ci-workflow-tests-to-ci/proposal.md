# give-ci-workflow-tests-to-ci

Tracker: none
Queued: 2026-10-22

## The ask, verbatim

> cli's Governs: claims release_workflow_test.go, label_workflow_test.go and gates_test.go, which test ci's workflows.

(Found while measuring report-drift-trend-in-metrics, 2026-10-22; captured by the orchestrator under the queue-run authorisation.)

## Reading

Exclude the three from cli and claim them in ci's Governs:, so a workflow test change points at the ci spec.
