# add-change-queue

Targets: payload

Two new commands: `/libretto-queue` captures feature ideas into `.agents/changes/` as
queued proposals, and `/libretto-next` takes the oldest one and runs it through the
flow. Capture and execution are separate on purpose — an idea written down is not work
started.

## Outcomes

- `/libretto-queue` asks for ideas, one at a time, and each one becomes
  `.agents/changes/<verb-led-name>/proposal.md` with `Tracker: none`, a
  `Queued: <ISO date>` line, and the ask in the words it was asked in. It keeps asking
  until the user says stop. Capture only: no spec, no plan, no branch, no code.
- Queued proposals are committed on the current branch as docs-only commits, so the
  queue is visible from `main`. The branch-per-change rule fires when work *starts*,
  not when an idea is captured.
- `/libretto-next` lists the queued proposals oldest-first, recommends the oldest, lets
  the user pick another, then: creates the branch, removes the `Queued:` line, and
  enters the flow at phase 2 — phase 1's artifact already exists.
- **A queued idea is not work in flight.** `/libretto-flow <task>` never diverts to the
  queue; find-work's source 1 keeps meaning *started* work (plan boxes, unmerged
  branches). `/libretto-status` reports queued ideas as their own section, after
  in-flight.
- Both commands are payload: installed by `libretto install`, self-sufficient, and they
  delegate scanning to `find-work` the way `libretto-status` does rather than walking
  the directories their own way.

## Scope boundaries

**In:** the two command files, the queued-idea convention (`Queued:` line), the
find-work amendment that teaches it to report queued proposals without treating them as
in flight, and the `libretto-status` amendment that shows the queue.

**Out:**

- **priorities, reordering, tags.** The queue is FIFO by `Queued:` date; `/libretto-next`
  letting the user pick a different one *is* the reordering mechanism. A priority field
  returns the day FIFO measurably hurts.
- **editing or deleting queued ideas via a command.** They are markdown files; edit or
  delete the folder. A CRUD surface over three files is ceremony.
- **a separate queue file or directory.** `.agents/changes/` is already where proposals
  live; a `queue.md` beside it would be a second source of truth about what is queued.
- **the Go CLI knowing about the queue.** This is payload, not `libretto` the binary.
- **batch execution.** `/libretto-next` runs one idea and the flow's own stops apply;
  "do the whole queue unattended" is a different feature with different risks.

## Constraints

- A skill or command may only invoke what gets installed — no `scripts/`, no `docs/`
  paths.
- Frontmatter `name:` equals the filename; every referenced skill exists
  (`scripts/check-payload` enforces both).
- One author per file: the commands route, `find-work` owns the scan.
- `/libretto-queue` never accepts a tracker key as a queued idea — a key means the
  tracker is the source of truth and the flow already handles it.

## Prior decisions

- **Execution is a dedicated command, not `/libretto-flow` overloading** — the user's
  call, 2026-08-11: "le pasas una tarea… no tiene sentido pasarle la tarea y que haga
  otra". `/libretto-flow` with an argument does that argument, always; the queue drains
  through `/libretto-next`.
- **Queued ≠ in flight.** Source 1's home-first rule exists so *started* work does not
  get abandoned. Ideas are cheap and abandoning one costs nothing; blocking a Jira task
  because four ideas sit captured would make capture punitive and the queue would go
  unused.
- **FIFO by a `Queued:` line, not git archaeology.** The date is in the file so ordering
  survives rebases and needs no `git log` walk. Ceiling: no priorities; the replacement
  is a priority field the day someone reorders constantly.
- **Capture commits on the current branch.** A branch per captured idea would scatter
  the queue across N branches nobody can see from `main`. The existing "branch before
  the first write" rule is about the change's own work; a queued proposal is not the
  change's work yet, and `/libretto-next` creates the branch at pickup, which is that
  change's first write.

## Task breakdown

- [ ] `commands/libretto-queue.md` — the capture loop
- [ ] `commands/libretto-next.md` — pick, branch, de-queue, enter phase 2
- [ ] `skills/find-work/SKILL.md` — queued proposals: reported, never blocking, never
      source 1
- [ ] `commands/libretto-status.md` — the queue section
- [ ] `.agents/specs/payload/spec.md` delta lands with the change (phase 8)

## Verification criteria

- both commands parse and every skill they reference exists
  Proof: scripts/check-payload
- no command or skill references a path that does not get installed
  Proof: scripts/check-payload
- every `Proof:` citation in this delta resolves
  Proof: skills/record-work/spec-drift --anchors

**Behaviour is prose and is checked by running it.** Stated as to-observe, the same
standing the payload spec uses: capturing two ideas leaves two committed proposals and
no branch; `/libretto-next` offers the oldest first and enters phase 2 on a fresh
branch; `/libretto-flow EUCAR-123` with ideas queued reads the key and never mentions
the queue.
