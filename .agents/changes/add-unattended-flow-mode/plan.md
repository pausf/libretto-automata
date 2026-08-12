# add-unattended-flow-mode — plan

Spec: `.agents/changes/add-unattended-flow-mode/spec.md` (Targets: payload)
Branch: `feat/add-unattended-flow-mode`

**Goal:** a second door onto the flow, `/libretto-attacca`, that runs phases 1–8 without
stopping at a question and ends at a pushed request.

**Shape:** every task here is prose in the payload. Nothing compiles, so the only
mechanical proof is `scripts/check-payload`; the mode itself is verified by running it,
which is task 5's job and the spec says so in its own words.

**Execution is phase 6, `build-and-check`** — not `subagent-driven-development`. The
payload's own flow owns the build, and phase 6 does not fan out.

## Global constraints

- **The classification lives in exactly one file.** `commands/libretto-attacca.md` states
  what the invocation answers, what stays a stop and what an underivable question becomes.
  Every other file points at it and restates none of it.
- **No file describes a phase that is not its own.** The command names the skills it
  delegates to and explains none of them.
- Six gates before every commit: `gofmt -l .`, `go vet ./...`, `go test ./... -count=1`,
  `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`.
- Conventional commits, no AI attribution, one commit per task.
- `𝄞` and `♩♪♫♬` never appear in payload prose — README only. The command is named after a
  score instruction; it does not draw one.

## Task 1 — the command

**Files:** create `commands/libretto-attacca.md`

Depends on: nothing. Everything else depends on this.

- [ ] frontmatter with `description:` only, matching the shape of `commands/libretto-flow.md`.
      The description carries the gloss — `attacca` is the score instruction for *go on
      without pausing* — because the name is opaque to a reader who does not read music.
- [ ] the classification table: what the invocation answers (the spec stop, the plan stop,
      phase 8's push-and-request), what stays a stop (a failing gate, two stopped tasks, a
      missing or unauthenticated `jira`/`gh`/`glab`), and what an underivable question
      becomes (an assumption written into the spec's prior decisions, the phase 7 report
      and the request description).
- [ ] the delegation: phase 1 through 8 by `Skill(...)` name, in order, with no description
      of any phase.
- [ ] what it never does: merge, tag, release, apply a `release:` label, skip a gate, or
      drain the queue.
- [ ] gates, then commit.

**Closes:** *every skill and command referenced by the new prose exists* — Proof:
`scripts/check-payload`. And *frontmatter still parses* — same proof.

## Task 2 — the five stop-owning skills

**Files:** modify `skills/find-work/SKILL.md`, `skills/write-spec/SKILL.md`,
`skills/write-plan/SKILL.md`, `skills/record-work/SKILL.md`, `skills/present-work/SKILL.md`

Depends on: task 1 — each of these names the command, so the command must exist.

One commit, not five. These are five instances of one sentence, they land together, and a
reviewer who rejects one rejects the shape rather than the file.

- [ ] `find-work` — beside its two stops: under attacca the in-flight choice resolves to
      the oldest change with open boxes, and a missing tracker CLI still stops.
- [ ] `write-spec` — beside "Then stop, and wait": under attacca the contract stop carries
      on, and step 4's question is written into prior decisions marked as assumed instead
      of being asked.
- [ ] `write-plan` — beside "Then stop": the order stop carries on.
- [ ] `record-work` — beside "Pushing is the user's decision": under attacca the answer
      arrived when the command was typed, it covers this branch and this request and
      nothing past it, and the request description carries every assumption the run made.
- [ ] `present-work` — the report names what the mode answered and every assumption made.
- [ ] each edit says *what* changes and points at `/libretto-attacca` for *why*. A skill
      that restates the classification is a second copy of it.
- [ ] gates, then commit.

**Closes:** *no skill invokes a path that does not get installed* — Proof:
`scripts/check-payload`.

## Task 3 — the other door, named

**Files:** modify `commands/libretto-flow.md`

Depends on: task 1.

- [ ] one line, in the stops section, naming `/libretto-attacca` as the invocation that
      answers them — the way the file already names `/libretto-status`. One line, not a
      section: the classification is task 1's.
- [ ] gates, then commit.

**Closes:** *every referenced skill and command exists* — Proof: `scripts/check-payload`.

## Task 4 — the flow document

**Files:** modify `docs/FLOW.md`

Depends on: task 1. Independent of tasks 2 and 3.

- [ ] the mode beside the stops it removes, and the line that is the point: a gate is not
      a stop and cannot be answered.
- [ ] gates, then commit.

**Closes:** nothing mechanical. `docs/` is governed by no spec, and `check-payload` does
not read it — stated here so the gap is deliberate rather than missed.

## Task 5 — run it, then land it

Depends on: tasks 1–4.

- [ ] **run `/libretto-attacca` against a real, small piece of work** and record what was
      observed, not what was expected. The spec lists five claims; each becomes an
      observation or stays a claim, in the capability spec's own words.
- [ ] apply the delta onto `.agents/specs/payload/spec.md`: the mode in Outcomes, its scope
      boundaries folded into the existing Out list, both prior decisions carried across
      with their dates, the stops table gaining the attacca column or the sentence that
      replaces it.
- [ ] delete `.agents/changes/add-unattended-flow-mode/` — proposal, spec, plan.
- [ ] `spec-drift --anchors`, then the six gates.
- [ ] one commit: the final prose, the applied delta and the deleted folder together.

**Closes:** *every `Proof:` citation resolves, file and test name* — Proof:
`skills/record-work/spec-drift --anchors`.

## What can start now

Task 1. Tasks 2, 3 and 4 are independent of each other and all wait on it; task 5 waits on
all four.
