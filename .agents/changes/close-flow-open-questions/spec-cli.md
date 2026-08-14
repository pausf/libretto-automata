# close-flow-open-questions — cli delta

Targets: cli

`libretto metrics` learns to say *where* corrections surface, not only how many.
The instrument the ask-more-questions debate gets decided with: corrections at
phase 2 are questions working, findings at the 6→7 seam are questions that were
missing, and corrections after phase 8 are the expensive kind.

## Outcomes

1. The corpus-wide report gains a per-phase breakdown of ledger entries: one row per
   distinct phase value found in `.agents/lessons.md`, with its count. Phases are
   reported as the ledger spells them — `phase 2`, `6→7`, `phase 8` — never
   normalised, because the parser's contract is three plain fields and a normaliser
   is a second spelling that drifts.
2. An absent ledger keeps its existing meaning: the breakdown says capture is not in
   use, never a row of zeros. A ledger that exists with no parseable entries renders
   no phase rows and says so in one line — present-and-empty is a third fact, not
   either of the other two.
3. Orphan entries (change field `-`) keep being counted apart, and appear in the
   breakdown by their phase like any other entry — an orphan still surfaced
   somewhere.
4. Entries whose phase is `6→7` are reviewer findings by the payload delta's
   contract: they appear in the phase breakdown and are **excluded from the
   per-change corrections column**, whose meaning stays user corrections only.

## Scope boundaries

In: the breakdown above, in the corpus-wide report. Out, named:

- **No per-change phase breakdown.** The per-change table keeps its single
  corrections column; slicing it by phase multiplies columns for a question nobody
  asked yet. Add it when a real reading needs it.
- **No transcript reading for this.** Token cost reads transcripts; correction
  attribution reads only the ledger. Two sources for one column is drift.
- **No verdicts.** The report counts; the reading — "too many late findings, ask
  more at phase 2" — stays human.

## Constraints

- `corrections()` already parses the header; the breakdown extends what it returns
  rather than re-reading the file.
- Read-only, like all of metrics: derived from the ledger and git, instrumented
  nowhere.
- The report explains its own columns, as the existing sections do.

## Prior decisions

- Reviewer findings reach the ledger via review-work (payload delta), so this
  breakdown sees all three sources through one file. (2026-08-14, AskUserQuestion)
- Counts, not judgments: metrics stays a report. (standing, `libretto metrics`)

## Task breakdown

1. Extend `corrections()` to also return per-phase counts, with tests.
2. Render the breakdown in the corpus-wide report with its explanation line, with
   golden/behaviour tests.

## Verification criteria

- Per-phase counts come out of the ledger exactly as spelled, bad lines skipped.
  Proof: cmd/libretto/metrics_test.go TestCorrectionsCountsByPhase
- The corpus report renders one row per phase with its count and an explanation.
  Proof: cmd/libretto/metrics_test.go TestMetricsReportsCorrectionsByPhase
- An absent ledger renders the not-in-use line, no zero rows; a present-but-empty
  ledger renders the one-line empty statement, no phase rows.
  Proof: cmd/libretto/metrics_test.go TestMetricsPhaseBreakdownAbsentLedger
- `6→7` entries count in the phase breakdown and never in the per-change
  corrections column.
  Proof: cmd/libretto/metrics_test.go TestReviewerFindingsStayOutOfCorrections
