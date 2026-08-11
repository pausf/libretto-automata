# CLI — delta

Targets: cli

`go install` becomes a supported way to get this tool. That is three changes to the
command surface, and one of them is a bug that only shows up once the binary can live
outside the clone.

## Outcomes

### The clone is found, or made

- **`repoRoot()` requires a git repository, not just a `go.mod`.** Today it accepts any
  directory with a `go.mod` beside the compile-time source path — under `go install` that
  is the read-only module cache. Every operation that matters (`update`, the rebuild, the
  release check) needs git, so the check is `.git`, and the module cache stops
  masquerading as the clone. **This is the change that makes the rest correct.**
- **Resolution order, stated once:**

  | 1 | `LIBRETTO_ROOT` | absolute override, no validation, the escape hatch |
  | 2 | the compile-time source directory | when it has `.git` — development, `make build` |
  | 3 | the current directory | when it has `.git` and a `go.mod` naming this module |
  | 4 | `~/.libretto-automata` | the bootstrap clone |

- **No clone at any of those means bootstrap, and bootstrap announces itself first.** The
  destination is printed before `git clone` runs, not after. A tool that has just written
  a directory into someone's home and then tells them is a tool they cannot decline.
- **Bootstrap refuses a destination it did not create.** `~/.libretto-automata` that
  exists and is not our clone is reported and nothing is touched — the same promise the
  linker makes, applied to the tool's own directory.
- **Bootstrap is not silent success either.** It reports the destination, the tag it
  landed on, and then goes on to do what was asked. One invocation, no second round trip:
  the user typed `go install` for this tool and there is nothing else it can do without
  the payload.

### The version is knowable without ldflags

- **`go install` produces a binary that reports its real version.** `debug.ReadBuildInfo()
  .Main.Version` carries the tag for a `pkg@version` install, so the ldflags value stays
  authoritative when set and build info fills in when it is not. `dev` remains the last
  resort, and it remains honest: a binary that cannot prove its version says so.
- **The order is ldflags → build info → `dev`.** Never the reverse: a `make build` inside
  a dirty tree must keep reporting `-dirty`, which build info cannot know.

### The rebuild replaces the binary that is running

- **`update` rebuilds over `os.Executable()`, not over `bin/libretto`.** Once the binary
  can live in `$GOBIN`, rebuilding into the clone's `bin/` upgrades a file nobody runs —
  `update` would report success while every subsequent invocation stayed on the old
  version. This is the hole the bootstrapper opens, and closing it is not optional.
- **A symlinked executable is resolved to its target.** `make link` puts a symlink in
  `~/.local/bin` pointing at `bin/libretto`; the rebuild writes the resolved file so the
  development setup keeps working and the link is never replaced by a regular file.
- **Atomic rename, never a write over the running file.** Unchanged, and now it is load
  bearing: `$GOBIN/libretto` is the file executing.
- **An unwritable destination does not fail the update.** The pull happened and the links
  are correct; the report says where the new binary is and that the old one is still on
  `PATH`. Rolling back a successful pull because a rename failed loses more than it saves.

### The user is told a newer release exists

- **`doctor` checks live.** The user typed a diagnostic command, so it pays for the
  network call — with the deadline, and offline is a stated "could not check", never an
  error and never silence.
- **`doctor` stays read-only, with no exception to explain.** It asks the remote directly
  rather than through the cache, so it writes nothing at all — the capability spec's "it
  never writes" holds literally. Going through the cache would have meant writing a
  timestamp into `.git/`, and defending that in prose is more expensive than not doing it.
  It also could not distinguish "up to date" from "could not check", because the cache
  swallows the error by design.
- **The panel uses the cache.** See the `panel` delta.
- **The notice points at `update`, which already does the work.** No new subcommand. The
  machinery to move to a newer tag has existed since R3.

## Scope boundaries

**In:** `repoRoot` resolution, bootstrap, the version fallback, the rebuild destination,
`doctor`'s release line, the usage text, and `README` install instructions.

**Out:**

- **`libretto self-update` / re-running `go install` from inside the binary.** The clone is
  the source of truth for the payload; a `$GOBIN` binary that upgraded itself while the
  clone stayed at the old commit is two versions of one tool, and the links point at the
  stale one.
- **npm, npx, Homebrew, a tap, release archives, `goreleaser`.** One distribution channel
  added, working, before a second is discussed. Named here so it is a decision and not an
  omission.
- **auto-update, or bootstrap without saying so.** Nothing moves the user's version
  without them asking.
- **`install.sh`.** Still on the "ask first" list in `AGENTS.md` and not touched here.
- **a config file.** `LIBRETTO_ROOT` is the whole configuration surface, and it already
  exists.
- **`--no-bootstrap`.** `LIBRETTO_ROOT` at an existing clone covers it, and a flag whose
  only user is a test is a flag with no user.

## Constraints

- `CLAUDE_HOME` still governs where links go, and every test that touches a target still
  points it at a temporary directory. Bootstrap adds a second one: **no test clones to a
  real `~/.libretto-automata`.** `LIBRETTO_ROOT`, or the bootstrap path taken from an
  overridable seam.
- The binary keeps carrying no hardcoded version.
- `runtime.Caller(0)` under a module-mode install is either the module cache or a trimmed
  path with nothing on disk. Requiring `.git` is correct for both, and the test asserts the
  behaviour rather than which of the two happens.
- Every path here is reachable without a TTY. Bootstrap prints and proceeds; it does not
  prompt.

## Prior decisions

- **Bootstrapper shape.** Asked and answered this session. See the `repo-sync` delta.
- **`~/.libretto-automata`, a dotdir in `$HOME`.** Not `~/.local/share`: the clone is a
  working git repository the user is expected to `cd` into and edit, not opaque
  application data. `LIBRETTO_ROOT` moves it.
- **The panel and `doctor` differ on purpose.** Live in `doctor`, cached in the panel. A
  panel that waits on the network before painting is a panel that hangs on a bad DNS.

## Task breakdown

1. `repoRoot()`: require `.git`, add the `~/.libretto-automata` rung, keep
   `LIBRETTO_ROOT` first.
2. Bootstrap: announce, clone via `repo.Clone`, refuse a foreign destination, then
   continue into the requested command.
3. Version fallback through `debug.ReadBuildInfo()`.
4. `rebuild`: resolve `os.Executable()`, write there, report an unwritable destination
   without failing `update`.
5. `doctor`: a live release check with the deadline, and a "could not check" line when it
   fails.
6. Usage text: `go install` line, `LIBRETTO_ROOT` default documented.
7. `README`: `go install` as the install and update path.

## Verification criteria

```
Proof: cmd/libretto/root_test.go TestRepoRootRequiresGitDirectory
Proof: cmd/libretto/root_test.go TestRepoRootPrefersEnvOverride
Proof: cmd/libretto/root_test.go TestRepoRootFallsBackToBootstrapPath
Proof: cmd/libretto/bootstrap_test.go TestBootstrapAnnouncesDestinationBeforeCloning
Proof: cmd/libretto/bootstrap_test.go TestBootstrapRefusesForeignDestination
Proof: cmd/libretto/bootstrap_test.go TestBootstrapContinuesIntoRequestedCommand
Proof: cmd/libretto/version_test.go TestVersionPrefersLdflagsOverBuildInfo
Proof: cmd/libretto/version_test.go TestVersionFallsBackToBuildInfo
Proof: cmd/libretto/update_test.go TestRebuildReplacesRunningExecutable
Proof: cmd/libretto/update_test.go TestRebuildResolvesSymlinkedExecutable
Proof: cmd/libretto/update_test.go TestRebuildReportsUnwritableDestinationWithoutFailing
Proof: cmd/libretto/doctor_test.go TestDoctorReportsNewerRelease
Proof: cmd/libretto/doctor_test.go TestDoctorSaysSoWhenTheCheckFails
```
