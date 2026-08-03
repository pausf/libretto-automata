# Linking

Governs: internal/link/plan.go internal/link/apply.go internal/link/plan_test.go internal/link/apply_test.go

The only code in the program that writes. Planning is separated from applying so that
every question about what the tool *would* do is answerable without a directory in
sight.

## Outcomes

**Planning** turns a scan into a list of intended actions and touches nothing:

| State | Action |
|---|---|
| `missing` | `create` |
| `wrong target` | `repoint` |
| `conflict` | `skip` — present in the plan so the refusal is stated, never silent |
| `linked` | nothing |
| `stale` | nothing — removal is `prune`'s job, never a side effect of install |

**Applying** carries out a plan and reports the outcome of every action, including the
ones it declined. Three outcomes exist and they are not the same thing:

- **done** — the filesystem changed
- **refused** — the world no longer matched the plan, so nothing was touched. A safety
  outcome, not an error.
- **failed** — the operation was attempted and the filesystem said no

An already-correct tree produces an empty plan and therefore performs no writes. That
is what makes the command safe to run repeatedly rather than something to be careful
with.

## Scope boundaries

**In:** planning, creating links, repointing our own, removing our own stale ones,
creating missing parent directories.

**Out:**

- **resolving a conflict.** Never, under any flag. There is no `--force` and adding
  one would remove the only guarantee that makes this tool safe to install.
- deciding ownership — that is `ownership`, and this capability only consumes it
- following a link to its destination. Removing a symlink removes the link.
- git, remotes, rebuilds — that is `repo-sync`
- prompting. Confirmation is `cli`'s concern.

## Constraints

**Ownership is re-checked at the moment of writing.** A plan is a snapshot of a
moment; between the scan and the write another tool may have replaced a link, or a
user may have dropped a real directory in its place. Acting on a stale classification
is how a tool destroys something it never examined. When the re-check disagrees with
the plan, the action is refused and reported.

**One failure does not abandon the rest.** A plan of ten with one impossible item
performs the nine and reports the one. An early exit tells the user about neither.

**A stale link's destination is never followed.** `os.Remove` on a symlink removes the
link; the ownership re-check keeps a real directory from ever reaching that call.

**Planning is pure.** No filesystem access at all — it must produce a plan for paths
that do not exist. This is the property that makes the dangerous half testable in
memory.

## Prior decisions

- `skip` is a first-class action rather than an omission, because a tool whose promise
  is "never clobbers your work" should report what it declined.
- Links are written with absolute destinations. `ownership` accepts either, but one
  form written consistently is one fewer thing to reason about.
- Removal is separated from installation as a distinct plan, so no run of `install`
  can ever delete anything.

## Task breakdown

Complete. Phase 3.1 planning, 3.2 applying, 3.3 prune.

## Verification criteria

- each state maps to its action, and `linked`/`stale` map to nothing
  Proof: internal/link/plan_test.go TestPlan
- **an already-correct tree produces an empty plan**
  Proof: internal/link/plan_test.go TestPlanIsEmptyForACorrectTree
- planning never touches the filesystem, even for paths that do not exist
  Proof: internal/link/plan_test.go TestPlanIsPure
- prune plans only stale entries and nothing adjacent
  Proof: internal/link/plan_test.go TestPrunePlanTakesOnlyStale
- skips are excluded from the writing set
  Proof: internal/link/plan_test.go TestWritesExcludesSkip
- missing links are created, for both directory and file kinds
  Proof: internal/link/apply_test.go TestApplyCreatesMissingLinks
- **running install twice changes nothing the second time**
  Proof: internal/link/apply_test.go TestApplyIsIdempotent
- our own misaimed link is repointed
  Proof: internal/link/apply_test.go TestApplyRepointsOurOwnWrongLink
- **a real directory and a foreign symlink both survive untouched**
  Proof: internal/link/apply_test.go TestApplyNeverTouchesAConflict
- **a destination that changed after the scan is refused, not overwritten**
  Proof: internal/link/apply_test.go TestApplyRefusesWhenTheWorldChangedAfterTheScan
- prune removes ours and leaves a foreign link beside it
  Proof: internal/link/apply_test.go TestPruneRemovesOnlyOurOwnStaleLinks
- **prune removes the link and not what it pointed at**
  Proof: internal/link/apply_test.go TestPruneRemovesTheLinkAndNotItsDestination
- one impossible action does not abandon the others
  Proof: internal/link/apply_test.go TestApplyContinuesAfterAFailure
