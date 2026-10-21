# Payload

Governs: skills/** commands/** agents/** scripts/** THIRD-PARTY.md licenses/**

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
- **a plan that holds the how, and a checklist that holds the state — two files, not
  one.** Phase 5 produced a checklist and called it a plan, so the technical reasoning
  behind a change was never written anywhere: the approach chosen, the alternatives it
  beat, the risks and how it gets validated all died with the session that had them. That
  was reported from outside — specs reading like plans, plans reading like task lists.
  `plan.md` is now the approach and holds no state; `tasks.md` is the checklist and holds
  nothing else. **The alternatives table is the pillar that pays**: a diff shows what was
  built, and nothing in a repository shows what was not built and why
- **the contract and the plan are co-authored with the user and drafted by agents that
  never saw the conversation.** Phase 2 interviews one question at a time and phase 5
  forks on the approach; every answer lands verbatim in the change's `decisions.md`, and
  the documents are then drafted from that log — the spec by `spec-writer` (single case
  included, `[NEEDS CLARIFICATION]` markers where the inputs run out), the plan by
  `plan-writer` (read-only, gaps in its return). The user's words reach the artifacts
  unparaphrased, and the drafting session cannot write down what was merely decided —
  the same property the task cut already had, arriving two phases earlier
- **the checklist is cut by an agent that never saw the design.** One fresh `task-cutter`,
  given the spec and the plan by path and nothing else, in the seam between phases 5 and 6.
  The session that argued its way to an approach cannot distinguish what it wrote down
  from what it merely decided, and cuts boxes against the argument — which the next fresh
  session cannot read. The cutter returns the checklist **and what those two documents
  failed to answer**, and the second half is the only check this flow has that a plan says
  enough to be built from. It writes no file: one writer is still the orchestrator
- **a finished change that never landed is reported, not walked past.** Phase 1 runs two
  `rg -c` scans over every change's checklist and `rg -c` is silent for a file that did
  not match, so a change present in the closed scan and absent from the open one has every
  box closed and its folder still on disk — its landing did not finish. That signal sat in
  the output unnamed while two changes shipped half-landed and phase 1 walked over the
  state twice. It is written as a table of three outcomes rather than a paragraph, so the
  empty open scan is a case with a name instead of an absence somebody has to notice, and
  `/libretto-status` delegates to the same scan rather than re-deriving it. **A report,
  never a gate**: `--retired` fires on a deletion, and the failure here is that no deletion
  happened — there is no commit to refuse
- **a plan cannot be deleted taking its reasoning with it.** The landing commit removes
  the change folder, and with it the alternatives table whose whole argument is that a
  diff shows what was built and nothing shows what was not built and why. `spec-drift
  --retired` — inside `--anchors` — fails that commit unless a capability spec's *Prior
  decisions* section moved in it. **That section, never the file**: the delta lands in the
  same commit by definition, so "a spec was edited" is green on every landing and measures
  nothing. The escape is a declaration in the plan being deleted, `Durable decisions:
  none`, and not a flag — a flag is typed by whoever wants the commit through, a line in
  the plan is written by the person who knew. **What no mechanism can stop is that line
  becoming a reflex**, and `review-work` is what reads a full alternatives table sitting
  beside a `none`
- **a criterion that can be failed, not merely read.** Verification criteria are written
  in one of the five EARS patterns, and `spec-drift --ears` — inside `--anchors`, so the
  gate count stays at six — fails a change delta whose criterion carries no `shall`. A
  criterion in prose can only be interpreted, and one nobody can fail is treated as
  satisfied by default: the same failure the `Proof:` anchor exists to prevent, arriving
  one line higher up. **Hard on deltas, a warning on capability specs**, because 545
  criteria predate the syntax and rewriting them in one unreviewable diff is 545 chances
  to change a promise that works today. A capability migrates when a delta lands on it.
  The gate checks the marker and never the sentence — *"the system shall behave
  correctly"* passes it, and `review-spec` is what asks whether the response is concrete
- **the stop moved rather than multiplied.** Splitting the plan from the checklist invited
  a third stop — one for the approach, one for the order — so phase 5 runs through and the
  5→6 seam stops for both at once, with the cutter's gaps on the table. The count stays at
  three, two of them inside the work
- **a checklist whose boxes are cut to stand alone, not ordered layers.** Phase 5 required a plan
  to be ordered and derived; it never required a box to be worth closing on its own, so the
  checklist inherited phase 2's cut along capabilities. `write-tasks` now states that a box is one
  end-to-end change — the code and the proof that closes it — and that two boxes where the
  first only makes sense once the second lands are **one badly cut box**, to be merged rather
  than ordered. A box that cannot be cut that way is a finding about the spec, routed back
  through the existing derived-not-invented rule. **The capability spec is not part of a
  box**: a delta still lands once, in the final commit. What it buys is `libretto loop`,
  which runs one fresh session per open box and has no way to detect a tree left broken
  between them
- **the same flow with its three stops answered in advance** — `libretto-attacca`, the
  score's instruction to go on without pausing, running phases 1 to 8 and ending at a
  pushed branch with a request open on it. It writes everything the attended flow writes;
  what the invocation answers is the waiting. **What it cannot answer is a gate**, and a
  question no reading of the code settles becomes an assumption recorded in the spec, the
  report and the request rather than a prompt. **The request's description carries three
  things** — what the invocation answered, what the run assumed, and the bump a person
  chose — because without the first a reviewer cannot tell an agreed contract from an
  assumed one
- **one question an unattended run does ask, and it is the last thing that happens** — the
  `release:` bump, after the report, the request and the return to the base branch. It is
  not a stop: nothing is downstream of it and the work is complete and reviewable whether
  it is answered or not. **The run never decides the bump; it asks and types the answer**,
  which is the split `AGENTS.md` already draws — the reading is the user's, the typing is
  not. `release:major` is selectable and never recommended, there is no default when it
  goes unanswered, and a repository that defines none of the three labels is not asked at
  all. It exists because a designed refusal nobody predicts reads as breakage: the red
  `require-release-label` check *is* the bump question, and it arrived once as a broken
  pipeline
- **the finished work reviewed by someone who did not write it, then repaired** — in the
  seam between build and present, `review-work` launches one fresh `work-reviewer`
  subagent with none of the session's context; it re-runs every proof the change
  touches and returns findings, which phase 7 carries attributed and unedited **with the
  repair beside each one**
- three standing rules: **commit** and **evidence** hold at every phase; **ask** holds at
  phases 1, 2 and the 5→6 seam, and nowhere after
- **ceremony proportional to the change** — two stops for a change with a contract, none
  for a change too small to have one, and phase 8's question in both
- the vendored delegates the thin skills depend on, so a thin skill is never a broken
  one
- **the companions the flow calls by name ship vendored too** — `ponytail`,
  `ponytail-debt`, `caveman`, `caveman-commit` — so "how much gets built" and "how
  much gets said" work on a machine that installed nothing else. Only what the flow
  calls by name: the rest of both plugins stays upstream.
- drift detection that ships with the skill that uses it
- **the landing verified before the commit that performs it.** `record-work`'s landing
  step, where `libretto` is on PATH, runs `libretto land` once the landing commit is
  staged — after the delta application, the retirement and the folder deletion are in
  the index, and after `libretto wiki` so the refreshed index rides the same diff. A
  non-zero exit names the missing part; the fix is made and the command re-run, never
  committed past. Where the binary is absent — or too old to know `land`, failing with
  *unknown command* and naming no part — the landing is said to be unverified and the
  step continues: a missing convenience never blocks a landing. The clause names the
  command and reads its exit status; what `land` checks — flags, parts, discovery,
  output — is the `cli` capability's contract, not this one's
- **the flow learns from its own corrections.** `evidence` captures every user correction
  into `.agents/lessons.md` while the flow runs — a correction is work already done being
  wrong; a changed ask is new work, not a lesson — as append-only entries under the
  countable header `## <date> · <change> · <phase>`, the user's words kept verbatim.
  `libretto-retro` routes to `retro`, which spends the open entries: **project knowledge**
  is recorded in the project's own contract, a **flow defect** becomes an exact diff
  proposed against the payload skill and never applied, and a **one-off** is marked as
  such. Spent entries gain a `Resolved:` line and are never edited otherwise.

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
- **an unattended run merging, tagging or releasing the request, or labelling it with a
  bump the user did not choose.** It ends at a request open for review. The bump is a
  reading of `.agents/specs/` rather than of the commits, `release:major` is
  asked-and-waited-for by standing rule, and a version number cannot be recalled once the
  proxy has cached it — so the reading stays the user's, always. **What the run may do is
  type an answer it was given**, once, at the very end. That is one word of this non-goal
  reversed and no more: merging, tagging and releasing stay absolute, and a run that
  labelled without asking would be deciding the version, which is the thing forbidden here.
- **the bump question in the attended `/libretto-flow`.** Its phase 8 already stops with
  the user present; the red check ambushes nobody there. The user's call, 2026-08-13 —
  back the day an attended run pays the same round trip.
- **creating the `release:` labels where a repository defines none**, defaulting the bump
  when the question goes unanswered, or re-asking a request that already carries one.
  The first decides somebody else's release convention; the second is the silently-wrong
  bump wearing a politer hat; the third defends a state attacca cannot reach, because it
  opens the request itself in the same run.
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
- **a cross-project retro.** The retro runs where the flow ran, on that project's ledger.
  Back the day one payload lesson shows up in two projects' ledgers.
- **the retro applying a payload diff.** Propose only; the payload is the product and does
  not get edited as a side effect of a retro. Back, if ever, as an explicit flag after the
  proposals have earned trust — never the default.
- **capturing the model's self-detected failures.** Gates and `evidence` already own
  those; the ledger records *user* corrections, the signal that is unambiguous and
  otherwise lost.
- **a hook or automation for capture.** Recognising "the user is correcting me" is
  judgment, and it lives in a skill instruction.
- **editing or deleting ledger entries**, from any skill. Append and mark. History is the
  point, and an edited entry is history that lies.
- **the Go binary writing the ledger.** Delivery reads it (the corrections column in
  `metrics`); only the payload writes it.

### A stop is a place where the user changes something

That is the whole test, and it is what the count is derived from.

| After | Stops | What is being changed |
|---|---|---|
| 1 · find-work | no | — the reading is stated, and the spec is where it gets corrected |
| 2–3 · write-spec | **yes** | the contract |
| 5→6 · write-tasks | **yes** | the approach, the order, and what waits on what |
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

**Phase 2 interviews before it writes the spec** — one question per native-prompt call,
the recommendation carrying its reason, "no more questions" always an option, a soft
bound around five that judgment crosses in either direction — so the contract is built
by two people rather than handed over finished, and each answer can redirect the next
question, which a batched call never could. **Phase 5 asks exactly one: the approach**,
chosen from two or three with named tradeoffs. Every answer lands verbatim in the
change's `decisions.md`, the log the writer subagents read. **Zero is a legitimate
answer, reported in one line:** a quota manufactures questions the code already
answers, which is the rubber-stamp round trip removed from three other phases arriving
back through the door marked collaboration. The 5→6 seam stops for the approach and the
order and opens no second tranche, and **the stop count does not move** — the questions
ride the phases that already have them. The trivial lane asks nothing, because *no spec
needed* means there is no contract to disagree about.

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
| a change already in flight | unchecked boxes in `.agents/changes/*/tasks.md`, and in `plan.md` for a change created before the rename |
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

**A branch is work in flight too.** Scanning `.agents/changes/*/tasks.md` cannot see the
trivial lane: a change that needed no spec has no plan to scan, so it exists only as
commits on a branch. The lane's first real run produced exactly that, and the next phase 1
reported an empty house. Phase 1 also reads `git branch --no-merged` and the forge's open
requests, and names the state worth naming: **unpushed and un-requested** is work nobody
but this machine has.

**Landed changes are deleted, not archived.** gentle-ai moves them to `changes/archive/`;
git history is the archive here, and a directory nobody reopens is growth. A decision,
not an oversight.

## Constraints

**No skill, agent or command references `docs/PAYLOAD.md`, and nothing gates that.** The
uninstalled-path check is scoped to executables — `scripts/` and `bin/` — because prose mentioning
`docs/` is describing rather than instructing, and failing those turns the check into noise nobody
reads. So the index sits in the same position as `docs/FLOW.md`: a path this repository has that the
user's project does not, held by this constraint rather than by a gate.

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

**That record is `THIRD-PARTY.md`, and this capability governs it.** The three upstream
licence texts live in `licenses/`, so root-level `LICENSE` is the only licence file a reader
meets there — none of the three is an alternative to it, which is how the root listing read.
**No vendored licence text is ever deleted, reworded or merged**: a vendored copy has to carry
its own licence, and tidying one out of a directory listing is a licensing failure rather than
a cleanup.

`THIRD-PARTY.md`'s relative links resolve, and that is checked — the gate above already
*parses* that file's table to derive the vendored list, so it was load-bearing while no
capability's `Governs:` claimed it. A path nobody owns is a path where drift is nobody's
finding, which is the sentence the `readme` capability was created to answer about `README.md`.

## Prior decisions

- **The finished-and-not-landed signal is the difference between the two scans phase 1
  already runs, never a third command.** A third scan is a third thing to keep in step
  with the other two, and the first time they disagree the report is wrong in a way
  nobody can see.
- **The report stops at the signal and never checks whether the delta was applied.** That
  means reading a capability spec and deciding whether a delta is *present* in it, which
  is a reading — wrong in one direction it accuses a correct landing, wrong in the other
  it clears a broken one. Zero open boxes is the signal; the user interprets it.
- **A plan's `Durable decisions:` line is a claim about whether the list is empty, never
  the list.** The list is the change's `spec.md` *Prior decisions*. Both readings looked
  right and the tree held the contradiction to prove it: one plan declared "the two" over
  a section of three, another carried no line at all. Found by a 5→6 cutter reading both
  documents cold.
- **`plan.md` was reused for the technical approach rather than adding a `design.md`.**
  The complaint was that the file named `plan` is not a plan; a `design.md` beside an
  unchanged `plan.md` leaves that exact file unchanged and the complaint standing.
  Answered by the user, 2026-08-17.
- **The task cut is a subagent, not a numbered phase.** Independence comes from the fresh
  context; renumbering everything that says eight buys nothing, which is the answer the
  6→7 review seam already gave for itself.
- **Phase 5 asks, since 2026-08-19 — retiring the 2026-08-12 decision that it asked
  nothing.** The user reversed it for the fork alone: the approach is the one phase-5
  decision that is genuinely the user's, and a finished decision can only be
  rubber-stamped while a fork with named costs must be decided. The no-third-stop
  argument the original decision was made for survives: the fork rides inside the phase.
- **Four answers from the user, 2026-08-19, verbatim options:** the interview bound is
  *"Blando ~5, con juicio"* — never a hard cap; the single-spec case gets a brief always
  (*"Brief siempre"*) so `spec-writer` keeps one prompt contract; markers are patched by
  the orchestrator inline (*"Orquestador parchea inline"*), scoped to the bracket
  expression, the one declared exception to one-author-per-file; attacca's assumptions
  land in `decisions.md` marked `(assumed)` (*"Sí, marcadas 'assumed'"*), one home for
  decisions in both modes.
- **Deliberately not built, with the conditions that bring them back:**
  section-by-section plan validation, when plans outgrow one read; markers in plans —
  `plan-writer` returns gaps in its reply, like the cutter; a hard cap of five, if the
  soft bound degenerates into interrogation.
- **The specs wiki regenerates at the landing step, not on a schedule and not by a
  daemon.** The landing is the moment the specification moves, so the refreshed index
  rides the same commit as the delta that changed it; a generated view left behind by a
  landing is drift wearing a marker comment. Assumed 2026-08-18 rather than asked — if
  manual-only regeneration is wanted, this instruction is the part to drop.
- **The retirement gate compares the *section*, never the file.** Requiring any edit to a
  capability spec passes on the delta application alone — that happens in the landing
  commit by definition — so the gate would be green on every landing and measure nothing.
  It is the version that looks like a gate and is not one.
- **Its escape is a declaration in the plan, not a flag.** A flag is typed by whoever
  wants the commit through, at the moment they want it to stop complaining; a line in the
  plan is written while the plan is, by the person who knew. No mechanism can stop that
  line becoming a reflex — `review-work` reads the plan and the diff together.
- **New checks ride inside `--anchors` rather than becoming new gates.** The count "six
  gates" is written in ten places across `AGENTS.md`, the CI spec, the contributing spec
  and the workflow, and a number kept in ten places is a number that drifts. Both the
  EARS half and the retirement half were added this way.
- The tracker is read through its CLI, never MCP, never the REST API.
- The API token lives in the OS keyring, put there by `jira init`, run by the user in
  their own shell. It never enters a conversation.
- Specs are per **capability**, never per ticket: a capability spec accumulates and
  stays true, a ticket spec is dead the day the ticket closes.
- Deltas live in `.agents/changes/<change>/` and are applied onto the capability spec in
  the commit that lands the change, which then deletes the change folder.
- Drift and landing checks in this flow **warn or stop the author, never block someone
  else's commit** unless they opted in. A check that stops a commit in someone else's
  project is a check that gets deleted; a clause that stops the author mid-flow —
  `libretto land`'s fix-and-re-run — instructs the agent running the flow and installs
  no hook.
- **The payload learns about `libretto land`, minimally — one guarded clause in
  `record-work`, the wiki clause's shape.** Assumed 2026-08-21 under attacca: a
  verifier nothing invokes verifies nothing, so the landing step gains "where
  `libretto` is on PATH, run `libretto land` before the landing commit", absent-binary
  path unchanged. The guard is what keeps the skill self-sufficient once installed —
  `libretto` stays delivery, never a dependency. If wrong: the clause is one sentence
  to remove.
- Phase 2 may decide no spec is needed. Skipping the phase is a legitimate outcome of
  it, and the "no" collapses phase 7's gate with it.
- **The vertical-slicing rule lands in `write-plan` alone**, though the horizontal cut
  originates in phase 2's task breakdown. `write-plan` already owns the authority to reject a
  task the spec cut wrong, so the rule lives where it is enforced rather than in two files
  that then disagree — the failure `CLAUDE.md` opens by naming. If phase 5 starts sending work
  back to phase 2 repeatedly, the answer is a sentence in `write-spec`'s task-breakdown pillar
  and this rule stays as it is.
- **Nothing checks a plan's content, and nothing can.** Whether a box genuinely merges alone
  is judgment, and a check that exercises judgment drifts. The gate verifies the mandate is
  present in `skills/write-plan/SKILL.md` and no more. A box cut horizontally under the
  mandate surfaces at phase 6 or in the 6→7 review, never in a gate — the replacement, the
  day that bites, is a reviewer finding and not a longer regex.
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
- **the bump is asked at attacca's end rather than left as a red check to be discovered** —
  the user's call, 2026-08-13, and the reversal of one word of the labelling non-goal above.
  It holds only because the run asks and types rather than reads: `AGENTS.md` already splits
  those, and every rule around the question exists to keep the split intact. Two calls taken
  with it, both with what changes if they are wrong: **attacca only**, so an attended run
  keeps paying the round trip if that turns out to matter; and **detect the labels rather
  than hardcode this repository**, so a project whose `release:*` labels mean something else
  gets a question that does not apply and answers by ignoring it. **Ceiling named:** the
  question is only as reachable as the terminal it is asked in — a scheduled or piped run
  never sees it and lands unlabeled, which is today's behaviour and therefore the fallback
  rather than an error. The replacement, if that becomes common, is a `gh pr comment`
  carrying the three commands.
- **a wiring proof cannot tell a sentence from a step, and this change paid for the lesson** —
  the phase 6→7 reviewer found the request's description promised the chosen bump with
  nothing in the flow writing it there, while the criterion's row confirmed the sentence
  existed and passed. Four of its five findings were defects that shipped green. The rule
  that follows is not a new check but a reading of the existing one: a criterion proved by
  `check-payload` is proved to be *present*, and any criterion whose value is that something
  *runs* needs the step named in the same file, not only the promise.
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

- **The lessons ledger lives per project at `.agents/lessons.md`** — the user's call,
  2026-08-13. Central-in-`~/.claude` was rejected (outside git, no review possible, mixes
  projects); per-change files were rejected because phase 8 deletes the change folder at
  landing and lessons must outlive the change that taught them.
- **The retro proposes payload diffs and never applies them** — the user's call,
  2026-08-13, choosing this over an automatic retro inside phase 8, which was rejected as
  too much power without eyes on it.
- **Capture lives in `evidence`, not in eight phase skills.** One place, already invoked
  at every phase; eight copies of one rule is eight things that drift.
- **A lesson is classified by where the fix lives.** Project knowledge goes into that
  project's contract; a flow defect goes to the payload skill as a proposal. A retro that
  mixes them writes one project's manias into the flow and breaks it for every other
  project. An entry that could read both ways is project knowledge until the same lesson
  appears somewhere a second time — the cheaper wrong guess.

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
- [x] the capture rule in `evidence`, with the entry format stated once
- [x] `retro` — read, classify, record or propose, mark resolved
- [x] `libretto-retro` — route to the skill, describe nothing
- [x] the retro wiring guarded in `check-payload`, observed failing on a hand-broken
      copy before landing

## Verification criteria

- **When** a change folder holds `tasks.md` or `plan.md` with at least one box in it and
  none of those boxes open, phase 1 **shall** report that change as finished and not
  landed, naming it — and the status command **shall** carry the same report by delegating
  to that scan. **Where** the folder holds neither file, it **shall not** be reported that
  way: a captured idea is not an unlanded change. The plan skill **shall** state that a
  `Durable decisions:` line is a claim about whether the list is empty and not the list
  itself. **Ceiling named:** the gate proves each mandate is *present* in its file, never
  that a session obeyed it — that is behaviour, checked by running it, and it surfaces in
  the 6→7 review or not at all.
  Proof: scripts/check-payload
- **If** a staged commit deletes a change's `plan.md` and no capability spec's *Prior
  decisions* section differs between `HEAD` and the index, **then** `spec-drift` **shall**
  fail and name the change being landed — **where** the deleted plan declared `Durable
  decisions: none`, it **shall** pass instead. **When** no such deletion is staged, it
  **shall** report nothing and **shall not** fail. Driven through real temporary
  repositories, because a fixture made of strings would prove the awk and not the git
  plumbing, and the plumbing is where a gate goes silent — which reads as green.
  Proof: skills/record-work/spec-drift --self-test
- **If** a verification criterion carries no EARS `shall`, **then** `spec-drift` **shall**
  read it as unfailable — **where** emphasis or backticks wrap the keyword, it **shall**
  still be read as present. **Ceiling named:** the self-test drives the marker through the
  same `is_ears` the gate calls, so a matcher that breaks fails here. What it does not
  cover is the hard-on-deltas, soft-on-capabilities asymmetry, which has no fixture and is
  observed only by running the gate.
  Proof: skills/record-work/spec-drift --self-test
- **Where** `libretto` is on PATH and the project holds a consolidated specs directory,
  the `record-work` skill **shall** instruct the landing step to run `libretto wiki` and
  include the refreshed index in the landing commit; **where** the binary is absent, it
  **shall** say the wiki may be stale and move on rather than block the landing.
  **Ceiling named:** the anchor keeps the instruction findable in its file, never proves
  a session obeyed prose — the same limit every skill criterion here lives with.
  Proof: skills/record-work/SKILL.md
- **Where** `libretto` is on PATH, the `record-work` skill **shall** instruct
  the landing step to run `libretto land` before the landing commit; **if**
  the command exits non-zero naming a missing part, **then** the skill
  **shall** instruct fixing that part and re-running, never committing past
  it; **where** the binary is absent — or too old to know `land`, failing
  with *unknown command* and naming no part — it **shall** say the landing
  is unverified and continue rather than block. (Assumed 2026-08-21, from a
  cutter finding: an old binary must not wedge a landing; if wrong, this
  clause is one sentence to tighten.) **Ceiling named:** the anchor keeps the
  instruction findable in its file, never proves a session obeyed prose — the
  same limit the wiki-clause criterion beside it lives with.
  Proof: skills/record-work/SKILL.md
- **When** the `libretto land` clause lands, `scripts/check-payload` **shall** pass: the
  reference is to a binary on PATH, not an uninstalled repository path, and
  the guard is what makes that true. This proves the reference is legal, not
  that the clause is followed.
  Proof: scripts/check-payload
- frontmatter parses, and `name:` matches the directory or filename
  Proof: scripts/check-payload
- no stray file sits where the linker would install it as an item
  Proof: scripts/check-payload
- every referenced skill exists
  Proof: scripts/check-payload
- **no skill invokes a path that does not get installed**
  Proof: scripts/check-payload
- **prose addresses the agent, never Claude by name.** The payload installs into Codex
  and OpenCode too, and "ask Claude" read there is an instruction about somebody else.
  Every occurrence of `Claude` in skill and command bodies must match the allowlist of
  factual uses — `CLAUDE_HOME`, `CLAUDE.md`, real `.claude` paths, the product name
  `Claude Code` — or the gate fails; the check deletes allowlisted tokens per hit
  before re-searching, so an addressee sharing a line with a factual use still fails.
  The classification lives entirely in the allowlist; the check exercises no judgment.
  Proof: scripts/check-payload
- **a Claude-only mechanism names the host-neutral capability, in every file that mandates
  one.** The dependency stays — `AskUserQuestion` and `Skill(skill="…")` keep Claude
  Code's names, because a capability described and never named is a capability nobody can
  invoke. What the file owes alongside is the capability, so a model on Codex or OpenCode
  reaches for its own equivalent: a **marker phrase**, `native prompt` for a mandate to
  ask and `host's own` for a mandate to load a skill. Either satisfies the check.
  `record-work` remains the canonical statement of the rule; the per-file line points at
  it and never restates it.

  **This narrowed a promise.** It read "at each mandate *site*" and the gate it cited
  checked no part of it — the addressee half was all `check-payload` ever tested, so the
  criterion was green from the day it was written while seven files broke it. Site-level is
  now file-level, which is what a check can enforce against prose that wraps.

  **Two markers and not one**, because the four skills that already complied all carry
  `native prompt` and none carries anything else; a single new marker would have failed the
  files that got it right. **Newlines are squashed before searching** — the phrase wraps
  across line breaks at this width, and a line-scoped search reported `write-plan` as
  broken while this check was being written. `agents/**` is excluded: agents install into
  Claude alone, so there is no second host for their prose to be wrong for.

  **Ceiling named, and it is wider than "a file that already complied".** The check asks
  whether the marker appears anywhere in the file, so **any** occurrence satisfies it —
  including one in prose that has nothing to do with a mandate. A file whose only match is
  the sentence "talk about a host's own dog here" passes, having never named a capability
  at all; measured with a probe file, not reasoned about. A new bare mention riding into a
  legitimately compliant file passes for the same reason.

  This is accepted rather than closed, and the first draft of this bullet described only
  the narrower half — which is the worse error of the two, because a ceiling stated in a
  criterion is the one place a reader is entitled to trust about what the gate does not
  catch. **The gate stops whole files with no pointer at all**, which is the failure that
  actually happened: seven of eleven. It does not verify that a marker is doing its job,
  and it cannot — "is this sentence about the mandate" is judgment, and the addressee check
  beside it exists precisely because a check that exercises judgment drifts.

  The replacement, the day a file drifts internally or somebody games it, is a check scoped
  to the paragraph — never a longer regex.
  Proof: scripts/check-payload
- **no skill hardcodes the install layout.** A `~/.claude/` path is only true under
  `install --global`; a skill's tools resolve from its own base directory, which every
  invocation announces — `record-work` reaches `spec-drift` as its sibling, `write-spec`
  hops to `../record-work/`. Both layouts keep skills side by side.
  Proof: scripts/check-payload
- **`skills/write-plan/SKILL.md` carries the vertical-slicing mandate.** The gate searches the
  file for the literal phrase `one badly cut box`. `check_wiring` is `rg -qN`, so the match is
  **line-scoped and not newline-squashed** — unlike the marker check above it, and the phrase
  therefore has to sit on one line. The criterion says so because its first draft claimed
  squashing this gate does not do, and a criterion describing a stronger check than the one it
  cites is how a promise goes green for the wrong reason.

  **Watched red before green, and then again by the reviewer.** The gate went in first and
  failed naming the absent phrase; the 6→7 reviewer independently broke the phrase to
  `one badly-cut box` in a throwaway copy and got the same failure, which is the run that
  proves the gate is sensitive to what it guards rather than merely present.
  Proof: scripts/check-payload
- **the mandate introduces no bare `Claude` addressee**, and **a host-neutral marker is still
  present in `skills/write-plan/SKILL.md`** after the edit. Two criteria and not one joined by
  `and`: separate checks, separate failure modes, and joined either could be half-met and read
  as done.
  Proof: scripts/check-payload
- **`docs/PAYLOAD.md` lists every skill, agent and command that ships, and the gate fails when it
  has drifted.** `scripts/check-payload --index` writes it; the default run regenerates in memory
  and compares. One parse serves both, so a page and a gate cannot disagree about what an item is
  called — and `--index` writes while the default run only compares, because a gate that repairs
  what it measures can never fail.

  **Generated, never typed.** A hand-written list of directories is the failure this repository has
  paid for twice; `docs/SPEC.md` is the only place the capability list lives for exactly that
  reason, and a typed `PAYLOAD.md` reproduces it one directory deeper.

  **Watched biting, not merely passing**: absent page, a hand-edited description, an added item, a
  removed row — all four red with the fix named in the message.

  **Two defects the 6→7 reviewer found here, both of which had the gate green:**

  - **A folded `description:` produced an empty cell.** Four items — `caveman`, `caveman-commit`,
    `ponytail`, `ponytail-debt` — write `description: >`, so the page carried a literal `>` for
    four of thirty-six rows. The delta had named single-line descriptions as a *ceiling* and
    deferred the fix to "the day an item needs a folded description": **that day was already four
    items in the past.** A ceiling is a claim about the future, and one line of `rg` would have
    checked whether the present already broke it. Both block forms are read now; anything more
    exotic must `fail` rather than reach a YAML parser.
  - **The sort collation was unpinned.** With empty cells the order turned on punctuation alone,
    and glibc collates that differently from byte order — so the gate would report drift on a page
    nobody edited, on CI rather than on the author's machine. `LC_ALL=C` is pinned inside the
    pipeline; output verified byte-identical across three locales.

  Also: **the page's own boilerplate claimed everything is installed by symlink**, which is false
  for an OpenCode agent — a derived file this tool writes — and for Codex, which takes skills only.
  A generated file's prose is the part nobody re-reads on regeneration, so a false claim in it
  survives every future run.
  Proof: scripts/check-payload
- **every relative link in `THIRD-PARTY.md` resolves**, scanned over the flattened document
  because two guards in that test file have already shipped unable to fire on a wrapped
  phrase; **and the three vendored licence texts are under `licenses/` with none left at the
  root**, asserted directly because a link nobody updated still points at a file that is still
  there — so the link half alone passes with one moved and two forgotten.
  Proof: cmd/libretto/readme_test.go TestThirdPartyLinksResolve

  **Both halves watched red, separately.** Layout first, against the files still at the root.
  Then the links, against the files moved and the two lines untouched — which is precisely the
  silent breakage this criterion exists for.

  **`payload`'s `Governs:` is now wider than the payload directories**, so a licence-only
  commit reads as payload drift. That is the intent; the replacement, if it becomes noise, is a
  `vendoring` capability rather than a narrower glob.

  **No criterion asserts the `Governs:` widening itself.** The obvious citation was
  `spec-drift --trace`, and it cannot fail — `--trace` is a map and returns 0 whatever it
  finds, which its own line 42 says. A criterion citing a gate that can never be red reads as
  proven and is worse than one with no citation, so the widening is an observation in the
  report instead. The alternative was a test asserting that one line of one spec file contains
  two strings, which is machinery to check a document the next reader checks for free.
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
  **The proof is the wiring, not the behaviour**, and the difference is stated because
  citing a static checker for a prompt's conduct is a citation that can never fail: review
  deleted every `review-spec` line from the flow and this script stayed green. What is
  checked is that the phase exists, that both commands route to it, and that the decisive
  words are still in the file that owns them.
  Proof: scripts/check-payload
- **a bug amends the contract before it touches the code.** A bug is a hole in the
  specification — behaviour happened that some capability permitted, failed to forbid, or
  forbade in words no run could have failed — so phase 2 names the criterion that would have
  caught it and writes its `Proof:` as a test that **fails when written, for the reason the
  bug exists**. Watching it fail is evidence obtainable only before the fix. Reversed, the
  criterion is written to describe the fix rather than the failure: it passes on its first
  run, nothing ever proved it could fail, and the spec gained a sentence that cannot catch
  the bug it was written for. Almost always an amendment to the capability whose `Governs:`
  owns the broken path — a new spec directory per failure is a ticket spec, dead the day the
  fix ships. No capability owning the path is itself the finding. **Proved as wiring only** —
  that the branch and its red-test-first demand are still in `write-spec`, not that a session
  obeyed them.
  Proof: scripts/check-payload
- **phase 1 says whether the work is a bug, and never infers it from a summary.** The branch
  above cannot fire on work that arrived as generic, and "discount wrong on bundles" is a bug
  or a change of intent — opposite pieces of work behind one sentence. The reading names it
  in one line and `proposal.md` records the observed behaviour, the expected behaviour and
  what produced it, in the reporter's words: paraphrasing costs the reproduction, and with no
  reproduction there is no failing test, so the criterion is a sentence somebody hopes is true.
  **Proved as wiring only**, same as above.
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
- **a `Queued:` line is necessary and not sufficient: a change dispatched by a branch is
  reported in flight, never in the queue.** The scan reads the working tree, and the
  working tree is one branch's opinion — the line comes out in the pickup commit, which
  lives on the feature branch, so from the base branch a change that has been built,
  reviewed and pushed still reads as captured-and-not-started. The queue result is
  filtered against the branch scan the same phase already runs, matching whole names with
  a conventional prefix removed; an unrecognised prefix fails safe towards leaving the
  change in the queue rather than losing a captured idea.
  **Written from the failure**: `/libretto-status` reported a finished change under *the
  queue*, and the next command typed was `/libretto-attacca` on it — which would have
  branched from the base, found the proposal still queued there, and rebuilt work that
  already existed into a second request and a conflict on every spec it had amended. **A
  report that induces rework is worse than one that says nothing**, because silence costs
  a question and this costs the work twice.
  **Proved as wiring only** — that the rule is still in the skill that owns it, not that a
  session obeyed it.
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

- **a phase opens the spec that governs its work, never the corpus.** `write-spec` step 2
  looks up the owning capability in a fixed order — the change's `Targets:`, then the
  `Governs:` line claiming the path, then the project's index — and names reading the whole
  spec root as the thing not to do. A specification is written to accumulate, so a phase
  that opens all of it to learn one name pays for the corpus to answer a question one line
  settles. None of the three found means no capability owns the path, which is a finding to
  report rather than a licence to sweep. Stated without naming any repository's layout,
  because the skill installs into projects that organise their specification differently —
  **that property is a review-time reading and deliberately uncited**: no row can prove an
  absence over an unbounded set.
  **Proved as wiring only.**
  Proof: scripts/check-payload
- **a fan-out writer is given named brief sections, and always the vocabulary.** The brief's
  five contents are a fixed, enumerated heading set, so the prompt can name sections back to
  the writer and both sides hold the same list. Each writer opens what its subtask touches;
  the vocabulary is never trimmed, because it is the whole mechanism stopping two writers
  from giving one concept two names, and a brief read in slices without it is worse than a
  brief read whole. **Sections are named, never excerpted** — excerpting puts N copies in N
  contexts, which is the cost the brief was introduced to remove. The return contract is
  unchanged and carries no criterion: `agents/spec-writer.md` already promises deltas rather
  than restatements and already forbids the empty return, so a row over it would have been
  green the moment it was written.
  **Proved as wiring only.**
  Proof: scripts/check-payload
- **model and effort do not move part-way through a fan-out.** A phase is billed a fraction
  of input price while its prefix stays byte-identical between calls; switching the tier
  invalidates that prefix and rebills the whole context at full price, and a fan-out pays it
  once per writer. `write-spec` step 2b states it beside the fan-out that costs it, and
  `docs/FLOW.md` carries the same reasoning under *Delegation* — **uncontracted and cited by
  nothing, because no capability governs `docs/`**, and a `Proof:` over an unowned file
  anchors to nothing.
  **Proved as wiring only**, and this one cannot be more even in principle: the dial is the
  session's, and no skill can read it or stop a hand from moving it.
  Proof: scripts/check-payload

- **the retro wiring holds: capture in `evidence`, verbatim `Said:`, `Resolved:`-only
  marking, propose-never-apply, and the command routing to the skill.** Wiring only — a
  prompt is checked by running it, and these rows are what keeps the decisive words in
  the files that own them.
  Proof: scripts/check-payload
- **`spec-drift --block` turns the warning into exit 1, and only for whoever asked.** The
  same checks as default mode; the default still always exits 0, because a gate that
  surprises someone in their own project gets deleted and a deleted check finds nothing.
  The flag exists so opting in costs one paste — `record-work` documents the pre-commit
  snippet and names it opt-in — and nothing in the flow ever installs the hook. Proven
  end-to-end in a throwaway repository inside the self-test: exit 1 on drift, exit 0
  when the governing spec moved too, default unchanged.
  Proof: skills/record-work/spec-drift --self-test
- **the block hook snippet is documented as opt-in**, in the file that owns the
  warn-never-block reasoning, so the gate and the reason it is not the default live one
  paragraph apart.
  Proof: scripts/check-payload
- **if judging the change means looking at it, phase 6 renders it and looks before the
  review seam.** A palette, a layout, a panel row, an image qualify; output that is read,
  not looked at, does not. The project's own render path, what-was-seen in the evidence,
  measured contrast where the change is about colour. A rule inside the phase, never a
  stop — the seam's reviewer reads specs, diffs and test output, not pixels, so a render
  nobody looked at in phase 6 is a render nobody looks at anywhere. The 1.4:1 palette
  that satisfied its spec is the incident behind the rule. **Proved as wiring only.**
  Proof: scripts/check-payload
- **the 6→7 seam ledgers each finding under the phase value `6→7`**, one entry per
  finding, fixed or not, on the same header contract `evidence` writes — which is what
  makes findings countable and lets `libretto metrics` keep them out of the
  user-corrections column. The write is the seam's; the reviewer subagent still writes
  nothing. **Proved as wiring only.**
  Proof: scripts/check-payload
- **phase 2's questions are an interview: judgment around a soft five, biased to ask.**
  The hard cap of three was lifted by the user on 2026-08-14 — better asked out of
  caution than swallowed out of fear — and the one-call bound was retired by the user on
  2026-08-19: one question per call, each answer able to redirect the next, "no more
  questions" always offered. Both edges stay named: every question one a wrong guess
  would make expensive, never a form-length interrogation of things the code already
  answers, zero still legitimate and said in one line. The promise lives in three homes —
  `write-spec`, `docs/FLOW.md`, `commands/libretto-flow.md` — and moves in step or not at
  all. **Proved as wiring only**, three rows: the interview shape, the judgment rule and
  its upper edge.
  Proof: scripts/check-payload
- **Where** a change reaches phase 2 and a spec is written, `write-spec` **shall** direct
  the session to create `decisions.md` in the change folder — first write creates it —
  and record each answer verbatim, dated, `(assumed)` when nobody gave it, before the
  spec file is written. One writer, the orchestrator; the durable copy is the delta's
  *Prior decisions*, made per answer. **Proved as wiring only.**
  Proof: scripts/check-payload
- **the spec's author is `spec-writer`, single case included** — a single writer is a
  fan-out of one, launched with `brief.md` and `decisions.md` by path. **Where** the
  inputs do not settle a decision, the writer **shall** leave `[NEEDS CLARIFICATION:
  question]` and never guess; the orchestrator **shall** resolve each marker by asking,
  logging, and replacing the bracket expression only — the one declared exception to
  one-author-per-file, anything beyond it a relaunch. **Proved as wiring only**, three
  rows across the skill and the agent.
  Proof: scripts/check-payload
- **When** phase 5 runs on a change with a contract, `write-plan` **shall** present two
  or three approaches with tradeoffs as one native question, recommended first with its
  reason, and **shall** record chosen and rejected with why in `decisions.md`. `plan.md`
  **shall** be drafted by `plan-writer` — tools `Read, Grep, Glob, Skill` and nothing
  else, pinned by its own wiring row — returning markdown the orchestrator writes.
  **Proved as wiring only**, three rows: the fork, the log, the launch.
  Proof: scripts/check-payload
- **While** a run is `/libretto-attacca`, the interview, the fork and every marker
  **shall** become `(assumed)` entries in `decisions.md`, each naming what changes if
  wrong — logged where the writers read, so an assumed answer reaches them the same way
  a given one does. **Proved as wiring only.**
  Proof: scripts/check-payload
- **the three stops are each asked natively, and the check is the string search across
  stop-owning skills** — one row per stop: `write-spec`, `write-tasks`, `record-work`.
  That a fourth stop does not exist stays uncheckable, the same ceiling the stop table
  above has always named — a script cannot tell a prompt from a paragraph.
  Proof: scripts/check-payload

The bump question's rows, one per condition. **Each was written before the prose it
describes and observed red**, except where noted in the change's own record — a row added
to match prose that already exists has never proved it could fail, and this change found
one of those in its own plan before the reviewer found four more in its skill.

- **the bump is asked at the end of the attacca path in `skills/record-work/`**, which is
  the file that owns phase 8
  Proof: scripts/check-payload
- **the applied label is read back off the request.** A command that printed no error is
  not a change the forge accepted — the rule the push already carries
  Proof: scripts/check-payload
- **the chosen bump is written back into the request's description.** The description
  exists before the question is asked, so it reaches it only by being written back. Found
  by the reviewer as a promise with no step behind it, passing its own row
  Proof: scripts/check-payload
- **`release:major` is present in the prompt and is never the first option.** Recommending
  it is the announcement `AGENTS.md` forbids, and announcing it is what published `v1.0.0`
  Proof: scripts/check-payload
- **label detection is searched and bounded, never a bare `gh label list`.** The default
  fetches 30 ascending, so `release:*` falls off the page and the check answers "none of
  the three" for a repository that defines all three — inverting silently
  Proof: scripts/check-payload
- **labels are matched by whole name, never by prefix or substring.** `release:patch-hotfix`
  contains `release:patch` and is not it, and a detection looser than the workflow it feeds
  finds labels `gh pr edit` is then refused
  Proof: scripts/check-payload
- **a repository defining none of the three is not asked**, and **the closing report's
  red-check line survives an unanswered question** — the line is written before the question
  and is never withdrawn by it
  Proof: scripts/check-payload
- **the push consent does not extend to a label**, in `record-work` as well as in the
  command. The same absolute sentence lived in both; only one was amended and only one was
  guarded, in the file the contract does not designate as the owner
  Proof: scripts/check-payload
- **`commands/libretto-attacca.md` restates none of it** — no prompt, no `gh pr edit` — and
  its `Never` list forbids labelling only with a bump the user did not choose, while merging,
  tagging and releasing stay absolute
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

**The lessons ledger and the retro are prose and none of it has run.** Claims, not
facts: a correction mid-flow producing an entry without interrupting the phase; a retro
classifying a real ledger and putting each fix where it belongs; a proposed payload diff
a user could apply verbatim. The first real flow after this lands is the test.

**The queue is prose and none of it has run.** Claims, not facts: capturing two ideas
leaves two committed proposals and no branch; `/libretto-next` offers the oldest first and
enters phase 2 on a fresh branch; `/libretto-flow` handed a key never mentions the queue;
`/libretto-status` shows the queue as its own section. What *was* observed on the run that
landed it is the reviewer catching three defects in the prose — a duplicated scan, a commit
step stated in a never-does bullet instead of instructed, and a queued idea defined two
ways in adjacent lines.
