# Link state

Governs: internal/link/scan.go internal/link/state.go internal/link/scan_test.go internal/link/state_test.go

Enumerating what the repository holds and reporting, per item per target, exactly what
the situation is. Read-only, always.

## Outcomes

Given a repo root and a target, a scan produces one entry per item per accepted kind,
each carrying exactly one of five states, plus an entry for every owned link in the
target that no current item explains.

| State | Meaning | Remedy |
|---|---|---|
| `linked` | owned symlink, right destination | none |
| `missing` | in the repo, absent from the target | `install` |
| `wrong target` | owned symlink, wrong destination | `install` repoints it |
| `conflict` | something foreign in the way | **none** — reported, never touched |
| `stale` | owned symlink with no item behind it | `prune` |

Five states, mutually exclusive and exhaustive. An entry with two states, or none, is
a bug in the classifier and not a situation to handle downstream.

Results are ordered by kind then by name, so two runs over an unchanged tree produce
byte-identical output — which is what makes the output diffable and the tests
golden-able.

## Scope boundaries

**In:** enumerating items, classifying them, finding stale links, counting.

**Out:**

- **any write.** This is the guarantee the whole capability rests on.
- deciding what to do about a state — that is `linking`
- resolving a conflict. Reported, never resolved, in this version and by design.
- git state, remote state, permissions. `cli` composes those in.

## Constraints

**An item is a shape, not a name.** A skill is a directory; an agent or a command is a
`.md` file. Anything of the wrong shape is somebody's note, not an item, and is
skipped rather than reported as broken.

**Dotfiles are never items.** This is what keeps `.gitkeep` out of the count, and it
means an empty payload directory reports zero rather than one.

**A missing kind directory is not an error.** A repo may simply hold no agents yet.
Likewise a target directory that does not exist yet holds nothing stale, because
nothing can be stale in a directory that is absent.

**`stale` deliberately subsumes "broken link".** An owned link whose source was
deleted has no item behind it; so does an owned link aimed at something in the repo
that was never an item. Same diagnosis, same remedy, one concept — splitting them
would be two names for one situation and two code paths to keep in agreement.

## Prior decisions

- **Conflict covers three different intrusions** — a real file, a real directory, and
  a symlink leaving the repo — because the remedy is identical for all three: leave
  it alone and say so.
- Counts distinguish "none" from "not applicable": a kind the target rejects is
  absent from the map rather than present as zero.
- `NeedsAttention()` is a property of the state, not of the caller, so exit codes
  cannot drift between commands that ask the same question.

## Task breakdown

Complete. Shipped in phase 2.

## Verification criteria
- **a generated file whose bytes match the transform is `linked`**, and its `Actual`
  carries the marker's source — the same meaning the field has for a symlink
  Proof: internal/link/generated_test.go TestGeneratedMatchingContentIsLinked
- **a generated file whose bytes differ is `wrong target`**, not a sixth state: ours, at
  the right path, with the wrong content, fixable by rewriting — which is what `wrong
  target` already means
  Proof: internal/link/generated_test.go TestGeneratedDriftIsWrongTarget
- **a markerless file at a generated item's path is a `conflict`**, reported and never
  overwritten
  Proof: internal/link/generated_test.go TestMarkerlessFileIsAConflict
- **a generated file whose source item is gone is `stale`** — prune's, not install's
  Proof: internal/link/generated_test.go TestGeneratedOrphanIsStale
- **a source whose frontmatter cannot be transformed is a `conflict`**, never `linked`
  and never a crash: we cannot say what belongs there, so nothing is touched
  Proof: internal/link/generated_test.go TestUntransformableSourceIsAConflict
- **a target that does not transform classifies exactly as before**
  Proof: internal/link/generated_test.go TestNonTransformingTargetIsUnaffected
- **a transforming target still symlinks the kinds it does not transform** — the bug the
  first `Transformer` interface had, where every skill in the opencode destination read
  as a conflict
  Proof: internal/link/generated_test.go TestTransformingTargetStillLinksItsOtherKinds

- every state is produced for the situation that defines it
  Proof: internal/link/state_test.go TestScanStates
- **the scan writes nothing** — the guarantee, asserted rather than assumed
  Proof: internal/link/state_test.go TestScanIsReadOnly
- foreign entries beside our own are ignored, not miscounted
  Proof: internal/link/state_test.go TestScanIgnoresForeignEntries
- ordering is stable between runs
  Proof: internal/link/state_test.go TestScanOrderIsStable
- every kind the target accepts is covered
  Proof: internal/link/state_test.go TestScanCoversEveryAcceptedKind
- a repo reached through an alias classifies the same way
  Proof: internal/link/state_test.go TestScanThroughAnAliasedRepoPath
- item shape is enforced, and dotfiles are excluded
  Proof: internal/link/scan_test.go TestItems
- a missing kind directory is not an error
  Proof: internal/link/scan_test.go TestItemsMissingKindDirIsNotAnError
- items carry absolute repo paths, so callers never re-derive them
  Proof: internal/link/scan_test.go TestItemsCarryAbsoluteRepoPaths
- a rejected kind is absent from counts rather than zero
  Proof: internal/link/scan_test.go TestCountsOmitsRejectedKinds
- tallying and filtering preserve scan order
  Proof: internal/link/state_test.go TestTallyAndByState
- attention is a property of the state
  Proof: internal/link/state_test.go TestNeedsAttention
