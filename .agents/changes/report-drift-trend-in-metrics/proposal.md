# report-drift-trend-in-metrics

Tracker: none
Targets: cli
Durable decisions: the one entry in Prior decisions below — spec drift is not measured in metrics

## The ask, verbatim

> Drift trend en `metrics`: cuántas veces el spec de cada capability se movió DESPUÉS
> del código. Convierte `spec-drift` (warn-only) en una señal de salud medible. Ninguna
> tool de la industria tiene esto todavía.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Derive per-capability drift events from git history — spec edits landing after the code
they govern — and report the trend in `libretto metrics`.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed or measured)

- Q: what is a drift event, precisely? → A: **the ask's own reading — a commit touching a
  capability's `spec.md` after earlier commits touched its governed code without it — was
  prototyped on this repository's whole history, and it is mostly noise.** At commit
  granularity (498 non-merge commits) it yields 47 events, and **37 of the 47 catch-up
  commits are landing commits**: the flow working as designed, since a change's code
  commits precede the commit that applies its delta. At landing granularity (70
  first-parent commits of `main`) it yields 13 events; all 13 were read, and **one** —
  `readme`, PR #72 after PR #63 — moved the spec for a reason related to the code before
  it. Measured, not assumed.
- Q: does the "moved after" count ship? → A: no (assumed). In 12 of 13 landing-level events
  the later spec edit was about something else — a badge, a diagram, a new command.
- Q: what ships instead? → A: *(superseded by the C decision below)* `spec-drift`'s own
  warning replayed over every landing — per capability, the landings touching governed
  code, those `spec-drift` would have warned on, and those since the spec last moved
  (assumed). Measured: 27 warned of 120 capability-landings; of ten read against their
  diffs, **none** changed a promise without its spec — six touched a path whose semantic
  owner is another capability whose spec moved in the same landing, three were edits inside
  an existing contract, one is arguable.
- Q: the fork — how is it built? → A: *(superseded by the C decision below)* A, a Go replay
  inside `metrics`: one `git log --first-parent --reverse --name-only` over `HEAD`; the
  `Governs:` reading and bash `case` globbing (`*` and `?` crossing `/`, never `path.Match`)
  mirrored in Go and pinned by a test running `spec-drift --trace` on a fixture; one
  `specDrift` field on `metricsReport` rendered as a text block after the loop block and a
  `spec_drift` JSON key (`capabilities`: `name`, `landings`, `warned`, `outstanding`;
  `name` alone without `Governs:`; `error` when history is unreadable; absent without
  specs); a parity test mapping every metrics block header to its JSON key; 25 criteria
  (assumed, the recommended option at the time). B — a `--history` mode in `spec-drift`,
  exec'd by `metrics` — lost to A because `metrics` would need bash and `rg` at runtime. C —
  ship nothing — lost to A on "the replay is exact about what it counts".
- Q: the sub-decisions under A? → A: *(superseded by the C decision below; none is built)* a
  landing is a first-parent commit of `HEAD`; a co-claimed path counts for each claimant
  (the shared alternative halved `cli` from 10 to 5); a change folder whose `Targets:` names
  the capability counts as the spec moving (no number changed here); landings before the
  spec's first appearance do not count (removes six `readme` warnings); `HEAD`'s `Governs:`
  applies to all history; specs found with `findSpecsDir`; the block corpus-wide whatever
  `<change>` names; the agreement test failing, not skipping, without `rg`; no `AGENTS.md`
  or README edit.

### Session 2026-10-22, after the first draft (orchestrator, under the user's standing authorisation)

- Q: A, B or C? → A: **C, ship nothing.** "Your own numbers decide it. The asked-for count
  is 1 related event in 13. The replay flagged 27 landings, and 0 of the 10 hand-checked
  ones changed a promise without its spec; 6 were Governs overlaps. A metric that is mostly
  noise teaches people to ignore it, and A costs 25 tests plus a Go mirror of bash globbing
  to report it." What changes if wrong: A is rebuilt from this document's history — the
  two superseded entries above carry its mechanism and its JSON shape, and the definition
  and numbers in *Prior decisions* are enough to rewrite the prototype. The A draft was
  never committed; a copy and the prototype sit in this session's scratchpad
  (`/private/tmp/claude-501/-Users-pau-sanchez-gitrepos-libretto-automata/4ad75385-9673-4a78-8ec2-0826fbdbf663/scratchpad/drift/`,
  `proposal-option-A.md` and `drift.py`), which is not durable.
- Q: the two `Governs:` overlaps the measurement found — fixed here? → A: no. They are not
  this change's; listed under *Scope boundaries* for a separate capture.

## Outcomes

- The measurement is recorded where the next proposal of this idea will meet it: one entry
  in `cli`'s *Prior decisions*, carrying the definition prototyped, the numbers, the
  hand-checked sample, and what would reopen it.
- `libretto metrics`, its `--json` document and `spec-drift` are unchanged.

## Scope boundaries

**In:**

- `.agents/specs/cli/spec.md`, at landing: one *Prior decisions* entry, "spec drift is not
  measured in metrics".

**Out:**

- **the `metrics` spec-drift block** — no text block, no legend line, no ceiling bullet.
- **the `spec_drift` JSON key** — `metrics --json` keeps its current shape.
- **the Go mirror of `spec-drift`'s `Governs:` reading and bash `case` globbing**, and the
  test that would have pinned it against `spec-drift --trace`.
- **the parity test over every metrics block** (`TestMetricsJSONCarriesEveryTextBlock`).
  The recorded "whether it was mapped is review's to see, not a gate's" stands unchanged.
- **the "moved after" event count**, in any form.
- **any change to `spec-drift`.**
- **two `Governs:` overlaps the measurement found, left for a separate capture:**
  `cmd/libretto/readme_test.go` is claimed by both `cli` (`cmd/libretto/**`) and `readme`;
  and `cli` claims three tests of `ci`'s workflows — `release_workflow_test.go`,
  `label_workflow_test.go`, `gates_test.go`. Together they account for all ten of `cli`'s
  warned landings in the replay.

## Constraints

- No code. The change touches `.agents/specs/cli/spec.md` and this folder, nothing else.
- The *Prior decisions* entry carries the numbers in full, because the prototype lives in a
  scratchpad that will not survive the session.

## Prior decisions

- **Spec drift is not measured in `metrics`.** Assumed 2026-10-22, attacca run; decided as
  option C by the orchestrator under the user's standing authorisation, on a measurement of
  this repository at `dd6683d` (565 commits, 498 non-merge, 70 first-parent; 14 capability
  specs, all with `Governs:`).
  **The definition prototyped:** a *drift event* is a commit touching a capability's
  `spec.md` after one or more earlier commits touched paths its `Governs:` globs claim —
  read as `spec-drift`'s `governs_of` reads them and matched as its `matches_any` matches,
  where `*` crosses `/` — without touching that spec; paths under the specs directory and
  `.agents/changes/` are never code, and a commit touching a change folder whose `Targets:`
  names the capability counts as the spec moving. **The numbers:** at commit granularity,
  47 events, **37 of whose catch-up commits were landing commits** — the flow applying a
  delta after its code, by design. At landing granularity (a first-parent commit, a merge
  read against its first parent, counting from the first landing that touches the spec),
  **13 events, all read, 1 related** to the code before it (`readme`, PR #72 after #63).
  The fallback — `spec-drift`'s default-mode warning replayed per landing — flagged **27 of
  120** capability-landings (`cli` 10, `readme` 12, `review-project` 2, `payload` 2,
  `link-state` 1); of **10 hand-checked, 0** changed a promise without its spec: **6 were
  `Governs:` overlaps** (a path whose semantic owner is another capability whose spec moved
  in the same landing), 3 were edits inside an existing contract, 1 arguable.
  **Why nothing ships:** a metric that is mostly noise teaches people to ignore it, and the
  replay cost 25 tests and a Go mirror of bash globbing to report it. **B** — a `--history`
  mode in `spec-drift` exec'd by `metrics` — and **A** — a Go replay inside `metrics` with
  a `spec_drift` JSON key — both lost to **C**, ship nothing. **What would reopen it:** a
  corpus where the replay's hand-checked precision is high — most warned landings turning
  out to be promises moved without their spec — or a way to tie a spec edit to the code it
  answers (a commit trailer, a delta naming its landing), which the flow does not write
  today.

## Task breakdown

1. Land the one *Prior decisions* entry onto `cli`'s spec and delete this folder, in one
   commit.

## Verification criteria

This change ships no behaviour, so it carries no verification criteria: there is nothing a
test could fail.

## Approach

### Summary

Measurement only. The ask was prototyped against real history, its signal was 1 related
event in 13, and the warning-replay fallback was 0 real of 10 hand-checked; the definition
and the numbers land in `cli`'s *Prior decisions* so the idea is not re-proposed without
new evidence. The measurement itself is under *Decision log* and *Prior decisions* above.

### Technical context

One file at landing, `.agents/specs/cli/spec.md`, plus the deletion of this folder. The
landing commit must move `cli`'s *Prior decisions* section, which is what `spec-drift
--retired` checks when a change document carrying `Targets:` is deleted.

### The approach

Copy the *Prior decisions* entry above, verbatim, into `cli`'s *Prior decisions*; delete
`.agents/changes/report-drift-trend-in-metrics/`; run the six gates; commit.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| **A — a Go replay inside `metrics`, with a `spec_drift` JSON key** | 0 of 10 hand-checked warned landings were real; 25 tests and a Go mirror of bash globbing to report noise |
| **B — a `--history` mode in `spec-drift`, exec'd by `metrics`** | the same noise, plus bash and `rg` as runtime needs of `metrics` |
| **The ask as worded — a "moved after" count** | 37 of 47 commit-level events are landings; 1 of 13 landing-level events is related |

### Risks

| Risk | What catches it |
|---|---|
| The idea is re-proposed without the numbers being seen | the entry lives in `cli`'s *Prior decisions*, which `write-spec` reads for any change targeting `cli` |
| The landing commit forgets the entry | `spec-drift --retired` fails a deleted change document whose commit moves no *Prior decisions* section |

### Validation and rollback

The six gates at landing; `--retired` is the one that bears on it. Rollback: one revert.

### Complexity deliberately kept

None.

## Checklist

- [ ] **Land the measurement into cli's Prior decisions.** The entry above copied verbatim
      into `.agents/specs/cli/spec.md`'s *Prior decisions*, this folder deleted, same
      commit, six gates green. Closes: the outcome — the measurement recorded where the
      next proposal of this idea will meet it, checked by `spec-drift --retired` inside
      `--anchors`. Waits on: nothing.
