---
name: write-plan
description: "Trigger: specs are written and the work needs an ordered checklist; tracking task state across agents; marking a task done. Phase 5 of the flow. Turns the task breakdown into the one file that holds live state."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Phase 5 of the Libretto flow: turn the specs' task breakdown into the single file that
says what is done and what is next.

**Delegate the shape.** `skills/writing-plans/` already covers task right-sizing, plan
structure, the no-placeholders rule and self-review. Invoke it for how a plan is
written — it ships with this repository, see `THIRD-PARTY.md`.

**Two things it decides that this flow overrides.** Delegation without stated
overrides is two skills quietly disagreeing:

| It says | Here | Why |
|---|---|---|
| save to `docs/superpowers/plans/<date>-<name>.md` | `.agents/changes/<change>/plan.md` | the plan belongs inside the change it belongs to, beside the proposal and the delta, and it disappears with them when the change lands |
| the plan is the handoff to execution | the plan is **live state** | several agents read it while the work runs, so it is written continuously, not handed over once |

Everything else it says stands. What follows is only what it does not cover.

## The plan is derived, not invented

Every task traces to a spec. If a task is in the plan and in no spec, either the
spec is incomplete or the task is scope that arrived without asking — find out
which before writing it down.

Carry the link both ways: each task names the spec it comes from, and each task
names **the verification criterion that closes it**. A task whose criterion is
already met is not a task.

## One writer

This is the rule the whole file exists for.

**The orchestrator owns the plan.** Sub-agents never edit it. They report what they
finished; the orchestrator marks the box.

Several agents editing one markdown concurrently lose updates — two reads of the
same content, two writes, the second silently discarding the first. The symptom is
a finished task whose box is empty, and nobody notices until something downstream
waits forever for work that was already done.

One writer removes the problem instead of managing it. No locks, no merge logic, no
retry.

## State is written when it changes, not at the end

A box is marked the moment its task is genuinely finished — verified, not hoped, per
`skills/evidence/`. Not batched, not at the end of a session.

An agent joining late reads this file and believes it. That is its whole purpose,
and a plan updated in batches is a plan that was wrong for most of its life.

**A marked box that was never committed was never marked.** It ships in the same commit
as the task that closed it, exactly as the spec does — a `git add` scoped to the code
leaves the plan behind in the working tree, and a change that lands deletes the folder
with every unrecorded mark still in it. That is not hypothetical: the reviewer of the run
that added this line read six commits and found 0/24 boxes ticked on a plan whose work was
finished, because every mark lived only in a working tree.

Phase 1 reads unchecked boxes to decide what is in flight, so a plan that never moves
reports as fully open right up until it disappears.

Record three things per task and nothing more: done or not, where the evidence is,
and — if it was stopped — why. Discussion belongs in the spec.

## Order by dependency, not by convenience

Tasks that unblock others come first. Shared foundations before the things that
build on them.

Mark what each task waits on. Without that, two agents pick tasks that touch the
same ground and one of them wastes its work. With it, "what can start now?" has an
answer anyone can read off the file.

Independent tasks may be marked as such. Parallel by default is how a plan becomes
a race.

## Where it lives

`.agents/changes/<change-name>/plan.md` — inside the change, beside the proposal
and the spec delta it implements.

That placement is deliberate: the plan is as temporary as the change. When the
change lands and its delta is folded into the capability spec, the plan goes with
it. A plan that outlives its change becomes a list of things that may or may not
still be true.

Follow the project's layout if it differs; `skills/write-spec/` has the detection
order.

One plan per change. A second file tracking the same work is two sources of truth,
and they will disagree exactly when it matters.

## Output

Report where the plan is, how many tasks, and which can start immediately.

Then stop. Writing the plan is not starting the work.

**Ask for the go-ahead with `AskUserQuestion`, never as a sentence at the end of the
report.** Three options: start the work — recommended, and saying which task runs first —
change the order first, or go back to the contract.

Same rule as phase 2's stop and phase 8's question, for the reason `record-work` gives
once: a question in prose is a paragraph the reader skims, and a flow waiting on an answer
nobody realised was a question looks hung.

Unless the run is `/libretto-attacca`, where the invocation already agreed the order. The
plan is written and committed exactly the same; only the wait is answered.
