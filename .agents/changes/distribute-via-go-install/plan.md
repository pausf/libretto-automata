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

- [ ] failing tests: field-numeric ordering (`v0.10.0` > `v0.9.0`), prerelease and
      non-semver tags rejected by `parseSemver`, unparseable *running* version never
      yields true
- [ ] run them, watch them fail
- [ ] implement; `ponytail:` comment naming the prerelease ceiling at the comparison
- [ ] `make gates`, commit

**Closes:** `TestNewerComparesFieldsNumerically` ·
`TestNewerIsFalseForUnparseableRunningVersion` ·
`TestLatestTagIgnoresPrereleaseAndNonSemverTags` (the parse half)

---

### T2 · `Clone` and `ModuleURL`

**Spec:** `spec-repo-sync.md` → outcomes, *Clone*
**Files:** create `internal/repo/clone.go`, `internal/repo/clone_test.go`
**Blocked by:** nothing
**Produces:** `Clone(ctx, url, dest string) error` · `ModuleURL() string`

- [ ] failing tests: refuses a destination with anything in it; creates a missing
      destination; an existing empty directory is accepted; `ModuleURL` derives
      `https://<module path>.git` from build info
- [ ] run them, watch them fail
- [ ] implement — `exec.CommandContext`, the module-path fallback literal in exactly one
      place
- [ ] mark the real-git test `-short`-skippable, matching `git_test.go`
- [ ] `make gates`, commit

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

- [ ] failing tests: a `go.mod` without `.git` is rejected; `LIBRETTO_ROOT` wins over
      everything; nothing found falls back to the bootstrap path
- [ ] run them, watch them fail
- [ ] move `repoRoot` into `root.go`, change the probe from `go.mod` to `.git`, add the
      bootstrap rung
- [ ] `make gates`, commit

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

- [ ] failing tests: a set ldflags value wins over build info; `dev` falls through to
      build info; neither means `dev`
- [ ] run them, watch them fail
- [ ] implement, call it once at startup
- [ ] `make gates`, commit

**Closes:** `TestVersionPrefersLdflagsOverBuildInfo` · `TestVersionFallsBackToBuildInfo`

---

### T5 · the notice row

**Spec:** `spec-panel.md` → all outcomes except the seam
**Files:** modify `internal/ui/panel.go`; modify `internal/ui/panel_test.go`,
`internal/ui/fluid_test.go`
**Blocked by:** nothing
**Produces:** `Panel.UpdateNotice string`

- [ ] failing tests: rendered as its own row between the menu and the strip; absent when
      empty; present in `renderNarrow`; the fluid frame does not tear with it set
- [ ] run them, watch them fail
- [ ] add the field and the row, attention colour, both layouts
- [ ] `make gates`, commit

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

- [ ] failing tests: writes to the given executable path, not `bin/libretto`; a symlinked
      path is resolved and the link survives as a link; an unwritable destination is
      reported and `update` still succeeds
- [ ] run them, watch them fail
- [ ] implement — `filepath.EvalSymlinks`, temp file, atomic rename, the unwritable branch
      reporting where the new binary is
- [ ] `make gates`, commit

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

- [ ] failing tests: picks the highest plain semver from `git ls-remote --tags` output;
      an expired deadline returns no answer and no user-facing error
- [ ] run them, watch them fail
- [ ] implement, add to the interface, extend the fake
- [ ] `make gates`, commit

**Closes:** `TestLatestTagPicksHighestPlainSemver` ·
`TestLatestTagIgnoresPrereleaseAndNonSemverTags` · `TestLatestTagHonoursDeadline`

---

### T8 · the check cache

**Spec:** `spec-repo-sync.md` → outcomes, *Newer release* (last two bullets)
**Files:** modify `internal/repo/release.go`, `internal/repo/release_test.go`
**Blocked by:** **T7**
**Produces:** `CheckedLatest(ctx, root string, ttl time.Duration) (string, error)` —
cache-aware, and the only entry point the CLI calls

- [ ] failing tests: no second call inside the TTL; a failure is cached too, so an offline
      machine does not retry every launch
- [ ] run them, watch them fail
- [ ] implement over `.git/libretto-update-check`
- [ ] `make gates`, commit

**Closes:** `TestCheckCacheSuppressesCallsInsideTTL` ·
`TestCheckCacheRecordsFailureSoOfflineDoesNotRetry`

---

### T9 · bootstrap

**Spec:** `spec-cli.md` → outcomes, *The clone is found, or made*
**Files:** create `cmd/libretto/bootstrap.go`, `cmd/libretto/bootstrap_test.go`; modify
`cmd/libretto/main.go` (call it when `repoRoot` finds nothing)
**Blocked by:** **T2, T3**
**Consumes:** `repo.Clone`, `repo.ModuleURL`, `bootstrapPath`

- [ ] failing tests: the destination is printed **before** the clone runs; a destination
      that exists and is not our clone is refused and nothing is touched; after a
      successful clone the requested command runs
- [ ] run them, watch them fail
- [ ] implement; `LIBRETTO_ROOT` in every test so no real `$HOME` is reachable
- [ ] `make gates`, commit

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

- [ ] failing tests: `Init` returns a command when the check is set and `nil` when it is
      not; the message sets `UpdateNotice`; an action's feedback does not overwrite it
- [ ] run them, watch them fail
- [ ] implement
- [ ] `make gates`, commit

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

- [ ] pass the callback into `NewModel(...).WithReleaseCheck(...)`
- [ ] verify against the existing teatest flow in `cmd/libretto/panelrun_test.go`: the
      panel paints without the check having answered
- [ ] `make gates`, commit

**Closes:** covered by T10's proofs plus the existing panel-run flow. No new criterion —
this task is wiring, and a task with no criterion of its own is wiring or it is scope.

---

### T12 · `doctor` checks live

**Spec:** `spec-cli.md` → *The user is told a newer release exists*
**Files:** modify `cmd/libretto/main.go` (`doctor`); create `cmd/libretto/doctor_test.go`
**Blocked by:** **T7, T8**

- [ ] failing tests: a newer release is reported; a failed check says "could not check"
      and does not set the exit code
- [ ] run them, watch them fail
- [ ] implement with the deadline; the cache is refreshed and nothing in a target is
      touched
- [ ] `make gates`, commit

**Closes:** `TestDoctorReportsNewerRelease` · `TestDoctorSaysSoWhenTheCheckFails`

---

### T13 · say so

**Spec:** `spec-cli.md` → scope, *In*
**Files:** modify `cmd/libretto/main.go` (`usage`), `README.md`
**Blocked by:** **T9, T11, T12**

- [ ] usage: the `go install` line, `LIBRETTO_ROOT`'s default documented as
      `~/.libretto-automata`
- [ ] `README.md`: `go install` as the install and update path, the clone location, and
      that `update` is still how you move forward
- [ ] `scripts/check-payload` — no payload item changed, but run it and read it
- [ ] `make gates`, commit

**Closes:** no test. Documentation that agrees with the binary is checked by reading both,
and this is the one task in the plan whose proof is a human.

---

## Landing

Phase 8's job, listed here so it is not forgotten: apply the three deltas onto
`.agents/specs/repo-sync/spec.md`, `cli/spec.md`, `panel/spec.md`, delete this change
folder, and confirm `spec-drift --anchors` resolves the 31 citations that arrive with them.
