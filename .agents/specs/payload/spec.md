# Payload

Governs: skills/** commands/** agents/** scripts/** THIRD-PARTY.md licenses/**

The reason the project exists. Everything else is delivery.

The flow itself — its eight phases and its reasoning — is [docs/FLOW.md](../../../docs/FLOW.md).
This spec is the contract the payload's *artifacts* have to satisfy.

## Outcomes

Installing this repository gives a working flow on a machine that has nothing else.

- one command, `libretto-flow`, that routes and never implements
- `libretto-status`, read-only, reporting what is in flight
- **`libretto-converge`, which re-derives a checklist that stopped describing the code.**
  `/libretto-converge [<change>]` routes to the `converge` skill, which reads one change's
  open boxes against the codebase and ends with each open box in exactly one state —
  **confirmed open**, **closed with observed proof**, **reported as drifted**, or
  **reported as unmapped** — and with one appended box per criterion of the change's own
  that no box closes and whose proof does not pass. `loop` drives a checklist; this one
  repairs one that drifted from what got built.
  - the judgment is made by one fresh `converge-reader`, which never had the conversation,
    reads the tree as it is and cannot run anything; the checklist is written only by the
    session
  - a box is closed only after the session ran every proof its criteria cite, in the
    foreground, and read each one passing by the rule for its kind; the closed box's last
    line carries every citation that passed and the date
  - a drifted box is untouched in the file and named in the report, with what the box says
    and what the code does instead
  - it works on both layouts — the `## Checklist` of `proposal.md`, and `tasks.md` or
    `plan.md` for a change begun before
  - when it marked or appended anything it ends with one commit holding the checklist file
    alone; when it did neither it commits nothing; it never pushes. A door the user opens,
    like `libretto-status` — neither `libretto-flow` nor `libretto-attacca` routes to it
- `libretto-review`, which reviews a forge PR/MR in a workspace that restores itself
  — its contract is the `review-project` capability spec
- **a queue for ideas that arrive faster than they get built** — `libretto-queue`
  captures them one after another as proposals carrying a `Queued:` date, committing
  each on the current branch and creating no branch, no spec and no plan;
  `libretto-next` picks one, oldest first with the user free to choose another, creates
  the branch, removes the `Queued:` line and enters the flow at phase 2, because phase
  1's artifact already exists
- one skill per phase, each of which stops where its phase stops
- **one change document.** A contract-lane change keeps everything in
  `.agents/changes/<change>/proposal.md` — the file `find-work` or `/libretto-queue`
  captured, grown in place. No `brief.md`, `decisions.md`, `spec.md`, `plan.md` or
  `tasks.md` is created for it. Top to bottom:

  | Part | Where | Written by phase |
  |---|---|---|
  | `Tracker:` · `Targets:` · `Durable decisions:` lines | top, column 0, outside any fence | 1 (`Tracker:`), 2 (`Targets:`), 5 (`Durable decisions:`) |
  | the ask, verbatim · the reading | the headings the capture writes | 1 or the queue |
  | the decision log | `## Decision log` | 2 and 5. One `- Q: … → A: …` line per answer, verbatim, under dated `### Session` headings, `(assumed)` kept |
  | the contract | `## Outcomes` · `## Scope boundaries` · `## Constraints` · `## Prior decisions` · `## Task breakdown` · `## Verification criteria` | 2–3 |
  | the approach | `## Approach`, the parts `write-plan` names as `###` subsections | 5 |
  | the checklist | `## Checklist`: boxes and nothing else, the document's last section | the 5→6 seam |

  **The checklist is the boxes under the level-2 `## Checklist` heading, and nothing else**;
  the *Task breakdown* pillar is a list, never boxes. The approach holds no box and the
  checklist holds nothing but boxes — the separation that used to be two files, `plan.md`
  and `tasks.md`, survives as two sections. It exists because phase 5 once produced a
  checklist and called it a plan, so the approach chosen, the alternatives it beat, the
  risks and how it gets validated died with the session that had them. **The alternatives
  table is the pillar that pays**: a diff shows what was built, and nothing in a repository
  shows what was not built and why. **A change amending two or more capabilities** carries
  the primary capability's contract in `proposal.md` and one `spec-<capability>.md` beside
  it per extra capability, each with its own `Targets:` line and six pillars; the decision
  log, the approach, the `Durable decisions:` line and the checklist exist once, in
  `proposal.md`
- **the change document is co-authored with the user and written inline by the session
  that held the conversation.** Phase 2 interviews one question at a time and phase 5 forks
  on the approach; every answer lands verbatim in the document's `## Decision log` and is
  carried into `## Prior decisions` as it is answered. Nothing launches a writer for
  `proposal.md`, the approach or the checklist. **The one surviving writer is
  `write-spec`'s fan-out**: when a change has extra capabilities **and** the session's
  context cannot hold their research, one `spec-writer` per extra `spec-<capability>.md`,
  given `proposal.md` by path — its decision log, vocabulary and primary contract are the
  shared ground a brief used to carry — writing that one file and nothing else. Context is
  the only reason to fan out
- **the change document is read by one reviewer that never saw the conversation, before
  the stop.** Once the checklist is written, `review-spec` launches one fresh
  `spec-reviewer`, given by path the change document, every `spec-<capability>.md` beside
  it and the capability specs their `Targets:` name. It runs `spec-drift --trace` and
  `--ears`, reads the whole document and writes nothing. It reads **the contract** —
  ambiguity, testability, `Governs:`, cross-capability contradiction — and **the
  checklist, box by box**: what the contract and the approach together fail to answer for
  building that box, naming the gap and the part that should have answered it, never
  proposing the answer. The session that argued its way to an approach cannot distinguish
  what it wrote down from what it merely decided; the fresh eye moved from drafting to
  reading, and the per-box question is the only check this flow has that a document says
  enough to be built from
- **both layouts are read, by one rule.** Every payload reader that counts boxes uses the
  CLI's definition: a line matching `^\s*[-*]\s+\[( |x|X)\]`, none inside a fence; in
  `proposal.md` only boxes in a section opened by `## Checklist` and closed by the next
  `^##\s` heading, in `tasks.md` and `plan.md` every box outside fences. A change's
  checklist is `proposal.md` when its `## Checklist` holds a box, otherwise `tasks.md`,
  otherwise `plan.md` — never decided by a file existing, because every change in both
  layouts has a `proposal.md`. A queued capture has no `## Checklist`; the `Queued:` scan
  does not change. Nothing in flight is migrated: a change started in the old layout
  finishes in it, and each old-layout reading says when it can go
- **a finished change that never landed is reported, not walked past.** Phase 1 runs two
  scans over every change's checklist, open boxes and closed ones, so a change present in
  the closed scan and absent from the open one has every box closed and its folder still
  on disk — its landing did not finish. That signal sat in
  the output unnamed while two changes shipped half-landed and phase 1 walked over the
  state twice. It is written as a table of three outcomes rather than a paragraph, so the
  empty open scan is a case with a name instead of an absence somebody has to notice, and
  `/libretto-status` delegates to the same scan rather than re-deriving it. **A report,
  never a gate**: `--retired` fires on a deletion, and the failure here is that no deletion
  happened — there is no commit to refuse
- **an approach cannot be deleted taking its reasoning with it.** The landing commit
  removes the change folder, and with it the alternatives table whose whole argument is
  that a diff shows what was built and nothing shows what was not built and why. `spec-drift
  --retired` — inside `--anchors` — fails that commit unless a capability spec's *Prior
  decisions* section moved in it. It fires on a deleted `proposal.md` whose `HEAD` copy
  carries `Targets:` outside fences, and on a deleted old-layout `plan.md`; a change deleting
  both is reported once. **`Targets:` is the only thing that makes a `proposal.md` a
  contract** — the same discriminator `libretto land` uses, so the two landing tools cannot
  disagree — and deleting a queued or trivial-lane capture, which carries none, fails
  nothing. **That section, never the file**: the delta lands in the same commit by
  definition, so "a spec was edited" is green on every landing and measures nothing. The
  escape is a declaration in the document being deleted, `Durable decisions: none`, and not
  a flag — a flag is typed by whoever wants the commit through, a line in the document is
  written by the person who knew. **What no mechanism can stop is that line becoming a
  reflex**, and `review-work` is what reads a full alternatives table sitting beside a
  `none`
- **marking a box is not the contract moving.** The change document carries `Targets:`, so
  default `spec-drift` would read every box-marking commit as the capability's spec moving
  and say nothing about the code staged beside it. A staged `proposal.md` counts as its
  targets moving only when something outside its `## Checklist` section differs between
  `HEAD` and the index
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
- **two promises on one proof are named, never failed.** A criterion carrying two or more
  `shall` on exactly one `Proof:` line, naming no ceiling, is *joined*: one test rarely
  holds two promises, and a criterion whose citation exercises half its clause is this
  flow's most recurrent miss — eight entries in the lessons ledger. `spec-drift --trace`
  lists every joined criterion, in capability specs and change documents, with its file and
  a count in its summary line; `--anchors` prints one line counting them per scope and
  decides its exit status exactly as before, and with nothing joined its output is byte for
  byte what it was. The author answers by splitting the criterion, citing a second proof,
  or naming its ceiling in it — and the line goes away. **Structural, not semantic**: it
  reads the criterion and counts its citations, never the cited test, so it runs before the
  test exists — and it cannot see one `shall` promising a list whose proof asserts part of
  it, which stays a reading for `work-reviewer`
- **one stop before code, on the whole document.** The contract lane stops exactly once
  before code: after the fresh reviewer returns, on contract, approach and checklist
  together, with the reviewer's findings relayed unedited. `review-spec` owns it and asks
  it as a native question — carry on to phase 6 (recommended, naming the first box), change
  the document first, or stop here, with room to answer differently. Contract, approach and
  order are argued over one document, so they are agreed over one; a second stop inside
  the work would be a round trip for the same reading. The count is two, one of them inside
  the work
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
- **the same flow with its two stops — the one before code and the push — answered in
  advance** — `libretto-attacca`, the
  score's instruction to go on without pausing, running phases 1 to 8 and ending at a
  pushed branch with a request open on it. It writes everything the attended flow writes;
  what the invocation answers is the waiting. Every question it would have asked becomes an
  `(assumed)` entry in the change document's `## Decision log`, an assumed entry under
  `## Prior decisions`, and a line in the report and the request; because the run authors
  the document, it acts on the fresh reviewer's findings before phase 6 and names each one. **What it cannot answer is a gate**, and a
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
  phases 1, 2 and 5, up to the stop, and nowhere after
- **ceremony proportional to the change** — one stop before code on the contract lane,
  none on the trivial lane, and phase 8's question in both
- **the trivial lane's capture does not outlive it.** Phase 1 writes `proposal.md` for every
  source-3 request before phase 2's step 0 decides the lane, so a trivial-lane change always
  has a capture on its branch; the trivial lane deletes that folder in its own commit. The
  directory is created, and it is gone by the time the work is committed
- **no borrowed method's name in what installs.** No file under `skills/`, `agents/` or
  `commands/` names the method this flow's 2026-10-21 shape was inspired by, as an acronym
  or spelled out. The flow is *the Libretto flow*, its lanes the contract lane and the
  trivial lane; `docs/FLOW.md` credits gentle-ai in one line and nowhere else. A spec
  quoting a decision verbatim may contain the acronym
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
- **implicit creation of a `changes/` directory in a project that has none.** A
  staging area nobody empties is a second source of truth. Bootstrap has exactly
  one owner — explicit, user-invoked `libretto init` — and a skill or command that
  finds the directory missing still refuses, pointing at the owner where `libretto`
  is on PATH and refusing identically where it is absent.
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
- **a reviewer writing.** `work-reviewer` and `spec-reviewer` hold the same grant — `Read`,
  `Grep`, `Glob`, `Bash`, `Skill` — and nothing that edits: a reviewer that edits what it
  reviews is the author twice. `Skill` is there because its
  first instruction is to invoke `evidence`, `Bash` because running the proofs is its
  job. The repair happens in `review-work`, on the other side of the seam: an agent that
  fixes what it found is grading its own repair.
- **a reader that can run what it judges.** `converge-reader` holds less than either
  reviewer — `Read`, `Grep`, `Glob`, `Skill`, and no `Bash` — because its verdict decides
  whether a box closes: a reader that cannot run a proof cannot report one as passed, so
  its word is never mistaken for evidence. The session runs every proof it closes on.
- **a second fresh reader before the stop.** One `spec-reviewer`, not a panel.
- **a contract block per capability inside `proposal.md`.** Extra capabilities go in
  `spec-<capability>.md` files beside it.
- **migrating a change in flight.** No change folder is rewritten into the new layout; a
  change holding `spec.md`, `plan.md` or `tasks.md` finishes in the layout it started in,
  and what it still lacks the session writes inline in that layout.
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
- **the Go binary *defining* the queue.** find-work's queue section is the definition;
  `libretto queue` and the wiki implement it, and where they disagree the binary is the
  bug.
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
- **searching a cited test's body for the criterion's terms**, in any form, and opening,
  parsing or running a cited test to judge a criterion. Rejected on a measurement — see
  *Prior decisions*.
- **failing on a joined criterion**, in a change document or a capability spec. It is a
  count in `--anchors` and a list in `--trace`.
- **telling an enumerated response from a sentence.** One `shall` promising a list whose
  proof asserts part of it is six of the ledger's eight half-clause cases, and recognising
  an enumeration is a reading; `work-reviewer` caught all eight and stays the mechanism.
- **converge rewriting or annotating a drifted box.** Reported only, its text and state
  left as they are; resolving drift is the user's, at the contract.
- **converge editing the contract or the approach** — any pillar, `## Approach`, the
  decision log, `Targets:`, `Durable decisions:`. It writes to the checklist and nothing
  else.
- **converge unchecking or re-running a closed box.** A closed box whose proof fails today
  is a regression, and `libretto metrics` already counts a reopened box as its signal.
- **converge pushing, creating a branch, or staging anything but the checklist file.**
- **converge guessing a box-to-criterion mapping.** A box whose `Closes:` matches no
  criterion unambiguously, or that carries none, is reported `unmapped`, never closed.
- **a targeted capability spec's criteria as converge's uncovered work**, and **a diff
  range as the reader's input.** Uncovered means the change's own criteria; the reader
  reads the tree as it is.
- **converge boxing a contract gap** — an outcome or scope item no criterion covers. It is
  reported; writing the criterion is editing the contract.
- **a `libretto converge` subcommand**, `loop`'s stall message naming converge, and a
  `recommend.go` row for `converge-reader`. Each is a `cli` follow-up.
- **routing to converge from `libretto-flow` or `libretto-attacca`**, and a line for it in
  `docs/FLOW.md`. A door, not a phase.
- **a second definition**, inside converge, of a box, of which file holds the checklist,
  or of how a citation is run. `find-work`, `review-work` and `work-reviewer` own those;
  converge points to them.

### A stop is a place where the user changes something

That is the whole test, and it is what the count is derived from.

| After | Stops | What is being changed |
|---|---|---|
| 1 · find-work | no | — the reading is stated, and the spec is where it gets corrected |
| 2–3 · write-spec | no | — the contract is written, and agreed at the one stop |
| 5 · write-plan, 5→6 · write-tasks | no | — the approach and the checklist join the same document |
| before 6 · review-spec | **yes** | the whole change document: contract, approach, order, and what waits on what — with the fresh reviewer's findings on the table |
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
second reading of this same list rather than a feature of its own: the two stops are
answered by the invocation, phase 1's in-flight choice resolves to the oldest change with
open boxes, and everything not in the table is untouched. **A gate is not a stop and
cannot be answered** — a failing one still stops the run, two on one task still stop the
task, and a missing credential still stops the phase, because that is the input failing to
arrive and no invocation supplies it. A mode that answered a gate would not be unattended,
it would be unverified.

When phase 2 answers **no spec needed** the one stop before code collapses with it — there
is no document to agree — and phases 6, 7 and 8 run in one turn. **Phase 8's question survives
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

### Asking is bounded to before the stop

Phases 1, 2 and 5, up to the one stop, where nothing has been built on the answer yet —
which is why the stop sits exactly there. **After the stop, an unsettled question becomes
a finding**: it
reaches the phase 7 report with what was assumed in the meantime and what changes if the
assumption is wrong, and the user meets it at phase 8 with everything else.

**Phase 2 interviews before it writes the spec** — one question per native-prompt call,
the recommendation carrying its reason, "no more questions" always an option, a soft
bound around five that judgment crosses in either direction — so the contract is built
by two people rather than handed over finished, and each answer can redirect the next
question, which a batched call never could. **Phase 5 asks exactly one: the approach**,
chosen from two or three with named tradeoffs. Every answer lands verbatim in the
change document's `## Decision log`, which is also what a fan-out writer reads. **Zero is
a legitimate answer, reported in one line:** a quota manufactures questions the code
already answers, which is the rubber-stamp round trip removed from three other phases
arriving back through the door marked collaboration. The stop covers the approach and the
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
| a change already in flight | unchecked boxes in a change's checklist — `proposal.md`'s `## Checklist`, else `tasks.md`, else `plan.md` for a change created before either rename |
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

**Queued is three conditions, and find-work states all three**: a `Queued:` line at column
0, no unmerged local branch carrying the change's name, and no box under the proposal's
`## Checklist` — a box means the work started. Ties on one `Queued:` date go by change
name, ascending, so `/libretto-next`'s recommended pick and the first line of
`libretto queue` are the same change. The binary and the wiki implement that section;
they never extend it, and where the two disagree the section is the rule. Neither
`/libretto-status` nor `/libretto-next` calls the binary — they delegate to find-work and
keep working on a machine with the skills and no binary.

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

**A branch is work in flight too.** Scanning checklists cannot see the trivial lane: a
change that needed no spec has no checklist to scan, and its capture is deleted in its own
commit, so it exists only as commits on a branch. The lane's first real run produced exactly that, and the next phase 1
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

**A marked box that was never committed was never marked.** The checklist ships in the same
commit as the task that closed it, exactly as the spec does — or, when `/libretto-converge`
closes a box whose code had already landed, in its own `docs(<change>): converge the
checklist` commit holding the checklist file alone. Either way it is committed, never left
for the landing to find. A `git add` scoped to the code
leaves it behind in the working tree, and the change that lands then deletes the folder with
every unrecorded mark still in it — measured on the run that added this line: six commits,
0/24 boxes, work finished. Phase 1 reads unchecked boxes to decide what is in flight, so
the failure is silent until the checklist disappears.

**Every file has exactly one author.** The change document has one writer, the session;
each `spec-<capability>.md` has one, the session or a fan-out writer; the fresh reviewer
writes nothing, `converge-reader` writes nothing, and `build-and-check` and any subagent
report a finished box for the session to mark. This is what makes concurrency safe without locks.

**Sub-agents start with no rules loaded.** Every launch carries them explicitly, or the
default takes over — and the default is to guess plausibly.

**An agent that ships is named by the skill that launches it.** `review-spec` launches
`spec-reviewer`, named in it; `converge` launches `converge-reader`, exactly one, named in
it; `write-spec`'s fan-out names `spec-writer`, which nothing
outside that section and the agent file may name. An agent nothing references is dead payload —
`spec-writer` sat unreferenced from its landing until 2026-08-10, and no gate noticed.

**The parsers fix the change document's heading levels.** `spec-drift`'s `criteria_of`,
`criteria_without_proof` and `section_of` open and close a section on `^##[[:space:]]`
only, so `## Verification criteria`, `## Prior decisions` and `## Checklist` are level-2
headings, `## Checklist` its own section — a box under the criteria with no level-2 heading
between them is read as a criterion and fails `--ears`. **Every line starting `Proof:` or
`Targets:` outside a fence counts**, in every `.md` under the change roots: the approach
and the checklist start a line with `Proof:` only to cite a real test, `proposal.md`
writes exactly one `Targets:`, for the primary capability, and each `spec-<capability>.md`
one, for its own.

**A test-named citation's exit status alone is not a pass.** A named test that does not
exist passes under `go test -run`: measured 2026-10-22, `go test ./cmd/libretto/ -run
'^TestDoesNotExistAnywhere$' -v` exits 0 printing `PASS` and `[no tests to run]`, and no
`--- PASS:` line. So `converge` runs `spec-drift --anchors` before any proof — reached from
the skill as `<skill-base>/../record-work/spec-drift`, the form `write-spec` uses, never
through `scripts/` — and a test-named citation passes only on its `--- PASS: <TestName>`
line in `-v` output. A file-only citation passes on exit 0 with no `FAIL` line.

**converge finds the checklist and counts boxes by `find-work`'s rule**, the same answer
`checklistPath` and `readPlan` give in `cmd/libretto/loop.go`, and restates none of it.
**Its dirty tree is tracked modifications only** — `git status --porcelain
--untracked-files=no` listing anything — because this tree carries an untracked generated
`.agents/specs/wiki.html` on `main`, and stopping on it would stop converge on every run.

**Old-layout readers stay until nothing old is in flight**, and every place that reads an
old-layout name says when that reading can go — a `ponytail:` comment in code, a sentence
in prose — the precedent the `plan.md` → `tasks.md` rename set on 2026-08-17.

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
- **A `Durable decisions:` line is a claim about whether the list is empty, never the
  list.** The list is the change document's own `## Prior decisions` — the change's
  `spec.md` *Prior decisions* in the old layout, where the line sits in `plan.md`. Both readings looked
  right and the tree held the contradiction to prove it: one plan declared "the two" over
  a section of three, another carried no line at all. Found by a 5→6 cutter reading both
  documents cold.
- **`plan.md` was reused for the technical approach rather than adding a `design.md`.**
  The complaint was that the file named `plan` is not a plan; a `design.md` beside an
  unchanged `plan.md` leaves that exact file unchanged and the complaint standing.
  Answered by the user, 2026-08-17. Since 2026-10-21 a new change writes the approach as
  `proposal.md`'s `## Approach` section; the old-layout `plan.md` keeps this meaning.
- **The checklist is written inline in the 5→6 seam, which stays unnumbered.** Until
  2026-10-21 a fresh `task-cutter` cut it — *the task cut is a subagent, not a numbered
  phase* — and its second output, what the documents failed to answer, moved to the fresh
  `spec-reviewer`, reported per box. Renumbering everything that says eight still buys
  nothing, the answer the 6→7 review seam already gave for itself.
- **Phase 5 asks, since 2026-08-19 — retiring the 2026-08-12 decision that it asked
  nothing.** The user reversed it for the fork alone: the approach is the one phase-5
  decision that is genuinely the user's, and a finished decision can only be
  rubber-stamped while a fork with named costs must be decided. The no-third-stop
  argument the original decision was made for survives: the fork rides inside the phase.
- **Four answers from the user, 2026-08-19, verbatim options, as they stand after
  2026-10-21:** the interview bound is *"Blando ~5, con juicio"* — never a hard cap,
  untouched. *"Brief siempre"* is **retired** for every new change, fan-out included: no
  `brief.md` is written in the new layout, and a fan-out writer is given `proposal.md` by
  path instead. *"Orquestador parchea inline"* survives **for the fan-out only**, where a
  writer can still leave a `[NEEDS CLARIFICATION]` marker for the session to resolve,
  scoped to the bracket expression, the one declared exception to one-author-per-file; a
  single-capability change has no writer and so no marker. *"Sí, marcadas 'assumed'"*
  survives in a new place: assumed answers land in the change document's `## Decision log`,
  one home for decisions in both modes.
- **Deliberately not built, with the conditions that bring them back:**
  section-by-section document validation, when change documents outgrow one read; markers
  in the approach — the fresh reviewer reports gaps per box instead; a hard cap of five, if
  the soft bound degenerates into interrogation.
- **The specs wiki regenerates at the landing step, not on a schedule and not by a
  daemon.** The landing is the moment the specification moves, so the refreshed index
  rides the same commit as the delta that changed it; a generated view left behind by a
  landing is drift wearing a marker comment. Assumed 2026-08-18 rather than asked — if
  manual-only regeneration is wanted, this instruction is the part to drop.
- **The retirement gate compares the *section*, never the file.** Requiring any edit to a
  capability spec passes on the delta application alone — that happens in the landing
  commit by definition — so the gate would be green on every landing and measure nothing.
  It is the version that looks like a gate and is not one.
- **Its escape is a declaration in the change document, not a flag.** A flag is typed by
  whoever wants the commit through, at the moment they want it to stop complaining; a line
  in the document is written while the approach is, by the person who knew. No mechanism
  can stop that line becoming a reflex — `review-work` reads the approach and the diff
  together.
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
  usually the only one after the stop: a question written as a sentence is a paragraph
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
- **One change document, one stop, one fresh reviewer** — answered by the user,
  2026-10-21, from `adopt-organic-driven-development`, verbatim from its decision log:
  - Q: for substantial work, what replaces the proposal → brief → spec → plan → tasks chain? → A: Un doc de cambio + specs vivos (Recommended) — one document per change (intent, EARS criteria, approach, checklist); capability specs with Governs:/Proof: and the landing stay
  - Q: with one change document, where does the flow stop for the go-ahead before the push? → A: Un stop: el doc entero (Recommended) — contract, approach and checklist reviewed together, once, before the first line of code; then only the push question
  - Q: who writes the change document, and where does the fresh eye stay? → A: Sesión escribe, revisor fresco lee (Recommended) — the session writes the document inline; one fresh reviewer (review-spec) reads it before the stop; spec-writer, plan-writer and task-cutter launches go from three to one review
  - Q: what happens to changes and tooling that read proposal/plan/tasks today? → A: Leer ambos formatos (Recommended) — new changes are born as one document; loop/metrics/wiki/find-work keep reading the old layout for what already exists; nothing in flight is touched
  - Q: is the change document the same file growing from the capture, or a new one? → A: proposal.md crece (Recommended) — the queue/find-work capture is the document; criteria, approach and checklist are added to the same file; Queued: and proposal.md scans do not change
  - Q: what is this called in Libretto? → A: Sin sigla, vocabulario propio (Recommended) — no "ODD" in the payload; the flow stays "the Libretto flow", lanes named in its own vocabulary; FLOW.md credits gentle-ai as inspiration in one line
  - Q: without task-cutter, who checks that the document is enough to build every box? → A: El revisor fresco lo hereda (Recommended) — the one fresh reviewer reads contract and checklist and reports, per box, what the document fails to answer
  - Q: how does a change touching two or more capabilities carry its contract? → A: proposal.md + spec-<cap>.md extra (Recommended) — proposal.md carries the primary capability contract; each extra capability gets a spec-<cap>.md beside it with its own Targets:
  - Q: when an extra spec-<cap>.md needs its own writer, who writes it and what does it read? → A: spec-writer, lee proposal.md (Recommended) — agents/spec-writer.md survives for this case only, given proposal.md by path; no brief.md

  **Reversed by it, each named:** *two stops with a contract, none without, and phase 8's
  question in both* → one stop and phase 8's question; the stop rows `2–3 · write-spec` and
  `5→6 · write-tasks` → one row, `before 6 · review-spec`; *the stop moved rather than
  multiplied … the count stays at three* → two; *the contract and the plan are drafted by
  agents that never saw the conversation* and *the spec's author is `spec-writer`, single
  case included* → the session drafts inline and the fresh eye reads; *the checklist is cut
  by an agent that never saw the design* → cut inline, its gaps question moved to the
  reviewer; *a plan and a checklist, two files* → two sections of one file; *every file has
  one author … the brief, the plan* → the document has one, the session; *`ask` holds at the
  5→6 seam* → up to the stop; attacca's *three stops* → two; *`review-spec` runs in the 3→5
  seam* → after the checklist, before the stop, as a fresh subagent; attacca's assumptions
  *in `decisions.md`, where the writers read* → in `## Decision log`.
- **Derived from those answers, 2026-10-21, by the orchestrator — logged, not user
  answers:**
  - the checklist is only the boxes under a level-2 `## Checklist`: settled by the code,
    since spec-drift's parsers react to `## ` only and boxes elsewhere would be read as
    criteria by `--ears`
  - `spec-writer` keeps its `libretto models` recommendation (sonnet/high); only
    `plan-writer` and `task-cutter` keys are forbidden
  - a deleted `proposal.md` is a contract for `--retired` and `land` by a `Targets:` line
    outside fences and nothing else. *`## Verification criteria` as a second arm* lost: it
    let the two landing tools give opposite verdicts on one commit. **Ceiling named:** a
    contract-lane document written without `Targets:` is invisible to the gate — the blind
    spot an old-layout change has when its `plan.md` was never written; the fix that day is
    the fresh reviewer reporting a contract with no `Targets:`
  - `find-work` counts boxes as the CLI does: `planLine`'s shape, fences excluded, `^##\s`
    section boundaries, in every reader
  - `proposal.md` wins as the checklist only by holding a box under `## Checklist`. *Its
    existence selects it* lost: every change in both layouts has one, so existence would
    hide every old change's `tasks.md`/`plan.md`
  - the decision log is a section of the change document and a new change writes no brief:
    `brief.md` and `decisions.md` existed to feed writers that never had the conversation,
    and with the session writing inline a second file would contradict *one document*
  - `review-spec` owns the stop: it is the last thing before code, and the findings have to
    be on the table when the agreement is made. Its opening sections are the session's
    conduct (launch, relay, ask); the reviewer's instructions start at the reading
  - the fresh reviewer is an agent file, `agents/spec-reviewer.md`, not a generic subagent:
    every fresh seat here has one, and an agent that ships is named by the skill that
    launches it. `Bash` because `spec-drift --trace` and `--ears` are its first step,
    `Skill` because its first instruction is `evidence`
  - `write-spec`, `write-plan` and `write-tasks` stay as skills, each writing its part of
    the change document against the rules it already held; the phase numbers stay
  - the acronym rule's scope is `skills/`, `agents/` and `commands/`, the files that
    install; specs quoting a decision verbatim may contain it
- **The bootstrap non-goal was given an owner, not deleted.** Retired from
  `add-init-command`, assumed 2026-08-21: implicit creation of `changes/` by a skill stays
  forbidden, and explicit, user-invoked `libretto init` (the `cli` capability) is the one
  sanctioned bootstrap. *Deleting the non-goal* lost — it would have let any skill invent
  a staging area nobody empties. The two refusal sites, `write-spec` and `libretto-queue`,
  each gained one guarded sentence in record-work's grammar — where `libretto` is on PATH,
  name it; where it is absent, refuse unchanged — and neither restates the rule nor runs
  the command. No third site gains a pointer: no other payload document creates or
  declines to create `.agents/changes/`. If wrong: each pointer is one sentence to remove.
- **The queue-blind non-goal was reversed, not deleted silently.** Retired from
  `add-queue-command`, 2026-10-22, attacca run. "The Go binary knowing about the queue"
  became "the Go binary *defining* the queue": find-work's queue section is the
  definition, `libretto queue` and the wiki implement it, and where they disagree the
  binary is the bug. The original was already broken by the wiki's flow board, which read
  `Queued:` lines from the day it landed and nobody reversed the sentence; this says so.
- **A checklist box disqualifies a `Queued:` proposal.** Retired from `add-queue-command`,
  assumed 2026-10-22 on review-spec's finding. Queued is three conditions — the `Queued:`
  line at column 0, no unmerged branch taking the name, no box under `## Checklist` —
  stated in find-work and applied by `queuedIdeas` (`cli`). If wrong, the clause leaves
  both together.
- **Ties on one `Queued:` date go by change name, ascending.** Retired from
  `add-queue-command`, assumed 2026-10-22. *Capture-commit order* lost: a `git log` per
  proposal, against a design that put the date in the file to avoid exactly that. If
  wrong, the sentence and the binary change together.
- **`/libretto-status` and `/libretto-next` never call `libretto queue`, and find-work
  never runs it in place of its own scan.** Retired from `add-queue-command`, assumed
  2026-10-22. *Status calling the binary where it is on PATH* lost: two paths to "what is
  queued" is the thing find-work's single-owner rule exists to prevent, and status would
  behave differently on a skills-only machine. If wrong, status calls the binary where it
  is on PATH, with the skill's scan as the fallback.
- **A cited test's body is not searched for the criterion's terms — rejected on a
  measurement, so do not re-propose it without a new one.** Retired from
  `check-clause-coverage-in-anchors`, 2026-10-22, attacca run. The ask was exactly this
  heuristic: the EARS trigger clause's key terms (backticked spans, quoted strings,
  `--flags`) looked up in the body of the cited Go test. A throwaway prototype ran it over
  every criterion in `.agents/specs/*/spec.md` — 706 criteria, 719 citations (637 Go tests,
  61 `check-payload`, 16 `spec-drift`, 5 other). The trigger-clause variant flagged 15 of
  the 52 Go-proved criteria carrying a term (2.1% of all 706); the whole-criterion variant
  flagged 106 of 304 (15.0%). **Ten flags hand-checked from each were twenty false
  positives**: terms reach a test through helpers, constants, fixtures and `filepath.Join`
  segments, a negated trigger is *correctly* absent from the body, and a `check-payload`
  citation has no body at all — where two of the eight ledger cases live. Run against each
  of the ledger's six Go-proved cases at its fix commit's parent, the trigger variant
  caught none for the right reason and the whole-criterion variant one. And of the eight
  ledger cases only one missed the trigger; six missed an item of an enumerated response,
  which a trigger heuristic cannot see. **What would reopen it:** a resolver that follows
  a test into its helpers — a Go call graph, a dependency this bash script refuses — or a
  new measurement beating these numbers. If it ever ships, it ships behind its own flag,
  warn-only, never inside `--anchors`, its false-positive rate in the help line.
- **A criterion is joined when it carries two or more `shall`, exactly one `Proof:` line,
  and no occurrence of *ceiling*.** Retired from `check-clause-coverage-in-anchors`,
  assumed 2026-10-22. Structural on purpose: it reads the document and never guesses about
  the test, so it runs before the test exists — at `review-spec`, where splitting a
  criterion costs a line. `shall` is counted the way `is_ears` matches it, so the two
  checks never disagree about what a `shall` is. A criterion with no `Proof:` is not
  listed: it is already `--trace`'s *unproven criterion*. If wrong, the predicate is one
  function, `is_joined`.
- **Joined criteria warn and never fail, in change documents and capability specs alike.**
  Retired from `check-clause-coverage-in-anchors`, assumed 2026-10-22. Measured: 72 of the
  706 capability criteria (10.2%) and 36 of the 99 criteria in the ten landed EARS-era
  deltas history recovers are joined; 6 of 26 hand-checked flags were real gaps (23%). A
  gate failing one criterion in three, wrong three times in four, gets muted. *Failing on
  change documents, like `--ears`* lost on those numbers. **Ceiling named:** of the
  ledger's eight entries it saw one — wiki-flow-board, whose two joined criteria were both
  real misses. Three predate EARS and carry no `shall`; four rest one `shall` on an
  enumerated response, invisible to it. If wrong, `joined` returns 1 on change documents
  only.
- **`--anchors` carries a count, `--trace` carries the list.** Retired from
  `check-clause-coverage-in-anchors`, assumed 2026-10-22. The gate's output is read on
  every commit; one line there costs nothing, and 72 per-criterion lines would be the
  warning that gets muted. A separate `--joined` flag lost: the list already has a home in
  the map `review-spec`'s reader runs first, and a flag nobody runs finds nothing.
- **The escape is a declaration in the criterion — split it, a second `Proof:` line, or
  the word *ceiling*, by house phrase "Ceiling named:" — never a flag.** Retired from
  `check-clause-coverage-in-anchors`, assumed 2026-10-22, on the precedent of `Durable
  decisions: none`: written by whoever knew, at the moment they knew. The accepted risk is
  that *ceiling* mentioned for another reason silences a real joined criterion — an author
  who wrote the word was thinking about what the proof misses.
- **converge is a payload command and skill launching one fresh reader; no binary
  subcommand.** Retired from `add-converge-command`, assumed 2026-10-22, attacca run.
  Whether a box is done is a reading, and the binary cannot read prose — `loop.go` refuses
  to parse a task's text for that reason. *A `libretto converge` subcommand spawning a
  session as `loop` does* lost: `loop` sits outside the session because each box needs a
  fresh context, and converge is one reading. *A mechanical subcommand running each box's
  proofs* lost: it would parse `Closes:` prose, and would have to know how to run a
  citation in a project that is not Go. **What would reopen it:** wanting converge
  headless — a subcommand that spawns `/libretto-converge` exactly as `loop` spawns
  `/libretto-attacca`.
- **A test-named citation passes only on its `--- PASS:` line, and `spec-drift --anchors`
  runs before any proof.** Retired from `add-converge-command`, 2026-10-22, on a
  measurement: `go test -run` on a name that does not exist exits 0 and prints `PASS`, so
  an exit status would close a box on a test nobody wrote. Either guard alone misses a
  case — anchors does not run the test, and the pass rule is the session reading output.
  A file-only citation passes on exit 0 with no `FAIL` line. If another runner prints a
  different pass line, the skill names the line per runner, still never trusting an exit
  status alone for a named test.
- **The `converge-reader` has no `Bash`, and a reader's verdict never closes a box.**
  Retired from `add-converge-command`, assumed 2026-10-22, attacca run. A reader that
  cannot run a proof cannot report one as passed; the session closes a box only on runs it
  read itself, citations taken from the contract rather than the reader's quote. *The
  reader runs the proofs* lost: the checklist's one writer must be the one that observed
  the pass. If wrong, the reader gains `Bash` and its runs still close nothing.
- **Drift is reported, never rewritten.** Retired from `add-converge-command`, assumed
  2026-10-22, attacca run. Rewriting a box to match the code ratifies the drift without
  the contract moving — an edit to what the change promises, made by the back door. The
  accepted risk is that the next `loop` session builds a drifted box as written; the
  report names every drifted box. If wrong, an indented `Drift:` line beneath the box, the
  box itself untouched.
- **converge commits the checklist file alone — reversed from "converge does not
  commit".** Retired from `add-converge-command`, assumed 2026-10-22 after review-spec,
  reversing the attacca run's first answer. A mark never committed was never marked, so
  converge ends with one `docs(<change>): converge the checklist` commit staging nothing
  else when it marked or appended anything, and no commit when it did neither; no push, no
  branch. It stops on tracked modifications, because a proof against unrecorded code proves
  nothing about what was recorded and a staged file would ride the commit; untracked files
  do not stop it. If wrong — the user wants to read the edit first — the checklist is left
  staged and the report says the marks are unrecorded.
- **A box maps to criteria only through an unambiguous `Closes:` match, and uncovered means
  the change's own criteria.** Retired from `add-converge-command`, assumed 2026-10-22
  after review-spec. Anything else is `unmapped`, open and reported — a guessed mapping is
  the reading a verdict may not stand in for. A targeted capability spec's criteria belong
  to changes already landed; boxing them would turn converge into a regression sweep. A
  box is appended only for an uncovered criterion whose proof does not pass, because a
  criterion already met is not a task. If wrong, an exact quote from the reader is
  accepted as a mapping.

- **An unclosed fence fails `--anchors`.** Retired from `gate-recurring-defect-classes`,
  decided 2026-10-23 on the queue run's evidence: a criterion that wrote three backticks in
  prose opened a fence `defenced` never closed, and 157 `Proof:` citations after it went
  unchecked while the gate stayed green. An odd count is always a defect, so it fails rather
  than warns; with every fence closed the output does not move by a byte. The fence pattern
  lives once, in `FENCE`, read by `defenced` and `fences`. Lost: a warning (a gate that lies is
  still green) and fixing only the instance (the class came back).
- **Nothing in the payload's scripts pipes into `rg -q`.** Retired from
  `gate-recurring-defect-classes`, decided 2026-10-23: under `set -o pipefail`, `rg -q` stops
  reading at its first match, the writer before it dies of SIGPIPE, and the pipeline's 141
  reads as no match — it made `--retired` fail open on a long change document (fixed v0.37.0).
  Measured: 40,000 lines through `| rg -q` exit 141, through `rg … >/dev/null` exit 0, and a
  newline-squashed single line exits 0 either way, since the reader must finish the line —
  which is why `check-payload`'s marker search, rewritten here, was never actually vulnerable.
  Ceiling named: a pipe split across lines, or a `|` inside a quoted argument, is not seen.
- **A decided bump rides `gh pr create --label`.** Retired from
  `gate-recurring-defect-classes`, decided 2026-10-23: three requests opened unlabelled with
  the bump already authorised failed the label check's `opened` run and mailed the user a red
  check for an answered question. attacca with no bump decided keeps its designed red check,
  which is the question.

- **`Governs:` can exclude a path with a `!` entry.** Retired from
  `let-governs-exclude-a-path`, decided 2026-10-23: four test files under `cmd/libretto/**`
  prove ci's and readme's promises, and a glob list with no subtraction could only keep
  claiming them, so drift warnings pointed at cli for other capabilities' tests — 6 of 10
  hand-checked warned landings in report-drift-trend-in-metrics were overlaps of this kind.
  One rule, `governed_by`, decides for every reader: a plain entry claims, a `!` entry gives
  back, the exclusion wins, and entries are one capability's, never pooled. An exclusion
  matching nothing is a dead claim. Lost: a separate `Excludes:` key (a second line to keep
  beside the first) and listing cli's files one by one (a glob nobody maintains).

## Task breakdown

- [x] `find-work` — phase 1, three sources
- [x] `libretto-status` — read-only, what is in flight
- [x] `libretto-converge`, `converge` and `converge-reader` — re-derive a drifted
      checklist, closing only on proof the session observed and committing the checklist
      alone
- [x] `write-spec` — phases 0, 2, 3, including the fan-out
- [x] `write-plan` — phase 5
- [x] `build-and-check` — phase 6
- [x] `present-work` — phase 7
- [x] `record-work` — phase 8, including consolidation
- [x] `evidence` — the standing rules
- [x] `libretto-flow` — the routing command
- [x] `libretto-queue` and `libretto-next` — capture the queue, drain it one at a time,
      with `find-work` owning the scan both read
- [x] every stop asked with `AskUserQuestion` — the stop before code and phase 8, plus
      phase 1's in-flight choice, with the rule stated in each skill that owns one
- [x] `libretto-attacca` — the stops answered by the invocation, the classification
      in one file, and the five stop-owning skills each stating what happens to their own
- [x] vendored delegates with attribution
- [x] ponytail, ponytail-debt, caveman and caveman-commit vendored, callers' prose
      reconciled, THIRD-PARTY.md and docs recording the reversal
- [x] `spec-drift` and `check-payload`
- [x] one change document, one stop, one fresh reviewer — `spec-reviewer` added,
      `plan-writer` and `task-cutter` retired, both layouts read
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
- [x] the bootstrap non-goal given its owner — `libretto init` named, guarded, at both
      refusal sites
- [x] the retro wiring guarded in `check-payload`, observed failing on a hand-broken
      copy before landing

## Verification criteria

- **When** a change folder's checklist — `proposal.md`'s `## Checklist`, else `tasks.md`,
  else `plan.md` — holds at least one box and none of those boxes open, phase 1 **shall**
  report that change as finished and not landed, naming it — and the status command
  **shall** carry the same report by delegating to that scan. **Where** the folder holds no
  checklist box, it **shall not** be reported that way: a captured idea is not an unlanded
  change. The plan skill **shall** state that a
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
  plumbing, and the plumbing is where a gate goes silent — which reads as green. The three
  old-layout cases kept their expected exit codes, unedited, when the change document
  arrived.
  Proof: skills/record-work/spec-drift --self-test
- **If** a staged commit deletes a `proposal.md` whose `HEAD` copy carries a `Targets:` line
  outside fences, and no capability spec's *Prior decisions* section differs between `HEAD`
  and the index, and the deleted copy carries no `Durable decisions: none` line, **then**
  `spec-drift --retired` **shall** exit 1 naming the change. **Where** the deleted copy
  declares `Durable decisions: none`, or a *Prior decisions* section moved, it **shall**
  exit 0.
  Proof: skills/record-work/spec-drift --self-test
- **When** a staged commit deletes a `proposal.md` whose `HEAD` copy carries a `Queued:` line
  and no `Targets:` line outside fences, `spec-drift --retired` **shall** exit 0 and **shall
  not** report that change as retired. The case broken first: removing `Targets:` from the
  contract-lane fixture turns its exit 1 into 0, which proves the discriminator is what
  fires.
  Proof: skills/record-work/spec-drift --self-test
- **When** a staged commit deletes both a change's `plan.md` and its contract-lane
  `proposal.md`, `spec-drift --retired` **shall** report that change once.
  Proof: skills/record-work/spec-drift --self-test
- **If** the only staged difference in a `proposal.md` carrying `Targets:` lies inside its
  `## Checklist` section, **then** default `spec-drift` **shall** still report staged code
  governed by that capability as code whose specification did not move. **When** anything
  outside `## Checklist` differs, it **shall** count as the capability moving.
  Proof: skills/record-work/spec-drift --self-test
- **Where** a change document holds `## Verification criteria` followed by a `## Checklist`
  section of `- [ ]` boxes, `criteria_of` **shall** return the criteria and none of the
  boxes.
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
  and names orphan code (a tracked file no `Governs:` line governs), a `Governs:` entry
  matching nothing, a criterion under
  *Verification criteria* with no `Proof:` beneath it, and a joined criterion — two or
  more `shall` resting on one `Proof:` with no ceiling named. **It exits 0 unconditionally**: the
  first run here found sixteen orphans and sixteen unproven bullets, which is the honest
  state rather than a regression, and a report that fails the build the day it lands is a
  report somebody deletes.
  Proof: skills/record-work/spec-drift --self-test
- **Where** a capability's `Governs:` holds a `!` entry matching a path, `governed_by`
  **shall** answer that the capability does not govern it, even when a plain entry matches.
  Proof: skills/record-work/spec-drift --self-test
- **While** a `Governs:` line holds no `!` entry, spec-drift **shall** treat its paths as it
  did before exclusions existed.
  Proof: skills/record-work/spec-drift --self-test
- **If** a path is excluded by the only capability whose plain entries match it, **then**
  `spec-drift --trace` **shall** list it as an orphan.
  Proof: skills/record-work/spec-drift --self-test
- **When** a path one capability excludes and another claims is staged alone, spec-drift's
  default mode **shall** name only the claiming capability's spec.
  Proof: skills/record-work/spec-drift --self-test
- **If** a `!` entry in `Governs:` matches no tracked file, **then** `spec-drift --trace`
  **shall** report it as a dead claim.
  Proof: skills/record-work/spec-drift --self-test
- **When** `write-spec` defines `Governs:` in "And the spec declares what it governs", it
  **shall** name the `!` exclusion. Proved as wiring only, pinned on `a leading .!. excludes`.
  Ceiling named: the row pins the phrase.
  Proof: scripts/check-payload
- **Where** a verification criterion carries the EARS marker two or more times — counted as
  `is_ears` matches it — sits above exactly one `Proof:` line, and contains no
  case-insensitive `ceiling`, `is_joined` **shall** return true for it, and false for every
  other criterion. The self-test drives the real `is_joined` through six cases: two markers
  on one proof, two markers with "Ceiling named:", two markers on two proofs, one marker, a
  bolded pair, and `marshalling` beside one real marker.
  Proof: skills/record-work/spec-drift --self-test
- **Where** a change document holds `## Verification criteria`, holding one criterion with
  a `Proof:` and one without, then `## Checklist`, `joined_of` **shall** print one line,
  carrying the proven criterion's `Proof:` count. It prints nothing for the unproven
  criterion, which is `--trace`'s *unproven criterion* already, and nothing for the boxes.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --trace` runs in a repository whose capability spec and change
  document each hold one joined criterion, it **shall** list both under the joined-criteria
  section, each with its file.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --trace` runs in a repository holding exactly two joined criteria,
  its summary line **shall** end in `, 2 joined criterion(a)`. The self-test asserts that
  exact text, so a missing count and a wrong one both fail it.
  Proof: skills/record-work/spec-drift --self-test
- **When** `spec-drift --anchors` runs in a repository holding a joined criterion whose
  citation resolves, it **shall** print exactly one line counting joined criteria in change
  documents and in capability specs.
  Proof: skills/record-work/spec-drift --self-test
- **While** no criterion in the repository is joined, `spec-drift --anchors` **shall**
  print no joined line, and its output is byte-identical to `anchors`, `ears` and `retired`
  run directly in the same repository with the separators `anchors_and_ears` prints between
  them. That composition stands in for the gate as it was before the joined count existed,
  which the self-test cannot run.
  Proof: skills/record-work/spec-drift --self-test
- **If** joined criteria are the only finding, **then** `spec-drift --anchors` **shall**
  exit 0.
  Proof: skills/record-work/spec-drift --self-test
- **If** a capability spec or change document holds an odd number of fence lines — the
  lines `defenced` toggles on, defined once in spec-drift's `FENCE` — **then**
  `spec-drift --anchors` **shall** exit non-zero naming that file.
  Proof: skills/record-work/spec-drift --self-test
- **While** every capability spec and change document holds an even number of fence lines,
  `spec-drift --anchors` **shall** print no fence line.
  Proof: skills/record-work/spec-drift --self-test
- **If** a file under `skills/` or `scripts/` pipes into `rg` with a short-flag group
  containing `q`, or with `--quiet`, **then** `scripts/check-payload` **shall** fail naming
  that file and line.
  Proof: scripts/check-payload
- **Where** the bump is decided before the request exists, `record-work` **shall** say to
  pass the label to `gh pr create`. Proved as wiring only, pinned on the phrase `the label
  rides gh pr create`. Ceiling named: the row pins the phrase, not that a session follows it.
  Proof: scripts/check-payload
- The `record-work` skill **shall** carry joined criteria as a row of its `--trace` table,
  pinned by a `check_wiring` row on the pattern `^\| \*\*joined criterion\*\* \|`. Proved
  as wiring only. **Ceiling named:** the row proves the table names joined criteria. It
  does not prove the row's description matches what `is_joined` does, which is a reading.
  Proof: scripts/check-payload
- The `review-spec` skill **shall** state that `--trace` lists joined criteria as a prompt
  to "split the criterion or name its ceiling", pinned by a `check_wiring` row on the
  pattern `split the criterion or name its ceiling`. Proved as wiring only. **Ceiling
  named:** the row proves the wording is present, not that a reviewer acts on it.
  Proof: scripts/check-payload
- **a criterion is a column-0 bullet under the criteria heading, and nothing else.** Its
  `Proof:` survives the continuation lines beneath it and ends at the next such bullet or
  the next heading. Reading every bullet in the document instead buries the real findings
  under the Outcomes section, which has no `Proof:` by design.
  Proof: skills/record-work/spec-drift --self-test
- **the change document is reviewed before code is built on it.** Every other review in the
  flow reads code and all of them run after the code exists, so a criterion no run could
  fail is agreed at the stop, becomes boxes, becomes code, and is then measured against a
  sentence nobody could have failed — the work passes and the promise was never real.
  `review-spec` runs once the checklist is written, through one fresh `spec-reviewer`, over
  ambiguity that forks the work, criteria that cannot fail, `Governs:` boundaries
  describing nothing, a criterion another capability contradicts, and per box what the
  document fails to answer. **Before the stop, not after**, because the findings have to be
  on the table while the agreement is still being made. It reports and never rewrites the
  spec — a reviewer editing the contract it reviews is the author twice over. Skipped
  entirely on the trivial lane, which has no contract.
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
- **find-work's queue section shall state that a proposal holding a box under its
  `## Checklist` is not queued.** Proved as wiring only.
  Proof: scripts/check-payload
- **find-work's queue section shall order queued changes that carry the same `Queued:`
  date by change name, ascending.** Proved as wiring only.
  Proof: scripts/check-payload
- **find-work's queue section shall name `libretto queue` as an implementation of its
  definition, and shall state that where the two disagree the section is the rule.**
  Proved as wiring only.
  Proof: scripts/check-payload
- **`/libretto-status` and `/libretto-next` shall not invoke `libretto queue`.** Proved as
  absence of the string.
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
- **a fan-out writer is given `proposal.md` by path, and always the vocabulary.** The change
  document is the shared ground a brief used to be: its decision log, its vocabulary and
  the primary capability's contract. Each writer opens what its subtask touches; the
  vocabulary is never trimmed, because it is the whole mechanism stopping two writers from
  giving one concept two names. **Parts are named, never excerpted** — excerpting puts N
  copies in N contexts. The return contract carries no criterion: `agents/spec-writer.md`
  already promises deltas rather than restatements and already forbids the empty return, so
  a row over it would have been green the moment it was written.
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
- **Where** a change reaches phase 2 on the contract lane, `write-spec` **shall** direct the
  session to write the primary capability's contract into `proposal.md` itself: `Targets:`,
  the six pillars as level-2 headings, and every answer verbatim, dated, `(assumed)` when
  nobody gave it, under `## Decision log`, carried into `## Prior decisions` per answer. It
  **shall not** direct the session to create `brief.md`, `decisions.md` or `spec.md` for a
  single-capability change. **Proved as wiring only.**
  Proof: scripts/check-payload
- **Where** a change amends more than one capability, `write-spec` **shall** direct that
  each extra capability's contract is written to a `spec-<capability>.md` beside
  `proposal.md` with its own `Targets:` line, and **shall** launch a writer for one only when
  the session's context cannot hold its research. **Where** the inputs do not settle a
  decision, that writer **shall** leave `[NEEDS CLARIFICATION: question]` and never guess;
  the session **shall** resolve each marker by asking, logging, and replacing the bracket
  expression only — the one declared exception to one-author-per-file, anything beyond it a
  relaunch. **Proved as wiring only.**
  Proof: scripts/check-payload
- **When** phase 5 runs on the contract lane, `write-plan` **shall** present two or three
  approaches with tradeoffs as one native question, recommended first with its reason,
  logging chosen and rejected verbatim in the decision log, and **shall** direct the session
  to write `## Approach` and the `Durable decisions:` line into the change document, with no
  box in the approach. **Proved as wiring only**: the fork row, the `No checkboxes.` row,
  and the log row naming the decision-log section.
  Proof: scripts/check-payload
- **When** the 5→6 seam runs on the contract lane, `write-tasks` **shall** direct the session
  to write the checklist under a level-2 `## Checklist` as the change document's last
  section, keeping `one badly cut box` and the one-writer rule. **Proved as wiring only.**
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` names `plan-writer` or
  `task-cutter`, or `agents/plan-writer.md` or `agents/task-cutter.md` exists, **then**
  `check-payload` **shall** fail naming it.
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` names `spec-writer` outside
  `write-spec`'s fan-out section and the agent file itself, **then** `check-payload`
  **shall** fail naming it. A section-aware row, because the rows are otherwise
  line-scoped.
  Proof: scripts/check-payload
- **If** `agents/spec-writer.md` names `brief.md` or `decisions.md` as an input, **then**
  `check-payload` **shall** fail naming the file. Observed red on the agent file as it read
  both, before its prose moved to `proposal.md`.
  Proof: scripts/check-payload
- **When** the contract lane's checklist is written, `review-spec` **shall** launch exactly
  one `spec-reviewer` with, by path, the change document, every `spec-<capability>.md`
  beside it, and the capability spec(s) their `Targets:` lines name, and with none of the
  conversation, and **shall** relay the reviewer's findings unedited. **Proved as wiring
  only**: the launch sentence and the agent name in `review-spec`.
  Proof: scripts/check-payload
- The `spec-reviewer` agent **shall** declare exactly the tools `Read, Grep, Glob, Bash,
  Skill`, and **shall** carry the instruction never to edit what it reviews; `review-spec`
  **shall** keep `never rewrites the spec`.
  Proof: scripts/check-payload
- The `spec-reviewer` agent **shall** report, for each box under `## Checklist`, what the
  contract and the approach fail to answer for building it, naming the gap and the part of
  the document that should have answered it, and **shall not** propose the answer. **Proved
  as wiring only**: one row on the agent's per-box sentence.
  Proof: scripts/check-payload
- **When** the fresh reviewer has returned on a contract-lane change, `review-spec` **shall**
  ask the one stop natively, on the whole change document, and `write-spec` and
  `write-tasks` **shall not** ask for a go-ahead; phase 8's push stays asked natively in
  `record-work`. One row each, written against the skills' own wording and observed red
  before the prose moved. **Ceiling named:** that no other stop exists stays uncheckable,
  the limit the stop table has always named — a script cannot tell a prompt from a
  paragraph.
  Proof: scripts/check-payload
- **When** `/libretto-flow` routes the contract lane, it **shall** route `write-spec`,
  `write-plan`, `write-tasks` and then `review-spec`, keep `before the stop, not after`, and
  state the gear as one stop on the contract lane, none on the trivial lane, and phase 8's
  question in both. `/libretto-attacca` **shall** still route to `review-spec`. **Proved as
  wiring only.**
  Proof: scripts/check-payload
- **While** a run is `/libretto-attacca`, the interview, the fork and the stop **shall**
  become `(assumed)` entries in the change document's `## Decision log`, each naming what
  changes if it is wrong — logged where a fan-out writer reads, so an assumed answer
  reaches it the same way a given one does — and the command **shall** name the one stop
  and the push as what the invocation answered. **Proved as wiring only.**
  Proof: scripts/check-payload
- **When** phase 8 lands a change-document change, `record-work` **shall** direct the
  session to do three things in one commit: apply `proposal.md`'s contract and each
  `spec-<capability>.md` onto the capability spec its `Targets:` names, retire the
  `## Prior decisions` entries that outlive the change, and delete the folder. For an
  old-layout change it **shall** read `spec.md` and `plan.md` as before. **Proved as wiring
  only.**
  Proof: scripts/check-payload
- **When** a change folder's `proposal.md` holds at least one box under `## Checklist`,
  `find-work` **shall** read it as that change's checklist in both scans, ahead of
  `tasks.md` and `plan.md`; **where** it holds none, the skill **shall** fall back to
  `tasks.md`, then `plan.md`. It **shall** count boxes by the CLI's definition: a line
  matching `^\s*[-*]\s+\[( |x|X)\]`, none inside a fence, in `proposal.md` only within a
  `## Checklist` section closed by the next `^##\s` heading. **Proved as wiring only**: the
  scan lines, the selection sentence and the box definition are present in the skill.
  Proof: scripts/check-payload
- **When** `/libretto-status` reports, it **shall** delegate to `find-work`'s scan, and so
  count boxes by the same definition as the CLI, and **shall** name `proposal.md` among the
  files that scan reads, never restating the scan. **Proved as wiring only.**
  Proof: scripts/check-payload
- **When** step 0 sends a change down the trivial lane, `write-spec` **shall** direct that
  the change folder phase 1 created is deleted in the trivial lane's own commit, and
  `find-work` **shall not** claim the trivial lane never creates a `changes/` directory.
  **Proved as wiring only**: the instruction is present and the old sentence is absent.
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` contains the acronym quoted in
  the 2026-10-21 naming answer as a case-sensitive whole word, or the method's spelled-out
  name in any case, **then** `check-payload` **shall** fail naming the file. Observed red on
  a hand-broken copy before it was trusted, since the tree it shipped against held neither.
  Proof: scripts/check-payload
- **When** `write-spec` or `libretto-queue` finds no `.agents/changes/` in the
  project, it **shall** still refuse to create the directory; the payload spec
  **shall** name explicit, user-invoked `libretto init` as the sole owner of
  bootstrap and **shall** keep implicit creation by skills a non-goal.
  **Ceiling named:** the anchor keeps the sentences findable in their files, never
  proves a session obeyed prose — the same limit every skill criterion here lives with.
  Proof: .agents/specs/payload/spec.md
- **Where** `libretto` is on PATH, `skills/write-spec/SKILL.md`'s refusal paragraph
  **shall** point at `libretto init` as the sanctioned bootstrap; **where** the binary
  is absent, the paragraph **shall** leave the refusal unchanged — the pointer only
  says where bootstrap lives.
  Proof: skills/write-spec/SKILL.md
- **Where** `libretto` is on PATH, `commands/libretto-queue.md`'s refusal section
  **shall** point at `libretto init` with the same guarded grammar; **where** the
  binary is absent, it **shall** say so and stop, exactly as it does today.
  Proof: commands/libretto-queue.md
- **When** the bootstrap pointer lands in `skills/write-spec/SKILL.md`, its frontmatter
  **shall** carry `version: "2.1"`, and its frontmatter **shall** still parse with
  `name:` matching the directory.
  Proof: scripts/check-payload
- The `libretto-converge` command **shall** route to the `converge` skill with
  `Skill(skill="converge")`. Pinned by `check_wiring` on `Skill\(skill="converge"\)`.
  **Ceiling named:** proves the route is written, not that a host follows it.
  Proof: scripts/check-payload
- The `converge` skill **shall** launch exactly one `converge-reader`, given by path the
  tree as it is now, the change folder and the capability specs its `Targets:` lines name,
  never a diff range and none of the conversation. Pinned by `check_wiring` on
  `exactly one .converge-reader.`, on `none of the conversation` and on
  `never a diff range`. **Ceiling named:** the launch and the paths actually handed over
  are behaviour no row sees.
  Proof: scripts/check-payload
- The `converge-reader` agent **shall** declare exactly `Read, Grep, Glob, Skill` as its
  tools. Pinned by `check_wiring` on `^tools: Read, Grep, Glob, Skill$`, so a `Bash` added
  later fails the gate.
  Proof: scripts/check-payload
- The `converge-reader` agent **shall** return every open box as exactly one of open,
  closable or drifted. Pinned by `check_wiring` on `exactly one of open, closable or
  drifted`. **Ceiling named:** the return shape is stated, not checked against a real
  return.
  Proof: scripts/check-payload
- **If** `git status --porcelain --untracked-files=no` lists any file when converge starts,
  **then** the `converge` skill **shall** stop before launching the reader and name those
  files. An untracked file never stops it. Pinned by `check_wiring` on
  `a dirty tree stops converge before the reader` and on `--untracked-files=no`.
  **Ceiling named:** naming the files is the instruction's, not observed.
  Proof: scripts/check-payload
- **If** no change is named and `find-work`'s source-1 scan returns more than one change
  with an open box, **then** the `converge` skill **shall** ask which one with the native
  prompt. The queue scan is not run. Pinned by `check_wiring` on `ask which one` and on
  `source 1 only, never the queue`. **Ceiling named:** that nothing starts before the
  answer is behaviour no row sees.
  Proof: scripts/check-payload
- The `converge` skill **shall** name `find-work` as the definition of a box and of which
  file holds the checklist. Pinned by `check_wiring` on `find-work. defines the box`.
  **Ceiling named:** the pointer, not a reading of it.
  Proof: scripts/check-payload
- The `converge` skill **shall not** restate the box pattern. Pinned by `check_absent` on
  the escaped `planLine` regex, so a copy that would drift from the CLI's fails the gate.
  Proof: scripts/check-payload
- **When** the reader returns a box as drifted, the `converge` skill **shall** leave that
  box's text and state unchanged and report what it says beside what the code does
  instead. Pinned by `check_wiring` on `drift is reported, never rewritten` and on
  `beside the reader.s file and line`. **Ceiling named:** an untouched file is behaviour.
  Proof: scripts/check-payload
- **If** any proof cited by a closable box's criteria has not been run by the session in
  the foreground and read passing, **then** the `converge` skill **shall** leave that box
  open. Pinned by `check_wiring` on `a reader.s verdict never closes a box`. **Ceiling
  named:** the rule's presence, not a run that honoured it.
  Proof: scripts/check-payload
- Before running any proof, the `converge` skill **shall** run `spec-drift --anchors` and
  count a citation it reports unresolved as not passed. Pinned by `check_wiring` on
  `spec-drift --anchors. before any proof`. **Ceiling named:** the order and the counting
  are written, not observed.
  Proof: scripts/check-payload
- **If** the `-v` output of a test-named citation carries no `--- PASS: <TestName>` line
  for that test, **then** the `converge` skill **shall** count it as not passed, whatever
  its exit status. Pinned by `check_wiring` on `--- PASS: <TestName>` and on
  `whatever its exit status`. **Ceiling named:** the rule, not a session reading the
  output.
  Proof: scripts/check-payload
- **If** a file-only citation exits non-zero or prints a `FAIL` line, **then** the
  `converge` skill **shall** count it as not passed. Pinned by `check_wiring` on
  `exits 0 and prints no FAIL line`. **Ceiling named:** the rule, not a session reading
  the output.
  Proof: scripts/check-payload
- **If** a box's `Closes:` text matches no criterion of the contract unambiguously, or the
  box carries no `Closes:`, **then** the `converge` skill **shall** leave it open and
  report it as `unmapped`. Pinned by `check_wiring` on `reported as .unmapped.`.
  **Ceiling named:** whether a match was unambiguous is a reading no row sees.
  Proof: scripts/check-payload
- **When** the session closes a box on observed proof, the `converge` skill **shall** have
  it append `— converged <date>: <citation>, <citation> passed`, listing every citation
  that passed, to the box's last line. Pinned by `check_wiring` on
  `converged <date>: <citation>, <citation> passed` and on `the box.s last line`.
  **Ceiling named:** the format and place, not a written box.
  Proof: scripts/check-payload
- The `converge` skill **shall** count as uncovered only criteria under the
  `## Verification criteria` of the change's `proposal.md` and of each
  `spec-<capability>.md` beside it that no box's `Closes:` names. Pinned by `check_wiring`
  on `uncovered means the change.s own criteria`. **Ceiling named:** the definition, not
  a session applying it.
  Proof: scripts/check-payload
- **Where** an uncovered criterion's proof does not pass when the session runs it, the
  `converge` skill **shall** append one open box marked `(converge <date>)` and naming
  that criterion in `Closes:` after the checklist's last box. Pinned by `check_wiring` on
  `appended after the checklist.s last box` and on `\(converge <date>\)`. **Ceiling
  named:** an appended box is behaviour.
  Proof: scripts/check-payload
- **Where** an uncovered criterion's proof passes when the session runs it, the `converge`
  skill **shall** append no box for it. Pinned by `check_wiring` on
  `passes already, so no box is appended`. **Ceiling named:** the instruction only.
  Proof: scripts/check-payload
- The `converge` skill **shall** write to the checklist and nothing else. Pinned by
  `check_wiring` on `writes to the checklist and nothing else`. **Ceiling named:** no gate
  can see what a session wrote; the 6→7 reviewer reads the diff.
  Proof: scripts/check-payload
- The `converge` skill **shall** never uncheck a closed box. Pinned by `check_wiring` on
  `never unchecks a closed box`. **Ceiling named:** the instruction only.
  Proof: scripts/check-payload
- **When** converge marked or appended any box, the `converge` skill **shall** end with one
  `docs(<change>): converge the checklist` commit holding the checklist file alone. Pinned
  by `check_wiring` on `docs\(<change>\): converge the checklist` and on
  `stages nothing else`. **Ceiling named:** the commit's contents are behaviour, seen in
  `git show` and not by the gate.
  Proof: scripts/check-payload
- **If** converge marked and appended nothing, **then** the `converge` skill **shall**
  make no commit. Pinned by `check_wiring` on `nothing was marked or appended, no commit`.
  **Ceiling named:** the instruction only.
  Proof: scripts/check-payload
- The `converge` skill **shall** never push or create a branch. Pinned by `check_wiring` on
  `never pushes and never creates a branch`. **Ceiling named:** the instruction only.
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
watching. Claims, not facts: that a run passes the stop before code with the change
document on disk; that it ends at a pushed branch with a request whose description names
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
