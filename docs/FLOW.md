# The flow

The payload: what Libretto Automata installs into Claude Code (and, skills and commands
only, into Codex and OpenCode). Eight phases, two fresh readers, three standing rules, as
skills, agents and commands.

The CLI that delivers it is specified in
[`.agents/specs/cli/spec.md`](../.agents/specs/cli/spec.md); [SPEC.md](SPEC.md) lists every
capability. The two are independent — a phase is a markdown file with frontmatter.

## The whole flow on one page

```
 1  find the work          find-work        in flight · tracker · what you said
 2–3 write the contract    write-spec       proposal.md: the ask, decision log, six pillars
 5  choose the approach    write-plan       ## Approach, alternatives, Durable decisions:
 5→6 cut the checklist     write-tasks      ## Checklist, last — boxes with Closes:
    a cold read            review-spec      fresh spec-reviewer, then THE ONE STOP
 6  build and check        build-and-check  branch first, proportionate proof, commit per box
    a second cold read     review-work      fresh work-reviewer re-runs every proof; fixes, no ask
 7  present                present-work     in the contract's terms, with what was left out
 8  record and land        record-work      contract onto the spec, folder deleted; push? (the door)
```

The trivial lane skips 2 through 7: the change, its check, a commit, and the push question.
`/libretto-attacca` runs the same rows and answers both stops in advance.

## 1 · Find the work

Phase 1 asks three sources, in this order, and loads exactly one piece of work:

| Source | When |
|---|---|
| **a change already in flight** | some `.agents/changes/<c>/` has an open box in its checklist, or an unmerged branch carries work |
| **a tracker key or URL** | one was given |
| **what you said** | anything else — the request *is* the input |

**Phase 1 is not finished until its artifact exists.** A request made in conversation leaves
`.agents/changes/<name>/proposal.md` on disk — `Tracker: none`, the ask verbatim, a verb-led
name like `add-relative-discounts` — committed on a fresh branch, before the reading is even
stated. A phase that only described the work had not run, and that happened.

Three more things it reads, because each one was once missed:

- **a branch is work in flight too.** The trivial lane leaves no checklist behind, so
  unmerged branches and their open requests are listed beside the open boxes.
- **finished and not landed.** A change whose boxes are all closed and whose folder is still
  there is reported by name: the landing did not finish. It is not work to continue.
- **whether it is a bug.** A report that something built behaved wrong changes phase 2's
  whole shape — the criterion that would have caught it comes first, with a red `Proof:` — so
  the reading names it, in the reporter's own words.

### The tracker, when one is given

Everything goes through the `jira` CLI — never MCP, never the REST API — and the details
live in `find-work/tracker.md`, which the skill reads only when a key or URL was given. A
conversation request never pays for them. The preflight, because this ships to people who
have configured nothing:

```
which jira                                          → missing: brew install ankitpokhrel/jira-cli/jira-cli
test -f "${JIRA_CONFIG_FILE:-~/.config/.jira/.config.yml}"   → missing: the user runs `! jira init`
jira project list                                   → a real authenticated call; `jira me` proves nothing
```

Then the task is read with `--raw` and `jq` — `--plain` emits colour codes and omits the
description, which is the whole point of reading the task.

**The API token never enters the conversation.** `jira init` is interactive and the user
runs it in their own shell. A token pasted to an agent lands in the transcript, the logs and
the history. The domain may be asked for. The secret may not.

`libretto doctor` runs the same preflight outside the flow: `jira` and its config, `gh` or
`glab` for the forge, `rg`, `jq`, and the two vendored companions.

### Why the queue is two commands and not one

Ideas arrive faster than they get built. `/libretto-queue` captures them one after
another — a proposal with a `Queued:` date and your words verbatim, no branch and no spec
— and `/libretto-next` picks one up later, oldest first, and takes it into the flow.

Two commands, because `/libretto-flow EUCAR-1234` does *that* ticket, always. A flow that
quietly substitutes different work for what you handed it is the surprise nobody wants,
and a flow that reads the queue on its own would have to decide when to override you.

**Queued is not in flight.** A captured idea never blocks the first source, because
abandoning an idea costs nothing and a queue that is expensive to add to is a queue nobody
uses.

What counts as queued is written once, in `find-work/queue.md`, and `libretto queue` reads
the same rule from the shell: a proposal carrying `Queued:`, oldest first, ties by name —
minus any whose `## Checklist` already holds a box, and minus any an unmerged branch already
took. That last filter is a bug fix: a status report once listed as queued a change that was
open as a pull request, and the next command would have rebuilt it.

`/libretto-next` asks which idea, branches, removes the `Queued:` line, commits that, and
enters at phase 2. Neither command ever creates `.agents/changes/` in a project that has none;
`libretto init` is the bootstrap, run on purpose.

### Why the source order is the point, and not merely an order

Phase 1 asks three sources — a change already in flight, a tracker key or URL, and what you
said — and it asks them in that sequence rather than picking whichever is easiest to read.

**Home first, because starting something new while a change sits half-finished is how the
half-finished thing gets abandoned.** The cost is not the wasted work. It is a
`.agents/changes/` directory nobody trusts any more: once two abandoned changes are sitting
there, the next session cannot tell a live plan from a dead one, and the honest response is to
stop reading the directory at all. That is the point where the flow's memory stops working, and
it is reached by convenience rather than by decision.

Which is also why what you said is a **source and not a fallback**. It is how most work
arrives, and every change in this repository so far arrived that way.

## 2 · Write the spec

The first question is whether there is a contract at all, and the answer picks the lane:

| Lane | When | What it costs |
|---|---|---|
| **trivial** | nobody could reasonably disagree about what "done" means — a typo, a rename, a one-line fix | the change, its check, a commit, and phase 8's question. Nothing else |
| **contract** | somebody could | one change document, one fresh reader, **one stop before code**, then phase 8's question |

### Why two lanes and not a ladder

Until 2026-10-21 the flow had no middle gear. The trivial lane was free, and the moment a
change was not trivial it paid the full price: a proposal, a spec delta, `plan.md`,
`tasks.md` and a decision log — five artifacts — drafted by three fresh writers
(`spec-writer`, `plan-writer`, `task-cutter`), with two stops before code. A change that
needed four sentences of contract paid the same as one that needed forty.

**The cost was not the writing, it was the hand-offs.** Each fresh writer started from
files instead of the conversation, so each one re-read what the session already knew, and
each stop asked the user to agree to part of a decision before seeing the rest of it. The
contract lane now costs what its contract costs: the document grows only as long as the
change needs, and the fresh eye moved from drafting the document to reading it.

The proportional shape owes its idea to Gentleman Programming's gentle-ai and its Organic
Driven Development; the vocabulary and the living specs are this flow's own.

### The contract lane

An easy task is one spec and one session.

A large one is split — not for parallelism, but so each piece holds a single
coherent idea and fits in a session without the context going stale. Splitting for
speed and splitting for coherence give different cuts. Take the coherent one.

Six things explicit in every spec. Whatever is left vague is a gap an agent fills
by statistical association:

| | |
|---|---|
| **Outcomes** | the verifiable end state, in operational terms |
| **Scope boundaries** | what is in, and the non-goals |
| **Constraints** | topology, latency, versions, conventions |
| **Prior decisions** | what must not be relitigated |
| **Task breakdown** | atomic units, independently assignable |
| **Verification criteria** | how each unit is proven — in EARS (*When … the system shall …*), each one naming the test that proves it |

**Criteria carry `shall`, and that is checked.** `spec-drift --ears`, inside `--anchors`,
fails a change's criterion written in prose: one nobody could fail is treated as satisfied by
default. Each cites its test by file and name — `Proof: internal/link/apply_test.go
TestApplyIsIdempotent` — and `--anchors` fails a name that does not exist.

**Phase 2 opens the spec that governs the change, never the whole specification**: the
change's own `Targets:`, else the `Governs:` line that claims the path, else the index. Each
capability spec declares what it owns — `Governs: cmd/** !cmd/tool/gates_test.go`, where a
leading `!` gives a path back.

**A bug gets its criterion before its fix**: name what would have caught it, write its
`Proof:` as a test, watch it fail, then fix. Reversed, the criterion describes the fix that
exists and has never been seen to fail.

### Where it goes

```
.agents/
├── specs/                        the living, consolidated specification
│   └── bundle-products/
│       └── spec.md               Governs: src/bundles/**
└── changes/                      changes in flight
    ├── add-relative-discounts/   one capability
    │   └── proposal.md           Targets: bundle-products
    └── add-bundle-gift-cards/    two capabilities
        ├── proposal.md           Targets: bundle-products
        └── spec-gift-cards.md    Targets: gift-cards — that capability's six pillars
```

**One change document.** `proposal.md` is the file phase 1 or the queue captured, and it
grows in place rather than gaining siblings. Top to bottom:

| Part | Where | Written at |
|---|---|---|
| `Tracker:` · `Targets:` · `Durable decisions:` | the top lines | 1, 2, 5 |
| the ask, verbatim, and the reading | the capture's own headings | 1, or the queue |
| every answer, verbatim, `(assumed)` kept | `## Decision log` | 2 and 5 |
| the contract | the six pillars as `##` headings | 2–3 |
| the how | `## Approach` | 5 |
| the boxes, and nothing else | `## Checklist`, always last | 5→6 |

A change that amends a second capability adds one `spec-<capability>.md` beside it, with
its own `Targets:` and its own six pillars. The decision log, the approach and the
checklist still exist once, in `proposal.md` — a second file about the same change is a
second source of truth.

**The old layout is still read.** A change begun before 2026-10-21 — a `spec.md` delta,
`plan.md`, `tasks.md` beside its proposal — finishes in the layout it started in. The
`brief.md` and `decisions.md` of the layout before that are neither written nor read any
more; the change document is both. Every
reader picks a change's checklist the same way: `proposal.md` when it holds a box under
`## Checklist`, otherwise `tasks.md`, otherwise `plan.md`. Nothing in flight is migrated,
because migrating a half-built change is how its reasoning gets lost in transit.

**One spec per capability, never per ticket.** A capability spec accumulates and
stays true; a spec named after a ticket is dead the day the ticket closes, and a
directory of dead specs is worse than none because people still read them.

Amendments live in the change while the work is in flight, and phase 8 applies them
onto the capability spec and deletes the change. Exactly one current specification,
always.

## 3 · One contract per capability

When the task touches several capabilities, each gets its own contract: the primary one in
`proposal.md`, every other one in a `spec-<capability>.md` beside it. The cut follows
coherence, not the tracker's subtasks — subtasks cut along who does the work, and a spec cut
that way describes half an idea.

**The session writes them inline**, the same session that held the interview. It already
holds the answers; a fresh writer would have to rebuild them from files and would get the
ones nobody wrote down wrong.

**Context is the only reason to fan out.** When the extra capabilities' research will not
fit in this session — ten subtasks researched in one thread — one `spec-writer` per extra
`spec-<capability>.md` gathers what that file needs and writes it, and nothing else.
Research fans out; the files stay separate so sessions never collide.

What keeps a fanned-out set coherent is `proposal.md`, **written first**: its decision log,
the primary contract and the vocabulary every file will use. Each writer is given it by
path and reports back what it got wrong. Those reports are where the real findings are.

All of that lives in `write-spec/fan-out.md`, read only when a second capability is
involved. The common case — one capability — never reads a line of it.

## 4 · Ask

When the task is complex the agent asks: something it does not understand,
something it found that contradicts the spec, or anything where guessing wrong
breaks working code.

Every question offers three things — **the option the agent recommends**, real
alternatives, and room to write a different answer.

The answer is recorded in the spec next to the decision it produced. An
undocumented answer gets asked again next session.

`AskUserQuestion` is native. Do not build a prompt system.

### As many as a wrong guess would cost, at phase 2, and zero is allowed

The user's call, 2026-08-12: *"me gustaría que hagas más preguntas del estilo Claude, para
que el plan se cree entre los 2 y no solo tú"*. A phase that hands over a finished spec has
decided alone everything it did not ask about, and the decisions it made quietly are
exactly the ones nobody reviews.

So phase 2 asks before it writes the file — **as an interview, one question per call,
each answer able to redirect the next.** It was one batched call until 2026-08-19: the
batch was cheap but it interviewed nobody — all its questions came from the same draft,
and the user met them as a form. What keeps the sequence from becoming the serial
ceremony this flow argues against is the soft bound around five and a standing "no more
questions" option, and the answers still have to be *in* the contract rather than bolted
on afterwards — they land verbatim under the change document's `## Decision log` as they
arrive, and in its *Prior decisions* beside what they settled.

The count was capped at three until 2026-08-14, when the user lifted it: better asked out
of caution than swallowed out of fear. What replaced the cap is judgment with both edges
named — every question one a wrong guess would make expensive, and never a form-length
interrogation of things the code already answers. The bias, when in doubt, is to ask.

**Phase 5 asks too, since 2026-08-19, and it asks exactly one thing: the approach.** Two
or three real ones with named tradeoffs, the recommendation first with its reason, the
user choosing — chosen and rejected logged with why, which is what makes the approach's
alternatives table reconstructible outside the session that argued it. "Phase 2 alone,
not phase 5" was the rule from 2026-08-12, declined on the no-third-stop argument; the
user retired it, and the argument survives the retirement — the fork rides inside phase
5, and asking N questions before the first line of code is still how a flow with one stop
before code grows back into one with nine when the questions stop being cheap.

**Zero is legitimate, and it is said in one line.** The alternative — always three, stretched
if necessary — was offered and declined for the reason `AGENTS.md` already gives: do not ask
what the code can tell you. A quota manufactures the rubber-stamp question this flow removed
from three other phases.

Note what does *not* change: the stop count. The questions ride the phases that already
have them.

## 5 · The plan

**The how.** Until 2026-08-17 this phase produced a checklist and called it a plan,
which is what a Lead reading the output reported: the specs read like plans and the
plans read like task lists. The diagnosis was one line long — the spec's task breakdown
was transcribed into checkboxes, so the *technical reasoning* never got written down
anywhere at all.

`## Approach`, in the change document after the contract, holds seven things, and the fourth
is the one this flow had no home for:

| | |
|---|---|
| summary | the shape of the approach, in a paragraph |
| technical context | what it is built against, and the blast radius |
| the approach | the mechanism, in enough detail to be built from |
| **the alternatives it beat** | **one row each, and why it lost** |
| risks | what breaks, and the mechanism that catches it |
| validation and rollback | which gates carry it, which test would pass for the wrong reason, how it comes back out |
| complexity deliberately kept | where it is not the simplest thing that works, and why the simpler thing was wrong |

Before writing it, phase 5 asks its one question — the fork: two or three real approaches,
the recommended one first with its reason. Under `/libretto-attacca` the recommendation is
taken and logged `(assumed)`.

A diff shows what was built. Nothing in a repository shows what was **not** built and
why — so if it is not here it is gone, and the decision gets made again by somebody with
less context in the middle of phase 6.

**The session writes it inline**, right after the fork: the user chooses the approach, the
choice and the losers land in the decision log, and the session that heard the argument
writes the section. It also writes the `Durable decisions:` line at the top of the
document — which *Prior decisions* outlive the change, or `none` — because the landing
commit deletes the change and `--anchors` fails it when reasoning leaves with the folder
undeclared.

From 2026-08-19 until 2026-10-21 a fresh `plan-writer` drafted this from the files alone,
so it could only draft what the fork recorded. That check did not disappear, it moved: the
reader before the stop now asks the same question of the whole document at once.

No checkboxes. No state. That is the next seam.

## Between 5 and 6 · The checklist, then a cold reader

`## Checklist`, the change document's last section — with everything the checklist always
guaranteed: one writer, derived from the contract, ordered by dependency, every box cut to
merge on its own. **Only the boxes under that heading are the checklist**; the *Task
breakdown* pillar is a list, never boxes, so nothing counts a box twice.

A box is one line: what it is, which box it waits on, and a `Closes:` naming the criteria it
closes — the token `converge` maps a box on. No headers, no sub-boxes, no code inside a box;
`libretto loop` counts every box as a task.

```
- [ ] one field, stored, served and proved · Closes: the field round-trips
- [ ] the second field · waits on box 1 · Closes: both fields validate
```

**A marked box that was never committed was never marked.** The mark ships in the same
commit as the work that closed it; a reviewer once found 0/24 boxes ticked on finished work,
every mark living only in a working tree.

**The session cuts it inline, and one fresh `spec-reviewer` reads it cold.** The session
that argued its way to an approach cannot tell what it wrote down from what it merely
decided, and it cuts boxes against the argument: it writes "wire up the adapter" and knows
what that means, and nobody else does — including the fresh session `libretto loop`
starts tomorrow.

So `review-spec` launches one `spec-reviewer` with the change document and the capability
specs it targets, and none of the conversation. It reads two things. **The contract**:
ambiguity that forks the work, criteria no run could fail, `Governs:` that describes
nothing, a criterion another capability contradicts. **The checklist, box by box**: what the
contract and the approach fail to answer for building that box — a mechanism never
described, a criterion with no proof saying what closes it. It names the gap and where the
answer belonged, and never proposes the answer. It writes nothing; its findings reach the
user unedited. Before reading, it runs `spec-drift --trace` and `--ears`, so what a script
can find is not spent on a reader. A reviewer that returns nothing is relaunched once; twice
silent is a stopped item, never a clean review.

Until 2026-10-21 that per-box question was a `task-cutter`'s, asked while cutting. Cutting
and reading do not need to be the same act; what needs to be fresh is the reading.

**One writer still.** The reviewer returns; the session owns the file and marks the boxes.
Several agents editing one markdown concurrently lose updates, and a checklist that
silently forgets a finished task is worse than no checklist.

**Then the one stop before code**, with the contract, the approach, the boxes and the
reviewer's findings all on the table at the same time — see *Where the flow stops*.

Not a numbered phase, on the same grounds as the 6→7 review: independence comes from the
fresh context, not from renumbering everything that says eight.

## 6 · Build and test

**The branch comes before the first write**, not before the first commit. `git checkout -b`
carries uncommitted work along, so editing on the base and branching later appears to work
until the base moves. Phase 1 usually made the branch already; phase 6 ensures it.

Logic that can break leaves a check behind. **Proportionately** — the smallest
thing that fails when the logic fails, not a suite per function:

| Change | Check |
|---|---|
| a branch, a loop, a parser, money, security, a trust boundary | one runnable check, minimum |
| behaviour a user depends on | end-to-end, once |
| a fix for something that broke before | a regression test, always |
| a one-line change with no logic in it | none — YAGNI applies to tests too |

Where the shape is already known, the check comes first.

Proportion is about **how many**, never about how honest. A test that exists is
never weakened, skipped or deleted to get a green run — see `skills/evidence/`.

**Commit before verifying.** A gate proves something about the tree it ran against,
so if the tree was dirty it proved nothing about what got recorded.
Work-in-progress commits are cheap and squash into the task's commit in phase 8.
Verifying uncommitted work and then committing it is a claim about code nobody
tested.

Branch per parent task, or per subtask when subtasks are genuinely independent.
Independent branches need chaining, not eight of them racing at the trunk.

Worktree when isolation is cheap. It is not always cheap: unversioned `.env`,
`node_modules`, vendored dependencies and generated files all have to be
reproduced before a worktree can build. Check that the tree builds from a clean
checkout before assuming isolation is free.

## Between 6 and 7 · Review

The flow's own rule is that nothing is true until observed — and until here, phase 7
was a self-report by the same agent that wrote the code. So before presenting, the
work goes to someone who did not write it: `review-work` launches one fresh
`work-reviewer` subagent with none of the session's context, and the builder's
beliefs about the code become the claim under review.

The reviewer reads the contract and the diff, **re-runs every proof the change
touches** rather than trusting phase 6's report, and returns findings — each one
citing a pillar or a proof, never taste.

Then the seam **fixes every one of them, and asks nothing**. That is not the reviewer
softening into an author: the reviewer is read-only and stays that way, because an
agent that repairs what it found is grading its own repair and its next verdict is no
longer independent of its own hands. The reading and the writing sit on opposite sides
of the seam, and that is the whole design.

The reason there is no question here is that there is nothing to ask. A finding cites a
pillar or a proof *by contract* — so it is a defect against something the user already
agreed to at phase 2. "Shall I fix the thing that violates the spec you approved?" has
one answer, and charging a round trip for it is charging for a rubber stamp.

Two bounds, both named rather than discovered later:

- **one fix pass, no re-review.** A fix that introduces a new defect is caught by the
  proofs or not at all. The replacement, the day that bites, is one bounded second look
  at the fix diff — never a loop.
- **two failures on one finding stops that finding**, per `skills/evidence/`. It reaches
  phase 7 as found-and-not-fixed, and so does anything needing a decision that is not
  ours. Reported, not asked.

Every finding, fixed or not, appends one entry to `.agents/lessons.md` with the phase field
`6→7`. That is what makes them countable: `libretto metrics` reads the field to say where
corrections surface, and `/libretto-retro` spends the ledger — fixing what belongs to the
project, proposing what belongs to the payload, never applying a payload diff itself.

`spec-drift` keeps the older standing — warn, never block — and the two are not
inconsistent. Drift is a question about whether a contract still describes the code, and
that answer is the user's. A finding is a breach of a contract already settled.

No spec, no review: the trivial lane collapsed the ceremony because there was no
contract to disagree with, and a reviewer with no contract has nothing to check. One
line, no wait.

Not a numbered phase, on purpose — independence comes from the fresh context, not
from renumbering everything that says eight.

## 7 · Present

Show what was done, before it is recorded and in the same turn. Every box has already
committed its own work; what phase 8 adds is the landing.

Three things, and the third is the one that usually goes missing:

- what was done, in the terms the spec used
- where the evidence is — the run, the test, the commit
- **what was deliberately left out, and when it should be added**

And beside them, attributed and unedited, the reviewer's verdict from the seam
above — the proofs it ran and what it found, next to the builder's own account.

That last line is what turns a simplification into a decision. "Did the single-user
case; the shared-state version needs a lock — say so if you need it now" is
reviewable. Silence about the same choice is a surprise waiting for whoever reads
the code next.

If the explanation is longer than the change, the explanation is the problem.

## 8 · Commit

Every finished task commits. Per task, not batched at the end, so the history
records the work and a bisect finds the break. Conventional commits, no AI
attribution.

**And the spec ships with the code.** Implementation always learns something the spec
did not know; that correction goes in the same commit as the code that taught it. Not
a follow-up, not a cleanup pass. A spec updated separately was wrong for however long
the gap lasted, and anyone who read it during that window was misled by a document
that looked current.

In flight, "the spec" means the change's contract: the pillars in `proposal.md` and any
`spec-<capability>.md` beside it. **When the change lands, one commit does four things**:
whatever code remains (often none — every box committed its own), the contract applied onto
each capability spec its `Targets:` names, the durable decisions retired into that spec's
*Prior decisions*, and the change folder deleted. All four or none — a contract applied
without removing its change leaves two documents describing one capability, and a folder
deleted without retiring its decisions takes the reasoning with it.

The third part is the one that gets skipped, so it is gated: `spec-drift --retired`, inside
`--anchors`, fails the landing commit unless some *Prior decisions* section moved in it. A
change that retires nothing says `Durable decisions: none` at the top of its document,
written with the approach, never added to get past the gate.

Before that commit, where the binary is on PATH: `libretto wiki` refreshes the specs index
so it rides the same commit, and `libretto land <change>` checks the staged landing — the
folder gone, the capability spec moved — and names any part that is missing.

This is the whole of Spec-Anchored, and it is one question before each commit: did
this change teach the spec anything?

`skills/record-work/spec-drift` answers it in several modes:

| Mode | What it checks | Exit |
|---|---|---|
| *(none)* | staged code whose governing spec did not move | always 0 |
| `--anchors` | every `Proof:` resolves to a test that exists · delta criteria carry `shall` · no reasoning retired silently · fences closed | **fails** — one of the six gates |
| `--trace` | the whole tree: orphan code, dead `Governs:`, unproven and joined criteria | always 0 |
| `--block` | the default mode, but exit 1 — for a `pre-commit` hook somebody chose to install | 1 on drift |

The default mode **warns rather than blocks**. Enforcement that surprises someone in their
own project gets uninstalled, and an uninstalled gate catches nothing. A check that stops a
commit in someone else's project is a check that gets deleted, and a deleted check finds
nothing. Whoever wants it to block installs the opt-in hook — one line, documented in
`record-work/drift.md`, never installed by the flow. That hook is a git hook, so the old
objection — that the barrier would need `settings.json` — no longer applies.

**The push is asked at the very end, once, and never assumed** — push and open the request,
no, or push without the request. With no remote there is no question: the phase ends at the
commit and says so. Everything that needs a remote — which forge from `git remote get-url
origin` (`gh` or `glab`), the attacca exception, the bump, the return to base — lives in
`record-work/forge.md`, read only when a remote exists. When the bump is already decided, the
label rides `gh pr create --label` instead of a second command after it.

### And then the tree goes home

When the answer is yes and the request is confirmed open, phase 8 checks out the base
branch and fast-forwards it. Not a convenience: **a session starts wherever the last one
left the working directory**, so a flow that ends parked on a merged feature branch hands
the next phase 1 a base that is behind the remote — and phase 1's whole job is reading what
is in flight off exactly that.

That is measured, not theorised. Phase 1 of the run that added this reported a branch as
work in flight, offered the user a choice about it, and was wrong on both counts: the branch
had already been merged and tagged `v0.6.1`, and local `main` was seven commits behind. The
wrong reading was stated as fact.

Three things it deliberately does not do:

- **fetch the base ref without leaving the branch.** `git fetch origin main:main` is cheaper
  and keeps you where you are, and it was offered and declined — it fixes the ref, not the
  place the next session starts.
- **delete the feature branch.** The request is open, not merged. A branch deleted there
  takes the only local copy of unmerged work with it.
- **merge anything.** `--ff-only`, always. A merge commit manufactured on somebody's base
  branch by the bookkeeping phase is the surprise nobody wants, and a diverged base is a
  fact worth seeing rather than resolving unasked.

On the no path it does nothing at all. Nothing was pushed, so the branch is the only place
the work exists, and moving off it buys nothing. The base is read off the repository
(`origin/HEAD`), never hardcoded, and a dirty tree or a refused fast-forward stops it in
place.

**The merge is not the flow's.** It ends at a green request; merging, tagging and releasing
are the user's, and the `release:` label on the request is what turns a merge into a tag.

## Where the flow stops, and why only there

Two stops on the contract lane: one before code, one at the door. The trivial lane keeps
only the door.

| After | Stops | What the user is changing |
|---|---|---|
| 1 · find the work | no | — the reading is stated, and the contract is where it gets corrected |
| 2–3 · the contract | no | — agreed below, beside what building it costs |
| 5 · the approach | no | — agreed below, beside the boxes it produced |
| 5→6 · the checklist | no | — agreed below, in the same decision |
| before 6 · the change document, after the cold read | **yes** | the whole change document: contract, approach, order |
| 6 · build | no | — |
| 6→7 · review | no | — the seam fixes what it finds |
| 7 · present | no | — |
| 8 · commit | **yes, last** | whether the world sees it |

**A stop is a place where the user changes something.** That is the entire test. A stop
whose only available answer is "yes, carry on" is a round trip charged for a rubber
stamp, and it does not become one by being called a decision point.

**One stop before code, not one per document.** Until 2026-10-21 there were two stops
before code — the contract after phase 2, then the approach and the order after the tasks
were cut. The first asked the user to agree to a contract before seeing what building it
would cost, and the second asked again once it was visible. One stop with everything on
the table is one decision, made once, with the information that decides it.

**It is the cheap place to disagree, and that is why it comes before code rather than
later.** A wrong sentence in the contract costs a line to fix. The same misunderstanding
costs a day once it has become code, and by phase 7 it costs the reviewer measuring the work
against a promise nobody could have failed. Every stop after this one is more expensive than
this one, which is the argument for making the contract the first thing anybody reads.

### What `attacca` answers, and what it cannot

`/libretto-attacca` answers both stops in advance — the invocation *is* the answer, given
before the work starts. What it will not answer is a **gate**: a failing check stops it where
it stands, and it never merges, tags or releases. That boundary is the whole of what makes the
mode safe rather than merely fast, because a mode that answered a gate would not be unattended,
it would be unverified.

The version bump is the sharpest case. It turns on whether a promise in `.agents/specs/` moved,
which is a reading of the specification rather than of the commit types — and a version number
cannot be recalled once the module proxy has cached it. So the reading stays the user's even
under a mode built to answer everything else: it is **the one question an unattended run asks,
last of all**, after the request is open. The run recommends its own reading, never
recommends `release:major`, types the label the user chose, reads it back off the forge, and
writes who chose it into the request's description. Unanswered — a headless `libretto loop`
— the request stays unlabelled, and its red check *is* the question. Attended runs do not ask
it; the user is already at the door.

**And every one of them is a native question**, `AskUserQuestion`, never a sentence at the
end of a report — along with phase 1's choice between work already in flight and something
new. Two stops, two prompts, the same shape each time: the recommended option first
saying what will actually run, the real alternatives, and room to answer differently.

The argument was written once for the push, in `skills/record-work/`, and was never
push-specific. It stays written once — this paragraph states the rule and points at it,
because a reason copied into six files is six things to keep in sync and five of them will
not be.

Phase 1 used to stop, and what it bought was the user reading a paraphrase of their own
sentence back to them. Phase 7 used to stop, and what it bought was permission to commit
to a local branch nobody had seen — the cheapest possible place to change your mind,
guarded as if it were a deployment.

The cost of getting this wrong is not abstract. Four stops for a change is four turns of
latency, and the flow had already conceded that a typo should not pay them. A concession
held for typos and refused everywhere else is not a proportionate gear, it is a loophole,
and what people do with a flow that charges too much is route around it.

**Two exceptions in phase 1, and neither is ceremony:** work already in flight, where
continuing it or not is a choice about the user's priorities, and a missing or
unconfigured tracker, where nothing downstream exists. Both are the input failing to
arrive, not a phase boundary asking to be blessed.

### How an unattended run behaves

`/libretto-attacca` — the score's instruction to go on to the next movement without
pausing — runs the same phases and ends at a pushed branch with a request open on it. A
failing gate still stops the run, twice on one task still stops the task, a missing
credential still stops the phase, and no run reaches for `--force` or a weakened test. A run
that stops opens no request: there is nothing reviewable to put in one.

Given no task, it takes the oldest change with open boxes. It never drains the queue — that
is `/libretto-next`, which asks.

A question the flow cannot settle from the code is not asked either — it becomes an
assumption written into the decision log and *Prior decisions* marked as assumed, into the
report, and into the request's description, with what changes if it is wrong. That is the
same rule the flow already applies after the plan, moved to the front.

**The cold reader matters more here, not less.** Nobody reads the document at the stop, so
the `spec-reviewer` is the only thing between an unfailable criterion and the code built on
it. It still only reports. Unattended, the run is the document's author, so it acts on the
findings before phase 6 and names each one it acted on — in the report and in the request,
exactly like an assumption. A contract rewritten between review and code without saying so
is a contract nobody agreed to.

## Three rules, not phases

**Ask** (4) and **commit** (8) are written above as phases because that is where
they first appear, but neither happens once. Committing happens per task throughout
phase 6, not after phase 7. A flow that treats it as a phase will do it once and
consider it done.

**Asking is bounded to before the plan.** Phases 1, 2 and 5, where nothing has been
built on the answer yet — which is why the one stop before code sits right after them. After the plan
an unsettled question becomes a *finding*: it goes into the phase 7 report with what was
assumed in the meantime and what changes if the assumption is wrong, and the user meets
it at phase 8 alongside everything else.

That bound is the load-bearing part. "Ask whenever you need to" is reasonable in every
individual case and collapses the whole promise, one reasonable exception at a time.

**Evidence** is the third, and it is not numbered above at all: nothing is true
until it has been observed. Failing tests get fixed or reported, never edited into
silence. Commands whose result decides the next step run in the foreground with
their output read, because a pipe destroys the exit code. Nothing is reported as
done that was not seen. Two failed gates stop an item; two stopped items stop the
session. See `skills/evidence/`.

## Delegation

A sub-agent starts with none of this. Not the flow, not the standing rules, not what
was decided three turns ago. Whatever it is not told, it does not know — and it will
fill the gap with something plausible.

**So the session writes, and sub-agents read.** The change document is written inline by
the session that held the conversation, because that is where the answers are. A fresh
context is spent where freshness is the point — reading what was written without the
memory of what was meant:

| Agent | Launched by | Writes |
|---|---|---|
| `spec-reviewer` | `review-spec`, once, before the stop | nothing — findings come back |
| `work-reviewer` | `review-work`, once, between 6 and 7 | nothing — findings come back |
| `spec-writer` | `write-spec`, only when an extra capability's research will not fit in context | its one `spec-<capability>.md` |
| `review-lens-*` (five) | `review-project`, behind `/libretto-review` — a PR/MR review, outside the flow | nothing — findings per lens |

A launched reader may run in the background — the host decides that. What the flow requires
is that nothing proceeds until its return is read.

**Context goes down through a file, findings come back through the return value.**
The shared ground is `proposal.md`, read by path; a writer writes exactly one file that
nobody else writes; nothing is reported back by editing shared state. Every file has one
author. That is what makes running several at once safe, and it is the same rule the
checklist follows with the session as its single writer.

Four prohibitions travel with every launch, because each one is invisible from inside
a sub-agent's view of the work:

- **no commit** — only the orchestrator sees the whole task
- **no push** — that is a question for the user, and a sub-agent cannot ask it
- **no writing the change document or a sibling's file**
- **read-only everywhere except the one path assigned**

And the orchestrator owes them two things in return.

**Silence is not success.** An agent that died, timed out or returned nothing did not
finish its work. Marking a box because nobody complained invents a "done" that no one
observed.

**A way to be stuck.** There is no channel while a sub-agent runs, so being blocked
has to be a legal outcome or the default takes over — and the default is to guess
something plausible. A blocked sub-agent writes what it can, marks the gap where the
answer belongs, returns the question, and invents nothing past it. **Questions funnel
through the orchestrator**: several agents asking directly would interleave questions
from work the user cannot see, each answer arriving without the context that made it a
question.

**And the tier is chosen before the fan-out, not during it.** A phase's context is billed
at a fraction of the input price while its prefix stays byte-identical between calls;
changing the model or the effort level invalidates that prefix and rebills everything at
full price. Fanning out is where that is dearest — N writers is N contexts, and a switch
part-way through pays for all of them twice. Between phases the prefix is new anyway, so
that is where a change of tier is free.

The rule is stated in `skills/write-spec/` too, beside the fan-out it costs. It stays a
statement in both places rather than a check: the dial is the session's, and nothing in
the payload can read it.

Phase 6 does not fan out. Parallel implementation needs isolation per task, a serial
queue for merges, and a conflict protocol — without those three, concurrency
manufactures races. The order to add them, when it is time, is worktrees, then the
queue, then the parallelism.

## What ships, and what is merely used

Six skills written by other people **ship with this repository** —
`test-driven-development` and `using-git-worktrees` from
[obra/superpowers](https://github.com/obra/superpowers); `ponytail` and
`ponytail-debt` from [DietrichGebert/ponytail](https://github.com/DietrichGebert/ponytail);
`caveman` and `caveman-commit` from
[JuliusBrussee/caveman](https://github.com/JuliusBrussee/caveman). All MIT. The flow's
own skills are thin because they delegate to those, and a thin skill whose delegate
is missing is not thin, it is broken. Being installed on the author's machine is not
the same as being installed on yours. See [`THIRD-PARTY.md`](../THIRD-PARTY.md).

## The companions, and why they ship now

**ponytail** decides how much should be built. Its ladder runs from "does this need
to exist at all?" down to "only then, the minimum that works", and it carries the
list of things that are never trimmed — trust boundaries, data loss, security,
accessibility. This flow invokes it in phase 2, on requirements, which is where
removing work is cheapest. Its `ponytail:` comments and their harvested ledger —
`/ponytail-debt` — feed the prior-decisions pillar. It also sets its own intensity;
this flow reads that, it does not define one.

**caveman** decides how much gets said. It compresses prose, ponytail compresses
what gets built, and they do not overlap. `caveman-commit` is the same compression
applied to commit messages, offered by phase 8.

Until 2026-08-10 both were companions: called when present, never vendored, on the
grounds that the user may already have chosen a version. That assumed a user who has
versions of things — the flow's target is a machine with nothing on it, where every
"if installed" was a conditional that never came true. Now the two cores, the ledger
and the commit generator ship vendored; the rest of both plugins does not, because
the flow calls none of it. **Shipped is still not required**: nothing fails without
them, and `prune`/`uninstall` removes them like any other item. The upstream plugins
remain the path to what is deliberately not vendored — always-on hook mode, and both
copies coexist by namespace.

## Reading cost

A skill's body loads whole every time it is invoked, so guidance that applies only on a
condition lives in a sibling file that `SKILL.md` names with its condition:

| Skill | Sibling | Read when |
|---|---|---|
| `find-work` | `tracker.md` | a tracker key or URL was given |
| `find-work` | `queue.md` | the caller asks for the queue |
| `write-spec` | `fan-out.md` | the change amends more than one capability |
| `record-work` | `forge.md` | the repository has a remote |
| `record-work` | `drift.md` | running `--trace` or wiring `--block` |

A trial run on 2026-10-23 read about 3,000 lines of skill text to build 15 lines of code; the
split skills carry a line budget in `check-payload` so the bodies do not grow back.

## The flow outside a session

The CLI touches the flow in a few places, all of them reading what the payload writes:

| | |
|---|---|
| `libretto init` | bootstrap `.agents/` in a project with none |
| `libretto loop <change>` | one fresh `/libretto-attacca` session per open box, until none are open; stops after two rounds that close nothing, at `--max` (10), or at `--max-budget-usd`; never pushes. `--opencode` drives OpenCode instead |
| `libretto queue` | the queue, by `find-work/queue.md`'s rule |
| `libretto metrics [<change>]` | what each change cost — commits, span, boxes closed and reopened, where corrections surfaced — from git and the ledger |
| `libretto land <change>` | check a staged landing before its commit |
| `libretto wiki [--html]` | the specs as a README index, or one self-contained page |

`status`, `doctor`, `metrics` and `queue` take `--json`.

## Settled questions

Recorded so the reasoning is not lost.

**Where the artifact gets looked at.** This repository's first
palette satisfied its spec and was unreadable — 1.4:1 on borders. What caught it was
`docs/preview.py`, a throwaway that rendered the panel before the theme existed, and
a WCAG measurement of that render. The answer became a rule inside phase 6, not a
phase of its own: **if judging the change means looking at it, the builder renders it
and looks before the review seam**, measured contrast where the change is about
colour, with what was seen carried in the evidence. See `skills/build-and-check/`.
The reviewer in the 6→7 seam still reads specs, diffs and test output, never pixels —
which is exactly why the look happens before it.

**Who keeps the spec true.** Phase 8 commits the spec alongside the
code that taught it — see `skills/record-work/`. The three divergences in this
repository's own `SPEC.md`, accumulated over one phase of work, are what the
alternative looks like.
