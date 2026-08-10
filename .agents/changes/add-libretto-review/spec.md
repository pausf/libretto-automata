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
- the review itself is launched in one fresh subagent with none of the session's
  context; it reads the diff range `<base>...<head>` and returns findings that are
  relayed attributed and unedited — the same seam contract as `review-work`
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
- **One reviewer, not a panel** — inherited; two runs of one model are correlation.

## Task breakdown

- [ ] `skills/review-project/SKILL.md` — URL parsing, forge derivation, clone
      verification, worktree default, fallback with restore, reviewer launch, relay
- [ ] `commands/libretto-review.md` — routes to the skill, nothing else
- [ ] payload spec delta: `libretto-review` joins the outcomes list; capability spec
      `review-project` created on landing

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
