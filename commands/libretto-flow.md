---
description: Run the Libretto flow — read a task, spec it, plan it, build it, record it
---

Entry point for the flow in `docs/FLOW.md`.

Input: `$ARGUMENTS` — a Jira issue key, an issue URL, a board URL, or empty.

Invoke each skill with the `Skill` tool. Do not read their `SKILL.md` files directly
and do not inline their logic here.

## Standing rules

`Skill(skill="evidence")` applies at every phase, not at one of them. Nothing is
reported that was not observed, failing tests are fixed or reported and never edited
into silence, and commands whose result decides the next step run in the foreground
with their output read.

If `ponytail` is installed it applies throughout too — how much gets built is its
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

Report what was found, then **wait**. The user may want to look before a contract gets
written, and if something was already in flight they choose whether to continue it.

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

Then **wait** for the go-ahead.

### The trivial lane

**A "no spec needed" collapses the waits too, not only the spec.** That route is
phase 6, then 7 and 8 **in the same turn**, and exactly one question at the end: push
and open the request.

Everything still gets said. Phase 7 reports what was done, its evidence and what was
left out; phase 8 commits per task. What disappears is the *stop* between them, and
only because there is no contract for the user to disagree with — that is what phase 2
just established.

The four stops are for a change with a spec. Charging them to a typo is how a flow
gets routed around, and it gets routed around for typos first.

## 5 · Plan it

```
Skill(skill="write-plan")
```

Turns the task breakdown into the single file holding live state. One writer: the
orchestrator marks boxes, sub-agents report.

Then **wait**.

## 6 · Build it

```
Skill(skill="build-and-check")
```

Implements one task at a time with proportionate checks. Decides branch versus
worktree. Commits work-in-progress before verifying, because a gate proves something
only about the tree it ran against.

Two failed gates on one task stops the task. Two stopped tasks in a row stops the
session and updates `docs/STATE.md`.

## 7 · Present

```
Skill(skill="present-work")
```

Shows what was done in the spec's terms, names the evidence, and states **what was
deliberately left out with the condition that would bring it back**.

Then **wait**. Phase 8 begins when the user says it does.

## 8 · Record it

```
Skill(skill="record-work")
```

One commit per task, the spec updated in the same commit as the code that taught it,
conventional messages, no AI attribution.

Pushing is asked once at the very end. Never assumed.

## 4 · Asking — at every phase, not at one

When something is genuinely not yours to decide — a product tradeoff, two live
precedents in the codebase, anything where guessing quietly breaks working behaviour
— ask with `AskUserQuestion`: the option you recommend, the real alternatives, and
room to answer differently. One question, then stop and wait.

The answer goes into the spec, under prior decisions, next to what it settled. An
answer that lives only in the conversation gets asked again next session.

Do not ask what the code can tell you.
