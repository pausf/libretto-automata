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
lands in a single **unattributed** bucket that is printed, never discarded. **Measured on
this repository once it ran: 62% of entries attribute to no change.** A total that
silently omitted that is worse than no total, and 62% is also why the bucket is a
headline number rather than a footnote.

**3 · Three states are distinguishable, and each prints differently.**

This replaces an earlier sentence claiming the corrections column's rule "applied
unchanged". It does not apply unchanged, and the difference matters: `corrCell` dashes on
the **ledger being absent**, which is one global state, and prints `0` for a change the
ledger simply does not mention. Token cost has three states where corrections have two,
and collapsing them would hide the very miss rate this change exists to expose.

| State | Prints |
|---|---|
| no transcript root, or none for this repository | the whole token block is replaced by one line saying the measurement was unavailable |
| a root exists, and nothing attributed to this change | **a dash** — sessions were read and none reached it |
| attributed, and the numbers are genuinely zero | **`0`** |

The third is **not** unreachable, which an earlier draft of this spec claimed. A
`<synthetic>` entry carries an all-zero usage object and attributes like any other, so a
change whose only entry is one renders zeros. The distinction is therefore **presence in
the attribution map, never a zero total** — and the per-phase block is shown for such a
change too, because its entries named phases even though they cost nothing.

**4 · Cost is attributed per phase, from `attributionSkill`, at the coverage it actually
has.**

Entries carry `attributionSkill`, and its values are the flow's phases —
`build-and-check`, `record-work`, `write-spec`, `find-work`, `review-work`. Measured
coverage is **4,728 entries carrying it against 4,216 without** in the main files, so the
per-phase block carries an explicit unattributed row and never distributes the remainder.

**5 · Subagent transcripts are counted.**

They live at `<project>/<sessionId>/subagents/agent-*.jsonl`, carry `isSidechain: true`,
and carry their own `gitBranch`. One review-lens subagent file alone holds 3,026,036
cache-read tokens; there are 41 of them. **Ignoring them undercounts by more than the
number being reported.**

**6 · The ceiling gains what is true, and nothing in it is retracted.**

An earlier draft of this outcome said `flowCeiling` had become half-false. **It had not,
and the build is what established that.** The ceiling names per-phase *duration* and
`review-work` findings; both still need a phase to write them down, and both are still
unmeasurable. Cost was never on that list — it simply was not measured at all.

So the ceiling is extended rather than corrected, with three things it did not say:

- per-phase **cost** is measured now, because the transcripts record `attributionSkill`
- **duration stays off**, and now with its own reason: the same entries carry timestamps,
  but a phase's wall clock includes every wait for a human and would report attention the
  work never had
- **the token block's own ceiling** — a change is recognised by its branch name, so the
  unattributed row is the measurement's error bar and is to be read before the rest

## What each command prints

**Pinned here rather than described, because "the single-change view" named a surface that
does not exist.** `metrics <change>` today is the same table filtered to one row, not a
different rendering. Both commands gain a **token block**, printed after the existing
footer and before `flowLegend`.

`libretto metrics` — corpus-wide, one block:

```
  tokens                    input         output        cache-w          cache-r
  attributed                54784        2866759       10467947        654722879
  unattributed              84493        5029953       19145370       1736078669

  62% of session entries could not be attributed to a change
```

**Plain digits, no thousands separator and no exponent form.** An earlier draft of this
block pinned `21 480` and `1.42e+09`; both were rejected during the build. A space inside
a number breaks `awk`, `cut` and every other thing a person pipes a report into, and an
exponent throws away the precision that makes two runs comparable — which is the entire
use this measurement was queued for.

**And the miss rate never rounds itself away.** Integer division truncates, so one
unattributed entry in two hundred reads as `0%`. Anything above zero and below one
percent prints `<1%`: a number the ceiling calls an error bar must not report the
opposite of its own subject.

`libretto metrics <change>` — the same block scoped to that change, then the phases:

```
  by change                 input         output        cache-w          cache-r
  add-flow-retro              691         311444         757004         47033497

  by phase                  input         output        cache-w          cache-r
  build-and-check             244          93518         151654         20572299
  find-work                    68          18920          37962          3385831
  record-work                  50          25223          63643          5616550
  unattributed                  —              —              —                —
```

The corpus block prints under both commands, above this. That is what makes
*attributed + unattributed = corpus* readable whatever the filter, and it is why the
per-change rows are labelled `by change` rather than repeating the word `tokens`.

**The dash lives in that per-change block and nowhere else.** Under plain `libretto
metrics` there is no per-change token surface, so a change with nothing attributable is
not represented there at all — which is the direct consequence of taking no new table
column, and it is stated rather than discovered. Whoever wants to find the expensive
change without naming it first is asking for the column, and the boundary below says what
brings it back.

**The token footer is corpus-wide under both commands, and the per-change row is
filtered.** The two are labelled differently for exactly this reason: the invariant
*attributed + unattributed = corpus* has to be readable off the output, and it cannot be
if the totals move when a filter is applied.

## Scope boundaries

**In:**

- `cmd/libretto/usage.go` — transcript discovery, streaming parse, attribution
- `cmd/libretto/usage_test.go` — its proofs, against fixtures under `CLAUDE_HOME`
- `cmd/libretto/metrics.go` — the token block, the per-phase block, `flowCeiling`
- `cmd/libretto/metrics_test.go` — the rendering proofs
- `cmd/libretto/main.go` — resolving the transcript root and passing it in
- `.agents/specs/cli/spec.md` — the criteria, applied at phase 8

**Out, and named so it cannot be quietly added:**

- **`internal/target/claude.go` is not touched.** It was in scope until review asked who
  governs it: `targets`, not `cli` — so touching it would make this change span two
  capabilities with one delta declared. Resolved by reading rather than by declaring a
  second delta: `Root() string` is already exported, so `filepath.Join(claude.Root(),
  "projects")` in `main.go` needs no new accessor. **This is why `Targets: cli` is
  correct**, and it would not be if that file moved.
- **No new column in the change table.** It is already seven columns wide, and the tests
  pick the corrections cell by counting from the right (`corrField`), so a column on that
  side breaks a proof for a cosmetic reason. *Brings it back:* wanting to rank changes by
  cost without naming each one — at which point it goes **left of `corr`**, and it has to
  pick one of four numbers or invent a composite, which is the decision that was avoided
  here rather than made.
- **No pricing, no currency, no model price table.** Prices change under the repository
  and a stale table reports confident nonsense. Tokens are what was measured; money is a
  multiplication the reader can do with today's numbers.
- **No new dependency.** `encoding/json` from the standard library, as `internal/dist`
  already does. Every direct dependency in `go.mod` today is a Charm package, and adding a
  JSON library to parse JSON the standard library parses is the ladder rung `AGENTS.md`
  puts behind an ask.
- **No new package.** One file beside `metrics.go`, its only consumer. An `internal/`
  package for one caller is an abstraction with one implementation.
- **No instrumentation, no new artifact, nothing written under the transcript root.** The
  no-instrumentation criterion in `cli` is the bar: this command free-rides on a file
  somebody else already writes, exactly as the corrections count free-rides on the lessons
  ledger.
- **No decoding a project directory name back into a path.** The encoding replaces every
  `/` with `-` and is lossy: `-Users-x-gitrepos-promofarma-v3` could decode to two
  different paths. Encode the repository root forward and compare; never invert.
- **No git-side join to rescue the 14 unmatched changes.** Matching a change folder's
  commits to the branches containing them would raise the hit rate and is substantially
  larger work resting on history that gets squashed. *Brings it back:* the unattributed
  percentage — now printed — staying high enough to make the attributed numbers
  untrustworthy. The decision will be made on a number rather than a feeling.
- **No per-phase duration.** Still genuinely unmeasurable: phases 1 to 7 happen in one
  session and leave one commit. Only the *cost* claim is corrected.

## Constraints

- **`CLAUDE_HOME` in every test that touches a transcript root, pointed at `t.TempDir()`
  with fixture files written into it.** `AGENTS.md` forbids writing to a real `~/.claude`
  in a test; reading one is not forbidden by that sentence but is machine-dependent, which
  makes a test that passes here and fails on any other checkout. The repository already
  has the pattern — `notice_test.go`, `update_release_test.go`.
- **The transcript root is injected, never resolved inside the logic.** `metrics.go`
  already takes its `gitRunner` as a parameter and `internal/dist` takes its path inputs as
  parameters with a thin environment-reading wrapper. Same seam, same reason.
- **Every field is optional.** Whole entry types — `mode`, `last-prompt`, `ai-title`,
  `file-history-delta` — carry no `cwd`, no `gitBranch`, no `timestamp` and no `uuid`. A
  struct that assumes any field crashes on a real file.
- **A malformed line is skipped, never fatal.** An 81 MB corpus written by another program
  across many versions will contain something unparseable, and a metrics command that dies
  on one bad line reports nothing about the other 8,943.
- **Streamed line by line.** Files reach 6.8 MB and the directory is 81 MB.
- **`model` may be `<synthetic>`**, carrying an all-zero usage object with null
  `service_tier`. It must not be looked up in any table assuming a known name.
- **No transcript root, or none for this repository, is a state and not an error.**
- All six gates pass before any commit.

## Prior decisions

**Settled by measurement, this change** — each number came from the recon, not from
reasoning:

- **The four are never summed.** 2.3 billion cache reads against 33 thousand input tokens
  means a total is a cache-read count with extra steps.
- **Subagent files are in scope.** 41 files, one alone carrying 3M cache reads.
- **`gitBranch` is read per entry, not per file.** One session file was measured spanning
  four branches; a per-file reading would misattribute every entry after a checkout.
- **`sessionId` is preferred over `session_id`.** Both appear and are identical on
  `assistant` entries, but `sessionId` is present on entry types where the other is not.

**Settled by reading, after review asked:**

- **`target.Claude.Root()` suffices**, so this change touches only `cmd/libretto/**` and
  `Targets: cli` covers everything in scope.

**Assumed under `/libretto-attacca`, because nobody was there to answer.** Each names what
changes if it is wrong:

- **A1 · Branch-name matching is the attribution mechanism, and its miss rate is printed
  rather than engineered away.** *If wrong:* a git-side join replaces the matcher behind
  the same reporting surface; the unattributed bucket is what the report is built around
  either way, so nothing rendered has to move.
- **A2 · The change table gains no column; the four numbers go in a block after the
  footer.** *If wrong:* one column is additive, must sit left of `corr`, and forces the
  composite-metric decision this deliberately avoided.
- **A3 · Attribution strips a conventional branch prefix and then matches whole names.**
  `feat/`, `fix/`, `docs/`, `chore/`, `refactor/`. *If wrong:* the list is one line, and it
  fails safe — an unrecognised prefix lands in the unattributed bucket rather than
  matching the wrong change.
- **A4 · The corrected ceiling claim gets its own test rather than extending the existing
  one.** *If wrong:* the assertion moves into `TestTheReportNamesWhatItCannotMeasure` and
  the new test goes away; what must not happen either way is a criterion whose only proof
  passes before the work starts.

**Amended after `review-spec`, before phase 5 read this.** Six findings, all acted on:

- the rendering is pinned with worked output for both commands — "the single-change view"
  named a surface that does not exist
- the dash's home is stated, and so is the consequence of having no table column
- the corrections analogy is withdrawn; three states are named where it offered two
- `internal/target/claude.go` left the scope after `Root()` was read and found sufficient
- the ceiling correction got its own test, because the one it named was green before the
  work and stays green after
- the never-written proof names its witness

## Task breakdown

1. **Fixtures and the parse.** A fixture transcript root under `t.TempDir()`, and a
   streaming reader returning the four totals per entry with `gitBranch`,
   `attributionSkill` and `isSidechain`. Test-first: the shapes are known from real files,
   so the test is red before the parser exists.
2. **Malformed and absent.** A bad line, an entry with no `usage`, an entry with no
   `gitBranch`, a `<synthetic>` model, an empty root, a missing root.
3. **Discovery.** Encode the repository root into a project directory name; find the
   top-level `*.jsonl` and the `*/subagents/*.jsonl` beneath it.
4. **Attribution.** Prefix strip, whole-name match, everything else to the unattributed
   bucket.
5. **Rendering.** The corpus block, the per-change block, the per-phase rows, the dash.
6. **`flowCeiling`.** Correct the cost claim, leave the duration claim, new test.
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
  `<sessionId>/subagents/agent-*.jsonl` beneath it reports the sum of both; a reader that
  walked only the top level fails this.
  Proof: cmd/libretto/usage_test.go TestSubagentTranscriptsAreCounted

- **`gitBranch` is read per entry.** A single fixture file whose entries name two
  different branches attributes each entry to its own.
  Proof: cmd/libretto/usage_test.go TestBranchIsReadPerEntryNotPerFile

- **a conventional branch prefix is stripped and the rest is matched whole.**
  `feat/add-thing` attributes to change `add-thing`; `feat/add-thing-extra` does not.
  Proof: cmd/libretto/usage_test.go TestAPrefixIsStrippedAndTheNameMatchedWhole

- **what cannot be attributed is bucketed and printed, never dropped.** Entries on `main`,
  on `HEAD`, and on a branch matching no change all reach the unattributed total, and the
  attributed totals plus the unattributed total equal the corpus total.
  Proof: cmd/libretto/usage_test.go TestUnattributedTokensAreReportedNotDiscarded

- **nothing under the transcript root is written.** The witness is a **snapshot of the
  fixture tree — every path, its size, and the SHA-256 of its contents — taken before the
  command runs and compared after**. That goes red on a create, a delete, a truncation and
  an in-place rewrite alike, which an mtime comparison does not, and it is an observation
  of effect rather than a grep for `os.Create` in the source.
  Proof: cmd/libretto/usage_test.go TestTheTranscriptRootIsNeverWritten

- **a change with a transcript root and nothing attributed prints a dash; a change with no
  root at all is not represented as zero either.** The three states of outcome 3, each
  distinguishable in the output.
  Proof: cmd/libretto/metrics_test.go TestAChangeWithNoTokensReportsADashNotAZero

- **a change that was attributed and genuinely cost nothing prints zeros, not a dash, and
  keeps its per-phase block.** The distinction is presence in the attribution map, never a
  zero total — a `<synthetic>` entry reaches this state, so it is reachable rather than
  theoretical, and the first implementation collapsed it into the dash.
  Proof: cmd/libretto/metrics_test.go TestAnAttributedChangeWithZeroTokensIsNotADash

- **a miss rate above zero never prints as `0%`.** Integer division truncates; one
  unattributed entry in two hundred read as `0%` until this was fixed, which is the
  opposite of what a number the ceiling calls an error bar is for.
  Proof: cmd/libretto/metrics_test.go TestASmallMissRateDoesNotRoundAwayToZero

- **a dot in the repository path encodes to a dash, like a separator.** A home directory
  called `pau.sanchez` maps to `pau-sanchez`, so a reader that preserved dots found no
  transcripts at all for that user — which this did, and only running it against a real
  tree said so. Every fixture used a dotless path and none could have caught it.
  Proof: cmd/libretto/usage_test.go TestADotInThePathBecomesADashToo

- **the per-phase block carries its own unattributed row**, and never distributes entries
  that named no skill across the phases that did.
  Proof: cmd/libretto/metrics_test.go TestPerPhaseCostCarriesAnUnattributedRow

- **the corpus totals do not move when a change filter is applied**, so *attributed +
  unattributed = corpus* is readable off either command's output.
  Proof: cmd/libretto/metrics_test.go TestTheTokenFooterIsCorpusWideUnderAFilter

- **no transcript root is a state, not an error.** The git-derived report still prints in
  full and the token block is replaced by one line saying the measurement was unavailable.
  Proof: cmd/libretto/metrics_test.go TestNoTranscriptRootStillReportsTheGitMetrics

- **the ceiling says per-phase cost is measured and duration is not, and names the token
  block's own limit.** Two presence assertions, both red before this change: neither
  `attributionSkill` nor the unattributed row's description as an error bar appeared
  anywhere in the report.
  An earlier draft demanded an *absence* assertion instead, on the theory that the old
  claim had become false. It had not — the ceiling never spoke about cost — so the
  assertion would have been red for the wrong reason and green for a change that retracted
  a true sentence. The existing `TestTheReportNamesWhatItCannotMeasure` is left exactly as
  it is, still asserting the three things that are still unmeasured.
  Proof: cmd/libretto/metrics_test.go TestTheCeilingSeparatesCostFromDuration
