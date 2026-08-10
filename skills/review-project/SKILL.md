---
name: review-project
description: "Trigger: a GitHub or GitLab PR/MR URL to review in a repository already cloned on this machine. Prepares an isolated workspace — a worktree by default, a branch switch with guaranteed restore when the project is too large — runs three fresh review lenses (intent, security, design) over one frozen diff, and relays their findings per lens. Reports, never blocks."
license: MIT
metadata:
  author: pausf
  version: "1.0"
---

## What this does

Review somebody else's change without disturbing your own: take a PR/MR URL, put its
head in a workspace that costs the user nothing, hand one frozen diff to three
review lenses that carry none of this session — **intent** (does it do what the MR
says?), **security** (`review-security`) and **design** (`review-design`) — and
bring the findings back, per lens.

`evidence` governs here: nothing reported that was not observed, every command whose
exit code decides the next step run in the foreground and read, findings relayed as
returned.

The user's state is the thing this skill must not damage. Every path below ends with
the repository exactly as it was found — that is the contract, and it holds on
failure too.

## Step 0 — Parse the URL

Derive the forge from the URL: `github.com` → `gh`, a host containing `gitlab` →
`glab`. Extract the repository path and the PR/MR number.

Anything that is not a PR/MR URL — a bare repo URL, a branch link, an issue — is a
stop: say what was expected, do not guess a number.

Ceiling, named: a substring test does not survive a self-hosted forge on a neutral
domain, or Gitea and Forgejo. The replacement that day is an explicit setting read
from the repository, not a longer list of guesses here.

## Step 1 — Preconditions, each one a stop

**The forge CLI exists.**

```
which gh    # or glab
```

Missing → report it with its install line and stop:

```
brew install gh      # github.com
brew install glab    # gitlab
```

Never the other forge's CLI as a fallback, never a hand-built API call, never a
token. A skill handed a token says it is now exposed and continues without it.

**The clone is already here.** This skill reviews a repository the user has;
cloning is not its job.

```
git remote get-url origin
```

The remote must name the same repository as the URL — compare host and path,
ignoring `.git` and protocol. A mismatch is a stop that states both values and the
line the user would run themselves:

```
git clone <url>
```

**Exit codes are read.** Never pipe any command here into `head` or anything else
when its status decides the next step — the pipeline reports the last command's
status, so a failure reads as success.

## Step 2 — Resolve the PR/MR

Head branch and base, through the forge CLI, JSON out:

```
gh pr view <number> --json headRefName,baseRefName,title > /tmp/review-pr.json; echo $?
glab mr view <number> --output json > /tmp/review-mr.json; echo $?
```

Non-zero exit → read stderr, report it verbatim, stop. Then fetch:

```
git fetch origin
```

## Step 3 — Choose the workspace

**Worktree by default.** The fallback exists for the project too large to stand up
twice, and that is a judgment made from the repository, not a flag:

- the build needs unversioned files — `.env`, local databases, generated artifacts
  that exist only in this checkout
- the dependency install is too heavy to duplicate for one review
- the user said so

Say which path was chosen and why, in one line. When none of the signals apply, the
worktree is nearly free — take it.

## Step 4a — The worktree path

```
git worktree add /tmp/review-<number> origin/<headRefName>
```

Outside the project's working tree, or in a directory git already ignores — never
an unignored directory inside the repository.

Review runs there. Afterwards, always:

```
git worktree remove /tmp/review-<number>
```

`git worktree remove`, never a bare `rm -rf` — a removal that follows a symlink
destroys what the link pointed at, and the tool that made the tree is the tool that
knows how to unmake it.

## Step 4b — The fallback path

**A dirty tree stops the fallback before any checkout.**

```
git status --porcelain
```

Output means uncommitted work, and deciding what happens to it is the user's —
report the state and stop. Never stash on the user's behalf, never discard.

Clean tree:

```
git branch --show-current        # record it — this is the way back
git checkout <headRefName>
```

**The recorded branch is restored on every exit path.** Review finished, review
failed, reviewer died twice, anything — the last act of this path is:

```
git checkout <recorded-branch>
```

A fallback that can leave the user stranded on a branch they never chose has lost
the thing this skill exists to protect.

## Step 5 — Freeze the diff, launch the three lenses

**Freeze once, before any lens runs**: the diff range
`origin/<baseRefName>...origin/<headRefName>` and the commit list. Every lens
reviews those exact bytes; three lenses over three readings of a moving target are
three reviews of different things.

Three fresh subagents, in parallel, each carrying one brief and none of this
conversation. Telling a reviewer what to expect is priming the witness, so each
gets the workspace path, the frozen range, and its brief — nothing else:

- **intent** — the PR/MR title and description, and one question: does the diff do
  what it says? Report requirements stated but missing or partial, behaviour
  present but never asked for, and requirements implemented but implemented wrong —
  each finding quoting the stated intent it fails.
- **security** — apply the `review-security` skill to the frozen range.
- **design** — apply the `review-design` skill to the frozen range.

Each lens reviews the project on its own terms — external repository, so the
absence of a spec directory is not a finding; conventions the project writes down
override any baseline. If the project has a test suite, the intent lens runs it
against the head and reports what it observed, never what the MR claims.

No lens edits, commits or pushes in the reviewed repository.

**Silence is not a clean review.** A lens that died or returned nothing did not
run — say so and relaunch that lens, once. Twice is a stop, reported as that lens
missing, never papered over by the other two.

## Step 6 — Relay per lens, then restore

The three reports stay three, attributed and unedited — never merged, never
reranked against each other, never filtered by "that one is minor". A change can
pass two lenses and fail the third, and one lens passing must not mask another
failing. A lens's explicit "nothing found" is a statement, not an absence.

The review reports; it never blocks anything. Posting it to the forge is the user's
action — offer the `gh pr comment` / `glab mr note` line, never run it unasked.

Then the workspace goes: worktree removed, or the recorded branch restored. The
report is not finished while the user's repository is not back to where it was.
