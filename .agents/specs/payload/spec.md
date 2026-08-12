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
- **a queue for ideas that arrive faster than they get built** — `libretto-queue`
  captures them one after another as proposals carrying a `Queued:` date, committing
  each on the current branch and creating no branch, no spec and no plan;
  `libretto-next` picks one, oldest first with the user free to choose another, creates
  the branch, removes the `Queued:` line and enters the flow at phase 2, because phase
  1's artifact already exists
- one skill per phase, each of which stops where its phase stops
- **the same flow with its three stops answered in advance** — `libretto-attacca`, the
  score's instruction to go on without pausing, running phases 1 to 8 and ending at a
  pushed branch with a request open on it. It writes everything the attended flow writes;
  what the invocation answers is the waiting. **What it cannot answer is a gate**, and a
  question no reading of the code settles becomes an assumption recorded in the spec, the
  report and the request rather than a prompt. **The request's description carries both
  halves** — what the invocation answered and what the run assumed — because without the
  first a reviewer cannot tell an agreed contract from an assumed one
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
- **the companions the flow calls by name ship vendored too** — `ponytail`,
  `ponytail-debt`, `caveman`, `caveman-commit` — so "how much gets built" and "how
  much gets said" work on a machine that installed nothing else. Only what the flow
  calls by name: the rest of both plugins stays upstream.
- drift detection that ships with the skill that uses it

**Every skill is self-sufficient once installed.** A skill that only works inside this
repository is a skill that works for nobody.

## Scope boundaries

**In:** skills, agents, commands, the drift-detection tooling, the vendored
third-party items and their attribution.

**Out:**

- **requiring ponytail or caveman.** They ship vendored, but shipped is not
  required: nothing fails without them, and `prune`/`uninstall` removes them like
  any other item.
- **the companions' always-on mode.** As plugins both can inject themselves into
  every session via hooks; this payload does not manage `settings.json` or hooks.
  The vendored skills activate when the flow invokes them and when their
  descriptions trigger. The upstream plugin remains the path to always-on, and the
  two copies coexist by namespace.
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
- **priorities, tags or reordering on the queue.** FIFO by `Queued:` date, and
  `libretto-next` letting the user pick a different one *is* the reordering mechanism. A
  priority field returns the day FIFO measurably hurts.
- **a command that edits or deletes a queued idea.** They are markdown files in a
  directory; a CRUD surface over three files is ceremony.
- **a separate queue file.** `.agents/changes/` already holds proposals; a `queue.md`
  beside it would be a second source of truth about what is queued.
- **draining the queue unattended.** `libretto-next` runs one idea and the flow's own
  stops apply, and `libretto-attacca` does not reopen it: one invocation is one piece of
  work. Batch execution is a different feature with different risks.
- **an unattended run merging, tagging, releasing, or labelling the request.** It ends at
  a request open for review. The bump is a reading of `.agents/specs/` rather than of the
  commits, `release:major` is asked-and-waited-for by standing rule, and a version number
  cannot be recalled once the proxy has cached it.
- **skipping, reordering or softening a gate**, and `--force`, `--no-verify` or anything
  else that buys a green result. Unattended removes waits, never checks.
- **a second door into the unattended mode.** One command, no flag on `libretto-flow`
  doing the same thing: two entry points drift, and the argument-shaped one is the one
  that gets typed by accident. For the same reason `libretto-attacca` describes no phase —
  it delegates to the same skills and restates none of them.
- **an unattended `libretto-next` or `libretto-review`**, and **a setting, profile or env
  var for the mode.** Consent is for one run, and consent that persists is consent nobody
  remembers giving.
- **the Go binary knowing about the queue.** It is payload, not delivery.

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

**Every stop is a native question — `AskUserQuestion`, never a sentence at the end of a
report** — and so is phase 1's choice between work already in flight and something new. The
recommended option first, saying what will actually run; the real alternatives; room to
answer differently.

**The argument for it lives once, in `skills/record-work/`**, written for the push on
2026-08-10 and never push-specific. Every other place states the rule and points there — a
reason copied into six files is six things to keep in sync, and the reviewer of the run that
landed this caught exactly that: the constraint saying *it lives once* shipped in the same
commit as five copies of it.

Phases 2 and 5 waited in prose until 2026-08-12, including on the run that specced this,
which is where the evidence came from. **Ceiling named:** nothing mechanical checks it. `scripts/check-payload`
cannot tell a native prompt from a paragraph. The replacement, the day it drifts, is a
search for `AskUserQuestion` across the skills that own a stop — a string match, and it
would have caught this one. Deliberately not built: a guard against a failure that has
happened once.

**A stop whose only available answer is "yes, carry on" is a round trip charged for a
rubber stamp**, and calling it a decision point does not make it one. Phase 1's wait
bought the user reading a paraphrase of their own sentence back. Phase 7's bought
permission to commit to a local branch nobody had seen — the cheapest place in the flow
to change your mind, guarded as though it were a deployment.

**Two exceptions in phase 1, and neither is ceremony:** work already in flight, where
continuing it or not is a choice about the user's priorities, and a missing or
unconfigured tracker, where nothing downstream exists. Both are the input failing to
arrive rather than a phase boundary asking to be blessed.

**And the table is what `libretto-attacca` is defined over.** The unattended mode is a
second reading of this same list rather than a feature of its own: the three stops are
answered by the invocation, phase 1's in-flight choice resolves to the oldest change with
open boxes, and everything not in the table is untouched. **A gate is not a stop and
cannot be answered** — a failing one still stops the run, two on one task still stop the
task, and a missing credential still stops the phase, because that is the input failing to
arrive and no invocation supplies it. A mode that answered a gate would not be unattended,
it would be unverified.

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

**Phase 2 asks up to three, in one call, before it writes the spec** — so the contract is
built by two people rather than handed over finished, and the answers sit inside it rather
than bolted on. Three is where the questions worth a round trip run out. **Zero is a
legitimate answer, reported in one line:** a quota manufactures questions the code already
answers, which is the rubber-stamp round trip removed from three other phases arriving
back through the door marked collaboration. Phase 5 stops for the order and opens no
second tranche, and **the stop count does not move** — the questions ride the stop phase 2
already has. The trivial lane asks nothing, because *no spec needed* means there is no
contract to disagree about.

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

**A queued idea is not a fourth source.** `find-work` owns the `Queued:` scan so that
`/libretto-status` can report the queue and `/libretto-next` can pick from it without
either walking the directory itself — the same single-owner rule, and the reviewer caught
`libretto-next` breaking it on the run that landed the queue. The scan sits outside the
Source 1 heading for a reason: nested under it, "run source 1 only" left
`/libretto-status` unable to tell whether the queue was in scope.

**Queued is reported, never blocking.** Home first exists so *started* work does not get
abandoned; an idea costs nothing to abandon because nothing has been built on it. Four
captured ideas standing between the user and a Jira task would make capture punitive, and
a queue that is expensive to add to is a queue nobody uses.

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

**A marked box that was never committed was never marked.** The plan ships in the same
commit as the task that closed it, exactly as the spec does. A `git add` scoped to the code
leaves it behind in the working tree, and the change that lands then deletes the folder with
every unrecorded mark still in it — measured on the run that added this line: six commits,
0/24 boxes, work finished. Phase 1 reads unchecked boxes to decide what is in flight, so
the failure is silent until the plan disappears.

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
- **Phase 8 returns to the base branch when the request opens**, checking it out and
  fast-forwarding it, on the yes path only. The user's call, 2026-08-12. A session starts
  wherever the last one left the working directory, so a flow parked on a merged feature
  branch hands the next phase 1 a stale base — and reading what is in flight off that base
  is phase 1's whole job. Measured on the run that added it: phase 1 offered the user a
  choice about a branch already merged and tagged `v0.6.1`, because local `main` was seven
  commits behind. Fetching the ref without leaving the branch was offered and declined —
  it fixes the ref, not the place the next session starts. **The feature branch is never
  deleted** (the request is open, not merged, and the local copy is the only one), nothing
  is ever merged (`--ff-only`, and a diverged base is reported rather than resolved), and a
  dirty tree stops rather than being carried across by `git checkout`. On the no path it
  does nothing: the branch is the only place the work exists.
- **Up to three questions at phase 2, none forced.** The user's call, 2026-08-12: *"para
  que el plan se cree entre los 2 y no solo tú"*. Always-three and questions-at-phase-5
  were both offered and declined.
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
- **capturing an idea and running it are two commands, not one overloaded flow** — the
  user's call, 2026-08-11: "le pasas una tarea… no tiene sentido pasarle la tarea y que
  haga otra". `/libretto-flow <task>` does *that* task, always; the queue drains through
  `/libretto-next`. A flow that silently substitutes different work for what it was handed
  is the surprise nobody wants.
- **FIFO by a `Queued:` line in the file, not by git archaeology.** The date survives a
  rebase and needs no `git log` walk, and its presence is also what marks the change as
  not-started — `/libretto-next` removes it at pickup. **Ceiling named:** no priorities;
  the replacement is a priority field the day someone reorders constantly.
- **Captured proposals commit on the current branch, and the branch waits for
  `/libretto-next`.** A branch per captured idea scatters the queue across N branches
  nobody can see from the base. This does not weaken "the branch exists before the first
  write": a queued proposal is not the change's work yet, and the change's first write is
  the pickup.
- **the unattended mode is its own command, not a flag** — the user's call, 2026-08-12,
  reversing the reading the change was first specced with. "No second flow for small work"
  above refuses a second *flow*; this is a second *door* onto the same one, and the
  precedent that governs it is `libretto-next`, which is its own command because an
  invocation that behaves differently must be unmistakable in the history rather than an
  argument that can be typed by accident. What the refused rule still binds is the shape:
  it delegates and describes nothing.
- **`attacca` because it is the instruction, not a metaphor for it** — what a score writes
  to mean *go on to the next movement without pausing*. The cost is opacity to a reader
  who does not read music, paid by one line of `description:`. `auto` was rejected for
  saying something is automatic without saying what still stops, `solo` for pointing at a
  prominent voice rather than an absent pause.
- **Under the unattended mode a question the flow cannot derive is assumed, recorded and
  carried on from** — the user's call, 2026-08-12. Stopping to ask makes the mode
  unreliable in exactly the case it exists for; stopping without asking trades a wait for
  a dead run. The mechanism is not new — the flow already turns an unsettled post-plan
  question into a finding carrying what was assumed and what changes if it is wrong — and
  this moves that rule to the front. **Ceiling named:** an assumption is only as visible as
  the report and the request carrying it, so a run nobody reads has bought silence rather
  than speed. The replacement that day is refusing to open the request when an assumption
  was made, never a prompt mid-run.
- **Push and the request are answered by the invocation, never assumed.** "Never push
  unasked" is intact: the asking happened at the prompt. The consent covers that branch and
  that request and nothing past it.
- **ponytail and caveman are vendored, reversing THIRD-PARTY.md's original
  not-vendored entry — the user's explicit call, 2026-08-10.** The old rationale (a
  second copy of something the user may already have chosen a version of) assumed a
  user who has versions of things; the target is a machine with nothing on it, where
  "if installed" was a conditional that never came true. The collision half was
  already answered by naming: plugins namespace, vendored copies do not, both
  coexist. "The installer prints the install commands" was specced first and
  discarded — it left the fresh user one manual step from a flow that works as
  written. The selection rule is ponytail's own first rung: only what the flow calls
  by name. Pinned versions and update procedure live in THIRD-PARTY.md; drift from
  upstream is the accepted cost, and `diff` against a fresh clone at the pinned
  commit is the check when doubt arises.

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
- [x] `libretto-queue` and `libretto-next` — capture the queue, drain it one at a time,
      with `find-work` owning the scan both read
- [x] every stop asked with `AskUserQuestion` — phases 2, 5 and 8, plus phase 1's in-flight
      choice, with the rule stated in each skill that owns one
- [x] `libretto-attacca` — the three stops answered by the invocation, the classification
      in one file, and the five stop-owning skills each stating what happens to their own
- [x] vendored delegates with attribution
- [x] ponytail, ponytail-debt, caveman and caveman-commit vendored, callers' prose
      reconciled, THIRD-PARTY.md and docs recording the reversal
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
- **the whole tree is traceable, not only what is staged.** Both existing modes are scoped
  to a commit — one to staged paths, one to citations that already exist — so **code that
  was never staged while a spec existed is invisible to both, forever.** In a repository
  whose specs arrived after its code that is most of the repository, and the contract can
  pass every check while governing a third of the tree. `--trace` reads the tracked tree
  and names orphan code, a `Governs:` glob matching nothing, and a criterion under
  *Verification criteria* with no `Proof:` beneath it. **It exits 0 unconditionally**: the
  first run here found sixteen orphans and sixteen unproven bullets, which is the honest
  state rather than a regression, and a report that fails the build the day it lands is a
  report somebody deletes.
  Proof: skills/record-work/spec-drift --self-test
- **a criterion is a column-0 bullet under the criteria heading, and nothing else.** Its
  `Proof:` survives the continuation lines beneath it and ends at the next such bullet or
  the next heading. Reading every bullet in the document instead buries the real findings
  under the Outcomes section, which has no `Proof:` by design.
  Proof: skills/record-work/spec-drift --self-test
- **the contract is reviewed before a plan is built on it.** Every other review in the flow
  reads code and all of them run after the code exists, so a criterion no run could fail is
  agreed at the phase-3 stop, becomes tasks, becomes code, and is then measured against a
  sentence nobody could have failed — the work passes and the promise was never real.
  `review-spec` runs in the 3→5 seam over ambiguity that forks the work, criteria that
  cannot fail, `Governs:` boundaries describing nothing, and a criterion another capability
  contradicts. **Before the stop, not after**, because the findings have to be on the table
  while the agreement is still being made. It reports and never rewrites the spec — a
  reviewer editing the contract it reviews is the author twice over. Skipped entirely on the
  trivial lane, which has no contract.
  Proof: scripts/check-payload
- **the mechanical half is delegated, not restated.** `review-spec` runs `--trace` first and
  reads the spec for what a script cannot see: whether a `Proof:` names a test that could
  actually fail. A citation pointing at a test that only asserts the code ran is untestable
  with a citation attached, which reads as proven and is therefore worse than none.
  Proof: scripts/check-payload
- **the status command delegates rather than restating the scan**, and every referenced
  skill survived the rename
  Proof: scripts/check-payload
- **the queue commands delegate the scan too**, and reference only skills that exist
  Proof: scripts/check-payload
- **a vendored skill's namespaced reference is either shipped or answered.** `writing-plans`
  opens every plan it writes with a header demanding `superpowers:subagent-driven-development`
  and offers `superpowers:executing-plans`; `test-driven-development` cites
  `superpowers:writing-skills`. **The payload ships none of the three, and neither does a
  machine that installed only this** — so the flow's own phase-5 output has instructed its
  reader to invoke a skill that does not exist since the day the copies landed. Nothing read
  the namespace, so nothing noticed. The fix is never an edit to the vendored copy, because
  `THIRD-PARTY.md` promises it stays comparable with upstream and a divergence needs exactly
  one home: the Libretto skill that delegates to it states the override — `write-plan` for
  the two execution skills and the worktree one, `build-and-check` for the test-prose one.
  The check reads refs only inside vendored directories, listed from `THIRD-PARTY.md` rather
  than hardcoded, and requires each to resolve to a shipped skill or be named by a skill that
  is not the vendored one. **A ref inside a Libretto skill is the answer, never the fault** —
  scoping it the other way made the override itself the failure.
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

**The two 2026-08-12 amendments are prose in a skill, and no criterion above cites a test
for them** — a skill is a prompt, checked by running it, and a Go test named for "phase 8
checks out the base branch" would be the fabrication `--anchors` exists to catch. What was
observed on the run that added them: phase 2 asked three questions in one call and all
three answers reached this spec under prior decisions. Still claims: a yes at phase 8
ending with the tree on a current base and the feature branch intact, a no ending on the
branch unchanged, and a phase 2 with nothing to ask saying so in one line.

**The unattended mode is prose and none of it has run either**, and it is the one where
that gap costs the most, because every claim is about what happens when nobody is
watching. Claims, not facts: that a run passes the spec stop and the plan stop with both
artifacts on disk; that it ends at a pushed branch with a request whose description names
what the invocation answered; that **a failing gate stops it on its branch with no request
opened**; that an absent `gh` stops it with the install line and nothing else; and that a
question it could not derive appears as a marked assumption in the spec, the report and
the request, and nowhere as a prompt. The third of those is the one to run first — it is
the case where the mode must look exactly like the attended flow, and the only one whose
failure is silent.

**The queue is prose and none of it has run.** Claims, not facts: capturing two ideas
leaves two committed proposals and no branch; `/libretto-next` offers the oldest first and
enters phase 2 on a fresh branch; `/libretto-flow` handed a key never mentions the queue;
`/libretto-status` shows the queue as its own section. What *was* observed on the run that
landed it is the reviewer catching three defects in the prose — a duplicated scan, a commit
step stated in a never-does bullet instead of instructed, and a queued idea defined two
ways in adjacent lines.
