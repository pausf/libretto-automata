# The flow

The payload. What Libretto Automata installs into `~/.claude`: eight phases, two
standing rules, as skills and commands.

Read [SPEC.md](SPEC.md) for the CLI that delivers this. The two are independent —
a phase is a markdown file with frontmatter and can be written today.

## 1 · Read the task

Start from the tracker. One task, its subtasks, nothing invented.

Everything through the CLI. No MCP.

```
jira issue view EUCAR-1234 --plain
jira issue list -q "parent = EUCAR-1234" --plain
```

Before that runs, a preflight — this ships to people who have configured nothing:

```
which jira                            → missing: brew install ankitpokhrel/jira-cli/jira-cli
test -f ~/.config/.jira/.config.yml   → missing: the user runs `jira init`
                                      → present: read `server`, continue
```

Ask for the domain once and let `jira init` persist it. Do not build a second
configuration store, and do not write state into the skill — the skill is a
symlink into this repository.

**The API token never enters the conversation.** `jira init` is interactive and
the user runs it in their own shell. A token pasted to an agent lands in the
transcript, the logs and the history. The domain may be asked for. The secret may
not.

Detect the git host, never assume it. `gh` and `glab` are not interchangeable and
a machine may have either, both, or neither.

Preflight belongs in `libretto doctor`, which already answers "is everything present?"
for symlinks.

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

## 2 · Write the spec

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
| **Verification criteria** | how each unit is proven — each one naming the test that proves it |

### Where it goes

```
.agents/
├── specs/                        the living, consolidated specification
│   └── bundle-products/
│       └── spec.md               Governs: src/bundles/**
└── changes/                      proposals in flight
    └── add-relative-discounts/
        ├── proposal.md
        ├── spec.md               Targets: bundle-products
        └── plan.md
```

**One spec per capability, never per ticket.** A capability spec accumulates and
stays true; a spec named after a ticket is dead the day the ticket closes, and a
directory of dead specs is worse than none because people still read them.

Amendments live in the change while the work is in flight, and phase 8 applies them
onto the capability spec and deletes the change. Exactly one current specification,
always.

## 3 · One spec per subtask

When the task has subtasks, each gets its own spec.

When there are many, a sub-agent per spec gathers the context that spec needs and
writes it out as markdown. Research fans out; the specs stay separate files so
sessions never collide.

What keeps a fanned-out set coherent is a **brief written once before any sub-agent
starts** — shared conventions, inherited constraints, settled decisions, and the
vocabulary every spec will use. Each sub-agent writes against it and reports back
what the brief got wrong. Those reports are where the real findings are.

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

So phase 2 asks before it writes the file — **one call, never a string of turns.** Serial
round trips to build one contract is the ceremony this flow spends its length arguing
against, and the answers have to be *in* the contract rather than bolted on afterwards.

The count was capped at three until 2026-08-14, when the user lifted it: better asked out
of caution than swallowed out of fear. What replaced the cap is judgment with both edges
named — every question one a wrong guess would make expensive, and never a form-length
interrogation of things the code already answers. The bias, when in doubt, is to ask.

**Phase 2 alone, not phase 5**, which was offered and declined. The contract is where an
answer changes the most and costs the least to change; phase 5 already stops for the order
and what waits on what. Asking twice before the first line of code is how a three-stop flow
grows back into a nine-stop one, one defensible exception at a time.

**Zero is legitimate, and it is said in one line.** The alternative — always three, stretched
if necessary — was offered and declined for the reason `AGENTS.md` already gives: do not ask
what the code can tell you. A quota manufactures the rubber-stamp question this flow removed
from three other phases.

Note what does *not* change: the stop count. The questions ride the stop phase 2 already
has.

## 5 · The plan

One markdown file. One line per task, checked off as work lands, so an agent
joining late reads current state instead of guessing.

**One writer.** The orchestrator owns this file. Sub-agents report completion and
the orchestrator marks the box. Several agents editing one markdown concurrently
lose updates, and a checklist that silently forgets a finished task is worse than
no checklist.

## 6 · Build and test

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

`spec-drift` keeps the older standing — warn, never block — and the two are not
inconsistent. Drift is a question about whether a contract still describes the code, and
that answer is the user's. A finding is a breach of a contract already settled.

No spec, no review: the trivial lane collapsed the ceremony because there was no
contract to disagree with, and a reviewer with no contract has nothing to check. One
line, no wait.

Not a numbered phase, on purpose — independence comes from the fresh context, not
from renumbering everything that says eight.

## 7 · Present

Show what was done. Then commit it, same turn.

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

In flight, that means the delta in the change. **When the change lands, one commit
carries the final code, the delta applied onto the capability spec, and the deletion
of the change folder.** All three or none — a delta applied without removing its
change leaves two documents describing one capability, and nobody can tell which is
current.

This is the whole of Spec-Anchored, and it is one question before each commit: did
this change teach the spec anything?

`skills/record-work/spec-drift` asks it mechanically, from the staged index, and **warns rather
than blocks** — always exit 0. Enforcement that surprises someone in their own project
gets uninstalled, and an uninstalled gate catches nothing. A check that stops a commit in
someone else's project is a check that gets deleted, and a deleted check finds nothing.
Whoever wants it to block can wire it into a `pre-commit` hook or CI; that choice belongs
to them.

Being honest about the level this reaches: a warning is Spec-First discipline with a
reminder, not the automatic barrier the Spec-Anchored definition asks for. The barrier
would be a hook, and `settings.json` is out of scope for this project by an earlier
decision. The gap is deliberate and worth naming rather than papering over.

The push is asked at the very end — yes or no — and never assumed.

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
the work exists, and moving off it buys nothing.

## Where the flow stops, and why only there

Three stops. Two inside the work, one at the door.

| After | Stops | What the user is changing |
|---|---|---|
| 1 · read the task | no | — the reading is stated, and the spec is where it gets corrected |
| 2 · the spec | **yes** | the contract |
| 5 · the plan | **yes** | the order, and what waits on what |
| 6 · build | no | — |
| 6→7 · review | no | — the seam fixes what it finds |
| 7 · present | no | — |
| 8 · commit | **yes, last** | whether the world sees it |

**A stop is a place where the user changes something.** That is the entire test. A stop
whose only available answer is "yes, carry on" is a round trip charged for a rubber
stamp, and it does not become one by being called a decision point.

**And every one of them is a native question**, `AskUserQuestion`, never a sentence at the
end of a report — along with phase 1's choice between work already in flight and something
new. Three stops, three prompts, the same shape each time: the recommended option first
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

### All three answered in advance

`/libretto-attacca` — the score's instruction to go on to the next movement without
pausing — runs the same flow with the three stops answered by the invocation itself, and
ends at a pushed branch with a request open on it. Nothing about the phases changes: the
spec, the plan, the commits and the report are all still written.

**What it cannot answer is a gate**, and the distinction is the whole of the feature. A
stop is where the user changes something; a gate is where the code is measured. A failing
gate still stops the run, twice on one task still stops the task, a missing credential
still stops the phase, and no run reaches for `--force` or a weakened test to get past
one. **A mode that answers a gate is not unattended, it is unverified.**

A question the flow cannot settle from the code is not asked either — it becomes an
assumption written into the spec marked as assumed, into the report, and into the
request's description, with what changes if it is wrong. That is the same rule the flow
already applies after the plan, moved to the front.

## Three rules, not phases

**Ask** (4) and **commit** (8) are written above as phases because that is where
they first appear, but neither happens once. Committing happens per task throughout
phase 6, not after phase 7. A flow that treats it as a phase will do it once and
consider it done.

**Asking is bounded to before the plan.** Phases 1, 2 and 5, where nothing has been
built on the answer yet — which is why the two stops sit exactly there. After the plan
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

**Context goes down through a file, findings come back through the return value.**
The shared ground is written once to a brief and read by all of them; each writes
exactly one file that nobody else writes; nothing is reported back by editing shared
state. Every file has one author. That is what makes running several at once safe,
and it is the same rule the plan follows with a different single writer.

Four prohibitions travel with every launch, because each one is invisible from inside
a sub-agent's view of the work:

- **no commit** — only the orchestrator sees the whole task
- **no push** — that is a question for the user, and a sub-agent cannot ask it
- **no writing the plan, the brief, or a sibling's file**
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

Seven skills written by other people **ship with this repository** — `writing-plans`,
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

## Open

Not decided yet. Recorded so they are not lost.

**Settled since — where the artifact gets looked at.** This repository's first
palette satisfied its spec and was unreadable — 1.4:1 on borders. What caught it was
`docs/preview.py`, a throwaway that rendered the panel before the theme existed, and
a WCAG measurement of that render. The answer became a rule inside phase 6, not a
phase of its own: **if judging the change means looking at it, the builder renders it
and looks before the review seam**, measured contrast where the change is about
colour, with what was seen carried in the evidence. See `skills/build-and-check/`.
The reviewer in the 6→7 seam still reads specs, diffs and test output, never pixels —
which is exactly why the look happens before it.

**Settled since:** who keeps the spec true. Phase 8 commits the spec alongside the
code that taught it — see `skills/record-work/`. The three divergences in this
repository's own `SPEC.md`, accumulated over one phase of work, are what the
alternative looks like.
