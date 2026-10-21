Targets: payload

# Payload — one change document, one stop

The delta that moves the payload's **contract lane** from a chain of documents written by
three fresh writers, with two stops, to **one change document**: `proposal.md`, grown from
its capture, written inline by the session, read once by **one fresh reviewer**, and
agreed at **one stop**. The **trivial lane**, the bug rule, capability specs, `--anchors`
and the landing do not move. Every payload reader keeps reading changes in the **old
layout**.

The Go CLI (`loop`, `metrics`, `wiki`, `land`, `cmd/libretto/recommend.go`) is the sibling
`cli` delta, and the README's promise is the sibling `readme` delta. Neither is specified
here. This delta covers everything under `skills/`, `agents/` and `commands/`, plus
`scripts/check-payload`, which `payload` governs.

## Outcomes

**The reversal, stated once.** This delta retires "two stops with a contract, none without,
plus phase 8's question". It also retires the chain proposal → brief → spec delta (fresh
`spec-writer`) → `plan.md` (fresh `plan-writer`) → `tasks.md` (fresh `task-cutter`). It
replaces both. It does not sit beside them. After it lands, for a change started after it
lands:

1. **One change document.**
   A contract-lane change keeps its contract in `.agents/changes/<change>/proposal.md`. That
   is the file `find-work` or `/libretto-queue` captured, and it grows in place. No
   `brief.md`, `decisions.md`, `spec.md`, `plan.md` or `tasks.md` is created for the change.
   Written top to bottom, the file holds:

   | Part | Where | Written by phase |
   |---|---|---|
   | `Tracker:` · `Targets:` · `Durable decisions:` lines | top, column 0, outside any fence | 1 (`Tracker:`), 2 (`Targets:`), 5 (`Durable decisions:`) |
   | the ask, verbatim · the reading | the headings the capture already writes | 1 or the queue, unchanged from today |
   | the decision log | `## Decision log` | 2 and 5. One `- Q: … → A: …` line per answer, verbatim, under dated `### Session` headings, `(assumed)` kept |
   | the contract | `## Outcomes` · `## Scope boundaries` · `## Constraints` · `## Prior decisions` · `## Task breakdown` · `## Verification criteria` | 2–3 |
   | the approach | `## Approach`. The parts `write-plan` names today sit under it as `###` subsections | 5 |
   | the checklist | `## Checklist`: boxes and nothing else, the document's last section | the 5→6 seam |

   **The checklist is the boxes under the level-2 `## Checklist` heading, and nothing else.**
   A box anywhere else in the document is not a checklist box, and the *Task breakdown*
   pillar is written as a list, never as boxes.

   **A change amending two or more capabilities** carries the primary capability's contract
   in `proposal.md` (its `Targets:` line plus the six pillars). It adds one
   `spec-<capability>.md` beside it for each extra capability, each with its own `Targets:`
   line and its own six pillars. The decision log, the approach, the `Durable decisions:`
   line and the checklist exist once, in `proposal.md`.

2. **Written inline.** The session that ran the interview and the fork writes `proposal.md`
   itself, every part of it. Nothing launches `plan-writer` or `task-cutter` any more, and
   nothing launches a writer for `proposal.md`. The one exception is `write-spec`'s fan-out,
   which survives in exactly one form. When a change has extra capabilities **and** this
   session's context cannot hold their research, it launches one writer per extra
   `spec-<capability>.md`, each writing that one file and nothing else. With room in
   context, the session writes those files inline as well. Context is the only reason to
   fan out, which is the rule `write-spec` already states.
   That writer is `agents/spec-writer.md`, which survives for this case only. It is given
   `proposal.md` by path as its shared ground: the decision log, the vocabulary and the
   primary capability's contract. A fan-out change writes no `brief.md` and no
   `decisions.md`.

3. **One fresh reviewer, before the stop.** Once the checklist is written, `review-spec`
   launches exactly one `spec-reviewer` subagent. It is given by path the change document,
   every `spec-<capability>.md` beside it, and the capability spec(s) their `Targets:` lines
   name. It is given none of the conversation. It runs `spec-drift --trace` and `--ears`,
   then reads the whole change document and returns findings. It writes nothing. It covers
   two kinds of finding:
   - **the contract**: ambiguity, testability, `Governs:` and cross-capability
     contradiction, as `review-spec` reads today
   - **the checklist, box by box**: what the contract and the approach together fail to
     answer for building that box. This is the question `task-cutter` asked: a task implied
     with no mechanism described, a criterion with no proof that says what closes it, two
     statements that cannot both hold, a boundary that does not say which side something
     falls on. It names the gap and the part of the document that should have answered it.
     It never proposes the answer.

4. **One stop.** The contract lane stops exactly once before code. The stop comes after the
   fresh reviewer returns. It covers the whole change document (contract, approach and
   checklist), with the reviewer's findings relayed unedited. `review-spec` owns it and asks
   it as a native question. The options are: carry on to phase 6 (recommended, naming the
   first box), change the document first, or stop here, plus room to answer differently.
   After that comes phase 8's push question, unchanged. On the contract lane the flow has
   **one stop plus phase 8's question**. On the trivial lane it has **phase 8's question
   alone**.

5. **Each question still belongs to the phase that asks it.** Phase 2's interview and phase
   5's fork keep their shape:
   - the interview: one question per native-prompt call, a soft bound of about five, "no
     more questions" always offered, and zero questions legitimate
   - the fork: two or three approaches, the recommended one first

   Both happen before the stop. Each answer is written verbatim into the change document's
   `## Decision log` and carried into `## Prior decisions` as it is answered. That is the
   same double write `write-spec` asks for today, now inside one file.

6. **`/libretto-attacca` answers the stop and the push.** The bump question at its end is
   unchanged. Every question it would have asked becomes three things: an `(assumed)` entry
   in the change document's `## Decision log`, an assumed entry under `## Prior decisions`,
   and a line in the phase 7 report and the request's description. Because the run authors
   the document, it acts on the fresh reviewer's findings before phase 6 and names each one
   it acted on.

7. **Both layouts read, by one rule.** Every payload reader that counts boxes uses the CLI's
   definition. A box is a line matching `^\s*[-*]\s+\[( |x|X)\]`, and a box inside a fence
   does not count. In `proposal.md` only boxes in a section opened by a `## Checklist` heading
   and closed by the next `^##\s` heading count. In `tasks.md` and `plan.md` every box outside
   fences counts. Every payload reader that finds a change's checklist picks the file this
   way:
   - `proposal.md` is the checklist when it holds at least one box under `## Checklist`
   - otherwise `tasks.md`
   - otherwise `plan.md`

   A file existing never decides it on its own, because every change in both layouts has a
   `proposal.md`. Under this rule, `find-work`'s two scans, its three-row table (in flight ·
   finished and not landed · no checklist yet) and `/libretto-status` report old-layout and
   new-layout changes the same way. A queued capture has no `## Checklist`, so it stays in
   the third row. The queue scan (`^Queued:` in `proposal.md`) does not change.

8. **The landing gate reads the change document.** `spec-drift --retired` treats a deleted
   `proposal.md` as a contract only when its `HEAD` copy carries a `Targets:` line outside
   fences, and nothing else makes it one. This is the same discriminator `libretto land` uses
   (the sibling `cli` delta's N3), so a `proposal.md` with no `Targets:` is a contract for
   neither tool. It fails a staged commit that deletes
   such a file when both of these hold:
   - no capability spec's *Prior decisions* section moved
   - the deleted copy does not declare `Durable decisions: none`

   Deleting a queued capture or a trivial-lane capture fails nothing and reports nothing as
   retired, because neither carries a contract. The old-layout rule over `plan.md` stays
   exactly as it is. A change whose `plan.md` and contract `proposal.md` are both deleted is
   reported once.

9. **Marking a box is not the contract moving.** The change document carries `Targets:`, so
   default `spec-drift` would read every box-marking commit as "the capability's spec moved"
   and say nothing about the code staged beside it. A staged `proposal.md` counts as its
   targets moving only when something outside its `## Checklist` section differs between
   `HEAD` and the index.

10. **The trivial lane's capture does not outlive it.** Phase 1 writes `proposal.md` for
    every source-3 request before phase 2's step 0 decides the lane, so a trivial-lane change
    always has a capture on its branch. The trivial lane deletes that folder **in its own
    commit**. This resolves a contradiction inside `find-work`: its artifact table says
    `proposal.md` must exist, and its branch section says "a change that needed no spec never
    creates a `changes/` directory". The directory is created, and it is gone by the time the
    work is committed.

11. **No acronym in what installs.** No file under `skills/`, `agents/` or `commands/` names
    the method this change was inspired by, either as an acronym or spelled out. The flow is
    "the Libretto flow", and its lanes are the contract lane and the trivial lane.
    `docs/FLOW.md` credits gentle-ai in one line and nowhere else. Specs quoting a decision
    verbatim may contain the acronym.

## Scope boundaries

**In:**

- `skills/write-spec`, `write-plan`, `write-tasks`, `review-spec`, `find-work`,
  `record-work` (its prose and `spec-drift`), `review-work`, `build-and-check` and
  `present-work`, each where it names a document this delta renames, moves or retires, or a
  stop it counts
- `agents/spec-reviewer.md` (new); `agents/plan-writer.md` and `agents/task-cutter.md`
  (deleted); `agents/spec-writer.md` (narrowed to the fan-out, reading `proposal.md` by path
  instead of `brief.md` and `decisions.md`, per *Outcomes* 2); `agents/work-reviewer.md` (reads the contract in either layout, and its
  aside "like spec-writer's write scope" is reworded)
- `commands/libretto-flow.md`, `libretto-attacca.md`, `libretto-next.md` and
  `libretto-status.md`. `libretto-queue.md` is read and left unchanged, because its capture
  is already the first state of the change document
- `scripts/check-payload`: the wiring rows this delta retires or rewords, and the rows it
  adds
- `docs/PAYLOAD.md`, regenerated by `scripts/check-payload --index` because the agent set
  changes

**Out:**

- **any Go behaviour.** `loop`, `metrics`, `wiki`, `land` and `cmd/libretto/recommend.go`
  belong to the sibling `cli` delta. This delta states the checklist rule (*Outcomes* 7) and
  the agent set only so that the two deltas agree.
- **the README's promise and its test.** `readme/spec.md` and
  `cmd/libretto/readme_test.go` belong to the sibling `readme` delta.
- **`AGENTS.md` and `docs/*` as contract.** No capability governs them. The lines that go
  stale are listed under *Task breakdown* as plain doc edits, with no criterion.
- **migrating anything in flight.** No change folder is rewritten into the new layout. A
  change that already holds `spec.md`, `plan.md` or `tasks.md` (on this tree none does, but
  `feat/add-init-command` does) finishes in the layout it started in. Whatever it still
  lacks, the session writes inline in that layout, because no writer is launched for either
  layout.
- **a second flow, a size flag or a profile.** The gear stays knowable from the change:
  step 0's question decides the lane, as it does today.
- **renumbering the phases.** There are still eight. The phases that wrote three documents
  now write three parts of one. Renumbering everything that says eight buys nothing, which
  is the answer the review seams already gave.
- **changing any of these:** the trivial lane's test, the bug rule, `--anchors`' citation
  and EARS checks, the all-or-nothing landing commit, `review-work`'s fresh
  `work-reviewer`, phase 8's push question, or the bump question. They are read, and
  restated where a path moved. They are never altered.
- **changing what `/libretto-queue` writes or what `Queued:` means.**
- **the reviewer writing.** `spec-reviewer` gets no `Write` or `Edit`, for the same reason
  `work-reviewer` has none: a reviewer that edits what it reviews is the author twice.
- **a second fresh reader.** One reviewer, not a panel.
- **a contract block per capability inside `proposal.md`.** Extra capabilities go in
  `spec-<capability>.md` files, per the logged answer.

## Constraints

- **The parsers fix the heading levels.** In `spec-drift`, `criteria_of`,
  `criteria_without_proof` and `section_of` open and close a section on `^##[[:space:]]`
  only. A `###` heading neither opens nor closes one. So the change document needs
  `## Verification criteria`, `## Prior decisions` and `## Checklist` as level-2 headings,
  with `## Checklist` as its own section. A box under `## Verification criteria` with no
  level-2 heading between them is a column-0 `- ` bullet. `criteria_of` reads it as a
  criterion, and `--ears` fails it for carrying no `shall`.
- **Every reader counts boxes by the same definition, the CLI's.** A box is a line matching
  `^\s*[-*]\s+\[( |x|X)\]` (`planLine`'s shape). A box inside a fence does not count. In
  `proposal.md` only boxes in a section opened by a `## Checklist` heading and closed by the
  next `^##\s` heading count. In `tasks.md` and `plan.md` every box outside fences counts.
  `find-work` and `/libretto-status` state this definition in the skill itself; a bare count
  of `- [ ]` lines over the whole file is not it, because it counts boxes in fences and, in
  `proposal.md`, boxes outside `## Checklist`.
- **Every `Proof:` line in the change folder is resolved.** `proofs_of` reads every
  `^\s*Proof:` line outside fences in every `.md` under the change roots. The approach and
  the checklist start a line with `Proof:` only when they cite a real test.
- **Every `Targets:` at line start counts.** `targets_of` reads every such line outside
  fences. `proposal.md` writes exactly one, for the primary capability. Each
  `spec-<capability>.md` writes exactly one, for its own capability. No other line starts
  with `Targets:`.
- **`--anchors` and `--ears` already read `proposal.md` and `spec-<capability>.md`.**
  `change_files` is every `.md` under `.agents/changes`, `changes` and `openspec/changes`,
  and both checks already run over it (spec-drift:99, :239, :434). Two things need to
  change: `--retired` (spec-drift:365, which matches `plan.md` only) and how default mode
  treats a staged `proposal.md` (spec-drift:165-169).
- **`--retired` stays inside `--anchors`.** No seventh gate, because the count "six" is
  written in ten places.
- **A skill never references `scripts/` or `docs/`**, and any tool a skill needs ships
  inside the skill. `spec-reviewer` reaches `spec-drift` the way `review-spec` does today:
  `<skill-base>/../record-work/spec-drift`.
- **`check-payload` rows are line-scoped `rg -qN`.** The phrase a row searches for sits on
  one line of the file it guards. Before a new row is trusted, it is observed failing on a
  hand-broken copy, per the existing *green on first run is not evidence* practice.
- **One writer per file.** The session writes `proposal.md`. A fan-out writer writes one
  `spec-<capability>.md`. The fresh reviewer writes nothing. `build-and-check` and any
  subagent report a finished box, and the session marks it.
- **Every stop is a native prompt**: `AskUserQuestion` on Claude Code, the host's own native
  prompt elsewhere, and in conversation where there is none. `record-work` keeps the one
  statement of why.
- **An agent's frontmatter `name:` equals its filename**, and an agent ships only if the
  skill that launches it names it.
- **A skill whose instructions change meaningfully bumps its own `metadata.version`.** The
  repo tag does not force the bump, and this delta changes the contract of every skill in
  *Scope boundaries · In*.
- **No new dependency.** bash, `rg`, `git` and `awk`, as `spec-drift` and `check-payload`
  use today.
- **This delta removes and reverses promises.** The commit that lands it carries `!` per
  `AGENTS.md`. While the version starts with `0.`, that is still a minor bump, never a major.
- **Old-layout readers stay until nothing old is in flight.** Every place that reads an
  old-layout name says when that reading can go: a `ponytail:` comment in code, a sentence
  in prose. This is the precedent the `plan.md` → `tasks.md` rename of 2026-08-17 set.

## Prior decisions

### Answered by the user, 2026-10-21 (verbatim from the change's decision log)

- Q: for substantial work, what replaces the proposal → brief → spec → plan → tasks chain? → A: Un doc de cambio + specs vivos (Recommended) — one document per change (intent, EARS criteria, approach, checklist); capability specs with Governs:/Proof: and the landing stay
- Q: with one change document, where does the flow stop for the go-ahead before the push? → A: Un stop: el doc entero (Recommended) — contract, approach and checklist reviewed together, once, before the first line of code; then only the push question
- Q: who writes the change document, and where does the fresh eye stay? → A: Sesión escribe, revisor fresco lee (Recommended) — the session writes the document inline; one fresh reviewer (review-spec) reads it before the stop; spec-writer, plan-writer and task-cutter launches go from three to one review
- Q: what happens to changes and tooling that read proposal/plan/tasks today? → A: Leer ambos formatos (Recommended) — new changes are born as one document; loop/metrics/wiki/find-work keep reading the old layout for what already exists; nothing in flight is touched
- Q: is the change document the same file growing from the capture, or a new one? → A: proposal.md crece (Recommended) — the queue/find-work capture is the document; criteria, approach and checklist are added to the same file; Queued: and proposal.md scans do not change
- Q: what is this called in Libretto? → A: Sin sigla, vocabulario propio (Recommended) — no "ODD" in the payload; the flow stays "the Libretto flow", lanes named in its own vocabulary; FLOW.md credits gentle-ai as inspiration in one line
- Q: PR #71 (libretto land) reads Targets: from old deltas; how are the two changes ordered? → A: #71 primero, este se adapta (Recommended) — #71 merges as is; this change rebases on it and its cli delta makes land read Targets: from proposal.md
- Q: without task-cutter, who checks that the document is enough to build every box? → A: El revisor fresco lo hereda (Recommended) — the one fresh reviewer reads contract and checklist and reports, per box, what the document fails to answer
- Q: how does a change touching two or more capabilities carry its contract? → A: proposal.md + spec-<cap>.md extra (Recommended) — proposal.md carries the primary capability contract; each extra capability gets a spec-<cap>.md beside it with its own Targets:
- Q: when an extra spec-<cap>.md needs its own writer, who writes it and what does it read? → A: spec-writer, lee proposal.md (Recommended) — agents/spec-writer.md survives for this case only, given proposal.md by path; no brief.md

### Derived by the orchestrator, 2026-10-21. Logged, not a user answer

- Q: which boxes in proposal.md count as the checklist? → A: only boxes under a level-2 `## Checklist` heading — settled by the code, not asked: spec-drift parsers react to `## ` only, and boxes elsewhere would be read as criteria by --ears (derived, orchestrator)
- Q: spec-writer survives for the fan-out (decision 11); does it keep its `libretto models` recommendation? → A: yes, sonnet/high as today; R1 forbids only plan-writer and task-cutter keys (derived, orchestrator, from review-spec finding 1)
- Q: what makes a deleted proposal.md a contract, for both `--retired` and `land`? → A: a `Targets:` line outside fences, and nothing else; the `## Verification criteria` arm is dropped so the two landing tools cannot disagree (derived, orchestrator, from review-spec finding 2)
- Q: does find-work count boxes the way the CLI does? → A: yes — planLine's box shape, fences excluded, `^##\s` section boundaries, in every reader (derived, orchestrator, from review-spec finding 3)

### Derived from those answers. Not user answers, and open to challenge

- **The decision log becomes a section of the change document, and a new change writes no
  brief.** `brief.md` and `decisions.md` existed to feed writers that never had the
  conversation. With the session writing inline, that reason is gone, and a second file
  would contradict answer 1's *one document*. Old-layout changes keep theirs. **Ceiling
  named:** a fan-out writer has no conversation either. It reads `proposal.md` by path, whose
  decision log and contract are the shared ground a brief used to carry (*Outcomes* 2).
- **`review-spec` owns the stop.** It is the last thing that runs before code, the findings
  have to be on the table when the agreement is made, and one file asks the question. The
  split mirrors `write-plan` today: the skill's opening sections are the session's conduct
  (launch, relay, ask), and the reviewer's instructions start at the reading.
- **The fresh reviewer is an agent file, `agents/spec-reviewer.md`**, not a generic
  subagent. Every fresh seat in this payload has one (`work-reviewer` and the five review
  lenses), and *an agent that ships is named by the skill that launches it*. Its tools are
  `Read, Grep, Glob, Bash, Skill`. It needs `Bash` because `spec-drift --trace` and `--ears`
  are its first step, and `Skill` because its first instruction is `evidence`.
- **`write-spec`, `write-plan` and `write-tasks` stay as skills.** Each now writes its part
  of the change document. Their rules (the six pillars and EARS, the fork and the
  alternatives table, the standalone box) are what those parts are written against. The
  existing *renumbering buys nothing* decision keeps the phase numbers.
- **A contract-lane `proposal.md` is recognised by its contract, never by its name.** A
  `Targets:` line outside fences marks one, and nothing else does: the same discriminator
  `libretto land` uses, so the two landing tools cannot disagree. A queued capture and a
  trivial-lane capture carry none, so `--retired` needs no `Queued:` special case to spare
  them. **Ceiling named:** a contract-lane document written without `Targets:` is invisible
  to the gate. That is the same blind spot an old-layout change has when its `plan.md` was
  never written. The day it happens, the fix is the fresh reviewer reporting a contract with
  no `Targets:`.
- **`proposal.md` wins as the checklist only by holding a box under `## Checklist`.** The
  precedent's *newer wins* survives. Its *exists* test cannot, because `proposal.md` exists
  in both layouts.

### Reversed by this delta: each one named, none softened

- **"Two stops with a contract, none without, and phase 8's question in both"** (*ceremony
  proportional to the change*) → one stop on the contract lane, none on the trivial lane,
  phase 8's question in both.
- **The stop table's rows `2–3 · write-spec · yes` and `5→6 · write-tasks · yes`** → no stop
  at either. They are replaced by one row, `before 6 · review-spec · yes`, covering the
  whole change document.
- **"The stop moved rather than multiplied … the count stays at three, two of them inside
  the work"** → the count is two, one of them inside the work.
- **"The contract and the plan are … drafted by agents that never saw the conversation"**
  (`spec-writer`, single case included; `plan-writer`, read-only) → the user still
  co-authors them, and the session that held the conversation drafts them inline. The
  fresh eye moves from drafting to reading. The fan-out is the only place a writer
  survives, per *Outcomes* 2.
- **"The spec's author is `spec-writer`, single case included — a single writer is a
  fan-out of one"** → retired. A single-capability change has no writer subagent.
- **"The checklist is cut by an agent that never saw the design"** (`task-cutter`) → the
  session cuts it inline. The cutter's second output, *what the documents failed to
  answer*, moves to the fresh reviewer and is reported per box.
- **"A plan that holds the how, and a checklist that holds the state — two files, not
  one"** → two sections of one file. The separation survives as sections: the approach
  holds no box, and the checklist holds nothing else.
- **The 2026-08-19 answer *"Brief siempre"*** → retired for every new change, fan-out
  included. No `brief.md` is written in the new layout, ever: a fan-out writer is given
  `proposal.md` by path instead. *"Orquestador parchea inline"* → kept for the fan-out only,
  where a writer can still leave a `NEEDS CLARIFICATION` marker for the session to resolve.
  A single-capability change has no writer and so no marker to patch. *"Sí, marcadas
  'assumed'"* survives in a new place: assumed answers land in the change document's
  `## Decision log`. *"Blando ~5, con juicio"* is untouched.
- **"The task cut is a subagent, not a numbered phase"** → there is no cut subagent. The
  checklist is written inline in the 5→6 seam, which stays unnumbered.
- **"A plan's `Durable decisions:` line is a claim … the list is the change's `spec.md`
  *Prior decisions*"** → for a new-layout change, the line lives in `proposal.md` and the
  list is that document's own `## Prior decisions`. Reading the line as a claim, not the
  list, is kept.
- **"Every file has exactly one author. The brief is written once …; each spec delta has one
  writer; the plan has one writer, the orchestrator"** → the change document has one
  writer, the session. Each `spec-<capability>.md` has one writer, the session or a fan-out
  writer. No new single-capability change has a brief.
- **"An agent that ships is named by the skill that launches it. Phase 3's fan-out runs
  `spec-writer`"** → the rule stands. Its example becomes `review-spec` launching
  `spec-reviewer`.
- **"`ask` holds at phases 1, 2 and the 5→6 seam, and nowhere after"** → phases 1, 2 and 5,
  up to the stop, and nowhere after.
- **"The same flow with its three stops answered in advance"** (`libretto-attacca`) → its two
  stops, the one stop and the push, answered in advance.
- **"`review-spec` runs in the 3→5 seam … before a plan is built on it"** → it runs after
  the checklist and before the stop, over the whole change document. It runs as a fresh
  subagent, not inline in the session.
- **"While a run is `/libretto-attacca`, … every marker shall become `(assumed)` entries in
  `decisions.md` … logged where the writers read"** → assumed entries go to the change
  document's `## Decision log`, which is where a fan-out writer reads them, in
  `proposal.md`.

### Unchanged, and not relitigated here

All of these stay as they are:

- one flow with no size flag
- the trivial lane and its test (*could anyone reasonably disagree about what done means*)
- the bug rule (criterion and red test before the fix)
- capability specs with `Governs:`/`Proof:`/EARS
- `--anchors`
- landing in one commit
- `review-work`'s fresh `work-reviewer`
- attacca answering as `(assumed)`
- phase 8's push question
- the bump question
- *phase 5 asks* (the fork)
- *one reviewer, not a panel*
- *the retirement gate compares the section, never the file*
- *its escape is a declaration, not a flag*
- *context is the only reason to fan out*
- *model and effort do not move part-way through a fan-out*

## Task breakdown

The units are cut by concern, so each can be built, gated and committed on its own. The
capability spec is not part of any unit. This delta lands on `.agents/specs/payload/spec.md`
once, in the landing commit.

1. **`spec-drift` reads both layouts.** Waits on nothing.
   - `--retired` widens to a contract-lane `proposal.md` (*Outcomes* 8), with a change that
     deletes both files reported once.
   - Default mode ignores a staged `proposal.md` whose only difference is inside
     `## Checklist` (*Outcomes* 9).
   - Self-test cases cover both, driven through real temporary repositories the way the
     existing `land` helper is.
   - A self-test case proves that a `## Checklist` after `## Verification criteria` is not
     read as criteria.
   - The header comment's layout tree (spec-drift:10-16) shows both layouts and
     `spec-<capability>.md`.
2. **Phase 1 and status read both layouts.** Waits on nothing.
   - `find-work`'s two scans and its source table apply the *Outcomes* 7 rule, counting boxes
     by the CLI's definition: `^\s*[-*]\s+\[( |x|X)\]`, fences excluded, `proposal.md` boxes
     only in a `## Checklist` section closed by the next `^##\s` heading, `tasks.md` and
     `plan.md` boxes anywhere outside fences.
   - The trivial-lane contradiction is resolved (*Outcomes* 10) in `find-work` (lines 16-32
     against 132-134) and in `write-spec`'s step 0. That includes reworking step 0's "Six
     pillars, a brief and a plan around a two-character fix".
   - `/libretto-status`'s scan description names the three files.
3. **The change document.** Waits on 1, so a new-layout landing is gated from the first
   commit that can produce one, and on 2, so a new-layout change is visible while in flight.
   - `write-spec` grows `proposal.md` in place: `Targets:`, the `## Decision log` section
     replacing `decisions.md`, and the six pillars as level-2 headings. When phase 1 came
     from a tracker and wrote no `proposal.md`, it creates one carrying `Tracker: <key>`.
   - `write-spec`'s *Where specs live* tree shows both layouts, and its multi-capability
     paragraph names `spec-<capability>.md`.
   - `write-plan` writes `## Approach` and the `Durable decisions:` line into `proposal.md`.
   - `write-tasks` writes `## Checklist` into it and keeps the one-writer, standalone-box and
     mark-as-you-go rules.
   - `record-work` (lines 56-82), `review-work` (line 42), `agents/work-reviewer.md` and
     `build-and-check` find the contract, the approach and the durable decisions in either
     layout.
4. **Inline writing.** Waits on 3, because inline writing is written against the parts unit
   3 defines.
   - No skill launches `plan-writer` or `task-cutter`, and their agent files are deleted.
   - `write-spec` no longer launches a writer for `proposal.md`. Step 4b and *Resolving the
     markers* are retired for the single case.
   - Step 2b is narrowed to one writer per extra `spec-<capability>.md`, and only when
     context requires it, each writer given `proposal.md` by path (*Outcomes* 2).
   - `agents/work-reviewer.md` line 54 ("like spec-writer's write scope") is reworded.
   - `docs/PAYLOAD.md` is regenerated.
   - `check-payload` rows retired or reworded:
     - **retire:** 283 (single writer is a fan-out of one), 288 (`plan-writer` launch), 296
       (`plan-writer` tools), 324 (`task-cutter`)
     - **reword:** 281 and 287 (`decisions\.md` → the decision-log section's name) and 335
       (`never the list` stays, in whichever file holds the sentence)
     - **survive, reworded to `proposal.md`:** 274 (the fan-out vocabulary row, whose
       vocabulary now lives in `proposal.md` rather than a brief section) and 275 (the
       fan-out cost row)
     - **survive unchanged:** 284 and 285 (the `NEEDS CLARIFICATION` marker rows in
       `write-spec` and `agents/spec-writer.md`), because a fan-out writer can still leave a
       marker
5. **One stop, one fresh reviewer.** Waits on 3. It is one unit, not five, because removing
   the two old stops before the new one exists leaves a tree with no pre-code stop at all.
   - `agents/spec-reviewer.md` is added.
   - `review-spec` launches it, relays its findings unedited, carries the per-box gap
     question, and asks the stop natively. Its description moves from "the seam between
     phases 3 and 5" to before the stop.
   - `write-spec` (lines 713-722) and `write-tasks` (around line 208) stop asking for a
     go-ahead.
   - `commands/libretto-flow.md` routes write-spec → write-plan → write-tasks → review-spec,
     and carries the new stop table (lines 49-79) and gear sentence (line 133).
   - `commands/libretto-attacca.md` follows: its table (lines 24-25), line 36, its
     `decisions.md` entry (line 47), the `plan-writer` aside (line 86) and its *who acts on
     the findings* section.
   - `commands/libretto-next.md` (lines 60-61) names the one stop.
   - `check-payload` rows 293-294 (the stop rows for `write-spec` and `write-tasks`) are
     replaced by one row in `review-spec` plus two `check_absent` rows. Row 289 (attacca
     `decisions\.md`) is reworded.
6. **Vocabulary.** Waits on 3, 4 and 5. The `check-payload` row that fails on the acronym or
   the spelled-out method name anywhere under `skills/`, `agents/` or `commands/`. Today
   that row finds zero matches (observed).
7. **Plain doc edits, no criterion.** No capability governs these files, so they carry no
   `Proof:`. Waits on 3, 4 and 5, whose prose they describe.
   - `docs/FLOW.md`: rewritten where it describes phases 2–5, the stops, delegation and
     attacca, plus the one credit line for gentle-ai.
   - `docs/STATE.md:18`: the agent count. It says "7 items (spec-writer, work-reviewer, 5
     review lenses)", which is already stale today.
   - `docs/STATE.md:68`: names `spec-writer`.
   - `AGENTS.md:112`: "unchecked boxes in `.agents/changes/*/plan.md`".
   - `AGENTS.md:131`: "A plan that is 11/11 open".
   - `AGENTS.md:159-163`: "`plan.md` holds why … A plan that retires nothing declares it".
   - `AGENTS.md:58`: the `loop` comment "open box in that change's plan.md". The wording
     follows whatever the `cli` delta settles for `loop`.
8. **`spec-writer` reads the change document.** Waits on 3, which defines the parts it reads,
   and on 4, which narrows step 2b to the fan-out.
   - `agents/spec-writer.md`'s inputs are rewritten: it reads `proposal.md` by path (its
     `## Decision log`, its vocabulary and the primary capability's contract) instead of
     `brief.md` and `decisions.md`. Its description, its *Read the brief first* section and
     its "what the brief did not know" return are reworded to match.
   - The `check-payload` row that fails when `agents/spec-writer.md` names `brief.md` or
     `decisions.md` as an input is added, and observed red on the current file before the
     prose moves.

## Verification criteria

- **If** a staged commit deletes a `proposal.md` whose `HEAD` copy carries a `Targets:` line
  outside fences, and no capability spec's *Prior
  decisions* section differs between `HEAD` and the index, and the deleted copy carries no
  `Durable decisions: none` line, **then** `spec-drift --retired` **shall** exit 1 naming
  the change. **Where** the deleted copy declares `Durable decisions: none`, or a *Prior
  decisions* section moved, it **shall** exit 0. The test drives real temporary
  repositories, because the git plumbing is where a gate goes silent.
  Proof: skills/record-work/spec-drift --self-test
- **When** a staged commit deletes a `proposal.md` whose `HEAD` copy carries a `Queued:`
  line and no `Targets:` line outside fences, `spec-drift --retired` **shall** exit 0 and **shall not** report that change as retired. This is the
  case to break first: removing `Targets:` from the contract-lane fixture must turn its exit
  1 into exit 0, which proves the discriminator is what fires.
  Proof: skills/record-work/spec-drift --self-test
- **When** a staged commit deletes an old-layout `plan.md`, `spec-drift --retired` **shall**
  behave exactly as it did before this delta. The three existing `land` cases keep their
  expected exit codes, unedited.
  Proof: skills/record-work/spec-drift --self-test
- **When** a staged commit deletes both a change's `plan.md` and its contract-lane
  `proposal.md`, `spec-drift --retired` **shall** report that change once.
  Proof: skills/record-work/spec-drift --self-test
- **If** the only staged difference in a `proposal.md` carrying `Targets:` lies inside its
  `## Checklist` section, **then** default `spec-drift` **shall** still report staged code
  governed by that capability as code whose specification did not move. **When** anything
  outside `## Checklist` differs, it **shall** count as the capability moving, as today.
  Proof: skills/record-work/spec-drift --self-test
- **Where** a change document holds `## Verification criteria` followed by a `## Checklist`
  section of `- [ ]` boxes, `criteria_of` **shall** return the criteria and none of the
  boxes.
  Proof: skills/record-work/spec-drift --self-test
- **When** a change folder's `proposal.md` holds at least one box under `## Checklist`,
  `find-work` **shall** read it as that change's checklist in both scans, ahead of
  `tasks.md` and `plan.md`. **Where** it holds none, the skill **shall** fall back to
  `tasks.md`, then `plan.md`, as today. It **shall** count boxes by the CLI's definition: a
  line matching `^\s*[-*]\s+\[( |x|X)\]`, none inside a fence, in `proposal.md` only within
  a `## Checklist` section closed by the next `^##\s` heading. Proved as wiring only: the
  scan lines, the selection sentence and the box definition are present in the skill. Whether a session obeys them is behaviour,
  checked by running it.
  Proof: scripts/check-payload
- **When** `/libretto-status` reports, it **shall** delegate to `find-work`'s scan, and so
  count boxes by the same definition as the CLI, and **shall** name `proposal.md` among the
  files that scan reads, never restating the scan.
  Proved as wiring only.
  Proof: scripts/check-payload
- **When** step 0 sends a change down the trivial lane, `write-spec` **shall** direct that
  the change folder phase 1 created is deleted in the trivial lane's own commit, and
  `find-work` **shall not** claim the trivial lane never creates a `changes/` directory.
  Proved as wiring only: the instruction is present and the old sentence is absent.
  Proof: scripts/check-payload
- **Where** a change reaches phase 2 on the contract lane, `write-spec` **shall** direct the
  session to write the primary capability's contract into `proposal.md` itself: `Targets:`,
  the six pillars as level-2 headings, and every answer verbatim under `## Decision log`. It
  **shall not** direct the session to create `brief.md`, `decisions.md` or `spec.md` for a
  single-capability change. Proved as wiring only.
  Proof: scripts/check-payload
- **Where** a change amends more than one capability, `write-spec` **shall** direct that
  each extra capability's contract is written to a `spec-<capability>.md` beside
  `proposal.md` with its own `Targets:` line, and **shall** launch a writer for one only when
  the session's context cannot hold its research. Proved as wiring only.
  Proof: scripts/check-payload
- **When** phase 5 runs on the contract lane, `write-plan` **shall** keep the fork (two or
  three approaches as one native question, with the chosen and rejected logged verbatim).
  It **shall** direct the session to write `## Approach` and the `Durable decisions:` line
  into the change document, with no box in the approach. Proved as wiring only: the fork row
  and the `No checkboxes.` row survive, and the log row moves from `decisions.md` to the
  decision-log section.
  Proof: scripts/check-payload
- **When** the 5→6 seam runs on the contract lane, `write-tasks` **shall** direct the session
  to write the checklist under a level-2 `## Checklist` as the change document's last
  section, keeping `one badly cut box` and the one-writer rule. Proved as wiring only.
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` names `plan-writer` or
  `task-cutter`, or `agents/plan-writer.md` or `agents/task-cutter.md` exists, **then**
  `check-payload` **shall** fail naming it.
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` names `spec-writer` outside
  `write-spec`'s fan-out section and the agent file itself, **then** `check-payload`
  **shall** fail naming it. Answered 2026-10-21 (spec-writer survives for the fan-out case only), so this criterion stands as written.
  Proof: scripts/check-payload
- **If** `agents/spec-writer.md` names `brief.md` or `decisions.md` as an input, **then**
  `check-payload` **shall** fail naming the file. The row is observed red on today's agent
  file, which reads both, before its prose moves to `proposal.md`.
  Proof: scripts/check-payload
- **When** the contract lane's checklist is written, `review-spec` **shall** launch exactly
  one `spec-reviewer` with, by path, the change document, every `spec-<capability>.md`
  beside it, and the capability spec(s) their `Targets:` lines name, and with none of the
  conversation. It **shall** relay the reviewer's findings unedited. Proved as wiring only:
  the launch sentence and the agent name in `review-spec`.
  Proof: scripts/check-payload
- The `spec-reviewer` agent **shall** declare exactly the tools `Read, Grep, Glob, Bash,
  Skill`, and **shall** carry the instruction never to edit what it reviews. `review-spec`
  **shall** keep `never rewrites the spec`.
  Proof: scripts/check-payload
- The `spec-reviewer` agent **shall** report, for each box under `## Checklist`, what the
  contract and the approach fail to answer for building it, naming the gap and the part of
  the document that should have answered it, and **shall not** propose the answer. Proved as
  wiring only: one row on the agent's per-box sentence.
  Proof: scripts/check-payload
- **When** the fresh reviewer has returned on a contract-lane change, `review-spec` **shall**
  ask the one stop natively, on the whole change document, and `write-spec` and
  `write-tasks` **shall not** ask for a go-ahead. One row each: the native question is
  present in `review-spec`, and each old go-ahead sentence is absent from the other two.
  The rows are written against the skills' own wording and observed red on the current
  files before the prose moves. **Ceiling named:** that no other stop exists stays
  uncheckable, the limit the stop table has always named.
  Proof: scripts/check-payload
- **When** `/libretto-flow` routes the contract lane, it **shall** route `write-spec`,
  `write-plan`, `write-tasks` and then `review-spec`, keep `before the stop, not after`, and
  state the gear as one stop on the contract lane, none on the trivial lane, and phase 8's
  question in both. `/libretto-attacca` **shall** still route to `review-spec`. Proved as
  wiring only.
  Proof: scripts/check-payload
- **While** a run is `/libretto-attacca`, the interview, the fork and the stop **shall**
  become `(assumed)` entries in the change document's `## Decision log`, each naming what
  changes if it is wrong. The command **shall** name the one stop and the push as what the
  invocation answered. Proved as wiring only: the existing `decisions\.md` row moves to the
  decision-log section's name.
  Proof: scripts/check-payload
- **When** phase 8 lands a new-layout change, `record-work` **shall** direct the session to
  do three things in one commit: apply `proposal.md`'s contract and each
  `spec-<capability>.md` onto the capability spec its `Targets:` names, retire the
  `## Prior decisions` entries that outlive the change, and delete the folder. For an
  old-layout change it **shall** read `spec.md` and `plan.md` as today. Proved as wiring
  only.
  Proof: scripts/check-payload
- **If** any file under `skills/`, `agents/` or `commands/` contains the acronym quoted in
  the 2026-10-21 naming answer as a case-sensitive whole word, or the method's spelled-out
  name in any case, **then** `check-payload` **shall** fail naming the file. The row is
  observed red on a hand-broken copy before it is trusted.
  Proof: scripts/check-payload
- **When** an agent is added to or removed from `agents/`, `check-payload` **shall** fail
  until `docs/PAYLOAD.md` is regenerated. This is the existing index criterion, exercised by
  this delta's changes to the agent set rather than changed by them.
  Proof: scripts/check-payload
