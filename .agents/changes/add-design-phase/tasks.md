# Tasks — split the plan from the task list

Spec: [spec.md](spec.md) · Plan: [plan.md](plan.md)

Execution: `build-and-check`. One writer — the orchestrator marks the boxes.

- [x] **The checklist path, both sides** — `checklistPath()` in `loop.go` preferring
      `tasks.md` and falling back to `plan.md`; the second `git log` in `metrics.go`;
      both `ponytail:`-marked with what removes them. Tests written failing first.
      Spec: *`libretto loop` and `libretto metrics` read `tasks.md`*
      Closes on: `TestChecklistPathFallsBackToPlan` + `TestMetricsFallsBackToLegacyPlan` green,
      and `TestLoopStopsWhenEveryBoxIsClosed` still green against `tasks.md`.
      Waits on: nothing.

- [ ] **Phase 5 becomes the how** — `skills/write-plan/SKILL.md` rewritten: technical
      context, the decision and what it beat, risks, validation, rollback, complexity
      kept. Drops every line about checkboxes, state and one-writer.
      Spec: *`plan.md` is the technical approach*
      Closes on: `check-payload` green, and no reference in the file to a box or to state.
      Waits on: nothing.

- [ ] **The tasks seam** — `skills/write-tasks/SKILL.md` carrying everything the old
      `write-plan` guaranteed about the checklist, plus `agents/task-cutter.md` and the
      launch protocol. `libretto-flow` and `libretto-attacca` route to it.
      Spec: *The cutter has no memory of the design conversation*
      Closes on: `check-payload` green — the new skill and agent parse and are reachable.
      Waits on: the phase-5 rewrite, which is what the cutter reads.

- [ ] **Every remaining reference to the checklist by its old name** — `find-work`,
      `commands/libretto-status.md`, the `spec-drift` header comment, `root.go`'s
      comment.
      Spec: *`tasks.md` is the checklist*
      Closes on: `rg 'plan\.md' skills/ commands/ cmd/` returns only the fallback sites
      and the deliberate history references.
      Waits on: nothing. Independent of the three above.

- [ ] **The documents** — `docs/FLOW.md` §5 rewritten and the 5→6 seam added; the two
      capability spec deltas applied; `docs/PAYLOAD.md` regenerated with `--index`.
      Spec: all four outcomes
      Closes on: the six gates green, `spec-drift --anchors` resolving every new `Proof:`.
      Waits on: everything. It describes what the others built.
