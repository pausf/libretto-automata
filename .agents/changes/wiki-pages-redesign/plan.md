# Plan — wiki-pages-redesign

Durable decisions: the ones in Prior decisions of spec.md

## Summary

Restructure `renderWikiHTML` into home-then-articles, add `renderCards` (the
mini-bar reuses the chart's integer arithmetic), add the `wikiGitDate` seam and
thread its results through `wikiCapability`, replace the router/search script,
delete `renderChart` and the spy with their tests, and rewrite the two
superseded capability criteria in the same commit their tests disappear.

## Technical context

- Branch `feat/add-wiki-html-output` (PR #67). Blast radius:
  `cmd/libretto/wiki.go`, `cmd/libretto/wiki_test.go`,
  `.agents/specs/cli/spec.md` (the superseded criteria rewritten early — the
  same `--anchors`-forces-it precedent as the row criterion last change). Three
  files.
- The date is data for cards, so `wikiCapability` gains `changed string`,
  filled in `wiki()` after parsing: `wikiGitDate(projectDir, specPath)`.
- Gates: the six; the four new citations and the two removals land in one box.

## The approach

1. **Seam**: `var wikiGitDate = func(projectDir, specPath string) string` —
   `git -C projectDir log -1 --format=%as -- specPath`, trimmed; any error or
   empty output → `""`. Tests inject fixed dates or emptiness; the default is
   exercised only by the hand look.
2. **Markup**: `<section id="home">` first — h1, two inline totals, search
   input, `<div class="cards">` of `<a class="card" href="#name">` (name, count,
   `<span class="bar"><i style="width:N%">`, date when present) — then one
   `<article class="cap" id="name">` per capability: `<a href="#home">← home</a>`,
   h2, date line, governs, intro, criteria. `.cap` keeps its class so the landed
   motion criterion and test hold unchanged.
3. **Router**: inline script on `hashchange` + load — `location.hash` names an
   article id → `body` gets class `paged` and only that article shows; anything
   else shows home. All visibility via CSS classes; no JS → no `paged` class →
   everything visible stacked.
4. **Search**: filters cards by name + that capability's criteria text (a
   `data-crit` attribute on each card carries the lowercased criteria, built at
   render time — no cross-page DOM reads).
5. **Deletions**: `renderChart`, the spy block, `TestWikiHTMLCarriesTheBentoHeroAndChart`,
   `TestWikiHTMLScrollSpyIsAnEnhancement`; the cli spec's chart and spy criteria
   rewritten to the card/router criteria in the same commit.

## The alternatives it beat

| Considered | Why it lost |
|---|---|
| Multi-file output (index + page per spec) | N+1 generated files to own, mark, refresh and prune; the recorded return path if hashes ever fall short. |
| CSS-only `:target` routing | No JS at all, but no home-by-default: with no hash nothing is targeted, so home needs JS anyway or everything shows; the hybrid (JS toggling, CSS classes) keeps the no-JS story *better* — complete page — and the JS story exact. |
| File mtime for the date | Unstable across checkouts and clones — every fresh clone resets it; git history is the only date that survives travel. |
| `%aI` full timestamps | The card answers "how fresh", not "what minute"; `%as` (YYYY-MM-DD) is the whole need. |

## Risks

| Risk | What catches it |
|---|---|
| The git exec makes output non-deterministic across runs | Same history → same date; the determinism test runs with the real seam on a git-less TempDir (dates absent both runs) and `TestWikiDatesComeFromGitAndDegrade` pins injected-date stability. |
| Deleting the spy/chart tests leaves capability citations dead mid-branch | Rewritten in the same commit — the gate enforces it, the plan schedules it. |
| The `data-crit` attribute reopens escaping | It is attribute context: built from `html.EscapeString`-ed text like everything else; the escaping test's fixture criteria flow into it. Force-red target: `TestWikiDatesComeFromGitAndDegrade` by making the seam error surface as an error, watch red, restore. |

## Validation and rollback

`go test ./cmd/libretto/` — the suite lands at twenty-five wiki-named tests
(twenty-three today, two removed, four added; the cutter corrected both of this
plan's counts, which had been written from memory instead of `rg -c`). Force red on purpose: the degrade arm of
`TestWikiDatesComeFromGitAndDegrade` (make a seam failure fatal, watch it fail).
Render-and-look with the real seam against this repo — real dates on real cards,
opened in the browser. Rollback: one revert.

## Complexity deliberately kept

The `data-crit` attribute duplicates criteria text into the card. Kept: the
alternative is the search script walking hidden articles per keystroke, which
couples the filter to the router's visibility state — data where the query runs
is smaller than logic crossing pages.
