# Plan — split the plan from the task list

Spec: [spec.md](spec.md) · Proposal: [proposal.md](proposal.md)

## Summary

Phase 5 keeps its number and changes its output: `plan.md` becomes the technical
approach. The checklist it used to be moves to `tasks.md`, cut in a new unnumbered seam
by a subagent that never saw the design conversation.

Three edges: the payload writes the artifacts, the Go binary reads one of them by path,
and the documents describe both.

## Technical context

| | |
|---|---|
| Language | Go 1.26.5 — but only `cmd/libretto/` moves; no new package, no new dependency |
| Payload | markdown with YAML frontmatter, symlinked by the CLI |
| Gates | `gofmt -l .`, `go vet`, `go test ./... -count=1`, `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors` |
| Generated | `docs/PAYLOAD.md` — `scripts/check-payload --index`, never by hand |
| Blast radius | 3 Go files + 2 test files · 4 skills · 1 new agent · 3 commands · 4 docs · 2 capability specs |

The constraint that shapes everything below: **`libretto metrics` reads git history for
changes whose folders were deleted when they landed.** Every one of those has `plan.md`
and will never have `tasks.md`. A rename that ignores this silently blanks the churn
column for the entire history — which is exactly the bug `.agents/specs/cli/spec.md:865`
records having already been paid for once.

## The approach

### One helper owns the checklist path

`planPath()` becomes `checklistPath()` and returns `tasks.md` when it exists, `plan.md`
otherwise. Every caller in `loop.go` already goes through it, so the fallback is written
once and no call site learns about the transition.

`metrics.go` does not use that helper — it builds a path for `git log` against history,
where `os.Stat` is meaningless because the file is gone from disk. It gets its own
fallback: run the log against `tasks.md`, and if that returns nothing, run it against
`plan.md`.

### Alternatives rejected

| Considered | Why it lost |
|---|---|
| a new `design.md`, leave `plan.md` alone | zero Go churn, and it leaves the reported problem exactly where it was: a Lead opens `plan.md` and still finds a checklist. The complaint is about the file named `plan` |
| rename and drop the fallback | `metrics` is retroactive by design. Every landed change reads `—` and the tool's most useful column dies for the whole history to save eight lines |
| one `git log` pathspec covering both files | fewer calls, but it counts checkbox churn out of `plan.md` — which is now a prose document that may legitimately contain a box in an example. The metric would inflate quietly, which is the failure mode metrics can least afford |
| a numbered phase 6, shifting 6/7/8 up | touches README, AGENTS.md, four docs and every skill naming a phase. `docs/FLOW.md:251` already rejected this once for the review seam, for the same reason |
| the cutter runs inline in the orchestrator | it is what the user asked against, and the 6→7 seam already documents why: the session that argued for a design cuts boxes against the argument, not the document |

### The cutter's contract

`agents/task-cutter.md`, modelled on `agents/spec-writer.md` — read-only on the repo,
`Write` on one file, and told explicitly that what it was not given it does not know.

It reads `spec.md` and `plan.md` and returns two things: the checklist, and **what those
two documents failed to answer.** The second half is the load-bearing one — a cutter
that cannot produce a box because the plan never said how something gets built has found
a defect in the plan, and the current flow has no way to surface that before phase 6
runs into it.

## Risks

| Risk | Mitigation |
|---|---|
| a change in flight elsewhere has `plan.md` as its checklist, and `loop` stops finding it | the `os.Stat` fallback in `checklistPath()` reads it unchanged |
| the churn column blanks for landed history | second `git log` against the legacy path, proved by `TestMetricsFallsBackToLegacyPlan` |
| `plan.md` now means two things depending on when a change was created | the fallback is transitional and says so in a comment naming what removes it: no `.agents/changes/*/plan.md` containing a checkbox anywhere in the corpus |
| `docs/PAYLOAD.md` drifts and the gate fails on a page nobody edited | regenerate with `--index` as the last step, before the gates run |
| the new skill is unreachable — nothing routes to `write-tasks` | `check-payload` checks reachability; `libretto-flow` and `libretto-attacca` both get the route |

## Validation

The six gates, and `check-payload` is the one that carries this change: a new skill, a
new agent, three renamed references and a generated index all fail there if any of them
is half-done.

Two new tests, and both are written to fail first: `TestChecklistPathFallsBackToPlan`
proves the legacy read, `TestMetricsFallsBackToLegacyPlan` proves the history is not
lost. Green on a first run is not evidence here — the fallback is exactly the kind of
code that passes because the primary path never ran.

## Rollback

One revert. Nothing migrates, nothing is deleted, and no on-disk state changes shape:
a tree that reverts this commit reads `plan.md` as a checklist again, which is what
every existing change already has.

## Complexity deliberately kept

**Two fallbacks rather than one.** `loop` and `metrics` answer different questions —
what is on disk now, and what was in history — and a shared abstraction over both would
be an interface with two implementations that never converge. `ponytail:` marked at both
sites with the condition that removes them.
