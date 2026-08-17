---
name: write-tasks
description: "Trigger: the spec and the plan are written and the work needs an ordered checklist; tracking task state across agents; marking a task done. The seam between phases 5 and 6. Launches one fresh cutter, then owns the file it returns."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

The seam between phase 5 and phase 6: turn the spec and the plan into the single file
that says what is done and what is next.

`tasks.md`. Not `plan.md` — that is the technical approach now, and this file used to
carry its name.

**Delegate the shape.** `skills/writing-plans/` already covers task right-sizing,
checklist structure, the no-placeholders rule and self-review. Invoke it for how a
checklist is written — it ships with this repository, see `THIRD-PARTY.md`.

**What it decides that this flow overrides.** Delegation without stated overrides is two
skills quietly disagreeing:

| It says | Here | Why |
|---|---|---|
| save to `docs/superpowers/plans/<date>-<name>.md` | `.agents/changes/<change>/tasks.md` | the checklist belongs inside the change it belongs to, beside the proposal, the delta and the plan, and it disappears with them when the change lands |
| the file is the handoff to execution | the file is **live state** | several agents read it while the work runs, so it is written continuously, not handed over once |
| the worktree came from `superpowers:using-git-worktrees` | `using-git-worktrees` | the payload ships that skill unnamespaced. The namespace names a plugin the user may not have; the skill itself is right here |
| execution needs `superpowers:subagent-driven-development` or `superpowers:executing-plans` | `build-and-check`, phase 6 | **neither is shipped, by anybody, on a machine that installed only this.** A checklist whose header demands a skill that does not exist is a checklist that stalls on its own first line |

**Write the header without those two names.** `writing-plans` mandates a header naming
`superpowers:subagent-driven-development`; here the file is live state that phase 6
reads, so the header names `build-and-check` and nothing else. Same for the
execution-handoff question at the end — the flow already routed, so there is no choice
left to offer.

Everything else it says stands. What follows is only what it does not cover.

**The vendored copy is not edited to say this.** `THIRD-PARTY.md` keeps it
byte-comparable with upstream, so a divergence has exactly one place to live and it is
this table. The repository's payload gate enforces it: a namespaced skill a vendored copy
cites, which the payload does not ship, has to be answered here or the gate fails.

## The cut runs in a fresh subagent

**This is the rule this seam was created for.**

Launch one `task-cutter` subagent. It is given the spec and the plan by path, and
nothing else — none of the conversation, none of the design argument, none of what was
considered and dropped.

The reason is the same one the 6→7 review seam already documents, one phase earlier. The
session that argued its way to an approach remembers what it *considered* as vividly as
what it *chose*, and cuts boxes against the argument rather than against the document. It
writes "wire up the adapter" and knows what that means; nobody else does, including the
fresh session `libretto loop` starts tomorrow.

A cutter with no memory can only cut what was actually written down. **That doubles as
the only check this flow has that the plan says enough to be built from** — and it is
the reason the cutter returns two things, not one.

### What comes back

The checklist, and **what the spec and the plan failed to answer.**

The second half is not a courtesy. A cutter that cannot produce a box because the plan
never said how something gets built has found a defect in the plan, at the last moment
it is cheap to fix. Read it before writing the file. A gap goes back to phase 5 or to
the contract; it never gets papered over with a box that says "figure out X".

### The orchestrator writes the file

The cutter returns the checklist. The orchestrator writes it, and from that moment owns
it.

## One writer

**The orchestrator owns the file.** Sub-agents never edit it — not the cutter, not the
builders. They report what they finished; the orchestrator marks the box.

Several agents editing one markdown concurrently lose updates — two reads of the same
content, two writes, the second silently discarding the first. The symptom is a finished
task whose box is empty, and nobody notices until something downstream waits forever for
work that was already done.

One writer removes the problem instead of managing it. No locks, no merge logic, no
retry.

## The checklist is derived, not invented

Every task traces to the spec, through the plan. If a task is in the file and in neither,
either the contract is incomplete or the task is scope that arrived without asking — find
out which before writing it down.

Carry the link both ways: each task names what it comes from, and each task names **the
verification criterion that closes it**. A task whose criterion is already met is not a
task.

## State is written when it changes, not at the end

A box is marked the moment its task is genuinely finished — verified, not hoped, per
`skills/evidence/`. Not batched, not at the end of a session.

An agent joining late reads this file and believes it. That is its whole purpose, and a
checklist updated in batches is a checklist that was wrong for most of its life.

**A marked box that was never committed was never marked.** It ships in the same commit
as the task that closed it, exactly as the spec does — a `git add` scoped to the code
leaves this file behind in the working tree, and a change that lands deletes the folder
with every unrecorded mark still in it. That is not hypothetical: the reviewer of the run
that added this line read six commits and found 0/24 boxes ticked on a checklist whose
work was finished, because every mark lived only in a working tree.

Phase 1 reads unchecked boxes to decide what is in flight, so a file that never moves
reports as fully open right up until it disappears.

Record three things per task and nothing more: done or not, where the evidence is, and —
if it was stopped — why. Discussion belongs in the spec; reasoning belongs in the plan.

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

That is a finding about the **spec or the plan**, not a licence to ship a layer. The rule
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

`.agents/changes/<change-name>/tasks.md` — inside the change, beside the proposal, the
spec delta and the plan it implements.

That placement is deliberate: it is as temporary as the change. When the change lands and
its delta is folded into the capability spec, this goes with it. A checklist that
outlived its change becomes a list of things that may or may not still be true.

**`libretto loop` and `libretto metrics` read a change created before 2026-08-17 by its
old name, `plan.md`.** Nothing writes that name any more. A change in flight keeps
working; a new one never produces it.

Follow the project's layout if it differs; `skills/write-spec/` has the detection order.

One checklist per change. A second file tracking the same work is two sources of truth,
and they will disagree exactly when it matters.

## Output

Report where the file is, how many tasks, which can start immediately, and **anything
the cutter said the spec or the plan failed to answer** — that last one unedited, in the
cutter's words.

Then stop. Cutting the tasks is not starting the work.

**This is the flow's second stop and the last one before code runs**, and it carries
phase 5's decision as well as this one — `write-plan` deliberately does not stop, so the
approach is agreed here, alongside the boxes it produced. That is why the report has to
carry the plan's approach in a sentence too: a user agreeing to an order without seeing
the reasoning it came from is agreeing to half of it.

**Ask for the go-ahead with `AskUserQuestion` — in conversation where the native prompt
does not exist — never as a sentence at the end of the report.** Start the work —
recommended, and saying which task runs first — change the approach or the order first,
or go back to the contract. **And room to answer differently**: the thing wrong with an
order is often not one of the three ways it could be wrong.

Same rule as every stop. The reason lives once, in `skills/record-work/`.

Unless the run is `/libretto-attacca`, where the invocation already agreed the order. The
file is written and committed exactly the same; only the wait is answered.
