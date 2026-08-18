# Tasks — modernize-wiki-design

Execute with `build-and-check` (phase 6), one fresh session per open box.
Branch: `feat/add-wiki-html-output`. Every box lands code and proof in one
commit, tree green, all six gates. The capability spec delta lands separately
in the final commit, as always.

- [x] 1. Rewrite the wiki template — palette, hero bento, SVG chart, motion,
      scroll-spy — and land the four structure tests with it
  - Traces: spec delta *Task breakdown* item 1 and all five verification
    criteria; plan *The approach* steps 1–5.
  - Invoke the `dataviz` skill before writing any chart markup (plan,
    *Technical context*).
  - In `cmd/libretto/wiki.go`: replace `wikiCSS` with the token-complete
    dark-first palette from the plan's measured table (dark on bare `:root`,
    light behind `prefers-color-scheme: light`, alpha tokens for glass and
    chart track, colour only via `var()`); add the hero bento grid (capability
    count, criteria total, wide chart card); add `renderChart(caps)` returning
    `<svg role="img">` with one `<rect>` + `<text>` per capability, width
    `count * 100 / max` in viewBox units (integer math), zero-count bars at
    width 0 with label present, lexical order, `fill="var(--accent)"`; declare
    `#progress` and `.cap` reveal animations only inside
    `@media (prefers-reduced-motion: no-preference)` using
    `animation-timeline: scroll()` / `view()`, no scroll listeners; add the
    IntersectionObserver scroll-spy toggling `aria-current` on nav items
    beside the existing filter script, anchors working with JS off.
  - In `cmd/libretto/wiki_test.go`, same commit: add
    `TestWikiHTMLCarriesTheBentoHeroAndChart` (hero counts; one labelled bar
    per capability, proportional widths; the zero-criteria `checkout`
    fixture's label present at width 0),
    `TestWikiHTMLMotionRespectsReducedMotion` (progress and reveal animations
    only inside `prefers-reduced-motion: no-preference`; `animation-timeline`
    referenced, no scroll listeners for either),
    `TestWikiHTMLThemesAreTokenComplete` (full palette as `:root` tokens plus
    a complete redefinition block; no hex literal outside the two token
    blocks), and `TestWikiHTMLScrollSpyIsAnEnhancement` (nav entries remain
    plain #-anchors; inline script references IntersectionObserver and
    aria-current) — the cutter's gap, closed in the contract.
  - Closes when: `go test ./cmd/libretto/` is green with all nineteen wiki
    tests — the existing fifteen unchanged and the four new Proofs passing;
    `TestWikiHTMLThemesAreTokenComplete` forced red once via a stray hex in a
    component rule, observed failing, restored; the page regenerated against
    this repo's specs and opened — the look is the deliverable, so the look is
    the check; all six gates pass, including `spec-drift --anchors` resolving
    the new citations.
  - Waits on: nothing. Can start now.
  - Evidence: nineteen wiki tests green on commit "feat(cli): the wiki viewer's
    modern face"; TestWikiHTMLThemesAreTokenComplete forced red with a stray
    #ff0000 in a component rule, observed failing with its own message,
    restored; regenerated against this repo's 14 capabilities, opened in the
    browser and looked at (14 bars, motion present), artifact republished for
    the user; palette ratios measured in the plan, all ≥ 4.5:1.
