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
