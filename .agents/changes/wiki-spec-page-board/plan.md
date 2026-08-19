# Plan — wiki-spec-page-board

Durable decisions: the ones in Prior decisions of spec.md

## Summary

One seam (`wikiGitDates` — all commit dates for a spec), one pure function
(`monthBuckets`), two parser additions (`decisions []string` and the raw text
kept for the mention scan), and the page render: sparkline, chips from the
criterion data already carried, decisions block, related row.

## Technical context

Branch `feat/add-wiki-html-output` (PR #67). Blast radius: `cmd/libretto/wiki.go`,
`cmd/libretto/wiki_test.go`. Gates: the six; citations red until the one box
lands, as every change on this branch.

## The approach

1. `var wikiGitDates(projectDir, specPath) []string` — `git log --format=%as --
   <spec>`, newest first; failure → nil.
2. `monthBuckets(dates []string) [8]int` — bucket by `YYYY-MM`, anchor the last
   bucket at the newest date's month, walk back seven; pure, clock-free.
3. `parseSpec` grows: capture bullets under a `## Prior decisions` heading into
   `decisions []string` (first line of each bullet), and keep `raw` (the whole
   file) for the mention scan.
4. `relatedTo(c, caps)` — word-boundary regex per other capability name over
   `c.raw`, self excluded, order lexical (caps order), each once.
5. Render on the article: sparkline SVG (bars height `count*18/max`, integer),
   chips `<span class="chip ok|warn">file · test</span>` under each criterion,
   decisions block (≤3 bullets + "· N total"), related row of `#` links.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Anchoring buckets at time.Now | Bytes change at month boundaries with unchanged history; determinism criterion forbids it. |
| Parsing decisions with a markdown library | The section is a heading and bullets; the scanner already walks both. |
| Mentions via [[wiki-links]] syntax | Specs don't use it; inventing syntax to link them is a convention change, not a wiki feature. |

## Risks

| Risk | What catches it |
|---|---|
| A capability named like a common word over-links (e.g. `ci`) | Word-boundary match keeps `ci` from matching inside words, but prose "CI" still links — accepted, recorded assumption; the test pins the boundary behaviour. Force-red target: drop the boundary anchors, watch the related test fail on the substring fixture. |
| Sparkline empty-history div-by-zero | max floors at 1; the absence arm is tested. |

## Validation and rollback

`go test ./cmd/libretto/`; six gates; force red: remove the `\b` anchors in the
mention scan, watch `TestWikiPageLinksRelatedCapabilities` fail on its
substring fixture, restore. Render-and-look at the end. Rollback: one revert.

## Complexity deliberately kept

`raw` kept on the capability struct (whole file in memory per spec) — the
mention scan needs the text the parser otherwise discards, and a specs corpus
is kilobytes.
