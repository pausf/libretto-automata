---
description: Take the next queued idea and run it through the flow
---

Drains the queue `/libretto-queue` filled, one idea at a time.

This is a separate command and not `/libretto-flow` reading the queue on its own.
`/libretto-flow <task>` does *that* task, always — handing it a Jira key and having it
work on something else instead is the surprise nobody wants.

## Step 1 — Read the queue

```
Skill(skill="find-work")
```

Claude Code's spelling, and the same for `AskUserQuestion` below. Load the skill with the
host's own skill tool, and ask with the host's own native prompt — or in conversation where
there is none.

Ask it for **the queue only** — its *captured, never started* scan, oldest first. Not
source 1, not the tracker, and no phase begun: this command chooses the work itself in
step 2.

The scan is the skill's, not this command's, exactly as `/libretto-status` delegates it.
A second walk of the same directory here would be a second answer to "what is queued",
and the one that disagrees is always the one nobody is reading.

Nothing queued is a state, not an error — say so in one line, point at
`/libretto-queue`, and stop.

## Step 2 — Which one

Ask with `AskUserQuestion`: the oldest as the recommendation, the others as the
alternatives, room to name a different one. **Never choose.** FIFO is the default order,
not a rule the user has to argue with — picking a different one here *is* how the queue
gets reordered, which is why there is no priority field.

## Step 3 — Start it

In this order, because the first write needs the branch to already exist:

1. `git checkout -b <type>/<name>` from the base branch. This is the change's first
   write, so it is where the branch belongs — the same rule as `find-work` and
   `build-and-check` step 0.
2. Remove the `Queued:` line from `proposal.md`. The line means *not started*; leaving it
   in would keep offering work that is already underway.
3. Commit both together.

## Step 4 — Enter the flow at phase 2

**Phase 1 is already done** — its artifact is the proposal on disk, which is the whole
point of the queue. Do not re-run `find-work`: it would find the change in flight and ask
whether to continue the thing just chosen.

```
Skill(skill="write-spec")
```

From there the flow is the flow, exactly as `/libretto-flow` runs it: the spec stop, the
plan stop, build, review, present, and phase 8's question about pushing. Each phase's
skill owns its own rules; this command only knows how the work got picked.

## One idea, then stop

The flow's stops apply. Draining the queue unattended is a different feature with
different risks, and it is deliberately not this one: come back and run
`/libretto-next` again.
