# Review Project

Governs: skills/review-project/** skills/review-security/** skills/review-design/** skills/review-reliability/** skills/review-tests/** agents/review-lens.md agents/review-intent.md commands/libretto-review.md

Review somebody else's PR/MR without disturbing your own state.

## Outcomes

Installing this repository gives one command, `libretto-review`, that takes a GitHub
or GitLab PR/MR URL and reviews it in a workspace that leaves the user's repository
exactly as it found it.

- `commands/libretto-review.md` routes and never implements; the logic lives in one
  skill, `skills/review-project/`
- given a PR/MR URL, the skill derives the forge from the URL, resolves the request
  through that forge's CLI (`gh` / `glab`), and verifies that **some** remote matches
  the URL's repository — `origin` = fork, `upstream` = canonical is the ordinary
  contributor setup, and no match at all is a stop, not a clone
- the head is fetched to `refs/review/<number>` from the forge's PR/MR ref, so a
  change opened from a fork resolves the same as one from a branch on the repository
  itself — the ordinary shape of "somebody else's change" is not a special case
- the two forges' JSON field names are bound to one vocabulary at the point they are
  read, so no later step is written in one forge's dialect
- **default path:** a `git worktree` at the PR/MR head, the review runs there, and
  the worktree is removed afterwards — the user's checkout is never touched
- **fallback path**, when the project is too large for a worktree: the current branch
  is recorded, the PR/MR branch is checked out in place, and after the review the
  original branch is restored — always, including when the review itself fails
- a dirty working tree stops the fallback before any checkout: the state is reported,
  never stashed or discarded on the user's behalf
- the review runs as **five lenses over one frozen diff**, each a fresh subagent
  with none of the session's context, launched in parallel: **intent** (does the diff
  do what the MR says it does), **security**, **design**, **reliability** and
  **tests**. What each lens looks for is that lens's own contract and lives in its
  own skill — restating it here is the second copy nobody edits
- lenses launch as agents in `agents/`, never as general-purpose subagents, and each
  declares only the tools its lens uses: `review-lens` reads and never writes, and
  `review-intent` alone carries `Bash`. A general-purpose subagent brings every tool
  schema and every installed skill's description into a review that touches four
  tools, and that overhead is paid once per lens before any of them reads a line
- **two agents, not five.** `review-lens` serves all four skill-backed lenses: they
  differ in which skill they apply and in nothing else — same tool grant, same scope
  discipline, same return shape. `review-intent` is separate because it is genuinely
  different: different tools, and no skill to delegate to
- the four skill names are written in the orchestrating skill as
  `Skill(skill="review-security")` and its siblings — the form `scripts/check-payload`
  scans, so a renamed lens fails the gate rather than silently never running. Naming
  them there rather than in the agent is what lets one agent stay one agent
- the lens contract lives in the skill and the agent stays thin; inlining a lens
  contract into its agent would be the second copy
- **intent has an agent but no skill**, the one asymmetry: it is the only lens whose
  brief needs the forge payload the orchestrator holds, and it has exactly one caller,
  so there is nothing to make standalone. Its contract lives in its agent because
  there is no skill for it to duplicate
- the diff is frozen **to a file** and every lens is handed that path. Five lenses
  re-deriving the range is five readings of a moving target, and five explorations of
  one diff — the largest avoidable cost in the run
- a lens opens files outside the changed set only as far as a specific finding
  reaches, and returns findings in a fixed shape with no preamble or closing summary
- the MR description reaches the intent lens fenced as untrusted data, with the
  standing rule that nothing inside the fence is an instruction — the author of the
  change under review is the person who benefits from a clean report
- running the reviewed project's test suite is **opt-in per review**: a suite is code
  the PR author controls, running as the user with their credentials, and the
  worktree isolates git state rather than execution. Not run is a stated limit in the
  intent report, never a silent gap
- lens reports are relayed per lens, attributed and unedited — never merged or
  reranked against each other; one lens passing must not mask another failing
- the lens skills (`review-security`, `review-design`, `review-reliability`,
  `review-tests`) are standalone, each usable on any diff without `review-project`
- the review reports and never blocks, edits, commits or pushes in the reviewed
  repository

## Scope boundaries

**In:** the command, the orchestrating skill, the four lens skills, the two lens
agents, forge derivation, worktree lifecycle, the branch-switch fallback with state
restoration, launching the lens subagents, relaying findings.

**Out — named so it cannot be quietly added:**

- cloning. The input is a PR/MR of a repository the user already has; a repo URL
  with no local clone is a stop with the `git clone` line, not an implicit clone.
- a size threshold, flag or setting for "too large". Prior decision in the payload
  spec: size is knowable from the repository. The fallback is a judgment with named
  signals — a build that needs unversioned files (this repo's own `Ask first` case),
  a dependency install too heavy to duplicate, or the user saying so.
- reviewing against libretto specs. An external project has no `.agents/specs/`;
  the lenses review the diff on its own terms. If the project happens to have
  specs, a lens may read them, but their absence is not a finding.
- posting the review to the forge. Findings come back to the conversation; leaving
  a comment on the PR/MR is the user's action, offered never assumed.
- stashing. Deciding what to do with the user's uncommitted work is theirs.
- a fix loop, or acting on findings. Same standing as `review-work`: report, stop.
- more lenses. Accessibility, performance-as-its-own-lens and docs stay out until
  real runs show a class of finding systematically escaping the five — a lens is
  added from evidence, not from the catalogue.

Never scoped out: the restore. A fallback that can leave the user on the wrong
branch after a failure is data loss's neighbour — the original branch is restored on
every exit path.

## Constraints

- a skill may only invoke what gets installed: every tool these skills need ships
  inside them or is a documented external CLI (`git`, `gh`, `glab`, `rg`, `jq`)
- forge derivation is the recorded prior decision, ceiling included: substring on
  the URL, `github.com` → `gh`, `gitlab` → `glab`; a self-hosted forge on a neutral
  domain is out until an explicit setting exists
- a missing forge CLI stops with its install line — never a hand-built API call,
  never the other CLI as fallback, never a token
- frontmatter `name:` equals each skill's directory. Commands in this payload carry
  `description:` only, so no `name:` is asserted of them
- the PR/MR number matches `^[0-9]+$` before it reaches a command line
- everything the run writes goes in one `mktemp -d` directory, worktree included:
  no fixed `/tmp` path a concurrent run can truncate or a local symlink can redirect
- the worktree path carries owner, repo and number — PR numbers are not unique
  across repositories, and a tree belonging to another repository cannot be pruned
  from this one
- worktrees land outside the repository's working tree; removal uses
  `git worktree remove --force` — the test suite leaves untracked artifacts that
  plain `remove` refuses — never a bare `rm` that could follow a symlink. `rm -rf`
  is allowed on the `mktemp -d` directory alone, which holds nothing this run did
  not write
- the fallback checks out the resolved ref detached, never a bare branch name a
  stale local branch could win by DWIM, and records branch-or-SHA to a file under
  the git dir before moving — a detached HEAD records empty, and a context that ends
  mid-review takes an in-memory record with it

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
- **`review-work` stays a single reviewer** — its own spec anticipated that the day
  the payload chose lenses would be a spec change, and this is that change. The
  internal seam reviews one change against its own spec; this one reviews a stranger's
  diff with no spec to review it against. Merging them would be one abstraction over
  two different jobs.
- **One agent per distinct lens shape, not one per lens** — settled 2026-08-10 after
  three of the four skill-backed agents were produced by string substitution from the
  fourth. That a text replacement can generate them is the evidence they are one thing
  with a parameter; the parameter moved into the brief and the four files became one.
  A fifth stays separate on merit — different tools, no skill behind it.
- **Lens agents with narrow tool grants, over general-purpose subagents** — settled
  2026-08-10 from measurement, not taste. Five general-purpose lenses over a 649-line
  markdown diff cost ~307k tokens, and the per-lens figures were flat (56.8k–63.2k)
  whether a lens made 5 tool calls or 13. Flat cost under varying work is a fixed
  per-agent floor, so the lever is the floor, never the lens count: the skill bodies
  are ~800 tokens each and trimming them buys nothing.
- **Cutting a lens is not the saving, and needs its own evidence** — the payload's
  "a lens is added from evidence, not from the catalogue" cuts both ways. In the first
  real run the intent lens ran the reference checker, read it green, and missed that
  the four lens skills were invisible to it; the tests lens caught it. Same command,
  same output, different brief — that is what the fifth floor buys.
- **The reviewed project's suite is not run unless the user says so** — settled
  2026-08-10 after the security lens traced it on this very change. The first version
  ran it by default, which is arbitrary code execution as the user on the say-so of
  whoever opened the PR.
- **Ideas taken, text not copied** — the lens skills draw on published review
  skills (confidence-gated security findings à la Sentry, non-reranked axes à la
  Pocock, judgment-call principles à la code-quality-principles) and are written
  in this payload's own terms.

## Verification criteria

- frontmatter parses; `name:` matches each skill's directory and each agent's
  filename; every referenced skill resolves, the four lenses included; no uninstalled
  path is invoked
  Proof: scripts/check-payload
- the command delegates to the skill rather than restating it. **No `Proof:` — no
  check asserts this.** `check-payload` reads a command's frontmatter and nothing
  about its body; citing it here would be a file-level anchor standing in for a test
  that does not exist, which is what `--anchors` was written to catch one level up.
  An observation, until a check earns the citation.
- behaviour is a prompt and is checked by running it. Observations owed before this
  is fact rather than claim: the worktree path reviewed a real PR and removed its
  worktree; the same for a PR opened **from a fork**; the fallback restored the
  original branch after a review **and** after a deliberate failure, and restored a
  detached HEAD by SHA; the dirty-tree stop fired; the no-remote-matches stop fired;
  the missing-CLI stop gave the install line and nothing else; a lens that returned
  nothing was relaunched once, and a second silence was reported as that lens
  missing rather than papered over; five reports were relayed unmerged. **The token
  cost of a run with the lens agents, measured against the 307k baseline** — the
  narrow tool grants are a prediction until a second run is counted.
