# Delta — the plan is the how, the tasks are the checklist

Applies to: `.agents/specs/payload/spec.md` and `.agents/specs/cli/spec.md`

## Outcomes

The flow writes three artifacts per change where it wrote two, and each one's name says
what it holds.

- **`plan.md` is the technical approach**, written at phase 5: the context it is built
  against, the decision taken, the alternatives it beat and why they lost, what could go
  wrong, and how the change gets validated. It holds no checkboxes and no state.
- **`tasks.md` is the checklist**, cut in the seam between 5 and 6. Everything phase 5
  guaranteed about it before — one writer, derived from the spec, ordered by dependency,
  every box standing alone — is unchanged and now lives in `write-tasks`.
- **The cutter has no memory of the design conversation.** One fresh `task-cutter`
  subagent reads the spec and the plan, and returns the checklist plus what those two
  documents failed to say.
- **`libretto loop` and `libretto metrics` read `tasks.md`**, and fall back to `plan.md`
  when it is absent, so every change already landed keeps its measurable history.

## Scope boundaries

In: `skills/write-plan/`, a new `skills/write-tasks/`, a new `agents/task-cutter.md`,
the `plan.md` path in `cmd/libretto/`, `docs/FLOW.md`, and every payload reference to
the checklist by that name.

Out: renumbering the phases. Out: removing the spec's `Task breakdown` pillar. Out:
EARS syntax in verification criteria. Out: anything about how phase 6 builds.

## Constraints

- The phase count stays at eight. The task cut is a seam, per the precedent
  `docs/FLOW.md:251` set for the 6→7 review.
- `metrics` reads history for changes whose folders no longer exist. It cannot require
  `tasks.md` to have existed.
- No skill may reference `scripts/` or `docs/`.
- `docs/PAYLOAD.md` is generated, never hand-edited.

## Prior decisions

- **`plan.md` is reused for the new document rather than a new `design.md`.** The
  complaint was that the file named `plan` is not a plan; a `design.md` beside an
  unchanged `plan.md` leaves that exact file unchanged and the complaint standing.
  Answered by the user, 2026-08-17.
- **The cutter is a subagent, not a phase.** Asked for directly, and it matches the
  6→7 seam: independence comes from the fresh context.
- **The fallback is a second `git log`, not a two-path pathspec.** One pathspec covering
  both files would count checkbox churn out of `plan.md` too — which, after this change,
  is a document that may legitimately contain a box in prose.

## Task breakdown

The write side (`write-plan` rewritten, `write-tasks` and `task-cutter` new), the read
side (`loop`, `metrics`, `root`), the references (`find-work`, `libretto-status`,
`spec-drift`, the two commands), and the documents (`FLOW.md`, the two capability
specs, the generated index).

## Verification criteria

- `loop` drives a change whose checklist is `tasks.md`
  Proof: cmd/libretto/loop_test.go TestLoopStopsWhenEveryBoxIsClosed
- `loop` reads a legacy `plan.md` when no `tasks.md` is present
  Proof: cmd/libretto/loop_test.go TestChecklistPathFallsBackToPlan
- `metrics` reports churn for a change that only ever had `plan.md`
  Proof: cmd/libretto/metrics_test.go TestMetricsFallsBackToLegacyPlan
- every skill, agent and command parses, resolves and is reachable
  Proof: scripts/check-payload
