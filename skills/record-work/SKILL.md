---
name: record-work
description: "Trigger: a task is finished and needs committing; writing a commit message; the spec no longer matches the code; deciding whether to push. Phase 8 of the flow. Records the work and reconciles the contract with what was actually built."
license: MIT
metadata:
  author: pausf
  version: "1.8"
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

While the work is in flight, "the spec" means the change's contract: the level-2 pillars in
`.agents/changes/<change>/proposal.md` and any `spec-<capability>.md` beside it — or, for a
change begun in the old layout, its delta `spec.md`. That is where amendments belong until
they land.

## Landing a change consolidates it

The last commit of a change does four things together, in one commit:

1. whatever code remains — often none, when every box committed its own
2. **the contract applied onto the capability spec** its `Targets:` names —
   `proposal.md`'s contract and each `spec-<capability>.md` beside it, each onto its
   own `.agents/specs/<capability>/spec.md`
3. **the durable decisions retired** into that spec's *Prior decisions* — the
   `## Prior decisions` entries that outlive the change
4. **the change folder deleted** — the change document and everything beside it

**An old-layout change lands as it always has**: its `spec.md` delta is what gets applied,
and its `plan.md` is where the durable decisions are read. Every reader still finds an
old-layout change, and this rule stays until nothing in flight uses that layout.

All four or none. A delta applied without deleting the change leaves two documents
describing the same capability, and the next reader has no way to tell which one is
current. A change folder deleted without applying the delta loses the work
outright.

**Step 3 is the one that gets skipped, and it is gated.** The approach holds why the
change was built the way it was — and the alternatives it beat — and step 4 destroys that.
It survives in git history, which is to say nobody reads it again.

So before deleting: read the approach's *alternatives* and the contract's *prior
decisions*, and move anything that will still constrain work after this lands. Not the
whole approach — a decision about how *this* change was sequenced dies with it correctly.
The one worth keeping is the one somebody would otherwise make again, differently.

`spec-drift --retired`, inside `--anchors`, fails this commit when a change document
carrying `Targets:` — or an old-layout `plan.md` — is deleted and no capability spec's
*Prior decisions* section moved with it. A change that genuinely retires nothing says so
on its own line — `Durable decisions: none` — written when the approach was, per
`skills/write-plan/`. **Adding that line here, to get past the
gate, is the one thing this step is designed to stop.**

Applying is not copying. The contract says what changes; the capability spec has to
read afterwards as though the feature had always been there. Requirements merge into
the existing numbering, `Governs:` gains any new paths, `Proof:` citations come
across, and contradicted sentences are rewritten rather than left beside their
replacement.

Then verify the anchors before the commit lands. `<skill-base>` below is this
skill's base directory, announced when the skill is invoked — the script is a
sibling of this file, wherever the payload was installed:

```
<skill-base>/spec-drift --anchors
```

A citation that survived consolidation pointing at a test that did not is the most
common way this goes wrong.

If a change spans several capabilities, every contract is applied in that same commit.
Half-consolidated is the one state with no honest description.

**Where `libretto` is on PATH and the project holds a consolidated specs directory,
run `libretto wiki` before the landing commit**, so the refreshed index rides the
same commit as the delta that changed it — a generated view left behind by a landing
is drift wearing a marker comment. Where the binary is absent, say the wiki may be
stale and continue; a missing convenience never blocks a landing.

**Where `libretto` is on PATH, run `libretto land` before the landing commit** — after
everything, including the refreshed wiki, is staged; it verifies the staged index. On a
non-zero exit naming a missing part, fix that part and re-run — never commit past it. A
binary too old to know `land` fails with *unknown command* and names no part; that is
the absent case, not a blocked landing. Where the binary is absent, say the landing is
unverified and continue — the same rule as the wiki above.

**`spec-drift` asks it for you**, mechanically, in three directions. It ships beside
this file — `<skill-base>/spec-drift`, whatever directory this skill was installed
into — so it is there in any project, not only in the repository it came from. Run
it from the root of the project being worked on:

```
<skill-base>/spec-drift             staged code whose spec did not move
<skill-base>/spec-drift --anchors    every Proof: citation resolves
<skill-base>/spec-drift --trace      the whole tree, not just what is staged
```

The first reads each spec's `Governs:` globs, so it names the spec and the path
rather than guessing about the repository as a whole. Paths no spec governs are
reported separately and softly — not everything needs a contract.

The second is the reverse direction, and it is the one that catches rot: a
criterion citing a test that was renamed, deleted or never written. It checks the
**test name**, not just the file, because a file-level check passes an invented
name. Run it before the commit that would ship the lie.

**When you run `--trace` or wire `--block` into a hook, read `drift.md` beside this file.**
It carries the `--trace` table and why it never fails, and the opt-in pre-commit hook.

The checklist is updated in the same breath, per `skills/write-tasks/` — by the
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

`caveman-commit` ships with this flow; if the user prefers it, use it for the
message — it produces the same shape, compressed.

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

## Pushing is the user's decision — and it is one question

**Never push unasked.** Not as a convenience, not because it seems obviously wanted,
not because the previous task was pushed.

At the very end — after everything is committed and reported — ask once: **push, and
open the pull request?** One question, one answer. A no ends it and needs no
follow-up.

**Ask it with `AskUserQuestion`, never as a sentence at the end of a report.** This is
the last question in the flow and usually the only one after the plan, so it is the one
that must not be missable — a question in prose is a paragraph the reader can skim past,
and the flow then sits waiting for an answer to something that read as a summary. The
options are yes, no, and push-without-the-request for whoever genuinely wants it; the
recommendation goes first and says what will actually run.

Phase 4 already requires the native prompt for anything the user must settle. This is
the same rule arriving at the phase that used to be written as prose, and it is stated
here because a rule written only where the *asking* is described gets read as belonging
to phase 4 alone.

`AskUserQuestion` is Claude Code's native prompt, and this flow also runs in agents
that have no such tool. **The rule is about the stop, not the widget**: where the
native prompt does not exist, ask the same question in conversation — options listed,
recommendation first — and wait for the answer. That fallback holds everywhere this
document says `AskUserQuestion`, and everywhere another skill in this flow does.

**When the repository has a remote, read `forge.md` beside this file** before asking:
it carries the `/libretto-attacca` exception, the return to the base branch, the bump
question and which forge to use.

**No remote at all means no question.** There is nothing to push to, so the phase ends
at the commit and says so in one line.

## Before the last word

Nothing is reported as recorded without having been seen:

- the commit exists — `git log` was read, not assumed
- the tree is clean, or what remains uncommitted is named and explained
- the spec matches the code, or the divergence is stated
- if pushed, the remote tip matches, and the tree came back to a current base branch —
  or the reason it did not is named

Then one line per task: what it was, its commit, where its evidence is.

## Output

What was committed, on what branch, and whether the spec moved with it.

Then the one question — push and open the request — or, with no remote, the line that
says there is nowhere to push. Then stop.
