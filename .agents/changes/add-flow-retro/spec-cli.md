# add-flow-retro — cli delta

Targets: cli

`libretto metrics` gains the number the user asked for: how many times the AI got
corrected, per change.

## Outcomes

- `libretto metrics` shows a corrections column: for each change, the count of
  `.agents/lessons.md` entries whose header names that change.
- The count is read from the ledger file in the working tree — the ledger exists
  for the retro, and metrics reads it for free. `flowCeiling`'s bar is met: the
  artifact was wanted badly enough that somebody said what for, and it is written
  by one skill, not eight.
- No ledger file: the column reads `-` for every change, same as a change with no
  plan — absent is not zero. A ledger with no entries for a change reads `0`.
- Entries whose change field is `-` (no change open) are counted nowhere in the
  per-change rows; the report names them in one line so they are not silently lost.

## Scope boundaries

**In:** the column, its parsing, its tests.

**Out, named:**

- **parsing anything below the header line.** `Said:`/`Did:`/`Resolved:` are the
  retro's business; metrics counts `## ` headers and reads three fields.
- **resolved-versus-open breakdown.** One count. A second column returns the day
  somebody wants to see the retro's backlog from here.
- **git archaeology on the ledger.** The file is append-only by contract; the
  working tree is the whole history. If entries ever get deleted, that is a payload
  contract breach, not something metrics compensates for.

## Constraints

- The header format is owned by the payload delta: `## <date> · <change> · <phase>`,
  fields separated by ` · `. **A line matches when it starts with `## ` and carries
  exactly two ` · ` separators with three non-empty fields.** The date is not
  validated — a misspelled date is still a countable lesson. Anything else is
  skipped, never a crash.
- Same conventions as every metrics column: dash for not-measurable, derived at
  read time, `flowCeiling` honesty about what is not seen.

## Prior decisions

- **Reading a file the payload writes does not break "derived, never recorded".**
  That rule refuses instrumentation invented for the metric; this artifact exists
  for the retro, and metrics is a free rider. The comment atop `metrics.go` says
  the condition; this change is the condition arriving.

## Task breakdown

- [ ] parse the ledger header, count per change, dash when the file is absent
- [ ] the column in the report, plus the one-line note for `-`-change entries
- [ ] tests, same fixture style as the existing metrics tests

## Verification criteria

- entries are counted per change, read off the fixed header format
  Proof: cmd/libretto/metrics_test.go TestCorrectionsAreCountedPerChange
- no ledger file reports a dash, never a zero
  Proof: cmd/libretto/metrics_test.go TestNoLedgerReportsADashNotAZero
- a malformed header line is skipped, and entries with no change land in the note
  Proof: cmd/libretto/metrics_test.go TestMalformedAndChangelessEntriesDoNotCrashTheCount
