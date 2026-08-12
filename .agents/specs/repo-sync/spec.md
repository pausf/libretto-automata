# Repo sync

Governs: internal/repo/**

Getting the repository the tool lives in, refreshing it, deciding when the compiled
binary has been invalidated by what arrived, and knowing whether a newer release exists.

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

### It does not get the clone

`Clone` and `ModuleURL` lived here for one session, to bootstrap a checkout for a binary
installed with `go install`. Both are gone.

The clone bootstrap was the wrong shape: it made `update` announce `git pull` to somebody who
only wanted to use the tool. The payload now arrives as a published release asset — see the
`distribution` capability — so nothing clones, and a function with no caller is a function with
tests that prove nothing.

**It shipped in no tag.** `@latest` could never resolve to it, so no promise was withdrawn from
anybody and no version needed a major bump for the removal.

### Knowing a newer release exists

- **`LatestTag` is the highest plain release tag on the remote that is not retracted**,
  read from `git ls-remote --tags`. No token, no rate limit, no JSON, and it works against
  whatever remote the user's git can reach.
- **A tag that retracts itself is not a release, and this had to be taught.** `v1.0.2` here
  has no Release and exists only to carry a `retract` block covering the two bad versions
  before it and itself. It is valid semver, so it sorted highest and was offered as an
  update — measured, on a real clone. **The module proxy honours `retract` and answered
  `v0.7.0` correctly the whole time; a git ref carries no retraction at all**, so the two
  install modes disagreed and only the checkout was wrong.
  - The candidate is opened — `git show <tag>:go.mod` — and dropped if its own `retract`
    directives name it. All three grammar forms are read: a single directive, a
    parenthesised block, and a `[low, high]` range, inclusive at both ends.
  - **Hand-parsed, not through `golang.org/x/mod/modfile`.** A sixth direct dependency for
    one predicate, when the ladder puts stdlib first.
  - **An unreadable `go.mod` is not evidence of a retraction** and the tag is offered.
    `ls-remote` can name a tag the local object store has never seen — one pushed since the
    last fetch — and a tag that new is a release far more often than a tombstone. Hiding it
    would mean a genuine release going unannounced to everyone whose tags are a day old.
  - **No fetch to settle it.** Reaching the network for a speculative background check is
    the hang the cache exists to prevent, arriving from another direction, and it would
    write to the user's object store for a question nobody asked.
  - **The walk down the tags is bounded.** Each step past the first is a `git show` on a
    tag that turned out to be retracted, and an unbounded loop over a remote-controlled
    list is a subprocess per line.
- **Only plain `vX.Y.Z` counts.** A prerelease, a `git describe` string and build metadata
  are rejected rather than ranked. Invisible is the safe direction: a prerelease cannot
  claim to be newer than the release it precedes.
- **Comparison is numeric per field.** `v0.10.0` is ahead of `v0.9.0`; a string sort has
  that backwards, which is the whole reason this is not one line.
- **Either side unparseable is never "newer".** Telling somebody whose binary reports
  `dev` that they are out of date is a guess presented as a fact.
- **The peeled `^{}` ref is stripped.** Releases here are `git tag -a`, and an annotated
  tag makes `ls-remote` emit two lines per tag; the second names `v0.2.0^{}`, which parses
  as nothing.
- **An empty answer with no error means the remote has no release to offer.** Only
  prereleases, or no tags at all, is a state.
- **The call cannot hang**, and **the answer is cached for a day — failure included.**
  Caching only successes means a machine with no network pays the timeout on every launch,
  which is the hang this exists to prevent arriving once per invocation instead of once a
  day. The cache lives in `.git/`, so it is never committed, needs no `.gitignore` line
  anybody has to remember, and goes with the checkout.
- **No `.git` means ask without caching**, rather than fail. That is an installed copy, which
  has no checkout to write into — and `distribution` is what it asks instead.

## Scope boundaries

**In:** working-tree cleanliness, remote presence, HEAD, fast-forward pull, changed
paths, the rebuild decision, the latest release tag, the semver comparison, and the
check cache.

**Out:**

- **merges.** The pull is `--ff-only`. Diverged histories are the user's to resolve, in
  their own shell, seeing what they are doing. A merge commit created by a background
  command is a merge commit nobody chose.
- pushing, branching, committing, stashing. Read and fast-forward only.
- linking — that is `linking`, composed after the pull.
- authentication. Handled by the user's git, deliberately (see constraints).
- **deciding when to check, or what to say about it.** This package answers; `cli` and
  `panel` decide. A package that shells to git and also owns presentation cannot be tested
  without one of the two.
- **cloning.** It lived here for one session and is gone; nothing clones. See above.
- **downloading a release.** That is `distribution`, which has no git in it at all.
- **the GitHub API, the Go module proxy, and releases as an endpoint.** `ls-remote` needs
  no auth for a public repo and no second opinion about which versions exist.
- **fetching or checking out the newer tag.** `update` already fast-forwards; a second
  path to the same commit is a second thing to get wrong.
- **auto-update.** Nothing here moves the user's version.
- **prerelease ordering.** *Ceiling:* the first `v1.0.0-rc.1` meant to be offered needs
  real semver §11 precedence, which is thirty lines that currently prove nothing. Marked
  `ponytail:` at the comparison.

## Constraints

**Shell out to git; do not link a git library.** The pull has to work with the user's
ssh agent, credential helper, signing keys, proxy and `.gitconfig`, and the only
implementation guaranteed to honour all of those is the git that made the repository.
A library covers a subset, and the day it fails on a credential real git resolves, the
bug is unfixable from inside this program.

**Everything behind one interface**, so a caller can be handed something other than
`Shell`. **There is no fake, and there is not meant to be one** — see *What the tests are*
below. This constraint used to read "so the update flow is exercisable against a fake",
which promised a type this repository has never had and contradicted its own verification
section three headings later.

**What the constraint actually protects is that no test reaches the network**, and that
holds: the git-backed tests build a repository in a `t.TempDir()` and use a local path as
`origin`, and `checkedLatest` takes its clock and its asker as parameters.

**No new dependency for the version comparison.** `golang.org/x/mod/semver` would do it and
would be the sixth direct dependency for fifteen lines — the ladder's fourth rung losing to
its fifth.

**A repository with no commits is dirty.** Not a corner case to tidy away: every file
is untracked, so a pull could not be reconciled with anything.

**No commits is a state, not a failure.** `Head()` returns empty rather than an error,
because a fresh repository is a legitimate thing to be.

**The rebuild writes to a temporary file and renames.** Writing over the running
executable produces "text file busy", and a half-written binary is worse than a stale
one. Rename is atomic.

## Prior decisions

- The interface carries six questions and nothing more. Anything else the update flow
  turns out to need is a change to this spec first.
- `NeedsRebuild` is a pure function over paths, so the decision is testable without a
  repository. `.go`, `go.mod` and `go.sum` invalidate the binary; nothing else does.
- **Not an embedded payload, and — after one session — not a bootstrapped clone either.**
  Embedding was rejected first and stays rejected: `//go:embed` has no path on disk, so every
  symlink would become a copy and the ownership model would go with it. The clone bootstrap
  replaced it and was itself replaced by a release asset, because `git pull` is the wrong
  thing to say to somebody who only wanted to use the tool. `distribution` carries that
  decision now.
- **`git ls-remote` over the GitHub API and the Go module proxy**, for the release check.
  No token, no rate limit, no JSON, and it works against a fork.

## Task breakdown

- [x] the `Git` interface and its shell implementation
- [x] `NeedsRebuild`
- [x] the update flow composed in `cli`
- [x] ~~`Clone` and `ModuleURL`~~ removed in the same session — nothing clones
- [x] `LatestTag`, the semver comparison, and the check cache
- [ ] **the flow's own tests.** `update`'s composed behaviour, listed under *Still owed*
      below. Not a fake — the entry that used to say "a fake `Git` and the flow's own
      tests" asked for a type this spec now says should not exist.

## Verification criteria

- only Go sources and module files invalidate the binary
  Proof: internal/repo/git_test.go TestNeedsRebuild

- a repository with an uncommitted change reads as dirty, and a clean one does not
  Proof: internal/repo/git_test.go TestDirtyReportsTheWorkingTree
- a repository with no remote says so rather than failing
  Proof: internal/repo/git_test.go TestHasRemoteDistinguishesNoRemoteFromAnError
- `Head` returns the commit the repository is actually on
  Proof: internal/repo/git_test.go TestHeadIsTheCurrentCommit
- **a repository with no commits gives `Head` an empty answer, not an error** — the one
  call that deliberately swallows a git failure, pinned so a later tidy-up cannot turn
  it into an error and break the empty-repository path in silence
  Proof: internal/repo/git_test.go TestHeadOnARepositoryWithNoCommitsIsEmptyNotAnError
- `ChangedSince` names the paths a commit touched and no others
  Proof: internal/repo/git_test.go TestChangedSinceNamesOnlyWhatChanged
- an empty revision means "nothing to compare against", not a diff against everything
  Proof: internal/repo/git_test.go TestChangedSinceWithNoRevisionIsEmpty
- **`Pull` brings a commit across from a local bare remote** — no network, so the test
  cannot flake on one
  Proof: internal/repo/git_test.go TestPullFetchesFromALocalRemote
- **`Pull` refuses a diverged history**, which is what `--ff-only` is for: a merge
  commit made by a background command is a merge commit nobody chose
  Proof: internal/repo/git_test.go TestPullRefusesADivergedHistory
- outside a repository the reads error rather than answering "clean" and "no remote"
  Proof: internal/repo/git_test.go TestOutsideARepositoryTheReadsError

### The release check

- the highest plain semver wins, from real `ls-remote` output over annotated tags
  Proof: internal/repo/release_test.go TestLatestTagPicksHighestPlainSemver
- **a tag whose own `go.mod` retracts it is skipped** — the tombstone, offered as an update
  until it was not
  Proof: internal/repo/release_test.go TestLatestTagSkipsATagThatRetractsItself
- **a tag whose `go.mod` cannot be read is offered**, because unreadable is not evidence
  Proof: internal/repo/release_test.go TestLatestTagOffersATagWhoseGoModCannotBeRead
- every `retract` form is read: single, block, and an inclusive `[low, high]` range
  Proof: internal/repo/release_test.go TestRetractedReadsEveryDirectiveForm
- a remote whose only tags are prereleases has no release to offer, and that is not an
  error
  Proof: internal/repo/release_test.go TestLatestTagIgnoresPrereleaseAndNonSemverTags
- nor does a remote with no tags at all
  Proof: internal/repo/release_test.go TestLatestTagOnARemoteWithNoTags
- **a cancelled context stops the call** — an unanswering network must not hold the panel's
  first paint
  Proof: internal/repo/release_test.go TestLatestTagHonoursDeadline
- no remote configured cannot be asked
  Proof: internal/repo/release_test.go TestLatestTagWithNoRemote
- ordering is numeric per field, so `v0.10.0` beats `v0.9.0`
  Proof: internal/repo/release_test.go TestNewerComparesFieldsNumerically
- **only plain `vX.Y.Z` parses** — prereleases, `git describe` strings, build metadata and
  `dev` are all rejected rather than ranked
  Proof: internal/repo/release_test.go TestParseSemverAcceptsOnlyPlainReleases
- **a binary that cannot say what it is is never told it is out of date**
  Proof: internal/repo/release_test.go TestNewerIsFalseForUnparseableRunningVersion
- and neither side being parseable settles it
  Proof: internal/repo/release_test.go TestNewerIsFalseForUnparseableLatest
- the remote is asked once per TTL, not once per launch
  Proof: internal/repo/release_test.go TestCheckCacheSuppressesCallsInsideTTL
- and asked again once the TTL expires
  Proof: internal/repo/release_test.go TestCheckCacheAsksAgainOnceTheTTLExpires
- **failure is cached too**, so an offline machine does not pay the timeout every launch
  Proof: internal/repo/release_test.go TestCheckCacheRecordsFailureSoOfflineDoesNotRetry
- no `.git` means ask without caching, not fail
  Proof: internal/repo/release_test.go TestCheckCacheWithoutAGitDirectoryStillAnswers

### What the tests are, and what they are not

**Real git, not a fake.** The interface was built so callers can be handed something other
than `Shell`, and no fake was ever written; the tests build a repository in a `t.TempDir()`
instead, and use a local path as `origin` where a remote is needed. `Shell` exists so the
git invocation lives in one place, and replacing it in tests would prove the fake works.
The cost is tests that need `git` on the machine and are slower than the rest — which is
what `-short` is for.

**Where a seam was needed, it is a parameter and not an interface.** `checkedLatest` takes
its clock and its asker; `moduleURL` takes the build info. A function that reads the wall
clock behaves differently at midnight, and one that calls `debug.ReadBuildInfo()` inside
itself can only be tested for whichever answer the test binary happens to give.

**They are the real-git integration `make test-short` has always claimed to skip.**
That claim was empty before this: every test ran in both modes because nothing in the
repository called `testing.Short()`. Nine skip now. The teatest half of the claim is
still unimplemented.

**Still owed, and not by the tests above.** These are `update`'s behaviour, not
`Shell`'s — they need the command, not the git wrapper:

- a dirty tree refuses to pull and changes no links
- a missing remote skips the pull and relinks anyway
- a pull carrying only markdown does not rebuild
- a pull carrying Go sources rebuilds, and the notice is printed
- a failed rebuild after a successful pull reports both facts
