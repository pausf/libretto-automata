# CLI — delta

Targets: cli

Two words, and the second one changes meaning. `install` links the payload; `update` brings the
installation up to date. There is no third command.

## Outcomes

### One `go install`, and there is nothing else to do

```bash
go install github.com/pausf/libretto-automata/cmd/libretto@latest
libretto install
```

The payload ships inside the module, so the first command downloads it along with the binary —
see `distribution`. No tarball, no bootstrap, no second fetch.

### Two roots, not one

This tool used one `root` for two unrelated things, and got away with it for as long as a clone
happened to be both:

| The **payload root** | where `skills/`, `agents/`, `commands/` are — what links point at |
| The **checkout** | a git repository, needed only by `update`'s git route |

`go install` breaks the coincidence: there is a payload and no checkout.

**The payload root resolves in four rungs:**

| 1 | `LIBRETTO_ROOT` | absolute override, taken as given |
| 2 | the compile-time source directory | when it has `.git` — development |
| 3 | the working directory | when it has `.git` **and** a `go.mod` naming this module |
| 4 | `$GOMODCACHE/<module>@<version>` | the payload `go install` brought down |

- **The probe is `.git`, not `go.mod`.** It used to be `go.mod`, and under `go install` that
  matches the module cache — which would win for the wrong reason. Rungs 2 and 3 answer both
  questions at once; rung 4 answers only the first, which is why `update` probes for a checkout
  before pulling. A `.git` *file* counts too: that is what a worktree has.
- **Rung 3 checks the module, not just the presence of a repository.** Any git repository
  satisfied the old `$PWD` fallback, so `libretto install` inside an unrelated project went
  looking for a payload that project does not have.
- **Rung 4 is this binary's own version**, from build info. So the payload linked is the one that
  shipped with the command doing the linking, and the two can never be a version apart.
- **Two homes were tried and abandoned before this one:** `~/.libretto-automata` with a `git
  clone` into it, then `~/.local/share/libretto/<version>` with a tarball extracted into it.
  Neither shipped in a tag, so nothing migrates and nobody has one.

### `update` brings the installation up to date

| in a checkout | pull, rebuild when Go moved, relink |
| anywhere else | `go install` the newest version, relink |

- **One command, because for the person typing it there is one meaning:** put me on the newest
  version. Which mechanism runs is a consequence of how the tool got onto the machine, which
  they know.
- **An earlier draft split this into `update` and `upgrade`** and defended it: *a command whose
  mechanism depends on invisible state has unpredictable failures*. The state is not invisible —
  it is "did you clone this or install it" — and the split cost two commands, two menu rows, a
  pair of mutual refusals, a help table that changed per machine, and a function whose only job
  was choosing which of the two names to print in a notice.
- **The release route says `git` nowhere**, in success or in any failure. That is the complaint
  that started this change: the word appearing in front of somebody who only wanted to use the
  tool.
- **It relinks the version it just installed**, not `payloadRoot()`. This process still reports
  the old version, so resolving again would link the payload it is already on — and the update
  would report success having changed no link.
- **Relinking is not redundant.** The new version is a new directory, and a release that *adds*
  an item leaves that item unlinked with nothing to say so. That was a queued complaint of its
  own before this command existed.
- **Already newest is a success that changes nothing and says so.** Nothing is installed and
  nothing relinked.
- **Nothing installed yet upgrades rather than refusing.** `update` is how the payload arrives if
  `go install` was interrupted, or the module cache was cleaned.

### The rebuild replaces the binary that is running

`update`'s git route rebuilds over `os.Executable()`, not over the clone's `bin/libretto`. Once
the command can live in `$GOBIN`, rebuilding into the clone upgrades a file nobody executes.

- **A symlinked executable is resolved and written through.** `make link` puts one in
  `~/.local/bin`; replacing it with a regular file would sever the development setup silently.
- **The write is probed before the compile.** A refused destination otherwise surfaces as a `go
  build` error string, indistinguishable from broken code, after paying for the compile.
- **An unwritable destination does not fail the update**, and the report names both paths: what
  could not be written, and where the new binary is.
- **The advice matches what happened.** "Run it again to use the new one" is printed only when
  the running binary was actually replaced.

### A missing payload is a stop, not an empty report

A fresh install with an interrupted download, or a `go clean -modcache`, leaves a binary and no
payload. Every command that reads the tree says so and points at `update`.

**`prune` is why this is a stop.** With no payload every link in the target resolves to nothing,
so a scan reports all of them `stale` — and `prune --yes` would remove every item the user has.
A destructive command doing exactly what it promises, on a premise that is false.

- **`models` is exempt**: it reads the *target's* agents directory, not the payload, so it works
  on a machine with nothing installed — which is when somebody might most want it.
- **`update` is exempt**, because it is what fixes the state.
- **`version` and `help` are answered before the payload is even located**, and write nothing
  anywhere. Neither reads a skill.

### Saying a newer version exists

- **`doctor` checks live**, with a five-second ceiling, and names the mode it is in — "up to
  date" means something different in each: a checkout three commits past a tag is current with
  its remote and behind the release.
- **`doctor` always says something.** Newer available, up to date, none published, could not
  check, or "running X, the latest is Y" for a binary ahead or unparseable. Printing nothing
  would read as up to date, which is a claim nobody verified.
- **It never sets the exit code.**
- **The panel uses the cache**, and each mode keeps its answer somewhere it has: a checkout in
  `.git/`, an installed copy in the module cache root. One TTL and one failure policy, `repo`'s.
- **The notice names `update`**, which is the only command there is.

## Scope boundaries

**In:** the payload-root rungs, `update`'s two routes, the missing-payload stop, the rebuild
destination, the release line, usage and `README`.

**Out:**

- **downloading or verifying anything.** `distribution` runs `go install`; the go command does
  the rest.
- **`install.sh`, brew, a tap, npm, per-platform binaries.** One entry point. Named so each is a
  decision.
- **a `sync` command.** `gentle-ai` needs one because its assets live inside the binary and have
  to be written out. Here the payload is already files and relinking is what `install` does.
- **auto-update, and prompting to update.** The notice is a row; pressing it is the user's move.
- **migrating either abandoned payload home.** Neither shipped in a tag.
- **`update --to <version>`.** No reported need. *Ceiling:* the first bad release to back out.

## Constraints

- `CLAUDE_HOME` governs targets; `LIBRETTO_ROOT` governs the payload root. No third variable.
- **No test writes to a real `~/.claude`, and none runs a real `go install`.** The release
  route's four effects are parameters.
- The mode check is `isRepo(root)` — the same one-line probe the rungs use, not a second notion
  of what development means.
- The module path is written once, in `moduleLine`, and everything derives from it.

## Prior decisions

- **One `update`, not `update` and `upgrade`.** Asked and answered by the user, and the reasoning
  above is recorded because the split had an argument behind it that sounded better than it was.
- **`go install` is the entry point.** Which means **Go is required**, and that is a real limit
  on who can use this, not an oversight. Named here and in `distribution`.
- **The clone bootstrap and the release tarball were both wrong**, in that order, in one session.
  See `distribution` for why each lost.
- **A checkout you are standing in wins over the module cache.** That is what keeps editing a
  skill and seeing it live working.
- **The compile-time location is rung 2, not the definition of the root.** The old bullet said
  the root *comes from* there with `LIBRETTO_ROOT` as the only override; that stopped being true
  when the binary could be installed from outside a checkout. Its reasoning survives: a tool that
  guesses its repo from the working directory installs the wrong thing from the wrong place,
  which is why rung 3 requires the `go.mod` to name this module.

## Verification criteria

```
Proof: cmd/libretto/root_test.go TestPayloadRootRequiresGitDirectory
Proof: cmd/libretto/root_test.go TestPayloadRootAcceptsAWorktreeGitFile
Proof: cmd/libretto/root_test.go TestPayloadRootPrefersEnvOverride
Proof: cmd/libretto/root_test.go TestPayloadRootReadsTheOverrideFromTheEnvironment
Proof: cmd/libretto/root_test.go TestPayloadRootPrefersTheCompileTimePathOverTheWorkingDirectory
Proof: cmd/libretto/root_test.go TestPayloadRootAcceptsTheWorkingDirectoryOnlyForThisModule
Proof: cmd/libretto/root_test.go TestPayloadRootFallsBackToTheActivatedRelease
Proof: cmd/libretto/root_test.go TestPayloadRootStillPrefersACheckoutYouAreStandingIn
Proof: cmd/libretto/root_test.go TestVersionAndHelpTouchNothing
Proof: cmd/libretto/update_release_test.go TestUpdateInstallsThenRelinksTheVersionItInstalled
Proof: cmd/libretto/update_release_test.go TestUpdateRelinksTheNewVersionNotTheRunningOne
Proof: cmd/libretto/update_release_test.go TestUpdateRelinksSoNewItemsAppear
Proof: cmd/libretto/update_release_test.go TestUpdateFromAReleaseNeverMentionsGit
Proof: cmd/libretto/update_release_test.go TestUpdateReportsWhichStepFailed
Proof: cmd/libretto/update_release_test.go TestAFailedUpdateLeavesThePreviousVersionActive
Proof: cmd/libretto/update_release_test.go TestUpdateOnTheNewestVersionDoesNothing
Proof: cmd/libretto/update_release_test.go TestUpdateFromNothingInstalled
Proof: cmd/libretto/update_release_test.go TestUpdateTakesTheRouteThisInstallationCameBy
Proof: cmd/libretto/update_release_test.go TestTheUpdateRowNamesTheMechanism
Proof: cmd/libretto/update_release_test.go TestInstallWithNoPayloadPointsAtUpdate
Proof: cmd/libretto/update_release_test.go TestPruneWithNoPayloadRefusesInsteadOfDeletingEverything
Proof: cmd/libretto/update_release_test.go TestModelsWorksWithNoPayload
Proof: cmd/libretto/update_test.go TestRebuildReplacesRunningExecutable
Proof: cmd/libretto/update_test.go TestRebuildResolvesSymlinkedExecutable
Proof: cmd/libretto/update_test.go TestRebuildReportsUnwritableDestinationWithoutFailing
Proof: cmd/libretto/version_test.go TestVersionPrefersLdflagsOverBuildInfo
Proof: cmd/libretto/version_test.go TestVersionFallsBackToBuildInfo
Proof: cmd/libretto/version_test.go TestVersionIsDevWhenNothingKnows
Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesBothVersionsAndTheAction
Proof: cmd/libretto/version_test.go TestReleaseNoticeIsSilentWhenThereIsNothingToSay
Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesTheCommandForTheMode
Proof: cmd/libretto/doctor_test.go TestDoctorReportsNewerRelease
Proof: cmd/libretto/doctor_test.go TestDoctorSaysUpToDate
Proof: cmd/libretto/doctor_test.go TestDoctorSaysSoWhenTheCheckFails
Proof: cmd/libretto/doctor_test.go TestDoctorHandlesARemoteWithNoReleases
Proof: cmd/libretto/doctor_test.go TestDoctorDoesNotRankAnUnidentifiableBinary
Proof: cmd/libretto/doctor_test.go TestDoctorNamesTheModeItIsIn
```

`TestUpdateFromAReleaseNeverMentionsGit` is not a joke criterion. The complaint that produced
this change was `git pull` appearing in front of somebody who only wanted to use the tool, and a
promise about output is kept by asserting on output — across success and every failure.
