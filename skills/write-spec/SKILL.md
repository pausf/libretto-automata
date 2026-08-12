---
name: write-spec
description: "Trigger: a task has been read and needs a specification; writing or extending a spec; splitting a large task into several specs; recording a decision or an answer into a spec. Phase 2 of the flow. Writes the contract before any code exists."
license: MIT
metadata:
  author: pausf
  version: "1.2"
---

## What this does

Phase 2 of the Libretto flow: turn a task into a contract.

Input is a task already loaded by `find-work` — key, summary, description,
subtasks. Output is one or more spec files and nothing else.

**No code. No branch. No commit. No plan.** Those are phases 5, 6 and 8. A spec
session that also writes code has stopped being a spec session, and the contract
becomes a description of whatever got typed.

Read `skills/evidence/` first. It applies here too: the criteria this skill writes
are what the whole flow later gets measured against.

## Step 0 — Does this spec need to exist?

Ask it before anything else, and be willing to hear no.

A one-line change, a typo, a version bump, a rename with no behaviour behind it —
none of these need a contract. Six pillars, a brief and a plan around a
two-character fix is the exact over-engineering this project exists to prevent,
applied to itself.

When the task is that small, say so, hand back to the flow, and let it go straight
to the change, its check and a commit. **Skipping this phase is a legitimate
outcome of it.**

**And the "no" collapses the ceremony with it.** A change that needed no spec does not
get phase 7's gate either: the report and the commit land in the same turn, and the
only question asked is the last one — push and open the request. Phase 8's question
survives every collapse, because it is the user's decision and not ceremony.

The reason is that the gate exists to catch a disagreement about "done", and this
phase just established there is nothing to disagree about. Keeping the stop anyway
charges a typo the price of a feature, and a flow that costs four round trips for a
typo gets routed around — for typos first, then for small features, until what is
left is a ritual reserved for work important enough to deserve it.

**Say the "no" out loud, in one line.** A phase that declines and reports it is not
the same as a phase nobody ran, and from outside they look identical. That is a real
failure with a history: a session skipped this phase and phase 6 entirely, made the
same judgment this phase would have made, and the user could not tell the difference
between a considered skip and a forgotten one.

The test is not the diff size — it is whether anyone could reasonably disagree about
what "done" means. If they could, write the spec. If they could not, do not.

This is `ponytail`'s first rung applied to our own artifacts — it ships with this
flow. It states it as: does this need to exist at all?

## Step 1 — One spec, or several?

Ask what the shape is before writing anything.

**One spec** when the task is a single coherent change that fits one session with
room to spare.

**Several** when it does not. Split by **coherence**, not by parallelism — one spec
holds one idea that can be argued about on its own. Splitting to make agents run
concurrently produces specs that each describe half an idea and cannot be reviewed.
If the subtasks in the tracker already cut along coherent lines, use them; if they
cut along assignment lines, ignore them and cut again.

A spec whose scope you cannot state in one sentence is two specs.

## Step 2 — Find the ground truth before writing

A spec written without looking at the code is fiction, and the pillar it fails
first is **prior decisions**: existing conventions, patterns and constraints that
the change has to live inside.

So read the code first. Specifically:

- how this area is structured today and what it already does
- the conventions actually in use — naming, layout, error shape, test style
- what will break if this change is naive
- whether a decision recorded elsewhere already settles part of this

With one spec, read the code inline. A sub-agent for a change you can hold in your
head is a round trip for nothing.

## Step 2b — Several specs: the shared brief, then fan out

Ten subtasks researched in this thread will exhaust the context, which is the exact
problem fanning out exists to solve. So sub-agents write their own spec files.

**The context is the only reason.** Fanning out because the answer is not yet
obvious turns a reflex into a research project — if two readings of the code
already settle it, settle it and move on. Ten agents dispatched to confirm
something visible from one file is the expensive kind of thoroughness.

The risk that creates is a set of specs that each make sense alone and contradict
each other — different names for one concept, incompatible assumptions, the same
decision made twice two ways. Consistency has to come from somewhere other than a
single author.

**It comes from a brief written once, here, before any sub-agent starts.** Read the
shared ground and write it to `.agents/changes/<change>/brief.md`, inside the change
whose deltas it will produce:

- the conventions actually in use — naming, layout, error shape, test style
- the constraints everyone inherits
- the decisions already settled, and by whom
- **the vocabulary** for every concept the specs will share, so two agents cannot
  give one thing two names
- the six-pillar structure they all fill in

A file, not a paragraph held in this conversation. Three reasons, and the third is
the one that matters:

1. Ten sub-agents given the text verbatim is ten copies consuming ten contexts. One
   path is one line.
2. It is auditable. What they were told is on disk, so a spec that came out wrong
   can be traced to the instruction that produced it.
3. **It can be re-run.** When the brief turns out to be wrong, the fix is to correct
   the brief and regenerate the affected specs — which is impossible if the brief
   only ever existed inside a prompt.

Commit it with the deltas. It is the provenance of the whole change, and it is as
temporary as the change — when the deltas land on their capability specs, the brief
goes with the folder. What survives is what the deltas carried into the specs.

### The bridge runs one way

```
orchestrator ──writes once──▶ brief.md ──read-only──▶ N sub-agents
sub-agent    ──writes──▶ its own spec file       (one author per file)
sub-agent    ──returns──▶ findings in its reply  (never to a file)
```

**Sub-agents never write to the brief.** N writers on one file is the lost-update
race that `skills/write-plan/` exists to avoid — two reads of the same content, two
writes, the second silently discarding the first. Findings come back in the return
value.

Every file in this phase has exactly one author. That is the property that makes the
fan-out safe, and it is the same rule as the plan's, with a different single writer.

### What each sub-agent is given

Each one is launched as the `spec-writer` agent — that contract exists for exactly
this seat, and a generic sub-agent would start without it. In its prompt:

- the path to the brief, and the instruction to read it first
- its own subtask — key, summary, description
- the boundary between its spec and its siblings', stated explicitly
- the one path it may write, and nothing else

And the standing rules, because **a sub-agent starts with none of them**:

- invoke `evidence` before anything: nothing reported that was not observed, no gate
  in the background, no test weakened to go green
- **do not commit.** The orchestrator commits, because only it sees the whole task
- **do not push.** Ever. Pushing is a question for the user, and a sub-agent cannot
  ask it
- **do not touch the plan**, the brief, or a sibling's spec
- read-only everywhere except the one file assigned

### When a sub-agent gets stuck

There is no channel while it runs. It cannot ask the user — it does not have the
conversation — and it cannot ask the orchestrator, because the return value happens
once, at the end.

So the rule has to be explicit, or the default takes over: **a model with no
instruction guesses plausibly**, which is the exact failure this flow exists to
prevent.

When a sub-agent hits something it cannot legitimately decide:

1. **Write everything it can.** A spec with five solid pillars and one open question
   is useful. Nothing is not.
2. **Mark the gap in the file itself**, where the answer belongs, stating what is
   unresolved and what turns on it. Not a placeholder that reads like an answer.
3. **Return the question**, phrased so the orchestrator can ask it without having to
   reconstruct the context.
4. **Never invent past it.** Not a default, not the more common of two precedents,
   not "the obvious choice". A guess that looks like a decision is worse than a hole,
   because a hole gets filled and a guess gets built on.

Questions funnel through the orchestrator, always. If sub-agents could ask directly,
several would interleave questions from work the user cannot see, and each answer
would arrive without the context that made it a question.

### Sub-agents do not talk to each other

By design — no shared mutable state, no negotiation, no ordering to get wrong.

The consequence is that a boundary drawn badly between two specs is not fixed by the
two of them noticing. It surfaces as two returns describing overlapping ground, and
**the orchestrator resolves it after the fact**: one decision, applied to both files,
recorded once.

Both sub-agents were right about their own scope. The brief was wrong about the line
between them, and that is a finding, not a failure.

### What comes back

The path it wrote, and a short list of what it found that the brief did not know — a
convention the brief missed, a constraint that contradicts it, a question only the
user can settle.

**Read every one of those returns before accepting the set.** They are the
interesting part: a sub-agent that found the brief wrong found something real. When
two returns disagree about the same thing, that disagreement is a decision, and it is
resolved here — once, then applied to both files. Never left as two specs each
quietly assuming it won.

**Silence is not success.** A sub-agent that died, returned nothing, or came back
empty did not finish. Do not accept its spec because it failed to complain.

If the brief turns out to be wrong on something structural, correct the brief and
re-run the affected specs. Patching the files individually reproduces the drift the
brief was there to prevent.

## Step 3 — Write the six pillars

Every spec has all six, explicitly, with headings. A pillar left vague is not
brevity — it is a gap that gets filled later by whoever is typing, from
association rather than intent.

**Draft them, then hold the file until step 4 has its answers.** Drafting is what
surfaces the questions worth asking — a boundary the request left open, two precedents
with nothing choosing between them — and an answer that arrives after the file exists
gets bolted onto a contract instead of shaping it. The steps are numbered 3 then 4
because the *thinking* runs in that order; the writing happens once, after both.

### Outcomes

The verifiable end state, in operational terms. Not "add login" but what is true
when it works: what the user does, what the system emits, what persists, what
happens on a refresh.

### Scope boundaries

What is in. Then, at least as carefully, **what is out** — named, so it cannot be
quietly added. "Federated login is out of scope for this change" prevents an
afternoon of work nobody asked for.

This is the cheapest place in the entire flow to remove work. Deleting a
requirement here costs a line; deleting the code it would have produced costs a
day. So every candidate requirement gets asked whether it needs to exist at all,
and the ones that do not become named non-goals rather than silent omissions.

`ponytail`, which ships with this flow, is the discipline for this — its ladder runs
from "does this need to exist?" down to "only then, the minimum that works". Use it
here, on requirements, not only later on code. Applied at build time it can shrink
how; applied here it shrinks what.

**Five things are never scoped out, whatever the ladder says:**

- validation at a trust boundary
- error handling that prevents data loss
- security — and a secret that must not be stored or transmitted is part of the
  contract, not a detail of the implementation
- accessibility basics, including contrast and keyboard reachability
- anything the user asked for explicitly

A spec that trims one of these is not lean, it is incomplete. The palette that
satisfied its spec at 1.4:1 contrast failed on the fourth of these, and the spec
that permitted it was the thing at fault.

### Constraints

What the solution has to live inside: existing versions, data shape, latency
budgets, conventions, anything the surrounding system imposes.

### Prior decisions

What must not be relitigated, and why — from step 2, plus anything the user has
already settled. This is where an answer given once stops being asked again.

A deliberate shortcut belongs here too, and it belongs here **with its ceiling
named**: what it will not survive, and what replaces it when that day comes. "A
global lock, adequate below N writes a second, per-account locks if throughput
matters" is a decision. The same shortcut with no ceiling written down is hidden
debt, and the next session cannot tell the two apart.

This is the same information ponytail's `ponytail:` code comments carry, and
`/ponytail-debt` — shipped alongside it — harvests them into a ledger. Read that ledger
when writing this pillar — a shortcut already marked in the code is a prior
decision that should not be rediscovered from scratch.

### Task breakdown

Atomic units. Each one independently assignable, each one small enough that a
single failure does not take the rest with it. This becomes the plan in phase 5;
it is not the plan yet.

### Verification criteria

How each unit is proven — **each criterion naming the test that proves it**, by
file and case.

The tests do not exist yet. Name them anyway. That citation is what makes them get
written, and it is what makes this pillar falsifiable instead of decorative. A
criterion with no named test is a hope, and it will be treated as satisfied by
default.

Write citations in a form a tool can check:

```
Proof: internal/link/apply_test.go TestApplyIsIdempotent
```

`spec-drift --anchors` — the script beside the sibling `record-work` skill, at
`<skill-base>/../record-work/spec-drift` from this skill's announced base
directory — resolves every one of them and fails on a citation
whose file or whose **test** is absent. The test name is the part that matters: a
file-level check passes an invented test name as long as some file sits at that
path, which is precisely the fabrication the anchor exists to prevent. It caught
one in this repository's own spec.

### And the spec declares what it governs

One line per spec file:

```
Governs: cmd/** internal/** scripts/**
```

That is the other half of the anchor. Without it, "code changed but the spec did
not" is a guess about the whole repository. With it, the question is precise: this
path is governed by that spec, and that spec did not move.

A path no spec claims is reported separately and more softly — not everything needs
a contract, and treating every unanchored file as an error trains people to ignore
the warning.

## Step 4 — Ask what cannot be settled here. Up to three.

Some things are genuinely not yours to decide: a product tradeoff, a convention
with two live precedents in the codebase, anything where guessing wrong quietly
breaks working behaviour.

**Ask up to three of them, in one `AskUserQuestion` call, before the file is
written.** Each with the option you recommend, the real alternatives, and room to
answer differently.

Drafting the six pillars is what surfaces them — a boundary the request left open, a
scope decision that changes what gets built, two precedents with nothing choosing
between them. So the questions come after step 3's thinking and before step 3's
file. **The answers belong inside the contract, not bolted onto it afterwards.**

**Three, and never a fourth.** The spec is a contract built by two people, and a
phase that hands one over finished has decided alone everything it did not ask
about. Three is where the questions worth a round trip run out.

**Zero is a legitimate answer, and it is reported in one line.** When the code and
the proposal already settle everything, say so and write the spec. A quota
manufactures questions the code answers, which `AGENTS.md` forbids in as many words,
and a question whose only available answer is "yes, carry on" is a round trip charged
for a rubber stamp.

**One call, not three turns.** Three round trips to build one contract is the
ceremony this flow spends its length arguing against.

**The answers go into the spec**, under prior decisions, next to what each one
settled, attributed and dated. An answer that lives only in the conversation gets
asked again next session, and the second answer will not always match the first.

**Phase 2 is the only phase that asks this way.** Phase 5 stops for the order and
what waits on what; it does not open a second tranche of questions. The contract is
where an answer changes the most and costs the least to change, and asking twice
before the first line of code is how a three-stop flow starts growing back.

Do not ask what the code can tell you. Read it.

The trivial lane asks nothing: step 0 answered *no spec needed*, so there is no
contract to disagree about and nothing to ask about it.

## Step 5 — Check it before handing it over

Read the finished spec as an adversary and look for exactly these:

- a pillar with a heading and nothing real under it
- an outcome that no criterion tests
- a criterion with no named test
- two clauses that cannot both be true
- a scope boundary vague enough to argue with
- an answer from step 4 that never made it into the file

Fix them here. Every one of these is cheap now and expensive in phase 6.

## Where specs live

```
.agents/
├── specs/                        the living, consolidated specification
│   └── bundle-products/
│       └── spec.md               Governs: src/bundles/** src/pricing/*.ts
└── changes/                      proposals in flight
    └── add-relative-discounts/
        ├── proposal.md           the ticket, the intent, the non-goals
        ├── spec.md               Targets: bundle-products — the delta
        └── plan.md               the checklist
```

**One spec per capability, never per ticket.** This matters more than it looks. A
capability spec accumulates and stays true for years; a spec named after a ticket
is dead the day the ticket closes, and a directory of dead specs is worse than no
specs because people still read them.

So the first question is not "what shall I call this spec" but **which capability
is this change amending?** If none of them fit, that is the discovery: a new
capability, and it gets its own directory.

**Edits go in the change while the work is in flight.** Not into the capability
spec directly — the capability spec is what is true now, and the change is a
proposal until it lands. `Targets:` names the capability so a tool can tell that an
in-flight delta *is* that spec moving.

When it lands, phase 8 applies the delta onto the capability spec and the change
folder goes away. Exactly one current specification, always. See
`skills/record-work/`.

A change spanning several capabilities carries one delta per capability inside the
same change folder, each with its own `Targets:`. One change, several amendments —
not several changes racing each other.

### If the project already has a convention

Whatever exists wins. `spec-drift` (at `<skill-base>/../record-work/spec-drift`,
next to the sibling `record-work` skill) looks for consolidated specs under
`.agents/specs`, `specs`, `openspec`, `docs/specs`, `spec` — in that order — plus a
single-file `docs/SPEC.md`, and for changes under `.agents/changes`, `changes`,
`openspec/changes`.

A single-file spec is legitimate for a small project and the tooling supports it.
This repository is one.

**Never create a `changes/` directory in a project that has none** without also
doing the consolidation. A staging area nobody empties is just a second source of
truth.

## Output

Report which specs were written, where, and anything still unresolved that phase 5
onward will hit.

Then stop, and wait for the go-ahead before the plan.
