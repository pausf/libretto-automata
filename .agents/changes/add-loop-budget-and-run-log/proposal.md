# add-loop-budget-and-run-log

Tracker: none
Targets: cli
Durable decisions: the ten in Prior decisions below

## The ask, verbatim

> Mejoras a `loop`: presupuesto de coste por ronda (alimentando `metrics`), un log
> persistido por ronda (hoy imprime y se pierde). Los rieles de seguridad (spend limits,
> checkpoint por iteración) son exactamente lo que la comunidad de Ralph-loops añadió en
> 2025-26.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Give `libretto loop` a per-round cost budget and a persisted per-round run log, so a run
is inspectable after the fact and `metrics` can consume it.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below but the first is assumed)

- Q: is there a real cost signal, or must the budget fall back to rounds, wall time or
  tokens? → A: there is one. `claude --help` (2.1.307) lists `--max-budget-usd <amount>`
  ("Maximum dollar amount to spend on API calls (only works with --print)") and
  `--output-format json`; one observed `claude -p … --output-format json` run returned a
  single object carrying `total_cost_usd`, `usage`, `session_id`, `subtype`, `is_error`,
  `num_turns`, `duration_ms` and `result`. The budget is dollars, enforced by the host.
  Observed, not assumed — the entries below are the assumptions built on it.
- Q: what does "per round" mean? → A: per session. One round is one `claude -p` session,
  so the host's per-session `--max-budget-usd` is exactly a per-round budget (assumed). If
  wrong — the ask meant a ceiling across the whole loop — the loop sums `cost_usd` from
  its own records and stops at a total: one more flag and one more stop reason.
- Q: what is the flag called? → A: `--max-budget-usd`, the host's own name, passed through
  unchanged (assumed). A second name for the same number is a translation somebody has to
  learn. If wrong, it is a rename in `parseLoopArgs` and the help line.
- Q: a loop-wide total budget as well? → A: no. The header prints the worst case,
  `--max` × budget, which the user reads before anything is spent (assumed). If wrong, see
  the "per round" entry: summing the log's own records is the mechanism.
- Q: when a session stops on its budget, does the loop stop? → A: no special case. The
  host's result subtype is recorded verbatim; a round cut short usually closes no box, and
  two such rounds already stop the loop through stuck detection (assumed). If wrong — the
  first budget hit should end the run — it is one stop reason keyed on the recorded
  subtype, whose spelling has to be observed first rather than guessed.
- Q: where does the run log live? → A: `<git-common-dir>/libretto/loop.jsonl` (assumed).
  Inside git's own directory nothing is ever tracked or listed by `git status`, no
  `.gitignore` is edited in the user's project, and the common dir is shared by every
  worktree — so a loop run in a worktree survives the worktree's removal and is visible to
  `metrics` run from the main checkout. If wrong — the log should travel with the
  repository — it moves under `.agents/`, at the cost of a tracked file or a `.gitignore`
  the tool writes into a project it does not own.
- Q: one file, or one per change? → A: one JSON-lines file, one record per round, each
  carrying its change name, written with stdlib `encoding/json` (assumed). One file is one
  read for `metrics`. If wrong, the path gains the change name and `metrics` globs.
- Q: if the log cannot be located (not a git repository) or written? → A: one warning,
  and the loop keeps driving the checklist (assumed). The checklist is the state the loop
  decides on; the log is a record of it, and losing the record is no reason to stop work.
  If wrong — the user would rather not spend without a record — it becomes a refusal
  before the first session.
- Q: where does the cost number come from? → A: the host's `total_cost_usd`, recorded as
  reported and omitted when absent, never zero; the binary holds no price table (assumed,
  and forced by the existing *no pricing* decision in `cli`). The one observed run carried
  `"costBasis": "list"` under `modelUsage`, so the figure is the host's list-price estimate
  and is labelled as the host's figure. If wrong — the user wants billed cost — no local
  source has it, and the label already says so.
- Q: how does `metrics` show it? → A: a separate block after the token block, one row per
  change, never a new column in the change table (assumed, and forced by the existing *no
  new column* decision in `cli`). If wrong, the block's rows move; nothing upstream does.
- Q: is the log ever pruned or rotated? → A: no. A record is a few hundred bytes; ten
  thousand rounds is a few megabytes (assumed). If wrong, rotation is a later change with
  its own ceiling.
- Q: what about `codex` and `opencode`, which the queued `support-codex-opencode-in-loop`
  will drive? → A: out of scope. `--max-budget-usd` and the JSON result are `claude`'s;
  that change decides what a budget means for a tool without them (assumed). If wrong,
  this change grows a runner abstraction nothing needs yet.
- Q: the fork — how is it built? → A: **A, the host's flag plus a JSONL log in the git
  common dir** (assumed, the recommended option). B — a self-ignoring `.agents/loop/`
  directory carrying its own `.gitignore` of `*` — lost because it is per worktree: the
  log dies with the worktree the loop ran in, and `metrics` from the main checkout never
  sees it. C — no new output; `metrics` joins the host transcripts by `session_id` and
  reports tokens — lost because it produces no budget at all, and dollars are not in the
  transcripts. If wrong, B is a path change and C is a `metrics`-only change.

## Outcomes

- `libretto loop <change> --max-budget-usd 2.50` starts every session with
  `--max-budget-usd 2.50`, and its header states the per-round budget and the worst case
  for the whole run, `--max` × budget.
- Every round that runs leaves one JSON line in `<git-common-dir>/libretto/loop.jsonl`:
  when, which change, which round, boxes closed that round, boxes still open, the host's
  reported cost when it reported one, duration, session id and result subtype, the budget
  that applied, a session error when there was one, and — on the last round — why the loop
  stopped.
- `git status` in the project shows nothing new after a loop.
- A loop that cannot write its log says so once and carries on.
- What the session printed still reaches the terminal: the `result` text of the JSON, or
  the raw output when it is not JSON.
- `libretto metrics` prints a loop block: per change, rounds, boxes closed by the loop,
  the summed reported cost, and how many rounds had no reported cost. Absent log, one
  line; filtered, one change.

## Scope boundaries

**In:**

- `--max-budget-usd <amount>` on `loop`, validated and passed through to `claude`
- `--output-format json` on every session, parsed for cost, duration, session id, subtype
- the run log: location, record shape, append-only writes, the warn-and-continue path
- the header's worst-case line
- the `metrics` loop block, its absent state and its malformed-line count
- the `loop` help line in `main.go`

**Out:**

- **a loop-wide spending cap.** The worst case is printed instead. Returns if the user asks
  for a ceiling across rounds; the log already holds what it would sum.
- **a budget measured in rounds, wall time or tokens.** `--max` is already the round cap,
  and the host enforces dollars; a second unit is a second budget.
- **a price table, or any cost the binary computes.** Prices move under the repository.
- **stopping the loop on the first budget hit.** Stuck detection already bounds it.
- **reading the run log anywhere but `metrics`.** No `status` line, no panel row, no
  `--json`. `add-json-output` can carry the block if it wants it.
- **budget and cost for `codex` and `opencode`.** `support-codex-opencode-in-loop` owns them.
- **rotating, pruning or migrating the log.**
- **checkpointing.** The ask's "checkpoint por iteración" already exists: each session
  commits its box, and the checklist is the state. Nothing new is written for it.
- **writing the log into the change folder or anywhere under `.agents/`.**

## Constraints

- Go 1.26.5, stdlib only — `encoding/json`, `os`, `os/exec`, `strconv`. No dependency.
- `loop` stays outside `needsPayload` and must keep working with only the binary installed.
- `--dry-run` keeps needing no `claude` on PATH and starts no session.
- `runLoop` stays drivable without a real `claude`: the runner seam carries the result.
- Tests point every log path at `t.TempDir()`; nothing writes into this repository's
  `.git` or a real `~/.claude`.
- `metrics` stays read-only; it opens the log for reading and never creates it.
- Every existing `loop` and `metrics` criterion in `.agents/specs/cli/spec.md` holds
  unchanged — stuck detection, the cap, a session error not failing the loop, the missing
  plan reported before the missing dependency, the token block's three states.

## Prior decisions

- **The budget is the host's `--max-budget-usd`, per session, which is per round.**
  Assumed 2026-10-22, attacca run. The host enforces it; the loop passes it through and
  records what was reported, so an overshoot is visible in the log rather than hidden.
  If wrong (a loop-wide cap was meant): sum the log's `cost_usd` and stop at a total.
- **The flag carries the host's name, `--max-budget-usd`.** Assumed 2026-10-22, attacca
  run. If wrong: a rename in one parser and one help line.
- **No loop-wide cap; the header prints `--max` × budget.** Assumed 2026-10-22, attacca
  run. **Ceiling:** adequate while somebody reads the header; the day a run is left
  unwatched with a large `--max`, a total cap is the replacement.
- **A budget-cut session gets no special case.** Assumed 2026-10-22, attacca run. Its
  subtype is recorded verbatim and never interpreted, because its exact spelling was not
  observed. If wrong: one stop reason keyed on the recorded subtype, once observed.
- **The run log lives at `<git-common-dir>/libretto/loop.jsonl`**, located with
  `git rev-parse --path-format=absolute --git-common-dir`. Assumed 2026-10-22, attacca
  run. Never tracked, never in `git status`, no `.gitignore` written, shared by every
  worktree. **Ceiling:** it is per machine — a fresh clone starts with no history, like the
  host's transcripts. If that stops being acceptable, it moves under `.agents/`.
- **The run log is the one instrumented artifact, and its writer is the binary.** This
  amends the `metrics` rule *nothing is instrumented* without reversing its reason: that
  rule rejected artifacts eight skills would have to remember to write. The loop is one
  writer that runs every round by construction and cannot forget; nothing lands in a
  change folder, and no skill learns anything. 2026-10-22, attacca run.
- **Cost is the host's `total_cost_usd`, recorded verbatim, absent when unreported, never
  zero.** Assumed 2026-10-22, attacca run. It extends the existing *no pricing* decision
  rather than relitigating it — the binary still prices nothing. The observed run labels it
  `"costBasis": "list"`, so it is shown as the host's estimate.
- **`metrics` shows the log as its own block, never as a table column.** Assumed
  2026-10-22, attacca run; the existing *no new column* decision already settles the table.
- **A log that cannot be written warns once and never stops the loop.** Assumed
  2026-10-22, attacca run. If wrong: a refusal before the first session.
- **The fork: A — host flag plus JSONL in the git common dir.** Assumed 2026-10-22,
  attacca run. B, a self-ignoring `.agents/loop/`, lost on worktrees; C, transcripts joined
  by session id, lost because it gives no budget and no dollars.

## Task breakdown

- ask every session for its JSON result; parse it into cost, duration, session id,
  subtype and text; print the text
- locate the log through git's common dir; append one record per round; warn and continue
  on any failure; write nothing under `--dry-run`
- mark the last record with why the loop stopped
- parse and validate `--max-budget-usd`; pass it through; state the per-round budget and
  the worst case in the header; record the budget
- read the log in `metrics`; one row per change; absent, filtered and malformed states
- the `loop` help line names the new flag

## Verification criteria

- **When the loop starts a session, the `claude` argument list it builds shall carry
  `-p <prompt>` and `--output-format json`.**
  Proof: cmd/libretto/loop_test.go TestEverySessionIsAskedForItsJSONResult
- **Where `--max-budget-usd <amount>` is given, the `claude` argument list shall carry
  `--max-budget-usd <amount>` with the amount exactly as typed; where it is not given, the
  list shall carry no `--max-budget-usd`.**
  Proof: cmd/libretto/loop_test.go TestClaudeArgsCarryTheBudgetOnlyWhenGiven
- **If `--max-budget-usd` is missing its value, or its value is not a positive finite
  decimal, then `parseLoopArgs` shall return an error naming `--max-budget-usd`.**
  Proof: cmd/libretto/loop_test.go TestParseLoopArgsRefusesABadBudget
- **While a budget is set, the loop header shall print the per-round budget and the
  product of `--max` and that budget as the worst case for the run.**
  Proof: cmd/libretto/loop_test.go TestTheHeaderStatesTheWorstCaseSpend
- **When the session's output is the host's JSON result object, the session parser shall
  return its `total_cost_usd`, `duration_ms`, `session_id`, `subtype` and `result`; if the
  output is not JSON or carries no `total_cost_usd`, then the parser shall report the cost
  as unknown, never as zero.**
  Proof: cmd/libretto/loop_test.go TestParseSessionResultReadsTheHostsCost
- **When a session ends, the loop shall print the `result` text of its JSON to stdout; if
  the output does not parse as JSON, then the loop shall print it unchanged.**
  Proof: cmd/libretto/loop_test.go TestTheSessionsTextStillReachesTheTerminal
- **When a round ends, the loop shall append exactly one JSON line to the run log carrying
  `change`, `round`, `closed` (boxes closed that round), `open` (boxes open after it) and
  the session's `subtype`.**
  Proof: cmd/libretto/loop_test.go TestEveryRoundAppendsOneRecord
- **Where a budget is set, each round's record shall carry it as `budget_usd`; where none
  is set, the record shall omit `budget_usd`.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheBudgetThatApplied
- **If the session reported no cost, then the round's record shall omit `cost_usd`
  entirely rather than write zero.**
  Proof: cmd/libretto/loop_test.go TestAnUnpricedRoundOmitsCostRatherThanWritingZero
- **When the loop stops because every box closed, two rounds closed nothing, or the cap
  was reached, the last record it wrote shall carry `stopped` set to `done`, `stuck` or
  `cap` respectively, and no earlier record shall carry `stopped`.**
  Proof: cmd/libretto/loop_test.go TestTheLastRecordSaysWhyTheLoopStopped
- **The loop shall write the run log at `<git-common-dir>/libretto/loop.jsonl`, and
  `git status --porcelain` in the project shall print nothing after a loop over a clean
  tree.**
  Proof: cmd/libretto/loop_test.go TestTheRunLogLeavesGitStatusClean
- **If the run log cannot be located or written, then the loop shall print one warning
  naming the log and shall run every round it would have run with a writable log.**
  Proof: cmd/libretto/loop_test.go TestAnUnwritableLogDoesNotStopTheLoop
- **While `--dry-run` is set, the loop shall write no record and create no log file.**
  Proof: cmd/libretto/loop_test.go TestDryRunWritesNoRecord
- **Where a run log exists, `metrics` shall print a loop block with one row per change
  holding its rounds, the boxes its rounds closed, the sum of its `cost_usd` values, and
  the count of its rounds with no `cost_usd`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsReportsLoopRoundsPerChange
- **If a change's rounds all lack `cost_usd`, then its cost cell shall print a dash, never
  `0`.**
  Proof: cmd/libretto/metrics_test.go TestAnUnpricedChangeReportsADashNotAZero
- **If no run log exists, then `metrics` shall print one line saying no loop runs are
  recorded in this checkout, and shall print every other block as it does today.**
  Proof: cmd/libretto/metrics_test.go TestNoRunLogIsAStateNotAnError
- **If a run log line does not parse as a record, then `metrics` shall skip it, keep every
  other row, and print the number of skipped lines.**
  Proof: cmd/libretto/metrics_test.go TestAMalformedRunLogLineIsSkippedAndCounted
- **Where a change is named, the loop block shall print only that change's row.**
  Proof: cmd/libretto/metrics_test.go TestTheLoopBlockFollowsTheChangeFilter
- **The loop block shall label its cost column as the host's reported estimate in
  dollars.**
  Proof: cmd/libretto/metrics_test.go TestTheLoopBlockLabelsCostAsTheHostsFigure

## Approach

### Summary

Two small seams in `loop.go` and one reader in `metrics.go`. The loop stops treating a
session as an opaque process: it asks `claude` for its JSON result, passes the user's
budget straight through as the host's own flag, prints the result text as before, and
appends one JSON line per round to a file inside git's common directory. `metrics` reads
that file and prints a block. No price table, no new dependency, no file in the project
tree.

### Technical context

- Go 1.26.5, stdlib. `claude` 2.1.307 observed: `--max-budget-usd` (print mode only) and
  `--output-format json` exist; the JSON result object's keys were read off one real run.
- Blast radius: `cmd/libretto/loop.go`, `cmd/libretto/loop_test.go`,
  `cmd/libretto/metrics.go`, `cmd/libretto/metrics_test.go`, `cmd/libretto/main.go` (the
  log path into `loop` and `metrics`, the help line). Five files, plus the delta onto
  `.agents/specs/cli/spec.md` at landing.
- `loopRunner` changes from `func(prompt string) error` to returning a parsed session
  alongside the error; every existing `loop_test.go` fake adapts to it.
- The log path is resolved outside the logic and passed in, the way `transcriptProjects()`
  is for the token block: `main.go` resolves it through `execGit(projectDir)`, and the
  tests hand in a `t.TempDir()` path. `fakeGit` in `metrics_test.go` needs no new arm.
- `TestTheRunLogLeavesGitStatusClean` builds a real temporary repository, so `make
  test-short` skips it, per the suite's convention for real-git tests.
- Gates: the six in `AGENTS.md`.

### The approach

1. **Arguments.** `parseLoopArgs` gains `--max-budget-usd <amount>`, checked with
   `strconv.ParseFloat`, refused unless positive and finite, and kept as the string the
   user typed so it reaches `claude` unchanged. `claudeArgs(prompt, budget string)
   []string` returns `-p <prompt> --output-format json`, plus the budget pair when set;
   `execClaude` uses it.
2. **The session result.** `execClaude` captures stdout into a buffer — stderr still goes
   straight to the terminal — and `parseSession(out []byte) session` decodes the one object
   into `{cost *float64, durationMS int64, sessionID, subtype, text string}`. A pointer,
   because unknown and zero are different facts. Output that is not JSON yields a session
   whose `text` is the raw output and whose cost is nil. The text is printed to stdout when
   the session ends, which is when `-p` printed it before.
3. **The header.** With a budget, one more line: the per-round budget, and `--max` ×
   budget as the run's worst case.
4. **The log.** `loopLogPath(git gitRunner) string` runs
   `rev-parse --path-format=absolute --git-common-dir` and returns
   `<dir>/libretto/loop.jsonl`, or "" on failure. `runLoop` receives the path in
   `loopOpts`. After each round it builds a record — `time`, `change`, `round`, `closed`,
   `open`, `cost_usd` (omitempty, pointer), `duration_ms`, `session_id`, `subtype`,
   `budget_usd` (omitempty), `error` (omitempty), `stopped` (omitempty) — and appends it
   with one `Write` on a file opened `O_APPEND|O_CREATE|O_WRONLY`, after `MkdirAll` of the
   `libretto/` directory. The stop reason is known only once the round's boxes are read,
   so the record is written after the stuck and cap decisions, not before; a round that
   closes the last box writes `stopped: done`. Any failure prints one warning, the first
   only, and the loop carries on. `--dry-run` returns before any record exists.
5. **`metrics`.** `metrics(...)` gains the log path as a parameter. `loopBlock` reads the
   file with `bufio.Reader` — the transcript reader's rule, no line-length cap — decodes
   each line, aggregates per change, applies the table's `only` filter, and prints rounds,
   closed, cost sum and unpriced count. Absent file, one line. Malformed lines counted and
   named. The column header says the cost is the host's estimate in dollars.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — a self-ignoring `.agents/loop/` holding a `.gitignore` of `*` | per worktree: the log dies with the worktree the loop ran in, and `metrics` from the main checkout never sees it |
| C — no new output; `metrics` joins host transcripts by `session_id` | gives no budget at all, and transcripts carry tokens, never dollars |
| a log appended inside the change folder | the session commits the folder, so the loop would race it; and it is deleted at landing — printed and lost again |
| `~/.claude/projects/<encoded>/` beside the transcripts | writes into the host's directory, which the token block's spec promises is never written |
| a loop-wide budget summed by the binary | answers a question the ask did not ask; the printed worst case covers it until somebody does |
| tokens as the budget unit | the host enforces dollars natively; a token cap would be a second enforcement the binary implements and the host does not know about |

### Risks

| Risk | What catches it |
|---|---|
| `--output-format json` hides output the user saw before | the loop prints `result` after each session; `TestTheSessionsTextStillReachesTheTerminal`. **Unobserved:** the text mode's exact output was not run this session, so equivalence is assumed |
| a host version without `--output-format json` or `total_cost_usd` | non-JSON output is printed raw and recorded with no cost; `TestParseSessionResultReadsTheHostsCost` covers both |
| the host overshoots its own budget | not caught, and not ours to enforce; the record carries the reported cost, so the overshoot is visible in `metrics`. Accepted |
| two loops in two worktrees append to one file at once | one `Write` per record on an `O_APPEND` descriptor, a record being a few hundred bytes. Accepted, not tested |
| the budget-hit subtype is never interpreted | deliberate; stuck detection bounds the spend at two barren rounds |
| a test writes into this repository's `.git` | the path is a parameter; only `TestTheRunLogLeavesGitStatusClean` resolves one, inside its own temporary repository |
| an old `git` without `--path-format` | `loopLogPath` returns "" and the loop warns and continues. Accepted |

### Validation and rollback

New tests: the nineteen named under Verification criteria. Two must be forced red before
they are believed, because they pass on the obvious wrong build:

- `TestAnUnpricedRoundOmitsCostRatherThanWritingZero` — a plain `float64` field writes
  `0`, and a test that only checks the record exists passes. Assert the key is absent,
  then break the pointer on purpose and watch it fail.
- `TestTheRunLogLeavesGitStatusClean` — passes trivially if no log was written at all.
  Assert the file exists and holds a record before asserting the clean status.

Every existing `loop_test.go` test must stay green across the runner signature change;
`TestASessionErrorDoesNotAbortTheLoop` and `TestDryRunNeedsNoClaudeOnPath` are the ones
most likely to break.

Rollback: one revert. Nothing migrates; a log left in `.git/libretto/` is inert.

### Complexity deliberately kept

- `cost *float64`, not `float64`: absent and zero are different facts — the rule
  `metrics` already keeps for corrections and tokens.
- The log path resolved outside the logic and passed in: the `transcriptProjects()` seam,
  and what keeps the suite off this repository's `.git`.
- The record written after the stop decision rather than before it: one more branch, and
  the only way the last record can say why the loop stopped.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [ ] **The loop records every round.** `claudeArgs` with `--output-format json`,
      `parseSession`, the result text printed, `loopLogPath`, one record per round,
      `stopped` on the last, warn and continue, nothing under `--dry-run`, every existing
      loop test adapted to the new runner. Closes: TestEverySessionIsAskedForItsJSONResult,
      TestParseSessionResultReadsTheHostsCost, TestTheSessionsTextStillReachesTheTerminal,
      TestEveryRoundAppendsOneRecord, TestAnUnpricedRoundOmitsCostRatherThanWritingZero,
      TestTheLastRecordSaysWhyTheLoopStopped, TestTheRunLogLeavesGitStatusClean,
      TestAnUnwritableLogDoesNotStopTheLoop, TestDryRunWritesNoRecord. Waits on: nothing.
- [ ] **The loop takes a per-round budget.** `--max-budget-usd` parsed and passed through,
      the worst-case header line, `budget_usd` in the record, the help line in `main.go`.
      Closes: TestClaudeArgsCarryTheBudgetOnlyWhenGiven, TestParseLoopArgsRefusesABadBudget,
      TestTheHeaderStatesTheWorstCaseSpend, TestTheRecordCarriesTheBudgetThatApplied.
      Waits on: box 1 (`claudeArgs` and the record).
- [ ] **`metrics` reads the run log.** The path passed from `main.go`, the loop block, its
      absent, filtered and malformed states, the cost label. Closes:
      TestMetricsReportsLoopRoundsPerChange, TestAnUnpricedChangeReportsADashNotAZero,
      TestNoRunLogIsAStateNotAnError, TestAMalformedRunLogLineIsSkippedAndCounted,
      TestTheLoopBlockFollowsTheChangeFilter, TestTheLoopBlockLabelsCostAsTheHostsFigure.
      Waits on: box 1 (the record shape). Independent of box 2.
