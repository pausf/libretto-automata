---
name: record-work
description: "Trigger: a task is finished and needs committing; writing a commit message; the spec no longer matches the code; deciding whether to push. Phase 8 of the flow. Records the work and reconciles the contract with what was actually built."
license: MIT
metadata:
  author: pausf
  version: "1.2"
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

While the work is in flight, "the spec" means the delta inside
`.agents/changes/<change>/spec.md`. That is where amendments belong until they land.

## Landing a change consolidates it

The last commit of a change does three things together, in one commit:

1. the final code
2. **the delta applied onto the capability spec** it targets, in
   `.agents/specs/<capability>/spec.md`
3. **the change folder deleted** — proposal, delta and plan

All three or none. A delta applied without deleting the change leaves two documents
describing the same capability, and the next reader has no way to tell which one is
current. A change folder deleted without applying the delta loses the work
outright.

Applying is not copying. The delta says what changes; the capability spec has to
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

If a change spans several capabilities, every delta is applied in that same commit.
Half-consolidated is the one state with no honest description.

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

### `--trace` — the map, not a gate

The first two are both scoped to a commit: one to what is staged, one to citations that
already exist. **Neither can see code that was never staged while a spec existed**, and
in any repository whose specs arrived after its code that is most of the repository. So
the contract can be entirely intact by both checks and still govern a third of the tree.

`--trace` reads the whole tracked tree and answers three questions:

| | |
|---|---|
| **orphan code** | a tracked file no `Governs:` glob claims |
| **dead ownership** | a `Governs:` glob matching no tracked file — a package that moved, and a spec still pointing where it was |
| **unproven criterion** | a bullet under *Verification criteria* with no `Proof:` beneath it |

**It always exits 0, and that is deliberate.** A first run on a real project reports
dozens of orphans, because the honest answer is that most code was never specified. A
report that fails the build on the day it is introduced is a report somebody deletes,
and a deleted report finds nothing. Read it, pick the ones that matter, leave the rest.

**What it will not tell you is whether a `Proof:` names a *good* test.** That is a
reading of the test, and a script that guessed would be wrong exactly where the answer
matters. `skills/review-spec/` asks it.

**It warns; it does not block.** Always exits 0. A check that stops a commit in
someone else's project gets removed the same day, and then there is no check at all.
A deliberate no is a valid answer — say it out loud in the report so the next reader
knows the question was asked and not skipped.

Anyone who wants it to be a gate can wire it into their own `pre-commit` hook or CI.
That is their decision to make, not this flow's to make for them.

The plan is updated in the same breath, per `skills/write-plan/` — by the
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

### The one exception, and it is the whole of `/libretto-attacca`

**Everything above — the question, the native prompt, the waiting — is the attended
flow.** Under `/libretto-attacca` the question is not asked at all, because it was
answered when the command was typed. "Never push unasked" is intact rather than
overridden: the asking happened at the prompt.

The answer covers **this branch and this request and nothing past it** — no merge, no tag,
no release, no `release:` label.

And it is paid for in the request's description, which carries two things or the run has
bought silence rather than speed:

- **what the invocation answered** — the two stops and the push — so a reader can tell
  which of the decisions in front of them a person made
- **every question the run assumed past**, each with what changes if it is wrong

This section sits **below** the native-prompt paragraph deliberately. Written above it, the
paragraph that followed read as reinstating the ask unconditionally — the reviewer's
finding on the run that landed the mode.

The yes path is the yes path: what follows below about returning to the base branch happens
here too, because the push and the request were both confirmed.

The two used to be asked separately. That bought a second round trip and no safety:
pushing a branch and then declining to open the request for it is a state almost
nobody wants, and the user who wants exactly that says so in the same breath. Asking
twice about one intent is the ceremony this flow is supposed to spend elsewhere.

When the answer is yes, confirm it landed: **the remote tip matches the local tip.** A
push that printed no error is not a push that was accepted. Then read the created
request back from the forge rather than trusting the output of the command that made
it.

### Then return to the base branch, current

Once the push and the request are both confirmed, and only then:

```
git checkout <base>
git pull --ff-only
```

The working tree ends on the base branch, up to date with the remote. **The feature
branch is left alone** — not deleted, not merged, not rebased. The request is open, not
merged, and a branch deleted at that moment takes the only local copy of unmerged work
with it.

The base is what the repository says it is — `git symbolic-ref refs/remotes/origin/HEAD`,
or the branch the request was opened against. Never hardcoded to `main`.

**Two conditions, and both are stops rather than fixes:**

- **a dirty tree.** `git checkout` carries uncommitted work onto the base branch
  silently, which is the "branch before the first write" failure pointed the other way.
  Report what is uncommitted and stay put.
- **a fast-forward that refuses.** A diverged base is a fact worth seeing. Never
  `git pull` without `--ff-only`, and never resolve it here: a merge commit manufactured
  on somebody's base branch by a phase whose job is bookkeeping is the surprise nobody
  wants.

**Only on the yes path.** No push means the work exists nowhere but this branch, and
moving off it buys nothing.

Why this is here at all: a session starts wherever the last one left the working
directory, so a flow that ends on a merged feature branch hands the next phase 1 a stale
base to read. That is not hypothetical — phase 1 of the run that added this reported a
branch as work in flight when it had already been merged and tagged, because local `main`
was seven commits behind. The reading was wrong, it was reported as fact, and the
correction cost a round trip.

### Which forge, and whether there is one

Derive it. Do not assume, and do not ask what the repository can answer:

```
git remote get-url origin
```

| The URL contains | The tool |
|---|---|
| `github.com` | `gh pr create` |
| `gitlab` | `glab mr create` |

**No remote at all means no question.** There is nothing to push to, so the phase ends
at the commit and says so in one line.

A missing or unauthenticated CLI **stops**, with the install line and nothing else —
the shape `skills/find-work/` already uses for `jira`:

```
brew install gh          # then: gh auth login
brew install glab        # then: glab auth login
```

`gh auth login` and `glab auth login` are interactive. The user runs them
themselves — and **no token ever comes through the conversation.** A credential
pasted to an agent is in the transcript and the logs from that moment on, and
rotating it is the only remedy.

Do not offer a workaround. A hand-built API call with a token found in the
environment is how a stop becomes an exposure, and having the other forge's CLI
installed is not a fallback: `glab` cannot open a request on GitHub.

**Ceiling named:** this is a substring test on one URL. It does not survive a
self-hosted forge on a neutral domain, or Gitea and Forgejo. When that day comes the
answer is an explicit setting read from the repository, not a longer list of guesses.

**The description is written, not omitted.** One or three bullets of what the change
does, the evidence that it works, and what was deliberately left out — the phase 7
report already contains all three. `gitlab-mr-description` or an equivalent, if
installed, may shape it.

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
