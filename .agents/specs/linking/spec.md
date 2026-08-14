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

### Three plans, one set of actions

| State | `install` | `prune` | `uninstall` |
|---|---|---|---|
| `missing` | `create` | — | — |
| `linked` | — | — | `remove` |
| `wrong target` | `repoint` | — | `remove` |
| `stale` | — | `remove` | `remove` |
| `conflict` | `skip` | — | `skip` |

**`uninstall` is the pair of `install`, not of `prune`.** Prune cleans up after the
*repo* changed — a renamed item leaves a link pointing at nothing — and deliberately
spares links that are correct. That sparing is the guarantee that makes it safe to run:
you clean one broken link without risking a whole installation. Uninstall removes links
that are **working**, because the user changed their mind rather than because the repo
did.

Merging the two would turn a cleanup command into a delete command, and the day somebody
runs it expecting to tidy one broken link and loses their installation is the day this
tool stops being trustworthy.

**No new applying logic for any of them.** All three produce the same actions over the
same `Apply`, which already re-checks ownership at write time and already removes a link
without following it. A new plan is a new question about *which* entries, answered by a
pure function over machinery that is already proven.

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

### The four acts cover a written file too

`Create` makes the destination correct, `Repoint` makes an owned destination correct,
`Remove` deletes an owned destination with nothing behind it, `Skip` refuses. **What
varies is how, and the target decides — not the plan.** For a kind installed by
transform, `Create` and `Repoint` write bytes instead of making a link, and the bytes
travel on the entry: `link-state` already computed them to decide the state, so `Plan`
stays a pure function of entries with no target to ask.

**`Plan`, `PrunePlan` and `UninstallPlan` are untouched by that**, which is the payoff for
`link-state` refusing a sixth state — and it is why `uninstall` takes generated agents
back out with nothing added: `Linked`, `WrongTarget` and `Stale` already map to `Remove`
and `Conflict` to `Skip`. That is a guarantee from the state vocabulary, not a case
somebody remembered to add.

**A generated write is atomic**: a temporary file in the destination directory, renamed
over the target. Not a nicety — OpenCode throws on a malformed agent rather than skipping
it, so a torn write does not degrade one agent, it breaks the host's config load. The temp
file never goes in the system temp directory, because `os.Rename` fails across
filesystems. **Ceiling named:** atomic per file, not per plan — an `Apply` interrupted
between two files leaves one new and one old, the same guarantee linking already gives for
symlinks.

## Constraints

**Ownership is re-checked at the moment of writing.** A plan is a snapshot of a
moment; between the scan and the write another tool may have replaced a link, or a
user may have dropped a real directory in its place. Acting on a stale classification
is how a tool destroys something it never examined. When the re-check disagrees with
the plan, the action is refused and reported.

**One failure does not abandon the rest.** A plan of ten with one impossible item
performs the nine and reports the one. An early exit tells the user about neither.

**An emptied destination directory is left in place.** `~/.claude/skills/` is shared
with other tooling; removing it because our last item left would be deleting something
we did not create.

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

Complete. Phase 3.1 planning, 3.2 applying, 3.3 prune, and `uninstall` after it.

## Verification criteria
- **`Create` on a generated kind writes the file** — a regular file at 0644, carrying the
  marker, and recognised as ours immediately afterwards
  Proof: internal/link/apply_generated_test.go TestCreateWritesGeneratedContent
- **`Repoint` rewrites a drifted generated file whole.** Nothing in it is the user's;
  every byte was derived
  Proof: internal/link/apply_generated_test.go TestRepointRewritesGeneratedContent
- **`Repoint` refuses a file that stopped being ours between the scan and the apply**, and
  leaves it byte-identical. The plan is computed from a scan that is already stale by the
  time it runs, and this re-check is the last thing between a race and somebody's file
  Proof: internal/link/apply_generated_test.go TestRepointRefusesAForeignFileAtApplyTime
- **`Remove` deletes an owned generated file and spares a markerless one** at the same
  path
  Proof: internal/link/apply_generated_test.go TestRemoveDeletesAnOwnedGeneratedFile
- **a second install of a generated tree plans nothing** — idempotence is the promise that
  makes install safe to re-run, and a generated tree that rewrote itself every run would
  also mean `status` never reads clean
  Proof: internal/link/apply_generated_test.go TestGeneratedApplyIsIdempotent
- **no temporary file survives a write**, so a scan never finds a half-written neighbour
  and calls it a conflict
  Proof: internal/link/apply_generated_test.go TestGeneratedWriteLeavesNoTempFile
- **the temporary file is created in the destination directory**, never the system temp
  directory: `os.Rename` fails across filesystems and a target root is often on another
  one
  Proof: internal/link/apply_generated_test.go TestGeneratedWriteUsesTheDestinationDirectory

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
- every owned state becomes `remove` for uninstall, and a conflict becomes `skip`
  Proof: internal/link/plan_test.go TestUninstallPlan
- **a correct tree still gives uninstall a full plan** — nothing wrong is not nothing to do
  Proof: internal/link/plan_test.go TestUninstallPlanRemovesWorkingLinks
- uninstall never plans a write against a conflict
  Proof: internal/link/plan_test.go TestUninstallPlanNeverPlansAConflict
- **a real foreign directory and a foreign symlink both survive an uninstall**
  Proof: internal/link/apply_test.go TestUninstallLeavesForeignEntriesAlone
- **uninstall removes the link and not the item in the repo**
  Proof: internal/link/apply_test.go TestUninstallRemovesLinksNotSources
- an emptied destination directory survives
  Proof: internal/link/apply_test.go TestUninstallLeavesTheDirectoryInPlace
- a link that stopped being ours is refused here too
  Proof: internal/link/apply_test.go TestUninstallRefusesALinkThatStoppedBeingOurs
