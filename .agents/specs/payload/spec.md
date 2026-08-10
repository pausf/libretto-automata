# Payload

Governs: skills/** commands/** agents/** scripts/**

The reason the project exists. Everything else is delivery.

The flow itself — its eight phases and its reasoning — is [docs/FLOW.md](../../../docs/FLOW.md).
This spec is the contract the payload's *artifacts* have to satisfy.

## Outcomes

Installing this repository gives a working flow on a machine that has nothing else.

- one command, `libretto-flow`, that routes and never implements
- `libretto-status`, read-only, reporting what is in flight
- `libretto-review`, which reviews a forge PR/MR in a workspace that restores itself
  — its contract is the `review-project` capability spec
- one skill per phase, each of which stops where its phase stops
- **the finished work reviewed by someone who did not write it, then repaired** — in the
  seam between build and present, `review-work` launches one fresh `work-reviewer`
  subagent with none of the session's context; it re-runs every proof the change
  touches and returns findings, which phase 7 carries attributed and unedited **with the
  repair beside each one**
- three standing rules: **commit** and **evidence** hold at every phase; **ask** holds at
  phases 1, 2 and 5 and nowhere after
- **ceremony proportional to the change** — two stops for a change with a contract, none
  for a change too small to have one, and phase 8's question in both
- the vendored delegates the thin skills depend on, so a thin skill is never a broken
  one
- drift detection that ships with the skill that uses it

**Every skill is self-sufficient once installed.** A skill that only works inside this
repository is a skill that works for nobody.

## Scope boundaries

**In:** skills, agents, commands, the drift-detection tooling, the vendored
third-party items and their attribution.

**Out:**

- **requiring ponytail or caveman.** Called when present, never required. `doctor`
  reports them; nothing fails without them.
- creating a `changes/` directory in a project that has none. A staging area nobody
  empties is a second source of truth.
- accepting a secret. No skill asks for a token, and any skill that is handed one says
  it is now exposed and continues without it.
- pushing, or committing from a sub-agent.
- **a flag, profile or setting for how big a change is.** It is knowable from the change.
  Asking the user to declare it invents a second source of truth about a fact already on
  disk, and the two will disagree exactly when it matters.
- **a second flow for small work.** One flow with a proportionate gear, not two to keep
  in sync.
- **a review that blocks or opines.** Style and taste are not findings; a finding cites a
  pillar or a proof, and nothing a reviewer says stops the flow.
- **the reviewer writing.** `work-reviewer` keeps the grant it has — `Read`, `Grep`,
  `Glob`, `Bash`, `Skill` — and gains nothing that edits. `Skill` is there because its
  first instruction is to invoke `evidence`, `Bash` because running the proofs is its
  job. The repair happens in `review-work`, on the other side of the seam: an agent that
  fixes what it found is grading its own repair.
- **a fix loop.** One pass. Two failures on one finding stops that finding, per
  `skills/evidence/`, and it reaches phase 7 as found-and-not-fixed.
- **asking about a finding**, including one that looks like a product decision. It is
  reported to phase 7, never turned back into a question.

### A stop is a place where the user changes something

That is the whole test, and it is what the count is derived from.

| After | Stops | What is being changed |
|---|---|---|
| 1 · find-work | no | — the reading is stated, and the spec is where it gets corrected |
| 2–3 · write-spec | **yes** | the contract |
| 5 · write-plan | **yes** | the order, and what waits on what |
| 6 · build-and-check | no | — |
| 6→7 · review-work | no | — the seam fixes what it finds |
| 7 · present-work | no | — |
| 8 · record-work | **yes, last** | whether the world sees it |

**A stop whose only available answer is "yes, carry on" is a round trip charged for a
rubber stamp**, and calling it a decision point does not make it one. Phase 1's wait
bought the user reading a paraphrase of their own sentence back. Phase 7's bought
permission to commit to a local branch nobody had seen — the cheapest place in the flow
to change your mind, guarded as though it were a deployment.

**Two exceptions in phase 1, and neither is ceremony:** work already in flight, where
continuing it or not is a choice about the user's priorities, and a missing or
unconfigured tracker, where nothing downstream exists. Both are the input failing to
arrive rather than a phase boundary asking to be blessed.

When phase 2 answers **no spec needed** both remaining stops collapse with it — there is
no plan either — and phases 6, 7 and 8 run in one turn. **Phase 8's question survives
every collapse**: pushing is the user's decision, not ceremony.

**What collapses is the wait, never the saying.** Phase 7 still reports what was done, its
evidence, and what was left out with the condition that brings it back. A phase that skips
the report because the change looked small is how the one omission that mattered goes
unmentioned.

The cost of getting this wrong is measured, not theoretical: a session spent four round
trips updating two documentation files, and every one of the four was mandated by this
payload. A flow that charges a typo the price of a feature gets routed around — for typos
first, then for small features, until what is left is a ritual reserved for work important
enough to deserve it.

### Asking is bounded to before the plan

Phases 1, 2 and 5, where nothing has been built on the answer yet — which is why the two
stops sit exactly there. **After the plan, an unsettled question becomes a finding**: it
reaches the phase 7 report with what was assumed in the meantime and what changes if the
assumption is wrong, and the user meets it at phase 8 with everything else.

This reverses an earlier promise that asking held at every phase, as often as needed. That
was reasonable in every individual case and that is precisely the problem — it returns the
flow to five, six, nine stops one defensible exception at a time. The bound is the
load-bearing half of the count above; without it the table is decorative.

### A phase is invoked even when it has nothing to do

**A decision belongs to the phase that owns it.** `write-spec` decides whether a spec is
needed; `build-and-check` decides where the work lands and how much to check. An
orchestrator that makes those calls itself and proceeds has not saved a step — it has
moved a decision somewhere nobody can audit.

So every phase the work reaches is invoked, **including the ones whose answer is "nothing
here"**, and the declining is reported in one line. Invoking is not gating: a phase that
declines adds no wait.

The cost is one line. What it buys is the difference between a skip and an omission, which
from outside are identical — a session made exactly the right call about a small change,
never invoked phases 2 or 6 to make it, and the user asked why the flow had not run. It
had, in substance. Nothing said so.

### Work is found, not fetched

**The flow does not begin at a tracker.** Phase 1 asks three sources, in this order:

| Source | When |
|---|---|
| a change already in flight | unchecked boxes in `.agents/changes/*/plan.md` |
| a tracker key or URL | one was given |
| what the user said | anything else — the request *is* the input |

Home first, and the order carries the reason: starting something new while a change sits
half-finished is how the half-finished thing gets abandoned, and the cost is not the lost
work but a `.agents/changes/` directory nobody trusts.

**A request in conversation is a legitimate input.** Every change in this repository so
far arrived that way, and the phase used to treat it as a fallback in a table. With no key
the change is named from the request, verb-led, and `proposal.md` records `Tracker: none`
plus what was asked in the words it was asked in — a fake key would imply a tracker
somebody could consult.

Two other implementations agree that the tracker is not the door. gentle-ai mentions none
at all in its propose or explore phases; CodelyTV/agent-harness takes the task as an
argument and offers GitHub Issues as an optional variant.

**Reporting and choosing are one skill.** `/libretto-status` invokes phase 1 in reporting
mode rather than scanning the same directories its own way. Two answers to "what work
exists" is one too many, and the one that disagrees is always the one nobody is reading.

**Phase 1 produces an artifact, not only a report.** A source-3 request leaves
`.agents/changes/<name>/proposal.md` on disk — `Tracker: none`, and the ask in the words
it arrived in — **before** the phase reports, and the reading is confirmed after that file
exists rather than instead of it. A reading agreed in conversation and recorded nowhere is
re-derived next session and does not come out the same.

The rule was already written, in prose, under the source-3 heading. It read as a
description of what a proposal contains rather than as a step, and it was skipped: the ask
was reported, confirmation was asked for, and no file existed until the user pointed it
out. **A step stated as prose is a step that gets read and not done**, so phase 1 now
states its outputs as a table of what must exist before it may report.

Writing that file means committing it, and committing means a branch — so **phase 1
branches when it writes the proposal**, and phase 6's step 0 *ensures* a branch instead of
creating one. A step that assumes it is first makes a second branch and splits one change
across two.

**A branch is work in flight too.** Scanning `.agents/changes/*/plan.md` cannot see the
trivial lane: a change that needed no spec has no plan to scan, so it exists only as
commits on a branch. The lane's first real run produced exactly that, and the next phase 1
reported an empty house. Phase 1 also reads `git branch --no-merged` and the forge's open
requests, and names the state worth naming: **unpushed and un-requested** is work nobody
but this machine has.

**Landed changes are deleted, not archived.** gentle-ai moves them to `changes/archive/`;
git history is the archive here, and a directory nobody reopens is growth. A decision,
not an oversight.

## Constraints

**A skill may only invoke what gets installed.** `libretto install` links `skills/`,
`agents/` and `commands/` — nothing else. A skill telling the user to run `scripts/foo`
is broken for everyone who installed it, because `scripts/` never reaches their
machine. **A tool a skill needs ships inside the skill's own directory.**

For the same reason no skill names a path from this repository as though it existed in
the user's project. `docs/FLOW.md` is not somewhere they have.

**Frontmatter `name:` must equal the directory or filename.** A mismatch means the
caller reads one name off disk and invokes another, and the failure is a sub-agent that
silently never runs.

**Every file has exactly one author.** The brief is written once and read many times;
each spec delta has one writer; the plan has one writer, the orchestrator. This is what
makes concurrency safe without locks.

**Sub-agents start with no rules loaded.** Every launch carries them explicitly, or the
default takes over — and the default is to guess plausibly.

**An agent that ships is named by the skill that launches it.** Phase 3's fan-out runs
`spec-writer`, named in `write-spec`. An agent nothing references is dead payload —
`spec-writer` sat unreferenced from its landing until 2026-08-10, and no gate noticed.

**Vendored items are copied unmodified**, with their licence and version recorded. A
change needed for this flow goes in the skill that calls them, never into the copy, so
the copy stays comparable with upstream.

## Prior decisions

- The tracker is read through its CLI, never MCP, never the REST API.
- The API token lives in the OS keyring, put there by `jira init`, run by the user in
  their own shell. It never enters a conversation.
- Specs are per **capability**, never per ticket: a capability spec accumulates and
  stays true, a ticket spec is dead the day the ticket closes.
- Deltas live in `.agents/changes/<change>/` and are applied onto the capability spec in
  the commit that lands the change, which then deletes the change folder.
- Drift detection **warns and never blocks**. A check that stops a commit in someone
  else's project is a check that gets deleted.
- Phase 2 may decide no spec is needed. Skipping the phase is a legitimate outcome of
  it, and the "no" collapses phase 7's gate with it.
- Phase 6 does not fan out. Parallel implementation needs isolation, a serial merge
  queue and a conflict protocol; without those three, concurrency manufactures races.
- **The branch exists before the first write** — not before the first commit.
  `git checkout -b` carries uncommitted work, so editing on the base branch and branching
  at commit time succeeds until the base has moved or touches one of your files. Whichever
  phase writes first creates it: phase 1 when it writes a proposal, phase 6's step 0
  otherwise, which *ensures* rather than creates. Phase 8 keeps the same check as a
  backstop, names the writing phase as its owner,
  and reports rather than silently fixes: a backstop that covers for the rule it backs up
  is how the rule stops being followed.
- **Push and the pull request are one question**, and it is asked with `AskUserQuestion`.
  Asked separately they bought a second round trip and no safety — a pushed branch with
  no request opened is a state almost nobody wants, and whoever wants it says so in the
  same breath. Native rather than prose because it is the last question in the flow and
  usually the only one after the plan: a question written as a sentence is a paragraph
  the reader skims, and the flow then waits on an answer to something that read as a
  summary. Observed 2026-08-10, on the run that landed this.
- **The forge is derived, never assumed:** `git remote get-url origin`, `github.com` →
  `gh`, `gitlab` → `glab`. No remote means no question. **Ceiling named:** a substring
  test on one URL, which does not survive a self-hosted forge on a neutral domain, or
  Gitea and Forgejo. The replacement that day is an explicit setting read from the
  repository, not a longer list of guesses in a prompt.
- **A missing forge CLI stops**, with its install line and nothing else — the shape phase
  1 already uses for `jira`. Never a hand-built API call using a token found in the
  environment: that turns a stop into an exposure. The other forge's CLI is not a
  fallback.
- **Whether a change needs a spec stays a judgment, with no heuristic behind it** — not a
  diff-size threshold, not a file count. The real test is whether two people could
  reasonably disagree about what "done" means. A number would be wrong in both directions
  and trusted anyway.
- **The reviewer never blocks, and the seam repairs rather than asks.** The first half
  was decided when the seam was added; the second replaced it on 2026-08-10, when
  "acting on a finding is the user's call and a new pass through phase 6" turned out to
  be a round trip charged for a rubber stamp. A finding cites a pillar or a proof *by
  contract*, so it is a defect against something the user agreed to at phase 2 — there
  is no version of "yes, leave it broken" worth a stop.
  **Ceiling named:** one fix pass, no re-review, so a fix that introduces a new defect is
  caught by the proofs or not at all. The replacement, the day that bites, is one bounded
  second look at the fix diff — never a loop.
  `spec-drift` keeps the older standing, warn-never-block, and the two are consistent:
  drift asks whether a contract still describes the code, which is the user's call. A
  finding is a breach of a contract already settled.
- **The review seam is not a numbered phase.** Independence comes from the fresh
  context, not from renumbering everything that says eight.
- **The reviewer re-runs the proofs itself.** Phase 6 having reported them green is
  not evidence in a review — trusting the builder's report is the failure the seam
  exists to remove.
- **No spec, no review.** The trivial lane collapsed the ceremony because there was
  no contract to disagree with; a reviewer with no contract has nothing to check. One
  line, no wait.
- **One reviewer, not a panel.** Two runs of the same model over the same diff are
  correlation dressed as corroboration. The day reviews need lenses, that is a spec
  change, not a quiet doubling.

## Task breakdown

- [x] `find-work` — phase 1, three sources
- [x] `libretto-status` — read-only, what is in flight
- [x] `write-spec` — phases 0, 2, 3, including the fan-out and its brief
- [x] `write-plan` — phase 5
- [x] `build-and-check` — phase 6
- [x] `present-work` — phase 7
- [x] `record-work` — phase 8, including consolidation
- [x] `evidence` — the standing rules
- [x] `libretto-flow` — the routing command
- [x] vendored delegates with attribution
- [x] `spec-drift` and `check-payload`
- [x] **the flow run end to end against a real task.** `right-size-the-flow` — phase 1
      found it in flight, 2 wrote the delta, 5 the plan, 6 closed six tasks with a commit
      each, 7 reported, 8 consolidated. The forge stop was exercised for real, not
      simulated: `gh` was absent, phase 8 stopped with the install line.
- [ ] **the failure paths.** Still written and never run: an unconfigured tracker, a board
      URL where a key was expected, the trivial lane actually collapsing on a one-line
      change, and a sub-agent hitting a question it must not answer itself.
- [x] an independent verifier: check the implementation against the spec's criteria,
      never run by whoever wrote the code — `review-work` and `work-reviewer`, in the
      seam between phases 6 and 7

## Verification criteria

- frontmatter parses, and `name:` matches the directory or filename
  Proof: scripts/check-payload
- no stray file sits where the linker would install it as an item
  Proof: scripts/check-payload
- every referenced skill exists
  Proof: scripts/check-payload
- **no skill invokes a path that does not get installed**
  Proof: scripts/check-payload
- **no skill hardcodes the install layout.** A `~/.claude/` path is only true under
  `install --global`; a skill's tools resolve from its own base directory, which every
  invocation announces — `record-work` reaches `spec-drift` as its sibling, `write-spec`
  hops to `../record-work/`. Both layouts keep skills side by side.
  Proof: scripts/check-payload
- glob matching, capability derivation and citation extraction behave
  Proof: skills/record-work/spec-drift --self-test
- **an anchor inside a fenced code block is an illustration, not a declaration.** Any
  document explaining the convention shows it in a fence; reading that literally made
  the index page claim to govern the linking package. Knowing the difference is the
  extractor's job, not the author's job to avoid illustrating.
  Proof: skills/record-work/spec-drift --self-test
- **a glob that matches real directories is still a pattern.** Unquoted, the shell
  path-expands it before it is ever compared, so every glob matching real files
  silently stopped working and the check reported drift nowhere. `set -f` off, and a
  false negative in a checker is worse than no checker.
  Proof: skills/record-work/spec-drift --self-test
- **an invented test name is rejected rather than accepted**
  Proof: skills/record-work/spec-drift --self-test
- **spec-drift without rg refuses loudly.** Every question it asks goes through rg;
  missing, each match came back empty and default mode exited 0 having checked
  nothing. Now: exit 2 on stderr naming the tool, distinct from `--anchors`' failure
  exit 1. Drift findings keep exiting 0 — warn-never-block is about findings, not
  about being unable to look.
  Proof: skills/record-work/spec-drift --self-test
- every `Proof:` citation in every spec resolves, file and test name
  Proof: skills/record-work/spec-drift --anchors
- **the status command delegates rather than restating the scan**, and every referenced
  skill survived the rename
  Proof: scripts/check-payload

**What none of this verifies is behaviour.** A skill is a prompt, and a prompt is
checked by running it. The static checks above catch what silently degrades one — a
broken reference, an unreachable tool, frontmatter a host cannot parse — and nothing
more.

The failure paths are the ones worth exercising first, because they are the ones written
and never run: an unconfigured tracker, a board URL where a key was expected, a trivial
task that should skip phase 2, and a sub-agent that hits a question it must not answer
itself.

**What the first real run observed**, stated as observations rather than citations, because
a criterion citing a test that cannot exist is the fabrication the anchor exists to
prevent:

- phase 6 created the branch before the first edit; nothing landed on `main`
- the missing-CLI stop fired for real — `gh` absent, install line given, no workaround
  offered and no token asked for
- a phase declining was reported in one line, and cost no wait
- the four stops of the day still happened, because this change had a spec — the gear was
  proportionate, not removed. Left as recorded: it is what that run saw, and the count
  became three on 2026-08-10
- `prune` removed exactly the one stale link it planned to and left the other thirteen,
  the first exercise of that path outside a temporary directory

**What the reviewer's first two runs observed**, 2026-08-07:

- a fresh reviewer with no session context reviewed the change that built it — ran
  both cited proofs itself, verified every outcome against the diff, and returned a
  real finding: phase 7's own skill never carried the verdict, only `libretto-flow`
  did. Fixed in the same change. Caveat: the agent type registers at session start,
  so both runs used a fresh general-purpose subagent instructed to read and obey
  `agents/work-reviewer.md` first — same zero context, indirect binding. A fresh
  session exercises the registration itself.
- in a fixture repository, a criterion cited a test the committed code failed — the
  reviewer caught it by running it (`CarryDigits(45,55) = 1, want 2`, exit 1), named
  the unimplemented outcome clause, and never consulted the builder's claim.

Still unobserved, and therefore still claims rather than facts: the collapsed lane on a
change that needs no spec, the review seam's one-line decline on that same lane, and
every remaining failure path above.
