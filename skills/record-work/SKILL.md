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
3. **the plan's durable decisions retired** into that spec's *Prior decisions*
4. **the change folder deleted** — proposal, delta, plan and tasks

All four or none. A delta applied without deleting the change leaves two documents
describing the same capability, and the next reader has no way to tell which one is
current. A change folder deleted without applying the delta loses the work
outright.

**Step 3 is the one that gets skipped, and it is gated.** The plan holds why the change
was built the way it was — the approach, and the alternatives it beat — and step 4
destroys that. It survives in git history, which is to say nobody reads it again.

So before deleting: read the plan's *alternatives* and *prior decisions*, and move
anything that will still constrain work after this lands. Not the whole plan — a
decision about how *this* change was sequenced dies with it correctly. The one worth
keeping is the one somebody would otherwise make again, differently.

`spec-drift --retired`, inside `--anchors`, fails this commit when a `plan.md` is deleted
and no capability spec's *Prior decisions* section moved with it. A plan that genuinely
retires nothing says so on its own line — `Durable decisions: none` — written when the
plan was written, per `skills/write-plan/`. **Adding that line here, to get past the
gate, is the one thing this step is designed to stop.**

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
That is their decision to make, not this flow's to make for them — and once made, it
is one paste. `--block` runs the same checks and turns the warning into exit 1:

```bash
printf '#!/bin/sh\nexec <skill-base>/spec-drift --block\n' > .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

With `<skill-base>` spelled as the real installed path. **This is opt-in and stays
opt-in**: nothing in the flow installs the hook, the default mode still always exits
0, and removing the gate is deleting the hook file.

The checklist is updated in the same breath, per `skills/write-tasks/` — by the
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

`AskUserQuestion` is Claude Code's native prompt, and this flow also runs in agents
that have no such tool. **The rule is about the stop, not the widget**: where the
native prompt does not exist, ask the same question in conversation — options listed,
recommendation first — and wait for the answer. That fallback holds everywhere this
document says `AskUserQuestion`, and everywhere another skill in this flow does.

### The one exception, and it is the whole of `/libretto-attacca`

**Everything above — the question, the native prompt, the waiting — is the attended
flow.** Under `/libretto-attacca` the question is not asked at all, because it was
answered when the command was typed. "Never push unasked" is intact rather than
overridden: the asking happened at the prompt.

The answer covers **this branch and this request and nothing past it** — no merge, no tag,
no release, and no `release:` label the user did not choose. The bump is asked separately
at the very end and is never covered by this consent: typing an answer the user gave is
not the same act as assuming one, and the invocation answered the push, not the version.

And it is paid for in the request's description, which carries three things or the run has
bought silence rather than speed:

- **what the invocation answered** — the two stops and the push — so a reader can tell
  which of the decisions in front of them a person made
- **every question the run assumed past**, each with what changes if it is wrong
- **the bump a person chose**, and that a person chose it — see the bump section below. A
  reviewer has to be able to tell an agreed bump from an assumed one, and that distinction
  is what this whole command is built out of. An unlabeled request says the question went
  unanswered, which is its own honest state and not an omission to be tidied away.

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

### Last of all, under `/libretto-attacca`: the bump

**This is the only question an unattended run asks, and it is the last thing that
happens** — after the report, after the request is confirmed open, after the return to
the base branch. There is nothing downstream of it, and that is what makes it not a
stop: the work is complete and reviewable whether it is answered or not.

It exists because a designed refusal that nobody predicts reads as breakage. The
`release:` label is required before a merge, so an unattended run's request carries a red
check by design — that red *is* the bump question arriving. It arrived as a broken
pipeline once, and cost an alarm and a round trip for a decision the run was standing
next to.

**The run never decides the bump. It asks, and it types the answer.** `AGENTS.md` states
the split: *the reading is yours, the typing is not*. Everything below holds that line.

First, does the repository define the labels at all?

```
gh label list --search 'release:' --limit 100 --json name
```

**Never a bare `gh label list`.** It fetches 30 by default, ascending, so `release:*`
sorts late and falls off the page in any repository past thirty labels — and the answer
comes back "none of the three" for a repository that defines all three. That inverts the
whole check silently, which is the worst way for it to be wrong.

Then match **whole names** against exactly `release:patch`, `release:minor` and
`release:major`. Never a prefix or a substring: `release:patch-hotfix` contains
`release:patch` and is not it. The workflow that reads these labels matches whole names
for the same reason, and a detection looser than the check it feeds is a detection that
finds a label `gh pr edit` will then be refused.

A repository that **defines none of the three** is not asked and is not told why — the
convention does not exist there, so neither does the question. Nothing is created:
inventing `release:minor` in somebody else's repository is deciding that repository's
release convention on their behalf.

Where they do exist, **the bump is asked once** with `AskUserQuestion` — the run's own
reading recommended first, then the others, and room to answer differently.

**`release:major` is present and is never the first option.** Selectable, because a
native question the user answers *is* the asking its standing rule demands. Never
recommended, because recommending it is the announcement that rule forbids — and
announcing a major three times and proceeding is what published `v1.0.0` and `v1.0.1`
from a table read without the paragraph above it. A version number cannot be recalled
once the proxy has cached it.

Then apply exactly one — the workflow refuses two — and **read the request back** off the
forge to confirm the label is on it:

```
gh pr edit <n> --add-label <label>
gh pr view <n> --json labels
```

The same rule the push already carries, for the same reason: a command that printed no
error is not a change the forge accepted.

**Then put it in the description, and this step is not optional.** The description was
written before the question was asked — the request had to exist for there to be anything
to label — so the bump reaches it only by being written back:

```
gh pr edit <n> --body "<the description, plus the bump a person chose>"
```

Without this the third bullet above is a promise the run never keeps: the description
lists what the invocation answered and what it assumed, and a reader has no way to tell
that the bump was neither. The label alone does not say who chose it.

**Unanswered, the run ends exactly as it does today** — unlabeled, and the closing
report's red-check line *is never withdrawn* by the question. That line is written before
the question is asked. A report that promised a red check and then quietly labelled the
request has lied about the state the user will find, which is the same failure as the
silence, wearing the opposite face.

**No default. Ever.** Not patch, not "the safe one". A silently-wrong bump is the failure
that published `v1.0.0` with a politer name on it. Headless runs — `libretto loop` — have
no prompt for `AskUserQuestion` to arrive in, so unanswered is the normal path there and
must be the quiet one.

**Ceiling named:** the question is only as reachable as the terminal it is asked in. A
scheduled or piped run never sees it and lands unlabeled — today's behaviour, which is
why today's behaviour is the fallback rather than an error. The replacement, if that
becomes the common case, is a `gh pr comment` carrying the three commands. Deliberately
not built for a case that has not happened.

**Attended runs do not ask this.** `/libretto-flow` already stops at phase 8 with the user
present and watching; the red check ambushes nobody there. The user's call, 2026-08-13 —
back the day an attended run pays the same round trip.

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
