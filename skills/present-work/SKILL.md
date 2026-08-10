---
name: present-work
description: "Trigger: a task is finished and needs reporting; summarising what was done; explaining what was left out; before committing or asking to push. Phase 7 of the flow. Shows the work in the spec's own terms, names its evidence, and states what was deliberately not built."
license: MIT
metadata:
  author: pausf
  version: "1.1"
---

## What this does

Phase 7 of the Libretto flow: show what was done, then commit it.

It runs after `skills/build-and-check/` and **before** `skills/record-work/`, in the
same turn as both. The order is the point — a report written after the commit is
assembled from the commit, and it will only ever say what the commit already said.
Written before, it is still an account of the work rather than a description of the
diff, and what got left out is still in the author's head.

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

And when `skills/review-work/` ran in the seam before this phase, a fourth: **the
reviewer's verdict, attributed and unedited** — the proofs it ran, what it found, and
**what was done about each finding**, in its own section next to the account above.
The builder's report and the reviewer's are two sources and stay two; merging them is
how a finding gets softened into the prose it was questioning.

By the time this phase runs, most of those findings are already fixed. **Report them
anyway, with the repair beside each one** — fixed and re-run, or stopped after two
attempts, or reported-not-fixed because the decision was not ours. A section showing
only the repaired state hides that the builder got it wrong, and getting it wrong is
the thing the seam exists to surface.

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

## This phase says everything and stops for nothing

**The report is the whole job, and phase 8 follows it in the same turn.**

This phase used to end by waiting, on the reasoning that presenting before the commit
is only a decision point if "no, not that" is still available. It still is — it just
arrives one beat later, at phase 8's own question, with the commits written and the
report in front of the user. A commit on a local branch nobody has seen is not a thing
that needs undoing; it is the cheapest possible place to change your mind.

What the wait actually cost was a round trip on every single change, and the flow had
already conceded the point for changes with no spec. Conceding it once for typos and
holding it everywhere else is how a rule becomes a thing people route around.

The two stops that survive are earlier and better placed: the spec, where the contract
can still be argued with, and the plan, where the order can. By here both were agreed
and the work matches them — `review-work` fixed what it did not.

**What was never conditional is the content.** All three things still get said — what
was done, where the evidence is, and what was left out with the condition that brings
it back. Nothing about removing the wait removes a sentence. A phase that skips the
saying because the change looked small is how the one omission that mattered goes
unmentioned.

## Output

Per task: what it was, in the spec's terms · its evidence · closed or stopped.

Then the omissions, each with its condition.

Then carry straight on to phase 8, in the same turn. Do not push — that is phase 8's
one question and the last one in the flow.
