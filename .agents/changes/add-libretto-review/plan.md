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

### Task 4: the review lenses (scope added by the user 2026-08-10)

Depends on: Tasks 1–3 (rewrites review-project's step 5–6).
From spec: task breakdown items 4–6. Closes: check-payload criterion.

**Files:** create `skills/review-security/SKILL.md`, `skills/review-design/SKILL.md`;
modify `skills/review-project/SKILL.md` steps 5–6.

- security lens: high-confidence findings only, attacker-controlled input traced,
  framework mitigations checked before flagging, severity classified, explicit
  do-not-flag list; standalone on any diff
- design lens: YAGNI, KISS, SOLID, smells — every finding a labelled judgment call,
  the reviewed project's conventions override the baseline; standalone on any diff
- review-project: freeze the diff once, launch three fresh lens subagents in
  parallel (intent brief inline; security and design by skill name), relay per lens,
  never merge or rerank across lenses
- [x] **4.1** write `skills/review-security/SKILL.md`
- [x] **4.2** write `skills/review-design/SKILL.md`
- [x] **4.3** rewrite review-project steps 5–6 for the three lenses
- [x] **4.4** all six gates green, outputs read (gofmt empty, vet 0, tests 0,
      check-payload ok incl. both new skills, self-test ok, 168 anchors resolve)
- [x] **4.5** committed: a625b87 `feat(payload): review lenses — intent, security, design`

### Task 5: two more lenses — reliability and tests (scope added by the user 2026-08-10)

Depends on: Task 4 (extends the lens set).
From spec: task breakdown items for `review-reliability`, `review-tests`, and the
five-lens wiring. Closes: check-payload criterion.

**Files:** create `skills/review-reliability/SKILL.md`,
`skills/review-tests/SKILL.md`; modify `skills/review-project/SKILL.md` step 5.

- reliability lens: trigger-traced runtime bugs (races, unbounded work, leaked
  resources, data-losing error paths), severity by impact, blind spots named
- tests lens: change carries its proof; invariant erosion (skips added, assertions
  loosened, tests deleted) always severe until the MR's intent explains it
- [ ] **5.1** write `skills/review-reliability/SKILL.md`
- [ ] **5.2** write `skills/review-tests/SKILL.md`
- [ ] **5.3** wire five lenses in review-project step 5
- [ ] **5.4** run all six gates foreground, exit codes read
- [ ] **5.5** commit
