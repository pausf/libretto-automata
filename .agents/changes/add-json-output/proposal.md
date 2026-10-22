# add-json-output

Tracker: none
Targets: cli
Durable decisions: every entry in Prior decisions below

## The ask, verbatim

> `--json` en `status`/`doctor`/`metrics` — hoy nada externo puede consumirlos.

(Proposed in the 2026-08-21 feature-analysis session; accepted with "todas".)

## Reading

Machine-readable output for the read-only commands, already listed as pending in
`.agents/specs/cli/spec.md` for status and doctor.

## Decision log

### Session 2026-10-22 (attacca — nobody was asked; every answer below is assumed unless it says measured)

- Q: which pending text does this fulfil? → A: **the cli spec's Task breakdown line
  `5.2 --json for status and doctor`, the same box in `docs/PLAN.md` (5.2), and the
  `docs/STATE.md` row naming it** (measured: those three are the only places the spec and
  the docs name it). The queue non-goal in the cli spec's Scope boundaries — "no `--json`
  … machine-readable output for the read-only commands is `add-json-output`'s" — is the
  fourth pointer, and it is fulfilled here rather than left dangling. No parallel text is
  written: landing ticks 5.2 and rewrites the queue non-goal in place.
- Q: is `queue` in scope? → A: **yes** (assumed). The cli spec already assigns its
  `--json` to this change by name, the ask's reason — nothing external can consume it —
  holds for it word for word, and its data is one struct that already exists. Leaving it
  out would turn a pointer in the spec into a promise nobody keeps. If wrong, the queue box
  is dropped and the queue non-goal is rewritten to name a later change instead.
- Q: is `models` in scope, or any other read-only command? → A: **no** (assumed). `models`
  reads agent frontmatter a consumer can read itself, its listing carries a
  recommendations table whose wording moves with every agent added, and nobody asked.
  `land` already answers in its exit code; `version` is one line; `prune` and `uninstall`
  dry runs are plans for destructive commands, not reports; `wiki`, `init` and `loop`
  write. If wrong, `models --json` is its own change with its own shape — the helper this
  change adds makes it one box.
- Q: a top-level version field? → A: **yes — `"schema": 1` on every document, an integer**
  (assumed). While the tool is `0.x` the binary's version cannot be the contract, and a
  consumer needs one number to check before trusting a field. Adding a field never bumps
  it; removing, renaming or retyping one does, and is a promise removed — `!` on the
  commit. If wrong, the field is dropped before landing and consumers key on
  `libretto version` instead.
- Q: unknown versus zero? → A: **an unknown value is an absent key, never `0`, `""` or
  `null`; a known zero is present; a known-empty list is `[]`** (assumed, extending the
  house rule the loop run log set: cost "absent when unreported, never zero"). If wrong,
  unknowns become explicit `null` — a shape change, so it is decided before landing.
- Q: does `--json` change exit codes? → A: **no.** Every command exits exactly as it does
  without the flag — `doctor` non-zero when something needs attention, after its full
  JSON report (assumed). The exit code already carries the verdict and a second meaning
  would make the first unreadable — the reason the release notice was kept off it. If
  wrong, nothing here changes; a consumer that wants zero can ignore the code.
- Q: where do errors go, and in what shape? → A: **stderr, as the same plain
  `<invoked>: <message>` line `main` prints today, and nothing on stdout when the report
  never completed** (assumed). A consumer reads stdout as data and the exit code as the
  verdict; JSON on stderr would be a second error format to keep in step with the first.
  `doctor`'s "N item(s) need attention" stays on stderr after a complete report on stdout.
  An unreadable loop run log is the one error *inside* a report, because `metrics` exits 0
  on it today. If wrong, errors become `{"schema":1,"error":"…"}` on stderr, one helper in
  `main`.
- Q: colour and the TTY? → A: **JSON never carries an escape byte, TTY or not** (assumed).
  None of the four commands colours today (measured: `TestStatusOutputHasNoEscapeCodes`
  already pins status), so this is a guard on the JSON path against a future colour path,
  not a fix. If wrong — nothing; there is no other honest answer.
- Q: does the release notice stay on under `--json`? → A: **yes, on stderr as now**
  (assumed). Stdout stays pure, which is the promise the notice was placed on stderr to
  keep. If wrong, `carriesNotice` learns the flag and skips it.
- Q: pretty or compact? → A: **two-space indent through `json.Encoder`, one document, one
  trailing newline; whitespace is not part of the contract** (assumed). `jq` reads either,
  and a human reading a bug report reads the indented one. If wrong, `SetIndent` goes.
- Q: HTML escaping? → A: **off — `SetEscapeHTML(false)`** (assumed). A reading or a phase
  name carrying `<`, `>` or `&` reaches the consumer as written, not as `<`. If wrong,
  the default comes back and only the escaped bytes differ.
- Q: field naming? → A: **snake_case** (assumed), the convention the loop run log already
  set (`cost_usd`, `session_id`). If wrong, a rename before landing — after it, a schema
  bump.
- Q: where does `--json` go on the command line? → A: **after the command name, anywhere
  among its arguments; `libretto --json status` stays an unknown command** (assumed). It is
  a per-command flag, and only four commands honour it; a global flag would need every
  other command to refuse it. If wrong, `scopeFlags` strips it and the four read it from
  there.
- Q: what may `status` and `doctor` be given besides `--json`? → A: **nothing — any other
  argument exits non-zero naming it** (assumed). Today both ignore arguments silently, so
  `status --jsno` would print text into a consumer's parser; `metrics` and `queue` already
  refuse. If wrong, unknown arguments go back to being ignored. **Reversed after
  review-spec — see the next session.**
- Q: what goes into each shape? → A: **the data the text report shows, never its prose,
  and nothing a consumer can compute from the rest** (assumed). Out: `doctor`'s remedy
  text and prerequisite notes, `status`'s tally, `metrics`' footer totals, the token
  percentage, the legend and the ceiling. Adding a field later costs nothing; removing one
  is a broken promise, so the shapes start small. If wrong, each one is an additive field
  with no schema bump.
- Q: golden files, as `docs/PLAN.md` 5.2 asks? → A: **no — tests decode the output and
  assert fields** (assumed). There is no `testdata/` in this repository, and a golden
  file pins whitespace this contract says is not promised. If wrong, a `testdata/` golden
  per command, regenerated by a flag.
- Q: does the README move? → A: **no, and no `spec-readme.md`** (assumed). `readme`'s
  promise is one line per command, and no command is added; `--json` reaches users
  through the help text, which this change does move. If wrong, one README line and a
  `spec-readme.md` delta.
- Q: the fork — how is each command built? → A: **A, gather once and render twice: each
  command's derivation moves into one function returning a report struct, and the text
  renderer and the JSON encoder both read that struct** (assumed, the recommended option).
  B — a second function per command deriving the JSON from the same inputs — lost: it is a
  second derivation, a second source of truth, and the two drift the first time somebody
  fixes one. C — marshalling the existing internal types directly, tags on
  `link.Entry`, `changeMetrics`, `loopTally` — lost: their fields are unexported, carry
  bytes and paths no consumer needs (`Generated`, `RepoPath`), and a tag on an internal
  type makes every refactor of it a contract change.

### Session 2026-10-22, after review-spec (the run acting as author on the reviewer's findings)

- Q: how does `queue`'s existing refusal criterion move? → A: **amended in place at
  landing.** The cli spec's "If the command is given any argument, then it … exit
  non-zero naming the argument" becomes "any argument other than `--json`", its wording
  otherwise unchanged. `TestQueueRefusesArguments` keeps refusing every other argument and
  gains a row accepting `--json` (assumed). A second criterion beside the old one would
  leave the spec promising both that `--json` is refused and that it is read. If wrong,
  the old criterion is retired and a new one written in its place, with the same test.
- Q: do `status` and `doctor` still refuse arguments other than `--json`? → A: **no —
  reversed. They keep ignoring every other argument, as today; only `--json` is read**
  (assumed). Refusing would change an exit code on a path that has nothing to do with
  JSON, and "exit codes and text unchanged on every path" then holds with no exception.
  If wrong, a refusal criterion and test come back, in a change of their own.
- Q: what does `metrics --json` carry when git has seen no change? → A: **gather computes
  the ledger, the tokens and the loop anyway; the JSON carries `changes: []` plus
  whichever of `corrections`, `tokens` and `loop` exist; the text still prints only the
  no-changes line, byte-identical** (assumed). Those blocks do not depend on a change
  having been seen — a run log or a ledger can exist before any change folder — and
  omitting them would report as unknown what was read. If wrong, the JSON carries
  `changes: []` alone, as the text does.
- Q: when is `loop` omitted? → A: **only when no run-log file exists.** A log that exists
  with no records, or with only malformed lines, gives `loop: {changes: [], skipped: N}`
  (assumed). The file existing is a known fact; the text's "no loop runs are recorded"
  covers both cases because a human does not need the difference, and a consumer does.
  If wrong, `loop` is omitted whenever it holds no record, and `skipped` is lost with it.
- Q: what does `queue --json` carry without a `.agents/changes/`? → A: **`queued: []` and
  no `dispatched` key**, because the command returns before git is asked for branches,
  so which changes are dispatched is unknown (assumed) — the absent-means-unknown rule
  applied. If wrong, `dispatched: []` is written there and git is asked first.
- Q: how is "the release line is built from the same facts" observed? → A: **one stubbed
  check feeds both outputs.** With latest `v9.9.9` against running `v0.1.0`, the text line
  names `v9.9.9` and the JSON carries `latest: "v9.9.9"` and `newer: true`; with the stub
  changed to an unrankable pair, the text line carries no newer-release wording and the
  JSON omits `newer`. `TestReleaseLineIsRenderedFromTheCheck` asserts both outputs from
  the same stub (assumed). If wrong, two separate tests, and the shared derivation is
  review's to see.

## Outcomes

`libretto status --json`, `doctor --json`, `metrics --json` and `queue --json` each write
one JSON object to stdout and nothing else, and a script can read the same facts the text
report shows a human. Each object carries `"schema": 1`. Errors, warnings and the release
notice stay on stderr exactly as they are; exit codes do not move. Without `--json` every
command's output is unchanged.

The four shapes are the contract. Every key below is a promise; a key marked *omitted
when* is absent — never zero — in that state.

### `status --json`

    {
      "schema": 1,
      "target": { "name": "claude", "root": "/Users/u/.claude", "configured": true },
      "items": [
        { "kind": "skills", "name": "evidence", "state": "linked" },
        { "kind": "agents", "name": "spec-reviewer", "state": "wrong target" }
      ]
    }

- `target.name` — the target's name as the text header prints it: `claude`, `project`,
  `codex`, `opencode`.
- `target.root` — its root directory; *omitted when* the target has none (a project
  target with no directory).
- `target.configured` — whether the root exists on disk; the text prints `not configured`
  for `false`.
- `items` — every scanned entry in the scan's order; `[]` when there is nothing to link.
- `kind` — `skills`, `agents` or `commands`. `state` — `linked`, `missing`,
  `wrong target`, `conflict` or `stale`, spelled as `link.State` spells them.

### `doctor --json`

    {
      "schema": 1,
      "target": { "name": "claude", "root": "/Users/u/.claude", "configured": true },
      "problems": [
        { "kind": "skills", "name": "alpha", "state": "missing" }
      ],
      "release": {
        "mode": "installed",
        "payload_root": "/Users/u/.local/share/libretto/current",
        "running": "v0.5.2",
        "latest": "v0.6.0",
        "newer": true
      },
      "prerequisites": [
        { "name": "jira", "found": false },
        { "name": "rg", "found": true }
      ]
    }

- `target` — the same object as `status`.
- `problems` — every entry whose state needs attention, same item shape as `status`;
  `[]` when the text says `no problems`.
- `release.mode` — `checkout` or `installed`. `release.payload_root` — the payload root
  the text names after the mode.
- `release.running` — the version the binary reports, always present.
- `release.latest` — the newest release the check found; *omitted when* the check failed
  or the remote has no releases.
- `release.newer` — `true` when `latest` ranks above `running`, `false` when they are
  equal; *omitted when* there is no `latest` or the two do not rank (a `dev` build, or a
  running version ahead of the remote).
- `release.error` — the failed check's first line; *omitted when* the check answered.
- `prerequisites` — one object per prerequisite in the text's order, `name` and `found`.

### `metrics --json`

    {
      "schema": 1,
      "changes": [
        {
          "name": "add-converge-command", "state": "landed",
          "commits": 9, "first": "2026-10-21T18:02:11Z", "last": "2026-10-22T09:40:03Z",
          "closed": 3, "boxes": 3, "reopened": 0, "corrections": 1
        },
        { "name": "broken-one", "state": "unreadable", "corrections": 0 }
      ],
      "corrections": { "by_phase": { "build-and-check": 4, "6→7": 2 }, "outside_changes": 1 },
      "tokens": {
        "attributed":   { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
        "unattributed": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
        "entries": 8944, "unattributed_entries": 312,
        "change": {
          "name": "add-converge-command",
          "totals":   { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 },
          "by_phase": { "write-spec": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 } },
          "no_phase": { "input": 1, "output": 2, "cache_write": 3, "cache_read": 4 }
        }
      },
      "loop": {
        "changes": [
          { "name": "add-converge-command", "rounds": 3, "net_closed": 3, "cost_usd": 4.12, "unpriced": 0 }
        ],
        "skipped": 0
      }
    }

- `changes` — one object per row of the table, in its order; `[]` when git has seen no
  change, and then `corrections`, `tokens` and `loop` are still present whenever they
  exist — the text prints only its no-changes line. A `<change>` argument narrows it
  exactly as it narrows the table.
- `name`, `state` (`landed`, `in flight`, `unreadable` — the table's spelling) — always.
- `commits`, `first`, `last` — *omitted when* `unreadable`. Times are RFC 3339 in UTC,
  from the commit times `measure` already reads; the span is `last − first` and is not
  repeated.
- `closed` (boxes closed now), `boxes` (the denominator), `reopened` — *omitted when* the
  change has no checklist in its history (the text's `—`) or is `unreadable`.
- `corrections` on a change — *omitted when* there is no `.agents/lessons.md`; `0`
  included when there is.
- top-level `corrections` — *omitted when* there is no ledger. `by_phase` keys spelled as
  the ledger spells them, `{}` for an empty ledger; `outside_changes` is the orphan count.
- `tokens` — *omitted when* no transcripts are configured or none exist for this
  repository. `attributed`, `unattributed`, `entries` and `unattributed_entries` are
  corpus-wide whatever the filter, as the text's are.
- `tokens.change` — *omitted when* no `<change>` was given. `totals` and `by_phase` —
  *omitted when* no session was attributed to that change, leaving `name` alone.
  `no_phase` — *omitted when* every attributed entry named a phase.
- `loop` — *omitted when* no run-log file exists, and only then. `changes` — per change in
  name order, following the filter; `[]` when the log holds none for the filter, holds no
  record, or holds only malformed lines. `cost_usd` — *omitted
  when* none of that change's rounds carried a cost. `skipped` — malformed lines, `0`
  included. `loop.error` — present, in place of `changes` and `skipped`, only when the log
  exists and cannot be read.
- The legend and the ceiling are prose for a human reading the table and are not in the
  JSON; this section is their machine-readable counterpart.

### `queue --json`

    {
      "schema": 1,
      "queued": [
        { "name": "report-drift-trend-in-metrics", "date": "2026-08-21",
          "reading": "Derive per-capability drift events from git history." }
      ],
      "dispatched": [
        { "name": "add-json-output", "branch": "feat/add-json-output" }
      ]
    }

- `queued` — oldest date first, ties by name, as the text; `[]` when nothing is queued
  or no `.agents/changes/` exists.
- `reading` — *omitted when* the proposal has no `## Reading`.
- `dispatched` — changes an unmerged branch already took, `[]` when none; *omitted when*
  git could not list branches, or when no `.agents/changes/` exists and git was never
  asked — either way nobody knows which are dispatched.

## Scope boundaries

In:

- `--json` on `status`, `doctor`, `metrics` and `queue`, with the four shapes above
- `queue`'s refusal criterion amended in place to accept `--json`
- one help line naming the flag and the four commands
- the refactor each needs so text and JSON read one gathered report

Out, each named so it cannot arrive quietly:

- **`--json` on `models`, `land`, `version`, `prune`, `uninstall`, `wiki`, `init`,
  `loop`, `preview`.** Returns when a consumer names one; the helper makes each a box.
- **JSON errors on stderr.** Plain lines stay; returns if a consumer cannot use the exit
  code.
- **a global `--json`, a `--format` flag, or JSON Lines streaming.** One flag, four
  commands, one document each.
- **a published JSON Schema file.** The shapes live in this spec; a schema file is a
  second copy.
- **prose in the JSON**: `doctor`'s remedies and prerequisite notes, `metrics`' legend,
  ceiling, footer totals and token percentage, `status`'s tally. Each is additive later.
- **golden files.** Tests decode and assert.
- **README changes** — `readme`'s capability, untouched.
- **any change to the text output, exit codes, the release notice, or `needsPayload`, on
  any path.** `status` and `doctor` still need the payload; `metrics` and `queue` still
  do not.
- **`status` and `doctor` validating their arguments.** They read `--json` and ignore the
  rest, as they ignore everything today. Returns as its own change if a typo printing
  text into a parser bites somebody.
- **the drift trend** `report-drift-trend-in-metrics` queues for `metrics` — when it
  lands, its block gets a JSON key under this change's rules, in that change.

## Constraints

- Go 1.26.5, **stdlib `encoding/json` only** — no dependency, as `distribution` and every
  other JSON reader in this binary already do.
- The six gates in `AGENTS.md`.
- Every test against `t.TempDir()` with `CLAUDE_HOME` set; `metrics` and `queue` tests
  drive the injected `gitRunner`, never this repository's history.
- **No test reaches the module proxy.** `doctor` asks it live today, and
  `TestDoctorExitsNonZeroWhenSomethingNeedsAttention` reaches it from `go test` (measured:
  `doctor` calls `dist.Latest` with no seam). This change adds the seam and stubs it in
  `TestMain`, the way `noticeSource` already is.
- Text output byte-identical to today's for every command, with or without the refactor.
- Exit codes and stderr unchanged for every path; `main` still prints the error.
- `status` and `doctor` print through `fmt.Printf` to `os.Stdout`; `metrics` and `queue`
  take an `io.Writer`. The JSON path follows whichever the command already uses.

## Prior decisions

Every entry below is assumed under attacca on 2026-10-22 unless it says measured, and
names what changes if it is wrong.

- **Every JSON document carries `"schema": 1`.** Assumed. A removed, renamed or retyped
  key bumps it and is a promise removed; an added key does not. If wrong, it is dropped
  before landing and consumers key on `libretto version`.
- **An unknown value is an absent key; a known zero is present; a known-empty list is
  `[]`.** Assumed, extending *loop cost … absent when unreported, never zero*. If wrong,
  unknowns become `null`, decided before landing.
- **`--json` changes stdout only.** Exit codes, stderr lines and the release notice are
  what they are without it. Assumed. If wrong, nothing moves here.
- **Errors are plain lines on stderr, and a report that never completed writes nothing to
  stdout.** Assumed. The one exception is `metrics`' unreadable run log, an error inside a
  report that exits 0 today, which becomes `loop.error`. If wrong, an
  `{"schema":1,"error":…}` object on stderr.
- **The JSON carries the report's data, never its prose, and nothing derivable from the
  rest.** Assumed. If wrong, each missing field is additive.
- **snake_case keys, two-space indent, HTML escaping off.** Assumed; the first from the
  loop run log, the rest because whitespace is not promised and a reading must arrive as
  written. If wrong, a rename is a schema bump after landing.
- **One gathered report per command, rendered twice.** Assumed, the fork's answer. A
  second derivation for JSON lost as a second source of truth; tags on internal types
  lost because every refactor would become a contract change. **Ceiling:** the text
  renderers may still compute a display value from the report — a span, a percentage —
  but never a fact the report does not hold.
- **`--json` is a per-command flag after the command name.** Assumed. If wrong,
  `scopeFlags` strips it globally.
- **`status` and `doctor` read `--json` and ignore every other argument, as today.**
  Assumed after review-spec, reversing the first session's refusal: refusing would move
  an exit code on a path unrelated to JSON. If wrong, a refusal lands as its own change.
- **`queue`'s refusal criterion is amended in place, never duplicated.** Assumed after
  review-spec: "any argument" becomes "any argument other than `--json`", and
  `TestQueueRefusesArguments` gains a `--json` row. If wrong, the old criterion is retired
  and replaced, with the same test.
- **With no change in git's history, `metrics --json` still carries the ledger, tokens and
  loop it read.** Assumed after review-spec; the text keeps its single no-changes line.
  If wrong, the JSON carries `changes: []` alone.
- **`loop` is absent only when no run-log file exists.** Assumed after review-spec: an
  existing log with no records, or only malformed lines, is `{changes: [], skipped: N}`.
  If wrong, an empty log omits `loop` and its `skipped` count.
- **`queue --json` without a `.agents/changes/` omits `dispatched`**, since git was never
  asked. Assumed after review-spec. If wrong, `dispatched: []` and git asked first.
- **The release text and the release JSON come from one `releaseCheck`, observed through
  one stub.** Assumed after review-spec. If wrong, two separate tests.
- **`queue` is in scope; `models` and every other read-only command are out.** Assumed;
  the cli spec already assigned `queue` here. If wrong, a box is added or dropped.
- **`doctor`'s release check is injected, and the suite stubs it.** Assumed, the
  `noticeSource` precedent. If wrong, nothing user-facing moves.

## Task breakdown

- a shared helper: detect `--json` among a command's arguments, and one encoder with
  indent and HTML escaping off
- `status`: gather target and entries once; text and JSON from it; other arguments
  ignored as today
- `doctor`: the release check returning facts rather than a sentence, behind a seam; the
  report gathered once; text and JSON from it; other arguments ignored as today
- `queue`: kept and dropped ideas gathered once, with whether branches were asked and
  read; text and JSON from it; the refusal accepting `--json`
- `metrics`: rows, ledger, token totals and loop tallies gathered once into one report,
  with or without any change in history; every text block rendered from it; JSON from it
- the help line

## Verification criteria

Shared — each test runs one row per command, and each command's box adds its row.

- **When `status`, `doctor`, `metrics` or `queue` runs with `--json`, the command shall
  write exactly one JSON object to stdout followed by one newline.**
  Proof: cmd/libretto/json_test.go TestJSONIsOneObjectOnStdout
- **The JSON object each command writes shall carry a top-level `schema` key whose value
  is the integer 1.**
  Proof: cmd/libretto/json_test.go TestJSONCarriesSchemaOne
- **Where `--json` is given, the command shall exit with the same code it exits with on
  the same state without `--json`.**
  Proof: cmd/libretto/json_test.go TestJSONDoesNotChangeTheExitCode
- **If a command run with `--json` fails before its report is complete, then it shall
  write nothing to stdout.**
  Proof: cmd/libretto/json_test.go TestJSONFailureWritesNothingToStdout
- **If a command fails, then the error it returns under `--json` shall equal the error it
  returns without `--json`.**
  Proof: cmd/libretto/json_test.go TestJSONErrorsMatchTheTextErrors
- **The JSON a command writes shall contain no escape byte (0x1b).**
  Proof: cmd/libretto/json_test.go TestJSONHasNoEscapeCodes
- **While a newer release is cached, the stdout of `status --json` shall parse as one JSON
  object.**
  Proof: cmd/libretto/notice_test.go TestJSONStdoutParsesWithANoticePending
- **The help text shall carry a `--json` line naming `status`, `doctor`, `metrics` and
  `queue`.**
  Proof: cmd/libretto/main_test.go TestHelpNamesJSON

`status`

- **When `status --json` runs, the object shall carry `target` with the target's `name`,
  its `root` and `configured`.**
  Proof: cmd/libretto/json_test.go TestStatusJSONNamesTheTarget
- **If the target has no root, then `target` shall omit `root`.**
  Proof: cmd/libretto/json_test.go TestStatusJSONOmitsAnUnknownRoot
- **When `status --json` runs, `items` shall list one object per scanned entry, carrying
  `kind`, `name` and `state`, in the scan's order.**
  Proof: cmd/libretto/json_test.go TestStatusJSONListsEveryItem
- **The `state` of every item shall be spelled exactly as the text report spells it.**
  Proof: cmd/libretto/json_test.go TestStatusJSONStatesMatchTheText
- **While the payload holds nothing to link, `items` shall be an empty array.**
  Proof: cmd/libretto/json_test.go TestStatusJSONEmptyIsAnEmptyArray

`doctor`

- **When `doctor --json` runs, `problems` shall list exactly the entries whose state needs
  attention, each carrying `kind`, `name` and `state`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONListsOnlyProblems
- **When `doctor --json` runs, `release.mode` shall be `checkout` for a payload root
  holding `.git` and `installed` otherwise.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONNamesTheMode
- **The `release` object shall carry `running`, the version the binary reports.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONCarriesTheRunningVersion
- **When the release check answers with a release, `release.latest` shall carry it.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONCarriesTheLatestRelease
- **Where `latest` and `running` rank, `release.newer` shall be true when `latest` ranks
  above `running` and false when they are equal.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONRanksWhatRanks
- **If `latest` is absent or the two versions do not rank, then `release` shall omit
  `newer`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONOmitsAnUnrankableComparison
- **If the release check fails, then `release` shall carry `error`, the failure's first
  line, in place of `latest`.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONReportsAFailedCheck
- **When `doctor --json` runs, `prerequisites` shall list one object per prerequisite,
  carrying `name` and `found`, in the text's order.**
  Proof: cmd/libretto/json_test.go TestDoctorJSONListsPrerequisites
- **When one stubbed release check feeds both outputs, the release text line and the
  `release` object shall agree: for latest `v9.9.9` against running `v0.1.0` the line
  names `v9.9.9` while the object carries `latest` `v9.9.9` and `newer` true, and for an
  unrankable pair the line carries no newer-release wording while the object omits
  `newer`.**
  Proof: cmd/libretto/doctor_test.go TestReleaseLineIsRenderedFromTheCheck

`metrics`

- **When `metrics --json` runs, `changes` shall list one object per row of the table, in
  the table's order.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONListsTheTableRows
- **For every readable change, `commits` in the JSON shall equal the table's commits
  column on the same history.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAgreesWithTheTable
- **Every readable change object shall carry `first` and `last` as RFC 3339 times in
  UTC.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCarriesTimesInUTC
- **Where a change has a checklist in its history, its object shall carry `closed`,
  `boxes` and `reopened`, zero included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONKeepsAKnownZero
- **If a change has no checklist in its history, then its object shall omit `closed`,
  `boxes` and `reopened`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsChurnWithoutAChecklist
- **If a change's history cannot be read, then its object shall carry no key beyond
  `name`, `state` set to `unreadable`, and `corrections`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnreadableRowCarriesOnlyWhatIsKnown
- **While `.agents/lessons.md` exists, every change object shall carry `corrections`, zero
  included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCorrectionsZeroIsPresent
- **If no `.agents/lessons.md` exists, then the report shall carry no `corrections` key
  at any level.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAbsentLedgerOmitsCorrections
- **While the ledger exists, the top-level `corrections` shall carry `by_phase`, keyed by
  each phase as the ledger spells it, and `outside_changes`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCorrectionsByPhase
- **If no session transcripts are configured or found for the repository, then the
  report shall omit `tokens`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsTokensWithoutTranscripts
- **When transcripts are found, `tokens` shall carry `attributed`, `unattributed`,
  `entries` and `unattributed_entries` computed over every change git has seen, whatever
  the filter.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONTokensAreCorpusWide
- **Where a `<change>` filter names a change sessions were attributed to, `tokens.change`
  shall carry its `name`, `totals` and `by_phase`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONTokensForOneChange
- **If the filtered change has no attributed session, then `tokens.change` shall carry its
  `name` alone.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnattributedChangeCarriesItsNameAlone
- **Where an attributed entry of the filtered change names no phase, `tokens.change`
  shall carry those entries' sum as `no_phase`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONNoPhaseIsNeverDistributed
- **If no run-log file exists, then the report shall omit `loop`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONOmitsLoopWithoutARunLog
- **If the run-log file exists and holds no loop record, then `loop` shall be
  `{changes: [], skipped: N}`, with N the count of its malformed lines.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONAnEmptyRunLogIsKnown
- **When the run log exists, `loop.changes` shall list one object per change it records,
  following the filter, in name order, carrying `name`, `rounds`, `net_closed` and
  `unpriced`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONLoopRows
- **If none of a change's recorded rounds carried a cost, then its loop object shall omit
  `cost_usd`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnpricedOmitsCost
- **When the run log exists, `loop` shall carry `skipped`, the count of lines that are not
  loop records, zero included.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONCountsSkippedLines
- **If the run log exists and cannot be read, then `loop` shall carry `error` in place of
  `changes` and `skipped`.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONUnreadableRunLogNamesTheError
- **While git has seen no change, `changes` shall be an empty array.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONEmptyHistoryIsAnEmptyArray
- **While git has seen no change, the JSON report shall still carry each of
  `corrections`, `tokens` and `loop` whose source exists.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONEmptyHistoryKeepsTheOtherBlocks
- **While git has seen no change, the text report shall print only its no-changes line.**
  Proof: cmd/libretto/metrics_test.go TestMetricsTextWithNoHistoryIsUnchanged
- **When `metrics` is given `--json` and a change name, the command shall accept them in
  either order.**
  Proof: cmd/libretto/metrics_test.go TestMetricsJSONFlagAndChangeInEitherOrder

`queue`

- **When `queue --json` runs, `queued` shall list one object per queued change, carrying
  `name`, `date` and `reading`, in the text's order.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONListsOldestFirst
- **If a queued proposal has no `## Reading`, then its object shall omit `reading`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONOmitsAMissingReading
- **A `reading` carrying `<`, `>` or `&` shall reach the JSON unescaped.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONKeepsReadingTextVerbatim
- **When git lists branches, `dispatched` shall list each change a branch took, carrying
  `name` and `branch`, as an empty array when there is none.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONNamesDispatchedChanges
- **If git cannot list branches, then the report shall omit `dispatched`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONUnreadableBranchesOmitsDispatched
- **While no `.agents/changes/` exists at the repository root, `queued` shall be an empty
  array.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONWithoutAChangesDirIsEmpty
- **While no `.agents/changes/` exists at the repository root, the report shall omit
  `dispatched`.**
  Proof: cmd/libretto/queue_test.go TestQueueJSONWithoutAChangesDirIsEmpty
- **If `queue` is given any argument other than `--json`, then it shall exit non-zero
  naming the argument.** This is the cli spec's existing refusal criterion, amended in
  place at landing; its test gains a row accepting `--json`.
  Proof: cmd/libretto/queue_test.go TestQueueRefusesArguments

## Approach

### Summary

Each of the four commands is split into a gather step that returns one report struct and
two renderers that read it: the existing text, moved without changing a byte, and a JSON
encoder. The JSON structs are new, small, exported-field types with explicit tags, built
from the report — never tags on `link.Entry`, `changeMetrics` or `loopTally`. One shared
helper finds `--json` among a command's arguments and owns the encoder settings. Four
boxes, one per command; the first carries the helper.

### Technical context

- Go 1.26.5, stdlib only: `encoding/json` is already imported by `metrics.go`, `loop.go`
  and `usage.go`.
- Tests live in `cmd/libretto/*_test.go`, package `main`, standard library only.
  `capture` in `helpers_test.go` swaps `os.Stdout`/`os.Stderr` for `status` and `doctor`;
  `metrics` and `queue` already write to an `io.Writer`. `TestMain` in `notice_test.go`
  stubs `noticeSource` for the package.
- Gates: the six in `AGENTS.md`. `spec-drift --anchors` resolves every citation in the
  criteria above down to the test name, so each box's tests must exist under exactly the
  names cited.
- Blast radius: `cmd/libretto/main.go` (`run`'s dispatch passes `args[1:]` to `status`
  and `doctor`; `status`, `doctor`, `usage`), `cmd/libretto/version.go` (`releaseLine`),
  `cmd/libretto/metrics.go` (`metricsWith`, `phaseBreakdown`, `loopBlock`, `tokenBlock`),
  `cmd/libretto/queue.go` (`queue`), one new `cmd/libretto/json.go`, and tests in
  `json_test.go` (new), `doctor_test.go`, `metrics_test.go`, `queue_test.go`,
  `main_test.go`, `notice_test.go`. No `internal/` package changes.
- At landing (phase 8): the cli spec's Task breakdown box `5.2` ticks and gains `metrics`
  and `queue`; the queue non-goal loses "no `--json`"; in the same edit, the queue
  criterion "If the command is given any argument, then it shall exit non-zero naming the
  argument, and shall print no queue." is replaced in place by "If the command is given
  any argument other than `--json`, then it shall exit non-zero naming the argument, and
  shall print no queue.", keeping its `TestQueueRefusesArguments` citation; the criteria above land under a new
  `### --json — the read-only reports, machine-readable` section; `docs/PLAN.md` 5.2 ticks;
  the `docs/STATE.md` row for `5.2` goes.

### The approach

1. **`json.go`** — `jsonFlag(args []string) (bool, []string)` returns whether `--json`
   was present and the remaining arguments; `writeJSON(w io.Writer, v any) error` builds a
   `json.Encoder` with `SetIndent("", "  ")` and `SetEscapeHTML(false)` and encodes once.
   Each document struct carries `Schema int` with tag `json:"schema"`, set to the constant
   `jsonSchema = 1`.
2. **Unknown versus zero, in types.** A key that may be unknown is a pointer or slice with
   `omitempty`; a key that is always known is a plain value without it. A known zero that
   must stay present is a non-pointer `int` with no `omitempty`, or a pointer that is set.
   `omitempty` on a plain `int` is the bug this rule exists to catch: it turns a known `0`
   into an absent key. A known-empty list is initialised to a non-nil empty slice so it
   encodes as `[]`, never `null`.
3. **`status`** — `run` passes `args[1:]`; `status` calls `jsonFlag` and ignores
   whatever is left over, as it ignores every argument today, then `link.Scan` once. Text: unchanged lines. JSON: `statusDoc{Schema,
   Target, Items}` built from the same `entries`; `Target.Root` empty means omitted.
4. **`doctor`** — `releaseLine` becomes two functions: `checkRelease(running, ask)
   releaseCheck` returning `{running, latest string; err error}`, and `releaseLine(c
   releaseCheck) string` formatting it with today's exact sentences. The ask is chosen by a
   package variable `releaseSource func(root string) func(context.Context) (string,
   error)` — the `noticeSource` pattern — whose default is today's choice between
   `repo.Shell.LatestTag` and `dist.Latest`, and which `TestMain` stubs. `newer` is
   `repo.IsNewer(latest, running)` when true, `false` when `running == latest`, absent
   otherwise — the same switch `releaseLine` uses. The scan, the release check and
   `prerequisites()` are gathered once; text and JSON render from them; the returned error
   for problems is unchanged and comes after the JSON is written. Arguments other than
   `--json` are ignored, as today. `TestReleaseLineIsRenderedFromTheCheck` sets one stub —
   `v9.9.9` against a running `v0.1.0`, then an unrankable pair — and asserts the text line
   and the JSON `release` object from that same stub.
5. **`queue`** — gather `kept`, `dropped` and `branchesRead bool`, false both when git
   failed and when no `.agents/changes/` exists, because the command returns before
   asking; text renders as today; JSON maps `queuedIdea` to `{name, date, reading}` and
   `dispatchedIdea` to `{name, branch}`, with `dispatched` nil (omitted) whenever
   `branchesRead` is false. The argument refusal accepts `--json` and still refuses
   anything else; `TestQueueRefusesArguments` gains the `--json` row.
6. **`metrics`** — `metricsWith` grows a gather step, `gatherMetrics(...) (metricsReport,
   error)`, holding: the rows (`changeMetrics` or an unreadable marker), the ledger
   (`corrections()` results and `seen`), the token figures (`readUsage`, `attribute`, the
   per-phase split now inside `tokenBlock`), and the loop tallies (the parsing now inside
   `loopBlock`, plus `skipped`, plus whether the file exists, plus an open error). The
   gather runs the ledger, token and loop steps even when git has seen no change; the text
   renderer checks for no rows first and prints only its no-changes line, exactly as now,
   while the JSON carries `changes: []` and whichever blocks exist. `loop` is omitted only
   when the file does not exist. `phaseBreakdown`, `tokenBlock` and
   `loopBlock` become renderers over the report's parts and print exactly what they print
   now. The JSON doc maps from the same parts. `--json` is taken out of the arguments
   before the existing unknown-flag check, so `--json` and `<change>` work in either
   order and any other flag is still refused.
7. **Help** — one usage line: `--json` with `status`, `doctor`, `metrics`, `queue`: one
   JSON object on stdout.
8. **Tests** — every JSON test decodes stdout with `json.Unmarshal` into
   `map[string]any` (so an absent key is observable, which a typed struct would hide) and
   asserts keys and values. The shared tests in `json_test.go` are tables with one row per
   command; each box adds its row.

### The alternatives it beat

| Considered | Why it lost |
|---|---|
| B — a separate JSON function per command, deriving from the same inputs | a second derivation is a second source of truth; the first fix to one side drifts the other, and nothing would notice |
| C — `json` tags on the internal types (`link.Entry`, `changeMetrics`, `loopTally`, `queuedIdea`) | their fields are unexported, they carry bytes and paths no consumer needs, and every refactor of an internal type would become a contract change |
| JSON errors on stderr | a second error format beside `main`'s line, for a consumer that already has the exit code |
| No `schema` key, keyed on `libretto version` | in `0.x` the binary's version moves on every merge and says nothing about the shape |
| Golden files, as `docs/PLAN.md` 5.2 says | there is no `testdata/` here, and a golden file pins indentation this contract does not promise |
| `null` for unknown values | the loop run log already set absent-not-zero, and two conventions in one binary is one too many |

### Risks

| Risk | What catches it |
|---|---|
| the refactor changes a byte of the text output | the existing text tests in `metrics_test.go`, `queue_test.go` and `main_test.go` stay untouched and green; a box that edits one to pass is a stop, per `evidence` |
| `omitempty` on a plain `int` drops a known zero | `TestMetricsJSONKeepsAKnownZero` and `TestMetricsJSONCorrectionsZeroIsPresent` assert the zero is present, decoding into `map[string]any` |
| a nil slice encodes as `null` where `[]` is promised | `TestStatusJSONEmptyIsAnEmptyArray`, `TestMetricsJSONEmptyHistoryIsAnEmptyArray`, `TestQueueJSONWithoutAChangesDirIsEmpty` |
| `doctor` tests still reach the proxy | `releaseSource` stubbed in `TestMain`; a test wanting an answer sets its own |
| the notice lands on stdout under `--json` | `TestJSONStdoutParsesWithANoticePending`, with `stubNotice` set |
| a later block in `metrics` (the queued drift trend) is added to the text and forgotten in JSON | not caught by a gate. The approach makes the report struct the only place a block's data lives, so a block rendered as text exists in the report; whether it is mapped to JSON is review's to see |

### Validation and rollback

Each box's named tests, then the six gates. Three tests would pass for the wrong reason
and are forced red once before they are believed: `TestJSONHasNoEscapeCodes` (no command
colours today — inject an escape into one rendered string, watch it fail, revert);
`TestMetricsJSONKeepsAKnownZero` (add `omitempty` to `closed`, watch it fail, revert);
`TestJSONStdoutParsesWithANoticePending` (point the notice at stdout, watch it fail,
revert). Rollback is one revert per box; nothing is written to disk or migrated.

### Complexity deliberately kept

- **A separate JSON struct per command, beside the report struct.** Marshalling the report
  directly would be shorter, and it would make every field of an internal type a promise.
  The mapping is the contract's boundary.
- **The `releaseSource` seam.** One package variable only tests use; it exists because
  `doctor` reaches the network from `go test` today.
- **`map[string]any` in the tests** rather than the doc structs. A typed decode cannot see
  whether a key was absent or zero, and that distinction is half the contract.

## Checklist

Executed by `build-and-check`. Box 1 comes first; boxes 2, 3 and 4 each wait only on it.

- [x] **`status --json`, end to end, with the shared helper.** `json.go` with `jsonFlag`
      and `writeJSON`; `run` passes `args[1:]` to `status`; `status` reads `--json` and
      ignores other arguments as today, gathers once, renders text unchanged or the
      `status` shape; the help line;
      `json_test.go` with the shared tables holding the `status` row. Closes:
      TestJSONIsOneObjectOnStdout, TestJSONCarriesSchemaOne,
      TestJSONDoesNotChangeTheExitCode, TestJSONFailureWritesNothingToStdout,
      TestJSONErrorsMatchTheTextErrors, TestJSONHasNoEscapeCodes (each for its `status`
      row), TestJSONStdoutParsesWithANoticePending, TestHelpNamesJSON,
      TestStatusJSONNamesTheTarget, TestStatusJSONOmitsAnUnknownRoot,
      TestStatusJSONListsEveryItem, TestStatusJSONStatesMatchTheText,
      TestStatusJSONEmptyIsAnEmptyArray. Waits on: nothing.
- [x] **`doctor --json`, end to end.** `checkRelease` and `releaseLine(releaseCheck)` in
      `version.go`, existing doctor tests adjusted only in how they call it, never in what
      they assert; `releaseSource` stubbed in `TestMain`; `doctor` reads `--json` and
      ignores other arguments as today, gathers once, renders text unchanged or the
      `doctor` shape; `TestReleaseLineIsRenderedFromTheCheck` asserting text and JSON from
      one stub; the `doctor` rows of the shared tables. Closes: TestDoctorJSONListsOnlyProblems, TestDoctorJSONNamesTheMode,
      TestDoctorJSONCarriesTheRunningVersion, TestDoctorJSONCarriesTheLatestRelease,
      TestDoctorJSONRanksWhatRanks, TestDoctorJSONOmitsAnUnrankableComparison,
      TestDoctorJSONReportsAFailedCheck, TestDoctorJSONListsPrerequisites,
      TestReleaseLineIsRenderedFromTheCheck, and the shared criteria's `doctor` rows.
      Waits on: box 1.
- [ ] **`queue --json`, end to end.** `queue` accepts `--json` and still refuses any other
      argument, with `TestQueueRefusesArguments` gaining the `--json` row; kept, dropped
      and whether branches were asked and read gathered once; text unchanged or the
      `queue` shape; the `queue` rows of the shared tables. Closes:
      TestQueueJSONListsOldestFirst, TestQueueJSONOmitsAMissingReading,
      TestQueueJSONKeepsReadingTextVerbatim, TestQueueJSONNamesDispatchedChanges,
      TestQueueJSONUnreadableBranchesOmitsDispatched,
      TestQueueJSONWithoutAChangesDirIsEmpty (both the empty `queued` and the absent
      `dispatched`), the amended refusal through TestQueueRefusesArguments, and the shared
      criteria's `queue` rows.
      Waits on: box 1.
- [ ] **`metrics --json`, end to end.** `gatherMetrics` returning one report holding rows,
      ledger, tokens and loop tallies — computed even with no change in history, and with
      whether the run-log file exists; `phaseBreakdown`, `tokenBlock` and `loopBlock`
      rendering from it byte-identically, the no-history text still one line; `--json` taken out before the unknown-flag
      check; the `metrics` shape; the `metrics` rows of the shared tables. Closes:
      TestMetricsJSONListsTheTableRows, TestMetricsJSONAgreesWithTheTable,
      TestMetricsJSONCarriesTimesInUTC, TestMetricsJSONKeepsAKnownZero,
      TestMetricsJSONOmitsChurnWithoutAChecklist,
      TestMetricsJSONUnreadableRowCarriesOnlyWhatIsKnown,
      TestMetricsJSONCorrectionsZeroIsPresent, TestMetricsJSONAbsentLedgerOmitsCorrections,
      TestMetricsJSONCorrectionsByPhase, TestMetricsJSONOmitsTokensWithoutTranscripts,
      TestMetricsJSONTokensAreCorpusWide, TestMetricsJSONTokensForOneChange,
      TestMetricsJSONUnattributedChangeCarriesItsNameAlone,
      TestMetricsJSONNoPhaseIsNeverDistributed, TestMetricsJSONOmitsLoopWithoutARunLog,
      TestMetricsJSONAnEmptyRunLogIsKnown, TestMetricsJSONLoopRows,
      TestMetricsJSONUnpricedOmitsCost, TestMetricsJSONCountsSkippedLines,
      TestMetricsJSONUnreadableRunLogNamesTheError,
      TestMetricsJSONEmptyHistoryIsAnEmptyArray,
      TestMetricsJSONEmptyHistoryKeepsTheOtherBlocks,
      TestMetricsTextWithNoHistoryIsUnchanged, TestMetricsJSONFlagAndChangeInEitherOrder,
      and the shared criteria's `metrics` rows. Waits on: box 1.
