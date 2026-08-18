# Tasks — wiki-pages-redesign

Execution: `build-and-check` (phase 6). One writer owns this file; builders report,
the orchestrator marks.

- [x] **Home-and-pages redesign, one commit** — restructure `renderWikiHTML` in
  `cmd/libretto/wiki.go` into home-then-articles: the `wikiGitDate` seam
  (`var wikiGitDate = func(projectDir, specPath string) string`, `git -C <dir>
  log -1 --format=%as -- <spec>`, any failure → `""`), `changed` threaded through
  `wikiCapability` from `wiki()`; `renderCards` with mini-bars reusing the
  chart's integer arithmetic; `<section id="home">` first then one
  `<article class="cap" id="name">` per capability (`.cap` class kept so the
  landed motion proof holds); inline `hashchange` router toggling a `paged`
  body class, CSS-only visibility, plain anchors throughout; home search over
  card name + render-time `data-crit` (escaped like all other output). Delete
  `renderChart`, the scroll-spy block, `TestWikiHTMLCarriesTheBentoHeroAndChart`
  and `TestWikiHTMLScrollSpyIsAnEnhancement`; rewrite the cli capability spec's
  chart and spy criteria to the new card/router criteria in this same commit.
  - Traces: spec.md Task breakdown 1 and all five verification criteria;
    plan.md approach 1–5.
  - Closes when: `TestWikiHTMLIsHomeAndPages`, `TestWikiDatesComeFromGitAndDegrade`,
    `TestWikiHTMLRouterIsAnEnhancement`, `TestWikiHTMLHomeSearchIsInline` exist
    and pass; `TestWikiHTMLWritesTheViewer` and every other non-superseded wiki
    proof green unchanged (determinism proof on a git-less TempDir, dates absent
    both runs); the degrade arm of `TestWikiDatesComeFromGitAndDegrade` forced
    red on purpose (seam failure made fatal), observed red, restored; no
    surviving citation of the two removed tests anywhere; render-and-look with
    the real seam against this repo — real dates on real cards, opened in the
    browser, observed; all six gates green on the commit.
  - Waits on: nothing. Can start now.
  - Evidence: six gates green on commit "feat(cli): the wiki becomes home and
    pages"; degrade arm forced red by making an absent date fatal, observed
    failing, restored; rendered via a throwaway git worktree so the real seam
    saw real history — 14 cards with genuine dates (2026-08-10…), 14 articles,
    opened in the browser and republished to the artifact; worktree removed.

The remainder of the capability-spec delta lands at the landing commit as always;
it is not a box.
