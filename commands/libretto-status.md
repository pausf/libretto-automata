---
description: What work is in flight, how much is left, and what can start now
---

Read-only. Reports; chooses nothing, changes nothing.

```
Skill(skill="find-work")
```

Invoke it in **reporting mode**: run source 1 only — scan `.agents/changes/*/plan.md`,
count open boxes, report what remains and what each plan says can start. Then stop.

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

## What not to do

**Do not offer to pick one up.** That is `/libretto-flow`'s job, and it asks first.
Reporting that quietly becomes an offer is how a read-only command starts writing.
