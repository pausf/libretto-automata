# add-token-cost-to-metrics

Targets: cli

`libretto metrics` derives everything from git and measures no tokens, so any context
diet is guesswork. This teaches it to read the Claude Code session transcripts and report
what a change actually cost — **the four usage numbers kept apart**, because cache reads
bill at a fraction of input and summing them makes cheap cached context look catastrophic.

**Two facts were measured before this was written, and both move the design:**

- **Cache reads are 99% of everything.** Across this repository's 38 session files:
  2,309,422,022 cache-read tokens against 33,052 plain input tokens — a ratio near
  70,000:1. A single "tokens" number is a cache-read number wearing a disguise.
- **The transcript names a branch, and a branch names a change only sometimes.** Every
  entry carries `gitBranch`, but 4,300 entries say `main`, 2,806 say `HEAD`, and one
  change was built on a branch with no string relationship to its name. Matching
  prefix-stripped branch names against the 32 changes git has ever seen: **18 match, 14 do
  not.**

That second number is the whole shape of this change. It is not a defect to engineer
away — it is the honest ceiling of deriving attribution from an artifact nobody wrote for
this purpose, and the report says so rather than implying a precision it does not have.

## Outcomes

**1 · The four usage numbers are read, kept apart, and never summed into one.**

Input, output, cache creation and cache read, from `.message.usage` on entries whose
`type` is `assistant` — the only type that carries usage, measured at 8,944 of 8,944.
`int64` throughout: one repository's cache-read total is 2.3 billion.

`.message.usage.iterations[]` repeats the same four numbers and is **not** added to them.
Adding both double-counts every entry that has it.

**2 · Cost is attributed per change, by branch, and what cannot be attributed says so.**

An entry's `gitBranch`, with a leading `feat/`, `fix/`, `docs/`, `chore/` or `refactor/`
segment removed, matched **whole** against the change names git has seen. Never a prefix
match and never fuzzy: a substring rule turns one wrong guess into a number nobody can
audit.

Everything else — `main`, `HEAD`, a branch named unlike its change, an absent field —
lands in a single **unattributed** bucket that is printed, never discarded. A total that
silently omits 44% of the history is worse than no total.

**A change with no attributable session prints a dash, never a zero.** This is the
existing rule for the corrections column, applied unchanged: absent means the measurement
could not reach it; zero means it ran and cost nothing. Printing `0` for the first claims
the second.

**3 · Cost is attributed per phase, from `attributionSkill`, at the coverage it actually
has.**

Entries carry `attributionSkill`, and its values are the flow's phases —
`build-and-check`, `record-work`, `write-spec`, `find-work`, `review-work`. Measured
coverage is **4,728 entries carrying it against 4,216 without** in the main files, so the
per-phase block carries an explicit unattributed row and never distributes the remainder.

**4 · Subagent transcripts are counted.**

They live at `<project>/<sessionId>/subagents/agent-*.jsonl`, carry `isSidechain: true`,
and carry their own `gitBranch`. One review-lens subagent file alone holds 3,026,036
cache-read tokens; there are 41 of them. **Ignoring them undercounts by more than the
number being reported.**

**5 · The report's claim about what it cannot measure is corrected, not left standing.**

`flowCeiling` and the `cli` criterion behind it say per-phase measurement "needs a phase
to write them down". For **duration** that is still true. For **cost** it is now false —
`attributionSkill` writes it down already. A criterion that stays wrong because it was
convenient is the drift this repository's spec rules exist to prevent.

## Scope boundaries

**In:**

- `cmd/libretto/usage.go` — transcript discovery, streaming parse, attribution
- `cmd/libretto/usage_test.go` — its proofs, against fixtures under `CLAUDE_HOME`
- `cmd/libretto/metrics.go` — the footer totals, the per-change block, `flowCeiling`
- `cmd/libretto/metrics_test.go` — the rendering proofs
- `cmd/libretto/main.go` — resolving the transcript root and passing it in
- `internal/target/claude.go` — an accessor for the `projects/` directory, if `Root()`
  does not already suffice
- `.agents/specs/cli/spec.md` — the criteria, applied at phase 8

**Out, and named so it cannot be quietly added:**

- **No new column in the change table.** It is already seven columns wide; four more
  makes eleven, and the tests pick the corrections cell by counting from the right
  (`corrField`), so a column added on that side breaks a proof for a cosmetic reason. The
  four numbers go in the footer as totals and in the single-change view per phase.
  *Brings it back:* wanting to rank changes by cost without running the command once each.
- **No pricing, no currency, no model price table.** Prices change under the repository
  and a stale table reports confident nonsense. Tokens are what was measured; money is a
  multiplication the reader can do with today's numbers.
- **No new dependency.** `encoding/json` from the standard library, as `internal/dist`
  already does. Every direct dependency in `go.mod` today is a Charm package and adding a
  JSON library to parse JSON the standard library parses is the exact ladder rung
  `AGENTS.md` puts behind an ask.
- **No new package.** One file beside `metrics.go`, its only consumer. An `internal/`
  package for one caller is an abstraction with one implementation.
- **No instrumentation, no new artifact, no writing to the transcripts.** The
  no-instrumentation criterion in `cli` is the bar: this command is a free rider on a file
  somebody else already writes, exactly as the corrections count free-rides on the
  lessons ledger. **Read-only, always** — nothing under `~/.claude` is written, moved or
  deleted.
- **No decoding a project directory name back into a path.** The encoding replaces every
  `/` with `-` and is lossy: `-Users-x-gitrepos-promofarma-v3` could decode to two
  different paths. Encode the repository root forward and compare; never invert.
- **No git-side join to rescue the 14 unmatched changes.** Matching a change folder's
  commits to the branches containing them would raise the hit rate and is a substantially
  larger piece of work resting on history that gets squashed. *Brings it back:* the
  unattributed bucket staying large enough to make the attributed numbers untrustworthy —
  which the report now prints, so the decision will be made on a number rather than a
  feeling.
- **No per-phase duration.** Still genuinely unmeasurable: phases 1 to 7 happen in one
  session and leave one commit. Only the *cost* claim is corrected.

## Constraints

- **`CLAUDE_HOME` in every test that touches a transcript root, pointed at `t.TempDir()`
  with fixture files written into it.** `AGENTS.md` forbids writing to a real `~/.claude`
  in a test; reading one is not forbidden by that sentence but is machine-dependent, which
  makes a test that passes here and fails on any other checkout. The repository already
  has the pattern — `notice_test.go`, `update_release_test.go`.
- **The transcript root is injected, never resolved inside the logic.** `metrics.go`
  already takes its `gitRunner` as a parameter and `internal/dist` takes its three path
  inputs as parameters with a thin environment-reading wrapper. Same seam, same reason.
- **Every field is optional.** Whole entry types — `mode`, `last-prompt`, `ai-title`,
  `file-history-delta` — carry no `cwd`, no `gitBranch`, no `timestamp` and no `uuid`. A
  struct that assumes any field crashes on a real file.
- **A malformed line is skipped, never fatal.** An 81 MB corpus written by another program
  across many versions will contain something unparseable, and a metrics command that dies
  on one bad line reports nothing about the other 8,943.
- **Streamed line by line.** Files reach 6.8 MB and the directory is 81 MB. Reading whole
  files into memory to sum four integers is waste with no upside.
- **`model` may be `<synthetic>`**, carrying an all-zero usage object with null
  `service_tier`. It must not be looked up in any table that assumes a known name.
- **No transcript root, or none for this repository, is a state and not an error.** The
  command still prints everything it derives from git; the token block says the
  measurement was not available. Someone running this on a checkout that never hosted a
  session is the normal case, not a broken one.
- All six gates pass before any commit.

## Prior decisions

**Settled by measurement, this change** — each number came from the recon, not from
reasoning:

- **The four are never summed.** 2.3 billion cache reads against 33 thousand input tokens
  means a total is a cache-read count with extra steps.
- **Subagent files are in scope.** 41 files, one of them alone carrying 3M cache reads.
- **`gitBranch` is read per entry, not per file.** One session file was measured spanning
  four branches; a per-file reading would misattribute every entry after a checkout.
- **`sessionId` is preferred over `session_id`.** Both appear and are identical on
  `assistant` entries, but `sessionId` is present on entry types where the other is not.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · Branch-name matching is the attribution mechanism, and its miss rate is printed
  rather than engineered away.** The alternative — joining a change folder's commits to
  the branches that contain them — is a much larger change resting on history that gets
  squashed on merge. *If wrong:* the join replaces the matcher behind the same reporting
  surface; the unattributed bucket is what the report is built around either way, so
  nothing rendered has to move.
- **A2 · The change table gains no column; the four numbers go in the footer and in the
  single-change per-phase block.** *If wrong:* one column of one derived number is
  additive, but it must go left of `corr` or `corrField` breaks.
- **A3 · Attribution strips a conventional branch prefix and then matches whole names.**
  `feat/`, `fix/`, `docs/`, `chore/`, `refactor/`. *If wrong:* the list is one line, and
  it fails safe — an unrecognised prefix lands in the unattributed bucket rather than
  matching the wrong change.
- **A4 · The `cli` criterion claiming per-phase measurement needs instrumentation is
  amended to hold for duration and not for cost.** *If wrong:* the criterion reverts, but
  it would then be a sentence the code contradicts.

## Task breakdown

1. **Fixtures and the parse.** A fixture transcript root under `t.TempDir()`, and a
   streaming reader returning the four totals per entry with `gitBranch`,
   `attributionSkill` and `isSidechain`. Written test-first: the shapes are known from
   real files, so the test can be red before the parser exists.
2. **Malformed and absent.** A bad line, an entry with no `usage`, an entry with no
   `gitBranch`, a `<synthetic>` model, an empty root, a missing root. Each proved not to
   crash and to be counted where the outcome says.
3. **Discovery.** Encode the repository root into a project directory name, find the
   top-level `*.jsonl` and the `*/subagents/*.jsonl` beneath it.
4. **Attribution.** Prefix strip, whole-name match against git's change names, everything
   else to the unattributed bucket.
5. **Rendering.** The footer totals, the single-change per-phase block, the dash for a
   change with nothing attributable.
6. **`flowCeiling` and the criterion.** Correct the cost claim, leave the duration claim.
7. Six gates, then apply the delta onto `.agents/specs/cli/spec.md` and delete the change
   folder, in the commit with the final code.

## Verification criteria

- **the four usage numbers are read apart and never summed**, from `assistant` entries
  only, with `iterations[]` ignored so no entry is counted twice.
  Proof: cmd/libretto/usage_test.go TestTheFourUsageNumbersAreKeptApart

- **a malformed line, a missing usage object and an absent field are all survivable.** One
  unparseable line does not cost the other entries in the file.
  Proof: cmd/libretto/usage_test.go TestAMalformedLineDoesNotCostTheRestOfTheFile

- **subagent transcripts are counted.** A fixture with a top-level file and a
  `<sessionId>/subagents/agent-*.jsonl` beneath it reports the sum of both, and a reader
  that walked only the top level fails this.
  Proof: cmd/libretto/usage_test.go TestSubagentTranscriptsAreCounted

- **`gitBranch` is read per entry.** A single fixture file whose entries name two
  different branches attributes each entry to its own.
  Proof: cmd/libretto/usage_test.go TestBranchIsReadPerEntryNotPerFile

- **a conventional branch prefix is stripped and the rest is matched whole.**
  `feat/add-thing` attributes to change `add-thing`; `feat/add-thing-extra` does not
  attribute to `add-thing`.
  Proof: cmd/libretto/usage_test.go TestAPrefixIsStrippedAndTheNameMatchedWhole

- **what cannot be attributed is bucketed and printed, never dropped.** Entries on `main`,
  on `HEAD`, and on a branch matching no change all reach the unattributed total, and the
  attributed totals plus the unattributed total equal the corpus total.
  Proof: cmd/libretto/usage_test.go TestUnattributedTokensAreReportedNotDiscarded

- **a change with no attributable session prints a dash, never a zero.** The same rule the
  corrections column already keeps.
  Proof: cmd/libretto/metrics_test.go TestAChangeWithNoTokensReportsADashNotAZero

- **the per-phase block carries its own unattributed row**, and never distributes entries
  that named no skill across the phases that did.
  Proof: cmd/libretto/metrics_test.go TestPerPhaseCostCarriesAnUnattributedRow

- **no transcript root is a state, not an error.** The git-derived report still prints in
  full and the token block says the measurement was unavailable.
  Proof: cmd/libretto/metrics_test.go TestNoTranscriptRootStillReportsTheGitMetrics

- **the report no longer claims per-phase cost is unmeasurable, and still claims it for
  duration.** `flowCeiling` was true when it was written and half of it stopped being
  true when `attributionSkill` started being recorded.
  Proof: cmd/libretto/metrics_test.go TestTheReportNamesWhatItCannotMeasure

- **nothing under the transcript root is written.** The command opens files for reading
  and creates, moves and removes nothing — the free-rider bar the no-instrumentation
  criterion sets.
  Proof: cmd/libretto/usage_test.go TestTheTranscriptRootIsNeverWritten
