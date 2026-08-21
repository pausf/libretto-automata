---
name: write-plan
description: "Trigger: the spec is written and the work needs an approach before any code; choosing between two ways to build something; recording why an alternative lost; naming what could go wrong. Phase 5 of the flow. Writes the how — the document that used to not exist."
license: MIT
metadata:
  author: pausf
  version: "2.0"
---

## What this does

Phase 5 of the Libretto flow: decide **how** the change gets built, and write the
reasoning down before it evaporates with the session that had it.

The spec says what and why. This says how. The checklist that says in what order is
`write-tasks`, in the seam after this one — **it is not this file, and until 2026-08-17
it was.**

## Why this phase exists at all

It was reported missing. A Lead read this flow's output and found specs that read like
plans and plans that read like task lists, and the diagnosis was one line long: phase 5
transcribed the spec's task breakdown into checkboxes and called it a plan.

So the technical reasoning — the approach chosen, the two that were rejected, the thing
that will probably go wrong — happened in a session and died there. Every later reader
got the *what* twice and the *how* never, and the first person to ask "why is it built
like this?" got an answer reconstructed from the diff.

**A decision nobody wrote down gets made again**, and the second time it is made by
somebody with less context, under more pressure, in the middle of phase 6.

## The fork — the user chooses the approach

Before any document exists, phase 5 asks its one question: **two or three approaches
with tradeoffs, in one `AskUserQuestion` call — the host's native prompt, or in
conversation where one does not exist** — the recommended one first *with the reason
it is recommended*, the real alternatives with what each would cost, and room to answer
differently. Never a fait accompli: a finished decision can only be rubber-stamped,
while a fork with named costs must be decided, and deciding it is the user's.

Two or three means real ones. An approach nobody could have chosen is a strawman
propping up the recommendation, and a change with genuinely one way to build it says so
in one line and asks nothing — zero is as legitimate here as it is in phase 2.

**The chosen and the rejected go to `decisions.md`** — the choice verbatim, each
rejected approach with why it lost. That is what makes the alternatives table below
reconstructible outside this conversation, which is the whole reason the table used to
be the pillar that died with the session.

**This does not add a stop.** The fork rides inside phase 5 the way the interview rides
inside phase 2; the flow still stops in three places, and the approach is still agreed
at the 5→6 seam with the boxes it produced. Until 2026-08-19 this phase asked nothing —
a decision made deliberately and retired deliberately: the no-third-stop reason it was
made for survives, the silence did not.

**Under `/libretto-attacca` the fork is pre-answered**: the recommended approach is
taken and logged in `decisions.md` with the `(assumed)` suffix, naming what changes if
it is wrong — the same rule as every question attacca answers in advance.

## The drafting runs in a fresh subagent

The fork above and this launch are **the orchestrator's conduct** — a `plan-writer`
invoking this skill for the document's shape starts at *What goes in it* and follows
nothing above that heading.

With the approach chosen and logged, **launch one `plan-writer`**: the change's spec
delta, the capability spec it targets, `proposal.md` and `decisions.md`, each by path,
plus the standing rules restated verbatim — no commit, no push, no file written, its
questions in its return value. It drafts the plan the same way `task-cutter` cuts the
checklist one seam later, and for the same reason: the session that argued its way to
an approach cannot tell what it wrote down from what it merely decided, and it drafts
against the argument. The writer can only draft what the log actually carries — which
is the first check the fork's record gets.

**Read what it failed to answer before writing anything.** A gap goes back to the fork
or to the contract; it never gets drafted over. Then the orchestrator writes `plan.md`
from the returned markdown — one writer, the same rule as `tasks.md`.

## What goes in it

Six things. The fourth and the sixth are the ones that get skipped, and they are the two
worth the most in six months.

### Summary

What is being built and the shape of the approach, in a paragraph. Somebody who reads
only this should be able to say whether the rest is worth their time.

### Technical context

What the change is built against, concretely enough to be wrong: language and version,
the dependencies in play, where the tests live, what gates it has to pass, what is
generated rather than written.

**And the blast radius** — which files, how many. A change that names its own edges is
a change whose review can be scoped.

Vague context is worse than none. "Uses the existing test setup" tells a reader nothing
they could contradict.

### The approach

The decision, and the mechanism. Enough that somebody could build it from this document
without the conversation that produced it — which is exactly the test the task cutter
applies in the next seam, because it has nothing else.

### The alternatives it beat

**This is the pillar the flow was missing, and it is the one that pays.**

One row per real alternative: what it was, and why it lost. Not a strawman — something
that was genuinely considered, by somebody who could have chosen it.

| Considered | Why it lost |
|---|---|

The fork's rejected approaches are rows here by definition — they are in
`decisions.md` with why they lost, in the user's words, and the table carries them
forward. Alternatives that surfaced later, while the plan was being drafted, join them.

A diff can show what was built. Nothing in a repository can show what was *not* built
and why, so if it is not here it is gone. And the alternative that lost for a reason
that later stops being true is the single most valuable thing in this file: it is the
change somebody should make next, and without this table they will never know it was
already on the table.

**"It was obvious" is not a row.** If it was obvious there was no decision, and the
table is shorter. An empty table is legitimate and says so in one line.

### Risks

What could go wrong, and what catches it. One row each — the mitigation names a
mechanism, never an intention.

"Be careful with the migration" is not a mitigation. "The migration runs behind a
feature flag, and `TestMigrationIsReversible` proves the rollback" is.

The risk with no mitigation is still worth writing down. It becomes the thing phase 7
reports as accepted rather than the thing that surprises somebody.

### Validation and rollback

How the change gets proved, and how it gets taken back out.

Which gates carry it, which tests are new, and — the part that goes missing — **which
of them would pass for the wrong reason.** Code that only runs on a path the tests do
not take passes green and proves nothing. Say which test has to be forced red on purpose
before it is believed, per `skills/evidence/`.

Rollback is usually one line ("one revert, nothing migrates"). When it is not one line,
that is the finding.

### Complexity deliberately kept

Where the change is not the simplest thing that works, and why the simpler thing was
wrong. `ponytail:` in the code marks the site; this says what the site is for.

Absent this, the next reader deletes it — correctly, on the evidence they have.

## What does not go in it

- **No checkboxes.** No state, no progress, no "done". That is `tasks.md`, and a plan
  that grows a checkbox becomes a second source of truth about what is finished.
- **No task list.** Ordering, dependency and what-can-start-now belong to the next seam.
- **No restatement of the spec.** Link it. A copied requirement is a requirement that
  will disagree with its original.

## Where it lives

`.agents/changes/<change-name>/plan.md` — beside the proposal and the spec delta.

It is as temporary as the change. When the change lands, the delta folds into the
capability spec and this goes with it. What survives is what the *spec* records; a plan
that outlived its change would be a list of decisions that may or may not still hold.

**The decision worth keeping outlives the file, and this is now measured.** A choice that
will still constrain work after this change lands belongs in the capability spec's *Prior
decisions*, and phase 8 is where it moves. A plan is where a decision is made, never where
it retires.

`spec-drift --retired` — inside `--anchors`, one of the six gates — fails the landing
commit when a `plan.md` is deleted and no capability spec's *Prior decisions* section
moved with it. Not "the spec was edited": that happens in the same commit by definition,
since the delta lands there. **That section.**

So a plan carries one line, near the top, and it is written while the plan is:

```
Durable decisions: the two in Prior decisions below
Durable decisions: none
```

**The line is a claim about whether the list is empty. It is never the list.** What gets
retired is read off the change's `spec.md` *Prior decisions* — the line says only whether
there is anything there, which is all `--retired` greps it for. Written because the two
readings both looked right and produced a contradiction on the first landing: one plan
declared "the two" over a section holding three, and the other carried no line at all,
having been written before the line existed. Found by a 5→6 cutter reading both documents
cold, which is what that seam is for.

**`none` is legitimate and it is a claim.** A rename, a typo fix that grew a plan, a
change whose only alternative was doing nothing — those genuinely retire nothing. What
the gate cannot stop is `none` becoming a reflex, and no mechanism can: a plan with a
full alternatives table and `none` beside it is a contradiction only a reader catches.
`review-work` has the plan and the diff. Write it honestly, because the cost of a wrong
`none` is paid by whoever asks, two changes from now, why it was built this way.

Follow the project's layout if it differs; `skills/write-spec/` has the detection order.

## Output

Report where the plan is, the approach in a sentence, and any risk with no mitigation.

**Then carry on to `write-tasks`. This phase does not stop, and that is deliberate.**

Splitting the plan out of the checklist invited a third stop — one for the approach, one
for the order — and this flow spends its length arguing that a stop where the only
available answer is "yes, carry on" is a round trip charged for a rubber stamp. So the
stop moved one seam later rather than multiplying: the user is asked once, with the
approach *and* the boxes it produced *and* whatever the cutter found the documents failed
to answer, all on the table at the same time. One decision, better informed.

**Writing the plan is not starting the work**, and nothing here touches code.
