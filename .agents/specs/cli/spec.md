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
| `update` | install the newest version and relink; in a checkout, pull and rebuild |
| `status` | every item's state, read-only |
| `doctor` | what needs attention, what the payload expects on this machine, and the release state |
| `prune` | show links whose source is gone; `--yes` removes them |
| `uninstall` | show what this repo installed here; `--yes` removes it |
| `preview` | print the panel once, no TUI |
| `models` | list every agent with the model it runs on and the effort it runs at, read-only |
| `models set <model> <agent>…` | write that model into each named agent; `--all` for every one |
| `models effort <level> <agent>…` | write that effort level into each named agent; `--all` for every one |
| `version`, `help` | say so |

### Finding the payload

```bash
go install github.com/pausf/libretto-automata/cmd/libretto@latest
```

**That is the whole install.** The payload ships inside the Go module, so this brings it down with
the binary and Go verifies it against the checksum database — see `distribution`.

**The payload root and the checkout are two different things**, and this tool conflated them
for as long as a clone happened to be both. `go install` breaks the coincidence: there is a
payload and no checkout. Rungs 2 and 3 below answer both questions at once; rung 4 answers
only the first, which is why `update` probes for a checkout before choosing its route.

**The payload root is resolved in four rungs:**

| 1 | `LIBRETTO_ROOT` | absolute override, taken as given, unvalidated |
| 2 | the compile-time source directory | when it has `.git` — development, `make build` |
| 3 | the working directory | when it has `.git` **and** a `go.mod` naming this module |
| 4 | `$GOMODCACHE/<module>@<version>` | the payload `go install` brought down |

**The probe is `.git`, not `go.mod`.** A `go.mod` matches the module cache, which is rung 4's
answer and not rung 2's — so the two would be indistinguishable and the wrong one would win. A
`.git` **file** counts too, because that is what a worktree has and this flow recommends
worktrees.

**Rung 4 is this binary's own version**, from build info. The payload it links is the one that
shipped with the command doing the linking, so the two can never be a version apart — and the
version being *in the path* is what makes an update a new directory rather than an overwrite of
the one the links currently point at.

Rung 3 checks the module, not just the presence of a repository. Any git repository
satisfied the old `$PWD` fallback, so `libretto install` inside an unrelated project went
looking for a payload that project does not have.

**Nothing found means nothing is installed yet**, and rung 4 returns a path that may not exist
— an interrupted `go install`, or a `go clean -modcache`. `update` is what fixes it.

**Two other homes were tried in one session and both are gone.** `~/.libretto-automata` reached
by a `git clone`, then `~/.local/share/libretto/<version>` with a release tarball extracted into
it. Neither shipped in a tag, so nothing migrates and no machine has one. Why each lost is in
`distribution`, recorded because a rejected design with no reason attached gets proposed again.

- **`version` and `help` are answered before the payload is even located.** Neither reads a
  skill, so neither should care whether one is installed — and neither writes anything
  anywhere.
- **Nothing prompts.** Every path works without a TTY, or `install` in CI breaks.

### A missing payload is a stop, not an empty report

Every command that reads the tree refuses and points at `update`, rather than describing an empty
one. "nothing to link" is what a correctly linked machine also says.

**`prune` is why this is a stop and not a warning.** With no payload every link in the target
resolves to nothing, so a scan reports all of them `stale` — and `prune --yes` would remove every
item the user has. A destructive command doing exactly what it promises, on a premise that is
false.

- **`models` is exempt.** It reads the *target's* agents directory, not the payload, so it works
  on a machine with nothing installed — which is when somebody might most want it.
- **`update` is exempt**, because it is what fixes the state.

### `update` brings the installation up to date

| in a checkout | pull, rebuild when Go moved, relink |
| anywhere else | `go install` the newest version, relink |

**One command, because for the person typing it there is one meaning:** put me on the newest
version. Which mechanism runs is a consequence of how the tool got onto the machine, which they
know.

- **The release route says `git` nowhere**, in success or in any failure. A pull announced to
  somebody who only wanted to use the tool is the complaint that produced this whole shape.
- **It relinks the version it just installed**, not `payloadRoot()`. This process still reports
  the old version, so resolving again would link the payload it is already on — and the update
  would report success having changed no link at all.
- **Relinking is not redundant.** A new version is a new directory, and a release that *adds* an
  item leaves that item unlinked with nothing to say so.
- **Already newest changes nothing and says so.** Nothing installed, nothing relinked.
- **Nothing installed yet updates rather than refusing.** That is the repair path for a cleaned
  module cache.

**A draft split this into `update` and `upgrade`** and defended it: *a command whose mechanism
depends on invisible state has unpredictable failures*. The state is not invisible — it is "did
you clone this or install it" — and the split cost two commands, two menu rows, a pair of mutual
refusals, a help table that changed per machine, and a function whose only job was choosing which
of the two names to print in a notice.

### The rebuild replaces the binary that is running

`update`'s git route rebuilds over `os.Executable()`, not over the clone's `bin/libretto`. Once the
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
- **`doctor` names the mode it is in**, because "up to date" means something different in each: a
  checkout three commits past a tag is current with its remote and behind the release.
- **`doctor` always says something.** Newer available, up to date, none published,
  could not check, or "running X, the latest is Y" for a binary ahead of the remote or one
  whose version will not parse. Printing nothing on a failed check would read as "you are up
  to date", which is a claim nobody verified. That last case gets the facts and **no
  ranking**.
- **It never sets the exit code.** Being a release behind is news, and an unreachable remote
  is not this tool's fault.
- **The panel uses the cache, not a live call** — see `panel`. A panel that waits on the
  network before its first paint hangs on bad DNS. Each mode keeps its answer somewhere it has: a
  checkout in `.git/`, an installed copy in the module cache root. One TTL and one failure policy,
  `repo`'s.
- **Which source is asked depends on the mode, and each mode has exactly one.** A checkout asks
  its remote what it has tagged; an installed copy asks the module proxy what the project has
  published. Two *questions*, not two answers to one, and only one is ever asked on a machine.
- **The notice names `update`**, which is the only command there is.
- **And nothing in the panel names a mechanism.** The `update` row reads `bring this installation
  up to date` on every machine — see `panel`. What ran is in the command's output, where a reader
  who wants it will look; a menu label that says `pull` is the original complaint in quieter words.
- **Every subcommand that reads the payload ends by saying it too** — `status`, `preview`,
  `install`, `prune`, `uninstall`, `models`. The panel needs a TTY and `doctor` needs somebody to
  already suspect there is something to diagnose; the payload is used inside Claude Code, so after
  installing most people run neither for months and the notice sat where the user was not. Now
  *any* run says it.
- **On stderr, after the command's own output.** `status` output is parseable and somebody may be
  parsing it, so stdout stays byte-identical to a run with nothing to say. After, not before,
  because the notice is news about something else.
- **It never changes the exit code.** Being a release behind is not an error, and `install`
  already spends a non-zero exit on a conflict — a second meaning would make the first unreadable.
- **Cached, and silent when there is nothing to say** — not newer, no cached answer, or a failed
  check all print nothing. A subcommand is not a diagnostic: it does not get to pay for a network
  call the panel would not have made, and "could not check" from a command about something else is
  the noise that gets a notice ignored.
- **`doctor` and `update` are excluded.** `doctor` already says something on every path, live, and
  a second line would print the same fact twice with one copy stale. `update` *is* the update.
- **The line is never coloured.** The panel's row is gold because it sits inside a rendered frame;
  a line appended to a command's output is not, and colouring it buys a branch, a theme dependency
  and a test for something nobody asked for.
- **Ceiling, named: this reaches nobody who does not run the binary.** What is guaranteed is that
  every run tells them. The only thing that reaches somebody who never runs it is the payload
  saying it, and that costs `payload`'s promise that an installed skill works in a project which
  has never heard of `libretto`. Not paid.

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

**A verb the dispatch accepts, `help` offers.** `models effort` shipped working — the
command, the panel key, the footer legend — and was reported as missing, because `help`
listed `models set` and stopped there. A feature nobody can find is a feature that was
not delivered, and the help text is the door people try. The check reads the verbs out of
the dispatch's own unknown-command error rather than from a second list, so a third verb
added without a help line fails rather than passing quietly.

**The model and the effort are two verbs, not one verb with a flag.** `models set` writes
the tier; `models effort` writes the depth. `set opus --effort xhigh` was the alternative
and it forces the model to be restated to change the effort, which is a write nobody
asked for — the keys are independent and the surface says so. Both take agent names or
`--all`, both refuse when given neither, and that refusal is one implementation rather
than two: it is the thing standing between a forgotten argument and every agent on the
machine becoming the same thing.

Three things the listing says that it would be dishonest to leave out:

- **`(session)` in both value columns** for an agent declaring neither key. An empty
  column reads as a bug, and it is one word for one state.
- **which model has no effort levels at all, and which has only some.** `haiku` carries
  `— no effort levels`; a model with four of the five carries `— effort: low medium high
  max`. A row sitting silently among others that support effort reads as a fourth that
  does, and the refusal further down would then arrive as a surprise.
- **which provider the versions were resolved for**, on its own line. The version column
  is only checkable against something: a reader who disagrees with `sonnet · Sonnet 4.5`
  needs to know it was read off the Amazon Bedrock row, and the provider is that row. It
  is derived from the environment — see `agent-models` — so it changes with the machine
  rather than with this table.
- **when the effort table was last checked**, under the same `Resolved` date as the model
  versions. It decays for the same reason — an organisation can cap which levels are
  available and nothing here can ask — and a trailer that carries the date next to one
  that does not is a claim about which of the two is current.

**`models set` reports a cleared effort on the row it cleared.** Moving an agent to a
model with no effort support drops its `effort:` key, and a drop the user cannot see is a
silent edit to a prompt file.

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
- **npm, npx, Homebrew, a tap, `install.sh`, per-platform binaries.** `go install` is the one
  entry point, and it leaves nothing else to do. Each of these was considered; each adds a second
  way in.
- **a `sync` command.** `gentle-ai` needs one because its assets live inside its binary and have
  to be written out. Here the payload is already files, and relinking is what `install` does.
- **downloading or verifying anything.** `distribution` runs `go install`; the go command does
  the rest.
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

**Go is required, and that is a limit on who can use this rather than an oversight.** `go install`
is the entry point and the update mechanism both. The audience for the payload is not necessarily
a Go developer; named here so the constraint is visible.

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

**Nothing in the test suite writes to a real payload home, and nothing runs a real `go install`.**
`LIBRETTO_ROOT` is the second half of what `CLAUDE_HOME` does for targets, and `resolveRoot` takes
its inputs rather than reading them — `runtime.Caller` is fixed at compile time, so inside this
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
- **The payload ships inside the Go module.** Not embedded in the binary, not a release asset,
  not a clone. See `distribution` for all three and why the first two lost.
- **One `update`, not `update` and `upgrade`.** Asked and answered by the user, with the split's
  own argument recorded above so it is not re-made.
- After a rebuild the process states that it is still the old binary, rather than
  implying the update took effect mid-run — **and only when the binary was really
  replaced.** On the refused path it says the opposite.
- **A checkout you are standing in wins over the module cache.** That is what keeps editing a
  skill and seeing it live working, and it is the reason the payload is not compiled in.
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
- [x] ~~bootstrap: announce, clone, refuse a foreign destination, clean up a failure~~
- [x] ~~a release tarball, downloaded, verified and extracted~~
      both removed in the same session — see `distribution` for what replaced them
- [x] `update`'s two routes, and the missing-payload stop
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
- **the release notice goes to stderr, and stdout is byte-identical to a run with nothing
  to say** — verified by mutation: pointing the dispatch's `defer` at `os.Stdout` fails it
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeGoesToStderr
- **nothing is printed when there is no newer release, no cached answer, or a failed
  check** — through the formatter itself, against a pre-answered cache rather than a stub
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeIsSilentWithNothingToSay
- **the exit code is whatever the subcommand returned**, on a succeeding path and on
  `models nonesuch` — a command that both carries the notice and fails
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeDoesNotChangeTheExitCode
- **`doctor` and `update` do not carry it**, so the fact is never printed twice
  Proof: cmd/libretto/notice_test.go TestDoctorAndUpdateDoNotRepeatTheNotice
- **the line carries no escape codes**, checked against what the formatter really builds
  Proof: cmd/libretto/notice_test.go TestSubcommandNoticeHasNoEscapeCodes
- **the suite never reaches the module proxy for it.** The dispatch's `defer` fires for
  every test that calls `run`, so a test that did not know about the notice inherited a
  live lookup and a write of `.update-check` into the user's own module cache. `TestMain`
  silences the source package-wide and tests opt in — the same reason `CLAUDE_HOME` exists,
  applied to a directory it does not cover.
  Proof: cmd/libretto/notice_test.go TestMain
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

Finding the payload:

- **a directory with a `go.mod` and no `.git` is not rung 2** — the module cache is exactly that,
  and it is rung 4's answer, so accepting it earlier would make the two indistinguishable
  Proof: cmd/libretto/root_test.go TestPayloadRootRequiresGitDirectory
- a worktree's `.git` *file* counts as a checkout
  Proof: cmd/libretto/root_test.go TestPayloadRootAcceptsAWorktreeGitFile
- the override wins over every rung and is not validated
  Proof: cmd/libretto/root_test.go TestPayloadRootPrefersEnvOverride
- and `payloadRoot` really reads it from the environment
  Proof: cmd/libretto/root_test.go TestPayloadRootReadsTheOverrideFromTheEnvironment
- the compile-time path beats the working directory
  Proof: cmd/libretto/root_test.go TestPayloadRootPrefersTheCompileTimePathOverTheWorkingDirectory
- **the working directory is accepted only when its `go.mod` names this module**
  Proof: cmd/libretto/root_test.go TestPayloadRootAcceptsTheWorkingDirectoryOnlyForThisModule
- **nothing found resolves to the module cache entry for this binary's own version**, with the
  version in the path — not to `$PWD`, and not to either abandoned payload home
  Proof: cmd/libretto/root_test.go TestPayloadRootFallsBackToTheActivatedRelease
- **a checkout you are standing in still wins**, which is what keeps editing a skill and seeing it
  live working
  Proof: cmd/libretto/root_test.go TestPayloadRootStillPrefersACheckoutYouAreStandingIn

`update`:

- **it installs, then relinks the version it installed** — not the one the running process
  reports, which would link the payload it is already on
  Proof: cmd/libretto/update_release_test.go TestUpdateInstallsThenRelinksTheVersionItInstalled
  Proof: cmd/libretto/update_release_test.go TestUpdateRelinksTheNewVersionNotTheRunningOne
- relinking is last, so a version that adds an item does not leave it unlinked
  Proof: cmd/libretto/update_release_test.go TestUpdateRelinksSoNewItemsAppear
- **the release route says `git` nowhere**, in success or in any of its failures
  Proof: cmd/libretto/update_release_test.go TestUpdateFromAReleaseNeverMentionsGit
- a failure names which step, and changes nothing
  Proof: cmd/libretto/update_release_test.go TestUpdateReportsWhichStepFailed
  Proof: cmd/libretto/update_release_test.go TestAFailedUpdateLeavesThePreviousVersionActive
- already newest is a success that changes nothing
  Proof: cmd/libretto/update_release_test.go TestUpdateOnTheNewestVersionDoesNothing
- nothing installed yet updates rather than refusing
  Proof: cmd/libretto/update_release_test.go TestUpdateFromNothingInstalled
- **the route follows how this installation arrived**
  Proof: cmd/libretto/update_release_test.go TestUpdateTakesTheRouteThisInstallationCameBy

A missing payload:

- `install` explains itself instead of reporting an empty tree
  Proof: cmd/libretto/update_release_test.go TestInstallWithNoPayloadPointsAtUpdate
- **`prune --yes` refuses instead of deleting every link**, which is what a scan with no payload
  would otherwise licence it to do
  Proof: cmd/libretto/update_release_test.go TestPruneWithNoPayloadRefusesInsteadOfDeletingEverything
- `models` still works, because it reads the target and not the payload
  Proof: cmd/libretto/update_release_test.go TestModelsWorksWithNoPayload
- **`version` and `help` write nothing anywhere**
  Proof: cmd/libretto/root_test.go TestVersionAndHelpTouchNothing

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

`models effort`:

- the listing shows each agent's effort, and `(session)` when it declares none
  Proof: cmd/libretto/models_test.go TestModelsListsEffortBesideTheModel
- the trailer names the five levels, says which model has none, and states when it was
  last checked
  Proof: cmd/libretto/models_test.go TestModelsListsTheEffortCatalogue
- **the listing names the provider it resolved the versions for**
  Proof: cmd/libretto/models_test.go TestModelsNamesTheProviderItResolvedFor
- **under a provider where `sonnet` is older, the listing says so and offers it no levels**
  Proof: cmd/libretto/models_test.go TestModelsReflectsTheProviderInTheVersionColumn
- **the refusal follows the provider and not the alias** — `xhigh` on a sonnet agent is
  accepted on the Anthropic API and refused under Amazon Bedrock
  Proof: cmd/libretto/models_test.go TestModelsEffortFollowsTheProviderNotTheAlias
- `effort` writes the level to the named agents and no others
  Proof: cmd/libretto/models_test.go TestModelsEffortWritesOnlyTheNamedAgents
- `effort --all` reaches every agent in the destination
  Proof: cmd/libretto/models_test.go TestModelsEffortAllReachesEveryAgent
- **no agents named and no `--all` is refused, and nothing is written**
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesWithNothingNamed
- an unknown level is refused and the refusal names the five
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesAnUnknownLevel
- **an agent on Haiku is refused by name and the whole set is left alone**
  Proof: cmd/libretto/models_test.go TestModelsEffortRefusesAnAgentThatCannotRunTheLevel
- `default` removes the key rather than writing a word
  Proof: cmd/libretto/models_test.go TestModelsEffortDefaultRemovesTheKey
- **`models set haiku` on an agent declaring an effort reports the clearing**
  Proof: cmd/libretto/models_test.go TestModelsSetReportsAClearedEffort
- **every verb the dispatch accepts is offered by `help`**, read out of the dispatch's
  own error rather than from a list kept here
  Proof: cmd/libretto/models_test.go TestEveryModelsVerbIsInTheHelp

### `loop` — one fresh session per open box

- **the engine is the binary, not a skill, and that is not a preference.** A loop that
  relaunches a session cannot *be* that session: the whole value of the pattern is a
  fresh context per task, and a skill runs inside the context it would be trying to
  discard. The Ralph playbook ships a `loop.sh` beside four other files, and every one of
  those four already has a richer equivalent here — `write-plan`, `build-and-check`,
  `AGENTS.md` and the plan's own boxes. Copying the structure would duplicate state; the
  engine was the only piece genuinely missing.
  Proof: cmd/libretto/loop_test.go TestTheLoopPromptRoutesAndNeverAuthorises
- **it owns no state.** `.agents/changes/<change>/plan.md` is the state — the file phase 5
  writes and phase 6 marks — and the loop reads the boxes and nothing else. Only the box is
  parsed, never the task text: the flow owns the plan's shape, and a runner that understood
  it would be a second opinion about what a task is.
  Proof: cmd/libretto/loop_test.go TestCountBoxesReadsOnlyTheBox
- one session per open box, and it stops when the last one closes
  Proof: cmd/libretto/loop_test.go TestLoopStopsWhenEveryBoxIsClosed
- **two consecutive rounds that close nothing stop the loop.** This is the guardrail the
  whole command exists for. A session that finished no task leaves a plan the next fresh
  session reads identically, so it makes the identically non-existent progress — forever,
  or until the cap. One barren round is a hiccup and the counter resets; two is a loop that
  has stopped being one, and the message says why a third is pointless.
  Proof: cmd/libretto/loop_test.go TestLoopStopsAfterTwoRoundsThatCloseNothing
- a single barren round does not stop it
  Proof: cmd/libretto/loop_test.go TestOneBarrenRoundDoesNotStopTheLoop
- **progress is boxes closed, never boxes remaining.** A session that finishes one task and
  splits another in two leaves the open count unchanged, and reading that as no progress
  stops the loop reporting "two rounds closed nothing" when two things were done. A plan is
  allowed to grow: phase 6 discovering a task is the flow working.
  Proof: cmd/libretto/loop_test.go TestAPlanThatGrowsIsProgressNotAStall
- **the missing plan is reported before the missing dependency.** `loop <typo>` on a machine
  without `claude` named the dependency rather than the mistake actually made.
  Proof: cmd/libretto/loop_test.go TestTheMissingPlanIsReportedBeforeTheMissingDependency
- **`--dry-run` works with an empty PATH**, proved through `loop` rather than `runLoop` —
  the guard sits above `runLoop`, so a test calling it directly proves nothing about it.
  Proof: cmd/libretto/loop_test.go TestDryRunNeedsNoClaudeOnPath
- **the cap is the second ceiling, and hitting it is an error naming the flag that raises
  it.** A loop with no ceiling that never converges spends a budget silently.
  Proof: cmd/libretto/loop_test.go TestLoopStopsAtTheCapAndSaysSo
- **a session exiting non-zero is not a failed loop.** The plan is the state, so the next
  round carries on from whatever this one finished; an exit code cannot tell a crash from a
  deliberate stop, and the boxes can.
  Proof: cmd/libretto/loop_test.go TestASessionErrorDoesNotAbortTheLoop
- a missing plan is refused, and the refusal names the path it looked for
  Proof: cmd/libretto/loop_test.go TestLoopRefusesAPlanThatIsNotThere
- **a plan with prose and no boxes is refused rather than read as finished.** Zero open
  boxes means done everywhere else in the flow; here it also means a plan phase 5 never
  wrote, and reporting success having run nothing is the worse of the two readings.
  Proof: cmd/libretto/loop_test.go TestLoopRefusesAPlanWithNoBoxes
- `--dry-run` prints the prompt and starts no session, and needs no `claude` on PATH
  Proof: cmd/libretto/loop_test.go TestDryRunPrintsThePromptAndStartsNoSession
- flags are parsed strictly: one change, a positive `--max`, no unknown flags
  Proof: cmd/libretto/loop_test.go TestParseLoopArgs
- **it is not gated on the payload.** `loop` reads the *project's* plan and relaunches a
  session that resolves its own skills; refusing without this repository's payload tree
  would refuse on exactly the machine that installed the binary and nothing else.
  Proof: cmd/libretto/loop_test.go TestLoopIsNotGatedOnThePayload
- **it never pushes, merges or tags, and the prompt forbids rather than omits.** Those are
  what `libretto-attacca` answers for one branch after a person typed it; nothing here was
  asked, and an omission reads as an oversight to whatever is on the other end.
  Proof: cmd/libretto/loop_test.go TestTheLoopPromptRoutesAndNeverAuthorises

### `metrics` — what the flow cost, derived and never recorded

- **nothing is instrumented.** The obvious build has each phase write a timestamp and the
  report read them back; that costs a new artifact in every change folder which eight
  skills have to remember, and a metric collected only when somebody remembers has holes
  exactly where the interesting runs are. Git already holds the answer: a change folder's
  first and last commits are the wall clock, the commits between are the iterations, and
  `plan.md`'s diffs are the churn. It works retroactively on every change ever landed.
  Proof: cmd/libretto/metrics_test.go TestChangeNamesSeesLandedChangesNotJustOpenOnes
- **`--full-history` on every query, and it is load-bearing.** Git's default history
  simplification prunes commits that do not change a path's *final* state, and a landed
  change's final state is "does not exist" — so the entire history of every change this
  report exists to measure is simplified away. Measured, not reasoned: without it this
  repository reported 12 changes instead of 22, two of them with no history at all, and
  **no `plan.md` was found for any of them**, so every churn column read as "no plan" on
  changes that plainly had one.
  Proof: cmd/libretto/metrics_test.go TestChangeNamesSeesLandedChangesNotJustOpenOnes
- **the churn query excludes deletions, or every number inverts.** A change lands by
  deleting its folder, and that commit's diff removes every line of `plan.md` — so every
  box ever closed also appears as a removed `[x]`, and the report reads as 0 closed and 52
  reopened on a change that had neither. `--diff-filter=AM`. **Both of these were found by
  running it against real history, not by the fixtures**, which is why the fixtures now
  refuse a query missing either flag.
  Proof: cmd/libretto/metrics_test.go TestReopenedBoxesAreCountedSeparatelyFromClosedOnes
- **a reopened box is counted apart from a closed one.** Closed-then-reopened means a task
  was called done before it was, which is the one thing in a plan's history worth knowing,
  and a net count is precisely what hides it.
  Proof: cmd/libretto/metrics_test.go TestReopenedBoxesAreCountedSeparatelyFromClosedOnes
- git log is newest-first, and reading it as oldest-first inverts the span into a negative
  duration that prints as a plausible number
  Proof: cmd/libretto/metrics_test.go TestMeasureReadsSpanAndCommitsOldestLast
- **a change with no plan reports a dash, never a zero.** "No plan existed" and "a plan
  existed and nothing was finished" are different facts, and one of them is an accusation.
  Proof: cmd/libretto/metrics_test.go TestAChangeWithNoPlanReportsADashNotAZero
- **an unreadable change prints a row rather than being skipped.** The footer counts names,
  so a silent skip makes the total disagree with the rows above it — a report claiming
  twelve while showing ten is worse than one admitting it could not read two.
  Proof: cmd/libretto/metrics_test.go TestAnUnreadableChangePrintsARowRatherThanVanishing
- **a reworded or deleted finished task is not a reopening.** The churn arithmetic is net
  per commit, not summed per line, and that is the whole correctness argument: rewording a
  closed box emits a removed `[x]` and an added one in the same commit, so line-counting
  reported one close and one reopen for a typo fix. A reopening is the accusation "this was
  called done before it was" and it has to be earned. Net zero within a commit means the
  text moved and the state did not.
  Proof: cmd/libretto/metrics_test.go TestRewordingAndDeletingAFinishedTaskAreNotReopenings
- **the two commands read a checkbox identically.** `metrics` counted a box anywhere in a
  line while `loop` anchors to the start, so one `plan.md` gave two different totals
  depending on which command read it. `metrics` strips the diff marker and asks `loop`'s
  parser; there is one definition of a box.
  Proof: cmd/libretto/metrics_test.go TestBoxInReadsOnlyRealCheckboxes
- **a truncated diff line does not panic.** The bounds guard was off by one — `i+4 > len(l)`
  where it needed `>=` — so a line ending one character after `- [` indexed past the end.
  The table that was meant to cover it stopped exactly one character short of the crash,
  which is a boundary case that reads as covered and is not.
  Proof: cmd/libretto/metrics_test.go TestBoxInReadsOnlyRealCheckboxes
- **it asks git for the repository root rather than trusting the cwd.** git pathspecs and
  `os.Stat` are both cwd-relative, so run from a subdirectory this reported "no changes in
  this repository's history yet" — plausible, and false.
  Proof: cmd/libretto/metrics_test.go TestMetricsFiltersToOneChangeAndRefusesAnUnknownOne
- **the report names the two things it cannot measure.** Per-phase duration and
  `review-work` finding counts both need a phase to write them down: phases 1 to 7 happen
  in one session and leave one commit, and a finding is repaired before anything lands so
  only the repair is in the diff. A metrics command silently omitting two of the three
  things asked of it reads as having measured them and found nothing.
  Proof: cmd/libretto/metrics_test.go TestTheReportNamesWhatItCannotMeasure
- a span drops precision it does not have — days, not minutes, on a multi-day change
  Proof: cmd/libretto/metrics_test.go TestHumanSpanDropsPrecisionItDoesNotHave
- **the footer total merges overlapping spans.** It says wall clock, so each calendar
  hour counts once however many changes were open during it — the plain sum reported two
  weeks for two changes open the same week, a number that was nobody's clock. Per-row
  spans are untouched; only the total merges.
  Proof: cmd/libretto/metrics_test.go TestTotalSpanMergesOverlappingChanges
- **the closed cell carries its denominator.** `n/total` — boxes closed now over boxes
  the plan holds — because a bare 5 hides whether the plan had 5 boxes or 18, and those
  are opposite facts about a change in flight. The numerator is net current state;
  cumulative churn stays in the reopen column. A reword moves neither number, and a
  plan-less change keeps its dash, never `0/0`.
  Proof: cmd/libretto/metrics_test.go TestClosedShowsItsDenominator
- **the report explains its own columns.** A legend beside the not-measured note names
  the six measured facts, including the `unreadable` state and the `—` cell. It already
  printed what it cannot measure; what it does measure deserves no less.
  Proof: cmd/libretto/metrics_test.go TestTheReportExplainsItsColumns
- one change may be named, by full name or unambiguous prefix — a name is typed from
  memory of how it starts. Exact wins even when it prefixes a sibling, an ambiguous
  prefix is refused naming the candidates, and an unknown name and an unknown flag keep
  their refusals.
  Proof: cmd/libretto/metrics_test.go TestAPrefixSelectsAChangeUnlessAmbiguous
  Proof: cmd/libretto/metrics_test.go TestMetricsFiltersToOneChangeAndRefusesAnUnknownOne
- it reads the project's git history, so it is not gated on the payload
  Proof: cmd/libretto/metrics_test.go TestMetricsIsNotGatedOnThePayload
- **corrections are counted per change, read off the lessons ledger.** `.agents/lessons.md`
  is the one artifact the payload writes and this command only counts — it exists for the
  retro, written by one skill, and metrics is a free rider, which is the bar the
  no-instrumentation rule above sets for an artifact. The header is the contract: a line
  matches when it starts with `## ` and carries exactly two ` · ` separators with three
  non-empty fields; the date is not validated, because a misspelled date is still a
  countable lesson.
  Proof: cmd/libretto/metrics_test.go TestCorrectionsAreCountedPerChange
- **no ledger reports a dash, never a zero.** Absent means capture is not in use; zero
  means the flow ran and was never corrected. Printing 0 for the first claims the second.
  Proof: cmd/libretto/metrics_test.go TestNoLedgerReportsADashNotAZero
- **a malformed header is skipped and a changeless correction is named, never lost.** The
  ledger is written by prompts, so a parser that dies on one bad line loses the whole
  report; entries whose change field is `-` (no change open) belong to no row and are
  reported in one line instead of silently dropped.
  Proof: cmd/libretto/metrics_test.go TestMalformedAndChangelessEntriesDoNotCrashTheCount
