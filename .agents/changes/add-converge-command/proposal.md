# add-converge-command

Tracker: none
Targets: payload
Durable decisions: every entry in Prior decisions below

## The ask, verbatim

> `converge` estilo spec-kit: re-derivar los boxes abiertos comparando el codebase real
> contra spec/plan — el complemento natural de `loop` cuando el plan derrapó.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Assess the actual codebase against a change's spec/plan/tasks and append the remaining
work to the checklist, for plans that drifted from what got built. Prior art:
`/speckit.converge`.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed unless it says measured)

- Q: where does converge live — a payload command and skill launching one fresh reader, a
  CLI that spawns a session like `loop`, or a mechanical CLI plus a reading? → A: **a
  payload command, `/libretto-converge`, routing to a `converge` skill that launches one
  fresh `converge-reader` agent.** The mechanical part — running the cited proofs — is done
  by the session with `spec-drift --anchors` and the project's own test runner, never by a
  new binary subcommand (assumed). Whether a box is done is a reading, and the binary
  cannot read prose; `loop.go` already records that a runner which parsed a task's text
  "would be a second opinion about what a task is". If wrong — the user wants it
  scriptable — a `libretto converge` subcommand spawns one session with
  `/libretto-converge <change>` as its prompt, the way `loopPrompt` does, and this skill
  stays the rule.
- Q: a new agent, or reuse `work-reviewer` or `spec-reviewer`? → A: **a new agent,
  `converge-reader`** (assumed). `work-reviewer` reviews a finished diff against three
  fixed questions and `spec-reviewer` reads a contract before code exists; handing either
  a per-box brief contradicts the agent file it runs under. If wrong, `work-reviewer`
  gains a second mode and its file says so.
- Q: does the reader run anything? → A: **no. Its tools are `Read, Grep, Glob, Skill`, and
  no `Bash`** (assumed). A reader that cannot run a proof cannot report one as passed, so
  its verdict can never be mistaken for evidence. If wrong, it gains `Bash`, and its runs
  still close nothing: the session re-runs every proof it closes on.
- Q: who runs the proofs that close a box? → A: **the session, in the foreground, reading
  each `Proof:` line from the contract itself — never from the reader's quote of it**
  (assumed). The checklist has one writer, and the writer is the one that must have
  observed the pass. If wrong, the reader's runs are relayed like `work-reviewer`'s, and a
  box still closes only on a run the session read.
- Q: what happens to a drifted box? → A: **reported, never rewritten.** Its text and its
  state stay as they are, and the report gives both readings: what the box says, and
  what the code did instead, with a file and line (assumed). A box traces to the contract,
  and rewriting it to match the code ratifies the drift without the contract moving — an
  edit to what the change promises, made by the back door. If wrong, converge writes one
  indented `Drift:` line under the box, still never changing the box itself.
- Q: does converge add boxes? → A: **yes, one per criterion no box closes, and only when
  that criterion's proof does not pass when the session runs it.** Appended after the
  checklist's last box, open, carrying `Closes:` and `(converge <date>)` (assumed). A
  criterion whose proof already passes gets no box: `write-tasks` already says a task whose
  criterion is met is not a task. If wrong, converge appends nothing and lists uncovered
  criteria in its report only.
- Q: what about an outcome or a scope item no criterion covers? → A: **reported as a
  contract gap, never boxed** (assumed). A box names the criterion that closes it; a box
  with no criterion is a box nobody can close, and writing the criterion is editing the
  contract, which converge never does. If wrong, such a box is appended with
  `Closes: none — contract gap` and phase 2 is named as its owner.
- Q: closed boxes — re-checked, or left alone? → A: **left alone, never re-run, never
  unchecked** (assumed). The ask is about open boxes. A closed box whose proof fails today
  is a regression, and `libretto metrics` already counts a reopened box as the signal of
  one. If wrong, converge re-runs closed boxes' proofs and reports failures, still never
  unchecking.
- Q: does converge commit? → A: **no.** It writes the checklist and stops, and its report
  says the marks are not recorded until the checklist is committed (assumed). **Reversed
  after review-spec — see the next session.**
- Q: what if the working tree is dirty? → A: **converge stops before the reader, naming
  the files** (assumed). A proof run against uncommitted code proves nothing about what
  was recorded — `evidence`'s rule — and a dirty checklist file would mix converge's
  edits with somebody else's. If wrong, a dirty tree runs report-only: the reader runs,
  nothing is closed and nothing is appended. *Narrowed to tracked files after
  review-spec.*
- Q: which change, when none is named? → A: **the one `find-work`'s in-flight scan
  returns; if it returns several, one native question asking which; if none, say so and
  stop** (assumed). The scan is `find-work`'s, as `/libretto-status` delegates it. If wrong,
  the change name becomes a required argument.
- Q: is converge wired into the flow or into `loop`? → A: **neither.** It is a door the
  user opens, like `/libretto-status`; `libretto-flow` and `libretto-attacca` do not route
  to it, and `loop`'s stall message does not name it (assumed). A hint in `loop` is a
  `cli` change, and widening this into a second capability is what the scope exists to
  stop. If wrong, `loop`'s two-rounds stop message gains one line naming
  `/libretto-converge <change>`, under the `cli` spec.
- Q: does `cmd/libretto/recommend.go` gain a model for `converge-reader`? → A: no — a `cli`
  change, and an agent with no recommendation is silence by design (assumed). If wrong,
  one `sonnet`/`high` row, "a cheap miss closes nothing but drifts the report".
- Q: does another capability move? → A: **no `spec-<capability>.md`.** `readme`'s Outcome 6
  already requires every command in `README.md`, and that promise does not change — it
  gains a line to satisfy (measured: the promise is
  cmd/libretto/readme_test.go TestEveryCommandIsInTheReadme). If wrong, a `spec-readme.md`
  carries the new row.
- Q: the names? → A: `/libretto-converge`, skill `converge`, agent `converge-reader` — the
  `retro`/`libretto-retro` precedent (assumed). If wrong, three renames before landing.
- Q: the fork — how is it built? → A: **A, a payload command and skill with one fresh
  reader; the session runs every proof it closes on** (assumed, the recommended option).
  B — a `libretto converge` subcommand spawning a session, like `loop` — lost: `loop`
  sits outside the session because its value is a fresh context per box, and converge is
  one reading, which one fresh agent already gives inside the session that must be the
  writer anyway. C — a mechanical CLI that runs the cited proofs plus a reading — lost: the
  binary would have to parse `Closes:` prose to map a box to its criteria, and know how to
  run a citation in a project that is not Go; the payload installs into projects of any
  language, and the session already knows how its project runs a test.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: what counts as a passing citation? → A: **two kinds, two rules.** A test-named
  citation, `<file> TestX`, is run with `-v` and passes only when the output carries
  `--- PASS: TestX`; a name that does not exist prints no such line. A file-only citation —
  `scripts/check-payload`, `skills/record-work/spec-drift --self-test` — passes when that
  command exits 0 and prints no `FAIL` line. "Whatever its exit status" applies to the
  test-named kind only (assumed). If wrong — another runner prints a different pass line —
  the skill names the line per runner, still never trusting an exit status alone for a
  named test.
- Q: what is an uncovered criterion? → A: **one of the change's own criteria only** — the
  `## Verification criteria` of `proposal.md` and of every `spec-<capability>.md` beside
  it — that no box's `Closes:` names. Never a criterion of a targeted capability spec
  (assumed). The capability spec is what is true already; its criteria belong to changes
  that landed. If wrong, converge would also read the targeted specs and box every failing
  promise there, which turns a checklist tool into a regression sweep.
- Q: what does the reader read? → A: **the whole tree as it is now, the change folder, and
  the capability specs the change's `Targets:` lines name. Never a diff range** (assumed).
  Of `review-work`'s step-1 world, the contract and the folder apply; the diff range does
  not, because drift is a property of the code as it stands, wherever it was written. If
  wrong, the diff `<base>...HEAD` is added as a hint and the tree stays the subject.
- Q: which of `find-work`'s scans does converge use? → A: **source 1 only — the open-box
  scan and the unmerged branches. Never the queue scan**, because a queued capture has no
  checklist to converge (assumed). If wrong, nothing changes in practice: a queued capture
  holds no box and would be refused anyway.
- Q: what makes the tree dirty? → A: **tracked modifications only:
  `git status --porcelain --untracked-files=no` listing anything.** An untracked file — a
  generated `wiki.html`, a scratch note — never stops converge (assumed). Untracked files
  are not code any proof runs against unless something imports them, and stopping on one
  this tree carries today would stop converge on every run. If wrong, untracked files under
  `Governs:` paths also stop it.
- Q: where does the converged mark go? → A: **on the box's last line**, listing every
  citation that passed: `— converged <date>: <citation>, <citation> passed` (assumed). A
  box wraps over several lines, and the mark on its last one is where a reader finishes
  the box. If wrong, the mark goes on an indented line of its own beneath the box.
- Q: how is a box mapped to its criteria? → A: **through the box's `Closes:` text,
  matched against the contract's criteria.** Where the match is not unambiguous — no
  criterion, or more than one candidate — the box stays open and is reported as
  `unmapped`. An old-layout box with no `Closes:` stays open and is reported as `unmapped`
  (assumed). Guessing a mapping is the reading a reader's verdict is not allowed to stand
  in for. If wrong, the reader proposes a mapping and the session accepts it only when the
  cited criterion's text is quoted exactly.
- Q: does converge commit? → A: **yes — reversed.** payload's rule is that a marked box
  that was never committed was never marked, and the checklist ships in the same commit
  as the task that closed it. So converge ends by committing the checklist file alone, in
  one `docs(<change>): converge the checklist` commit, staging nothing else. No push, no
  branch. Nothing marked or appended: no commit (assumed). If wrong — the user wants to
  read the edit first — converge leaves the checklist file staged and the report says the
  marks are unrecorded until it is committed.
- Q: does every criterion with a half no row pins get a second row? → A: **yes, where a
  stable phrase exists; otherwise a "Ceiling named:" that names the unpinned half.** No
  half is left silent (assumed). If wrong, the extra rows go and each ceiling widens.
- Q: where does box 3's hand-run live? → A: **in a throwaway git repository made with
  `mktemp -d`, never on this branch.** Its failing proof is a cited Go test that asserts
  false, or a cited script that exits 1 (assumed). A fixture committed here would be a
  change folder with open boxes that every phase 1 reports as work in flight. If wrong,
  the fixture becomes a self-test case, which needs a runner this payload does not ship.

## Outcomes

- `/libretto-converge [<change>]` reads one change's open boxes against the codebase and
  ends with each open box in exactly one state — **confirmed open**, **closed with observed
  proof**, **reported as drifted**, or **reported as unmapped** — and with one appended box
  per criterion of the change's own that no box closes and whose proof does not pass.
- A box is closed by converge only after the session ran every proof its criteria cite, in
  the foreground, and read each one passing by the rule for its kind. The closed box's last
  line carries where that evidence is: every citation that passed and the date.
- A drifted box is untouched in the file and named in the report, with what the box says
  and what the code does instead.
- The judgment is made by one fresh `converge-reader`, which never had the conversation,
  reads the tree as it is, and cannot run anything; the checklist is written only by the
  session.
- It works on both layouts: the `## Checklist` of `proposal.md`, and `tasks.md` or
  `plan.md` for a change begun before.
- When it marked or appended anything, it ends with one commit holding the checklist file
  alone. When it did neither, it commits nothing. It never pushes.

## Scope boundaries

**In:**

- `commands/libretto-converge.md` — routes to the skill, describes none of it
- `skills/converge/SKILL.md` — version 1.0: the preconditions, which change, the world
  given to the reader, the mapping, the proof runs, the writes, the commit, the report
- `agents/converge-reader.md` — the reading, and its return shape
- `scripts/check-payload` — `check_wiring` and `check_absent` rows pinning the criteria
  below
- `README.md` — one line for the command, keeping the file at or under its 380-line
  ceiling
- `docs/PAYLOAD.md` — regenerated by `scripts/check-payload --index`

**Out:**

- **rewriting a drifted box**, or annotating it. Reported only.
- **editing the contract or the approach** — every pillar, `## Approach`, the decision log,
  `Targets:`, `Durable decisions:`. converge writes to the checklist and nothing else.
- **unchecking or re-running a closed box.**
- **pushing, creating a branch, or staging anything but the checklist file.** Its one
  commit holds that file alone.
- **guessing a mapping.** A box whose `Closes:` matches no criterion unambiguously is
  reported `unmapped`, never closed.
- **the targeted capability specs' criteria as uncovered work.** Uncovered means the
  change's own criteria.
- **a diff range as the reader's input.** The tree as it is.
- **a binary subcommand**, and any change under `cmd/` or `internal/` — including `loop`'s
  stall message and `recommend.go`. Each is a `cli` follow-up, named in the decision log.
- **routing from `libretto-flow` or `libretto-attacca`.** A door, not a phase.
- **boxing a contract gap** — an outcome or scope item no criterion covers. Reported.
- **a second definition** of a box, of which file holds the checklist, or of how a
  citation is run. `find-work`, `review-work` and `work-reviewer` own those; converge
  points to them.
- **`docs/FLOW.md`.** It describes the eight phases; a side door is not one. One line
  there is free to add later.

## Constraints

- Payload only: markdown under `skills/`, `agents/`, `commands/`, plus `scripts/check-payload`
  rows, `README.md` and the generated `docs/PAYLOAD.md`. No dependency is added.
- `scripts/check-payload` exits 0: frontmatter whose `name:` matches its directory or file,
  every `Skill(skill="…")` and `skills/<name>/` reference resolving, no `scripts/` or
  `bin/` path in a skill, no `~/.claude/` path, no bare product name addressed in prose,
  the host-neutral marker (`native prompt`, `host's own`) in every file mandating
  `AskUserQuestion` or `Skill(skill=…)`, and **the borrowed method's acronym nowhere** under
  `skills/`, `agents/` or `commands/`.
- `docs/PAYLOAD.md` is regenerated, never hand-edited; the gate fails on drift.
- `README.md` is at its 380-line ceiling today (measured: `wc -l` reads 380), and
  cmd/libretto/readme_test.go TestReadmeStaysShort fails above it. The command's line is
  added without growing the file.
- `spec-drift` is reached from a skill as `<skill-base>/../record-work/spec-drift`, the
  form `write-spec` uses — never through `scripts/`.
- `check_wiring` rows are line-scoped: each pinned phrase sits on one line of its file.
- The checklist is found and boxes are counted by `find-work`'s rule — a box is a line
  matching the CLI's `planLine`, outside fences, under `## Checklist` in `proposal.md`, else
  every box in `tasks.md`, else `plan.md` — the same answer `checklistPath` and `readPlan`
  give in `cmd/libretto/loop.go`.
- A named test that does not exist passes under `go test -run`: measured 2026-10-22,
  `go test ./cmd/libretto/ -run '^TestDoesNotExistAnywhere$' -v` exits 0 printing
  `PASS` and `[no tests to run]`, and no `--- PASS:` line. A test-named citation's exit
  status alone is not a pass.
- This tree carries an untracked generated file today (`.agents/specs/wiki.html` on
  `main`), which is why dirtiness reads tracked files only.

## Prior decisions

- **converge is a payload command and skill launching one fresh reader; no binary
  subcommand.** Assumed 2026-10-22, attacca run. Whether a box is done is a reading, and
  the binary cannot read prose — `loop.go` refuses to parse a task's text for that reason.
  The mechanical half is the session running each cited proof, which needs no new code.
  **What would reopen it:** wanting converge headless, which is a subcommand that spawns
  `/libretto-converge` exactly as `loop` spawns `/libretto-attacca`.
- **A reader's verdict never closes a box.** Assumed 2026-10-22, attacca run. The
  `converge-reader` has no shell, so it cannot run a proof, and the session closes a box
  only on runs it read itself, citations taken from the contract rather than the
  reader's quote. If wrong, the reader gains `Bash` and nothing else changes.
- **A citation passes by the rule for its kind.** Assumed 2026-10-22, after review-spec. A
  test-named citation passes only on a `--- PASS: <TestName>` line in `-v` output; a
  file-only citation passes on exit 0 with no `FAIL` line. If wrong, the pass line is
  named per runner.
- **A box maps to criteria only through an unambiguous `Closes:` match.** Assumed
  2026-10-22, after review-spec. Anything else is `unmapped`, open and reported. If wrong,
  an exact quote from the reader is accepted as a mapping.
- **The reader reads the tree as it is, never a diff range; uncovered means the change's
  own criteria.** Assumed 2026-10-22, after review-spec. Drift is a property of the code
  that exists, and a capability spec's criteria belong to changes already landed. If
  wrong, the diff becomes a hint and the targeted specs stay out.
- **Drift is reported, never rewritten.** Assumed 2026-10-22, attacca run. Rewriting a box
  to match the code is changing the change's promise without touching its contract. If
  wrong, an indented `Drift:` line beneath the box, the box itself untouched.
- **A box is appended only for a criterion no box closes and whose proof does not pass.**
  Assumed 2026-10-22, attacca run. A criterion that passes already is not a task, and an
  outcome with no criterion is a contract gap, reported. If wrong, nothing is appended and
  the report lists uncovered criteria.
- **converge commits the checklist file alone, never unchecks, and stops on tracked
  modifications.** Assumed 2026-10-22, after review-spec, reversing the run's first
  answer. One `docs(<change>): converge the checklist` commit when anything was marked or
  appended, because a mark never committed was never marked; no push, no branch. No
  uncheck, because a closed box's regression is `metrics`' reopened count. Tracked
  modifications stop it, because a proof against unrecorded code proves nothing about
  what was recorded; untracked files do not. If wrong, the checklist is left staged and
  the report says it is unrecorded.
- **Every criterion below is proved as wiring.** Assumed 2026-10-22, attacca run, the way
  every payload prose criterion is. **Ceiling named:** `scripts/check-payload` proves each
  mandate is *present* in its file, never that a session obeyed it; obedience is
  behaviour, observed by running converge on a drifted change, and surfaces in the 6→7
  review or not at all.

## Task breakdown

1. The read-only converge: the command, the skill's preconditions (no tracked
   modification, which change), the world handed to one `converge-reader`, the reader's
   return shape, and the report — drift and the reader's three verdicts named, nothing
   written. The README line and the regenerated index.
2. Closing on observed proof: `spec-drift --anchors` first, each closable box mapped
   through its `Closes:` or reported `unmapped`, its cited proofs run by the session and
   read by the rule for their kind, the mark on the box's last line, and the one commit of
   the checklist file alone.
3. Appending the uncovered: each of the change's own criteria no box closes, its proof run,
   and one box appended when it does not pass; a contract gap reported, never boxed.

## Verification criteria

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

## Approach

### Summary

A side door into the flow, built like `/libretto-retro` and `/libretto-status`: a command
that routes, a skill that holds the rule, and — like `review-work` — one fresh agent that
reads and never writes. The `converge-reader` reads the tree as it is and returns, per open
box, one of three verdicts plus the change's own criteria no box closes. The session then
does the only part that counts as evidence: it maps each closable box through its
`Closes:`, runs the cited proofs itself after `spec-drift --anchors`, reads each by the
rule for its kind, and writes the checklist — closing what passed, appending a box for
each uncovered criterion that does not, leaving every drifted or unmapped box as it was
and naming it in the report. When it wrote anything, it commits the checklist file alone.

### Technical context

- Payload markdown, read by the host as a skill, an agent and a command. Installed by
  symlink; on OpenCode the agent is a derived file the linker writes, which needs nothing
  from this change.
- **Reused, never redefined.** `find-work`'s source 1 — its open-box scan and its unmerged
  branches, never the queue scan — finds the in-flight changes and defines the box and
  which file holds the checklist: the CLI's `planLine`, and `checklistPath`'s order,
  `proposal.md`'s `## Checklist` when it holds a box, else `tasks.md`, else `plan.md`.
  `work-reviewer` shows how a citation is run. `write-tasks` owns *One writer* and what a
  box records.
- **The reader's world is part of `review-work`'s, not all of it.** From `review-work`'s
  step 1, the contract applies — `proposal.md`, every `spec-<capability>.md` beside it, or
  an old-layout delta, and the capability specs their `Targets:` lines name — and so does
  the change folder, for the approach. **The diff range does not**: the reader reads the
  whole tree as it is now, because drift is a property of the code that exists, wherever
  it was written.
- **Uncovered means the change's own criteria**: the `## Verification criteria` of
  `proposal.md` and of each `spec-<capability>.md` beside it — an old-layout change's
  `spec.md` delta in their place. The targeted capability specs are read for context and
  never contribute an uncovered criterion.
- **Two kinds of citation.** Test-named, `<file> TestX`: run with `-v`, passed only on a
  `--- PASS: TestX` line. File-only, `scripts/check-payload` or
  `skills/record-work/spec-drift --self-test`: passed on exit 0 with no `FAIL` line.
- `spec-drift --anchors` resolves every `Proof:` down to the test name, and fails on one
  whose test is absent.
- `scripts/check-payload`: `check_wiring <file> <pattern> <what>` and `check_absent`, one
  line each, plus the generic checks under *Constraints*.
- Blast radius: three new files (`commands/libretto-converge.md`, `skills/converge/SKILL.md`,
  `agents/converge-reader.md`), and `scripts/check-payload`, `README.md`,
  `docs/PAYLOAD.md`. At landing, the delta onto `.agents/specs/payload/spec.md`: an
  *Outcomes* bullet for the command, these criteria, the *Prior decisions*.
- Gates: the six in `AGENTS.md`. The Go suite moves only through
  `TestEveryCommandIsInTheReadme` and `TestReadmeStaysShort`.

### The approach

1. **The command**, `commands/libretto-converge.md`: frontmatter `description:`, the
   `Skill(skill="converge")` line, the host-neutral sentence `libretto-retro` uses, and
   one line of why — "`loop` drives a checklist; this re-derives one that stopped
   describing the code". Input: `$ARGUMENTS`, a change name or empty.
2. **The skill**, `skills/converge/SKILL.md`, in this order:
   1. **Preconditions.** `git status --porcelain --untracked-files=no`; anything listed
      means *a dirty tree stops converge before the reader*, naming the files. Untracked
      files are not read.
   2. **Which change.** Named: that one, refused if its folder holds no checklist box.
      Empty: `Skill(skill="find-work")`, *source 1 only, never the queue* — its open-box
      scan and its unmerged branches; a queued capture has no checklist. One in flight:
      that one. Several: *ask which one*, one `AskUserQuestion` call — the host's native
      prompt, or in conversation where there is none — and start nothing until answered.
      None: say so in one line and stop. *`find-work` defines the box* and which file
      holds the checklist; this skill restates neither.
   3. **The reader.** Launch *exactly one `converge-reader`*, foreground, given by path the
      tree as it is now, the change folder and the capability specs its `Targets:` lines
      name — *never a diff range*, and *none of the conversation*. Silence is not a
      verdict: relaunch once, then stop — `review-work`'s rule.
   4. **The mechanical half.** Run `<skill-base>/../record-work/spec-drift --anchors`,
      redirected to a file with its exit read — *`spec-drift --anchors` before any proof*;
      any citation of this change it names as unresolved is not passed.
   5. **Closable boxes.** For each:
      - **Map it.** Match the box's `Closes:` text against the contract's criteria. No
        match, more than one candidate, or no `Closes:` at all — an old-layout box — and
        the box stays open, *reported as `unmapped`*.
      - **Run it.** Take each mapped criterion's `Proof:` line from the contract — not from
        the reader — and run it in the foreground the way `work-reviewer` runs a citation,
        output to a file and read. A test-named citation runs with `-v` and passes only
        when its output carries `--- PASS: <TestName>` for that test, *whatever its exit
        status*. A file-only citation passes when it *exits 0 and prints no FAIL line*.
      - **Mark it.** Every citation passed: mark `[x]` and append, on *the box's last
        line*, `— converged <date>: <citation>, <citation> passed`. Otherwise the box
        stays open and the report names the citation that did not pass. *A reader's
        verdict never closes a box.*
   6. **Drifted boxes.** *Drift is reported, never rewritten*: the box is not touched, and
      the report carries the box's words *beside the reader's file and line*.
   7. **Uncovered criteria.** *Uncovered means the change's own criteria* — the
      `## Verification criteria` of `proposal.md` and of each `spec-<capability>.md`
      beside it that no box's `Closes:` names. For each, run its proof by the same rules.
      Not passed: one open box *appended after the checklist's last box* —
      `- [ ] **<criterion, in a phrase>** (converge <date>). Closes: <criterion> on
      <citation>. Waits on: nothing.` Passed: it *passes already, so no box is appended*,
      and the report says so. A contract gap the reader names is reported only.
   8. **The writes.** converge *writes to the checklist and nothing else*, and *never
      unchecks a closed box*.
   9. **The commit.** Anything marked or appended: `git add` the checklist file alone and
      commit it as *`docs(<change>): converge the checklist`* — it *stages nothing else*.
      If *nothing was marked or appended, no commit*. It *never pushes and never creates a
      branch*.
   10. **The report**: one line per open box with its state — open, closed with its
       citations, drifted, unmapped — the proofs run with pass or fail, boxes appended,
       contract gaps, and the commit's hash or "no commit".
3. **The agent**, `agents/converge-reader.md`: `tools: Read, Grep, Glob, Skill`; invokes
   `evidence` first; reads the tree as it is, the change folder and the targeted capability
   specs, and never a diff. Returns, for every open box, *exactly one of open, closable or
   drifted*:
   - `open — <box head>: <what is missing, file or symbol>`
   - `closable — <box head>: <criteria it closes>; code at <file:line>`
   - `drifted — <box head>: says <…>; the code does <…> at <file:line>`

   then `uncovered — <criterion head>: <its Proof: line as written>` for each of **the
   change's own criteria** no box's `Closes:` names — never a targeted capability spec's —
   then `gap — <outcome or scope item>` for each outcome no criterion covers, or the
   sentence **"nothing found"**. It never edits anything, never commits, never judges
   style; a suspicion it cannot settle is returned marked unverified.
4. **`scripts/check-payload`**: one row per pinned phrase named under *Verification
   criteria*, after the existing payload rows, each forced red on a throwaway copy of the
   file it pins.
5. **`README.md`**: `/libretto-converge <change>` joins the door list with a comment
   "re-derive a drifted checklist against the code", and the file stays at 380 lines or
   fewer — a line folded elsewhere in the same section, never the ceiling raised.
6. **`docs/PAYLOAD.md`**: `scripts/check-payload --index`.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — `libretto converge`, a subcommand spawning one session as `loop` does | `loop` is outside the session because each box needs a fresh context; converge is one reading, which one fresh agent gives inside the session that has to write the checklist anyway. Adds a `cli` capability change and a dependency on a headless host on PATH |
| C — a mechanical subcommand running each box's proofs, plus a reading | the binary would parse `Closes:` prose to map a box to its criteria — the second opinion about a task `loop.go` refuses to be — and would have to know how to run a citation in a project that is not Go |
| reuse `work-reviewer` | its file fixes three questions about a finished diff; a per-box brief contradicts the agent it runs under |
| reuse `spec-reviewer` | it reads a contract before code exists and reports what the document fails to answer — the opposite direction |
| the reader runs the proofs | then a reader's report could close a box; the one writer must be the one that observed the pass |
| the reader reads the diff range, as `work-reviewer` does | drift is in the code as it stands; a box can be satisfied, or contradicted, by code the branch never touched |
| a test-named citation passes on exit 0 | a test that does not exist exits 0 (measured) |
| the reader proposes the box-to-criterion mapping | a guessed mapping is the reading that is not allowed to close a box; `unmapped` is honest |
| uncovered includes the targeted capability specs' criteria | turns converge into a regression sweep of everything already landed |
| rewrite a drifted box to match the code | ratifies drift without the contract moving |
| a `Drift:` line under a drifted box | an edit to a box whose meaning converge does not own; the report carries it. Kept as the fallback if report-only proves too easy to lose |
| append a box for every uncovered criterion, passing or not | a box whose criterion is met is not a task (`write-tasks`) |
| leave the checklist uncommitted | a mark never committed was never marked — `write-tasks` and `payload` |
| commit the checklist with whatever else is staged | a converge commit must be reviewable alone; it stages nothing else |
| stop on untracked files too | this tree carries an untracked generated `wiki.html` today; converge would never run |
| re-run closed boxes and uncheck failures | widens the ask from open boxes to regressions; `metrics` already counts reopened boxes |
| no reader — the session reads the code inline | the session running converge is often the one whose plan drifted, and reading a whole codebase into the writer's context is what the reader is for |

### Risks

| Risk | What catches it |
|---|---|
| a cited Go test that does not exist passes `go test -run` and closes a box | `spec-drift --anchors` runs first and fails on an absent test name; and a test-named citation passes only on its `--- PASS:` line. Both pinned by a criterion |
| a file-only citation exits 0 while reporting a failure | it must also print no `FAIL` line; pinned |
| the reader calls "closable" what is half-built | it cannot close anything; the session's run of the cited proof decides |
| a box's `Closes:` is vague, or an old-layout box has none | it is not mapped by guess: it stays open, reported `unmapped`, and the user sharpens the `Closes:` |
| a proof that passes for the wrong reason closes a box | not caught here — the same weakness the proof has everywhere. `review-work` re-runs the proofs at 6→7 |
| a drifted box stays open and the next `loop` session builds it as written, against the code | the report names every drifted box and says `loop` would rebuild it; resolving drift is the user's, at the contract. **Accepted**, with the `Drift:` line as the named replacement |
| the converge commit sweeps in something else the user staged | a staged file is a tracked modification, so the precondition stops converge before it begins; and the skill commits with the checklist path as the commit's only pathspec |
| `README.md` crosses 380 lines | `TestReadmeStaysShort` fails the gate |
| a `check_wiring` row green the moment it is written | each row forced red on a throwaway copy before it is believed — the discipline the existing rows record |

### Validation and rollback

`scripts/check-payload` carries every criterion; each new row is proved by deleting its
pinned phrase from a throwaway copy and watching the gate fail, then restoring it. The
`check_absent` row is forced red by pasting the `planLine` regex into the copy. The README
change is proved by `TestEveryCommandIsInTheReadme` going red before the line is added —
it will, the moment `commands/libretto-converge.md` exists — and green after, with
`TestReadmeStaysShort` green beside it.

**What would pass for the wrong reason:** every row here proves a sentence is present. None
proves a session obeys it. Behaviour is checked once by hand during phase 6, **in a
throwaway git repository made with `mktemp -d`, never on this branch**: a change folder
with three boxes — one whose proof passes, one whose proof fails, one describing a function
the code named differently — and one criterion no box closes, run through
`/libretto-converge`, its report and its commit read against the expected outcomes. The
failing proof is a cited Go test that asserts false, or a cited script that exits 1. The
repository is discarded after; what was observed is reported in phase 7, not gated.

Rollback: one revert; three new files go, three edited files return. Nothing migrates.

### Complexity deliberately kept

- **A new agent rather than a brief.** A fresh agent with its own file is three paragraphs
  more than an inline reading; it is what keeps the reader's verdict out of the writer's
  context, and what makes "no `Bash`" a declared, gated fact rather than a promise.
- **`spec-drift --anchors` and a pass rule per citation kind, both.** Either alone misses a
  case: anchors does not run the test, and the pass rule is the session reading output.
  Two cheap guards on the one failure that would close a box falsely.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [x] **A read-only converge reports every open box's state.** The command, the skill's
      preconditions on tracked modifications, change selection by `find-work`'s source 1,
      the reader launch on the tree as it is and the report; the `converge-reader` agent
      and its return shape; the README door-list line within 380 lines; `docs/PAYLOAD.md`
      regenerated; every `check_wiring` and `check_absent` row its criteria name, each
      forced red once. Writes and commits nothing yet. Closes: the command-routing
      criterion, the exactly-one-reader criterion, the reader-tools criterion, the
      three-verdicts criterion, the dirty-tree criterion, the ask-which-one criterion, the
      find-work-defines-the-box criterion, the no-restated-box-pattern criterion and the
      drift-is-reported criterion — all on scripts/check-payload. Waits on: nothing.
- [x] **Converge closes a box only on proof the session observed, and commits it.** In the
      skill: `spec-drift --anchors` first; each closable box mapped through its `Closes:`
      or reported `unmapped`; its cited proofs run from the contract and read by the rule
      for their kind; the mark on the box's last line listing every citation; writes
      limited to the checklist; no uncheck; the one commit of the checklist file alone,
      none when nothing changed, no push and no branch. Every row its criteria name, each
      forced red once. Closes: the reader-never-closes criterion, the anchors-first
      criterion, the test-named-citation criterion, the file-only-citation criterion, the
      unmapped criterion, the converged-mark criterion, the checklist-only criterion, the
      never-uncheck criterion, the commit-the-checklist-alone criterion, the
      no-commit-when-nothing-changed criterion and the never-push criterion — all on
      scripts/check-payload. Waits on: box 1.
- [ ] **Converge appends a box for each of the change's own uncovered criteria whose proof
      does not pass.** In the skill: uncovered defined as the change's own criteria only;
      the reader's `uncovered` lines run by the same rules; an open box marked
      `(converge <date>)` appended after the last box when the proof fails, nothing when
      it passes, contract gaps reported. Every row its criteria name, each forced red once.
      Then the one hand-run in a throwaway `mktemp -d` git repository with a passing, a
      failing and a drifted box plus one uncovered criterion, its report and commit read
      against the expected outcomes and carried to phase 7. Closes: the uncovered-means-own
      criterion, the append-when-failing criterion and the no-box-when-passing criterion —
      all on scripts/check-payload. Waits on: box 2, which defines how a proof is mapped,
      run, read and committed.
