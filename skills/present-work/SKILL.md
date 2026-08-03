---
name: present-work
description: "Trigger: a task is finished and needs reporting; summarising what was done; explaining what was left out; before committing or asking to push. Phase 7 of the flow. Shows the work in the spec's own terms, names its evidence, and states what was deliberately not built."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Phase 7 of the Libretto flow: show what was done, then stop.

It runs after `skills/build-and-check/` and **before** `skills/record-work/`. That
order is the point — presenting after the commit is an announcement; presenting
before it is a decision point where the answer can still be "no, not that".

`skills/evidence/` governs this phase more directly than any other. Read it first.
Every other phase produces work; this one makes claims about work, and a claim is
where an unobserved result becomes a lie.

## Three things

- **what was done, in the terms the spec used.** Not the terms of the code. If the
  spec called it "relative discounts" the report says relative discounts, not
  `DiscountPolicy`. A report that only makes sense to whoever wrote it has not
  reported anything.
- **where the evidence is** — the run, the test, the commit. Named, not described.
- **what was deliberately left out, and the condition that brings it back.**

## The third one is the reason this phase exists

It is the one that goes missing, and it is the one that turns a simplification into
a decision.

> Did the single-user case; the shared-state version needs a lock — say so if you
> need it now.

That is reviewable. Silence about the same choice is a surprise waiting for whoever
reads the code next. The difference between the two is a sentence.

Each omission gets both halves — **what was not built, and what would change that**:

| Not built | Brings it back |
|---|---|
| retry on the upload | a second report of a dropped file |
| pagination on the list endpoint | more than ~200 rows in production |
| the migration for old rows | anything reading rows written before today |

An omission with no condition is a confession, not a decision. It tells the reader
something is missing and leaves them to work out whether it matters.

If `ponytail` is installed this is where its skipped work surfaces. It answers how
much gets built; this phase is where what it declined becomes visible instead of
silently absent.

## Every claim points at something that happened

The whole report is subject to `skills/evidence/`:

- "tests pass" names the run and its output was read
- "it works" names what was executed, not what should work
- a task reported closed is one whose plan criterion was met and watched
- anything not observed is reported as not observed — **not omitted**

A task that was stopped is presented as stopped, with what failed and what is still
unknown. Two failed gates stop a task by `skills/build-and-check/`; hiding that in a
report undoes the rule.

**A phase that reports work is the last place where an unverified claim is still
cheap to correct.** After this it is in a commit message and then in someone's head.

## Length

If the explanation is longer than the change, the explanation is the problem.

A three-line diff does not get four paragraphs of rationale. Long prose defending a
small change is usually a change that is not defensible, or a change nobody actually
understood.

Nothing here is a tour of the code. The diff is available; the reader can read it.
What they cannot read is what was decided and what was declined.

If `caveman` is installed, it governs how much gets said.

## The orchestrator presents

Sub-agents never do. Only the orchestrator saw the whole task, and a report assembled
from fragments that each only saw their own file is a report that cannot say what was
left out — because no fragment knows.

Sub-agents return findings; per `skills/write-spec/`, those returns are read before
the set is accepted, and what they surfaced belongs in this report.

## The saying is unconditional; the stopping is not

**Normally this phase ends by stopping.** Phase 8 begins when the user says it does,
because presenting before the commit is only a decision point if the answer can still
be "no, not that".

**One exception, and it is narrow: a change `skills/write-spec/` decided needed no
spec.** There, the report and the commit land in the same turn, and the only question
asked is phase 8's last one.

The reason the exception is safe is that there is no contract to disagree with. A
change too small to argue about what "done" means is a change too small for a gate
whose whole purpose is to catch a disagreement about "done". Four stops for a typo
teaches people to route around the flow, and they route around it for typos first.

**What is never conditional is the content.** All three things still get said — what
was done, where the evidence is, and what was left out with the condition that brings
it back. The collapse removes a wait, never a sentence. A phase that skips the saying
because the change was small is how the one omission that mattered goes unmentioned.

## Output

Per task: what it was, in the spec's terms · its evidence · closed or stopped.

Then the omissions, each with its condition.

Then stop — unless there was no spec, in which case carry straight on to phase 8.

Then **stop and wait**. Do not commit, do not push, do not start the next task.
Phase 8 begins when the user says it does — presenting exists so that "no" is still
an available answer.
