---
name: build-and-check
description: "Trigger: a plan exists and a task is ready to implement; writing code for a spec; deciding how much to test; choosing a branch or a worktree. Phase 6 of the flow. Writes the code and leaves proportionate proof behind."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Phase 6 of the Libretto flow: implement one task from the plan and leave behind the
proof that it works.

**Delegate the parts that are already solved:**

- `skills/test-driven-development/` — the red-green-refactor discipline, and watching
  the test fail before making it pass
- `skills/using-git-worktrees/` — creating isolation, detecting whether you are
  already in it, and the per-language project setup a fresh checkout needs

Both ship with this repository; see `THIRD-PARTY.md`.

`skills/evidence/` governs everything here. Read it first.

What follows is only what those do not decide: where the work lands, how much to
test, and whether to isolate at all.

## Step 0 — Ensure the branch, before the first write

**Before editing a single file, look at the branch.** If it is the base branch,
create one now.

Often it already exists: phase 1 branches when it writes a proposal, because that file
has to be committed too. Then this step confirms and moves on. **Ensure, not create** —
a step that assumes it is first will make a second branch and split one change across
two.

```
git branch --show-current
git checkout -b <type>/<description> <base>
```

Not before the first commit — before the first **write**. The two look
interchangeable and are not: `git checkout -b` carries uncommitted work along, so
editing on the base branch and branching at commit time appears to work, and keeps
appearing to work until the base has moved or a file it changed is one of yours.
Then the branch cannot be created without a stash, a merge or a loss, and the choice
arrives at the worst moment — after the work, with the diff already in hand.

This happened in the session that produced this rule. Two files were edited on
`main`; the branch was created at phase 8. Nothing broke, which is the reason the
habit survives.

`record-work` checks the same invariant before committing. That is a backstop on a
cheap question, **not** the place the decision belongs: by phase 8 the work already
exists, and a check that can only say "too late" is not a check.

## How much to check

Proportionate. The smallest thing that fails when the logic fails — not a suite per
function:

| The change | The check |
|---|---|
| a branch, a loop, a parser, money, security, a trust boundary | one runnable check, minimum |
| behaviour a user depends on | end-to-end, once |
| a fix for something that broke before | a regression test, always |
| a one-line change with no logic in it | none |

That last row is not laziness. A test for a change with no logic in it tests the
language, fails only when something unrelated moves, and gets deleted by whoever is
annoyed by it first.

**Proportion is about how many, never about how honest.** A test that exists is never
weakened, skipped, deleted or wrapped in a false condition to get a green run. When
one fails: fix the cause, or stop and say why. Those are the only two moves.

If `ponytail` is installed, this is its rule that non-trivial logic leaves one
runnable check behind and trivial one-liners leave none.

## Isolate only when isolation is cheap

A worktree is worth it when several tasks run at once and would collide.

It is not always cheap. Unversioned `.env` files, `node_modules`, vendored
dependencies, generated artifacts and local databases all have to exist again before
the new tree can build. `using-git-worktrees` handles the setup — the decision of
whether to pay for it is here.

The test: **can this project build from a clean checkout without a manual step?** If
yes, a worktree is nearly free. If no, the honest options are a branch in place, or
reproducing the setup deliberately and saying how long it took.

One task at a time on one machine does not need a worktree at all.

## Commit as you go

The branch already exists — step 0. What is left is how many, and when.

A branch per parent task, or per subtask when subtasks are genuinely independent.
Independent branches get chained — `chained-pr` if it is installed — never eight of
them racing at the trunk.

**Commit before verifying, not after.** A gate proves something about the tree it ran
against; if the tree was dirty, it proved nothing about what was recorded.
Work-in-progress commits are cheap and squash later. Verifying uncommitted work and
then committing it is a claim about code nobody tested.

After a rebase or a conflict resolution the tree is new, and everything that was
green before it is unproven again.

## Closing a task

A task is done when the criterion the plan named for it is met, and the run that
proves it was watched. Then the box gets marked — by the orchestrator, per
`skills/write-plan/`.

Two failed gates on one task: stop the task. Not a third attempt. Report what
failed, what was tried, and what is still unknown.

## Output

Per task, one line: what it was, its evidence, and whether it is closed or stopped.

Anything deliberately left out is named here with the condition that would bring it
back, so it reaches the phase 7 report instead of the reader's surprise.
