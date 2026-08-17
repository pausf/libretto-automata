# Contributing

Thanks for looking. This is a small, opinionated repository, and the fastest way to have a
change accepted is to know four things it does differently from most projects.

**[`AGENTS.md`](AGENTS.md) is the rulebook** — the commit convention, the boundaries, what
never happens here. This file does not repeat it. It tells you the parts you would not think
to look for, and then sends you there.

## 1 · Run the six gates before you open anything

All six pass, or the change is not ready. There is no partial credit, and a failing test is
fixed or explained — never weakened, skipped or deleted to get a green result.

```bash
gofmt -l .                                       # must print nothing
go vet ./...
go test ./... -count=1
scripts/check-payload                            # frontmatter, references, reachability
skills/record-work/spec-drift --self-test
skills/record-work/spec-drift --anchors          # every Proof: citation must resolve
```

One trap worth knowing, because it has bitten this repository: **never pipe a gate into `head`
when the exit code matters.** The pipeline reports the last command's status, so a failure reads
as success. Redirect to a file, check `$?`, read the file.

`spec-drift` with no flag warns about staged code whose spec did not move. It never blocks —
run it, read it, answer the question out loud.

## 2 · Your pull request needs a `release:` label, and the run refuses without one

This is the thing that surprises everybody. Every merge to `main` gets a version tag, and the
tag comes from a label on the request:

| Label | When |
|---|---|
| `release:patch` | fixes, refactors, docs, chores — no promise moved |
| `release:minor` | a new feature, a new capability, or a new promise in a spec |
| `release:major` | **ask first, always.** Nobody has decided this project is at 1.0.0 |

**With no label the release run refuses and names the three.** That red check is not a broken
pipeline — it is the question being asked. It refuses rather than guessing, and it never assumes
patch.

**You propose the bump; you do not decide it.** Which one applies turns on whether a promise in
[`.agents/specs/`](.agents/specs/) moved, and that is a reading of the specification rather than
of your commit types. Say which you think it is and why. If it is arguable, say that too — the
maintainer settles it.

Two reasons this is stricter than it looks: a version number **cannot be recalled** once the Go
module proxy has cached it, and `1.0.0` here is a decision somebody makes rather than an output
of the table. Until it is made, nothing about this tool's interface is stable.

## 3 · Work does not come from a tracker

There is no board to pick from. Work arrives from one of three places, asked in that order:
a change already in flight under `.agents/changes/`, a tracker key if one
was handed over, or **what somebody said** — which is a legitimate input here, not a fallback,
and is how every change in this repository so far arrived.

So an issue describing what you want to change **is** the useful contribution. You do not need a
ticket id, and please do not invent one.

## 4 · The specification is per capability, and a change lands on it

[`docs/SPEC.md`](docs/SPEC.md) indexes the capabilities; each one declares the paths it governs
and cites, for every criterion, a test that exists. If your change adds a promise, it adds a
criterion, and that criterion names its proof by file and test name.

A criterion whose test does not exist fails
[`spec-drift --anchors`](skills/record-work/spec-drift) — by test name, not just by file, because
a file-level check passes an invented name.

## What to expect on review

Findings arrive attributed and unedited, including the ones that were already fixed. That is
deliberate: what was wrong is part of the record even when it is no longer true.

If a review says something you think is wrong, say so with evidence. Agreement without
verification is not useful to anybody here.

## Where the reasoning lives

- [`AGENTS.md`](AGENTS.md) — the rules, and why each one exists
- [`docs/FLOW.md`](docs/FLOW.md) — the eight-phase flow, and why it is those eight
- [`docs/DESIGN.md`](docs/DESIGN.md) — the design arguments
- [`docs/SPEC.md`](docs/SPEC.md) — the capability index
- [`THIRD-PARTY.md`](THIRD-PARTY.md) — the vendored skills and their licences

Licensed MIT — see [`LICENSE`](LICENSE).
