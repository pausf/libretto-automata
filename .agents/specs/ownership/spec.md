# Ownership

Governs: internal/link/own.go internal/link/own_test.go

Deciding whether an entry in a target directory belongs to this repository — a
symlink by its destination, a generated file by a marker it carries.

**This is the single most consequential predicate in the program.** Say "yes" about
somebody else's link and `install` overwrites their work. Say "no" about our own and
the tool can never repair anything. Every write in `linking` is gated on it.

## Outcomes

**`Owned(repoRoot, path)` is symlink-only, and `OwnedEither` is the widened question.**
Two predicates, not one, and the split is a safety property rather than tidiness — see
below.

- **A symlink** is ours when its destination resolves inside `repoRoot`. **True even
  when the destination is gone** — an owned link whose source was deleted is still ours,
  and recognising that is the only thing that makes `prune` possible. `Owned` answers
  this and nothing else.
- **A regular file** is ours when its frontmatter carries `x-libretto-source: <path>`
  and that path resolves **strictly** inside `repoRoot` — never the root itself, never a
  directory. `OwnedGenerated` is that arm alone; `OwnedEither` is the union.
- **The marker arm is asked only for a kind whose target installs it by transform.** With
  one predicate covering both, a review found `prune --claude --yes` offering to delete a
  hand-written file in a Claude destination that merely carried a marker line. **A
  destination that never writes a generated file cannot own one**, so asking there adds a
  way to destroy work and buys nothing. `link-state` and `linking` carry the kind on the
  entry and ask the narrower question everywhere else.
- **False for everything else**: a file with no marker, a directory, a symlink leaving
  the repo, a path that does not exist, a file that cannot be read.

**The marker arm is deliberately narrower than the symlink arm, not looser.** A symlink
is ours when its destination lands anywhere in a whole subtree; a generated file is ours
only when it names an exact path by an exact key. **A file with no marker is foreign,
always** — whatever its name, whatever its content.

Three things bound it, and each is a criterion below:

- **Only the frontmatter block is read** — between the opening `---` on line 1 and the
  next `---`. A marker key in an agent's prose is prose. Reading the whole file would let
  anybody grant us ownership of their file by typing one line into it.
- **The value must be absolute.** A symlink may be relative because there is an
  unambiguous base — the directory holding the link. A marker has none: the file sits in
  the target, so resolving a relative source against it would name something inside the
  target rather than the repo. The consequence is stated rather than discovered:
  **the first draft of this sentence got it wrong**: moving the repository makes every
  generated file **foreign**, not stale, so `prune` skips it as somebody else's and can
  never remove it — measured, not reasoned about. Every symlink behaves identically for
  the same reason. **The remedy is `uninstall` from the old checkout before moving it**;
  after the fact it is a manual delete. "Prune is the remedy" would have sent somebody
  to a command that reports nothing and changes nothing.
- **`within` compares path segments, never string prefixes**, so a marker naming
  `/repo-backup/agents/x.md` is not inside `/repo`.

**A marker whose source no longer exists is still ours.** Orphaned is `Stale`, which
`prune` removes; calling it foreign would strand it forever.

`GeneratedSource(path)` reports the marker's source and whether there was one, mirroring
`LinkTarget`.

**The predicate now reads file content, which is new** — it needed only `Lstat` and
`Readlink` before. **Ceiling named:** a symlink cannot lie about its destination and a
marker can be typed by hand. The containment test still bounds what a lie can claim, and
what the marker buys is that ownership is *legible* — the file says where it came from,
in a form a person can read.

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
- **a symlink's ownership answer is unchanged by the generated arm** — reaching it
  requires the path not to be a symlink, so no existing answer can move
  Proof: internal/link/own_generated_test.go TestOwnedSymlinkArmIsUnchanged
- **a regular file carrying a marker inside the repo is ours**
  Proof: internal/link/own_generated_test.go TestOwnedGeneratedFile
- **everything it cannot prove is foreign** — no marker, a marker outside the repo, a
  marker naming a prefix sibling, a marker in the prose body, a marker after the
  frontmatter closes, an unclosed block, no frontmatter, a relative marker, an empty
  marker, a directory, a missing path
  Proof: internal/link/own_generated_test.go TestGeneratedOwnershipRefusesWhatItCannotProve
- **a marker naming a source that no longer exists is still ours**, or `prune` could
  never remove it
  Proof: internal/link/own_generated_test.go TestGeneratedOwnershipSurvivesAMissingSource
- **`GeneratedSource` reports the source and whether a marker was there**
  Proof: internal/link/own_generated_test.go TestGeneratedSourceReporting
- **an unreadable file is foreign rather than an error**
  Proof: internal/link/own_generated_test.go TestUnreadableFileIsForeign
- **the marker arm is never asked about a kind the target does not generate** — `Owned`
  refuses a marked regular file, `OwnedEither` accepts it, and `prune` on a
  non-generating destination spares it
  Proof: internal/link/own_generated_test.go TestOwnedIgnoresTheMarkerForNonGeneratedKinds
  Proof: internal/link/apply_generated_test.go TestPruneSparesAMarkedFileInANonGeneratingDestination
- **a marker naming the repository root, or a directory inside it, is refused** — a marker
  names an item, and both were accepted and then deleted by `prune --yes` before a review
  found them
  Proof: internal/link/own_generated_test.go TestGeneratedOwnershipRefusesTheRootAndDirectories
- **a quoted marker and a bare one both read** — the transform quotes, and refusing a bare
  value would turn a file written before quoting landed into a foreign one nothing can
  clean up
  Proof: internal/link/own_generated_test.go TestGeneratedSourceAcceptsQuotedAndBareMarkers

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
