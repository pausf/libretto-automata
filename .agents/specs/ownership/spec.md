# Ownership

Governs: internal/link/own.go internal/link/own_test.go

Deciding whether a symlink in a target directory belongs to this repository.

**This is the single most consequential predicate in the program.** Say "yes" about
somebody else's link and `install` overwrites their work. Say "no" about our own and
the tool can never repair anything. Every write in `linking` is gated on it.

## Outcomes

`Owned(repoRoot, path)` answers one question with no ambiguity: does this path name a
symlink whose destination resolves inside `repoRoot`?

- **True** only when the path is a symlink *and* its destination lies within the repo.
- **False** for everything else: a regular file, a directory, a symlink leaving the
  repo, a path that does not exist.
- **True even when the destination is gone.** An owned link whose source was deleted
  is still ours, and recognising that is the only thing that makes `prune` possible.

`LinkTarget(path)` reports where a symlink resolves, normalised, and whether it was a
symlink at all — so callers can distinguish "not ours" from "not a link".

## Scope boundaries

**In:** the predicate, symlink reading, path normalisation, containment testing.

**Out:**

- deciding what to *do* about an answer — that is `linking`
- classifying items into states — that is `link-state`
- any filesystem write. This file never mutates anything.
- permissions. Whether we *may* write is a different question from whether we *should*.

## Constraints

Conservative by construction: **anything it cannot prove is ours is treated as
foreign.** The cost of a false negative is a repair that does not happen and a
message the user can read. The cost of a false positive is destroyed work.

Three properties are not optional, each one a way real filesystems differ from the
naive model:

1. **Relative links resolve against the link's own directory**, never the process
   working directory. Getting this backwards is the classic way to misjudge
   ownership.

2. **Paths are normalised through symlinks before comparison.** On macOS `/tmp` is
   `/private/tmp` and `t.TempDir()` sits under a symlinked `/var`, so one directory
   has two spellings. A naive compare calls our own links foreign. Normalisation must
   work on paths that do not exist, because a broken link is still ours — resolve as
   far as possible and append the unresolvable remainder verbatim.

3. **Containment compares path segments, not string prefixes.** `strings.HasPrefix`
   places `/repo-backup` inside `/repo`. That is how a tool deletes from the wrong
   tree.

## Prior decisions

- Absolute or relative links are both accepted. The repo writes absolute ones, but a
  link a human made by hand is still ours if it resolves inside.
- No allow-list, no marker file, no extended attribute. Ownership is derived from
  where the link points, because any marker can be lost, copied, or forged, and a
  derived answer cannot go stale.
- `within(root, root)` is true. A link pointing at the repo root itself is ours.

## Task breakdown

Complete. This capability shipped in phase 2 and has not needed a change since.

## Verification criteria

Each of the three constraints above has a dedicated test, because each one is a bug
that a plausible implementation would contain.

- the predicate distinguishes ours, foreign, absent and broken
  Proof: internal/link/own_test.go TestOwned
- a sibling directory sharing a name prefix is not inside
  Proof: internal/link/own_test.go TestOwnedRejectsPrefixSiblings
- a repo reached by an aliased path is still recognised
  Proof: internal/link/own_test.go TestOwnedWhenTheRepoPathIsItselfSymlinked
- a link written through an aliased path is still recognised
  Proof: internal/link/own_test.go TestOwnedWhenTheLinkUsesAnAliasedPath
- normalisation survives components that do not exist
  Proof: internal/link/own_test.go TestNormaliseHandlesMissingComponents
- containment is by segment
  Proof: internal/link/own_test.go TestWithin
- a non-symlink is reported as such rather than guessed at
  Proof: internal/link/own_test.go TestLinkTarget
