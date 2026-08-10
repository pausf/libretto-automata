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
citing a pillar or a proof, never taste. **It reports and never blocks**: the same
standing as `spec-drift`, because the stops in this flow exist so the user can say
no, not so a machine can. Acting on a finding is a new pass through phase 6.

No spec, no review: the trivial lane collapsed the ceremony because there was no
contract to disagree with, and a reviewer with no contract has nothing to check. One
line, no wait.

Not a numbered phase, on purpose — independence comes from the fresh context, not
from renumbering everything that says eight.

## 7 · Present

Show what was done. Then stop.

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
gets uninstalled, and an uninstalled gate catches nothing. Whoever wants it to block
can wire it into a `pre-commit` hook or CI; that choice belongs to them.

Being honest about the level this reaches: a warning is Spec-First discipline with a
reminder, not the automatic barrier the Spec-Anchored definition asks for. The barrier
would be a hook, and `settings.json` is out of scope for this project by an earlier
decision. The gap is deliberate and worth naming rather than papering over.

The push is asked at the very end — yes or no — and never assumed.

## Three rules, not phases

**Ask** (4) and **commit** (8) are written above as phases because that is where
they first appear, but neither happens once. Asking applies at every phase, as
often as needed. Committing happens per task throughout phase 6, not after phase
7. A flow that treats them as phases will do them once and consider them done.

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

**Where does the artifact get looked at?** This repository's first palette
satisfied its spec and was unreadable — 1.4:1 on borders. What caught it was
`docs/preview.py`, a throwaway that rendered the panel in a terminal before
`internal/ui/theme.go` existed, and then a WCAG measurement of that render. The
requirements for fluid width, centring and the single-colour menu all came from
looking at it, not from the ticket. Whether that becomes a phase of its own, or a
rule inside 2 and 6, is undecided. The reviewer in the 6→7 seam does not answer it:
it reads specs, diffs and test output, not pixels.

**Settled since:** who keeps the spec true. Phase 8 commits the spec alongside the
code that taught it — see `skills/record-work/`. The three divergences in this
repository's own `SPEC.md`, accumulated over one phase of work, are what the
alternative looks like.
