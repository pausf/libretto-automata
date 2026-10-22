# add-queue-command

Tracker: none
Targets: cli
Durable decisions: every entry in Prior decisions below, and the four in spec-payload.md

## The ask, verbatim

> `libretto queue` en el binario. El CLI es deliberadamente ciego a la cola, pero `wiki`
> YA parsea las líneas `Queued:` — el dato está adentro, solo falta el comando de
> lectura. Barato y mejora el flujo diario.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

A read-only `libretto queue` listing queued ideas oldest-first, reusing the `Queued:`
parsing the wiki flow board already does. Reverses the queue-blind non-goal in
`.agents/specs/payload/spec.md`.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below but the first two is assumed)

- Q: which spec carries the queue-blind non-goal? → A: `payload` only. Its *Scope
  boundaries* say "the Go binary knowing about the queue. It is payload, not delivery."
  `cli` has no such sentence — it already specifies the wiki counting and listing
  `Queued:` proposals (`TestWikiHomeCarriesTheInFlightStrip`,
  `TestWikiFlowBoardListsTheQueue`). Observed by searching both specs, not assumed. **So
  the non-goal was already broken when the wiki landed**, and nobody reversed it: this
  change reverses it in `spec-payload.md`.
- Q: does the wiki apply find-work's branch filter today? → A: no. `readInFlight` in
  `cmd/libretto/wiki.go` counts every `proposal.md` holding a line that starts with the
  `Queued:` prefix and runs no git at all. find-work's *The queue — captured, never
  started* says a queued name matching an unmerged branch is in flight, never queued. **The
  two definitions disagree today**. Observed on this worktree's base before it was rebased:
  `add-loop-budget-and-run-log` still carried its `Queued:` line there although its work
  had merged as PR #74. After the rebase onto `main` that folder is gone, landed. A finding.
- Q: does the binary honour the branch filter? → A: yes, the skill's rule exactly — local
  branches, one leading `feat/` `fix/` `docs/` `chore/` `refactor/` removed, the change
  name matched whole (assumed). If wrong — the binary should list the raw `Queued:` scan —
  the filter comes out and the command reproduces the rework report find-work was patched
  against.
- Q: unmerged into what? find-work says "the base branch", and the binary does not know it.
  → A: into `HEAD`: `git branch --format=%(refname:short) --no-merged HEAD` (assumed). Run on
  the base branch it is the skill's command exactly; run elsewhere it can only drop more
  changes, never fewer, and only ones some branch carries by name. If wrong, a base is
  resolved from `refs/remotes/origin/HEAD` with `main` as fallback, in one function.
- Q: local branches only, or remote ones too? → A: local only, as find-work's own command
  reads (assumed). If wrong, both find-work and the binary add `-r` and strip the remote
  name — together, or they disagree again.
- Q: does this change align the wiki with the skill? → A: yes. The scan moves out of
  `readInFlight` into one function both the wiki and `queue` call, and the filter applies
  to both (assumed). If wrong, the wiki keeps its unfiltered count and the binary holds two
  definitions of "queued" — the thing the find-work section exists to prevent.
- Q: tie-breaker for one `Queued:` date? → A: the change name, ascending (assumed). The wiki
  already sorts that way, it needs no git, and the date in the file was chosen over git
  archaeology on purpose. find-work gains the same sentence so `/libretto-next`'s
  recommendation matches the binary's first line. If wrong (capture order was meant), one
  `git log --diff-filter=A` per proposal, changed in the skill and the binary together.
- Q: what does one line say? → A: the `Queued:` date, the change name, and the first
  sentence of its `## Reading` through the wiki's existing `firstSentence`; the name and
  date alone when there is no `## Reading` (assumed). That is find-work's "name and what its
  proposal says it is for". If wrong, the reading drops off the line.
- Q: what happens to a change the filter drops? → A: one closing line names each, with the
  branch that took it (assumed). An idea that vanishes from a list without a word reads as
  lost. If wrong, the line goes.
- Q: git cannot list branches (not a repository, no commit yet)? → A: list every
  `Queued:` proposal unfiltered, say so in one stderr line, exit 0 (assumed). find-work
  fails safe towards keeping an idea in the queue; this is the same direction.
- Q: where are the changes read from? → A: the repository root, from
  `git rev-parse --show-toplevel`, falling back to the cwd when git fails (assumed). `cli`
  records `metrics` reporting nothing from a subdirectory because it trusted the cwd, and
  `init` was specified against that. If wrong, the cwd, as `wiki` and `loop` read it.
- Q: arguments? → A: none accepted; any argument is refused (assumed). If wrong, a name
  filter is a later change.
- Q: is `queue` exempt from the missing-payload stop? → A: yes, like `loop`, `metrics`,
  `land` and `init` — it reads the project, never the payload (assumed).
- Q: does `/libretto-status` change? → A: no. It and `/libretto-next` keep delegating the
  scan to find-work and never call the binary, so they work where only the skills are
  installed (assumed). If wrong, status runs `libretto queue` where it is on PATH and keeps
  the skill's scan as the fallback — two paths to keep in agreement.
- Q: what changes in the skills? → A: find-work's queue section gains three sentences — the
  checklist clause (added after review-spec, below), the tie-breaker, and that `libretto queue` implements this section, the section winning any
  disagreement (assumed). No other skill moves.
- Q: a panel row, `--json`, a README row? → A: none (assumed). `--json` is the queued
  `add-json-output`'s; a panel row earns its place when missed; the README's command table
  already omits `loop`, `metrics` and `wiki`. The help line and `AGENTS.md`'s command
  listing carry it.
- Q: the fork — how is it built? → A: **A, one scan in `queue.go` that the wiki also calls**
  (assumed, the recommended option). B — `queue` filters on top of `readInFlight` and the
  wiki stays as it is — lost because it leaves two definitions inside one binary and keeps
  the wiki's wrong count. C — `queue` shells out to find-work's `rg` line — lost because
  the binary must not need `rg`, and the filter would still be written in Go.
- Q: the stop before code → A: carried into phase 6 by the invocation (assumed);
  `review-spec`'s findings, if any, are acted on and named in the report.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: is "no `.agents/changes/`" the same state as "nothing queued"? → A: no. With no
  directory the command prints `no .agents/changes/ here — nothing is queued`; with a
  directory holding no queued idea it prints `nothing is queued`. Both exit 0 (assumed).
  A missing directory usually means a project that never ran `init`, and saying so points
  at the fix. If wrong, the two collapse into one line again.
- Q: is "nothing is queued" measured before or after the branch filter? → A: after. When
  every queued change is dispatched, the command prints `nothing is queued` and then the
  closing line naming the dispatched ones (assumed). Before the filter it would print an
  empty list under a claim that something is queued. If wrong, the empty line is
  suppressed whenever anything was dropped.
- Q: which stream carries the closing dispatched line? → A: stdout — it is information,
  not an error. Only the unreadable-branches line goes to stderr (assumed). If wrong, the
  dispatched line moves to stderr and a pipe sees only the queue.
- Q: does a `Queued:` proposal holding a box under `## Checklist` count as queued? → A: no.
  A proposal is queued when it carries a `Queued:` line at column 0, **and** no unmerged
  local branch takes its name, **and** it has no box under its `## Checklist`, because a
  box means the work started. `queuedIdeas` applies all three, the amended `cli` wiki
  criterion keeps its checklist clause and adds the branch clause, and find-work's queue
  section gains the same sentence so the skill and the binary keep one definition
  (assumed). If wrong — a box should not disqualify — the clause comes out of all three
  places together.
- Q: is "the wiki lists what the command lists" an outcome without a test? → A: it gets its
  own criterion, `TestWikiQueueMatchesTheCommand`: one fixture, and the wiki's flow-board
  queue names the same changes, in the same order, as `queuedIdeas` (assumed). Box 2
  closes it. If wrong, the outcome is dropped rather than left unproven.
- Q: does `spec-payload.md`'s third criterion also promise that status and next delegate to
  find-work? → A: no. The existing payload criterion "the queue commands delegate the scan
  too" already owns that; the new one keeps only the absence of `libretto queue`
  (assumed). A duplicate drifts. If wrong, the clause returns and the two criteria are
  merged at landing.

## Outcomes

- `libretto queue` prints the project's queued ideas, one per line — date, name, what the
  proposal says it is for — oldest `Queued:` date first, ties by name.
- Queued means three things at once: a `Queued:` line at column 0, no unmerged local
  branch carrying the name, and no box under the proposal's `## Checklist`.
- A change whose name an unmerged local branch carries is not in that list; a closing line
  on stdout names it with the branch.
- A missing `.agents/changes/` prints `no .agents/changes/ here — nothing is queued`; a
  changes directory with no queued idea left after the filter prints `nothing is queued`,
  followed by the dispatched line when anything was dropped. A project where git cannot
  list branches gets the unfiltered queue and one stderr line. All of these exit 0.
- Nothing on disk changes when it runs.
- The wiki's in-flight strip count and its flow-board queue list exactly what
  `libretto queue` lists, in the same order.
- `libretto help` and `AGENTS.md` name the command.

## Scope boundaries

**In:**

- `cmd/libretto/queue.go`: the scan, the branch filter, the command
- the wiki reading the queue through that scan
- dispatch, the `needsPayload` exemption, the help line, one line in `AGENTS.md`
- find-work's three sentences and their wiring checks (`spec-payload.md`)
- reversing `payload`'s non-goal (`spec-payload.md`)

**Out:**

- **editing, deleting, reordering or picking up a queued idea.** `/libretto-next` picks
  up; files are edited by hand. The payload's non-goals on that stand.
- **`--json`.** `add-json-output` owns machine-readable output for read-only commands.
- **a name or date filter, or any argument.**
- **a panel row.**
- **remote branches, and resolving a base branch.** `HEAD` stands in for the base.
- **`/libretto-status` or `/libretto-next` calling the binary.**
- **listing changes in flight.** That is source 1 and `/libretto-status`; this command
  lists the queue and nothing else.
- **reading a tracker.**

## Constraints

- Go 1.26.5, stdlib only. No dependency is added, so `AGENTS.md`'s ask-first on
  dependencies is not reached.
- Read-only: no file written, no git mutation. Two git invocations —
  `rev-parse --show-toplevel` and `branch --format=%(refname:short) --no-merged HEAD`.
- Git runs through the `gitRunner` seam `metrics`, `land` and `init` use, so tests drive
  it with fixed output and build no repository.
- The `Queued:` definition stays find-work's. The binary implements it; it does not
  extend it.
- A line counts only when it starts with the `Queued:` prefix at column 0 — the wiki's
  existing rule — so a quoted mention inside an ask never queues a change.
- Every existing wiki criterion in `.agents/specs/cli/spec.md` holds, except the one this
  change amends (see Prior decisions).
- Tests point at `t.TempDir()` and a fake git; nothing touches this repository's changes.

## Prior decisions

- **The binary may read the queue; find-work still defines it.** 2026-10-22, attacca run.
  Reverses `payload`'s "the Go binary knowing about the queue" — which the wiki had
  already broken without saying so. The binary implementing a rule the skill owns is
  delivery; the binary deciding what counts as queued would not be. See `spec-payload.md`.
- **The branch filter is find-work's, applied by the binary**: local branches unmerged
  into `HEAD`, one leading `feat/` `fix/` `docs/` `chore/` `refactor/` removed, the name
  matched whole. Assumed 2026-10-22, attacca run. **Ceiling:** `HEAD` stands in for the
  base branch; run from a branch behind the base, a landed change can be dropped as
  dispatched rather than shown — which drops, never invents. The day that misleads, a base
  is resolved from `refs/remotes/origin/HEAD`.
- **Local branches only.** Assumed 2026-10-22, attacca run. A change dispatched from
  another machine and not fetched still shows as queued, exactly as in find-work. If
  wrong, both add remote branches together.
- **The wiki reads the same filtered scan.** Assumed 2026-10-22, attacca run. This amends
  `cli`'s criterion "Where a change's `proposal.md` carries a `Queued:` line and no box
  under `## Checklist`, the wiki shall count it in the queue": at landing it keeps its
  checklist clause and gains "and no unmerged local branch carries its name". If wrong,
  the wiki's unfiltered count stays and one binary holds two definitions.
- **Queued is three conditions, applied in one place.** A `Queued:` line at column 0, no
  unmerged local branch taking the name, and no box under `## Checklist` — a box means the
  work started. `queuedIdeas` applies all three for both the wiki and the command, and
  find-work's queue section states the same three. Assumed 2026-10-22, attacca run, on
  review-spec's finding. If wrong, the checklist clause leaves all three places together.
- **The empty states are told apart, and measured after the filter.** No
  `.agents/changes/` prints `no .agents/changes/ here — nothing is queued`; nothing left
  after the filter prints `nothing is queued`, then the dispatched line when anything was
  dropped. Exit 0 either way. Assumed 2026-10-22, attacca run.
- **The dispatched line is stdout; only the unreadable-branches line is stderr.** Assumed
  2026-10-22, attacca run — information, not an error.
- **Ties on one `Queued:` date go by change name, ascending.** Assumed 2026-10-22, attacca
  run. Already the wiki's order; find-work gains the sentence. If wrong (capture order),
  both move to one `git log` per proposal.
- **Unreadable branches list the queue unfiltered, with one stderr line, exit 0.** Assumed
  2026-10-22, attacca run — find-work's fail-safe direction: keep the idea.
- **Changes are read at the repository root, the cwd only when git fails.** Assumed
  2026-10-22, attacca run, from `metrics`' subdirectory bug. If wrong, the cwd.
- **A dropped change is named on a closing line with its branch.** Assumed 2026-10-22,
  attacca run.
- **No arguments; `queue` is exempt from the missing-payload stop.** Assumed 2026-10-22,
  attacca run.
- **The fork: A — one scan in `queue.go`, called by the wiki and the command.** Assumed
  2026-10-22, attacca run. B, filtering only in the command, lost on two definitions in
  one binary; C, shelling out to `rg`, lost on a new runtime requirement.

## Task breakdown

- extract the `Queued:` scan from `readInFlight` into one function returning name, date
  and reading, sorted by date then name, skipping any proposal with a box under its
  `## Checklist`
- list unmerged local branches through `gitRunner`; match names after removing one known
  prefix; split the scan into queued and dispatched
- the command: repository root, the filter, one line per idea, the closing line, the
  empty and unreadable states, refusing arguments, writing nothing
- dispatch, the `needsPayload` exemption, the help line, `AGENTS.md`
- the wiki: branches through a seam like its other git seams; strip count and flow board
  read the filtered scan
- find-work's three sentences; `check-payload` wiring for them and for status and next
  not calling the binary

## Verification criteria

- **When `libretto queue` runs, it shall print one line per queued change carrying its
  `Queued:` date, its name and the first sentence of its `## Reading`, ordered by date,
  oldest first.**
  Proof: cmd/libretto/queue_test.go TestQueueListsOldestFirst
- **Where two queued changes carry the same `Queued:` date, the command shall order them
  by change name, ascending.**
  Proof: cmd/libretto/queue_test.go TestQueueBreaksDateTiesByName
- **Where a queued proposal has no `## Reading`, its line shall carry the date and the name
  only.**
  Proof: cmd/libretto/queue_test.go TestQueueWithoutAReadingPrintsTheNameAlone
- **If a proposal mentions `Queued:` anywhere but at the start of a line, then the command
  shall not list that change.**
  Proof: cmd/libretto/queue_test.go TestQueueReadsOnlyALineStartingWithQueued
- **If a proposal carrying a `Queued:` line holds a box under its `## Checklist`, then the
  command shall not list that change.**
  Proof: cmd/libretto/queue_test.go TestQueueSkipsAProposalWithAChecklistBox
- **If an unmerged local branch's name, with one leading `feat/`, `fix/`, `docs/`,
  `chore/` or `refactor/` removed, equals a queued change's name, then the command shall
  leave that change out of the list and shall name it, with the branch, on a closing
  line written to stdout.**
  Proof: cmd/libretto/queue_test.go TestQueueDropsAChangeItsBranchAlreadyTook
- **If a branch matches a queued change's name only in part, or only after removing a
  prefix outside that list, then the command shall keep the change in the list.**
  Proof: cmd/libretto/queue_test.go TestQueueKeepsAChangeOnAPartialOrForeignMatch
- **The command shall ask git for branches with exactly
  `branch --format=%(refname:short) --no-merged HEAD`.**
  Proof: cmd/libretto/queue_test.go TestQueueAsksGitForUnmergedLocalBranches
- **If git cannot list branches, then the command shall list every queued change
  unfiltered, shall say so in one line on stderr, and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueUnreadableBranchesListsUnfiltered
- **When run from a subdirectory of a repository, the command shall read the
  `.agents/changes/` at the repository root.**
  Proof: cmd/libretto/queue_test.go TestQueueReadsTheRepositoryRoot
- **While no `.agents/changes/` exists at the repository root, the command shall print
  `no .agents/changes/ here — nothing is queued` and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueWithoutAChangesDirSaysSo
- **While `.agents/changes/` exists and no change in it is queued after the branch filter,
  the command shall print `nothing is queued` and shall exit 0.**
  Proof: cmd/libretto/queue_test.go TestQueueEmptyIsAStateNotAnError
- **If every queued change is dispatched by a branch, then the command shall print
  `nothing is queued` followed by the closing line naming each dispatched change with its
  branch, both on stdout.**
  Proof: cmd/libretto/queue_test.go TestQueueAllDispatchedSaysNothingQueuedThenNamesThem
- **The command shall leave every file under the project byte-identical.**
  Proof: cmd/libretto/queue_test.go TestQueueWritesNothing
- **If the command is given any argument, then it shall exit non-zero naming the argument,
  and shall print no queue.**
  Proof: cmd/libretto/queue_test.go TestQueueRefusesArguments
- **The missing-payload stop shall not apply to `queue`.**
  Proof: cmd/libretto/queue_test.go TestQueueNeedsNoPayload
- **The help text shall carry a `queue` line saying it lists queued ideas, read-only.**
  Proof: cmd/libretto/main_test.go TestHelpNamesQueue
- **If an unmerged local branch carries a queued change's name, then the wiki shall
  neither list that change in the flow board's queue nor count it in the in-flight strip.**
  Proof: cmd/libretto/wiki_test.go TestWikiQueueDropsAChangeItsBranchAlreadyTook
- **The wiki's flow-board queue shall list the same change names, in the same order, as
  `queuedIdeas` returns for the same changes directory and branch list.**
  Proof: cmd/libretto/wiki_test.go TestWikiQueueMatchesTheCommand

## Approach

### Summary

One new file, `cmd/libretto/queue.go`, holding the scan the wiki already performs —
lifted out of `readInFlight` — plus find-work's branch filter, and a short command that
prints the result. The wiki then calls the same scan, so the binary has one answer to
"what is queued" and it is the skill's. find-work gains the tie-breaker and a sentence
naming the binary as its implementation; `payload`'s non-goal is reversed.

### Technical context

- Go 1.26.5, stdlib: `os`, `path/filepath`, `sort`, `strings`, `fmt`, `io`.
- The `gitRunner` seam and `execGit` live in `cmd/libretto/metrics.go`; `land` and `init`
  take one the same way from `main.go`'s dispatch.
- The wiki reaches git through package-level function variables (`wikiGitDates`,
  `wikiGitTracked`) that its tests replace; the branch list joins them as one more.
- `firstSentence` in `wiki.go` already cuts a paragraph at its first `. `.
- `needsPayload` in `cmd/libretto/root.go` lists the exempt commands; `usage()` in
  `main.go` is the help text, and `main_test.go` has `TestHelpNamesLand` and
  `TestHelpNamesInit` as precedent.
- Blast radius: `cmd/libretto/queue.go` (new), `cmd/libretto/queue_test.go` (new),
  `cmd/libretto/wiki.go`, `cmd/libretto/wiki_test.go`, `cmd/libretto/main.go`,
  `cmd/libretto/main_test.go`, `cmd/libretto/root.go`, `AGENTS.md`,
  `skills/find-work/SKILL.md`, `scripts/check-payload`; at landing the deltas onto
  `.agents/specs/cli/spec.md` and `.agents/specs/payload/spec.md`.
- Gates: the six in `AGENTS.md`.

### The approach

1. **The scan.** `queuedIdeas(dir string) []queuedIdea` walks `<dir>/.agents/changes/*/`,
   reads each `proposal.md`, and keeps those holding a line with the `Queued:` prefix at
   column 0 — `readInFlight`'s loop, moved — and drops any whose `## Checklist` holds a
   box, counted with `readPlan`'s box definition so the loop and the scan agree on what a
   box is. `queuedIdea` is `{name, date, reading string}`;
   `reading` is `firstSentence` over the non-empty lines between `## Reading` and the next
   `## ` heading. Sorted by date, then name — the sort `readInFlight` already does. A
   missing directory returns nil.
2. **The filter.** `unmergedBranches(git gitRunner) ([]string, error)` runs
   `branch --format=%(refname:short) --no-merged HEAD` and splits non-empty lines.
   `dispatchedBy(name string, branches []string) string` returns the first branch whose
   name, after removing the first matching prefix of `feat/ fix/ docs/ chore/ refactor/`,
   equals `name`; else "". `splitDispatched(ideas, branches)` returns the kept ideas in
   order and the dropped ones paired with their branch.
3. **The command.** `queue(w, errw io.Writer, args []string, cwd string, git gitRunner)
   error`. Any argument: an error naming it. Root: `rev-parse --show-toplevel`, trimmed,
   else `cwd`. No `<root>/.agents/changes/`: `no .agents/changes/ here — nothing is
   queued`. Branches: on error, one stderr line and no filter. After the filter it prints
   `<date>  <name>  <reading>` per idea, or `nothing is queued` when none is left; then,
   when anything was dropped, one closing line on stdout naming each `name (branch)`.
   `main.go`
   dispatches `queue` with `execGit(projectDir)`; `root.go` adds it to the exemption list;
   `usage()` gains a line; `AGENTS.md`'s command block gains one.
4. **The wiki.** `readInFlight` replaces its inline loop with `queuedIdeas(projectDir)`
   and, from box 2, filters through a new seam, `wikiGitBranches`, defined as
   `unmergedBranches(execGit(dir))`. `wikiQueued` gives way to `queuedIdea`. Its error
   path is no filter, so every existing wiki test — temporary directories that are not
   repositories — keeps its current output.
5. **The payload.** find-work's queue section gains: a proposal with a box under its
   `## Checklist` is not queued, because the box means the work started; ties on one
   `Queued:` date go by change name; and `libretto queue` implements this section, which
   is the rule where the two disagree. `check-payload` gains a `check_wiring` for each of
   the three, and a `check_absent` of
   `libretto queue` in `commands/libretto-status.md` and `commands/libretto-next.md`.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — `queue` filters on top of `readInFlight`, the wiki untouched | two definitions of "queued" in one binary, and the wiki keeps counting dispatched changes |
| C — `queue` runs find-work's `rg` line | puts `rg` on the binary's runtime path, and the branch filter would still be Go |
| unmerged into a resolved base branch | one more lookup with its own failure modes, for a case `HEAD` gets right whenever the command runs on the base |
| capture-commit order for ties | a `git log` per proposal, against a design that put the date in the file to avoid exactly that |
| the cwd, as the wiki reads it | the bug `metrics` recorded: a subdirectory run reports an empty queue |
| `/libretto-status` calling `libretto queue` where it is on PATH | a second path to "what is queued", and status would behave differently on a skills-only machine |

### Risks

| Risk | What catches it |
|---|---|
| the extraction changes the wiki's queue for an unfiltered project | `TestWikiFlowBoardListsTheQueue`, `TestWikiHomeCarriesTheInFlightStrip` and `TestWikiInFlightReadsTheChangeDocument` run unchanged |
| a wiki test resolves real branches | the seam is a variable; `TestWikiQueueDropsAChangeItsBranchAlreadyTook` sets it, and the others run in temporary directories that are not repositories |
| find-work's prose and the Go prefix list drift apart | `check_wiring` pins the find-work sentences; the prefix list itself is not cross-checked. Accepted — the skill names it in one sentence and the code in one slice |
| `HEAD` on a branch behind the base drops a landed change as dispatched | named in Prior decisions as the ceiling; the closing line shows what was dropped and why |
| the reading's first sentence runs long | not truncated; one line per idea, whatever its width. Accepted |

### Validation and rollback

New tests: the nineteen named under Verification criteria, plus the five `check-payload`
lines. Two must be forced red before they are believed:

- `TestQueueDropsAChangeItsBranchAlreadyTook` passes if the fixture's change is simply
  missing from the scan. Assert it is listed with the fake git returning no branches
  first, then dropped with the branch present — and break `dispatchedBy` on purpose to
  watch the second half fail.
- `TestQueueWritesNothing` passes if the command printed nothing at all. Assert the queue
  was printed before comparing the tree.

Rollback: one revert. Nothing migrates, nothing is written.

### Complexity deliberately kept

- The closing line for dropped changes: one more output path, and the only thing that
  stops a dispatched idea from reading as lost.
- The repository-root lookup: one more git call than the wiki makes, for the bug `metrics`
  already paid for.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [x] **`libretto queue` lists the queue, filtered as find-work filters it.** `queue.go`
      with `queuedIdeas`, `unmergedBranches`, `dispatchedBy`, `splitDispatched` and the
      command, the checklist clause inside `queuedIdeas`; `readInFlight` calls
      `queuedIdeas` with no branch filter, so the wiki's output changes only for a queued
      proposal holding a checklist box; dispatch, the `needsPayload` exemption, the help
      line, the `AGENTS.md` line. Closes: TestQueueListsOldestFirst, TestQueueBreaksDateTiesByName,
      TestQueueWithoutAReadingPrintsTheNameAlone, TestQueueReadsOnlyALineStartingWithQueued,
      TestQueueSkipsAProposalWithAChecklistBox, TestQueueDropsAChangeItsBranchAlreadyTook,
      TestQueueKeepsAChangeOnAPartialOrForeignMatch,
      TestQueueAsksGitForUnmergedLocalBranches,
      TestQueueUnreadableBranchesListsUnfiltered, TestQueueReadsTheRepositoryRoot,
      TestQueueWithoutAChangesDirSaysSo, TestQueueEmptyIsAStateNotAnError,
      TestQueueAllDispatchedSaysNothingQueuedThenNamesThem, TestQueueWritesNothing,
      TestQueueRefusesArguments, TestQueueNeedsNoPayload, TestHelpNamesQueue. Waits on: nothing.
- [x] **The wiki reads the same filtered queue.** The `wikiGitBranches` seam;
      `readInFlight` filters through `splitDispatched`; every existing wiki test unchanged.
      Closes: TestWikiQueueDropsAChangeItsBranchAlreadyTook, TestWikiQueueMatchesTheCommand.
      Waits on: box 1.
- [x] **find-work states the checklist clause, the tie-breaker and the binary.** The three
      sentences in `skills/find-work/SKILL.md`; `check-payload` gains their three
      `check_wiring` lines and the two `check_absent` lines for `/libretto-status` and
      `/libretto-next`. Closes: the four criteria in `spec-payload.md`, all on
      scripts/check-payload.
      Waits on: box 1, so the sentence never names a command that does not exist.
      Independent of box 2.
