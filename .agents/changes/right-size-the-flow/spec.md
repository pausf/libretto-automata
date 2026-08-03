# Right-size the flow — delta

Targets: payload

Amends the payload capability: `skills/**`, `commands/**`.

Two ideas, one sentence: **the flow's ceremony sits in the wrong place — in time, by
charging every change the same four stops, and in space, by putting a phase's decision
in a later phase.**

## Outcomes

**A change that needed no spec costs one stop, not four.**

- `write-spec` answering "no spec needed" also collapses the phase-7 gate. Present and
  record land in the same turn, and the only question asked is the last one.
- The push question survives that collapse. It is the one stop that is not ceremony.

**The branch exists before the first write, not before the first commit.**

- `build-and-check` checks the current branch as its first act. On the base branch it
  branches immediately, before any file is edited.
- `record-work` keeps the same check as a backstop and says so — a second look at a
  cheap invariant, not a duplicated decision.

**Push and the pull request are one question.**

- Phase 8 asks once: push, and open the PR. Not two round trips for one intent.
- A `no` to the pair is still a complete answer needing no follow-up.

**Phase 8 knows which host it is talking to.**

- The forge is derived from `git remote get-url origin`: `github.com` → `gh`,
  `gitlab` → `glab`.
- A missing or unauthenticated CLI stops with the install line and nothing else —
  the shape `find-work` already uses for `jira`.
- No remote means no question. There is nothing to push to.

**A phase that declines is a phase that ran.**

- `write-spec` and `build-and-check` are invoked even when the answer is "nothing to do
  here", and the declining is reported in one line.
- Its "no" costs no wait. Announcing a skip and gating on it are different things, and
  conflating them is what made the previous session look like it forgot the flow.

## Scope boundaries

**In:** `skills/write-spec/SKILL.md`, `skills/build-and-check/SKILL.md`,
`skills/record-work/SKILL.md`, `skills/present-work/SKILL.md`,
`commands/libretto-flow.md`, and the payload capability spec these amend.

**Out:**

- **any new flag, env var or config.** How big the change is, is knowable from the
  change. Asking the user to declare it invents a second source of truth about a fact
  already on disk.
- **the push question.** Never scoped out — it is the user's decision, and it is
  something the user asked for explicitly.
- **`spec-drift` becoming a gate.** It warns and exits 0. A check that blocks a commit
  in somebody else's project gets deleted the same day.
- **sub-agents committing, pushing, or asking the user.** Unchanged in every respect.
- **`present-work`'s content.** What was left out and the condition that brings it back
  still gets said. Only the *stop* after saying it is conditional, never the saying.
- **a `--fast` mode, a profile, or a second flow.** One flow with a proportionate gear,
  not two flows to keep in sync.

## Constraints

- **A skill may only invoke what gets installed.** `install` links `skills/`, `agents/`
  and `commands/` and nothing else, so no instruction added here may name
  `scripts/`, `docs/FLOW.md` or any path that exists only in this repository.
  `scripts/check-payload` enforces it.
- **Every skill stays self-sufficient once installed.** A rule that only works inside
  this repository works for nobody.
- **`ponytail` and `caveman` are called when present, never required.**
- These are prompts, not code. The only mechanical checks available are frontmatter
  validity, reference reachability and citation resolution. **Behaviour is checked by
  running the flow**, which is why the criteria below split into two kinds and say
  which is which.

## Prior decisions

- **The four stops were correct behaviour, not a bug in the agent.** Every one is
  written in the payload. This change edits the instructions, not the discipline that
  followed them.
- **`Huh` is not a dependency** and the confirmation lives in the model. Recorded in
  `docs/STATE.md`; not to be relitigated here.
- **`docs/PLAN.md` is superseded** by these per-capability task breakdowns. It is not
  updated by this change.
- **Work is found, not fetched**, and `/libretto-status` delegates to phase 1 rather
  than scanning on its own. Unchanged.
- **Landed changes are deleted, not archived.** This folder goes away when the delta
  lands on `.agents/specs/payload/spec.md`.
- **This change lands on `feat/right-size-the-flow`, branched from `main`**, not on
  `docs/state-refresh`. Asked and answered 2026-08-03. The two run in parallel and share
  no file: PR #1 is documentation only, this is a process change, and mixing them would
  hand a reviewer something they did not ask for — the thing this change exists to stop.
- **Ceiling named:** the forge detection covers `github.com` and `gitlab` by hostname
  match, which is a substring test on one URL. It will not survive a self-hosted forge
  on a neutral domain, or Gitea and Forgejo. The upgrade path when that day comes is
  an explicit override read from the repository, not a longer list of guesses in a
  prompt.
- **Ceiling named:** "did this change need a spec" stays a judgment made by reading the
  task, with no heuristic behind it — not a diff-size threshold, not a file count.
  `write-spec` Step 0 already states the real test: could two people reasonably
  disagree about what "done" means. A number would be wrong in both directions and
  would be trusted anyway.

## Task breakdown

1. **`build-and-check` branches before the first write.** Move the base-branch check to
   the top of the skill and state the cost of getting it wrong.
2. **`record-work` keeps the check as a backstop**, cross-referencing phase 6 as the
   owner so the two are not read as two decisions.
3. **Push and PR become one question**, with the forge derived from the remote and the
   missing-CLI stop shaped like `find-work`'s.
4. **The trivial lane**: `write-spec`'s "no" collapses the phase-7 gate;
   `present-work` says the stop is conditional; `commands/libretto-flow.md` routes it.
5. **Invoke to decline**: `commands/libretto-flow.md` requires phases 2 and 6 to be
   invoked and their declining reported in one line.
6. **Prune the stale link the rename left behind** — `.claude/skills/read-task-jira`
   points at a directory that no longer exists. Real `stale` state, in the author's own
   project, and the first honest exercise of `libretto prune`.

## Verification criteria

Mechanically checkable:

- frontmatter still parses and every `name:` matches its directory
  Proof: scripts/check-payload
- no edited skill invokes a path that does not get installed
  Proof: scripts/check-payload
- every referenced skill and path still exists after the edits
  Proof: scripts/check-payload
- every `Proof:` citation in every spec and delta resolves, file and test name
  Proof: skills/record-work/spec-drift --anchors
- the drift tooling itself still behaves after the payload moves
  Proof: skills/record-work/spec-drift --self-test
- `prune` removes only stale links this repository owns, and takes nothing else
  Proof: internal/link/apply_test.go TestPruneRemovesOnlyOurOwnStaleLinks
- the prune plan is built from stale entries alone
  Proof: internal/link/plan_test.go TestPrunePlanTakesOnlyStale

Checked only by running the flow — **stated as observations, not as citations, because
a criterion citing a test that cannot exist is the fabrication the anchor exists to
prevent**:

- Observed: a documentation-only change reaches a commit with exactly one question
  asked, and that question is the push.
- Observed: on a task that needs a spec, all four stops still happen. The gear is
  proportionate, not removed.
- Observed: phase 6 reports the branch it created before reporting the first edit, and
  no edit lands on the base branch.
- Observed: `push` and `open the PR` arrive as one question, and a `no` ends it.
- Observed: with `gh` absent, phase 8 stops with the install line and does not offer a
  workaround. With no remote, it does not ask at all.
- Observed: phases 2 and 6 each report themselves, including when what they report is
  that there was nothing for them to do.
