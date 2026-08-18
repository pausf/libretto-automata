# AGENTS.md

Instructions for agents working **in this repository**. Auto-loaded by most agent
tooling. If you are working in someone else's project *using* this repo's payload,
you want the installed skills instead, not this file.

## What this is

Two things in one repository, and confusing them is the first mistake to avoid:

- **The payload** — an eight-phase spec-driven flow, shipped as skills, agents and
  commands. **This is the point of the project.**
- **The CLI** — a Go binary that symlinks the payload into `~/.claude`, one item at a
  time, never over anything it did not create. **This is delivery.**

A change that makes the CLI prettier and the payload no better is a change that missed
the point.

## Stack

| | |
|---|---|
| Go | 1.26.5 |
| TUI | bubbletea 1.3.10, lipgloss 1.1.0, termenv 0.16.0 |
| TTY | go-isatty 0.0.24, charmbracelet/x/term 0.2.2 |
| Tracker | `jira` CLI (ankitpokhrel/jira-cli), never MCP, never the REST API |
| Shell tooling | bash, `rg`, `fd`, `jq` |

No test framework beyond the standard library. No linter beyond `gofmt` and `go vet`.

## Commands

```bash
make build              # -ldflags stamps the version from git describe
make test               # go test ./...
make test-short         # skips real-git integration and teatest flows
make fmt vet
make preview            # print the panel once, colour forced
make link               # symlink `libretto` and `libretto-automata` into ~/.local/bin
make unlink
make clean
```

```bash
libretto                # the panel (needs a TTY)
libretto status         # every item's state, read-only
libretto install        # link everything; idempotent; non-zero if anything was skipped
libretto doctor         # what needs attention + what the payload expects on this machine
libretto prune          # links whose source is gone; change nothing
libretto prune --yes    # remove them
libretto uninstall      # what this repo installed here; change nothing
libretto uninstall --yes # take it back out
libretto preview
libretto models         # which model each agent runs on and how hard it thinks, read-only
libretto models set haiku review-lens-design review-lens-tests   # --all for every agent
libretto models effort xhigh review-lens-security   # the depth, not the tier. Haiku has none

libretto loop <change>  # one fresh session per open box in that change's plan.md,
                        # until every box is closed. Stops after two rounds that close
                        # nothing, and at --max (10). Never pushes.
libretto loop <change> --dry-run --max 3
libretto metrics        # what every change cost — commits, span, boxes closed and
                        # reopened — derived from git, instrumented nowhere. Read-only.
libretto metrics <change>
libretto wiki           # render this project's specs into <specs-dir>/README.md,
                        # behind a marker; refuses a README it did not generate

libretto update         # install the newest version and relink; pulls in a checkout
libretto install --project   # <cwd>/.claude instead of ~/.claude
libretto install --global    # the default; both flags at once is an error
```

## Gates — all six pass before any commit

```bash
gofmt -l .                                       # must print nothing
go vet ./...
go test ./... -count=1
scripts/check-payload                            # frontmatter, references, reachability
skills/record-work/spec-drift --self-test
skills/record-work/spec-drift --anchors          # Proof: resolves · delta criteria are EARS · no plan retired silently
```

`spec-drift` with no flag warns about staged code whose spec did not move. It never
blocks — run it, read it, answer the question out loud.

Never pipe a gate into `head` when the exit code matters: the pipeline reports the last
command's status, so a failure reads as success. Redirect to a file, check `$?`, read
the file.

## Structure

```
cmd/libretto/           the CLI, dispatch and the scope flags
internal/target/        what an installable destination is
internal/link/          own.go (ownership) · scan.go state.go (read) · plan.go apply.go (write)
internal/repo/          git, the rebuild decision, and the newer-version check
internal/ui/            logo, theme, panel, model, models.go (the selector screen)
internal/agentmodel/    the model: key in agents/*.md, and the catalogue of values
skills/ agents/ commands/   THE PAYLOAD — what gets symlinked
scripts/check-payload   repo-only tooling. Never referenced from a skill.
.agents/specs/          the specification, one directory per capability. docs/SPEC.md lists them.
docs/                   FLOW.md (the flow) · DESIGN.md · PLAN.md · STATE.md · SPEC.md (index only)
```

## Where work comes from

Not from a tracker. Phase 1 asks three sources in order:

1. **a change already in flight** — unchecked boxes in `.agents/changes/*/plan.md`
2. a tracker key or URL, if one was given
3. **what the user said** — a legitimate input, not a fallback

Home first. Starting new work while something sits half-finished is how the half-finished
thing gets abandoned, and the cost is a `.agents/changes/` nobody trusts.

Creating a task without a tracker *is* creating the change folder: `proposal.md` with
`Tracker: none` and the request in the words it was asked in, named verb-led —
`add-relative-discounts`, never an invented ticket id.

`/libretto-status` reports what is open and changes nothing.

An idea you do not want to build yet goes in the queue: `/libretto-queue` captures it as a
proposal carrying a `Queued:` date and nothing else, and `/libretto-next` picks the oldest
one up when you are ready. **Queued is not in flight** — a captured idea never blocks
source 1, because abandoning an idea costs nothing and a queue that is expensive to add to
is a queue nobody uses.

**Mark the boxes as you go.** A plan that is 11/11 open while the work is nearly done is
documentation pretending to be state — which happened while building this very feature.

## Specs

Per **capability**, never per ticket, in `.agents/specs/<capability>/spec.md`. `docs/SPEC.md`
is the index and has no requirements in it — **it is also the only place the list lives.**
This paragraph used to name them and count them, and by the time anybody noticed it said
"ten" over eleven directories. A number kept in two places is a number that drifts, and the
copy nobody edited is the one that reads as authoritative.

Each declares what it owns and cites the test behind each criterion:

```
Governs: internal/link/plan.go internal/link/apply.go internal/link/plan_test.go
Proof:   internal/link/apply_test.go TestApplyIsIdempotent
```

**Criteria are written in EARS**, and `--anchors` fails a change delta whose criterion
carries no `shall` — one of the five patterns in `skills/write-spec/`. A criterion in
prose cannot be failed, only interpreted, and one nobody can fail is treated as satisfied
by default. Capability specs written before 2026-08-17 are warned about and never failed:
545 criteria predate the syntax, and rewriting them in one unreviewable diff is 545
chances to change a promise that works today. A capability migrates when a delta lands on
it. **This is inside `--anchors` and not a seventh gate**, because the count "six" is
written in ten places here and a number kept in ten places is a number that drifts.

**A plan cannot be deleted taking its reasoning with it.** The landing commit deletes the
change folder, and `plan.md` holds why the change was built the way it was. `--anchors`
fails that commit unless some capability spec's *Prior decisions* section moved in it —
that section, not the file, because the delta lands in the same commit by definition and
"the spec was edited" would be green every time. A plan that retires nothing declares it:
`Durable decisions: none`, written when the plan was, not when the gate complains.

A `Proof:` citation must name a **test that exists**. `--anchors` checks the test name,
not just the file, because a file-level check passes an invented name — that happened
here, twice.

Work in flight goes in `.agents/changes/<change>/`, and lands by applying the delta onto
the capability spec and deleting the change folder, in the same commit as the final code.

## Versioning

Semver, from git tags. **The binary never carries a hardcoded version.**

```bash
make build                     # -ldflags stamps `git describe --tags --always --dirty`
./bin/libretto version         # v0.5.0  ·  v0.5.0-3-gabc123  ·  v0.5.0-dirty
```

**There is nothing to edit to change the reported version.** It is `git describe` at build
time, and `debug.ReadBuildInfo()` for a binary from `go install`. `v0.5.0-3-gabc123` is not a
bug to fix — it is a correct report of three commits past a tag. Tag, rebuild, and it says
`v0.5.0`. A constant in a source file desynchronises the moment somebody forgets it, and it
does it silently.

A plain `go build` with no ldflags reports `dev`. That is deliberate: a binary that
cannot prove its version says so rather than claiming one.

**Until a tag exists, `go install ...@latest` installs the last release and not `main`.**

### Every merge to `main` gets a tag. No merge leaves `main` untagged.

**This is the rule, and it is not optional.** The bump comes from what the change was, read
off the commits it lands:

| The merge contains | Bump | Example |
|---|---|---|
| only `fix:` / `refactor:` / `docs:` / `chore:` with no contract change | patch | `v0.5.0` → `v0.5.1` |
| a `feat:`, a new capability, or a new promise in an existing spec | minor | `v0.5.0` → `v0.6.0` |
| **a promise removed or reversed — while this is `0.x`** | **minor** | `v0.5.0` → `v0.6.0` |

Mixed merge takes the highest: one `feat:` among nine `fix:` is a minor.

### `1.0.0` is a decision somebody makes. It is never the output of this table.

**While the version starts with `0.`, nothing in this table produces a major**, and that is
semver rather than a local preference: in `0.x` anything may break, and the number that carries
a breaking change is the minor. `1.0.0` is a *declaration* — the CLI's flags, `install` and
`prune` behaviour, and the payload's skill contracts are stable and will not move without a
major. **Nobody has decided that here.** Until somebody does, the answer to "does this deserve
a major" is no.

Two more things follow, and the second is the expensive one:

- **A promise in `.agents/specs/` is not the same kind of promise as a promise to a user.**
  The table means the tool's contract: what `install` does, what `prune` spares, what a skill
  guarantees. A spec describing *how this repository releases itself* can be reversed without
  anything a user depends on moving, and reversing one is not a major even after `1.0`.
- **A version number cannot be recalled.** `proxy.golang.org` caches a version the first time
  anybody resolves it and serves it indefinitely by design, and `sum.golang.org` records its
  hash permanently. Deleting the tag does not undo either. So a wrong number is not a mistake
  you fix — it is a number burned for that content forever.

**This is written from having done it.** `v1.0.0` and `v1.0.1` were published from this table
read mechanically: a workflow reversing two promises in `.agents/specs/ci/spec.md` was labelled
`release:major`, and the tool's contract had not moved by a line. Both were cached by the proxy
within minutes. The tags and Releases are gone from the remote; the versions are not gone from
the ecosystem, and neither number can ever name different content. That is the cost of reading
this table without the paragraph above it, which is why the paragraph now exists.

### Do not delete the `v1.0.2` tag. It is load-bearing.

**It has no Release and it looks like junk. It is what keeps `@latest` on the `0.5.x` line.**

Deleting the two bad tags was not enough: `proxy.golang.org` had cached them, so
`go install ...@latest` kept resolving to `v1.0.1` with the tags already gone from GitHub —
measured, not assumed. Nothing removes a version from the proxy or from `sum.golang.org`; both
are append-only by design, and that is the same guarantee that stops anybody's dependency
changing underneath them.

What is available is `retract`, and only the **highest** version can retract — including
itself. So `v1.0.2` exists as a tombstone: no Release, a `retract` block in its `go.mod`
covering `v1.0.0`, `v1.0.1` and `v1.0.2`, and **a commit deliberately on no branch**, reachable
through the tag alone so `git describe --tags` from `main` keeps answering `0.5.x` instead of
climbing to `v1.0.3`.

Delete that tag and the retraction goes with it: the two bad versions become selectable again
and `@latest` returns to `v1.0.1`. Measured after publishing it:

```
go list -m …@latest              → v0.5.2
go list -m -versions …           → … v0.5.1 v0.5.2      (no v1 offered)
go install …/cmd/libretto@latest → libretto-automata v0.5.2
```

`retract` hides them; it does not delete them. `@v1.0.1` still resolves for anyone who asks for
it by name, marked `(retracted)` with the reason from `go.mod`. That is the whole of what the
ecosystem allows, and wanting more of it is what re-tagging a published version looks like —
which would change a hash the checksum database has pinned and turn a wrong number into a
failed install for everybody.

**The reading is yours. The typing is not.** Put a `release:patch`, `release:minor` or
`release:major` label on the request before it merges, and `.github/workflows/release.yml`
does the rest: gates, tag `main`, push it, open the Release.

```bash
gh pr edit <n> --add-label release:patch   # the one decision that is not mechanical
```

**With no label the run refuses** and names the three. That is deliberate: the bump turns on
whether a promise moved, which is a reading of `.agents/specs/` rather than of the commit
types, so a workflow that guessed would be wrong precisely when a contract moved. It asks
instead of assuming, and it never assumes patch.

`make release` is still there for a tag made outside this path — see below — but it is no
longer the route.

**The push comes before `make release`, and the order is load-bearing.** `gh release create`
creates the tag itself when it is not on the remote, at the default branch's HEAD — so running
it first gives you a second tag with your name on it pointing somewhere you did not choose, and
yours is the one that loses. The target passes `--verify-tag`, which refuses instead.

**A tag is not a Release**, and here only the tag is load-bearing. A tag is a git ref; a Release
is a GitHub object. The module proxy resolves `@latest` from tags, so `go install` works with no
Releases at all — which was this repository's state at `v0.4.0`: four tags, zero Releases,
verified. `make release` opens a Release page carrying the tag's notes, for humans; skipping it
costs nothing mechanical.

**Branch pushes do not tag.** Push a feature branch as often as you like; push it again after
review; none of that is a release. Only the merge is — and that is now enforced rather than
trusted: the workflow fires on a merged pull request and on nothing else.

That boundary is the whole rule, and it is where the old version of this section went wrong in
the other direction. It said "tag when the work goes out" and left *when* to judgment — so
work went out untagged and `main` sat ahead of its last release. And the failure it was written
to prevent is real too: four tags for one unpushed feature is four releases nobody could
install, because the fixes were to code that had never existed for anyone. **Tagging every
branch push reproduces exactly that.** One tag per merge is the line that avoids both.

**A tag is still a release, not a commit marker.** The difference now is that a merge *is* the
work going out, so there is nothing left to judge.

Retagging to move a tag onto a tip is gone as a practice — with one tag per merge there is
never a stale tag to move.

Retagging is only safe while the tag is unpublished. Once it is on the remote it is
somebody's reference point and it stays where it is.

**Skill frontmatter `version:` is separate** and tracks that skill's own contract. A
skill whose instructions change meaningfully gets a bump; the repo tag does not force
one.

Pre-1.0, `install`/`prune` behaviour may still change. Say so in the tag message when it
does — and note that this is the same fact the bump table depends on: while that is true, `1.0.0`
is not available to be reached by accident.

## Commits

Conventional commits. `type(scope): subject` — imperative, lowercase after the type, no
trailing period, ≤72 chars.

The body explains **why**. The diff already says what changed; it cannot say what the
alternative was or why it lost.

**A commit that removes or reverses a promise says so in the commit itself**: `!` after
the scope — `feat(cli)!: prune no longer spares …` — or a `BREAKING CHANGE:` footer when
the subject has no room. This changes nothing in the bump table: while this is `0.x` a
broken promise is still a minor. What it buys is that `git log` carries the record with
no forge attached — the label and the tag message live on GitHub, the `!` lives in every
clone — and when the `1.0.0` debate arrives, `git log --oneline | rg '!'` lists every
promise that moved.

**No AI attribution.** No `Co-Authored-By` for a model, no generated-with trailer.

One commit per finished unit of work, so a bisect lands on something meaningful. The
spec ships in the same commit as the code that taught it.

## Boundaries

### Always

- run all six gates before committing
- write the test in the same commit as the logic it proves
- name a `Proof:` in the spec for a new criterion, then make it exist
- point `CLAUDE_HOME` at a temporary directory in anything that touches a target
- **label the request with its bump before it merges** — `release:patch`, `release:minor` or
  `release:major`, read off what the merge contains. See *Versioning*. The tag and the Release
  then happen on merge; **with no label the run refuses and `main` is left untagged**, which is
  a release nobody can install and a `libretto version` naming a tag the code is past.
- state what was deliberately left out, and what would bring it back
- write `ponytail:` comments in English, whatever language the session is in. The comment
  lives in the source next to English identifiers, and `ponytail-debt` harvests it into a
  ledger that has to read as one document.

### Ask first

- adding a dependency. Five direct ones today, and the ladder is: stdlib, then a native
  feature, then something already here, then — last — something new.
- changing anything in `docs/STATE.md` under *decisions not to relitigate*
- deleting `install.sh`. It goes only when `libretto install` has been verified against
  a real `~/.claude` with a throwaway item, and that has not happened.
- a git worktree in a project whose build needs unversioned files
- any push
- **which bump a merge deserves, when it is arguable.** Patch versus minor turns on whether a
  promise moved, and that is a reading of the specs rather than of the commit types. Say which
  one and why; do not pick the smaller one to avoid the question.
- **`release:major`, always, with no exception while this is `0.x`.** Not announced — *asked*,
  and waited for. Announcing a major three times and proceeding is what published `v1.0.0` off
  a misread table, and the number cannot be recalled once the proxy has it. If the answer has
  not arrived, the label does not go on.

### Never

- **write to a real `~/.claude` in a test.** `CLAUDE_HOME` exists for this and it is
  what makes the suite safe to run twice.
- **weaken, skip or delete a failing test to get a green gate.** Fix the cause or stop
  and say why. Two moves, no third.
- **overwrite anything the tool did not create.** No `--force` for conflicts, ever.
  Reported, never resolved.
- **follow a symlink when removing it.** Remove the link. Never its destination.
- accept, store, echo or write an API token. `jira init` puts it in the OS keyring and
  the user runs that themselves.
- commit `CLAUDE.md` or `settings.json` handling — other tooling rewrites them and this
  is out of scope by decision.
- put `𝄞` anywhere a terminal will see it. README only: it is outside the BMP and
  renders as tofu. `♩♪♫♬` are banned too — ambiguous width tears the layout.
- reference `scripts/` or `docs/` from a skill. Those are not installed, so the skill
  would work only inside this repo. A tool a skill needs ships inside the skill.
- end a turn with a gate still running in the background, or report as done anything
  that was not observed.
