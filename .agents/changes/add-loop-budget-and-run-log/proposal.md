# add-loop-budget-and-run-log

Tracker: none
Targets: cli
Durable decisions: the eleven in Prior decisions below

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

### Session 2026-10-22 · review-spec findings (attacca — the run acted on every one; each answer assumed)

- Q: Constraints said every existing `metrics` criterion holds, while Prior decisions
  amended *nothing is instrumented* — which wins? → A: Constraints exempts exactly that
  bullet and the heading above it, and the delta states the replacement text for both:
  the change folders and the flow stay uninstrumented; the only recorded data is the
  loop's own run log, written by the binary and nothing else. Its proof stays
  `TestChangeNamesSeesLandedChangesNotJustOpenOnes` plus the run-log tests, and a box
  rewrites both at landing (assumed). If wrong — the bullet should stand untouched — the
  run log contradicts a capability promise and the change cannot land as designed.
- Q: what does `--max-budget-usd` accept? → A: plain decimal only, `^[0-9]+(\.[0-9]+)?$`,
  greater than zero; `1e2`, `0x…`, `+3`, `inf` and everything else is refused; the value
  reaches `claude` exactly as typed (assumed). If wrong — the host accepts a wider syntax
  the user wants — the regex widens; the pass-through does not change.
- Q: how is the header's worst case printed? → A: `%.2f`, cents — `--max 3
  --max-budget-usd 0.1` prints `0.30` (assumed). If wrong, a format string.
- Q: how is the `metrics` cost sum printed? → A: `%.2f` — no exponent and no float noise,
  consistent with the token block's *plain digits, no exponent* decision (assumed). If
  wrong, a format string; sub-cent rounds show only once they add up to a cent.
- Q: which record fields does a test pin? → A: every one — `time`, `change`, `round`,
  `closed`, `open`, `cost_usd`, `duration_ms`, `session_id`, `subtype`, `is_error`,
  `budget_usd`, `error`, `stopped` — so a build missing any field fails a test (assumed).
  If wrong, a field nobody reads gets dropped along with its criterion.
- Q: one error field or two? → A: two. `error` is the exec error's message on a non-zero
  exit, omitted otherwise; `is_error` is the host JSON's `is_error`, recorded as reported
  when present and omitted when absent. Neither is interpreted (assumed). If wrong, one of
  them goes; the loop's decisions never read either.
- Q: what does the no-log line say? → A: "no loop runs are recorded in this repository" —
  the log is shared across worktrees, so "checkout" would be wrong (assumed). If wrong,
  one string.
- Q: the checklist disappears mid-loop — is that round recorded? → A: yes. It gets a
  record with `closed` and `open` omitted, `stopped: "error"` and the error text in
  `error`; the loop still returns the error it returns today (assumed). If wrong, the
  last round of a failed run is the one missing from the log.
- Q: where in `loop()` is the log path resolved? → A: after the checklist refusal and
  after the PATH check, so *the missing plan is reported before the missing dependency*
  holds unchanged; `--dry-run` never resolves it and prints no warning (assumed). If
  wrong, a refusal could be preceded by a log warning about a run that never started.
- Q: what do the two log warnings say? → A: unlocatable: `run log not recorded: git common
  dir not found` — there is no path to name; unwritable: the warning names the path
  (assumed). If wrong, two strings.
- Q: is `closed` signed? → A: yes, `after.done − before.done`, so a reopened box reads
  negative; `metrics` sums it as a net figure and the column header says net (assumed). If
  wrong — closes and reopens should be counted apart — the record gains a second field,
  the way the change table already splits closed from reopen.
- Q: does box 2's help line have a criterion? → A: yes — `TestHelpNamesTheLoopBudget`
  asserts `--max-budget-usd` in the loop help (assumed). If wrong, the flag ships
  undiscoverable.
- Q: `metrics <change>` for a change with no records? → A: one line, `no loop rounds
  recorded for <change>` (assumed). If wrong, an empty block that reads as broken.
- Q: in what order are the loop block's rows? → A: sorted by change name (assumed). If
  wrong, a sort key.
- Q: does a change filter narrow the malformed-line count? → A: no — the count always
  covers the whole file, because a malformed line has no change to filter by (assumed). If
  wrong, malformed lines vanish under every filter, which is the silent skip the count
  exists to prevent.

## Outcomes

- `libretto loop <change> --max-budget-usd 2.50` starts every session with
  `--max-budget-usd 2.50`, and its header states the per-round budget and the worst case
  for the whole run, `--max` × budget, in cents.
- Every round that runs leaves one JSON line in `<git-common-dir>/libretto/loop.jsonl`:
  `time`, `change`, `round`, `closed` (signed), `open`, `cost_usd` when the host reported
  one, `duration_ms`, `session_id`, `subtype`, `is_error` when reported, `budget_usd` when
  set, `error` on a non-zero exit, and — on the last round only — `stopped`, one of
  `done`, `stuck`, `cap`, `error`.
- A round whose checklist vanished still leaves a record, marked `stopped: "error"`.
- `git status` in the project shows nothing new after a loop.
- A loop that cannot locate or write its log says so once and carries on; `--dry-run`
  never touches the log and never warns about it.
- What the session printed still reaches the terminal: the `result` text of the JSON, or
  the raw output when it is not JSON.
- `libretto metrics` prints a loop block, sorted by change: rounds, net boxes closed by
  the loop, the summed reported cost in dollars to the cent, and how many rounds had no
  reported cost. No log: one line. A named change with no rounds: one line. Malformed
  lines: counted over the whole file.

## Scope boundaries

**In:**

- `--max-budget-usd <amount>` on `loop`, validated as a plain decimal and passed through
  to `claude` as typed
- `--output-format json` on every session, parsed for cost, duration, session id,
  subtype, `is_error` and text
- the run log: location, every record field, append-only writes, the two warnings, the
  record for a vanished checklist
- the header's worst-case line
- the `metrics` loop block: rows, sort, absent, filtered-empty and malformed states
- the `loop` help line in `main.go`
- rewriting the `metrics` heading and its *nothing is instrumented* bullet in
  `.agents/specs/cli/spec.md` at landing

**Out:**

- **a loop-wide spending cap.** The worst case is printed instead. Returns if the user asks
  for a ceiling across rounds; the log already holds what it would sum.
- **a budget measured in rounds, wall time or tokens.** `--max` is already the round cap,
  and the host enforces dollars; a second unit is a second budget.
- **a price table, or any cost the binary computes.** Prices move under the repository.
- **stopping the loop on the first budget hit.** Stuck detection already bounds it.
- **interpreting `subtype`, `is_error` or `error`.** Recorded, never read back by the loop.
- **reading the run log anywhere but `metrics`.** No `status` line, no panel row, no
  `--json`. `add-json-output` can carry the block if it wants it.
- **budget and cost for `codex` and `opencode`.** `support-codex-opencode-in-loop` owns them.
- **rotating, pruning or migrating the log.**
- **checkpointing.** The ask's "checkpoint por iteración" already exists: each session
  commits its box, and the checklist is the state. Nothing new is written for it.
- **writing the log into the change folder or anywhere under `.agents/`.**

## Constraints

- Go 1.26.5, stdlib only — `encoding/json`, `os`, `os/exec`, `regexp`, `strconv`. No
  dependency.
- `loop` stays outside `needsPayload` and must keep working with only the binary installed.
- `--dry-run` keeps needing no `claude` on PATH and starts no session.
- `runLoop` stays drivable without a real `claude`: the runner seam carries the result.
- Tests point every log path at `t.TempDir()`; nothing writes into this repository's
  `.git` or a real `~/.claude`.
- `metrics` stays read-only; it opens the log for reading and never creates it.
- Every existing `loop` and `metrics` criterion in `.agents/specs/cli/spec.md` holds
  unchanged — stuck detection, the cap, a session error not failing the loop, the missing
  plan reported before the missing dependency, the token block's three states — **with
  exactly one exemption:** the `metrics` heading *what the flow cost, derived and never
  recorded* and its first bullet, *nothing is instrumented*, are rewritten as stated under
  Prior decisions. Nothing else in either section moves.

## Prior decisions

- **The budget is the host's `--max-budget-usd`, per session, which is per round.**
  Assumed 2026-10-22, attacca run. The host enforces it; the loop passes it through and
  records what was reported, so an overshoot is visible in the log rather than hidden.
  If wrong (a loop-wide cap was meant): sum the log's `cost_usd` and stop at a total.
- **The flag carries the host's name, `--max-budget-usd`, and accepts a plain decimal
  only** — `^[0-9]+(\.[0-9]+)?$`, greater than zero — passed through exactly as typed.
  Assumed 2026-10-22, attacca run and spec review. If wrong: a rename or a wider regex in
  one parser, and one help line.
- **No loop-wide cap; the header prints `--max` × budget to the cent.** Assumed
  2026-10-22, attacca run. **Ceiling:** adequate while somebody reads the header; the day a
  run is left unwatched with a large `--max`, a total cap is the replacement.
- **A budget-cut session gets no special case, and nothing the host reports is
  interpreted.** Assumed 2026-10-22, attacca run. `subtype` and `is_error` are recorded
  verbatim, and `error` carries the exec error's message, because the budget-hit spelling
  was not observed. If wrong: one stop reason keyed on the recorded subtype, once observed.
- **The run log lives at `<git-common-dir>/libretto/loop.jsonl`**, located with
  `git rev-parse --path-format=absolute --git-common-dir` after the checklist refusal and
  the PATH check, and never under `--dry-run`. Assumed 2026-10-22, attacca run and spec
  review. Never tracked, never in `git status`, no `.gitignore` written, shared by every
  worktree — which is why its messages say *repository*, never *checkout*. **Ceiling:** it
  is per machine — a fresh clone starts with no history, like the host's transcripts. If
  that stops being acceptable, it moves under `.agents/`.
- **The run log is the one recorded artifact, and its writer is the binary.** This
  replaces the `metrics` section's heading and first bullet in `.agents/specs/cli/spec.md`,
  which promised nothing was recorded. The heading becomes *`metrics` — what the flow
  cost, derived from git and the loop's own log*. The bullet becomes: **the change folders
  and the flow stay uninstrumented; the only recorded data is the loop's own run log,
  written by the binary and nothing else.** Its reason survives: the rule rejected
  artifacts eight skills would have to remember to write, and the loop is one writer that
  runs every round by construction and cannot forget — nothing lands in a change folder,
  and no skill learns anything. Assumed 2026-10-22, attacca run and spec review.
- **Cost is the host's `total_cost_usd`, recorded verbatim, absent when unreported, never
  zero.** Assumed 2026-10-22, attacca run. It extends the existing *no pricing* decision
  rather than relitigating it — the binary still prices nothing. The observed run labels it
  `"costBasis": "list"`, so it is shown as the host's estimate, summed with `%.2f` under the
  existing *plain digits, no exponent* rule.
- **`closed` is signed, and `metrics` reports it as a net figure.** Assumed 2026-10-22,
  spec review. If wrong: a second field splits closes from reopens.
- **`metrics` shows the log as its own block, sorted by change, never as a table column.**
  Assumed 2026-10-22, attacca run and spec review; the existing *no new column* decision
  already settles the table. The malformed-line count covers the whole file under any
  filter, because a malformed line names no change.
- **A log that cannot be located or written warns once and never stops the loop; a round
  whose checklist vanished is still recorded, as `stopped: "error"`.** Assumed 2026-10-22,
  attacca run and spec review. If wrong: a refusal before the first session.
- **The fork: A — host flag plus JSONL in the git common dir.** Assumed 2026-10-22,
  attacca run. B, a self-ignoring `.agents/loop/`, lost on worktrees; C, transcripts joined
  by session id, lost because it gives no budget and no dollars.

## Task breakdown

- ask every session for its JSON result; parse it into cost, duration, session id,
  subtype, `is_error` and text; print the text
- locate the log through git's common dir, after the existing refusals and never under
  `--dry-run`; append one record per round with every field; the two warnings; carry on
- mark the last record with why the loop stopped, including a vanished checklist
- parse and validate `--max-budget-usd`; pass it through; state the per-round budget and
  the worst case in the header; record the budget; name the flag in the help
- read the log in `metrics`; sorted rows with a net closed and a cent-rounded cost; the
  absent, filtered-empty and malformed states
- at landing, rewrite the `metrics` heading and its *nothing is instrumented* bullet

## Verification criteria

- **When the loop starts a session, the `claude` argument list it builds shall carry
  `-p <prompt>` and `--output-format json`.**
  Proof: cmd/libretto/loop_test.go TestEverySessionIsAskedForItsJSONResult
- **Where `--max-budget-usd <amount>` is given, the `claude` argument list shall carry
  `--max-budget-usd <amount>` with the amount exactly as typed; where it is not given, the
  list shall carry no `--max-budget-usd`.**
  Proof: cmd/libretto/loop_test.go TestClaudeArgsCarryTheBudgetOnlyWhenGiven
- **If `--max-budget-usd` is missing its value, or its value does not match
  `^[0-9]+(\.[0-9]+)?$`, or it equals zero, then `parseLoopArgs` shall return an error
  naming `--max-budget-usd`** — the table covers `1e2`, `0x1`, `+3`, `-1`, `inf`, `.5`,
  `0` and `0.00`.
  Proof: cmd/libretto/loop_test.go TestParseLoopArgsRefusesABadBudget
- **While a budget is set, the loop header shall print the per-round budget and the
  product of `--max` and that budget formatted with `%.2f` as the worst case for the run**
  — `--max 3 --max-budget-usd 0.1` prints `0.30`.
  Proof: cmd/libretto/loop_test.go TestTheHeaderStatesTheWorstCaseSpend
- **The loop help shall name `--max-budget-usd`.**
  Proof: cmd/libretto/main_test.go TestHelpNamesTheLoopBudget
- **When the session's output is the host's JSON result object, the session parser shall
  return its `total_cost_usd`, `duration_ms`, `session_id`, `subtype`, `is_error` and
  `result`; if the output is not JSON or carries no `total_cost_usd`, then the parser
  shall report the cost as unknown, never as zero.**
  Proof: cmd/libretto/loop_test.go TestParseSessionResultReadsTheHostsCost
- **When a session ends, the loop shall print the `result` text of its JSON to stdout; if
  the output does not parse as JSON, then the loop shall print it unchanged.**
  Proof: cmd/libretto/loop_test.go TestTheSessionsTextStillReachesTheTerminal
- **When a round ends, the loop shall append exactly one JSON line to the run log.**
  Proof: cmd/libretto/loop_test.go TestEveryRoundAppendsOneRecord
- **Each record shall carry `time` as an RFC 3339 timestamp of when the round ended,
  `change` as the change name, and `round` as the 1-based round number.**
  Proof: cmd/libretto/loop_test.go TestTheRecordSaysWhenWhichChangeAndWhichRound
- **Each record of a round whose checklist was read after the session shall carry
  `closed` as `after.done − before.done`, negative when a box reopened, and `open` as the
  open boxes after the round.**
  Proof: cmd/libretto/loop_test.go TestClosedIsSignedSoAReopenedBoxIsNegative
- **Where the session's JSON carried them, the record shall carry `cost_usd`,
  `duration_ms`, `session_id`, `subtype` and `is_error` with the host's values
  unchanged.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheHostsReportVerbatim
- **If the session reported no cost, then the round's record shall omit `cost_usd`
  entirely rather than write zero; if it reported no `is_error`, then the record shall
  omit `is_error`.**
  Proof: cmd/libretto/loop_test.go TestAnUnpricedRoundOmitsCostRatherThanWritingZero
- **Where a budget is set, each round's record shall carry it as `budget_usd`; where none
  is set, the record shall omit `budget_usd`.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheBudgetThatApplied
- **If the session exited non-zero, then the record shall carry the exec error's message
  as `error`; when it exited zero, the record shall omit `error`.**
  Proof: cmd/libretto/loop_test.go TestTheRecordCarriesTheExecErrorOnlyOnFailure
- **When the loop stops because every box closed, two rounds closed nothing, the cap was
  reached, or the checklist could not be read, the last record it wrote shall carry
  `stopped` set to `done`, `stuck`, `cap` or `error` respectively, and no earlier record
  shall carry `stopped`.**
  Proof: cmd/libretto/loop_test.go TestTheLastRecordSaysWhyTheLoopStopped
- **If the checklist cannot be read after a session, then the loop shall append that
  round's record with `closed` and `open` omitted, `stopped` set to `error` and the read
  error in `error`, and shall return the error it returns today.**
  Proof: cmd/libretto/loop_test.go TestAVanishedChecklistIsStillRecorded
- **The loop shall write the run log at `<git-common-dir>/libretto/loop.jsonl`, and
  `git status --porcelain` in the project shall print nothing after a loop over a clean
  tree.**
  Proof: cmd/libretto/loop_test.go TestTheRunLogLeavesGitStatusClean
- **If the git common dir cannot be located, then the loop shall print the one warning
  `run log not recorded: git common dir not found` and shall run every round it would
  have run with a log.**
  Proof: cmd/libretto/loop_test.go TestAnUnlocatableLogWarnsAndCarriesOn
- **If the run log cannot be written, then the loop shall print one warning naming the
  log's path, however many rounds fail to write, and shall run every round it would have
  run with a writable log.**
  Proof: cmd/libretto/loop_test.go TestAnUnwritableLogDoesNotStopTheLoop
- **When `loop` refuses a missing checklist or a missing `claude`, it shall not have
  resolved the run log path; while `--dry-run` is set, the loop shall not resolve the
  path, shall print no log warning, and shall write no record.**
  Proof: cmd/libretto/loop_test.go TestDryRunWritesNoRecord
- **The change folders and the flow shall stay uninstrumented; the only recorded data
  shall be the loop's own run log, written by the binary and nothing else.**
  Proof: cmd/libretto/metrics_test.go TestChangeNamesSeesLandedChangesNotJustOpenOnes
  Proof: cmd/libretto/metrics_test.go TestMetricsNeverWritesTheRunLog
- **Where a run log exists, `metrics` shall print a loop block with one row per change,
  sorted by change name, holding its rounds, the net sum of its `closed` values under a
  header that says net, the sum of its `cost_usd` values formatted with `%.2f`, and the
  count of its rounds with no `cost_usd`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsReportsLoopRoundsPerChange
- **If a change's rounds all lack `cost_usd`, then its cost cell shall print a dash, never
  `0.00`.**
  Proof: cmd/libretto/metrics_test.go TestAnUnpricedChangeReportsADashNotAZero
- **If no run log exists, then `metrics` shall print the one line `no loop runs are
  recorded in this repository`, and shall print every other block as it does today.**
  Proof: cmd/libretto/metrics_test.go TestNoRunLogIsAStateNotAnError
- **If a run log line does not parse as a record, then `metrics` shall skip it, keep every
  other row, and print the number of skipped lines counted over the whole file, whether
  or not a change is named.**
  Proof: cmd/libretto/metrics_test.go TestAMalformedRunLogLineIsSkippedAndCounted
- **Where a change is named, the loop block shall print only that change's row; if the
  named change has no record, then it shall print the one line `no loop rounds recorded
  for <change>`.**
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
  log path into `loop` and `metrics`, the help line), `cmd/libretto/main_test.go`. Six
  files, plus the delta onto `.agents/specs/cli/spec.md` at landing.
- `loopRunner` changes from `func(prompt string) error` to returning a parsed session
  alongside the error; every existing `loop_test.go` fake adapts to it.
- The log path is resolved outside the logic and passed in, the way `transcriptProjects()`
  is for the token block. For `loop` it is a resolver function in `loopOpts`, not a
  string, so `loop()` can call it late; tests hand in a resolver returning a `t.TempDir()`
  path, an empty string, or one that fails the test if called. `fakeGit` in
  `metrics_test.go` needs no new arm.
- `TestTheRunLogLeavesGitStatusClean` builds a real temporary repository, so `make
  test-short` skips it, per the suite's convention for real-git tests.
- Gates: the six in `AGENTS.md`.

### The approach

1. **Arguments.** `parseLoopArgs` gains `--max-budget-usd <amount>`, matched against
   `^[0-9]+(\.[0-9]+)?$`, then `strconv.ParseFloat` to refuse zero, and kept as the string
   the user typed so it reaches `claude` unchanged. `claudeArgs(prompt, budget string)
   []string` returns `-p <prompt> --output-format json`, plus the budget pair when set;
   `execClaude` uses it. The help line in `main.go` names the flag.
2. **The session result.** `execClaude` captures stdout into a buffer — stderr still goes
   straight to the terminal — and `parseSession(out []byte) session` decodes the one object
   into `{cost *float64, durationMS *int64, sessionID, subtype string, isError *bool,
   text string}`. Pointers, because unknown and zero are different facts. Output that is
   not JSON yields a session whose `text` is the raw output and whose other fields are
   absent. The text is printed to stdout when the session ends, which is when `-p` printed
   it before.
3. **The header.** With a budget, one more line: the per-round budget as typed, and
   `--max` × budget formatted `%.2f` as the run's worst case.
4. **Where the path is resolved.** `loop()` keeps its order — `openChecklist` refusal,
   then the PATH check — and only then, and only without `--dry-run`, calls
   `loopLogPath(git gitRunner) string`, which runs
   `rev-parse --path-format=absolute --git-common-dir` and returns
   `<dir>/libretto/loop.jsonl`, or "" on failure. An empty path prints
   `run log not recorded: git common dir not found` once, before the first round.
   `runLoop` returns from `--dry-run` before any of this exists.
5. **The record.** After each round, once the stuck and cap decisions are made, the loop
   builds `{time, change, round, closed *int, open *int, cost_usd, duration_ms,
   session_id, subtype, is_error, budget_usd, error, stopped}` — every optional field
   `omitempty` over a pointer or empty string — and appends it with one `Write` on a file
   opened `O_APPEND|O_CREATE|O_WRONLY`, after `MkdirAll` of `libretto/`. `closed` is
   `after.done − before.done`, signed. `budget_usd` is the parsed float. `error` is the
   exec error's message. `stopped` is set only on the round that ends the loop: `done`,
   `stuck`, `cap`, or `error` when `readPlan` fails after a session — that round's
   `closed` and `open` are nil, its `error` holds the read error, and the loop returns the
   error it returns today. A write failure prints one warning naming the path, the first
   only, and the loop carries on.
6. **`metrics`.** `metrics(...)` gains the log path as a parameter. `loopBlock` reads the
   whole file with `bufio.Reader` — the transcript reader's rule, no line-length cap —
   counting every line that fails to decode, then aggregates per change, applies the
   table's `only` filter, sorts by name, and prints rounds, net closed, cost sum (`%.2f`,
   or a dash when no round was priced) and the unpriced count. The malformed count is
   taken before the filter, so it always covers the whole file. Absent file: `no loop runs
   are recorded in this repository`. A named change with no record: `no loop rounds
   recorded for <change>`. The cost header says the figure is the host's estimate in
   dollars.
7. **At landing.** The delta rewrites the `metrics` section's heading and first bullet in
   `.agents/specs/cli/spec.md` to the text under Prior decisions, keeping
   `TestChangeNamesSeesLandedChangesNotJustOpenOnes` and adding
   `TestMetricsNeverWritesTheRunLog` as its proofs, and lands every criterion above
   beside the existing ones.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — a self-ignoring `.agents/loop/` holding a `.gitignore` of `*` | per worktree: the log dies with the worktree the loop ran in, and `metrics` from the main checkout never sees it |
| C — no new output; `metrics` joins host transcripts by `session_id` | gives no budget at all, and transcripts carry tokens, never dollars |
| a log appended inside the change folder | the session commits the folder, so the loop would race it; and it is deleted at landing — printed and lost again |
| `~/.claude/projects/<encoded>/` beside the transcripts | writes into the host's directory, which the token block's spec promises is never written |
| a loop-wide budget summed by the binary | answers a question the ask did not ask; the printed worst case covers it until somebody does |
| tokens as the budget unit | the host enforces dollars natively; a token cap would be a second enforcement the binary implements and the host does not know about |
| `strconv.ParseFloat` alone as the budget check | accepts `1e2`, `0x1p-2`, `+3` and `inf`, which then reach `claude` as typed with meanings nobody checked |
| one `error` field holding both the exec error and the host's `is_error` | two sources with different meanings; merging them is interpreting them |

### Risks

| Risk | What catches it |
|---|---|
| `--output-format json` hides output the user saw before | the loop prints `result` after each session; `TestTheSessionsTextStillReachesTheTerminal`. **Unobserved:** the text mode's exact output was not run this session, so equivalence is assumed |
| a host version without `--output-format json` or `total_cost_usd` | non-JSON output is printed raw and recorded with no cost; `TestParseSessionResultReadsTheHostsCost` covers both |
| the host overshoots its own budget | not caught, and not ours to enforce; the record carries the reported cost, so the overshoot is visible in `metrics`. Accepted |
| two loops in two worktrees append to one file at once | one `Write` per record on an `O_APPEND` descriptor, a record being a few hundred bytes. Accepted, not tested |
| the budget-hit subtype is never interpreted | deliberate; stuck detection bounds the spend at two barren rounds |
| resolving the log path moves a refusal | the resolver is called after both refusals; `TestDryRunWritesNoRecord` hands a resolver that fails the test if called, and `TestTheMissingPlanIsReportedBeforeTheMissingDependency` stays green unchanged |
| a test writes into this repository's `.git` | the path comes from an injected resolver; only `TestTheRunLogLeavesGitStatusClean` resolves a real one, inside its own temporary repository |
| an old `git` without `--path-format` | `loopLogPath` returns "" and the loop warns and continues. Accepted |

### Validation and rollback

New tests: twenty-seven, all named under Verification criteria; the twenty-eighth citation,
`TestChangeNamesSeesLandedChangesNotJustOpenOnes`, already exists. Three must be forced
red before they are believed, because they pass on the obvious wrong build:

- `TestAnUnpricedRoundOmitsCostRatherThanWritingZero` — a plain `float64` field writes
  `0`, and a test that only checks the record exists passes. Assert the key is absent,
  then break the pointer on purpose and watch it fail.
- `TestTheRunLogLeavesGitStatusClean` — passes trivially if no log was written at all.
  Assert the file exists and holds a record before asserting the clean status.
- `TestClosedIsSignedSoAReopenedBoxIsNegative` — passes on an unsigned build if the
  fixture never reopens a box. The fixture must reopen one, and an unsigned build must
  fail it.

Every existing `loop_test.go` test must stay green across the runner signature change;
`TestASessionErrorDoesNotAbortTheLoop`, `TestDryRunNeedsNoClaudeOnPath` and
`TestTheMissingPlanIsReportedBeforeTheMissingDependency` are the ones most likely to break.

Rollback: one revert. Nothing migrates; a log left in `.git/libretto/` is inert.

### Complexity deliberately kept

- Pointers for `cost_usd`, `duration_ms`, `is_error`, `closed` and `open`: absent and zero
  are different facts — the rule `metrics` already keeps for corrections and tokens.
- The log path as an injected resolver called late, not a string computed early: the
  `transcriptProjects()` seam, what keeps the suite off this repository's `.git`, and what
  keeps the existing refusal order intact.
- The record written after the stop decision rather than before it: one more branch, and
  the only way the last record can say why the loop stopped.
- A regex in front of `ParseFloat`: the parser alone accepts syntaxes that would reach the
  host unexamined.

## Checklist

Phase 6 runs these with `build-and-check`. Each box leaves the tree green on its own.

- [ ] **The loop records every round.** `claudeArgs` with `--output-format json`,
      `parseSession`, the result text printed, the late-resolved `loopLogPath`, one record
      per round with every field but `budget_usd`, `stopped` on the last including a
      vanished checklist, both warnings, nothing under `--dry-run`, every existing loop
      test adapted to the new runner. Closes: TestEverySessionIsAskedForItsJSONResult,
      TestParseSessionResultReadsTheHostsCost, TestTheSessionsTextStillReachesTheTerminal,
      TestEveryRoundAppendsOneRecord, TestTheRecordSaysWhenWhichChangeAndWhichRound,
      TestClosedIsSignedSoAReopenedBoxIsNegative, TestTheRecordCarriesTheHostsReportVerbatim,
      TestAnUnpricedRoundOmitsCostRatherThanWritingZero,
      TestTheRecordCarriesTheExecErrorOnlyOnFailure, TestTheLastRecordSaysWhyTheLoopStopped,
      TestAVanishedChecklistIsStillRecorded, TestTheRunLogLeavesGitStatusClean,
      TestAnUnlocatableLogWarnsAndCarriesOn, TestAnUnwritableLogDoesNotStopTheLoop,
      TestDryRunWritesNoRecord. Waits on: nothing.
- [ ] **The loop takes a per-round budget.** `--max-budget-usd` validated as a plain
      decimal and passed through as typed, the `%.2f` worst-case header line, `budget_usd`
      in the record, the help line in `main.go`. Closes:
      TestClaudeArgsCarryTheBudgetOnlyWhenGiven, TestParseLoopArgsRefusesABadBudget,
      TestTheHeaderStatesTheWorstCaseSpend, TestTheRecordCarriesTheBudgetThatApplied,
      TestHelpNamesTheLoopBudget. Waits on: box 1 (`claudeArgs` and the record).
- [ ] **`metrics` reads the run log.** The path passed from `main.go`, the sorted loop
      block with net closed and a `%.2f` cost, the absent, filtered-empty and malformed
      states, the cost label, and the proof that `metrics` never writes the log. Closes:
      TestMetricsReportsLoopRoundsPerChange, TestAnUnpricedChangeReportsADashNotAZero,
      TestNoRunLogIsAStateNotAnError, TestAMalformedRunLogLineIsSkippedAndCounted,
      TestTheLoopBlockFollowsTheChangeFilter, TestTheLoopBlockLabelsCostAsTheHostsFigure,
      TestMetricsNeverWritesTheRunLog. Waits on: box 1 (the record shape). Independent of
      box 2.
- [ ] **The `metrics` promise is rewritten at landing.** In `.agents/specs/cli/spec.md`,
      the heading *`metrics` — what the flow cost, derived and never recorded* and its
      bullet *nothing is instrumented* are replaced with the text under Prior decisions,
      citing TestChangeNamesSeesLandedChangesNotJustOpenOnes and
      TestMetricsNeverWritesTheRunLog. Closes: the uninstrumented-flow criterion. Waits on:
      boxes 1 and 3. Lands in the landing commit with the rest of the delta.
