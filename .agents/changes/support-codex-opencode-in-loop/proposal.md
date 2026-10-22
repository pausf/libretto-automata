# support-codex-opencode-in-loop

Tracker: none
Targets: cli
Durable decisions: every entry in Prior decisions below

## The ask, verbatim

> soporte codex/opencode — instalás a tres tools pero el loop solo maneja `claude`.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

`libretto loop` hands a fixed prompt to `claude` (`cmd/libretto/loop.go`) while install
already targets codex and opencode. Let the loop drive the other two tools it installs
to.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed unless it says measured)

- Q: can each tool run one headless session that reaches `/libretto-attacca`? → A:
  **opencode yes, codex no** (measured — see *Approach · Technical context* for the help
  lines and the source read). OpenCode's `run` takes `--command`, and a symlinked
  `.opencode/commands/libretto-attacca.md` resolved as the command `libretto-attacca` in
  an isolated `opencode debug config`. Codex installs skills only, its command equivalent
  (custom prompts) is deprecated and not targeted by the `targets` capability, and
  `codex debug prompt-input` showed `/libretto-attacca` reaching the model as literal text.
- Q: what happens on `--codex`? → A: **the loop refuses, before reading the checklist,
  naming codex, `/libretto-attacca` and why** (assumed). Faking support — a codex prompt
  that names the skills one by one — would restate in the loop what the attacca command
  owns, which `loopPrompt` exists to never do, and the 6→7 review's fresh reviewer is an
  agent codex has no destination for. If wrong (codex should run *something*): a codex
  prompt routed to `$build-and-check` becomes its own change, and its document has to say
  what replaces attacca's assume-never-ask rule there.
- Q: a new `--tool claude|codex|opencode` flag, or the existing tool axis? → A: **the
  existing axis — `libretto loop <change> --opencode`; no `--tool`** (assumed). `scopeFlags`
  already strips `--claude`, `--codex` and `--opencode` from every subcommand's arguments
  (measured, `cmd/libretto/main.go`), so today `loop c --opencode` silently drives
  `claude`. A `--tool` beside it would be a second spelling of one axis — `loop c --codex
  --tool opencode` is two answers to one question — and the `targets` rule is to extend an
  axis, never append a second. If wrong (`--tool` was wanted): it becomes an alias that
  errors when it disagrees with the axis flag, and the axis still decides.
- Q: `--max-budget-usd` with opencode? → A: **refused, naming `--max-budget-usd`,
  `opencode`, and that `opencode run` has no budget flag** (assumed; the absence measured
  in `opencode run --help`). The prior decision makes the budget the host's own flag, and
  the loop prices nothing itself, so there is nothing to pass it to. If wrong: the
  loop-wide cap the cli spec already names as its replacement, summing what the log holds —
  which for opencode holds no cost, so that change would first have to decide the cost.
- Q: what does an opencode round record for cost? → A: **nothing — `cost_usd` is omitted,
  and so are `duration_ms`, `subtype` and `is_error`** (assumed). OpenCode does report a
  cost, per step, on `step_finish` events (measured in the source), but `run` emits only
  the parts of its own session and skips every subagent's, so a sum would be a figure built
  here and short by exactly the 6→7 reviewer's spend. The prior decision is "the host's
  figure or absent"; a constructed partial figure is neither. If wrong (a partial cost is
  wanted): sum `part.cost` over the run's own `step_finish` events into a field that is not
  `cost_usd`, so `metrics` never adds it to the host's own totals.
- Q: does an opencode session get a permission flag? → A: **no — no `--auto`, no
  `--dangerously-skip-permissions`** (assumed). The loop passes `claude` none today; the
  user's own permission configuration governs both tools the same way. Measured in source:
  without `--auto`, `run` auto-rejects every permission ask and denies the `question` tool.
  If wrong: a round whose writes are rejected closes no box, stuck detection stops the
  loop after two rounds, and the fix is the user's opencode config, never a flag here.
- Q: how does the prompt reach opencode? → A: **`--command libretto-attacca`, with the
  prompt's arguments on stdin, never as a positional message** (measured in source):
  `run` wraps any positional argument containing a space in double quotes and escapes the
  quotes inside it, while piped stdin with no positional message becomes the message
  verbatim. If wrong (a future opencode changes stdin handling): the dry run shows the
  text, and the session's own output shows what arrived.
- Q: what does `--dry-run` print for opencode? → A: **the same prompt text it prints for
  claude, with the header naming the tool** (assumed). The text is the route and the task;
  the split into `--command` and stdin is plumbing the criteria pin down. If wrong: the
  dry run prints the opencode argument list as well.
- Q: does every record name its tool? → A: **yes, `tool` on every record the loop writes
  from now on** (assumed). One log is shared by every worktree and now by two tools, and an
  unpriced claude round must stay distinguishable from an opencode round that is unpriced
  by design. Records already written carry none and mean claude. If wrong: the field goes,
  and the two kinds of unpriced round become indistinguishable in `metrics`.
- Q: does `metrics` learn the tool? → A: **no** (assumed). It decodes `loopRecord`, an
  unknown field is ignored, and an opencode round already counts as unpriced. If wrong: a
  tool column is its own change.
- Q: the stop before code → A: **carried by the invocation** — attacca answers it; the
  document went to `review-spec` and the run acts on its findings as the author.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: how is "claude runs `claude`, opencode runs `opencode`" proved without a real CLI?
  → A: **a pure seam, `sessionFor(tool)`, returning the executable name and the argument
  builder (and, for the runner, the stdin and the parser); `loopWith` builds its runner
  from it, and the test asserts both tools' pairs, so a swapped wiring fails** (assumed).
  The first draft's criterion named `claudeArgs` and could not see which executable
  `loopWith` started. If wrong (the seam is too much for two tools): the switch stays in
  `loopWith` and the criterion is downgraded to the argument lists alone, with the wiring
  named as an accepted unobserved ceiling.
- Q: what does the stdin criterion prove? → A: **only the pure function
  `opencodeStdin`; that `execOpencode` actually sets `Stdin` from it is unobserved in
  tests, an accepted ceiling, with a Risks row** (assumed). No test starts `opencode`, so
  the claim "the session's stdin" could not be failed as written. If wrong (the wiring must
  be proved): `execOpencode` takes an injected command factory and a test reads the
  `exec.Cmd`'s `Stdin` back.
- Q: which spec text goes stale at landing besides the loop section? → A: **the cli
  spec's `init` paragraph**, "Scope flags are stripped by `scopeFlags` before dispatch and
  mean nothing to it, the same as `metrics` and `loop`." It becomes "Scope flags are
  stripped by `scopeFlags` before dispatch and mean nothing to it, the same as `metrics`;
  `loop` alone reads the tool flag." (assumed). If wrong (the sentence should not mention
  `loop` at all): it ends at `metrics`, and the loop section carries the tool flag alone.
- Q: what does the PATH refusal say on its second line? → A: **claude's stays word for
  word — `install Claude Code, or use --dry-run to see what would run` — so claude output
  is unchanged; opencode's is `install OpenCode, or use --dry-run to see what would run`**
  (assumed). The first draft generalised the line and would have moved claude's message.
  If wrong: one generic line for both, and the claude message change is named in the
  Outcomes.
- Q: what happens to a host `error` event? → A: **printed to the terminal after any text,
  prefixed `opencode error: `; when there was no exec error, its message also goes into
  the record's `error` field, as reported and never interpreted** (assumed). The message is
  the event's `error.data.message` when present and its `error.name` otherwise — the rule
  `run.ts` itself applies to print it — and several events join with a newline, as `run.ts`
  joins them. When the session exited non-zero, `error` keeps the exec error, unchanged
  from the claude path. If wrong (the event should not reach the record): the record keeps
  only the exec error, and an event on a zero exit is visible on the terminal alone.
- Q: what does the header print when `loopOpts` carries no tool? → A: **exactly what it
  printed before this change — no ` · via` suffix — the same reading the record gives an
  absent `tool`: claude** (assumed). Every existing test calling `runLoop` without a tool
  then prints what it printed before. If wrong (the header must always name a tool): an
  empty tool prints ` · via claude`, and the "unchanged output" promise for those tests
  goes.

## Outcomes

`libretto loop <change> --opencode` drives the change's checklist with one fresh
`opencode run` per round, exactly as the plain command drives `claude`: same checklist,
same stop rules, same run log. Each round hands OpenCode the installed
`/libretto-attacca` command, with the loop's task text as its arguments, and prints the
session's text to the terminal.

`libretto loop <change>` and `libretto loop <change> --claude` are unchanged in every
argument they hand `claude`, every message they print — the missing-`claude` refusal
included, word for word — and every field they record except the new `tool`; the header
gains ` · via claude`.

An OpenCode `error` event reaches the terminal after the session's text, as
`opencode error: <message>`, and reaches the record's `error` when the session exited
zero.

`libretto loop <change> --codex` refuses before doing anything, and the refusal says why:
Codex installs skills only, and every round is routed to a command Codex cannot run.

`libretto loop <change> --opencode --max-budget-usd 2` refuses before doing anything,
because `opencode run` has no per-session budget to pass the amount to.

Every run-log record names the tool that ran the round. An opencode record carries
`session_id`, `error` when the session exited non-zero or reported an `error` event, and never a cost.

## Scope boundaries

**In:**

- the existing tool axis (`--claude` · `--opencode` · `--codex`) reaching `loop`
- an opencode runner: its argument list, its stdin, its JSON event parser
- refusing codex, and refusing a budget for opencode
- `tool` on every run-log record, and the tool in the loop header and the help line
- the docs that name the loop's host: the `AGENTS.md` command block, and at landing the
  cli spec's loop prose, its codex/opencode non-goal, and its `init` paragraph's
  sentence on scope flags

**Out, named:**

- **codex in the loop.** Scoped out on evidence; see *Prior decisions* for what brings it
  back.
- **a `--tool` flag.** The axis already exists; see *Prior decisions*.
- **any cost for opencode**, summed, priced, or read back from `opencode export` or
  `opencode stats` after the round. The host's own figure for the session or nothing.
- **a budget for opencode**, by any mechanism — a wall-time kill, a token count, a
  loop-side sum.
- **permission flags for either tool.** The user's configuration governs.
- **`-m`/`--model`, `--agent` or `--variant` for opencode**, and any model choice for
  `claude`. The loop chooses no model today; OpenCode's default governs.
- **`--project` meaning anything to `loop`.** `scopeFlags` strips it today and the loop
  ignores it; that stays as it is.
- **checking that the payload is installed into the chosen tool** before the first round.
  The loop does not check it for `claude` either. A missing command surfaces as the
  session's own error and as two rounds that closed nothing.
- **`metrics` reading `tool`**, and any tool column in its report.
- **`--json` on `loop`**, still out for the reason the cli spec gives.

## Constraints

- Go 1.26.5, stdlib only. No dependency.
- **No real CLI in any test.** Every session is a `loopRunner` the test injects, as the
  loop's tests already do; argument lists, stdin and parsing are pure functions tested
  directly.
- `claudeArgs`, `parseSession` and `loopPrompt`'s text do not change by a byte —
  the existing tests that pin them stay as they are.
- `TestParseLoopArgs` compares `loopOpts` with `!=`, so every field added to it stays
  comparable.
- `CLAUDE_HOME` and the other roots are irrelevant here — `loop` writes nothing under any
  tool's root — but no test may reach a real `~/.config/opencode` either: the tests never
  start `opencode`.
- The vendor facts this rests on were read off **opencode 1.17.15** — its `--help`, and
  `sst/opencode` at tag `v1.17.15` (`packages/opencode/src/cli/cmd/run.ts`,
  `packages/schema/src/v1/session.ts`) — and off **codex-cli 0.148.0-alpha.9**. A later
  release can move them; the paths are written here so the next reader re-checks rather
  than re-guesses.
- An opencode event line has no length bound: a `tool_use` event carries the tool's whole
  output. The parser cannot use a line scanner with a fixed buffer.

## Prior decisions

From the cli spec, honoured unchanged:

- The budget is the host's own per-session flag, passed through as typed; the loop prices
  nothing.
- Nothing the host reports is interpreted; a field the host did not report is absent,
  never zero.
- The prompt routes and never instructs; it forbids pushing rather than omitting it.
- The missing checklist is reported before the missing dependency.
- `loopWith` exists so the run log path is injected without moving a call site.
- Destinations are tool × scope, two orthogonal axes; a new tool extends an axis, never a
  list (`targets`, and the flag table in the cli spec's *Scope: where it writes*).

Settled by this change, 2026-10-22, attacca run:

- **The loop's tool is the existing tool axis, never a `--tool` flag.** (assumed) Retires
  into the cli spec. `scopeFlags` already strips the three tool flags from every
  subcommand, so the loop silently drove `claude` under `--opencode` before this change;
  reading the axis fixes that rather than adding a second spelling. `--tool` lost because
  two flags answering one question is the error the axis rule exists to report. If wrong:
  `--tool` becomes an alias that errors on disagreement.
- **Codex is refused by the loop, on evidence.** (assumed scope, measured facts) Codex
  installs skills only (`targets`: "commands for codex" is out — custom prompts are
  deprecated); `codex debug prompt-input` rendered `/libretto-attacca` as literal user
  text, while it listed project `.agents/skills/` by name; and codex has no agent
  destination, so the flow's fresh reviewers have nothing to launch. **What brings it back:
  the payload shipping the attacca route as a skill Codex lists, or Codex installing
  commands** — then the refusal becomes a runner, and this entry says which.
- **OpenCode runs `opencode run --command libretto-attacca --format json`, with the loop's
  task text on stdin.** (measured) Positional arguments containing a space are re-quoted
  by `run`; piped stdin with no positional message is the message verbatim.
- **OpenCode gets no `--max-budget-usd`: the flag is refused under `--opencode`.**
  (assumed) `opencode run --help` lists no budget. If wrong, a loop-wide cap is the
  replacement the cli spec already names.
- **An opencode round records no cost, though OpenCode reports one per step.** (assumed)
  `run` emits only its own session's parts (`if (part.sessionID !== sessionID) continue`),
  so a subagent's `step_finish` — the 6→7 reviewer's — never reaches the stream; a sum
  would be a figure made here, and short. **Ceiling:** the day OpenCode's `run` emits a
  session total, or the children's steps, it is recorded verbatim as `cost_usd`.
- **No permission flag for either host.** (assumed) Parity with `claude` today, where the
  loop passes none. Without `--auto`, `opencode run` auto-rejects permission asks and
  denies `question`; a round so blocked shows up as a stuck loop, never as a silent pass.
- **Every record carries `tool`; a record without it means `claude`.** (assumed) Older
  records are not rewritten. The header reads an empty tool the same way: no suffix,
  exactly as before this change.
- **The tool's session is chosen through one pure seam, `sessionFor(tool)`.** (assumed,
  after review-spec) It is the only place an executable name is paired with its argument
  builder, so the pairing is tested rather than trusted. **Ceiling:** that
  `execOpencode` sets `Stdin` from `opencodeStdin` is not observed by any test, since none
  starts `opencode`; an injected command factory is the replacement if it ever bites.
- **A host `error` event is printed as `opencode error: <message>` after the text, and
  recorded in `error` only when the exec error is absent.** (assumed, after review-spec)
  The message is chosen by `run.ts`'s own rule, so nothing is interpreted here. If wrong,
  the record keeps the exec error alone.
- **The missing-`claude` refusal is unchanged word for word; OpenCode's names its own
  install.** (assumed, after review-spec)

## Task breakdown

- the tool reaching `loopWith` from `run`, and `loop()` keeping its signature on claude
- the codex refusal and the opencode budget refusal, before the checklist is read
- the PATH check naming the tool it looked for
- `sessionFor`, the one place an executable is paired with its argument builder
- the opencode argument list and the stdin it is handed
- the host `error` event, printed and recorded
- the opencode event parser: session id, text, raw fallback, unbounded lines
- `tool` on the record and in the header
- the help line and the `AGENTS.md` command block

## Verification criteria

### Which tool runs

- **When `run` receives `loop <change> --opencode --dry-run`, the loop shall print a header
  naming `opencode`.** Proved through `run`, because `scopeFlags` strips the flag before
  dispatch and a test of `loopWith` alone would prove nothing about that.
  Proof: cmd/libretto/loop_test.go TestTheToolFlagReachesTheLoop
- **Where `loopOpts` carries a tool, the loop header shall end its first line with
  ` · via <tool>`.**
  Proof: cmd/libretto/loop_test.go TestTheHeaderNamesTheTool
- **Where `loopOpts` carries no tool, the loop header shall be byte-identical to the
  header this change found** — the same reading the record gives an absent `tool`, so
  every test calling `runLoop` without a tool prints what it printed before.
  Proof: cmd/libretto/loop_test.go TestAnEmptyToolPrintsTheHeaderItAlwaysPrinted
- **`sessionFor("claude")` shall return the executable `claude` and an argument builder
  whose output equals `claudeArgs` for the same prompt and budget.** `loopWith` builds its
  runner from `sessionFor`, so a swapped wiring fails here.
  Proof: cmd/libretto/loop_test.go TestLoopDrivesClaudeWhenNoOtherToolIsNamed
- **`sessionFor("opencode")` shall return the executable `opencode` and an argument
  builder whose output equals `opencodeArgs`.**
  Proof: cmd/libretto/loop_test.go TestLoopDrivesClaudeWhenNoOtherToolIsNamed
- **If the tool is `codex`, then `loop` shall return an error containing `codex` and
  `/libretto-attacca` for a change that does not exist, with or without `--dry-run`,
  while PATH is empty and without resolving the run log path.** The refusal comes before
  the checklist because no checklist makes codex able to run it.
  Proof: cmd/libretto/loop_test.go TestLoopRefusesCodexAndSaysWhy
- **If `--max-budget-usd` is given while the tool is `opencode`, then `loop` shall return
  an error containing `--max-budget-usd` and `opencode` before reading the checklist.**
  Proof: cmd/libretto/loop_test.go TestOpencodeRefusesABudgetItCannotEnforce
- **If the tool is `opencode` and `opencode` is not on PATH, then `loop` shall return an
  error containing `opencode is not on PATH` and `install OpenCode, or use --dry-run to
  see what would run`.**
  Proof: cmd/libretto/loop_test.go TestTheMissingToolIsNamed
- **If the tool is `claude` and `claude` is not on PATH, then `loop` shall return the
  error it returned before this change, word for word: `claude is not on PATH — the loop
  relaunches it once per task`, a newline, seven spaces, and `install Claude Code, or use
  --dry-run to see what would run`.**
  Proof: cmd/libretto/loop_test.go TestTheMissingClaudeMessageIsUnchanged
- **If the tool is `opencode` and the change has no checklist, then `loop` shall return
  the `no checklist at` error while PATH is empty.** The missing plan is still reported
  before the missing dependency.
  Proof: cmd/libretto/loop_test.go TestTheMissingChecklistIsReportedBeforeTheMissingOpencode
- **While `--dry-run` is set and the tool is `opencode`, `loop` shall succeed with an empty
  PATH.**
  Proof: cmd/libretto/loop_test.go TestDryRunNeedsNoOpencodeOnPath

### What OpenCode is handed

- **The opencode argument list shall be exactly `run`, `--command`, `libretto-attacca`,
  `--format`, `json`, in that order.** No positional message, no budget, no model.
  Proof: cmd/libretto/loop_test.go TestOpencodeArgsRunTheAttaccaCommand
- **`opencodeStdin` shall return the loop prompt with its leading `/libretto-attacca `
  removed, byte for byte** — newlines and the double quotes around the change name
  included. Ceiling, named: this proves the pure function only; that `execOpencode` sets
  the session's `Stdin` from it is unobserved, because no test starts `opencode`.
  Proof: cmd/libretto/loop_test.go TestOpencodeArgumentsArriveOnStdinVerbatim
- **The opencode argument list shall contain none of `--auto`, `--yolo` and
  `--dangerously-skip-permissions`.**
  Proof: cmd/libretto/loop_test.go TestOpencodeArgsGrantNoPermissions

### What OpenCode reports

- **When the session's output is OpenCode's JSON event stream, the opencode parser shall
  return the `sessionID` its events carry as the session id.**
  Proof: cmd/libretto/loop_test.go TestParseOpencodeReadsTheSessionID
- **When the stream carries `text` events, the loop shall print their `part.text`, in the
  order they arrived, to stdout.**
  Proof: cmd/libretto/loop_test.go TestOpencodeTextReachesTheTerminal
- **If the output carries neither a `text` nor an `error` event, or its first line does
  not parse as JSON, then the loop shall print the raw output unchanged.** A session never
  reaches the terminal as nothing — the rule `parseSession` already keeps for `claude`.
  Proof: cmd/libretto/loop_test.go TestOpencodeWithoutTextPrintsTheRawOutput
- **When the stream carries an `error` event, the loop shall print `opencode error: `
  followed by its message after the session's text**, the message being the event's
  `error.data.message` when present and its `error.name` otherwise, several joined with a
  newline.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeErrorEventIsPrintedAndRecorded
- **If the stream carries an `error` event and the session exited zero, then the round's
  record shall carry that message as `error`, as reported.** On a non-zero exit `error`
  keeps the exec error, as it does for `claude`.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeErrorEventIsPrintedAndRecorded
- **If an event in the stream is longer than 1 MiB, then the opencode parser shall still
  return the session id and the text of the events after it.**
  Proof: cmd/libretto/loop_test.go TestALongOpencodeEventDoesNotEndTheParse
- **While the tool is `opencode`, each record shall omit `cost_usd`, `duration_ms`,
  `subtype` and `is_error`, even when the stream's `step_finish` events carry a `cost`.**
  The fixture must carry a non-zero `part.cost`, or an implementation that summed it would
  pass.
  Proof: cmd/libretto/loop_test.go TestAnOpencodeRoundRecordsNoCost
- **Each record the loop writes shall carry `tool` set to the tool that ran the round.**
  Proof: cmd/libretto/loop_test.go TestTheRecordNamesTheTool

### Said where it is read

- **The loop help shall name `--opencode` on a `loop` line.**
  Proof: cmd/libretto/main_test.go TestHelpNamesTheLoopTool

## Approach

### Summary

The loop already has the one seam this needs: a `loopRunner` that takes a prompt and
returns a `session`. This change gives it a second runner, for OpenCode, picked by the tool
axis `run` already parses, and refuses the two combinations no runner can honour — codex
at all, and a budget for opencode. Nothing in `runLoop`'s decisions moves: the checklist
is still the state, and the stop rules do not know which tool ran.

### Technical context

Go 1.26.5, stdlib. Tests in `cmd/libretto/loop_test.go` and `cmd/libretto/main_test.go`;
gates are the six in `AGENTS.md`.

**Blast radius:** `cmd/libretto/loop.go`, `cmd/libretto/main.go` (one call site, one help
line), `cmd/libretto/loop_test.go` (new tests, and the three `loopWith` calls in
`TestDryRunWritesNoRecord` gaining the tool argument), `cmd/libretto/main_test.go` (one
test), `AGENTS.md` (one line), and at landing `.agents/specs/cli/spec.md`.

**What the binaries said, observed 2026-10-22:**

OpenCode 1.17.15 (`opencode --version`). `opencode run --help`:

- `opencode run [message..]` — "run opencode with a message"
- `--command  the command to run, use message for args`
- `--format  format: default (formatted) or json (raw JSON events)`
- `--auto  auto-approve permissions that are not explicitly denied (dangerous!)`
- no option naming a budget, a cost or a spend limit
- top-level `--help`: `opencode stats  show token usage and cost statistics`,
  `opencode export [sessionID]  export session data as JSON`

Read off `sst/opencode` at tag `v1.17.15`, because the help does not say these:

- `run.ts`: "Non-interactive (default): sends a single prompt, streams events to stdout,
  and exits when the session goes idle." With `--command` it calls
  `client.session.command({ command: args.command, arguments: message })`.
- `run.ts`: each event goes out as one line, `JSON.stringify({ type, timestamp, sessionID,
  ...data })`; types emitted are `tool_use`, `step_start`, `step_finish`, `text`,
  `reasoning`, `error`. Parts of any other session are skipped: `if (part.sessionID !==
  sessionID) continue`.
- `run.ts`: positional words are joined after `arg.includes(" ") ? "\"" + escaped + "\""`;
  `resolveRunInput` returns piped stdin unchanged when no message was given.
- `run.ts`, exit codes: `process.exit(1)` on argument errors (no message and no command,
  session not found, a missing `--file`); `process.exitCode = 1` when the command or prompt
  call returns an error, when a `session.error` event arrived, or when the event loop
  threw; otherwise 0 when the session goes idle.
- `run.ts`: non-interactive sessions deny `question`, `plan_enter` and `plan_exit`; a
  `permission.asked` event is answered `reject` unless `--auto`.
- `packages/schema/src/v1/session.ts`: `StepFinishPart` carries `cost: Schema.Finite` and
  `tokens { input, output, reasoning, cache { read, write } }`.

And observed with the binary, isolated (`XDG_CONFIG_HOME`, `XDG_DATA_HOME`,
`XDG_STATE_HOME`, `XDG_CACHE_HOME` all under a scratch directory, nothing written to the
user's config): in a probe project holding `.opencode/commands/libretto-attacca.md` as a
symlink to this repository's `commands/libretto-attacca.md`, `opencode debug config` listed
the command `libretto-attacca` with the file's body as its template, `$ARGUMENTS` intact.

Codex. **The `codex` on PATH does not run on this machine**: `/opt/homebrew/bin/codex`
(npm `@openai/codex` 0.128.0) fails with `spawn …/vendor/aarch64-apple-darwin/codex/codex
ENOENT` for `--version` and `exec --help` alike. The help below is from the binary bundled
in ChatGPT.app, `codex-cli 0.148.0-alpha.9`. `codex exec --help`:

- `Run Codex non-interactively` · `Usage: codex exec [OPTIONS] [PROMPT]`
- `[PROMPT]  Initial instructions for the agent. If not provided as an argument (or if '-'
  is used), instructions are read from stdin.`
- `--json  Print events to stdout as JSONL` · `-o, --output-last-message <FILE>`
- no option naming a budget, a cost or a spend limit; exit codes not documented
- no option for a command or a slash command. The top-level command list has `exec`,
  `review`, `resume`, `fork`, `debug` — no prompts or commands subcommand.

And `codex debug prompt-input` ("Render the model-visible prompt input list as JSON") with
`CODEX_HOME` pointed at a scratch directory, in a probe project holding
`.agents/skills/build-and-check` and `.agents/skills/find-work` as symlinks: the developer
message listed both under `### Available skills` with their file paths, and the user
message was `/libretto-attacca continue the change "x". Use $build-and-check.` —
unexpanded. Codex reaches skills by name; a slash command reaches it as text.

No billed session was run for either tool: the help and the source settled every fact
the design depends on. What stays unobserved: OpenCode's real event stream end to end, and
codex's exit codes and usage events.

### The approach

1. **The tool reaches the loop.** `loopWith(projectDir string, tool target.Tool, args
   []string, logPath func() string)`; `main.go` passes the `tool` `scopeFlags` returned;
   `loop(projectDir, args)` keeps its signature and passes `target.ClaudeTool`.
   `loopOpts` gains `tool target.Tool`, set by `loopWith`, never by `parseLoopArgs` — so
   `TestParseLoopArgs`'s table needs no edit.
2. **Refusals, in order, all before `openChecklist`:** codex → an error naming codex,
   `/libretto-attacca`, that Codex installs skills only, and `--claude or --opencode`; a
   budget under opencode → an error naming `--max-budget-usd`, `opencode`, and that
   `opencode run` has no budget flag. Then the checklist, then — not under `--dry-run` —
   `exec.LookPath(string(o.tool))`. Its first line becomes `<tool> is not on PATH — the
   loop relaunches it once per task`, which for claude is today's text; its second line
   is per tool — `install Claude Code, or use --dry-run to see what would run` kept word
   for word, `install OpenCode, or use --dry-run to see what would run` for opencode.
3. **The prompt splits without changing.** `const attaccaRoute = "/libretto-attacca "`;
   `loopArguments(change, checklist)` returns today's text after it; `loopPrompt` returns
   `attaccaRoute + loopArguments(...)`, byte-identical to today. `opencodeStdin(prompt)`
   is `strings.TrimPrefix(prompt, attaccaRoute)`.
4. **The opencode runner** mirrors claude's three pieces: `opencodeArgs() []string`
   returns the fixed list; `parseOpencode(out []byte) session` walks the output with a
   `json.Decoder` — never a line scanner, so a long `tool_use` line cannot end the parse
   — decoding each value into `{type, sessionID, part{text}}`, keeping the first non-empty
   `sessionID` and joining every `text` event's `part.text` with a newline. A first value
   that fails to decode, or a stream with neither a `text` nor an `error` event, sets
   `text` to the raw output. An `error` event's message — `error.data.message` when
   present, else `error.name`, `run.ts`'s own rule — goes into a new `session.hostError`,
   several joined with a newline. A decode error after the first value stops the walk and
   keeps what was read. Cost, duration, subtype and is_error stay nil.
   `execOpencode(prompt)` runs `opencode` with those arguments, `Stdin` set to
   `opencodeStdin(prompt)`, stderr inherited.
5. **The runner is built from one pure seam.** `sessionFor(tool) (hostSession, bool)`,
   where `hostSession` holds `name string`, `args func(prompt, budget string) []string`,
   `stdin func(prompt string) string` (nil for claude) and `parse func([]byte) session`.
   claude → `claude`, `claudeArgs`, nil, `parseSession`; opencode → `opencode`, a builder
   returning `opencodeArgs()`, `opencodeStdin`, `parseOpencode`; anything else → false,
   never reached because codex is refused first. One `execHost(h, prompt, budget)` runs
   whatever `sessionFor` returned — `execClaude` and `execOpencode` collapse into it —
   and `loopWith` builds its runner from that pair alone.
6. **Header, error event and record.** Where `o.tool` is non-empty, the header's first
   line gains ` · via <tool>`; where it is empty the header is byte-identical to today's.
   After printing the session's text, `runLoop` prints `s.hostError`, when non-empty, as
   `<tool> error: <message>` — only opencode's parser ever sets it. `loopRecord` gains
   `Tool string` tagged `json:"tool,omitempty"`, set by `newRecord` from `o.tool`; and
   where `runErr` is nil and `s.hostError` is not, `newRecord` sets `error` to it. Tests
   that build `loopOpts` without a tool write records without one, which is the "absent
   means claude" reading.
7. **Help.** One line under the existing loop lines: `loop <change> --opencode   drive
   OpenCode instead of Claude Code; codex cannot run the flow`. `AGENTS.md`'s command
   block gains the same example.
8. **At landing**, besides applying this contract onto the cli spec's loop section and
   rewriting its non-goal "a budget or a cost for `codex` and `opencode`" into the
   decisions above, the `init` paragraph's sentence "Scope flags are stripped by
   `scopeFlags` before dispatch and mean nothing to it, the same as `metrics` and `loop`."
   becomes "Scope flags are stripped by `scopeFlags` before dispatch and mean nothing to
   it, the same as `metrics`; `loop` alone reads the tool flag." The loop section's
   "A round is one `claude -p` session" becomes "a round is one session of the chosen
   tool", with the budget paragraph scoped to `claude`.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| `--tool claude\|codex\|opencode` on `loop` | a second spelling of an axis `run` already parses and strips; two flags answering one question |
| codex driven with a prompt naming `$find-work`/`$build-and-check` | restates attacca's rules inside the loop, which `loopPrompt` exists never to do, and the flow's fresh reviewers have no agent destination in codex |
| the opencode prompt as one positional argument | `run` wraps a spaced argument in quotes and escapes its quotes — the session would see different text from the one `claude` sees |
| the opencode prompt as whitespace-split positional words | every newline in the prompt collapses to a space |
| `--format default` for opencode, reading plain text | no session id, so a round could not be matched to OpenCode's own record of it |
| summing `step_finish` costs into `cost_usd` | a figure built here, and short by every subagent's spend; `metrics` would add it to `claude`'s host totals as if it were one |
| reading the cost back with `opencode export <id>` after each round | a second call per round whose output shape this change has not observed, to produce a number the prior decision forbids constructing |
| `--auto` on opencode | the loop grants `claude` nothing; a permission the user did not configure is not the loop's to give |

### Risks

| Risk | What catches it |
|---|---|
| The payload is not installed into OpenCode, so `--command libretto-attacca` fails | `run` exits 1 with the error on stdout as an `error` event; the loop prints the raw output, records `error`, and stuck detection stops it after two rounds |
| A future OpenCode changes the event shape | the raw-output fallback prints whatever arrived; the record loses `session_id`, nothing else depends on it |
| A future OpenCode starts reading positional and stdin differently | the dry run shows the intended text; the session's own text shows what arrived. Not caught by a test, since no test runs OpenCode — accepted |
| Permission rejections leave a round unable to write | stuck detection, as above; the session's output names `auto-rejecting` |
| `execOpencode` builds the command but never sets `Stdin`, so OpenCode receives no arguments | no test catches it: `TestOpencodeArgumentsArriveOnStdinVerbatim` proves `opencodeStdin` only, and no test starts `opencode` — **accepted, unobserved**. On a real run it shows as an attacca session handed no change name, which closes no box, so stuck detection stops the loop in two rounds. Replacement if it bites: an injected command factory whose `exec.Cmd` a test reads back |
| The npm `codex` on PATH is broken on this machine | irrelevant to this change — codex is refused before PATH is consulted |
| An implementation reading `--opencode` itself in `parseLoopArgs` | never reached: `scopeFlags` strips it first; `TestTheToolFlagReachesTheLoop` goes through `run` so the stripping is in the path |

### Validation and rollback

The new tests above, plus every existing loop test unchanged — they are the proof that the
claude path did not move. Five need forcing red on purpose before they are believed:
`TestAnOpencodeRoundRecordsNoCost` (temporarily sum `part.cost` into `cost`; it must fail),
`TestALongOpencodeEventDoesNotEndTheParse` (swap the decoder for a default
`bufio.Scanner`; it must fail), `TestTheToolFlagReachesTheLoop` (pass
`target.ClaudeTool` from `main.go`; it must fail), `TestLoopDrivesClaudeWhenNoOtherToolIsNamed`
(swap the two names in `sessionFor`; it must fail), and
`TestAnOpencodeErrorEventIsPrintedAndRecorded` (drop the `runErr == nil` guard so a
non-zero exit's record takes the event's message; it must fail).

Rollback is one revert. Records already written with `tool` stay valid: `metrics` ignores
the field.

### Complexity deliberately kept

- **The codex refusal is code, not an absence.** Leaving codex out of the switch would make
  `loop c --codex` fall through to `claude` — the exact silent mis-drive this change found.
  The refusal is what makes the axis honest.
- **A `json.Decoder` instead of a scanner.** The scanner is shorter and fails on the first
  tool output over its buffer, which is the common case in a build round, not an edge.

## Checklist

Executed by `build-and-check`. Box 1 comes first; boxes 2 and 4 wait only on it; box 3
waits on box 2. The landing edits to the cli spec — the loop section, its codex/opencode
non-goal, and the `init` paragraph's scope-flag sentence, with the replacement text in
*Approach* step 8 — ride the landing commit and are not a box.

- [x] **`loop --opencode` drives OpenCode through `sessionFor`, `--codex` is refused, the
      header names the tool.** `loopWith` gains the tool parameter and `loopOpts.tool`;
      `loop()` passes claude; `main.go` passes `scopeFlags`' tool; the three `loopWith`
      calls in `TestDryRunWritesNoRecord` pass `target.ClaudeTool`; the codex refusal and
      the opencode budget refusal before `openChecklist`; the PATH check with the tool's
      name on line one and the per-tool install line, claude's word for word;
      `attaccaRoute`, `loopArguments`, `loopPrompt` rebuilt on them with identical output;
      `opencodeStdin`, `opencodeArgs`; `hostSession`, `sessionFor` and `execHost`, with
      opencode's parser returning the raw output as the session text until box 2; the
      header's ` · via <tool>` only where a tool is set. Closes:
      TestTheToolFlagReachesTheLoop, TestTheHeaderNamesTheTool,
      TestAnEmptyToolPrintsTheHeaderItAlwaysPrinted,
      TestLoopDrivesClaudeWhenNoOtherToolIsNamed, TestLoopRefusesCodexAndSaysWhy,
      TestOpencodeRefusesABudgetItCannotEnforce, TestTheMissingToolIsNamed,
      TestTheMissingClaudeMessageIsUnchanged,
      TestTheMissingChecklistIsReportedBeforeTheMissingOpencode,
      TestDryRunNeedsNoOpencodeOnPath, TestOpencodeArgsRunTheAttaccaCommand,
      TestOpencodeArgumentsArriveOnStdinVerbatim, TestOpencodeArgsGrantNoPermissions.
- [x] **OpenCode's event stream is read: session id, text, raw fallback, no line bound.**
      `parseOpencode` with a `json.Decoder`, set as opencode's parser in `sessionFor`.
      Waits on box 1. Closes: TestParseOpencodeReadsTheSessionID,
      TestOpencodeTextReachesTheTerminal, TestOpencodeWithoutTextPrintsTheRawOutput,
      TestALongOpencodeEventDoesNotEndTheParse.
- [x] **Every record names its tool, an opencode round records no cost, and a host error
      event is printed and recorded.** `loopRecord.Tool`, set in `newRecord`;
      `session.hostError` filled by `parseOpencode`, printed by `runLoop` after the text as
      `<tool> error: <message>`, and copied into `error` by `newRecord` only when `runErr`
      is nil; fixture streams whose `step_finish` carries a non-zero `part.cost` and whose
      `error` event carries `data.message`. Waits on box 2. Closes:
      TestTheRecordNamesTheTool, TestAnOpencodeRoundRecordsNoCost,
      TestAnOpencodeErrorEventIsPrintedAndRecorded.
- [x] **The help line and the `AGENTS.md` command block name `--opencode` for `loop`.**
      Waits on box 1. Closes: TestHelpNamesTheLoopTool.
