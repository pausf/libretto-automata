# Repo sync — delta

Targets: repo-sync

Two additions to what `internal/repo` answers: **where the clone comes from when there
isn't one**, and **whether a newer release exists**.

Both are git questions, both belong behind the `Git` seam that already makes the update
flow testable against a fake, and neither one changes the pull.

## Outcomes

### Clone

- **`Clone(url, dest)` creates the clone the tool needs.** It runs real `git clone`, for
  the same reason `Pull` does: the user's credential helper, ssh agent and proxy are the
  only implementation guaranteed to work everywhere.
- **It refuses a destination that already has something in it.** Not `--force`, not
  merge-into, not "it looked empty enough". Reported, never resolved — the same promise
  the linker makes. A missing destination is created; an existing empty directory is
  fine.
- **The URL is derived, not hardcoded.** `debug.ReadBuildInfo().Main.Path` is the module
  path, and `https://<module path>.git` is the clone URL. A constant beside a module path
  is two declarations of one fact, and the one that drifts is the one nobody edits. When
  build info is absent the literal `github.com/pausf/libretto-automata` is the fallback,
  and it is named in exactly one place.

### Newer release

- **`LatestTag()` reports the highest release tag on the remote**, from
  `git ls-remote --tags`. No network library, no API token, no rate limit — and it works
  against any remote the user's git can reach, including a fork.
- **Only plain `vX.Y.Z` tags count as releases.** A prerelease or a tag that is not
  semver is ignored rather than ranked. This repository tags plain semver and says so in
  `AGENTS.md`; ranking prerelease correctly is code with no caller.
  *Ceiling:* the first `v1.0.0-rc.1` that is meant to be offered needs real prerelease
  ordering here. Marked `ponytail:` at the comparison.
- **Comparison is numeric, per field.** `v0.10.0` is newer than `v0.9.0`. String
  comparison gets that backwards, which is the whole reason this is not one line.
- **A version that cannot be parsed is never "older".** `dev`, `v0.2.0-3-gabc123`, empty
  — none of them produce a notice. Telling someone with an unidentifiable binary that
  they are out of date is a guess presented as a fact.
- **The call cannot hang.** A context with a deadline, and a deadline that expires is the
  same as offline: no answer, no error surfaced to the user.
- **The answer is cached, and failure is cached too.** `.git/libretto-update-check` holds
  the tag and the time it was asked, and it is not consulted again inside the TTL. Caching
  only successes means a machine with no network asks on every single launch, which is the
  hang risk arriving once per invocation instead of once a day.
- **The cache lives inside `.git/`.** Never committed, removed with the clone, and it does
  not need a `.gitignore` entry that someone has to remember.

## Scope boundaries

**In:** `Clone`, `LatestTag`, the semver comparison, the check cache and its TTL.

**Out:**

- **deciding when to check, or what to say.** This package answers; `cmd/libretto` and
  `internal/ui` decide. A package that shells to git and also owns presentation is a
  package that cannot be tested without one of the two.
- **anything but https for the clone.** An ssh URL depends on a key this tool cannot see
  the state of; the user who wants one sets `LIBRETTO_ROOT` at a clone they made
  themselves.
- **the GitHub API, the Go module proxy, and releases.** `git ls-remote` needs no auth for
  a public repo, no JSON, and no second opinion about which versions exist.
- **fetching or checking out the newer tag.** `update` already fast-forwards. A second
  path to the same commit is a second thing to get wrong.
- **auto-update.** Nothing here moves the user's version. See the `cli` delta.

## Constraints

- Shell out to git. Unchanged, and for the reasons already in the capability spec.
- Every new call goes on the `Git` interface, so a caller can be handed something else.
  **There is no fake, and there is not meant to be one** — `git_test.go` says why:
  `Shell` exists so the invocation lives in one place, and replacing it in tests would
  prove the fake works. An earlier draft of this delta said "so the fake covers it" and
  added a task to extend one, which described a type this repository has never had.
- What the constraint is actually protecting is that **no test reaches the network**, and
  that is met by other means: `LatestTag` runs against a local path used as `origin`, and
  `checkedLatest` takes its asker and its clock as parameters. The `internal/repo` suite
  runs with no network at all.
- No new dependency. `golang.org/x/mod/semver` would do the comparison and would be the
  sixth direct dependency for fifteen lines — that is the ladder's fourth rung losing to
  its fifth.
- `Clone` is not a method on `Shell{Root}`: there is no root yet. It is a function taking
  the destination.

## Prior decisions

- **Bootstrapper, not an embedded payload.** Asked and answered this session: `go install`
  delivers the binary, the binary gets a clone, and links point into the clone exactly as
  they do today. Embedding the payload would end "edit a skill and see it live", which is
  how this payload is developed.
- **`--ff-only` stays.** Diverged histories remain the user's to resolve.
- **A dirty tree still stops `update`.** Unchanged, and a newer release available does not
  weaken it.

## Task breakdown

1. `Clone(ctx, url, dest)` in `internal/repo/clone.go`, refusing a non-empty destination.
2. `ModuleURL()` deriving the clone URL from build info, with the single named fallback.
3. `parseSemver` / `newer` — plain `vX.Y.Z` only, numeric per field, unparseable is never
   older.
4. `LatestTag(ctx)` on `Git`, from `git ls-remote --tags`, with a deadline.
5. The check cache in `.git/libretto-update-check`: read, TTL, write on both success and
   failure.
6. ~~Extend the test fake with `LatestTag`.~~ There is no fake — see constraints.

## Verification criteria

```
Proof: internal/repo/clone_test.go TestCloneRefusesNonEmptyDestination
Proof: internal/repo/clone_test.go TestCloneCreatesMissingDestination
Proof: internal/repo/clone_test.go TestModuleURLDerivesFromBuildInfo
Proof: internal/repo/release_test.go TestLatestTagPicksHighestPlainSemver
Proof: internal/repo/release_test.go TestLatestTagIgnoresPrereleaseAndNonSemverTags
Proof: internal/repo/release_test.go TestNewerComparesFieldsNumerically
Proof: internal/repo/release_test.go TestNewerIsFalseForUnparseableRunningVersion
Proof: internal/repo/release_test.go TestLatestTagHonoursDeadline
Proof: internal/repo/release_test.go TestCheckCacheSuppressesCallsInsideTTL
Proof: internal/repo/release_test.go TestCheckCacheRecordsFailureSoOfflineDoesNotRetry
```

`TestCloneCreatesMissingDestination` and `TestLatestTagPicksHighestPlainSemver` touch
real git and belong behind the same `-short` skip the existing integration tests use.
