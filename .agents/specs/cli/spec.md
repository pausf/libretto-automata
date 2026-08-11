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
| `doctor` | what needs attention, what the payload expects on this machine, and the release state |
| `prune` | show links whose source is gone; `--yes` removes them |
| `uninstall` | show what this repo installed here; `--yes` removes it |
| `preview` | print the panel once, no TUI |
| `models` | list every agent with the model it runs on, read-only |
| `models set <model> <agent>…` | write that model into each named agent; `--all` for every one |
| `version`, `help` | say so |

### Finding the clone, or making one

`go install github.com/pausf/libretto-automata/cmd/libretto@latest` is a supported way in.
The binary arrives with no payload, and every link it makes points into a clone, so it can
get one.

**The clone is resolved in four rungs:**

| 1 | `LIBRETTO_ROOT` | absolute override, taken as given, unvalidated |
| 2 | the compile-time source directory | when it has `.git` — development, `make build` |
| 3 | the working directory | when it has `.git` **and** a `go.mod` naming this module |
| 4 | `~/.libretto-automata` | the bootstrap clone, whether or not it exists yet |

**The probe is `.git`, not `go.mod`.** It used to be `go.mod`, and under `go install` that
matches the read-only module cache — which would win, and every link would point into a
versioned cache directory that the next install orphans. Everything the root is used for
needs git: the pull, the rebuild decision, the release check. A `.git` **file** counts too,
because that is what a worktree has and this flow recommends worktrees.

Rung 3 checks the module, not just the presence of a repository. Any git repository
satisfied the old `$PWD` fallback, so `libretto install` inside an unrelated project went
looking for a payload that project does not have.

**Nothing found means bootstrap, and bootstrap announces itself first.** The destination is
printed before `git clone` runs. A tool that has already written a directory into somebody's
home and then mentions it is a tool they had no chance to decline.

Three outcomes and no fourth:

| Already our clone | nothing happens, nothing is printed — narrating a no-op every launch is noise |
| Nothing there | announce the destination and the URL, clone, report |
| Something else | refused, named, nothing touched — the linker's promise, applied to the tool's own directory |

- **A failed clone is cleaned up.** Half a clone in `~/.libretto-automata` would be a
  foreign destination forever, and the user would have to work out for themselves that
  deleting it is the fix.
- **It does not report a tag.** `git clone` checks out the default branch's tip, which is
  normally *past* the last tag, so any tag named there would describe a release the clone is
  not at. `libretto version` answers that and cannot be wrong about it.
- **`version` and `help` are answered before the clone is even looked for.** Cloning a
  repository into somebody's home because they asked what version they were running would be
  indefensible.
- **Nothing prompts.** Every path works without a TTY, or `install` in CI breaks the day it
  needs to bootstrap.

### The rebuild replaces the binary that is running

`update` rebuilds over `os.Executable()`, not over the clone's `bin/libretto`. Once the
command can live in `$GOBIN`, rebuilding into the clone upgrades a file nobody executes —
`update` would pull, relink, print "rebuilt" and change nothing the user runs.

- **A symlinked executable is resolved and written through.** `make link` puts one in
  `~/.local/bin`; replacing it with a regular file would sever the development setup
  silently.
- **The write is probed before the compile.** A refused destination otherwise surfaces as a
  `go build` error string — indistinguishable from broken code, and only after paying for the
  compile.
- **An unwritable destination does not fail the update.** The pull happened and the links are
  correct; rolling those back over a refused rename loses more than it saves. The report names
  **both** paths: what could not be written, and where the new binary actually is — which
  means building it somewhere, so the fallback is the clone's `bin/`, the one place a write
  cannot be refused because the user owns the checkout.
- **And the advice matches what happened.** "Run it again to use the new one" is printed only
  when the running binary was actually replaced. On the refused path the line says the
  opposite, because advice that does nothing reads as though the upgrade took.

### Saying a newer release exists

- **`doctor` checks live**, with a five-second ceiling: the user typed a diagnostic and can
  afford it. It stays read-only — it asks the remote directly rather than through the cache,
  so it writes nothing at all. Going through the cache would also swallow the very error this
  line exists to report.
- **`doctor` always says something.** Newer available, up to date, the remote has no releases,
  could not check, or "running X, the latest is Y" for a binary ahead of the remote or one
  whose version will not parse. Printing nothing on a failed check would read as "you are up
  to date", which is a claim nobody verified. That last case gets the facts and **no
  ranking**.
- **It never sets the exit code.** Being a release behind is news, and an unreachable remote
  is not this tool's fault.
- **The panel uses the cache, not a live call** — see `panel`. A panel that waits on the
  network before its first paint hangs on bad DNS.
- **The action is `update`, which already exists.** No new subcommand, no accelerator key: the
  machinery to move to a newer tag has been there since the update flow landed.

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

**`models` acts on the agents directory of the selected target** — `~/.claude/agents`
under `--global`, `<cwd>/.claude/agents` under `--project`. Every `*.md` there is
listed, whether libretto created it or not.

An agent file in a target is one of two things, and they behave differently:

| The file is | Writing it affects | Marked |
|---|---|---|
| a symlink into this repository | **every project on the machine** — one file, many targets | `shared` |
| a real file in that target | that target only | *(nothing)* |

The listing marks the shared ones, and `set` says which kind it just wrote.

**`set` edits real files this tool did not create, and does so without asking.**
Chosen over a per-write confirmation: the user's own agents are the point of the
feature, and a prompt on every one of twenty-odd rows is a prompt people learn to
dismiss. This is not `--force` by the back door — `AGENTS.md` forbids *overwriting*
what the tool did not create, meaning replacing a foreign item with our symlink and
destroying what was there. This replaces one frontmatter line, in a file the user
named, at their request. `install` still refuses a conflict exactly as before.

`models` names any link in the target that points at nothing and sends the reader to
`doctor` — skipping in silence would trade a loud failure for a quiet one, and a
rename leaves exactly these behind. A target with no `agents/` directory says so and
exits zero. A `*.md` there with no
frontmatter is refused by name and **nothing in the set is written** — the
all-or-nothing guarantee, written for this repository's own tidy files, now protects
foreign ones.

**An earlier version of this command read `<repo>/agents` regardless of scope**, and
told every user their write was machine-wide. That was true while the subject was
always the repository's own file. It listed seven agents installed nowhere on a
machine holding twenty-two, and the blanket warning is now a per-row fact: saying it
about a target-local write is the same class of error as the silence it replaced.

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
variables, clone discovery and bootstrap, the rebuild destination, the release line, the
prerequisite report.

**Out:**

- classification and mutation — `link-state` and `linking` own those
- git — `repo-sync` owns it, cloning and the release check included
- rendering — `panel` owns it
- **`libretto self-update`, or re-running `go install` from inside the binary.** The clone
  is the source of truth for the payload; a `$GOBIN` binary that upgraded itself while the
  clone stayed at the old commit is two versions of one tool, with the links pointing at the
  stale one.
- **npm, npx, Homebrew, a tap, release archives, `goreleaser`.** One distribution channel
  added and working before a second is discussed. Named so it is a decision and not an
  omission.
- **auto-update, and bootstrap without saying so.** Nothing moves the user's version without
  them asking.
- **`--no-bootstrap`.** `LIBRETTO_ROOT` at an existing clone covers it, and a flag whose only
  user is a test has no user.
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
| `LIBRETTO_ROOT` | the payload clone, instead of resolving it; default `~/.libretto-automata`. **The whole configuration surface for where the clone lives**, which is why no `--no-bootstrap` flag exists |
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

**But ldflags is no longer the only source.** `go install pkg@version` gets no ldflags and
does record the module version, so the order is **ldflags → build info → `dev`** and never
the reverse: build info knows neither the dirt nor the commits past the tag, so preferring
it would turn `v0.2.0-3-gabc123-dirty` back into a clean `v0.2.0`. `(devel)` falls through
to `dev` — it is what the toolchain records for a build from a working tree, not a version.

A build with neither reports `dev`, not a version number. **A binary that cannot prove its
version says so rather than claiming one**, which is the same rule as everywhere else here:
nothing is asserted that was not observed.

**Two timeouts, deliberately different.** The bootstrap clone gets five minutes because the
user is waiting for the payload they asked for; the release check gets five seconds because
nobody is waiting to be told they are up to date.

**Nothing in the test suite clones to a real `~/.libretto-automata`.** `LIBRETTO_ROOT` is
the second half of what `CLAUDE_HOME` does for targets, and `resolveRoot` takes its four
inputs rather than reading them — `runtime.Caller` is fixed at compile time, so inside this
repository rung 2 always wins and rungs 3 and 4 are otherwise unreachable. Those are exactly
the rungs that exist only for a binary living somewhere else.

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
- **The clone is resolved in four rungs, and the compile-time location is only the second.**
  This bullet used to say the root *comes from* the compile-time location with
  `LIBRETTO_ROOT` as the only override, which stopped being true when the binary could be
  installed from outside a checkout. Its reasoning survives intact: a tool that guesses its
  repo from the working directory installs the wrong thing from the wrong place, which is why
  rung 3 requires the `go.mod` to name this module.
- **`go install` delivers a bootstrapper, not an embedded payload.** Asked and answered. See
  `repo-sync` for the reasoning; the consequence here is that a clone always exists behind
  the links.
- **`~/.libretto-automata`, a dotdir in `$HOME`** — not `~/.local/share`. It is a working git
  repository the user is expected to `cd` into and edit, not opaque application data.
- After a rebuild the process states that it is still the old binary, rather than
  implying the upgrade took effect mid-run — **and only when the binary was really
  replaced.** On the refused path it says the opposite.
- **A clone you are standing in wins over `~/.libretto-automata`.** That is what keeps
  editing a skill and seeing it live working, and it is the reason the payload is not
  compiled in.
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
- [x] `repoRoot` by `.git`, in four rungs, behind a testable seam
- [x] bootstrap: announce, clone, refuse a foreign destination, clean up a failure
- [x] the version fallback through build info
- [x] the rebuild over the running executable
- [x] `doctor`'s release line, and the panel's cached notice
- [ ] 5.2 `--json` for `status` and `doctor`
- [ ] 5.3 the panel path under a real TTY
- [ ] `doctor`: target directory missing or unwritable
- [ ] `doctor`: repo state — uncommitted changes, behind the remote *(the release half of
      this is done; the working-tree half is not)*
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
- **the report includes the gate tools** — `rg` and `jq`, each attributed to the skill
  that runs on it, with its install line. The payload's own gates ask their questions
  through them, and a doctor that omits them leaves a silent false negative looking
  healthy.
  Proof: cmd/libretto/main_test.go TestPrerequisitesIncludeTheGateTools
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

Finding the clone:

- **a directory with a `go.mod` and no `.git` is not the clone** — the module cache is
  exactly that, and accepting it is what would link `~/.claude` into a versioned cache
  Proof: cmd/libretto/root_test.go TestRepoRootRequiresGitDirectory
- a worktree's `.git` *file* counts as a repository
  Proof: cmd/libretto/root_test.go TestRepoRootAcceptsAWorktreeGitFile
- the override wins over every rung and is not validated
  Proof: cmd/libretto/root_test.go TestRepoRootPrefersEnvOverride
- and `repoRoot` really reads it from the environment
  Proof: cmd/libretto/root_test.go TestRepoRootReadsTheOverrideFromTheEnvironment
- the compile-time path beats the working directory
  Proof: cmd/libretto/root_test.go TestRepoRootPrefersTheCompileTimePathOverTheWorkingDirectory
- **the working directory is accepted only when its `go.mod` names this module**
  Proof: cmd/libretto/root_test.go TestRepoRootAcceptsTheWorkingDirectoryOnlyForThisModule
- nothing found resolves to the bootstrap path, not to `$PWD`
  Proof: cmd/libretto/root_test.go TestRepoRootFallsBackToBootstrapPath

Bootstrap:

- **the destination is announced before the clone starts**, not after
  Proof: cmd/libretto/bootstrap_test.go TestBootstrapAnnouncesDestinationBeforeCloning
- a destination we did not create is refused, named, and left untouched
  Proof: cmd/libretto/bootstrap_test.go TestBootstrapRefusesForeignDestination
- an existing clone is a silent no-op, not a re-clone and not noise
  Proof: cmd/libretto/bootstrap_test.go TestBootstrapIsIdempotent
- **a failed clone leaves nothing behind** that a later run would refuse
  Proof: cmd/libretto/bootstrap_test.go TestBootstrapCleansUpAfterAFailedClone
- the requested command runs after the clone — one invocation, no "now run it again"
  Proof: cmd/libretto/bootstrap_test.go TestBootstrapContinuesIntoRequestedCommand
- **`version` and `help` write nothing into the home directory**
  Proof: cmd/libretto/bootstrap_test.go TestVersionAndHelpDoNotBootstrap

The version, and the release state:

- a stamped version wins over build info, dirt and commit count included
  Proof: cmd/libretto/version_test.go TestVersionPrefersLdflagsOverBuildInfo
- `go install`'s module version fills in when ldflags said nothing
  Proof: cmd/libretto/version_test.go TestVersionFallsBackToBuildInfo
- neither source, or `(devel)`, means `dev`
  Proof: cmd/libretto/version_test.go TestVersionIsDevWhenNothingKnows
- the panel's notice names both versions and the action
  Proof: cmd/libretto/version_test.go TestReleaseNoticeNamesBothVersionsAndTheAction
- **and is silent in every case where a row would be a guess** — up to date, ahead of the
  remote, unparseable, a dirty build, nothing cached
  Proof: cmd/libretto/version_test.go TestReleaseNoticeIsSilentWhenThereIsNothingToSay
- `doctor` reports a newer release
  Proof: cmd/libretto/doctor_test.go TestDoctorReportsNewerRelease
- and says so when it is up to date
  Proof: cmd/libretto/doctor_test.go TestDoctorSaysUpToDate
- **a failed check is stated, never silent** — printing nothing would read as up to date
  Proof: cmd/libretto/doctor_test.go TestDoctorSaysSoWhenTheCheckFails
- a remote with no releases is neither a failure nor an update
  Proof: cmd/libretto/doctor_test.go TestDoctorHandlesARemoteWithNoReleases
- **a binary whose version will not parse is given the facts and not ranked**
  Proof: cmd/libretto/doctor_test.go TestDoctorDoesNotRankAnUnidentifiableBinary

The rebuild:

- **it replaces the executable it was given, and leaves the clone's `bin/` alone**
  Proof: cmd/libretto/update_test.go TestRebuildReplacesRunningExecutable
- a symlinked executable is written through, and stays a symlink
  Proof: cmd/libretto/update_test.go TestRebuildResolvesSymlinkedExecutable
- **an unwritable destination is reported with both paths and does not fail the update** —
  what could not be written, and where the new binary is, which is stat'd
  Proof: cmd/libretto/update_test.go TestRebuildReportsUnwritableDestinationWithoutFailing

`models`:

- the listing is the target's agents, not the repository's
  Proof: cmd/libretto/models_test.go TestModelsListsTheTargetsAgents
- **an agent the target has and the repository does not is listed and editable**
  Proof: cmd/libretto/models_test.go TestModelsEditsAnAgentTheRepositoryDoesNotShip
- a symlink into the repository is marked shared; a real file is not
  Proof: cmd/libretto/models_test.go TestModelsMarksSharedAgents
- writing a shared agent says the effect reaches every project
  Proof: cmd/libretto/models_test.go TestModelsSetSaysWhenAWriteIsShared
- **writing a target-local agent does not claim to be machine-wide**
  Proof: cmd/libretto/models_test.go TestModelsSetDoesNotOverclaimALocalWrite
- the two scopes list different agents
  Proof: cmd/libretto/models_test.go TestModelsListingDiffersBetweenScopes
- an agent with no declared model is listed as running the session's
  Proof: cmd/libretto/models_test.go TestModelsShowsDefaultForAnUndeclaredAgent
- `set` applies one model to several named agents
  Proof: cmd/libretto/models_test.go TestModelsSetAppliesToEveryNamedAgent
- `set --all` reaches every agent in the target
  Proof: cmd/libretto/models_test.go TestModelsSetAllReachesEveryAgent
- **`set` with no agents and no `--all` is an error and writes nothing** — a
  destructive default that fires on a forgotten argument is how every agent on the
  machine silently becomes the same model
  Proof: cmd/libretto/models_test.go TestModelsSetWithoutAgentsIsAnError
- an unknown model exits non-zero and writes nothing
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownModel
- an unknown agent name exits non-zero and leaves the valid ones untouched
  Proof: cmd/libretto/models_test.go TestModelsSetRejectsAnUnknownAgentAndWritesNothing
- **a stray file with no frontmatter is refused and the whole set is left alone**
  Proof: cmd/libretto/models_test.go TestModelsSetRefusesAStrayFileAndWritesNothing
- a target with no agents directory says so and exits zero
  Proof: cmd/libretto/models_test.go TestModelsWithNoAgentsSaysSo
- piped `models` carries no escape codes
  Proof: cmd/libretto/models_test.go TestModelsOutputHasNoEscapeCodes
