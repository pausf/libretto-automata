# add-libretto-review — plan

**Goal:** a `libretto-review` command that takes a GitHub/GitLab PR/MR URL and reviews
it in a worktree, with a branch-switch-and-restore fallback for projects too large to
worktree.

**Architecture:** payload only — one command that routes, one skill that implements,
one reviewer subagent launched fresh. No Go changes.

Spec: `.agents/changes/add-libretto-review/spec.md`. One writer: the orchestrator
marks these boxes; sub-agents report.

## Global constraints

- every criterion closes through the six gates; the payload ones here are
  `scripts/check-payload`, `spec-drift --self-test`, `spec-drift --anchors`
- frontmatter `name:` equals directory (`review-project`) / filename (`libretto-review`)
- no skill references `scripts/` or `docs/`; no token handling; no push from any step
- release `v0.5.0` is tagged at ship time (phase 8), not during these tasks

---

### Task 1: `skills/review-project/SKILL.md`

Depends on: nothing. Start here.
From spec: task breakdown item 1. Closes: check-payload criterion.

**Files:** create `skills/review-project/SKILL.md`

The skill carries, in this order:

- frontmatter: `name: review-project`, trigger description (PR/MR URL to review)
- step 0 — parse the URL: forge by substring (`github.com` → `gh`, `gitlab` →
  `glab`, ceiling named), repo path, PR/MR number; anything else is a stop
- step 1 — preconditions: forge CLI installed (missing → install line, stop; never
  the other CLI, never an API call, never a token); `git remote get-url origin`
  matches the URL's repo (mismatch → stop with the `git clone` line, no clone)
- step 2 — resolve the PR/MR: head branch and base via `gh pr view` / `glab mr view`
  (`--json` / JSON output, exit code read, never piped into `head`)
- step 3 — choose the workspace: worktree by default; fallback is a judgment with
  the spec's named signals (unversioned-files build, heavy install, user says so)
- step 4a — worktree path: `git worktree add` outside the working tree or in an
  ignored dir, review runs there, `git worktree remove` after — never bare `rm`
- step 4b — fallback path: dirty tree → stop and report, never stash; record
  current branch; checkout PR branch; **restore the original branch on every exit
  path**, review failed or not
- step 5 — launch exactly one fresh reviewer subagent: diff range `base...head`,
  no session context, no priming; reviews the diff on its own terms (external
  project, no libretto specs expected); silence ≠ clean review, relaunch once
- step 6 — relay findings attributed and unedited; report never blocks; posting to
  the forge is offered as the user's action, never done
- [x] **1.1** write the file per the outline above
- [x] **1.2** check-payload exit 0 (/tmp/gate1.out, all checks passed)
- [x] **1.3** committed: 7f9893c `feat(payload): add review-project skill`

### Task 2: `commands/libretto-review.md`

Depends on: Task 1 (references the skill by name).
From spec: task breakdown item 2. Closes: command-delegates criterion.

**Files:** create `commands/libretto-review.md`

- frontmatter `description:` only (matches `libretto-status.md` convention)
- takes the URL as input, invokes `Skill(skill="review-project")`, restates nothing
- [x] **2.1** write the file
- [x] **2.2** check-payload exit 0 (/tmp/gate2.out; command and skill both ok)
- [x] **2.3** committed: 65cccaa `feat(payload): add libretto-review command`

### Task 3: spec updates

Depends on: Tasks 1–2 (describes what now exists).
From spec: task breakdown item 3. Closes: `spec-drift --anchors` (166+ citations
resolve).

**Files:** modify `.agents/specs/payload/spec.md` (add `libretto-review` to the
outcomes list, one line); the `review-project` capability spec itself is created by
phase 8 from this change's delta, not here.

- [x] **3.1** add the outcome line to the payload spec
- [x] **3.2** all six gates green, foreground, outputs read (gofmt empty, vet 0,
      tests 0, check-payload passed, self-test passed, 168 anchors resolve)
- [x] **3.3** committed: c4dd023 `docs(spec): record libretto-review in payload outcomes`
