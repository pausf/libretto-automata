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

`status`, `doctor`, `metrics` and `queue` also answer with `--json`: one JSON object on
stdout in place of the text report, with exit codes, stderr and the release notice
unchanged — see *`--json` — the read-only reports, machine-readable*.

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

Two axes, combinable — a tool flag and a scope flag answer different questions:

| Invocation | Meaning |
|---|---|
| `--claude` | Claude Code — the default tool |
| `--codex` | Codex CLI, skills only |
| `--opencode` | OpenCode, all three kinds — agents generated, not linked |
| `--global` / `-g` | the tool's machine-wide root — the default scope |
| `--project` / `-p` | the tool's directory inside `<cwd>` |
| none | **claude/global** for every subcommand — what every invocation meant before the axes existed. The panel opens where it was left |

`--codex --project` acts on `<cwd>/.agents`; `--opencode` alone on
`~/.config/opencode`. Two flags on **one axis** is an **error**, not a precedence
rule — two answers to one question is a mistake worth reporting rather than
resolving by picking the last one and hoping. No short flags for the tools:
`-g`/`-p` predate them, and a `-c` that might one day mean something else is not
worth squatting.

**A destination installs the kinds it accepts and nothing else**, and the kinds it does
not accept are absent from the run and from the summary, never errors. `--codex` is the
skills-only one: `install --codex` links skills alone, and it reaches OpenCode too
because OpenCode reads `~/.agents/skills` — a fact the README carries.
`install --opencode` links skills and commands and **writes** agents: that destination
takes all three, and one of them is generated rather than linked. `status` and `doctor`
report a generated agent through the same five states with no new vocabulary — drift reads
`wrong target`, a file of somebody's own at that path reads `conflict`, an agent whose
source is gone reads `stale`. `uninstall --opencode --yes` takes ours back out and keeps
theirs.

The paragraph used to say "a skills-only destination" and mean both new tools. The rule
was right and the list was not, which is the shape a sentence goes stale in: it names
the members instead of the property, so adding a kind to one tool makes it false without
touching a word of it.

The remembered panel destination is a (tool, scope) pair; a legacy one-word file
still reads, and anything unrecognised falls back to claude/global.

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
prerequisite report, `init`'s bootstrap of a project's `.agents/` layout, and the
`--json` shapes of `status`, `doctor`, `metrics` and `queue`.

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
- **`land` performing the landing.** No write mode, no staging, no `git` mutation of
  any kind. Verify-only is forward-compatible: a write mode is a later change, not a
  hedge.
- **`land` checking part 1 or part 3 of the landing contract.** The final code is
  unverifiable mechanically, and part 3 is owned by `spec-drift --retired` — `land`
  names the owner and checks nothing.
- **`land` running `spec-drift`, or any gate.** One more voice, not an orchestrator.
  It verifies no delta content either — EARS, `Proof:` citations, pillar structure are
  `--anchors`' to fail.
- **a `--commit <sha>` mode for `land`.** The staged index is the default, and the
  decision log records that a commit mode can be added later without moving it.
- **`init` writing anything beyond its three artifacts** — no sample spec, no
  `docs/SPEC.md`, no `CLAUDE.md` or `AGENTS.md`, and never `.agents/specs/README.md`,
  which is `wiki`'s by marker. The list grows one file at a time, if ever.
- **`init` adopting or migrating a foreign layout**, and any `--layout` flag to do so.
  Refusal is the whole answer; a flag is a later change if somebody asks for one.
- **`init` writing to git.** No commit, no staging: one `rev-parse --show-toplevel` is
  the only git it runs. Committing the user's repository is the user's act.
- **`--force` for `init`.** There is no overwrite path to force.
- **a loop-wide spending cap.** The header prints the worst case instead. Returns if a
  ceiling across rounds is asked for; the run log already holds what it would sum.
- **a loop budget in rounds, wall time or tokens.** `--max` is already the round cap and
  the host enforces dollars; a second unit is a second budget.
- **stopping the loop on the first budget hit**, or reading back anything the host
  reports. Stuck detection already bounds a run of budget-cut rounds.
- **reading the run log anywhere but `metrics`** — no `status` line, no panel row.
- **codex in the loop.** `--codex` is refused, never driven; see *Prior decisions* for
  the evidence and for what brings it back.
- **a `--tool` flag on `loop`.** The tool axis already answers the question, and a second
  spelling of it would be two answers to one question.
- **a budget or a cost for `opencode`, by any mechanism** — a wall-time kill, a token
  count, a loop-side sum, a figure read back from `opencode export` or `opencode stats`.
  `--max-budget-usd` and the JSON result are `claude`'s; `opencode run` has no per-session
  budget to pass an amount to, and its stream omits every other session's cost. See
  *Prior decisions*.
- **permission flags for either tool, and any model, agent or variant choice.** The user's
  own configuration governs both tools; the loop grants nothing and chooses nothing.
- **`--project` meaning anything to `loop`.** `scopeFlags` strips it and the loop ignores
  it.
- **checking that the payload is installed into the chosen tool** before the first round,
  for either tool. A missing command surfaces as the session's own error and as two rounds
  that closed nothing.
- **`metrics` reading `tool`**, and any tool column in its report. An opencode round
  already counts as unpriced.
- **rotating, pruning or migrating the run log.** A record is a few hundred bytes; rotation
  is a later change with its own ceiling.
- **a checkpoint written by the loop.** Each session commits its box, and the checklist is
  the state; nothing new is written for it.
- **a panel row for `init`.** It is a once-per-project command typed on purpose; a row
  earns its place when somebody misses it.
- **the payload half of bootstrap.** The non-goal that keeps skills from creating
  `changes/` implicitly, and the pointers to `init` at its two refusal sites, are
  `payload`'s.
- **`queue` editing, deleting, reordering or picking up a queued idea.** `/libretto-next`
  picks up; the files are edited by hand. The payload's non-goals on that stand.
- **`queue` taking any argument but `--json`** — no name or date filter. A filter is a
  later change.
- **`queue` reading remote branches, or resolving a base branch.** Local branches unmerged
  into `HEAD` are find-work's filter as the skill runs it; see *Prior decisions* for the
  ceiling.
- **`queue` listing changes in flight, or reading a tracker.** That is source 1 and
  `/libretto-status`; the command lists the queue and nothing else.
- **a panel row for `queue`**, for the reason `init` has none.
- **`--json` on `models`, `land`, `version`, `prune`, `uninstall`, `wiki`, `init`, `loop`
  or `preview`.** Returns when a consumer names one; the shared helper makes each one box,
  with its own shape.
- **JSON errors on stderr.** Errors stay `main`'s plain line; returns if a consumer cannot
  use the exit code.
- **a global `--json`, a `--format` flag, or JSON Lines streaming.** One flag, four
  commands, one document each.
- **a published JSON Schema file.** The shapes live in this spec; a schema file is a second
  copy of them.
- **prose in the JSON**: `doctor`'s remedies and prerequisite notes, `metrics`' legend,
  ceiling, footer totals and token percentage, `status`'s tally. Each is an additive field
  later.
- **golden files for the JSON.** Tests decode the output and assert keys; a golden file
  pins whitespace the contract does not promise.
- **`status` and `doctor` validating their arguments.** They read `--json` and ignore the
  rest, as they always have. Returns as its own change if a typo printing text into a
  parser bites somebody.
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

**Seven environment variables, each with a working default:**

| Variable | Effect |
|---|---|
| `CLAUDE_HOME` | Claude Code's root instead of `~/.claude`. What makes the test suite safe. |
| `AGENTS_HOME` | the codex destination's root instead of `~/.agents`. Libretto-only; Codex does not read it |
| `OPENCODE_HOME` | the opencode destination's root instead of `~/.config/opencode`. Libretto-only |
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

**`land` is read-only, absolutely.** No file written, no index touched, no ref moved —
the repository's bytes before and after a run are identical, and that is tested rather
than assumed, because the command runs immediately before the most destructive commit in
the flow. Stdlib plus exec'd git only, through the same `gitRunner`-style seam `metrics`
drives, so tests parse real output and the integration tests build real temporary
repositories (`make test-short` skips those, per the suite's convention). It needs no
payload: it reads the project being landed, not `~/.claude` and not the payload tree, so
it joins the `needsPayload` exemptions beside `models`, `update`, `loop` and `metrics` —
a machine that installed the binary and nothing else is exactly the machine record-work
invokes it on. Outside a git repository it refuses with the same shape as `metrics`: an
error, not an empty report. Warnings go to stderr and the report to stdout; the exit
code means the contract and nothing else, and the wiki warning never changes it.

**`init` never overwrites anything, at all.** The repository rule is "never overwrite what
the tool did not create"; for `init` it tightens to never overwriting anything, because
`init` cannot tell its own `.gitkeep` from anyone else's and does not need to — present is
present, and a path that is a symlink, dangling or not, is present and never written
through. Stdlib only. The repository root comes from one `git rev-parse --show-toplevel`
through the `gitRunner` seam `metrics` defined, injected so tests drive it with fixed
output, and a failure refuses with `metrics`' wording — not a git repository, or git is
unavailable — on stderr with exit 1. The foreign-root lists it probes mirror `SPEC_ROOTS`
and `CHANGE_ROOTS` in `spec-drift`, minus the `.agents/` pair it bootstraps; the script is
the authority and a divergence is a bug in the Go copy. It needs no payload and joins the
`needsPayload` exemptions beside `land` — a machine with nothing installed is exactly the
machine a fresh project is bootstrapped on. Scope flags are stripped by `scopeFlags` before
dispatch and mean nothing to it, the same as `metrics`; `loop` alone reads the tool flag.

**The run log is written only by `loop`, and only outside the project tree.** Its path is
`<git-common-dir>/libretto/loop.jsonl`, located with one `git rev-parse
--path-format=absolute --git-common-dir`; no `.gitignore` is written into a project the
tool does not own. The path is resolved behind an injected function, so the suite hands
in a `t.TempDir()` path and never touches this repository's `.git` — the one exception,
`TestTheRunLogLeavesGitStatusClean`, builds its own temporary repository and is skipped by
`make test-short`. Records are appended with one `Write` on an `O_APPEND` descriptor, and
two loops in two worktrees writing at once is accepted untested. `metrics` opens the log
for reading and never creates it. Stdlib only, and `loop` keeps needing no payload.

**No test starts a real tool.** Every session is a `loopRunner` the test injects; argument
lists, stdin and parsing are pure functions tested directly, so no test reaches
`~/.config/opencode` or runs `opencode` or `claude`. `claudeArgs`, `parseSession` and
`loopPrompt`'s text do not change in behaviour for a second tool — the tests that pin
them are the proof the claude path did not move. `TestParseLoopArgs` compares `loopOpts`
with `!=`, so every field it gains stays comparable. An OpenCode event line has no length
bound — a `tool_use` event carries the tool's whole output — so its parser walks the
output with a `json.Decoder`, never a line scanner with a fixed buffer. The vendor facts
the loop rests on were read off **opencode 1.17.15** — its `--help`, and `sst/opencode` at
tag `v1.17.15` (`packages/opencode/src/cli/cmd/run.ts`, `packages/schema/src/v1/session.ts`)
— and off **codex-cli 0.148.0-alpha.9**; a later release can move them, and the paths are
written so the next reader re-checks rather than re-guesses.

**`queue` is read-only and implements find-work's definition; it never extends it.** No
file written, no git mutation: its two git invocations are `rev-parse --show-toplevel` and
`branch --format=%(refname:short) --no-merged HEAD`, both through the `gitRunner` seam
`metrics`, `land` and `init` use, so tests drive it with fixed output, point at
`t.TempDir()`, and never read this repository's changes. Stdlib only. It reads
`.agents/changes/` at the repository root — the cwd only when git cannot answer — and joins
the `needsPayload` exemptions beside `init`. A line counts only when it starts with the
`Queued:` prefix at column 0, so a quoted mention inside an ask never queues a change. The
wiki reads the same scan, `queuedIdeas`, through the same filter, so the binary holds one
answer to "what is queued". **Its ceiling, marked `ponytail:` in `readInFlight`:** the
wiki asks git for branches once per `readInFlight` call and a render calls it three times;
the day a repository's branch count makes that slow, the render passes the branches down.

**`--json` is stdlib `encoding/json` and changes stdout only.** No dependency, as every
other JSON reader in this binary. Text output is byte-identical with or without the
gathered-report refactor behind it; exit codes, stderr lines and the release notice are
what they are without the flag; `main` still prints the error. `status` and `doctor`
print to `os.Stdout`, `metrics` and `queue` to an injected `io.Writer`, and the JSON path
follows whichever the command already uses.

**`doctor` prints nothing until its release check has answered.** It gathers the scan,
the release check and the prerequisites before rendering either output, so the problems
list that used to appear at once now waits on the live check — up to the five-second
ceiling on a slow or unreachable remote. Named as a cost, not hidden: one gathered report
read by both renderers is what keeps the text and the JSON from disagreeing, and printing
the problems first would need the text path to keep its own order of work.

**No test reaches the network for a release.** `doctor`'s check is chosen through the
package variable `releaseSource` — the `noticeSource` pattern — whose default asks a
checkout's remote or the module proxy, and `TestMain` stubs it for the whole package; a
test wanting an answer sets its own stub.

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

- **`metrics` reads churn with one `git log` call per file name, never one pathspec
  covering several.** The checklist was renamed from `plan.md` to `tasks.md` on 2026-08-17,
  and this report is retroactive over changes whose folders are already deleted — so the
  legacy path can never go. One pathspec naming both files is shorter and wrong: `plan.md` is
  the technical approach now, a prose document allowed to show a checkbox in an example,
  and counting its diffs would inflate the churn of every change from here on. Retired
  out of `add-design-phase`'s plan, which `payload` could not hold: `cmd/libretto/**` is
  governed here. The calls became three on 2026-10-21, `proposal.md` first.
- **One change document, both layouts read.** User, 2026-10-21, from
  `adopt-organic-driven-development`: a new change is born as one `proposal.md` — the
  capture grown in place to hold the contract, the approach and, last, a `## Checklist`
  section — and `loop`, `metrics`, `wiki` and `land` keep reading the old layout
  (`tasks.md`, `plan.md`, `spec.md`) for what already exists. Nothing in flight is
  migrated, and no reader rewrites a file. A change amending several capabilities carries
  the primary one's contract in `proposal.md` and one `spec-<capability>.md` per extra
  capability, each with its own `Targets:`.
- **The checklist is the boxes under a level-2 `## Checklist`, and nothing else.** Derived,
  2026-10-21: spec-drift's parsers open and close a section on `^##\s` only, and a box
  anywhere else in the document would be read as a criterion by `--ears`. So the CLI uses
  spec-drift's boundaries exactly — fences dropped, `###` and `#` closing nothing.
- **`proposal.md` is the checklist when, and only when, its `## Checklist` section holds a
  box.** Retired from `adopt-organic-driven-development`, 2026-10-21. *Existence selects it*
  lost: every change in both layouts has a `proposal.md`, so that rule would hide every old
  change's `tasks.md`/`plan.md` behind a prose capture and blank its loop, its strip entry
  and its churn. The 2026-08-17 precedent's *newer wins* survives; its *exists* test cannot.
- **The `proposal.md` churn query carries whole-file context.** Retired from
  `adopt-organic-driven-development`, 2026-10-21. *Querying with git's default 3-line diff*
  lost: a hunk does not carry its section heading, so whether a box sits under
  `## Checklist` cannot be decided from it. **Named ceiling, accepted:** old-layout churn
  still reads 3-line hunks with no fence state, so a fenced box in an old `tasks.md`/`plan.md`
  history counts there. Widening those queries would change what reaches git for every old
  change, and the real history holds no fenced box.
- **`land` gained no discovery code for the change document.** Retired from
  `adopt-organic-driven-development`, 2026-10-21. `headTargets` already read `Targets:` from
  every `.md` the folder holds at `HEAD`, fences stripped, so `proposal.md` and
  `spec-<capability>.md` were covered before anyone asked; *new `Targets:` discovery in
  `land`* lost to three tests proving what was already true.
- **`Targets:` outside fences is the one contract discriminator**, shared with
  `spec-drift --retired`. Derived, 2026-10-21: a second arm (`## Verification criteria`)
  was considered and dropped because it let `spec-drift` and `land` give opposite verdicts
  on one commit.
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
- **The flow board's cost table is deliberately not rendered.** Assumed
  2026-08-19 under attacca: `libretto metrics` answers it at the CLI, and the
  wiki version means reusing that git walker — it returns as its own change
  when somebody wants the numbers on the page. Ponytail-debt counts share the
  condition.
- **The palette is substring search over one JSON block; the card filter keeps
  its own gesture.** Assumed 2026-08-19 under attacca: no fuzzy ranking until
  substring demonstrably fails a real search, and the palette opens by shortcut
  or its own button — the reviewer caught a click wiring that collapsed the two
  gestures, and the separate button is the recorded repair.
- **Sparkline buckets anchor to the last commit's month, never to today.**
  A now-anchored histogram changes bytes at each month boundary with unchanged
  history. 2026-08-19.
- **The wiki's numbers are advisory mirrors; spec-drift and the loop stay the
  authorities.** Assumed 2026-08-19 under attacca: the Go checks reimplement the
  minimal definitions the wiki displays (criterion = bullet with `Proof:`
  beneath; EARS = `shall` with emphasis stripped; proof resolves = file exists +
  Go test declared), and a divergence from the script is a wiki bug, never a
  second truth. Absolute dates everywhere — a relative date is a function of the
  clock and byte-determinism forbids one. Extraction lands with its consumer
  (vertical cut, four changes not five).
- **The wiki viewer is home-and-pages: single file, pages by hash.** Asked
  2026-08-18; multi-file lost on N+1 generated files to own, mark and clean, and
  stays the recorded return path if hashes ever fall short (crawlers, per-page
  HTTP). Key info on a card is the last-changed date (git, `%as`, via an
  injectable seam — mtime lost because it does not survive a clone) and the
  criteria count with a mini-bar; commit-count was offered and not chosen.
  Governs lives on the page, not the card. This superseded the bento hero chart
  and the scroll-spy the same day they landed, unreleased — the cards absorbed
  the chart's information and the home absorbed the spy's job.
- **The wiki viewer is bento-and-data, dark-first, its palette measured.** Chosen
  2026-08-18 after web research, over editorial-minimal and glass-total; glass
  confined to the sidebar, which is what the 2026 retrospectives themselves
  recommend. The chart is criteria-per-capability bars because it answers the
  wiki's own question from data already extracted; richer charts wait for a
  question they answer. Chart libraries lost to rectangles; canvas lost to
  markup that prints and works with JS off; scroll listeners lost to
  animation-timeline.
- **The specs wiki is generated by this binary, in markdown, into the specs directory
  itself.** Asked and answered 2026-08-18. A skill-shipped bash generator lost on
  testability (comes back if a host without the binary needs the wiki); an install-time
  node script lost on timing — install happens once, specs move on every landing — before
  it lost on stack; HTML and a serve mode lost to a forge that already renders markdown
  in-tree; a hand-maintained index lost to this repository's own "ten over eleven
  directories" history. A goldmark dependency lost to twenty lines of `bufio.Scanner`:
  the extraction is line-shaped, and what counts as a criterion is `spec-drift`'s own
  definition — a bullet with a `Proof:` beneath it — because two tools with two
  definitions would disagree about what a capability promises.
- **A generated index is not a second source of truth.** The one-place rule forbids a
  second *hand-maintained* copy; the wiki is a marked, regenerated view of the specs
  themselves — a cache, not a source. The drift risk is why regeneration rides
  `record-work`'s landing step (see `payload`).
- **The HTML viewer is a static file, never a server.** The ask was a page reached by
  a link, and `file://` reaches it; `--serve` comes back only if someone needs a URL
  other people open. Its template ships inside the binary as a const — a template
  file would be payload the linker has to place, for exactly one consumer.
- **The panel's wiki row is project scope only, by instruction.** The user drew the
  line 2026-08-18: under global it never appears — `~/.claude` has no
  specification. And no specs, no row (the `models` precedent), rather than a
  disabled row with an explanation; that alternative comes back if absence
  confuses. The row runs `wiki --open` — asked for 2026-08-18: a menu press ends
  with the viewer on screen. It still dispatches the same subcommand; the plain
  run stays what landings use.
- **The viewer opens via `file://`, never a localhost server.** The page is
  self-contained by contract, so the browser tab is identical and there is no
  daemon lifetime to manage inside a panel action. Bring-back condition for
  `--serve`: a URL someone else must reach, or a browser policy blocking
  `file://`. The opener is the platform's own, detached with `Start` so a hung
  handler cannot block the panel. 2026-08-18.
- **One plain `wiki` run refreshes every marked view present.** The alternative was
  teaching `record-work` to name `--html` at landings — a second payload sentence
  that drifts from the set of outputs. The binary refreshing what it owns keeps that
  instruction closed under new views, and it is why a foreign `wiki.html` is skipped
  silently on a plain run: erroring there would block every landing regeneration.
  2026-08-18.
- **`land` is verify-only — a read-only gate over the staged index, never a performer.**
  Applying a delta is a semantic merge a binary cannot do, and the observed failures
  were *incomplete* landings, which a verifier catches. If wrong: the command grows a
  write mode later; verify-only is forward-compatible. Assumed 2026-08-21, attacca run.
- **`land` checks parts 2 and 4 of the landing; part 3 stays owned by `spec-drift
  --retired`.** Every file of the change folder deleted with no partial deletion, and
  each delta's `Targets:` capability spec modified in the same staged diff. Duplicating
  the retirement check in the binary is two sources of truth about what "retired"
  means. If wrong: the check list grows, the ownership split does not. Assumed
  2026-08-21, attacca run.
- **The staged index, pre-commit — before the mistake is history.** Change name
  optional: with none given, `land` infers from staged change-folder deletions. If
  wrong: a `--commit <sha>` mode can be added without moving the default. Assumed
  2026-08-21, attacca run.
- **A stale wiki does not block `land` — warn only**, the same rule as record-work's
  "a missing convenience never blocks a landing". Assumed 2026-08-21, attacca run.
- **The payload learns about the command, minimally** — one guarded clause in
  `record-work`, the wiki clause's shape, owned by the `payload` capability. A verifier
  nothing invokes verifies nothing. Recorded here only so this spec is not read as the
  whole change; if wrong, the clause is one sentence to remove. Assumed 2026-08-21,
  attacca run.
- **`land` is one self-contained file, `cmd/libretto/land.go`.** An `internal/land`
  package lost on one consumer — a package boundary for a single command is a premature
  abstraction, and it lifts into `internal/` the day a second consumer exists. A payload
  bash script beside `spec-drift` lost because the command must run on machines that
  installed the binary and nothing else — exactly the machine record-work invokes it on
  — and payload never implements delivery. go-git lost to the ladder — stdlib first —
  and to the exec'd-git seam already here (`metrics`, `wiki`): a library's index
  semantics are a second implementation of git to trust. 2026-08-21.
- **`landChangeRoots` duplicates spec-drift's `CHANGE_ROOTS`, and the script is the
  authority.** Three strings in Go mirroring three in bash; a shared source would mean
  the binary parsing the script or the script reading Go, and both cost more than three
  strings. Ceiling: unify — or add a cross-checking test that greps the script — the
  day the two lists diverge. 2026-08-21.
- **A rename out of the change folder counts its source as removed** — the contract is
  that nothing under the folder survives the commit, and a rename-out leaves the folder
  empty. The diff base is ordinary `HEAD`; in-progress merges are out of scope. One
  exit code covers every failure kind — distinct codes wait for a caller that needs
  them. All three assumed 2026-08-21, attacca run.
- **`init` creates the `.agents/` layout and nothing else: two keep-files and the lessons
  ledger.** Retired from `add-init-command`, assumed 2026-08-21. Guided onboarding content —
  a sample spec, `docs/SPEC.md`, `CLAUDE.md` — lost to ponytail: three artifacts are the
  whole contract, and the list grows one file at a time when someone asks.
- **A `.gitkeep`, never a placeholder `README.md`, carries `changes/` through git.** Retired
  from `add-init-command`'s plan, 2026-08-21. A README must say something, be maintained,
  and invites into `changes/` the marker-ownership ambiguity `wiki` resolves for
  `.agents/specs/README.md`; a zero-byte keep-file is universally read as "carry this
  directory" and can never drift. The hole it closes is real: git drops an empty
  `changes/` on every fresh checkout, which is how a queued idea was lost. This repository
  carries its own `.agents/changes/.gitkeep` for the same reason — deleting it reopens the
  hole.
- **A foreign consolidation root refuses, and the refusal wins over a partial `.agents/`.**
  Assumed 2026-08-21, `add-init-command`: a project already consolidating under `specs/`,
  `openspec/`, `docs/specs/`, `spec/`, `changes/` or `openspec/changes/` has a convention,
  and inventing a second one is two sources of truth — which is also why filling a partial
  `.agents/` beside a foreign root is refused rather than done. With no foreign root, a
  partial `.agents/` gets its missing pieces. The probe runs before any stat-and-create, so
  the precedence is ordering, not a flag. If wrong: a flag to adopt another layout is a
  later change.
- **`init` is idempotent and refuses outside a git repository.** Assumed 2026-08-21,
  `add-init-command`: a second run creates nothing and exits zero, and a directory nothing
  tracks has nothing for a keep-file to survive. If wrong: a warn-and-continue mode is one
  branch.
- **`init` is one self-contained file, `cmd/libretto/init.go`.** Retired from
  `add-init-command`, 2026-08-21: an `internal/bootstrap` package lost on one consumer, the
  `land` precedent; it lifts the day a second consumer exists.
- **`init`'s foreign-root lists duplicate spec-drift's, and the script is the authority.**
  Retired from `add-init-command`'s plan, 2026-08-21: two short lists that change roughly
  never, against a shared source that would mean the binary parsing the script. The same
  ceiling as `landChangeRoots`: unify, or add a cross-checking test, the day they diverge.
- **A symlink at an artifact path is present, dangling or not.** Found by the 6→7 review,
  2026-10-21: `os.Stat` followed a dangling link and the write landed outside `.agents/`.
  The check is `Lstat`, and a link is never written through.
- **The loop's budget is the host's `--max-budget-usd`, per session, which is per round —
  passed through under the host's own name, as typed.** Retired from
  `add-loop-budget-and-run-log`, assumed 2026-10-22, attacca run. The host enforces it and
  the record carries what was reported, so an overshoot is visible rather than hidden. The
  value must match `^[0-9]+(\.[0-9]+)?$` and exceed zero: `strconv.ParseFloat` alone lost
  because it accepts `1e2`, `0x1p-2`, `+3` and `inf`, which would reach the host with
  meanings nobody checked. Tokens as the unit lost to a cap the host enforces natively. If
  wrong (a loop-wide cap was meant): sum the log's `cost_usd` and stop at a total.
- **No loop-wide cap; the header prints `--max` × budget to the cent.** Retired from
  `add-loop-budget-and-run-log`, assumed 2026-10-22. **Ceiling:** adequate while somebody
  reads the header; the day a run is left unwatched with a large `--max`, a total cap is
  the replacement.
- **Nothing the host reports is interpreted, and a budget-cut session gets no special
  case.** Retired from `add-loop-budget-and-run-log`, assumed 2026-10-22. `subtype` and
  `is_error` are recorded verbatim and `error` carries the exec error's message — two
  fields, because one field holding both merges two sources with different meanings, which
  is interpreting them. The budget-hit spelling was never observed; if the first hit should
  end a run, it becomes one stop reason keyed on the recorded subtype, once observed.
- **The run log lives at `<git-common-dir>/libretto/loop.jsonl`, one JSON-lines file for
  every change.** Retired from `add-loop-budget-and-run-log`, assumed 2026-10-22. Never
  tracked, never in `git status`, no `.gitignore` written, shared by every worktree — which
  is why its messages say *repository*, never *checkout*. What lost: a self-ignoring
  `.agents/loop/`, because it is per worktree and the log would die with the worktree the
  loop ran in; a log inside the change folder, which the session commits — a race — and
  the landing deletes; and `~/.claude/projects/`, the host's directory the token block
  promises never to write. **Ceiling:** it is per machine — a fresh clone starts with no
  history, like the host's transcripts. If that stops being acceptable, it moves under
  `.agents/`.
- **The run log is the one recorded artifact, and its writer is the binary.** Retired from
  `add-loop-budget-and-run-log`, 2026-10-22, amending *nothing is instrumented*: the rule
  rejected artifacts eight skills would have to remember to write, and the loop is one
  writer that runs every round by construction and cannot forget. A phase writing to it,
  or anything landing in a change folder, still breaks the rule.
- **Loop cost is the host's `total_cost_usd`, recorded verbatim, absent when unreported,
  never zero.** Retired from `add-loop-budget-and-run-log`, assumed 2026-10-22. It extends
  *no pricing* rather than relitigating it — the binary still prices nothing, and the
  observed run labels the figure `"costBasis": "list"`, so it is shown as the host's
  estimate. Joining the host transcripts by `session_id` instead lost: transcripts carry
  tokens, never dollars, and give no budget at all.
- **`closed` is signed, and `metrics` reports it as a net figure.** Retired from
  `add-loop-budget-and-run-log`, assumed 2026-10-22. If closes and reopens should be
  counted apart, the record gains a second field, the way the change table already splits
  them.
- **A run log that cannot be located or written warns once and never stops the loop; a
  round whose checklist vanished is still recorded, as `stopped: "error"`.** Retired from
  `add-loop-budget-and-run-log`, assumed 2026-10-22. If the user would rather not spend
  without a record, it becomes a refusal before the first session.
- **`loopWith` and `metricsWith` exist so the run log path is injected without moving a
  call site.** Retired from `add-loop-budget-and-run-log`, 2026-10-22: `loop` and `metrics`
  keep their signatures and delegate, and `main.go` calls the `With` forms. A resolver in
  `loopOpts` would have made it uncomparable for `TestParseLoopArgs`, and a new parameter
  would have edited four `loop` and twenty-eight `metrics` call sites.
- **The loop's tool is the existing tool axis, never a `--tool` flag.** Retired from
  `support-codex-opencode-in-loop`, assumed 2026-10-22, attacca run. `scopeFlags` already
  stripped `--claude`, `--codex` and `--opencode` from every subcommand's arguments, so
  before this `loop c --opencode` silently drove `claude` — the flag was accepted and
  thrown away. Reading the axis fixes that rather than adding a second spelling; `--tool`
  lost because `loop c --codex --tool opencode` is two answers to one question, the error
  the axis rule exists to report. If wrong (`--tool` was wanted): it becomes an alias that
  errors when it disagrees with the axis, and the axis still decides.
- **Codex is refused by the loop, on evidence.** Retired from
  `support-codex-opencode-in-loop`, 2026-10-22 — the scope assumed, the facts measured.
  Codex has no headless route to a command: `codex exec --help` (codex-cli
  0.148.0-alpha.9) offers no command or slash-command option, its top-level list has no
  prompts or commands subcommand, and `codex debug prompt-input`, with `CODEX_HOME` in a
  scratch directory, rendered `/libretto-attacca` as literal user text while it listed the
  project's `.agents/skills/` by name. And agents cannot be installed into Codex —
  `targets` gives it skills only, its custom prompts are deprecated — so the flow's fresh
  reviewers have nothing to launch. A codex prompt naming `$find-work` and
  `$build-and-check` lost: it restates attacca's rules inside the loop, which `loopPrompt`
  exists never to do. The refusal is code, not an absence — left out of the switch,
  `--codex` would fall through to `claude`, the silent mis-drive above. **What brings it
  back:** the payload shipping the attacca route as a skill Codex lists, or Codex
  installing commands; then the refusal becomes a runner, and this entry says which.
- **OpenCode runs `opencode run --command libretto-attacca --format json`, with the loop's
  task text on stdin.** Retired from `support-codex-opencode-in-loop`, measured
  2026-10-22. `run` wraps a positional argument holding a space in double quotes and
  escapes the quotes inside it, and `resolveRunInput` returns piped stdin unchanged when no
  message was given — so stdin is the one route on which OpenCode sees the text `claude`
  sees. Whitespace-split positional words lost because every newline collapses to a space;
  `--format default` lost because it carries no session id to match a round to OpenCode's
  own record of it.
- **`--max-budget-usd` is refused under `--opencode`.** Retired from
  `support-codex-opencode-in-loop`, assumed 2026-10-22; the absence measured.
  `opencode run --help` lists no option naming a budget, a cost or a spend limit, and the
  loop prices nothing itself, so there is nothing to pass the amount to. If wrong: the
  loop-wide cap named above as the replacement — which for opencode would first have to
  decide its cost.
- **An opencode round records no cost, though OpenCode reports one per step.** Retired
  from `support-codex-opencode-in-loop`, assumed 2026-10-22. `StepFinishPart` carries
  `cost` and `tokens` (`packages/schema/src/v1/session.ts`), but `run` emits only its own
  session's parts — `if (part.sessionID !== sessionID) continue` — so a subagent's
  `step_finish`, the 6→7 reviewer's among them, never reaches the stream. A sum would be a
  figure built here and short by exactly that spend, and *the host's figure or absent*
  allows neither. `duration_ms`, `subtype` and `is_error` are omitted for the same reason.
  Reading the cost back with `opencode export <id>` lost: a second call per round, of an
  unobserved shape, for a number the loop may not construct. **Ceiling:** the day `run`
  emits a session total, or the children's steps, it is recorded verbatim as `cost_usd`;
  a partial figure, if ever wanted, goes in a field that is not `cost_usd`, so `metrics`
  never adds it to the host totals.
- **No permission flag for either tool.** Retired from `support-codex-opencode-in-loop`,
  assumed 2026-10-22: parity with `claude`, to which the loop passes none. Without
  `--auto`, `opencode run` answers every `permission.asked` with `reject` and denies
  `question`, `plan_enter` and `plan_exit`; a round so blocked closes no box, and stuck
  detection stops the loop after two. The fix is the user's opencode configuration, never
  a flag here.
- **Every record carries `tool`; a record without it means `claude`, and so does a header
  without ` · via`.** Retired from `support-codex-opencode-in-loop`, assumed 2026-10-22.
  One log is shared by every worktree and now by two tools, and an unpriced claude round
  has to stay distinguishable from an opencode round unpriced by design. Older records are
  not rewritten. If wrong: the field goes, and the two kinds of unpriced round become
  indistinguishable in `metrics`.
- **The tool's session is chosen through one pure seam, `sessionFor(tool)`.** Retired
  from `support-codex-opencode-in-loop`, assumed 2026-10-22 after review-spec. It is the
  only place an executable name is paired with its argument builder, stdin and parser, and
  `loopWith` builds its runner from it alone, so the pairing is tested rather than trusted:
  a criterion on `claudeArgs` could not see which executable ran. **Ceiling:** that the
  runner sets the session's `Stdin` from `opencodeStdin` is observed by no test, since
  none starts `opencode`; an injected command factory whose `exec.Cmd` a test reads back
  is the replacement if it ever bites.
- **A host `error` event is printed as `opencode error: <message>` after the text, and
  recorded in `error` only when the exec error is absent.** Retired from
  `support-codex-opencode-in-loop`, assumed 2026-10-22 after review-spec. The message is
  `error.data.message` when present and `error.name` otherwise, several joined with a
  newline — the rule `run.ts` itself applies to print them, so nothing is interpreted
  here. On a non-zero exit `error` keeps the exec error, as it does for `claude`. If
  wrong: the record keeps the exec error alone, and an event on a zero exit is visible on
  the terminal only.
- **What OpenCode was observed to do, 2026-10-22.** Retired from
  `support-codex-opencode-in-loop`, so the next reader re-checks rather than re-guesses.
  `opencode run --help` (1.17.15): `--command` "the command to run, use message for args",
  `--format` "default (formatted) or json (raw JSON events)", `--auto` "auto-approve
  permissions that are not explicitly denied (dangerous!)", no budget. From `run.ts` at tag
  `v1.17.15`: one event per line as `{ type, timestamp, sessionID, ...data }`, the types
  `tool_use`, `step_start`, `step_finish`, `text`, `reasoning` and `error`; `--command`
  calls `session.command` with the message as its arguments; exit 1 on argument errors,
  `exitCode = 1` when the command call errors, a `session.error` arrives or the event loop
  throws, otherwise 0 once the session goes idle. Isolated under scratch `XDG_*`
  directories, `opencode debug config` resolved a symlinked
  `.opencode/commands/libretto-attacca.md` as the command `libretto-attacca`, `$ARGUMENTS`
  intact. **Never observed:** a real OpenCode event stream end to end, and codex's exit
  codes and usage events — no billed session was run for either tool.
- **The binary may read the queue; find-work still defines it.** Retired from
  `add-queue-command`, 2026-10-22, attacca run. It reverses `payload`'s "the Go binary
  knowing about the queue", which the wiki had already broken without saying so. The
  binary implementing a rule the skill owns is delivery; the binary deciding what counts as
  queued would not be.
- **Queued is three conditions, applied in one place, `queuedIdeas`, for both the wiki and
  `queue`.** Retired from `add-queue-command`, assumed 2026-10-22 on review-spec's finding:
  a `Queued:` line at column 0, no unmerged local branch taking the name, and no box under
  `## Checklist` — a box means the work started. find-work states the same three. *B —
  `queue` filtering on top of the wiki's own scan, the wiki untouched* lost: two
  definitions of "queued" in one binary, and the wiki kept counting dispatched changes. *C
  — `queue` running find-work's `rg` line* lost: it puts `rg` on the binary's runtime path,
  and the filter would still be written in Go. If the checklist clause is wrong, it leaves
  all three places together; if the wiki should not filter, it keeps an unfiltered count
  and the binary holds two definitions.
- **The branch filter is find-work's, applied by the binary**: local branches unmerged
  into `HEAD`, one leading `feat/` `fix/` `docs/` `chore/` `refactor/` removed, the name
  matched whole. Retired from `add-queue-command`, assumed 2026-10-22. **Ceiling:** `HEAD`
  stands in for the base branch; run from a branch behind the base, a landed change can be
  dropped as dispatched rather than shown — which drops, never invents. *Resolving the base
  branch* lost: one more lookup with its own failure modes, for a case `HEAD` gets right
  whenever the command runs on the base. The day that misleads, a base is resolved from
  `refs/remotes/origin/HEAD` with `main` as fallback, in one function. **Local branches
  only**: a change dispatched from another machine and not fetched still shows as queued,
  exactly as in find-work; if wrong, both add remote branches together.
- **Ties on one `Queued:` date go by change name, ascending.** Retired from
  `add-queue-command`, assumed 2026-10-22 — already the wiki's order. *Capture-commit
  order* lost: a `git log` per proposal, against a design that put the date in the file to
  avoid exactly that. If wrong, the skill and the binary move together.
- **`queue` reads the repository root, the cwd only when git fails.** Retired from
  `add-queue-command`, assumed 2026-10-22. *The cwd, as the wiki reads it* lost: it is the
  bug `metrics` recorded, a subdirectory run reporting an empty queue.
- **Unreadable branches list the queue unfiltered, with one stderr line, exit 0.** Retired
  from `add-queue-command`, assumed 2026-10-22 — find-work's fail-safe direction: keep the
  idea.
- **A dropped change is named on a closing line with its branch, on stdout; only the
  unreadable-branches line is stderr.** Retired from `add-queue-command`, assumed
  2026-10-22: an idea that vanishes from a list without a word reads as lost, and the
  closing line is information, not an error.
- **The empty states are told apart, and measured after the filter.** Retired from
  `add-queue-command`, assumed 2026-10-22. No `.agents/changes/` prints `no .agents/changes/
  here — nothing is queued`, which points at `init`; nothing left after the filter prints
  `nothing is queued`, then the dispatched line when anything was dropped. Exit 0 either way.
- **`queue` takes no argument but `--json` and is exempt from the missing-payload stop.**
  Retired from `add-queue-command`, assumed 2026-10-22: it reads the project, never the
  payload. `--json` was admitted by `add-json-output` the same day, by amending the refusal
  criterion in place rather than writing a second one beside it — two criteria would have
  promised both that `--json` is refused and that it is read.
- **One gathered report per command, rendered twice.** Retired from `add-json-output`,
  assumed 2026-10-22: each of `status`, `doctor`, `metrics` and `queue` derives its facts
  once into a report struct, and the text renderer and the JSON encoder both read it. *B —
  a second function per command deriving the JSON from the same inputs* lost: a second
  derivation is a second source of truth, and the first fix to one side drifts the other
  with nothing to notice. *C — `json` tags on the internal types* (`link.Entry`,
  `changeMetrics`, `loopTally`, `queuedIdea`) lost: their fields are unexported, they carry
  bytes and paths no consumer needs, and every refactor of an internal type would become a
  contract change. So the JSON structs are separate, small, exported-field types mapped from
  the report — the mapping is the contract's boundary — and **no internal type carries a
  `json` tag.** **Ceiling:** a text renderer may still compute a display value from the
  report — a span, a percentage — but never a fact the report does not hold. A block added
  to the text later (the queued drift trend for `metrics`) lands in the report and gets its
  JSON key in the same change; whether it was mapped is review's to see, not a gate's.
- **Every JSON document carries `"schema": 1`.** Retired from `add-json-output`, assumed
  2026-10-22: while the tool is `0.x` the binary's version moves on every merge and says
  nothing about the shape, so a consumer gets one integer to check. An added key never bumps
  it; a removed, renamed or retyped key does, and is a promise removed — `!` on the commit.
  *Keying on `libretto version`* lost for that reason.
- **An unknown value is an absent key; a known zero is present; a known-empty list is
  `[]`.** Retired from `add-json-output`, assumed 2026-10-22, extending the run log's *cost
  absent when unreported, never zero*. *`null` for unknowns* lost: two conventions in one
  binary is one too many. In the types: a key that may be unknown is a pointer or slice with
  `omitempty`; one always known is a plain value without it — `omitempty` on a plain `int`
  turns a known `0` into an absent key, which is the bug the rule exists to catch.
- **`--json` changes stdout only, and errors stay plain lines on stderr.** Retired from
  `add-json-output`, assumed 2026-10-22: the exit code already carries the verdict, and a
  second error format would have to be kept in step with `main`'s. A report that never
  completed writes nothing to stdout. The one error inside a report is `metrics`' unreadable
  run log, which exits 0 and becomes `loop.error`.
- **The JSON carries the report's data, never its prose, and nothing derivable from the
  rest; keys are snake_case, indented two spaces, HTML escaping off.** Retired from
  `add-json-output`, assumed 2026-10-22: the shapes start small because adding a key is free
  and removing one is a broken promise; snake_case is the run log's convention; whitespace
  is not promised; a reading must arrive as written, not as `<`.
- **`--json` is a per-command flag after the command name**, so `libretto --json status`
  stays an unknown command. Retired from `add-json-output`, assumed 2026-10-22: only four
  commands honour it, and a global flag would need every other command to refuse it.
- **`status` and `doctor` read `--json` and ignore every other argument, as they always
  have.** Retired from `add-json-output`, assumed 2026-10-22 after review-spec, **reversing
  the first session's answer** that both would refuse any other argument: refusing would
  move an exit code on a path unrelated to JSON, and "exit codes unchanged on every path"
  then holds with no exception. If a typo printing text into a parser bites somebody, a
  refusal lands as its own change.
- **`doctor`'s release check sits behind the `releaseSource` seam, and the suite stubs
  it.** Retired from `add-json-output`, assumed 2026-10-22 on the `noticeSource` precedent:
  `doctor` called `dist.Latest` with no seam, so
  `TestDoctorExitsNonZeroWhenSomethingNeedsAttention` reached the module proxy from
  `go test` (measured). The check now returns facts — running, latest, error — and the text
  line and the JSON `release` object are both rendered from one `releaseCheck`, observed
  through one stub in `TestReleaseLineIsRenderedFromTheCheck`.
- **Spec drift is not measured in `metrics`.** Retired from `report-drift-trend-in-metrics`, assumed 2026-10-22; decided as
  option C by the orchestrator under the user's standing authorisation, on a measurement of
  this repository at `dd6683d` (565 commits, 498 non-merge, 70 first-parent; 14 capability
  specs, all with `Governs:`).
  **The definition prototyped:** a *drift event* is a commit touching a capability's
  `spec.md` after one or more earlier commits touched paths its `Governs:` globs claim —
  read as `spec-drift`'s `governs_of` reads them and matched as its `matches_any` matches,
  where `*` crosses `/` — without touching that spec; paths under the specs directory and
  `.agents/changes/` are never code, and a commit touching a change folder whose `Targets:`
  names the capability counts as the spec moving. **The numbers:** at commit granularity,
  47 events, **37 of whose catch-up commits were landing commits** — the flow applying a
  delta after its code, by design. At landing granularity (a first-parent commit, a merge
  read against its first parent, counting from the first landing that touches the spec),
  **13 events, all read, 1 related** to the code before it (`readme`, PR #72 after #63).
  The fallback — `spec-drift`'s default-mode warning replayed per landing — flagged **27 of
  120** capability-landings (`cli` 10, `readme` 12, `review-project` 2, `payload` 2,
  `link-state` 1); of **10 hand-checked, 0** changed a promise without its spec: **6 were
  `Governs:` overlaps** (a path whose semantic owner is another capability whose spec moved
  in the same landing), 3 were edits inside an existing contract, 1 arguable.
  **Why nothing ships:** a metric that is mostly noise teaches people to ignore it, and the
  replay cost 25 tests and a Go mirror of bash globbing to report it. **B** — a `--history`
  mode in `spec-drift` exec'd by `metrics` — and **A** — a Go replay inside `metrics` with
  a `spec_drift` JSON key — both lost to **C**, ship nothing. **What would reopen it:** a
  corpus where the replay's hand-checked precision is high — most warned landings turning
  out to be promises moved without their spec — or a way to tie a spec edit to the code it
  answers (a commit trailer, a delta naming its landing), which the flow does not write
  today.

- **The wiki's title names the repository, not the run directory.** Retired from
  `name-the-wiki-after-the-repository`, fixed 2026-10-23: `filepath.Base(projectDir)` wrote a
  linked worktree's directory name into the committed view seven landings in a row, each
  worked around with a temporary symlink. The name is the parent of git's common dir, shared
  by every worktree; only a common dir called `.git` names it, so a submodule
  (`.git/modules/<name>`) or a separate git dir keeps the run directory's name rather than
  naming `modules`. Lost: the origin remote's name (a repository need not have one, and a fork
  renames it) and `--show-toplevel` (the worktree's own root — the bug itself).

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
- [x] 5.2 `--json` for `status`, `doctor`, `metrics` and `queue`
- [ ] 5.3 the panel path under a real TTY
- [ ] `doctor`: target directory missing or unwritable
- [ ] `doctor`: repo state — uncommitted changes, behind the remote *(the release half of
      this is done; the working-tree half is not)*
- [ ] `install`: present the plan before applying, as `prune` already does
- [ ] 7.3 delete `install.sh`, once `libretto install` is verified against a real
      `~/.claude` with a throwaway item
- [x] `wiki` — the project's specs rendered into one marked README
- [x] `wiki --html` — the same specs as one self-contained, marker-owned viewer
- [x] the `wiki` panel row — project scope only, absent without specs
- [x] `init` — the `.agents/` layout bootstrapped at the repository root, never over anything
- [x] `loop --max-budget-usd` and the per-round run log, read back by `metrics`
- [x] `loop --opencode`, and `--codex` refused with the reason
- [x] `queue` — find-work's queue, read-only, and the wiki reading the same scan

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
- **`--codex` and `--opencode` resolve their targets, and any two destination flags
  together is an error**
  Proof: cmd/libretto/scope_test.go TestDestinationFlags
- **a codex install links skills only and leaves every other destination alone** —
  with a command present in the repo, so the absent `commands/` proves something
  Proof: cmd/libretto/scope_test.go TestInstallCodexLeavesOthersAlone
- **an opencode install links a command as a symlink into `<root>/commands`**, links
  the skill too, **writes the agent as a regular file with no `tools:` and a marker in
  it**, and leaves every other destination alone
  Proof: cmd/libretto/scope_test.go TestInstallOpencodeLeavesOthersAlone
- **`uninstall --opencode --yes` removes a generated agent and keeps a file somebody else
  put at the same path** — the promise applied to a kind that is written rather than
  linked, and it needed no code
  Proof: cmd/libretto/scope_test.go TestUninstallOpencodeRemovesGeneratedAgents

  **The symlink half is the load-bearing word.** "Links commands" is satisfied by a
  copy, and a copy stops tracking the repo the moment either side moves — the one
  outcome this promises not to produce. Four properties of one install run, joined on
  purpose because one test asserts all four; a criterion joined by *and* whose proof
  checks one half is the failure this shape usually hides. **Watched red before green**
  by reverting `Opencode.Kinds()` alone: the command assertion failed and nothing else
  did.
- the remembered destination round-trips the new words and falls back to global
  Proof: cmd/libretto/remembered_test.go TestRememberedDestinationRecognisesNewTargets
- help names the new flags and both env overrides
  Proof: cmd/libretto/main_test.go TestHelpNamesEveryDestination

  **Ceiling named:** it proves each flag appears, never that the row beside it says
  something true — `--opencode` described as skills-only would still satisfy it. The
  row's accuracy is held by reading, and the replacement the day that bites is a
  criterion about what a row must contain.

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
- **it owns no state.** The change's checklist is the state — the `## Checklist` section
  of `.agents/changes/<change>/proposal.md`, or `tasks.md` (`plan.md` before 2026-08-17) for
  a change in the old layout: the boxes the 5→6 seam writes and phase 6 marks — and the loop
  reads the boxes and nothing else. Only the box is parsed, never the task text: the flow
  owns the checklist's shape, and a runner that understood it would be a second opinion
  about what a task is. The section boundary is the one addition, and it decides which
  boxes are in scope, never what they say. The run log below is a record the loop writes
  and never reads back: no decision about the next round is taken from it.
  Proof: cmd/libretto/loop_test.go TestCountBoxesReadsOnlyTheBox
- **When counting a change document, the section counter shall count only `planLine` boxes
  that lie after a `## Checklist` line, before the next line matching `^##\s`, and outside
  fenced code blocks.** A box under `## Task breakdown`, a box in a fence inside
  `## Checklist`, and a box under a `### ` subheading inside `## Checklist` count as zero,
  zero and one respectively — spec-drift's boundaries exactly, so the two cannot disagree
  about what the checklist is.
  Proof: cmd/libretto/loop_test.go TestOnlyTheChecklistSectionCounts
- **Where a change document has no `## Checklist` heading, the section counter shall return
  zero open and zero done**, whatever boxes the file holds elsewhere.
  Proof: cmd/libretto/loop_test.go TestOnlyTheChecklistSectionCounts
- **When `countBoxes` reads a `tasks.md` or `plan.md`, it shall not count a box inside a
  fenced code block** — one box definition in every reader, `find-work` included.
  Proof: cmd/libretto/loop_test.go TestCountBoxesSkipsFences
- **Where a change's `proposal.md` holds at least one box under `## Checklist`,
  `checklistPath` shall return that `proposal.md`**, whether or not `tasks.md` or `plan.md`
  also exist in the folder. Existence is never the signal: every change in both layouts
  has a `proposal.md`.
  Proof: cmd/libretto/loop_test.go TestChecklistPathPrefersTheChangeDocument
- **Where a change's `proposal.md` holds no box under `## Checklist`, `checklistPath` shall
  return `tasks.md` when it exists and otherwise `plan.md` when it exists; where neither
  exists, it shall return `proposal.md`**, including when `proposal.md` holds boxes outside
  `## Checklist`. It reads a change created before either rename: the checklist carried the
  name `plan.md` until 2026-08-17 and `tasks.md` until 2026-10-21, and a rename that
  stranded work already in flight would have cost more than the name was worth. With no
  checklist anywhere the error names the file the flow now writes, so it teaches the layout
  the flow actually uses. The fallbacks carry a `ponytail:` saying when they go: once no
  `tasks.md` or `plan.md` under any change folder holds a box.
  Proof: cmd/libretto/loop_test.go TestChecklistPathFallsBackToPlan
- **If a change holds a `proposal.md` with no box under `## Checklist` and has neither
  `tasks.md` nor `plan.md`, then `loop` shall return an error containing `no task boxes` and
  the `proposal.md` path while PATH is empty, and shall start no session.** A queued capture
  is named as having no checklist rather than as missing `claude`; the empty PATH proves the
  refusal comes before the dependency check.
  Proof: cmd/libretto/loop_test.go TestLoopRefusesAQueuedProposal
- **While `proposal.md` is the selected checklist, the loop shall count progress from
  `## Checklist` boxes only, and the prompt printed by `--dry-run` and handed to each session
  shall name `.agents/changes/<change>/proposal.md` and `## Checklist`.** A box closed outside
  the section is not progress.
  Proof: cmd/libretto/loop_test.go TestTheLoopDrivesTheChangeDocumentsChecklist
- **The `no checklist at` and `no task boxes` refusals shall not contain `cut into tasks`.**
  No agent cuts the checklist any more; the session writes it in the 5→6 seam.
  Proof: cmd/libretto/loop_test.go TestLoopRefusalsDoNotNameTheCutter
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

#### Which tool runs a round

The tool is the axis every subcommand already parses — `--claude`, the default, or
`--opencode` — and `loop` is the one command that reads it after `scopeFlags` strips it.
There is no `--tool`. Each round hands the chosen tool the installed `/libretto-attacca`
command with the loop's task text as its arguments: `claude -p` takes the whole prompt;
`opencode run --command libretto-attacca --format json` takes the arguments on stdin,
because `run` re-quotes a positional argument holding a space and passes piped stdin
through verbatim. `--codex` is refused: Codex installs skills only, and every round is
routed to a command Codex cannot run. The executable is paired with its argument builder,
its stdin and its parser in one place, `sessionFor`, and the loop builds its runner from
that alone.

- **When `run` receives `loop <change> --opencode --dry-run`, the loop shall print a header
  naming `opencode`.** Proved through `run`, because `scopeFlags` strips the flag before
  dispatch and a test of `loopWith` alone would prove nothing about that.
  Proof: cmd/libretto/loop_test.go TestTheToolFlagReachesTheLoop
- **Where `loopOpts` carries a tool, the loop header shall end its first line with
  ` · via <tool>`.**
  Proof: cmd/libretto/loop_test.go TestTheHeaderNamesTheTool
- **Where `loopOpts` carries no tool, the loop header shall carry no ` · via` suffix** —
  the same reading the record gives an absent `tool`, claude, so a `runLoop` call without
  a tool prints the header it always printed.
  Proof: cmd/libretto/loop_test.go TestAnEmptyToolPrintsTheHeaderItAlwaysPrinted
- **`sessionFor("claude")` shall return the executable `claude` and an argument builder
  whose output equals `claudeArgs` for the same prompt and budget.** `loopWith` builds its
  runner from `sessionFor`, so a swapped wiring fails here.
  Proof: cmd/libretto/loop_test.go TestLoopDrivesClaudeWhenNoOtherToolIsNamed
- **`sessionFor("opencode")` shall return the executable `opencode` and an argument
  builder whose output equals `opencodeArgs`.**
  Proof: cmd/libretto/loop_test.go TestLoopDrivesClaudeWhenNoOtherToolIsNamed
- **If the tool is `codex`, then `loop` shall return an error containing `codex` and
  `/libretto-attacca` for a change that does not exist, with or without `--dry-run`,
  while PATH is empty and without resolving the run log path.** The refusal comes before
  the checklist because no checklist makes codex able to run it.
  Proof: cmd/libretto/loop_test.go TestLoopRefusesCodexAndSaysWhy
- **If `--max-budget-usd` is given while the tool is `opencode`, then `loop` shall return
  an error containing `--max-budget-usd` and `opencode` before reading the checklist.**
  Proof: cmd/libretto/loop_test.go TestOpencodeRefusesABudgetItCannotEnforce
- **If the tool is `opencode` and `opencode` is not on PATH, then `loop` shall return an
  error containing `opencode is not on PATH` and `install OpenCode, or use --dry-run to
  see what would run`.**
  Proof: cmd/libretto/loop_test.go TestTheMissingToolIsNamed
- **If the tool is `claude` and `claude` is not on PATH, then `loop` shall return the
  error `claude is not on PATH — the loop relaunches it once per task`, a newline, seven
  spaces, and `install Claude Code, or use --dry-run to see what would run`.** Each tool's
  install line names its own install.
  Proof: cmd/libretto/loop_test.go TestTheMissingClaudeMessageIsUnchanged
- **If the tool is `opencode` and the change has no checklist, then `loop` shall return
  the `no checklist at` error while PATH is empty.** The missing plan is still reported
  before the missing dependency.
  Proof: cmd/libretto/loop_test.go TestTheMissingChecklistIsReportedBeforeTheMissingOpencode
- **While `--dry-run` is set and the tool is `opencode`, `loop` shall succeed with an empty
  PATH.**
  Proof: cmd/libretto/loop_test.go TestDryRunNeedsNoOpencodeOnPath
- **The loop help shall name `--opencode` on a `loop` line.**
  Proof: cmd/libretto/main_test.go TestHelpNamesTheLoopTool

#### What OpenCode is handed, and what it reports

- **The opencode argument list shall be exactly `run`, `--command`, `libretto-attacca`,
  `--format`, `json`, in that order.** No positional message, no budget, no model.
  Proof: cmd/libretto/loop_test.go TestOpencodeArgsRunTheAttaccaCommand
- **`opencodeStdin` shall return the loop prompt with its leading `/libretto-attacca `
  removed, byte for byte** — newlines and the double quotes around the change name
  included. Ceiling, named: this proves the pure function only; that the runner sets the
  session's `Stdin` from it is unobserved, because no test starts `opencode`.
  Proof: cmd/libretto/loop_test.go TestOpencodeArgumentsArriveOnStdinVerbatim
- **The opencode argument list shall contain none of `--auto`, `--yolo` and
  `--dangerously-skip-permissions`.**
  Proof: cmd/libretto/loop_test.go TestOpencodeArgsGrantNoPermissions
- **When the session's output is OpenCode's JSON event stream, the opencode parser shall
  return the `sessionID` its events carry as the session id.**
  Proof: cmd/libretto/loop_test.go TestParseOpencodeReadsTheSessionID
- **When the stream carries `text` events, the loop shall print their `part.text`, in the
  order they arrived, to stdout.**
  Proof: cmd/libretto/loop_test.go TestOpencodeTextReachesTheTerminal
- **If the output carries neither a `text` nor an `error` event, or its first line does
  not parse as JSON, then the loop shall print the raw output unchanged.** A session never
  reaches the terminal as nothing — the rule `parseSession` keeps for `claude`.
  Proof: cmd/libretto/loop_test.go TestOpencodeWithoutTextPrintsTheRawOutput
- **When the stream carries an `error` event, the loop shall print `opencode error: `
  followed by its message after the session's text**, the message being the event's
  `error.data.message` when present and its `error.name` otherwise, several joined with a
  newline.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeErrorEventIsPrintedAndRecorded
- **If an event in the stream is longer than 1 MiB, then the opencode parser shall still
  return the session id and the text of the events after it.**
  Proof: cmd/libretto/loop_test.go TestALongOpencodeEventDoesNotEndTheParse

#### The per-round budget and the run log

A round is one session of the chosen tool, and every round that runs, under either tool,
leaves one JSON line in `<git-common-dir>/libretto/loop.jsonl` — inside git's own
directory, so nothing is tracked, `git status` shows nothing new, and every worktree of the
repository shares one log. What the host reports is recorded verbatim and never
interpreted. Under `claude` a round is one `claude -p` session, so the host's own
per-session `--max-budget-usd` is exactly a per-round budget: the loop passes it through
and the host enforces it, and every session is asked for its JSON result. `opencode run`
has no per-session budget, so the flag is refused under `--opencode`, and its stream
carries no cost the loop may record.

- **When the loop starts a session, the `claude` argument list it builds shall carry
  `-p <prompt>` and `--output-format json`.**
  Proof: cmd/libretto/loop_test.go TestEverySessionIsAskedForItsJSONResult
- **Where `--max-budget-usd <amount>` is given, the `claude` argument list shall carry
  `--max-budget-usd <amount>` with the amount exactly as typed; where it is not given, the
  list shall carry no `--max-budget-usd`.**
  Proof: cmd/libretto/loop_test.go TestClaudeArgsCarryTheBudgetOnlyWhenGiven
- **If `--max-budget-usd` is missing its value, or its value does not match
  `^[0-9]+(\.[0-9]+)?$`, or it equals zero, then `parseLoopArgs` shall return an error
  naming `--max-budget-usd`** — the table covers `1e2`, `0x1`, `+3`, `-1`, `inf`, `.5`,
  `0` and `0.00`.
  Proof: cmd/libretto/loop_test.go TestParseLoopArgsRefusesABadBudget
- **While a budget is set, the loop header shall print the per-round budget and the
  product of `--max` and that budget formatted with `%.2f` as the worst case for the run**
  — `--max 3 --max-budget-usd 0.1` prints `0.30`.
  Proof: cmd/libretto/loop_test.go TestTheHeaderStatesTheWorstCaseSpend
- **The loop help shall name `--max-budget-usd`.**
  Proof: cmd/libretto/main_test.go TestHelpNamesTheLoopBudget
- **When a `claude` session's output is the host's JSON result object, the session
  parser shall return its `total_cost_usd`, `duration_ms`, `session_id`, `subtype`,
  `is_error` and `result`; if the output is not JSON or carries no `total_cost_usd`, then the parser
  shall report the cost as unknown, never as zero.**
  Proof: cmd/libretto/loop_test.go TestParseSessionResultReadsTheHostsCost
- **When a `claude` session ends, the loop shall print the `result` text of its JSON to
  stdout; if the output does not parse as JSON, then the loop shall print it unchanged.**
  JSON that
  carries no `result` field — a bare `null`, `{}` — is not the host's report either, and
  its raw output is what gets printed, so a session never reaches the terminal as nothing.
  Proof: cmd/libretto/loop_test.go TestTheSessionsTextStillReachesTheTerminal
  Proof: cmd/libretto/loop_test.go TestJSONWithoutAResultPrintsTheRawOutput
- **When a round ends, the loop shall append exactly one JSON line to the run log.**
  Proof: cmd/libretto/loop_test.go TestEveryRoundAppendsOneRecord
- **Each record shall carry `time` as an RFC 3339 timestamp of when the round ended,
  `change` as the change name, and `round` as the 1-based round number.**
  Proof: cmd/libretto/loop_test.go TestTheRecordSaysWhenWhichChangeAndWhichRound
- **Each record of a round whose checklist was read after the session shall carry
  `closed` as `after.done − before.done`, negative when a box reopened, and `open` as the
  open boxes after the round.**
  Proof: cmd/libretto/loop_test.go TestClosedIsSignedSoAReopenedBoxIsNegative
- **Where the session's JSON carried them, the record shall carry `cost_usd`,
  `duration_ms`, `session_id`, `subtype` and `is_error` with the host's values
  unchanged.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheHostsReportVerbatim
- **If the session reported no cost, then the round's record shall omit `cost_usd`
  entirely rather than write zero; if it reported no `is_error`, then the record shall
  omit `is_error`.** Unknown and zero are different facts.
  Proof: cmd/libretto/loop_test.go TestAnUnpricedRoundOmitsCostRatherThanWritingZero
- **Where a budget is set, each round's record shall carry it as `budget_usd`; where none
  is set, the record shall omit `budget_usd`.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheBudgetThatApplied
- **If the session exited non-zero, then the record shall carry the exec error's message
  as `error`; when it exited zero and the host reported no `error` event, the record shall
  omit `error`.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheExecErrorOnlyOnFailure
- **If an OpenCode stream carries an `error` event and the session exited zero, then the
  round's record shall carry that message as `error`, as reported.** On a non-zero exit
  `error` keeps the exec error, as it does for `claude`.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeErrorEventIsPrintedAndRecorded
- **While the tool is `opencode`, each record shall omit `cost_usd`, `duration_ms`,
  `subtype` and `is_error`, even when the stream's `step_finish` events carry a `cost`.**
  The fixture carries a non-zero `part.cost`, or an implementation that summed it would
  pass.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeRoundRecordsNoCost
- **Each record the loop writes shall carry `tool` set to the tool that ran the round.**
  A record without one was written before the field existed, and means `claude`.
  Proof: cmd/libretto/loop_test.go TestTheRecordNamesTheTool
- **When the loop stops because every box closed, two rounds closed nothing, the cap was
  reached, or the checklist could not be read, the last record it wrote shall carry
  `stopped` set to `done`, `stuck`, `cap` or `error` respectively, and no earlier record
  shall carry `stopped`.**
  Proof: cmd/libretto/loop_test.go TestTheLastRecordSaysWhyTheLoopStopped
- **If the checklist cannot be read after a session, then the loop shall append that
  round's record with `closed` and `open` omitted, `stopped` set to `error` and the read
  error in `error`, and shall return the error it returns today.**
  Proof: cmd/libretto/loop_test.go TestAVanishedChecklistIsStillRecorded
- **The loop shall write the run log at `<git-common-dir>/libretto/loop.jsonl`, and
  `git status --porcelain` in the project shall print nothing after a loop over a clean
  tree.**
  Proof: cmd/libretto/loop_test.go TestTheRunLogLeavesGitStatusClean
- **If the git common dir cannot be located, then the loop shall print the one warning
  `run log not recorded: git common dir not found` and shall run every round it would
  have run with a log.** The checklist is the state the loop decides on; the log is a
  record of it, and losing the record is no reason to stop work.
  Proof: cmd/libretto/loop_test.go TestAnUnlocatableLogWarnsAndCarriesOn
- **If the run log cannot be written, then the loop shall print one warning naming the
  log's path, however many rounds fail to write, and shall run every round it would have
  run with a writable log.**
  Proof: cmd/libretto/loop_test.go TestAnUnwritableLogDoesNotStopTheLoop
- **When `loop` refuses a missing checklist or a missing `claude`, it shall not have
  resolved the run log path; while `--dry-run` is set, the loop shall not resolve the
  path, shall print no log warning, and shall write no record.** The path is resolved
  after both refusals, so *the missing plan is reported before the missing dependency*
  holds unchanged.
  Proof: cmd/libretto/loop_test.go TestDryRunWritesNoRecord

### `metrics` — what the flow cost, derived from git and the loop's own log

- **The change folders and the flow shall stay uninstrumented; the only recorded data
  shall be the loop's own run log, written by the binary and nothing else.** The obvious
  build has each phase write a timestamp and the report read them back; that costs a new
  artifact in every change folder which eight skills have to remember, and a metric
  collected only when somebody remembers has holes exactly where the interesting runs are.
  Git already holds the answer: a change folder's first and last commits are the wall
  clock, the commits between are the iterations, and the checklist's diffs are the churn.
  It works retroactively on every change ever landed. The run log does not reopen that
  rule: it has one writer, the loop, which runs every round by construction and cannot
  forget — nothing lands in a change folder, and no skill learns anything. `metrics` only
  reads it.
  Proof: cmd/libretto/metrics_test.go TestChangeNamesSeesLandedChangesNotJustOpenOnes
  Proof: cmd/libretto/metrics_test.go TestMetricsNeverWritesTheRunLog
- **`--full-history` on every query, and it is load-bearing.** Git's default history
  simplification prunes commits that do not change a path's *final* state, and a landed
  change's final state is "does not exist" — so the entire history of every change this
  report exists to measure is simplified away. Measured, not reasoned: without it this
  repository reported 12 changes instead of 22, two of them with no history at all, and
  **no `plan.md` was found for any of them**, so every churn column read as "no plan" on
  changes that plainly had one.
  Proof: cmd/libretto/metrics_test.go TestChangeNamesSeesLandedChangesNotJustOpenOnes
- **churn is read from `proposal.md`'s `## Checklist`, then from `tasks.md`, then from
  `plan.md`, and the old-layout queries are not optional.** Every change that landed before
  2026-10-21 carries its checklist in `tasks.md` — or, before 2026-08-17, `plan.md` — and a
  folder that no longer exists, so it can never acquire the new one; and this report is
  retroactive by design. Dropping a legacy path would blank the churn column for the whole
  history to tidy a filename, which is the same failure the bullet above already records
  paying for once. Separate queries and not one pathspec: `plan.md` is the technical
  approach now, a prose document allowed to show a checkbox in an example, and counting its
  diffs would inflate every change from here on. Each change is read from one file; counts
  from two are never added together.
  Proof: cmd/libretto/metrics_test.go TestMetricsFallsBackToLegacyPlan
- **Where a change's `proposal.md` history (`--full-history --diff-filter=AM -p`) carries at
  least one box line under `## Checklist`, `measure` shall report `checked`, `uncheck` and
  `boxes` from that history alone**, even when the `tasks.md` history for the same change
  also carries boxes. The fixture gives the two files different numbers, so a sum or the
  wrong file fails.
  Proof: cmd/libretto/metrics_test.go TestMetricsReadsChurnFromTheChangeDocument
- **Where a change's `proposal.md` history is non-empty but carries no box line under
  `## Checklist`, `measure` shall report churn from `tasks.md`, or from `plan.md` when
  `tasks.md` has none; where neither carries a box, it shall leave `planSeen` false so the
  row prints a dash.** A first non-empty diff is not a signal: the capture adds every
  `proposal.md`, so testing for a non-empty diff with `proposal.md` queried first would blank
  the churn of every old-layout change — watched failing under exactly that order.
  Proof: cmd/libretto/metrics_test.go TestAProposalWithoutAChecklistDoesNotShadowTheOldLayout
- **When `measure` queries `proposal.md` history, it shall pass `--unified=`, and `churn`
  shall count a `+` box only when its post-image section is `## Checklist` and a `-` box only
  when its pre-image section is `## Checklist`.** The post-image is context and `+` lines,
  the pre-image context and `-` lines, and the section resets at each commit. `fakeGit`
  refuses a `proposal.md` query without `--unified=`. **Named limit:** a commit that renames
  a heading so that untouched boxes move into or out of `## Checklist` changes no box line,
  so it moves neither the churn nor the denominator.
  Proof: cmd/libretto/metrics_test.go TestChangeDocumentChurnIsScopedToTheChecklist
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
  printed what it cannot measure; what it does measure deserves no less. **The legend's `—`
  line shall state that the change has no checklist in its history and name `proposal.md`,
  `tasks.md` and `plan.md`** — it said "no plan.md in its history", stale since 2026-08-17.
  Proof: cmd/libretto/metrics_test.go TestTheReportExplainsItsColumns
- one change may be named, by full name or unambiguous prefix — a name is typed from
  memory of how it starts. Exact wins even when it prefixes a sibling, an ambiguous
  prefix is refused naming the candidates, and an unknown name and an unknown flag keep
  their refusals. `--json` is the one flag it reads, before or after the name.
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
- **corrections are broken down by phase, spelled as the ledger spells them.** One row per
  distinct phase value, count beside it, sorted count-then-name so the biggest number is
  the first thing read — the instrument the ask-more-questions debate gets decided with:
  corrections at phase 2 are questions working, findings at 6→7 are questions that were
  missing, corrections after phase 8 are the expensive kind. No normaliser — a second
  spelling drifts — and no verdicts: the counting is here, the reading stays human.
  Proof: cmd/libretto/metrics_test.go TestCorrectionsCountsByPhase
- **the breakdown renders in the corpus report with its explanation line.**
  Proof: cmd/libretto/metrics_test.go TestMetricsReportsCorrectionsByPhase
- **`6→7` entries are reviewer findings, not user corrections.** They count in the phase
  breakdown and never in the per-change corrections column, whose meaning does not move.
  Proof: cmd/libretto/metrics_test.go TestReviewerFindingsStayOutOfCorrections
- **a `6→7` entry with no change open is still a finding.** The exclusion outranks the
  orphan count, so it never lands in the corrections-outside-any-change line — surfaced
  by this change's own 6→7 reviewer, on the first entries the rule ever produced.
  Proof: cmd/libretto/metrics_test.go TestOrphanFindingIsNotAnOrphanCorrection
- **absent, empty and populated ledgers are three facts.** Absent says capture is not in
  use; present-but-empty says so in one line and renders no phase rows.
  Proof: cmd/libretto/metrics_test.go TestMetricsPhaseBreakdownAbsentLedger

#### Token cost, read off the session transcripts

The git half above answers what a change took. It cannot answer what a change **cost**,
so a context diet was guesswork — which bit once: five review lenses cost ~307k and the
lens count turned out not to be the lever.

The transcripts under the host's `projects/` directory carry `usage` on every assistant
turn. This is a free rider on them, exactly as the corrections count free-rides on the
lessons ledger, and it clears the same bar the no-instrumentation rule sets: nothing new
is written, and nothing under the transcript root is written at all.

- **the four usage numbers are read apart and never summed**, from `assistant` entries
  only — the sole type carrying usage, measured at 8,944 of 8,944 — with
  `.message.usage.iterations[]` ignored, because it repeats the same four numbers and
  adding both double-counts every entry that has one. `int64` throughout: this
  repository's cache-read total is 2.3 billion against 33 thousand plain input tokens, a
  ratio near 70,000:1, **which is the whole reason a single "tokens" number would be a
  cache-read number in disguise**.
  Proof: cmd/libretto/usage_test.go TestTheFourUsageNumbersAreKeptApart
- **a malformed line, a missing usage object and an absent field are all survivable.** An
  81 MB corpus written by another program across many versions will contain something
  unparseable, and whole entry types carry no envelope at all — no `cwd`, no `gitBranch`,
  no `timestamp`. One bad line must not cost the other nine thousand.
  Proof: cmd/libretto/usage_test.go TestAMalformedLineDoesNotCostTheRestOfTheFile
- **subagent transcripts are counted.** They sit one level down, under
  `<sessionId>/subagents/`, and one review-lens file alone holds 3,026,036 cache-read
  tokens. **Ignoring them undercounts by more than the number being reported.** The two
  locations are enumerated rather than walked for every `*.jsonl`, because the host keeps
  `tool-results/` and `memory/` in the same tree and counting a mirror would inflate the
  total silently.
  Proof: cmd/libretto/usage_test.go TestSubagentTranscriptsAreCounted
- **a dot in the repository path encodes to a dash, like a separator.** A home directory
  called `pau.sanchez` maps to `pau-sanchez`, so a reader that preserved dots finds no
  transcripts at all for that user. **Every fixture used a dotless path and none could
  have caught it** — running the built binary against a real tree is what did. The
  encoding is therefore lossy twice over and is never inverted.
  Proof: cmd/libretto/usage_test.go TestADotInThePathBecomesADashToo
- **`gitBranch` is read per entry, never per file.** One real session file was measured
  spanning four branches; a per-file reading misattributes every entry after a checkout.
  Proof: cmd/libretto/usage_test.go TestBranchIsReadPerEntryNotPerFile
- **a conventional branch prefix is stripped and the rest is matched whole.** `feat/`,
  `fix/`, `docs/`, `chore/`, `refactor/`, then an exact match against the change names git
  has seen — never a prefix and never a substring, because a looser rule turns one wrong
  guess into a number nobody can audit. An unrecognised prefix fails safe: unattributed
  rather than the wrong change.
  Proof: cmd/libretto/usage_test.go TestAPrefixIsStrippedAndTheNameMatchedWhole
- **what cannot be attributed is bucketed and printed, never dropped, and the invariant is
  readable off the output.** Attributed plus unattributed equals the corpus, on all four
  numbers. **This is the criterion that carries the feature's honesty**: attribution is a
  heuristic, not a mechanism — measured on this repository, 62% of entries attribute to no
  change, because work happens on `main`, in a detached `HEAD`, and on branches named
  unlike their change. The bucket is a headline rather than a footnote for that reason,
  and the totals stay corpus-wide under a change filter so the invariant survives scoping.
  Proof: cmd/libretto/usage_test.go TestUnattributedTokensAreReportedNotDiscarded
  Proof: cmd/libretto/metrics_test.go TestTheTokenFooterIsCorpusWideUnderAFilter
- **a miss rate above zero never prints as `0%`.** Integer division truncates, so one
  unattributed entry in two hundred read as `0%` — the opposite of what a number the
  report calls its own error bar is for. Under one percent it says `<1%`.
  Proof: cmd/libretto/metrics_test.go TestASmallMissRateDoesNotRoundAwayToZero
- **three states, and each prints differently.** No transcript root replaces the block
  with one line; a root that reached nothing for a change prints a **dash**; a change that
  was attributed and genuinely cost nothing prints **zeros**. The distinction is presence
  in the attribution map, never a zero total — a `<synthetic>` entry carries an all-zero
  usage object and reaches the third state, so it is reachable rather than theoretical,
  and the first implementation collapsed it into the dash.
  **This extends the ledger's dash rule rather than reusing it**: absent-versus-zero is
  two states there and three here, and collapsing them would hide the miss rate above.
  Proof: cmd/libretto/metrics_test.go TestAChangeWithNoTokensReportsADashNotAZero
  Proof: cmd/libretto/metrics_test.go TestAnAttributedChangeWithZeroTokensIsNotADash
- **the per-phase block carries its own unattributed row.** `attributionSkill` records the
  flow phase, at a measured 4,728 entries against 4,216 without, so the remainder is shown
  rather than distributed across the phases that did name one.
  Proof: cmd/libretto/metrics_test.go TestPerPhaseCostCarriesAnUnattributedRow
- **no transcript root is a state, not an error.** A checkout that never hosted a session
  is the normal case; the git-derived report still prints in full.
  Proof: cmd/libretto/metrics_test.go TestNoTranscriptRootStillReportsTheGitMetrics
- **nothing under the transcript root is written.** The witness is a snapshot of every
  path, its size and the SHA-256 of its contents, taken before the read and compared
  after — red on a create, a delete, a truncation and an in-place rewrite alike, which an
  mtime comparison is not.
  Proof: cmd/libretto/usage_test.go TestTheTranscriptRootIsNeverWritten
- **the ceiling says cost is measured, duration is not, and names the token block's own
  limit.** Nothing in `flowCeiling` was retracted: it named per-phase *duration* and
  `review-work` findings, and both still need a phase to write them down. Cost was never
  on that list — it was simply not measured. Duration stays off with its own reason now:
  the entries carry timestamps, but a phase's wall clock includes every wait for a human
  and would report attention the work never had.
  Proof: cmd/libretto/metrics_test.go TestTheCeilingSeparatesCostFromDuration

**No new column in the change table, and no pricing.** The table is seven columns wide and
the tests pick the corrections cell by counting from the right, so a column on that side
breaks a proof for a cosmetic reason — and it would force a choice of one number among
four, or a composite. Prices change under the repository, and a stale price table reports
confident nonsense; tokens are what was measured. The binary prices nothing anywhere: the
only dollars in the report are the loop block's below, and those are the host's own
figure, recorded as it reported them.

**Plain digits, no thousands separator and no exponent.** A space inside a number breaks
every pipe a report gets read through, and an exponent throws away the precision that
makes two runs comparable, which is the entire use this was built for. Dollars follow the
same rule at `%.2f`.

#### Loop rounds, read off the run log

The loop's run log is the one artifact `metrics` reads that git does not hold. It gets its
own block after the token block, one row per change, never a column in the change table.
The cost is the host's `total_cost_usd` — a list-price estimate by the host's own label,
not a bill — and the column says so.

- **Where a run log exists, `metrics` shall print a loop block with one row per change,
  sorted by change name, holding its rounds, the net sum of its `closed` values under a
  header that says net, the sum of its `cost_usd` values formatted with `%.2f`, and the
  count of its rounds with no `cost_usd`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsReportsLoopRoundsPerChange
- **If a change's rounds all lack `cost_usd`, then its cost cell shall print a dash, never
  `0.00`.**
  Proof: cmd/libretto/metrics_test.go TestAnUnpricedChangeReportsADashNotAZero
- **If no run log exists, then `metrics` shall print the one line `no loop runs are
  recorded in this repository`, and shall print every other block as it does today.** A
  log that holds no record reads the same: an empty file is no runs, never an empty block.
  The log is shared by every worktree, which is why the line says *repository*.
  Proof: cmd/libretto/metrics_test.go TestNoRunLogIsAStateNotAnError
- **If a run log line does not parse as a record, then `metrics` shall skip it, keep every
  other row, and print the number of skipped lines counted over the whole file, whether
  or not a change is named.** A malformed line names no change, so a filtered count would
  hide it under every filter.
  Proof: cmd/libretto/metrics_test.go TestAMalformedRunLogLineIsSkippedAndCounted
- **Where a change is named, the loop block shall print only that change's row; if the
  named change has no record, then it shall print the one line `no loop rounds recorded
  for <change>`.**
  Proof: cmd/libretto/metrics_test.go TestTheLoopBlockFollowsTheChangeFilter
- **The loop block shall label its cost column as the host's reported estimate in
  dollars.**
  Proof: cmd/libretto/metrics_test.go TestTheLoopBlockLabelsCostAsTheHostsFigure

#### A per-agent model recommendation, and it is never applied

`libretto models set` and `models effort` have existed since the selector did, but which
agent suits which tier lived in one paragraph of a skill and on nobody's machine. This
ships it as data and prints it.

**Recommend only.** Nothing writes a model or an effort because of the table — not a flag,
not a first run, not the panel. The reading is the repository's; the typing is the user's,
which is the same split `AGENTS.md` draws for the release bump and is here for the same
reason: a tool that acts on its own reading is a tool whose reading nobody audits.

**The table lives in `cmd/libretto`, and that is a layering decision rather than a
convenience.** `agent-models` settled, from a real failure, that *the subject is a
directory, not the repository* — it is handed a directory and works on every `*.md` in it,
whoever created them. A map from `review-lens-security` to a tier **is** this payload's
agent list, so putting it there would reintroduce the coupling that decision exists
against, in the one package whose layering claim is that it is testable against a bare
`t.TempDir()`. This binary already knows what the payload is; it gates commands on it.

- **every recommendation names a model the catalogue accepts**, and an effort the
  recommended model actually supports. Walked, not reviewed: an entry naming `sonnet-5`,
  or a level on `haiku`, is a recommendation this binary would refuse to type — the one
  thing worse than saying nothing.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationIsTypeable
- **nothing is recommended onto the two priciest tiers.** The catalogue orders itself
  cheapest-first because this exists to lower a bill; a table steering onto `opus` or
  `fable` inverts what it was built for. An agent that genuinely needs more is a decision
  a person makes at the screen.
  Proof: cmd/libretto/recommend_test.go TestNothingIsRecommendedOntoThePriciestTiers
- **an agent the table does not know gets no recommendation and no reason.** Silence, never
  a guess — every agent a user writes themselves lands here, and an invented opinion about
  somebody's own agent is worse than none.
  Proof: cmd/libretto/recommend_test.go TestAnUnknownAgentGetsNoRecommendation
- **every recommendation carries a reason.** A verdict with no reason is an instruction,
  and nothing here gives instructions. The model catalogue's labels already answer "when
  would I pick this"; this is the same idea keyed by agent.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationCarriesAReason
- **the listing carries the reasons, as a trailer rather than a fourth column**, and an
  agent with no recommendation is absent from it rather than blank in it — an empty cell
  reads as *none recommended* where the truth is *not known*. **This listing is the only
  surface wide enough for them**: the panel's selector narrows to 58 columns and the
  reasons run to seventy runes, so what reaches the screen there is a mark and this is
  where it is explained.
  Proof: cmd/libretto/models_test.go TestModelsListsTheRecommendationAndItsReason
- **it says when an agent runs against its recommendation**, and only then. A
  recommendation nobody can tell they are ignoring is a recommendation that changes
  nothing. Four of the seven diverge on the day this ships, `review-lens-design` among
  them deliberately: the skill's own reasoning puts design on the cheap end and the file
  declares `sonnet`, so **the disagreement is shown rather than settled silently here**.
  Proof: cmd/libretto/models_test.go TestModelsMarksAgentsRunningAgainstTheRecommendation
- **the panel is handed the answer through the existing adapter**, exactly as it was
  handed `Efforts`, so `internal/ui` still imports no `agentmodel`.
  Proof: cmd/libretto/models_test.go TestAgentRowsCarryTheRecommendation
- **For every key in the recommendation table, the repository shall hold `agents/<key>.md`.**
  The CLI does not own the agent set, so a key is valid only when the payload ships its file,
  read from `../../agents` as `gates_test.go` already reads repository files. `spec-reviewer`
  is recommended sonnet/high with `work-reviewer`'s reason — a cheap miss there is a false
  green (user, 2026-10-21) — and `spec-writer` keeps sonnet/high, because it survives for the
  fan-out.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationNamesAShippedAgent
- **The recommendation table shall have no `plan-writer` or `task-cutter` key.** Both agents
  were retired on 2026-10-21. Watched red with a `plan-writer` key added, then reverted.
  Proof: cmd/libretto/recommend_test.go TestEveryRecommendationNamesAShippedAgent

**The values themselves are reviewed and never tested**, and review caught one: the table
shipped `high` for `review-lens-intent` while the change's own table omitted it, and every
guard above stayed green — they walk the table for what a machine can decide and none of
them can read a paragraph.

### `wiki` — the specs rendered browsable

Every project the flow runs in accumulates capability specs, and a directory of
`spec.md` files is a specification only its author browses. `wiki` renders them into
one `README.md` inside the specs directory — index table, then a section per
capability — so the forge shows the specification to anyone who opens the directory.
The file is generated output behind a marker, never payload: nothing about it is
symlinked, and it regenerates wherever the specs move.

- **the `wiki` subcommand shall discover the consolidated specs directory in the order
  `.agents/specs`, `specs`, `openspec`, `docs/specs`, `spec`, first hit wins** — the
  same walk `spec-drift` does, verbatim, because two tools with two discovery orders
  would disagree about which specification a project has.
  Proof: cmd/libretto/wiki_test.go TestWikiDiscoversSpecsDirInDriftOrder
- **when run in a project whose specs directory holds `*/spec.md` files, `wiki` shall
  write `README.md` inside that directory** with an index row and a section per
  capability, each section carrying the capability's first paragraph, its `Governs:`
  line, its criteria, and a relative link to its spec. A spec with no prose before its
  first heading gets a section without an intro — absent is absent, never an error.
  Proof: cmd/libretto/wiki_test.go TestWikiWritesIndexAndSections
- **the generated file shall open with a marker comment** naming `libretto wiki` as its
  generator and the command that refreshes it.
  Proof: cmd/libretto/wiki_test.go TestGeneratedReadmeCarriesTheMarker
- **if the target `README.md` exists and does not carry the marker, then `wiki` shall
  refuse**, report the conflict, exit non-zero, and leave the file untouched — the
  never-overwrite rule applied to output the tool writes into a directory people also
  write into.
  Proof: cmd/libretto/wiki_test.go TestWikiNeverOverwritesAHandWrittenReadme
- **if no specs directory is found, or the directory found holds no `*/spec.md`, then
  `wiki` shall report which of the two it is in one line, exit non-zero, and write
  nothing** — an empty specification has nothing to index, and a marker-only README
  would be noise pretending to be a wiki.
  Proof: cmd/libretto/wiki_test.go TestWikiReportsNoSpecsAndExitsNonZero
- **the `wiki` subcommand shall produce byte-identical output for unchanged input.**
  Nothing in the render path passes through a map; determinism is structural.
  Proof: cmd/libretto/wiki_test.go TestWikiOutputIsDeterministic
- **the `wiki` subcommand shall write nothing but the one `README.md`.**
  Proof: cmd/libretto/wiki_test.go TestWikiWritesNothingButTheReadme

The same specification renders as a second view: `wiki --html` writes one
self-contained page — sidebar index with criteria counts, a section per
capability, an inline criteria filter, light and dark via `prefers-color-scheme` —
that a browser opens from `file://` with no server and no build step. Both views
are generated output behind markers; the tool refreshes its own and refuses
everything else.

- **when invoked with `--html`, the `wiki` subcommand shall write `wiki.html` into
  the discovered specs directory** — a single self-contained page carrying a
  navigation entry and a section per capability, each section with the capability's
  intro, its `Governs:` line and its criteria — and shall not write or modify
  `README.md` in that run. A `**bold**` and a backticked span in spec prose arrive
  as `<strong>` and `<code>` markup.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesTheViewer
- **the generated `wiki.html` shall open with a first-line HTML-comment marker**
  naming `libretto wiki` as its generator and the command that refreshes it.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLCarriesTheMarker
- **if the target `wiki.html` exists and does not carry the marker, then `wiki
  --html` shall refuse**, report the conflict, exit non-zero, and leave the file
  untouched.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLNeverOverwritesAForeignFile
- **the `--html` output shall carry every spec-sourced string HTML-escaped**, so
  markup or script inside a spec renders as text and never executes. Escaping runs
  before the inline conversions, always — that order is the criterion.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLEscapesSpecContent
- **the `--html` output shall reference no external resource beyond font
  stylesheets from `fonts.googleapis.com`** (and the `fonts.gstatic.com` files they
  pull) — no external scripts, no other stylesheets, no remote images — so the page
  works offline, fonts degrading to their declared fallback stacks.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsSelfContained
- **the `--html` run shall produce byte-identical output for unchanged input.**
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsDeterministic
- **when run without `--html` and the specs directory holds a `wiki.html` carrying
  the marker, the subcommand shall refresh that file in the same run**; where the
  `wiki.html` present carries no marker, the plain run shall leave it alone and not
  fail — refusal belongs only to `--html` targeting it explicitly.
  Proof: cmd/libretto/wiki_test.go TestPlainWikiRefreshesAMarkedHTMLView
- **the `--html` run shall write nothing but the one `wiki.html`.**
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLWritesNothingButTheOneFile
- **if `wiki` is given an argument other than `--html`, then it shall report the
  unknown argument and exit non-zero, writing nothing.**
  Proof: cmd/libretto/wiki_test.go TestWikiRejectsAnUnknownFlag

**Untested by decision, declared:** the filter's behaviour and the theming run in a
browser Go tests cannot drive; the builder renders and looks, and the evidence names
what was looked at. The inline conversion is not on that list — it is proven above.

The panel offers the same command as a row, and only where pressing it can
succeed. Project scope is the user's line, drawn explicitly; the missing-specs
arm is the `models` precedent — an entry that opens an empty screen is a promise
the panel cannot keep.

- **Where the panel scope is project and the project holds a consolidated specs
  directory, `panelData` shall include an enabled `wiki` row; where the scope is
  global, or the project holds no specs directory, it shall not include one.**
  Visibility shares the command's own discovery through `findSpecsDir` — never a
  second list.
  Proof: cmd/libretto/wiki_test.go TestPanelOffersWikiOnlyInAProjectWithSpecs
- **When the `wiki` row is dispatched, the panel shall run `wiki --open` against
  the project directory** — the press ends with the viewer in the browser, by
  instruction 2026-08-18. The project directory reaches `dispatch` from `run`'s
  single resolution — a second `os.Getwd` would be the two-answers bug this spec
  already records.
  Proof: cmd/libretto/wiki_test.go TestPanelWikiRowOpensTheViewer

- **When invoked with `--open`, the `wiki` subcommand shall write the HTML view
  exactly as `--html` does and then invoke the opener seam with the written
  file's path, after the write succeeds**; if the opener errors, the run shall
  surface that error naming the written path, non-zero — the file stays, and the
  message hands the user what to open by hand.
  Proof: cmd/libretto/wiki_test.go TestWikiOpenGeneratesAndOpensTheViewer
- **If generation refuses or fails, then the `--open` run shall not invoke the
  opener.**
  Proof: cmd/libretto/wiki_test.go TestWikiOpenDoesNotOpenOnFailure
- **The opener argv shall name `open` on darwin and `xdg-open` on every other
  platform**, with the file's path as the argument — a pure function, so the
  mapping is testable without executing it.
  Proof: cmd/libretto/wiki_test.go TestOpenerArgvPerPlatform

The viewer is a home of capability cards and one page per capability — not one
infinite scroll. Minimalist and visual at once: the cards carry the numbers, a
mini-bar and the last-changed date from git; CSS scroll-driven motion stays
behind `prefers-reduced-motion`, all self-contained.

- **The generated page shall carry a home section with one card per
  capability** — each card a plain `#`-anchor showing the capability's name,
  its criteria count, and a mini-bar whose width is proportional to the
  project's maximum count — followed by one article per capability carrying a
  home link, its `Governs:` line, intro and criteria; zero-criteria
  capabilities get a card and a page like every other.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLIsHomeAndPages
- **Where the git seam yields a last-changed date for a spec, the card and its
  article shall both carry it; where it yields nothing, neither shall render a
  date and the run shall succeed unchanged** — no git, no repository, no
  history is absence, never an error.
  Proof: cmd/libretto/wiki_test.go TestWikiDatesComeFromGitAndDegrade
- **The page shall declare its reading-progress and reveal animations only
  inside `prefers-reduced-motion: no-preference`**, referencing
  `animation-timeline` rather than scroll listeners for both.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLMotionRespectsReducedMotion
- **The stylesheet shall define the complete dark palette as tokens on the bare
  `:root` (dark-first), with a full light redefinition in a
  `prefers-color-scheme: light` block carrying the same token names**, and no
  hex colour literal outside the two token blocks — components take colour only
  through `var()`.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLThemesAreTokenComplete
- **The routing shall ship as an inline `hashchange` script toggling visibility
  classes only — document order home-first, every article present in the
  markup, every navigation a plain anchor** — so a JS-less render shows the
  whole wiki and every link still lands. Its runtime behaviour is browser-land,
  untested by decision beside the search and the theming; the structure is what
  Go pins.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLRouterIsAnEnhancement
- **The home search shall filter cards client-side by capability name and
  criteria text**, inline and offline, over escaped render-time `data-crit`
  attributes.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLHomeSearchIsInline

The home is the project board: what moved, what is being built, how healthy the
contract is — advisory mirrors of what spec-drift and the loop already define,
never a second authority and never a gate.

- **The home shall carry a recently-changed rail of up to three of the most
  recently changed capabilities** — most recent first, each with its date and
  the last commit subject from the git seam — omitting the rail entirely when no
  capability has a date.
  Proof: cmd/libretto/wiki_test.go TestWikiHomeCarriesTheRecentRail
- **Where the changes directory holds a checklist with open boxes — chosen per change by
  `checklistPath`, the loop's rule — the home shall carry an in-flight strip** naming each
  such change with its closed/total box count and the count of `Queued:` proposals; where
  there are none it shall be absent — a queue alone does not summon the strip (assumed
  2026-08-19 under attacca; if wrong, the fix is rendering the count solo).
  Proof: cmd/libretto/wiki_test.go TestWikiHomeCarriesTheInFlightStrip
- **Where a change's `proposal.md` holds an open box under `## Checklist`, the in-flight
  strip shall name that change with its closed/total count read from that section**,
  including when the folder also holds a `tasks.md` with different counts and `proposal.md`
  holds boxes outside the section.
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightReadsTheChangeDocument
- **Where a change's `proposal.md` carries a `Queued:` line and no box under `## Checklist`,
  and no unmerged local branch carries its name, the wiki shall count it in the queue and
  shall not name it in the in-flight strip.**
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightReadsTheChangeDocument
- **The in-flight strip shall count boxes with the loop's definition: a `* [ ]` line and a
  `- [X]` line in a `tasks.md` shall count as one open box and one closed box.** It matched
  `- [ ]` and `- [x]` by prefix while its comment claimed the loop's definition.
  Proof: cmd/libretto/wiki_test.go TestWikiInFlightCountsBoxesAsTheLoopDoes
- **The home shall carry a segmented health bar whose green width is the integer
  division of `shall`-carrying criteria over the total and whose amber width is
  the remainder to 100** — the bar always closes — with the count of criteria
  whose cited proof does not resolve (file exists; a Go citation also declares
  the named test).
  Proof: cmd/libretto/wiki_test.go TestWikiHomeMeasuresContractHealth
- **Each capability card shall carry a two-state health marker**: green where
  every criterion carries `shall` and every cited proof resolves, amber
  otherwise. No red — the wiki reports, gates block.
  Proof: cmd/libretto/wiki_test.go TestWikiCardsCarryTheHealthDot
- **The home footer shall state how many tracked files some `Governs:` glob
  claims and how many none does**, over `git ls-files` behind a seam; where git
  is unavailable the footer shall be absent. Glob matching mirrors spec-drift's
  bash-case semantics — `*` crosses directory separators — because two matching
  authorities would disagree exactly when it matters.
  Proof: cmd/libretto/wiki_test.go TestWikiFooterMeasuresGovernedTree

The capability page carries its own history and its proof, from the same
sources, absent when the data is.

- **Where the dates seam yields commit dates for a spec, the page shall carry a
  sparkline of monthly commit counts** — eight buckets ending at the most recent
  commit's month, never the clock's, bar heights proportional — and shall omit
  it when the seam yields nothing.
  Proof: cmd/libretto/wiki_test.go TestWikiPageCarriesTheActivitySparkline
- **Each criterion on a page shall carry a chip naming its cited proof file
  and, when one is named, its test** — green where the criterion carries
  `shall` and its proof resolves, amber otherwise, one definition of resolving
  shared with the home board.
  Proof: cmd/libretto/wiki_test.go TestWikiCriteriaCarryProofChips
- **Where a spec holds a `Prior decisions` section, the page shall render up to
  its first three bullets with the section's total count, above the criteria**;
  where it holds none, the block shall be absent.
  Proof: cmd/libretto/wiki_test.go TestWikiPageSurfacesPriorDecisions
- **Where a spec file's raw text mentions another capability's name on a
  case-insensitive word boundary, the page shall carry a Related row linking
  each mentioned capability once, never itself**; with no mentions the row
  shall be absent. (The match rule is assumed 2026-08-19 under attacca; smarter
  linking waits for a real false positive.)
  Proof: cmd/libretto/wiki_test.go TestWikiPageLinksRelatedCapabilities

The palette searches the whole contract, offline, from one block.

- **The generated page shall carry one `application/json` index block that
  parses as JSON and holds three groups — criteria, decisions, capabilities —
  each entry naming its owning capability and its text**, built from the same
  extraction the pages render.
  Proof: cmd/libretto/wiki_test.go TestWikiCarriesTheSearchIndex
- **The page shall carry the palette overlay hidden by default and an inline
  script that reads the index block, listens on `keydown`, and navigates by
  setting the location hash** — no external script, classes over style
  juggling; opened by shortcut or by its own button, never by a click on the
  card filter, which keeps its own gesture.
  Proof: cmd/libretto/wiki_test.go TestWikiPaletteIsInlineAndDormant
- **Spec content in the index shall arrive JSON-encoded such that a criterion
  containing a closing script tag cannot terminate the block or execute.**
  Proof: cmd/libretto/wiki_test.go TestWikiSearchIndexEscapesContent

The flow itself gets a page, from the ledgers, absent when they are.

- **Where the project holds a lessons ledger with `## <date> · <change> ·
  <phase>` entry headers, the page shall carry a flow article at `#flow`** with
  one labeled bar per phase, widths proportional to the counts — the phase read
  from the last ` · ` field, because a change name may carry the separator;
  where the ledger is absent or empty of valid entries, the article shall be
  absent.
  Proof: cmd/libretto/wiki_test.go TestWikiFlowBoardCountsCorrections
- **Where queued proposals exist, the flow article shall list each with its
  `Queued:` date, oldest first**; with none the queue block shall be absent —
  and a queue alone shall still summon the article, which is where the queue
  lives on its own.
  Proof: cmd/libretto/wiki_test.go TestWikiFlowBoardListsTheQueue
- **If an unmerged local branch carries a queued change's name, then the wiki shall
  neither list that change in the flow board's queue nor count it in the in-flight strip.**
  Proof: cmd/libretto/wiki_test.go TestWikiQueueDropsAChangeItsBranchAlreadyTook
- **The wiki's flow-board queue shall list the same change names, in the same order, as
  `queuedIdeas` returns for the same changes directory and branch list.**
  Proof: cmd/libretto/wiki_test.go TestWikiQueueMatchesTheCommand
- **Where the flow article exists, the home shall link `#flow`; where it does
  not, no such link.**
  Proof: cmd/libretto/wiki_test.go TestWikiHomeLinksTheFlowBoard

#### The title names the repository

The title on both views — `README.md`'s heading, `wiki.html`'s `<title>` and its home
`<h1>` — is the repository's name: the parent of git's common dir when that dir is `.git`,
so every linked worktree of one repository writes the same title.

- **When** `libretto wiki` runs in a linked worktree, the `README.md` heading **shall** name
  the main checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiFromAWorktreeNamesTheRepository
- **When** `libretto wiki --html` runs in a linked worktree, `wiki.html`'s `<title>` and its
  home page's `<h1>` **shall** both name the main checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiHTMLFromAWorktreeNamesTheRepository
- **When** `libretto wiki` runs in the main checkout, the `README.md` heading **shall** name
  that checkout's directory.
  Proof: cmd/libretto/wiki_test.go TestWikiFromTheMainCheckoutNamesItsDirectory
- **If** git cannot report a common dir, **then** the title **shall** be the run directory's
  name.
  Proof: cmd/libretto/wiki_test.go TestWikiOutsideGitKeepsTheDirectoryName
- **If** the common dir git reports is not named `.git`, **then** the title **shall** be the
  run directory's name.
  Proof: cmd/libretto/wiki_test.go TestWikiProjectNameFallsBackOffAPlainGitDir

### `land` — the landing commit verified before it exists

A landing is one commit doing four things — final code, delta applied onto each
`Targets:` capability spec, durable decisions retired, change folder deleted — and it
failed twice by half-finishing silently. `land [<change>]` inspects `git diff --cached`,
before the mistake is history, and exits zero only when the two parts a binary can check
both hold: every tracked file under the change folder a staged deletion with nothing
under the folder surviving the commit (part 4 — a tracked file not staged for deletion
and an untracked file still on disk both fail, by name), and each delta's `Targets:`
capability spec added or modified in the same staged diff (part 2). It verifies; it
never performs — applying a delta is a semantic merge a binary cannot do, and the
observed failures were incomplete landings, which a verifier catches.

**Part 3 is deliberately not checked.** `spec-drift --retired`, inside `--anchors`, owns
it, and a second implementation is two sources of truth about what "retired" means. The
failure report says so — it attributes part 3 to `spec-drift --anchors` rather than
staying silent about a part it does not verify, so a green `land` is never read as the
whole contract passing.

**The deltas are read from `HEAD`, not from disk or the index.** A delta is any `.md` in
the folder carrying `Targets:` — `proposal.md` and each `spec-<capability>.md` in the
change-document layout, the old layout's delta files otherwise; a capture with no
`Targets:` contributes no capability. A staged deletion
removes the file from the index and the landing may already have removed it from the
working tree, so `HEAD` is the one place the `Targets:` lines are guaranteed to still
be. `Targets:` lines inside fenced code blocks do not count, matching `spec-drift`'s
own fence-stripping — the two tools must not disagree about what a delta targets.

**The change name is optional.** With none given, the change is inferred from staged
deletions under a change root; with more than one folder landing, each is verified —
each has its own contract, and refusing an unusual-but-legal commit helps nobody. With
none inferable, or with a named change that has no staged deletions, the exit is
non-zero saying nothing is landing: a verifier that exits zero having verified nothing
is the silent half-landing wearing a green light. A folder whose deletion carries no
delta passes part 2 vacuously, and the report says so — deleting a queued proposal is
abandoning an idea, which the flow says costs nothing; it is not a landing and must not
be failed as a broken one.

**Discovery mirrors what already exists, on both axes.** Change roots are the three
`spec-drift` reads — `.agents/changes`, `changes`, `openspec/changes` — and the specs
directory is `findSpecsDir`'s order (`.agents/specs`, `specs`, `openspec`, `docs/specs`,
`spec`). A fourth list in the binary would disagree with the others exactly when it
matters.

**A stale wiki warns and never blocks.** Stale, operationally: a marked view
(`README.md` under the wiki marker, or `wiki.html` under the HTML marker) exists in the
specs directory, a capability `spec.md` is in the staged diff, and that view is not.
That is the mechanical reading of record-work's own clause — the refreshed index rides
the same commit as the delta that changed it. The warning goes to stderr and the exit
code is untouched; a foreign view (no marker) is ignored silently, matching `wiki`'s
own precedent. No marked view present means nothing to say.

**Every failure names its part.** The non-zero exit carries one line per missing part —
which file survived, which capability spec did not move — so the fix is legible from
the output. All failures are reported in one run; stopping at the first would make the
repair iterative for no reason. Remedy lines carry `invokedAs()`, as everywhere else.

These run against real temporary git repositories built per case — an index is a git
fact, and a faked one proves nothing about the command that reads the real thing.

- **When the staged index deletes every file of the change folder and modifies each
  delta's `Targets:` capability spec, `land` shall exit zero.**
  Proof: cmd/libretto/land_test.go TestLandPassesACompleteLanding
- **If a tracked file under the change folder is not a staged deletion, then `land`
  shall exit non-zero and name that file under part 4.**
  Proof: cmd/libretto/land_test.go TestLandFailsAPartialFolderDeletion
- **If an untracked file remains on disk under the change folder, then `land` shall
  exit non-zero and name it** — the commit cannot delete a folder it leaves a file in.
  Proof: cmd/libretto/land_test.go TestLandFailsAnUntrackedLeftover
- **If a delta's `Targets:` capability spec is neither added nor modified in the staged
  diff, then `land` shall exit non-zero and name the capability under part 2.**
  Proof: cmd/libretto/land_test.go TestLandFailsWhenTheCapabilitySpecDidNotMove
- **When both checked parts are missing, `land` shall name both in one run**, not stop
  at the first.
  Proof: cmd/libretto/land_test.go TestLandNamesEveryMissingPart
- **`land` shall read `Targets:` lines from `HEAD`**, so a delta already deleted from
  the working tree is still verified, **and shall ignore a `Targets:` inside a fenced
  code block**, matching spec-drift's fence-stripping.
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromHeadAndSkipsFences
- **When several delta files in one folder carry `Targets:` lines, `land` shall check
  every named capability**, failing on any one whose spec did not move.
  Proof: cmd/libretto/land_test.go TestLandChecksEveryDeltasTarget
- **When a change's `proposal.md` at `HEAD` carries `Targets: <cap>` outside fences, `land`
  shall treat `<cap>` as one of the change's capabilities**, including when the index has
  deleted the change folder.
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromTheChangeDocument
- **When a change folder at `HEAD` holds a `spec-<cap>.md` carrying `Targets: <cap>`, `land`
  shall treat `<cap>` as one of the change's capabilities** alongside those named by
  `proposal.md`.
  Proof: cmd/libretto/land_test.go TestLandReadsTargetsFromAnExtraCapabilitySpec
- **If a change's `proposal.md` carries no `Targets:` line outside fences, then `land` shall
  derive no capability from it, and `land` on an old-layout change shall find the same
  capabilities it found before the change document existed.** Forced red by excluding
  `proposal.md` and `spec-*.md` from `headTargets`, then reverted — all three were green on
  their first run, which is not evidence.
  Proof: cmd/libretto/land_test.go TestLandOldLayoutTargetsAreUnchanged
- **When no change name is given, `land` shall infer the change from staged deletions
  under a change root**, and shall verify each folder when more than one is landing.
  Proof: cmd/libretto/land_test.go TestLandInfersTheChangeFromStagedDeletions
  Proof: cmd/libretto/land_test.go TestLandVerifiesEveryStagedChangeFolder
- **If no staged deletion lies under any change root and no name was given, then `land`
  shall exit non-zero saying nothing is landing**; a named change with no staged
  deletions shall be refused the same way.
  Proof: cmd/libretto/land_test.go TestLandWithNothingStagedRefuses
  Proof: cmd/libretto/land_test.go TestLandRefusesANamedChangeWithNothingStaged
- **When a change folder's deletion carries no delta, `land` shall pass part 2
  vacuously and say so** — an abandoned queued proposal is not a broken landing.
  Proof: cmd/libretto/land_test.go TestLandAllowsADeltalessFolderDeletion
- **`land` shall discover change folders under `.agents/changes`, `changes` and
  `openspec/changes`** — the same three roots spec-drift reads.
  Proof: cmd/libretto/land_test.go TestLandDiscoversEveryChangeRoot
- **`land` shall not check part 3, and its report shall attribute it to
  `spec-drift --anchors`** — a landing that would fail `--retired` still passes `land`.
  Proof: cmd/libretto/land_test.go TestLandLeavesPartThreeToSpecDrift
- **When a marked wiki view exists and a capability spec is staged while that view is
  not, `land` shall warn on stderr and shall not change the exit code.**
  Proof: cmd/libretto/land_test.go TestLandWarnsOnAStaleWikiWithoutBlocking
- **When no marked view exists, or the view rides the same staged diff, `land` shall
  print no warning**, and a foreign unmarked view shall be ignored silently.
  Proof: cmd/libretto/land_test.go TestLandStaysSilentWhenTheWikiIsCurrentOrForeign
- **`land` shall write nothing** — the repository's files, index and refs are
  byte-identical before and after a run, on the passing path and on every failing one.
  Proof: cmd/libretto/land_test.go TestLandChangesNothing
- **`land` shall run with no payload installed** — it reads the project, not the tree.
  Proof: cmd/libretto/land_test.go TestLandWorksWithNoPayload
- **If the working directory is not inside a git repository, then `land` shall exit
  non-zero with an error naming git**, never an empty report.
  Proof: cmd/libretto/land_test.go TestLandOutsideARepositoryFails
- **`help` shall name `land`**, or the feature was not delivered.
  Proof: cmd/libretto/main_test.go TestHelpNamesLand

### `init` — bootstrap the `.agents/` layout

`init` is the explicit, user-invoked owner of bootstrap: the skills refuse to create a
project's `changes/` directory implicitly, and point here instead (see `payload`). Without
it the layout appears only as a side effect of the first change, and git drops the empty
`changes/` on every fresh checkout — which is how a queued idea was lost.

Run inside a git repository, it ends with three artifacts present at the **repository
root** — never the cwd, the bug `metrics` records having had:

| Artifact | Content |
|---|---|
| `.agents/specs/.gitkeep` | empty — a keep-file so git carries the directory |
| `.agents/changes/.gitkeep` | empty — the same, and the fix for the vanishing queue |
| `.agents/lessons.md` | the standard append-only header: a `# Lessons` heading and the paragraph saying entries are spent by `libretto-retro`, gain a `Resolved:` line there, and are never edited otherwise |

Only what is missing is created; everything already present is left byte-identical. The
report names the root it acted on, then one line per artifact saying `created` or
`exists`, plain text with no escape codes. A second run reports all three as present and
exits zero.

**Two refusals**, each an error on stderr with exit 1 and nothing created:

- **outside a git repository, or with git unavailable** — the `metrics` shape, an error
  naming git. The keep-files exist to survive git; a directory nothing tracks has nothing
  to survive.
- **a foreign layout** — any other consolidation root spec-drift reads, present as a
  directory: `specs`, `openspec`, `docs/specs`, `spec`, `changes`, `openspec/changes`. The
  refusal names what it found, and it wins over a partial `.agents/`: filling one beside
  another convention is the two sources of truth the refusal exists to prevent.

`init` reads the project, never the payload and never an install target — no
`CLAUDE_HOME`, no `~/.claude` — so it runs on a machine with nothing installed.

- **When run in a git repository with no `.agents/` layout, `init` shall create
  `.agents/specs/.gitkeep`, `.agents/changes/.gitkeep` and `.agents/lessons.md` at
  the repository root, report each as created, and exit zero.**
  Proof: cmd/libretto/init_test.go TestInitBootstrapsAFreshProject
- **When run a second time over its own output, `init` shall create nothing, modify
  nothing, report all three artifacts as already present, and exit zero.**
  Proof: cmd/libretto/init_test.go TestInitIsIdempotent
- **If a partial `.agents/` layout exists, then `init` shall create only the missing
  artifacts and shall leave every existing file byte-identical** — an existing
  `lessons.md` with entries in it is the case that matters.
  Proof: cmd/libretto/init_test.go TestInitFillsAPartialLayoutWithoutTouchingWhatExists
- **If a foreign consolidation root is present (`specs`, `openspec`, `docs/specs`,
  `spec`, `changes`, or `openspec/changes` as a directory), then `init` shall refuse,
  name what it found, create nothing, and exit non-zero.**
  Proof: cmd/libretto/init_test.go TestInitRefusesAForeignLayoutAndNamesIt
- **If a foreign root and a partial `.agents/` are both present, then `init` shall
  still refuse and create nothing** — the refusal wins over the fill.
  Proof: cmd/libretto/init_test.go TestInitRefusalWinsOverAPartialLayout
- **If the directory is not inside a git repository or git is unavailable, then
  `init` shall refuse with an error naming git and shall create nothing.**
  Proof: cmd/libretto/init_test.go TestInitRefusesOutsideAGitRepository
- **When run from a subdirectory of a repository, `init` shall act on the repository
  root, never the cwd.**
  Proof: cmd/libretto/init_test.go TestInitActsOnTheRepositoryRootNotTheCwd
- **`init` shall run exactly one git invocation, `rev-parse --show-toplevel`, and
  shall never commit** — asserted through the fake `gitRunner`'s call log.
  Proof: cmd/libretto/init_test.go TestInitNeverCommitsAndRunsOnlyRevParse
- **A created `lessons.md` shall carry the standard append-only header** — the
  `# Lessons` heading and the append-only paragraph naming `libretto-retro` and
  `Resolved:`.
  Proof: cmd/libretto/init_test.go TestInitLessonsHeaderIsTheStandardOne
- **`init` shall write nothing under `.agents/specs/` beyond the keep-file** — in
  particular no `README.md`, which `libretto wiki` owns by marker.
  Proof: cmd/libretto/init_test.go TestInitWritesNoReadmeIntoSpecs
- **`init` shall work with no payload installed** — the `needsPayload` exemption,
  proven the way `models` proves its own.
  Proof: cmd/libretto/init_test.go TestInitWorksWithNoPayload
- **The report shall name the repository root it acted on before the per-artifact
  lines, and shall carry no escape codes when piped.**
  Proof: cmd/libretto/init_test.go TestInitNamesTheRootAndCarriesNoEscapeCodes
- **If an artifact path is a symlink, dangling or not, then `init` shall report it as
  present and shall write nothing through it.**
  Proof: cmd/libretto/init_test.go TestInitNeverWritesThroughASymlink
- **`help` shall list `init` with a line saying what it bootstraps.**
  Proof: cmd/libretto/main_test.go TestHelpNamesInit

### `queue` — the queued ideas, read-only

`libretto queue` prints the project's queued ideas, one per line — the `Queued:` date, the
change name, and the first sentence of its `## Reading` — oldest date first, ties by name.
Queued is find-work's definition, and the binary implements it rather than owning it
(see `payload`): a `Queued:` line at column 0, no unmerged local branch carrying the name,
and no box under the proposal's `## Checklist`. A change an unmerged branch carries is
left out and named, with that branch, on a closing line on stdout — an idea that vanishes
without a word reads as lost. The wiki's in-flight strip and flow-board queue read the
same scan, so the two always list the same changes.

- **When `libretto queue` runs, it shall print one line per queued change carrying its
  `Queued:` date, its name and the first sentence of its `## Reading`, ordered by date,
  oldest first.**
  Proof: cmd/libretto/queue_test.go TestQueueListsOldestFirst
- **Where two queued changes carry the same `Queued:` date, the command shall order them
  by change name, ascending.**
  Proof: cmd/libretto/queue_test.go TestQueueBreaksDateTiesByName
- **Where a queued proposal has no `## Reading`, its line shall carry the date and the name
  only.**
  Proof: cmd/libretto/queue_test.go TestQueueWithoutAReadingPrintsTheNameAlone
- **If a proposal mentions `Queued:` anywhere but at the start of a line, then the command
  shall not list that change.**
  Proof: cmd/libretto/queue_test.go TestQueueReadsOnlyALineStartingWithQueued
- **If a proposal carrying a `Queued:` line holds a box under its `## Checklist`, then the
  command shall not list that change.**
  Proof: cmd/libretto/queue_test.go TestQueueSkipsAProposalWithAChecklistBox
- **If an unmerged local branch's name, with one leading `feat/`, `fix/`, `docs/`,
  `chore/` or `refactor/` removed, equals a queued change's name, then the command shall
  leave that change out of the list and shall name it, with the branch, on a closing
  line written to stdout.**
  Proof: cmd/libretto/queue_test.go TestQueueDropsAChangeItsBranchAlreadyTook
- **If a branch matches a queued change's name only in part, or only after removing a
  prefix outside that list, then the command shall keep the change in the list.**
  Proof: cmd/libretto/queue_test.go TestQueueKeepsAChangeOnAPartialOrForeignMatch
- **The command shall ask git for branches with exactly
  `branch --format=%(refname:short) --no-merged HEAD`.**
  Proof: cmd/libretto/queue_test.go TestQueueAsksGitForUnmergedLocalBranches
- **If git cannot list branches, then the command shall list every queued change
  unfiltered, shall say so in one line on stderr, and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueUnreadableBranchesListsUnfiltered
- **When run from a subdirectory of a repository, the command shall read the
  `.agents/changes/` at the repository root.**
  Proof: cmd/libretto/queue_test.go TestQueueReadsTheRepositoryRoot
- **While no `.agents/changes/` exists at the repository root, the command shall print
  `no .agents/changes/ here — nothing is queued` and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueWithoutAChangesDirSaysSo
- **While `.agents/changes/` exists and no change in it is queued after the branch filter,
  the command shall print `nothing is queued` and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueEmptyIsAStateNotAnError
- **If every queued change is dispatched by a branch, then the command shall print
  `nothing is queued` followed by the closing line naming each dispatched change with its
  branch, both on stdout.**
  Proof: cmd/libretto/queue_test.go TestQueueAllDispatchedSaysNothingQueuedThenNamesThem
- **The command shall leave every file under the project byte-identical.**
  Proof: cmd/libretto/queue_test.go TestQueueWritesNothing
- **If the command is given any argument other than `--json`, then it shall exit non-zero
  naming the argument, and shall print no queue.**
  Proof: cmd/libretto/queue_test.go TestQueueRefusesArguments
- **The missing-payload stop shall not apply to `queue`.**
  Proof: cmd/libretto/queue_test.go TestQueueNeedsNoPayload
- **The help text shall carry a `queue` line saying it lists queued ideas, read-only.**
  Proof: cmd/libretto/main_test.go TestHelpNamesQueue

### `--json` — the read-only reports, machine-readable

`libretto status --json`, `doctor --json`, `metrics --json` and `queue --json` each write
one JSON object to stdout and nothing else, so a script reads the same facts the text
report shows a human. Each object carries `"schema": 1`. Errors, warnings and the release
notice stay on stderr exactly as they are, and exit codes do not move — `doctor` still
exits non-zero when something needs attention, after its complete report. Without
`--json` every command's output is unchanged.

The flag goes after the command name, anywhere among its arguments; `libretto --json
status` is an unknown command. `status` and `doctor` ignore every other argument, as they
always have; `metrics` takes it beside a `<change>` in either order; `queue` refuses
anything else. The help text names the flag and the four commands.

**The four shapes are the contract.** Every key below is a promise. A key marked *omitted
when* is absent — never `0`, `""` or `null` — in that state: an unknown value is an absent
key, a known zero is present, a known-empty list is `[]`. Keys are snake_case; the
document is indented two spaces with one trailing newline, and whitespace is not part of
the contract; HTML escaping is off, so `<`, `>` and `&` arrive as written; no escape byte
appears, TTY or not. Adding a key never bumps `schema`; removing, renaming or retyping one
does, and is a promise removed. The JSON carries the report's data, never its prose, and
nothing a consumer can compute from the rest.

#### `status --json`

    {
      "schema": 1,
      "target": { "name": "claude", "root": "/Users/u/.claude", "configured": true },
      "items": [
        { "kind": "skills", "name": "evidence", "state": "linked" },
        { "kind": "agents", "name": "spec-reviewer", "state": "wrong target" }
      ]
    }

| Key | Meaning |
|---|---|
| `target.name` | the target's name as the text header prints it: `claude`, `project`, `codex`, `opencode` |
| `target.root` | its root directory; *omitted when* the target has none (a project target with no directory) |
| `target.configured` | whether the root exists on disk; the text prints `not configured` for `false` |
| `items` | every scanned entry in the scan's order; `[]` when there is nothing to link |
| `kind` | `skills`, `agents` or `commands` |
| `state` | `linked`, `missing`, `wrong target`, `conflict` or `stale`, spelled as `link.State` spells them |

#### `doctor --json`

    {
      "schema": 1,
      "target": { "name": "claude", "root": "/Users/u/.claude", "configured": true },
      "problems": [
        { "kind": "skills", "name": "alpha", "state": "missing" }
      ],
      "release": {
        "mode": "installed",
        "payload_root": "/Users/u/.local/share/libretto/current",
        "running": "v0.5.2",
        "latest": "v0.6.0",
        "newer": true
      },
      "prerequisites": [
        { "name": "jira", "found": false },
        { "name": "rg", "found": true }
      ]
    }

| Key | Meaning |
|---|---|
| `target` | the same object as `status` |
| `problems` | every entry whose state needs attention, same item shape as `status`; `[]` when the text says `no problems` |
| `release.mode` | `checkout` or `installed` |
| `release.payload_root` | the payload root the text names after the mode |
| `release.running` | the version the binary reports, always present |
| `release.latest` | the newest release the check found; *omitted when* the check failed or the remote has no releases |
| `release.newer` | `true` when `latest` ranks above `running`, `false` when they are equal; *omitted when* there is no `latest` or the two do not rank (a `dev` build, or a running version ahead of the remote) |
| `release.error` | the failed check's first line; *omitted when* the check answered |
| `prerequisites` | one object per prerequisite in the text's order, `name` and `found` |

#### `metrics --json`

    {
      "schema": 1,
      "changes": [
        {
          "name": "add-converge-command", "state": "landed",
          "commits": 9, "first": "2026-10-21T18:02:11Z", "last": "2026-10-22T09:40:03Z",
          "closed": 3, "boxes": 3, "reopened": 0, "corrections": 1
        },
        { "name": "broken-one", "state": "unreadable", "corrections": 0 }
      ],
      "corrections": { "by_phase": { "build-and-check": 4, "6→7": 2 }, "outside_changes": 1 },
      "tokens": {
        "attributed":   { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
        "unattributed": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
        "entries": 8944, "unattributed_entries": 312,
        "change": {
          "name": "add-converge-command",
          "totals":   { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
          "by_phase": { "write-spec": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 } },
          "no_phase": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 }
        }
      },
      "loop": {
        "changes": [
          { "name": "add-converge-command", "rounds": 3, "net_closed": 3, "cost_usd": 4.12, "unpriced": 0 }
        ],
        "skipped": 0
      }
    }

| Key | Meaning |
|---|---|
| `changes` | one object per row of the table, in its order, narrowed by a `<change>` exactly as the table is; `[]` when git has seen no change — and then `corrections`, `tokens` and `loop` are still present whenever their sources exist, while the text prints only its no-changes line |
| `name`, `state` | always; `state` is `landed`, `in flight` or `unreadable`, the table's spelling |
| `commits`, `first`, `last` | *omitted when* `unreadable`. Times are RFC 3339 in UTC, from the commit times `measure` reads; the span is `last − first` and is not repeated |
| `closed`, `boxes`, `reopened` | boxes closed now, the denominator, and reopened; *omitted when* the change has no checklist in its history (the text's `—`) or is `unreadable` |
| `corrections` on a change | *omitted when* there is no `.agents/lessons.md`; `0` included when there is |
| top-level `corrections` | *omitted when* there is no ledger. `by_phase` keyed as the ledger spells each phase, `{}` for an empty ledger; `outside_changes` is the orphan count |
| `tokens` | *omitted when* no transcripts are configured or none exist for this repository. `attributed`, `unattributed`, `entries` and `unattributed_entries` are corpus-wide whatever the filter, as the text's are |
| `tokens.change` | *omitted when* no `<change>` was given. `totals` and `by_phase` *omitted when* no session was attributed to that change, leaving `name` alone; `no_phase` *omitted when* every attributed entry named a phase |
| `loop` | *omitted when* no run-log file exists, and only then |
| `loop.changes` | per change in name order, following the filter; `[]` when the log holds none for the filter, holds no record, or holds only malformed lines. `cost_usd` *omitted when* none of that change's rounds carried a cost |
| `loop.skipped` | malformed lines, `0` included |
| `loop.error` | present, in place of `changes` and `skipped`, only when the log exists and cannot be read |

The legend and the ceiling are prose for a human reading the table and are not in the
JSON; this table is their machine-readable counterpart.

#### `queue --json`

    {
      "schema": 1,
      "queued": [
        { "name": "report-drift-trend-in-metrics", "date": "2026-08-21",
          "reading": "Derive per-capability drift events from git history." }
      ],
      "dispatched": [
        { "name": "add-json-output", "branch": "feat/add-json-output" }
      ]
    }

| Key | Meaning |
|---|---|
| `queued` | oldest date first, ties by name, as the text; `[]` when nothing is queued or no `.agents/changes/` exists |
| `reading` | *omitted when* the proposal has no `## Reading` |
| `dispatched` | changes an unmerged branch already took, `[]` when none; *omitted when* git could not list branches, or when no `.agents/changes/` exists and git was never asked — either way nobody knows which are dispatched |

#### Criteria

The shared criteria run one row per command.

- **When `status`, `doctor`, `metrics` or `queue` runs with `--json`, the command shall
  write exactly one JSON object to stdout followed by one newline.**
  Proof: cmd/libretto/json_test.go TestJSONIsOneObjectOnStdout
- **The JSON object each command writes shall carry a top-level `schema` key whose value
  is the integer 1.**
  Proof: cmd/libretto/json_test.go TestJSONCarriesSchemaOne
- **Where `--json` is given, the command shall exit with the same code it exits with on
  the same state without `--json`.**
  Proof: cmd/libretto/json_test.go TestJSONDoesNotChangeTheExitCode
- **If a command run with `--json` fails before its report is complete, then it shall
  write nothing to stdout.**
  Proof: cmd/libretto/json_test.go TestJSONFailureWritesNothingToStdout
- **If a command fails, then the error it returns under `--json` shall equal the error it
  returns without `--json`.**
  Proof: cmd/libretto/json_test.go TestJSONErrorsMatchTheTextErrors
- **The JSON a command writes shall contain no escape byte (0x1b).**
  Proof: cmd/libretto/json_test.go TestJSONHasNoEscapeCodes
- **While a newer release is cached, the stdout of `status --json` shall parse as one JSON
  object.**
  Proof: cmd/libretto/notice_test.go TestJSONStdoutParsesWithANoticePending
- **The help text shall carry a `--json` line naming `status`, `doctor`, `metrics` and
  `queue`.**
  Proof: cmd/libretto/main_test.go TestHelpNamesJSON

`status`

- **When `status --json` runs, the object shall carry `target` with the target's `name`,
  its `root` and `configured`.**
  Proof: cmd/libretto/json_test.go TestStatusJSONNamesTheTarget
- **If the target has no root, then `target` shall omit `root`.**
  Proof: cmd/libretto/json_test.go TestStatusJSONOmitsAnUnknownRoot
- **When `status --json` runs, `items` shall list one object per scanned entry, carrying
  `kind`, `name` and `state`, in the scan's order.**
  Proof: cmd/libretto/json_test.go TestStatusJSONListsEveryItem
- **The `state` of every item shall be spelled exactly as the text report spells it.**
  Proof: cmd/libretto/json_test.go TestStatusJSONStatesMatchTheText
- **While the payload holds nothing to link, `items` shall be an empty array.**
  Proof: cmd/libretto/json_test.go TestStatusJSONEmptyIsAnEmptyArray

`doctor`

- **When `doctor --json` runs, `problems` shall list exactly the entries whose state needs
  attention, each carrying `kind`, `name` and `state`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONListsOnlyProblems
- **When `doctor --json` runs, `release.mode` shall be `checkout` for a payload root
  holding `.git` and `installed` otherwise.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONNamesTheMode
- **The `release` object shall carry `running`, the version the binary reports.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONCarriesTheRunningVersion
- **When the release check answers with a release, `release.latest` shall carry it.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONCarriesTheLatestRelease
- **Where `latest` and `running` rank, `release.newer` shall be true when `latest` ranks
  above `running` and false when they are equal.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONRanksWhatRanks
- **If `latest` is absent or the two versions do not rank, then `release` shall omit
  `newer`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONOmitsAnUnrankableComparison
- **If the release check fails, then `release` shall carry `error`, the failure's first
  line, in place of `latest`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONReportsAFailedCheck
- **When `doctor --json` runs, `prerequisites` shall list one object per prerequisite,
  carrying `name` and `found`, in the text's order.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONListsPrerequisites
- **When one stubbed release check feeds both outputs, the release text line and the
  `release` object shall agree: for latest `v9.9.9` against running `v0.1.0` the line
  names `v9.9.9` while the object carries `latest` `v9.9.9` and `newer` true, and for an
  unrankable pair the line carries no newer-release wording while the object omits
  `newer`.**
  Proof: cmd/libretto/doctor_test.go TestReleaseLineIsRenderedFromTheCheck

`metrics`

- **When `metrics --json` runs, `changes` shall list one object per row of the table, in
  the table's order.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONListsTheTableRows
- **For every readable change, `commits` in the JSON shall equal the table's commits
  column on the same history.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAgreesWithTheTable
- **Every readable change object shall carry `first` and `last` as RFC 3339 times in
  UTC.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCarriesTimesInUTC
- **Where a change has a checklist in its history, its object shall carry `closed`,
  `boxes` and `reopened`, zero included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONKeepsAKnownZero
- **If a change has no checklist in its history, then its object shall omit `closed`,
  `boxes` and `reopened`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsChurnWithoutAChecklist
- **If a change's history cannot be read, then its object shall carry no key beyond
  `name`, `state` set to `unreadable`, and `corrections`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnreadableRowCarriesOnlyWhatIsKnown
- **While `.agents/lessons.md` exists, every change object shall carry `corrections`, zero
  included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCorrectionsZeroIsPresent
- **If no `.agents/lessons.md` exists, then the report shall carry no `corrections` key
  at any level.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAbsentLedgerOmitsCorrections
- **While the ledger exists, the top-level `corrections` shall carry `by_phase`, keyed by
  each phase as the ledger spells it, and `outside_changes`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCorrectionsByPhase
- **If no session transcripts are configured or found for the repository, then the
  report shall omit `tokens`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsTokensWithoutTranscripts
- **When transcripts are found, `tokens` shall carry `attributed`, `unattributed`,
  `entries` and `unattributed_entries` computed over every change git has seen, whatever
  the filter.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONTokensAreCorpusWide
- **Where a `<change>` filter names a change sessions were attributed to, `tokens.change`
  shall carry its `name`, `totals` and `by_phase`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONTokensForOneChange
- **If the filtered change has no attributed session, then `tokens.change` shall carry its
  `name` alone.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnattributedChangeCarriesItsNameAlone
- **Where an attributed entry of the filtered change names no phase, `tokens.change`
  shall carry those entries' sum as `no_phase`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONNoPhaseIsNeverDistributed
- **If no run-log file exists, then the report shall omit `loop`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsLoopWithoutARunLog
- **If the run-log file exists and holds no loop record, then `loop` shall be
  `{changes: [], skipped: N}`, with N the count of its malformed lines.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAnEmptyRunLogIsKnown
- **When the run log exists, `loop.changes` shall list one object per change it records,
  following the filter, in name order, carrying `name`, `rounds`, `net_closed` and
  `unpriced`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONLoopRows
- **If none of a change's recorded rounds carried a cost, then its loop object shall omit
  `cost_usd`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnpricedOmitsCost
- **When the run log exists, `loop` shall carry `skipped`, the count of lines that are not
  loop records, zero included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCountsSkippedLines
- **If the run log exists and cannot be read, then `loop` shall carry `error` in place of
  `changes` and `skipped`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnreadableRunLogNamesTheError
- **While git has seen no change, `changes` shall be an empty array.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONEmptyHistoryIsAnEmptyArray
- **While git has seen no change, the JSON report shall still carry each of
  `corrections`, `tokens` and `loop` whose source exists.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONEmptyHistoryKeepsTheOtherBlocks
- **While git has seen no change, the text report shall print only its no-changes line.**
  Proof: cmd/libretto/metrics_test.go TestMetricsTextWithNoHistoryIsUnchanged
- **When `metrics` is given `--json` and a change name, the command shall accept them in
  either order.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONFlagAndChangeInEitherOrder

`queue` — its refusal of every argument other than `--json` is the criterion under
*`queue` — the queued ideas, read-only*, proven by `TestQueueRefusesArguments`.

- **When `queue --json` runs, `queued` shall list one object per queued change, carrying
  `name`, `date` and `reading`, in the text's order.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONListsOldestFirst
- **If a queued proposal has no `## Reading`, then its object shall omit `reading`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONOmitsAMissingReading
- **A `reading` carrying `<`, `>` or `&` shall reach the JSON unescaped.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONKeepsReadingTextVerbatim
- **When git lists branches, `dispatched` shall list each change a branch took, carrying
  `name` and `branch`, as an empty array when there is none.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONNamesDispatchedChanges
- **If git cannot list branches, then the report shall omit `dispatched`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONUnreadableBranchesOmitsDispatched
- **While no `.agents/changes/` exists at the repository root, `queued` shall be an empty
  array.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONWithoutAChangesDirIsEmpty
- **While no `.agents/changes/` exists at the repository root, the report shall omit
  `dispatched`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONWithoutAChangesDirIsEmpty
