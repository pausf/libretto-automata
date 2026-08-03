# CLI

Governs: cmd/libretto/** install.sh

The command surface. Every action is reachable without a terminal, so scripts and CI
can use it, and the same binary shows a panel when a human is driving.

## Outcomes

`libretto` with a TTY and no subcommand shows the panel. Otherwise it behaves as a
plain command.

| Command | Does |
|---|---|
| `install` | link every item into every target |
| `update` | pull, relink, rebuild when Go changed, report |
| `status` | every item's state, read-only |
| `doctor` | what needs attention, plus what the payload expects on this machine |
| `prune` | show links whose source is gone; `--yes` removes them |
| `uninstall` | show what this repo installed here; `--yes` removes it |
| `preview` | print the panel once, no TUI |
| `version`, `help` | say so |

### Scope: where it writes

| Invocation | Acts on |
|---|---|
| `--project` / `-p` | `<cwd>/.claude` |
| `--global` / `-g` | `~/.claude` |
| neither | **global** — what every invocation meant before scopes existed |

Both flags at once is an **error**, not a precedence rule. Two answers to one question
is a mistake worth reporting rather than resolving by picking the last one and hoping.

Repeating the same flag is fine. It is unambiguous, and rejecting it would be pedantry.

**Scope flags are removed from the arguments before the subcommand sees them.** A flag
left in place would reach `prune` and be read as a confirmation nobody gave.

**Every command names the root it acted on, before acting.** A run whose output does
not say where it wrote is a run you have to reconstruct later.

A command acts on exactly one scope. Nothing iterates destinations.

**Where the project is, is resolved once per invocation and threaded down.** Two
lookups give two answers that agree only by accident, and they did disagree — the
panel's strip read one root while the action wrote to another.

### Running an action for the panel

The panel needs an action's report as **lines**, not as output on the terminal it is
drawing over. So stdout is redirected for the duration of the action and handed back
split into lines, blanks dropped.

Three properties make that safe rather than clever:

- **stdout is put back afterwards**, always. Left redirected, every later print goes
  to a dead pipe.
- **the pipe is drained concurrently.** A report longer than the pipe's buffer would
  otherwise block the action half-way, with the links half-written.
- **a failing action still returns its lines.** What half-happened is the part that
  explains the failure.

The report is the command's own words. A second rendering of the same facts is a
second thing that can disagree with the first.

**Confirmation is passed in, not inferred.** `prune` is dry unless told otherwise, so
the panel can show a plan on the first press and carry it out on the second without
the dispatcher guessing which press it is on.

`prune`'s dry report shortens the link's resolved destination. It is an absolute path
repeating what the item name already says, and left whole it pushed the interesting
part of the line out of view.

**Exit codes carry meaning.** Non-zero whenever something was not linked or needs
attention, so a script can gate on it. A conflict counts: an item that did not get
linked is an incomplete install, whatever the reason.

**Without a TTY and without a subcommand it prints usage and exits non-zero.** It never
opens a TUI it cannot drive.

**Help names the command it was invoked as.** The binary can be linked under any name,
so help that hardcodes one lies to everyone using another.

## Scope boundaries

**In:** argument parsing, dispatch, output formatting, exit codes, environment
variables, repo-root discovery, the prerequisite report.

**Out:**

- classification and mutation — `link-state` and `linking` own those
- git — `repo-sync` owns it
- rendering — `panel` owns it
- **requiring any optional companion.** The prerequisite report says what is present.
  It never blocks.
- configuration files. There are none, deliberately.

## Constraints

**`uninstall` is dry by default too**, and exits non-zero when anything it meant to
remove survived — a refusal or a failure means the destination is not in the state the
report implies. A kept conflict is not a failure: it was never ours, and saying so is
the correct outcome.

**`prune` is dry by default.** With no flag it prints what it would remove and changes
nothing; `--yes` carries it out. This is stricter than confirming only when
interactive — a pipe is no reason to be less careful, and a destructive command that
acts before being asked twice eventually deletes the wrong thing.

**The prerequisite report never affects the exit code.** None of those tools is
required, and failing on an optional absence trains people to ignore the section.

**A companion counts as present wherever it legitimately lives** — the plugin tree,
`skills/`, or `commands/`. Checking only one place reports an installed tool as missing,
which is worse than not checking, because it sends people to install what they have.

**Five environment variables, each with a working default:**

| Variable | Effect |
|---|---|
| `CLAUDE_HOME` | Claude Code's root instead of `~/.claude`. What makes the test suite safe. |
| `LIBRETTO_ROOT` | the repo location, instead of deriving it from the binary |
| `LIBRETTO_ASCII` | `safe` swaps quadrant glyphs for half blocks |
| `LIBRETTO_THEME` | `dark` or `light`, instead of detecting |
| `COLUMNS` | layout width when stdout is not a terminal |

A root too long for the panel is elided from the left with `…`, keeping whole trailing
segments — the tail says *which* directory this is, and the marker says plainly that
something was removed. Left unbounded it pushes the frame's right border out of
alignment, which `panel` forbids at every width; a temporary directory is long enough
to do it, and it took looking at a render to notice.

A value that never varies is not configuration, so nothing else is exposed.

**The version is stamped from git at build time, never held in a source constant.**
`make build` passes `-ldflags "-X main.version=$(git describe --tags --always --dirty)"`,
so the binary reports `v0.2.0`, or `v0.2.0-3-gabc123` past a tag, or `v0.2.0-dirty` over
uncommitted changes.

A constant drifts from the tag the moment someone forgets to bump it, and it drifts
**silently** — the binary keeps announcing a version nobody released. Asking git means
the answer cannot be wrong.

A build without those flags reports `dev`, not a version number. **A binary that cannot
prove its version says so rather than claiming one**, which is the same rule as
everywhere else here: nothing is asserted that was not observed.

**`install.sh` still exists** as the bootstrap for a machine without Go. It is a
prototype and it lacks the ownership re-check, so it must not be presented as
equivalent. It goes when `libretto install` has been verified against a real
`~/.claude`.

## Prior decisions

- The command is `libretto`; `libretto-automata` is linked too. `lib` was rejected: it
  reads as a system directory.
- `make link` symlinks rather than copies, so `make build` updates the installed
  command and no stale binary can pretend to be current. It refuses to overwrite a
  link that is not ours — the tool does not get an exception from its own rule.
- The repo root comes from the binary's compile-time location, with `LIBRETTO_ROOT` as
  the override. A tool that guesses its own repo from the working directory installs
  the wrong thing from the wrong place.
- After a rebuild the process states that it is still the old binary, rather than
  implying the upgrade took effect mid-run.

## Task breakdown

- [x] dispatch, `status`, `preview`, `version`, `help`
- [x] `install`, `doctor`, `prune --yes`
- [x] `update` composed over `repo-sync`
- [x] tests for the composition — exit codes, dispatch, flags, prerequisites
- [ ] 5.2 `--json` for `status` and `doctor`
- [ ] 5.3 the panel path under a real TTY
- [ ] `doctor`: target directory missing or unwritable
- [ ] `doctor`: repo state — uncommitted changes, behind the remote
- [ ] `install`: present the plan before applying, as `prune` already does
- [ ] 7.3 delete `install.sh`, once `libretto install` is verified against a real
      `~/.claude` with a throwaway item

## Verification criteria

These cover the **composition**, not the logic composed. `linking`, `link-state`,
`targets` and `panel` prove the pieces behave; these prove the CLI wires them to the
right exit codes and refuses what it promised to refuse.

`CLAUDE_HOME` is what makes them safe: every one runs against `t.TempDir()`, so
commands that write and delete never see a real `~/.claude`.

- **a conflict makes `install` exit non-zero, and the foreign file is unchanged**
  Proof: cmd/libretto/main_test.go TestInstallExitsNonZeroOnConflict
- everything linking exits zero, and the link is really there
  Proof: cmd/libretto/main_test.go TestInstallExitsZeroWhenEverythingLinks
- a second install reports an empty plan rather than doing the work again
  Proof: cmd/libretto/main_test.go TestInstallIsIdempotent
- anything needing attention makes `doctor` exit non-zero
  Proof: cmd/libretto/main_test.go TestDoctorExitsNonZeroWhenSomethingNeedsAttention
- **`prune` without `--yes` changes nothing, and says what it would do**
  Proof: cmd/libretto/main_test.go TestPruneWithoutYesChangesNothing
- **`prune --yes` removes only what the plan named** — a live item's link and a foreign
  entry both survive
  Proof: cmd/libretto/main_test.go TestPruneYesRemovesOnlyWhatThePlanNamed
- `version` and `help` succeed; an unknown command fails
  Proof: cmd/libretto/main_test.go TestRunDispatch
- **no TTY and no subcommand prints usage and exits 2** — checked by re-executing the
  test binary, because `os.Exit` cannot be observed in-process
  Proof: cmd/libretto/main_test.go TestNoTTYAndNoSubcommandExitsNonZero
- piped `status` carries no escape codes
  Proof: cmd/libretto/main_test.go TestStatusOutputHasNoEscapeCodes
- **help and remedies name the command they were invoked as**, not a fixed string
  Proof: cmd/libretto/main_test.go TestOutputNamesTheInvokedCommand
- **the prerequisite report never changes the exit code** — with an empty `PATH` and an
  empty home, `doctor` still exits zero on a correct tree
  Proof: cmd/libretto/main_test.go TestPrerequisitesDoNotAffectTheExitCode
- a companion counts as present in all four places it can legitimately live
  Proof: cmd/libretto/main_test.go TestCompanionFoundWhereverItLegitimatelyLives
- an absent companion is reported absent
  Proof: cmd/libretto/main_test.go TestCompanionAbsentIsAbsent
- **a project install leaves the global config untouched**
  Proof: cmd/libretto/scope_test.go TestInstallProjectScopeLeavesGlobalAlone
- **a global install leaves the project untouched**
  Proof: cmd/libretto/scope_test.go TestInstallGlobalScopeLeavesProjectAlone
- **`prune --project --yes` removes from the project only**
  Proof: cmd/libretto/scope_test.go TestPruneProjectScopeLeavesGlobalAlone
- no flag means global
  Proof: cmd/libretto/scope_test.go TestDefaultScopeIsGlobal
- scope flags never reach the subcommand
  Proof: cmd/libretto/scope_test.go TestScopeFlagsAreRemovedFromTheArguments
- both flags at once is an error
  Proof: cmd/libretto/scope_test.go TestBothScopeFlagsIsAnError
- repeating one flag is accepted
  Proof: cmd/libretto/scope_test.go TestRepeatingTheSameScopeFlagIsFine
- the output names the root it acted on
  Proof: cmd/libretto/scope_test.go TestOutputNamesTheScopeRoot
- **a long root is elided rather than tearing the frame**
  Proof: cmd/libretto/scope_test.go TestShortenKeepsRootsInsideTheBudget
- **the strip and the action agree on where the project is**
  Proof: cmd/libretto/scope_test.go TestStripAndRunnerAgreeOnTheProjectRoot
- **the capture returns lines and puts stdout back**
  Proof: cmd/libretto/panelrun_test.go TestRunCapturedReturnsLinesAndRestoresStdout
- **a failing action still returns what it printed**
  Proof: cmd/libretto/panelrun_test.go TestRunCapturedKeepsOutputOnFailure
- every enabled menu label has a dispatch case
  Proof: cmd/libretto/scope_test.go TestEveryMenuLabelDispatches
- **an unconfirmed prune removes nothing**
  Proof: cmd/libretto/scope_test.go TestDispatchedPruneIsDry
- **an unconfirmed uninstall removes nothing**
  Proof: cmd/libretto/uninstall_test.go TestUninstallWithoutYesChangesNothing
- `uninstall --yes` removes what the plan named
  Proof: cmd/libretto/uninstall_test.go TestUninstallYesRemovesOurLinks
- **it acts on one destination only**
  Proof: cmd/libretto/uninstall_test.go TestUninstallProjectScopeLeavesGlobalAlone
- **a conflict is kept, reported, and not an error on its own**
  Proof: cmd/libretto/uninstall_test.go TestUninstallReportsConflicts
- nothing of ours installed is a state, not an error
  Proof: cmd/libretto/uninstall_test.go TestUninstallOnACleanDestinationSaysSo
- **install then uninstall is a round trip**
  Proof: cmd/libretto/uninstall_test.go TestInstallThenUninstallIsARoundTrip
- a root inside the budget is left alone
  Proof: cmd/libretto/scope_test.go TestShortenLeavesShortPathsAlone
- the tally renders in a fixed order and omits zeros, so one situation gives one line
  Proof: cmd/libretto/main_test.go TestSummarise

And the payload the CLI installs is checked statically:

  Proof: scripts/check-payload

**Two of these were mutation-checked rather than trusted.** Making `prune` always
confirm, and dropping conflicts from the exit-code sum, each turned the matching test
red. A test that has never been seen to fail is a test nobody has reason to believe.

**Still not covered:** the panel itself under a TTY, and `update`'s three git paths —
those are `repo-sync`'s to prove, and its `Git` interface still has no fake.
