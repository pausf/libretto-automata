---
name: write-tasks
description: "Trigger: the contract and the approach are written and the work needs an ordered checklist; tracking task state across agents; marking a task done. The seam between phases 5 and 6. Writes the change document's checklist inline, then owns it."
license: MIT
metadata:
  author: pausf
  version: "2.0"
---

## What this does

The seam between phase 5 and phase 6: turn the contract and the approach into the one
list that says what is done and what is next.

It is written under a level-2 `## Checklist`, the change document's last section, in
`proposal.md` — boxes and nothing else. Not `## Approach`, which holds no box, and not the
*Task breakdown* pillar, which is a list the checklist is cut against.

**The shape, whole.** One box per end-to-end unit of work, in dependency order. Each box
says what it is and carries a `Closes:` naming the criteria it closes, as written under
`## Verification criteria` — the token `converge` maps a box on — and, when it waits on
another box, which one:

```
- [ ] one field, stored, served and proved · Closes: the field round-trips
- [ ] the second field · waits on box 1 · Closes: both fields validate
```

No headers, no per-step sub-boxes, no code blocks inside a box: the approach says how,
and `libretto loop` counts every box as a task.

## Written inline, by the session

The session that held the interview and the fork cuts the boxes itself and writes them
under `## Checklist`, the change document's last section. No cutter is launched.

The blind spot that used to justify one is real: the session that argued its way to an
approach remembers what it *considered* as vividly as what it *chose*, and cuts boxes
against the argument. It writes "wire up the adapter" and knows what that means; nobody
else does, including the fresh session `libretto loop` starts tomorrow. **It is caught by
reading, not by a second author.** Once this section exists, `review-spec` hands the whole
document to one reader that never had the conversation, and what the contract and the
approach fail to answer comes back per box, before the stop.

A box you cannot cut because the approach never said how something gets built is the
same finding, found early: take it back to phase 5 or to the contract. It never gets
papered over with a box that says "figure out X".

## One writer

**The session owns the checklist.** Sub-agents never edit it — not a writer, not the
reviewer, not the builders. They report what they finished; the session marks the box.

Several agents editing one markdown concurrently lose updates — two reads of the same
content, two writes, the second silently discarding the first. The symptom is a finished
task whose box is empty, and nobody notices until something downstream waits forever for
work that was already done.

One writer removes the problem instead of managing it. No locks, no merge logic, no
retry.

## The checklist is derived, not invented

Every task traces to the contract, through the approach. If a task is in the checklist and
in neither, either the contract is incomplete or the task is scope that arrived without
asking — find out which before writing it down.

Carry the link both ways: each task names what it comes from, and each task names **the
verification criterion that closes it**. A task whose criterion is already met is not a
task.

## State is written when it changes, not at the end

A box is marked the moment its task is genuinely finished — verified, not hoped, per
`skills/evidence/`. Not batched, not at the end of a session.

An agent joining late reads this checklist and believes it. That is its whole purpose, and a
checklist updated in batches is a checklist that was wrong for most of its life.

**A marked box that was never committed was never marked.** It ships in the same commit
as the task that closed it, exactly as the spec does — a `git add` scoped to the code
leaves the document behind in the working tree, and a change that lands deletes the folder
with every unrecorded mark still in it. That is not hypothetical: the reviewer of the run
that added this line read six commits and found 0/24 boxes ticked on a checklist whose
work was finished, because every mark lived only in a working tree.

Phase 1 reads unchecked boxes to decide what is in flight, so a checklist that never
moves reports as fully open right up until it disappears.

Record three things per task and nothing more: done or not, where the evidence is, and —
if it was stopped — why. Discussion belongs in the contract; reasoning belongs in the
approach.

## Order by dependency, not by convenience

Tasks that unblock others come first. Shared foundations before the things that build on
them.

Mark what each task waits on. Without that, two agents pick tasks that touch the same
ground and one of them wastes its work. With it, "what can start now?" has an answer
anyone can read off the file.

Independent tasks may be marked as such. Parallel by default is how a checklist becomes
a race.

## Cut each box so it stands alone

Ordering says which box comes first. **This says what belongs inside one**, and it is the
harder half: a checklist can be perfectly ordered and still be cut wrong.

A box is one end-to-end change — the code and the proof that closes it, together. Not a
layer of one. Closing a box leaves the tree green and mergeable **on its own**, with
nothing half-built waiting for the next box to make sense of it.

So the test is not "is this task small" but: **could this box merge by itself?** If the
answer is no because the box after it is what makes it work, they are not two boxes. Two
boxes where the first only makes sense once the second lands are one badly cut box, and
the fix is to merge them rather than to order them.

What the horizontal cut looks like, because it is the default and it reads as tidy:

| Cut | Boxes | Merge the first one alone and you get |
|---|---|---|
| by layer | the model · the handler · the tests | a model nothing calls, then a handler nothing proves |
| end to end | one field, stored, served and proved · the next field | a working field |

**The capability spec is not part of a box.** A delta lands on it once, in the final
commit, exactly as it always has — what a box owes is its own proof, never a slice of the
delta.

### When a box cannot be cut that way

That is a finding about the **contract or the approach**, not a licence to ship a layer. The rule
above about tasks that trace to neither already says what to do: find out which end is
wrong before writing it down. Take it back.

### What this buys, and it is not tidiness

`libretto loop` runs one fresh session per open box. A box that does not stand alone
leaves the tree broken between sessions — the next session opens on work it did not do
and cannot see the reasoning for, which is the failure the loop is most exposed to and
the one it has no way to detect.

**Ceiling named:** nothing here can check a checklist. Whether a box genuinely merges
alone is judgment, and the payload's own gate only verifies that this mandate is present
in this file. A box cut horizontally under a mandate sitting three paragraphs above it
surfaces at phase 6 or in the 6→7 review, never in a gate.

## Where it lives

`## Checklist` in `.agents/changes/<change-name>/proposal.md` — the last section, after the
contract and the approach it implements.

That placement is deliberate: it is as temporary as the change. When the change lands and
its contract is folded into the capability spec, this goes with the folder. A checklist
that outlived its change becomes a list of things that may or may not still be true.

**Only the boxes under `## Checklist` are the checklist.** `libretto loop`,
`libretto metrics` and phase 1 count nothing else in `proposal.md`, so a box anywhere above
it is not a task — and under `## Verification criteria` it is read as a criterion.

**A change begun in the old layout keeps its checklist where it started**: `tasks.md`, or
`plan.md` for one created before 2026-08-17. Every reader still finds them, and nothing
writes either name for a new change.

Follow the project's layout if it differs; `skills/write-spec/` has the detection order.

One checklist per change. A second list tracking the same work is two sources of truth,
and they will disagree exactly when it matters.

## Output

Report where the checklist is, how many boxes, and which can start immediately.

**Then carry on to `review-spec`. This seam does not stop**, and the reason is the one
stop that replaced it: the user agrees to the contract, the approach and the order once,
with a fresh reader's findings beside them, and `review-spec` asks that question. Cutting
the boxes is not starting the work either way.
