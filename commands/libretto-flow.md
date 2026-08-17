---
description: Run the Libretto flow — read a task, spec it, plan it, build it, record it
---

Entry point for the flow in `docs/FLOW.md`.

Input: `$ARGUMENTS` — a Jira issue key, an issue URL, a board URL, or empty.

Invoke each skill with the `Skill` tool. Do not read their `SKILL.md` files directly
and do not inline their logic here.

**`Skill(skill="…")` and `AskUserQuestion` below are Claude Code's spellings, and this
payload installs into three hosts.** Every one of them means a capability, not a product:
load that skill with the host's own skill tool, and put a question to the user with the
host's own native prompt — or in conversation where there is none. **The stop is the rule,
never the widget**, and `record-work` states that once for the whole flow. Said here once
rather than beside each of the twelve mentions below: a routing table annotated twelve
times stops being scannable, which is the only reason it is a table.

## Standing rules

`Skill(skill="evidence")` applies at every phase, not at one of them. Nothing is
reported that was not observed, failing tests are fixed or reported and never edited
into silence, and commands whose result decides the next step run in the foreground
with their output read.

`ponytail` ships with this flow and applies throughout — how much gets built is its
question, and phase 2 is where it earns the most.

### A phase is invoked even when it has nothing to do

**Never pre-empt a phase's own judgment.** `write-spec` decides whether a spec is
needed; `build-and-check` decides how much to check and where the work lands. Those
decisions belong to those skills, and an orchestrator that makes them itself and
proceeds has not saved a step — it has moved a decision somewhere nobody can see it.

So invoke every phase the work reaches, **including the ones whose answer is "nothing
here"**, and report the declining in one line: *phase 2: no spec — a two-line
documentation edit, nothing to disagree about.*

The cost is one line. What it buys is the difference between a skip and an omission,
which from outside are identical. That is not hypothetical: a session made exactly the
right call about a small change, never invoked phases 2 or 6 to make it, and the user
asked why the flow had not run. It had, in substance. Nothing said so.

Invoking is not gating. A phase that declines does not add a wait — see the trivial
lane below.

### The flow stops in three places

**Two of them, and only two, are inside the work:** after the spec, and after the plan's
tasks are cut. Everything else runs through.

| After | Stops | Because |
|---|---|---|
| 1 · find-work | no | the reading is stated; the place to disagree with it is the spec |
| 2–3 · write-spec | **yes** | the contract, and it is cheapest to change here |
| 5 · write-plan | no | the approach is agreed at the stop below, with the boxes it produced |
| 5→6 · write-tasks | **yes** | the approach, the order, and what the cutter found missing |
| 6 · build-and-check | no | |
| 6→7 · review-work | no | it fixes what it finds rather than asking |
| 7 · present-work | no | the report and the commits land together |
| 8 · record-work | **yes, last** | push and open the request — the user's, never assumed |

Two exceptions in phase 1, neither of them ceremony: **work already in flight** is a
choice about the user's priorities, and a **missing or unconfigured tracker** leaves
nothing downstream to do. Both are the input failing to arrive, not a phase transition.

A stop is where the user changes something. A stop where the only available answer is
"yes, carry on" is a round trip charged for a rubber stamp, and a flow that charges
those gets routed around.

**All three are asked with `AskUserQuestion`, never as a sentence at the end of a report** —
and so is phase 1's choice between work in flight and something new. The recommended option
first, saying what will actually run; the real alternatives; room to answer differently. The
reason lives once, in `skills/record-work/`.

`/libretto-attacca` is the same flow with all three answered in advance, and it owns that
classification — including what stays a stop under it, because a gate is not one.

## 1 · Find the work

```
Skill(skill="find-work")
```

Three sources, in order: **a change already in flight**, then a tracker key or URL, then
what the user said. Home first — starting something new while a change sits half-finished
is how the half-finished thing gets abandoned.

A request in conversation is a legitimate input, not a fallback. It is how most work
arrives.

If it stops — a tracker was named and its CLI is missing, unconfigured or unauthorised —
stop here. Nothing downstream works without knowing what the work is.

Report what was found and **carry on into phase 2 in the same turn**. Stating the reading
is what makes a wrong one visible; waiting for a yes is what makes it expensive, and the
spec is where a wrong reading gets caught.

Two things it stops for, and both are the input failing to arrive: something already in
flight, where continuing it or not is the user's call, and a tracker whose CLI is missing.

`/libretto-status` runs the same source 1 and stops there, for when the question is only
"what is open?"

## 2–3 · Spec it

```
Skill(skill="write-spec")
```

Decides whether the task needs a spec at all, whether it is one spec or several,
reads the code before writing, fans out with a shared brief when there are many,
asks what it cannot settle, and writes the six pillars.

**It may report that no spec is needed.** That is a legitimate outcome for a
one-line change — go straight to phase 6, then 8.

Then **wait** for the go-ahead. This is the first of the two stops, and the one that
matters most: everything downstream is measured against what gets agreed here.

### The trivial lane

**A "no spec needed" collapses both remaining stops, not only the spec.** That route is
phase 6, then 7 and 8 **in the same turn**, and exactly one question at the end: push
and open the request. There is no plan either — a change with no contract has nothing to
break into ordered tasks.

Everything still gets said. Phase 7 reports what was done, its evidence and what was
left out; phase 8 commits per task.

So the gear is: **two stops with a contract, none without, and phase 8's question in
both.**

### Before the stop, read the contract

```
Skill(skill="review-spec")
```

**Only when a spec was written** — the trivial lane has no contract to review, and
running it there is ceremony on a one-line change.

Every other review in this flow reads code, and all of them run after the code exists.
This one reads the promise: ambiguity that forks the work, criteria no run could fail,
`Governs:` boundaries that describe nothing, a criterion another capability contradicts.
It reports; phase 3's author decides.

**It runs before the stop, not after.** The stop is where the contract is agreed, and
agreeing to a criterion nothing can fail is the failure this phase exists to prevent —
so the findings have to be on the table while the agreement is still being made.

## 5 · Plan it

```
Skill(skill="write-plan")
```

Writes the **how**: technical context, the approach and the alternatives it beat, risks,
validation, rollback. No boxes — this file holds no state.

**Do not wait here.** Splitting the plan from the checklist invited a third stop, and the
flow spends its length arguing against exactly that. The stop moves one seam later
instead, where the user sees the approach *and* the boxes it produced *and* whatever the
cutter found missing — the same single decision, made with more on the table.

## 5→6 · Cut the tasks

```
Skill(skill="write-tasks")
```

Launches one fresh `task-cutter` with the spec and the plan and nothing else, then writes
the checklist it returns. One writer: the orchestrator marks boxes, sub-agents report.

**Not a numbered phase, on purpose** — the same answer the 6→7 review already gives.
Independence comes from the fresh context, not from renumbering everything that says
eight.

Report what the cutter said the documents failed to answer, in its words. A gap goes back
to phase 5 or to the contract; it never becomes a box that says "figure out X".

Then **wait**. The second stop and the last one before the work runs — what is being
agreed is the approach, the order, and what waits on what.

## 6 · Build it

```
Skill(skill="build-and-check")
```

Implements one task at a time with proportionate checks. Decides branch versus
worktree. Commits work-in-progress before verifying, because a gate proves something
only about the tree it ran against.

Two failed gates on one task stops the task. Two stopped tasks in a row stops the
session and updates `docs/STATE.md`.

## 6→7 · Review it

```
Skill(skill="review-work")
```

Hands the finished work to a fresh `work-reviewer` subagent that saw none of this
session — it reads the contract and the diff, re-runs the proofs itself, and returns
findings. **Then the seam fixes every one of them, without asking**, and re-runs the
proofs each fix touched.

The reviewer itself never writes. Read-only is what makes the second pair of eyes worth
having, and the repair happens on the other side of the seam.

A finding that fails to fix twice, or that needs a decision that is not ours, goes to
phase 7 as found-and-not-fixed. It is reported there, never turned back into a question.

A change with no spec gets a one-line decline and no wait: no contract, nothing to
review against. The same invoked-even-when-empty rule as every phase.

## 7 · Present

```
Skill(skill="present-work")
```

Shows what was done in the spec's terms, names the evidence, and states **what was
deliberately left out with the condition that would bring it back** — and carries the
reviewer's verdict, attributed, next to the builder's own account — what it found, and
what was done about each one.

Then **carry straight into phase 8**, same turn. The report is not a gate; it is the
context the last question gets asked in.

## 8 · Record it

```
Skill(skill="record-work")
```

One commit per task, the spec updated in the same commit as the code that taught it,
conventional messages, no AI attribution.

Pushing is asked once at the very end, with `AskUserQuestion` — never as a sentence at
the bottom of a report. Never assumed.

**When the answer is yes and the request is open, the tree returns to the base branch and
brings it up to date** — `--ff-only`, the feature branch left alone. A session starts
wherever the last one left the working directory, and a stale base is what phase 1 reads
next time.

## 4 · Asking — before the plan, not after it

When something is genuinely not yours to decide — a product tradeoff, two live
precedents in the codebase, anything where guessing quietly breaks working behaviour —
ask with `AskUserQuestion`: the option you recommend, the real alternatives, and room to
answer differently. One question, then stop and wait.

**Phase 2 is the exception, and it asks as many as a wrong guess would cost, in one
call, before it writes the spec** — the contract gets built by two people rather than
handed over finished, and the bias when in doubt is to ask. No hard cap, but judgment
both ways: never a form-length interrogation of things the code already answers. Zero
is legitimate and is said in one line; a quota manufactures questions the code already
answers. The 5→6 seam stops for the approach and the order, and opens no second tranche, and no stop is
added: the questions ride the one phase 2 already has.

**Ask it before the plan is agreed.** Phases 1, 2 and 5 are where a question is cheap,
because nothing has been built on the answer yet. That is not a coincidence — the two
stops sit there for the same reason.

**After the plan, a question becomes a finding.** Anything discovered in phase 6 or in
the review seam that cannot be settled from the code goes into the phase 7 report as an
open decision, with what was assumed in the meantime and what would change if the
assumption is wrong. The user meets it at phase 8, in one place, with everything else.

The alternative is the flow interrupting mid-build to ask something the user cannot
answer without the context the build just produced — which is how a three-stop flow
becomes a nine-stop one, one reasonable exception at a time.

The answer goes into the spec, under prior decisions, next to what it settled. An answer
that lives only in the conversation gets asked again next session.

Do not ask what the code can tell you.
