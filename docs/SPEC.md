# Specification — index

**The specification lives in [`.agents/specs/`](../.agents/specs/), one directory per
capability.** This file is a signpost, not a spec: it carries no requirements and no
`Governs:` line, so there is still exactly one place where the contract is written.

It exists because `DESIGN.md`, `PLAN.md` and `STATE.md` link here, and because the old
`R1…R11` numbering needs somewhere to resolve.

## Purpose

Keep one person's AI agent configuration — skills, agents, commands — versioned in git
and installed by symlink, so it is editable in place, reviewable in history, and safe
from being overwritten by other tooling that writes into the same directories.

## Vocabulary

| Term | Meaning |
|---|---|
| **repo** | this repository, the source of truth |
| **item** | one directory under `skills/`, or one `.md` file under `agents/` or `commands/` |
| **target** | an agent installation that consumes items, e.g. Claude Code at `~/.claude` |
| **kind** | a category of item: `skills`, `agents`, `commands` |
| **owned link** | a symlink in a target whose destination resolves inside the repo |
| **foreign entry** | anything else in the target's directory — another tool's files, or the user's own |
| **capability** | one coherent area of behaviour, with its own spec and its own `Governs:` paths |

## Capabilities

| Capability | Governs | What it owns |
|---|---|---|
| [ownership](../.agents/specs/ownership/spec.md) | `internal/link/own.go` | whether a symlink is ours. The most consequential predicate in the program. |
| [link-state](../.agents/specs/link-state/spec.md) | `internal/link/scan.go`, `state.go` | enumeration and the five states. Read-only. |
| [linking](../.agents/specs/linking/spec.md) | `internal/link/plan.go`, `apply.go` | plan then apply. The only code that writes. |
| [targets](../.agents/specs/targets/spec.md) | `internal/target/**` | what an installable destination is |
| [repo-sync](../.agents/specs/repo-sync/spec.md) | `internal/repo/**` | pull, and when the binary is invalidated |
| [panel](../.agents/specs/panel/spec.md) | `internal/ui/**` | logo, palette, layout, menu |
| [cli](../.agents/specs/cli/spec.md) | `cmd/libretto/**`, `install.sh` | the command surface, exit codes, environment |
| [payload](../.agents/specs/payload/spec.md) | `skills/**`, `commands/**`, `agents/**`, `scripts/**` | the flow itself — the reason the project exists |
| [agent-models](../.agents/specs/agent-models/spec.md) | `internal/agentmodel/**` | which model each payload agent runs on |
| [review-project](../.agents/specs/review-project/spec.md) | `skills/review-*/**`, `agents/review-*.md`, `commands/libretto-review.md` | reviewing somebody else's PR/MR without disturbing your own state |

The flow's own reasoning — its eight phases and why they are those — is
[FLOW.md](FLOW.md), not a spec.

## Why capabilities instead of requirement numbers

`R1…R11` described the tool as one thing with eleven promises. That reads well and
maintains badly: a change to the ownership predicate touched R1, R2, R4, R6 and R9 at
once, so the whole document had to be re-read to learn what a two-line change affected.

A capability spec declares `Governs:`. `spec-drift` reads it and names the spec that owns
a staged file, which turns "did the spec move?" from a judgement into a lookup.

## Where the old numbers went

| Was | Now |
|---|---|
| R1 link each item individually | linking · outcomes |
| R2 never overwrite a foreign entry | linking · scope boundaries, and ownership entire |
| R3 `update` refreshes and relinks | repo-sync |
| R4 `status` reports the truth | link-state |
| R5 `doctor` diagnoses | cli · outcomes |
| R6 `prune` removes stale links | linking, and cli for the `--yes` gate |
| R7 interactive panel | panel |
| R8 every action is a subcommand | cli |
| R9 mutations are planned, then applied | linking · constraints |
| R10 targets are extensible | targets |
| R11 configuration is environment only | cli · constraints |

`PLAN.md` still carries `R…` tags. They resolve through this table.

## The anchors

Every capability spec declares what it governs and cites the test behind each criterion:

```
Governs: internal/link/plan.go internal/link/apply.go
Proof:   internal/link/apply_test.go TestApplyIsIdempotent
```

```
skills/record-work/spec-drift --anchors   every citation resolves, test name included
skills/record-work/spec-drift             staged code whose owning spec did not move
```

**70 citations resolve today, and the check earned its keep immediately.** The first pass
named a test that did not exist, and ten more named the right test in the wrong file. It
also caught itself: indented `Proof:` lines were invisible to the extractor, so it
reported "all citations resolve" having read seven of seventy — a false green, which is
the worst kind.

## The largest gaps

Each capability lists its own in its task breakdown. Two are worth naming here:

- **`cmd/libretto` has no tests.** 662 lines, `no test files`. Exit codes, dispatch and
  flag handling are held up by manual runs only.
- **`internal/repo` has one test for 155 lines.** The `Git` interface exists to be faked
  and no fake exists, so the dirty-tree refusal and the rebuild notice are promises
  nothing keeps.

## Out of scope

- **`CLAUDE.md` and `settings.json`.** Other tooling rewrites regions of these files.
  They stay hand-managed.
- **Targets other than Claude Code.** The interface exists; the implementations do not.
- **Migrating the author's existing skills into this repo.** After `install` is proven.
- **Resolving conflicts.** Reported, never resolved.
- **Publishing.** Homebrew formula, release automation, cross-compilation — later.
