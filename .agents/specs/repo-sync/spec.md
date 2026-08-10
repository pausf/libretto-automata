# Repo sync

Governs: internal/repo/**

Refreshing the repository the tool lives in, and deciding when the compiled binary has
been invalidated by what arrived.

## Outcomes

`update` pulls, relinks, and reports — in that order, and only when pulling is safe.

- **A dirty working tree stops it.** No pull, no links changed, the reason stated, exit
  non-zero. Losing the user's uncommitted work is worse than being out of date, and no
  flag exists to argue otherwise.
- **No remote is not an error.** The pull is skipped with a clear message and relinking
  still happens, because relinking is useful on its own.
- **A pull that carried Go changes triggers a rebuild.** A pull carrying only markdown
  does not, which is most pulls in a project whose payload is prose.
- **After a rebuild, the running process says so.** It is still the old binary; the
  upgrade takes effect next invocation. Pretending otherwise would be a claim about
  code that is not executing.

## Scope boundaries

**In:** working-tree cleanliness, remote presence, HEAD, fast-forward pull, changed
paths, the rebuild decision.

**Out:**

- **merges.** The pull is `--ff-only`. Diverged histories are the user's to resolve, in
  their own shell, seeing what they are doing. A merge commit created by a background
  command is a merge commit nobody chose.
- pushing, branching, committing, stashing. Read and fast-forward only.
- linking — that is `linking`, composed after the pull.
- authentication. Handled by the user's git, deliberately (see constraints).

## Constraints

**Shell out to git; do not link a git library.** The pull has to work with the user's
ssh agent, credential helper, signing keys, proxy and `.gitconfig`, and the only
implementation guaranteed to honour all of those is the git that made the repository.
A library covers a subset, and the day it fails on a credential real git resolves, the
bug is unfixable from inside this program.

**Everything behind one interface**, so the update flow is exercisable against a fake —
no network, no temporary repository, no flakiness in the suite.

**A repository with no commits is dirty.** Not a corner case to tidy away: every file
is untracked, so a pull could not be reconciled with anything.

**No commits is a state, not a failure.** `Head()` returns empty rather than an error,
because a fresh repository is a legitimate thing to be.

**The rebuild writes to a temporary file and renames.** Writing over the running
executable produces "text file busy", and a half-written binary is worse than a stale
one. Rename is atomic.

## Prior decisions

- The interface carries five questions and nothing more. Anything else the update flow
  turns out to need is a change to this spec first.
- `NeedsRebuild` is a pure function over paths, so the decision is testable without a
  repository. `.go`, `go.mod` and `go.sum` invalidate the binary; nothing else does.

## Task breakdown

- [x] the `Git` interface and its shell implementation
- [x] `NeedsRebuild`
- [x] the update flow composed in `cli`
- [ ] **a fake `Git` and the flow's own tests.** The interface exists precisely to make
      this possible and it has not been used yet.

## Verification criteria

- only Go sources and module files invalidate the binary
  Proof: internal/repo/git_test.go TestNeedsRebuild

**This capability is under-tested and the gap is named rather than hidden.** One test
covers 155 lines. The interface was built to be faked and no fake exists, so none of
the outcomes above — the dirty-tree refusal, the missing-remote path, the rebuild
notice — is currently held up by anything. Every one of them is a promise this
document makes and the suite does not keep.

Criteria still owed a test:

- a dirty tree refuses to pull and changes no links
- a missing remote skips the pull and relinks anyway
- a pull carrying only markdown does not rebuild
- a pull carrying Go sources rebuilds, and the notice is printed
- a failed rebuild after a successful pull reports both facts
