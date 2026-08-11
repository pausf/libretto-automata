# distribute-via-go-install — plan

**This file is live state.** A box is marked the moment its task is verified, never
batched. The orchestrator is the only writer; sub-agents report and the box gets marked
here.

**Goal:** `go install` becomes a supported way to get and update this tool, and the panel
says when a newer release exists.

**Architecture:** the bootstrapper shape — the installed binary clones the repo to
`~/.libretto-automata` and links from there, so `update` stays `git pull` + rebuild +
relink. `internal/repo` answers "clone this" and "what is the latest tag"; `cmd/libretto`
decides when to ask; `internal/ui` renders one row and performs no I/O.

**Specs:** `spec-repo-sync.md` · `spec-cli.md` · `spec-panel.md`

## Global constraints

- Go 1.26.5. **No new dependency** — five direct ones, and `AGENTS.md` puts adding a
  sixth on the ask-first list.
- Every gate before every commit: `make gates`.
- `CLAUDE_HOME` at a temp dir in anything touching a target. **`LIBRETTO_ROOT` at a temp
  dir in anything that could bootstrap** — no test clones to a real `~/.libretto-automata`.
- `internal/ui` reads no file, opens no socket, runs no subprocess.
- Every new call on the `Git` interface, so the fake covers it.
- Semver ordering exists in exactly one place: `internal/repo`.
- `ponytail:` comments in English.

## Files

| File | Responsibility |
|---|---|
| `internal/repo/release.go` *(new)* | semver parse/compare, `LatestTag`, the check cache |
| `internal/repo/release_test.go` *(new)* | its proof |
| `internal/repo/clone.go` *(new)* | `Clone`, `ModuleURL` |
| `internal/repo/clone_test.go` *(new)* | its proof |
| `internal/repo/git.go` | `LatestTag` on the interface |
| `cmd/libretto/root.go` *(new)* | `repoRoot` moved out of `main.go` and grown a rung |
| `cmd/libretto/root_test.go` *(new)* | resolution order |
| `cmd/libretto/bootstrap.go` *(new)* | announce, clone, refuse a foreign destination |
| `cmd/libretto/version.go` *(new)* | ldflags → build info → `dev` |
| `cmd/libretto/main.go` | `rebuild` destination, `doctor`'s release line, usage, wiring |
| `internal/ui/panel.go` | `UpdateNotice` field and its row, both layouts |
| `internal/ui/model.go` | `WithReleaseCheck`, `Init`, the message |
| `internal/ui/notice_test.go` *(new)* | the seam and the message |
| `README.md` | `go install` as the install path |

---

## Can start now

**T1, T2, T3, T4, T5, T6** — six independent tasks, no shared file between any two.

---

### T1 · semver: parse and compare

**Spec:** `spec-repo-sync.md` → outcomes, *Newer release*
**Files:** create `internal/repo/release.go`, `internal/repo/release_test.go`
**Blocked by:** nothing
**Produces:** `parseSemver(string) ([3]int, bool)` ·
`IsNewer(latest, running string) bool` — exported, because the CLI formats the row and
must not own a second answer to the same question

- [x] failing tests: field-numeric ordering (`v0.10.0` > `v0.9.0`), prerelease and
      non-semver tags rejected by `parseSemver`, unparseable *running* version never
      yields true
- [x] run them, watch them fail — `undefined: IsNewer`, `undefined: parseSemver`
- [x] implement; `ponytail:` comment naming the prerelease ceiling at the comparison
- [x] `make gates`, commit — exit 0, 223 citations resolve · `7e0…` *(T1 closed)*

Added beyond the plan: `TestParseSemverAcceptsOnlyPlainReleases` (13 cases) and
`TestNewerIsFalseForUnparseableLatest`. A digit-count guard went in too — a tag is a
remote-controlled string, and an overflowing field compares as negative, which reads as
*older* and swallows the notice silently.

**Closes:** `TestNewerComparesFieldsNumerically` ·
`TestNewerIsFalseForUnparseableRunningVersion` ·
`TestLatestTagIgnoresPrereleaseAndNonSemverTags` (the parse half)

---

### T2 · `Clone` and `ModuleURL`

**Spec:** `spec-repo-sync.md` → outcomes, *Clone*
**Files:** create `internal/repo/clone.go`, `internal/repo/clone_test.go`
**Blocked by:** nothing
**Produces:** `Clone(ctx, url, dest string) error` · `ModuleURL() string`

- [x] failing tests: refuses a destination with anything in it; creates a missing
      destination; an existing empty directory is accepted; `ModuleURL` derives
      `https://<module path>.git` from build info
- [x] run them, watch them fail — `undefined: Clone`, `undefined: moduleURL`
- [x] implement — `exec.CommandContext`, the module-path fallback literal in exactly one
      place
- [x] mark the real-git test `-short`-skippable, matching `git_test.go` — verified:
      `go test -short -run TestClone` passes with the git-backed cases skipped
- [x] `make gates`, commit — exit 0 *(T2 closed)*

`ModuleURL()` reads the build info; `moduleURL(info)` takes it, so the fallback branch is
testable at all. Deriving the URL also means a fork installed from its own module path
bootstraps from the fork — behaviour worth having, and free.

**Closes:** `TestCloneRefusesNonEmptyDestination` · `TestCloneCreatesMissingDestination` ·
`TestModuleURLDerivesFromBuildInfo`

---

### T3 · `repoRoot` requires `.git`

**Spec:** `spec-cli.md` → outcomes, *The clone is found, or made*
**Files:** create `cmd/libretto/root.go`, `cmd/libretto/root_test.go`; modify
`cmd/libretto/main.go` (remove `repoRoot`)
**Blocked by:** nothing
**Produces:** `repoRoot() (string, error)` with the four-rung order, and
`bootstrapPath()` as an overridable seam so no test can reach a real `$HOME`

**This is the task that makes the rest correct.** The module cache has a `go.mod` and
would otherwise win.

- [x] failing tests: a `go.mod` without `.git` is rejected; `LIBRETTO_ROOT` wins over
      everything; nothing found falls back to the bootstrap path
- [x] run them, watch them fail — `undefined: isRepo`, `EnvRoot`, `BootstrapDir`
- [x] move `repoRoot` into `root.go`, change the probe from `go.mod` to `.git`, add the
      bootstrap rung
- [x] `make gates`, commit — exit 0 *(T3 closed)*

**The first attempt failed, and the failure was the design's, not the test's.**
`TestRepoRootFallsBackToBootstrapPath` could not pass: `runtime.Caller` is fixed at
compile time, so inside this repository rung 2 always wins and rungs 3 and 4 are
unreachable — the exact rungs that exist only for a binary living somewhere else. Split
into `resolveRoot(override, compileTime, wd, home)` taking its inputs, with `repoRoot()`
gathering them. Without that seam the `go install` behaviour is the part with no proof.

Two rungs went in rather than one. The working directory is accepted only when its
`go.mod` names **this** module — any git repository satisfied the old `$PWD` fallback, so
`libretto install` inside an unrelated project would look for a payload that project does
not have. A `.git` *file* counts as a repo too: that is what a worktree has, and refusing
it would break the tool exactly where this flow's own advice puts you.

**Closes:** `TestRepoRootRequiresGitDirectory` · `TestRepoRootPrefersEnvOverride` ·
`TestRepoRootFallsBackToBootstrapPath`

---

### T4 · version from build info

**Spec:** `spec-cli.md` → outcomes, *The version is knowable without ldflags*
**Files:** create `cmd/libretto/version.go`, `cmd/libretto/version_test.go`; modify
`cmd/libretto/main.go` (the `version` var keeps its ldflags target)
**Blocked by:** nothing
**Produces:** `resolveVersion(ldflags string, info *debug.BuildInfo) string`

Pure, taking the build info rather than reading it, so both branches are testable — a
function calling `debug.ReadBuildInfo()` directly cannot be tested for the fallback.

- [x] failing tests: a set ldflags value wins over build info; `dev` falls through to
      build info; neither means `dev`
- [x] run them, watch them fail — `undefined: resolveVersion`
- [x] implement, call it once at startup
- [x] `make gates`, commit — exit 0 *(T4 closed)*

`(devel)` falls through to `dev` too — it is what the toolchain records for a build from a
working tree, not a version. `buildVersion(stamped)` wraps `resolveVersion` and is called
once in `main`, into the same package var everything already reads, so the footer and
`libretto version` both pick it up with no other change.

**Closes:** `TestVersionPrefersLdflagsOverBuildInfo` · `TestVersionFallsBackToBuildInfo`

---

### T5 · the notice row

**Spec:** `spec-panel.md` → all outcomes except the seam
**Files:** modify `internal/ui/panel.go`; modify `internal/ui/panel_test.go`,
`internal/ui/fluid_test.go`
**Blocked by:** nothing
**Produces:** `Panel.UpdateNotice string`

- [x] failing tests: rendered as its own row between the menu and the strip; absent when
      empty; present in `renderNarrow`; the fluid frame does not tear with it set
- [x] run them, watch them fail — `p.UpdateNotice undefined`
- [x] add the field and the row, attention colour, both layouts
- [x] `make gates`, commit — exit 0 *(T5 closed)*

The row is elided to the content area like a report line, so a long notice cannot tear the
frame — `TestFrameHoldsWithUpdateNotice` is the guard. `TestUpdateNoticeAndActionFeedbackCoexist`
replaced the planned `TestActionFeedbackDoesNotOverwriteUpdateNotice`: at this layer the
claim is that both render, which is what the two-fields decision buys. The model-level
version of it lands in T10.

**Closes:** `TestPanelRendersUpdateNoticeBetweenMenuAndStrip` ·
`TestPanelOmitsUpdateNoticeWhenEmpty` · `TestNarrowLayoutKeepsUpdateNotice` ·
`TestFrameHoldsWithUpdateNotice`

---

### T6 · rebuild replaces the running binary

**Spec:** `spec-cli.md` → outcomes, *The rebuild replaces the binary that is running*
**Files:** modify `cmd/libretto/main.go` (`rebuild`); create
`cmd/libretto/update_test.go`
**Blocked by:** nothing
**Produces:** `rebuild(root, exe string) error` — the destination passed in, not read
from `os.Executable()` inside, so the test does not have to be the binary under test

**The hole the bootstrapper opens.** Without this, `update` reports success and every
later invocation stays on the old version.

- [x] failing tests: writes to the given executable path, not `bin/libretto`; a symlinked
      path is resolved and the link survives as a link; an unwritable destination is
      reported and `update` still succeeds
- [x] run them, watch them fail — `too many arguments in call to rebuild`
- [x] implement — `filepath.EvalSymlinks`, temp file, atomic rename, the unwritable branch
      reporting where the new binary is
- [x] `make gates`, commit — exit 0 *(T6 closed)*

**The unwritable branch failed on the first attempt, and the fix is a real finding.** A
locked destination surfaces as a `go build` error *string*, not a wrapped
`os.ErrPermission` — the temp file is created by the compiler, inside the directory that
refused it. So `rebuildOrReport` could not tell "you cannot write there" from "the code is
broken", and only after paying three seconds for the compile. The write is probed with
`os.OpenFile` first, which also means the failure is instant.

The temp file stays beside the destination rather than in `TMPDIR`: rename across
filesystems fails, and `$GOBIN` and `/tmp` are routinely on different ones.

One test assertion was wrong too — "nothing was written to the clone's `bin/`" cannot be
an existence check, because a development checkout has usually run `make build`. Compared
by modification time instead.

**Closes:** `TestRebuildReplacesRunningExecutable` · `TestRebuildResolvesSymlinkedExecutable`
· `TestRebuildReportsUnwritableDestinationWithoutFailing`

---

### T7 · `LatestTag`

**Spec:** `spec-repo-sync.md` → outcomes, *Newer release*
**Files:** modify `internal/repo/release.go`, `internal/repo/git.go`,
`internal/repo/release_test.go`; modify the fake in `internal/repo/git_test.go`
**Blocked by:** **T1**
**Consumes:** `parseSemver`, `newer`
**Produces:** `LatestTag(ctx) (string, error)` on `Git`

- [x] failing tests: picks the highest plain semver from `git ls-remote --tags` output;
      an expired deadline returns no answer and no user-facing error
- [x] run them, watch them fail — `undefined: checkedLatest`
- [x] implement, add to the interface, ~~extend the fake~~
- [x] `make gates`, commit — exit 0 *(T7 closed, together with T8)*

**There is no fake.** The plan said to extend one; `Git`'s only implementation is `Shell`,
and `git_test.go` says why in a comment — replacing it in tests would prove the fake works.
`LatestTag` is tested against a real local repository used as a remote, which needs no
network. Nothing to extend, so nothing was.

**A comment claimed something the test did not prove.** `highestTag` strips the peeled
`^{}` ref, and the first test used lightweight tags — which never emit one. Verified
against real git that annotated tags emit two lines per tag, then switched the fixture to
`git tag -a`, which is what `AGENTS.md` says releases are anyway.

T8's tests share the file, so the two landed in one commit rather than committing a file
whose tests do not build.

**Closes:** `TestLatestTagPicksHighestPlainSemver` ·
`TestLatestTagIgnoresPrereleaseAndNonSemverTags` · `TestLatestTagHonoursDeadline`

---

### T8 · the check cache

**Spec:** `spec-repo-sync.md` → outcomes, *Newer release* (last two bullets)
**Files:** modify `internal/repo/release.go`, `internal/repo/release_test.go`
**Blocked by:** **T7**
**Produces:** `CheckedLatest(ctx, root string, ttl time.Duration) (string, error)` —
cache-aware, and the only entry point the CLI calls

- [x] failing tests: no second call inside the TTL; a failure is cached too, so an offline
      machine does not retry every launch
- [x] run them, watch them fail
- [x] implement over `.git/libretto-update-check`
- [x] `make gates`, commit — exit 0, same commit as T7 *(T8 closed)*

`CheckedLatest(ctx, root, ttl)` is the exported entry point; `checkedLatest` takes the clock
and the asker so neither the wall clock nor the network is in the test. Two cases went in
beyond the plan: the TTL expiring asks again, and **no `.git` means ask without caching**
rather than fail — that is the bootstrap case, where there is nowhere to write yet.

**Closes:** `TestCheckCacheSuppressesCallsInsideTTL` ·
`TestCheckCacheRecordsFailureSoOfflineDoesNotRetry`

---

### T9 · bootstrap

**Spec:** `spec-cli.md` → outcomes, *The clone is found, or made*
**Files:** create `cmd/libretto/bootstrap.go`, `cmd/libretto/bootstrap_test.go`; modify
`cmd/libretto/main.go` (call it when `repoRoot` finds nothing)
**Blocked by:** **T2, T3**
**Consumes:** `repo.Clone`, `repo.ModuleURL`, `bootstrapPath`

- [x] failing tests: the destination is printed **before** the clone runs; a destination
      that exists and is not our clone is refused and nothing is touched; after a
      successful clone the requested command runs
- [x] run them, watch them fail — `undefined: bootstrap`
- [x] implement; `LIBRETTO_ROOT` in every test so no real `$HOME` is reachable
- [x] `make gates`, commit — exit 0 *(T9 closed)*

Two cases beyond the plan, both real. **A failed clone is cleaned up** — half a clone in
`~/.libretto-automata` would be a foreign destination forever, and the user would have to
work out that deleting it is the fix. And **`version`/`help` are answered before the clone
is even looked for**, with a test asserting the home directory stays empty: cloning a
repository into somebody's home because they asked what version they were running would be
indefensible. That meant moving those two cases above the root resolution in `run`.

`bootstrap` takes the clone function and an `io.Writer`, so no test needs git or a network.
`ensureClone()` is what commands call now instead of `repoRoot()`.

**Closes:** `TestBootstrapAnnouncesDestinationBeforeCloning` ·
`TestBootstrapRefusesForeignDestination` · `TestBootstrapContinuesIntoRequestedCommand`

---

### T10 · the panel's seam

**Spec:** `spec-panel.md` → outcomes, last two bullets, and task 4
**Files:** modify `internal/ui/model.go`; create `internal/ui/notice_test.go`
**Blocked by:** **T5**
**Produces:** `WithReleaseCheck(func() string) Model` — **returns the finished row or
empty**, so the comparison stays in `internal/repo` · `releaseMsg string` ·
`Init() tea.Cmd`

- [x] failing tests: `Init` returns a command when the check is set and `nil` when it is
      not; the message sets `UpdateNotice`; an action's feedback does not overwrite it
- [x] run them, watch them fail — `undefined: releaseMsg`
- [x] implement
- [x] `make gates`, commit — exit 0 *(T10 closed)*

`TestNavigationDoesNotClearUpdateNotice` went in as well: moving the cursor clears action
feedback — existing behaviour — and must not take the news with it.

**Closes:** `TestInitReturnsReleaseCheckCommand` · `TestUpdateNoticeSetFromMessage` ·
`TestActionFeedbackDoesNotOverwriteUpdateNotice` ·
`TestNoReleaseCheckMeansNoCommandAndNoNotice`

---

### T11 · wire the panel's check

**Spec:** `spec-cli.md` → *The panel uses the cache*; `spec-panel.md` → *Cached, not live*
**Files:** modify `cmd/libretto/main.go`
**Blocked by:** **T8, T10**
**Consumes:** `repo.CheckedLatest`, `repo.IsNewer`, `resolveVersion`

The callback composes them and formats `v0.2.0 → v0.3.0 available · choose update`.

- [x] pass the callback into `NewModel(...).WithReleaseCheck(...)`
- [x] verify against the existing teatest flow in `cmd/libretto/panelrun_test.go`: the
      panel paints without the check having answered — those flows build the model
      directly, so the suite never reaches a remote
- [x] `make gates`, commit — exit 0 *(T11 closed)*

**It earned criteria after all.** The plan called this wiring with no proof of its own, and
that was wrong: `releaseNotice` decides when to stay silent, and silence has five distinct
cases worth pinning — up to date, ahead of the remote, unparseable, a dirty build, nothing
cached. Tests build the notice from a pre-seeded cache file, so nothing reaches a network.

```
Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesBothVersionsAndTheAction
Proof: cmd/libretto/version_test.go TestReleaseNoticeIsSilentWhenThereIsNothingToSay
```

---

### T12 · `doctor` checks live

**Spec:** `spec-cli.md` → *The user is told a newer release exists*
**Files:** modify `cmd/libretto/main.go` (`doctor`); create `cmd/libretto/doctor_test.go`
**Blocked by:** **T7, T8**

- [x] failing tests: a newer release is reported; a failed check says "could not check"
      and does not set the exit code
- [x] run them, watch them fail — `undefined: releaseLine`
- [x] implement with the deadline; ~~the cache is refreshed~~ **nothing is written at all**
- [x] `make gates`, commit — exit 0 *(T12 closed)*

**The spec changed here, and for the better.** It said `doctor` refreshes the check cache,
which meant writing a timestamp into `.git/` and then defending in prose why "it never
writes" still held. Going live instead removes the caveat *and* fixes a second problem: the
cache swallows the ask error by design, so a cached `doctor` could not tell "up to date"
from "could not check" — the exact distinction the line exists to make. The `cli` delta was
amended to match.

Five outcomes, five sentences, and the fifth is the interesting one: a binary ahead of the
remote or reporting `dev` gets the facts and **no ranking**.

**Closes:** `TestDoctorReportsNewerRelease` · `TestDoctorSaysSoWhenTheCheckFails`

---

### T13 · say so

**Spec:** `spec-cli.md` → scope, *In*
**Files:** modify `cmd/libretto/main.go` (`usage`), `README.md`
**Blocked by:** **T9, T11, T12**

- [x] usage: the `go install` line, `LIBRETTO_ROOT`'s default documented as
      `~/.libretto-automata` — taken from `BootstrapDir` rather than retyped, so the help
      cannot drift from the code that creates the directory
- [x] `README.md`: `go install` as the install and update path, the clone location, and
      that `update` is still how you move forward
- [x] `scripts/check-payload` — exit 0, all checks passed
- [x] `make gates`, commit — exit 0 *(T13 closed)*

The clone survives in the README as the route for working *on* the payload, with the reason
attached: a clone you are standing in wins over `~/.libretto-automata`, and that is what
keeps edit-a-skill-and-see-it-live working.

**Closes:** no test. Documentation that agrees with the binary is checked by reading both,
and this is the one task in the plan whose proof is a human.

---

## Landing

Phase 8's job, listed here so it is not forgotten: apply the three deltas onto
`.agents/specs/repo-sync/spec.md`, `cli/spec.md`, `panel/spec.md`, delete this change
folder, and confirm `spec-drift --anchors` resolves the **34** citations that arrive with
them.

**The count was wrong here, and the reviewer caught why it mattered.** It said 31, and one
of the citations named a test that did not exist — `TestBootstrapContinuesIntoRequestedCommand`.
`--anchors` does not scan change deltas, so nothing would have failed until the delta landed
in `.agents/specs/cli/spec.md` and broke the gate there. Worse, `go test -run` on a name with
no match exits 0 silently, so the citation read as satisfied. The test was written; all 34
now resolve, verified by hand.
