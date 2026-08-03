---
name: record-work
description: "Trigger: a task is finished and needs committing; writing a commit message; the spec no longer matches the code; deciding whether to push. Phase 8 of the flow. Records the work and reconciles the contract with what was actually built."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Phase 8 of the Libretto flow: put the work into history, and make the spec true again.

Read `skills/evidence/` first. Nothing here is reported as recorded that was not
observed to be recorded.

## One commit per task

Every finished task commits. Not batched at the end of a session, not one commit for
six tasks.

The reason is not tidiness. A history of one commit per task is a history you can
bisect: when something breaks two weeks from now, the commit that introduced it is
the commit that named it. A single commit containing six tasks tells you the day and
nothing else.

Work-in-progress commits during phase 6 are fine and expected — the tree under test
has to be the tree that was recorded. They squash into the task's commit.

## The spec ships with the code

**This is the rule that makes the rest of the flow mean anything.**

Implementation always learns something the spec did not know: a constraint that was
wrong, an outcome that needed splitting, a criterion that turned out untestable as
written. When that happens, the spec is updated **in the same commit as the code that
taught it**.

Not in a follow-up. Not in a cleanup pass. Not "noted for later". The same commit.

A spec updated separately is a spec that was wrong for however long the gap lasted,
and anyone reading it during that window was misled by a document that looked
current. This repository already carries three such divergences accumulated over one
phase of work — which is what it looks like when this rule is aspirational instead of
mechanical.

Before committing, ask it explicitly: **did this change teach the spec anything?** If
yes, the spec change is part of the commit. If no, say so and move on — the question
takes five seconds and the drift it prevents takes an afternoon.

While the work is in flight, "the spec" means the delta inside
`.agents/changes/<change>/spec.md`. That is where amendments belong until they land.

## Landing a change consolidates it

The last commit of a change does three things together, in one commit:

1. the final code
2. **the delta applied onto the capability spec** it targets, in
   `.agents/specs/<capability>/spec.md`
3. **the change folder deleted** — proposal, delta and plan

All three or none. A delta applied without deleting the change leaves two documents
describing the same capability, and the next reader has no way to tell which one is
current. A change folder deleted without applying the delta loses the work
outright.

Applying is not copying. The delta says what changes; the capability spec has to
read afterwards as though the feature had always been there. Requirements merge into
the existing numbering, `Governs:` gains any new paths, `Proof:` citations come
across, and contradicted sentences are rewritten rather than left beside their
replacement.

Then verify the anchors before the commit lands:

```
~/.claude/skills/record-work/spec-drift --anchors
```

A citation that survived consolidation pointing at a test that did not is the most
common way this goes wrong.

If a change spans several capabilities, every delta is applied in that same commit.
Half-consolidated is the one state with no honest description.

**`spec-drift` asks it for you**, mechanically, in two directions. It ships beside
this file — `~/.claude/skills/record-work/spec-drift` once installed — so it is
there in any project, not only in the repository it came from. Run it from the root
of the project being worked on:

```
~/.claude/skills/record-work/spec-drift             staged code whose spec did not move
~/.claude/skills/record-work/spec-drift --anchors    every Proof: citation resolves
```

The first reads each spec's `Governs:` globs, so it names the spec and the path
rather than guessing about the repository as a whole. Paths no spec governs are
reported separately and softly — not everything needs a contract.

The second is the reverse direction, and it is the one that catches rot: a
criterion citing a test that was renamed, deleted or never written. It checks the
**test name**, not just the file, because a file-level check passes an invented
name. Run it before the commit that would ship the lie.

**It warns; it does not block.** Always exits 0. A check that stops a commit in
someone else's project gets removed the same day, and then there is no check at all.
A deliberate no is a valid answer — say it out loud in the report so the next reader
knows the question was asked and not skipped.

Anyone who wants it to be a gate can wire it into their own `pre-commit` hook or CI.
That is their decision to make, not this flow's to make for them.

The plan is updated in the same breath, per `skills/write-plan/` — by the
orchestrator, never by a sub-agent.

## Messages

Conventional commits. `type(scope): subject` — imperative, lowercase after the type,
no trailing period.

The body explains **why**, only when the why is not obvious from the subject. What
changed is already in the diff; the diff cannot say what the alternative was or why
it lost.

Reference the tracker key so the commit points back at its origin.

**No AI attribution.** No `Co-Authored-By` for a model, no generated-with trailer.
The work is the author's.

If `caveman-commit` is installed and the user prefers it, use it for the message —
it produces the same shape, compressed.

## Branches — the backstop, not the decision

**`skills/build-and-check/` owns the branch**, at its step 0, before the first file is
written. By the time this phase runs the work already exists, so what happens here is
a second look at a cheap invariant — not a decision being made for the first time.

Never commit directly to the base branch. If the current branch is the base, stop and
create one, because the cost of getting it wrong is a rewrite of shared history.

But say plainly that it should not have got this far: reaching phase 8 on the base
branch means phase 6 skipped its own first step, and that is worth reporting rather
than quietly fixing. A backstop that silently covers for the rule it backs up is how
the rule stops being followed.

A branch per parent task, or per subtask when subtasks are genuinely independent.
Independent branches get chained rather than raced at the trunk. `chained-pr` if it is
installed.

## Pushing is the user's decision

**Never push unasked.** Not as a convenience, not because it seems obviously wanted,
not because the previous task was pushed.

At the very end — after everything is committed and reported — ask once: push, yes or
no. One question, then respect the answer. A no is a complete answer and needs no
follow-up.

When the answer is yes, confirm it landed: the remote tip matches the local tip. A
push that printed no error is not a push that was accepted.

Opening a pull request is a separate question, asked separately.

## Before the last word

Nothing is reported as recorded without having been seen:

- the commit exists — `git log` was read, not assumed
- the tree is clean, or what remains uncommitted is named and explained
- the spec matches the code, or the divergence is stated
- if pushed, the remote tip matches

Then one line per task: what it was, its commit, where its evidence is.

## Output

What was committed, on what branch, and whether the spec moved with it.

Then the push question. Then stop.
