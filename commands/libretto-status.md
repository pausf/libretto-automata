---
description: What work is in flight, how much is left, and what can start now
---

Read-only. Reports; chooses nothing, changes nothing.

```
Skill(skill="find-work")
```

Invoke it in **reporting mode**: two of its scans and no more — source 1
(`.agents/changes/*/plan.md`, open boxes, what each plan says can start, and the unmerged
branches) and **the queue** (its *captured, never started* scan). Then stop.

Both are named because "source 1 only" left the queue ambiguous: the queue is not a
source, so an agent reading that literally reported an empty house while ideas sat
captured on disk.

Do not read the tracker, do not ask what to work on, do not begin a phase.

The scan is the skill's, not this command's. A status command that walked the same
directories its own way would be a second answer to "what work exists", and the one that
disagrees is always the one nobody is reading.

## What to report

Per change in flight:

- its name, and what its `proposal.md` says it is for, in one line
- boxes open out of total
- **what can start now**, from the plan's own dependency notes

Then one line for the whole picture: how many changes are open, how many boxes between
them.

Nothing in flight is a state, not an emptiness to apologise for. Say so in one line.

## Then the queue, as its own section

Ideas captured by `/libretto-queue` come **after** in flight, oldest first, name and one
line each. The skill finds them; this command only keeps them separate.

**What counts as queued is the skill's to decide, and it is not just the `Queued:` line.**
This sentence used to restate that definition and the restatement went stale the day the
rule grew: a change dispatched by a branch keeps the line on the base branch and is not
queued. Two descriptions of one thing is one too many, and the copy that drifts is always
the one nobody is editing.

Two sections, never one list. A captured idea and a change half-built are different kinds
of thing, and merging them makes the queue look like a backlog of stalled work.

## What not to do

**Do not offer to pick one up.** That is `/libretto-flow`'s job, and it asks first.
Reporting that quietly becomes an offer is how a read-only command starts writing.
