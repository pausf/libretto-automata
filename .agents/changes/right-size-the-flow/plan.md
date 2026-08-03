# Plan — right-size the flow

Live state. The orchestrator writes this file; sub-agents report and never edit it.

Spec: [spec.md](spec.md) · Proposal: [proposal.md](proposal.md) · Targets: `payload`

Six tasks. Every one traces to a task in the spec's breakdown and names the criterion
that closes it.

---

- [x] **1 · `build-and-check` branches before the first write**
  `skills/build-and-check/SKILL.md`. Move the base-branch check to the top, ahead of
  "How much to check", and name the cost: two files were edited on `main` in the
  session that produced this change, and it only worked because `git checkout -b`
  carries uncommitted work.
  Closes: *phase 6 reports the branch it created before reporting the first edit, and
  no edit lands on the base branch.*
  Waits on: nothing. **Can start now.**

- [x] **2 · `record-work` keeps the check as a backstop**
  `skills/record-work/SKILL.md`. Same invariant, explicitly cross-referenced to phase 6
  as its owner, so a reader does not see two decisions about one thing.
  Closes: *no edit lands on the base branch* — the phase-8 half.
  Waits on: **1**. The wording has to point at what task 1 leaves behind.

- [x] **3 · Push and the PR become one question**
  `skills/record-work/SKILL.md`. Delete "Opening a pull request is a separate question,
  asked separately". Derive the forge from `git remote get-url origin` — `github.com` →
  `gh`, `gitlab` → `glab`. A missing or unauthenticated CLI stops with the install line,
  shaped like `find-work`'s `jira` stop. No remote, no question.
  Closes: *push and open the PR arrive as one question* · *with `gh` absent, phase 8
  stops with the install line* · *with no remote, it does not ask.*
  Waits on: nothing. **Can start now** — independent of 1 and 2, different section of
  the same file, so do not run it concurrently with **2**.

- [x] **4 · The trivial lane**
  Three files, one idea. `skills/write-spec/SKILL.md`: a "no" collapses the phase-7 gate
  as well as the spec. `skills/present-work/SKILL.md`: the *stop* after presenting is
  conditional — what gets said never is. `commands/libretto-flow.md`: route it.
  Closes: *a documentation-only change reaches a commit with exactly one question, and
  that question is the push* · *on a task that needs a spec, all four stops still
  happen.*
  Waits on: **3**. The collapsed lane ends at the merged question, so it has to exist
  first.

- [x] **5 · Invoke to decline**
  `commands/libretto-flow.md`. Phases 2 and 6 are invoked even when the answer is
  "nothing here", and the declining is reported in one line. Announcing a skip and
  gating on it are different things.
  Closes: *phases 2 and 6 each report themselves, including when what they report is
  that there was nothing to do.*
  Waits on: **4**. Same file, and 4 sets the routing this task adds a rule to.

- [x] **6 · Prune the stale link the rename left behind**
  `.claude/skills/read-task-jira` → `skills/read-task-jira`, a directory that no longer
  exists. Left by the `read-task-jira` → `find-work` rename, gitignored, and invisible
  to `check-payload`. Run `libretto status`, then `libretto prune`, and read what they
  say.
  Closes: *`prune` removes only stale links this repository owns, and takes nothing
  else.*
  Waits on: nothing. **Can start now.** Independent of every other task — no file
  overlap.

---

## Evidence

All six closed on `feat/right-size-the-flow`, branched from `main`.

| | Commit | Proof watched |
|---|---|---|
| 1 | `fa8ff42` | `scripts/check-payload` — all checks passed |
| 2 | `f2cc78f` | `scripts/check-payload` — all checks passed |
| 3 | `faea7a8` | `scripts/check-payload` — all checks passed |
| 4 | `efbbdc5` | `scripts/check-payload` — all checks passed |
| 5 | `3d1f31d` | `scripts/check-payload` — all checks passed |
| 6 | no diff | `libretto status --project`: `13 linked · 1 stale` → `libretto prune --project --yes`: `1 removed · 0 refused · 0 failed` → `13 linked`. The thirteen surviving links are the proof it took nothing else |

Task **6** produced no commit and that is correct: `.claude/` is gitignored, so the
prune is a change to the working machine, not to the repository. Its evidence is the
run, which is why the run is quoted here rather than described.

## What can start now

**1**, **3** and **6**. Task **6** touches no file the others touch and is the only one
that exercises the Go binary rather than the prompts.

**2 → after 1. 4 → after 3. 5 → after 4.**

`3` and `2` edit the same file in different sections: sequential, not parallel. One
task at a time on one machine needs no worktree.

## Standing rules for this change

- Every task commits on its own — `record-work`, one commit per task.
- The delta in `spec.md` is amended in the same commit as any code that contradicts it.
- Landing this change applies `spec.md` onto `.agents/specs/payload/spec.md` and
  **deletes this folder**, in one commit. All three or none.
- After the payload moves, re-run `scripts/check-payload`,
  `skills/record-work/spec-drift --self-test` and `--anchors`. Editing the skills can
  break a reference the checks are there to catch.
