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

libretto install --project   # <cwd>/.claude instead of ~/.claude
libretto install --global    # the default; both flags at once is an error
```

## Gates — all six pass before any commit

```bash
gofmt -l .                                       # must print nothing
go vet ./...
go test ./... -count=1                           # 117 tests
scripts/check-payload                            # frontmatter, references, reachability
skills/record-work/spec-drift --self-test        # 17 checks
skills/record-work/spec-drift --anchors          # 105 citations must resolve
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
internal/repo/          git, and the rebuild decision. ONE test for 155 lines — the largest gap.
internal/ui/            logo, theme, panel, model
skills/ agents/ commands/   THE PAYLOAD — what gets symlinked
scripts/check-payload   repo-only tooling. Never referenced from a skill.
.agents/specs/          the specification, one directory per capability
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

**Mark the boxes as you go.** A plan that is 11/11 open while the work is nearly done is
documentation pretending to be state — which happened while building this very feature.

## Specs

Per **capability**, never per ticket, in `.agents/specs/<capability>/spec.md`. Eight of
them: `ownership`, `link-state`, `linking`, `targets`, `repo-sync`, `panel`, `cli`,
`payload`. `docs/SPEC.md` is an index with no requirements in it.

Each declares what it owns and cites the test behind each criterion:

```
Governs: internal/link/plan.go internal/link/apply.go internal/link/plan_test.go
Proof:   internal/link/apply_test.go TestApplyIsIdempotent
```

A `Proof:` citation must name a **test that exists**. `--anchors` checks the test name,
not just the file, because a file-level check passes an invented name — that happened
here, twice.

Work in flight goes in `.agents/changes/<change>/`, and lands by applying the delta onto
the capability spec and deleting the change folder, in the same commit as the final code.

## Versioning

Semver, from git tags. **The binary never carries a hardcoded version.**

```bash
git tag -a v0.2.0 -m "..."     # the tag is the release
make build                     # -ldflags stamps `git describe --tags --always --dirty`
./bin/libretto version         # v0.2.0  ·  v0.2.0-3-gabc123  ·  v0.2.0-dirty
```

A plain `go build` with no ldflags reports `dev`. That is deliberate: a binary that
cannot prove its version says so rather than claiming one.

| Bump | When |
|---|---|
| patch | a fix with no contract change |
| minor | a new capability, or a new promise in an existing spec |
| major | a promise removed or reversed |

**A tag is a release, not a commit marker.** Tag when the work goes out, not on the way
there. Four tags for one unpushed feature and three fixes to it is four releases nobody
could install — the fixes were to code that had never existed for anyone, so they need
no patch numbers of their own. Move the tag to the tip and let it describe what actually
shipped.

Retagging is only safe while the tag is unpublished. Once it is on the remote it is
somebody's reference point and it stays where it is.

**Skill frontmatter `version:` is separate** and tracks that skill's own contract. A
skill whose instructions change meaningfully gets a bump; the repo tag does not force
one.

Pre-1.0, `install`/`prune` behaviour may still change. Say so in the tag message when it
does.

## Commits

Conventional commits. `type(scope): subject` — imperative, lowercase after the type, no
trailing period, ≤72 chars.

The body explains **why**. The diff already says what changed; it cannot say what the
alternative was or why it lost.

**No AI attribution.** No `Co-Authored-By` for a model, no generated-with trailer.

One commit per finished unit of work, so a bisect lands on something meaningful. The
spec ships in the same commit as the code that taught it.

## Boundaries

### Always

- run all six gates before committing
- write the test in the same commit as the logic it proves
- name a `Proof:` in the spec for a new criterion, then make it exist
- point `CLAUDE_HOME` at a temporary directory in anything that touches a target
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
