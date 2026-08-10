# Delta: add-libretto-review

Targets: review-project — a new capability; no existing spec covers reviewing an
external project's PR/MR.

## Outcomes

Installing this repository gives one more command, `libretto-review`, that takes a
GitHub or GitLab PR/MR URL and reviews it in a workspace that leaves the user's own
state exactly as it found it.

- `commands/libretto-review.md` routes and never implements; the logic lives in one
  skill, `skills/review-project/`
- given a PR/MR URL, the skill derives the forge from the URL, resolves the request
  through that forge's CLI (`gh` / `glab`), and verifies the current repository's
  `origin` matches the URL's repository — a mismatch is a stop, not a clone
- **default path:** a `git worktree` at the PR/MR head, the review runs there, and
  the worktree is removed afterwards — the user's checkout is never touched
- **fallback path**, when the project is too large for a worktree: the current branch
  is recorded, the PR/MR branch is checked out in place, and after the review the
  original branch is restored — always, including when the review itself fails
- a dirty working tree stops the fallback before any checkout: the state is reported,
  never stashed or discarded on the user's behalf
- the review runs as **five lenses over one frozen diff**, each a fresh subagent
  with none of the session's context, launched in parallel:
  - **intent** — does the diff do what the MR says it does: missing requirements,
    scope creep, implemented-but-wrong; every finding quotes the MR's stated
    intent. If the project has a test suite, this lens runs it against the head
    and reports what it observed, never what the MR claims
  - **security** — the `review-security` skill: confirmed findings with
    attacker-controlled input traced and severity classified; a dangerous pattern
    with one unreadable hop is reported as a marked question, never dressed up as
    confirmed; theoretical issues are dropped
  - **design** — the `review-design` skill: YAGNI, KISS, SOLID and code smells as
    labelled judgment calls; the reviewed project's own conventions override the
    baseline
  - **reliability** — the `review-reliability` skill: what breaks at runtime —
    logic errors, edge cases, races, unbounded work, leaked resources, error paths
    that lose data; every finding names the input or state that triggers it, and
    the report closes by naming what could not be verified
  - **tests** — the `review-tests` skill: does the change carry its proof, and did
    any existing proof get quietly weakened — skips added, assertions broadened,
    invariant tests deleted; test-tampering findings are always severe
- lens reports are relayed per lens, attributed and unedited — never merged or
  reranked against each other; one lens passing must not mask another failing
- the lens skills (`review-security`, `review-design`, `review-reliability`,
  `review-tests`) are standalone, each usable on any diff without `review-project`
- the review reports and never blocks, edits, commits or pushes in the reviewed
  repository

## Scope boundaries

**In:** the command, the skill, forge derivation, worktree lifecycle, the
branch-switch fallback with state restoration, launching one reviewer subagent,
relaying findings.

**Out — named so it cannot be quietly added:**

- cloning. The input is a PR/MR of a repository the user already has; a repo URL
  with no local clone is a stop with the `git clone` line, not an implicit clone.
- a size threshold, flag or setting for "too large". Prior decision in the payload
  spec: size is knowable from the repository. The fallback is a judgment with named
  signals — a build that needs unversioned files (this repo's own `Ask first` case),
  a dependency install too heavy to duplicate, or the user saying so.
- reviewing against libretto specs. An external project has no `.agents/specs/`;
  the reviewer reviews the diff on its own terms. If the project happens to have
  specs, the reviewer may read them, but their absence is not a finding.
- posting the review to the forge. Findings come back to the conversation; leaving
  a comment on the PR/MR is the user's action, offered never assumed.
- stashing. Deciding what to do with the user's uncommitted work is theirs.
- a second reviewer, a fix loop, or acting on findings. Same standing as
  `review-work`: report, stop.

Never scoped out: the restore. A fallback that can leave the user on the wrong
branch after a failure is data loss's neighbour — the original branch is restored on
every exit path.

## Constraints

- a skill may only invoke what gets installed: every tool the skill needs ships in
  `skills/review-project/` or is a documented external CLI (`git`, `gh`, `glab`,
  `rg`, `jq`)
- forge derivation is the recorded prior decision, ceiling included: substring on
  the URL, `github.com` → `gh`, `gitlab` → `glab`; a self-hosted forge on a neutral
  domain is out until an explicit setting exists
- a missing forge CLI stops with its install line — never a hand-built API call,
  never the other CLI as fallback, never a token
- frontmatter `name:` equals the directory (`review-project`) and the command
  filename (`libretto-review`)
- worktrees land outside the repository's working tree or in an ignored directory;
  removal uses `git worktree remove`, never a bare `rm` that could follow a symlink
- release: this is a new capability — a minor bump, `v0.5.0`, tagged when it ships

## Prior decisions

- **Input is a PR/MR URL of an already-cloned repository** — settled by the user
  2026-08-10, over clone-always and both. The branch-switch fallback only makes
  sense against a clone the user owns.
- **Payload, not CLI** — settled by the user 2026-08-10. A command + skill in
  `~/.claude`, not a Go subcommand.
- **"Too large" stays a judgment** — inherited from the payload spec's "no flag for
  how big a change is". Signals named under scope boundaries.
- **The reviewer reports and never blocks** — inherited from the review seam.
- **Lenses, not one reviewer and not N copies** — settled by the user 2026-08-10,
  after reviewing published review skills; extended the same day from three to five
  (reliability and tests). The payload's "one reviewer, not a panel" stands for the
  internal seam and anticipated exactly this: lenses are a spec change, and this is
  that change. Diversity comes from distinct briefs (intent, security, design,
  reliability, tests), never from re-running one brief twice — two runs of the same
  model over the same diff with the same brief are correlation.
- **Ideas taken, text not copied** — the lens skills draw on published review
  skills (confidence-gated security findings à la Sentry, non-reranked axes à la
  Pocock, judgment-call principles à la code-quality-principles) and are written
  in this payload's own terms.

## Task breakdown

- [x] `skills/review-project/SKILL.md` — URL parsing, forge derivation, clone
      verification, worktree default, fallback with restore, reviewer launch, relay
- [x] `commands/libretto-review.md` — routes to the skill, nothing else
- [x] payload spec delta: `libretto-review` joins the outcomes list; capability spec
      `review-project` created on landing
- [x] `skills/review-security/SKILL.md` — the security lens, standalone
- [ ] `skills/review-reliability/SKILL.md` — the runtime-bugs lens, standalone
- [ ] `skills/review-tests/SKILL.md` — the proof lens, standalone
- [ ] `review-project` step 5 wired for five lenses
- [x] `skills/review-design/SKILL.md` — the design lens, standalone
- [x] `review-project` step 5–6 rewritten: freeze the diff, three lenses in
      parallel, relay per lens without reranking

## Verification criteria

- frontmatter parses; `name:` matches directory and filename; every reference
  resolves; no uninstalled path is invoked
  Proof: scripts/check-payload
- the command delegates to the skill rather than restating it
  Proof: scripts/check-payload
- behaviour is a prompt and is checked by running it. Observations owed before this
  is fact rather than claim: the worktree path reviewed a real PR and removed its
  worktree; the fallback restored the original branch after a review **and** after a
  deliberate failure; the dirty-tree stop fired; the origin-mismatch stop fired; the
  missing-CLI stop gave the install line and nothing else.
