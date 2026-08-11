<div align="center">

# 𝄞 Libretto Automata

**The libretto is written first. The automaton performs it.**

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-1.3-FF6AC1?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-1.1-FFB3E8?logo=charm&logoColor=white)](https://github.com/charmbracelet/lipgloss)
[![Claude Code](https://img.shields.io/badge/Claude%20Code-payload-D97757?logo=anthropic&logoColor=white)](https://claude.com/claude-code)
[![Jira CLI](https://img.shields.io/badge/Jira%20CLI-tracker-0052CC?logo=jira&logoColor=white)](https://github.com/ankitpokhrel/jira-cli)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-367%20passing-brightgreen.svg)](#gates)

<img src="docs/panel.svg" alt="The libretto panel: a menu with install, uninstall, update, status, models, doctor and prune, acting on a chosen destination — here ~/.claude with all 26 items linked" width="800">

*The panel — real terminal output, captured from the binary. `tab` switches where it acts.*

</div>

---

An 18th-century automaton was a machine that played music by reading a score. A human
wrote the notes; the machine executed them. It never improvised — it did exactly what
the paper said. **If the performance was wrong, the paper was wrong.**

That is how this works. You write the spec. The agent performs it. When the output is
bad, you fix the libretto, not the automaton.

## Two things in one repository

| | |
|---|---|
| **The payload** | an eight-phase spec-driven flow, as skills and commands. **This is the point.** |
| **The CLI** | a Go binary that symlinks the payload into `~/.claude`. This is delivery. |

Links are made **per item, never per directory**, so this coexists with anything else
installed into the same folders. Anything already there that this tool did not create is
left untouched and reported — there is no `--force`, by design.

## Install

Requires [Go 1.26+](https://go.dev/dl/).

```bash
go install github.com/pausf/libretto-automata/cmd/libretto@latest

libretto install   # symlink every item into ~/.claude
libretto           # the panel
```

**That is the whole install.** The payload — skills, agents and commands — ships *inside* the Go
module, so `go install` downloads it along with the binary and verifies it against Go's checksum
database. There is no tarball to fetch and no bootstrap step.

It lands here:

```
~/go/bin/libretto                                          the command
$GOMODCACHE/github.com/pausf/libretto-automata@v0.5.0/     the payload
  ├── skills/  agents/  commands/

~/.claude/skills/write-spec  →  …@v0.5.0/skills/write-spec
```

The version is in the path, which is what makes an update a new directory rather than an
overwrite of the one your links point at. `LIBRETTO_ROOT` points the tool somewhere else.

### Updating

```bash
libretto update
```

Installs the newest version and relinks. The relink is not redundant — a new version is a new
directory, and a release that *adds* a skill would otherwise leave it unlinked with nothing to say
so.

The panel also tells you when a newer version exists, checked once a day and silent when it
cannot check.

### From a checkout instead

For working *on* the payload rather than with it:

```bash
git clone git@github.com:pausf/libretto-automata.git ~/gitrepos/libretto-automata
cd ~/gitrepos/libretto-automata

make build      # stamps the version from git describe
make link       # puts `libretto` on your PATH via ~/.local/bin
```

`make link` symlinks rather than copies, so `make build` updates the installed command and no
stale binary can pretend to be current. It refuses to overwrite a `libretto` it did not create.

A checkout you are standing in wins over the module cache, so editing a skill and seeing it live
still works — which is the reason the payload is not compiled into the binary. `update` pulls
there instead of downloading: same command, and the tree you are standing in is the installation.

### Releasing

```bash
git tag -a v0.5.0 -m "..."
git push origin v0.5.0
make release       # gates, then a Release page carrying the tag's notes
```

Nothing is attached to the Release — the module proxy resolves `@latest` from **tags**, so
installing works without one. `make release` exists so a human can read what changed.

## Commands

| | |
|---|---|
| `libretto` | the panel — needs a terminal |
| `libretto status` | every item's state. Read-only, always. |

| `libretto install` | link everything. Idempotent; non-zero exit if anything was skipped. |
| `libretto uninstall` | show what this repo installed here — **changes nothing** |
| `libretto uninstall --yes` | take it back out |
| `libretto update` | install the newest version and relink — or pull, in a checkout |
| `libretto doctor` | what needs attention, plus what the payload expects here |
| `libretto prune` | show links whose source is gone — **changes nothing** |
| `libretto prune --yes` | remove them |
| `libretto preview` | print the panel once, no TUI |
| `libretto models` | which model each agent runs on. Read-only. |
| `libretto models set <model> <agent>…` | declare it; `--all` for every agent |

### Choosing a model per agent

Not every agent needs the same model. The review lenses that pattern-match over prose —
design, tests — do the job on a cheap one; the lens looking for what an attacker can
reach should not.

```bash
libretto models                                  # who runs on what
libretto models set haiku review-lens-design review-lens-tests
libretto models set default work-reviewer        # back to the session's model
```

The panel has the same thing behind its `models` row: mark agents with `space`, or all
of them with `a`, pick a model with `m`. One gesture for the whole set, because making
the prose lenses cheap is one decision, not four.

`models` acts on **the agents of the destination you name** — `~/.claude/agents` under
`--global`, `<cwd>/.claude/agents` under `--project`. Every agent there is listed and
editable, not just the ones this repository ships.

A row marked `shared` is a file this repository owns, reached from more than one
destination: writing it changes every project on the machine. An unmarked row is a real
file in that destination and changing it changes nothing else.

`default` means no `model:` key at all — an absent key is already how the format says
"whatever the session runs on", and two spellings of one state is a difference somebody
eventually treats as meaningful.

The values are **aliases**, so an agent file does not need editing the day a new model
ships. The listing shows what each one resolves to, and when that was last checked:

```
models available (aliases; versions as of 2026-08):
  default                 the session's model — whatever you are running
  haiku      Haiku 4.5    cheapest; fine for pattern-matching over prose
  sonnet     Sonnet 5     the everyday working model
  opus       Opus 5       most capable; Max plans, metered on Pro
```

### Where it installs

| | |
|---|---|
| `--global`, `-g` | `~/.claude` — **the default** |
| `--project`, `-p` | `<this directory>/.claude` |

A project-local install keeps a flow scoped to one repository without editing the
configuration every other project shares. In the panel, the strip shows both
destinations and `tab` switches which one the keys act on — the active one is marked
`◉` in gold.

Passing both flags is an error. Two answers to one question is a mistake worth
reporting, not one worth resolving by guessing.

`prune` and `uninstall` are both dry by default. A destructive command that acts before
being asked twice eventually deletes the wrong thing, and a pipe is no reason to be less
careful. In the panel they show the plan and then ask — `y` to go ahead, `n` to cancel,
and no other key carries them out.

**`prune` and `uninstall` are not the same thing.** Prune cleans up after *the repo*
changed — rename an item and the old link points at nothing, which is `stale`. Uninstall
removes links that are **working**, because you changed your mind. Prune deliberately
spares correct links, and that is what makes it safe to run: you clean one broken link
without risking a whole installation.

### The five states

`status` reports one of these per item, per target:

| State | Meaning | Remedy |
|---|---|---|
| `linked` | our symlink, right destination | none |
| `missing` | in the repo, absent from the target | `install` |
| `wrong target` | our symlink, wrong destination | `install` repoints it |
| `conflict` | something foreign in the way | **none** — reported, never touched |
| `stale` | our symlink with no item behind it | `prune` |

### Environment

| | |
|---|---|
| `CLAUDE_HOME` | Claude Code's root instead of `~/.claude`. What makes the test suite safe. |
| `LIBRETTO_ROOT` | the repo location, instead of deriving it from the binary |
| `LIBRETTO_ASCII=safe` | swap the clef's quadrant glyphs for half blocks |
| `LIBRETTO_THEME` | `dark` or `light`, instead of detecting |
| `COLUMNS` | layout width when stdout is not a terminal |

There is no configuration file. A value that never varies is not configuration.

## The flow

Eight phases, installed as skills. Start it with `/libretto-flow`, with or without a
tracker key.

**The flow does not begin at a tracker.** Phase 1 asks three sources in order — a change
already in flight, a tracker key or URL, and what you said — and the order is the point:
starting something new while a change sits half-finished is how the half-finished thing
gets abandoned.

`/libretto-status` runs only the first of those and stops, for when the question is just
"what is open?"

**Ideas arrive faster than they get built.** `/libretto-queue` captures them one after
another — a proposal with a `Queued:` date and your words verbatim, no branch and no spec
— and `/libretto-next` picks one up later, oldest first, and takes it into the flow.

Two commands and not one, because `/libretto-flow EUCAR-1234` does *that* ticket, always.
A flow that quietly substitutes different work for what you handed it is the surprise
nobody wants.

```bash
/libretto-status              # what is in flight and what is queued
/libretto-flow                # find the work and take it through the phases
/libretto-flow EUCAR-1234     # …starting from a ticket
/libretto-queue               # capture ideas, one after another, and build none of them
/libretto-next                # take the oldest queued idea into the flow
/libretto-review <pr-url>     # review a PR/MR in a workspace that restores itself
```

| | Phase | Skill |
|---|---|---|
| 1 | find the work — in flight, tracker, or asked for | `find-work` |
| 0·2·3 | does a spec even need to exist · the six pillars · one per subtask | `write-spec` |
| 5 | the plan — live state, one writer | `write-plan` |
| 6 | build, with proportionate checks | `build-and-check` |
| 7 | present, including what was left out | `present-work` |
| 8 | commit, and make the spec true again | `record-work` |

Three rules hold at **every** phase, not one of them:

- **ask** — with a recommended option, real alternatives, and room to answer otherwise
- **commit** — per task, so a bisect lands somewhere meaningful
- **evidence** — nothing is true until it has been observed

Details and reasoning: [docs/FLOW.md](docs/FLOW.md).

## Specs

Per **capability**, never per ticket, under [`.agents/specs/`](.agents/specs/). A
capability spec accumulates and stays true; a spec named after a ticket is dead the day
the ticket closes.

Each declares what it owns and cites the test behind every criterion:

```
Governs: internal/link/plan.go internal/link/apply.go
Proof:   internal/link/apply_test.go TestApplyIsIdempotent
```

Two anchors, checked mechanically:

```bash
skills/record-work/spec-drift --anchors   # every citation resolves, test name included
skills/record-work/spec-drift             # staged code whose owning spec did not move
```

It warns; it never blocks. A check that stops a commit in someone else's project is a
check that gets deleted, and a deleted check finds nothing.

Work in flight lives in `.agents/changes/<change>/` and lands by applying its delta onto
the capability spec and deleting the change folder — in the same commit as the code.

## Gates

All six pass before any commit — `make gates` runs them, and so does
[GitHub Actions](.github/workflows/gates.yml) on every push and pull request.

```bash
gofmt -l .                                  # must print nothing
go vet ./...
go test ./... -count=1                      # 251 tests
scripts/check-payload                       # frontmatter, references, reachability
skills/record-work/spec-drift --self-test   # 17 checks
skills/record-work/spec-drift --anchors     # 208 citations
```

## Built with

| | |
|---|---|
| [Go](https://go.dev) | 1.26 |
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) · [Lip Gloss](https://github.com/charmbracelet/lipgloss) · [termenv](https://github.com/muesli/termenv) | the panel |
| [jira-cli](https://github.com/ankitpokhrel/jira-cli) | the tracker, by CLI — never MCP, never the REST API |

## Standing on other people's work

Seven skills **ship with this repository**, so the flow works on a machine that has
nothing else installed. Copied unmodified, licence and version recorded in
[THIRD-PARTY.md](THIRD-PARTY.md):

- [**obra/superpowers**](https://github.com/obra/superpowers) — `writing-plans`,
  `test-driven-development`, `using-git-worktrees`. The flow's own skills are thin
  because they delegate to these, and a thin skill whose delegate is missing is not
  thin, it is broken.
- [**DietrichGebert/ponytail**](https://github.com/DietrichGebert/ponytail) —
  `ponytail`, `ponytail-debt`. Decides how much gets built: the ladder runs from
  *does this need to exist at all?* down to *only then, the minimum that works*, and
  it carries the list of things that are never trimmed: trust boundaries, data loss,
  security, accessibility. The flow invokes it in phase 2, on requirements, because
  that is where removing work is cheapest.
- [**JuliusBrussee/caveman**](https://github.com/JuliusBrussee/caveman) — `caveman`,
  `caveman-commit`. Decides how much gets said. Compresses prose; ponytail compresses
  what gets built. No overlap.

Shipped is not required: nothing fails without them, and they prune like any other
item. Only what the flow calls by name is vendored — the rest of both plugins,
including always-on hook mode, stays with the upstream plugins, and the two coexist
by namespace.

## Not managed here

`CLAUDE.md` and `settings.json`. Other tooling rewrites regions of those files, so they
stay hand-managed. Linking them would start a fight this tool cannot win.

## Licence

[MIT](LICENSE). Vendored items keep their own — see [THIRD-PARTY.md](THIRD-PARTY.md).
