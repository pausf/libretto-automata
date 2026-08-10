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
| `models` | list every agent with the model it runs on, read-only |
| `models set <model> <agent>…` | write that model into each named agent; `--all` for every one |
| `version`, `help` | say so |

### Scope: where it writes

| Invocation | Acts on |
|---|---|
| `--project` / `-p` | `<cwd>/.claude` |
| `--global` / `-g` | `~/.claude` |
| neither | **global** for every subcommand — what every invocation meant before scopes existed. The panel opens where it was left |

Both flags at once is an **error**, not a precedence rule. Two answers to one question
is a mistake worth reporting rather than resolving by picking the last one and hoping.

Repeating the same flag is fine. It is unambiguous, and rejecting it would be pedantry.

**Scope flags are removed from the arguments before the subcommand sees them.** A flag
left in place would reach `prune` and be read as a confirmation nobody gave.

**Every command names the root it acted on, before acting.** A run whose output does
not say where it wrote is a run you have to reconstruct later.

A command acts on exactly one scope. Nothing iterates destinations.

**`models` is the one command whose scope does not decide where it writes.** An
agent's model is one line in the repository's own `agents/*.md`, and both targets
symlink to those files — a write through a symlink writes its destination, so the
change is shared whichever flag was typed. `models set` therefore says so, once, at
the moment it writes: a user who passed `--project` has every reason to expect a
project-local effect, and leaving that expectation unchallenged is how a shared
setting gets discovered by accident weeks later.

What the flag *does* change is the listing. `models` marks every agent the repository
has but the target does not, and **only `linked` counts**: a conflict is somebody
else's file in our slot, and an agent whose slot is occupied does not reach that
target however much the repository wishes it did.

That clause shipped once with no criterion behind it and no code implementing it —
both listings were byte-identical and the suite had no way to notice. It was caught in
review, not by a gate. An outcome with no `Proof:` is a promise nothing keeps.

**Where the project is, is resolved once per invocation and threaded down.** Two
lookups give two answers that agree only by accident, and they did disagree — the
panel's strip read one root while the action wrote to another.

### The panel remembers its destination; nothing else does

**With no flag the panel opens on the destination `tab` last moved it to**, recorded as
one word in `libretto-scope` at the global root. One machine-wide value, not one per
directory.

**Every subcommand with no flag still acts on global.** A command typed into a terminal
does not change meaning because of state left by a session the reader cannot see, and
that includes `preview` — `make preview` prints the same panel whatever a previous
session pressed.

That boundary is also what makes a single machine-wide value safe enough to keep: the
remembered word decides which side the panel *opens* on, and the destination strip is
gold end to end and cursor-marked before any key that acts. The panel can open on
`project` inside a repository that never asked for it; what it cannot do is install
there without that destination being on screen first.

**An explicit flag wins and overwrites nothing.** The flag is about this run; the file
is about where the user works, and a flag typed once must not silently re-answer the
question for every future session.

**`tab` writes it only once the switch has really happened** — after the new view is
assembled, never before. A failed refresh leaves the panel where it was, and a file that
disagrees with the screen is the same class of lie as a strip showing one destination's
counts under another's name.

**Unreadable or unrecognised is `global`**, following `target.Resolve` rather than
inventing a second rule: a value nobody recognises must not produce a destination nobody
chose. **An unwritable one changes nothing else** — remembering is a convenience, and a
convenience that can fail the thing it decorates is worse than no convenience.

**`prune` and `uninstall` never touch it.** It is a regular file at the target's root,
not a link under `skills/`, `agents/` or `commands/`, so it is not ours to remove — and
that is tested rather than assumed, because the failure mode is a destructive path
deleting a file it does not own.

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
- **configuration files. There are none, deliberately** — nothing on disk changes what a
  command does, and the environment-variable table below is the whole of what is
  configurable.

  `libretto-scope` is not one of them, and the difference is the point: it is one word
  the tool writes for itself, recording a position the user already put on screen with
  `tab`, read only to decide which side the *panel* opens on. No command's behaviour can
  be altered by editing it.

  **Its ceiling:** the moment a second thing wants remembering, a one-word file is being
  asked to be a format, and the answer is a real file with a real shape — not a second
  one-word file beside it, and not a delimiter invented in place.

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

**The remembered destination lives inside the global target**, via
`target.Global().Root()` and never a second lookup of the home directory. So
`CLAUDE_HOME` relocates it too, and no new variable is needed — which is also what keeps
a suite that runs `install` and `prune` away from a real `~/.claude`.

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
- **The remembered destination is one machine-wide value, and only the panel reads it.**
  Per-directory was the alternative and was declined with its counter-case stated: a
  single value lets the panel open on `project` inside a repository that never asked for
  it. What bounds that is the same decision's other half — no subcommand consults it —
  and the two are not to be separated. Loosening either one alone reopens a case that was
  closed deliberately.
- **`scopeFlags` returns which flag was seen, not just the scope it resolved to.** "No
  flag" and `--global` produce the same scope and mean different things: one is a default,
  the other an instruction. A second scanner over the arguments would be two places
  deciding what a scope flag is, which is the bug this spec already records happening with
  the project root.
- **`openingScope` and `panelRefresh` are named functions with one caller each.** `run`'s
  panel branch checks `isatty` and so cannot be entered from a test; a decision left
  inline there is a decision no criterion can reach. They exist to be callable, not to
  abstract anything, and a second caller is not expected.

## Task breakdown

- [x] dispatch, `status`, `preview`, `version`, `help`
- [x] `install`, `doctor`, `prune --yes`
- [x] `update` composed over `repo-sync`
- [x] tests for the composition — exit codes, dispatch, flags, prerequisites
- [x] the panel opens on the destination it was left on
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

**Keys are fed to the panel one at a time, with a pause.** Bubbletea batches whatever
is available on a single read into one message, so handing it `y` and `q` together
arrives as `Runes{'y','q'}` — a key called `yq` that matches nothing, and a program that
never exits. The test hung for ten seconds and looked like a deadlock in the code. Real
typing has gaps in it; the harness has to as well.

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

`models`:

- it lists every agent with its current model and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsListsEveryAgentAndChangesNothing
- an agent with no declared model is listed as running the session's
  Proof: cmd/libretto/models_test.go TestModelsShowsDefaultForAnUndeclaredAgent
- `set` applies one model to several named agents
  Proof: cmd/libretto/models_test.go TestModelsSetAppliesToEveryNamedAgent
- `set --all` reaches every agent
  Proof: cmd/libretto/models_test.go TestModelsSetAllReachesEveryAgent
- **`set` with no agents and no `--all` is an error and writes nothing** — a
  destructive default that fires on a forgotten argument is how every agent on the
  machine silently becomes the same model
  Proof: cmd/libretto/models_test.go TestModelsSetWithoutAgentsIsAnError
- an unknown model exits non-zero and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownModel
- an unknown agent name exits non-zero and leaves the valid ones untouched
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownAgentAndWritesNothing
- an agent the repository has but this target does not is marked as such
  Proof: cmd/libretto/models_test.go TestModelsMarksAgentsThatDoNotReachThisScope
- **the two scopes do not produce the same listing** — the flag changes something
  Proof: cmd/libretto/models_test.go TestModelsListingDiffersBetweenScopes
- writing under `--project` says the effect is not project-local
  Proof: cmd/libretto/models_test.go TestModelsSetUnderProjectScopeSaysTheEffectIsShared
- a repository with no agents says so rather than failing
  Proof: cmd/libretto/models_test.go TestModelsWithNoAgentsSaysSo
- piped `models` carries no escape codes
  Proof: cmd/libretto/models_test.go TestModelsOutputHasNoEscapeCodes
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

The remembered destination:

- **with nothing remembered, the panel opens on global**
  Proof: cmd/libretto/remembered_test.go TestNothingRememberedOpensGlobal
- what is written is what is read back, including a value padded by hand
  Proof: cmd/libretto/remembered_test.go TestRememberThenReadIsARoundTrip
- **a remembered `project` is the scope the panel is opened with**
  Proof: cmd/libretto/remembered_test.go TestTheRememberedDestinationOpensThePanel
- **an explicit flag beats it, and does not overwrite it**
  Proof: cmd/libretto/remembered_test.go TestAnExplicitFlagBeatsTheRememberedDestination
- **a successful switch is remembered**
  Proof: cmd/libretto/remembered_test.go TestSwitchingDestinationRemembersIt
- **a switch whose refresh failed remembers nothing** — the file must agree with the
  screen, and the screen did not move
  Proof: cmd/libretto/remembered_test.go TestAFailedSwitchRemembersNothing
- **garbage is global** — empty, whitespace, wrong case, a trailing word
  Proof: cmd/libretto/remembered_test.go TestUnrecognisedPreferenceIsGlobal
- **an unwritable destination does not fail the switch**
  Proof: cmd/libretto/remembered_test.go TestAFailedWriteDoesNotFailTheSwitch
- **subcommands ignore it** — a remembered `project` still leaves a flagless invocation
  resolving to global
  Proof: cmd/libretto/remembered_test.go TestSubcommandsIgnoreTheRememberedDestination
- **`uninstall --yes` does not remove it**, because it was never ours to remove
  Proof: cmd/libretto/remembered_test.go TestUninstallLeavesThePreferenceAlone

**`CLAUDE_HOME` must point at a directory, even in a test that wants a write to fail.**
Aimed at a regular file it breaks `panelData` too — that reads the target's `skills/` —
so the whole refresh errors and the test proves the opposite of its criterion. It did:
`an unwritable preference failed the switch: open …/not-a-directory/skills: not a
directory`. The isolating setup is a **read-only** root, plus an assertion that the write
really was blocked, so a run as root fails loudly rather than passing for the wrong
reason.

**Two of these were mutation-checked.** Defaulting to the remembered value inside
`scopeFlags` turns `TestSubcommandsIgnoreTheRememberedDestination` red; moving the write
above the range check in `panelRefresh` turns `TestAFailedSwitchRemembersNothing` red.

And the payload the CLI installs is checked statically:

  Proof: scripts/check-payload

**Two of these were mutation-checked rather than trusted.** Making `prune` always
confirm, and dropping conflicts from the exit-code sum, each turned the matching test
red. A test that has never been seen to fail is a test nobody has reason to believe.

**Still not covered:** the panel itself under a TTY — which now includes the one line in
`run` handing `openingScope`'s answer to `panelUI`; the decision is proven, the wiring is
read in review. And `update`'s three git paths —
those are `repo-sync`'s to prove, and its `Git` interface still has no fake.
