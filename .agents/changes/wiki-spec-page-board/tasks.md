# Tasks — wiki-spec-page-board

Execution: `build-and-check` (phase 6). One writer owns this file.

- [x] 1. The page carries its history and its proof — seam, parsing, render, five
  tests, one commit. Traces: all five criteria; plan steps 1–5. wikiGitDates seam;
  pure monthBuckets anchored at the newest date's month; parseSpec grows decisions
  + raw; relatedTo case-insensitive word-boundary (the cutter's gap, recorded);
  sparkline SVG, chip per criterion reusing isEARS/proofResolves, decisions block
  (≤3 + total), Related row. Absence arms throughout. Closes when the five cited
  tests pass, existing suite green, the force-red (\b anchors dropped, related test
  red on its substring fixture, restored) observed, six gates green.
  - Waits on: nothing.
  - Evidence: six gates green on commit "feat(cli): the capability page carries
    its history and its proof"; boundary force-red observed on the substring
    fixture's own message, restored.

- [x] 2. The look — generate against this repo, observe a page (sparkline, chips in
  both colours, decisions, related; absence where data is absent). Fix visual
  defects in place, gates rerun. Observation-only close per the plan.
  - Waits on: 1.
  - Evidence: rendered via throwaway worktree — 14 sparklines, 36 ok / 525 warn
    chips (the pre-EARS corpus, honestly amber), 14 decisions blocks, 14 related
    rows; 274 KB page. No visual defects found.
